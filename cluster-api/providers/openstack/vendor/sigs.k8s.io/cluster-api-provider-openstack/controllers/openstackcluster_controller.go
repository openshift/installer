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
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/loadbalancer"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	utils "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/controllers"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
)

const (
	waitForBastionToReconcile  = 15 * time.Second
	waitForOctaviaPortsCleanup = 15 * time.Second
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

func (r *OpenStackClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, reterr error) {
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
			result = ctrl.Result{}
			reterr = kerrors.NewAggregate([]error{reterr, fmt.Errorf("error patching OpenStackCluster %s/%s: %w", openStackCluster.Namespace, openStackCluster.Name, err)})
		}
	}()

	clientScope, err := r.ScopeFactory.NewClientScopeFromObject(ctx, r.Client, r.CaCertificates, log, openStackCluster)
	if err != nil {
		return reconcile.Result{}, err
	}
	scope := scope.NewWithLogger(clientScope, log)

	// Handle deleted clusters
	if !openStackCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, scope, cluster, openStackCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, scope, cluster, openStackCluster)
}

func (r *OpenStackClusterReconciler) reconcileDelete(ctx context.Context, scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (ctrl.Result, error) {
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

	clusterResourceName := names.ClusterResourceName(cluster)

	// A bastion may have been created if cluster initialisation previously reached populating the network status
	// We attempt to delete it even if no status was written, just in case
	if openStackCluster.Status.Network != nil {
		if err := r.deleteBastion(ctx, scope, cluster, openStackCluster); err != nil {
			return reconcile.Result{}, err
		}
	}

	// If a bastion server was found, we need to reconcile now until it's actually deleted.
	// We don't want to remove the cluster finalizer until the associated OpenStackServer resource is deleted.
	bastionServer, err := r.getBastionServer(ctx, openStackCluster, cluster)
	if client.IgnoreNotFound(err) != nil {
		return reconcile.Result{}, err
	}
	if bastionServer != nil {
		scope.Logger().Info("Waiting for the bastion OpenStackServer object to be deleted", "openStackServer", bastionServer.Name)
		return ctrl.Result{Requeue: true}, nil
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return reconcile.Result{}, err
	}

	if openStackCluster.Spec.APIServerLoadBalancer.IsEnabled() {
		loadBalancerService, err := loadbalancer.NewService(scope)
		if err != nil {
			return reconcile.Result{}, err
		}

		result, err := loadBalancerService.DeleteLoadBalancer(openStackCluster, clusterResourceName)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete load balancer: %w", err), false)
			return reconcile.Result{}, fmt.Errorf("failed to delete load balancer: %w", err)
		}
		if result != nil {
			return *result, nil
		}
	}

	// if ManagedSubnets was not set, no network was created.
	if len(openStackCluster.Spec.ManagedSubnets) > 0 {
		if err = networkingService.DeleteRouter(openStackCluster, clusterResourceName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete router: %w", err), false)
			return ctrl.Result{}, fmt.Errorf("failed to delete router: %w", err)
		}

		if err = networkingService.DeleteClusterPorts(openStackCluster); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete ports: %w", err), false)
			return reconcile.Result{}, fmt.Errorf("failed to delete ports: %w", err)
		}

		if err = networkingService.DeleteNetwork(openStackCluster, clusterResourceName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete network: %w", err), false)
			return ctrl.Result{}, fmt.Errorf("failed to delete network: %w", err)
		}
	}

	if err = networkingService.DeleteSecurityGroups(openStackCluster, clusterResourceName); err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete security groups: %w", err), false)
		return reconcile.Result{}, fmt.Errorf("failed to delete security groups: %w", err)
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

