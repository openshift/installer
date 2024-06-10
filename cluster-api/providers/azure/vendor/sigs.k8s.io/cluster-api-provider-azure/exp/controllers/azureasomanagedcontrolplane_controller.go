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

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	infracontroller "sigs.k8s.io/cluster-api-provider-azure/controllers"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-azure/exp/mutators"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/secret"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var errInvalidClusterKind = errors.New("AzureASOManagedControlPlane cannot be used without AzureASOManagedCluster")

// AzureASOManagedControlPlaneReconciler reconciles a AzureASOManagedControlPlane object.
type AzureASOManagedControlPlaneReconciler struct {
	client.Client
	WatchFilterValue string

	newResourceReconciler func(*infrav1exp.AzureASOManagedControlPlane, []*unstructured.Unstructured) resourceReconciler
}

// SetupWithManager sets up the controller with the Manager.
func (r *AzureASOManagedControlPlaneReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	_, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedControlPlaneReconciler.SetupWithManager",
		tele.KVP("controller", infrav1exp.AzureASOManagedControlPlaneKind),
	)
	defer done()

	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1exp.AzureASOManagedControlPlane{}).
		WithEventFilter(predicates.ResourceHasFilterLabel(log, r.WatchFilterValue)).
		Watches(&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToAzureASOManagedControlPlane),
			builder.WithPredicates(
				predicates.ResourceHasFilterLabel(log, r.WatchFilterValue),
				infracontroller.ClusterPauseChangeAndInfrastructureReady(log),
			),
		).
		// User errors that CAPZ passes through agentPoolProfiles on create must be fixed in the
		// AzureASOManagedMachinePool, so trigger a reconciliation to consume those fixes.
		Watches(
			&infrav1exp.AzureASOManagedMachinePool{},
			handler.EnqueueRequestsFromMapFunc(r.azureASOManagedMachinePoolToAzureASOManagedControlPlane),
		).
		Owns(&corev1.Secret{}).
		Build(r)
	if err != nil {
		return err
	}

	externalTracker := &external.ObjectTracker{
		Cache:      mgr.GetCache(),
		Controller: c,
	}

	r.newResourceReconciler = func(asoManagedCluster *infrav1exp.AzureASOManagedControlPlane, resources []*unstructured.Unstructured) resourceReconciler {
		return &ResourceReconciler{
			Client:    r.Client,
			resources: resources,
			owner:     asoManagedCluster,
			watcher:   externalTracker,
		}
	}

	return nil
}

func clusterToAzureASOManagedControlPlane(_ context.Context, o client.Object) []ctrl.Request {
	controlPlaneRef := o.(*clusterv1.Cluster).Spec.ControlPlaneRef
	if controlPlaneRef != nil &&
		controlPlaneRef.APIVersion == infrav1exp.GroupVersion.Identifier() &&
		controlPlaneRef.Kind == infrav1exp.AzureASOManagedControlPlaneKind {
		return []ctrl.Request{{NamespacedName: client.ObjectKey{Namespace: controlPlaneRef.Namespace, Name: controlPlaneRef.Name}}}
	}
	return nil
}

func (r *AzureASOManagedControlPlaneReconciler) azureASOManagedMachinePoolToAzureASOManagedControlPlane(ctx context.Context, o client.Object) []ctrl.Request {
	asoManagedMachinePool := o.(*infrav1exp.AzureASOManagedMachinePool)
	clusterName := asoManagedMachinePool.Labels[clusterv1.ClusterNameLabel]
	if clusterName == "" {
		return nil
	}
	cluster, err := util.GetClusterByName(ctx, r.Client, asoManagedMachinePool.Namespace, clusterName)
	if client.IgnoreNotFound(err) != nil || cluster == nil {
		return nil
	}
	return clusterToAzureASOManagedControlPlane(ctx, cluster)
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedcontrolplanes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedcontrolplanes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedcontrolplanes/finalizers,verbs=update

// Reconcile reconciles an AzureASOManagedControlPlane.
func (r *AzureASOManagedControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, resultErr error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedControlPlaneReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", infrav1exp.AzureASOManagedControlPlaneKind),
	)
	defer done()

	asoManagedControlPlane := &infrav1exp.AzureASOManagedControlPlane{}
	err := r.Get(ctx, req.NamespacedName, asoManagedControlPlane)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	patchHelper, err := patch.NewHelper(asoManagedControlPlane, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create patch helper: %w", err)
	}
	defer func() {
		err := patchHelper.Patch(ctx, asoManagedControlPlane)
		if err != nil && resultErr == nil {
			resultErr = err
			result = ctrl.Result{}
		}
	}()

	asoManagedControlPlane.Status.Ready = false
	asoManagedControlPlane.Status.Initialized = false

	cluster, err := util.GetOwnerCluster(ctx, r.Client, asoManagedControlPlane.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}

	if cluster != nil && cluster.Spec.Paused ||
		annotations.HasPaused(asoManagedControlPlane) {
		return r.reconcilePaused(ctx, asoManagedControlPlane)
	}

	if !asoManagedControlPlane.GetDeletionTimestamp().IsZero() {
		return r.reconcileDelete(ctx, asoManagedControlPlane)
	}

	return r.reconcileNormal(ctx, asoManagedControlPlane, cluster)
}

