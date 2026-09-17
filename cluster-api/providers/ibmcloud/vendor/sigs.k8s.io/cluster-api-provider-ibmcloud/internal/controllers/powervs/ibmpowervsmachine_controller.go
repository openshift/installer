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

package powervs

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	"sigs.k8s.io/cluster-api/util/finalizers"
	clog "sigs.k8s.io/cluster-api/util/log"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/paused"
	"sigs.k8s.io/cluster-api/util/predicates"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	powervsscope "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/scope/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
)

// loadBalancerSettleRequeueInterval is the poll interval while waiting for a VPC load balancer pool
// update to leave the update_pending state. IBM VPC load balancers typically settle in 10–30 s;
// 15 s avoids hammering the API while still reacting promptly.
const loadBalancerSettleRequeueInterval = 15 * time.Second

// IBMPowerVSMachineReconciler reconciles a IBMPowerVSMachine object.
type IBMPowerVSMachineReconciler struct {
	client.Client
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme

	// WatchFilterValue is the label value used to filter events prior to reconciliation.
	WatchFilterValue string

	ClientBuilder powervsscope.ClientBuilder
}

// dhcpCacheStore is a cache store to hold the Power VS VM DHCP IP.
var dhcpCacheStore cache.Store

func init() {
	dhcpCacheStore = powervs.InitialiseDHCPCacheStore()
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines/status,verbs=get;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconcileation logic for IBMPowerVSMachine.
func (r *IBMPowerVSMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMPowerVSMachine")
	defer log.Info("Finished reconciling IBMPowerVSMachine")

	// 1. Fetch the IBMPowerVSMachine instance.
	ibmPowerVSMachine := &infrav1.IBMPowerVSMachine{}
	err := r.Client.Get(ctx, req.NamespacedName, ibmPowerVSMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("IBMPowerVSMachine not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSMachine: %w", err)
	}

	// 2. Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmPowerVSMachine, infrav1.IBMPowerVSMachineFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	// 3. Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, ibmPowerVSMachine.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get machine for IBMPowerVSMachine: %w", err)
	}
	if machine == nil {
		log.Info("Waiting for machine controller to set owner ref on IBMPowerVSMachine")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("Machine", klog.KObj(machine))
	ctx = ctrl.LoggerInto(ctx, log)

	// AddOwners adds the owners of IBMPowerVSMachine as k/v pairs to the logger.
	if ctx, log, err = clog.AddOwners(ctx, r.Client, ibmPowerVSMachine); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to add owners to log: %w", err)
	}

	// 4. Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Info("IBMPowerVSMachine owner Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}
	if cluster == nil {
		log.Info(fmt.Sprintf("Please associate this machine with a cluster using the label %s: <name of cluster>", clusterv1.ClusterNameLabel))
		return ctrl.Result{}, nil
	}

	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	// 5. Ensure Infrastructure is defined.
	if !cluster.Spec.InfrastructureRef.IsDefined() {
		log.Info("Cluster infrastructureRef is not available yet")
		return ctrl.Result{}, nil
	}

	// 6. Check if cluster is paused
	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, ibmPowerVSMachine); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	// 7. Fetch the IBMPowerVSCluster.
	ibmPowerVSCluster := &infrav1.IBMPowerVSCluster{}
	ibmPowerVSClusterName := client.ObjectKey{
		Namespace: ibmPowerVSMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, ibmPowerVSClusterName, ibmPowerVSCluster); err != nil {
		log.Info("IBMPowerVSCluster is not available yet")
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSCluster: %w", err)
	}

	log = log.WithValues(ibmPowerVSClusterKind, klog.KObj(ibmPowerVSCluster))
	ctx = ctrl.LoggerInto(ctx, log)

	// 8. Fetch the IBMPowerVSImage.
	var ibmPowerVSImage *infrav1.IBMPowerVSImage
	if ibmPowerVSMachine.Spec.Image.Type == infrav1.ImageSourceTypeImport && ibmPowerVSMachine.Spec.Image.Import.Name != "" {
		ibmPowerVSImage = &infrav1.IBMPowerVSImage{}
		ibmPowerVSImageName := client.ObjectKey{
			Namespace: ibmPowerVSMachine.Namespace,
			Name:      ibmPowerVSMachine.Spec.Image.Import.Name,
		}
		if err := r.Client.Get(ctx, ibmPowerVSImageName, ibmPowerVSImage); err != nil {
			log.Info("IBMPowerVSImage is not available yet", "imageName", ibmPowerVSImageName.Name)
			return ctrl.Result{}, nil
		}
	}

	// 9. Create the machine scope.
	machineScope, err := powervsscope.NewMachineScope(ctx, powervsscope.MachineScopeParams{
		Client:            r.Client,
		Cluster:           cluster,
		IBMPowerVSCluster: ibmPowerVSCluster,
		Machine:           machine,
		IBMPowerVSMachine: ibmPowerVSMachine,
		IBMPowerVSImage:   ibmPowerVSImage,
		ServiceEndpoint:   r.ServiceEndpoint,
		DHCPIPCacheStore:  dhcpCacheStore,
		ClientBuilder:     r.ClientBuilder,
		Recorder:          r.Recorder,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create IBMPowerVS machine scope: %w", err)
	}

	// 10. Initialize the patch helper
	patchHelper, err := patch.NewHelper(ibmPowerVSMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Always attempt to Patch the IBMPowerVSMachine object and status after each reconciliation.
	defer func() {
		if err := patchIBMPowerVSMachine(ctx, patchHelper, ibmPowerVSMachine); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// 11. Handle deleted machines.
	if !ibmPowerVSMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machineScope)
	}

	// 12. Handle non-deleted machines.
	return r.reconcileNormal(ctx, machineScope)
}