func (r *OpenStackClusterReconciler) deleteBastion(ctx context.Context, scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) error {
	scope.Logger().Info("Deleting Bastion")

	computeService, err := compute.NewService(scope)
	if err != nil {
		return err
	}
	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	bastionServer, err := r.getBastionServer(ctx, openStackCluster, cluster)
	if client.IgnoreNotFound(err) != nil {
		return err
	}

	var statusFloatingIP *string
	var specFloatingIP *string
	if openStackCluster.Status.Bastion != nil && openStackCluster.Status.Bastion.FloatingIP != "" {
		statusFloatingIP = &openStackCluster.Status.Bastion.FloatingIP
	}
	if openStackCluster.Spec.Bastion != nil && openStackCluster.Spec.Bastion.FloatingIP != nil {
		specFloatingIP = openStackCluster.Spec.Bastion.FloatingIP
	}

	// We only remove the bastion's floating IP if it exists and if it's not the same value defined both in the spec and in status.
	// This decision was made so if a user specifies a pre-created floating IP that is intended to only be used for the bastion, the floating IP won't get removed once the bastion is destroyed.
	if statusFloatingIP != nil && (specFloatingIP == nil || *statusFloatingIP != *specFloatingIP) {
		if err = networkingService.DeleteFloatingIP(openStackCluster, openStackCluster.Status.Bastion.FloatingIP); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete floating IP: %w", err), false)
			return fmt.Errorf("failed to delete floating IP: %w", err)
		}
	}

	bastionStatus := openStackCluster.Status.Bastion

	var instanceStatus *compute.InstanceStatus
	if bastionStatus != nil && bastionServer != nil && bastionServer.Status.InstanceID != nil {
		instanceStatus, err = computeService.GetInstanceStatus(*bastionServer.Status.InstanceID)
		if err != nil {
			return err
		}
	}

	if instanceStatus != nil {
		instanceNS, err := instanceStatus.NetworkStatus()
		if err != nil {
			return err
		}
		addresses := instanceNS.Addresses()

		for _, address := range addresses {
			if address.Type == corev1.NodeExternalIP {
				// If a floating IP retrieved is the same as what is set in the bastion spec, skip deleting it.
				// This decision was made so if a user specifies a pre-created floating IP that is intended to only be used for the bastion, the floating IP won't get removed once the bastion is destroyed.
				if specFloatingIP != nil && address.Address == *specFloatingIP {
					continue
				}
				// Floating IP may not have properly saved in bastion status (thus not deleted above), delete any remaining floating IP
				if err = networkingService.DeleteFloatingIP(openStackCluster, address.Address); err != nil {
					handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete floating IP: %w", err), false)
					return fmt.Errorf("failed to delete floating IP: %w", err)
				}
			}
		}
	}

	if err := r.reconcileDeleteBastionServer(ctx, scope, openStackCluster, cluster); err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete bastion: %w", err), false)
		return fmt.Errorf("failed to delete bastion: %w", err)
	}

	openStackCluster.Status.Bastion = nil
	scope.Logger().Info("Deleted Bastion")

	return nil
}

