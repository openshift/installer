/*
Copyright 2021 The Kubernetes Authors.

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
	goctx "context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	vmwarecontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=providerserviceaccounts,verbs=get;list;watch;
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=providerserviceaccounts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

const (
	// ProviderServiceAccountControllerName defines the controller used when creating clients.
	ProviderServiceAccountControllerName = "provider-serviceaccount-controller"
	kindProviderServiceAccount           = "ProviderServiceAccount"
	systemServiceAccountPrefix           = "system.serviceaccount"
)

// AddServiceAccountProviderControllerToManager adds this controller to the provided manager.
func AddServiceAccountProviderControllerToManager(ctx *context.ControllerManagerContext, mgr manager.Manager, tracker *remote.ClusterCacheTracker, options controller.Options) error {
	var (
		controlledType     = &vmwarev1.ProviderServiceAccount{}
		controlledTypeName = reflect.TypeOf(controlledType).Elem().Name()

		controllerNameShort = fmt.Sprintf("%s-controller", strings.ToLower(controlledTypeName))
		controllerNameLong  = fmt.Sprintf("%s/%s/%s", ctx.Namespace, ctx.Name, controllerNameShort)
	)

	controllerContext := &context.ControllerContext{
		ControllerManagerContext: ctx,
		Name:                     controllerNameShort,
		Recorder:                 record.New(mgr.GetEventRecorderFor(controllerNameLong)),
		Logger:                   ctx.Logger.WithName(controllerNameShort),
	}
	r := ServiceAccountReconciler{
		ControllerContext:         controllerContext,
		remoteClusterCacheTracker: tracker,
	}

	clusterToInfraFn := clusterToSupervisorInfrastructureMapFunc(ctx)

	return ctrl.NewControllerManagedBy(mgr).For(controlledType).
		WithOptions(options).
		// Watch a ProviderServiceAccount
		Watches(
			&vmwarev1.ProviderServiceAccount{},
			handler.EnqueueRequestsFromMapFunc(r.providerServiceAccountToVSphereCluster),
		).
		// Watches the secrets to re-enqueue once the service account token is set
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.secretToVSphereCluster),
		).
		// Watches clusters and reconciles the vSphereCluster
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(func(ctx goctx.Context, o client.Object) []reconcile.Request {
				requests := clusterToInfraFn(ctx, o)
				if requests == nil {
					return nil
				}

				c := &vmwarev1.VSphereCluster{}
				if err := r.Client.Get(ctx, requests[0].NamespacedName, c); err != nil {
					r.Logger.V(4).Error(err, "Failed to get VSphereCluster")
					return nil
				}

				if annotations.IsExternallyManaged(c) {
					r.Logger.V(4).Info("VSphereCluster is externally managed, skipping mapping.")
					return nil
				}
				return requests
			}),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), ctx.WatchFilterValue)).
		Complete(r)
}

func clusterToSupervisorInfrastructureMapFunc(managerContext *context.ControllerManagerContext) handler.MapFunc {
	gvk := vmwarev1.GroupVersion.WithKind(reflect.TypeOf(&vmwarev1.VSphereCluster{}).Elem().Name())
	return clusterutilv1.ClusterToInfrastructureMapFunc(managerContext, gvk, managerContext.Client, &vmwarev1.VSphereCluster{})
}

type ServiceAccountReconciler struct {
	*context.ControllerContext

	remoteClusterCacheTracker *remote.ClusterCacheTracker
}

func (r ServiceAccountReconciler) Reconcile(_ goctx.Context, req reconcile.Request) (_ reconcile.Result, reterr error) {
	r.ControllerContext.Logger.V(4).Info("Starting Reconcile")

	// Get the vSphereCluster for this request.
	vsphereCluster := &vmwarev1.VSphereCluster{}
	clusterKey := client.ObjectKey{Namespace: req.Namespace, Name: req.Name}
	if err := r.Client.Get(r, clusterKey, vsphereCluster); err != nil {
		if apierrors.IsNotFound(err) {
			r.Logger.Info("vSphereCluster not found, won't reconcile", "cluster", clusterKey)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(vsphereCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(
			err,
			"failed to init patch helper for %s %s/%s",
			vsphereCluster.GroupVersionKind(),
			vsphereCluster.Namespace,
			vsphereCluster.Name)
	}

	// Create the cluster context for this request.
	clusterContext := &vmwarecontext.ClusterContext{
		ControllerContext: r.ControllerContext,
		VSphereCluster:    vsphereCluster,
		Logger:            r.Logger.WithName(req.Namespace).WithName(req.Name),
		PatchHelper:       patchHelper,
	}

	// Always issue a patch when exiting this function so changes to the
	// resource are patched back to the API server.
	defer func() {
		if err := clusterContext.Patch(); err != nil {
			if reterr == nil {
				reterr = err
			} else {
				clusterContext.Logger.Error(err, "patch failed", "cluster", clusterContext.String())
			}
		}
	}()
	if !vsphereCluster.DeletionTimestamp.IsZero() {
		return r.ReconcileDelete(clusterContext)
	}

	cluster, err := clusterutilv1.GetClusterFromMetadata(r, r.Client, vsphereCluster.ObjectMeta)
	if err != nil {
		r.Logger.Info("unable to get capi cluster from vsphereCluster", "err", err)
		return reconcile.Result{}, nil
	}

	// Pause reconciliation if entire vSphereCluster or Cluster is paused
	if annotations.IsPaused(cluster, vsphereCluster) {
		r.Logger.V(4).Info("VSphereCluster %s/%s linked to a cluster that is paused",
			vsphereCluster.Namespace, vsphereCluster.Name)
		return reconcile.Result{}, nil
	}

	// We cannot proceed until we are able to access the target cluster. Until
	// then just return a no-op and wait for the next sync. This will occur when
	// the Cluster's status is updated with a reference to the secret that has
	// the Kubeconfig data used to access the target cluster.
	guestClient, err := r.remoteClusterCacheTracker.GetClient(clusterContext, client.ObjectKeyFromObject(cluster))
	if err != nil {
		if errors.Is(err, remote.ErrClusterLocked) {
			r.Logger.V(5).Info("Requeuing because another worker has the lock on the ClusterCacheTracker")
			return ctrl.Result{Requeue: true}, nil
		}
		clusterContext.Logger.Info("The control plane is not ready yet", "err", err)
		return reconcile.Result{RequeueAfter: clusterNotReadyRequeueTime}, nil
	}

	// Defer to the Reconciler for reconciling a non-delete event.
	return r.ReconcileNormal(&vmwarecontext.GuestClusterContext{
		ClusterContext: clusterContext,
		GuestClient:    guestClient,
	})
}

func (r ServiceAccountReconciler) ReconcileDelete(ctx *vmwarecontext.ClusterContext) (reconcile.Result, error) {
	ctx.Logger.V(4).Info("Reconciling deleting Provider ServiceAccounts", "cluster", ctx.VSphereCluster.Name)

	pSvcAccounts, err := getProviderServiceAccounts(ctx)
	if err != nil {
		ctx.Logger.Error(err, "Error fetching provider serviceaccounts")
		return reconcile.Result{}, err
	}

	for _, pSvcAccount := range pSvcAccounts {
		// Delete entries for configmap with serviceaccount
		if err := r.deleteServiceAccountConfigMap(ctx, pSvcAccount); err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "unable to delete configmap entry for provider serviceaccount %s", pSvcAccount.Name)
		}
	}

	return reconcile.Result{}, nil
}

func (r ServiceAccountReconciler) ReconcileNormal(ctx *vmwarecontext.GuestClusterContext) (_ reconcile.Result, reterr error) {
	ctx.Logger.V(4).Info("Reconciling Provider ServiceAccount", "cluster", ctx.VSphereCluster.Name)
	defer func() {
		if reterr != nil {
			conditions.MarkFalse(ctx.VSphereCluster, vmwarev1.ProviderServiceAccountsReadyCondition, vmwarev1.ProviderServiceAccountsReconciliationFailedReason,
				clusterv1.ConditionSeverityWarning, reterr.Error())
		} else {
			conditions.MarkTrue(ctx.VSphereCluster, vmwarev1.ProviderServiceAccountsReadyCondition)
		}
	}()

	pSvcAccounts, err := getProviderServiceAccounts(ctx.ClusterContext)
	if err != nil {
		ctx.Logger.Error(err, "Error fetching provider serviceaccounts")
		return reconcile.Result{}, err
	}
	err = r.ensureProviderServiceAccounts(ctx, pSvcAccounts)
	if err != nil {
		ctx.Logger.Error(err, "Error ensuring provider serviceaccounts")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// Ensure service accounts from provider spec is created.
func (r ServiceAccountReconciler) ensureProviderServiceAccounts(ctx *vmwarecontext.GuestClusterContext, pSvcAccounts []vmwarev1.ProviderServiceAccount) error {
	for i, pSvcAccount := range pSvcAccounts {
		if ctx.Cluster != nil && annotations.IsPaused(ctx.Cluster, &(pSvcAccounts[i])) {
			r.Logger.V(4).Info("ProviderServiceAccount %s/%s linked to a cluster that is paused or has pause annotation",
				pSvcAccount.Namespace, pSvcAccount.Name)
			continue
		}

		// 1. Create service accounts by the name specified in Provider Spec
		if err := r.ensureServiceAccount(ctx.ClusterContext, pSvcAccount); err != nil {
			return errors.Wrapf(err, "unable to create provider serviceaccount %s", pSvcAccount.Name)
		}
		// 2. Update configmap with serviceaccount
		if err := r.ensureServiceAccountConfigMap(ctx.ClusterContext, pSvcAccount); err != nil {
			return errors.Wrapf(err, "unable to sync configmap for provider serviceaccount %s", pSvcAccount.Name)
		}
		// 3. Create secret of Service account token type for the service account
		if err := r.ensureServiceAccountSecret(ctx.ClusterContext, pSvcAccount); err != nil {
			return errors.Wrapf(err, "unable to create provider serviceaccount secret %s", getServiceAccountSecretName(pSvcAccount))
		}

		// 4. Create the associated role for the service account
		if err := r.ensureRole(ctx.ClusterContext, pSvcAccount); err != nil {
			return errors.Wrapf(err, "unable to create role for provider serviceaccount %s", pSvcAccount.Name)
		}

		// 5. Create the associated roleBinding for the service account
		if err := r.ensureRoleBinding(ctx.ClusterContext, pSvcAccount); err != nil {
			return errors.Wrapf(err, "unable to create rolebinding for provider serviceaccount %s", pSvcAccount.Name)
		}

		// 6. Sync the service account with the target
		if err := r.syncServiceAccountSecret(ctx, pSvcAccount); err != nil {
			return errors.Wrapf(err, "unable to sync secret for provider serviceaccount %s", pSvcAccount.Name)
		}
	}
	return nil
}

func (r ServiceAccountReconciler) ensureServiceAccount(ctx *vmwarecontext.ClusterContext, pSvcAccount vmwarev1.ProviderServiceAccount) error {
	svcAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getServiceAccountName(pSvcAccount),
			Namespace: pSvcAccount.Namespace,
		},
	}
	logger := ctx.Logger.WithValues("providerserviceaccount", pSvcAccount.Name, "serviceaccount", svcAccount.Name)
	err := util.SetControllerReferenceWithOverride(&pSvcAccount, &svcAccount, ctx.Scheme)
	if err != nil {
		return err
	}
	logger.V(4).Info("Creating service account")
	err = ctx.Client.Create(ctx, &svcAccount)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		// Note: We skip updating the service account because the token controller updates the service account with a
		// secret and we don't want to overwrite it with an empty secret.
		return err
	}
	return nil
}

func (r ServiceAccountReconciler) ensureServiceAccountSecret(ctx *vmwarecontext.ClusterContext, pSvcAccount vmwarev1.ProviderServiceAccount) error {
	secret := corev1.Secret{
		Type: corev1.SecretTypeServiceAccountToken,
		ObjectMeta: metav1.ObjectMeta{
			Name:      getServiceAccountSecretName(pSvcAccount),
			Namespace: pSvcAccount.Namespace,
			Annotations: map[string]string{
				// denotes that this secret holds the token for the service account
				corev1.ServiceAccountNameKey: getServiceAccountName(pSvcAccount),
			},
		},
	}

	logger := ctx.Logger.WithValues("providerserviceaccount", pSvcAccount.Name, "secret", secret.Name)
	err := util.SetControllerReferenceWithOverride(&pSvcAccount, &secret, ctx.Scheme)
	if err != nil {
		return err
	}
	logger.V(4).Info("Creating service account secret")
	err = ctx.Client.Create(ctx, &secret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		// Note: We skip updating the service account because the token controller updates the service account with a
		// secret and we don't want to overwrite it with an empty secret.
		return err
	}
	return nil
}

func (r ServiceAccountReconciler) ensureRole(ctx *vmwarecontext.ClusterContext, pSvcAccount vmwarev1.ProviderServiceAccount) error {
	role := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRoleName(pSvcAccount),
			Namespace: pSvcAccount.Namespace,
		},
	}
	logger := ctx.Logger.WithValues("providerserviceaccount", pSvcAccount.Name, "role", role.Name)
	logger.V(4).Info("Creating or updating role")
	_, err := controllerutil.CreateOrPatch(ctx, ctx.Client, &role, func() error {
		if err := util.SetControllerReferenceWithOverride(&pSvcAccount, &role, ctx.Scheme); err != nil {
			return err
		}
		role.Rules = pSvcAccount.Spec.Rules
		return nil
	})
	return err
}

func (r ServiceAccountReconciler) ensureRoleBinding(ctx *vmwarecontext.ClusterContext, pSvcAccount vmwarev1.ProviderServiceAccount) error {
	roleName := getRoleName(pSvcAccount)
	svcAccountName := getServiceAccountName(pSvcAccount)
	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRoleBindingName(pSvcAccount),
			Namespace: pSvcAccount.Namespace,
		},
	}
	logger := ctx.Logger.WithValues("providerserviceaccount", pSvcAccount.Name, "rolebinding", roleBinding.Name)
	logger.V(4).Info("Creating or updating rolebinding")

	err := ctx.Client.Get(ctx, types.NamespacedName{Name: getRoleBindingName(pSvcAccount), Namespace: pSvcAccount.Namespace}, &roleBinding)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if err == nil {
		// If the roleRef needs changing, we have to delete the rolebinding and recreate it.
		if roleBinding.RoleRef.Name != roleName || roleBinding.RoleRef.Kind != "Role" || roleBinding.RoleRef.APIGroup != rbacv1.GroupName {
			if err := ctx.Client.Delete(ctx, &roleBinding); err != nil {
				return err
			}
		}
	}

	_, err = controllerutil.CreateOrPatch(ctx, ctx.Client, &roleBinding, func() error {
		if err := util.SetControllerReferenceWithOverride(&pSvcAccount, &roleBinding, ctx.Scheme); err != nil {
			return err
		}
		roleBinding.RoleRef = rbacv1.RoleRef{
			Name:     roleName,
			Kind:     "Role",
			APIGroup: rbacv1.GroupName,
		}
		roleBinding.Subjects = []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				APIGroup:  "",
				Name:      svcAccountName,
				Namespace: pSvcAccount.Namespace,
			},
		}
		return nil
	})
	return err
}

func (r ServiceAccountReconciler) syncServiceAccountSecret(ctx *vmwarecontext.GuestClusterContext, pSvcAccount vmwarev1.ProviderServiceAccount) error {
	logger := ctx.Logger.WithValues("providerserviceaccount", pSvcAccount.Name)
	logger.V(4).Info("Attempting to sync token secret for provider service account")

	secretName := getServiceAccountSecretName(pSvcAccount)
	logger.V(4).Info("Fetching secret for service account token details", "secret", secretName)
	var svcAccountTokenSecret corev1.Secret
	err := ctx.Client.Get(ctx, types.NamespacedName{Name: secretName, Namespace: pSvcAccount.Namespace}, &svcAccountTokenSecret)
	if err != nil {
		return err
	}
	// Check if token data exists
	if len(svcAccountTokenSecret.Data) == 0 {
		// Note: We don't have to requeue here because we have a watch on the secret and the cluster should be reconciled
		// when a secret has token data populated.
		logger.Info("Skipping sync secret for provider service account: secret has no data", "secret", secretName)
		return nil
	}

	// Create the target namespace if it is not existing
	targetNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: pSvcAccount.Spec.TargetNamespace,
		},
	}

	if err = ctx.GuestClient.Get(ctx, client.ObjectKey{Name: pSvcAccount.Spec.TargetNamespace}, targetNamespace); err != nil {
		if apierrors.IsNotFound(err) {
			err = ctx.GuestClient.Create(ctx, targetNamespace)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	targetSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pSvcAccount.Spec.TargetSecretName,
			Namespace: pSvcAccount.Spec.TargetNamespace,
		},
	}
	logger.V(4).Info("Creating or updating secret in cluster", "namespace", targetSecret.Namespace, "name", targetSecret.Name)
	_, err = controllerutil.CreateOrPatch(ctx, ctx.GuestClient, targetSecret, func() error {
		targetSecret.Data = svcAccountTokenSecret.Data
		return nil
	})
	return err
}

func (r ServiceAccountReconciler) getConfigMapAndBuffer(ctx *vmwarecontext.ClusterContext) (*corev1.ConfigMap, *corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}

	if err := ctx.Client.Get(ctx, GetCMNamespaceName(), configMap); err != nil {
		return nil, nil, err
	}

	configMapBuffer := &corev1.ConfigMap{}
	configMapBuffer.Name = configMap.Name
	configMapBuffer.Namespace = configMap.Namespace
	return configMapBuffer, configMap, nil
}

func (r ServiceAccountReconciler) deleteServiceAccountConfigMap(ctx *vmwarecontext.ClusterContext, svcAccount vmwarev1.ProviderServiceAccount) error {
	logger := ctx.Logger.WithValues("providerserviceaccount", svcAccount.Name)

	svcAccountName := getSystemServiceAccountFullName(svcAccount)
	configMapBuffer, configMap, err := r.getConfigMapAndBuffer(ctx)
	if err != nil {
		return err
	}
	if valid, exist := configMap.Data[svcAccountName]; !exist || valid != strconv.FormatBool(true) {
		// Service account name is not in the config map
		return nil
	}
	logger.Info("Deleting config map entry for provider service account")
	_, err = controllerutil.CreateOrPatch(ctx, ctx.Client, configMapBuffer, func() error {
		configMapBuffer.Data = configMap.Data
		delete(configMapBuffer.Data, svcAccountName)
		return nil
	})
	return err
}

func (r ServiceAccountReconciler) ensureServiceAccountConfigMap(ctx *vmwarecontext.ClusterContext, svcAccount vmwarev1.ProviderServiceAccount) error {
	logger := ctx.Logger.WithValues("providerserviceaccount", svcAccount.Name)

	svcAccountName := getSystemServiceAccountFullName(svcAccount)
	configMapBuffer, configMap, err := r.getConfigMapAndBuffer(ctx)
	if err != nil {
		return err
	}
	if valid, exist := configMap.Data[svcAccountName]; exist && valid == strconv.FormatBool(true) {
		// Service account name is already in the config map
		return nil
	}
	logger.Info("Updating config map for provider service account")
	_, err = controllerutil.CreateOrPatch(ctx, ctx.Client, configMapBuffer, func() error {
		configMapBuffer.Data = configMap.Data
		configMapBuffer.Data[svcAccountName] = "true"
		return nil
	})
	return err
}

func getProviderServiceAccounts(ctx *vmwarecontext.ClusterContext) ([]vmwarev1.ProviderServiceAccount, error) {
	var pSvcAccounts []vmwarev1.ProviderServiceAccount

	pSvcAccountList := vmwarev1.ProviderServiceAccountList{}
	if err := ctx.Client.List(ctx, &pSvcAccountList, client.InNamespace(ctx.VSphereCluster.Namespace)); err != nil {
		return nil, err
	}

	for _, pSvcAccount := range pSvcAccountList.Items {
		// step to clean up the target secret in the guest cluster. Note: when the provider service account is deleted
		// all the associated serviceaccounts are deleted. Hence, the bearer token in the target
		// secret will be rendered invalid. Still, it's a good practice to delete the secret that we created.
		if pSvcAccount.DeletionTimestamp != nil {
			continue
		}
		ref := pSvcAccount.Spec.Ref
		if ref != nil && ref.Name == ctx.VSphereCluster.Name {
			pSvcAccounts = append(pSvcAccounts, pSvcAccount)
		}
	}
	return pSvcAccounts, nil
}

func getRoleName(pSvcAccount vmwarev1.ProviderServiceAccount) string {
	return pSvcAccount.Name
}

func getRoleBindingName(pSvcAccount vmwarev1.ProviderServiceAccount) string {
	return pSvcAccount.Name
}

func getServiceAccountName(pSvcAccount vmwarev1.ProviderServiceAccount) string {
	return pSvcAccount.Name
}

func getServiceAccountSecretName(pSvcAccount vmwarev1.ProviderServiceAccount) string {
	return fmt.Sprintf("%s-secret", pSvcAccount.Name)
}

func getSystemServiceAccountFullName(pSvcAccount vmwarev1.ProviderServiceAccount) string {
	return fmt.Sprintf("%s.%s.%s", systemServiceAccountPrefix, getServiceAccountNamespace(pSvcAccount), getServiceAccountName(pSvcAccount))
}

func getServiceAccountNamespace(pSvcAccount vmwarev1.ProviderServiceAccount) string {
	return pSvcAccount.Namespace
}

// GetCMNamespaceName gets capi valid modifier configmap metadata.
func GetCMNamespaceName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: os.Getenv("SERVICE_ACCOUNTS_CM_NAMESPACE"),
		Name:      os.Getenv("SERVICE_ACCOUNTS_CM_NAME"),
	}
}

// secretToVSphereCluster is a mapper function used to enqueue reconcile.Request objects.
// It accepts a Secret object owned by the controller and fetches the service account
// that contains the token and creates a reconcile.Request for the vmwarev1.VSphereCluster object.
func (r ServiceAccountReconciler) secretToVSphereCluster(ctx goctx.Context, o client.Object) []reconcile.Request {
	secret, ok := o.(*corev1.Secret)
	if !ok {
		return nil
	}

	ownerRef := metav1.GetControllerOf(secret)
	if ownerRef != nil && ownerRef.Kind == kindProviderServiceAccount {
		if !metav1.HasAnnotation(secret.ObjectMeta, corev1.ServiceAccountNameKey) {
			return nil
		}
		svcAccountName := secret.GetAnnotations()[corev1.ServiceAccountNameKey]
		svcAccount := &corev1.ServiceAccount{}
		if err := r.Client.Get(ctx, client.ObjectKey{
			Namespace: secret.Namespace,
			Name:      svcAccountName,
		}, svcAccount); err != nil {
			return nil
		}
		return r.serviceAccountToVSphereCluster(svcAccount)
	}
	return nil
}

// serviceAccountToVSphereCluster is a mapper function used to enqueue reconcile.Request objects.
// From the watched object owned by this controller, it creates reconcile.Request object
// for the vmwarev1.VSphereCluster object that owns the watched object.
func (r ServiceAccountReconciler) serviceAccountToVSphereCluster(o client.Object) []reconcile.Request {
	// We do this because this controller is effectively a vSphereCluster controller that reconciles its
	// dependent ProviderServiceAccount objects.
	ownerRef := metav1.GetControllerOf(o)
	if ownerRef != nil && ownerRef.Kind == kindProviderServiceAccount {
		key := types.NamespacedName{Namespace: o.GetNamespace(), Name: ownerRef.Name}
		pSvcAccount := &vmwarev1.ProviderServiceAccount{}
		if err := r.Client.Get(r.Context, key, pSvcAccount); err != nil {
			return nil
		}
		return toVSphereClusterRequest(pSvcAccount)
	}
	return nil
}

// providerServiceAccountToVSphereCluster is a mapper function used to enqueue reconcile.Request objects.
func (r ServiceAccountReconciler) providerServiceAccountToVSphereCluster(_ goctx.Context, o client.Object) []reconcile.Request {
	providerServiceAccount, ok := o.(*vmwarev1.ProviderServiceAccount)
	if !ok {
		return nil
	}

	return toVSphereClusterRequest(providerServiceAccount)
}

func toVSphereClusterRequest(providerServiceAccount *vmwarev1.ProviderServiceAccount) []reconcile.Request {
	vsphereClusterRef := providerServiceAccount.Spec.Ref
	if vsphereClusterRef == nil || vsphereClusterRef.Name == "" {
		return nil
	}
	return []reconcile.Request{
		{NamespacedName: client.ObjectKey{Namespace: providerServiceAccount.Namespace, Name: vsphereClusterRef.Name}},
	}
}
