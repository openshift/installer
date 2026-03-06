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
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"

	"github.com/IBM/vpc-go-sdk/vpcv1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"         //nolint:staticcheck
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2" //nolint:staticcheck
	v1beta1patch "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"                   //nolint:staticcheck
	"sigs.k8s.io/cluster-api/util/deprecated/v1beta1/paused"
	"sigs.k8s.io/cluster-api/util/finalizers"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	capibmrecord "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

// IBMVPCMachineReconciler reconciles a IBMVPCMachine object.
type IBMVPCMachineReconciler struct {
	client.Client
	Log             logr.Logger
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconcileation logic for IBMVPCMachine.
func (r *IBMVPCMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconciling IBMVPCMachine")
	defer log.Info("Finished reconciling IBMVPCMachine")

	// Fetch the IBMVPCMachine instance.
	ibmVPCMachine := &infrav1.IBMVPCMachine{}
	err := r.Get(ctx, req.NamespacedName, ibmVPCMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, ibmVPCMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, ibmVPCMachine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	ibmVPCCluster := &infrav1.IBMVPCCluster{}
	ibmVPCClusterName := client.ObjectKey{
		Namespace: ibmVPCMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, ibmVPCClusterName, ibmVPCCluster); err != nil {
		log.Info("IBMVPCCluster is not available yet")
		return ctrl.Result{}, nil
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmVPCMachine, infrav1.MachineFinalizer); err != nil || finalizerAdded {
		return ctrl.Result{}, err
	}

	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	// Initialize the patch helper.
	patchHelper, err := v1beta1patch.NewHelper(ibmVPCMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to initialize patch helper: %w", err)
	}

	// Always attempt to Patch the IBMVPCMachine object and status after each reconciliation.
	defer func() {
		if err := patchIBMVPCMachine(ctx, patchHelper, ibmVPCMachine); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, ibmVPCMachine); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	// Create the machine scope.
	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Client:          r.Client,
		Cluster:         cluster,
		IBMVPCCluster:   ibmVPCCluster,
		Machine:         machine,
		IBMVPCMachine:   ibmVPCMachine,
		ServiceEndpoint: r.ServiceEndpoint,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	log = log.WithValues("IBMVPCMachine", klog.KObj(ibmVPCMachine))
	ctx = ctrl.LoggerInto(ctx, log)

	// Handle deleted machines.
	if !ibmVPCMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machineScope)
	}

	// Handle non-deleted machines.
	return r.reconcileNormal(ctx, machineScope)
}

// SetupWithManager creates a new IBMVPCMachine controller for a manager.
func (r *IBMVPCMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.IBMVPCMachine{}).
		Complete(r)
}

