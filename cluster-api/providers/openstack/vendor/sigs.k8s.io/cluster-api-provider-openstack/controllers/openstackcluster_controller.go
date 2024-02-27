/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/collections"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/loadbalancer"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
)

const (
	BastionInstanceHashAnnotation = "infrastructure.cluster.x-k8s.io/bastion-hash"
)

// OpenStackClusterReconciler reconciles a OpenStackCluster object.
type OpenStackClusterReconciler struct {
	Client           client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	ScopeFactory     scope.Factory
	CaCertificates   []byte // PEM encoded ca certificates.
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=openstackclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

func (r *OpenStackClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the OpenStackCluster instance
	openStackCluster := &infrav1.OpenStackCluster{}
	err := r.Client.Get(ctx, req.NamespacedName, openStackCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, openStackCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	if annotations.IsPaused(cluster, openStackCluster) {
		log.Info("OpenStackCluster or linked Cluster is marked as paused. Not reconciling")
		return reconcile.Result{}, nil
	}

	patchHelper, err := patch.NewHelper(openStackCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always patch the openStackCluster when exiting this function so we can persist any OpenStackCluster changes.
	defer func() {
		if err := patchHelper.Patch(ctx, openStackCluster); err != nil {
			if reterr == nil {
				reterr = fmt.Errorf("error patching OpenStackCluster %s/%s: %w", openStackCluster.Namespace, openStackCluster.Name, err)
			}
		}
	}()

	scope, err := r.ScopeFactory.NewClientScopeFromCluster(ctx, r.Client, openStackCluster, r.CaCertificates, log)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Handle deleted clusters
	if !openStackCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, scope, cluster, openStackCluster)
	}

	// Handle non-deleted clusters
	return reconcileNormal(scope, cluster, openStackCluster)
}

func (r *OpenStackClusterReconciler) reconcileDelete(ctx context.Context, scope scope.Scope, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (ctrl.Result, error) {
	scope.Logger().Info("Reconciling Cluster delete")

	// Wait for machines to be deleted before removing the finalizer as they
	// depend on this resource to deprovision.  Additionally it appears that
	// allowing the Kubernetes API to vanish too quickly will upset the capi
	// kubeadm control plane controller.
	machines, err := collections.GetFilteredMachinesForCluster(ctx, r.Client, cluster)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(machines) != 0 {
		scope.Logger().Info("Waiting for machines to be deleted", "remaining", len(machines))
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	if err := deleteBastion(scope, cluster, openStackCluster); err != nil {
		return reconcile.Result{}, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return reconcile.Result{}, err
	}

	clusterName := fmt.Sprintf("%s-%s", cluster.Namespace, cluster.Name)

	if err = networkingService.DeletePorts(openStackCluster); err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete ports: %w", err))
		return reconcile.Result{}, fmt.Errorf("failed to delete ports: %w", err)
	}

	if openStackCluster.Spec.APIServerLoadBalancer.Enabled {
		loadBalancerService, err := loadbalancer.NewService(scope)
		if err != nil {
			return reconcile.Result{}, err
		}

		if err = loadBalancerService.DeleteLoadBalancer(openStackCluster, clusterName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete load balancer: %w", err))
			return reconcile.Result{}, fmt.Errorf("failed to delete load balancer: %w", err)
		}
	}

	if err = networkingService.DeleteSecurityGroups(openStackCluster, clusterName); err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete security groups: %w", err))
		return reconcile.Result{}, fmt.Errorf("failed to delete security groups: %w", err)
	}

	// if NodeCIDR was not set, no network was created.
	if openStackCluster.Spec.NodeCIDR != "" {
		if err = networkingService.DeleteRouter(openStackCluster, clusterName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete router: %w", err))
			return ctrl.Result{}, fmt.Errorf("failed to delete router: %w", err)
		}

		if err = networkingService.DeleteNetwork(openStackCluster, clusterName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete network: %w", err))
			return ctrl.Result{}, fmt.Errorf("failed to delete network: %w", err)
		}
	}

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(openStackCluster, infrav1.ClusterFinalizer)
	scope.Logger().Info("Reconciled Cluster deleted successfully")
	return ctrl.Result{}, nil
}

func contains(arr []string, target string) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}