func (r *IBMPowerVSMachineReconciler) reconcileDelete(ctx context.Context, scope *powervsscope.MachineScope) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// 1. Signal that deletion is in progress
	r.markCondition(scope, metav1.ConditionFalse, infrav1.InstanceDeletingReason, "")

	// 2. Ensure finalizer is removed only on complete success
	defer func() {
		if reterr == nil {
			log.Info("PowerVS machine deleted, removing finalizer")
			controllerutil.RemoveFinalizer(scope.IBMPowerVSMachine, infrav1.IBMPowerVSMachineFinalizer)
		}
	}()

	// 3. Early exit if the instance was never actually provisioned in IBM Cloud
	if scope.IBMPowerVSMachine.Status.InstanceID == "" {
		log.Info("IBMPowerVSMachine instance ID is not yet set, skipping PowerVS API deletion")
		return ctrl.Result{}, nil
	}

	// 4. Delete the VM from PowerVS
	if err := scope.DeleteMachine(ctx); err != nil {
		log.Error(err, "error deleting IBMPowerVSMachine")
		r.markCondition(scope, metav1.ConditionFalse, infrav1.InstanceDeletingReason, fmt.Sprintf("failed to delete instance: %v", err))
		return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSMachine %v: %w", klog.KObj(scope.IBMPowerVSMachine), err)
	}

	// 5. Delete Ignition Bootstrap Data from COS (if applicable)
	if err := scope.DeleteMachineIgnition(ctx); err != nil {
		log.Error(err, "error deleting IBMPowerVSMachine ignition data")
		r.markCondition(scope, metav1.ConditionFalse, infrav1.InstanceDeletingReason, fmt.Sprintf("failed to delete ignition data: %v", err))
		return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSMachine ignition %v: %w", klog.KObj(scope.IBMPowerVSMachine), err)
	}

	// 6. Cleanup local caches
	if err := scope.DHCPIPCacheStore.Delete(powervs.VMip{Name: scope.IBMPowerVSMachine.Name}); err != nil {
		// This is non-fatal. We just log it and move on so we don't block the finalizer removal.
		log.Error(err, "failed to delete the machine entry from DHCP cache store")
	}

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSMachineReconciler) reconcileNormal(ctx context.Context, machineScope *powervsscope.MachineScope) (ctrl.Result, error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)

	// 1. Gate: Wait for Infrastructure
	if machineScope.Cluster.Status.Initialization.InfrastructureProvisioned == nil || !*machineScope.Cluster.Status.Initialization.InfrastructureProvisioned {
		log.Info("Cluster infrastructure is not ready yet, skipping reconciliation")
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceWaitingForClusterInfrastructureReadyReason, "")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	// 2. Gate: Wait for Image Import
	if machineScope.IBMPowerVSImage != nil && machineScope.IBMPowerVSImage.Status.ImageState != infrav1.PowerVSImageStateACTIVE {
		log.Info("IBMPowerVSImage is not active yet, skipping reconciliation", "imageState", machineScope.IBMPowerVSImage.Status.ImageState)
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceWaitingForImageReason, "")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	// 3. Gate: Wait for Bootstrap Data
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		if !util.IsControlPlaneMachine(machineScope.Machine) && !conditions.IsTrue(machineScope.Cluster, clusterv1.ClusterControlPlaneInitializedCondition) {
			log.Info("Waiting for the control plane to be initialized, skipping reconciliation")
			r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceWaitingForControlPlaneInitializedReason, "")
			return ctrl.Result{}, nil
		}
		log.Info("Waiting for bootstrap data to be ready, skipping reconciliation")
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceWaitingForBootstrapDataReason, "")
		return ctrl.Result{}, nil
	}

	// 4. Create or Get the Machine
	machine, err := machineScope.CreateMachine(ctx)
	if err != nil {
		log.Error(err, "Unable to create PowerVS machine")
		reason := infrav1.InstanceProvisionFailedReason
		var configErr *powervsscope.ConfigurationError

		if errors.As(err, &configErr) {
			reason = infrav1.InvalidMachineConfigurationReason
		}

		r.markCondition(machineScope, metav1.ConditionFalse, reason, err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to create IBMPowerVSMachine: %w", err)
	}

	if machine == nil {
		machineScope.SetNotReady()
		r.markCondition(machineScope, metav1.ConditionUnknown, infrav1.InstanceStateUnknownReason, "Machine creation returned nil without error")
		return ctrl.Result{}, nil
	}

	// 5. Sync Cloud State to Kubernetes Status
	instance, err := machineScope.IBMPowerVSClient.GetInstance(ctx, *machine.PvmInstanceID)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := machineScope.SetProviderID(*machine.PvmInstanceID); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to set provider ID: %w", err)
	}

	machineScope.SetInstanceID(instance.PvmInstanceID)

	machineScope.SetAddresses(ctx, instance)

	machineScope.SetHealth(instance.Health)
	machineScope.SetInstanceState(instance.Status)

	// 6. Evaluate PowerVS Instance Status
	switch machineScope.GetInstanceState() {
	case infrav1.PowerVSInstanceStateBUILD:
		machineScope.SetNotReady()
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceNotReadyReason, "Instance is building")

	case infrav1.PowerVSInstanceStateSHUTOFF:
		machineScope.SetNotReady()
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceStoppedReason, "Instance is shutoff")
		return ctrl.Result{}, nil

	case infrav1.PowerVSInstanceStateACTIVE:
		machineScope.SetReady()

	case infrav1.PowerVSInstanceStateERROR:
		msg := "Unknown error"
		if instance.Fault != nil {
			msg = instance.Fault.Details
		}
		machineScope.SetNotReady()
		// Note: No more SetFailureMessage/Reason! It goes straight into the condition.
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceErroredReason, msg)
		r.Recorder.Eventf(machineScope.IBMPowerVSMachine, corev1.EventTypeWarning, "FailedBuildInstance", "Failed to build the instance: %s", msg)
		return ctrl.Result{}, nil

	default:
		machineScope.SetNotReady()
		log.Info("PowerVS instance state is undefined", "state", *instance.Status, "instance-id", machineScope.GetInstanceID())
		r.markCondition(machineScope, metav1.ConditionUnknown, infrav1.InstanceStateUnknownReason, fmt.Sprintf("Unknown state: %s", *instance.Status))
	}

	// Requeue if instance is still booting up
	if !machineScope.IsReady() {
		log.Info("IBMPowerVSMachine instance is not ready, requeue", "state", *instance.Status)
		return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
	}

	// 7. Load Balancer Registration (If applicable)
	if machineScope.IBMPowerVSCluster.Spec.VPC.Region == "" {
		log.Info("Skipping configuring machine to load balancer as VPC is not set")
		r.markCondition(machineScope, metav1.ConditionTrue, infrav1.InstanceReadyReason, "")
		return ctrl.Result{}, nil
	}

	internalIP := machineScope.GetMachineInternalIP()
	if internalIP == "" {
		log.Info("Unable to update the load balancer, Machine internal IP not yet set")
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceWaitingForNetworkAddressReason, "Internal IP not yet set")
		return ctrl.Result{}, nil
	}

	log.Info("Configuring load balancer for machine", "IP", internalIP)
	result, err := r.handleLoadBalancerPoolMemberConfiguration(ctx, machineScope)
	if err != nil {
		r.markCondition(machineScope, metav1.ConditionFalse, infrav1.InstanceLoadBalancerConfigurationFailedReason, fmt.Sprintf("Failed to configure load balancer: %v", err))
		return result, fmt.Errorf("failed to configure load balancer: %w", err)
	}

	// 8. Mark conditions
	r.markCondition(machineScope, metav1.ConditionTrue, infrav1.InstanceReadyReason, "")
	return result, nil
}

