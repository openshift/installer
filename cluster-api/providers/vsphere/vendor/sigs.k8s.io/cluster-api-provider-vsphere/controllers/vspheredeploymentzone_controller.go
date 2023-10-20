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
	"reflect"
	"strings"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
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
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/identity"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheredeploymentzones,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheredeploymentzones/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspherefailuredomains,verbs=get;list;watch;create;update;patch;delete

// AddVSphereDeploymentZoneControllerToManager adds the VSphereDeploymentZone controller to the provided manager.
func AddVSphereDeploymentZoneControllerToManager(ctx *context.ControllerManagerContext, mgr manager.Manager, options controller.Options) error {
	var (
		controlledType     = &infrav1.VSphereDeploymentZone{}
		controlledTypeName = reflect.TypeOf(controlledType).Elem().Name()
		controlledTypeGVK  = infrav1.GroupVersion.WithKind(controlledTypeName)

		controllerNameShort = fmt.Sprintf("%s-controller", strings.ToLower(controlledTypeName))
		controllerNameLong  = fmt.Sprintf("%s/%s/%s", ctx.Namespace, ctx.Name, controllerNameShort)
	)

	// Build the controller context.
	controllerContext := &context.ControllerContext{
		ControllerManagerContext: ctx,
		Name:                     controllerNameShort,
		Recorder:                 record.New(mgr.GetEventRecorderFor(controllerNameLong)),
		Logger:                   ctx.Logger.WithName(controllerNameShort),
	}
	reconciler := vsphereDeploymentZoneReconciler{ControllerContext: controllerContext}

	return ctrl.NewControllerManagedBy(mgr).
		// Watch the controlled, infrastructure resource.
		For(controlledType).
		WithOptions(options).
		Watches(
			&infrav1.VSphereFailureDomain{},
			handler.EnqueueRequestsFromMapFunc(reconciler.failureDomainsToDeploymentZones)).
		// Watch a GenericEvent channel for the controlled resource.
		// This is useful when there are events outside of Kubernetes that
		// should cause a resource to be synchronized, such as a goroutine
		// waiting on some asynchronous, external task to complete.
		WatchesRawSource(
			&source.Channel{Source: ctx.GetGenericEventChannelFor(controlledTypeGVK)},
			&handler.EnqueueRequestForObject{},
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), ctx.WatchFilterValue)).
		Complete(reconciler)
}

type vsphereDeploymentZoneReconciler struct {
	*context.ControllerContext
}

