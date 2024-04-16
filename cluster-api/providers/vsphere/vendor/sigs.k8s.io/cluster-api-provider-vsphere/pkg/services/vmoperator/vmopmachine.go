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

package vmoperator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	infrautilv1 "sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// VmopMachineService reconciles VM Operator VM.
type VmopMachineService struct {
	Client client.Client
}

// GetMachinesInCluster returns a list of VSphereMachine objects belonging to the cluster.
func (v *VmopMachineService) GetMachinesInCluster(
	ctx context.Context,
	namespace, clusterName string) ([]client.Object, error) {
	labels := map[string]string{clusterv1.ClusterNameLabel: clusterName}
	machineList := &vmwarev1.VSphereMachineList{}

	if err := v.Client.List(
		ctx, machineList,
		client.InNamespace(namespace),
		client.MatchingLabels(labels)); err != nil {
		return nil, err
	}

	objects := []client.Object{}
	for _, machine := range machineList.Items {
		m := machine
		objects = append(objects, &m)
	}
	return objects, nil
}

// FetchVSphereMachine returns a MachineContext with a VSphereMachine for the passed NamespacedName.
func (v *VmopMachineService) FetchVSphereMachine(ctx context.Context, name apitypes.NamespacedName) (capvcontext.MachineContext, error) {
	vsphereMachine := &vmwarev1.VSphereMachine{}
	err := v.Client.Get(ctx, name, vsphereMachine)
	return &vmware.SupervisorMachineContext{VSphereMachine: vsphereMachine}, err
}

