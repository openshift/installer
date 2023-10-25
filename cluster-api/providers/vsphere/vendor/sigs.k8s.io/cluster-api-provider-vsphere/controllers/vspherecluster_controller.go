/*
Copyright 2019 The Kubernetes Authors.

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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/controllers/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	inframanager "sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/vmoperator"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;patch;update
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vsphereclusteridentities,verbs=get;list;watch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vsphereclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vsphereclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=topology.tanzu.vmware.com,resources=availabilityzones,verbs=get;list;watch
// +kubebuilder:rbac:groups=topology.tanzu.vmware.com,resources=availabilityzones/status,verbs=get;list;watch

// AddClusterControllerToManager adds the cluster controller to the provided
// manager.
func AddClusterControllerToManager(ctx *context.ControllerManagerContext, mgr manager.Manager, clusterControlledType client.Object, options controller.Options) error {
	supervisorBased, err := util.IsSupervisorType(clusterControlledType)
	if err != nil {
		return err
	}

	var (
		clusterControlledTypeName = reflect.TypeOf(clusterControlledType).Elem().Name()
		clusterControlledTypeGVK  = infrav1.GroupVersion.WithKind(clusterControlledTypeName)
		controllerNameShort       = fmt.Sprintf("%s-controller", strings.ToLower(clusterControlledTypeName))
		controllerNameLong        = fmt.Sprintf("%s/%s/%s", ctx.Namespace, ctx.Name, controllerNameShort)
	)
	if supervisorBased {
		clusterControlledTypeGVK = vmwarev1.GroupVersion.WithKind(clusterControlledTypeName)
		controllerNameShort = fmt.Sprintf("%s-supervisor-controller", strings.ToLower(clusterControlledTypeName))
		controllerNameLong = fmt.Sprintf("%s/%s/%s", ctx.Namespace, ctx.Name, controllerNameShort)
	}

	// Build the controller context.
	controllerContext := &context.ControllerContext{
		ControllerManagerContext: ctx,
		Name:                     controllerNameShort,
		Recorder:                 record.New(mgr.GetEventRecorderFor(controllerNameLong)),
		Logger:                   ctx.Logger.WithName(controllerNameShort),
	}

	if supervisorBased {
		networkProvider, err := inframanager.GetNetworkProvider(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to create a network provider")
		}
		reconciler := vmware.ClusterReconciler{
			ControllerContext:     controllerContext,
			ResourcePolicyService: vmoperator.RPService{},
			ControlPlaneService:   vmoperator.CPService{},
			NetworkProvider:       networkProvider,
		}
		return ctrl.NewControllerManagedBy(mgr).
			Named(controllerNameShort).
			For(clusterControlledType).
			WithOptions(options).
			Watches(
				&vmwarev1.VSphereMachine{},
				handler.EnqueueRequestsFromMapFunc(reconciler.VSphereMachineToCluster),
			).
			WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), ctx.WatchFilterValue)).
			Complete(reconciler)
	}

	reconciler := clusterReconciler{
		ControllerContext:       controllerContext,
		clusterModuleReconciler: NewReconciler(controllerContext),
	}
	clusterToInfraFn := clusterToInfrastructureMapFunc(ctx)
	c, err := ctrl.NewControllerManagedBy(mgr).
		// Watch the controlled, infrastructure resource.
		For(clusterControlledType).
		WithOptions(options).
		// Watch the CAPI resource that owns this infrastructure resource.
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(func(ctx goctx.Context, o client.Object) []reconcile.Request {
				requests := clusterToInfraFn(ctx, o)
				if requests == nil {
					return nil
				}

				c := &infrav1.VSphereCluster{}
				if err := reconciler.Client.Get(ctx, requests[0].NamespacedName, c); err != nil {
					reconciler.Logger.V(4).Error(err, "Failed to get VSphereCluster")
					return nil
				}

				if annotations.IsExternallyManaged(c) {
					reconciler.Logger.V(4).Info("VSphereCluster is externally managed, skipping mapping.")
					return nil
				}
				return requests
			}),
		).

		// Watch the infrastructure machine resources that belong to the control
		// plane. This controller needs to reconcile the infrastructure cluster
		// once a control plane machine has an IP address.
		Watches(
			&infrav1.VSphereMachine{},
			handler.EnqueueRequestsFromMapFunc(reconciler.controlPlaneMachineToCluster),
		).
		// Watch the Vsphere deployment zone with the Server field matching the
		// server field of the VSphereCluster.
		Watches(
			&infrav1.VSphereDeploymentZone{},
			handler.EnqueueRequestsFromMapFunc(reconciler.deploymentZoneToCluster),
		).
		// Watch a GenericEvent channel for the controlled resource.
		//
		// This is useful when there are events outside of Kubernetes that
		// should cause a resource to be synchronized, such as a goroutine
		// waiting on some asynchronous, external task to complete.
		WatchesRawSource(
			&source.Channel{Source: ctx.GetGenericEventChannelFor(clusterControlledTypeGVK)},
			&handler.EnqueueRequestForObject{},
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), ctx.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(reconciler.Logger)).
		Build(reconciler)
	if err != nil {
		return err
	}

	if feature.Gates.Enabled(feature.NodeAntiAffinity) {
		return reconciler.clusterModuleReconciler.PopulateWatchesOnController(mgr, c)
	}
	return nil
}

func clusterToInfrastructureMapFunc(managerContext *context.ControllerManagerContext) handler.MapFunc {
	gvk := infrav1.GroupVersion.WithKind(reflect.TypeOf(&infrav1.VSphereCluster{}).Elem().Name())
	return clusterutilv1.ClusterToInfrastructureMapFunc(managerContext, gvk, managerContext.Client, &infrav1.VSphereCluster{})
}
