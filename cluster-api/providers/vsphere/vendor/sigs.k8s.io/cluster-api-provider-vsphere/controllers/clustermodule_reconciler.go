/*
Copyright 2022 The Kubernetes Authors.

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
	"strings"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/clustermodule"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
)

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=kubeadmcontrolplanes,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vsphereclusters,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,verbs=get;list;watch

type Reconciler struct {
	*context.ControllerContext

	ClusterModuleService clustermodule.Service
}

func NewReconciler(ctx *context.ControllerContext) Reconciler {
	return Reconciler{
		ControllerContext:    ctx,
		ClusterModuleService: clustermodule.NewService(),
	}
}

func (r Reconciler) Reconcile(ctx *context.ClusterContext) (reconcile.Result, error) {
	ctx.Logger.Info("reconcile anti affinity setup")
	if !clustermodule.IsClusterCompatible(ctx) {
		conditions.MarkFalse(ctx.VSphereCluster, infrav1.ClusterModulesAvailableCondition, infrav1.VCenterVersionIncompatibleReason, clusterv1.ConditionSeverityInfo,
			"vCenter API version %s is not compatible with cluster modules", ctx.VSphereCluster.Status.VCenterVersion)
		ctx.Logger.Info("cluster is not compatible for anti affinity",
			"api version", ctx.VSphereCluster.Status.VCenterVersion)
		return reconcile.Result{}, nil
	}

	objectMap, err := r.fetchMachineOwnerObjects(ctx)
	if err != nil {
		return reconcile.Result{}, err
	}

	modErrs := []clusterModError{}

	clusterModuleSpecs := []infrav1.ClusterModule{}
	for _, mod := range ctx.VSphereCluster.Spec.ClusterModules {
		curr := mod.TargetObjectName
		if mod.ControlPlane {
			curr = appendKCPKey(curr)
		}
		if obj, ok := objectMap[curr]; !ok {
			// delete the cluster module as the object is marked for deletion
			// or already deleted.
			if err := r.ClusterModuleService.Remove(ctx, mod.ModuleUUID); err != nil {
				ctx.Logger.Error(err, "failed to delete cluster module for object",
					"name", mod.TargetObjectName, "moduleUUID", mod.ModuleUUID)
			}
			delete(objectMap, curr)
		} else {
			// verify the cluster module
			exists, err := r.ClusterModuleService.DoesExist(ctx, obj, mod.ModuleUUID)
			if err != nil {
				// Add the error to modErrs so it gets handled below.
				modErrs = append(modErrs, clusterModError{obj.GetName(), errors.Wrapf(err, "failed to verify cluster module %q", mod.ModuleUUID)})
				ctx.Logger.Error(err, "failed to verify cluster module for object",
					"name", mod.TargetObjectName, "moduleUUID", mod.ModuleUUID)
				// Append the module and remove it from objectMap to not create new ones instead.
				clusterModuleSpecs = append(clusterModuleSpecs, infrav1.ClusterModule{
					ControlPlane:     obj.IsControlPlane(),
					TargetObjectName: obj.GetName(),
					ModuleUUID:       mod.ModuleUUID,
				})
				delete(objectMap, curr)
				continue
			}

			// append the module and object info to the VSphereCluster object
			// and remove it from the object map since no new cluster module
			// needs to be created.
			if exists {
				clusterModuleSpecs = append(clusterModuleSpecs, infrav1.ClusterModule{
					ControlPlane:     obj.IsControlPlane(),
					TargetObjectName: obj.GetName(),
					ModuleUUID:       mod.ModuleUUID,
				})
				delete(objectMap, curr)
			} else {
				ctx.Logger.Info("module for object not found",
					"moduleUUID", mod.ModuleUUID,
					"object", mod.TargetObjectName)
			}
		}
	}

	for _, obj := range objectMap {
		moduleUUID, err := r.ClusterModuleService.Create(ctx, obj)
		if err != nil {
			ctx.Logger.Error(err, "failed to create cluster module for target object", "name", obj.GetName())
			modErrs = append(modErrs, clusterModError{obj.GetName(), err})
			continue
		}
		// module creation was skipped
		if moduleUUID == "" {
			continue
		}
		clusterModuleSpecs = append(clusterModuleSpecs, infrav1.ClusterModule{
			ControlPlane:     obj.IsControlPlane(),
			TargetObjectName: obj.GetName(),
			ModuleUUID:       moduleUUID,
		})
	}
	ctx.VSphereCluster.Spec.ClusterModules = clusterModuleSpecs

	switch {
	case len(modErrs) > 0:
		incompatibleOwnerErrs := incompatibleOwnerErrors(modErrs)
		// if cluster module creation is not possible due to incompatibility,
		// cluster creation should succeed with a warning condition
		if len(incompatibleOwnerErrs) > 0 && len(incompatibleOwnerErrs) == len(modErrs) {
			err = nil
		} else {
			err = errors.New(generateClusterModuleErrorMessage(modErrs))
		}
		conditions.MarkFalse(ctx.VSphereCluster, infrav1.ClusterModulesAvailableCondition, infrav1.ClusterModuleSetupFailedReason,
			clusterv1.ConditionSeverityWarning, generateClusterModuleErrorMessage(modErrs))
	case len(modErrs) == 0 && len(clusterModuleSpecs) > 0:
		conditions.MarkTrue(ctx.VSphereCluster, infrav1.ClusterModulesAvailableCondition)
	default:
		conditions.Delete(ctx.VSphereCluster, infrav1.ClusterModulesAvailableCondition)
	}
	return reconcile.Result{}, err
}

func (r Reconciler) toAffinityInput(ctx goctx.Context, obj client.Object) []reconcile.Request {
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, metav1.ObjectMeta{
		Namespace:       obj.GetNamespace(),
		Labels:          obj.GetLabels(),
		OwnerReferences: obj.GetOwnerReferences(),
	})
	if err != nil {
		r.Logger.Error(err, "failed to get owner cluster")
		return nil
	}

	if cluster.Spec.InfrastructureRef == nil {
		r.Logger.Error(err, "cluster infrastructureRef not set. Requeing",
			"namespace", cluster.Namespace, "cluster_name", cluster.Name)
		return nil
	}
	vsphereCluster := &infrav1.VSphereCluster{}
	if err := r.Client.Get(ctx, client.ObjectKey{
		Name:      cluster.Spec.InfrastructureRef.Name,
		Namespace: cluster.Namespace,
	}, vsphereCluster); err != nil {
		r.Logger.Error(err, "failed to get vSphereCluster object",
			"namespace", cluster.Namespace, "name", cluster.Spec.InfrastructureRef.Name)
		return nil
	}

	return []reconcile.Request{
		{NamespacedName: client.ObjectKeyFromObject(vsphereCluster)},
	}
}

func (r Reconciler) PopulateWatchesOnController(mgr manager.Manager, controller controller.Controller) error {
	if err := controller.Watch(
		source.Kind(mgr.GetCache(), &controlplanev1.KubeadmControlPlane{}),
		handler.EnqueueRequestsFromMapFunc(r.toAffinityInput),
		predicate.Funcs{
			GenericFunc: func(genericEvent event.GenericEvent) bool {
				return false
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				return false
			},
		},
	); err != nil {
		return err
	}

	return controller.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.MachineDeployment{}),
		handler.EnqueueRequestsFromMapFunc(r.toAffinityInput),
		predicate.Funcs{
			GenericFunc: func(genericEvent event.GenericEvent) bool {
				return false
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				return false
			},
		},
	)
}

func (r Reconciler) fetchMachineOwnerObjects(ctx *context.ClusterContext) (map[string]clustermodule.Wrapper, error) {
	objects := map[string]clustermodule.Wrapper{}

	name, ok := ctx.VSphereCluster.GetLabels()[clusterv1.ClusterNameLabel]
	if !ok {
		return nil, errors.Errorf("missing CAPI cluster label")
	}

	labels := map[string]string{clusterv1.ClusterNameLabel: name}
	kcpList := &controlplanev1.KubeadmControlPlaneList{}
	if err := r.Client.List(
		ctx, kcpList,
		client.InNamespace(ctx.VSphereCluster.GetNamespace()),
		client.MatchingLabels(labels)); err != nil {
		return nil, errors.Wrapf(err, "failed to list control plane objects")
	}
	if len(kcpList.Items) > 1 {
		return nil, errors.Errorf("multiple control plane objects found, expected 1, found %d", len(kcpList.Items))
	}

	if len(kcpList.Items) != 0 {
		if kcp := &kcpList.Items[0]; kcp.GetDeletionTimestamp().IsZero() {
			objects[appendKCPKey(kcp.GetName())] = clustermodule.NewWrapper(kcp)
		}
	}

	mdList := &clusterv1.MachineDeploymentList{}
	if err := r.Client.List(
		ctx, mdList,
		client.InNamespace(ctx.VSphereCluster.GetNamespace()),
		client.MatchingLabels(labels)); err != nil {
		return nil, errors.Wrapf(err, "failed to list machine deployment objects")
	}
	for _, md := range mdList.Items {
		if md.DeletionTimestamp.IsZero() {
			objects[md.GetName()] = clustermodule.NewWrapper(md.DeepCopy())
		}
	}
	return objects, nil
}

// appendKCPKey adds the prefix "kcp" to the name of the object
// This is used to separate a single KCP object from the Machine Deployment objects
// having the same name.
func appendKCPKey(name string) string {
	return "kcp" + name
}

func incompatibleOwnerErrors(errList []clusterModError) []clusterModError {
	toReport := []clusterModError{}
	for _, e := range errList {
		if clustermodule.IsIncompatibleOwnerError(e.err) {
			toReport = append(toReport, e)
		}
	}
	return toReport
}

type clusterModError struct {
	name string
	err  error
}

func generateClusterModuleErrorMessage(errList []clusterModError) string {
	sb := strings.Builder{}
	sb.WriteString("failed to create cluster modules for: ")

	for _, e := range errList {
		sb.WriteString(fmt.Sprintf("%s %s, ", e.name, e.err.Error()))
	}
	msg := sb.String()
	return msg[:len(msg)-2]
}
