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
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/clustercache"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"
	clusterutilv1 "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	clog "sigs.k8s.io/cluster-api/util/log"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlbldr "sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/clustermodule"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/identity"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspherevms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=vspherevms/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments;machinesets,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=kubeadmcontrolplanes,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;delete

// AddVMControllerToManager adds the VM controller to the provided manager.
func AddVMControllerToManager(ctx context.Context, controllerManagerCtx *capvcontext.ControllerManagerContext, mgr manager.Manager, clusterCache clustercache.ClusterCache, options controller.Options) error {
	r := vmReconciler{
		ControllerManagerContext: controllerManagerCtx,
		Recorder:                 mgr.GetEventRecorderFor("vspherevm-controller"),
		VMService:                &govmomi.VMService{},
		clusterCache:             clusterCache,
	}
	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "vspherevm")

	return ctrl.NewControllerManagedBy(mgr).
		// Watch the controlled, infrastructure resource.
		For(&infrav1.VSphereVM{}).
		WithOptions(options).
		// Watch a GenericEvent channel for the controlled resource.
		//
		// This is useful when there are events outside of Kubernetes that
		// should cause a resource to be synchronized, such as a goroutine
		// waiting on some asynchronous, external task to complete.
		WatchesRawSource(
			source.Channel(
				controllerManagerCtx.GetGenericEventChannelFor(infrav1.GroupVersion.WithKind("VSphereVM")),
				&handler.EnqueueRequestForObject{},
			),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), predicateLog, controllerManagerCtx.WatchFilterValue)).
		Watches(
			&clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(r.clusterToVSphereVMs),
			ctrlbldr.WithPredicates(
				predicate.Funcs{
					UpdateFunc: func(e event.UpdateEvent) bool {
						newCluster := e.ObjectNew.(*clusterv1.Cluster)
						// check whether cluster has either spec.paused or pasued annotation
						return !annotations.IsPaused(newCluster, newCluster)
					},
					CreateFunc: func(e event.CreateEvent) bool {
						cluster := e.Object.(*clusterv1.Cluster)
						// check whether cluster has either spec.paused or pasued annotation
						return annotations.IsPaused(cluster, cluster)
					},
				}),
		).
		Watches(
			&infrav1.VSphereCluster{},
			handler.EnqueueRequestsFromMapFunc(r.vsphereClusterToVSphereVMs),
			ctrlbldr.WithPredicates(
				predicate.Funcs{
					UpdateFunc: func(e event.UpdateEvent) bool {
						oldCluster := e.ObjectOld.(*infrav1.VSphereCluster)
						newCluster := e.ObjectNew.(*infrav1.VSphereCluster)
						return !clustermodule.Compare(oldCluster.Spec.ClusterModules, newCluster.Spec.ClusterModules)
					},
					CreateFunc:  func(event.CreateEvent) bool { return false },
					DeleteFunc:  func(event.DeleteEvent) bool { return false },
					GenericFunc: func(event.GenericEvent) bool { return false },
				}),
		).
		Watches(
			&ipamv1.IPAddressClaim{},
			handler.EnqueueRequestsFromMapFunc(r.ipAddressClaimToVSphereVM),
		).
		WatchesRawSource(r.clusterCache.GetClusterSource("vspherevm", r.clusterToVSphereVMs)).
		Complete(r)
}

type vmReconciler struct {
	Recorder record.EventRecorder
	*capvcontext.ControllerManagerContext
	VMService    services.VirtualMachineService
	clusterCache clustercache.ClusterCache
}

