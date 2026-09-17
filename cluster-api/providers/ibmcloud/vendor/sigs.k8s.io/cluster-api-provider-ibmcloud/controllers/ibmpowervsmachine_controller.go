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

package controllers

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"         //nolint:staticcheck
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2" //nolint:staticcheck
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"                   //nolint:staticcheck                   //nolint:staticcheck
	"sigs.k8s.io/cluster-api/util/deprecated/v1beta1/paused"
	"sigs.k8s.io/cluster-api/util/finalizers"
	clog "sigs.k8s.io/cluster-api/util/log"
	"sigs.k8s.io/cluster-api/util/predicates"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	capibmrecord "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

// IBMPowerVSMachineReconciler reconciles a IBMPowerVSMachine object.
type IBMPowerVSMachineReconciler struct {
	client.Client
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme

	// WatchFilterValue is the label value used to filter events prior to reconciliation.
	WatchFilterValue string
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

	// Fetch the IBMPowerVSMachine instance.
	ibmPowerVSMachine := &infrav1.IBMPowerVSMachine{}
	err := r.Client.Get(ctx, req.NamespacedName, ibmPowerVSMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("IBMPowerVSMachine not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSMachine: %w", err)
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmPowerVSMachine, infrav1.IBMPowerVSMachineFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	// Fetch the Machine.
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
	// Specifically, it will add KubeadmControlPlane, MachineSet and MachineDeployment.
	if ctx, log, err = clog.AddOwners(ctx, r.Client, ibmPowerVSMachine); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to add owners to log: %w", err)
	}

	// Fetch the Cluster.
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

	// Fetch the IBMPowerVSCluster.
	ibmPowerVSCluster := &infrav1.IBMPowerVSCluster{}
	ibmPowerVSClusterName := client.ObjectKey{
		Namespace: ibmPowerVSMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, ibmPowerVSClusterName, ibmPowerVSCluster); err != nil {
		log.Info("IBMPowerVSCluster is not available yet")
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSCluster: %w", err)
	}

	log = log.WithValues("IBMPowerVSCluster", klog.KObj(ibmPowerVSCluster))
	ctx = ctrl.LoggerInto(ctx, log)

	// Fetch the IBMPowerVSImage.
	var ibmPowerVSImage *infrav1.IBMPowerVSImage
	if ibmPowerVSMachine.Spec.ImageRef != nil {
		ibmPowerVSImage = &infrav1.IBMPowerVSImage{}
		ibmPowerVSImageName := client.ObjectKey{
			Namespace: ibmPowerVSMachine.Namespace,
			Name:      ibmPowerVSMachine.Spec.ImageRef.Name,
		}
		if err := r.Client.Get(ctx, ibmPowerVSImageName, ibmPowerVSImage); err != nil {
			log.Info("IBMPowerVSImage is not available yet", "IBMPowerVSImage", klog.KObj(ibmPowerVSImage))
			return ctrl.Result{}, nil
		}
	}

	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, ibmPowerVSMachine); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	if !cluster.Spec.InfrastructureRef.IsDefined() {
		log.Info("Cluster infrastructureRef is not available yet")
		return ctrl.Result{}, nil
	}

	// Create the machine scope.
	machineScope, err := scope.NewPowerVSMachineScope(scope.PowerVSMachineScopeParams{
		Client:            r.Client,
		Logger:            log,
		Cluster:           cluster,
		IBMPowerVSCluster: ibmPowerVSCluster,
		Machine:           machine,
		IBMPowerVSMachine: ibmPowerVSMachine,
		IBMPowerVSImage:   ibmPowerVSImage,
		ServiceEndpoint:   r.ServiceEndpoint,
		DHCPIPCacheStore:  dhcpCacheStore,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create IBMPowerVS machine scope: %w", err)
	}

	// Initialize the patch helper
	patchHelper, err := v1beta1patch.NewHelper(ibmPowerVSMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Always attempt to Patch the IBMPowerVSMachine object and status after each reconciliation.
	defer func() {
		if err := patchIBMPowerVSMachine(ctx, patchHelper, ibmPowerVSMachine); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Handle deleted machines.
	if !ibmPowerVSMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machineScope)
	}

	// Handle non-deleted machines.
	return r.reconcileNormal(ctx, machineScope)
}

