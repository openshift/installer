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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/collections"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/identity"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheredeploymentzones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheredeploymentzones/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspherefailuredomains,verbs=get;list;watch;create;update;patch;delete

// AddVSphereDeploymentZoneControllerToManager adds the VSphereDeploymentZone controller to the provided manager.
func AddVSphereDeploymentZoneControllerToManager(ctx context.Context, controllerManagerCtx *capvcontext.ControllerManagerContext, mgr manager.Manager, options controller.Options) error {
	// Build the controller context.
	reconciler := vsphereDeploymentZoneReconciler{
		ControllerManagerContext: controllerManagerCtx,
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.VSphereDeploymentZone{}).
		WithOptions(options).
		Watches(
			&infrav1.VSphereFailureDomain{},
			handler.EnqueueRequestsFromMapFunc(reconciler.failureDomainsToDeploymentZones)).
		// Watch a GenericEvent channel for the controlled resource.
		// This is useful when there are events outside of Kubernetes that
		// should cause a resource to be synchronized, such as a goroutine
		// waiting on some asynchronous, external task to complete.
		WatchesRawSource(
			&source.Channel{Source: controllerManagerCtx.GetGenericEventChannelFor(infrav1.GroupVersion.WithKind("VSphereDeploymentZone"))},
			&handler.EnqueueRequestForObject{},
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), controllerManagerCtx.WatchFilterValue)).
		Complete(reconciler)
}

type vsphereDeploymentZoneReconciler struct {
	*capvcontext.ControllerManagerContext
}

func (r vsphereDeploymentZoneReconciler) Reconcile(ctx context.Context, request reconcile.Request) (_ reconcile.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the VSphereDeploymentZone for this request.
	vsphereDeploymentZone := &infrav1.VSphereDeploymentZone{}
	if err := r.Client.Get(ctx, request.NamespacedName, vsphereDeploymentZone); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	log = log.WithValues("VSphereFailureDomain", klog.KRef("", vsphereDeploymentZone.Spec.FailureDomain))
	ctx = ctrl.LoggerInto(ctx, log)

	if annotations.HasPaused(vsphereDeploymentZone) {
		log.Info("Reconciliation is paused for this object")
		return reconcile.Result{}, nil
	}

	patchHelper, err := patch.NewHelper(vsphereDeploymentZone, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to initialize patch helper")
	}

	vsphereDeploymentZoneContext := &capvcontext.VSphereDeploymentZoneContext{
		ControllerManagerContext: r.ControllerManagerContext,
		VSphereDeploymentZone:    vsphereDeploymentZone,
		PatchHelper:              patchHelper,
	}
	defer func() {
		if err := vsphereDeploymentZoneContext.Patch(ctx); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	if !vsphereDeploymentZone.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, vsphereDeploymentZoneContext)
	}

	// If the VSphereDeploymentZone doesn't have our finalizer, add it.
	// Requeue immediately after adding finalizer to avoid the race condition between init and delete
	if !ctrlutil.ContainsFinalizer(vsphereDeploymentZone, infrav1.DeploymentZoneFinalizer) {
		ctrlutil.AddFinalizer(vsphereDeploymentZone, infrav1.DeploymentZoneFinalizer)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, r.reconcileNormal(ctx, vsphereDeploymentZoneContext)
}

func (r vsphereDeploymentZoneReconciler) reconcileNormal(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext) error {
	failureDomain := &infrav1.VSphereFailureDomain{}
	failureDomainKey := client.ObjectKey{Name: deploymentZoneCtx.VSphereDeploymentZone.Spec.FailureDomain}
	if err := r.Client.Get(ctx, failureDomainKey, failureDomain); err != nil {
		return errors.Wrapf(err, "failed to get VSphereFailureDomain %s", klog.KRef(failureDomainKey.Namespace, failureDomainKey.Name))
	}

	authSession, err := r.getVCenterSession(ctx, deploymentZoneCtx, failureDomain.Spec.Topology.Datacenter)
	if err != nil {
		conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VCenterAvailableCondition, infrav1.VCenterUnreachableReason, clusterv1.ConditionSeverityError, err.Error())
		deploymentZoneCtx.VSphereDeploymentZone.Status.Ready = ptr.To(false)
		return err
	}
	deploymentZoneCtx.AuthSession = authSession
	conditions.MarkTrue(deploymentZoneCtx.VSphereDeploymentZone, infrav1.VCenterAvailableCondition)

	if err := r.reconcilePlacementConstraint(ctx, deploymentZoneCtx); err != nil {
		deploymentZoneCtx.VSphereDeploymentZone.Status.Ready = ptr.To(false)
		return err
	}
	conditions.MarkTrue(deploymentZoneCtx.VSphereDeploymentZone, infrav1.PlacementConstraintMetCondition)

	// reconcile the failure domain
	if err := r.reconcileFailureDomain(ctx, deploymentZoneCtx, failureDomain); err != nil {
		deploymentZoneCtx.VSphereDeploymentZone.Status.Ready = ptr.To(false)
		return err
	}

	// Mark the deployment zone as ready.
	deploymentZoneCtx.VSphereDeploymentZone.Status.Ready = ptr.To(true)
	return nil
}

