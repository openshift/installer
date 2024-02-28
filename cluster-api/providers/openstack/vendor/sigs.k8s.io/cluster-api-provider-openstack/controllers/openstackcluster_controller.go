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

	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
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
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/compute"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/loadbalancer"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/services/networking"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/scope"
	utils "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/controllers"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
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

	clientScope, err := r.ScopeFactory.NewClientScopeFromCluster(ctx, r.Client, openStackCluster, r.CaCertificates, log)
	if err != nil {
		return reconcile.Result{}, err
	}
	scope := scope.NewWithLogger(clientScope, log)

	// Handle deleted clusters
	if !openStackCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, scope, cluster, openStackCluster)
	}

	// Handle non-deleted clusters
	return reconcileNormal(scope, cluster, openStackCluster)
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
		// Attempt to resolve bastion resources before delete. We don't need to worry about starting if the resources have changed on update.
		if _, err := resolveBastionResources(scope, clusterResourceName, openStackCluster); err != nil {
			return reconcile.Result{}, err
		}

		if err := deleteBastion(scope, cluster, openStackCluster); err != nil {
			return reconcile.Result{}, err
		}
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

		if err = loadBalancerService.DeleteLoadBalancer(openStackCluster, clusterResourceName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete load balancer: %w", err))
			return reconcile.Result{}, fmt.Errorf("failed to delete load balancer: %w", err)
		}
	}

	// if ManagedSubnets was not set, no network was created.
	if len(openStackCluster.Spec.ManagedSubnets) > 0 {
		if err = networkingService.DeleteRouter(openStackCluster, clusterResourceName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete router: %w", err))
			return ctrl.Result{}, fmt.Errorf("failed to delete router: %w", err)
		}

		if err = networkingService.DeleteClusterPorts(openStackCluster); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete ports: %w", err))
			return reconcile.Result{}, fmt.Errorf("failed to delete ports: %w", err)
		}

		if err = networkingService.DeleteNetwork(openStackCluster, clusterResourceName); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete network: %w", err))
			return ctrl.Result{}, fmt.Errorf("failed to delete network: %w", err)
		}
	}

	if err = networkingService.DeleteSecurityGroups(openStackCluster, clusterResourceName); err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete security groups: %w", err))
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