// FetchVSphereCluster adds the VSphereCluster for the cluster to the MachineContext.
func (v *VmopMachineService) FetchVSphereCluster(ctx context.Context, cluster *clusterv1.Cluster, machineContext capvcontext.MachineContext) (capvcontext.MachineContext, error) {
	machineCtx, ok := machineContext.(*vmware.SupervisorMachineContext)
	if !ok {
		return nil, errors.New("received unexpected SupervisorMachineContext type")
	}

	vsphereCluster := &vmwarev1.VSphereCluster{}
	key := client.ObjectKey{
		Namespace: machineContext.GetObjectMeta().Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	err := v.Client.Get(ctx, key, vsphereCluster)

	machineCtx.VSphereCluster = vsphereCluster
	return machineCtx, err
}

// ReconcileDelete reconciles delete events for VM Operator VM.
func (v *VmopMachineService) ReconcileDelete(ctx context.Context, machineCtx capvcontext.MachineContext) error {
	log := ctrl.LoggerFrom(ctx)
	supervisorMachineCtx, ok := machineCtx.(*vmware.SupervisorMachineContext)
	if !ok {
		return errors.New("received unexpected SupervisorMachineContext type")
	}
	log.Info("Destroying VM")

	// If debug logging is enabled, report the number of vms in the cluster before and after the reconcile
	if log.V(5).Enabled() {
		vms, err := v.getVirtualMachinesInCluster(ctx, supervisorMachineCtx)
		log.V(5).Info("Trace Destroy PRE: VirtualMachines", "vmcount", len(vms), "err", err)
		defer func() {
			vms, err := v.getVirtualMachinesInCluster(ctx, supervisorMachineCtx)
			log.V(5).Info("Trace Destroy POST: VirtualMachines", "vmcount", len(vms), "err", err)
		}()
	}

	// First, check to see if it's already deleted
	vmopVM := vmoprv1.VirtualMachine{}
	if err := v.Client.Get(ctx, apitypes.NamespacedName{Namespace: supervisorMachineCtx.Machine.Namespace, Name: supervisorMachineCtx.Machine.Name}, &vmopVM); err != nil {
		// If debug logging is enabled, report the number of vms in the cluster before and after the reconcile
		if apierrors.IsNotFound(err) {
			supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateNotFound
			return err
		}
		supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateError
		return err
	}

	// Next, check to see if it's in the process of being deleted
	if vmopVM.GetDeletionTimestamp() != nil {
		supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateDeleting
		return nil
	}

	// If none of the above are true, Delete the VM
	if err := v.Client.Delete(ctx, &vmopVM); err != nil {
		if apierrors.IsNotFound(err) {
			supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateNotFound
			return err
		}
		supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateError
		return err
	}
	supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateDeleting
	return nil
}

// SyncFailureReason returns true if there is a Failure on the VM Operator VM.
func (v *VmopMachineService) SyncFailureReason(_ context.Context, machineCtx capvcontext.MachineContext) (bool, error) {
	supervisorMachineCtx, ok := machineCtx.(*vmware.SupervisorMachineContext)
	if !ok {
		return false, errors.New("received unexpected SupervisorMachineContext type")
	}

	return supervisorMachineCtx.VSphereMachine.Status.FailureReason != nil || supervisorMachineCtx.VSphereMachine.Status.FailureMessage != nil, nil
}

// ReconcileNormal reconciles create and update events for VM Operator VMs.
func (v *VmopMachineService) ReconcileNormal(ctx context.Context, machineCtx capvcontext.MachineContext) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	supervisorMachineCtx, ok := machineCtx.(*vmware.SupervisorMachineContext)
	if !ok {
		return false, errors.New("received unexpected SupervisorMachineContext type")
	}

	supervisorMachineCtx.VSphereMachine.Spec.FailureDomain = supervisorMachineCtx.Machine.Spec.FailureDomain

	// If debug logging is enabled, report the number of vms in the cluster before and after the reconcile
	if log.V(5).Enabled() {
		vms, err := v.getVirtualMachinesInCluster(ctx, supervisorMachineCtx)
		log.V(5).Info("Trace ReconcileVM PRE: VirtualMachines", "vmcount", len(vms), "err", err)
		defer func() {
			vms, err = v.getVirtualMachinesInCluster(ctx, supervisorMachineCtx)
			log.V(5).Info("Trace ReconcileVM POST: VirtualMachines", "vmcount", len(vms), "err", err)
		}()
	}

	// Set the VM state. Will get reset throughout the reconcile
	supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStatePending

	// Check for the presence of an existing object
	vmOperatorVM := &vmoprv1.VirtualMachine{}
	if err := v.Client.Get(ctx, client.ObjectKey{
		Namespace: supervisorMachineCtx.Machine.Namespace,
		Name:      supervisorMachineCtx.Machine.Name,
	}, vmOperatorVM); err != nil {
		if !apierrors.IsNotFound(err) {
			return false, err
		}
		// Define the VM Operator VirtualMachine resource to reconcile.
		vmOperatorVM = v.newVMOperatorVM(supervisorMachineCtx)
	}

	// Reconcile the VM Operator VirtualMachine.
	if err := v.reconcileVMOperatorVM(ctx, supervisorMachineCtx, vmOperatorVM); err != nil {
		conditions.MarkFalse(supervisorMachineCtx.VSphereMachine, infrav1.VMProvisionedCondition, vmwarev1.VMCreationFailedReason, clusterv1.ConditionSeverityWarning,
			fmt.Sprintf("failed to create or update VirtualMachine: %v", err))
		// TODO: what to do if AlreadyExists error
		return false, err
	}

	// Update the VM's state to Pending
	supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStatePending

	// Since vm operator only has one condition for now, we can't set vspheremachine's condition fully based on virtualmachine's
	// condition. Once vm operator surfaces enough conditions in virtualmachine, we could simply mirror the conditions in vspheremachine.
	// For now, we set conditions based on the whole virtualmachine status.
	// TODO: vm-operator does not use the cluster-api condition type. so can't use cluster-api util functions to fetch the condition
	for _, cond := range vmOperatorVM.Status.Conditions {
		if cond.Type == vmoprv1.VirtualMachinePrereqReadyCondition && cond.Severity == vmoprv1.ConditionSeverityError {
			conditions.MarkFalse(supervisorMachineCtx.VSphereMachine, infrav1.VMProvisionedCondition, cond.Reason, clusterv1.ConditionSeverityError, cond.Message)
			return false, errors.Errorf("vm prerequisites check fails: %s", supervisorMachineCtx)
		}
	}

	// Requeue until the VM Operator VirtualMachine has:
	// * Been created
	// * Been powered on
	// * An IP address
	// * A BIOS UUID
	if vmOperatorVM.Status.Phase != vmoprv1.Created {
		conditions.MarkFalse(supervisorMachineCtx.VSphereMachine, infrav1.VMProvisionedCondition, vmwarev1.VMProvisionStartedReason, clusterv1.ConditionSeverityInfo, "")
		log.Info(fmt.Sprintf("VM is not yet created: %s", supervisorMachineCtx))
		return true, nil
	}
	// Mark the VM as created
	supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateCreated

	if vmOperatorVM.Status.PowerState != vmoprv1.VirtualMachinePoweredOn {
		conditions.MarkFalse(supervisorMachineCtx.VSphereMachine, infrav1.VMProvisionedCondition, vmwarev1.PoweringOnReason, clusterv1.ConditionSeverityInfo, "")
		log.Info(fmt.Sprintf("VM is not yet powered on: %s", supervisorMachineCtx))
		return true, nil
	}
	// Mark the VM as poweredOn
	supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStatePoweredOn

	if vmOperatorVM.Status.VmIp == "" {
		conditions.MarkFalse(supervisorMachineCtx.VSphereMachine, infrav1.VMProvisionedCondition, vmwarev1.WaitingForNetworkAddressReason, clusterv1.ConditionSeverityInfo, "")
		log.Info(fmt.Sprintf("VM does not have an IP address: %s", supervisorMachineCtx))
		return true, nil
	}

	if vmOperatorVM.Status.BiosUUID == "" {
		conditions.MarkFalse(supervisorMachineCtx.VSphereMachine, infrav1.VMProvisionedCondition, vmwarev1.WaitingForBIOSUUIDReason, clusterv1.ConditionSeverityInfo, "")
		log.Info(fmt.Sprintf("VM does not have a BIOS UUID: %s", supervisorMachineCtx))
		return true, nil
	}

	// Mark the VM as ready
	supervisorMachineCtx.VSphereMachine.Status.VMStatus = vmwarev1.VirtualMachineStateReady

	if ok := v.reconcileNetwork(supervisorMachineCtx, vmOperatorVM); !ok {
		log.Info("IP not yet assigned")
		return true, nil
	}

	v.reconcileProviderID(ctx, supervisorMachineCtx, vmOperatorVM)

	// Mark the VSphereMachine as Ready
	supervisorMachineCtx.VSphereMachine.Status.Ready = true
	conditions.MarkTrue(supervisorMachineCtx.VSphereMachine, infrav1.VMProvisionedCondition)
	return false, nil
}