func (r *AzureASOManagedControlPlaneReconciler) reconcileNormal(ctx context.Context, asoManagedControlPlane *infrav1exp.AzureASOManagedControlPlane, cluster *clusterv1.Cluster) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedControlPlaneReconciler.reconcileNormal",
	)
	defer done()
	log.V(4).Info("reconciling normally")

	if cluster == nil {
		log.V(4).Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}
	if cluster.Spec.InfrastructureRef == nil ||
		cluster.Spec.InfrastructureRef.APIVersion != infrav1exp.GroupVersion.Identifier() ||
		cluster.Spec.InfrastructureRef.Kind != infrav1exp.AzureASOManagedClusterKind {
		return ctrl.Result{}, reconcile.TerminalError(errInvalidClusterKind)
	}

	needsPatch := controllerutil.AddFinalizer(asoManagedControlPlane, infrav1exp.AzureASOManagedControlPlaneFinalizer)
	needsPatch = infracontroller.AddBlockMoveAnnotation(asoManagedControlPlane) || needsPatch
	if needsPatch {
		return ctrl.Result{Requeue: true}, nil
	}

	resources, err := mutators.ApplyMutators(ctx, asoManagedControlPlane.Spec.Resources, mutators.SetManagedClusterDefaults(r.Client, asoManagedControlPlane, cluster))
	if err != nil {
		return ctrl.Result{}, err
	}

	var managedClusterName string
	for _, resource := range resources {
		if resource.GroupVersionKind().Group == asocontainerservicev1.GroupVersion.Group &&
			resource.GroupVersionKind().Kind == "ManagedCluster" {
			managedClusterName = resource.GetName()
			break
		}
	}
	if managedClusterName == "" {
		return ctrl.Result{}, reconcile.TerminalError(mutators.ErrNoManagedClusterDefined)
	}

	resourceReconciler := r.newResourceReconciler(asoManagedControlPlane, resources)
	err = resourceReconciler.Reconcile(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile resources: %w", err)
	}
	for _, status := range asoManagedControlPlane.Status.Resources {
		if !status.Ready {
			return ctrl.Result{}, nil
		}
	}

	managedCluster := &asocontainerservicev1.ManagedCluster{}
	err = r.Get(ctx, client.ObjectKey{Namespace: asoManagedControlPlane.Namespace, Name: managedClusterName}, managedCluster)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting ManagedCluster: %w", err)
	}

	asoManagedControlPlane.Status.ControlPlaneEndpoint = getControlPlaneEndpoint(managedCluster)
	if managedCluster.Status.CurrentKubernetesVersion != nil {
		asoManagedControlPlane.Status.Version = "v" + *managedCluster.Status.CurrentKubernetesVersion
	}

	err = r.reconcileKubeconfig(ctx, asoManagedControlPlane, cluster, managedCluster)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile kubeconfig: %w", err)
	}

	asoManagedControlPlane.Status.Ready = !asoManagedControlPlane.Status.ControlPlaneEndpoint.IsZero()
	// The AKS API doesn't allow us to distinguish between CAPI's definitions of "initialized" and "ready" so
	// we treat them equivalently.
	asoManagedControlPlane.Status.Initialized = asoManagedControlPlane.Status.Ready

	return ctrl.Result{}, nil
}

