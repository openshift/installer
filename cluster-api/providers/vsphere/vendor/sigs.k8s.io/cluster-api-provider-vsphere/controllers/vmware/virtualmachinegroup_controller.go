/*
Copyright 2025 The Kubernetes Authors.

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

package vmware

import (
	"context"

	pkgerrors "github.com/pkg/errors"
	apitypes "k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	capicontrollerutil "sigs.k8s.io/cluster-api/util/controller"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
	conversionclient "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/client"
)

// AddVirtualMachineGroupControllerToManager adds the VirtualMachineGroup controller to the provided manager.
func AddVirtualMachineGroupControllerToManager(ctx context.Context, controllerManagerCtx *capvcontext.ControllerManagerContext, mgr manager.Manager, options controller.Options) error {
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "virtualmachinegroup")

	reconciler := &VirtualMachineGroupReconciler{
		Client:   controllerManagerCtx.Client,
		Recorder: mgr.GetEventRecorderFor("virtualmachinegroup-controller"),
	}

	// NOTE: use vm-operator native types for watches (the reconciler uses the internal hub version).
	vmGroup, err := conversionclient.WatchObject(reconciler.Client, &vmoprvhub.VirtualMachineGroup{})
	if err != nil {
		return pkgerrors.Wrapf(err, "failed to create watch object for VirtualMachineGroup")
	}

	builder := capicontrollerutil.NewControllerManagedBy(mgr, predicateLog).
		For(&clusterv1.Cluster{}).
		WithOptions(options).
		// Set the controller's name explicitly to virtualmachinegroup.
		Named("virtualmachinegroup").
		Watches(
			vmGroup,
			handler.EnqueueRequestForOwner(mgr.GetScheme(), reconciler.Client.RESTMapper(), &clusterv1.Cluster{}),
		).
		Watches(
			&vmwarev1.VSphereMachine{},
			handler.EnqueueRequestsFromMapFunc(reconciler.VSphereMachineToCluster),
			predicate.Funcs{
				UpdateFunc: func(event.UpdateEvent) bool { return false },
				CreateFunc: func(e event.CreateEvent) bool {
					// Only handle VSphereMachine which belongs to a MachineDeployment
					_, found := e.Object.GetLabels()[clusterv1.MachineDeploymentNameLabel]
					return found
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					// Only handle VSphereMachine which belongs to a MachineDeployment
					_, found := e.Object.GetLabels()[clusterv1.MachineDeploymentNameLabel]
					return found
				},
				GenericFunc: func(event.GenericEvent) bool { return false },
			},
		).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), predicateLog, controllerManagerCtx.WatchFilterValue))

	return builder.Complete(ctx, reconciler)
}

// VSphereMachineToCluster maps VSphereMachine events to Cluster reconcile requests.
func (r *VirtualMachineGroupReconciler) VSphereMachineToCluster(_ context.Context, a ctrlclient.Object) []reconcile.Request {
	vSphereMachine, ok := a.(*vmwarev1.VSphereMachine)
	if !ok {
		return nil
	}

	clusterName, ok := vSphereMachine.Labels[clusterv1.ClusterNameLabel]
	if !ok || clusterName == "" {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: apitypes.NamespacedName{
			Namespace: vSphereMachine.Namespace,
			Name:      clusterName,
		},
	}}
}