// GetHostInfo returns the hostname or IP address of the infrastructure host for the VM Operator VM.
func (v *VmopMachineService) GetHostInfo(ctx context.Context, machineCtx capvcontext.MachineContext) (string, error) {
	supervisorMachineCtx, ok := machineCtx.(*vmware.SupervisorMachineContext)
	if !ok {
		return "", errors.New("received unexpected SupervisorMachineContext type")
	}

	vmOperatorVM := &vmoprv1.VirtualMachine{}
	if err := v.Client.Get(ctx, client.ObjectKey{
		Name:      supervisorMachineCtx.Machine.Name,
		Namespace: supervisorMachineCtx.Machine.Namespace,
	}, vmOperatorVM); err != nil {
		return "", err
	}

	return vmOperatorVM.Status.Host, nil
}

func (v *VmopMachineService) newVMOperatorVM(supervisorMachineCtx *vmware.SupervisorMachineContext) *vmoprv1.VirtualMachine {
	return &vmoprv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      supervisorMachineCtx.Machine.Name,
			Namespace: supervisorMachineCtx.Machine.Namespace,
		},
	}
}

func (v *VmopMachineService) reconcileVMOperatorVM(ctx context.Context, supervisorMachineCtx *vmware.SupervisorMachineContext, vmOperatorVM *vmoprv1.VirtualMachine) error {
	// All Machine resources should define the version of Kubernetes to use.
	if supervisorMachineCtx.Machine.Spec.Version == nil || *supervisorMachineCtx.Machine.Spec.Version == "" {
		return errors.Errorf(
			"missing kubernetes version for %s %s/%s",
			supervisorMachineCtx.Machine.GroupVersionKind(),
			supervisorMachineCtx.Machine.Namespace,
			supervisorMachineCtx.Machine.Name)
	}

	var dataSecretName string
	if dsn := supervisorMachineCtx.Machine.Spec.Bootstrap.DataSecretName; dsn != nil {
		dataSecretName = *dsn
	}

	var minHardwareVersion int32
	if version := supervisorMachineCtx.VSphereMachine.Spec.MinHardwareVersion; version != "" {
		hwVersion, err := infrautilv1.ParseHardwareVersion(version)
		if err != nil {
			return err
		}
		minHardwareVersion = int32(hwVersion)
	}

	_, err := ctrlutil.CreateOrPatch(ctx, v.Client, vmOperatorVM, func() error {
		// Define a new VM Operator virtual machine.
		// NOTE: Set field-by-field in order to preserve changes made directly
		//  to the VirtualMachine spec by other sources (e.g. the cloud provider)
		if vmOperatorVM.Spec.ImageName == "" {
			vmOperatorVM.Spec.ImageName = supervisorMachineCtx.VSphereMachine.Spec.ImageName
		}
		if vmOperatorVM.Spec.ClassName == "" {
			vmOperatorVM.Spec.ClassName = supervisorMachineCtx.VSphereMachine.Spec.ClassName
		}
		if vmOperatorVM.Spec.StorageClass == "" {
			vmOperatorVM.Spec.StorageClass = supervisorMachineCtx.VSphereMachine.Spec.StorageClass
		}
		vmOperatorVM.Spec.PowerState = vmoprv1.VirtualMachinePoweredOn
		if vmOperatorVM.Spec.ResourcePolicyName == "" {
			vmOperatorVM.Spec.ResourcePolicyName = supervisorMachineCtx.VSphereCluster.Status.ResourcePolicyName
		}
		vmOperatorVM.Spec.VmMetadata = &vmoprv1.VirtualMachineMetadata{
			SecretName: dataSecretName,
			Transport:  vmoprv1.VirtualMachineMetadataCloudInitTransport,
		}
		vmOperatorVM.Spec.PowerOffMode = vmoprv1.VirtualMachinePowerOpMode(supervisorMachineCtx.VSphereMachine.Spec.PowerOffMode)
		if vmOperatorVM.Spec.MinHardwareVersion == 0 {
			vmOperatorVM.Spec.MinHardwareVersion = minHardwareVersion
		}

		// VMOperator supports readiness probe and will add/remove endpoints to a
		// VirtualMachineService based on the outcome of the readiness check.
		// When creating the initial control plane node, we do not declare a probe
		// in order to reduce the likelihood of a race between the VirtualMachineService
		// endpoint additions and the kubeadm commands run on the VM itself.
		// Once the initial control plane node is ready, we can re-add the probe so
		// that subsequent machines do not attempt to speak to a kube-apiserver
		// that is not yet ready.
		if infrautilv1.IsControlPlaneMachine(supervisorMachineCtx.Machine) && supervisorMachineCtx.Cluster.Status.ControlPlaneReady {
			vmOperatorVM.Spec.ReadinessProbe = &vmoprv1.Probe{
				TCPSocket: &vmoprv1.TCPSocketAction{
					Port: intstr.FromInt(defaultAPIBindPort),
				},
			}
		}

		// Assign the VM's labels.
		vmOperatorVM.Labels = getVMLabels(supervisorMachineCtx, vmOperatorVM.Labels)

		addResourcePolicyAnnotations(supervisorMachineCtx, vmOperatorVM)

		if err := v.addVolumes(ctx, supervisorMachineCtx, vmOperatorVM); err != nil {
			return err
		}

		// Apply hooks to modify the VM spec
		// The hooks are loosely typed so as to allow for different VirtualMachine backends
		for _, vmModifier := range supervisorMachineCtx.VMModifiers {
			modified, err := vmModifier(vmOperatorVM)
			if err != nil {
				return err
			}
			typedModified, ok := modified.(*vmoprv1.VirtualMachine)
			if !ok {
				return fmt.Errorf("VM modifier returned result of the wrong type: %T", typedModified)
			}
			vmOperatorVM = typedModified
		}

		// Make sure the VSphereMachine owns the VM Operator VirtualMachine.
		if err := ctrlutil.SetControllerReference(supervisorMachineCtx.VSphereMachine, vmOperatorVM, v.Client.Scheme()); err != nil {
			return errors.Wrapf(err, "failed to mark %s %s/%s as owner of %s %s/%s",
				supervisorMachineCtx.VSphereMachine.GroupVersionKind(),
				supervisorMachineCtx.VSphereMachine.Namespace,
				supervisorMachineCtx.VSphereMachine.Name,
				vmOperatorVM.GroupVersionKind(),
				vmOperatorVM.Namespace,
				vmOperatorVM.Name)
		}

		return nil
	})
	return err
}

