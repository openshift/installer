/*
Copyright 2022 Nutanix

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
	"cmp"
	"context"
	"errors"
	"fmt"
	"time"

	credentialTypes "github.com/nutanix-cloud-native/prism-go-client/environment/credentials"
	corev1 "k8s.io/api/core/v1"
	kapierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	kutilerrors "k8s.io/apimachinery/pkg/util/errors"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/utils/ptr"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capiutil "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	nutanixclient "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/pkg/client"
	nctx "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/pkg/context"
)

// NutanixClusterReconciler reconciles a NutanixCluster object
type NutanixClusterReconciler struct {
	Client            client.Client
	SecretInformer    coreinformers.SecretInformer
	ConfigMapInformer coreinformers.ConfigMapInformer
	Scheme            *runtime.Scheme
	controllerConfig  *ControllerConfig
}

func NewNutanixClusterReconciler(client client.Client, secretInformer coreinformers.SecretInformer, configMapInformer coreinformers.ConfigMapInformer, scheme *runtime.Scheme, copts ...ControllerConfigOpts) (*NutanixClusterReconciler, error) {
	controllerConf := &ControllerConfig{}
	for _, opt := range copts {
		if err := opt(controllerConf); err != nil {
			return nil, err
		}
	}
	return &NutanixClusterReconciler{
		Client:            client,
		SecretInformer:    secretInformer,
		ConfigMapInformer: configMapInformer,
		Scheme:            scheme,
		controllerConfig:  controllerConf,
	}, nil
}

// SetupWithManager sets up the NutanixCluster controller with the Manager.
func (r *NutanixClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	copts := controller.Options{
		MaxConcurrentReconciles: r.controllerConfig.MaxConcurrentReconciles,
		RateLimiter:             r.controllerConfig.RateLimiter,
		SkipNameValidation:      ptr.To(r.controllerConfig.SkipNameValidation),
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("nutanixcluster-controller").
		For(&infrav1.NutanixCluster{}). // Watch the controlled, infrastructure resource.
		Watches(
			&capiv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(
				capiutil.ClusterToInfrastructureMapFunc(
					ctx,
					infrav1.GroupVersion.WithKind(infrav1.NutanixClusterKind),
					mgr.GetClient(),
					&infrav1.NutanixCluster{},
				),
			),
			builder.WithPredicates(predicates.ClusterPausedTransitionsOrInfrastructureReady(r.Scheme, ctrl.LoggerFrom(ctx))),
		).
		Watches(
			&infrav1.NutanixFailureDomain{},
			handler.EnqueueRequestsFromMapFunc(
				r.mapNutanixFailureDomainToNutanixCluster(),
			),
		).
		WithOptions(copts).
		Complete(r)
}

func (r *NutanixClusterReconciler) mapNutanixFailureDomainToNutanixCluster() handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		log := ctrl.LoggerFrom(ctx)
		nfd, ok := o.(*infrav1.NutanixFailureDomain)
		if !ok {
			log.Error(fmt.Errorf("expected a NutanixFailureDomain object but was %T", o), "unexpected type")
			return nil
		}

		// Get all the NutanixClusters in the local namespace
		nclList := &infrav1.NutanixClusterList{}
		if err := r.Client.List(ctx, nclList, client.InNamespace(nfd.Namespace)); err != nil {
			log.Error(err, "Failed to list nutanix clusters for failure domain")
			return nil
		}

		// Return those NutanixCluster instances having reference to the failure domain
		reqs := make([]ctrl.Request, 0)
		for _, ncl := range nclList.Items {
			for _, fdRef := range ncl.Spec.ControlPlaneFailureDomains {
				if fdRef.Name == nfd.Name {
					objKey := client.ObjectKey{Name: ncl.Name, Namespace: ncl.Namespace}
					reqs = append(reqs, ctrl.Request{NamespacedName: objKey})
					continue
				}
			}
		}

		return reqs
	}
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update;delete
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NutanixCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NutanixClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Reconciling the NutanixCluster")

	var err error

	// Fetch the NutanixCluster instance
	cluster := &infrav1.NutanixCluster{}
	err = r.Client.Get(ctx, req.NamespacedName, cluster)
	if err != nil {
		if kapierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.V(1).Info("NutanixCluster not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "failed to fetch the NutanixCluster object")
		return reconcile.Result{}, err
	}

	// Fetch the CAPI Cluster.
	capiCluster, err := capiutil.GetOwnerCluster(ctx, r.Client, cluster.ObjectMeta)
	if err != nil {
		log.Error(err, "failed to fetch the owner CAPI Cluster object")
		return reconcile.Result{}, err
	}
	if capiCluster == nil {
		log.Info("Waiting for Cluster Controller to set OwnerRef for the NutanixCluster object")
		return reconcile.Result{}, nil
	}
	if annotations.IsPaused(capiCluster, cluster) {
		log.Info("The NutanixCluster object linked to a cluster that is paused")
		return reconcile.Result{}, nil
	}
	log.Info(fmt.Sprintf("Fetched the owner Cluster: %s", capiCluster.Name))

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(cluster, r.Client)
	if err != nil {
		log.Error(err, "Failed to configure the patch helper")
		return ctrl.Result{Requeue: true}, nil
	}

	defer func() {
		// Always attempt to Patch the NutanixCluster object and its status after each reconciliation.
		if err := patchHelper.Patch(ctx, cluster); err != nil {
			reterr = kutilerrors.NewAggregate([]error{reterr, err})
		}
		log.V(1).Info(fmt.Sprintf("Patched NutanixCluster. Status: %+v", cluster.Status))
	}()

	if err := r.reconcileCredentialRef(ctx, cluster); err != nil {
		log.Error(err, fmt.Sprintf("error occurred while reconciling credential ref for cluster %s", capiCluster.Name))
		conditions.MarkFalse(cluster, infrav1.CredentialRefSecretOwnerSetCondition, infrav1.CredentialRefSecretOwnerSetFailed, capiv1.ConditionSeverityError, "%s", err.Error())
		return reconcile.Result{}, err
	}
	conditions.MarkTrue(cluster, infrav1.CredentialRefSecretOwnerSetCondition)

	if err := r.reconcileTrustBundleRef(ctx, cluster); err != nil {
		log.Error(err, fmt.Sprintf("error occurred while reconciling trust bundle ref for cluster %s", capiCluster.Name))
		return reconcile.Result{}, err
	}

	v3Client, err := getPrismCentralClientForCluster(ctx, cluster, r.SecretInformer, r.ConfigMapInformer)
	if err != nil {
		log.Error(err, "error occurred while fetching prism central client")
		return reconcile.Result{}, err
	}

	rctx := &nctx.ClusterContext{
		Context:        ctx,
		Cluster:        capiCluster,
		NutanixCluster: cluster,
		NutanixClient:  v3Client,
	}
	// Check for request action
	if !cluster.DeletionTimestamp.IsZero() {
		// NutanixCluster is being deleted
		return r.reconcileDelete(rctx)
	}

	return r.reconcileNormal(rctx)
}

func (r *NutanixClusterReconciler) reconcileDelete(rctx *nctx.ClusterContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(rctx.Context)
	log.Info("Handling NutanixCluster deletion")
	// Check if there are nutanixmachine resources left. Only continue if all of them have been cleaned
	nutanixMachines, err := rctx.GetNutanixMachinesInCluster(r.Client)
	if err != nil {
		log.Error(err, "error occurred while checking nutanixmachines during cluster deletion")
		return reconcile.Result{}, err
	}

	if len(nutanixMachines) > 0 {
		log.Info(fmt.Sprintf("waiting for %d nutanixmachines to be deleted", len(nutanixMachines)))
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
	}

	log.V(1).Info("no existing nutanixMachine resources found. Continuing with deleting cluster")

	err = r.reconcileCategoriesDelete(rctx)
	if err != nil {
		log.Error(err, "error occurred while running deletion of categories")
		return reconcile.Result{}, err
	}

	// delete the client from the cache
	log.Info(fmt.Sprintf("deleting nutanix prism client for cluster %s from cache", rctx.NutanixCluster.GetNamespacedName()))
	nutanixclient.NutanixClientCache.Delete(&nutanixclient.CacheParams{NutanixCluster: rctx.NutanixCluster})
	nutanixclient.NutanixClientCacheV4.Delete(&nutanixclient.CacheParams{NutanixCluster: rctx.NutanixCluster})

	if err := r.reconcileCredentialRefDelete(rctx.Context, rctx.NutanixCluster); err != nil {
		log.Error(err, fmt.Sprintf("error occurred while reconciling credential ref deletion for cluster %s", rctx.Cluster.Name))
		return reconcile.Result{}, err
	}

	if err := r.reconcileTrustBundleRefDelete(rctx.Context, rctx.NutanixCluster); err != nil {
		log.Error(err, fmt.Sprintf("error occurred while reconciling trust bundle ref deletion for cluster %s", rctx.Cluster.Name))
		return reconcile.Result{}, err
	}

	// Remove the finalizer from the NutanixCluster object
	ctrlutil.RemoveFinalizer(rctx.NutanixCluster, infrav1.NutanixClusterFinalizer)
	ctrlutil.RemoveFinalizer(rctx.NutanixCluster, infrav1.DeprecatedNutanixClusterFinalizer)

	// Remove the workload cluster client from cache
	clusterKey := apitypes.NamespacedName{
		Namespace: rctx.Cluster.Namespace,
		Name:      rctx.Cluster.Name,
	}
	nctx.RemoveRemoteClient(clusterKey)

	return reconcile.Result{}, nil
}

func (r *NutanixClusterReconciler) reconcileNormal(rctx *nctx.ClusterContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(rctx.Context)
	if rctx.NutanixCluster.Status.FailureReason != nil || rctx.NutanixCluster.Status.FailureMessage != nil {
		log.Error(fmt.Errorf("nutanix cluster has failed. Will not reconcile %s", rctx.NutanixCluster.Name), "Nutanix Cluster failed")
		return reconcile.Result{}, nil
	}
	log.Info("Handling NutanixCluster reconciling")

	// Add finalizer first if not exist to avoid the race condition between init and delete
	if !ctrlutil.ContainsFinalizer(rctx.NutanixCluster, infrav1.NutanixClusterFinalizer) {
		ctrlutil.AddFinalizer(rctx.NutanixCluster, infrav1.NutanixClusterFinalizer)
	}
	ctrlutil.RemoveFinalizer(rctx.NutanixCluster, infrav1.DeprecatedNutanixClusterFinalizer)

	// Reconciling failure domains before Ready check to allow failure domains to be modified
	if err := r.reconcileFailureDomains(rctx); err != nil {
		log.Error(err, "failed to reconcile failure domains for cluster")
		return reconcile.Result{}, err
	}

	if rctx.NutanixCluster.Status.Ready {
		log.Info("NutanixCluster is already in ready status.")
		return reconcile.Result{}, nil
	}

	err := r.reconcileCategories(rctx)
	if err != nil {
		log.Error(err, "error occurred while reconciling categories")
		// Don't return fatal error but keep retrying until categories are created.
		return reconcile.Result{}, err
	}

	rctx.NutanixCluster.Status.Ready = true
	return reconcile.Result{}, nil
}

func (r *NutanixClusterReconciler) reconcileFailureDomains(rctx *nctx.ClusterContext) error {
	log := ctrl.LoggerFrom(rctx.Context)
	log.Info("Reconciling failure domains for cluster")

	failureDomains := capiv1.FailureDomains{}
	if len(rctx.NutanixCluster.Spec.ControlPlaneFailureDomains) == 0 && len(rctx.NutanixCluster.Spec.FailureDomains) == 0 { //nolint:staticcheck // suppress complaining on Deprecated field
		log.Info("No failure domains configured for cluster.")
		conditions.MarkTrue(rctx.NutanixCluster, infrav1.NoFailureDomainsConfiguredCondition)
		conditions.Delete(rctx.NutanixCluster, infrav1.FailureDomainsValidatedCondition)

		// Reset the failure domains for nutanixcluster status
		rctx.NutanixCluster.Status.FailureDomains = failureDomains
		return nil
	}

	// Clear NoFailureDomainsConfiguredCondition condition
	conditions.Delete(rctx.NutanixCluster, infrav1.NoFailureDomainsConfiguredCondition)

	validationErrs := []error{}
	for _, fdRef := range rctx.NutanixCluster.Spec.ControlPlaneFailureDomains {
		// Fetch the referent failure domain object
		fdObj := &infrav1.NutanixFailureDomain{}
		fdKey := client.ObjectKey{Name: fdRef.Name, Namespace: rctx.NutanixCluster.Namespace}
		if err := r.Client.Get(rctx.Context, fdKey, fdObj); err != nil {
			if kapierrors.IsNotFound(err) {
				validationErrs = append(validationErrs, fmt.Errorf("not found the failure domain object with name %q", fdRef.Name))
				continue
			}
			validationErrs = append(validationErrs, fmt.Errorf("failed to fetch the failure domain object with name %q: %w", fdRef.Name, err))
			continue
		}

		// Validate the failure domain configuration
		if err := r.validateFailureDomainSpec(rctx, fdObj); err != nil {
			validationErrs = append(validationErrs, fmt.Errorf("failed to validate the failure domain %q configuration: %w", fdRef.Name, err))
			continue
		}

		// The failure domain configuration passed validation. Add it to the result map.
		failureDomains[fdObj.Name] = capiv1.FailureDomainSpec{ControlPlane: true}
	}

	// Remove below when the Deprecated field NutanixCluster.Spec.FailureDomains is removed
	for _, fd := range rctx.NutanixCluster.Spec.FailureDomains { //nolint:staticcheck // suppress complaining on Deprecated field
		failureDomains[fd.Name] = capiv1.FailureDomainSpec{ControlPlane: fd.ControlPlane}
	}

	// Set the failure domains for nutanixcluster status
	rctx.NutanixCluster.Status.FailureDomains = failureDomains

	if len(validationErrs) != 0 {
		conditions.MarkFalse(rctx.NutanixCluster, infrav1.FailureDomainsValidatedCondition,
			infrav1.FailureDomainsMisconfiguredReason, capiv1.ConditionSeverityWarning, "%s", errors.Join(validationErrs...).Error())
		return nil
	}

	conditions.MarkTrue(rctx.NutanixCluster, infrav1.FailureDomainsValidatedCondition)
	return nil
}

// validateFailureDomainSpec is to validate the input failure domain's spec configuration.
// It returns error if validation fails, and returns nil if validation succeeds.
func (r *NutanixClusterReconciler) validateFailureDomainSpec(rctx *nctx.ClusterContext, fd *infrav1.NutanixFailureDomain) error {
	pe := fd.Spec.PrismElementCluster
	peUUID, err := GetPEUUID(rctx.Context, rctx.NutanixClient, pe.Name, pe.UUID)
	if err != nil {
		return err
	}

	subnets := fd.Spec.Subnets
	_, err = GetSubnetUUIDList(rctx.Context, rctx.NutanixClient, subnets, peUUID)
	if err != nil {
		return err
	}

	return nil
}

func (r *NutanixClusterReconciler) reconcileCategories(rctx *nctx.ClusterContext) error {
	log := ctrl.LoggerFrom(rctx.Context)
	log.Info("Reconciling categories for cluster")
	defaultCategories := GetDefaultCAPICategoryIdentifiers(rctx.Cluster.Name)
	_, err := GetOrCreateCategories(rctx.Context, rctx.NutanixClient, defaultCategories)
	if err != nil {
		conditions.MarkFalse(rctx.NutanixCluster, infrav1.ClusterCategoryCreatedCondition, infrav1.ClusterCategoryCreationFailed, capiv1.ConditionSeverityError, "%s", err.Error())
		return err
	}
	conditions.MarkTrue(rctx.NutanixCluster, infrav1.ClusterCategoryCreatedCondition)
	return nil
}

func (r *NutanixClusterReconciler) reconcileCategoriesDelete(rctx *nctx.ClusterContext) error {
	log := ctrl.LoggerFrom(rctx.Context)
	log.Info(fmt.Sprintf("Reconciling deletion of categories for cluster %s", rctx.Cluster.Name))
	if conditions.IsTrue(rctx.NutanixCluster, infrav1.ClusterCategoryCreatedCondition) ||
		conditions.GetReason(rctx.NutanixCluster, infrav1.ClusterCategoryCreatedCondition) == infrav1.DeletionFailed {
		defaultCategories := GetDefaultCAPICategoryIdentifiers(rctx.Cluster.Name)
		obsoleteCategories := GetObsoleteDefaultCAPICategoryIdentifiers(rctx.Cluster.Name)
		err := DeleteCategories(rctx.Context, rctx.NutanixClient, defaultCategories, obsoleteCategories)
		if err != nil {
			conditions.MarkFalse(rctx.NutanixCluster, infrav1.ClusterCategoryCreatedCondition, infrav1.DeletionFailed, capiv1.ConditionSeverityWarning, "%s", err.Error())
			return err
		}
	} else {
		log.V(1).Info(fmt.Sprintf("skipping category deletion since they were not created for cluster %s", rctx.Cluster.Name))
	}
	conditions.MarkFalse(rctx.NutanixCluster, infrav1.ClusterCategoryCreatedCondition, capiv1.DeletingReason, capiv1.ConditionSeverityInfo, "")
	return nil
}

func (r *NutanixClusterReconciler) reconcileCredentialRefDelete(ctx context.Context, nutanixCluster *infrav1.NutanixCluster) error {
	log := ctrl.LoggerFrom(ctx)
	credentialRef, err := getPrismCentralCredentialRefForCluster(nutanixCluster)
	if err != nil {
		log.Error(err, fmt.Sprintf("error occurred while getting credential ref for cluster %s", nutanixCluster.Name))
		return err
	}
	if credentialRef == nil {
		log.V(1).Info(fmt.Sprintf("Credential ref is nil for cluster %s. Ignoring since object must be deleted", nutanixCluster.Name))
		return nil
	}
	log.V(1).Info(fmt.Sprintf("Credential ref is kind Secret for cluster %s. Continue with deletion of secret", nutanixCluster.Name))
	secret := &corev1.Secret{}
	secretKey := client.ObjectKey{
		Namespace: nutanixCluster.Namespace,
		Name:      credentialRef.Name,
	}
	err = r.Client.Get(ctx, secretKey, secret)
	if err != nil {
		if kapierrors.IsNotFound(err) {
			log.V(1).Info(fmt.Sprintf("Secret %s in namespace %s for cluster %s not found. Ignoring since object must be deleted", secret.Name, secret.Namespace, nutanixCluster.Name))
			return nil
		}
		return err
	}

	ctrlutil.RemoveFinalizer(secret, infrav1.NutanixClusterCredentialFinalizer)
	ctrlutil.RemoveFinalizer(secret, infrav1.DeprecatedNutanixClusterCredentialFinalizer)
	log.V(1).Info(fmt.Sprintf("removing finalizers from secret %s in namespace %s for cluster %s", secret.Name, secret.Namespace, nutanixCluster.Name))
	if err := r.Client.Update(ctx, secret); err != nil {
		return err
	}

	if secret.DeletionTimestamp.IsZero() {
		log.Info(fmt.Sprintf("removing secret %s in namespace %s for cluster %s", secret.Name, secret.Namespace, nutanixCluster.Name))
		if err := r.Client.Delete(ctx, secret); err != nil && !kapierrors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func (r *NutanixClusterReconciler) reconcileTrustBundleRef(ctx context.Context, nutanixCluster *infrav1.NutanixCluster) error {
	log := ctrl.LoggerFrom(ctx)
	trustBundleRef := nutanixCluster.GetPrismCentralTrustBundle()
	if trustBundleRef == nil {
		log.Info(fmt.Sprintf("trust bundle ref is nil for cluster %s", nutanixCluster.Name))
		return nil
	}

	// get the trust bundle configmap
	configMap := &corev1.ConfigMap{}
	configMapKey := client.ObjectKey{
		Namespace: cmp.Or(trustBundleRef.Namespace, nutanixCluster.Namespace),
		Name:      trustBundleRef.Name,
	}
	if err := r.Client.Get(ctx, configMapKey, configMap); err != nil {
		log.Error(err, "error occurred while fetching trust bundle configmap", "nutanixCluster", nutanixCluster.Name)
		conditions.MarkFalse(nutanixCluster, infrav1.TrustBundleSecretOwnerSetCondition, infrav1.TrustBundleSecretOwnerSetFailed, capiv1.ConditionSeverityError, "%s", err.Error())
		return err
	}

	if !capiutil.IsOwnedByObject(configMap, nutanixCluster) {
		// Check if another nutanixCluster already has set ownerRef. Secret can only be owned by one nutanixCluster object
		if capiutil.HasOwner(configMap.OwnerReferences, infrav1.GroupVersion.String(), []string{nutanixCluster.Kind}) {
			return fmt.Errorf("configmap %s/%s already owned by another nutanixCluster object", configMap.Namespace, configMap.Name)
		}

		configMap.OwnerReferences = capiutil.EnsureOwnerRef(configMap.OwnerReferences, metav1.OwnerReference{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       nutanixCluster.Kind,
			UID:        nutanixCluster.UID,
			Name:       nutanixCluster.Name,
		})
	}

	if !ctrlutil.ContainsFinalizer(configMap, infrav1.NutanixClusterCredentialFinalizer) {
		ctrlutil.AddFinalizer(configMap, infrav1.NutanixClusterCredentialFinalizer)
	}
	ctrlutil.RemoveFinalizer(configMap, infrav1.DeprecatedNutanixClusterCredentialFinalizer)

	if err := r.Client.Update(ctx, configMap); err != nil {
		log.Error(err, "error occurred while updating trust bundle configmap", "nutanixCluster", nutanixCluster)
		conditions.MarkFalse(nutanixCluster, infrav1.TrustBundleSecretOwnerSetCondition, infrav1.TrustBundleSecretOwnerSetFailed, capiv1.ConditionSeverityError, "%s", err.Error())
		return err
	}

	conditions.MarkTrue(nutanixCluster, infrav1.TrustBundleSecretOwnerSetCondition)
	return nil
}

func (r *NutanixClusterReconciler) reconcileTrustBundleRefDelete(ctx context.Context, nutanixCluster *infrav1.NutanixCluster) error {
	log := ctrl.LoggerFrom(ctx)
	trustBundleRef := nutanixCluster.GetPrismCentralTrustBundle()
	if trustBundleRef == nil {
		log.Info(fmt.Sprintf("trust bundle ref is nil for cluster %s", nutanixCluster.Name))
		return nil
	}

	configMapKey := client.ObjectKey{
		Namespace: cmp.Or(trustBundleRef.Namespace, nutanixCluster.Namespace),
		Name:      trustBundleRef.Name,
	}

	configMap := &corev1.ConfigMap{}
	if err := r.Client.Get(ctx, configMapKey, configMap); err != nil {
		if kapierrors.IsNotFound(err) {
			log.Info(fmt.Sprintf("configmap %s/%s for cluster %s not found. Ignoring since object must be deleted", configMapKey.Namespace, configMapKey.Name, nutanixCluster.Name))
			return nil
		}

		return err
	}

	ctrlutil.RemoveFinalizer(configMap, infrav1.NutanixClusterCredentialFinalizer)
	ctrlutil.RemoveFinalizer(configMap, infrav1.DeprecatedNutanixClusterCredentialFinalizer)
	log.V(1).Info(fmt.Sprintf("removing finalizers from configmap %s/%s for cluster %s", configMap.Namespace, configMap.Name, nutanixCluster.Name))
	if err := r.Client.Update(ctx, configMap); err != nil {
		return err
	}

	if configMap.DeletionTimestamp.IsZero() {
		log.Info(fmt.Sprintf("removing configmap %s/%s for cluster %s", configMap.Namespace, configMap.Name, nutanixCluster.Name))
		if err := r.Client.Delete(ctx, configMap); err != nil && !kapierrors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func (r *NutanixClusterReconciler) reconcileCredentialRef(ctx context.Context, nutanixCluster *infrav1.NutanixCluster) error {
	log := ctrl.LoggerFrom(ctx)
	credentialRef, err := getPrismCentralCredentialRefForCluster(nutanixCluster)
	if err != nil {
		return err
	}

	secret := &corev1.Secret{}
	if credentialRef == nil {
		return nil
	}

	log.V(1).Info(fmt.Sprintf("credential ref is kind Secret for cluster %s", nutanixCluster.Name))
	secretKey := client.ObjectKey{
		Namespace: nutanixCluster.Namespace,
		Name:      credentialRef.Name,
	}

	if err := r.Client.Get(ctx, secretKey, secret); err != nil {
		errorMsg := fmt.Errorf("error occurred while fetching cluster %s secret for credential ref: %v", nutanixCluster.Name, err)
		log.Error(errorMsg, "error occurred fetching cluster")
		return errorMsg
	}

	// Check if ownerRef is already set on nutanixCluster object
	if !capiutil.IsOwnedByObject(secret, nutanixCluster) {
		// Check if another nutanixCluster already has set ownerRef. Secret can only be owned by one nutanixCluster object
		if capiutil.HasOwner(secret.OwnerReferences, infrav1.GroupVersion.String(), []string{
			nutanixCluster.Kind,
		}) {
			return fmt.Errorf("secret %s already owned by another nutanixCluster object", secret.Name)
		}
		// Set nutanixCluster ownerRef on the secret
		secret.OwnerReferences = capiutil.EnsureOwnerRef(secret.OwnerReferences, metav1.OwnerReference{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       nutanixCluster.Kind,
			UID:        nutanixCluster.UID,
			Name:       nutanixCluster.Name,
		})
	}

	if !ctrlutil.ContainsFinalizer(secret, infrav1.NutanixClusterCredentialFinalizer) {
		ctrlutil.AddFinalizer(secret, infrav1.NutanixClusterCredentialFinalizer)
	}
	ctrlutil.RemoveFinalizer(secret, infrav1.DeprecatedNutanixClusterCredentialFinalizer)

	err = r.Client.Update(ctx, secret)
	if err != nil {
		errorMsg := fmt.Errorf("failed to update secret for cluster %s: %v", nutanixCluster.Name, err)
		log.Error(errorMsg, "failed to update secret")
		return errorMsg
	}

	return nil
}

// getPrismCentralCredentialRefForCluster calls nutanixCluster.GetPrismCentralCredentialRef() function
// and returns an error if nutanixCluster is nil
func getPrismCentralCredentialRefForCluster(nutanixCluster *infrav1.NutanixCluster) (*credentialTypes.NutanixCredentialReference, error) {
	if nutanixCluster == nil {
		return nil, fmt.Errorf("cannot get credential reference if nutanix cluster object is nil")
	}
	return nutanixCluster.GetPrismCentralCredentialRef()
}
