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

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	capiexputil "sigs.k8s.io/cluster-api/exp/util"
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
	"sigs.k8s.io/cluster-api-provider-azure/pkg/coalescing"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AzureManagedControlPlaneReconciler reconciles an AzureManagedControlPlane object.
type AzureManagedControlPlaneReconciler struct {
	client.Client
	Recorder                                 record.EventRecorder
	Timeouts                                 reconciler.Timeouts
	WatchFilterValue                         string
	CredentialCache                          azure.CredentialCache
	getNewAzureManagedControlPlaneReconciler func(scope *scope.ManagedControlPlaneScope) (*azureManagedControlPlaneService, error)
}

// SetupWithManager initializes this controller with a manager.
func (amcpr *AzureManagedControlPlaneReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options Options) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureManagedControlPlaneReconciler.SetupWithManager",
		tele.KVP("controller", "AzureManagedControlPlane"),
	)
	defer done()

	amcpr.getNewAzureManagedControlPlaneReconciler = newAzureManagedControlPlaneReconciler
	var r reconcile.Reconciler = amcpr
	if options.Cache != nil {
		r = coalescing.NewReconciler(amcpr, options.Cache, log)
	}

	azManagedControlPlane := &infrav1.AzureManagedControlPlane{}
	// create mapper to transform incoming AzureManagedClusters into AzureManagedControlPlane requests
	azureManagedClusterMapper, err := AzureManagedClusterToAzureManagedControlPlaneMapper(ctx, amcpr.Client, log)
	if err != nil {
		return errors.Wrap(err, "failed to create AzureManagedCluster to AzureManagedControlPlane mapper")
	}

	// map requests for machine pools corresponding to AzureManagedControlPlane's defaultPool back to the corresponding AzureManagedControlPlane.
	azureManagedMachinePoolMapper := MachinePoolToAzureManagedControlPlaneMapFunc(ctx, amcpr.Client, infrav1.GroupVersion.WithKind(infrav1.AzureManagedControlPlaneKind), log)

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options.Options).
		For(azManagedControlPlane).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, amcpr.WatchFilterValue)).
		// watch AzureManagedCluster resources
		Watches(
			&infrav1.AzureManagedCluster{},
			handler.EnqueueRequestsFromMapFunc(azureManagedClusterMapper),
		).
		// watch MachinePool resources
		Watches(
			&expv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(azureManagedMachinePoolMapper),
		).
		// Add a watch on clusterv1.Cluster object for pause/unpause & ready notifications.
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(amcpr.ClusterToAzureManagedControlPlane),
			builder.WithPredicates(
				ClusterPauseChangeAndInfrastructureReady(mgr.GetScheme(), log),
				predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, amcpr.WatchFilterValue),
			),
		).
		Complete(r)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=resources.azure.com,resources=resourcegroups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=resources.azure.com,resources=resourcegroups/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=containerservice.azure.com,resources=managedclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=containerservice.azure.com,resources=managedclusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=network.azure.com,resources=privateendpoints;virtualnetworks;virtualnetworkssubnets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.azure.com,resources=privateendpoints/status;virtualnetworks/status;virtualnetworkssubnets/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=containerservice.azure.com,resources=fleetsmembers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=containerservice.azure.com,resources=fleetsmembers/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=kubernetesconfiguration.azure.com,resources=extensions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubernetesconfiguration.azure.com,resources=extensions/status,verbs=get;list;watch