func resolveBastionResources(scope *scope.WithLogger, clusterResourceName string, openStackCluster *infrav1.OpenStackCluster) (bool, error) {
	// Resolve and store resources for the bastion
	if openStackCluster.Spec.Bastion.IsEnabled() {
		if openStackCluster.Status.Bastion == nil {
			openStackCluster.Status.Bastion = &infrav1.BastionStatus{}
		}
		if openStackCluster.Spec.Bastion.Spec == nil {
			return false, fmt.Errorf("bastion spec is nil when bastion is enabled, this shouldn't happen")
		}
		resolved := openStackCluster.Status.Bastion.Resolved
		if resolved == nil {
			resolved = &infrav1.ResolvedMachineSpec{}
			openStackCluster.Status.Bastion.Resolved = resolved
		}
		changed, err := compute.ResolveMachineSpec(scope,
			openStackCluster.Spec.Bastion.Spec, resolved,
			clusterResourceName, bastionName(clusterResourceName),
			openStackCluster, getBastionSecurityGroupID(openStackCluster))
		if err != nil {
			return false, err
		}
		if changed {
			// If the resolved machine spec changed we need to restart the reconcile to avoid inconsistencies between reconciles.
			return true, nil
		}
		resources := openStackCluster.Status.Bastion.Resources
		if resources == nil {
			resources = &infrav1.MachineResources{}
			openStackCluster.Status.Bastion.Resources = resources
		}

		err = compute.AdoptMachineResources(scope, resolved, resources)
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func deleteBastion(scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) error {
	scope.Logger().Info("Deleting Bastion")

	computeService, err := compute.NewService(scope)
	if err != nil {
		return err
	}
	networkingService, err := networking.NewService(scope)
	if err != nil {
		return err
	}

	if openStackCluster.Status.Bastion != nil && openStackCluster.Status.Bastion.FloatingIP != "" {
		if err = networkingService.DeleteFloatingIP(openStackCluster, openStackCluster.Status.Bastion.FloatingIP); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete floating IP: %w", err))
			return fmt.Errorf("failed to delete floating IP: %w", err)
		}
	}

	bastionStatus := openStackCluster.Status.Bastion

	var instanceStatus *compute.InstanceStatus
	if bastionStatus != nil && bastionStatus.ID != "" {
		instanceStatus, err = computeService.GetInstanceStatus(openStackCluster.Status.Bastion.ID)
		if err != nil {
			return err
		}
	} else {
		instanceStatus, err = computeService.GetInstanceStatusByName(openStackCluster, bastionName(cluster.Name))
		if err != nil {
			return err
		}
	}

	// If no instance was created we currently need to check for orphaned
	// volumes.
	if instanceStatus == nil {
		bastion := openStackCluster.Spec.Bastion
		if bastion != nil && bastion.Spec != nil {
			if err := computeService.DeleteVolumes(bastionName(cluster.Name), bastion.Spec.RootVolume, bastion.Spec.AdditionalBlockDevices); err != nil {
				return fmt.Errorf("delete volumes: %w", err)
			}
		}
	} else {
		instanceNS, err := instanceStatus.NetworkStatus()
		if err != nil {
			return err
		}
		addresses := instanceNS.Addresses()

		for _, address := range addresses {
			if address.Type == corev1.NodeExternalIP {
				// Floating IP may not have properly saved in bastion status (thus not deleted above), delete any remaining floating IP
				if err = networkingService.DeleteFloatingIP(openStackCluster, address.Address); err != nil {
					handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete floating IP: %w", err))
					return fmt.Errorf("failed to delete floating IP: %w", err)
				}
			}
		}

		if err = computeService.DeleteInstance(openStackCluster, instanceStatus); err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete bastion: %w", err))
			return fmt.Errorf("failed to delete bastion: %w", err)
		}
	}

	if bastionStatus != nil && bastionStatus.Resources != nil {
		trunkSupported, err := networkingService.IsTrunkExtSupported()
		if err != nil {
			return err
		}
		for _, port := range bastionStatus.Resources.Ports {
			if err := networkingService.DeleteInstanceTrunkAndPort(openStackCluster, port, trunkSupported); err != nil {
				handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to delete port: %w", err))
				return fmt.Errorf("failed to delete port: %w", err)
			}
		}
		bastionStatus.Resources.Ports = nil
	}

	scope.Logger().Info("Deleted Bastion")

	openStackCluster.Status.Bastion = nil
	delete(openStackCluster.ObjectMeta.Annotations, BastionInstanceHashAnnotation)

	return nil
}