func (r *AzureASOManagedControlPlaneReconciler) reconcileKubeconfig(ctx context.Context, asoManagedControlPlane *infrav1exp.AzureASOManagedControlPlane, cluster *clusterv1.Cluster, managedCluster *asocontainerservicev1.ManagedCluster) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedControlPlaneReconciler.reconcileKubeconfig",
	)
	defer done()

	var secretRef *genruntime.SecretDestination
	if managedCluster.Spec.OperatorSpec != nil &&
		managedCluster.Spec.OperatorSpec.Secrets != nil {
		secretRef = managedCluster.Spec.OperatorSpec.Secrets.UserCredentials
		if managedCluster.Spec.OperatorSpec.Secrets.AdminCredentials != nil {
			secretRef = managedCluster.Spec.OperatorSpec.Secrets.AdminCredentials
		}
	}
	if secretRef == nil {
		return reconcile.TerminalError(fmt.Errorf("ManagedCluster must define at least one of spec.operatorSpec.secrets.{userCredentials,adminCredentials}"))
	}
	asoKubeconfig := &corev1.Secret{}
	err := r.Get(ctx, client.ObjectKey{Namespace: cluster.Namespace, Name: secretRef.Name}, asoKubeconfig)
	if err != nil {
		return fmt.Errorf("failed to fetch secret created by ASO: %w", err)
	}

	expectedSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name(cluster.Name, secret.Kubeconfig),
			Namespace: cluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(asoManagedControlPlane, infrav1exp.GroupVersion.WithKind(infrav1exp.AzureASOManagedControlPlaneKind)),
			},
			Labels: map[string]string{clusterv1.ClusterNameLabel: cluster.Name},
		},
		Data: map[string][]byte{
			secret.KubeconfigDataName: asoKubeconfig.Data[secretRef.Key],
		},
	}

	return r.Patch(ctx, expectedSecret, client.Apply, client.FieldOwner("capz-manager"), client.ForceOwnership)
}

//nolint:unparam // an empty ctrl.Result is always returned here, leaving it as-is to avoid churn in refactoring later if that changes.
func (r *AzureASOManagedControlPlaneReconciler) reconcilePaused(ctx context.Context, asoManagedControlPlane *infrav1exp.AzureASOManagedControlPlane) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureASOManagedControlPlaneReconciler.reconcilePaused")
	defer done()
	log.V(4).Info("reconciling pause")

	resources, err := mutators.ToUnstructured(ctx, asoManagedControlPlane.Spec.Resources)
	if err != nil {
		return ctrl.Result{}, err
	}
	resourceReconciler := r.newResourceReconciler(asoManagedControlPlane, resources)
	err = resourceReconciler.Pause(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to pause resources: %w", err)
	}

	infracontroller.RemoveBlockMoveAnnotation(asoManagedControlPlane)

	return ctrl.Result{}, nil
}

//nolint:unparam // an empty ctrl.Result is always returned here, leaving it as-is to avoid churn in refactoring later if that changes.
func (r *AzureASOManagedControlPlaneReconciler) reconcileDelete(ctx context.Context, asoManagedControlPlane *infrav1exp.AzureASOManagedControlPlane) (ctrl.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureASOManagedControlPlaneReconciler.reconcileDelete",
	)
	defer done()
	log.V(4).Info("reconciling delete")

	resources, err := mutators.ToUnstructured(ctx, asoManagedControlPlane.Spec.Resources)
	if err != nil {
		return ctrl.Result{}, err
	}
	resourceReconciler := r.newResourceReconciler(asoManagedControlPlane, resources)
	err = resourceReconciler.Delete(ctx)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile resources: %w", err)
	}
	if len(asoManagedControlPlane.Status.Resources) > 0 {
		return ctrl.Result{}, nil
	}

	controllerutil.RemoveFinalizer(asoManagedControlPlane, infrav1exp.AzureASOManagedControlPlaneFinalizer)
	return ctrl.Result{}, nil
}

func getControlPlaneEndpoint(managedCluster *asocontainerservicev1.ManagedCluster) clusterv1.APIEndpoint {
	if managedCluster.Status.PrivateFQDN != nil {
		return clusterv1.APIEndpoint{
			Host: *managedCluster.Status.PrivateFQDN,
			Port: 443,
		}
	}
	if managedCluster.Status.Fqdn != nil {
		return clusterv1.APIEndpoint{
			Host: *managedCluster.Status.Fqdn,
			Port: 443,
		}
	}
	return clusterv1.APIEndpoint{}
}
