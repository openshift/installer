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
	"fmt"
	"slices"

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/external"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	utilexp "sigs.k8s.io/cluster-api/exp/util"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/pkg/mutators"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// AzureASOManagedMachinePoolReconciler reconciles a AzureASOManagedMachinePool object.
type AzureASOManagedMachinePoolReconciler struct {
	client.Client
	WatchFilterValue string
	Tracker          ClusterTracker

	newResourceReconciler func(*infrav1.AzureASOManagedMachinePool, []*unstructured.Unstructured) resourceReconciler
}

// ClusterTracker wraps a CAPI remote.ClusterCacheTracker.
type ClusterTracker interface {
	GetClient(context.Context, types.NamespacedName) (client.Client, error)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AzureASOManagedMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	_, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedMachinePoolReconciler.SetupWithManager",
		tele.KVP("controller", infrav1.AzureASOManagedMachinePoolKind),
	)
	defer done()

	clusterToAzureASOManagedMachinePools, err := util.ClusterToTypedObjectsMapper(mgr.GetClient(), &infrav1.AzureASOManagedMachinePoolList{}, mgr.GetScheme())
	if err != nil {
		return fmt.Errorf("failed to get Cluster to AzureASOManagedMachinePool mapper: %w", err)
	}

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.AzureASOManagedMachinePool{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue)).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToAzureASOManagedMachinePools),
			builder.WithPredicates(
				predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue),
				predicates.Any(mgr.GetScheme(), log,
					predicates.ClusterControlPlaneInitialized(mgr.GetScheme(), log),
					ClusterUpdatePauseChange(log),
				),
			),
		).
		Watches(
			&expv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(utilexp.MachinePoolToInfrastructureMapFunc(ctx,
				infrav1.GroupVersion.WithKind(infrav1.AzureASOManagedMachinePoolKind)),
			),
			builder.WithPredicates(
				predicates.ResourceHasFilterLabel(mgr.GetScheme(), log, r.WatchFilterValue),
			),
		).
		Build(r)
	if err != nil {
		return err
	}

	externalTracker := &external.ObjectTracker{
		Cache:           mgr.GetCache(),
		Controller:      c,
		Scheme:          mgr.GetScheme(),
		PredicateLogger: &log,
	}

	r.newResourceReconciler = func(asoManagedCluster *infrav1.AzureASOManagedMachinePool, resources []*unstructured.Unstructured) resourceReconciler {
		return &ResourceReconciler{
			Client:    r.Client,
			resources: resources,
			owner:     asoManagedCluster,
			watcher:   externalTracker,
		}
	}

	return nil
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedmachinepools,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedmachinepools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedmachinepools/finalizers,verbs=update

// Reconcile reconciles an AzureASOManagedMachinePool.
func (r *AzureASOManagedMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, resultErr error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedMachinePoolReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", infrav1.AzureASOManagedMachinePoolKind),
	)
	defer done()

	asoManagedMachinePool := &infrav1.AzureASOManagedMachinePool{}
	err := r.Get(ctx, req.NamespacedName, asoManagedMachinePool)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	patchHelper, err := patch.NewHelper(asoManagedMachinePool, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create patch helper: %w", err)
	}
	defer func() {
		err := patchHelper.Patch(ctx, asoManagedMachinePool)
		if err != nil && resultErr == nil {
			resultErr = err
			result = ctrl.Result{}
		}
	}()

	asoManagedMachinePool.Status.Ready = false

	machinePool, err := utilexp.GetOwnerMachinePool(ctx, r.Client, asoManagedMachinePool.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machinePool == nil {
		log.V(4).Info("Waiting for MachinePool Controller to set OwnerRef on AzureASOManagedMachinePool")
		return ctrl.Result{}, nil
	}

	machinePoolBefore := machinePool.DeepCopy()
	defer func() {
		// Skip using a patch helper here because we will never modify the MachinePool status.
		err := r.Patch(ctx, machinePool, client.MergeFrom(machinePoolBefore))
		if err != nil && resultErr == nil {
			resultErr = err
			result = ctrl.Result{}
		}
	}()

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("AzureASOManagedMachinePool owner MachinePool is missing cluster label or cluster does not exist: %w", err)
	}
	if cluster == nil {
		log.Info(fmt.Sprintf("Waiting for MachinePool controller to set %s label on MachinePool", clusterv1.ClusterNameLabel))
		return ctrl.Result{}, nil
	}
	if cluster.Spec.ControlPlaneRef == nil ||
		!matchesASOManagedAPIGroup(cluster.Spec.ControlPlaneRef.APIVersion) ||
		cluster.Spec.ControlPlaneRef.Kind != infrav1.AzureASOManagedControlPlaneKind {
		return ctrl.Result{}, reconcile.TerminalError(fmt.Errorf("AzureASOManagedMachinePool cannot be used without AzureASOManagedControlPlane"))
	}

	if annotations.IsPaused(cluster, asoManagedMachinePool) {
		return r.reconcilePaused(ctx, asoManagedMachinePool)
	}

	if !asoManagedMachinePool.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, asoManagedMachinePool, cluster)
	}

	return r.reconcileNormal(ctx, asoManagedMachinePool, machinePool, cluster)
}