func (r *OpenStackClusterReconciler) reconcileNormal(ctx context.Context, scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (ctrl.Result, error) { //nolint:unparam
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

	availabilityZones, err := computeService.GetAvailabilityZones()
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create a new list in case any AZs have been removed from OpenStack
	openStackCluster.Status.FailureDomains = make(clusterv1beta1.FailureDomains)
	for _, az := range availabilityZones {
		// By default, the AZ is used or not used for control plane nodes depending on the flag
		found := !ptr.Deref(openStackCluster.Spec.ControlPlaneOmitAvailabilityZone, false)
		// If explicit AZs for control plane nodes are given, they override the value
		if len(openStackCluster.Spec.ControlPlaneAvailabilityZones) > 0 {
			found = contains(openStackCluster.Spec.ControlPlaneAvailabilityZones, az.ZoneName)
		}
		// Add the AZ object to the failure domains for the cluster

		openStackCluster.Status.FailureDomains[az.ZoneName] = clusterv1beta1.FailureDomainSpec{
			ControlPlane: found,
		}
	}

	openStackCluster.Status.Ready = true
	openStackCluster.Status.FailureMessage = nil
	openStackCluster.Status.FailureReason = nil
	scope.Logger().Info("Reconciled Cluster created successfully")

	result, err := r.reconcileBastion(ctx, scope, cluster, openStackCluster)
	if err != nil {
		return reconcile.Result{}, err
	}
	if result != nil {
		return *result, nil
	}
	scope.Logger().Info("Reconciled Bastion created successfully")

	return reconcile.Result{}, nil
}

func (r *OpenStackClusterReconciler) reconcileBastion(ctx context.Context, scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (*ctrl.Result, error) {
	scope.Logger().V(4).Info("Reconciling Bastion")

	clusterResourceName := names.ClusterResourceName(cluster)

	computeService, err := compute.NewService(scope)
	if err != nil {
		return nil, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return nil, err
	}

	bastionServer, waitingForServer, err := r.reconcileBastionServer(ctx, scope, openStackCluster, cluster)
	if err != nil || waitingForServer {
		return &reconcile.Result{RequeueAfter: waitForBastionToReconcile}, err
	}
	if bastionServer == nil {
		return nil, nil
	}

	var instanceStatus *compute.InstanceStatus
	if bastionServer != nil && bastionServer.Status.InstanceID != nil {
		if instanceStatus, err = computeService.GetInstanceStatus(*bastionServer.Status.InstanceID); err != nil {
			return nil, err
		}
	}
	if instanceStatus == nil {
		// At this point we return an error if we don't have an instance status
		return nil, fmt.Errorf("bastion instance status is nil")
	}

	// Save hash & status as soon as we know we have an instance
	instanceStatus.UpdateBastionStatus(openStackCluster)

	port, err := computeService.GetManagementPort(openStackCluster, instanceStatus)
	if err != nil {
		err = fmt.Errorf("getting management port for bastion: %w", err)
		handleUpdateOSCError(openStackCluster, err, false)
		return nil, err
	}

	if !ptr.Deref(openStackCluster.Spec.DisableExternalNetwork, false) {
		return bastionAddFloatingIP(openStackCluster, clusterResourceName, port, networkingService)
	}

	return nil, nil
}

func bastionAddFloatingIP(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string, port *ports.Port, networkingService *networking.Service) (*reconcile.Result, error) {
	fp, err := networkingService.GetFloatingIPByPortID(port.ID)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to get or create floating IP for bastion: %w", err), false)
		return nil, fmt.Errorf("failed to get floating IP for bastion port: %w", err)
	}
	if fp != nil {
		// Floating IP is already attached to bastion, no need to proceed
		openStackCluster.Status.Bastion.FloatingIP = fp.FloatingIP
		return nil, nil
	}

	var floatingIP *string
	switch {
	case openStackCluster.Status.Bastion.FloatingIP != "":
		// Some floating IP has already been created for this bastion, make sure we re-use it
		floatingIP = &openStackCluster.Status.Bastion.FloatingIP
	case openStackCluster.Spec.Bastion.FloatingIP != nil:
		// Use floating IP from the spec
		floatingIP = openStackCluster.Spec.Bastion.FloatingIP
	}
	// Check if there is an existing floating IP attached to bastion, in case where FloatingIP would not yet have been stored in cluster status
	fp, err = networkingService.GetOrCreateFloatingIP(openStackCluster, openStackCluster, clusterResourceName, floatingIP)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to get or create floating IP for bastion: %w", err), false)
		return nil, fmt.Errorf("failed to get or create floating IP for bastion: %w", err)
	}
	openStackCluster.Status.Bastion.FloatingIP = fp.FloatingIP

	err = networkingService.AssociateFloatingIP(openStackCluster, fp, port.ID)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to associate floating IP with bastion: %w", err), false)
		return nil, fmt.Errorf("failed to associate floating IP with bastion: %w", err)
	}

	return nil, nil
}

// reconcileDeleteBastionServer reconciles the OpenStackServer object for the OpenStackCluster bastion.
// It returns nil if the OpenStackServer object is not found, otherwise it returns an error if any.
func (r *OpenStackClusterReconciler) reconcileDeleteBastionServer(ctx context.Context, scope *scope.WithLogger, openStackCluster *infrav1.OpenStackCluster, cluster *clusterv1.Cluster) error {
	scope.Logger().Info("Reconciling Bastion delete server")
	server := &infrav1alpha1.OpenStackServer{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: openStackCluster.Namespace, Name: bastionName(cluster.Name)}, server)
	if client.IgnoreNotFound(err) != nil {
		return err
	}
	if apierrors.IsNotFound(err) {
		return nil
	}

	return r.Client.Delete(ctx, server)
}