func deleteBastion(scope scope.Scope, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) error {
	computeService, err := compute.NewService(scope)
	if err != nil {
		return err
	}
	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	instanceName := fmt.Sprintf("%s-bastion", cluster.Name)
	instanceStatus, err := computeService.GetInstanceStatusByName(openStackCluster, instanceName)
	if err != nil {
		return err
	}

	if instanceStatus != nil {
		instanceNS, err := instanceStatus.NetworkStatus()
		if err != nil {
			return err
		}
		addresses := instanceNS.Addresses()

		for _, address := range addresses {
			if address.Type == corev1.NodeExternalIP {
				if err = networkingService.DeleteFloatingIP(openStackCluster, address.Address); err != nil {
					handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete floating IP: %w", err))
					return fmt.Errorf("failed to delete floating IP: %w", err)
				}
			}
		}

		instanceSpec := bastionToInstanceSpec(openStackCluster, cluster.Name)

		if err = computeService.DeleteInstance(openStackCluster, openStackCluster, instanceStatus, instanceSpec); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete bastion: %w", err))
			return fmt.Errorf("failed to delete bastion: %w", err)
		}
	}

	openStackCluster.Status.Bastion = nil

	if err = networkingService.DeleteBastionSecurityGroup(openStackCluster, fmt.Sprintf("%s-%s", cluster.Namespace, cluster.Name)); err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete bastion security group: %w", err))
		return fmt.Errorf("failed to delete bastion security group: %w", err)
	}
	openStackCluster.Status.BastionSecurityGroup = nil

	delete(openStackCluster.ObjectMeta.Annotations, BastionInstanceHashAnnotation)

	return nil
}