func (r vsphereDeploymentZoneReconciler) reconcilePlacementConstraint(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext) error {
	placementConstraint := deploymentZoneCtx.VSphereDeploymentZone.Spec.PlacementConstraint

	if resourcePool := placementConstraint.ResourcePool; resourcePool != "" {
		if _, err := deploymentZoneCtx.AuthSession.Finder.ResourcePool(ctx, resourcePool); err != nil {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.PlacementConstraintMetCondition, infrav1.ResourcePoolNotFoundReason, clusterv1.ConditionSeverityError, "resource pool %s is misconfigured", resourcePool)
			return errors.Wrapf(err, "failed to reconcile placement contraint: unable to find resource pool %s", resourcePool)
		}
	}

	if folder := placementConstraint.Folder; folder != "" {
		if _, err := deploymentZoneCtx.AuthSession.Finder.Folder(ctx, placementConstraint.Folder); err != nil {
			conditions.MarkFalse(deploymentZoneCtx.VSphereDeploymentZone, infrav1.PlacementConstraintMetCondition, infrav1.FolderNotFoundReason, clusterv1.ConditionSeverityError, "folder %s is misconfigured", folder)
			return errors.Wrapf(err, "failed to reconcile placement contraint: unable to find folder %s", folder)
		}
	}
	return nil
}

func (r vsphereDeploymentZoneReconciler) getVCenterSession(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext, datacenter string) (*session.Session, error) {
	log := ctrl.LoggerFrom(ctx)

	params := session.NewParams().
		WithServer(deploymentZoneCtx.VSphereDeploymentZone.Spec.Server).
		WithDatacenter(datacenter).
		WithUserInfo(r.ControllerManagerContext.Username, r.ControllerManagerContext.Password).
		WithFeatures(session.Feature{
			EnableKeepAlive:   r.EnableKeepAlive,
			KeepAliveDuration: r.KeepAliveDuration,
		})

	clusterList := &infrav1.VSphereClusterList{}
	if err := r.Client.List(ctx, clusterList); err != nil {
		return nil, errors.Wrapf(err, "failed to list VSphereClusters")
	}

	for _, vsphereCluster := range clusterList.Items {
		if deploymentZoneCtx.VSphereDeploymentZone.Spec.Server != vsphereCluster.Spec.Server || vsphereCluster.Spec.IdentityRef == nil {
			continue
		}

		// Note: We have to use := here to not overwrite log & ctx outside the for loop.
		log := log.WithValues("VSphereCluster", klog.KRef(vsphereCluster.Namespace, vsphereCluster.Name))
		ctx := ctrl.LoggerInto(ctx, log)

		params = params.WithThumbprint(vsphereCluster.Spec.Thumbprint)
		vsphereCluster := vsphereCluster
		creds, err := identity.GetCredentials(ctx, r.Client, &vsphereCluster, r.Namespace)
		if err != nil {
			log.Error(err, "error retrieving credentials from IdentityRef")
			continue
		}
		log.V(4).Info("Using credentials from VSphereCluster IdentityRef to create the authenticated session")
		params = params.WithUserInfo(creds.Username, creds.Password)
		return session.GetOrCreate(ctx, params)
	}

	// Fallback to using credentials provided to the manager
	log.V(4).Info("Using credentials provided to the manager to create the authenticated session")
	return session.GetOrCreate(ctx, params)
}