// Reconcile ensures the back-end state reflects the Kubernetes resource state intent.
func (r vmReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Get the VSphereVM resource for this request.
	vsphereVM := &infrav1.VSphereVM{}
	if err := r.Client.Get(ctx, req.NamespacedName, vsphereVM); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	cluster, err := clusterutilv1.GetClusterFromMetadata(ctx, r.Client, vsphereVM.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to get Cluster from VSphereVM: Machine is missing cluster label or cluster does not exist")
	}
	if cluster != nil {
		log = log.WithValues("Cluster", klog.KObj(cluster))
		ctx = ctrl.LoggerInto(ctx, log)

		if annotations.IsPaused(cluster, vsphereVM) {
			log.Info("Reconciliation is paused for this object")
			return reconcile.Result{}, nil
		}
	} else if annotations.HasPaused(vsphereVM) {
		log.Info("Reconciliation is paused for this object")
		return reconcile.Result{}, nil
	}

	// Create the patch helper.
	patchHelper, err := patch.NewHelper(vsphereVM, r.Client)
	if err != nil {
		return reconcile.Result{}, err
	}

	authSession, err := r.retrieveVcenterSession(ctx, vsphereVM)
	if err != nil {
		conditions.MarkFalse(vsphereVM, infrav1.VCenterAvailableCondition, infrav1.VCenterUnreachableReason, clusterv1.ConditionSeverityError, err.Error())
		return reconcile.Result{}, err
	}
	conditions.MarkTrue(vsphereVM, infrav1.VCenterAvailableCondition)

	// Fetch the owner VSphereMachine.
	vsphereMachine, err := util.GetOwnerVSphereMachine(ctx, r.Client, vsphereVM.ObjectMeta)
	// vsphereMachine can be nil in cases where custom mover other than clusterctl
	// moves the resources without ownerreferences set
	// in that case nil vsphereMachine can cause panic and CrashLoopBackOff the pod
	// preventing vspheremachine_controller from setting the ownerref
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to get VSphereMachine for VSphereVM")
	}
	if vsphereMachine == nil {
		log.Info("Waiting for VSphereMachine controller to set OwnerRef on VSphereVM")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("VSphereMachine", klog.KObj(vsphereMachine))
	ctx = ctrl.LoggerInto(ctx, log)

	vsphereCluster, err := util.GetVSphereClusterFromVSphereMachine(ctx, r.Client, vsphereMachine)
	if err != nil || vsphereCluster == nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to get VSphereCluster from VSphereMachine")
	}

	log = log.WithValues("VSphereCluster", klog.KObj(vsphereCluster))
	ctx = ctrl.LoggerInto(ctx, log)

	// Fetch the CAPI Machine.
	machine, err := clusterutilv1.GetOwnerMachine(ctx, r.Client, vsphereMachine.ObjectMeta)
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

	var vsphereFailureDomain *infrav1.VSphereFailureDomain
	if failureDomain := machine.Spec.FailureDomain; failureDomain != nil {
		vsphereDeploymentZone := &infrav1.VSphereDeploymentZone{}
		if err := r.Client.Get(ctx, apitypes.NamespacedName{Name: *failureDomain}, vsphereDeploymentZone); err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to get VSphereDeploymentZone %s", *failureDomain)
		}

		vsphereFailureDomain = &infrav1.VSphereFailureDomain{}
		if err := r.Client.Get(ctx, apitypes.NamespacedName{Name: vsphereDeploymentZone.Spec.FailureDomain}, vsphereFailureDomain); err != nil {
			return reconcile.Result{}, errors.Wrapf(err, "failed to get VSphereFailureDomain %s", vsphereDeploymentZone.Spec.FailureDomain)
		}
	}

	// Create the VM context for this request.
	vmContext := &capvcontext.VMContext{
		ControllerManagerContext: r.ControllerManagerContext,
		VSphereVM:                vsphereVM,
		VSphereFailureDomain:     vsphereFailureDomain,
		Session:                  authSession,
		PatchHelper:              patchHelper,
	}

	// Print the task-ref upon entry and upon exit.
	log.V(4).Info("VSphereVM.Status.TaskRef OnEntry", "taskRef", vmContext.VSphereVM.Status.TaskRef)
	defer func() {
		log.V(4).Info("VSphereVM.Status.TaskRef OnExit", "taskRef", vmContext.VSphereVM.Status.TaskRef)
	}()
	originalTaskRef := vmContext.VSphereVM.Status.TaskRef

	// Always issue a patch when exiting this function so changes to the
	// resource are patched back to the API server.
	defer func() {
		// always update the readyCondition.
		conditions.SetSummary(vmContext.VSphereVM,
			conditions.WithConditions(
				infrav1.VCenterAvailableCondition,
				infrav1.IPAddressClaimedCondition,
				infrav1.VMProvisionedCondition,
			),
		)

		// Patch the VSphereVM resource.
		if err := vmContext.Patch(ctx); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}

		// Wait until VSphereVM is updated in the cache if the `.Status.TaskRef` field changes.
		// Note: We have to do this because otherwise using a cached client in current state could
		// return a stale state of a VSphereVM we just patched (because the cache might be stale).
		// This can lead to duplicate tasks being triggered (e.g. VM deletion) and make the controller
		// wait for longer then required.
		if vmContext.VSphereVM.Status.TaskRef != originalTaskRef {
			err = wait.PollUntilContextTimeout(ctx, 5*time.Millisecond, 5*time.Second, true, func(ctx context.Context) (bool, error) {
				key := ctrlclient.ObjectKey{Namespace: vmContext.VSphereVM.GetNamespace(), Name: vmContext.VSphereVM.GetName()}
				cachedVSphereVM := &infrav1.VSphereVM{}
				if err := r.Client.Get(ctx, key, cachedVSphereVM); err != nil {
					return false, err
				}
				return originalTaskRef != cachedVSphereVM.Status.TaskRef, nil
			})
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	if vsphereVM.ObjectMeta.DeletionTimestamp.IsZero() {
		// If the VSphereVM doesn't have our finalizer, add it.
		// Requeue immediately to avoid the race condition between init and delete
		if !ctrlutil.ContainsFinalizer(vsphereVM, infrav1.VMFinalizer) {
			ctrlutil.AddFinalizer(vsphereVM, infrav1.VMFinalizer)
			return reconcile.Result{}, nil
		}
	}

	return r.reconcile(ctx, vmContext, fetchClusterModuleInput{
		VSphereCluster: vsphereCluster,
		Machine:        machine,
	})
}

