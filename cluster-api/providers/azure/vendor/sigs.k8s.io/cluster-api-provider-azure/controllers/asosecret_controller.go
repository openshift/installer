/*
Copyright 2023 The Kubernetes Authors.

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

	asoconfig "github.com/Azure/azure-service-operator/v2/pkg/common/config"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/scope"
	"sigs.k8s.io/cluster-api-provider-azure/util/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ASOSecretReconciler reconciles ASO secrets associated with AzureCluster objects.
type ASOSecretReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	Timeouts         reconciler.Timeouts
	WatchFilterValue string
}

// SetupWithManager initializes this controller with a manager.
func (asos *ASOSecretReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	_, log, done := tele.StartSpanWithLogger(ctx,
		"controllers.ASOSecretReconciler.SetupWithManager",
		tele.KVP("controller", "ASOSecret"),
	)
	defer done()

	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.AzureCluster{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, asos.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		Named("ASOSecret").
		Owns(&corev1.Secret{}).
		Build(asos)
	if err != nil {
		return errors.Wrap(err, "error creating controller")
	}

	// Add a watch on infrav1.AzureManagedControlPlane.
	if err = c.Watch(
		source.Kind(mgr.GetCache(), &infrav1.AzureManagedControlPlane{}),
		&handler.EnqueueRequestForObject{},
		predicates.ResourceNotPausedAndHasFilterLabel(log, asos.WatchFilterValue),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready AzureManagedControlPlanes")
	}

	// Add a watch on ASO secrets owned by an AzureManagedControlPlane
	if err = c.Watch(
		source.Kind(mgr.GetCache(), &corev1.Secret{}),
		handler.EnqueueRequestForOwner(asos.Scheme(), asos.RESTMapper(), &infrav1.AzureManagedControlPlane{}, handler.OnlyControllerOwner()),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for secrets")
	}

	// Add a watch on clusterv1.Cluster object for unpause notifications.
	if err = c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, infrav1.GroupVersion.WithKind(infrav1.AzureClusterKind), mgr.GetClient(), &infrav1.AzureCluster{})),
		predicates.ClusterUnpaused(log),
		predicates.ResourceNotPausedAndHasFilterLabel(log, asos.WatchFilterValue),
	); err != nil {
		return errors.Wrap(err, "failed adding a watch for ready clusters")
	}

	return nil
}

// Reconcile reconciles the ASO secrets associated with AzureCluster objects.
func (asos *ASOSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, asos.Timeouts.DefaultedLoopTimeout())
	defer cancel()

	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.ASOSecret.Reconcile",
		tele.KVP("namespace", req.Namespace),
		tele.KVP("name", req.Name),
		tele.KVP("kind", infrav1.AzureClusterKind),
	)
	defer done()

	log = log.WithValues("namespace", req.Namespace)

	// asoSecretOwner is the resource that created the identity. This could be either an AzureCluster or AzureManagedControlPlane (if AKS is enabled).
	// check for AzureCluster first and if it is not found, check for AzureManagedControlPlane.
	var asoSecretOwner client.Object

	azureCluster := &infrav1.AzureCluster{}
	checkForManagedControlPlane := false
	// Fetch the AzureCluster or AzureManagedControlPlane instance
	asoSecretOwner = azureCluster
	err := asos.Get(ctx, req.NamespacedName, azureCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			checkForManagedControlPlane = true
		} else {
			return reconcile.Result{}, err
		}
	} else {
		log = log.WithValues("AzureCluster", req.Name)
	}

	if checkForManagedControlPlane {
		// Fetch the AzureManagedControlPlane instance instead
		azureManagedControlPlane := &infrav1.AzureManagedControlPlane{}
		asoSecretOwner = azureManagedControlPlane
		err = asos.Get(ctx, req.NamespacedName, azureManagedControlPlane)
		if err != nil {
			if apierrors.IsNotFound(err) {
				asos.Recorder.Eventf(azureCluster, corev1.EventTypeNormal, "AzureClusterObjectNotFound",
					fmt.Sprintf("AzureCluster object %s/%s not found", req.Namespace, req.Name))
				asos.Recorder.Eventf(azureManagedControlPlane, corev1.EventTypeNormal, "AzureManagedControlPlaneObjectNotFound",
					fmt.Sprintf("AzureManagedControlPlane object %s/%s not found", req.Namespace, req.Name))
				log.Info("object was not found")
				return reconcile.Result{}, nil
			} else {
				return reconcile.Result{}, err
			}
		} else {
			log = log.WithValues("AzureManagedControlPlane", req.Name)
		}
	}

	var clusterIdentity *corev1.ObjectReference
	var cluster *clusterv1.Cluster
	var azureClient scope.AzureClients

	switch ownerType := asoSecretOwner.(type) {
	case *infrav1.AzureCluster:
		clusterIdentity = ownerType.Spec.IdentityRef

		// Fetch the Cluster.
		cluster, err = util.GetOwnerCluster(ctx, asos.Client, ownerType.ObjectMeta)
		if err != nil {
			return reconcile.Result{}, err
		}
		if cluster == nil {
			log.Info("Cluster Controller has not yet set OwnerRef")
			return reconcile.Result{}, nil
		}

		// Create the scope.
		clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
			Client:       asos.Client,
			Cluster:      cluster,
			AzureCluster: ownerType,
			Timeouts:     asos.Timeouts,
		})
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to create scope")
		}

		azureClient = clusterScope.AzureClients

	case *infrav1.AzureManagedControlPlane:
		clusterIdentity = ownerType.Spec.IdentityRef

		// Fetch the Cluster.
		cluster, err = util.GetOwnerCluster(ctx, asos.Client, ownerType.ObjectMeta)
		if err != nil {
			return reconcile.Result{}, err
		}
		if cluster == nil {
			log.Info("Cluster Controller has not yet set OwnerRef")
			return reconcile.Result{}, nil
		}

		// Create the scope.
		clusterScope, err := scope.NewManagedControlPlaneScope(ctx, scope.ManagedControlPlaneScopeParams{
			Client:       asos.Client,
			Cluster:      cluster,
			ControlPlane: ownerType,
			Timeouts:     asos.Timeouts,
		})
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to create scope")
		}

		azureClient = clusterScope.AzureClients
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		asos.Recorder.Eventf(asoSecretOwner, corev1.EventTypeNormal, "OwnerRefNotFound",
			fmt.Sprintf("Cluster Controller has not yet set OwnerRef for object %s/%s", req.Namespace, req.Name))
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Return early if the ASO Secret Owner(AzureCluster or AzureManagedControlPlane) or Cluster is paused.
	if annotations.IsPaused(cluster, asoSecretOwner) {
		log.Info(fmt.Sprintf("%T or linked Cluster is marked as paused. Won't reconcile", asoSecretOwner))
		asos.Recorder.Eventf(asoSecretOwner, corev1.EventTypeNormal, "ClusterPaused",
			fmt.Sprintf("%T or linked Cluster is marked as paused. Won't reconcile", asoSecretOwner))
		return ctrl.Result{}, nil
	}

	// Construct the ASO secret for this Cluster
	newASOSecret, err := asos.createSecretFromClusterIdentity(ctx, clusterIdentity, cluster, azureClient)
	if err != nil {
		return reconcile.Result{}, err
	}

	gvk := asoSecretOwner.GetObjectKind().GroupVersionKind()
	owner := metav1.OwnerReference{
		APIVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Name:       asoSecretOwner.GetName(),
		UID:        asoSecretOwner.GetUID(),
		Controller: ptr.To(true),
	}

	newASOSecret.OwnerReferences = []metav1.OwnerReference{owner}

	if err := reconcileAzureSecret(ctx, asos.Client, owner, newASOSecret, cluster.GetName()); err != nil {
		asos.Recorder.Eventf(cluster, corev1.EventTypeWarning, "Error reconciling ASO secret", err.Error())
		return ctrl.Result{}, errors.Wrap(err, "failed to reconcile ASO secret")
	}

	return ctrl.Result{}, nil
}

func (asos *ASOSecretReconciler) createSecretFromClusterIdentity(ctx context.Context, clusterIdentity *corev1.ObjectReference, cluster *clusterv1.Cluster, azureClient scope.AzureClients) (*corev1.Secret, error) {
	newASOSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aso.GetASOSecretName(cluster.GetName()),
			Namespace: cluster.GetNamespace(),
			Labels: map[string]string{
				cluster.GetName(): string(infrav1.ResourceLifecycleOwned),
			},
		},
		Data: map[string][]byte{
			asoconfig.AzureSubscriptionID: []byte(azureClient.SubscriptionID()),
		},
	}

	// if the namespace isn't specified then assume it's in the same namespace as the Cluster's one
	namespace := clusterIdentity.Namespace
	if namespace == "" {
		namespace = cluster.GetNamespace()
	}
	identity := &infrav1.AzureClusterIdentity{}
	key := client.ObjectKey{
		Name:      clusterIdentity.Name,
		Namespace: namespace,
	}
	if err := asos.Get(ctx, key, identity); err != nil {
		return nil, errors.Wrap(err, "failed to retrieve AzureClusterIdentity")
	}

	newASOSecret.Data[asoconfig.AzureTenantID] = []byte(identity.Spec.TenantID)
	newASOSecret.Data[asoconfig.AzureClientID] = []byte(identity.Spec.ClientID)

	// If the identity type is WorkloadIdentity or UserAssignedMSI, then we don't need to fetch the secret so return early
	if identity.Spec.Type == infrav1.WorkloadIdentity {
		newASOSecret.Data[asoconfig.AuthMode] = []byte(asoconfig.WorkloadIdentityAuthMode)
		return newASOSecret, nil
	}
	if identity.Spec.Type == infrav1.UserAssignedMSI {
		newASOSecret.Data[asoconfig.AuthMode] = []byte(asoconfig.PodIdentityAuthMode)
		return newASOSecret, nil
	}

	// Fetch identity secret, if it exists
	key = types.NamespacedName{
		Namespace: identity.Spec.ClientSecret.Namespace,
		Name:      identity.Spec.ClientSecret.Name,
	}
	identitySecret := &corev1.Secret{}
	err := asos.Get(ctx, key, identitySecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch AzureClusterIdentity secret")
	}

	switch identity.Spec.Type {
	case infrav1.ServicePrincipal, infrav1.ManualServicePrincipal:
		newASOSecret.Data[asoconfig.AzureClientSecret] = identitySecret.Data[scope.AzureSecretKey]
	case infrav1.ServicePrincipalCertificate:
		newASOSecret.Data[asoconfig.AzureClientCertificate] = identitySecret.Data["certificate"]
		newASOSecret.Data[asoconfig.AzureClientCertificatePassword] = identitySecret.Data["password"]
	}
	return newASOSecret, nil
}
