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

package vmware

import (
	goctx "context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	topologyv1 "github.com/vmware-tanzu/vm-operator/external/tanzu-topology/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/collections"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

const (
	apiEndpointPort = 6443
)

type ClusterReconciler struct {
	*context.ControllerContext
	NetworkProvider       services.NetworkProvider
	ControlPlaneService   services.ControlPlaneEndpointService
	ResourcePolicyService services.ResourcePolicyService
}

// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclustertemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.com,resources=virtualnetworks;virtualnetworks/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=virtualmachinesetresourcepolicies;virtualmachinesetresourcepolicies/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=virtualmachineservices;virtualmachineservices/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=netoperator.vmware.com,resources=networks,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;update;create;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims/status,verbs=get;update;patch

func (r ClusterReconciler) Reconcile(_ goctx.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	logger := r.Logger.WithName(req.Namespace).WithName(req.Name)
	logger.V(3).Info("Starting Reconcile vsphereCluster")

	// Fetch the vsphereCluster instance
	vsphereCluster := &vmwarev1.VSphereCluster{}
	err := r.Client.Get(r, req.NamespacedName, vsphereCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			r.Logger.V(4).Info("VSphereCluster not found, won't reconcile", "key", req.NamespacedName)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := clusterutilv1.GetOwnerCluster(r, r.Client, vsphereCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Build the patch helper.
	patchHelper, err := patch.NewHelper(vsphereCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to initialize patch helper")
	}

	// Build the cluster context.
	clusterContext := &vmware.ClusterContext{
		ControllerContext: r.ControllerContext,
		Cluster:           cluster,
		Logger:            r.Logger.WithName(vsphereCluster.Namespace).WithName(vsphereCluster.Name),
		VSphereCluster:    vsphereCluster,
		PatchHelper:       patchHelper,
	}

	// Always close the context when exiting this function so we can persist any vsphereCluster changes.
	defer func() {
		r.Recorder.EmitEvent(vsphereCluster, "Reconcile", reterr, true)
		if err := clusterContext.Patch(); err != nil {
			clusterContext.Logger.Error(err, "failed to patch vspherecluster")
			if reterr == nil {
				reterr = err
			}
		}
	}()

	// Handle deleted clusters
	if !vsphereCluster.DeletionTimestamp.IsZero() {
		r.reconcileDelete(clusterContext)
		return ctrl.Result{}, nil
	}

	if cluster == nil {
		logger.V(2).Info("waiting on Cluster controller to set OwnerRef on infra cluster")
		return reconcile.Result{}, nil
	}

	// If the VSphereCluster doesn't have our finalizer, add it.
	// Requeue immediately after adding finalizer to avoid the race condition between init and delete
	if !controllerutil.ContainsFinalizer(vsphereCluster, vmwarev1.ClusterFinalizer) {
		controllerutil.AddFinalizer(vsphereCluster, vmwarev1.ClusterFinalizer)
		return ctrl.Result{}, nil
	}

	// Handle non-deleted clusters
	return ctrl.Result{}, r.reconcileNormal(clusterContext)
}

func (r *ClusterReconciler) reconcileDelete(ctx *vmware.ClusterContext) {
	ctx.Logger.Info("Reconciling vsphereCluster delete")

	deletingConditionTypes := []clusterv1.ConditionType{
		vmwarev1.ResourcePolicyReadyCondition,
		vmwarev1.ClusterNetworkReadyCondition,
		vmwarev1.LoadBalancerReadyCondition,
	}

	for _, t := range deletingConditionTypes {
		if c := conditions.Get(ctx.VSphereCluster, t); c != nil {
			conditions.MarkFalse(ctx.VSphereCluster, t, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
		}
	}

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(ctx.VSphereCluster, vmwarev1.ClusterFinalizer)
}

func (r *ClusterReconciler) reconcileNormal(ctx *vmware.ClusterContext) error {
	ctx.Logger.Info("Reconciling vsphereCluster")

	// Get any failure domains to report back to the CAPI core controller.
	failureDomains, err := r.getFailureDomains(ctx)
	if err != nil {
		return errors.Wrapf(
			err,
			"unexpected error while discovering failure domains for %s", ctx.VSphereCluster.Name)
	}
	ctx.VSphereCluster.Status.FailureDomains = failureDomains

	// Reconcile ResourcePolicy before we create the machines. If the ResourcePolicy is not reconciled before we create the Node VMs,
	//   it will be handled by vm operator by relocating the VMs to the ResourcePool and Folder specified by the ResourcePolicy.
	// Reconciling the ResourcePolicy early potentially saves us the extra relocate operation.
	resourcePolicyName, err := r.ResourcePolicyService.ReconcileResourcePolicy(ctx)
	if err != nil {
		conditions.MarkFalse(ctx.VSphereCluster, vmwarev1.ResourcePolicyReadyCondition, vmwarev1.ResourcePolicyCreationFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		return errors.Wrapf(err,
			"failed to configure resource policy for vsphereCluster %s/%s",
			ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name)
	}
	conditions.MarkTrue(ctx.VSphereCluster, vmwarev1.ResourcePolicyReadyCondition)
	ctx.VSphereCluster.Status.ResourcePolicyName = resourcePolicyName

	// Configure the cluster for the cluster network
	err = r.NetworkProvider.ProvisionClusterNetwork(ctx)
	if err != nil {
		return errors.Wrapf(err,
			"failed to configure cluster network for vsphereCluster %s/%s",
			ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name)
	}

	if ok, err := r.reconcileControlPlaneEndpoint(ctx); !ok {
		if err != nil {
			return errors.Wrapf(err, "unexpected error while reconciling control plane endpoint for %s", ctx.VSphereCluster.Name)
		}
	}

	ctx.VSphereCluster.Status.Ready = true
	ctx.Logger.V(2).Info("Reconcile completed, vsphereCluster is infrastructure-ready")
	return nil
}

func (r *ClusterReconciler) reconcileControlPlaneEndpoint(ctx *vmware.ClusterContext) (bool, error) {
	if !ctx.Cluster.Spec.ControlPlaneEndpoint.IsZero() {
		ctx.VSphereCluster.Spec.ControlPlaneEndpoint.Host = ctx.Cluster.Spec.ControlPlaneEndpoint.Host
		ctx.VSphereCluster.Spec.ControlPlaneEndpoint.Port = ctx.Cluster.Spec.ControlPlaneEndpoint.Port
		if r.NetworkProvider.HasLoadBalancer() {
			conditions.MarkTrue(ctx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition)
		}
		ctx.Logger.Info("skipping control plane endpoint reconciliation",
			"reason", "ControlPlaneEndpoint already set on Cluster",
			"controlPlaneEndpoint", ctx.Cluster.Spec.ControlPlaneEndpoint.String())
		return true, nil
	}

	if !ctx.VSphereCluster.Spec.ControlPlaneEndpoint.IsZero() {
		if r.NetworkProvider.HasLoadBalancer() {
			conditions.MarkTrue(ctx.VSphereCluster, vmwarev1.LoadBalancerReadyCondition)
		}
		ctx.Logger.Info("skipping control plane endpoint reconciliation",
			"reason", "ControlPlaneEndpoint already set on vsphereCluster",
			"controlPlaneEndpoint", ctx.VSphereCluster.Spec.ControlPlaneEndpoint.String())
		return true, nil
	}

	if r.NetworkProvider.HasLoadBalancer() {
		if err := r.reconcileLoadBalancedEndpoint(ctx); err != nil {
			return false, errors.Wrapf(err,
				"failed to reconcile loadbalanced endpoint for vsphereCluster %s/%s",
				ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name)
		}

		return true, nil
	}

	if err := r.reconcileAPIEndpoints(ctx); err != nil {
		return false, errors.Wrapf(err,
			"failed to reconcile API endpoints for vsphereCluster %s/%s",
			ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name)
	}

	return true, nil
}

func (r *ClusterReconciler) reconcileLoadBalancedEndpoint(ctx *vmware.ClusterContext) error {
	ctx.Logger.Info("Reconciling load-balanced control plane endpoint")

	// Will create a VirtualMachineService for a NetworkProvider that supports load balancing
	cpEndpoint, err := r.ControlPlaneService.ReconcileControlPlaneEndpointService(ctx, r.NetworkProvider)
	if err != nil {
		// Likely the endpoint is not ready. Keep retrying.
		return errors.Wrapf(err,
			"failed to get control plane endpoint for Cluster %s/%s",
			ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name)
	}

	if cpEndpoint == nil {
		return fmt.Errorf("control plane endpoint not available for Cluster %s/%s",
			ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name)
	}

	// If we've got here and we have a cpEndpoint, we're done.
	ctx.VSphereCluster.Spec.ControlPlaneEndpoint = *cpEndpoint
	ctx.Logger.V(3).Info(
		"found API endpoint via virtual machine service",
		"host", cpEndpoint.Host,
		"port", cpEndpoint.Port)
	return nil
}

func (r *ClusterReconciler) reconcileAPIEndpoints(ctx *vmware.ClusterContext) error {
	ctx.Logger.Info("Reconciling control plane endpoint")
	machines, err := collections.GetFilteredMachinesForCluster(ctx, r.Client, ctx.Cluster, collections.ControlPlaneMachines(ctx.Cluster.Name))
	if err != nil {
		return errors.Wrapf(err,
			"failed to get Machines for Cluster %s/%s",
			ctx.Cluster.Namespace, ctx.Cluster.Name)
	}

	// Define a variable to assign the API endpoints of control plane
	// machines as they are discovered.
	var apiEndpointList []clusterv1.APIEndpoint //nolint:prealloc

	// Iterate over the cluster's control plane CAPI machines.
	for _, machine := range machines {
		// Only machines with bootstrap data will have an IP address.
		if machine.Spec.Bootstrap.DataSecretName == nil {
			ctx.Logger.V(5).Info(
				"skipping machine while looking for IP address",
				"reason", "bootstrap.DataSecretName is nil",
				"machine-name", machine.Name)
			continue
		}

		// Get the vsphereMachine for the CAPI Machine resource.
		vsphereMachine, err := util.GetVSphereMachine(ctx, ctx.Client, machine.Namespace, machine.Name)
		if err != nil {
			return errors.Wrapf(err,
				"failed to get vsphereMachine for Machine %s/%s/%s",
				ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name, machine.Name)
		}

		// If the machine has no IP address then skip it.
		if vsphereMachine.Status.IPAddr == "" {
			ctx.Logger.V(5).Info("skipping machine without IP address")
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
		ctx.Logger.V(3).Info(
			"found API endpoint via control plane machine",
			"host", apiEndpoint.Host,
			"port", apiEndpoint.Port)
	}

	// The reconciliation is only successful if some API endpoints were
	// discovered. Otherwise return an error so the cluster is requeued
	// for reconciliation.
	if len(apiEndpointList) == 0 {
		return errors.Wrapf(err,
			"failed to reconcile API endpoints for %s/%s",
			ctx.VSphereCluster.Namespace, ctx.VSphereCluster.Name)
	}

	// Update the vsphereCluster's list of APIEndpoints.
	ctx.VSphereCluster.Spec.ControlPlaneEndpoint = apiEndpointList[0]

	return nil
}

func (r *ClusterReconciler) VSphereMachineToCluster(ctx goctx.Context, o client.Object) []reconcile.Request {
	vsphereMachine, ok := o.(*vmwarev1.VSphereMachine)
	if !ok {
		r.Logger.Error(errors.New("did not get vspheremachine"), "got", fmt.Sprintf("%T", o))
		return nil
	}
	if !util.IsControlPlaneMachine(vsphereMachine) {
		r.Logger.V(5).Info("rejecting vsphereCluster reconcile as not CP machine", "machineName", vsphereMachine.Name)
		return nil
	}
	// Only currently interested in updating Cluster from vsphereMachines with IP addresses
	if vsphereMachine.Status.IPAddr == "" {
		r.Logger.V(5).Info("rejecting vsphereCluster reconcile as no IP address", "machineName", vsphereMachine.Name)
		return nil
	}

	cluster, err := util.GetVSphereClusterFromVMwareMachine(ctx, r.Client, vsphereMachine)
	if err != nil {
		r.Logger.Error(err, "failed to get cluster", "machine", vsphereMachine.Name, "namespace", vsphereMachine.Namespace)
		return nil
	}

	// Can add further filters on Cluster state so that we don't keep reconciling Cluster
	r.Logger.V(3).Info("triggering VSphereCluster reconcile from VSphereMachine", "machineName", vsphereMachine.Name)
	return []ctrl.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: cluster.Namespace,
			Name:      cluster.Name,
		},
	}}
}

var isFaultDomainsFSSEnabled = func() bool {
	return os.Getenv("FSS_WCP_FAULTDOMAINS") == "true"
}

// Returns the failure domain information discovered on the cluster
// hosting this controller.
func (r *ClusterReconciler) getFailureDomains(ctx *vmware.ClusterContext) (clusterv1.FailureDomains, error) {
	if !isFaultDomainsFSSEnabled() {
		return nil, nil
	}

	availabilityZoneList := &topologyv1.AvailabilityZoneList{}
	if err := ctx.Client.List(ctx, availabilityZoneList); err != nil {
		return nil, err
	}

	if len(availabilityZoneList.Items) == 0 {
		return nil, nil
	}

	failureDomains := clusterv1.FailureDomains{}
	for _, az := range availabilityZoneList.Items {
		failureDomains[az.Name] = clusterv1.FailureDomainSpec{
			ControlPlane: true,
		}
	}

	return failureDomains, nil
}
