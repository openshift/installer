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

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/scalesetvms"
	infracontroller "sigs.k8s.io/cluster-api-provider-azure/controllers"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/coalescing"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type (
	azureMachinePoolMachineReconcilerFactory func(*scope.MachinePoolMachineScope) (azure.Reconciler, error)

	// AzureMachinePoolMachineController handles Kubernetes change events for AzureMachinePoolMachine resources.
	AzureMachinePoolMachineController struct {
		client.Client
		Scheme            *runtime.Scheme
		Recorder          record.EventRecorder
		Timeouts          reconciler.Timeouts
		WatchFilterValue  string
		reconcilerFactory azureMachinePoolMachineReconcilerFactory
	}

	azureMachinePoolMachineReconciler struct {
		Scope              *scope.MachinePoolMachineScope
		scalesetVMsService *scalesetvms.Service
	}
)

// NewAzureMachinePoolMachineController creates a new AzureMachinePoolMachineController to handle updates to Azure Machine Pool Machines.
func NewAzureMachinePoolMachineController(c client.Client, recorder record.EventRecorder, timeouts reconciler.Timeouts, watchFilterValue string) *AzureMachinePoolMachineController {
	return &AzureMachinePoolMachineController{
		Client:            c,
		Recorder:          recorder,
		Timeouts:          timeouts,
		WatchFilterValue:  watchFilterValue,
		reconcilerFactory: newAzureMachinePoolMachineReconciler,
	}
}

// SetupWithManager initializes this controller with a manager.
func (ampmr *AzureMachinePoolMachineController) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options infracontroller.Options) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureMachinePoolMachineController.SetupWithManager",
		tele.KVP("controller", "AzureMachinePoolMachine"),
	)
	defer done()

	var r reconcile.Reconciler = ampmr
	if options.Cache != nil {
		r = coalescing.NewReconciler(ampmr, options.Cache, log)
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options.Options).
		For(&infrav1exp.AzureMachinePoolMachine{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, ampmr.WatchFilterValue)).
		Build(r)
	if err != nil {
		return errors.Wrapf(err, "error creating controller")
	}

	// Add a watch on AzureMachinePool for model changes
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &infrav1exp.AzureMachinePool{}),
		handler.EnqueueRequestsFromMapFunc(AzureMachinePoolToAzureMachinePoolMachines(ctx, mgr.GetClient(), log)),
		MachinePoolModelHasChanged(log),
		predicates.ResourceNotPausedAndHasFilterLabel(log, ampmr.WatchFilterValue),
	); err != nil {
		return errors.Wrapf(err, "failed adding a watch for AzureMachinePool model changes")
	}

	// Add a watch on CAPI Machines for MachinePool Machines
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Machine{}),
		handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1exp.GroupVersion.WithKind("AzureMachinePoolMachine"))),
		predicates.ResourceNotPausedAndHasFilterLabel(log, ampmr.WatchFilterValue),
	); err != nil {
		return errors.Wrapf(err, "failed adding a watch for Machine model changes")
	}

	return nil
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepools/status,verbs=get
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepoolmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepoolmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch

