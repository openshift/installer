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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/identities"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AzureJSONTemplateReconciler reconciles Azure json secrets for AzureMachineTemplate objects.
type AzureJSONTemplateReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	Timeouts         reconciler.Timeouts
	WatchFilterValue string
}

// SetupWithManager initializes this controller with a manager.
func (r *AzureJSONTemplateReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	_, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.AzureJSONTemplateReconciler.SetupWithManager",
	)
	defer done()

	azureMachineTemplateMapper, err := util.ClusterToTypedObjectsMapper(r.Client, &infrav1.AzureMachineTemplateList{}, mgr.GetScheme())
	if err != nil {
		return errors.Wrap(err, "failed to create mapper for Cluster to AzureMachineTemplates")
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.AzureMachineTemplate{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		Owns(&corev1.Secret{}).
		// Add a watch on Clusters to requeue when the infraRef is set. This is needed because the infraRef is not initially
		// set in Clusters created from a ClusterClass.
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(azureMachineTemplateMapper),
			builder.WithPredicates(
				predicates.ClusterUnpausedAndInfrastructureReady(log),
				predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue),
			),
		).
		Complete(r)
}

// Reconcile reconciles Azure json secrets for Azure machine templates.
func (r *AzureJSONTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.AzureJSONTemplateReconciler.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", "AzureMachineTemplate"),
	)
	defer done()

	// Fetch the AzureMachineTemplate instance
	azureMachineTemplate := &infrav1.AzureMachineTemplate{}
	err := r.Get(ctx, req.NamespacedName, azureMachineTemplate)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("object was not found")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, azureMachineTemplate.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Return early if the object or Cluster is paused.
	if annotations.IsPaused(cluster, azureMachineTemplate) {
		log.Info("AzureMachineTemplate or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	// only look at azure clusters
	if cluster.Spec.InfrastructureRef == nil {
		log.Info("infra ref is nil")
		return ctrl.Result{}, nil
	}
	if cluster.Spec.InfrastructureRef.Kind != infrav1.AzureClusterKind {
		log.WithValues("kind", cluster.Spec.InfrastructureRef.Kind).Info("infra ref was not an AzureCluster")
		return ctrl.Result{}, nil
	}

	// fetch the corresponding azure cluster
	azureCluster := &infrav1.AzureCluster{}
	azureClusterName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}

	if err := r.Get(ctx, azureClusterName, azureCluster); err != nil {
		log.Error(err, "failed to fetch AzureCluster")
		return reconcile.Result{}, err
	}

	// Create the scope.
	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:       r.Client,
		Cluster:      cluster,
		AzureCluster: azureCluster,
		Timeouts:     r.Timeouts,
	})
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create scope")
	}

	apiVersion, kind := infrav1.GroupVersion.WithKind("AzureMachineTemplate").ToAPIVersionAndKind()
	owner := metav1.OwnerReference{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       azureMachineTemplate.GetName(),
		UID:        azureMachineTemplate.GetUID(),
		Controller: ptr.To(true),
	}

	// Construct secret for this machine template
	userAssignedIdentityIfExists := ""
	if len(azureMachineTemplate.Spec.Template.Spec.UserAssignedIdentities) > 0 {
		var identitiesClient identities.Client
		identitiesClient, err := identities.NewClient(clusterScope)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to create identities client")
		}
		parsed, err := azureutil.ParseResourceID(azureMachineTemplate.Spec.Template.Spec.UserAssignedIdentities[0].ProviderID)
		if err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to parse ProviderID %s", azureMachineTemplate.Spec.Template.Spec.UserAssignedIdentities[0].ProviderID)
		}
		if parsed.SubscriptionID != clusterScope.SubscriptionID() {
			identitiesClient, err = identities.NewClientBySub(clusterScope, parsed.SubscriptionID)
			if err != nil {
				return reconcile.Result{}, errors.Wrapf(err, "failed to create identities client from subscription ID %s", parsed.SubscriptionID)
			}
		}
		userAssignedIdentityIfExists, err = identitiesClient.GetClientID(
			ctx, azureMachineTemplate.Spec.Template.Spec.UserAssignedIdentities[0].ProviderID)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to get user-assigned identity ClientID")
		}
	}

	if azureMachineTemplate.Spec.Template.Spec.Identity == infrav1.VMIdentityNone {
		log.Info(fmt.Sprintf("WARNING, %s", spIdentityWarning))
		r.Recorder.Eventf(azureMachineTemplate, corev1.EventTypeWarning, "VMIdentityNone", spIdentityWarning)
	}

	newSecret, err := GetCloudProviderSecret(
		clusterScope,
		azureMachineTemplate.Namespace,
		azureMachineTemplate.Name,
		owner,
		azureMachineTemplate.Spec.Template.Spec.Identity,
		userAssignedIdentityIfExists,
	)

	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create cloud provider config")
	}

	if err := reconcileAzureSecret(ctx, r.Client, owner, newSecret, clusterScope.ClusterName()); err != nil {
		r.Recorder.Eventf(azureMachineTemplate, corev1.EventTypeWarning, "Error reconciling cloud provider secret for AzureMachineTemplate", err.Error())
		return ctrl.Result{}, errors.Wrap(err, "failed to reconcile azure secret")
	}

	return ctrl.Result{}, nil
}
