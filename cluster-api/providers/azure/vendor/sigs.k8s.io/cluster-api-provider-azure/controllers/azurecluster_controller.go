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
	"time"

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

// AzureClusterReconciler reconciles an AzureCluster object.
type AzureClusterReconciler struct {
	client.Client
	Recorder                  record.EventRecorder
	ReconcileTimeout          time.Duration
	WatchFilterValue          string
	createAzureClusterService azureClusterServiceCreator
}

type azureClusterServiceCreator func(clusterScope *scope.ClusterScope) (*azureClusterService, error)

// NewAzureClusterReconciler returns a new AzureClusterReconciler instance.
func NewAzureClusterReconciler(client client.Client, recorder record.EventRecorder, reconcileTimeout time.Duration, watchFilterValue string) *AzureClusterReconciler {
	acr := &AzureClusterReconciler{
		Client:           client,
		Recorder:         recorder,
		ReconcileTimeout: reconcileTimeout,
		WatchFilterValue: watchFilterValue,
	}

	acr.createAzureClusterService = newAzureClusterService

	return acr
}

// SetupWithManager initializes this controller with a manager.
func (acr *AzureClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options Options) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureClusterReconciler.SetupWithManager",
		tele.KVP("controller", "AzureCluster"),
	)
	defer done()

	var r reconcile.Reconciler = acr
	if options.Cache != nil {
		r = coalescing.NewReconciler(acr, options.Cache, log)
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options.Options).
		For(&infrav1.AzureCluster{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(log, acr.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	// Add a watch on clusterv1.Cluster object for pause/unpause notifications.
	if err = c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind("AzureCluster"), mgr.GetClient(), &infrav1.AzureCluster{})),
		ClusterUpdatePauseChange(log),
		predicates.ResourceHasFilterLabel(log, acr.WatchFilterValue),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready clusters")
	}

	return nil
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremachinetemplates;azuremachinetemplates/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureclusteridentities;azureclusteridentities/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=list;
// +kubebuilder:rbac:groups=resources.azure.com,resources=resourcegroups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=resources.azure.com,resources=resourcegroups/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=network.azure.com,resources=natgateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.azure.com,resources=natgateways/status,verbs=get;list;watch

// Reconcile idempotently gets, creates, and updates a cluster.
func (acr *AzureClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultedLoopTimeout(acr.ReconcileTimeout))
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(
		ctx,
		"controllers.AzureClusterReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "AzureCluster"),
	)
	defer done()

	// Fetch the AzureCluster instance
	azureCluster := &infrav1.AzureCluster{}
	err := acr.Get(ctx, req.NamespacedName, azureCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			acr.Recorder.Eventf(azureCluster, corev1.EventTypeNormal, "AzureClusterObjectNotFound", err.Error())
			log.Info("object was not found")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, acr.Client, azureCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		acr.Recorder.Eventf(azureCluster, corev1.EventTypeNormal, "OwnerRefNotSet", "Cluster Controller has not yet set OwnerRef")
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Create the scope.
	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:       acr.Client,
		Cluster:      cluster,
		AzureCluster: azureCluster,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to create scope")
		acr.Recorder.Eventf(azureCluster, corev1.EventTypeWarning, "CreateClusterScopeFailed", err.Error())
		return reconcile.Result{}, err
	}

	// Always close the scope when exiting this function so we can persist any AzureMachine changes.
	defer func() {
		if err := clusterScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, azureCluster) {
		acr.Recorder.Eventf(azureCluster, corev1.EventTypeNormal, "ClusterPaused", "AzureCluster or linked Cluster is marked as paused. Won't reconcile normally")
		log.Info("AzureCluster or linked Cluster is marked as paused. Won't reconcile normally")
		return acr.reconcilePause(ctx, clusterScope)
	}

	if azureCluster.Spec.IdentityRef != nil {
		err := EnsureClusterIdentity(ctx, acr.Client, azureCluster, azureCluster.Spec.IdentityRef, infrav1.ClusterFinalizer)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		log.Info(fmt.Sprintf("WARNING, %s", deprecatedManagerCredsWarning))
		acr.Recorder.Eventf(azureCluster, corev1.EventTypeWarning, "AzureClusterIdentity", deprecatedManagerCredsWarning)
	}

	// Handle deleted clusters
	if !azureCluster.DeletionTimestamp.IsZero() {
		return acr.reconcileDelete(ctx, clusterScope)
	}

	// Handle non-deleted clusters
	return acr.reconcileNormal(ctx, clusterScope)
}