func (r *AzureASOManagedMachinePoolReconciler) reconcileNormal(ctx context.Context, asoManagedMachinePool *infrav1.AzureASOManagedMachinePool, machinePool *expv1.MachinePool, cluster *clusterv1.Cluster) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedMachinePoolReconciler.reconcileNormal",
	)
	defer done()
	log.V(4).Info("reconciling normally")

	needsPatch := controllerutil.AddFinalizer(asoManagedMachinePool, clusterv1.ClusterFinalizer)
	needsPatch = AddBlockMoveAnnotation(asoManagedMachinePool) || needsPatch
	if needsPatch {
		return ctrl.Result{Requeue: true}, nil
	}

	resources, err := mutators.ApplyMutators(ctx, asoManagedMachinePool.Spec.Resources, mutators.SetAgentPoolDefaults(r.Client, machinePool))
	if err != nil {
		return ctrl.Result{}, err
	}

	var agentPoolName string
	for _, resource := range resources {
		if resource.GroupVersionKind().Group == asocontainerservicev1.GroupVersion.Group &&
			resource.GroupVersionKind().Kind == "ManagedClustersAgentPool" {
			agentPoolName = resource.GetName()
			break
		}
	}
	if agentPoolName == "" {
		return ctrl.Result{}, reconcile.TerminalError(mutators.ErrNoManagedClustersAgentPoolDefined)
	}

	resourceReconciler := r.newResourceReconciler(asoManagedMachinePool, resources)
	err = resourceReconciler.Reconcile(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile resources: %w", err)
	}
	for _, status := range asoManagedMachinePool.Status.Resources {
		if !status.Ready {
			return ctrl.Result{}, nil
		}
	}

	agentPool := &asocontainerservicev1.ManagedClustersAgentPool{}
	err = r.Get(ctx, client.ObjectKey{Namespace: asoManagedMachinePool.Namespace, Name: agentPoolName}, agentPool)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting ManagedClustersAgentPool: %w", err)
	}

	managedCluster := &asocontainerservicev1.ManagedCluster{}
	err = r.Get(ctx, client.ObjectKey{Namespace: agentPool.Namespace, Name: agentPool.Owner().Name}, managedCluster)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting ManagedCluster: %w", err)
	}
	if managedCluster.Status.NodeResourceGroup == nil {
		return ctrl.Result{}, nil
	}
	rg := *managedCluster.Status.NodeResourceGroup

	clusterClient, err := r.Tracker.GetClient(ctx, util.ObjectKey(cluster))
	if err != nil {
		return ctrl.Result{}, err
	}
	nodes := &corev1.NodeList{}
	err = clusterClient.List(ctx, nodes,
		client.MatchingLabels(expectedNodeLabels(agentPool.AzureName(), rg)),
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list nodes in workload cluster: %w", err)
	}
	providerIDs := make([]string, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		if node.Spec.ProviderID == "" {
			// the node will receive a provider id soon
			return ctrl.Result{Requeue: true}, nil
		}
		providerIDs = append(providerIDs, node.Spec.ProviderID)
	}
	// Prevent a different order from updating the spec.
	slices.Sort(providerIDs)
	asoManagedMachinePool.Spec.ProviderIDList = providerIDs
	asoManagedMachinePool.Status.Replicas = int32(ptr.Deref(agentPool.Status.Count, 0))
	if _, autoscaling := machinePool.Annotations[clusterv1.ReplicasManagedByAnnotation]; autoscaling {
		machinePool.Spec.Replicas = &asoManagedMachinePool.Status.Replicas
	}

	asoManagedMachinePool.Status.Ready = true

	return ctrl.Result{}, nil
}

func expectedNodeLabels(poolName, nodeRG string) map[string]string {
	if len(poolName) > validation.LabelValueMaxLength {
		poolName = poolName[:validation.LabelValueMaxLength]
	}
	if len(nodeRG) > validation.LabelValueMaxLength {
		nodeRG = nodeRG[:validation.LabelValueMaxLength]
	}
	return map[string]string{
		"kubernetes.azure.com/agentpool": poolName,
		"kubernetes.azure.com/cluster":   nodeRG,
	}
}

func (r *AzureASOManagedMachinePoolReconciler) reconcilePaused(ctx context.Context, asoManagedMachinePool *infrav1.AzureASOManagedMachinePool) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureASOManagedMachinePoolReconciler.reconcilePaused")
	defer done()
	log.V(4).Info("reconciling pause")

	resources, err := mutators.ToUnstructured(ctx, asoManagedMachinePool.Spec.Resources)
	if err != nil {
		return ctrl.Result{}, err
	}
	resourceReconciler := r.newResourceReconciler(asoManagedMachinePool, resources)
	err = resourceReconciler.Pause(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to pause resources: %w", err)
	}

	RemoveBlockMoveAnnotation(asoManagedMachinePool)

	return ctrl.Result{}, nil
}

func (r *AzureASOManagedMachinePoolReconciler) reconcileDelete(ctx context.Context, asoManagedMachinePool *infrav1.AzureASOManagedMachinePool, cluster *clusterv1.Cluster) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedMachinePoolReconciler.reconcileDelete",
	)
	defer done()
	log.V(4).Info("reconciling delete")

	// If the entire cluster is being deleted, this ASO ManagedClustersAgentPool will be deleted with the rest
	// of the ManagedCluster.
	if cluster.DeletionTimestamp.IsZero() {
		resourceReconciler := r.newResourceReconciler(asoManagedMachinePool, nil)
		err := resourceReconciler.Delete(ctx)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to reconcile resources: %w", err)
		}
		if len(asoManagedMachinePool.Status.Resources) > 0 {
			return ctrl.Result{}, nil
		}
	}

	controllerutil.RemoveFinalizer(asoManagedMachinePool, clusterv1.ClusterFinalizer)
	return ctrl.Result{}, nil
}
