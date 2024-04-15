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
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
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

// AzureMachineReconciler reconciles an AzureMachine object.
type AzureMachineReconciler struct {
	client.Client
	Recorder                  record.EventRecorder
	Timeouts                  reconciler.Timeouts
	WatchFilterValue          string
	createAzureMachineService azureMachineServiceCreator
}

type azureMachineServiceCreator func(machineScope *scope.MachineScope) (*azureMachineService, error)

// NewAzureMachineReconciler returns a new AzureMachineReconciler instance.
func NewAzureMachineReconciler(client client.Client, recorder record.EventRecorder, timeouts reconciler.Timeouts, watchFilterValue string) *AzureMachineReconciler {
	amr := &AzureMachineReconciler{
		Client:           client,
		Recorder:         recorder,
		Timeouts:         timeouts,
		WatchFilterValue: watchFilterValue,
	}

	amr.createAzureMachineService = newAzureMachineService

	return amr
}

// SetupWithManager initializes this controller with a manager.
func (amr *AzureMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options Options) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureMachineReconciler.SetupWithManager",
		tele.KVP("controller", "AzureMachine"),
	)
	defer done()

	var r reconcile.Reconciler = amr
	if options.Cache != nil {
		r = coalescing.NewReconciler(amr, options.Cache, log)
	}

	// create mapper to transform incoming AzureClusters into AzureMachine requests
	azureClusterToAzureMachinesMapper, err := AzureClusterToAzureMachinesMapper(ctx, amr.Client, &infrav1.AzureMachineList{}, mgr.GetScheme(), log)
	if err != nil {
		return errors.Wrap(err, "failed to create AzureCluster to AzureMachines mapper")
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options.Options).
		For(&infrav1.AzureMachine{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(log, amr.WatchFilterValue)).
		// watch for changes in CAPI Machine resources
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("AzureMachine"))),
		).
		// watch for changes in AzureCluster
		Watches(
			&infrav1.AzureCluster{},
			handler.EnqueueRequestsFromMapFunc(azureClusterToAzureMachinesMapper),
		).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	azureMachineMapper, err := util.ClusterToTypedObjectsMapper(amr.Client, &infrav1.AzureMachineList{}, mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "failed to create mapper for Cluster to AzureMachines")
	}

	// Add a watch on clusterv1.Cluster object for pause/unpause & ready notifications.
	if err := c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(azureMachineMapper),
		ClusterPauseChangeAndInfrastructureReady(log),
		predicates.ResourceHasFilterLabel(log, amr.WatchFilterValue),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready clusters")
	}

	return nil
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch

// Reconcile idempotently gets, creates, and updates a machine.
func (amr *AzureMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, amr.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(
		ctx,
		"controllers.AzureMachineReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "AzureMachine"),
	)
	defer done()

	// Fetch the AzureMachine VM.
	azureMachine := &infrav1.AzureMachine{}
	err := amr.Get(ctx, req.NamespacedName, azureMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, amr.Client, azureMachine.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if machine == nil {
		amr.Recorder.Eventf(azureMachine, corev1.EventTypeNormal, "Machine controller dependency not yet met", "Machine Controller has not yet set OwnerRef")
		log.Info("Machine Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("machine", machine.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, amr.Client, machine.ObjectMeta)
	if err != nil {
		amr.Recorder.Eventf(azureMachine, corev1.EventTypeNormal, "Unable to get cluster from metadata", "Machine is missing cluster label or cluster does not exist")
		log.Info("Machine is missing cluster label or cluster does not exist")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	log = log.WithValues("AzureCluster", cluster.Spec.InfrastructureRef.Name)
	azureClusterName := client.ObjectKey{
		Namespace: azureMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	azureCluster := &infrav1.AzureCluster{}
	if err := amr.Client.Get(ctx, azureClusterName, azureCluster); err != nil {
		amr.Recorder.Eventf(azureMachine, corev1.EventTypeNormal, "AzureCluster unavailable", "AzureCluster is not available yet")
		log.Info("AzureCluster is not available yet")
		return reconcile.Result{}, nil
	}

	// Create the cluster scope
	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:       amr.Client,
		Cluster:      cluster,
		AzureCluster: azureCluster,
		Timeouts:     amr.Timeouts,
	})
	if err != nil {
		amr.Recorder.Eventf(azureCluster, corev1.EventTypeWarning, "Error creating the cluster scope", err.Error())
		return reconcile.Result{}, err
	}

	// Create the machine scope
	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Client:       amr.Client,
		Machine:      machine,
		AzureMachine: azureMachine,
		ClusterScope: clusterScope,
	})
	if err != nil {
		amr.Recorder.Eventf(azureMachine, corev1.EventTypeWarning, "Error creating the machine scope", err.Error())
		return reconcile.Result{}, errors.Wrap(err, "failed to create scope")
	}

	// Always close the scope when exiting this function so we can persist any AzureMachine changes.
	defer func() {
		if err := machineScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, azureMachine) {
		log.Info("AzureMachine or linked Cluster is marked as paused. Won't reconcile normally")
		return amr.reconcilePause(ctx, machineScope)
	}

	// Handle deleted machines
	if !azureMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return amr.reconcileDelete(ctx, machineScope, clusterScope)
	}

	// Handle non-deleted machines
	return amr.reconcileNormal(ctx, machineScope, clusterScope)
}

