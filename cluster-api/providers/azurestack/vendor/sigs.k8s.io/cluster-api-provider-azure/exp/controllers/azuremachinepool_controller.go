/*
Copyright 2020 The Kubernetes Authors.

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
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/external"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	infracontroller "sigs.k8s.io/cluster-api-provider-azure/controllers"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/coalescing"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

type (
	// AzureMachinePoolReconciler reconciles an AzureMachinePool object.
	AzureMachinePoolReconciler struct {
		client.Client
		Scheme                        *runtime.Scheme
		Recorder                      record.EventRecorder
		Timeouts                      reconciler.Timeouts
		WatchFilterValue              string
		createAzureMachinePoolService azureMachinePoolServiceCreator
		externalTracker               external.ObjectTracker
		CredentialCache               azure.CredentialCache
	}

	// annotationReaderWriter provides an interface to read and write annotations.
	annotationReaderWriter interface {
		GetAnnotations() map[string]string
		SetAnnotations(annotations map[string]string)
	}
)

type azureMachinePoolServiceCreator func(machinePoolScope *scope.MachinePoolScope) (*azureMachinePoolService, error)

// NewAzureMachinePoolReconciler returns a new AzureMachinePoolReconciler instance.
func NewAzureMachinePoolReconciler(client client.Client, recorder record.EventRecorder, timeouts reconciler.Timeouts, watchFilterValue string, credCache azure.CredentialCache) *AzureMachinePoolReconciler {
	ampr := &AzureMachinePoolReconciler{
		Client:           client,
		Recorder:         recorder,
		Timeouts:         timeouts,
		WatchFilterValue: watchFilterValue,
		CredentialCache:  credCache,
	}

	ampr.createAzureMachinePoolService = newAzureMachinePoolService

	return ampr
}

// SetupWithManager initializes this controller with a manager.
func (ampr *AzureMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options infracontroller.Options) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureMachinePoolReconciler.SetupWithManager",
		tele.KVP("controller", "AzureMachinePool"),
	)
	defer done()

	var r reconcile.Reconciler = ampr
	if options.Cache != nil {
		r = coalescing.NewReconciler(ampr, options.Cache, log)
	}

	// create mappers to transform incoming AzureClusters and AzureManagedClusters into AzureMachinePool requests
	azureClusterMapper, err := AzureClusterToAzureMachinePoolsMapper(ctx, ampr.Client, mgr.GetScheme(), log)
	if err != nil {
		return errors.Wrapf(err, "failed to create AzureCluster to AzureMachinePools mapper")
	}
	azureManagedControlPlaneMapper, err := AzureManagedControlPlaneToAzureMachinePoolsMapper(ctx, ampr.Client, mgr.GetScheme(), log)
	if err != nil {
		return errors.Wrapf(err, "failed to create AzureManagedCluster to AzureMachinePools mapper")
	}

	azureMachinePoolMapper, err := util.ClusterToTypedObjectsMapper(ampr.Client, &infrav1exp.AzureMachinePoolList{}, mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "failed to create mapper for Cluster to AzureMachines")
	}

	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options.Options).
		For(&infrav1exp.AzureMachinePool{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, ampr.WatchFilterValue)).
		// watch for changes in CAPI MachinePool resources
		Watches(
			&expv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(MachinePoolToInfrastructureMapFunc(infrav1exp.GroupVersion.WithKind(infrav1.AzureMachinePoolKind), log)),
		).
		// watch for changes in AzureCluster resources
		Watches(
			&infrav1.AzureCluster{},
			handler.EnqueueRequestsFromMapFunc(azureClusterMapper),
		).
		// watch for changes in AzureManagedControlPlane resources
		Watches(
			&infrav1.AzureManagedControlPlane{},
			handler.EnqueueRequestsFromMapFunc(azureManagedControlPlaneMapper),
		).
		Watches(
			&infrav1exp.AzureMachinePoolMachine{},
			handler.EnqueueRequestsFromMapFunc(AzureMachinePoolMachineMapper(mgr.GetScheme(), log)),
			builder.WithPredicates(
				MachinePoolMachineHasStateOrVersionChange(log),
				predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, ampr.WatchFilterValue),
			),
		).
		// Add a watch on clusterv1.Cluster object for unpause & ready notifications.
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(azureMachinePoolMapper),
			builder.WithPredicates(
				infracontroller.ClusterPauseChangeAndInfrastructureReady(mgr.GetScheme(), log),
				predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, ampr.WatchFilterValue),
			),
		).Build(r)
	if err != nil {
		return fmt.Errorf("creating new controller manager: %w", err)
	}

	predicateLog := ptr.To(ctrl.LoggerFrom(ctx).WithValues("controller", "azuremachinepool"))
	ampr.externalTracker = external.ObjectTracker{
		Controller:      controller,
		Cache:           mgr.GetCache(),
		Scheme:          mgr.GetScheme(),
		PredicateLogger: predicateLog,
	}

	return nil
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepools,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepools/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=*,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepoolmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinepoolmachines/status,verbs=get
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch

// Reconcile idempotently gets, creates, and updates a machine pool.
func (ampr *AzureMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, logger, done := tele.StartSpanWithLogger(
		ctx,
		"controllers.AzureMachinePoolReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", infrav1.AzureMachinePoolKind),
	)
	defer done()
	ctx, cancel := context.WithTimeout(ctx, ampr.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	logger = logger.WithValues("namespace", req.Namespace, "azureMachinePool", req.Name)

	azMachinePool := &infrav1exp.AzureMachinePool{}
	err := ampr.Get(ctx, req.NamespacedName, azMachinePool)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the CAPI MachinePool.
	machinePool, err := infracontroller.GetOwnerMachinePool(ctx, ampr.Client, azMachinePool.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if machinePool == nil {
		logger.V(2).Info("MachinePool Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	logger = logger.WithValues("machinePool", machinePool.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, ampr.Client, machinePool.ObjectMeta)
	if err != nil {
		logger.V(2).Info("MachinePool is missing cluster label or cluster does not exist")
		return reconcile.Result{}, nil
	}

	logger = logger.WithValues("cluster", cluster.Name)

	clusterScope, err := infracontroller.GetClusterScoper(ctx, logger, ampr.Client, cluster, ampr.Timeouts, ampr.CredentialCache)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to create cluster scope for cluster %s/%s", cluster.Namespace, cluster.Name)
	}

	// Create the machine pool scope
	machinePoolScope, err := scope.NewMachinePoolScope(scope.MachinePoolScopeParams{
		Client:           ampr.Client,
		MachinePool:      machinePool,
		AzureMachinePool: azMachinePool,
		ClusterScope:     clusterScope,
	})
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create machinepool scope")
	}

	// Always close the scope when exiting this function so we can persist any AzureMachine changes.
	defer func() {
		if err := machinePoolScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, azMachinePool) {
		logger.V(2).Info("AzureMachinePool or linked Cluster is marked as paused. Won't reconcile normally")
		return ampr.reconcilePause(ctx, machinePoolScope)
	}

	// Handle deleted machine pools
	if !azMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
		return ampr.reconcileDelete(ctx, machinePoolScope, clusterScope)
	}

	// Handle non-deleted machine pools
	return ampr.reconcileNormal(ctx, machinePoolScope, cluster)
}

func (ampr *AzureMachinePoolReconciler) reconcileNormal(ctx context.Context, machinePoolScope *scope.MachinePoolScope, cluster *clusterv1.Cluster) (_ reconcile.Result, reterr error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachinePoolReconciler.reconcileNormal")
	defer done()

	log.Info("Reconciling AzureMachinePool")

	// If the AzureMachine is in an error state, return early.
	if machinePoolScope.AzureMachinePool.Status.FailureReason != nil || machinePoolScope.AzureMachinePool.Status.FailureMessage != nil {
		log.Info("Error state detected, skipping reconciliation")
		return reconcile.Result{}, nil
	}

	// Register the finalizer immediately to avoid orphaning Azure resources on delete
	needsPatch := controllerutil.AddFinalizer(machinePoolScope.AzureMachinePool, expv1.MachinePoolFinalizer)
	needsPatch = machinePoolScope.SetInfrastructureMachineKind() || needsPatch
	// Register the block-move annotation immediately to avoid moving un-paused ASO resources
	needsPatch = infracontroller.AddBlockMoveAnnotation(machinePoolScope.AzureMachinePool) || needsPatch
	if needsPatch {
		if err := machinePoolScope.PatchObject(ctx); err != nil {
			return reconcile.Result{}, err
		}
	}

	if !cluster.Status.InfrastructureReady {
		log.Info("Cluster infrastructure is not ready yet")
		return reconcile.Result{}, nil
	}

	// Add a Watch to the referenced Bootstrap Config
	if machinePoolScope.MachinePool.Spec.Template.Spec.Bootstrap.ConfigRef != nil {
		ref := machinePoolScope.MachinePool.Spec.Template.Spec.Bootstrap.ConfigRef
		obj, err := external.Get(ctx, ampr.Client, ref)
		if err != nil {
			if apierrors.IsNotFound(errors.Cause(err)) {
				return reconcile.Result{}, errors.Wrapf(err, "could not find %v %q for MachinePool %q in namespace %q, requeuing while searching for bootstrap ConfigRef",
					ref.GroupVersionKind(), ref.Name, machinePoolScope.MachinePool.Name, ref.Namespace)
			}
			return reconcile.Result{}, err
		}

		// Ensure we add a watch to the external object, if there isn't one already.
		if err := ampr.externalTracker.Watch(log, obj,
			handler.EnqueueRequestsFromMapFunc(BootstrapConfigToInfrastructureMapFunc(ampr.Client, *ampr.externalTracker.PredicateLogger)),
			predicates.ResourceIsChanged(ampr.Client.Scheme(), *ampr.externalTracker.PredicateLogger)); err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "could not add a watcher to the object %v %q for MachinePool %q in namespace %q, requeuing",
				ref.GroupVersionKind(), ref.Name, machinePoolScope.MachinePool.Name, ref.Namespace)
		}
	}

	// Make sure bootstrap data is available and populated.
	if machinePoolScope.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName == nil {
		log.Info("Bootstrap data secret reference is not yet available")
		return reconcile.Result{}, nil
	}

	var reconcileError azure.ReconcileError

	// Initialize the cache to be used by the AzureMachine services.
	err := machinePoolScope.InitMachinePoolCache(ctx)
	if err != nil {
		if errors.As(err, &reconcileError) && reconcileError.IsTerminal() {
			ampr.Recorder.Eventf(machinePoolScope.AzureMachinePool, corev1.EventTypeWarning, "SKUNotFound", errors.Wrap(err, "failed to initialize machinepool cache").Error())
			log.Error(err, "Failed to initialize machinepool cache")
			machinePoolScope.SetFailureReason(azure.InvalidConfiguration)
			machinePoolScope.SetFailureMessage(err)
			machinePoolScope.SetNotReady()
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, errors.Wrap(err, "failed to init machinepool scope cache")
	}

	ams, err := ampr.createAzureMachinePoolService(machinePoolScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed creating a newAzureMachinePoolService")
	}

	if err := ams.Reconcile(ctx); err != nil {
		// Handle transient and terminal errors
		var reconcileError azure.ReconcileError
		if errors.As(err, &reconcileError) {
			if reconcileError.IsTerminal() {
				log.Error(err, "failed to reconcile AzureMachinePool", "name", machinePoolScope.Name())
				return reconcile.Result{}, nil
			}

			if reconcileError.IsTransient() {
				if azure.IsOperationNotDoneError(reconcileError) {
					log.V(2).Info(fmt.Sprintf("AzureMachinePool reconcile not done: %s", reconcileError.Error()))
				} else {
					log.V(2).Info(fmt.Sprintf("transient failure to reconcile AzureMachinePool, retrying: %s", reconcileError.Error()))
				}
				return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
			}

			return reconcile.Result{}, errors.Wrap(err, "failed to reconcile AzureMachinePool")
		}

		return reconcile.Result{}, err
	}

	log.V(2).Info("Scale Set reconciled", "id",
		machinePoolScope.ProviderID(), "state", machinePoolScope.ProvisioningState())

	switch machinePoolScope.ProvisioningState() {
	case infrav1.Deleting:
		log.Info("Unexpected scale set deletion", "id", machinePoolScope.ProviderID())
		ampr.Recorder.Eventf(machinePoolScope.AzureMachinePool, corev1.EventTypeWarning, "UnexpectedVMDeletion", "Unexpected Azure scale set deletion")
	case infrav1.Failed:
		log.Info("Unexpected scale set failure", "id", machinePoolScope.ProviderID())
		ampr.Recorder.Eventf(machinePoolScope.AzureMachinePool, corev1.EventTypeWarning, "UnexpectedVMFailure", "Unexpected Azure scale set failure")
	}

	if machinePoolScope.NeedsRequeue() {
		return reconcile.Result{
			RequeueAfter: 30 * time.Second,
		}, nil
	}

	return reconcile.Result{}, nil
}

func (ampr *AzureMachinePoolReconciler) reconcilePause(ctx context.Context, machinePoolScope *scope.MachinePoolScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachinePoolReconciler.reconcilePause")
	defer done()

	log.Info("Reconciling AzureMachinePool pause")

	amps, err := ampr.createAzureMachinePoolService(machinePoolScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed creating a new AzureMachinePoolService")
	}

	if err := amps.Pause(ctx); err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "error deleting AzureMachinePool %s/%s", machinePoolScope.AzureMachinePool.Namespace, machinePoolScope.Name())
	}
	infracontroller.RemoveBlockMoveAnnotation(machinePoolScope.AzureMachinePool)

	return reconcile.Result{}, nil
}

func (ampr *AzureMachinePoolReconciler) reconcileDelete(ctx context.Context, machinePoolScope *scope.MachinePoolScope, clusterScope infracontroller.ClusterScoper) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureMachinePoolReconciler.reconcileDelete")
	defer done()

	log.V(2).Info("handling deleted AzureMachinePool")

	if infracontroller.ShouldDeleteIndividualResources(ctx, clusterScope) {
		amps, err := ampr.createAzureMachinePoolService(machinePoolScope)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed creating a new AzureMachinePoolService")
		}

		log.V(4).Info("deleting AzureMachinePool resource individually")
		if err := amps.Delete(ctx); err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "error deleting AzureMachinePool %s/%s", machinePoolScope.AzureMachinePool.Namespace, machinePoolScope.Name())
		}
	}

	// Block deletion until all AzureMachinePoolMachines are finished deleting.
	ampms, err := machinePoolScope.GetMachinePoolMachines(ctx)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "error finding AzureMachinePoolMachines while deleting AzureMachinePool %s/%s", machinePoolScope.AzureMachinePool.Namespace, machinePoolScope.Name())
	}

	if len(ampms) > 0 {
		log.Info("AzureMachinePool still has dependent AzureMachinePoolMachines, deleting them first and requeing", "count", len(ampms))

		var errs []error

		for _, ampm := range ampms {
			if !ampm.GetDeletionTimestamp().IsZero() {
				// Don't handle deleted child
				continue
			}

			if err := machinePoolScope.DeleteMachine(ctx, ampm); err != nil {
				err = errors.Wrapf(err, "error deleting AzureMachinePool %s/%s: failed to delete %s %s", machinePoolScope.AzureMachinePool.Namespace, machinePoolScope.AzureMachinePool.Name, ampm.Namespace, ampm.Name)
				log.Error(err, "Error deleting AzureMachinePoolMachine", "namespace", ampm.Namespace, "name", ampm.Name)
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return ctrl.Result{}, kerrors.NewAggregate(errs)
		}

		return reconcile.Result{}, nil
	}

	// Delete succeeded, remove finalizer
	log.V(4).Info("removing finalizer for AzureMachinePool")
	controllerutil.RemoveFinalizer(machinePoolScope.AzureMachinePool, expv1.MachinePoolFinalizer)
	return reconcile.Result{}, nil
}
