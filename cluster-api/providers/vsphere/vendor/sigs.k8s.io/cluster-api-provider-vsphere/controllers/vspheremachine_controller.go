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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	clog "sigs.k8s.io/cluster-api/util/log"
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
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	inframanager "sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/vmoperator"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

const (
	hostInfoErrStr = "host info cannot be used as a label value"
)

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
func AddMachineControllerToManager(ctx context.Context, controllerManagerContext *capvcontext.ControllerManagerContext, mgr manager.Manager, supervisorBased bool, options controller.Options) error {
	r := &machineReconciler{
		Client:          controllerManagerContext.Client,
		Recorder:        mgr.GetEventRecorderFor("vspheremachine-controller"),
		VMService:       &services.VimMachineService{Client: controllerManagerContext.Client},
		supervisorBased: supervisorBased,
	}
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "vspheremachine")

	if supervisorBased {
		networkProvider, err := inframanager.GetNetworkProvider(ctx, controllerManagerContext.Client, controllerManagerContext.NetworkProvider)
		if err != nil {
			return errors.Wrap(err, "failed to create a network provider")
		}
		r.networkProvider = networkProvider
		r.VMService = &vmoperator.VmopMachineService{Client: controllerManagerContext.Client, ConfigureControlPlaneVMReadinessProbe: r.networkProvider.SupportsVMReadinessProbe()}

		return ctrl.NewControllerManagedBy(mgr).
			// Watch the controlled, infrastructure resource.
			For(&vmwarev1.VSphereMachine{}).
			WithOptions(options).
			// Watch the CAPI resource that owns this infrastructure resource.
			Watches(
				&clusterv1.Machine{},
				handler.EnqueueRequestsFromMapFunc(clusterutilv1.MachineToInfrastructureMapFunc(vmwarev1.GroupVersion.WithKind("VSphereMachine"))),
			).
			Watches(
				&clusterv1.Cluster{},
				handler.EnqueueRequestsFromMapFunc(r.enqueueClusterToMachineRequests),
				ctrlbldr.WithPredicates(
					predicates.ClusterUnpausedAndInfrastructureReady(mgr.GetScheme(), predicateLog),
				),
			).
			// Watch a GenericEvent channel for the controlled resource.
			//
			// This is useful when there are events outside of Kubernetes that
			// should cause a resource to be synchronized, such as a goroutine
			// waiting on some asynchronous, external task to complete.
			WatchesRawSource(
				source.Channel(
					controllerManagerContext.GetGenericEventChannelFor(vmwarev1.GroupVersion.WithKind("VSphereMachine")),
					&handler.EnqueueRequestForObject{},
				),
			).
			WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), predicateLog, controllerManagerContext.WatchFilterValue)).
			// Watch any VirtualMachine resources owned by this VSphereMachine
			Owns(&vmoprv1.VirtualMachine{}).
			Complete(r)
	}

	return ctrl.NewControllerManagedBy(mgr).
		// Watch the controlled, infrastructure resource.
		For(&infrav1.VSphereMachine{}).
		WithOptions(options).
		// Watch the CAPI resource that owns this infrastructure resource.
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(clusterutilv1.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("VSphereMachine"))),
		).
		// Watch a GenericEvent channel for the controlled resource.
		//
		// This is useful when there are events outside of Kubernetes that
		// should cause a resource to be synchronized, such as a goroutine
		// waiting on some asynchronous, external task to complete.
		WatchesRawSource(
			source.Channel(
				controllerManagerContext.GetGenericEventChannelFor(infrav1.GroupVersion.WithKind("VSphereMachine")),
				&handler.EnqueueRequestForObject{},
			),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), predicateLog, controllerManagerContext.WatchFilterValue)).
		// Watch any VSphereVM resources owned by the controlled type.
		Watches(
			&infrav1.VSphereVM{},
			handler.EnqueueRequestForOwner(mgr.GetScheme(), mgr.GetRESTMapper(), &infrav1.VSphereMachine{}),
			ctrlbldr.WithPredicates(predicate.Funcs{
				// ignore creation events since this controller is responsible for
				// the creation of the type.
				CreateFunc: func(event.CreateEvent) bool {
					return false
				},
			}),
		).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(r.enqueueClusterToMachineRequests),
			ctrlbldr.WithPredicates(
				predicates.ClusterUnpausedAndInfrastructureReady(mgr.GetScheme(), predicateLog),
			),
		).Complete(r)
}

type machineReconciler struct {
	Client          client.Client
	Recorder        record.EventRecorder
	VMService       services.VSphereMachineService
	networkProvider services.NetworkProvider
	supervisorBased bool
}