// Reconcile idempotently gets, creates, and updates a machine pool.
func (ampmr *AzureMachinePoolMachineController) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, logger, done := tele.StartSpanWithLogger(
		ctx,
		"controllers.AzureMachinePoolMachineController.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "AzureMachinePoolMachine"),
	)
	defer done()

	logger = logger.WithValues("namespace", req.Namespace, "azureMachinePoolMachine", req.Name)

	ctx, cancel := context.WithTimeout(ctx, ampmr.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	azureMachine := &infrav1exp.AzureMachinePoolMachine{}
	err := ampmr.Get(ctx, req.NamespacedName, azureMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	logger.V(2).Info("Fetching cluster for AzureMachinePoolMachine", "ampm", azureMachine.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, ampmr.Client, azureMachine.ObjectMeta)
	if err != nil {
		logger.Info("AzureMachinePoolMachine is missing cluster label or cluster does not exist")
		return reconcile.Result{}, nil
	}

	logger = logger.WithValues("cluster", cluster.Name)

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, azureMachine) {
		logger.Info("AzureMachinePoolMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	clusterScope, err := infracontroller.GetClusterScoper(ctx, logger, ampmr.Client, cluster, ampmr.Timeouts)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to create cluster scope for cluster %s/%s", cluster.Namespace, cluster.Name)
	}

	logger.V(2).Info("Fetching AzureMachinePool with object meta", "meta", azureMachine.ObjectMeta)
	// Fetch the owning AzureMachinePool (VMSS)
	azureMachinePool, err := infracontroller.GetOwnerAzureMachinePool(ctx, ampmr.Client, azureMachine.ObjectMeta)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("AzureMachinePool not found error missing, removing finalizer", "azureMachinePoolMachine", azureMachine.Name)
			controllerutil.RemoveFinalizer(azureMachine, infrav1exp.AzureMachinePoolMachineFinalizer)
			return reconcile.Result{}, ampmr.Client.Update(ctx, azureMachine)
		}
		return reconcile.Result{}, err
	}
	if azureMachinePool == nil {
		logger.Info("AzureMachinePool not found error missing, removing finalizer", "azureMachinePoolMachine", azureMachine.Name)
		controllerutil.RemoveFinalizer(azureMachine, infrav1exp.AzureMachinePoolMachineFinalizer)
		return reconcile.Result{}, ampmr.Client.Update(ctx, azureMachine)
	}

	logger = logger.WithValues("azureMachinePool", azureMachinePool.Name)

	// Fetch the CAPI MachinePool.
	machinePool, err := infracontroller.GetOwnerMachinePool(ctx, ampmr.Client, azureMachinePool.ObjectMeta)
	if err != nil && !apierrors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	if machinePool != nil {
		logger = logger.WithValues("machinePool", machinePool.Name)
	}

	// Fetch the CAPI Machine.
	machine, err := util.GetOwnerMachine(ctx, ampmr.Client, azureMachine.ObjectMeta)
	if err != nil && !apierrors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	if machine != nil {
		logger = logger.WithValues("machine", machine.Name)
	} else {
		logger.Info("Waiting for Machine Controller to set OwnerRef on AzureMachinePoolMachine")
		return reconcile.Result{}, nil
	}

	// Create the machine pool scope
	machineScope, err := scope.NewMachinePoolMachineScope(scope.MachinePoolMachineScopeParams{
		Client:                  ampmr.Client,
		MachinePool:             machinePool,
		AzureMachinePool:        azureMachinePool,
		AzureMachinePoolMachine: azureMachine,
		Machine:                 machine,
		ClusterScope:            clusterScope,
	})
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create scope")
	}

	// Always close the scope when exiting this function so we can persist any AzureMachine changes.
	defer func() {
		if err := machineScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted machine pools machine
	if !azureMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return ampmr.reconcileDelete(ctx, machineScope, clusterScope)
	}

	if !cluster.Status.InfrastructureReady {
		logger.Info("Cluster infrastructure is not ready yet")
		return reconcile.Result{}, nil
	}

	// Handle non-deleted machine pools
	return ampmr.reconcileNormal(ctx, machineScope)
}

func (ampmr *AzureMachinePoolMachineController) reconcileNormal(ctx context.Context, machineScope *scope.MachinePoolMachineScope) (_ reconcile.Result, reterr error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachinePoolMachineController.reconcileNormal")
	defer done()

	log.Info("Reconciling AzureMachinePoolMachine")
	// If the AzureMachine is in an error state, return early.
	if machineScope.AzureMachinePool.Status.FailureReason != nil || machineScope.AzureMachinePool.Status.FailureMessage != nil {
		log.Info("Error state detected, skipping reconciliation")
		return reconcile.Result{}, nil
	}

	ampms, err := ampmr.reconcilerFactory(machineScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create AzureMachinePoolMachine reconciler")
	}
	if err := ampms.Reconcile(ctx); err != nil {
		// Handle transient and terminal errors
		var reconcileError azure.ReconcileError
		if errors.As(err, &reconcileError) {
			if reconcileError.IsTerminal() {
				log.Error(err, "failed to reconcile AzureMachinePool", "name", machineScope.Name())
				return reconcile.Result{}, nil
			}

			if reconcileError.IsTransient() {
				log.V(4).Info("failed to reconcile AzureMachinePoolMachine", "name", machineScope.Name(), "transient_error", err)
				return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
			}

			return reconcile.Result{}, errors.Wrap(err, "failed to reconcile AzureMachinePool")
		}

		return reconcile.Result{}, err
	}

	state := machineScope.ProvisioningState()
	switch state {
	case infrav1.Failed:
		ampmr.Recorder.Eventf(machineScope.AzureMachinePoolMachine, corev1.EventTypeWarning, "FailedVMState", "Azure scale set VM is in failed state")
		machineScope.SetFailureReason(capierrors.UpdateMachineError)
		machineScope.SetFailureMessage(errors.Errorf("Azure VM state is %s", state))
	case infrav1.Deleting:
		if err := ampmr.Client.Delete(ctx, machineScope.AzureMachinePoolMachine); err != nil {
			return reconcile.Result{}, errors.Wrap(err, "machine pool machine failed to be deleted when deleting")
		}
	}

	log.V(2).Info(fmt.Sprintf("Scale Set VM is %s", state), "id", machineScope.ProviderID())

	bootstrappingCondition := conditions.Get(machineScope.AzureMachinePoolMachine, infrav1.BootstrapSucceededCondition)
	if bootstrappingCondition != nil && bootstrappingCondition.Reason == infrav1.BootstrapFailedReason {
		return reconcile.Result{}, nil
	}

	if !infrav1.IsTerminalProvisioningState(state) || !machineScope.IsReady() {
		log.V(2).Info("Requeuing", "state", state, "ready", machineScope.IsReady())
		// we are in a non-terminal state, retry in a bit
		return reconcile.Result{
			RequeueAfter: 30 * time.Second,
		}, nil
	}

	return reconcile.Result{}, nil
}