func (r *IBMVPCMachineReconciler) reconcileNormal(ctx context.Context, machineScope *scope.MachineScope) (ctrl.Result, error) { //nolint:gocyclo
	log := ctrl.LoggerFrom(ctx)
	if controllerutil.AddFinalizer(machineScope.IBMVPCMachine, infrav1.MachineFinalizer) {
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Bootstrap data secret reference is not yet available")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	if machineScope.IBMVPCCluster.Status.Subnet.ID != nil {
		machineScope.IBMVPCMachine.Spec.PrimaryNetworkInterface = infrav1.NetworkInterface{
			Subnet: *machineScope.IBMVPCCluster.Status.Subnet.ID,
		}
	}

	instance, err := r.getOrCreate(ctx, machineScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile VSI for IBMVPCMachine %s/%s: %w", machineScope.IBMVPCMachine.Namespace, machineScope.IBMVPCMachine.Name, err)
	}

	machineRunning := false
	if instance != nil {
		// Attempt to tag the Instance.
		if err := machineScope.TagResource(machineScope.IBMVPCCluster.Name, *instance.CRN); err != nil {
			return ctrl.Result{}, fmt.Errorf("error failed to tag machine: %w", err)
		}

		// Set available status' for Machine.
		machineScope.SetInstanceID(*instance.ID)
		if err := machineScope.SetProviderID(instance.ID); err != nil {
			return ctrl.Result{}, fmt.Errorf("error failed to set machine provider id: %w", err)
		}
		machineScope.SetAddresses(instance)
		machineScope.SetInstanceStatus(*instance.Status)

		// Depending on the state of the Machine, update status, conditions, etc.
		switch machineScope.GetInstanceStatus() {
		case vpcv1.InstanceStatusPendingConst:
			machineScope.SetNotReady()
			v1beta1conditions.MarkFalse(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, clusterv1beta1.ConditionSeverityWarning, "")
			v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMVPCMachineInstanceNotReadyV1Beta2Reason,
			})
		case vpcv1.InstanceStatusStartingConst:
			machineScope.SetNotReady()
			v1beta1conditions.MarkFalse(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, clusterv1beta1.ConditionSeverityWarning, "")
			v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMVPCMachineInstanceNotReadyV1Beta2Reason,
			})
		case vpcv1.InstanceStatusStoppedConst:
			machineScope.SetNotReady()
			v1beta1conditions.MarkFalse(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStoppedReason, clusterv1beta1.ConditionSeverityError, "")
			v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.InstanceStoppedReason,
			})
		case vpcv1.InstanceStatusDeletingConst:
			v1beta1conditions.MarkFalse(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition, infrav1.InstanceDeletingReason, clusterv1beta1.ConditionSeverityError, "")
			v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.InstanceDeletingReason,
			})
		case vpcv1.InstanceStatusFailedConst:
			msg := ""
			healthReasonsLen := len(instance.HealthReasons)
			if healthReasonsLen > 0 {
				// Create a failure message using the last entry's Code and Message fields.
				// TODO(cjschaef): Consider adding the MoreInfo field as well, as it contains a link to IBM Cloud docs.
				msg = fmt.Sprintf("%s: %s", *instance.HealthReasons[healthReasonsLen-1].Code, *instance.HealthReasons[healthReasonsLen-1].Message)
			}
			machineScope.SetNotReady()
			machineScope.SetFailureReason(infrav1.UpdateMachineError)
			machineScope.SetFailureMessage(msg)
			v1beta1conditions.MarkFalse(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition, infrav1.InstanceErroredReason, clusterv1beta1.ConditionSeverityError, "%s", msg)
			v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.InstanceErroredReason,
			})
			capibmrecord.Warnf(machineScope.IBMVPCMachine, "FailedBuildInstance", "Failed to build the instance - %s", msg)
			return ctrl.Result{}, nil
		case vpcv1.InstanceStatusRunningConst:
			machineRunning = true
		default:
			machineScope.SetNotReady()
			log.V(3).Info("unexpected vpc instance status", "instanceStatus", *instance.Status, "instanceID", machineScope.GetInstanceID())
			v1beta1conditions.MarkUnknown(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition, "", "")
			v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.InstanceStateUnknownReason,
			})
		}
	} else {
		machineScope.SetNotReady()
		v1beta1conditions.MarkUnknown(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStateUnknownReason, "")
		v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
			Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.InstanceStateUnknownReason,
		})
	}

	// Check if the Machine is running.
	if !machineRunning {
		// Requeue after 1 minute if machine is not running.
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	// Rely on defined VPC Load Balancer Pool Members first before falling back to hardcoded defaults.
	if len(machineScope.IBMVPCMachine.Spec.LoadBalancerPoolMembers) > 0 {
		needsRequeue := false
		for _, poolMember := range machineScope.IBMVPCMachine.Spec.LoadBalancerPoolMembers {
			requeue, err := machineScope.ReconcileVPCLoadBalancerPoolMember(ctx, poolMember)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("error failed to reconcile machine's pool member: %w", err)
			} else if requeue {
				needsRequeue = true
			}
		}

		// If any VPC Load Balancer Pool Member needs reconciliation, requeue.
		if needsRequeue {
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		}
	} else {
		// Otherwise, default to previous Load Balancer Pool Member configuration.
		_, ok := machineScope.IBMVPCMachine.Labels[clusterv1.MachineControlPlaneNameLabel]
		if err = machineScope.SetProviderID(instance.ID); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to set provider id IBMVPCMachine %s/%s: %w", machineScope.IBMVPCMachine.Namespace, machineScope.IBMVPCMachine.Name, err)
		}
		if ok {
			if instance.PrimaryNetworkInterface.PrimaryIP.Address == nil || *instance.PrimaryNetworkInterface.PrimaryIP.Address == "0.0.0.0" {
				return ctrl.Result{}, fmt.Errorf("invalid primary ip address")
			}
			internalIP := instance.PrimaryNetworkInterface.PrimaryIP.Address
			port := int64(machineScope.APIServerPort())
			poolMember, err := machineScope.CreateVPCLoadBalancerPoolMember(ctx, internalIP, port)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to bind port %d to control plane %s/%s: %w", port, machineScope.IBMVPCMachine.Namespace, machineScope.IBMVPCMachine.Name, err)
			}
			if poolMember != nil && *poolMember.ProvisioningStatus != string(infrav1.VPCLoadBalancerStateActive) {
				return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
			}
		}
	}

	// Handle Additional Volumes
	var result ctrl.Result
	result, err = r.reconcileAdditionalVolumes(ctx, machineScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error reconciling additional volumes: %w", err)
	}

	// With a running machine and all Load Balancer Pool Members reconciled, mark machine as ready.
	machineScope.SetReady()
	v1beta1conditions.MarkTrue(machineScope.IBMVPCMachine, infrav1.InstanceReadyCondition)
	v1beta2conditions.Set(machineScope.IBMVPCMachine, metav1.Condition{
		Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: infrav1.IBMVPCMachineInstanceReadyV1Beta2Reason,
	})
	log.Info("Reconcile complete", "result", result)
	return result, nil
}