// Reconcile idempotently gets, creates, and updates a managed control plane.
func (amcpr *AzureManagedControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, amcpr.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureManagedControlPlaneReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", infrav1.AzureManagedControlPlaneKind),
	)
	defer done()

	// Fetch the AzureManagedControlPlane instance
	azureControlPlane := &infrav1.AzureManagedControlPlane{}
	err := amcpr.Get(ctx, req.NamespacedName, azureControlPlane)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, amcpr.Client, azureControlPlane.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Fetch all the ManagedMachinePools owned by this Cluster.
	opt1 := client.InNamespace(azureControlPlane.Namespace)
	opt2 := client.MatchingLabels(map[string]string{
		clusterv1.ClusterNameLabel: cluster.Name,
	})

	ammpList := &infrav1.AzureManagedMachinePoolList{}
	if err := amcpr.List(ctx, ammpList, opt1, opt2); err != nil {
		return reconcile.Result{}, err
	}

	var pools = make([]scope.ManagedMachinePool, len(ammpList.Items))

	for i, ammp := range ammpList.Items {
		// Fetch the owner MachinePool.
		ownerPool, err := capiexputil.GetOwnerMachinePool(ctx, amcpr.Client, ammp.ObjectMeta)
		if err != nil || ownerPool == nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to fetch owner MachinePool for AzureManagedMachinePool: %s", ammp.Name)
		}
		pools[i] = scope.ManagedMachinePool{
			InfraMachinePool: &ammpList.Items[i],
			MachinePool:      ownerPool,
		}
	}

	// Create the scope.
	mcpScope, err := scope.NewManagedControlPlaneScope(ctx, scope.ManagedControlPlaneScopeParams{
		Client:              amcpr.Client,
		Cluster:             cluster,
		ControlPlane:        azureControlPlane,
		ManagedMachinePools: pools,
		Timeouts:            amcpr.Timeouts,
		CredentialCache:     amcpr.CredentialCache,
	})
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create scope")
	}

	// Always patch when exiting so we can persist changes to finalizers and status
	defer func() {
		if err := mcpScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, azureControlPlane) {
		log.Info("AzureManagedControlPlane or linked Cluster is marked as paused. Won't reconcile normally")
		return amcpr.reconcilePause(ctx, mcpScope)
	}

	// check if the control plane's namespace is allowed for this identity and update owner references for the identity.
	if azureControlPlane.Spec.IdentityRef != nil {
		err := EnsureClusterIdentity(ctx, amcpr.Client, azureControlPlane, azureControlPlane.Spec.IdentityRef, infrav1.ManagedClusterFinalizer)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		warningMessage := "You're using deprecated functionality: " +
			"Using Azure credentials from the manager environment is deprecated and will be removed in future releases. " +
			"Please specify an AzureClusterIdentity for the AzureManagedControlPlane instead, see: https://capz.sigs.k8s.io/topics/multitenancy.html "
		log.Info(fmt.Sprintf("WARNING, %s", warningMessage))
		amcpr.Recorder.Eventf(azureControlPlane, corev1.EventTypeWarning, "AzureClusterIdentity", warningMessage)
	}

	// Handle deleted clusters
	if !azureControlPlane.DeletionTimestamp.IsZero() {
		return amcpr.reconcileDelete(ctx, mcpScope)
	}
	// Handle non-deleted clusters
	return amcpr.reconcileNormal(ctx, mcpScope)
}

func (amcpr *AzureManagedControlPlaneReconciler) reconcileNormal(ctx context.Context, scope *scope.ManagedControlPlaneScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureManagedControlPlaneReconciler.reconcileNormal")
	defer done()

	log.Info("Reconciling AzureManagedControlPlane")

	// Remove deprecated Cluster finalizer if it exists
	needsPatch := controllerutil.RemoveFinalizer(scope.ControlPlane, infrav1.ClusterFinalizer)
	// Register our finalizer immediately to avoid orphaning Azure resources on delete
	needsPatch = controllerutil.AddFinalizer(scope.ControlPlane, infrav1.ManagedClusterFinalizer) || needsPatch
	// Register the block-move annotation immediately to avoid moving un-paused ASO resources
	needsPatch = AddBlockMoveAnnotation(scope.ControlPlane) || needsPatch
	if needsPatch {
		if err := scope.PatchObject(ctx); err != nil {
			amcpr.Recorder.Eventf(scope.ControlPlane, corev1.EventTypeWarning, "AzureManagedControlPlane unavailable", "failed to patch resource: %s", err)
			return reconcile.Result{}, err
		}
	}

	svc, err := amcpr.getNewAzureManagedControlPlaneReconciler(scope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create azureManagedControlPlane service")
	}
	if err := svc.Reconcile(ctx); err != nil {
		// Handle transient and terminal errors
		log := log.WithValues("name", scope.ControlPlane.Name, "namespace", scope.ControlPlane.Namespace)
		var reconcileError azure.ReconcileError
		if errors.As(err, &reconcileError) {
			if reconcileError.IsTerminal() {
				log.Error(err, "failed to reconcile AzureManagedControlPlane")
				return reconcile.Result{}, nil
			}

			if reconcileError.IsTransient() {
				log.V(4).Info("requeuing due to transient failure", "error", err)
				return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
			}

			return reconcile.Result{}, errors.Wrap(err, "failed to reconcile AzureManagedControlPlane")
		}

		return reconcile.Result{}, errors.Wrapf(err, "error creating AzureManagedControlPlane %s/%s", scope.ControlPlane.Namespace, scope.ControlPlane.Name)
	}

	// No errors, so mark us ready so the Cluster API Cluster Controller can pull it
	scope.ControlPlane.Status.Ready = true
	scope.ControlPlane.Status.Initialized = true
	scope.ControlPlane.Status.Version = scope.ControlPlane.Spec.Version

	log.Info("Successfully reconciled")

	return reconcile.Result{}, nil
}