// Reconcile ensures the back-end state reflects the Kubernetes resource state intent.
func (r *machineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch VSphereMachine object and populate the machine context
	machineContext, err := r.VMService.FetchVSphereMachine(ctx, req.NamespacedName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Fetch the CAPI Machine.
	machine, err := clusterutilv1.GetOwnerMachine(ctx, r.Client, machineContext.GetObjectMeta())
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to get Machine for VSphereMachine")
	}
	if machine == nil {
		log.Info("Waiting for Machine controller to set OwnerRef on VSphereMachine")
		return reconcile.Result{}, nil
	}
	log = log.WithValues("Machine", klog.KObj(machine))
	ctx = ctrl.LoggerInto(ctx, log)

	// AddOwners adds the owners of Machine as k/v pairs to the logger.
	// Specifically, it will add KubeadmControlPlane, MachineSet and MachineDeployment.
	ctx, log, err = clog.AddOwners(ctx, r.Client, machine)
	if err != nil {
		return ctrl.Result{}, err
	}

	cluster, err := clusterutilv1.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to get Cluster from VSphereCluster: Machine is missing cluster label or cluster does not exist")
	}
	if cluster != nil {
		log = log.WithValues("Cluster", klog.KObj(cluster))
		if cluster.Spec.InfrastructureRef != nil {
			log = log.WithValues("VSphereCluster", klog.KRef(cluster.Namespace, cluster.Spec.InfrastructureRef.Name))
		}
		ctx = ctrl.LoggerInto(ctx, log)

		if annotations.IsPaused(cluster, machineContext.GetVSphereMachine()) {
			log.Info("Reconciliation is paused for this object")
			return reconcile.Result{}, nil
		}
	} else if annotations.HasPaused(machineContext.GetVSphereMachine()) {
		log.Info("Reconciliation is paused for this object")
		return reconcile.Result{}, nil
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(machineContext.GetVSphereMachine(), r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}
	machineContext.SetBaseMachineContext(&capvcontext.BaseMachineContext{
		Cluster:     cluster,
		Machine:     machine,
		PatchHelper: patchHelper,
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
		if err := machineContext.Patch(ctx); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	if !machineContext.GetObjectMeta().DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machineContext)
	}

	// Checking whether cluster is nil here as we still want to run reconcileDelete above even if cluster is not found.
	if cluster == nil {
		log.Info("Failed to get Cluster")
		return reconcile.Result{}, nil
	}

	if cluster.Spec.InfrastructureRef == nil {
		log.Info("Cluster.spec.infrastructureRef is not yet set")
		return reconcile.Result{}, nil
	}

	// If the VSphereCluster has been previously set as an ownerReference remove it. This may have been set in an older
	// version of CAPV to prevent VSphereMachines from being orphaned, but is no longer needed.
	// For more info see: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/issues/2054
	// TODO: This should be removed in a future release of CAPV.
	machineContext.GetVSphereMachine().SetOwnerReferences(
		clusterutilv1.RemoveOwnerRef(
			machineContext.GetVSphereMachine().GetOwnerReferences(),
			metav1.OwnerReference{
				Name:       cluster.Spec.InfrastructureRef.Name,
				APIVersion: cluster.Spec.InfrastructureRef.APIVersion,
				Kind:       cluster.Spec.InfrastructureRef.Kind,
			},
		),
	)

	// Fetch the VSphereCluster and update the machine context
	machineContext, err = r.VMService.FetchVSphereCluster(ctx, cluster, machineContext)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to get VSphereCluster")
	}

	// If the VSphereMachine doesn't have our finalizer, add it.
	// Requeue immediately after adding finalizer to avoid the race condition between init and delete
	if !ctrlutil.ContainsFinalizer(machineContext.GetVSphereMachine(), infrav1.MachineFinalizer) {
		ctrlutil.AddFinalizer(machineContext.GetVSphereMachine(), infrav1.MachineFinalizer)
		return reconcile.Result{}, nil
	}

	// Handle non-deleted machines
	return r.reconcileNormal(ctx, machineContext)
}

