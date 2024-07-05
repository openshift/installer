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
	"time"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/compute/firewalls"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/compute/loadbalancers"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/compute/networks"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/compute/subnets"
	"sigs.k8s.io/cluster-api-provider-gcp/util/reconciler"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// GCPClusterReconciler reconciles a GCPCluster object.
type GCPClusterReconciler struct {
	client.Client
	ReconcileTimeout time.Duration
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=gcpclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

func (r *GCPClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := log.FromContext(ctx).WithValues("controller", "GCPCluster")

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.GCPCluster{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		Build(r)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	clusterToInfraFn := util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind("GCPCluster"), mgr.GetClient(), &infrav1.GCPCluster{})
	if err = c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(func(mapCtx context.Context, o client.Object) []reconcile.Request {
			requests := clusterToInfraFn(mapCtx, o)
			if requests == nil {
				return nil
			}

			gcpCluster := &infrav1.GCPCluster{}
			if err := r.Get(ctx, requests[0].NamespacedName, gcpCluster); err != nil {
				log.V(4).Error(err, "Failed to get GCP cluster")
				return nil
			}

			if annotations.IsExternallyManaged(gcpCluster) {
				log.V(4).Info("GCPCluster is externally managed, skipping mapping.")
				return nil
			}
			return requests
		}),
		predicates.ClusterUnpaused(log),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready clusters")
	}

	return nil
}

func (r *GCPClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultedLoopTimeout(r.ReconcileTimeout))
	defer cancel()

	log := log.FromContext(ctx)
	gcpCluster := &infrav1.GCPCluster{}
	err := r.Get(ctx, req.NamespacedName, gcpCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("GCPCluster resource not found or already deleted")
			return ctrl.Result{}, nil
		}

		log.Error(err, "Unable to fetch GCPCluster resource")
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, gcpCluster.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to get owner cluster")
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, gcpCluster) {
		log.Info("GCPCluster of linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:     r.Client,
		Cluster:    cluster,
		GCPCluster: gcpCluster,
	})
	if err != nil {
		return ctrl.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope when exiting this function so we can persist any GCPMachine changes.
	defer func() {
		if err := clusterScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted clusters
	if !gcpCluster.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, clusterScope)
	}

	// Handle non-deleted clusters
	return r.reconcile(ctx, clusterScope)
}

func (r *GCPClusterReconciler) reconcile(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling GCPCluster")

	controllerutil.AddFinalizer(clusterScope.GCPCluster, infrav1.ClusterFinalizer)
	if err := clusterScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	region, err := clusterScope.Cloud().Regions().Get(ctx, meta.GlobalKey(clusterScope.Region()))
	if err != nil {
		return ctrl.Result{}, err
	}

	zones, err := clusterScope.Cloud().Zones().List(ctx, filter.Regexp("region", region.SelfLink))
	if err != nil {
		return ctrl.Result{}, err
	}

	failureDomains := make(clusterv1.FailureDomains, len(zones))
	for _, zone := range zones {
		if len(clusterScope.GCPCluster.Spec.FailureDomains) > 0 {
			for _, fd := range clusterScope.GCPCluster.Spec.FailureDomains {
				if fd == zone.Name {
					failureDomains[zone.Name] = clusterv1.FailureDomainSpec{
						ControlPlane: true,
					}
				}
			}
		} else {
			failureDomains[zone.Name] = clusterv1.FailureDomainSpec{
				ControlPlane: true,
			}
		}
	}

	clusterScope.SetFailureDomains(failureDomains)

	reconcilers := []cloud.Reconciler{
		networks.New(clusterScope),
		firewalls.New(clusterScope),
		// Reconcile subnets before loadbalancers since subnet is needed for internal LB
		subnets.New(clusterScope),
		loadbalancers.New(clusterScope),
	}

	for _, r := range reconcilers {
		if err := r.Reconcile(ctx); err != nil {
			log.Error(err, "Reconcile error")
			record.Warnf(clusterScope.GCPCluster, "GCPClusterReconcile", "Reconcile error - %v", err)
			return ctrl.Result{}, err
		}
	}

	controlPlaneEndpoint := clusterScope.ControlPlaneEndpoint()
	if controlPlaneEndpoint.Host == "" {
		log.Info("GCPCluster does not have control-plane endpoint yet. Reconciling")
		record.Event(clusterScope.GCPCluster, "GCPClusterReconcile", "Waiting for control-plane endpoint")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	record.Eventf(clusterScope.GCPCluster, "GCPClusterReconcile", "Got control-plane endpoint - %s", controlPlaneEndpoint.Host)
	clusterScope.SetReady()
	record.Event(clusterScope.GCPCluster, "GCPClusterReconcile", "Reconciled")
	return ctrl.Result{}, nil
}

func (r *GCPClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.ClusterScope) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling Delete GCPCluster")

	reconcilers := []cloud.Reconciler{
		loadbalancers.New(clusterScope),
		subnets.New(clusterScope),
		firewalls.New(clusterScope),
		networks.New(clusterScope),
	}

	for _, r := range reconcilers {
		if err := r.Delete(ctx); err != nil {
			log.Error(err, "Reconcile error")
			record.Warnf(clusterScope.GCPCluster, "GCPClusterReconcile", "Reconcile error - %v", err)
			return err
		}
	}

	controllerutil.RemoveFinalizer(clusterScope.GCPCluster, infrav1.ClusterFinalizer)
	record.Event(clusterScope.GCPCluster, "GCPClusterReconcile", "Reconciled")
	return nil
}