func (amcpr *AzureManagedControlPlaneReconciler) reconcilePause(ctx context.Context, scope *scope.ManagedControlPlaneScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureManagedControlPlane.reconcilePause")
	defer done()

	log.Info("Reconciling AzureManagedControlPlane pause")

	svc, err := amcpr.getNewAzureManagedControlPlaneReconciler(scope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create azureManagedControlPlane service")
	}
	if err := svc.Pause(ctx); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to pause control plane services")
	}
	RemoveBlockMoveAnnotation(scope.ControlPlane)

	return reconcile.Result{}, nil
}

func (amcpr *AzureManagedControlPlaneReconciler) reconcileDelete(ctx context.Context, scope *scope.ManagedControlPlaneScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureManagedControlPlaneReconciler.reconcileDelete")
	defer done()

	log.Info("Reconciling AzureManagedControlPlane delete")

	svc, err := amcpr.getNewAzureManagedControlPlaneReconciler(scope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create azureManagedControlPlane service")
	}
	if err := svc.Delete(ctx); err != nil {
		// Handle transient errors
		var reconcileError azure.ReconcileError
		if errors.As(err, &reconcileError) && reconcileError.IsTransient() {
			if azure.IsOperationNotDoneError(reconcileError) {
				log.V(2).Info(fmt.Sprintf("AzureManagedControlPlane delete not done: %s", reconcileError.Error()))
			} else {
				log.V(2).Info("transient failure to delete AzureManagedControlPlane, retrying")
			}
			return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
		}
		return reconcile.Result{}, errors.Wrapf(err, "error deleting AzureManagedControlPlane %s/%s", scope.ControlPlane.Namespace, scope.ControlPlane.Name)
	}

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(scope.ControlPlane, infrav1.ManagedClusterFinalizer)

	if scope.ControlPlane.Spec.IdentityRef != nil {
		err := RemoveClusterIdentityFinalizer(ctx, amcpr.Client, scope.ControlPlane, scope.ControlPlane.Spec.IdentityRef, infrav1.ManagedClusterFinalizer)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// ClusterToAzureManagedControlPlane is a handler.ToRequestsFunc to be used to enqueue requests for
// reconciliation for AzureManagedControlPlane based on updates to a Cluster.
func (amcpr *AzureManagedControlPlaneReconciler) ClusterToAzureManagedControlPlane(_ context.Context, o client.Object) []ctrl.Request {
	c, ok := o.(*clusterv1.Cluster)
	if !ok {
		panic(fmt.Sprintf("Expected a Cluster but got a %T", o))
	}

	controlPlaneRef := c.Spec.ControlPlaneRef
	if controlPlaneRef != nil && controlPlaneRef.Kind == infrav1.AzureManagedControlPlaneKind {
		return []ctrl.Request{{NamespacedName: client.ObjectKey{Namespace: controlPlaneRef.Namespace, Name: controlPlaneRef.Name}}}
	}

	return nil
}