func (v *VmopMachineService) reconcileNetwork(supervisorMachineCtx *vmware.SupervisorMachineContext, vm *vmoprv1.VirtualMachine) bool {
	if vm.Status.VmIp == "" {
		return false
	}

	supervisorMachineCtx.VSphereMachine.Status.IPAddr = vm.Status.VmIp

	return true
}

func (v *VmopMachineService) reconcileProviderID(ctx context.Context, supervisorMachineCtx *vmware.SupervisorMachineContext, vm *vmoprv1.VirtualMachine) {
	log := ctrl.LoggerFrom(ctx)
	providerID := fmt.Sprintf("vsphere://%s", vm.Status.BiosUUID)

	if supervisorMachineCtx.VSphereMachine.Spec.ProviderID == nil || *supervisorMachineCtx.VSphereMachine.Spec.ProviderID != providerID {
		supervisorMachineCtx.VSphereMachine.Spec.ProviderID = &providerID
		log.Info("Updated providerID", "providerID", providerID)
	}

	if supervisorMachineCtx.VSphereMachine.Status.ID == nil || *supervisorMachineCtx.VSphereMachine.Status.ID != vm.Status.BiosUUID {
		supervisorMachineCtx.VSphereMachine.Status.ID = &vm.Status.BiosUUID
		log.Info("Updated VM ID", "vmID", vm.Status.BiosUUID)
	}
}