// reconcile encases the behavior of the controller around cluster module information
// retrieval depending upon inputs passed.
//
// This logic was moved to a smaller function outside the main Reconcile() loop
// for the ease of testing.
func (r vmReconciler) reconcile(ctx context.Context, vmCtx *capvcontext.VMContext, input fetchClusterModuleInput) (reconcile.Result, error) {
	if feature.Gates.Enabled(feature.NodeAntiAffinity) && !input.VSphereCluster.Spec.DisableClusterModule {
		clusterModuleInfo, err := r.fetchClusterModuleInfo(ctx, input)
		// If cluster module information cannot be fetched for a VM being deleted,
		// we should not block VM deletion since the cluster module is updated
		// once the VM gets removed.
		if err != nil && vmCtx.VSphereVM.ObjectMeta.DeletionTimestamp.IsZero() {
			return reconcile.Result{}, err
		}
		vmCtx.ClusterModuleInfo = clusterModuleInfo
	}

	// Handle deleted machines
	if !vmCtx.VSphereVM.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, vmCtx)
	}

	// Handle non-deleted machines
	return r.reconcileNormal(ctx, vmCtx)
}

func (r vmReconciler) reconcileDelete(ctx context.Context, vmCtx *capvcontext.VMContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	result, vm, err := r.VMService.DestroyVM(ctx, vmCtx)
	if err != nil {
		conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, "DeletionFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return reconcile.Result{}, errors.Wrapf(err, "failed to destroy VM")
	}

	if !result.IsZero() {
		// a non-zero value means we need to requeue the request before proceed.
		return result, nil
	}

	// Requeue the operation until the VM is "notfound".
	if vm.State != infrav1.VirtualMachineStateNotFound {
		log.Info(fmt.Sprintf("VM state is %q, waiting for %q", vm.State, infrav1.VirtualMachineStateNotFound))
		return reconcile.Result{}, nil
	}

	// Attempt to delete the node corresponding to the vsphere VM
	err = r.deleteNode(ctx, vmCtx, vm.Name)
	if err != nil {
		log.Error(err, "Failed to delete Node (best-effort)")
	}

	if err := r.deleteIPAddressClaims(ctx, vmCtx); err != nil {
		return reconcile.Result{}, err
	}

	// The VM is deleted so remove the finalizer.
	if ctrlutil.RemoveFinalizer(vmCtx.VSphereVM, infrav1.VMFinalizer) {
		log.Info(fmt.Sprintf("Removing finalizer %s", infrav1.VMFinalizer))
	}

	return reconcile.Result{}, nil
}