// reconcileBastionServer reconciles the OpenStackServer object for the OpenStackCluster bastion.
// It returns the OpenStackServer object, a boolean indicating if the reconciliation should continue
// and an error if any.
func (r *OpenStackClusterReconciler) reconcileBastionServer(ctx context.Context, scope *scope.WithLogger, openStackCluster *infrav1.OpenStackCluster, cluster *clusterv1.Cluster) (*infrav1alpha1.OpenStackServer, bool, error) {
	server, err := r.getBastionServer(ctx, openStackCluster, cluster)
	if client.IgnoreNotFound(err) != nil {
		scope.Logger().Error(err, "Failed to get the bastion OpenStackServer object")
		return nil, true, err
	}
	bastionNotFound := apierrors.IsNotFound(err)

	// If the bastion is not enabled, we don't need to create it and continue with the reconciliation.
	if bastionNotFound && !openStackCluster.Spec.Bastion.IsEnabled() {
		return nil, false, nil
	}

	// If the bastion is found but is not enabled, we need to delete it and reconcile.
	if !bastionNotFound && !openStackCluster.Spec.Bastion.IsEnabled() {
		scope.Logger().Info("Bastion is not enabled, deleting the OpenStackServer object")
		if err := r.deleteBastion(ctx, scope, cluster, openStackCluster); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	}

	// If the bastion is found but the spec has changed, we need to delete it and reconcile.
	bastionServerSpec, err := bastionToOpenStackServerSpec(openStackCluster)
	if err != nil {
		return nil, true, err
	}
	if !bastionNotFound && server != nil && !apiequality.Semantic.DeepEqual(bastionServerSpec, &server.Spec) {
		scope.Logger().Info("Bastion spec has changed, re-creating the OpenStackServer object")
		if err := r.deleteBastion(ctx, scope, cluster, openStackCluster); err != nil {
			return nil, true, err
		}
		return nil, true, nil
	}

	// If the bastion is not found, we need to create it.
	if bastionNotFound {
		scope.Logger().Info("Creating the bastion OpenStackServer object")
		server, err = r.createBastionServer(ctx, openStackCluster, cluster)
		if err != nil {
			return nil, true, err
		}
		return server, true, nil
	}

	// If the bastion server is not ready, we need to wait for it to be ready and reconcile.
	if !server.Status.Ready {
		scope.Logger().Info("Waiting for the bastion OpenStackServer to be ready")
		return server, true, nil
	}

	return server, false, nil
}

// getBastionServer returns the OpenStackServer object for the bastion server.
// It returns the OpenStackServer object and an error if any.
func (r *OpenStackClusterReconciler) getBastionServer(ctx context.Context, openStackCluster *infrav1.OpenStackCluster, cluster *clusterv1.Cluster) (*infrav1alpha1.OpenStackServer, error) {
	bastionServer := &infrav1alpha1.OpenStackServer{}
	bastionServerName := client.ObjectKey{
		Namespace: openStackCluster.Namespace,
		Name:      bastionName(cluster.Name),
	}
	err := r.Client.Get(ctx, bastionServerName, bastionServer)
	if err != nil {
		return nil, err
	}
	return bastionServer, nil
}

// createBastionServer creates the OpenStackServer object for the bastion server.
// It returns the OpenStackServer object and an error if any.
func (r *OpenStackClusterReconciler) createBastionServer(ctx context.Context, openStackCluster *infrav1.OpenStackCluster, cluster *clusterv1.Cluster) (*infrav1alpha1.OpenStackServer, error) {
	bastionServerSpec, err := bastionToOpenStackServerSpec(openStackCluster)
	if err != nil {
		return nil, err
	}
	bastionServer := &infrav1alpha1.OpenStackServer{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: openStackCluster.Labels[clusterv1.ClusterNameLabel],
			},
			Name:      bastionName(cluster.Name),
			Namespace: openStackCluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: openStackCluster.APIVersion,
					Kind:       openStackCluster.Kind,
					Name:       openStackCluster.Name,
					UID:        openStackCluster.UID,
				},
			},
		},
		Spec: *bastionServerSpec,
	}

	if err := r.Client.Create(ctx, bastionServer); err != nil {
		return nil, fmt.Errorf("failed to create bastion server: %w", err)
	}
	return bastionServer, nil
}