func reconcileNormal(scope scope.Scope, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (ctrl.Result, error) { //nolint:unparam
	scope.Logger().Info("Reconciling Cluster")

	// If the OpenStackCluster doesn't have our finalizer, add it.
	if controllerutil.AddFinalizer(openStackCluster, infrav1.ClusterFinalizer) {
		// Register the finalizer immediately to avoid orphaning OpenStack resources on delete
		return reconcile.Result{}, nil
	}

	computeService, err := compute.NewService(scope)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = reconcileNetworkComponents(scope, cluster, openStackCluster)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = reconcileBastion(scope, cluster, openStackCluster); err != nil {
		return reconcile.Result{}, err
	}

	availabilityZones, err := computeService.GetAvailabilityZones()
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create a new list in case any AZs have been removed from OpenStack
	openStackCluster.Status.FailureDomains = make(clusterv1.FailureDomains)
	for _, az := range availabilityZones {
		// By default, the AZ is used or not used for control plane nodes depending on the flag
		found := !openStackCluster.Spec.ControlPlaneOmitAvailabilityZone
		// If explicit AZs for control plane nodes are given, they override the value
		if len(openStackCluster.Spec.ControlPlaneAvailabilityZones) > 0 {
			found = contains(openStackCluster.Spec.ControlPlaneAvailabilityZones, az.ZoneName)
		}
		// Add the AZ object to the failure domains for the cluster
		openStackCluster.Status.FailureDomains[az.ZoneName] = clusterv1.FailureDomainSpec{
			ControlPlane: found,
		}
	}

	openStackCluster.Status.Ready = true
	openStackCluster.Status.FailureMessage = nil
	openStackCluster.Status.FailureReason = nil
	scope.Logger().Info("Reconciled Cluster created successfully")
	return reconcile.Result{}, nil
}

func reconcileBastion(scope scope.Scope, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) error {
	scope.Logger().Info("Reconciling Bastion")

	if openStackCluster.Spec.Bastion == nil || !openStackCluster.Spec.Bastion.Enabled {
		return deleteBastion(scope, cluster, openStackCluster)
	}

	computeService, err := compute.NewService(scope)
	if err != nil {
		return err
	}

	instanceSpec := bastionToInstanceSpec(openStackCluster, cluster.Name)
	bastionHash, err := compute.HashInstanceSpec(instanceSpec)
	if err != nil {
		return fmt.Errorf("failed computing bastion hash from instance spec: %w", err)
	}

	instanceStatus, err := computeService.GetInstanceStatusByName(openStackCluster, fmt.Sprintf("%s-bastion", cluster.Name))
	if err != nil {
		return err
	}
	if instanceStatus != nil {
		if !bastionHashHasChanged(bastionHash, openStackCluster.ObjectMeta.Annotations) {
			bastion, err := instanceStatus.BastionStatus(openStackCluster)
			if err != nil {
				return err
			}
			// Add the current hash if no annotation is set.
			if _, ok := openStackCluster.ObjectMeta.Annotations[BastionInstanceHashAnnotation]; !ok {
				annotations.AddAnnotations(openStackCluster, map[string]string{BastionInstanceHashAnnotation: bastionHash})
			}
			openStackCluster.Status.Bastion = bastion
			return nil
		}

		if err := deleteBastion(scope, cluster, openStackCluster); err != nil {
			return err
		}
	}

	instanceStatus, err = computeService.CreateInstance(openStackCluster, openStackCluster, instanceSpec, cluster.Name, true)
	if err != nil {
		return fmt.Errorf("failed to reconcile bastion: %w", err)
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}
	clusterName := fmt.Sprintf("%s-%s", cluster.Namespace, cluster.Name)
	fp, err := networkingService.GetOrCreateFloatingIP(openStackCluster, openStackCluster, clusterName, openStackCluster.Spec.Bastion.Instance.FloatingIP)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to get or create floating IP for bastion: %w", err))
		return fmt.Errorf("failed to get or create floating IP for bastion: %w", err)
	}
	port, err := computeService.GetManagementPort(openStackCluster, instanceStatus)
	if err != nil {
		err = fmt.Errorf("getting management port for bastion: %w", err)
		handleUpdateOSCError(openStackCluster, err)
		return err
	}
	err = networkingService.AssociateFloatingIP(openStackCluster, fp, port.ID)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to associate floating IP with bastion: %w", err))
		return fmt.Errorf("failed to associate floating IP with bastion: %w", err)
	}

	bastion, err := instanceStatus.BastionStatus(openStackCluster)
	if err != nil {
		return err
	}
	bastion.FloatingIP = fp.FloatingIP
	openStackCluster.Status.Bastion = bastion
	annotations.AddAnnotations(openStackCluster, map[string]string{BastionInstanceHashAnnotation: bastionHash})
	return nil
}

func bastionToInstanceSpec(openStackCluster *infrav1.OpenStackCluster, clusterName string) *compute.InstanceSpec {
	name := fmt.Sprintf("%s-bastion", clusterName)
	instanceSpec := &compute.InstanceSpec{
		Name:          name,
		Flavor:        openStackCluster.Spec.Bastion.Instance.Flavor,
		SSHKeyName:    openStackCluster.Spec.Bastion.Instance.SSHKeyName,
		Image:         openStackCluster.Spec.Bastion.Instance.Image,
		ImageUUID:     openStackCluster.Spec.Bastion.Instance.ImageUUID,
		FailureDomain: openStackCluster.Spec.Bastion.AvailabilityZone,
		RootVolume:    openStackCluster.Spec.Bastion.Instance.RootVolume,
	}

	instanceSpec.SecurityGroups = openStackCluster.Spec.Bastion.Instance.SecurityGroups
	if openStackCluster.Spec.ManagedSecurityGroups {
		if openStackCluster.Status.BastionSecurityGroup != nil {
			instanceSpec.SecurityGroups = append(instanceSpec.SecurityGroups, infrav1.SecurityGroupFilter{
				ID: openStackCluster.Status.BastionSecurityGroup.ID,
			})
		}
	}

	instanceSpec.Ports = openStackCluster.Spec.Bastion.Instance.Ports

	return instanceSpec
}