// deleteNode attempts to find and best effort delete the node corresponding to the VM
// This is necessary since CAPI does not surface the nodeRef field on the owner Machine object
// until the node moves to Ready state. Hence, on Machine deletion it is unable to delete
// the kubernetes node corresponding to the VM.
// Note: If this fails, CPI normally cleans up orphaned nodes.
func (r vmReconciler) deleteNode(ctx context.Context, vmCtx *capvcontext.VMContext, name string) error {
	log := ctrl.LoggerFrom(ctx)
	// Fetching the cluster object from the VSphereVM object to create a remote client to the cluster
	cluster, err := clusterutilv1.GetClusterFromMetadata(ctx, r.Client, vmCtx.VSphereVM.ObjectMeta)
	if err != nil {
		return err
	}

	// Skip deleting the Node if the cluster is being deleted.
	if !cluster.DeletionTimestamp.IsZero() {
		return nil
	}

	clusterClient, err := r.clusterCache.GetClient(ctx, ctrlclient.ObjectKeyFromObject(cluster))
	if err != nil {
		if errors.Is(err, clustercache.ErrClusterNotConnected) {
			log.V(2).Info("Skipping node deletion because connection to the workload cluster is down")
			return nil
		}
		return err
	}

	// Attempt to delete the corresponding node
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return clusterClient.Delete(ctx, node)
}

func (r vmReconciler) reconcileNormal(ctx context.Context, vmCtx *capvcontext.VMContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	if vmCtx.VSphereVM.Status.FailureReason != nil || vmCtx.VSphereVM.Status.FailureMessage != nil {
		log.Info("VM is failed, won't reconcile")
		return reconcile.Result{}, nil
	}

	if r.isWaitingForStaticIPAllocation(vmCtx) {
		conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.WaitingForStaticIPAllocationReason, clusterv1.ConditionSeverityInfo, "")
		log.Info("VM is waiting for static ip to be available")
		return reconcile.Result{}, nil
	}

	if err := r.reconcileIPAddressClaims(ctx, vmCtx); err != nil {
		return reconcile.Result{}, err
	}

	// Get or create the VM.
	vm, err := r.VMService.ReconcileVM(ctx, vmCtx)
	if err != nil {
		return reconcile.Result{}, errors.Wrapf(err, "failed to reconcile VM")
	}

	// Do not proceed until the backend VM is marked ready.
	if vm.State != infrav1.VirtualMachineStateReady {
		log.Info(fmt.Sprintf("VM state is %q, waiting for %q", vm.State, infrav1.VirtualMachineStateReady))
		return reconcile.Result{}, nil
	}

	// Update the VSphereVM's BIOS UUID.
	// Defensive check to ensure we are not removing the biosUUID
	if vm.BiosUUID != "" {
		if vmCtx.VSphereVM.Spec.BiosUUID != vm.BiosUUID {
			log.Info("Update VM biosUUID", "biosUUID", vm.BiosUUID)
			vmCtx.VSphereVM.Spec.BiosUUID = vm.BiosUUID
		}
	} else {
		return reconcile.Result{}, errors.Errorf("biosUUID is empty while VM is ready")
	}

	// VMRef should be set just once. It is not supposed to change!
	if vm.VMRef != "" && vmCtx.VSphereVM.Status.VMRef == "" {
		log.Info("Update VM vmRef", "vmRef", vm.VMRef)
		vmCtx.VSphereVM.Status.VMRef = vm.VMRef
	}

	// Update the VSphereVM's network status.
	r.reconcileNetwork(vmCtx, vm)

	// we didn't get any addresses, requeue
	if len(vmCtx.VSphereVM.Status.Addresses) == 0 {
		conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.WaitingForIPAllocationReason, clusterv1.ConditionSeverityInfo, "")
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Once the network is online the VM is considered ready.
	vmCtx.VSphereVM.Status.Ready = true
	conditions.MarkTrue(vmCtx.VSphereVM, infrav1.VMProvisionedCondition)
	log.Info("VSphereVM is ready")
	return reconcile.Result{}, nil
}