// bastionToOpenStackServerSpec converts the OpenStackMachineSpec for the bastion to an OpenStackServerSpec.
// It returns the OpenStackServerSpec and an error if any.
func bastionToOpenStackServerSpec(openStackCluster *infrav1.OpenStackCluster) (*infrav1alpha1.OpenStackServerSpec, error) {
	bastion := openStackCluster.Spec.Bastion
	if bastion == nil {
		bastion = &infrav1.Bastion{}
	}
	if bastion.Spec == nil {
		// For the case when Bastion is deleted but we don't have spec, let's use an empty one.
		// v1beta1 API validations prevent this from happening in normal circumstances.
		bastion.Spec = &infrav1.OpenStackMachineSpec{}
	}

	az := ""
	if bastion.AvailabilityZone != nil {
		az = *bastion.AvailabilityZone
	}
	openStackServerSpec, err := openStackMachineSpecToOpenStackServerSpec(bastion.Spec, openStackCluster.Spec.IdentityRef, compute.InstanceTags(bastion.Spec, openStackCluster), az, nil, getBastionSecurityGroupID(openStackCluster), openStackCluster.Status.Network)
	if err != nil {
		return nil, err
	}

	return openStackServerSpec, nil
}

func bastionName(clusterResourceName string) string {
	return fmt.Sprintf("%s-bastion", clusterResourceName)
}

// getBastionSecurityGroupID returns the ID of the bastion security group if
// managed security groups is enabled.
func getBastionSecurityGroupID(openStackCluster *infrav1.OpenStackCluster) *string {
	if openStackCluster.Spec.ManagedSecurityGroups == nil {
		return nil
	}

	if openStackCluster.Status.BastionSecurityGroup != nil {
		return &openStackCluster.Status.BastionSecurityGroup.ID
	}
	return nil
}

func resolveLoadBalancerNetwork(openStackCluster *infrav1.OpenStackCluster, networkingService *networking.Service) error {
	lbSpec := openStackCluster.Spec.APIServerLoadBalancer
	if lbSpec.IsEnabled() {
		lbStatus := openStackCluster.Status.APIServerLoadBalancer
		if lbStatus == nil {
			lbStatus = &infrav1.LoadBalancer{}
			openStackCluster.Status.APIServerLoadBalancer = lbStatus
		}

		lbNetStatus := lbStatus.LoadBalancerNetwork
		if lbNetStatus == nil {
			lbNetStatus = &infrav1.NetworkStatusWithSubnets{
				NetworkStatus: infrav1.NetworkStatus{},
			}
		}

		if lbSpec.Network != nil {
			lbNet, err := networkingService.GetNetworkByParam(lbSpec.Network)
			if err != nil {
				if errors.Is(err, capoerrors.ErrFilterMatch) {
					handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to find loadbalancer network: %w", err), true)
				}
				return fmt.Errorf("failed to find network: %w", err)
			}

			lbNetStatus.Name = lbNet.Name
			lbNetStatus.ID = lbNet.ID
			lbNetStatus.Tags = lbNet.Tags

			// Filter out only relevant subnets specified by the spec
			lbNetStatus.Subnets = []infrav1.Subnet{}
			for _, s := range lbSpec.Subnets {
				matchFound := false
				for _, subnetID := range lbNet.Subnets {
					subnet, err := networkingService.GetSubnetByParam(&s)
					if s.ID != nil && subnetID == *s.ID && err == nil {
						matchFound = true
						lbNetStatus.Subnets = append(
							lbNetStatus.Subnets, infrav1.Subnet{
								ID:   subnet.ID,
								Name: subnet.Name,
								CIDR: subnet.CIDR,
								Tags: subnet.Tags,
							})
					}
				}
				if !matchFound {
					handleUpdateOSCError(openStackCluster, fmt.Errorf("no subnet match was found in the specified network (specified subnet: %v, available subnets: %v)", s, lbNet.Subnets), false)
					return fmt.Errorf("no subnet match was found in the specified network (specified subnet: %v, available subnets: %v)", s, lbNet.Subnets)
				}
			}

			openStackCluster.Status.APIServerLoadBalancer.LoadBalancerNetwork = lbNetStatus
		}
	}

	return nil
}