func (r *machineReconciler) reconcileDelete(ctx context.Context, machineCtx capvcontext.MachineContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	conditions.MarkFalse(machineCtx.GetVSphereMachine(), infrav1.VMProvisionedCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")

	if err := r.VMService.ReconcileDelete(ctx, machineCtx); err != nil {
		if apierrors.IsNotFound(err) {
			// The VM is deleted so remove the finalizer.
			if ctrlutil.RemoveFinalizer(machineCtx.GetVSphereMachine(), infrav1.MachineFinalizer) {
				log.Info(fmt.Sprintf("Removing finalizer %s", infrav1.MachineFinalizer))
			}
			return reconcile.Result{}, nil
		}
		conditions.MarkFalse(machineCtx.GetVSphereMachine(), infrav1.VMProvisionedCondition, clusterv1.DeletionFailedReason, clusterv1.ConditionSeverityWarning, "")
		return reconcile.Result{}, err
	}

	// VM is being deleted
	return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *machineReconciler) reconcileNormal(ctx context.Context, machineCtx capvcontext.MachineContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	machineFailed, err := r.VMService.SyncFailureReason(ctx, machineCtx)
	if err != nil && !apierrors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	// If the VSphereMachine is in an error state, return early.
	if machineFailed {
		log.Error(err, "Error state detected, skipping reconciliation")
		return reconcile.Result{}, nil
	}

	// Cluster `.status.infrastructureReady == false is handled differently depending on if the machine is supervisor based.
	// 1) If the Cluster is not supervisor-based mark the VMProvisionedCondition false and return nil.
	// 2) If the Cluster is supervisor-based continue to reconcile as InfrastructureReady is not set to true until after the kube apiserver is available.
	if !r.supervisorBased {
		// vmwarev1.VSphereCluster doesn't set Cluster.Status.Ready until the API endpoint is available.
		if !machineCtx.GetCluster().Status.InfrastructureReady {
			log.Info("Cluster infrastructure is not ready yet, skipping reconciliation")
			conditions.MarkFalse(machineCtx.GetVSphereMachine(), infrav1.VMProvisionedCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
			return reconcile.Result{}, nil
		}
	} else {
		if err := r.setVMModifiers(ctx, machineCtx); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Make sure bootstrap data is available and populated.
	if machineCtx.GetMachine().Spec.Bootstrap.DataSecretName == nil {
		if !util.IsControlPlaneMachine(machineCtx.GetVSphereMachine()) && !conditions.IsTrue(machineCtx.GetCluster(), clusterv1.ControlPlaneInitializedCondition) {
			log.Info("Waiting for the control plane to be initialized, skipping reconciliation")
			conditions.MarkFalse(machineCtx.GetVSphereMachine(), infrav1.VMProvisionedCondition, clusterv1.WaitingForControlPlaneAvailableReason, clusterv1.ConditionSeverityInfo, "")
			return ctrl.Result{}, nil
		}
		log.Info("Waiting for bootstrap data to be ready, skipping reconciliation")
		conditions.MarkFalse(machineCtx.GetVSphereMachine(), infrav1.VMProvisionedCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return reconcile.Result{}, nil
	}

	requeue, err := r.VMService.ReconcileNormal(ctx, machineCtx)
	if err != nil {
		return reconcile.Result{}, err
	} else if requeue {
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// The machine is patched at the last stage before marking the VM as provisioned
	// This makes sure that the VSphereMachine exists and is in a Running state
	// before attempting to patch.
	err = r.patchMachineLabelsWithHostInfo(ctx, machineCtx)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to patch Machine with host info label")
	}

	conditions.MarkTrue(machineCtx.GetVSphereMachine(), infrav1.VMProvisionedCondition)
	return reconcile.Result{}, nil
}

// patchMachineLabelsWithHostInfo adds the ESXi host information as a label to the Machine object.
// The ESXi host information is added with the CAPI node label prefix
// which would be added onto the node by the CAPI controllers.
func (r *machineReconciler) patchMachineLabelsWithHostInfo(ctx context.Context, machineCtx capvcontext.MachineContext) error {
	hostInfo, err := r.VMService.GetHostInfo(ctx, machineCtx)
	if err != nil {
		return err
	}

	info := util.SanitizeHostInfoLabel(hostInfo)
	errs := validation.IsValidLabelValue(info)
	if len(errs) > 0 {
		return errors.Errorf("%s (hostInfo: %s): %s", hostInfoErrStr, hostInfo, strings.Join(errs, ","))
	}

	machine := machineCtx.GetMachine()
	patchHelper, err := patch.NewHelper(machine, r.Client)
	if err != nil {
		return err
	}

	labels := machine.GetLabels()
	labels[constants.ESXiHostInfoLabel] = info
	machine.Labels = labels

	return patchHelper.Patch(ctx, machine)
}

// Return hooks that will be invoked when a VirtualMachine is created.
func (r *machineReconciler) setVMModifiers(ctx context.Context, machineCtx capvcontext.MachineContext) error {
	log := ctrl.LoggerFrom(ctx)
	supervisorMachineCtx, ok := machineCtx.(*vmware.SupervisorMachineContext)
	if !ok {
		return errors.New("received unexpected MachineContext. expecting SupervisorMachineContext type")
	}

	networkModifier := func(obj runtime.Object) (runtime.Object, error) {
		// No need to check the type. We know this will be a VirtualMachine
		vm, _ := obj.(*vmoprv1.VirtualMachine)
		log.V(3).Info("Applying network config to VM")
		err := r.networkProvider.ConfigureVirtualMachine(ctx, supervisorMachineCtx.GetClusterContext(), vm)
		if err != nil {
			return nil, errors.Errorf("failed to configure machine network: %+v", err)
		}
		return vm, nil
	}
	supervisorMachineCtx.VMModifiers = []vmware.VMModifier{networkModifier}
	return nil
}

// enqueueClusterToMachineRequests returns a list of VSphereMachine reconcile requests
// belonging to the cluster.
func (r *machineReconciler) enqueueClusterToMachineRequests(ctx context.Context, a client.Object) []reconcile.Request {
	requests := []reconcile.Request{}
	machines, err := r.VMService.GetMachinesInCluster(ctx, a.GetNamespace(), a.GetName())
	if err != nil {
		return requests
	}
	for _, m := range machines {
		r := reconcile.Request{
			NamespacedName: apitypes.NamespacedName{
				Name:      m.GetName(),
				Namespace: m.GetNamespace(),
			},
		}
		requests = append(requests, r)
	}
	return requests
}
