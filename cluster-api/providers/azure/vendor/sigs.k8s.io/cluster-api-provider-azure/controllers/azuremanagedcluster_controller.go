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
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/coalescing"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// AzureManagedClusterReconciler reconciles an AzureManagedCluster object.
type AzureManagedClusterReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	ReconcileTimeout time.Duration
	WatchFilterValue string
}

// SetupWithManager initializes this controller with a manager.
func (amcr *AzureManagedClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options Options) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureManagedClusterReconciler.SetupWithManager",
		tele.KVP("controller", "AzureManagedCluster"),
	)
	defer done()

	var r reconcile.Reconciler = amcr
	if options.Cache != nil {
		r = coalescing.NewReconciler(amcr, options.Cache, log)
	}

	azManagedCluster := &infrav1.AzureManagedCluster{}
	// create mapper to transform incoming AzureManagedControlPlanes into AzureManagedCluster requests
	azureManagedControlPlaneMapper, err := AzureManagedControlPlaneToAzureManagedClusterMapper(ctx, amcr.Client, log)
	if err != nil {
		return errors.Wrap(err, "failed to create AzureManagedControlPlane to AzureManagedClusters mapper")
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options.Options).
		For(azManagedCluster).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, amcr.WatchFilterValue)).
		// watch AzureManagedControlPlane resources
		Watches(
			&infrav1.AzureManagedControlPlane{},
			handler.EnqueueRequestsFromMapFunc(azureManagedControlPlaneMapper),
		).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	// Add a watch on clusterv1.Cluster object for unpause notifications.
	if err = c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind("AzureManagedCluster"), mgr.GetClient(), &infrav1.AzureManagedCluster{})),
		predicates.ClusterUnpaused(log),
		predicates.ResourceNotPausedAndHasFilterLabel(log, amcr.WatchFilterValue),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready clusters")
	}
	return nil
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile idempotently gets, creates, and updates a managed cluster.
func (amcr *AzureManagedClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultedLoopTimeout(amcr.ReconcileTimeout))
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(
		ctx,
		"controllers.AzureManagedClusterReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "AzureManagedCluster"),
	)
	defer done()

	// Fetch the AzureManagedCluster instance
	aksCluster := &infrav1.AzureManagedCluster{}
	err := amcr.Get(ctx, req.NamespacedName, aksCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, amcr.Client, aksCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	controlPlane := &infrav1.AzureManagedControlPlane{}
	controlPlaneRef := types.NamespacedName{
		Name:      cluster.Spec.ControlPlaneRef.Name,
		Namespace: cluster.Namespace,
	}

	log = log.WithValues("cluster", cluster.Name)

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, aksCluster) {
		log.Info("AzureManagedCluster or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	if err := amcr.Get(ctx, controlPlaneRef, controlPlane); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to get control plane ref")
	}

	log = log.WithValues("controlPlane", controlPlaneRef.Name)

	patchhelper, err := patch.NewHelper(aksCluster, amcr.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to init patch helper")
	}

	// Infrastructure must be ready before control plane. We should also enqueue
	// requests from control plane to infra cluster to keep control plane endpoint accurate.
	aksCluster.Status.Ready = true
	aksCluster.Spec.ControlPlaneEndpoint = controlPlane.Spec.ControlPlaneEndpoint

	if err := patchhelper.Patch(ctx, aksCluster); err != nil {
		return reconcile.Result{}, err
	}

	log.Info("Successfully reconciled")

	return reconcile.Result{}, nil
}