func reconcileNetworkComponents(scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) error {
	clusterResourceName := names.ClusterResourceName(cluster)

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	scope.Logger().Info("Reconciling network components")

	err = networkingService.ReconcileExternalNetwork(openStackCluster)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile external network: %w", err), false)
		return fmt.Errorf("failed to reconcile external network: %w", err)
	}

	if len(openStackCluster.Spec.ManagedSubnets) == 0 {
		if err := reconcilePreExistingNetworkComponents(scope, networkingService, openStackCluster); err != nil {
			return err
		}
	} else if len(openStackCluster.Spec.ManagedSubnets) == 1 {
		if err := reconcileProvisionedNetworkComponents(networkingService, openStackCluster, clusterResourceName); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("failed to reconcile network: ManagedSubnets only supports one element, %d provided", len(openStackCluster.Spec.ManagedSubnets))
	}

	err = resolveLoadBalancerNetwork(openStackCluster, networkingService)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile loadbalancer network: %w", err), false)
		return fmt.Errorf("failed to reconcile loadbalancer network: %w", err)
	}

	err = networkingService.ReconcileSecurityGroups(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile security groups: %w", err), false)
		return fmt.Errorf("failed to reconcile security groups: %w", err)
	}

	return reconcileControlPlaneEndpoint(scope, networkingService, openStackCluster, clusterResourceName)
}

// reconcilePreExistingNetworkComponents reconciles the cluster network status when the cluster is
// using pre-existing networks, subnets and router which are not provisioned by the
// cluster controller.
func reconcilePreExistingNetworkComponents(scope *scope.WithLogger, networkingService *networking.Service, openStackCluster *infrav1.OpenStackCluster) error {
	scope.Logger().V(4).Info("No need to reconcile network, searching network, subnet and router instead")

	if openStackCluster.Status.Network == nil {
		openStackCluster.Status.Network = &infrav1.NetworkStatusWithSubnets{}
	}

	if openStackCluster.Spec.Network != nil {
		network, err := networkingService.GetNetworkByParam(openStackCluster.Spec.Network)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to find network: %w", err), false)
			return fmt.Errorf("error fetching cluster network: %w", err)
		}
		setClusterNetwork(openStackCluster, network)
	}

	subnets, err := getClusterSubnets(networkingService, openStackCluster)
	if err != nil {
		return err
	}

	// Populate the cluster status with the cluster subnets
	capoSubnets := make([]infrav1.Subnet, len(subnets))
	for i := range subnets {
		subnet := &subnets[i]
		capoSubnets[i] = infrav1.Subnet{
			ID:   subnet.ID,
			Name: subnet.Name,
			CIDR: subnet.CIDR,
			Tags: subnet.Tags,
		}
	}
	if err := utils.ValidateSubnets(capoSubnets); err != nil {
		return err
	}
	openStackCluster.Status.Network.Subnets = capoSubnets

	// If network is not yet populated, use networkID defined on the first
	// cluster subnet to get the Network. Cluster subnets are constrained to
	// be in the same network.
	if openStackCluster.Status.Network.ID == "" && len(subnets) > 0 {
		network, err := networkingService.GetNetworkByID(subnets[0].NetworkID)
		if err != nil {
			return err
		}
		setClusterNetwork(openStackCluster, network)
	}

	if openStackCluster.Spec.Router != nil {
		router, err := networkingService.GetRouterByParam(openStackCluster.Spec.Router)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to find router: %w", err), false)
			return fmt.Errorf("error fetching cluster router: %w", err)
		}

		scope.Logger().V(4).Info("Found pre-existing router", "id", router.ID, "name", router.Name)

		routerIPs := []string{}
		for _, ip := range router.GatewayInfo.ExternalFixedIPs {
			routerIPs = append(routerIPs, ip.IPAddress)
		}

		openStackCluster.Status.Router = &infrav1.Router{
			Name: router.Name,
			ID:   router.ID,
			Tags: router.Tags,
			IPs:  routerIPs,
		}
	}

	return nil
}