func (r *IBMPowerVSMachineReconciler) reconcileDelete(ctx context.Context, scope *scope.PowerVSMachineScope) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	v1beta1conditions.MarkFalse(scope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, clusterv1.DeletingReason, clusterv1beta1.ConditionSeverityInfo, "")
	v1beta2conditions.Set(scope.IBMPowerVSMachine, metav1.Condition{
		Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
		Status: metav1.ConditionFalse,
		Reason: infrav1.IBMPowerVSMachineInstanceDeletingV1Beta2Reason,
	})

	defer func() {
		if reterr == nil {
			// PowerVS machine is deleted so remove the finalizer.
			controllerutil.RemoveFinalizer(scope.IBMPowerVSMachine, infrav1.IBMPowerVSMachineFinalizer)
		}
	}()

	if scope.IBMPowerVSMachine.Status.InstanceID == "" {
		log.Info("IBMPowerVSMachine instance id is not yet set, so not invoking the PowerVS API to delete the instance")
		return ctrl.Result{}, nil
	}
	if err := scope.DeleteMachine(); err != nil {
		log.Error(err, "error deleting IBMPowerVSMachine")
		v1beta1conditions.MarkFalse(scope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, clusterv1beta1.DeletionFailedReason, clusterv1beta1.ConditionSeverityWarning, "")
		v1beta2conditions.Set(scope.IBMPowerVSMachine, metav1.Condition{
			Type:    infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.IBMPowerVSMachineInstanceDeletingV1Beta2Reason,
			Message: fmt.Sprintf("failed to delete IBMPowerVSMachine: %v", err),
		})
		return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSMachine %v: %w", klog.KObj(scope.IBMPowerVSMachine), err)
	}
	if err := scope.DeleteMachineIgnition(ctx); err != nil {
		log.Info("error deleting IBMPowerVSMachine ignition")
		return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSMachine ignition %v: %w", klog.KObj(scope.IBMPowerVSMachine), err)
	}
	// Remove the cached VM IP
	err := scope.DHCPIPCacheStore.Delete(powervs.VMip{Name: scope.IBMPowerVSMachine.Name})
	if err != nil {
		log.Error(err, "failed to delete the machine entry from DHCP cache store")
	}
	return ctrl.Result{}, nil
}

// handleLoadBalancerPoolMemberConfiguration handles load balancer pool member creation flow.
func (r *IBMPowerVSMachineReconciler) handleLoadBalancerPoolMemberConfiguration(ctx context.Context, machineScope *scope.PowerVSMachineScope) (ctrl.Result, error) {
	poolMember, err := machineScope.CreateVPCLoadBalancerPoolMember(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create VPC load balancer pool member: %w", err)
	}
	if poolMember != nil && *poolMember.ProvisioningStatus != string(infrav1.VPCLoadBalancerStateActive) {
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}
	return ctrl.Result{}, nil
}