func (r *IBMVPCMachineReconciler) getOrCreate(ctx context.Context, scope *scope.MachineScope) (*vpcv1.Instance, error) {
	instance, err := scope.CreateMachine(ctx)
	return instance, err
}

func (r *IBMVPCMachineReconciler) reconcileDelete(ctx context.Context, scope *scope.MachineScope) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Handling deleted IBMVPCMachine")

	if _, ok := scope.IBMVPCMachine.Labels[clusterv1.MachineControlPlaneNameLabel]; ok {
		if err := scope.DeleteVPCLoadBalancerPoolMember(ctx); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to delete loadBalancer pool member: %w", err)
		}
	}

	if err := scope.DeleteMachine(); err != nil {
		log.Info("Error deleting IBMVPCMachine")
		return ctrl.Result{}, fmt.Errorf("error deleting IBMVPCMachine %s/%s: %w", scope.IBMVPCMachine.Namespace, scope.IBMVPCMachine.Spec.Name, err)
	}

	defer func() {
		if reterr == nil {
			// VSI is deleted so remove the finalizer.
			controllerutil.RemoveFinalizer(scope.IBMVPCMachine, infrav1.MachineFinalizer)
		}
	}()

	return ctrl.Result{}, nil
}

func patchIBMVPCMachine(ctx context.Context, patchHelper *v1beta1patch.Helper, ibmVPCMachine *infrav1.IBMVPCMachine) error {
	// Before computing ready condition, make sure that InstanceReady is always set.
	// NOTE: This is required because v1beta2 conditions comply to guideline requiring conditions to be set at the
	// first reconcile.
	if c := v1beta2conditions.Get(ibmVPCMachine, infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition); c == nil {
		if ibmVPCMachine.Status.Ready {
			v1beta2conditions.Set(ibmVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.IBMVPCMachineInstanceReadyV1Beta2Reason,
			})
		} else {
			v1beta2conditions.Set(ibmVPCMachine, metav1.Condition{
				Type:   infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMVPCMachineInstanceNotReadyV1Beta2Reason,
			})
		}
	}

	v1beta1conditions.SetSummary(ibmVPCMachine,
		v1beta1conditions.WithConditions(
			infrav1.InstanceReadyCondition,
		),
	)

	if err := v1beta2conditions.SetSummaryCondition(ibmVPCMachine, ibmVPCMachine, infrav1.IBMVPCMachineReadyV1Beta2Condition,
		v1beta2conditions.ForConditionTypes{
			infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		v1beta2conditions.CustomMergeStrategy{
			MergeStrategy: v1beta2conditions.DefaultMergeStrategy(
				// Use custom reasons.
				v1beta2conditions.ComputeReasonFunc(v1beta2conditions.GetDefaultComputeMergeReasonFunc(
					infrav1.IBMVPCMachineNotReadyV1Beta2Reason,
					infrav1.IBMVPCMachineReadyUnknownV1Beta2Reason,
					infrav1.IBMVPCMachineReadyV1Beta2Reason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1.IBMVPCMachineReadyV1Beta2Condition, err)
	}

	// Patch the IBMVPCMachine resource.
	return patchHelper.Patch(ctx, ibmVPCMachine, v1beta1patch.WithOwnedV1Beta2Conditions{Conditions: []string{
		infrav1.IBMVPCMachineReadyV1Beta2Condition,
		infrav1.IBMVPCMachineInstanceReadyV1Beta2Condition,
		clusterv1beta1.PausedV1Beta2Condition,
	}})
}
func (r *IBMVPCMachineReconciler) reconcileAdditionalVolumes(ctx context.Context, machineScope *scope.MachineScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	// Return immediately if no additional volumes exist
	if len(machineScope.IBMVPCMachine.Spec.AdditionalVolumes) == 0 {
		return ctrl.Result{}, nil
	}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: machineScope.IBMVPCMachine.Namespace,
		Name:      machineScope.IBMVPCMachine.Name,
	}, machineScope.IBMVPCMachine)
	if err != nil {
		log.Error(err, "Could not fetch machine status")
		return ctrl.Result{}, err
	}
	machineVolumes := machineScope.IBMVPCMachine.Spec.AdditionalVolumes
	result := ctrl.Result{}
	if machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs == nil {
		machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs = make([]string, len(machineScope.IBMVPCMachine.Spec.AdditionalVolumes))
	}
	volumeAttachmentList, err := machineScope.GetVolumeAttachments()
	if err != nil {
		return result, err
	}
	volumeAttachmentNames := sets.New[string]()
	for i := range volumeAttachmentList {
		sets.Insert(volumeAttachmentNames, *volumeAttachmentList[i].Name)
	}
	errList := []error{}
	// Read through the list, checking if volume exists and create volume if it does not
	for v := range machineVolumes {
		if volumeAttachmentNames.Has(machineVolumes[v].Name) {
			// volume attachment has been created so volume is already attached
			continue
		}
		if machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs[v] != "" {
			// volume was already created, fetch volume status and attach if possible
			state, err := machineScope.GetVolumeState(machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs[v])
			if err != nil {
				errList = append(errList, err)
			}
			switch state {
			case vpcv1.VolumeStatusPendingConst, vpcv1.VolumeStatusUpdatingConst:
				result = ctrl.Result{RequeueAfter: 10 * time.Second}
			case vpcv1.VolumeStatusFailedConst, vpcv1.VolumeStatusUnusableConst:
				errList = append(errList, fmt.Errorf("volume in unexpected state: %s", state))
			case vpcv1.VolumeStatusAvailableConst:
				log.Info("Volume is in available state, trying to attach it", "VolumeID", machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs[v])
				err = machineScope.AttachVolume(machineVolumes[v].DeleteVolumeOnInstanceDelete, machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs[v], machineVolumes[v].Name)
				if err != nil {
					log.Error(err, "Error while attaching volume", "VolumeID", machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs[v])
					errList = append(errList, err)
				}
				log.Info("Successfully attached volume", "VolumeID", machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs[v])
			}
		} else {
			// volume does not exist, create it and requeue so that it becomes available
			volumeID, err := machineScope.CreateVolume(machineVolumes[v])
			machineScope.IBMVPCMachine.Status.V1Beta2.AdditionalVolumeIDs[v] = volumeID
			if err != nil {
				log.Error(err, "Could not update Machine status. Created Volume needs to be cleaned up manually", "VolumeID", volumeID)
				errList = append(errList, err)
			}
			log.Info("Created new volume", "name", machineVolumes[v].Name, "VolumeID", volumeID)
			result = ctrl.Result{RequeueAfter: 10 * time.Second}
		}
	}
	return result, errors.Join(errList...)
}