func reconcileNormal(scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (ctrl.Result, error) { //nolint:unparam
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

	result, err := reconcileBastion(scope, cluster, openStackCluster)
	if err != nil {
		return reconcile.Result{}, err
	}
	if result != nil {
		return *result, nil
	}

	availabilityZones, err := computeService.GetAvailabilityZones()
	if err != nil {
		return ctrl.Result{}, err
	}

	// Create a new list in case any AZs have been removed from OpenStack
	openStackCluster.Status.FailureDomains = make(clusterv1.FailureDomains)
	for _, az := range availabilityZones {
		// By default, the AZ is used or not used for control plane nodes depending on the flag
		found := !ptr.Deref(openStackCluster.Spec.ControlPlaneOmitAvailabilityZone, false)
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

func reconcileBastion(scope *scope.WithLogger, cluster *clusterv1.Cluster, openStackCluster *infrav1.OpenStackCluster) (*ctrl.Result, error) {
	scope.Logger().V(4).Info("Reconciling Bastion")

	clusterResourceName := names.ClusterResourceName(cluster)
	changed, err := resolveBastionResources(scope, clusterResourceName, openStackCluster)
	if err != nil {
		return nil, err
	}
	if changed {
		return &reconcile.Result{}, nil
	}

	// No Bastion defined
	if !openStackCluster.Spec.Bastion.IsEnabled() {
		// Delete any existing bastion
		if openStackCluster.Status.Bastion != nil {
			if err := deleteBastion(scope, cluster, openStackCluster); err != nil {
				return nil, err
			}
			// Reconcile again before continuing
			return &reconcile.Result{}, nil
		}

		// Otherwise nothing to do
		return nil, nil
	}

	computeService, err := compute.NewService(scope)
	if err != nil {
		return nil, err
	}

	networkingService, err := networking.NewService(scope)
	if err != nil {
		return nil, err
	}

	instanceSpec, err := bastionToInstanceSpec(openStackCluster, cluster)
	if err != nil {
		return nil, err
	}

	bastionHash, err := compute.HashInstanceSpec(instanceSpec)
	if err != nil {
		return nil, fmt.Errorf("failed computing bastion hash from instance spec: %w", err)
	}
	if bastionHashHasChanged(bastionHash, openStackCluster.ObjectMeta.Annotations) {
		scope.Logger().Info("Bastion instance spec has changed, deleting existing bastion")

		if err := deleteBastion(scope, cluster, openStackCluster); err != nil {
			return nil, err
		}

		// Add the new annotation and reconcile again before continuing
		annotations.AddAnnotations(openStackCluster, map[string]string{BastionInstanceHashAnnotation: bastionHash})
		return &reconcile.Result{}, nil
	}

	err = getOrCreateBastionPorts(openStackCluster, networkingService)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to get or create ports for bastion: %w", err))
		return nil, fmt.Errorf("failed to get or create ports for bastion: %w", err)
	}
	bastionPortIDs := GetPortIDs(openStackCluster.Status.Bastion.Resources.Ports)

	var instanceStatus *compute.InstanceStatus
	if openStackCluster.Status.Bastion != nil && openStackCluster.Status.Bastion.ID != "" {
		if instanceStatus, err = computeService.GetInstanceStatus(openStackCluster.Status.Bastion.ID); err != nil {
			return nil, err
		}
	}
	if instanceStatus == nil {
		// Check if there is an existing instance with bastion name, in case where bastion ID would not have been properly stored in cluster status
		if instanceStatus, err = computeService.GetInstanceStatusByName(openStackCluster, instanceSpec.Name); err != nil {
			return nil, err
		}
	}
	if instanceStatus == nil {
		instanceStatus, err = computeService.CreateInstance(openStackCluster, instanceSpec, bastionPortIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to create bastion: %w", err)
		}
	}

	// Save hash & status as soon as we know we have an instance
	instanceStatus.UpdateBastionStatus(openStackCluster)

	// Make sure that bastion instance has a valid state
	switch instanceStatus.State() {
	case infrav1.InstanceStateError:
		return nil, fmt.Errorf("failed to reconcile bastion, instance state is ERROR")
	case infrav1.InstanceStateBuild, infrav1.InstanceStateUndefined:
		scope.Logger().Info("Waiting for bastion instance to become ACTIVE", "id", instanceStatus.ID(), "status", instanceStatus.State())
		return &reconcile.Result{RequeueAfter: waitForBuildingInstanceToReconcile}, nil
	case infrav1.InstanceStateDeleted:
		// Not clear why this would happen, so try to clean everything up before reconciling again
		if err := deleteBastion(scope, cluster, openStackCluster); err != nil {
			return nil, err
		}
		return &reconcile.Result{}, nil
	}

	port, err := computeService.GetManagementPort(openStackCluster, instanceStatus)
	if err != nil {
		err = fmt.Errorf("getting management port for bastion: %w", err)
		handleUpdateOSCError(openStackCluster, err)
		return nil, err
	}

	return bastionAddFloatingIP(openStackCluster, clusterResourceName, port, networkingService)
}

func bastionAddFloatingIP(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string, port *ports.Port, networkingService *networking.Service) (*reconcile.Result, error) {
	fp, err := networkingService.GetFloatingIPByPortID(port.ID)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to get or create floating IP for bastion: %w", err))
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
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to get or create floating IP for bastion: %w", err))
		return nil, fmt.Errorf("failed to get or create floating IP for bastion: %w", err)
	}
	openStackCluster.Status.Bastion.FloatingIP = fp.FloatingIP

	err = networkingService.AssociateFloatingIP(openStackCluster, fp, port.ID)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to associate floating IP with bastion: %w", err))
		return nil, fmt.Errorf("failed to associate floating IP with bastion: %w", err)
	}

	return nil, nil
}

