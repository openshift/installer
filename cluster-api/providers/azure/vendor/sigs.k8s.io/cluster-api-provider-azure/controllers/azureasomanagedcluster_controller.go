/*
Copyright 2024 The Kubernetes Authors.

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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	infrav1alpha "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/mutators"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var errInvalidControlPlaneKind = errors.New("AzureASOManagedCluster cannot be used without AzureASOManagedControlPlane")

// AzureASOManagedClusterReconciler reconciles a AzureASOManagedCluster object.
type AzureASOManagedClusterReconciler struct {
	client.Client
	WatchFilterValue string

	newResourceReconciler func(*infrav1alpha.AzureASOManagedCluster, []*unstructured.Unstructured) resourceReconciler
}

type resourceReconciler interface {
	// Reconcile reconciles resources defined by this object and updates this object's status to reflect the
	// state of the specified resources.
	Reconcile(context.Context) error

	// Pause stops ASO from continuously reconciling the specified resources.
	Pause(context.Context) error

	// Delete begins deleting the specified resources and updates the object's status to reflect the state of
	// the specified resources.
	Delete(context.Context) error
}

// SetupWithManager sets up the controller with the Manager.
func (r *AzureASOManagedClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedClusterReconciler.SetupWithManager",
		tele.KVP("controller", infrav1alpha.AzureASOManagedClusterKind),
	)
	defer done()

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1alpha.AzureASOManagedCluster{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(log, r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		// Watch clusters for pause/unpause notifications
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(
				util.ClusterToInfrastructureMapFunc(ctx, infrav1alpha.GroupVersion.WithKind(infrav1alpha.AzureASOManagedClusterKind), mgr.GetClient(), &infrav1alpha.AzureASOManagedCluster{}),
			),
			builder.WithPredicates(
				predicates.ResourceHasFilterLabel(log, r.WatchFilterValue),
				ClusterUpdatePauseChange(log),
			),
		).
		Watches(
			&infrav1alpha.AzureASOManagedControlPlane{},
			handler.EnqueueRequestsFromMapFunc(asoManagedControlPlaneToManagedClusterMap(r.Client)),
			builder.WithPredicates(
				predicates.ResourceHasFilterLabel(log, r.WatchFilterValue),
				predicate.Funcs{
					CreateFunc: func(ev event.CreateEvent) bool {
						controlPlane := ev.Object.(*infrav1alpha.AzureASOManagedControlPlane)
						return !controlPlane.Status.ControlPlaneEndpoint.IsZero()
					},
					UpdateFunc: func(ev event.UpdateEvent) bool {
						oldControlPlane := ev.ObjectOld.(*infrav1alpha.AzureASOManagedControlPlane)
						newControlPlane := ev.ObjectNew.(*infrav1alpha.AzureASOManagedControlPlane)
						return oldControlPlane.Status.ControlPlaneEndpoint !=
							newControlPlane.Status.ControlPlaneEndpoint
					},
				},
			),
		).
		Build(r)
	if err != nil {
		return err
	}

	externalTracker := &external.ObjectTracker{
		Cache:      mgr.GetCache(),
		Controller: c,
	}

	r.newResourceReconciler = func(asoManagedCluster *infrav1alpha.AzureASOManagedCluster, resources []*unstructured.Unstructured) resourceReconciler {
		return &ResourceReconciler{
			Client:    r.Client,
			resources: resources,
			owner:     asoManagedCluster,
			watcher:   externalTracker,
		}
	}

	return nil
}

func asoManagedControlPlaneToManagedClusterMap(c client.Client) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		asoManagedControlPlane := o.(*infrav1alpha.AzureASOManagedControlPlane)

		cluster, err := util.GetOwnerCluster(ctx, c, asoManagedControlPlane.ObjectMeta)
		if err != nil {
			return nil
		}

		if cluster == nil ||
			cluster.Spec.InfrastructureRef == nil ||
			cluster.Spec.InfrastructureRef.APIVersion != infrav1alpha.GroupVersion.Identifier() ||
			cluster.Spec.InfrastructureRef.Kind != infrav1alpha.AzureASOManagedClusterKind {
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: client.ObjectKey{
					Namespace: cluster.Spec.InfrastructureRef.Namespace,
					Name:      cluster.Spec.InfrastructureRef.Name,
				},
			},
		}
	}
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedclusters/finalizers,verbs=update

// Reconcile reconciles an AzureASOManagedCluster.
func (r *AzureASOManagedClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, resultErr error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedClusterReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", infrav1alpha.AzureASOManagedClusterKind),
	)
	defer done()

	asoManagedCluster := &infrav1alpha.AzureASOManagedCluster{}
	err := r.Get(ctx, req.NamespacedName, asoManagedCluster)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	patchHelper, err := patch.NewHelper(asoManagedCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create patch helper: %w", err)
	}
	defer func() {
		err := patchHelper.Patch(ctx, asoManagedCluster)
		if err != nil && resultErr == nil {
			resultErr = err
			result = ctrl.Result{}
		}
	}()

	asoManagedCluster.Status.Ready = false

	cluster, err := util.GetOwnerCluster(ctx, r.Client, asoManagedCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}

	if cluster != nil && cluster.Spec.Paused ||
		annotations.HasPaused(asoManagedCluster) {
		return r.reconcilePaused(ctx, asoManagedCluster)
	}

	if !asoManagedCluster.GetDeletionTimestamp().IsZero() {
		return r.reconcileDelete(ctx, asoManagedCluster)
	}

	return r.reconcileNormal(ctx, asoManagedCluster, cluster)
}

func (r *AzureASOManagedClusterReconciler) reconcileNormal(ctx context.Context, asoManagedCluster *infrav1alpha.AzureASOManagedCluster, cluster *clusterv1.Cluster) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedClusterReconciler.reconcileNormal",
	)
	defer done()
	log.V(4).Info("reconciling normally")

	if cluster == nil {
		log.V(4).Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}
	if cluster.Spec.ControlPlaneRef == nil ||
		cluster.Spec.ControlPlaneRef.APIVersion != infrav1alpha.GroupVersion.Identifier() ||
		cluster.Spec.ControlPlaneRef.Kind != infrav1alpha.AzureASOManagedControlPlaneKind {
		return ctrl.Result{}, reconcile.TerminalError(errInvalidControlPlaneKind)
	}

	needsPatch := controllerutil.AddFinalizer(asoManagedCluster, clusterv1.ClusterFinalizer)
	needsPatch = AddBlockMoveAnnotation(asoManagedCluster) || needsPatch
	if needsPatch {
		return ctrl.Result{Requeue: true}, nil
	}

	resources, err := mutators.ToUnstructured(ctx, asoManagedCluster.Spec.Resources)
	if err != nil {
		return ctrl.Result{}, err
	}
	resourceReconciler := r.newResourceReconciler(asoManagedCluster, resources)
	err = resourceReconciler.Reconcile(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile resources: %w", err)
	}
	for _, status := range asoManagedCluster.Status.Resources {
		if !status.Ready {
			return ctrl.Result{}, nil
		}
	}

	asoManagedControlPlane := &infrav1alpha.AzureASOManagedControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Spec.ControlPlaneRef.Namespace,
			Name:      cluster.Spec.ControlPlaneRef.Name,
		},
	}
	err = r.Get(ctx, client.ObjectKeyFromObject(asoManagedControlPlane), asoManagedControlPlane)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get AzureASOManagedControlPlane %s/%s: %w", asoManagedControlPlane.Namespace, asoManagedControlPlane.Name, err)
	}
	asoManagedCluster.Spec.ControlPlaneEndpoint = asoManagedControlPlane.Status.ControlPlaneEndpoint

	asoManagedCluster.Status.Ready = !asoManagedCluster.Spec.ControlPlaneEndpoint.IsZero()

	return ctrl.Result{}, nil
}

//nolint:unparam // an empty ctrl.Result is always returned here, leaving it as-is to avoid churn in refactoring later if that changes.
func (r *AzureASOManagedClusterReconciler) reconcilePaused(ctx context.Context, asoManagedCluster *infrav1alpha.AzureASOManagedCluster) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureASOManagedClusterReconciler.reconcilePaused")
	defer done()
	log.V(4).Info("reconciling pause")

	resources, err := mutators.ToUnstructured(ctx, asoManagedCluster.Spec.Resources)
	if err != nil {
		return ctrl.Result{}, err
	}
	resourceReconciler := r.newResourceReconciler(asoManagedCluster, resources)
	err = resourceReconciler.Pause(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to pause resources: %w", err)
	}

	RemoveBlockMoveAnnotation(asoManagedCluster)

	return ctrl.Result{}, nil
}

//nolint:unparam // an empty ctrl.Result is always returned here, leaving it as-is to avoid churn in refactoring later if that changes.
func (r *AzureASOManagedClusterReconciler) reconcileDelete(ctx context.Context, asoManagedCluster *infrav1alpha.AzureASOManagedCluster) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedClusterReconciler.reconcileDelete",
	)
	defer done()
	log.V(4).Info("reconciling delete")

	resources, err := mutators.ToUnstructured(ctx, asoManagedCluster.Spec.Resources)
	if err != nil {
		return ctrl.Result{}, err
	}
	resourceReconciler := r.newResourceReconciler(asoManagedCluster, resources)
	err = resourceReconciler.Delete(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile resources: %w", err)
	}
	if len(asoManagedCluster.Status.Resources) > 0 {
		return ctrl.Result{}, nil
	}

	controllerutil.RemoveFinalizer(asoManagedCluster, clusterv1.ClusterFinalizer)
	return ctrl.Result{}, nil
}