// getVirtualMachinesInCluster returns all VMOperator VirtualMachine objects in the current cluster.
// First filter by clusterSelectorKey. If the result is empty, they fall back to legacyClusterSelectorKey.
func (v *VmopMachineService) getVirtualMachinesInCluster(ctx context.Context, supervisorMachineCtx *vmware.SupervisorMachineContext) ([]*vmoprv1.VirtualMachine, error) {
	if supervisorMachineCtx.Cluster == nil {
		return []*vmoprv1.VirtualMachine{}, errors.Errorf("No cluster is set for machine %s in namespace %s", supervisorMachineCtx.GetVSphereMachine().GetName(), supervisorMachineCtx.GetVSphereMachine().GetNamespace())
	}
	labels := map[string]string{clusterSelectorKey: supervisorMachineCtx.Cluster.Name}
	vmList := &vmoprv1.VirtualMachineList{}

	if err := v.Client.List(
		ctx, vmList,
		client.InNamespace(supervisorMachineCtx.Cluster.Namespace),
		client.MatchingLabels(labels)); err != nil {
		return nil, errors.Wrapf(
			err, "error getting virtualmachines in cluster %s/%s",
			supervisorMachineCtx.Cluster.Namespace, supervisorMachineCtx.Cluster.Name)
	}

	// If the list is empty, fall back to usse legacy labels for filtering
	if len(vmList.Items) == 0 {
		legacyLabels := map[string]string{legacyClusterSelectorKey: supervisorMachineCtx.Cluster.Name}
		if err := v.Client.List(
			ctx, vmList,
			client.InNamespace(supervisorMachineCtx.Cluster.Namespace),
			client.MatchingLabels(legacyLabels)); err != nil {
			return nil, errors.Wrapf(
				err, "error getting virtualmachines in cluster %s/%s using legacy labels",
				supervisorMachineCtx.Cluster.Namespace, supervisorMachineCtx.Cluster.Name)
		}
	}

	vms := make([]*vmoprv1.VirtualMachine, len(vmList.Items))
	for i := range vmList.Items {
		vms[i] = &vmList.Items[i]
	}

	return vms, nil
}