func (r vsphereDeploymentZoneReconciler) reconcileDelete(ctx context.Context, deploymentZoneCtx *capvcontext.VSphereDeploymentZoneContext) error {
	log := ctrl.LoggerFrom(ctx)

	machines := &clusterv1.MachineList{}
	if err := r.Client.List(ctx, machines); err != nil {
		return errors.Wrapf(err, "failed to list Machines")
	}

	machinesUsingDeploymentZone := collections.FromMachineList(machines).Filter(collections.ActiveMachines, func(machine *clusterv1.Machine) bool {
		if machine.Spec.FailureDomain != nil {
			return *machine.Spec.FailureDomain == deploymentZoneCtx.VSphereDeploymentZone.Name
		}
		return false
	})
	if len(machinesUsingDeploymentZone) > 0 {
		machineNamesStr := util.MachinesAsString(machinesUsingDeploymentZone.SortedByCreationTimestamp())
		return errors.Errorf("blocking VSphereDeploymentZone deletion: currently in use by Machines %s", machineNamesStr)
	}

	failureDomain := &infrav1.VSphereFailureDomain{}
	failureDomainKey := client.ObjectKey{Name: deploymentZoneCtx.VSphereDeploymentZone.Spec.FailureDomain}
	// Return an error if the FailureDomain can not be retrieved.
	if err := r.Client.Get(ctx, failureDomainKey, failureDomain); err != nil {
		// If the VSphereFailureDomain is not found return early and remove the finalizer.
		// This prevents early deletion of the VSphereFailureDomain from blocking VSphereDeploymentZone deletion.
		if apierrors.IsNotFound(err) {
			ctrlutil.RemoveFinalizer(deploymentZoneCtx.VSphereDeploymentZone, infrav1.DeploymentZoneFinalizer)
			return nil
		}
		return errors.Wrapf(err, "failed to get VSphereFailureDomain")
	}

	// Reconcile the deletion of the VSphereFailureDomain by removing ownerReferences and deleting if necessary.
	if err := updateOwnerReferences(ctx, failureDomain, r.Client, func() []metav1.OwnerReference {
		return clusterutilv1.RemoveOwnerRef(failureDomain.OwnerReferences, metav1.OwnerReference{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       "VSphereDeploymentZone",
			Name:       deploymentZoneCtx.VSphereDeploymentZone.Name,
		})
	}); err != nil {
		return err
	}

	if len(failureDomain.OwnerReferences) == 0 && failureDomain.DeletionTimestamp.IsZero() {
		log.Info("Deleting VSphereFailureDomain")
		if err := r.Client.Delete(ctx, failureDomain); err != nil && !apierrors.IsNotFound(err) {
			return errors.Wrapf(err, "failed to delete VSphereFailureDomain %s", failureDomain.Name)
		}
	}

	ctrlutil.RemoveFinalizer(deploymentZoneCtx.VSphereDeploymentZone, infrav1.DeploymentZoneFinalizer)
	return nil
}

// updateOwnerReferences uses the ownerRef function to calculate the owner references
// to be set on the object and patches the object.
func updateOwnerReferences(ctx context.Context, obj client.Object, client client.Client, ownerRefFunc func() []metav1.OwnerReference) error {
	patchHelper, err := patch.NewHelper(obj, client)
	if err != nil {
		return errors.Wrapf(err, "failed to init patch helper for %s %s",
			obj.GetObjectKind(),
			obj.GetName())
	}

	obj.SetOwnerReferences(ownerRefFunc())
	if err := patchHelper.Patch(ctx, obj); err != nil {
		return errors.Wrapf(err, "failed to patch object %s %s",
			obj.GetObjectKind(),
			obj.GetName())
	}
	return nil
}

func (r vsphereDeploymentZoneReconciler) failureDomainsToDeploymentZones(ctx context.Context, a client.Object) []reconcile.Request {
	log := ctrl.LoggerFrom(ctx)

	failureDomain, ok := a.(*infrav1.VSphereFailureDomain)
	if !ok {
		log.Error(nil, fmt.Sprintf("Expected a VSphereFailureDomain but got a %T", a))
		return nil
	}

	var zones infrav1.VSphereDeploymentZoneList
	if err := r.Client.List(ctx, &zones); err != nil {
		log.V(4).Error(err, "Failed to list VSphereDeploymentZones")
		return nil
	}

	var requests []reconcile.Request
	for _, zone := range zones.Items {
		if zone.Spec.FailureDomain == failureDomain.Name {
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name: zone.Name,
				},
			})
		}
	}
	return requests
}