// isWaitingForStaticIPAllocation checks whether the VM should wait for a static IP
// to be allocated.
// It checks the state of both DHCP4 and DHCP6 for all the network devices and if
// any static IP addresses or IPAM Pools are specified.
func (r vmReconciler) isWaitingForStaticIPAllocation(vmCtx *capvcontext.VMContext) bool {
	devices := vmCtx.VSphereVM.Spec.Network.Devices
	for _, dev := range devices {
		// Ignore device if SkipIPAllocation is set.
		if dev.SkipIPAllocation {
			continue
		}

		// Ignore device if it is configured to use DHCP.
		if dev.DHCP4 || dev.DHCP6 {
			continue
		}

		if len(dev.IPAddrs) == 0 && len(dev.AddressesFromPools) == 0 {
			// One or more IPs are expected for the device but are not set yet.
			return true
		}
	}

	return false
}

func (r vmReconciler) reconcileNetwork(vmCtx *capvcontext.VMContext, vm infrav1.VirtualMachine) {
	vmCtx.VSphereVM.Status.Network = vm.Network
	ipAddrs := make([]string, 0, len(vm.Network))
	for _, netStatus := range vmCtx.VSphereVM.Status.Network {
		ipAddrs = append(ipAddrs, netStatus.IPAddrs...)
	}
	vmCtx.VSphereVM.Status.Addresses = ipAddrs
}

func (r vmReconciler) clusterToVSphereVMs(ctx context.Context, a ctrlclient.Object) []reconcile.Request {
	requests := []reconcile.Request{}
	vms := &infrav1.VSphereVMList{}
	err := r.Client.List(ctx, vms, ctrlclient.MatchingLabels(
		map[string]string{
			clusterv1.ClusterNameLabel: a.GetName(),
		},
	))
	if err != nil {
		return requests
	}
	for _, vm := range vms.Items {
		r := reconcile.Request{
			NamespacedName: apitypes.NamespacedName{
				Name:      vm.Name,
				Namespace: vm.Namespace,
			},
		}
		requests = append(requests, r)
	}
	return requests
}

func (r vmReconciler) vsphereClusterToVSphereVMs(ctx context.Context, a ctrlclient.Object) []reconcile.Request {
	vsphereCluster, ok := a.(*infrav1.VSphereCluster)
	if !ok {
		return nil
	}
	clusterName, ok := vsphereCluster.Labels[clusterv1.ClusterNameLabel]
	if !ok {
		return nil
	}

	requests := []reconcile.Request{}
	vms := &infrav1.VSphereVMList{}
	err := r.Client.List(ctx, vms, ctrlclient.MatchingLabels(
		map[string]string{
			clusterv1.ClusterNameLabel: clusterName,
		},
	))
	if err != nil {
		return requests
	}
	for _, vm := range vms.Items {
		r := reconcile.Request{
			NamespacedName: apitypes.NamespacedName{
				Name:      vm.Name,
				Namespace: vm.Namespace,
			},
		}
		requests = append(requests, r)
	}
	return requests
}