// Helper function to add annotations to indicate which tag vm-operator should add as well as which clusterModule VM
// should be associated.
func addResourcePolicyAnnotations(supervisorMachineCtx *vmware.SupervisorMachineContext, vm *vmoprv1.VirtualMachine) {
	annotations := vm.ObjectMeta.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	if infrautilv1.IsControlPlaneMachine(supervisorMachineCtx.Machine) {
		annotations[ProviderTagsAnnotationKey] = ControlPlaneVMVMAntiAffinityTagValue
		annotations[ClusterModuleNameAnnotationKey] = ControlPlaneVMClusterModuleGroupName
	} else {
		annotations[ProviderTagsAnnotationKey] = WorkerVMVMAntiAffinityTagValue
		annotations[ClusterModuleNameAnnotationKey] = getMachineDeploymentNameForCluster(supervisorMachineCtx.Cluster)
	}

	vm.ObjectMeta.SetAnnotations(annotations)
}

func volumeName(machine *vmwarev1.VSphereMachine, volume vmwarev1.VSphereMachineVolume) string {
	return machine.Name + "-" + volume.Name
}

// addVolume ensures volume is included in vm.Spec.Volumes.
func addVolume(vm *vmoprv1.VirtualMachine, name string) {
	for _, volume := range vm.Spec.Volumes {
		claim := volume.PersistentVolumeClaim
		if claim != nil && claim.ClaimName == name {
			return // volume already present in the spec
		}
	}

	vm.Spec.Volumes = append(vm.Spec.Volumes, vmoprv1.VirtualMachineVolume{
		Name: name,
		PersistentVolumeClaim: &vmoprv1.PersistentVolumeClaimVolumeSource{
			PersistentVolumeClaimVolumeSource: corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: name,
				ReadOnly:  false,
			},
		},
	})
}

