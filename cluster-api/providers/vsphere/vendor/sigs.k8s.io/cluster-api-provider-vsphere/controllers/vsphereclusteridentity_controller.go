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
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	pkgidentity "sigs.k8s.io/cluster-api-provider-vsphere/pkg/identity"
)

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vsphereclusteridentities,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vsphereclusteridentities/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;patch;update;delete

// AddVsphereClusterIdentityControllerToManager adds a VSphereClusterIdentity controller to the controller manager.
func AddVsphereClusterIdentityControllerToManager(ctx context.Context, controllerManagerCtx *capvcontext.ControllerManagerContext, mgr manager.Manager, options controller.Options) error {
	reconciler := clusterIdentityReconciler{
		ControllerManagerCtx: controllerManagerCtx,
		Client:               controllerManagerCtx.Client,
		Recorder:             mgr.GetEventRecorderFor("vsphereclusteridentity-controller"),
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.VSphereClusterIdentity{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), controllerManagerCtx.WatchFilterValue)).
		Complete(reconciler)
}

type clusterIdentityReconciler struct {
	ControllerManagerCtx *capvcontext.ControllerManagerContext
	Client               client.Client
	Recorder             record.EventRecorder
}

func (r clusterIdentityReconciler) Reconcile(ctx context.Context, req reconcile.Request) (_ reconcile.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	identity := &infrav1.VSphereClusterIdentity{}
	if err := r.Client.Get(ctx, req.NamespacedName, identity); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if annotations.HasPaused(identity) {
		log.Info("Reconciliation is paused for this object")
		return reconcile.Result{}, nil
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(identity, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to initialize patch helper")
	}

	defer func() {
		conditions.SetSummary(identity, conditions.WithConditions(infrav1.CredentialsAvailableCondidtion))

		if err := patchHelper.Patch(ctx, identity); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	if !identity.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, identity)
	}

	// Add a finalizer and requeue to ensure that the secret is deleted when the identity is deleted.
	if !ctrlutil.ContainsFinalizer(identity, infrav1.VSphereClusterIdentityFinalizer) {
		ctrlutil.AddFinalizer(identity, infrav1.VSphereClusterIdentityFinalizer)
		return reconcile.Result{}, nil
	}

	// fetch secret
	secret := &corev1.Secret{}
	secretKey := client.ObjectKey{
		Namespace: r.ControllerManagerCtx.Namespace,
		Name:      identity.Spec.SecretName,
	}
	if err := r.Client.Get(ctx, secretKey, secret); err != nil {
		conditions.MarkFalse(identity, infrav1.CredentialsAvailableCondidtion, infrav1.SecretNotAvailableReason, clusterv1.ConditionSeverityWarning, err.Error())
		return reconcile.Result{}, errors.Wrapf(err, "failed to get Secret %s", klog.KRef(secretKey.Namespace, secretKey.Name))
	}

	// If this secret is owned by a different VSphereClusterIdentity or a VSphereCluster, mark the identity as not ready and return an error.
	if !clusterutilv1.IsOwnedByObject(secret, identity) && pkgidentity.IsOwnedByIdentityOrCluster(secret.GetOwnerReferences()) {
		conditions.MarkFalse(identity, infrav1.CredentialsAvailableCondidtion, infrav1.SecretAlreadyInUseReason, clusterv1.ConditionSeverityError, "secret being used by another Cluster/VSphereIdentity")
		identity.Status.Ready = false
		return reconcile.Result{}, errors.New("secret being used by another Cluster/VSphereIdentity")
	}

	// Ensure the VSphereClusterIdentity is set as the owner of the secret, and that the reference has an up to date APIVersion.
	secret.SetOwnerReferences(
		clusterutilv1.EnsureOwnerRef(secret.GetOwnerReferences(),
			metav1.OwnerReference{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       "VSphereClusterIdentity",
				Name:       identity.Name,
				UID:        identity.UID,
			}))

	if !ctrlutil.ContainsFinalizer(secret, infrav1.SecretIdentitySetFinalizer) {
		ctrlutil.AddFinalizer(secret, infrav1.SecretIdentitySetFinalizer)
	}
	err = r.Client.Update(ctx, secret)
	if err != nil {
		conditions.MarkFalse(identity, infrav1.CredentialsAvailableCondidtion, infrav1.SecretOwnerReferenceFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		return reconcile.Result{}, err
	}

	conditions.MarkTrue(identity, infrav1.CredentialsAvailableCondidtion)
	identity.Status.Ready = true
	return reconcile.Result{}, nil
}

func (r clusterIdentityReconciler) reconcileDelete(ctx context.Context, identity *infrav1.VSphereClusterIdentity) error {
	log := ctrl.LoggerFrom(ctx)
	secret := &corev1.Secret{}
	secretKey := client.ObjectKey{
		Namespace: r.ControllerManagerCtx.Namespace,
		Name:      identity.Spec.SecretName,
	}
	err := r.Client.Get(ctx, secretKey, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// The secret no longer exists. Remove the finalizer from the VSphereClusterIdentity.
			ctrlutil.RemoveFinalizer(identity, infrav1.VSphereClusterIdentityFinalizer)
			return nil
		}
		return err
	}

	if ctrlutil.RemoveFinalizer(secret, infrav1.SecretIdentitySetFinalizer) {
		log.Info(fmt.Sprintf("Removing finalizer %s", infrav1.SecretIdentitySetFinalizer), "Secret", klog.KObj(secret))
		if err := r.Client.Update(ctx, secret); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("failed to update Secret %s", klog.KObj(secret)))
		}
	}

	if secret.DeletionTimestamp.IsZero() {
		log.Info("Deleting Secret", "Secret", klog.KObj(secret))
		if err := r.Client.Delete(ctx, secret); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("failed to delete Secret %s", klog.KObj(secret)))
		}
	}

	// Remove the finalizer from the identity as all cleanup is complete.
	ctrlutil.RemoveFinalizer(identity, infrav1.VSphereClusterIdentityFinalizer)
	return nil
}