func (amr *AzureMachineReconciler) reconcileNormal(ctx context.Context, machineScope *scope.MachineScope, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachineReconciler.reconcileNormal")
	defer done()

	log.Info("Reconciling AzureMachine")
	// If the AzureMachine is in an error state, return early.
	if machineScope.AzureMachine.Status.FailureReason != nil || machineScope.AzureMachine.Status.FailureMessage != nil {
		log.Info("Error state detected, skipping reconciliation")
		return reconcile.Result{}, nil
	}

	// Register our finalizer immediately to avoid orphaning Azure resources on delete
	needsPatch := controllerutil.AddFinalizer(machineScope.AzureMachine, infrav1.MachineFinalizer)
	// Register the block-move annotation immediately to avoid moving un-paused ASO resources
	needsPatch = AddBlockMoveAnnotation(machineScope.AzureMachine) || needsPatch
	if needsPatch {
		if err := machineScope.PatchObject(ctx); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Make sure the Cluster Infrastructure is ready.
	if !clusterScope.Cluster.Status.InfrastructureReady {
		log.Info("Cluster infrastructure is not ready yet")
		conditions.MarkFalse(machineScope.AzureMachine, infrav1.VMRunningCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
		return reconcile.Result{}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(machineScope.AzureMachine, infrav1.VMRunningCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return reconcile.Result{}, nil
	}

	var reconcileError azure.ReconcileError

	// Initialize the cache to be used by the AzureMachine services.
	err := machineScope.InitMachineCache(ctx)
	if err != nil {
		if errors.As(err, &reconcileError) && reconcileError.IsTerminal() {
			amr.Recorder.Eventf(machineScope.AzureMachine, corev1.EventTypeWarning, "SKUNotFound", errors.Wrap(err, "failed to initialize machine cache").Error())
			log.Error(err, "Failed to initialize machine cache")
			machineScope.SetFailureReason(capierrors.InvalidConfigurationMachineError)
			machineScope.SetFailureMessage(err)
			machineScope.SetNotReady()
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, errors.Wrap(err, "failed to init machine scope cache")
	}

	// Mark the AzureMachine as failed if the identities are not ready.
	cond := conditions.Get(machineScope.AzureMachine, infrav1.VMIdentitiesReadyCondition)
	if cond != nil && cond.Status == corev1.ConditionFalse && cond.Reason == infrav1.UserAssignedIdentityMissingReason {
		amr.Recorder.Eventf(machineScope.AzureMachine, corev1.EventTypeWarning, infrav1.UserAssignedIdentityMissingReason, "VM is unhealthy")
		machineScope.SetFailureReason(capierrors.UnsupportedChangeMachineError)
		machineScope.SetFailureMessage(errors.New("VM identities are not ready"))
		return reconcile.Result{}, errors.New("VM identities are not ready")
	}

	ams, err := amr.createAzureMachineService(machineScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create azure machine service")
	}

	if err := ams.Reconcile(ctx); err != nil {
		// This means that a VM was created and managed by this controller, but is not present anymore.
		// In this case, we mark it as failed and leave it to MHC for remediation
		if errors.As(err, &azure.VMDeletedError{}) {
			amr.Recorder.Eventf(machineScope.AzureMachine, corev1.EventTypeWarning, "VMDeleted", errors.Wrap(err, "failed to reconcile AzureMachine").Error())
			machineScope.SetFailureReason(capierrors.UpdateMachineError)
			machineScope.SetFailureMessage(err)
			machineScope.SetNotReady()
			machineScope.SetVMState(infrav1.Deleted)
			return reconcile.Result{}, errors.Wrap(err, "failed to reconcile AzureMachine")
		}

		// Handle transient and terminal errors
		if errors.As(err, &reconcileError) {
			if reconcileError.IsTerminal() {
				amr.Recorder.Eventf(machineScope.AzureMachine, corev1.EventTypeWarning, "ReconcileError", errors.Wrapf(err, "failed to reconcile AzureMachine").Error())
				log.Error(err, "failed to reconcile AzureMachine", "name", machineScope.Name())
				machineScope.SetFailureReason(capierrors.CreateMachineError)
				machineScope.SetFailureMessage(err)
				machineScope.SetNotReady()
				machineScope.SetVMState(infrav1.Failed)
				return reconcile.Result{}, nil
			}

			if reconcileError.IsTransient() {
				if azure.IsOperationNotDoneError(reconcileError) {
					log.V(2).Info(fmt.Sprintf("AzureMachine reconcile not done: %s", reconcileError.Error()))
				} else {
					log.V(2).Info(fmt.Sprintf("transient failure to reconcile AzureMachine, retrying: %s", reconcileError.Error()))
				}
				return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
			}
		}
		amr.Recorder.Eventf(machineScope.AzureMachine, corev1.EventTypeWarning, "ReconcileError", errors.Wrapf(err, "failed to reconcile AzureMachine").Error())
		return reconcile.Result{}, errors.Wrap(err, "failed to reconcile AzureMachine")
	}

	machineScope.SetReady()

	return reconcile.Result{}, nil
}

//nolint:unparam // Always returns an empty struct for reconcile.Result
func (amr *AzureMachineReconciler) reconcilePause(ctx context.Context, machineScope *scope.MachineScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachine.reconcilePause")
	defer done()

	log.Info("Reconciling AzureMachine pause")

	ams, err := amr.createAzureMachineService(machineScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create azure machine service")
	}

	if err := ams.Pause(ctx); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to pause azure machine services")
	}
	RemoveBlockMoveAnnotation(machineScope.AzureMachine)

	return reconcile.Result{}, nil
}

func (amr *AzureMachineReconciler) reconcileDelete(ctx context.Context, machineScope *scope.MachineScope, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachineReconciler.reconcileDelete")
	defer done()

	log.Info("Handling deleted AzureMachine")
	conditions.MarkFalse(machineScope.AzureMachine, infrav1.VMRunningCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := machineScope.PatchObject(ctx); err != nil {
		return reconcile.Result{}, err
	}

	if ShouldDeleteIndividualResources(ctx, clusterScope) {
		log.Info("Deleting AzureMachine")
		ams, err := amr.createAzureMachineService(machineScope)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to create azure machine service")
		}

		if err := ams.Delete(ctx); err != nil {
			// Handle transient errors
			var reconcileError azure.ReconcileError
			if errors.As(err, &reconcileError) {
				if reconcileError.IsTransient() {
					if azure.IsOperationNotDoneError(reconcileError) {
						log.V(2).Info(fmt.Sprintf("AzureMachine delete not done: %s", reconcileError.Error()))
					} else {
						log.V(2).Info("transient failure to delete AzureMachine, retrying")
					}
					return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
				}
			}

			amr.Recorder.Eventf(machineScope.AzureMachine, corev1.EventTypeWarning, "Error deleting AzureMachine", errors.Wrapf(err, "error deleting AzureMachine %s/%s", machineScope.Namespace(), machineScope.Name()).Error())
			return reconcile.Result{}, errors.Wrapf(err, "error deleting AzureMachine %s/%s", machineScope.Namespace(), machineScope.Name())
		}
	} else {
		log.Info("Skipping AzureMachine Deletion; will delete whole resource group.")
	}

	// we're done deleting this AzureMachine so remove the finalizer.
	log.Info("Removing finalizer from AzureMachine")
	controllerutil.RemoveFinalizer(machineScope.AzureMachine, infrav1.MachineFinalizer)

	return reconcile.Result{}, nil
}