func (ampmr *AzureMachinePoolMachineController) reconcileDelete(ctx context.Context, machineScope *scope.MachinePoolMachineScope, clusterScope infracontroller.ClusterScoper) (_ reconcile.Result, reterr error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachinePoolMachineController.reconcileDelete")
	defer done()

	if !infracontroller.ShouldDeleteIndividualResources(ctx, clusterScope) {
		log.Info("Skipping VMSS VM deletion as the whole resource group is being deleted")

		controllerutil.RemoveFinalizer(machineScope.AzureMachinePoolMachine, infrav1exp.AzureMachinePoolMachineFinalizer)
		return reconcile.Result{}, nil
	}

	if !machineScope.AzureMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Skipping VMSS VM deletion as VMSS delete will delete individual instances")

		controllerutil.RemoveFinalizer(machineScope.AzureMachinePoolMachine, infrav1exp.AzureMachinePoolMachineFinalizer)
		return reconcile.Result{}, nil
	}

	log.Info("Deleting AzureMachinePoolMachine")

	// deleting a single machine
	// 1) delete the infrastructure, node drain already done by owner Machine
	// 2) remove finalizer

	ampms, err := ampmr.reconcilerFactory(machineScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create AzureMachinePoolMachine reconciler")
	}
	if err := ampms.Delete(ctx); err != nil {
		// Handle transient and terminal errors
		var reconcileError azure.ReconcileError
		if errors.As(err, &reconcileError) {
			if reconcileError.IsTerminal() {
				log.Error(err, "failed to delete AzureMachinePoolMachine", "name", machineScope.Name())
				return reconcile.Result{}, nil
			}

			if reconcileError.IsTransient() {
				log.V(4).Info("failed to delete AzureMachinePoolMachine", "name", machineScope.Name(), "transient_error", err)
				return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
			}

			return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile AzureMachinePool")
		}

		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func newAzureMachinePoolMachineReconciler(scope *scope.MachinePoolMachineScope) (azure.Reconciler, error) {
	scaleSetVMsSvc, err := scalesetvms.NewService(scope)
	if err != nil {
		return nil, err
	}
	return &azureMachinePoolMachineReconciler{
		Scope:              scope,
		scalesetVMsService: scaleSetVMsSvc,
	}, nil
}

// Reconcile will reconcile the state of the Machine Pool Machine with the state of the Azure VMSS VM.
func (r *azureMachinePoolMachineReconciler) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "controllers.azureMachinePoolMachineReconciler.Reconcile")
	defer done()

	if err := r.scalesetVMsService.Reconcile(ctx); err != nil {
		return errors.Wrap(err, "failed to reconcile scalesetVMs")
	}

	if err := r.Scope.UpdateNodeStatus(ctx); err != nil {
		return errors.Wrap(err, "failed to update VMSS VM node status")
	}

	if err := r.Scope.UpdateInstanceStatus(ctx); err != nil {
		return errors.Wrap(err, "failed to update VMSS VM instance status")
	}

	return nil
}

// Delete will attempt to drain and delete the Azure VMSS VM.
func (r *azureMachinePoolMachineReconciler) Delete(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.azureMachinePoolMachineReconciler.Delete")
	defer done()

	defer func() {
		if err := r.Scope.UpdateNodeStatus(ctx); err != nil {
			log.V(4).Info("failed to update VMSS VM node status during delete")
		}

		if err := r.Scope.UpdateInstanceStatus(ctx); err != nil {
			log.V(4).Info("failed to update VMSS VM instance status during delete")
		}
	}()

	if err := r.scalesetVMsService.Delete(ctx); err != nil {
		return errors.Wrap(err, "failed to reconcile scalesetVMs")
	}

	// no long running operation, so we are finished deleting the resource. Remove the finalizer.
	controllerutil.RemoveFinalizer(r.Scope.AzureMachinePoolMachine, infrav1exp.AzureMachinePoolMachineFinalizer)

	return nil
}