func (r vsphereDeploymentZoneReconciler) Reconcile(ctx goctx.Context, request reconcile.Request) (_ reconcile.Result, reterr error) {
	logr := r.Logger.WithValues("vspheredeploymentzone", request.Name)
	// Fetch the VSphereDeploymentZone for this request.
	vsphereDeploymentZone := &infrav1.VSphereDeploymentZone{}
	if err := r.Client.Get(ctx, request.NamespacedName, vsphereDeploymentZone); err != nil {
		if apierrors.IsNotFound(err) {
			logr.V(4).Info("VSphereDeploymentZone not found, won't reconcile", "key", request.NamespacedName)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	failureDomain := &infrav1.VSphereFailureDomain{}
	failureDomainKey := client.ObjectKey{Name: vsphereDeploymentZone.Spec.FailureDomain}
	if err := r.Client.Get(ctx, failureDomainKey, failureDomain); err != nil {
		if apierrors.IsNotFound(err) {
			logr.V(4).Info("Failure Domain not found, won't reconcile", "key", failureDomainKey)
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	patchHelper, err := patch.NewHelper(vsphereDeploymentZone, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(
			err,
			"failed to init patch helper for %s %s/%s",
			vsphereDeploymentZone.GroupVersionKind(),
			vsphereDeploymentZone.Namespace,
			vsphereDeploymentZone.Name)
	}

	vsphereDeploymentZoneContext := &context.VSphereDeploymentZoneContext{
		ControllerContext:     r.ControllerContext,
		VSphereDeploymentZone: vsphereDeploymentZone,
		VSphereFailureDomain:  failureDomain,
		Logger:                logr,
		PatchHelper:           patchHelper,
	}
	defer func() {
		if err := vsphereDeploymentZoneContext.Patch(); err != nil {
			if reterr == nil {
				reterr = err
			}
			logr.Error(err, "patch failed", "vsphereDeploymentZone", vsphereDeploymentZoneContext.String())
		}
	}()

	if !vsphereDeploymentZone.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(vsphereDeploymentZoneContext)
	}

	// If the VSphereDeploymentZone doesn't have our finalizer, add it.
	// Requeue immediately after adding finalizer to avoid the race condition between init and delete
	if !ctrlutil.ContainsFinalizer(vsphereDeploymentZone, infrav1.DeploymentZoneFinalizer) {
		ctrlutil.AddFinalizer(vsphereDeploymentZone, infrav1.DeploymentZoneFinalizer)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, r.reconcileNormal(vsphereDeploymentZoneContext)
}

func (r vsphereDeploymentZoneReconciler) reconcileNormal(ctx *context.VSphereDeploymentZoneContext) error {
	authSession, err := r.getVCenterSession(ctx)
	if err != nil {
		ctx.Logger.V(4).Error(err, "unable to create session")
		conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.VCenterAvailableCondition, infrav1.VCenterUnreachableReason, clusterv1.ConditionSeverityError, err.Error())
		ctx.VSphereDeploymentZone.Status.Ready = pointer.Bool(false)
		return errors.Wrapf(err, "unable to create auth session")
	}
	ctx.AuthSession = authSession
	conditions.MarkTrue(ctx.VSphereDeploymentZone, infrav1.VCenterAvailableCondition)

	if err := r.reconcilePlacementConstraint(ctx); err != nil {
		ctx.VSphereDeploymentZone.Status.Ready = pointer.Bool(false)
		return errors.Wrap(err, "placement constraint is misconfigured")
	}
	conditions.MarkTrue(ctx.VSphereDeploymentZone, infrav1.PlacementConstraintMetCondition)

	// reconcile the failure domain
	if err := r.reconcileFailureDomain(ctx); err != nil {
		ctx.Logger.V(4).Error(err, "failed to reconcile failure domain", "failureDomain", ctx.VSphereDeploymentZone.Spec.FailureDomain)
		ctx.VSphereDeploymentZone.Status.Ready = pointer.Bool(false)
		return errors.Wrapf(err, "failed to reconcile failure domain")
	}
	conditions.MarkTrue(ctx.VSphereDeploymentZone, infrav1.VSphereFailureDomainValidatedCondition)

	// Mark the deployment zone as ready.
	ctx.VSphereDeploymentZone.Status.Ready = pointer.Bool(true)
	return nil
}

func (r vsphereDeploymentZoneReconciler) reconcilePlacementConstraint(ctx *context.VSphereDeploymentZoneContext) error {
	placementConstraint := ctx.VSphereDeploymentZone.Spec.PlacementConstraint

	if resourcePool := placementConstraint.ResourcePool; resourcePool != "" {
		if _, err := ctx.AuthSession.Finder.ResourcePool(ctx, resourcePool); err != nil {
			ctx.Logger.V(4).Error(err, "unable to find resource pool", "name", resourcePool)
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.PlacementConstraintMetCondition, infrav1.ResourcePoolNotFoundReason, clusterv1.ConditionSeverityError, "resource pool %s is misconfigured", resourcePool)
			return errors.Wrapf(err, "unable to find resource pool %s", resourcePool)
		}
	}

	if folder := placementConstraint.Folder; folder != "" {
		if _, err := ctx.AuthSession.Finder.Folder(ctx, placementConstraint.Folder); err != nil {
			ctx.Logger.V(4).Error(err, "unable to find folder", "name", folder)
			conditions.MarkFalse(ctx.VSphereDeploymentZone, infrav1.PlacementConstraintMetCondition, infrav1.FolderNotFoundReason, clusterv1.ConditionSeverityError, "datastore %s is misconfigured", folder)
			return errors.Wrapf(err, "unable to find folder %s", folder)
		}
	}
	return nil
}

func (r vsphereDeploymentZoneReconciler) getVCenterSession(ctx *context.VSphereDeploymentZoneContext) (*session.Session, error) {
	params := session.NewParams().
		WithServer(ctx.VSphereDeploymentZone.Spec.Server).
		WithDatacenter(ctx.VSphereFailureDomain.Spec.Topology.Datacenter).
		WithUserInfo(r.ControllerContext.Username, r.ControllerContext.Password).
		WithFeatures(session.Feature{
			EnableKeepAlive:   r.EnableKeepAlive,
			KeepAliveDuration: r.KeepAliveDuration,
		})

	clusterList := &infrav1.VSphereClusterList{}
	if err := r.Client.List(ctx, clusterList); err != nil {
		return nil, err
	}

	for _, vsphereCluster := range clusterList.Items {
		if ctx.VSphereDeploymentZone.Spec.Server == vsphereCluster.Spec.Server && vsphereCluster.Spec.IdentityRef != nil {
			logger := ctx.Logger.WithValues("cluster", vsphereCluster.Name)
			params = params.WithThumbprint(vsphereCluster.Spec.Thumbprint)
			clust := vsphereCluster
			creds, err := identity.GetCredentials(ctx, r.Client, &clust, r.Namespace)
			if err != nil {
				logger.Error(err, "error retrieving credentials from IdentityRef")
				continue
			}
			logger.Info("using server credentials to create the authenticated session")
			params = params.WithUserInfo(creds.Username, creds.Password)
			return session.GetOrCreate(r.Context,
				params)
		}
	}

	// Fallback to using credentials provided to the manager
	return session.GetOrCreate(r.Context,
		params)
}

func (r vsphereDeploymentZoneReconciler) reconcileDelete(ctx *context.VSphereDeploymentZoneContext) error {
	r.Logger.Info("Deleting VSphereDeploymentZone")

	machines := &clusterv1.MachineList{}
	if err := r.Client.List(ctx, machines); err != nil {
		r.Logger.Error(err, "unable to list machines")
		return errors.Wrapf(err, "unable to list machines")
	}

	machinesUsingDeploymentZone := collections.FromMachineList(machines).Filter(collections.ActiveMachines, func(machine *clusterv1.Machine) bool {
		if machine.Spec.FailureDomain != nil {
			return *machine.Spec.FailureDomain == ctx.VSphereDeploymentZone.Name
		}
		return false
	})
	if len(machinesUsingDeploymentZone) > 0 {
		machineNamesStr := util.MachinesAsString(machinesUsingDeploymentZone.SortedByCreationTimestamp())
		err := errors.Errorf("%s is currently in use by machines: %s", ctx.VSphereDeploymentZone.Name, machineNamesStr)
		r.Logger.Error(err, "Error deleting VSphereDeploymentZone", "name", ctx.VSphereDeploymentZone.Name)
		return err
	}

	if err := updateOwnerReferences(ctx, ctx.VSphereFailureDomain, r.Client, func() []metav1.OwnerReference {
		return clusterutilv1.RemoveOwnerRef(ctx.VSphereFailureDomain.OwnerReferences, metav1.OwnerReference{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       ctx.VSphereDeploymentZone.Kind,
			Name:       ctx.VSphereDeploymentZone.Name,
		})
	}); err != nil {
		return err
	}

	if len(ctx.VSphereFailureDomain.OwnerReferences) == 0 {
		ctx.Logger.Info("deleting vsphereFailureDomain", "name", ctx.VSphereFailureDomain.Name)
		if err := r.Client.Delete(ctx, ctx.VSphereFailureDomain); err != nil && !apierrors.IsNotFound(err) {
			ctx.Logger.Error(err, "failed to delete related %s %s", ctx.VSphereFailureDomain.GroupVersionKind(), ctx.VSphereFailureDomain.Name)
		}
	}

	ctrlutil.RemoveFinalizer(ctx.VSphereDeploymentZone, infrav1.DeploymentZoneFinalizer)

	return nil
}

// updateOwnerReferences uses the ownerRef function to calculate the owner references
// to be set on the object and patches the object.
func updateOwnerReferences(ctx goctx.Context, obj client.Object, client client.Client, ownerRefFunc func() []metav1.OwnerReference) error {
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

func (r vsphereDeploymentZoneReconciler) failureDomainsToDeploymentZones(ctx goctx.Context, a client.Object) []reconcile.Request {
	failureDomain, ok := a.(*infrav1.VSphereFailureDomain)
	if !ok {
		r.Logger.Error(nil, fmt.Sprintf("expected a VSphereFailureDomain but got a %T", a))
		return nil
	}

	var zones infrav1.VSphereDeploymentZoneList
	if err := r.Client.List(ctx, &zones); err != nil {
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