func (acr *AzureClusterReconciler) reconcileNormal(ctx context.Context, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureClusterReconciler.reconcileNormal")
	defer done()

	log.Info("Reconciling AzureCluster")
	azureCluster := clusterScope.AzureCluster

	// If the AzureCluster doesn't have our finalizer, add it.
	if controllerutil.AddFinalizer(azureCluster, infrav1.ClusterFinalizer) {
		// Register the finalizer immediately to avoid orphaning Azure resources on delete
		if err := clusterScope.PatchObject(ctx); err != nil {
			return reconcile.Result{}, err
		}
	}

	acs, err := acr.createAzureClusterService(clusterScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create a new AzureClusterReconciler")
	}

	if err := acs.Reconcile(ctx); err != nil {
		// Handle terminal & transient errors
		var reconcileError azure.ReconcileError
		if errors.As(err, &reconcileError) {
			if reconcileError.IsTerminal() {
				acr.Recorder.Eventf(clusterScope.AzureCluster, corev1.EventTypeWarning, "ReconcileError", errors.Wrapf(err, "failed to reconcile AzureCluster").Error())
				log.Error(err, "failed to reconcile AzureCluster", "name", clusterScope.ClusterName())
				conditions.MarkFalse(azureCluster, infrav1.NetworkInfrastructureReadyCondition, infrav1.FailedReason, clusterv1.ConditionSeverityError, "")
				return reconcile.Result{}, nil
			}
			if reconcileError.IsTransient() {
				if azure.IsOperationNotDoneError(reconcileError) {
					log.V(2).Info(fmt.Sprintf("AzureCluster reconcile not done: %s", reconcileError.Error()))
				} else {
					log.V(2).Info(fmt.Sprintf("transient failure to reconcile AzureCluster, retrying: %s", reconcileError.Error()))
				}
				return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
			}
		}

		wrappedErr := errors.Wrap(err, "failed to reconcile cluster services")
		acr.Recorder.Eventf(azureCluster, corev1.EventTypeWarning, "ClusterReconcilerNormalFailed", wrappedErr.Error())
		conditions.MarkFalse(azureCluster, infrav1.NetworkInfrastructureReadyCondition, infrav1.FailedReason, clusterv1.ConditionSeverityError, wrappedErr.Error())
		return reconcile.Result{}, wrappedErr
	}

	// Set APIEndpoints so the Cluster API Cluster Controller can pull them
	if azureCluster.Spec.ControlPlaneEndpoint.Host == "" {
		azureCluster.Spec.ControlPlaneEndpoint.Host = clusterScope.APIServerHost()
	}
	if azureCluster.Spec.ControlPlaneEndpoint.Port == 0 {
		azureCluster.Spec.ControlPlaneEndpoint.Port = clusterScope.APIServerPort()
	}

	// No errors, so mark us ready so the Cluster API Cluster Controller can pull it
	azureCluster.Status.Ready = true
	conditions.MarkTrue(azureCluster, infrav1.NetworkInfrastructureReadyCondition)

	return reconcile.Result{}, nil
}

func (acr *AzureClusterReconciler) reconcilePause(ctx context.Context, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureClusterReconciler.reconcilePause")
	defer done()

	log.Info("Reconciling AzureCluster pause")

	acs, err := acr.createAzureClusterService(clusterScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create a new azureClusterService")
	}

	if err := acs.Pause(ctx); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to pause cluster services")
	}

	return reconcile.Result{}, nil
}

func (acr *AzureClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureClusterReconciler.reconcileDelete")
	defer done()

	log.Info("Reconciling AzureCluster delete")

	azureCluster := clusterScope.AzureCluster

	acs, err := acr.createAzureClusterService(clusterScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create a new AzureClusterReconciler")
	}

	if err := acs.Delete(ctx); err != nil {
		// Handle transient errors
		var reconcileError azure.ReconcileError
		if errors.As(err, &reconcileError) {
			if reconcileError.IsTransient() {
				if azure.IsOperationNotDoneError(reconcileError) {
					log.V(2).Info(fmt.Sprintf("AzureCluster delete not done: %s", reconcileError.Error()))
				} else {
					log.V(2).Info("transient failure to delete AzureCluster, retrying")
				}
				return reconcile.Result{RequeueAfter: reconcileError.RequeueAfter()}, nil
			}
		}

		wrappedErr := errors.Wrapf(err, "error deleting AzureCluster %s/%s", azureCluster.Namespace, azureCluster.Name)
		acr.Recorder.Eventf(azureCluster, corev1.EventTypeWarning, "ClusterReconcilerDeleteFailed", wrappedErr.Error())
		conditions.MarkFalse(azureCluster, infrav1.NetworkInfrastructureReadyCondition, clusterv1.DeletionFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		return reconcile.Result{}, wrappedErr
	}

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(azureCluster, infrav1.ClusterFinalizer)

	if azureCluster.Spec.IdentityRef != nil {
		// Cluster is deleted so remove the identity finalizer.
		err := RemoveClusterIdentityFinalizer(ctx, acr.Client, azureCluster, azureCluster.Spec.IdentityRef, infrav1.ClusterFinalizer)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
