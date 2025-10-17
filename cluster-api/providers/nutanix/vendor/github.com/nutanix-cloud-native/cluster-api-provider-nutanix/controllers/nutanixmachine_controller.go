/*
Copyright 2022 Nutanix

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
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nutanix-cloud-native/prism-go-client/utils"
	prismclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	apitypes "k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/utils/ptr"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capiutil "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	nutanixclient "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/pkg/client"
	nctx "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/pkg/context"
)

const (
	projectKind = "project"

	deviceTypeCDROM = "CDROM"
	adapterTypeIDE  = "IDE"
)

var (
	minMachineSystemDiskSize resource.Quantity
	minMachineDataDiskSize   resource.Quantity
	minMachineMemorySize     resource.Quantity
	minVCPUsPerSocket        = 1
	minVCPUSockets           = 1
)

func init() {
	minMachineSystemDiskSize = resource.MustParse("20Gi")
	minMachineDataDiskSize = resource.MustParse("1Gi")
	minMachineMemorySize = resource.MustParse("2Gi")
}

// NutanixMachineReconciler reconciles a NutanixMachine object
type NutanixMachineReconciler struct {
	client.Client
	SecretInformer    coreinformers.SecretInformer
	ConfigMapInformer coreinformers.ConfigMapInformer
	Scheme            *runtime.Scheme
	controllerConfig  *ControllerConfig
}

func NewNutanixMachineReconciler(client client.Client, secretInformer coreinformers.SecretInformer, configMapInformer coreinformers.ConfigMapInformer, scheme *runtime.Scheme, copts ...ControllerConfigOpts) (*NutanixMachineReconciler, error) {
	controllerConf := &ControllerConfig{}
	for _, opt := range copts {
		if err := opt(controllerConf); err != nil {
			return nil, err
		}
	}

	return &NutanixMachineReconciler{
		Client:            client,
		SecretInformer:    secretInformer,
		ConfigMapInformer: configMapInformer,
		Scheme:            scheme,
		controllerConfig:  controllerConf,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NutanixMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	copts := controller.Options{
		MaxConcurrentReconciles: r.controllerConfig.MaxConcurrentReconciles,
		RateLimiter:             r.controllerConfig.RateLimiter,
		SkipNameValidation:      ptr.To(r.controllerConfig.SkipNameValidation),
	}

	clusterToObjectFunc, err := capiutil.ClusterToTypedObjectsMapper(r.Client, &infrav1.NutanixMachineList{}, mgr.GetScheme())
	if err != nil {
		return fmt.Errorf("failed to create mapper for Cluster to NutanixMachine: %s", err)
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("nutanixmachine-controller").
		For(&infrav1.NutanixMachine{}).
		// Watch the CAPI resource that owns this infrastructure resource.
		Watches(
			&capiv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(
				capiutil.MachineToInfrastructureMapFunc(
					infrav1.GroupVersion.WithKind("NutanixMachine"),
				),
			),
		).
		Watches(
			&infrav1.NutanixCluster{},
			handler.EnqueueRequestsFromMapFunc(
				r.mapNutanixClusterToNutanixMachines(),
			),
		).
		Watches(
			&capiv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(clusterToObjectFunc),
			builder.WithPredicates(predicates.ClusterPausedTransitionsOrInfrastructureReady(r.Scheme, ctrl.LoggerFrom(ctx))),
		).
		WithOptions(copts).
		Complete(r)
}

func (r *NutanixMachineReconciler) mapNutanixClusterToNutanixMachines() handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		log := ctrl.LoggerFrom(ctx)
		nutanixCluster, ok := o.(*infrav1.NutanixCluster)
		if !ok {
			log.Error(fmt.Errorf("expected a NutanixCluster object in mapNutanixClusterToNutanixMachines but was %T", o), "unexpected type")
			return nil
		}

		cluster, err := capiutil.GetOwnerCluster(ctx, r.Client, nutanixCluster.ObjectMeta)
		if apierrors.IsNotFound(err) || cluster == nil {
			log.V(1).Info(fmt.Sprintf("CAPI cluster for NutanixCluster %s not found", nutanixCluster.Name))
			return nil
		}
		if err != nil {
			log.Error(err, "error occurred finding CAPI cluster for NutanixCluster")
			return nil
		}
		searchLabels := map[string]string{capiv1.ClusterNameLabel: cluster.Name}
		machineList := &capiv1.MachineList{}
		if err := r.List(ctx, machineList, client.InNamespace(cluster.Namespace), client.MatchingLabels(searchLabels)); err != nil {
			log.V(1).Error(err, "failed to list machines for cluster")
			return nil
		}
		requests := make([]ctrl.Request, 0)
		for _, m := range machineList.Items {
			if m.Spec.InfrastructureRef.Name == "" || m.Spec.InfrastructureRef.GroupVersionKind().Kind != "NutanixMachine" {
				continue
			}

			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}
			requests = append(requests, ctrl.Request{NamespacedName: name})
		}

		return requests
	}
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update;delete
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update;delete
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixmachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixmachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixmachines/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixclusters,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=kubeadmconfigs,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NutanixMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NutanixMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling the NutanixMachine.")

	// Get the NutanixMachine resource for this request.
	ntxMachine := &infrav1.NutanixMachine{}
	err := r.Get(ctx, req.NamespacedName, ntxMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("NutanixMachine not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "Failed to fetch the NutanixMachine object")
		return reconcile.Result{}, err
	}

	// Fetch the CAPI Machine.
	machine, err := capiutil.GetOwnerMachine(ctx, r.Client, ntxMachine.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to fetch the owner CAPI Machine object")
		return reconcile.Result{}, err
	}
	if machine == nil {
		log.Info("Waiting for capi Machine Controller to set OwnerRef on NutanixMachine")
		return reconcile.Result{}, nil
	}
	log.Info(fmt.Sprintf("Fetched the owner Machine: %s", machine.Name))

	// Fetch the CAPI Cluster.
	cluster, err := capiutil.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Error(err, "Machine is missing cluster label or cluster does not exist")
		return reconcile.Result{}, nil
	}
	if annotations.IsPaused(cluster, machine) {
		log.V(1).Info("linked to a cluster that is paused")
		return reconcile.Result{}, nil
	}

	// Fetch the NutanixCluster
	ntxCluster := &infrav1.NutanixCluster{}
	nclKey := client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	err = r.Get(ctx, nclKey, ntxCluster)
	if err != nil {
		log.Error(err, "Waiting for NutanixCluster")
		return reconcile.Result{}, nil
	}

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(ntxMachine, r.Client)
	if err != nil {
		log.Error(err, "failed to configure the patch helper")
		return ctrl.Result{Requeue: true}, nil
	}

	log.Info(fmt.Sprintf("Reconciling NutanixMachine %s in namespace %s", ntxMachine.Name, ntxMachine.Namespace))
	// Create a Nutanix client for the NutanixCluster.
	v3Client, err := getPrismCentralClientForCluster(ctx, ntxCluster, r.SecretInformer, r.ConfigMapInformer)
	if err != nil {
		log.Error(err, "error occurred while fetching prism central client")
		return reconcile.Result{}, err
	}

	rctx := &nctx.MachineContext{
		Context:        ctx,
		Cluster:        cluster,
		Machine:        machine,
		NutanixCluster: ntxCluster,
		NutanixMachine: ntxMachine,
		NutanixClient:  v3Client,
	}

	defer func() {
		if err == nil {
			// Always attempt to Patch the NutanixMachine object and its status after each reconciliation.
			if err := patchHelper.Patch(ctx, ntxMachine); err != nil {
				log.Error(err, "failed to patch NutanixMachine")
				reterr = kerrors.NewAggregate([]error{reterr, err})
			}
			log.V(1).Info(fmt.Sprintf("Patched NutanixMachine. Spec: %+v. Status: %+v.",
				ntxMachine.Spec, ntxMachine.Status))
		} else {
			log.Error(err, "not patching NutanixMachine since error occurred")
		}
	}()

	// Handle deleted machines
	if !ntxMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(rctx)
	}

	// Handle non-deleted machines
	return r.reconcileNormal(rctx)
}

func (r *NutanixMachineReconciler) reconcileDelete(rctx *nctx.MachineContext) (reconcile.Result, error) {
	ctx := rctx.Context
	log := ctrl.LoggerFrom(ctx)
	v3Client := rctx.NutanixClient
	vmName := rctx.Machine.Name
	log.Info(fmt.Sprintf("Handling deletion of VM: %s", vmName))
	conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, capiv1.DeletingReason, capiv1.ConditionSeverityInfo, "")
	vmUUID, err := GetVMUUID(rctx.NutanixMachine)
	if err != nil {
		errorMsg := fmt.Errorf("failed to get VM UUID during delete: %v", err)
		log.Error(errorMsg, "failed to delete VM")
		return reconcile.Result{}, errorMsg
	}

	// Check if VMUUID is absent
	if vmUUID == "" {
		log.Info(fmt.Sprintf("VM UUID was not found in spec for VM %s. Skipping delete", vmName))
		log.Info(fmt.Sprintf("Removing finalizers for VM %s during delete reconciliation", vmName))
		ctrlutil.RemoveFinalizer(rctx.NutanixMachine, infrav1.NutanixMachineFinalizer)
		ctrlutil.RemoveFinalizer(rctx.NutanixMachine, infrav1.DeprecatedNutanixMachineFinalizer)
		return reconcile.Result{}, nil
	}

	vm, err := FindVMByUUID(ctx, v3Client, vmUUID)
	if err != nil {
		errorMsg := fmt.Errorf("error finding VM %s with UUID %s: %v", vmName, vmUUID, err)
		log.Error(errorMsg, "error finding VM")
		conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.DeletionFailed, capiv1.ConditionSeverityWarning, "%s", errorMsg.Error())
		return reconcile.Result{}, errorMsg
	}

	if vm == nil {
		log.Info(fmt.Sprintf("no VM found with UUID %s: assuming it is already deleted; skipping delete", vmUUID))
		log.Info(fmt.Sprintf("removing finalizers for VM %s during delete reconciliation", vmName))
		ctrlutil.RemoveFinalizer(rctx.NutanixMachine, infrav1.NutanixMachineFinalizer)
		ctrlutil.RemoveFinalizer(rctx.NutanixMachine, infrav1.DeprecatedNutanixMachineFinalizer)
		return reconcile.Result{}, nil
	}

	// Check if the VM name matches the Machine name or the NutanixMachine name.
	// Earlier, we were creating VMs with the same name as the NutanixMachine name.
	// Now, we create VMs with the same name as the Machine name in line with other CAPI providers.
	// This check is to ensure that we are deleting the correct VM for both cases as older CAPX VMs
	// will have the NutanixMachine name as the VM name.
	if *vm.Spec.Name != vmName && *vm.Spec.Name != rctx.NutanixMachine.Name {
		return reconcile.Result{}, fmt.Errorf("found VM with UUID %s but name %s did not match Machine name %s or NutanixMachineName %s", vmUUID, *vm.Spec.Name, vmName, rctx.NutanixMachine.Name)
	}

	log.V(1).Info(fmt.Sprintf("VM %s with UUID %s was found.", *vm.Spec.Name, vmUUID))
	lastTaskUUID, err := GetTaskUUIDFromVM(vm)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred fetching task UUID from VM: %v", err)
		log.Error(errorMsg, "error fetching task UUID")
		conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.DeletionFailed, capiv1.ConditionSeverityWarning, "%s", errorMsg.Error())
		return reconcile.Result{}, errorMsg
	}

	if lastTaskUUID != "" {
		log.Info(fmt.Sprintf("checking if VM %s with UUID %s has in progress tasks", vmName, vmUUID))
		taskInProgress, err := HasTaskInProgress(ctx, rctx.NutanixClient, lastTaskUUID)
		if err != nil {
			log.Error(err, fmt.Sprintf("error occurred while checking task %s for VM %s. Trying to delete VM", lastTaskUUID, vmName))
		}
		if taskInProgress {
			log.Info(fmt.Sprintf("VM %s task with UUID %s still in progress. Requeuing", vmName, vmUUID))
			return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
		}
		log.V(1).Info(fmt.Sprintf("no running tasks anymore... Initiating delete for VM %s with UUID %s", vmName, vmUUID))
	} else {
		log.V(1).Info(fmt.Sprintf("no task UUID found on VM %s. Starting delete.", vmName))
	}

	var vgDetachNeeded bool
	if vm.Spec.Resources != nil && vm.Spec.Resources.DiskList != nil {
		for _, disk := range vm.Spec.Resources.DiskList {
			if disk.VolumeGroupReference != nil {
				vgDetachNeeded = true
				break
			}
		}
	}

	if vgDetachNeeded {
		if err := r.detachVolumeGroups(rctx, vmName, vmUUID, vm.Spec.Resources.DiskList); err != nil {
			err := fmt.Errorf("failed to detach volume groups from VM %s with UUID %s: %v", vmName, vmUUID, err)
			log.Error(err, "failed to detach volume groups from VM")
			conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.VolumeGroupDetachFailed, capiv1.ConditionSeverityWarning, "%s", err.Error())

			return reconcile.Result{}, err
		}

		// Requeue to wait for volume group detach tasks to complete. This is done instead of blocking on task
		// completion to avoid long-running reconcile loops.
		log.Info(fmt.Sprintf("detaching volume groups from VM %s with UUID %s; requeueing again after %s", vmName, vmUUID, detachVGRequeueAfter))
		return reconcile.Result{RequeueAfter: detachVGRequeueAfter}, nil
	}

	// Delete the VM since the VM was found (err was nil)
	deleteTaskUUID, err := DeleteVM(ctx, v3Client, vmName, vmUUID)
	if err != nil {
		err := fmt.Errorf("failed to delete VM %s with UUID %s: %v", vmName, vmUUID, err)
		log.Error(err, "failed to delete VM")
		conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.DeletionFailed, capiv1.ConditionSeverityWarning, "%s", err.Error())

		return reconcile.Result{}, err
	}
	log.Info(fmt.Sprintf("Deletion task with UUID %s received for vm %s with UUID %s. Requeueing", deleteTaskUUID, vmName, vmUUID))
	return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
}

func (r *NutanixMachineReconciler) detachVolumeGroups(rctx *nctx.MachineContext, vmName string, vmUUID string, vmDiskList []*prismclientv3.VMDisk) error {
	v4Client, err := getPrismCentralV4ClientForCluster(rctx.Context, rctx.NutanixCluster, r.SecretInformer, r.ConfigMapInformer)
	if err != nil {
		return fmt.Errorf("error occurred while fetching Prism Central v4 client: %w", err)
	}

	if err := detachVolumeGroupsFromVM(rctx.Context, v4Client, vmName, vmUUID, vmDiskList); err != nil {
		return fmt.Errorf("failed to detach volume groups from VM %s with UUID %s: %w", vmName, vmUUID, err)
	}

	return nil
}

func (r *NutanixMachineReconciler) reconcileNormal(rctx *nctx.MachineContext) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(rctx.Context)
	if rctx.NutanixMachine.Status.FailureReason != nil || rctx.NutanixMachine.Status.FailureMessage != nil {
		log.Error(fmt.Errorf("nutanix machine has failed. Will not reconcile"), "nutanix machine failed")
		return reconcile.Result{}, nil
	}
	log.Info("Handling NutanixMachine reconciling")
	var err error

	// Add finalizer first if not exist to avoid the race condition between init and delete
	if !ctrlutil.ContainsFinalizer(rctx.NutanixMachine, infrav1.NutanixMachineFinalizer) {
		ctrlutil.AddFinalizer(rctx.NutanixMachine, infrav1.NutanixMachineFinalizer)
	}
	ctrlutil.RemoveFinalizer(rctx.NutanixMachine, infrav1.DeprecatedNutanixMachineFinalizer)

	log.V(1).Info(fmt.Sprintf("Checking current machine status for machine %s: Status %+v Spec %+v", rctx.NutanixMachine.Name, rctx.NutanixMachine.Status, rctx.NutanixMachine.Spec))
	if rctx.NutanixMachine.Status.Ready {
		if !rctx.Machine.Status.InfrastructureReady || rctx.Machine.Spec.ProviderID == nil {
			log.Info("The NutanixMachine is ready, wait for the owner Machine's update.")
			return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
		}
		log.Info(fmt.Sprintf("The NutanixMachine is ready, providerID: %s", rctx.NutanixMachine.Spec.ProviderID))

		return reconcile.Result{}, nil
	}

	// Make sure Cluster.Status.InfrastructureReady is true
	log.Info("Checking if cluster infrastructure is ready")
	if !rctx.Cluster.Status.InfrastructureReady {
		log.Info("The cluster infrastructure is not ready yet")
		conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.ClusterInfrastructureNotReady, capiv1.ConditionSeverityInfo, "")
		return reconcile.Result{}, nil
	}

	// Make sure bootstrap data is available and populated.
	if rctx.NutanixMachine.Spec.BootstrapRef == nil {
		if rctx.Machine.Spec.Bootstrap.DataSecretName == nil {
			if !nctx.IsControlPlaneMachine(rctx.NutanixMachine) &&
				!conditions.IsTrue(rctx.Cluster, capiv1.ControlPlaneInitializedCondition) {
				log.Info("Waiting for the control plane to be initialized")
				conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.ControlplaneNotInitialized, capiv1.ConditionSeverityInfo, "")
			} else {
				conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.BootstrapDataNotReady, capiv1.ConditionSeverityInfo, "")
				log.Info("Waiting for bootstrap data to be available")
			}
			return reconcile.Result{}, nil
		}

		rctx.NutanixMachine.Spec.BootstrapRef = &corev1.ObjectReference{
			APIVersion: "v1",
			Kind:       "Secret",
			Name:       *rctx.Machine.Spec.Bootstrap.DataSecretName,
			Namespace:  rctx.Machine.Namespace,
		}
		log.V(1).Info(fmt.Sprintf("Added the spec.bootstrapRef to NutanixMachine object: %v", rctx.NutanixMachine.Spec.BootstrapRef))
	}

	// Create or get existing VM
	vm, err := r.getOrCreateVM(rctx)
	if err != nil {
		log.Error(err, fmt.Sprintf("Failed to create VM %s.", rctx.Machine.Name))
		return reconcile.Result{}, err
	}
	log.V(1).Info(fmt.Sprintf("Found VM with name: %s, vmUUID: %s", rctx.Machine.Name, *vm.Metadata.UUID))
	rctx.NutanixMachine.Status.VmUUID = *vm.Metadata.UUID

	// Set the NutanixMachine.status.failureDomain if the Machine is created with failureDomain
	if err = r.checkFailureDomainStatus(rctx); err != nil {
		log.Error(err, "Failed to check/set status.failureDomain")
		return reconcile.Result{}, err
	}

	log.V(1).Info(fmt.Sprintf("Patching machine post creation vmUUID: %s", rctx.NutanixMachine.Status.VmUUID))
	if err := r.patchMachine(rctx); err != nil {
		errorMsg := fmt.Errorf("failed to patch NutanixMachine %s after creation. %v", rctx.NutanixMachine.Name, err)
		log.Error(errorMsg, "failed to patch")
		return reconcile.Result{}, errorMsg
	}

	log.Info(fmt.Sprintf("Assigning IP addresses to VM with name: %s, vmUUID: %s", rctx.NutanixMachine.Name, rctx.NutanixMachine.Status.VmUUID))
	if err := r.assignAddressesToMachine(rctx, vm); err != nil {
		errorMsg := fmt.Errorf("failed to assign addresses to VM %s with UUID %s...: %v", rctx.Machine.Name, rctx.NutanixMachine.Status.VmUUID, err)
		log.Error(errorMsg, "failed to assign addresses")
		conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMAddressesAssignedCondition, infrav1.VMAddressesFailed, capiv1.ConditionSeverityError, "%s", err.Error())
		return reconcile.Result{}, errorMsg
	}

	conditions.MarkTrue(rctx.NutanixMachine, infrav1.VMAddressesAssignedCondition)
	// Update the NutanixMachine Spec.ProviderID
	rctx.NutanixMachine.Spec.ProviderID = GenerateProviderID(rctx.NutanixMachine.Status.VmUUID)
	rctx.NutanixMachine.Status.Ready = true
	log.V(1).Info(fmt.Sprintf("Created VM %s for cluster %s, update NutanixMachine spec.providerID to %s, and machinespec %+v, vmUuid: %s",
		rctx.Machine.Name, rctx.NutanixCluster.Name, rctx.NutanixMachine.Spec.ProviderID,
		rctx.NutanixMachine, rctx.NutanixMachine.Status.VmUUID))
	return reconcile.Result{}, nil
}

// checkFailureDomainStatus checks and sets the NutanixMachine.status.failureDomain if necessary
func (r *NutanixMachineReconciler) checkFailureDomainStatus(rctx *nctx.MachineContext) error {
	if rctx.Machine.Spec.FailureDomain == nil || *rctx.Machine.Spec.FailureDomain == "" {
		return nil
	}

	fd := *rctx.Machine.Spec.FailureDomain
	// Fetch the referent failure domain object
	fdSpec, err := r.getFailureDomainSpec(rctx, fd)
	if err != nil {
		return err
	}

	// Validate the NutanixMachine machine spec is consistent with that in the failure domain spec
	// Note that when failure domain is used, the cluster/subnets fields of NutanixMachine spec are
	// replaced with that in the failure domain spec, when the machine VM is created.
	errMessages := []string{}
	if !rctx.NutanixMachine.Spec.Cluster.EqualTo(&fdSpec.PrismElementCluster) {
		errMessages = append(
			errMessages,
			fmt.Sprintf(
				"NutanixMachine.spec.cluster=%s, NutanixFailureDomain.spec.prismElementCluster=%s",
				rctx.NutanixMachine.Spec.Cluster.DisplayString(),
				fdSpec.PrismElementCluster.DisplayString(),
			),
		)
	}
	if !resourceIdsEquals(rctx.NutanixMachine.Spec.Subnets, fdSpec.Subnets) {
		errMessages = append(
			errMessages,
			fmt.Sprintf(
				"NutanixMachine.spec.subnets=%v, NutanixFailureDomain.spec.subnets=%v",
				rctx.NutanixMachine.Spec.Subnets,
				fdSpec.Subnets,
			),
		)
	}
	if len(errMessages) > 0 {
		return fmt.Errorf(
			"the NutanixMachine is not consistent with the referenced NutanixFailureDomain %q: %s",
			*rctx.Machine.Spec.FailureDomain,
			strings.Join(errMessages, "; "),
		)
	}

	// Set the NutanixMachine.status.failureDomain
	rctx.NutanixMachine.Status.FailureDomain = &fd

	return nil
}

func (r *NutanixMachineReconciler) getFailureDomainSpec(rctx *nctx.MachineContext, fdName string) (*infrav1.NutanixFailureDomainSpec, error) {
	// TODO: @faiq -- to handle the legacy failure domains this function checks to see if fdName
	// is present in the legacy embedded field. if it is, we return a "dummy" spec for the new failure domain
	// CR with the subnets and cluster info
	failureDomainName := *rctx.Machine.Spec.FailureDomain
	if rctx.NutanixCluster != nil && len(rctx.NutanixCluster.Spec.FailureDomains) > 0 { //nolint:staticcheck // this handles old field
		failureDomain := GetLegacyFailureDomainFromNutanixCluster(failureDomainName, rctx.NutanixCluster)
		if failureDomain != nil {
			cluster := failureDomain.Cluster
			subnets := failureDomain.Subnets
			fdSpec := &infrav1.NutanixFailureDomainSpec{
				PrismElementCluster: cluster,
				Subnets:             subnets,
			}
			return fdSpec, nil
		}
	}
	// if the old field wasn't set or the failure domain name referenced isn't present there, we
	// can assume that it is refering to the new CRD so we make a get
	fdObj := &infrav1.NutanixFailureDomain{}
	fdKey := client.ObjectKey{Name: fdName, Namespace: rctx.NutanixMachine.Namespace}
	if err := r.Get(rctx.Context, fdKey, fdObj); err != nil {
		return nil, fmt.Errorf("failed to fetch the referent failure domain object %q: %w", fdName, err)
	}
	return &fdObj.Spec, nil
}

func (r *NutanixMachineReconciler) validateFailureDomainSpec(rctx *nctx.MachineContext, fdSpec *infrav1.NutanixFailureDomainSpec) error {
	// Validate the failure domain configuration
	pe := fdSpec.PrismElementCluster
	peUUID, err := GetPEUUID(rctx.Context, rctx.NutanixClient, pe.Name, pe.UUID)
	if err != nil {
		return err
	}

	subnets := fdSpec.Subnets
	_, err = GetSubnetUUIDList(rctx.Context, rctx.NutanixClient, subnets, peUUID)
	if err != nil {
		return err
	}

	return nil
}

func (r *NutanixMachineReconciler) validateMachineConfig(rctx *nctx.MachineContext) error {
	log := ctrl.LoggerFrom(rctx.Context)
	fdName := rctx.Machine.Spec.FailureDomain
	if fdName != nil && *fdName != "" {
		log.WithValues("failureDomain", *fdName)
		fdSpec, err := r.getFailureDomainSpec(rctx, *fdName)
		if err != nil {
			log.Error(err, fmt.Sprintf("Failed to get the failure domain %s", *fdName))
			return err
		}
		if err := r.validateFailureDomainSpec(rctx, fdSpec); err != nil {
			log.Error(err, fmt.Sprintf("Failed to validate the failure domain %v", fdSpec))
			return err
		}
		// Update the NutanixMachine machine config based on the failure domain spec
		rctx.NutanixMachine.Spec.Cluster = fdSpec.PrismElementCluster
		rctx.NutanixMachine.Spec.Subnets = fdSpec.Subnets
		rctx.NutanixMachine.Status.FailureDomain = fdName
		log.Info(fmt.Sprintf("Updated the NutanixMachine %s machine config from the failure domain %s configuration.", rctx.NutanixMachine.Name, *fdName))
	}

	if len(rctx.NutanixMachine.Spec.Subnets) == 0 {
		return fmt.Errorf("at least one subnet is needed to create the VM %s", rctx.NutanixMachine.Name)
	}
	if (rctx.NutanixMachine.Spec.Cluster.Name == nil || *rctx.NutanixMachine.Spec.Cluster.Name == "") &&
		(rctx.NutanixMachine.Spec.Cluster.UUID == nil || *rctx.NutanixMachine.Spec.Cluster.UUID == "") {
		return fmt.Errorf("cluster name or uuid are required to create the VM %s", rctx.NutanixMachine.Name)
	}

	diskSize := rctx.NutanixMachine.Spec.SystemDiskSize
	// Validate disk size
	if diskSize.Cmp(minMachineSystemDiskSize) < 0 {
		diskSizeMib := GetMibValueOfQuantity(diskSize)
		minMachineSystemDiskSizeMib := GetMibValueOfQuantity(minMachineSystemDiskSize)
		return fmt.Errorf("minimum systemDiskSize is %vMib but given %vMib", minMachineSystemDiskSizeMib, diskSizeMib)
	}

	memorySize := rctx.NutanixMachine.Spec.MemorySize
	// Validate memory size
	if memorySize.Cmp(minMachineMemorySize) < 0 {
		memorySizeMib := GetMibValueOfQuantity(memorySize)
		minMachineMemorySizeMib := GetMibValueOfQuantity(minMachineMemorySize)
		return fmt.Errorf("minimum memorySize is %vMib but given %vMib", minMachineMemorySizeMib, memorySizeMib)
	}

	vcpusPerSocket := rctx.NutanixMachine.Spec.VCPUsPerSocket
	if vcpusPerSocket < int32(minVCPUsPerSocket) {
		return fmt.Errorf("minimum vcpus per socket is %v but given %v", minVCPUsPerSocket, vcpusPerSocket)
	}

	vcpuSockets := rctx.NutanixMachine.Spec.VCPUSockets
	if vcpuSockets < int32(minVCPUSockets) {
		return fmt.Errorf("minimum vcpu sockets is %v but given %v", minVCPUSockets, vcpuSockets)
	}

	dataDisks := rctx.NutanixMachine.Spec.DataDisks
	if dataDisks != nil {
		if err := r.validateDataDisks(dataDisks); err != nil {
			return err
		}
	}

	return nil
}

func (r *NutanixMachineReconciler) validateDataDisks(dataDisks []infrav1.NutanixMachineVMDisk) error {
	errors := []error{}
	for _, disk := range dataDisks {

		if disk.DiskSize.Cmp(minMachineDataDiskSize) < 0 {
			diskSizeMib := GetMibValueOfQuantity(disk.DiskSize)
			minMachineDataDiskSizeMib := GetMibValueOfQuantity(minMachineDataDiskSize)
			errors = append(errors, fmt.Errorf("minimum data disk size is %vMib but given %vMib", minMachineDataDiskSizeMib, diskSizeMib))
		}

		if disk.DeviceProperties != nil {
			errors = validateDataDiskDeviceProperties(disk, errors)
		}

		if disk.DataSource != nil {
			errors = validateDataDiskDataSource(disk, errors)
		}

		if disk.StorageConfig != nil {
			errors = validateDataDiskStorageConfig(disk, errors)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("data disks validation errors: %v", errors)
	}

	return nil
}

func validateDataDiskStorageConfig(disk infrav1.NutanixMachineVMDisk, errors []error) []error {
	if disk.StorageConfig.StorageContainer != nil && disk.StorageConfig.StorageContainer.IsUUID() {
		if disk.StorageConfig.StorageContainer.UUID == nil {
			errors = append(errors, fmt.Errorf("name or uuid is required for storage container in data disk"))
		} else {
			if _, err := uuid.Parse(*disk.StorageConfig.StorageContainer.UUID); err != nil {
				errors = append(errors, fmt.Errorf("invalid UUID for storage container in data disk: %v", err))
			}
		}
	}

	if disk.StorageConfig.StorageContainer != nil &&
		disk.StorageConfig.StorageContainer.IsName() &&
		disk.StorageConfig.StorageContainer.Name == nil {
		errors = append(errors, fmt.Errorf("name or uuid is required for storage container in data disk"))
	}

	if disk.StorageConfig.DiskMode != infrav1.NutanixMachineDiskModeFlash && disk.StorageConfig.DiskMode != infrav1.NutanixMachineDiskModeStandard {
		errors = append(errors, fmt.Errorf("invalid disk mode %s for data disk", disk.StorageConfig.DiskMode))
	}
	return errors
}

func validateDataDiskDataSource(disk infrav1.NutanixMachineVMDisk, errors []error) []error {
	if disk.DataSource.Type == infrav1.NutanixIdentifierUUID && disk.DataSource.UUID == nil {
		errors = append(errors, fmt.Errorf("UUID is required for data disk with UUID source"))
	}

	if disk.DataSource.Type == infrav1.NutanixIdentifierName && disk.DataSource.Name == nil {
		errors = append(errors, fmt.Errorf("name is required for data disk with name source"))
	}
	return errors
}

func validateDataDiskDeviceProperties(disk infrav1.NutanixMachineVMDisk, errors []error) []error {
	validAdapterTypes := map[infrav1.NutanixMachineDiskAdapterType]bool{
		infrav1.NutanixMachineDiskAdapterTypeIDE:   false,
		infrav1.NutanixMachineDiskAdapterTypeSCSI:  false,
		infrav1.NutanixMachineDiskAdapterTypeSATA:  false,
		infrav1.NutanixMachineDiskAdapterTypePCI:   false,
		infrav1.NutanixMachineDiskAdapterTypeSPAPR: false,
	}

	switch disk.DeviceProperties.DeviceType {
	case infrav1.NutanixMachineDiskDeviceTypeDisk:
		validAdapterTypes[infrav1.NutanixMachineDiskAdapterTypeSCSI] = true
		validAdapterTypes[infrav1.NutanixMachineDiskAdapterTypePCI] = true
		validAdapterTypes[infrav1.NutanixMachineDiskAdapterTypeSPAPR] = true
		validAdapterTypes[infrav1.NutanixMachineDiskAdapterTypeSATA] = true
		validAdapterTypes[infrav1.NutanixMachineDiskAdapterTypeIDE] = true
	case infrav1.NutanixMachineDiskDeviceTypeCDRom:
		validAdapterTypes[infrav1.NutanixMachineDiskAdapterTypeIDE] = true
		validAdapterTypes[infrav1.NutanixMachineDiskAdapterTypePCI] = true
	default:
		errors = append(errors, fmt.Errorf("invalid device type %s for data disk", disk.DeviceProperties.DeviceType))
	}

	if !validAdapterTypes[disk.DeviceProperties.AdapterType] {
		errors = append(errors, fmt.Errorf("invalid adapter type %s for data disk", disk.DeviceProperties.AdapterType))
	}

	if disk.DeviceProperties.DeviceIndex < 0 {
		errors = append(errors, fmt.Errorf("invalid device index %d for data disk", disk.DeviceProperties.DeviceIndex))
	}
	return errors
}

// GetOrCreateVM creates a VM and is invoked by the NutanixMachineReconciler
func (r *NutanixMachineReconciler) getOrCreateVM(rctx *nctx.MachineContext) (*prismclientv3.VMIntentResponse, error) {
	var err error
	var vm *prismclientv3.VMIntentResponse
	ctx := rctx.Context
	log := ctrl.LoggerFrom(ctx)
	vmName := rctx.Machine.Name
	v3Client := rctx.NutanixClient

	// Check if the VM already exists
	vm, err = FindVM(ctx, v3Client, rctx.NutanixMachine, vmName)
	if err != nil {
		log.Error(err, fmt.Sprintf("error occurred finding VM %s by name or uuid", vmName))
		return nil, err
	}

	// if VM exists
	if vm != nil {
		log.Info(fmt.Sprintf("vm %s found with UUID %s", *vm.Spec.Name, rctx.NutanixMachine.Status.VmUUID))
		conditions.MarkTrue(rctx.NutanixMachine, infrav1.VMProvisionedCondition)
		return vm, nil
	}

	log.Info(fmt.Sprintf("No existing VM found. Starting creation process of VM %s.", vmName))
	err = r.validateMachineConfig(rctx)
	if err != nil {
		rctx.SetFailureStatus(createErrorFailureReason, err)
		return nil, err
	}

	peUUID, subnetUUIDs, err := r.GetSubnetAndPEUUIDs(rctx)
	if err != nil {
		log.Error(err, fmt.Sprintf("failed to get the config for VM %s.", vmName))
		rctx.SetFailureStatus(createErrorFailureReason, err)
		return nil, err
	}

	vmInput := &prismclientv3.VMIntentInput{}
	vmSpec := &prismclientv3.VM{Name: utils.StringPtr(vmName)}

	nicList := make([]*prismclientv3.VMNic, len(subnetUUIDs))
	for idx, subnetUUID := range subnetUUIDs {
		nicList[idx] = &prismclientv3.VMNic{
			SubnetReference: &prismclientv3.Reference{
				UUID: utils.StringPtr(subnetUUID),
				Kind: utils.StringPtr("subnet"),
			},
		}
	}

	// Set Categories to VM Sepc before creating VM
	categories, err := GetCategoryVMSpec(ctx, v3Client, r.getMachineCategoryIdentifiers(rctx))
	if err != nil {
		errorMsg := fmt.Errorf("error occurred while creating category spec for vm %s: %v", vmName, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, errorMsg
	}

	vmMetadata := &prismclientv3.Metadata{
		Kind:        utils.StringPtr("vm"),
		SpecVersion: utils.Int64Ptr(1),
		Categories:  categories,
	}
	// Set Project in VM Spec before creating VM
	err = r.addVMToProject(rctx, vmMetadata)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred while trying to add VM %s to project: %v", vmName, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	// Get GPU list
	gpuList, err := GetGPUList(ctx, v3Client, rctx.NutanixMachine.Spec.GPUs, peUUID)
	if err != nil {
		errorMsg := fmt.Errorf("failed to get the GPU list to create the VM %s. %v", vmName, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	diskList, err := getDiskList(rctx, peUUID)
	if err != nil {
		errorMsg := fmt.Errorf("failed to get the disk list to create the VM %s. %v", vmName, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	memorySizeMib := GetMibValueOfQuantity(rctx.NutanixMachine.Spec.MemorySize)
	vmSpec.Resources = &prismclientv3.VMResources{
		PowerState:            utils.StringPtr("ON"),
		HardwareClockTimezone: utils.StringPtr("UTC"),
		NumVcpusPerSocket:     utils.Int64Ptr(int64(rctx.NutanixMachine.Spec.VCPUsPerSocket)),
		NumSockets:            utils.Int64Ptr(int64(rctx.NutanixMachine.Spec.VCPUSockets)),
		MemorySizeMib:         utils.Int64Ptr(memorySizeMib),
		NicList:               nicList,
		DiskList:              diskList,
		GpuList:               gpuList,
	}
	vmSpec.ClusterReference = &prismclientv3.Reference{
		Kind: utils.StringPtr("cluster"),
		UUID: utils.StringPtr(peUUID),
	}

	if err := r.addGuestCustomizationToVM(rctx, vmSpec); err != nil {
		errorMsg := fmt.Errorf("error occurred while adding guest customization to vm spec: %v", err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	// Set BootType in VM Spec before creating VM
	err = r.addBootTypeToVM(rctx, vmSpec)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred while adding boot type to vm spec: %v", err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	vmInput.Spec = vmSpec
	vmInput.Metadata = vmMetadata
	// Create the actual VM/Machine
	log.Info(fmt.Sprintf("Creating VM with name %s for cluster %s", vmName, rctx.NutanixCluster.Name))
	vmResponse, err := v3Client.V3.CreateVM(ctx, vmInput)
	if err != nil {
		errorMsg := fmt.Errorf("failed to create VM %s. error: %v", vmName, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	if vmResponse == nil || vmResponse.Metadata == nil || vmResponse.Metadata.UUID == nil || *vmResponse.Metadata.UUID == "" {
		errorMsg := fmt.Errorf("no valid VM UUID found in response after creating vm %s", rctx.Machine.Name)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, errorMsg
	}
	vmUuid := *vmResponse.Metadata.UUID
	// set the VM UUID on the nutanix machine as soon as it is available. VM UUID can be used for cleanup in case of failure
	rctx.NutanixMachine.Spec.ProviderID = GenerateProviderID(vmUuid)
	rctx.NutanixMachine.Status.VmUUID = vmUuid

	log.V(1).Info(fmt.Sprintf("Sent the post request to create VM %s. Got the vm UUID: %s, status.state: %s", vmName, vmUuid, *vmResponse.Status.State))
	log.V(1).Info(fmt.Sprintf("Getting task vmUUID for VM %s", vmName))
	lastTaskUUID, err := GetTaskUUIDFromVM(vmResponse)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred fetching task UUID from vm %s after creation: %v", rctx.Machine.Name, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, errorMsg
	}

	if lastTaskUUID == "" {
		errorMsg := fmt.Errorf("failed to retrieve task UUID for VM %s after creation", vmName)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, errorMsg
	}

	log.Info(fmt.Sprintf("Waiting for task %s to get completed for VM %s", lastTaskUUID, rctx.NutanixMachine.Name))
	if err := nutanixclient.WaitForTaskToSucceed(ctx, v3Client, lastTaskUUID); err != nil {
		errorMsg := fmt.Errorf("error occurred while waiting for task %s to start: %v", lastTaskUUID, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, errorMsg
	}

	log.Info("Fetching VM after creation")
	vm, err = FindVMByUUID(ctx, v3Client, vmUuid)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred while getting VM %s after creation: %v", vmName, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, errorMsg
	}

	conditions.MarkTrue(rctx.NutanixMachine, infrav1.VMProvisionedCondition)
	return vm, nil
}

func (r *NutanixMachineReconciler) addGuestCustomizationToVM(rctx *nctx.MachineContext, vmSpec *prismclientv3.VM) error {
	// Get the bootstrapData
	bootstrapRef := rctx.NutanixMachine.Spec.BootstrapRef
	if bootstrapRef.Kind == infrav1.NutanixMachineBootstrapRefKindSecret {
		bootstrapData, err := r.getBootstrapData(rctx)
		if err != nil {
			return err
		}

		// Encode the bootstrapData by base64
		bsdataEncoded := base64.StdEncoding.EncodeToString(bootstrapData)
		metadata := fmt.Sprintf("{\"hostname\": \"%s\", \"uuid\": \"%s\"}", rctx.Machine.Name, uuid.New())
		metadataEncoded := base64.StdEncoding.EncodeToString([]byte(metadata))

		vmSpec.Resources.GuestCustomization = &prismclientv3.GuestCustomization{
			IsOverridable: utils.BoolPtr(true),
			CloudInit: &prismclientv3.GuestCustomizationCloudInit{
				UserData: utils.StringPtr(bsdataEncoded),
				MetaData: utils.StringPtr(metadataEncoded),
			},
		}
	}

	return nil
}

func getDiskList(rctx *nctx.MachineContext, peUUID string) ([]*prismclientv3.VMDisk, error) {
	diskList := make([]*prismclientv3.VMDisk, 0)

	systemDisk, err := getSystemDisk(rctx)
	if err != nil {
		return nil, err
	}
	diskList = append(diskList, systemDisk)

	bootstrapRef := rctx.NutanixMachine.Spec.BootstrapRef
	if bootstrapRef.Kind == infrav1.NutanixMachineBootstrapRefKindImage {
		bootstrapDisk, err := getBootstrapDisk(rctx)
		if err != nil {
			return nil, err
		}

		diskList = append(diskList, bootstrapDisk)
	}

	dataDisks, err := getDataDisks(rctx, peUUID)
	if err != nil {
		return nil, err
	}
	diskList = append(diskList, dataDisks...)

	return diskList, nil
}

func getSystemDisk(rctx *nctx.MachineContext) (*prismclientv3.VMDisk, error) {
	var nodeOSImage *prismclientv3.ImageIntentResponse
	var err error
	if rctx.NutanixMachine.Spec.Image != nil {
		nodeOSImage, err = GetImage(
			rctx.Context,
			rctx.NutanixClient,
			*rctx.NutanixMachine.Spec.Image,
		)
	} else if rctx.NutanixMachine.Spec.ImageLookup != nil {
		nodeOSImage, err = GetImageByLookup(
			rctx.Context,
			rctx.NutanixClient,
			rctx.NutanixMachine.Spec.ImageLookup.Format,
			&rctx.NutanixMachine.Spec.ImageLookup.BaseOS,
			rctx.Machine.Spec.Version,
		)
	}
	if err != nil {
		errorMsg := fmt.Errorf("failed to get system disk image %q: %w", rctx.NutanixMachine.Spec.Image, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	// Consider this a precaution. If the image is marked for deletion after we
	// create the "VM create" task, then that task will fail. We will handle that
	// failure separately.
	if ImageMarkedForDeletion(nodeOSImage) {
		err := fmt.Errorf("system disk image %s is being deleted", *nodeOSImage.Metadata.UUID)
		rctx.SetFailureStatus(createErrorFailureReason, err)
		return nil, err
	}

	systemDiskSizeMib := GetMibValueOfQuantity(rctx.NutanixMachine.Spec.SystemDiskSize)
	systemDisk, err := CreateSystemDiskSpec(*nodeOSImage.Metadata.UUID, systemDiskSizeMib)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred while creating system disk spec: %w", err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	return systemDisk, nil
}

func getBootstrapDisk(rctx *nctx.MachineContext) (*prismclientv3.VMDisk, error) {
	bootstrapImageRef := infrav1.NutanixResourceIdentifier{
		Type: infrav1.NutanixIdentifierName,
		Name: ptr.To(rctx.NutanixMachine.Spec.BootstrapRef.Name),
	}
	bootstrapImage, err := GetImage(rctx.Context, rctx.NutanixClient, bootstrapImageRef)
	if err != nil {
		errorMsg := fmt.Errorf("failed to get bootstrap disk image %q: %w", bootstrapImageRef, err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	// Consider this a precaution. If the image is marked for deletion after we
	// create the "VM create" task, then that task will fail. We will handle that
	// failure separately.
	if ImageMarkedForDeletion(bootstrapImage) {
		err := fmt.Errorf("bootstrap disk image %s is being deleted", *bootstrapImage.Metadata.UUID)
		rctx.SetFailureStatus(createErrorFailureReason, err)
		return nil, err
	}

	bootstrapDisk := &prismclientv3.VMDisk{
		DeviceProperties: &prismclientv3.VMDiskDeviceProperties{
			DeviceType: ptr.To(deviceTypeCDROM),
			DiskAddress: &prismclientv3.DiskAddress{
				AdapterType: ptr.To(adapterTypeIDE),
				DeviceIndex: ptr.To(int64(0)),
			},
		},
		DataSourceReference: &prismclientv3.Reference{
			Kind: ptr.To(strings.ToLower(infrav1.NutanixMachineBootstrapRefKindImage)),
			UUID: bootstrapImage.Metadata.UUID,
		},
	}

	return bootstrapDisk, nil
}

func getDataDisks(rctx *nctx.MachineContext, peUUID string) ([]*prismclientv3.VMDisk, error) {
	dataDisks, err := CreateDataDiskList(rctx.Context, rctx.NutanixClient, rctx.NutanixMachine.Spec.DataDisks, peUUID)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred while creating data disk spec: %w", err)
		rctx.SetFailureStatus(createErrorFailureReason, errorMsg)
		return nil, err
	}

	return dataDisks, nil
}

// getBootstrapData returns the Bootstrap data from the ref secret
func (r *NutanixMachineReconciler) getBootstrapData(rctx *nctx.MachineContext) ([]byte, error) {
	if rctx.NutanixMachine.Spec.BootstrapRef == nil {
		return nil, errors.New("NutanixMachine spec.BootstrapRef is nil.")
	}

	secretName := rctx.NutanixMachine.Spec.BootstrapRef.Name
	secret := &corev1.Secret{}
	secretKey := apitypes.NamespacedName{
		Namespace: rctx.NutanixMachine.Spec.BootstrapRef.Namespace,
		Name:      secretName,
	}
	if err := r.Get(rctx.Context, secretKey, secret); err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve bootstrap data secret %s", secretName)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, nil
}

func (r *NutanixMachineReconciler) patchMachine(rctx *nctx.MachineContext) error {
	log := ctrl.LoggerFrom(rctx.Context)
	patchHelper, err := patch.NewHelper(rctx.NutanixMachine, r.Client)
	if err != nil {
		errorMsg := fmt.Errorf("failed to create patch helper to patch machine %s: %v", rctx.NutanixMachine.Name, err)
		return errorMsg
	}
	err = patchHelper.Patch(rctx.Context, rctx.NutanixMachine)
	if err != nil {
		errorMsg := fmt.Errorf("failed to patch machine %s: %v", rctx.NutanixMachine.Name, err)
		return errorMsg
	}
	log.V(1).Info(fmt.Sprintf("Patched machine %s: Status %+v Spec %+v", rctx.NutanixMachine.Name, rctx.NutanixMachine.Status, rctx.NutanixMachine.Spec))
	return nil
}

func (r *NutanixMachineReconciler) assignAddressesToMachine(rctx *nctx.MachineContext, vm *prismclientv3.VMIntentResponse) error {
	rctx.NutanixMachine.Status.Addresses = []capiv1.MachineAddress{}
	if vm.Status == nil || vm.Status.Resources == nil {
		return fmt.Errorf("unable to fetch network interfaces from VM. Retrying")
	}
	foundIPs := 0
	for _, nic := range vm.Status.Resources.NicList {
		for _, ipEndpoint := range nic.IPEndpointList {
			if ipEndpoint.IP != nil {
				rctx.NutanixMachine.Status.Addresses = append(rctx.NutanixMachine.Status.Addresses, capiv1.MachineAddress{
					Type:    capiv1.MachineInternalIP,
					Address: *ipEndpoint.IP,
				})
				foundIPs++
			}
		}
	}
	if foundIPs == 0 {
		return fmt.Errorf("unable to determine network interfaces from VM. Retrying")
	}
	rctx.IP = rctx.NutanixMachine.Status.Addresses[0].Address
	rctx.NutanixMachine.Status.Addresses = append(rctx.NutanixMachine.Status.Addresses, capiv1.MachineAddress{
		Type:    capiv1.MachineHostName,
		Address: *vm.Spec.Name,
	})
	return nil
}

func (r *NutanixMachineReconciler) getMachineCategoryIdentifiers(rctx *nctx.MachineContext) []*infrav1.NutanixCategoryIdentifier {
	log := ctrl.LoggerFrom(rctx.Context)
	categoryIdentifiers := GetDefaultCAPICategoryIdentifiers(rctx.Cluster.Name)
	// Only try to create default categories. ignoring error so that we can return all including
	// additionalCategories as well
	_, err := GetOrCreateCategories(rctx.Context, rctx.NutanixClient, categoryIdentifiers)
	if err != nil {
		log.Error(err, "Failed to getOrCreateCategories")
	}

	additionalCategories := rctx.NutanixMachine.Spec.AdditionalCategories
	if len(additionalCategories) > 0 {
		for _, at := range additionalCategories {
			additionalCat := at
			categoryIdentifiers = append(categoryIdentifiers, &additionalCat)
		}
	}

	return categoryIdentifiers
}

func (r *NutanixMachineReconciler) addBootTypeToVM(rctx *nctx.MachineContext, vmSpec *prismclientv3.VM) error {
	bootType := rctx.NutanixMachine.Spec.BootType
	// Defaults to legacy if boot type is not set.
	if bootType != "" {
		if bootType != infrav1.NutanixBootTypeLegacy && bootType != infrav1.NutanixBootTypeUEFI {
			errorMsg := fmt.Errorf("boot type must be %s or %s but was %s", string(infrav1.NutanixBootTypeLegacy), string(infrav1.NutanixBootTypeUEFI), bootType)
			conditions.MarkFalse(rctx.NutanixMachine, infrav1.VMProvisionedCondition, infrav1.VMBootTypeInvalid, capiv1.ConditionSeverityError, "%s", errorMsg.Error())
			return errorMsg
		}

		// Only modify VM spec if boot type is UEFI. Otherwise, assume default Legacy mode
		if bootType == infrav1.NutanixBootTypeUEFI {
			vmSpec.Resources.BootConfig = &prismclientv3.VMBootConfig{
				BootType: utils.StringPtr(strings.ToUpper(string(bootType))),
			}
		}
	}

	return nil
}

func (r *NutanixMachineReconciler) addVMToProject(rctx *nctx.MachineContext, vmMetadata *prismclientv3.Metadata) error {
	log := ctrl.LoggerFrom(rctx.Context)
	vmName := rctx.Machine.Name
	projectRef := rctx.NutanixMachine.Spec.Project
	if projectRef == nil {
		log.V(1).Info("Not linking VM to a project")
		return nil
	}

	if vmMetadata == nil {
		errorMsg := fmt.Errorf("metadata cannot be nil when adding VM %s to project", vmName)
		log.Error(errorMsg, "failed to add vm to project")
		conditions.MarkFalse(rctx.NutanixMachine, infrav1.ProjectAssignedCondition, infrav1.ProjectAssignationFailed, capiv1.ConditionSeverityError, "%s", errorMsg.Error())
		return errorMsg
	}

	projectUUID, err := GetProjectUUID(rctx.Context, rctx.NutanixClient, projectRef.Name, projectRef.UUID)
	if err != nil {
		errorMsg := fmt.Errorf("error occurred while searching for project for VM %s: %v", vmName, err)
		log.Error(errorMsg, "error occurred while searching for project")
		conditions.MarkFalse(rctx.NutanixMachine, infrav1.ProjectAssignedCondition, infrav1.ProjectAssignationFailed, capiv1.ConditionSeverityError, "%s", errorMsg.Error())
		return errorMsg
	}

	vmMetadata.ProjectReference = &prismclientv3.Reference{
		Kind: utils.StringPtr(projectKind),
		UUID: utils.StringPtr(projectUUID),
	}
	conditions.MarkTrue(rctx.NutanixMachine, infrav1.ProjectAssignedCondition)
	return nil
}

func (r *NutanixMachineReconciler) GetSubnetAndPEUUIDs(rctx *nctx.MachineContext) (string, []string, error) {
	if rctx == nil {
		return "", nil, fmt.Errorf("cannot create machine config if machine context is nil")
	}

	peUUID, err := GetPEUUID(rctx.Context, rctx.NutanixClient, rctx.NutanixMachine.Spec.Cluster.Name, rctx.NutanixMachine.Spec.Cluster.UUID)
	if err != nil {
		return "", nil, err
	}

	subnetUUIDs, err := GetSubnetUUIDList(rctx.Context, rctx.NutanixClient, rctx.NutanixMachine.Spec.Subnets, peUUID)
	if err != nil {
		return "", nil, err
	}

	return peUUID, subnetUUIDs, nil
}
