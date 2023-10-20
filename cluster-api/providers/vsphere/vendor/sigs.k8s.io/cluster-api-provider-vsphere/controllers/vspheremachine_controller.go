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
	"time"

	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlbldr "sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/constants"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	inframanager "sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/vmoperator"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

const hostInfoErrStr = "host info cannot be used as a label value"

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheremachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspheremachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmware.infrastructure.cluster.x-k8s.io,resources=vspheremachinetemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=virtualmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vmoperator.vmware.com,resources=virtualmachineimages;virtualmachineimages/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=nodes;events;configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get;update;patch

// AddMachineControllerToManager adds the machine controller to the provided
// manager.

func AddMachineControllerToManager(ctx *context.ControllerManagerContext, mgr manager.Manager, controlledType client.Object, options controller.Options) error {
	supervisorBased, err := util.IsSupervisorType(controlledType)
	if err != nil {
		return err
	}

	var (
		controlledTypeName  = reflect.TypeOf(controlledType).Elem().Name()
		controlledTypeGVK   = infrav1.GroupVersion.WithKind(controlledTypeName)
		controllerNameShort = fmt.Sprintf("%s-controller", strings.ToLower(controlledTypeName))
		controllerNameLong  = fmt.Sprintf("%s/%s/%s", ctx.Namespace, ctx.Name, controllerNameShort)
	)

	if supervisorBased {
		controlledTypeGVK = vmwarev1.GroupVersion.WithKind(controlledTypeName)
		controllerNameShort = fmt.Sprintf("%s-supervisor-controller", strings.ToLower(controlledTypeName))
		controllerNameLong = fmt.Sprintf("%s/%s/%s", ctx.Namespace, ctx.Name, controllerNameShort)
	}

	// Build the controller context.
	controllerContext := &context.ControllerContext{
		ControllerManagerContext: ctx,
		Name:                     controllerNameShort,
		Recorder:                 record.New(mgr.GetEventRecorderFor(controllerNameLong)),
		Logger:                   ctx.Logger.WithName(controllerNameShort),
	}

	builder := ctrl.NewControllerManagedBy(mgr).
		// Watch the controlled, infrastructure resource.
		For(controlledType).
		WithOptions(options).
		// Watch the CAPI resource that owns this infrastructure resource.
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(clusterutilv1.MachineToInfrastructureMapFunc(controlledTypeGVK)),
		).
		// Watch a GenericEvent channel for the controlled resource.
		//
		// This is useful when there are events outside of Kubernetes that
		// should cause a resource to be synchronized, such as a goroutine
		// waiting on some asynchronous, external task to complete.
		WatchesRawSource(
			&source.Channel{Source: ctx.GetGenericEventChannelFor(controlledTypeGVK)},
			&handler.EnqueueRequestForObject{},
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), ctx.WatchFilterValue))

	r := &machineReconciler{
		ControllerContext: controllerContext,
		VMService:         &services.VimMachineService{},
		supervisorBased:   supervisorBased,
	}

	if supervisorBased {
		// Watch any VirtualMachine resources owned by this VSphereMachine
		builder.Owns(&vmoprv1.VirtualMachine{})
		r.VMService = &vmoperator.VmopMachineService{}
		networkProvider, err := inframanager.GetNetworkProvider(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to create a network provider")
		}
		r.networkProvider = networkProvider
	} else {
		// Watch any VSphereVM resources owned by the controlled type.
		builder.Watches(
			&infrav1.VSphereVM{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), controlledType),
			ctrlbldr.WithPredicates(predicate.Funcs{
				// ignore creation events since this controller is responsible for
				// the creation of the type.
				CreateFunc: func(e event.CreateEvent) bool {
					return false
				},
			}),
		)
	}

	c, err := builder.Build(r)
	if err != nil {
		return err
	}

	if !supervisorBased {
		err = c.Watch(
			source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
			handler.EnqueueRequestsFromMapFunc(r.clusterToVSphereMachines),
			predicates.ClusterUnpausedAndInfrastructureReady(r.Logger))
		if err != nil {
			return err
		}
	}
	return nil
}

type machineReconciler struct {
	*context.ControllerContext
	VMService       services.VSphereMachineService
	networkProvider services.NetworkProvider
	supervisorBased bool
}