// markCondition safely sets both the modern and legacy conditions for the machine.
func (r *IBMPowerVSMachineReconciler) markCondition(machineScope *powervsscope.MachineScope, status metav1.ConditionStatus, reason, msg string) {
	// v1beta3 Condition
	conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
		Type:    infrav1.InstanceReadyCondition,
		Status:  status,
		Reason:  reason,
		Message: msg,
	})

	// Legacy v1beta2 mapping
	legacySeverity := clusterv1.ConditionSeverityInfo
	if status == metav1.ConditionFalse {
		switch reason {
		case infrav1.InstanceErroredReason, infrav1.InstanceProvisionFailedReason, infrav1.InstanceStoppedReason:
			legacySeverity = clusterv1.ConditionSeverityError
		case infrav1.InstanceWaitingForClusterInfrastructureReadyReason,
			infrav1.InstanceWaitingForControlPlaneInitializedReason,
			infrav1.InstanceWaitingForBootstrapDataReason,
			infrav1.InstanceWaitingForImageReason:
			legacySeverity = clusterv1.ConditionSeverityInfo
		default:
			legacySeverity = clusterv1.ConditionSeverityWarning
		}
	}

	switch status {
	case metav1.ConditionUnknown:
		deprecatedv1beta1conditions.MarkUnknown(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyV1Beta2Condition, reason, "%s", msg)
	case metav1.ConditionTrue:
		deprecatedv1beta1conditions.MarkTrue(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyV1Beta2Condition)
	default:
		deprecatedv1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyV1Beta2Condition, reason, legacySeverity, "%s", msg)
	}
}

