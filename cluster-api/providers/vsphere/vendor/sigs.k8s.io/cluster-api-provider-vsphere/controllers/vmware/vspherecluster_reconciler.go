/*
Copyright 2021 The Kubernetes Authors.

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

// Package vmware contains the VSphereCluster reconciler.
package vmware

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/collections"
	"sigs.k8s.io/cluster-api/util/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/conditions/v1beta2"
	"sigs.k8s.io/cluster-api/util/finalizers"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/paused"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	topologyv1 "sigs.k8s.io/cluster-api-provider-vsphere/internal/apis/topology/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

const (
	apiEndpointPort = 6443
)

// ClusterReconciler reconciles VSphereClusters.
type ClusterReconciler struct {
	Client                client.Client
	Recorder              record.EventRecorder
	NetworkProvider       services.NetworkProvider
	ControlPlaneService   services.ControlPlaneEndpointService
	ResourcePolicyService services.ResourcePolicyService
}

// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclustertemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crd.nsx.vmware.com,resources=subnetsets;subnetsets/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.com,resources=virtualnetworks;virtualnetworks/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=virtualmachinesetresourcepolicies;virtualmachinesetresourcepolicies/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=virtualmachineservices;virtualmachineservices/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=netoperator.vmware.com,resources=networks,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;update;create;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims/status,verbs=get;update;patch

func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the vsphereCluster instance.
	vsphereCluster := &vmwarev1.VSphereCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, vsphereCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, vsphereCluster, vmwarev1.ClusterFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := clusterutilv1.GetOwnerCluster(ctx, r.Client, vsphereCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster != nil {
		log = log.WithValues("Cluster", klog.KObj(cluster))
		ctx = ctrl.LoggerInto(ctx, log)
	}

	// Build the patch helper.
	patchHelper, err := patch.NewHelper(vsphereCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}

	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, vsphereCluster); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	// Build the cluster context.
	clusterContext := &vmware.ClusterContext{
		Cluster:        cluster,
		VSphereCluster: vsphereCluster,
		PatchHelper:    patchHelper,
	}

	// Always close the context when exiting this function so we can persist any vsphereCluster changes.
	defer func() {
		if err := r.patch(ctx, clusterContext); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Handle deleted clusters
	if !vsphereCluster.DeletionTimestamp.IsZero() {
		r.reconcileDelete(clusterContext)
		return ctrl.Result{}, nil
	}

	if cluster == nil {
		log.Info("Waiting for Cluster controller to set OwnerRef on VSphereCluster")
		return reconcile.Result{}, nil
	}

	// Handle non-deleted clusters
	return ctrl.Result{}, r.reconcileNormal(ctx, clusterContext)
}

func (r *ClusterReconciler) patch(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	// always update the readyCondition.
	conditions.SetSummary(clusterCtx.VSphereCluster,
		conditions.WithConditions(
			vmwarev1.ResourcePolicyReadyCondition,
			vmwarev1.ClusterNetworkReadyCondition,
			vmwarev1.LoadBalancerReadyCondition,
		),
	)

	if err := v1beta2conditions.SetSummaryCondition(clusterCtx.VSphereCluster, clusterCtx.VSphereCluster, vmwarev1.VSphereClusterReadyV1Beta2Condition,
		v1beta2conditions.ForConditionTypes{
			vmwarev1.VSphereClusterResourcePolicyReadyV1Beta2Condition,
			vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
			vmwarev1.VSphereClusterLoadBalancerReadyV1Beta2Condition,
			// ProviderServiceAccountsReady and ServiceDiscoveryReady will be set by other controllers after
			// the API server in the workload cluster is up and running.
			vmwarev1.VSphereClusterProviderServiceAccountsReadyV1Beta2Condition,
			vmwarev1.VSphereClusterServiceDiscoveryReadyV1Beta2Condition,
		},
		v1beta2conditions.IgnoreTypesIfMissing{
			vmwarev1.VSphereClusterProviderServiceAccountsReadyV1Beta2Condition,
			vmwarev1.VSphereClusterServiceDiscoveryReadyV1Beta2Condition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		v1beta2conditions.CustomMergeStrategy{
			MergeStrategy: v1beta2conditions.DefaultMergeStrategy(
				// Use custom reasons.
				v1beta2conditions.ComputeReasonFunc(v1beta2conditions.GetDefaultComputeMergeReasonFunc(
					vmwarev1.VSphereClusterNotReadyV1Beta2Reason,
					vmwarev1.VSphereClusterReadyUnknownV1Beta2Reason,
					vmwarev1.VSphereClusterReadyV1Beta2Reason,
				)),
			),
		},
	); err != nil {
		return errors.Wrapf(err, "failed to set %s condition", vmwarev1.VSphereClusterReadyV1Beta2Condition)
	}

	return clusterCtx.PatchHelper.Patch(ctx, clusterCtx.VSphereCluster,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			vmwarev1.ResourcePolicyReadyCondition,
			vmwarev1.ClusterNetworkReadyCondition,
			vmwarev1.LoadBalancerReadyCondition,
		}},
		patch.WithOwnedV1Beta2Conditions{Conditions: []string{
			clusterv1.PausedV1Beta2Condition,
			vmwarev1.VSphereClusterReadyV1Beta2Condition,
			vmwarev1.VSphereClusterResourcePolicyReadyV1Beta2Condition,
			vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
			vmwarev1.VSphereClusterLoadBalancerReadyV1Beta2Condition,
			// NOTE: ProviderServiceAccountsReady and ServiceDiscoveryReady are not owned by this controller
		}},
	)
}

func (r *ClusterReconciler) reconcileDelete(clusterCtx *vmware.ClusterContext) {
	deletingConditionTypes := []clusterv1.ConditionType{
		vmwarev1.ResourcePolicyReadyCondition,
		vmwarev1.ClusterNetworkReadyCondition,
		vmwarev1.LoadBalancerReadyCondition,
	}

	for _, t := range deletingConditionTypes {
		if c := conditions.Get(clusterCtx.VSphereCluster, t); c != nil {
			conditions.MarkFalse(clusterCtx.VSphereCluster, t, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
		}
	}

	v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterResourcePolicyReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: vmwarev1.VSphereClusterResourcePolicyReadyDeletingV1Beta2Reason,
	})

	v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: vmwarev1.VSphereClusterNetworkReadyDeletingV1Beta2Reason,
	})

	v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterLoadBalancerReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: vmwarev1.VSphereClusterLoadBalancerDeletingV1Beta2Reason,
	})

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(clusterCtx.VSphereCluster, vmwarev1.ClusterFinalizer)
}

func (r *ClusterReconciler) reconcileNormal(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	// Get any failure domains to report back to the CAPI core controller.
	failureDomains, err := r.getFailureDomains(ctx, clusterCtx.VSphereCluster.Namespace)
	if err != nil {
		return errors.Wrapf(
			err,
			"unexpected error while discovering failure domains for %s", clusterCtx.VSphereCluster.Name)
	}
	clusterCtx.VSphereCluster.Status.FailureDomains = failureDomains

	// Reconcile ResourcePolicy before we create the machines. If the ResourcePolicy is not reconciled before we create the Node VMs,
	// it will be handled by vm operator by relocating the VMs to the ResourcePool and Folder specified by the ResourcePolicy.
	// Reconciling the ResourcePolicy early potentially saves us the extra relocate operation.
	resourcePolicyName, err := r.ResourcePolicyService.ReconcileResourcePolicy(ctx, clusterCtx)
	if err != nil {
		conditions.MarkFalse(clusterCtx.VSphereCluster, vmwarev1.ResourcePolicyReadyCondition, vmwarev1.ResourcePolicyCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:    vmwarev1.VSphereClusterResourcePolicyReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  vmwarev1.VSphereClusterResourcePolicyNotReadyV1Beta2Reason,
			Message: err.Error(),
		})
		return errors.Wrapf(err,
			"failed to configure resource policy for vsphereCluster %s/%s",
			clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
	}
	conditions.MarkTrue(clusterCtx.VSphereCluster, vmwarev1.ResourcePolicyReadyCondition)
	v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterResourcePolicyReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: vmwarev1.VSphereClusterResourcePolicyReadyV1Beta2Reason,
	})

	clusterCtx.VSphereCluster.Status.ResourcePolicyName = resourcePolicyName

	// Configure the cluster for the cluster network
	err = r.NetworkProvider.ProvisionClusterNetwork(ctx, clusterCtx)
	if err != nil {
		return errors.Wrapf(err,
			"failed to configure cluster network for VSphereCluster %s/%s",
			clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
	}

	if err := r.reconcileControlPlaneEndpoint(ctx, clusterCtx); err != nil {
		return errors.Wrapf(err, "unexpected error while reconciling control plane endpoint for %s", clusterCtx.VSphereCluster.Name)
	}

	clusterCtx.VSphereCluster.Status.Ready = true
	return nil
}

func (r *ClusterReconciler) reconcileControlPlaneEndpoint(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	log := ctrl.LoggerFrom(ctx)

	if !clusterCtx.Cluster.Spec.ControlPlaneEndpoint.IsZero() {
		clusterCtx.VSphereCluster.Spec.ControlPlaneEndpoint.Host = clusterCtx.Cluster.Spec.ControlPlaneEndpoint.Host
		clusterCtx.VSphereCluster.Spec.ControlPlaneEndpoint.Port = clusterCtx.Cluster.Spec.ControlPlaneEndpoint.Port
		v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:   vmwarev1.VSphereClusterLoadBalancerReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: vmwarev1.VSphereClusterLoadBalancerReadyV1Beta2Reason,
		})
		if r.NetworkProvider.HasLoadBalancer() {
			conditions.MarkTrue(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition)
		}
		log.Info("Skipping control plane endpoint reconciliation",
			"reason", "ControlPlaneEndpoint already set on Cluster",
			"controlPlaneEndpoint", clusterCtx.Cluster.Spec.ControlPlaneEndpoint.String())
		return nil
	}

	if !clusterCtx.VSphereCluster.Spec.ControlPlaneEndpoint.IsZero() {
		v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
			Type:   vmwarev1.VSphereClusterLoadBalancerReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: vmwarev1.VSphereClusterLoadBalancerReadyV1Beta2Reason,
		})
		if r.NetworkProvider.HasLoadBalancer() {
			conditions.MarkTrue(clusterCtx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition)
		}
		log.Info("Skipping control plane endpoint reconciliation",
			"reason", "ControlPlaneEndpoint already set on VSphereCluster",
			"controlPlaneEndpoint", clusterCtx.VSphereCluster.Spec.ControlPlaneEndpoint.String())
		return nil
	}

	if r.NetworkProvider.HasLoadBalancer() {
		if err := r.reconcileLoadBalancedEndpoint(ctx, clusterCtx); err != nil {
			return errors.Wrapf(err,
				"failed to reconcile loadbalanced endpoint for VSphereCluster %s/%s",
				clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
		}

		return nil
	}

	if err := r.reconcileAPIEndpoints(ctx, clusterCtx); err != nil {
		return errors.Wrapf(err,
			"failed to reconcile API endpoints for VSphereCluster %s/%s",
			clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
	}

	return nil
}

func (r *ClusterReconciler) reconcileLoadBalancedEndpoint(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	log := ctrl.LoggerFrom(ctx)

	// Will create a VirtualMachineService for a NetworkProvider that supports load balancing
	cpEndpoint, err := r.ControlPlaneService.ReconcileControlPlaneEndpointService(ctx, clusterCtx, r.NetworkProvider)
	if err != nil {
		// Likely the endpoint is not ready. Keep retrying.
		return errors.Wrapf(err,
			"failed to get control plane endpoint for VSphereCluster %s/%s",
			clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
	}

	if cpEndpoint == nil {
		return fmt.Errorf("control plane endpoint not available for VSphereCluster %s/%s",
			clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
	}

	// If we've got here and we have a cpEndpoint, we're done.
	clusterCtx.VSphereCluster.Spec.ControlPlaneEndpoint = *cpEndpoint
	log.V(4).Info("Found API endpoint via virtual machine service", "host", cpEndpoint.Host, "port", cpEndpoint.Port)
	return nil
}

func (r *ClusterReconciler) reconcileAPIEndpoints(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	log := ctrl.LoggerFrom(ctx)

	machines, err := collections.GetFilteredMachinesForCluster(ctx, r.Client, clusterCtx.Cluster, collections.ControlPlaneMachines(clusterCtx.Cluster.Name))
	if err != nil {
		return errors.Wrapf(err,
			"failed to get Machines for Cluster %s/%s",
			clusterCtx.Cluster.Namespace, clusterCtx.Cluster.Name)
	}

	// Define a variable to assign the API endpoints of control plane
	// machines as they are discovered.
	apiEndpointList := []clusterv1.APIEndpoint{}

	// Iterate over the cluster's control plane CAPI machines.
	for _, machine := range machines {
		// Note: We have to use := here to create a new variable and not overwrite log & ctx outside the for loop.
		log := log.WithValues("Machine", klog.KObj(machine))
		ctx := ctrl.LoggerInto(ctx, log)

		// Only machines with bootstrap data will have an IP address.
		if machine.Spec.Bootstrap.DataSecretName == nil {
			log.V(4).Info("Skipping Machine while looking for IP address", "reason", "bootstrap.DataSecretName is nil")
			continue
		}

		// Get the vsphereMachine for the CAPI Machine resource.
		vsphereMachine, err := util.GetVSphereMachine(ctx, r.Client, machine.Namespace, machine.Name)
		if err != nil {
			return errors.Wrapf(err, "failed to get VSphereMachine for Machine %s/%s", machine.Namespace, machine.Name)
		}
		log = log.WithValues("VSphereMachine", klog.KObj(vsphereMachine))
		ctx = ctrl.LoggerInto(ctx, log) //nolint:ineffassign,staticcheck // ensure the logger is up-to-date in ctx, even if we currently don't use ctx below.

		// If the machine has no IP address then skip it.
		if vsphereMachine.Status.IPAddr == "" {
			log.V(4).Info("Skipping Machine without IP address")
			continue
		}

		// Append the control plane machine's IP address to the list of API
		// endpoints for this cluster so that they can be read into the
		// analogous CAPI cluster via an unstructured reader.
		apiEndpoint := clusterv1.APIEndpoint{
			Host: vsphereMachine.Status.IPAddr,
			Port: apiEndpointPort,
		}
		apiEndpointList = append(apiEndpointList, apiEndpoint)
		log.V(4).Info("Found API endpoint via control plane machine", "host", apiEndpoint.Host, "port", apiEndpoint.Port)
	}

	// The reconciliation is only successful if some API endpoints were
	// discovered. Otherwise return an error so the cluster is requeued
	// for reconciliation.
	if len(apiEndpointList) == 0 {
		return errors.Wrapf(err,
			"failed to reconcile API endpoints for %s/%s",
			clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
	}

	// Update the VSphereCluster's list of APIEndpoints.
	clusterCtx.VSphereCluster.Spec.ControlPlaneEndpoint = apiEndpointList[0]

	return nil
}

// VSphereMachineToCluster adds reconcile requests for a Cluster when one of its control plane machines has an event.
func (r *ClusterReconciler) VSphereMachineToCluster(ctx context.Context, o client.Object) []reconcile.Request {
	log := ctrl.LoggerFrom(ctx)

	vsphereMachine, ok := o.(*vmwarev1.VSphereMachine)
	if !ok {
		log.Error(nil, fmt.Sprintf("Expected a VSphereMachine but got a %T", o))
		return nil
	}
	log = log.WithValues("VSphereMachine", klog.KObj(vsphereMachine))
	ctx = ctrl.LoggerInto(ctx, log)

	if !util.IsControlPlaneMachine(vsphereMachine) {
		log.V(6).Info("Skipping VSphereCluster reconcile as Machine is not a control plane Machine")
		return nil
	}

	// Only currently interested in updating Cluster from VSphereMachines with IP addresses
	if vsphereMachine.Status.IPAddr == "" {
		log.V(6).Info("Skipping VSphereCluster reconcile as Machine does not have an IP address")
		return nil
	}

	vsphereCluster, err := util.GetVSphereClusterFromVMwareMachine(ctx, r.Client, vsphereMachine)
	if err != nil {
		log.V(4).Error(err, "Failed to get VSphereCluster from VSphereMachine")
		return nil
	}

	// Can add further filters on Cluster state so that we don't keep reconciling Cluster
	log.V(6).Info("Triggering VSphereCluster reconcile from VSphereMachine")
	return []ctrl.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: vsphereCluster.Namespace,
			Name:      vsphereCluster.Name,
		},
	}}
}

// ZoneToVSphereClusters adds reconcile requests for VSphereClusters when Zone has an event.
func (r *ClusterReconciler) ZoneToVSphereClusters(ctx context.Context, o client.Object) []reconcile.Request {
	log := ctrl.LoggerFrom(ctx)

	zone, ok := o.(*topologyv1.Zone)
	if !ok {
		log.Error(nil, fmt.Sprintf("Expected a Zone but got a %T", o))
		return nil
	}
	log = log.WithValues("Zone", klog.KObj(zone))
	ctx = ctrl.LoggerInto(ctx, log)

	vsphereClusters := &vmwarev1.VSphereClusterList{}
	err := r.Client.List(ctx, vsphereClusters, &client.ListOptions{Namespace: zone.Namespace})
	if err != nil {
		log.V(4).Error(err, "Failed to get VSphereClusters from Zone")
		return nil
	}

	log.V(6).Info("Triggering VSphereCluster reconcile for Zone")
	requests := []reconcile.Request{}
	for _, c := range vsphereClusters.Items {
		r := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      c.Name,
				Namespace: c.Namespace,
			},
		}
		requests = append(requests, r)
	}

	return requests
}

// Returns the failure domain information discovered on the cluster
// hosting this controller.
func (r *ClusterReconciler) getFailureDomains(ctx context.Context, namespace string) (clusterv1.FailureDomains, error) {
	failureDomains := clusterv1.FailureDomains{}
	// Determine the source of failure domain based on feature gates NamespaceScopedZones.
	// If NamespaceScopedZones is enabled, use Zone which is Namespace scoped,otherwise use
	// Availability Zone which is Cluster scoped.
	if feature.Gates.Enabled(feature.NamespaceScopedZones) {
		zoneList := &topologyv1.ZoneList{}
		listOptions := &client.ListOptions{Namespace: namespace}
		if err := r.Client.List(ctx, zoneList, listOptions); err != nil {
			return nil, errors.Wrapf(err, "failed to list Zones in namespace %s", namespace)
		}

		for _, zone := range zoneList.Items {
			// Skip zones which are in deletion
			if !zone.DeletionTimestamp.IsZero() {
				continue
			}
			failureDomains[zone.Name] = clusterv1.FailureDomainSpec{ControlPlane: true}
		}

		if len(failureDomains) == 0 {
			return nil, nil
		}

		return failureDomains, nil
	}
	availabilityZoneList := &topologyv1.AvailabilityZoneList{}
	if err := r.Client.List(ctx, availabilityZoneList); err != nil {
		return nil, err
	}

	if len(availabilityZoneList.Items) == 0 {
		return nil, nil
	}
	for _, az := range availabilityZoneList.Items {
		failureDomains[az.Name] = clusterv1.FailureDomainSpec{
			ControlPlane: true,
		}
	}

	return failureDomains, nil
}
