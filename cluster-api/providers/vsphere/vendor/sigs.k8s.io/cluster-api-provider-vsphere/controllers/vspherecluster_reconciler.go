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

// Package controllers contains controllers for CAPV objects.
package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	pkgerrors "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/identity"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
	infrautilv1 "sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

type clusterReconciler struct {
	ControllerManagerContext *capvcontext.ControllerManagerContext
	Client                   client.Client

	vmService               services.VimMachineService
	clusterModuleReconciler Reconciler
}

// Reconcile ensures the back-end state reflects the Kubernetes resource state intent.
func (r *clusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Get the VSphereCluster resource for this request.
	vsphereCluster := &infrav1.VSphereCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, vsphereCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the CAPI Cluster.
	cluster, err := clusterutilv1.GetOwnerCluster(ctx, r.Client, vsphereCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if cluster == nil {
		log.Info("Waiting for Cluster controller to set OwnerRef on VSphereCluster")
		return reconcile.Result{}, nil
	}
	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	if annotations.IsPaused(cluster, vsphereCluster) {
		log.Info("Reconciliation is paused for this object")
		return reconcile.Result{}, nil
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(vsphereCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, pkgerrors.Wrap(err, "failed to initialize patch helper")
	}

	// Create the cluster context for this request.
	clusterContext := &capvcontext.ClusterContext{
		Cluster:        cluster,
		VSphereCluster: vsphereCluster,
		PatchHelper:    patchHelper,
	}

	// Always issue a patch when exiting this function so changes to the
	// resource are patched back to the API server.
	defer func() {
		if err := clusterContext.Patch(ctx); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Handle deleted clusters
	if !vsphereCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterContext)
	}

	// If the VSphereCluster doesn't have our finalizer, add it.
	// Requeue immediately after adding finalizer to avoid the race condition between init and delete
	if !ctrlutil.ContainsFinalizer(vsphereCluster, infrav1.ClusterFinalizer) {
		ctrlutil.AddFinalizer(vsphereCluster, infrav1.ClusterFinalizer)
		return reconcile.Result{}, nil
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, clusterContext)
}

func (r *clusterReconciler) reconcileDelete(ctx context.Context, clusterCtx *capvcontext.ClusterContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	vsphereMachines, err := r.vmService.GetMachinesInCluster(ctx, clusterCtx.Cluster.Namespace, clusterCtx.Cluster.Name)
	if err != nil {
		return reconcile.Result{}, pkgerrors.Wrapf(err,
			"unable to list VSphereMachines part of VSphereCluster %s/%s", clusterCtx.VSphereCluster.Namespace, clusterCtx.VSphereCluster.Name)
	}

	if len(vsphereMachines) > 0 {
		log.Info("Waiting for VSphereMachines to be deleted", "count", len(vsphereMachines))
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// The cluster module info needs to be reconciled before the secret deletion
	// since it needs access to the vCenter instance to be able to perform LCM operations
	// on the cluster modules.
	affinityReconcileResult, err := r.reconcileClusterModules(ctx, clusterCtx)
	if err != nil {
		return affinityReconcileResult, err
	}

	// Remove finalizer on Identity Secret
	if identity.IsSecretIdentity(clusterCtx.VSphereCluster) {
		secret := &corev1.Secret{}
		secretKey := client.ObjectKey{
			Namespace: clusterCtx.VSphereCluster.Namespace,
			Name:      clusterCtx.VSphereCluster.Spec.IdentityRef.Name,
		}
		if err := r.Client.Get(ctx, secretKey, secret); err != nil {
			if apierrors.IsNotFound(err) {
				ctrlutil.RemoveFinalizer(clusterCtx.VSphereCluster, infrav1.ClusterFinalizer)
				return reconcile.Result{}, nil
			}
			return reconcile.Result{}, err
		}

		if ctrlutil.RemoveFinalizer(secret, infrav1.SecretIdentitySetFinalizer) {
			log.Info(fmt.Sprintf("Removing finalizer %s", infrav1.SecretIdentitySetFinalizer), "Secret", klog.KObj(secret))
			if err := r.Client.Update(ctx, secret); err != nil {
				return reconcile.Result{}, pkgerrors.Wrapf(err, fmt.Sprintf("failed to update Secret %s", klog.KObj(secret)))
			}
		}

		if secret.DeletionTimestamp.IsZero() {
			log.Info("Deleting Secret", "Secret", klog.KObj(secret))
			if err := r.Client.Delete(ctx, secret); err != nil {
				return reconcile.Result{}, pkgerrors.Wrapf(err, fmt.Sprintf("failed to delete Secret %s", klog.KObj(secret)))
			}
		}
	}

	// Cluster is deleted so remove the finalizer.
	ctrlutil.RemoveFinalizer(clusterCtx.VSphereCluster, infrav1.ClusterFinalizer)

	return reconcile.Result{}, nil
}

func (r *clusterReconciler) reconcileNormal(ctx context.Context, clusterCtx *capvcontext.ClusterContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	ok, err := r.reconcileDeploymentZones(ctx, clusterCtx)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !ok {
		log.Info("Waiting for failure domains to be reconciled")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	if err := r.reconcileIdentitySecret(ctx, clusterCtx); err != nil {
		conditions.MarkFalse(clusterCtx.VSphereCluster, infrav1.VCenterAvailableCondition, infrav1.VCenterUnreachableReason, clusterv1.ConditionSeverityError, err.Error())
		return reconcile.Result{}, err
	}

	vcenterSession, err := r.reconcileVCenterConnectivity(ctx, clusterCtx)
	if err != nil {
		conditions.MarkFalse(clusterCtx.VSphereCluster, infrav1.VCenterAvailableCondition, infrav1.VCenterUnreachableReason, clusterv1.ConditionSeverityError, err.Error())
		return reconcile.Result{}, pkgerrors.Wrapf(err,
			"unexpected error while probing vcenter for %s", clusterCtx)
	}
	conditions.MarkTrue(clusterCtx.VSphereCluster, infrav1.VCenterAvailableCondition)

	err = r.reconcileVCenterVersion(clusterCtx, vcenterSession)
	if err != nil || clusterCtx.VSphereCluster.Status.VCenterVersion == "" {
		conditions.MarkFalse(clusterCtx.VSphereCluster, infrav1.ClusterModulesAvailableCondition, infrav1.MissingVCenterVersionReason, clusterv1.ConditionSeverityWarning, "vCenter version not set")
		log.Error(err, "could not reconcile vCenter version")
	}

	affinityReconcileResult, err := r.reconcileClusterModules(ctx, clusterCtx)
	if err != nil {
		conditions.MarkFalse(clusterCtx.VSphereCluster, infrav1.ClusterModulesAvailableCondition, infrav1.ClusterModuleSetupFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
		return affinityReconcileResult, err
	}

	clusterCtx.VSphereCluster.Status.Ready = true

	return reconcile.Result{}, nil
}

func (r *clusterReconciler) reconcileIdentitySecret(ctx context.Context, clusterCtx *capvcontext.ClusterContext) error {
	vsphereCluster := clusterCtx.VSphereCluster
	if !identity.IsSecretIdentity(vsphereCluster) {
		return nil
	}
	secret := &corev1.Secret{}
	secretKey := client.ObjectKey{
		Namespace: vsphereCluster.Namespace,
		Name:      vsphereCluster.Spec.IdentityRef.Name,
	}
	err := r.Client.Get(ctx, secretKey, secret)
	if err != nil {
		return err
	}

	// If a different VSphereCluster is an owner return an error.
	if !clusterutilv1.IsOwnedByObject(secret, vsphereCluster) && identity.IsOwnedByIdentityOrCluster(secret.GetOwnerReferences()) {
		return fmt.Errorf("another cluster has set the OwnerRef for Secret %s/%s", secret.Namespace, secret.Name)
	}

	helper, err := patch.NewHelper(secret, r.Client)
	if err != nil {
		return err
	}

	// Ensure the VSphereCluster is an owner and that the APIVersion is up to date.
	secret.SetOwnerReferences(clusterutilv1.EnsureOwnerRef(secret.GetOwnerReferences(),
		metav1.OwnerReference{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       "VSphereCluster",
			Name:       vsphereCluster.Name,
			UID:        vsphereCluster.UID,
		},
	))

	// Ensure the finalizer is added.
	if !ctrlutil.ContainsFinalizer(secret, infrav1.SecretIdentitySetFinalizer) {
		ctrlutil.AddFinalizer(secret, infrav1.SecretIdentitySetFinalizer)
	}
	err = helper.Patch(ctx, secret)
	if err != nil {
		return pkgerrors.Wrapf(err, "Failed to patch secret %s", klog.KObj(secret))
	}

	return nil
}

func (r *clusterReconciler) reconcileVCenterConnectivity(ctx context.Context, clusterCtx *capvcontext.ClusterContext) (*session.Session, error) {
	params := session.NewParams().
		WithServer(clusterCtx.VSphereCluster.Spec.Server).
		WithThumbprint(clusterCtx.VSphereCluster.Spec.Thumbprint).
		WithFeatures(session.Feature{
			EnableKeepAlive:   r.ControllerManagerContext.EnableKeepAlive,
			KeepAliveDuration: r.ControllerManagerContext.KeepAliveDuration,
		})

	if clusterCtx.VSphereCluster.Spec.IdentityRef != nil {
		creds, err := identity.GetCredentials(ctx, r.Client, clusterCtx.VSphereCluster, r.ControllerManagerContext.Namespace)
		if err != nil {
			return nil, pkgerrors.Wrap(err, "failed to get credentials from IdentityRef")
		}

		params = params.WithUserInfo(creds.Username, creds.Password)
		return session.GetOrCreate(ctx, params)
	}

	params = params.WithUserInfo(r.ControllerManagerContext.Username, r.ControllerManagerContext.Password)
	return session.GetOrCreate(ctx, params)
}

func (r *clusterReconciler) reconcileVCenterVersion(clusterCtx *capvcontext.ClusterContext, s *session.Session) error {
	version, err := s.GetVersion()
	if err != nil {
		return err
	}
	clusterCtx.VSphereCluster.Status.VCenterVersion = version
	return nil
}

func (r *clusterReconciler) reconcileDeploymentZones(ctx context.Context, clusterCtx *capvcontext.ClusterContext) (bool, error) {
	// If there is no failure domain selector, skip reconciliation
	if clusterCtx.VSphereCluster.Spec.FailureDomainSelector == nil {
		return true, nil
	}

	var opts client.ListOptions
	var err error
	opts.LabelSelector, err = metav1.LabelSelectorAsSelector(clusterCtx.VSphereCluster.Spec.FailureDomainSelector)
	if err != nil {
		return false, pkgerrors.Wrapf(err, "zone label selector is misconfigured")
	}

	var deploymentZoneList infrav1.VSphereDeploymentZoneList
	err = r.Client.List(ctx, &deploymentZoneList, &opts)
	if err != nil {
		return false, pkgerrors.Wrap(err, "unable to list VSphereDeploymentZones")
	}

	readyNotReported, notReady := 0, 0
	failureDomains := clusterv1.FailureDomains{}
	for _, zone := range deploymentZoneList.Items {
		if zone.Spec.Server != clusterCtx.VSphereCluster.Spec.Server {
			continue
		}

		if zone.Status.Ready == nil {
			readyNotReported++
			failureDomains[zone.Name] = clusterv1.FailureDomainSpec{
				ControlPlane: ptr.Deref(zone.Spec.ControlPlane, true),
			}
			continue
		}

		if *zone.Status.Ready {
			failureDomains[zone.Name] = clusterv1.FailureDomainSpec{
				ControlPlane: ptr.Deref(zone.Spec.ControlPlane, true),
			}
			continue
		}
		notReady++
	}

	clusterCtx.VSphereCluster.Status.FailureDomains = failureDomains
	if readyNotReported > 0 {
		conditions.MarkFalse(clusterCtx.VSphereCluster, infrav1.FailureDomainsAvailableCondition, infrav1.WaitingForFailureDomainStatusReason, clusterv1.ConditionSeverityInfo, "waiting for failure domains to report ready status")
		return false, nil
	}

	if len(failureDomains) > 0 {
		if notReady > 0 {
			conditions.MarkFalse(clusterCtx.VSphereCluster, infrav1.FailureDomainsAvailableCondition, infrav1.FailureDomainsSkippedReason, clusterv1.ConditionSeverityInfo, "one or more failure domains are not ready")
		} else {
			conditions.MarkTrue(clusterCtx.VSphereCluster, infrav1.FailureDomainsAvailableCondition)
		}
	} else {
		// Remove the condition if failure domains do not exist
		conditions.Delete(clusterCtx.VSphereCluster, infrav1.FailureDomainsAvailableCondition)
	}
	return true, nil
}

func (r *clusterReconciler) reconcileClusterModules(ctx context.Context, clusterCtx *capvcontext.ClusterContext) (reconcile.Result, error) {
	if feature.Gates.Enabled(feature.NodeAntiAffinity) {
		return r.clusterModuleReconciler.Reconcile(ctx, clusterCtx)
	}
	return reconcile.Result{}, nil
}

// controlPlaneMachineToCluster is a handler.ToRequestsFunc to be used
// to enqueue requests for reconciliation for VSphereCluster to update
// its status.apiEndpoints field.
func (r *clusterReconciler) controlPlaneMachineToCluster(ctx context.Context, o client.Object) []ctrl.Request {
	log := ctrl.LoggerFrom(ctx)

	vsphereMachine, ok := o.(*infrav1.VSphereMachine)
	if !ok {
		log.Error(nil, fmt.Sprintf("Expected a VSphereMachine but got a %T", o))
		return nil
	}
	log = log.WithValues("VSphereMachine", klog.KObj(vsphereMachine))
	ctx = ctrl.LoggerInto(ctx, log)

	if !infrautilv1.IsControlPlaneMachine(vsphereMachine) {
		log.V(6).Info("Skipping VSphereCluster reconcile as Machine is not a control plane Machine")
		return nil
	}

	if len(vsphereMachine.Status.Addresses) == 0 {
		log.V(6).Info("Skipping VSphereCluster reconcile as Machine does not have an IP address")
		return nil
	}

	// Get the VSphereMachine's preferred IP address.
	if _, err := infrautilv1.GetMachinePreferredIPAddress(vsphereMachine); err != nil {
		if errors.Is(err, infrautilv1.ErrNoMachineIPAddr) {
			log.V(6).Info("Skipping VSphereCluster reconcile as Machine does not have a preferred IP address")
			return nil
		}
		log.V(4).Error(err, "Failed to get preferred IP address for VSphereMachine")
		return nil
	}

	// Fetch the CAPI Cluster.
	cluster, err := clusterutilv1.GetClusterFromMetadata(ctx, r.Client, vsphereMachine.ObjectMeta)
	if err != nil {
		log.V(4).Error(err, "VSphereMachine is missing cluster label or cluster does not exist")
		return nil
	}
	log = log.WithValues("Cluster", klog.KObj(cluster))
	if cluster.Spec.InfrastructureRef != nil {
		log = log.WithValues("VSphereCluster", klog.KRef(cluster.Namespace, cluster.Spec.InfrastructureRef.Name))
	}
	ctx = ctrl.LoggerInto(ctx, log)

	if conditions.IsTrue(cluster, clusterv1.ControlPlaneInitializedCondition) {
		log.V(6).Info("Skipping VSphereCluster reconcile as control plane is already initialized")
		return nil
	}

	if !cluster.Spec.ControlPlaneEndpoint.IsZero() {
		log.V(6).Info("Skipping VSphereCluster reconcile as Cluster control plane endpoint is already set")
		return nil
	}

	if cluster.Spec.InfrastructureRef == nil {
		log.Error(nil, "Failed to get VSphereCluster: Cluster.spec.infrastructureRef is not yet set")
		return nil
	}

	// Fetch the VSphereCluster
	vsphereCluster := &infrav1.VSphereCluster{}
	vsphereClusterKey := client.ObjectKey{
		Namespace: vsphereMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, vsphereClusterKey, vsphereCluster); err != nil {
		log.V(4).Error(err, "Failed to get VSphereCluster")
		return nil
	}

	if !vsphereCluster.Spec.ControlPlaneEndpoint.IsZero() {
		log.V(6).Info("Skipping VSphereCluster reconcile as VSphereCluster control plane endpoint is already set")
		return nil
	}

	return []ctrl.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: vsphereClusterKey.Namespace,
			Name:      vsphereClusterKey.Name,
		},
	}}
}

func (r *clusterReconciler) deploymentZoneToCluster(ctx context.Context, o client.Object) []ctrl.Request {
	log := ctrl.LoggerFrom(ctx)

	var requests []ctrl.Request
	obj, ok := o.(*infrav1.VSphereDeploymentZone)
	if !ok {
		log.Error(nil, fmt.Sprintf("Expected a VSphereDeploymentZone but got a %T", o))
		return nil
	}

	var clusterList infrav1.VSphereClusterList
	err := r.Client.List(ctx, &clusterList)
	if err != nil {
		log.V(4).Error(err, "Failed to list VSphereClusters")
		return requests
	}

	for _, cluster := range clusterList.Items {
		if obj.Spec.Server == cluster.Spec.Server {
			r := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      cluster.Name,
					Namespace: cluster.Namespace,
				},
			}
			requests = append(requests, r)
		}
	}
	return requests
}