func (r *IBMPowerVSMachineReconciler) reconcileNormal(ctx context.Context, machineScope *scope.PowerVSMachineScope) (ctrl.Result, error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)

	if machineScope.Cluster.Status.Initialization.InfrastructureProvisioned == nil || !*machineScope.Cluster.Status.Initialization.InfrastructureProvisioned {
		log.Info("Cluster infrastructure is not ready yet, skipping reconciliation")
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.IBMPowerVSMachineInstanceWaitingForClusterInfrastructureReadyV1Beta2Reason,
		})
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if machineScope.IBMPowerVSImage != nil {
		if !machineScope.IBMPowerVSImage.Status.Ready {
			log.Info("IBMPowerVSImage is not ready yet, skipping reconciliation")
			v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForIBMPowerVSImageReason, clusterv1beta1.ConditionSeverityInfo, "")
			v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
				Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.WaitingForIBMPowerVSImageReason,
			})
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		}
	}

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		if !util.IsControlPlaneMachine(machineScope.Machine) && !conditions.IsTrue(machineScope.Cluster, clusterv1.ClusterControlPlaneInitializedCondition) {
			log.Info("Waiting for the control plane to be initialized, skipping reconciliation")
			v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, clusterv1beta1.WaitingForControlPlaneAvailableReason, clusterv1beta1.ConditionSeverityInfo, "")
			v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
				Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMPowerVSMachineInstanceWaitingForControlPlaneInitializedV1Beta2Reason,
			})
			return ctrl.Result{}, nil
		}

		log.Info("Waiting for bootstrap data to be ready, skipping reconciliation")
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForBootstrapDataReason, clusterv1beta1.ConditionSeverityInfo, "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.IBMPowerVSMachineInstanceWaitingForBootstrapDataV1Beta2Reason,
		})
		return reconcile.Result{}, nil
	}

	machine, err := machineScope.CreateMachine(ctx)
	if err != nil {
		log.Error(err, "Unable to create PowerVS machine")
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceProvisionFailedReason, clusterv1beta1.ConditionSeverityError, "%s", err.Error())
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:    infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.InstanceProvisionFailedReason,
			Message: err.Error(),
		})
		return ctrl.Result{}, fmt.Errorf("failed to create IBMPowerVSMachine: %w", err)
	}

	if machine == nil {
		machineScope.SetNotReady()
		v1beta1conditions.MarkUnknown(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStateUnknownReason, "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionUnknown,
			Reason: infrav1.InstanceStateUnknownReason,
		})
		return ctrl.Result{}, nil
	}

	instance, err := machineScope.IBMPowerVSClient.GetInstance(*machine.PvmInstanceID)
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

	switch machineScope.GetInstanceState() {
	case infrav1.PowerVSInstanceStateBUILD:
		machineScope.SetNotReady()
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, clusterv1beta1.ConditionSeverityWarning, "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.InstanceNotReadyReason,
		})
	case infrav1.PowerVSInstanceStateSHUTOFF:
		machineScope.SetNotReady()
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStoppedReason, clusterv1beta1.ConditionSeverityError, "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.InstanceStoppedReason,
		})
		return ctrl.Result{}, nil
	case infrav1.PowerVSInstanceStateACTIVE:
		machineScope.SetReady()
	case infrav1.PowerVSInstanceStateERROR:
		msg := ""
		if instance.Fault != nil {
			msg = instance.Fault.Details
		}
		machineScope.SetNotReady()
		machineScope.SetFailureReason(infrav1.UpdateMachineError)
		machineScope.SetFailureMessage(msg)
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceErroredReason, clusterv1beta1.ConditionSeverityError, "%s", msg)
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:    infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.InstanceErroredReason,
			Message: msg,
		})
		capibmrecord.Warnf(machineScope.IBMPowerVSMachine, "FailedBuildInstance", "Failed to build the instance %s", msg)
		return ctrl.Result{}, nil
	default:
		machineScope.SetNotReady()
		log.Info("PowerVS instance state is undefined", "state", *instance.Status, "instance-id", machineScope.GetInstanceID())
		v1beta1conditions.MarkUnknown(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, "", "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionUnknown,
			Reason: infrav1.InstanceStateUnknownReason,
		})
	}

	// Requeue after 2 minute if machine is not ready to update status of the machine properly.
	if !machineScope.IsReady() {
		log.Info("IBMPowerVSMachine instance is not ready, requeue", "state", *instance.Status)
		return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
	}

	if machineScope.IBMPowerVSCluster.Spec.VPC == nil || machineScope.IBMPowerVSCluster.Spec.VPC.Region == nil {
		log.Info("Skipping configuring machine to load balancer as VPC is not set")
		v1beta1conditions.MarkTrue(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition)
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionTrue,
			Reason: infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Reason,
		})
		return ctrl.Result{}, nil
	}

	// Register instance with load balancer
	log.Info("Updating load balancer for machine")
	internalIP := machineScope.GetMachineInternalIP()
	if internalIP == "" {
		log.Info("Unable to update the load balancer, Machine internal IP not yet set")
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.IBMPowerVSMachineInstanceWaitingForNetworkAddressV1Beta2Reason, clusterv1beta1.ConditionSeverityWarning, "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:    infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.IBMPowerVSMachineInstanceWaitingForNetworkAddressV1Beta2Reason,
			Message: "Internal IP not yet set",
		})
		return ctrl.Result{}, nil
	}
	log.Info("Configuring load balancer for machine", "IP", internalIP)
	result, err := r.handleLoadBalancerPoolMemberConfiguration(ctx, machineScope)
	if err != nil {
		v1beta1conditions.MarkFalse(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition, infrav1.IBMPowerVSMachineInstanceLoadBalancerConfigurationFailedV1Beta2Reason, clusterv1beta1.ConditionSeverityWarning, "")
		v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
			Type:    infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.IBMPowerVSMachineInstanceLoadBalancerConfigurationFailedV1Beta2Reason,
			Message: fmt.Sprintf("Failed to configure load balancer: %v", err),
		})
		return result, fmt.Errorf("failed to configure load balancer: %w", err)
	}
	v1beta1conditions.MarkTrue(machineScope.IBMPowerVSMachine, infrav1.InstanceReadyCondition)
	v1beta2conditions.Set(machineScope.IBMPowerVSMachine, metav1.Condition{
		Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Reason,
	})
	return result, nil
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