func (r vmReconciler) ipAddressClaimToVSphereVM(_ context.Context, a ctrlclient.Object) []reconcile.Request {
	ipAddressClaim, ok := a.(*ipamv1.IPAddressClaim)
	if !ok {
		return nil
	}

	requests := []reconcile.Request{}
	if clusterutilv1.HasOwner(ipAddressClaim.OwnerReferences, infrav1.GroupVersion.String(), []string{"VSphereVM"}) {
		for _, ref := range ipAddressClaim.OwnerReferences {
			if ref.Kind == "VSphereVM" {
				requests = append(requests, reconcile.Request{
					NamespacedName: apitypes.NamespacedName{
						Name:      ref.Name,
						Namespace: ipAddressClaim.Namespace,
					},
				})
				break
			}
		}
	}
	return requests
}

func (r vmReconciler) retrieveVcenterSession(ctx context.Context, vsphereVM *infrav1.VSphereVM) (*session.Session, error) {
	log := ctrl.LoggerFrom(ctx)
	// Get cluster object and then get VSphereCluster object

	params := session.NewParams().
		WithServer(vsphereVM.Spec.Server).
		WithDatacenter(vsphereVM.Spec.Datacenter).
		WithUserInfo(r.ControllerManagerContext.Username, r.ControllerManagerContext.Password).
		WithThumbprint(vsphereVM.Spec.Thumbprint)

	cluster, err := clusterutilv1.GetClusterFromMetadata(ctx, r.Client, vsphereVM.ObjectMeta)
	if err != nil {
		log.V(4).Info("Using credentials provided to the manager to create the authenticated session, VSphereVM is missing cluster label or cluster does not exist")
		return session.GetOrCreate(ctx, params)
	}

	if cluster.Spec.InfrastructureRef == nil {
		return nil, errors.Errorf("cannot retrieve vCenter session for cluster %s: Cluster.spec.infrastructureRef is nil", klog.KObj(cluster))
	}
	key := ctrlclient.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	vsphereCluster := &infrav1.VSphereCluster{}
	err = r.Client.Get(ctx, key, vsphereCluster)
	if err != nil {
		log.V(4).Info("Using credentials provided to the manager to create the authenticated session, failed to get VSphereCluster")
		return session.GetOrCreate(ctx, params)
	}

	if vsphereCluster.Spec.IdentityRef != nil {
		creds, err := identity.GetCredentials(ctx, r.Client, vsphereCluster, r.ControllerManagerContext.Namespace)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get credentials from IdentityRef")
		}
		params = params.WithUserInfo(creds.Username, creds.Password)
		return session.GetOrCreate(ctx, params)
	}

	// Fallback to using credentials provided to the manager
	log.V(4).Info("Using credentials provided to the manager to create the authenticated session")
	return session.GetOrCreate(ctx, params)
}

func (r vmReconciler) fetchClusterModuleInfo(ctx context.Context, clusterModInput fetchClusterModuleInput) (*string, error) {
	var (
		owner ctrlclient.Object
		err   error
	)
	log := ctrl.LoggerFrom(ctx)
	machine := clusterModInput.Machine

	input := util.FetchObjectInput{
		Client: r.Client,
		Object: machine,
	}
	// TODO (srm09): Figure out a way to find the latest version of the CRD
	if util.IsControlPlaneMachine(machine) {
		owner, err = util.FetchControlPlaneOwnerObject(ctx, input)
	} else {
		owner, err = util.FetchMachineDeploymentOwnerObject(ctx, input)
	}
	if err != nil {
		// If the owner objects cannot be traced, we can assume that the objects
		// have been deleted in which case we do not want cluster module info populated
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, mod := range clusterModInput.VSphereCluster.Spec.ClusterModules {
		if mod.TargetObjectName == owner.GetName() {
			log.V(4).Info("Cluster module found", "moduleUUID", mod.ModuleUUID)
			return ptr.To(mod.ModuleUUID), nil
		}
	}
	log.V(4).Info("No cluster module found")
	return nil, nil
}

type fetchClusterModuleInput struct {
	VSphereCluster *infrav1.VSphereCluster
	Machine        *clusterv1.Machine
}