func bastionToInstanceSpec(openStackCluster *infrav1.OpenStackCluster, cluster *clusterv1.Cluster) (*compute.InstanceSpec, error) {
	bastion := openStackCluster.Spec.Bastion
	if bastion == nil {
		bastion = &infrav1.Bastion{}
	}
	if bastion.Spec == nil {
		// For the case when Bastion is deleted but we don't have spec, let's use an empty one.
		// v1beta1 API validations prevent this from happening in normal circumstances.
		bastion.Spec = &infrav1.OpenStackMachineSpec{}
	}
	resolved := openStackCluster.Status.Bastion.Resolved
	if resolved == nil {
		return nil, errors.New("bastion resolved is nil")
	}

	machineSpec := bastion.Spec
	instanceSpec := &compute.InstanceSpec{
		Name:          bastionName(cluster.Name),
		Flavor:        machineSpec.Flavor,
		SSHKeyName:    machineSpec.SSHKeyName,
		ImageID:       resolved.ImageID,
		RootVolume:    machineSpec.RootVolume,
		ServerGroupID: resolved.ServerGroupID,
		Tags:          compute.InstanceTags(machineSpec, openStackCluster),
	}
	if bastion.AvailabilityZone != nil {
		instanceSpec.FailureDomain = *bastion.AvailabilityZone
	}
	return instanceSpec, nil
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

func getOrCreateBastionPorts(openStackCluster *infrav1.OpenStackCluster, networkingService *networking.Service) error {
	desiredPorts := openStackCluster.Status.Bastion.Resolved.Ports
	resources := openStackCluster.Status.Bastion.Resources
	if resources == nil {
		return errors.New("bastion resources are nil")
	}

	if len(desiredPorts) == len(resources.Ports) {
		return nil
	}

	err := networkingService.CreatePorts(openStackCluster, desiredPorts, resources)
	if err != nil {
		return fmt.Errorf("failed to create ports for bastion %s: %w", bastionName(openStackCluster.Name), err)
	}

	return nil
}

// bastionHashHasChanged returns a boolean whether if the latest bastion hash, built from the instance spec, has changed or not.
func bastionHashHasChanged(computeHash string, clusterAnnotations map[string]string) bool {
	latestHash, ok := clusterAnnotations[BastionInstanceHashAnnotation]
	if !ok {
		return false
	}
	return latestHash != computeHash
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
				handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to find loadbalancer network: %w", err))
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
					if s.ID != nil && subnetID == *s.ID {
						matchFound = true
						lbNetStatus.Subnets = append(
							lbNetStatus.Subnets, infrav1.Subnet{
								ID: *s.ID,
							})
					}
				}
				if !matchFound {
					handleUpdateOSCError(openStackCluster, fmt.Errorf("no subnet match was found in the specified network (specified subnet: %v, available subnets: %v)", s, lbNet.Subnets))
					return fmt.Errorf("no subnet match was found in the specified network (specified subnet: %v, available subnets: %v)", s, lbNet.Subnets)
				}
			}
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
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile external network: %w", err))
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
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile loadbalancer network: %w", err))
		return fmt.Errorf("failed to reconcile loadbalancer network: %w", err)
	}

	err = networkingService.ReconcileSecurityGroups(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile security groups: %w", err))
		return fmt.Errorf("failed to reconcile security groups: %w", err)
	}

	return reconcileControlPlaneEndpoint(scope, networkingService, openStackCluster, clusterResourceName)
}