func (v *VmopMachineService) addVolumes(ctx context.Context, supervisorMachineCtx *vmware.SupervisorMachineContext, vm *vmoprv1.VirtualMachine) error {
	nvolumes := len(supervisorMachineCtx.VSphereMachine.Spec.Volumes)
	if nvolumes == 0 {
		return nil
	}

	for _, volume := range supervisorMachineCtx.VSphereMachine.Spec.Volumes {
		storageClassName := volume.StorageClass
		if volume.StorageClass == "" {
			storageClassName = supervisorMachineCtx.VSphereMachine.Spec.StorageClass
		}

		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      volumeName(supervisorMachineCtx.VSphereMachine, volume),
				Namespace: supervisorMachineCtx.VSphereMachine.Namespace,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.ResourceRequirements{
					Requests: volume.Capacity,
				},
				StorageClassName: &storageClassName,
			},
		}

		// The CSI zone annotation must be set when using a zonal storage class,
		// which is required when the cluster has multiple (3) zones.
		// Single zone clusters (legacy/default) do not support zonal storage and must not
		// have the zone annotation set.
		zonal := len(supervisorMachineCtx.VSphereCluster.Status.FailureDomains) > 1

		if zone := supervisorMachineCtx.VSphereMachine.Spec.FailureDomain; zonal && zone != nil {
			topology := []map[string]string{
				{kubeTopologyZoneLabelKey: *zone},
			}
			b, err := json.Marshal(topology)
			if err != nil {
				return errors.Errorf("failed to marshal zone topology %q: %s", *zone, err)
			}
			pvc.Annotations = map[string]string{
				"csi.vsphere.volume-requested-topology": string(b),
			}
		}

		if _, err := ctrlutil.CreateOrPatch(ctx, v.Client, pvc, func() error {
			if err := ctrlutil.SetOwnerReference(
				supervisorMachineCtx.VSphereMachine,
				pvc,
				v.Client.Scheme(),
			); err != nil {
				return errors.Wrapf(
					err,
					"error setting %s/%s as owner of %s/%s",
					supervisorMachineCtx.VSphereMachine.Namespace,
					supervisorMachineCtx.VSphereMachine.Name,
					pvc.Namespace,
					pvc.Name,
				)
			}
			return nil
		}); err != nil {
			return errors.Wrapf(
				err,
				"failed to create volume %s/%s",
				pvc.Namespace,
				pvc.Name)
		}

		addVolume(vm, pvc.Name)
	}

	return nil
}

// getVMLabels returns the labels applied to a VirtualMachine.
func getVMLabels(supervisorMachineCtx *vmware.SupervisorMachineContext, vmLabels map[string]string) map[string]string {
	if vmLabels == nil {
		vmLabels = map[string]string{}
	}

	// Get the labels for the VM that differ based on the cluster role of
	// the Kubernetes node hosted on this VM.
	clusterRoleLabels := clusterRoleVMLabels(supervisorMachineCtx.GetClusterContext(), infrautilv1.IsControlPlaneMachine(supervisorMachineCtx.Machine))
	for k, v := range clusterRoleLabels {
		vmLabels[k] = v
	}

	// Get the labels that determine the VM's placement inside of a stretched
	// cluster.
	topologyLabels := getTopologyLabels(supervisorMachineCtx)
	for k, v := range topologyLabels {
		vmLabels[k] = v
	}

	// Ensure the VM has a label that can be used when searching for
	// resources associated with the target cluster.
	vmLabels[clusterv1.ClusterNameLabel] = supervisorMachineCtx.GetClusterContext().Cluster.Name

	return vmLabels
}

// getTopologyLabels returns the labels related to a VM's topology.
//
// TODO(akutz): Currently this function just returns the availability zone,
//
//	and thus the code is optimized as such. However, in the future
//	this function may return a more diverse topology.
func getTopologyLabels(supervisorMachineCtx *vmware.SupervisorMachineContext) map[string]string {
	if fd := supervisorMachineCtx.VSphereMachine.Spec.FailureDomain; fd != nil && *fd != "" {
		return map[string]string{
			kubeTopologyZoneLabelKey: *fd,
		}
	}
	return nil
}

// getMachineDeploymentName returns the MachineDeployment name for a Cluster.
// This is also the name used by VSphereMachineTemplate and KubeadmConfigTemplate.
func getMachineDeploymentNameForCluster(cluster *clusterv1.Cluster) string {
	return fmt.Sprintf("%s-workers-0", cluster.Name)
}