// reconcileProvisionedNetworkComponents reconciles the cluster network status when the cluster is
// using networks, subnets and router provisioned by the cluster controller.
func reconcileProvisionedNetworkComponents(networkingService *networking.Service, openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	err := networkingService.ReconcileNetwork(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile network: %w", err), false)
		return fmt.Errorf("failed to reconcile network: %w", err)
	}
	err = networkingService.ReconcileSubnet(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile subnets: %w", err), false)
		return fmt.Errorf("failed to reconcile subnets: %w", err)
	}
	err = networkingService.ReconcileRouter(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile router: %w", err), false)
		return fmt.Errorf("failed to reconcile router: %w", err)
	}

	return nil
}

// reconcileControlPlaneEndpoint configures the control plane endpoint for the
// cluster, creating it if necessary, and updates ControlPlaneEndpoint in the
// cluster spec.
func reconcileControlPlaneEndpoint(scope *scope.WithLogger, networkingService *networking.Service, openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	// Calculate the port that we will use for the API server
	apiServerPort := getAPIServerPort(openStackCluster)

	// host must be set by a matching control plane endpoint provider below
	var host string

	switch {
	// API server load balancer is enabled. Create an Octavia load balancer.
	// Note that we reconcile the load balancer even if the control plane
	// endpoint is already set.
	case openStackCluster.Spec.APIServerLoadBalancer.IsEnabled():
		loadBalancerService, err := loadbalancer.NewService(scope)
		if err != nil {
			return err
		}

		terminalFailure, err := loadBalancerService.ReconcileLoadBalancer(openStackCluster, clusterResourceName, int(apiServerPort))
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile load balancer: %w", err), terminalFailure)
			return fmt.Errorf("failed to reconcile load balancer: %w", err)
		}

		// Control plane endpoint is the floating IP if one was defined, otherwise the VIP address
		if openStackCluster.Status.APIServerLoadBalancer.IP != "" {
			host = openStackCluster.Status.APIServerLoadBalancer.IP
		} else {
			host = openStackCluster.Status.APIServerLoadBalancer.InternalIP
		}

	// Control plane endpoint is already set
	// Note that checking this here means that we don't re-execute any of
	// the branches below if the control plane endpoint is already set.
	case openStackCluster.Spec.ControlPlaneEndpoint != nil && openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
		host = openStackCluster.Spec.ControlPlaneEndpoint.Host

	// API server load balancer is disabled, but external netowork and floating IP are not. Create
	// a floating IP to be attached directly to a control plane host.
	case !ptr.Deref(openStackCluster.Spec.DisableAPIServerFloatingIP, false) && !ptr.Deref(openStackCluster.Spec.DisableExternalNetwork, false):
		fp, err := networkingService.GetOrCreateFloatingIP(openStackCluster, openStackCluster, clusterResourceName, openStackCluster.Spec.APIServerFloatingIP)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("floating IP cannot be got or created: %w", err), false)
			return fmt.Errorf("floating IP cannot be got or created: %w", err)
		}
		host = fp.FloatingIP

	// API server load balancer is disabled and we aren't using a control
	// plane floating IP. In this case we configure APIServerFixedIP as the
	// control plane endpoint and leave it to the user to configure load
	// balancing.
	case openStackCluster.Spec.APIServerFixedIP != nil:
		host = *openStackCluster.Spec.APIServerFixedIP

	// Control plane endpoint is not set, and none can be created
	default:
		err := fmt.Errorf("unable to determine control plane endpoint")
		handleUpdateOSCError(openStackCluster, err, false)
		return err
	}

	openStackCluster.Spec.ControlPlaneEndpoint = &clusterv1beta1.APIEndpoint{
		Host: host,
		Port: apiServerPort,
	}

	return nil
}

