/*
Copyright 2018 The Kubernetes Authors.

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

package machine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	openshiftfeatures "github.com/openshift/api/features"
	machinev1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/machine-api-operator/pkg/metrics"
	"github.com/openshift/machine-api-operator/pkg/util"
	"github.com/openshift/machine-api-operator/pkg/util/conditions"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/component-base/featuregate"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	NodeNameEnvVar = "NODE_NAME"
	requeueAfter   = 30 * time.Second

	// ExcludeNodeDrainingAnnotation annotation explicitly skips node draining if set
	ExcludeNodeDrainingAnnotation = "machine.openshift.io/exclude-node-draining"

	// MachineRegionLabelName as annotation name for a machine region
	MachineRegionLabelName = "machine.openshift.io/region"

	// MachineAZLabelName as annotation name for a machine AZ
	MachineAZLabelName = "machine.openshift.io/zone"

	// MachineInstanceStateAnnotationName as annotation name for a machine instance state
	MachineInstanceStateAnnotationName = "machine.openshift.io/instance-state"

	// MachineInstanceTypeLabelName as annotation name for a machine instance type
	MachineInstanceTypeLabelName = "machine.openshift.io/instance-type"

	// MachineInterruptibleInstanceLabelName as annotaiton name for interruptible instances
	MachineInterruptibleInstanceLabelName = "machine.openshift.io/interruptible-instance"

	// Hardcoded instance state set on machine failure
	unknownInstanceState = "Unknown"

	skipWaitForDeleteTimeoutSeconds = 1
)

// We export the PausedCondition and reasons as they're shared
// across the Machine and MachineSet controllers.
const (
	PausedCondition machinev1.ConditionType = "Paused"

	PausedConditionReason = "AuthoritativeAPINotMachineAPI"

	NotPausedConditionReason = "AuthoritativeAPIMachineAPI"
)

var DefaultActuator Actuator

func AddWithActuator(mgr manager.Manager, actuator Actuator, gate featuregate.MutableFeatureGate) error {
	return AddWithActuatorOpts(mgr, actuator, controller.Options{}, gate)
}

func AddWithActuatorOpts(mgr manager.Manager, actuator Actuator, opts controller.Options, gate featuregate.MutableFeatureGate) error {
	machineControllerOpts := opts
	machineControllerOpts.Reconciler = newReconciler(mgr, actuator, gate)

	if err := addWithOpts(mgr, machineControllerOpts, "machine-controller"); err != nil {
		return err
	}

	if err := addWithOpts(mgr, controller.Options{
		Reconciler:  newDrainController(mgr),
		RateLimiter: newDrainRateLimiter(),
	}, "machine-drain-controller"); err != nil {
		return err
	}
	return nil
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, actuator Actuator, gate featuregate.MutableFeatureGate) reconcile.Reconciler {
	r := &ReconcileMachine{
		Client:        mgr.GetClient(),
		eventRecorder: mgr.GetEventRecorderFor("machine-controller"),
		config:        mgr.GetConfig(),
		scheme:        mgr.GetScheme(),
		actuator:      actuator,
		gate:          gate,
	}
	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func addWithOpts(mgr manager.Manager, opts controller.Options, controllerName string) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, opts)
	if err != nil {
		return err
	}

	// Watch for changes to Machine
	return c.Watch(
		source.Kind(mgr.GetCache(), &machinev1.Machine{},
			&handler.TypedEnqueueRequestForObject[*machinev1.Machine]{},
		))
}

// ReconcileMachine reconciles a Machine object
type ReconcileMachine struct {
	client.Client
	config *rest.Config
	scheme *runtime.Scheme

	eventRecorder record.EventRecorder

	actuator Actuator
	gate     featuregate.MutableFeatureGate

	// nowFunc is used to mock time in testing. It should be nil in production.
	nowFunc func() time.Time
}

// Reconcile reads that state of the cluster for a Machine object and makes changes based on the state read
// and what is in the Machine.Spec
// +kubebuilder:rbac:groups=machine.openshift.io,resources=machines;machines/status,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileMachine) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Machine instance
	m := &machinev1.Machine{}
	if err := r.Client.Get(ctx, request.NamespacedName, m); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Implement controller logic here
	machineName := m.GetName()
	klog.Infof("%v: reconciling Machine", machineName)

	// Get the original state of conditions now so that they can be used to calculate the patch later.
	// This must be a copy otherwise the referenced slice will be modified by later machine conditions changes.
	originalConditions := conditions.DeepCopyConditions(m.Status.Conditions)

	if r.gate.Enabled(featuregate.Feature(openshiftfeatures.FeatureGateMachineAPIMigration)) {
		switch m.Status.AuthoritativeAPI {
		case "":
			// An empty .status.authoritativeAPI normally means the resource has not yet been reconciled.
			// and that the value in .spec.authoritativeAPI has not been propagated to .status.authoritativeAPI yet.
			// This value can be set by two separate controllers, depending which one of them is running at that time,
			// or in case they are both running, which one gets to set it first (the operation is idempotent so there is no harm in racing).
			// - the cluster-capi-operator machine-api-migration's migration controller
			// - this controller

			klog.Infof("%v: machine .status.authoritativeAPI is not yet set, setting it to .spec.authoritativeAPI", m.Name)

			if err := r.patchStatusAuthoritativeAPI(ctx, m, m.Spec.AuthoritativeAPI); err != nil {
				klog.Errorf("%v: error patching status to set .status.authoritativeAPI for machine: %v", m.Name, err)
				return reconcile.Result{}, fmt.Errorf("error patching status to set .status.authoritativeAPI for machine %s: %w", m.Name, err)
			}

			// Return to give a chance to the changes to get propagated.
			return reconcile.Result{}, nil

		case machinev1.MachineAuthorityClusterAPI, machinev1.MachineAuthorityMigrating:
			// In cases when .status.authoritativeAPI is set to machinev1.MachineAuthorityClusterAPI, machinev1.MachineAuthorityMigrating
			// the resource should be paused and not reconciled further.
			desiredCondition := conditions.TrueConditionWithReason(
				PausedCondition, PausedConditionReason,
				"The AuthoritativeAPI status is set to '%s'", string(m.Status.AuthoritativeAPI),
			)

			if _, err := r.ensureUpdatedPausedCondition(ctx, m, desiredCondition,
				fmt.Sprintf("%v: machine .status.authoritativeAPI is set to '%s', ensuring machine is paused", machineName, m.Status.AuthoritativeAPI)); err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to ensure paused condition: %w", err)
			}

			klog.Infof("%v: machine is paused, taking no further action", machineName)

			return reconcile.Result{}, nil

		case machinev1.MachineAuthorityMachineAPI:
			// The authority is MachineAPI and the resource should not be paused.
			desiredCondition := conditions.FalseCondition(
				PausedCondition, NotPausedConditionReason, machinev1.ConditionSeverityInfo, "%s",
				fmt.Sprintf("The AuthoritativeAPI status is set to '%s'", string(m.Status.AuthoritativeAPI)),
			)

			if updated, err := r.ensureUpdatedPausedCondition(ctx, m, desiredCondition,
				fmt.Sprintf("%v: machine .status.authoritativeAPI is set to '%s', unpausing machine", machineName, m.Status.AuthoritativeAPI)); err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to ensure paused condition: %w", err)
			} else if updated {
				klog.Infof("%v: setting machine paused condition to false", machineName)
			}

			// Fallthrough and continue reconcilation.
		default:
			klog.Errorf("%v: invalid .status.authoritativeAPI '%s'", machineName, m.Status.AuthoritativeAPI)
			return reconcile.Result{}, nil // Do not return an error to avoid immediate requeue.
		}
	}

	if errList := validateMachine(m); len(errList) > 0 {
		err := fmt.Errorf("%v: machine validation failed: %v", machineName, errList.ToAggregate().Error())
		klog.Error(err)
		r.eventRecorder.Eventf(m, corev1.EventTypeWarning, "FailedValidate", err.Error())
		return reconcile.Result{}, err
	}

	// If object hasn't been deleted and doesn't have a finalizer, add one
	// Add a finalizer to newly created objects.
	if m.ObjectMeta.DeletionTimestamp.IsZero() {
		finalizerCount := len(m.Finalizers)

		if !util.Contains(m.Finalizers, machinev1.MachineFinalizer) {
			m.Finalizers = append(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer)
		}

		if len(m.Finalizers) > finalizerCount {
			if err := r.Client.Update(ctx, m); err != nil {
				klog.Infof("%v: failed to add finalizers to machine: %v", machineName, err)
				return reconcile.Result{}, err
			}

			// Since adding the finalizer updates the object return to avoid later update issues
			return reconcile.Result{}, nil
		}
	}

	if !m.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := r.updateStatus(ctx, m, machinev1.PhaseDeleting, nil, originalConditions); err != nil {
			return reconcile.Result{}, err
		}

		// no-op if finalizer has been removed.
		if !util.Contains(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer) {
			klog.Infof("%v: reconciling machine causes a no-op as there is no finalizer", machineName)
			return reconcile.Result{}, nil
		}

		klog.Infof("%v: reconciling machine triggers delete", machineName)
		// check if machine was already drained
		drainedCondition := conditions.Get(m, machinev1.MachineDrained)
		if drainedCondition == nil || drainedCondition.Status != corev1.ConditionTrue {
			klog.Infof("%s: waiting for node to be drained before deleting instance", machineName)
			// this will requeue and proceed when drain controller will set the condition
			return reconcile.Result{}, nil
		}

		// pre-term.delete lifecycle hook
		// Return early without error, will requeue if/when the hook owner removes the annotation.
		if len(m.Spec.LifecycleHooks.PreTerminate) > 0 {
			klog.Infof("%v: not deleting machine: lifecycle blocked by pre-terminate hook", machineName)
			return reconcile.Result{}, nil
		}

		if err := r.actuator.Delete(ctx, m); err != nil {
			// isInvalidMachineConfiguration will take care of the case where the
			// configuration is invalid from the beginning. len(m.Status.Addresses) > 0
			// will handle the case when a machine configuration was invalidated
			// after an instance was created. So only a small window is left when
			// we can loose instances, e.g. right after request to create one
			// was sent and before a list of node addresses was set.
			if len(m.Status.Addresses) > 0 || !isInvalidMachineConfigurationError(err) {
				klog.Errorf("%v: failed to delete machine: %v", machineName, err)
				return delayIfRequeueAfterError(err)
			}
		}

		instanceExists, err := r.actuator.Exists(ctx, m)
		if err != nil {
			klog.Errorf("%v: failed to check if machine exists: %v", machineName, err)
			return reconcile.Result{}, err
		}

		if instanceExists {
			klog.V(3).Infof("%v: can't proceed deleting machine while cloud instance is being terminated, requeuing", machineName)
			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		}

		if m.Status.NodeRef != nil {
			klog.Infof("%v: deleting node %q for machine", machineName, m.Status.NodeRef.Name)
			if err := r.deleteNode(ctx, m.Status.NodeRef.Name); err != nil {
				klog.Errorf("%v: error deleting node for machine: %v", machineName, err)
				return reconcile.Result{}, err
			}
		}

		// Remove finalizer on successful deletion.
		m.ObjectMeta.Finalizers = util.Filter(m.ObjectMeta.Finalizers, machinev1.MachineFinalizer)
		if err := r.Client.Update(ctx, m); err != nil {
			klog.Errorf("%v: failed to remove finalizer from machine: %v", machineName, err)
			return reconcile.Result{}, err
		}

		klog.Infof("%v: machine deletion successful", machineName)
		return reconcile.Result{}, nil
	}

	if machineIsFailed(m) {
		klog.Warningf("%v: machine has gone %q phase. It won't reconcile", machineName, machinev1.PhaseFailed)
		return reconcile.Result{}, nil
	}

	instanceExists, err := r.actuator.Exists(ctx, m)
	if err != nil {
		klog.Errorf("%v: failed to check if machine exists: %v", machineName, err)

		conditions.Set(m, conditions.UnknownCondition(
			machinev1.InstanceExistsCondition,
			machinev1.ErrorCheckingProviderReason,
			"Failed to check if machine exists: %v", err,
		))

		if patchErr := r.updateStatus(ctx, m, ptr.Deref(m.Status.Phase, ""), nil, originalConditions); patchErr != nil {
			klog.Errorf("%v: error patching status: %v", machineName, patchErr)
		}

		return reconcile.Result{}, err
	}

	if instanceExists {
		klog.Infof("%v: reconciling machine triggers idempotent update", machineName)
		if err := r.actuator.Update(ctx, m); err != nil {
			klog.Errorf("%v: error updating machine: %v, retrying in %v seconds", machineName, err, requeueAfter)

			if patchErr := r.updateStatus(ctx, m, ptr.Deref(m.Status.Phase, ""), nil, originalConditions); patchErr != nil {
				klog.Errorf("%v: error patching status: %v", machineName, patchErr)
			}

			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		}

		// Mark the instance exists condition true after actuator update else the update may overwrite changes
		conditions.MarkTrue(m, machinev1.InstanceExistsCondition)

		if !machineIsProvisioned(m) {
			klog.Errorf("%v: instance exists but providerID or addresses has not been given to the machine yet, requeuing", machineName)
			if patchErr := r.updateStatus(ctx, m, ptr.Deref(m.Status.Phase, ""), nil, originalConditions); patchErr != nil {
				klog.Errorf("%v: error patching status: %v", machineName, patchErr)
			}

			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		}

		if !machineHasNode(m) {
			// Requeue until we reach running phase
			if err := r.updateStatus(ctx, m, machinev1.PhaseProvisioned, nil, originalConditions); err != nil {
				return reconcile.Result{}, err
			}
			klog.Infof("%v: has no node yet, requeuing", machineName)
			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		}

		return reconcile.Result{}, r.updateStatus(ctx, m, machinev1.PhaseRunning, nil, originalConditions)
	}

	// Instance does not exist but the machine has been given a providerID/address.
	// This can only be reached if an instance was deleted outside the machine API
	if machineIsProvisioned(m) {
		conditions.Set(m, conditions.FalseCondition(
			machinev1.InstanceExistsCondition,
			machinev1.InstanceMissingReason,
			machinev1.ConditionSeverityWarning,
			"Instance not found on provider",
		))

		if err := r.updateStatus(ctx, m, machinev1.PhaseFailed, errors.New("can't find created instance"), originalConditions); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	conditions.Set(m, conditions.FalseCondition(
		machinev1.InstanceExistsCondition,
		machinev1.InstanceNotCreatedReason,
		machinev1.ConditionSeverityWarning,
		"Instance has not been created",
	))

	// Machine resource created and instance does not exist yet.
	if ptr.Deref(m.Status.Phase, "") == "" {
		klog.V(2).Infof("%v: setting phase to Provisioning and requeuing", machineName)
		if err := r.updateStatus(ctx, m, machinev1.PhaseProvisioning, nil, originalConditions); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	klog.Infof("%v: reconciling machine triggers idempotent create", machineName)
	if err := r.actuator.Create(ctx, m); err != nil {
		klog.Warningf("%v: failed to create machine: %v", machineName, err)
		if isInvalidMachineConfigurationError(err) {
			if err := r.updateStatus(ctx, m, machinev1.PhaseFailed, err, originalConditions); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		return delayIfRequeueAfterError(err)
	}

	klog.Infof("%v: created instance, requeuing", machineName)
	return reconcile.Result{RequeueAfter: requeueAfter}, nil
}

func (r *ReconcileMachine) deleteNode(ctx context.Context, name string) error {
	var node corev1.Node
	if err := r.Client.Get(ctx, client.ObjectKey{Name: name}, &node); err != nil {
		if apierrors.IsNotFound(err) {
			klog.V(2).Infof("Node %q not found", name)
			return nil
		}
		klog.Errorf("Failed to get node %q: %v", name, err)
		return err
	}
	return r.Client.Delete(ctx, &node)
}

// ensureUpdatedPausedCondition updates the paused condition if needed.
func (r *ReconcileMachine) ensureUpdatedPausedCondition(ctx context.Context, m *machinev1.Machine, desiredCondition *machinev1.Condition, logMessage string) (bool, error) {
	oldM := m.DeepCopy()
	if !conditions.IsEquivalentTo(conditions.Get(m, PausedCondition), desiredCondition) {
		klog.Info(logMessage)
		conditions.Set(m, desiredCondition)
		if err := r.updateStatus(ctx, m, ptr.Deref(m.Status.Phase, ""), nil, oldM.Status.Conditions); err != nil {
			klog.Errorf("%v: error updating status: %v", oldM.Name, err)
			return false, fmt.Errorf("error updating status for machine %s: %w", oldM.Name, err)
		}

		return true, nil
	}

	return false, nil
}

func delayIfRequeueAfterError(err error) (reconcile.Result, error) {
	var requeueAfterError *RequeueAfterError
	if errors.As(err, &requeueAfterError) {
		klog.Infof("Actuator returned requeue-after error: %v", requeueAfterError)
		return reconcile.Result{Requeue: true, RequeueAfter: requeueAfterError.RequeueAfter}, nil
	}
	return reconcile.Result{}, err
}

func isInvalidMachineConfigurationError(err error) bool {
	var machineError *MachineError
	if errors.As(err, &machineError) {
		if machineError.Reason == machinev1.InvalidConfigurationMachineError {
			klog.Infof("Actuator returned invalid configuration error: %v", machineError)
			return true
		}
	}
	return false
}

// updateStatus is intended to ensure that the status of the Machine reflects the input to this function.
// Because the conditions are set on the machine outside of this function, we must pass the original state of the
// machine conditions so that the diff can be calculated properly within this function.
func (r *ReconcileMachine) updateStatus(ctx context.Context, machine *machinev1.Machine, phase string, failureCause error, originalConditions []machinev1.Condition) error {
	phaseChanged := false
	if ptr.Deref(machine.Status.Phase, "") != phase {
		klog.V(3).Infof("%v: going into phase %q", machine.GetName(), phase)

		phaseChanged = true
	}

	// Ensure the lifecycle hook conditions are accurate whenever the status is updated
	setLifecycleHookConditions(machine)

	// Conditions need to be deep copied as they are set outside of this function.
	// They will be restored after any updates to the base (done by patching annotations).
	conditions := conditions.DeepCopyConditions(machine.Status.Conditions)

	// A call to Patch will mutate our local copy of the machine to match what is stored in the API.
	// Before we make any changes to the status subresource on our local copy, we need to patch the object first,
	// otherwise our local changes to the status subresource will be lost.
	if phase == machinev1.PhaseFailed {
		err := r.patchFailedMachineInstanceAnnotation(ctx, machine)
		if err != nil {
			klog.Errorf("Failed to update machine %q: %v", machine.GetName(), err)
			return err
		}
	}

	// To ensure conditions can be patched properly, set the original conditions on the baseMachine.
	// This allows the difference to be calculated as part of the patch.
	baseMachine := machine.DeepCopy()
	baseMachine.Status.Conditions = originalConditions
	machine.Status.Conditions = conditions

	// Since we may have mutated the local copy of the machine above, we need to calculate baseToPatch here.
	// Any updates to the status must be done after this point.
	baseToPatch := client.MergeFrom(baseMachine)

	if phase == machinev1.PhaseFailed {
		if err := r.overrideFailedMachineProviderStatusState(machine); err != nil {
			klog.Errorf("Failed to update machine provider status %q: %v", machine.GetName(), err)
			return err
		}
	}

	machine.Status.Phase = &phase
	machine.Status.ErrorReason = nil
	machine.Status.ErrorMessage = nil
	if phase == machinev1.PhaseFailed && failureCause != nil {
		var machineError *MachineError
		if errors.As(failureCause, &machineError) {
			machine.Status.ErrorReason = &machineError.Reason
			machine.Status.ErrorMessage = &machineError.Message
		} else {
			errorMessage := failureCause.Error()
			machine.Status.ErrorMessage = &errorMessage
		}
	}

	if !reflect.DeepEqual(baseMachine.Status, machine.Status) {
		// Something on the status has been changed this reconcile
		now := metav1.NewTime(r.now())
		machine.Status.LastUpdated = &now
	}

	if err := r.Client.Status().Patch(ctx, machine, baseToPatch); err != nil {
		klog.Errorf("Failed to update machine status %q: %v", machine.GetName(), err)
		return err
	}

	// Update the metric after everything else has succeeded to prevent duplicate
	// entries when there are failures.
	// Only update when there is a change to the phase to avoid duplicating entries for
	// individual machines.
	if phaseChanged && phase != machinev1.PhaseDeleting {
		// Apart from deleting, update the transition metric
		// Deleting would always end up in the infinite bucket
		timeElapsed := r.now().Sub(machine.GetCreationTimestamp().Time).Seconds()
		metrics.MachinePhaseTransitionSeconds.With(map[string]string{"phase": phase}).Observe(timeElapsed)
	}

	return nil
}

func (r *ReconcileMachine) patchStatusAuthoritativeAPI(ctx context.Context, machine *machinev1.Machine, authoritativeAPI machinev1.MachineAuthority) error {
	baseToPatch := client.MergeFrom(machine.DeepCopy())
	machine.Status.AuthoritativeAPI = authoritativeAPI

	if err := r.Client.Status().Patch(ctx, machine, baseToPatch); err != nil {
		return fmt.Errorf("error patching machine status: %w", err)
	}

	return nil
}

func (r *ReconcileMachine) patchFailedMachineInstanceAnnotation(ctx context.Context, machine *machinev1.Machine) error {
	baseToPatch := client.MergeFrom(machine.DeepCopy())
	if machine.Annotations == nil {
		machine.Annotations = map[string]string{}
	}
	machine.Annotations[MachineInstanceStateAnnotationName] = unknownInstanceState
	if err := r.Client.Patch(ctx, machine, baseToPatch); err != nil {
		return err
	}
	return nil
}

// overrideFailedMachineProviderStatusState patches the state of the VM in the provider status if it is set.
// Not all providers set a state, but AWS, Azure, GCP and vSphere do.
// If the machine has gone into the Failed phase, and the providerStatus has already been set,
// the VM is in an unknown state. This function overrides the state.
func (r *ReconcileMachine) overrideFailedMachineProviderStatusState(machine *machinev1.Machine) error {
	if machine.Status.ProviderStatus == nil {
		return nil
	}

	// instanceState is used by AWS, GCP and vSphere; vmState is used by Azure.
	const instanceStateField = "instanceState"
	const vmStateField = "vmState"

	providerStatus, err := runtime.DefaultUnstructuredConverter.ToUnstructured(machine.Status.ProviderStatus)
	if err != nil {
		return fmt.Errorf("could not covert provider status to unstructured: %v", err)
	}

	// if the instanceState is set already, update it to unknown
	if _, found, err := unstructured.NestedString(providerStatus, instanceStateField); err == nil && found {
		if err := unstructured.SetNestedField(providerStatus, unknownInstanceState, instanceStateField); err != nil {
			return fmt.Errorf("could not set %s: %v", instanceStateField, err)
		}
	}

	// if the vmState is set already, update it to unknown
	if _, found, err := unstructured.NestedString(providerStatus, vmStateField); err == nil && found {
		if err := unstructured.SetNestedField(providerStatus, unknownInstanceState, vmStateField); err != nil {
			return fmt.Errorf("could not set %s: %v", instanceStateField, err)
		}
	}

	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(providerStatus, machine.Status.ProviderStatus); err != nil {
		return fmt.Errorf("could not convert provider status from unstructured: %v", err)
	}

	return nil
}

func validateMachine(m *machinev1.Machine) field.ErrorList {
	errors := field.ErrorList{}

	// validate spec.labels
	fldPath := field.NewPath("spec")
	if m.Labels[machinev1.MachineClusterIDLabel] == "" {
		errors = append(errors, field.Invalid(fldPath.Child("labels"), m.Labels, fmt.Sprintf("missing %v label.", machinev1.MachineClusterIDLabel)))
	}

	// validate provider config is set
	if m.Spec.ProviderSpec.Value == nil {
		errors = append(errors, field.Invalid(fldPath.Child("spec").Child("providerspec"), m.Spec.ProviderSpec, "value field must be set"))
	}

	return errors
}

func setLifecycleHookConditions(m *machinev1.Machine) {
	if len(m.Spec.LifecycleHooks.PreDrain) > 0 {
		conditions.Set(m, conditions.FalseCondition(
			machinev1.MachineDrainable,
			machinev1.MachineHookPresent,
			machinev1.ConditionSeverityWarning,
			"Drain operation currently blocked by: %+v", m.Spec.LifecycleHooks.PreDrain,
		))
	} else {
		conditions.MarkTrue(m, machinev1.MachineDrainable)
	}

	if len(m.Spec.LifecycleHooks.PreTerminate) > 0 {
		conditions.Set(m, conditions.FalseCondition(
			machinev1.MachineTerminable,
			machinev1.MachineHookPresent,
			machinev1.ConditionSeverityWarning,
			"Terminate operation currently blocked by: %+v", m.Spec.LifecycleHooks.PreTerminate,
		))
	} else {
		conditions.MarkTrue(m, machinev1.MachineTerminable)
	}
}

// now is used to get the current time. If the reconciler nowFunc is no nil this will be used instead of time.Now().
// This is only here so that tests can modify the time to check time based assertions.
func (r *ReconcileMachine) now() time.Time {
	if r.nowFunc != nil {
		return r.nowFunc()
	}
	return time.Now()
}

func machineIsProvisioned(machine *machinev1.Machine) bool {
	return len(machine.Status.Addresses) > 0 || ptr.Deref(machine.Spec.ProviderID, "") != ""
}

func machineHasNode(machine *machinev1.Machine) bool {
	return machine.Status.NodeRef != nil
}

func machineIsFailed(machine *machinev1.Machine) bool {
	return ptr.Deref(machine.Status.Phase, "") == machinev1.PhaseFailed
}

func nodeIsUnreachable(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionUnknown {
			return true
		}
	}

	return false
}

// writer implements io.Writer interface as a pass-through for klog.
type writer struct {
	logFunc func(args ...interface{})
}

// Write passes string(p) into writer's logFunc and always returns len(p)
func (w writer) Write(p []byte) (n int, err error) {
	w.logFunc(string(p))
	return len(p), nil
}