// reconcilePreExistingNetworkComponents reconciles the cluster network status when the cluster is
// using pre-existing networks and subnets which are not provisioned by the
// cluster controller.
func reconcilePreExistingNetworkComponents(scope *scope.WithLogger, networkingService *networking.Service, openStackCluster *infrav1.OpenStackCluster) error {
	scope.Logger().V(4).Info("No need to reconcile network, searching network and subnet instead")

	if openStackCluster.Status.Network == nil {
		openStackCluster.Status.Network = &infrav1.NetworkStatusWithSubnets{}
	}

	if openStackCluster.Spec.Network != nil {
		network, err := networkingService.GetNetworkByParam(openStackCluster.Spec.Network)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to find network: %w", err))
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

	return nil
}

func reconcileProvisionedNetworkComponents(networkingService *networking.Service, openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	err := networkingService.ReconcileNetwork(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile network: %w", err))
		return fmt.Errorf("failed to reconcile network: %w", err)
	}
	err = networkingService.ReconcileSubnet(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile subnets: %w", err))
		return fmt.Errorf("failed to reconcile subnets: %w", err)
	}
	err = networkingService.ReconcileRouter(openStackCluster, clusterResourceName)
	if err != nil {
		handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile router: %w", err))
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

		terminalFailure, err := loadBalancerService.ReconcileLoadBalancer(openStackCluster, clusterResourceName, apiServerPort)
		if err != nil {
			// if it's terminalFailure (not Transient), set the Failure reason and message
			if terminalFailure {
				handleUpdateOSCError(openStackCluster, fmt.Errorf("failed to reconcile load balancer: %w", err))
			}
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

	// API server load balancer is disabled, but floating IP is not. Create
	// a floating IP to be attached directly to a control plane host.
	case !ptr.Deref(openStackCluster.Spec.DisableAPIServerFloatingIP, false):
		fp, err := networkingService.GetOrCreateFloatingIP(openStackCluster, openStackCluster, clusterResourceName, openStackCluster.Spec.APIServerFloatingIP)
		if err != nil {
			handleUpdateOSCError(openStackCluster, fmt.Errorf("floating IP cannot be got or created: %w", err))
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
		handleUpdateOSCError(openStackCluster, err)
		return err
	}

	openStackCluster.Spec.ControlPlaneEndpoint = &clusterv1.APIEndpoint{
		Host: host,
		Port: int32(apiServerPort),
	}

	return nil
}

// getAPIServerPort returns the port to use for the API server based on the cluster spec.
func getAPIServerPort(openStackCluster *infrav1.OpenStackCluster) int {
	switch {
	case openStackCluster.Spec.ControlPlaneEndpoint != nil && openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
		return int(openStackCluster.Spec.ControlPlaneEndpoint.Port)
	case openStackCluster.Spec.APIServerPort != nil:
		return *openStackCluster.Spec.APIServerPort
	}
	return 6443
}

func (r *OpenStackClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	clusterToInfraFn := util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind("OpenStackCluster"), mgr.GetClient(), &infrav1.OpenStackCluster{})
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
			builder.WithPredicates(predicates.ClusterUnpaused(ctrl.LoggerFrom(ctx))),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(ctrl.LoggerFrom(ctx))).
		Complete(r)
}

func handleUpdateOSCError(openstackCluster *infrav1.OpenStackCluster, message error) {
	err := capierrors.UpdateClusterError
	openstackCluster.Status.FailureReason = &err
	openstackCluster.Status.FailureMessage = ptr.To(message.Error())
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
			if errors.Is(err, networking.ErrFilterMatch) {
				handleUpdateOSCError(openStackCluster, err)
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
				if errors.Is(err, networking.ErrFilterMatch) {
					handleUpdateOSCError(openStackCluster, err)
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