// bastionHashHasChanged returns a boolean whether if the latest bastion hash, built from the instance spec, has changed or not.
func bastionHashHasChanged(computeHash string, clusterAnnotations map[string]string) bool {
	latestHash, ok := clusterAnnotations[BastionInstanceHashAnnotation]
	if !ok {
		return false
	}
	return latestHash != computeHash
}

func reconcileNetworkComponents(scope scope.Scope, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) error {
	clusterName := fmt.Sprintf("%s-%s", cluster.Namespace, cluster.Name)

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	scope.Logger().Info("Reconciling network components")

	err = networkingService.ReconcileExternalNetwork(openStackCluster)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile external network: %w", err))
		return fmt.Errorf("failed to reconcile external network: %w", err)
	}

	if openStackCluster.Spec.NodeCIDR == "" {
		scope.Logger().V(4).Info("No need to reconcile network, searching network and subnet instead")

		netOpts := openStackCluster.Spec.Network.ToListOpt()
		networkList, err := networkingService.GetNetworksByFilter(&netOpts)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to find network: %w", err))
			return fmt.Errorf("failed to find network: %w", err)
		}
		if len(networkList) == 0 {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to find any network"))
			return fmt.Errorf("failed to find any network")
		}
		if len(networkList) > 1 {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("found multiple networks (result: %v)", networkList))
			return fmt.Errorf("found multiple networks (result: %v)", networkList)
		}
		if openStackCluster.Status.Network == nil {
			openStackCluster.Status.Network = &infrav1.NetworkStatusWithSubnets{}
		}
		openStackCluster.Status.Network.ID = networkList[0].ID
		openStackCluster.Status.Network.Name = networkList[0].Name
		openStackCluster.Status.Network.Tags = networkList[0].Tags

		subnet, err := networkingService.GetNetworkSubnetByFilter(openStackCluster.Status.Network.ID, &openStackCluster.Spec.Subnet)
		if err != nil {
			err = fmt.Errorf("failed to find subnet: %w", err)

			// Set the cluster to failed if subnet filter is invalid
			if errors.Is(err, networking.ErrFilterMatch) {
				handleUpdateOSCError(openStackCluster, err)
			}

			return err
		}

		openStackCluster.Status.Network.Subnets = []infrav1.Subnet{
			{
				ID:   subnet.ID,
				Name: subnet.Name,
				CIDR: subnet.CIDR,
				Tags: subnet.Tags,
			},
		}
	} else {
		err := networkingService.ReconcileNetwork(openStackCluster, clusterName)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile network: %w", err))
			return fmt.Errorf("failed to reconcile network: %w", err)
		}
		err = networkingService.ReconcileSubnet(openStackCluster, clusterName)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile subnets: %w", err))
			return fmt.Errorf("failed to reconcile subnets: %w", err)
		}
		err = networkingService.ReconcileRouter(openStackCluster, clusterName)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile router: %w", err))
			return fmt.Errorf("failed to reconcile router: %w", err)
		}
	}

	err = networkingService.ReconcileSecurityGroups(openStackCluster, clusterName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile security groups: %w", err))
		return fmt.Errorf("failed to reconcile security groups: %w", err)
	}

	// Calculate the port that we will use for the API server
	var apiServerPort int
	switch {
	case openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
		apiServerPort = int(openStackCluster.Spec.ControlPlaneEndpoint.Port)
	case openStackCluster.Spec.APIServerPort != 0:
		apiServerPort = openStackCluster.Spec.APIServerPort
	default:
		apiServerPort = 6443
	}

	if openStackCluster.Spec.APIServerLoadBalancer.Enabled {
		loadBalancerService, err := loadbalancer.NewService(scope)
		if err != nil {
			return err
		}

		terminalFailure, err := loadBalancerService.ReconcileLoadBalancer(openStackCluster, clusterName, apiServerPort)
		if err != nil {
			// if it's terminalFailure (not Transient), set the Failure reason and message
			if terminalFailure {
				handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile load balancer: %w", err))
			}
			return fmt.Errorf("failed to reconcile load balancer: %w", err)
		}
	}

	if !openStackCluster.Spec.ControlPlaneEndpoint.IsValid() {
		var host string
		// If there is a load balancer use the floating IP for it if set, falling back to the internal IP
		switch {
		case openStackCluster.Spec.APIServerLoadBalancer.Enabled:
			if openStackCluster.Status.APIServerLoadBalancer.IP != "" {
				host = openStackCluster.Status.APIServerLoadBalancer.IP
			} else {
				host = openStackCluster.Status.APIServerLoadBalancer.InternalIP
			}
		case !openStackCluster.Spec.DisableAPIServerFloatingIP:
			// If floating IPs are not disabled, get one to use as the VIP for the control plane
			fp, err := networkingService.GetOrCreateFloatingIP(openStackCluster, openStackCluster, clusterName, openStackCluster.Spec.APIServerFloatingIP)
			if err != nil {
				handleUpdateOSCError(openStackCluster, fmt.Errorf("floating IP cannot be got or created: %w", err))
				return fmt.Errorf("floating IP cannot be got or created: %w", err)
			}
			host = fp.FloatingIP
		case openStackCluster.Spec.APIServerFixedIP != "":
			// If a fixed IP was specified, assume that the user is providing the extra configuration
			// to use that IP as the VIP for the API server, e.g. using keepalived or kube-vip
			host = openStackCluster.Spec.APIServerFixedIP
		default:
			// For now, we do not provide a managed VIP without either a load balancer or a floating IP
			// In the future, we could manage a VIP port on the cluster network and set allowedAddressPairs
			// accordingly when creating control plane machines
			// However this would require us to deploy software on the control plane hosts to manage the
			// VIP (e.g. keepalived/kube-vip)
			return fmt.Errorf("unable to determine VIP for API server")
		}

		// Set APIEndpoints so the Cluster API Cluster Controller can pull them
		openStackCluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
			Host: host,
			Port: int32(apiServerPort),
		}
	}

	return nil
}