// getAPIServerPort returns the port to use for the API server based on the cluster spec.
func getAPIServerPort(openStackCluster *infrav1.OpenStackCluster) int32 {
	switch {
	case openStackCluster.Spec.ControlPlaneEndpoint != nil && openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
		return openStackCluster.Spec.ControlPlaneEndpoint.Port
	case openStackCluster.Spec.APIServerPort != nil:
		return int32(*openStackCluster.Spec.APIServerPort)
	}
	return 6443
}

func (r *OpenStackClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	clusterToInfraFn := util.ClusterToInfrastructureMapFunc(ctx, infrav1.SchemeGroupVersion.WithKind("OpenStackCluster"), mgr.GetClient(), &infrav1.OpenStackCluster{})
	log := ctrl.LoggerFrom(ctx)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.OpenStackCluster{}).
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
			builder.WithPredicates(predicates.ClusterUnpaused(mgr.GetScheme(), ctrl.LoggerFrom(ctx))),
		).
		Watches(
			&infrav1alpha1.OpenStackServer{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &infrav1.OpenStackCluster{}),
			builder.WithPredicates(OpenStackServerReconcileComplete(log)),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(mgr.GetScheme(), ctrl.LoggerFrom(ctx))).
		Complete(r)
}

func handleUpdateOSCError(openstackCluster *infrav1.OpenStackCluster, message error, isFatal bool) {
	if isFatal {
		err := capoerrors.DeprecatedCAPOUpdateClusterError
		openstackCluster.Status.FailureReason = &err
		openstackCluster.Status.FailureMessage = ptr.To(message.Error())
	}
}

// getClusterSubnets retrieves the subnets based on the Subnet filters specified on OpenstackCluster.
func getClusterSubnets(networkingService *networking.Service, openStackCluster *infrav1.OpenStackCluster) ([]subnets.Subnet, error) {
	var clusterSubnets []subnets.Subnet
	var err error
	openStackClusterSubnets := openStackCluster.Spec.Subnets
	networkID := ""
	if openStackCluster.Status.Network != nil {
		networkID = openStackCluster.Status.Network.ID
	}

	if len(openStackClusterSubnets) == 0 {
		if networkID == "" {
			// This should be a validation error
			return nil, fmt.Errorf("no network or subnets specified in OpenStackCluster spec")
		}

		listOpts := subnets.ListOpts{
			NetworkID: networkID,
		}
		clusterSubnets, err = networkingService.GetSubnetsByFilter(listOpts)
		if err != nil {
			err = fmt.Errorf("failed to find subnets: %w", err)
			if errors.Is(err, capoerrors.ErrFilterMatch) {
				handleUpdateOSCError(openStackCluster, err, true)
			}
			return nil, err
		}
		if len(clusterSubnets) > 2 {
			return nil, fmt.Errorf("more than two subnets found in the Network. Specify the subnets in the OpenStackCluster.Spec instead")
		}
	} else {
		for subnet := range openStackClusterSubnets {
			filteredSubnet, err := networkingService.GetNetworkSubnetByParam(networkID, &openStackClusterSubnets[subnet])
			if err != nil {
				err = fmt.Errorf("failed to find subnet %d in network %s: %w", subnet, networkID, err)
				if errors.Is(err, capoerrors.ErrFilterMatch) {
					handleUpdateOSCError(openStackCluster, err, true)
				}
				return nil, err
			}
			clusterSubnets = append(clusterSubnets, *filteredSubnet)

			// Constrain the next search to the network of the first subnet
			networkID = filteredSubnet.NetworkID
		}
	}
	return clusterSubnets, nil
}

// setClusterNetwork sets network information in the cluster status from an OpenStack network.
func setClusterNetwork(openStackCluster *infrav1.OpenStackCluster, network *networks.Network) {
	openStackCluster.Status.Network.ID = network.ID
	openStackCluster.Status.Network.Name = network.Name
	openStackCluster.Status.Network.Tags = network.Tags
}