// handleLoadBalancerPoolMemberConfiguration handles load balancer pool member creation flow.
func (r *IBMPowerVSMachineReconciler) handleLoadBalancerPoolMemberConfiguration(ctx context.Context, machineScope *powervsscope.MachineScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	pendingUpdate, err := machineScope.CreateVPCLoadBalancerPoolMember(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to configure VPC load balancer pool member: %w", err)
	}
	if pendingUpdate {
		log.V(3).Info("VPC load balancer pool member registration incomplete, requeuing",
			"requeueAfter", loadBalancerSettleRequeueInterval)
		return ctrl.Result{RequeueAfter: loadBalancerSettleRequeueInterval}, nil
	}
	return ctrl.Result{}, nil
}

func patchIBMPowerVSMachine(ctx context.Context, patchHelper *patch.Helper, ibmPowerVSMachine *infrav1.IBMPowerVSMachine) error {
	// Before computing ready condition, make sure that InstanceReady is always set.
	// NOTE: This is required because v1beta2 conditions comply to guideline requiring conditions to be set at the
	// first reconcile.
	if c := conditions.Get(ibmPowerVSMachine, infrav1.InstanceReadyCondition); c == nil {
		if ptr.Deref(ibmPowerVSMachine.Status.Initialization.Provisioned, false) {
			conditions.Set(ibmPowerVSMachine, metav1.Condition{
				Type:   infrav1.InstanceReadyCondition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.InstanceReadyReason,
			})
		} else {
			conditions.Set(ibmPowerVSMachine, metav1.Condition{
				Type:   infrav1.InstanceReadyCondition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.InstanceNotReadyReason,
			})
		}
	}

	// always update the readyCondition.
	deprecatedv1beta1conditions.SetSummary(ibmPowerVSMachine,
		deprecatedv1beta1conditions.WithConditions(
			infrav1.InstanceReadyV1Beta2Condition,
		),
	)

	if err := conditions.SetSummaryCondition(ibmPowerVSMachine, ibmPowerVSMachine, infrav1.IBMPowerVSMachineReadyCondition,
		conditions.ForConditionTypes{
			infrav1.InstanceReadyCondition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		conditions.CustomMergeStrategy{
			MergeStrategy: conditions.DefaultMergeStrategy(
				// Use custom reasons.
				conditions.ComputeReasonFunc(conditions.GetDefaultComputeMergeReasonFunc(
					infrav1.IBMPowerVSMachineNotReadyReason,
					infrav1.IBMPowerVSMachineReadyUnknownReason,
					infrav1.IBMPowerVSMachineReadyReason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1.IBMPowerVSMachineReadyCondition, err)
	}

	// Patch the IBMPowerVSMachine resource.
	return patchHelper.Patch(ctx, ibmPowerVSMachine, patch.WithOwnedConditions{Conditions: []string{
		infrav1.IBMPowerVSMachineReadyCondition,
		infrav1.InstanceReadyCondition,
		clusterv1.PausedCondition,
	}}, patch.Clusterv1ConditionsFieldPath{statusField, deprecatedStatus, v1beta2Version, deprecatedConditionsField})
}

// ibmPowerVSClusterToIBMPowerVSMachines is a handler.ToRequestsFunc to be used to enqueue requests for reconciliation
// of IBMPowerVSMachines.
func (r *IBMPowerVSMachineReconciler) ibmPowerVSClusterToIBMPowerVSMachines(ctx context.Context, o client.Object) []ctrl.Request {
	log := ctrl.LoggerFrom(ctx)
	result := []ctrl.Request{}
	c, ok := o.(*infrav1.IBMPowerVSCluster)
	if !ok {
		log.Error(fmt.Errorf("expected a IBMPowerVSCluster but got a %T", o), "failed to get IBMPowerVSMachines for IBMPowerVSCluster")
		return nil
	}

	cluster, err := util.GetOwnerCluster(ctx, r.Client, c.ObjectMeta)
	switch {
	case apierrors.IsNotFound(err) || cluster == nil:
		return result
	case err != nil:
		log.Error(err, "failed to get owning cluster")
		return result
	}

	labels := map[string]string{clusterv1.ClusterNameLabel: cluster.Name}
	machineList := &clusterv1.MachineList{}
	if err := r.List(ctx, machineList, client.InNamespace(c.Namespace), client.MatchingLabels(labels)); err != nil {
		log.Error(err, "failed to list Machines")
		return nil
	}
	for _, m := range machineList.Items {
		if m.Spec.InfrastructureRef.Name == "" {
			continue
		}
		name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}

	return result
}

// SetupWithManager creates a new IBMVPCMachine controller for a manager.
func (r *IBMPowerVSMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	if r.ClientBuilder == nil {
		r.ClientBuilder = powervsscope.ProdClientBuilder{}
	}
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "ibmpowervsmachine")
	clusterToIBMPowerVSMachines, err := util.ClusterToTypedObjectsMapper(mgr.GetClient(), &infrav1.IBMPowerVSMachineList{}, mgr.GetScheme())
	if err != nil {
		return err
	}

	err = ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.IBMPowerVSMachine{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(r.Scheme, predicateLog, r.WatchFilterValue)).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("IBMPowerVSMachine"))),
			builder.WithPredicates(predicates.ResourceIsChanged(r.Scheme, predicateLog)),
		).
		Watches(
			&infrav1.IBMPowerVSCluster{},
			handler.EnqueueRequestsFromMapFunc(r.ibmPowerVSClusterToIBMPowerVSMachines),
			builder.WithPredicates(predicates.ResourceIsChanged(r.Scheme, predicateLog)),
		).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToIBMPowerVSMachines),
			builder.WithPredicates(predicates.All(r.Scheme, predicateLog,
				predicates.ResourceIsChanged(r.Scheme, predicateLog),
				predicates.ClusterPausedTransitionsOrInfrastructureProvisioned(r.Scheme, predicateLog),
			)),
		).
		Complete(r)
	if err != nil {
		return fmt.Errorf("could not set up controller for IBMPowerVSMachine: %w", err)
	}

	return nil
}