func patchIBMPowerVSMachine(ctx context.Context, patchHelper *v1beta1patch.Helper, ibmPowerVSMachine *infrav1.IBMPowerVSMachine) error {
	// Before computing ready condition, make sure that InstanceReady is always set.
	// NOTE: This is required because v1beta2 conditions comply to guideline requiring conditions to be set at the
	// first reconcile.
	if c := v1beta2conditions.Get(ibmPowerVSMachine, infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition); c == nil {
		if ibmPowerVSMachine.Status.Ready {
			v1beta2conditions.Set(ibmPowerVSMachine, metav1.Condition{
				Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Reason,
			})
		} else {
			v1beta2conditions.Set(ibmPowerVSMachine, metav1.Condition{
				Type:   infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMPowerVSMachineInstanceNotReadyV1Beta2Reason,
			})
		}
	}

	// always update the readyCondition.
	v1beta1conditions.SetSummary(ibmPowerVSMachine,
		v1beta1conditions.WithConditions(
			infrav1.InstanceReadyCondition,
		),
	)

	if err := v1beta2conditions.SetSummaryCondition(ibmPowerVSMachine, ibmPowerVSMachine, infrav1.IBMPowerVSMachineReadyV1Beta2Condition,
		v1beta2conditions.ForConditionTypes{
			infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		v1beta2conditions.CustomMergeStrategy{
			MergeStrategy: v1beta2conditions.DefaultMergeStrategy(
				// Use custom reasons.
				v1beta2conditions.ComputeReasonFunc(v1beta2conditions.GetDefaultComputeMergeReasonFunc(
					infrav1.IBMPowerVSMachineNotReadyV1Beta2Reason,
					infrav1.IBMPowerVSMachineReadyUnknownV1Beta2Reason,
					infrav1.IBMPowerVSMachineReadyV1Beta2Reason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1.IBMPowerVSMachineReadyV1Beta2Condition, err)
	}

	// Patch the IBMPowerVSMachine resource.
	return patchHelper.Patch(ctx, ibmPowerVSMachine, v1beta1patch.WithOwnedV1Beta2Conditions{Conditions: []string{
		infrav1.IBMPowerVSMachineReadyV1Beta2Condition,
		infrav1.IBMPowerVSMachineInstanceReadyV1Beta2Condition,
		clusterv1beta1.PausedV1Beta2Condition,
	}})
}