func (r *OpenStackClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	clusterToInfraFn := util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind("OpenStackCluster"), mgr.GetClient(), &infrav1.OpenStackCluster{})
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.OpenStackCluster{},
			builder.WithPredicates(
				predicate.Funcs{
					// Avoid reconciling if the event triggering the reconciliation is related to incremental status updates
					UpdateFunc: func(e event.UpdateEvent) bool {
						oldCluster := e.ObjectOld.(*infrav1.OpenStackCluster).DeepCopy()
						newCluster := e.ObjectNew.(*infrav1.OpenStackCluster).DeepCopy()
						oldCluster.Status = infrav1.OpenStackClusterStatus{}
						newCluster.Status = infrav1.OpenStackClusterStatus{}
						oldCluster.ObjectMeta.ResourceVersion = ""
						newCluster.ObjectMeta.ResourceVersion = ""
						return !reflect.DeepEqual(oldCluster, newCluster)
					},
				},
			),
		).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, o client.Object) []reconcile.Request {
				requests := clusterToInfraFn(ctx, o)
				if len(requests) < 1 {
					return nil
				}

				c := &infrav1.OpenStackCluster{}
				if err := r.Client.Get(ctx, requests[0].NamespacedName, c); err != nil {
					log.V(4).Error(err, "Failed to get OpenStack cluster")
					return nil
				}

				if annotations.IsExternallyManaged(c) {
					log.V(4).Info("OpenStackCluster is externally managed, skipping mapping")
					return nil
				}
				return requests
			}),
			builder.WithPredicates(predicates.ClusterUnpaused(ctrl.LoggerFrom(ctx))),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(ctrl.LoggerFrom(ctx))).
		Complete(r)
}

func handleUpdateOSCError(openstackCluster *infrav1.OpenStackCluster, message error) {
	err := capierrors.UpdateClusterError
	openstackCluster.Status.FailureReason = &err
	openstackCluster.Status.FailureMessage = pointer.String(message.Error())
}