// Reconcile ensures the back-end state reflects the Kubernetes resource state intent.
func (r *machineReconciler) Reconcile(_ goctx.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	var machineContext context.MachineContext
	logger := r.Logger.WithName(req.Namespace).WithName(req.Name)
	logger.V(3).Info("Starting Reconcile VSphereMachine")

	// Fetch VSphereMachine object and populate the machine context
	machineContext, err := r.VMService.FetchVSphereMachine(r.Client, req.NamespacedName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the CAPI Machine and CAPI Cluster.
	machine, err := clusterutilv1.GetOwnerMachine(r, r.Client, machineContext.GetObjectMeta())
	if err != nil {
		return reconcile.Result{}, err
	}
	if machine == nil {
		logger.V(2).Info("waiting on Machine controller to set OwnerRef on infra machine")
		return reconcile.Result{}, nil
	}

	cluster := r.fetchCAPICluster(machine, machineContext.GetVSphereMachine())

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(machineContext.GetVSphereMachine(), r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(
			err,
			"failed to init patch helper for %s/%s",
			machineContext.GetObjectMeta().Namespace,
			machineContext.GetObjectMeta().Name)
	}
	machineContext.SetBaseMachineContext(&context.BaseMachineContext{
		ControllerContext: r.ControllerContext,
		Cluster:           cluster,
		Machine:           machine,
		Logger:            logger,
		PatchHelper:       patchHelper,
	})
	// always patch the VSphereMachine object
	defer func() {
		// always update the readyCondition.
		conditions.SetSummary(machineContext.GetVSphereMachine(),
			conditions.WithConditions(
				infrav1.VMProvisionedCondition,
			),
		)

		// Patch the VSphereMachine resource.
		if err := machineContext.Patch(); err != nil {
			if reterr == nil {
				reterr = err
			}
			machineContext.GetLogger().Error(err, "patch failed", "machine", machineContext.String())
		}
	}()

	if !machineContext.GetObjectMeta().DeletionTimestamp.IsZero() {
		return r.reconcileDelete(machineContext)
	}

	// Checking whether cluster is nil here as we still want to allow delete even if cluster is not found.
	if cluster == nil {
		return reconcile.Result{}, nil
	}

	// Fetch the VSphereCluster and update the machine context
	machineContext, err = r.VMService.FetchVSphereCluster(r.Client, cluster, machineContext)
	if err != nil {
		logger.Info("unable to retrieve VSphereCluster", "error", err)
		return reconcile.Result{}, nil
	}

	// If the VSphereMachine doesn't have our finalizer, add it.
	// Requeue immediately after adding finalizer to avoid the race condition between init and delete
	if !ctrlutil.ContainsFinalizer(machineContext.GetVSphereMachine(), infrav1.MachineFinalizer) {
		ctrlutil.AddFinalizer(machineContext.GetVSphereMachine(), infrav1.MachineFinalizer)
		return reconcile.Result{}, nil
	}
	// Handle non-deleted machines
	return r.reconcileNormal(machineContext)
}

func (r *machineReconciler) reconcileDelete(ctx context.MachineContext) (reconcile.Result, error) {
	ctx.GetLogger().Info("Handling deleted VSphereMachine")
	conditions.MarkFalse(ctx.GetVSphereMachine(), infrav1.VMProvisionedCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")

	if err := r.VMService.ReconcileDelete(ctx); err != nil {
		if apierrors.IsNotFound(err) {
			// The VM is deleted so remove the finalizer.
			ctrlutil.RemoveFinalizer(ctx.GetVSphereMachine(), infrav1.MachineFinalizer)
			return reconcile.Result{}, nil
		}
		conditions.MarkFalse(ctx.GetVSphereMachine(), infrav1.VMProvisionedCondition, clusterv1.DeletionFailedReason, clusterv1.ConditionSeverityWarning, "")
		return reconcile.Result{}, err
	}

	// VM is being deleted
	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *machineReconciler) reconcileNormal(ctx context.MachineContext) (reconcile.Result, error) {
	machineFailed, err := r.VMService.SyncFailureReason(ctx)
	if err != nil && !apierrors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	// If the VSphereMachine is in an error state, return early.
	if machineFailed {
		ctx.GetLogger().Info("Error state detected, skipping reconciliation")
		return reconcile.Result{}, nil
	}

	//nolint:gocritic
	if r.supervisorBased {
		err := r.setVMModifiers(ctx)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		// vmwarev1.VSphereCluster doesn't set Cluster.Status.Ready until the API endpoint is available.
		if !ctx.GetCluster().Status.InfrastructureReady {
			ctx.GetLogger().Info("Cluster infrastructure is not ready yet")
			conditions.MarkFalse(ctx.GetVSphereMachine(), infrav1.VMProvisionedCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
			return reconcile.Result{}, nil
		}
	}

	// Make sure bootstrap data is available and populated.
	if ctx.GetMachine().Spec.Bootstrap.DataSecretName == nil {
		if !util.IsControlPlaneMachine(ctx.GetVSphereMachine()) && !conditions.IsTrue(ctx.GetCluster(), clusterv1.ControlPlaneInitializedCondition) {
			ctx.GetLogger().Info("Waiting for the control plane to be initialized")
			conditions.MarkFalse(ctx.GetVSphereMachine(), infrav1.VMProvisionedCondition, clusterv1.WaitingForControlPlaneAvailableReason, clusterv1.ConditionSeverityInfo, "")
			return ctrl.Result{}, nil
		}
		ctx.GetLogger().Info("Waiting for bootstrap data to be available")
		conditions.MarkFalse(ctx.GetVSphereMachine(), infrav1.VMProvisionedCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return reconcile.Result{}, nil
	}

	requeue, err := r.VMService.ReconcileNormal(ctx)
	if err != nil {
		return reconcile.Result{}, err
	} else if requeue {
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// The machine is patched at the last stage before marking the VM as provisioned
	// This makes sure that the VSphereMachine exists and is in a Running state
	// before attempting to patch.
	err = r.patchMachineLabelsWithHostInfo(ctx)
	if err != nil {
		r.Logger.Error(err, "failed to patch machine with host info label", "machine ", ctx.GetMachine().Name)
		return reconcile.Result{}, err
	}

	conditions.MarkTrue(ctx.GetVSphereMachine(), infrav1.VMProvisionedCondition)
	return reconcile.Result{}, nil
}

// patchMachineLabelsWithHostInfo adds the ESXi host information as a label to the Machine object.
// The ESXi host information is added with the CAPI node label prefix
// which would be added onto the node by the CAPI controllers.
func (r *machineReconciler) patchMachineLabelsWithHostInfo(ctx context.MachineContext) error {
	hostInfo, err := r.VMService.GetHostInfo(ctx)
	if err != nil {
		return err
	}

	info := util.SanitizeHostInfoLabel(hostInfo)
	errs := validation.IsValidLabelValue(info)
	if len(errs) > 0 {
		err := errors.Errorf("%s: %s", hostInfoErrStr, strings.Join(errs, ","))
		r.Logger.Error(err, hostInfoErrStr, "info", hostInfo)
		return err
	}

	machine := ctx.GetMachine()
	patchHelper, err := patch.NewHelper(machine, r.Client)
	if err != nil {
		return err
	}

	labels := machine.GetLabels()
	labels[constants.ESXiHostInfoLabel] = info
	machine.Labels = labels

	return patchHelper.Patch(r, machine)
}

func (r *machineReconciler) clusterToVSphereMachines(ctx goctx.Context, a client.Object) []reconcile.Request {
	requests := []reconcile.Request{}
	machines, err := util.GetVSphereMachinesInCluster(ctx, r.Client, a.GetNamespace(), a.GetName())
	if err != nil {
		return requests
	}
	for _, m := range machines {
		r := reconcile.Request{
			NamespacedName: apitypes.NamespacedName{
				Name:      m.Name,
				Namespace: m.Namespace,
			},
		}
		requests = append(requests, r)
	}
	return requests
}

func (r *machineReconciler) fetchCAPICluster(machine *clusterv1.Machine, vsphereMachine metav1.Object) *clusterv1.Cluster {
	cluster, err := clusterutilv1.GetClusterFromMetadata(r, r.Client, machine.ObjectMeta)
	if err != nil {
		r.Logger.Info("Machine is missing cluster label or cluster does not exist")
		return nil
	}
	if annotations.IsPaused(cluster, vsphereMachine) {
		r.Logger.V(4).Info("VSphereMachine %s/%s linked to a cluster that is paused", vsphereMachine.GetNamespace(), vsphereMachine.GetName())
		return nil
	}

	return cluster
}

// Return hooks that will be invoked when a VirtualMachine is created.
func (r *machineReconciler) setVMModifiers(c context.MachineContext) error {
	ctx, ok := c.(*vmware.SupervisorMachineContext)
	if !ok {
		return errors.New("received unexpected MachineContext. expecting SupervisorMachineContext type")
	}

	networkModifier := func(obj runtime.Object) (runtime.Object, error) {
		// No need to check the type. We know this will be a VirtualMachine
		vm, _ := obj.(*vmoprv1.VirtualMachine)
		ctx.Logger.V(3).Info("Applying network config to VM", "vm-name", vm.Name)
		err := r.networkProvider.ConfigureVirtualMachine(ctx.GetClusterContext(), vm)
		if err != nil {
			return nil, errors.Errorf("failed to configure machine network: %+v", err)
		}
		return vm, nil
	}
	ctx.VMModifiers = []vmware.VMModifier{networkModifier}
	return nil
}
