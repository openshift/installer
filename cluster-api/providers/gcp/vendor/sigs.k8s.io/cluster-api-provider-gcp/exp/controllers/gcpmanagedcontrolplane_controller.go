/*
Copyright The Kubernetes Authors.

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

	"sigs.k8s.io/cluster-api/util/annotations"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/container/clusters"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/util/reconciler"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/record"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// GCPManagedControlPlaneReconciler reconciles a GCPManagedControlPlane object.
type GCPManagedControlPlaneReconciler struct {
	client.Client
	ReconcileTimeout time.Duration
	WatchFilterValue string
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedcontrolplanes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedcontrolplanes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedcontrolplanes/finalizers,verbs=update
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedclusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;patch

// SetupWithManager sets up the controller with the Manager.
func (r *GCPManagedControlPlaneReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := log.FromContext(ctx).WithValues("controller", "GCPManagedControlPlane")

	gcpManagedControlPlane := &infrav1exp.GCPManagedControlPlane{}
	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(gcpManagedControlPlane).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue)).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	if err = c.Watch(
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, gcpManagedControlPlane.GroupVersionKind(), mgr.GetClient(), &infrav1exp.GCPManagedControlPlane{})),
			predicates.ClusterPausedTransitionsOrInfrastructureReady(mgr.GetScheme(), log),
		)); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	return nil
}

func (r *GCPManagedControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultedLoopTimeout(r.ReconcileTimeout))
	defer cancel()

	log := ctrl.LoggerFrom(ctx)

	// Get the control plane instance
	gcpManagedControlPlane := &infrav1exp.GCPManagedControlPlane{}
	if err := r.Client.Get(ctx, req.NamespacedName, gcpManagedControlPlane); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Get the cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, gcpManagedControlPlane.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to retrieve owner Cluster from the API Server")
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, gcpManagedControlPlane) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	// Get the managed cluster
	managedCluster := &infrav1exp.GCPManagedCluster{}
	key := client.ObjectKey{
		Namespace: gcpManagedControlPlane.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, key, managedCluster); err != nil {
		log.Error(err, "Failed to retrieve GCPManagedCluster from the API Server")
		return ctrl.Result{}, err
	}

	managedControlPlaneScope, err := scope.NewManagedControlPlaneScope(ctx, scope.ManagedControlPlaneScopeParams{
		Client:                 r.Client,
		Cluster:                cluster,
		GCPManagedCluster:      managedCluster,
		GCPManagedControlPlane: gcpManagedControlPlane,
	})
	if err != nil {
		return ctrl.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope when exiting this function so we can persist any GCPMachine changes.
	defer func() {
		if err := managedControlPlaneScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted clusters
	if !gcpManagedControlPlane.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, managedControlPlaneScope)
	}

	// Handle non-deleted clusters
	return r.reconcile(ctx, managedControlPlaneScope)
}

func (r *GCPManagedControlPlaneReconciler) reconcile(ctx context.Context, managedControlPlaneScope *scope.ManagedControlPlaneScope) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("controller", "gcpmanagedcontrolplane")
	log.Info("Reconciling GCPManagedControlPlane")

	controllerutil.AddFinalizer(managedControlPlaneScope.GCPManagedControlPlane, infrav1exp.ManagedControlPlaneFinalizer)
	if err := managedControlPlaneScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	if !managedControlPlaneScope.GCPManagedCluster.Status.Ready {
		log.Info("GCPManagedCluster not ready yet, retry later")
		return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
	}

	reconcilers := map[string]cloud.ReconcilerWithResult{
		"container_clusters": clusters.New(managedControlPlaneScope),
	}

	for name, r := range reconcilers {
		res, err := r.Reconcile(ctx)
		if err != nil {
			log.Error(err, "Reconcile error", "reconciler", name)
			record.Warnf(managedControlPlaneScope.GCPManagedControlPlane, "GCPManagedControlPlaneReconcile", "Reconcile error - %v", err)
			return ctrl.Result{}, err
		}
		if res.RequeueAfter > 0 {
			log.V(4).Info("Reconciler requested requeueAfter", "reconciler", name, "after", res.RequeueAfter)
			return res, nil
		}
		if res.Requeue {
			log.V(4).Info("Reconciler requested requeue", "reconciler", name)
			return res, nil
		}
	}

	return ctrl.Result{}, nil
}

func (r *GCPManagedControlPlaneReconciler) reconcileDelete(ctx context.Context, managedControlPlaneScope *scope.ManagedControlPlaneScope) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("controller", "gcpmanagedcontrolplane", "action", "delete")
	log.Info("Deleting GCPManagedControlPlane")

	reconcilers := map[string]cloud.ReconcilerWithResult{
		"container_clusters": clusters.New(managedControlPlaneScope),
	}

	for name, r := range reconcilers {
		res, err := r.Delete(ctx)
		if err != nil {
			log.Error(err, "Reconcile error", "reconciler", name)
			record.Warnf(managedControlPlaneScope.GCPManagedControlPlane, "GCPManagedControlPlaneReconcile", "Reconcile error - %v", err)
			return ctrl.Result{}, err
		}
		if res.RequeueAfter > 0 {
			log.V(4).Info("Reconciler requested requeueAfter", "reconciler", name, "after", res.RequeueAfter)
			return res, nil
		}
		if res.Requeue {
			log.V(4).Info("Reconciler requested requeue", "reconciler", name)
			return res, nil
		}
	}

	if conditions.Get(managedControlPlaneScope.GCPManagedControlPlane, infrav1exp.GKEControlPlaneDeletingCondition).Reason == infrav1exp.GKEControlPlaneDeletedReason {
		controllerutil.RemoveFinalizer(managedControlPlaneScope.GCPManagedControlPlane, infrav1exp.ManagedControlPlaneFinalizer)
	}

	return ctrl.Result{RequeueAfter: reconciler.DefaultRetryTime}, nil
}
