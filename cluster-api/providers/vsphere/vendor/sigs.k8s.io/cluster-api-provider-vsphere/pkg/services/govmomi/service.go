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

package govmomi

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/pbm"
	pbmTypes "github.com/vmware/govmomi/pbm/types"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	corev1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/cluster"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/clustermodules"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/extra"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/ipam"
	govmominet "sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/net"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/pci"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// VMService provdes API to interact with the VMs using govmomi.
type VMService struct{}

// ReconcileVM makes sure that the VM is in the desired state by:
//  1. Creating the VM if it does not exist, then...
//  2. Updating the VM with the bootstrap data, such as the cloud-init meta and user data, before...
//  3. Powering on the VM, and finally...
//  4. Returning the real-time state of the VM to the caller
func (vms *VMService) ReconcileVM(ctx context.Context, vmCtx *capvcontext.VMContext) (vm infrav1.VirtualMachine, _ error) {
	// Initialize the result.
	vm = infrav1.VirtualMachine{
		Name:  vmCtx.VSphereVM.Name,
		State: infrav1.VirtualMachineStatePending,
	}

	// If there is an in-flight task associated with this VM then do not
	// reconcile the VM until the task is completed.
	if inFlight, err := reconcileInFlightTask(ctx, vmCtx); err != nil || inFlight {
		return vm, err
	}

	// This deferred function will trigger a reconcile event for the
	// VSphereVM resource once its associated task completes. If
	// there is no task for the VSphereVM resource then no reconcile
	// event is triggered.
	defer reconcileVSphereVMOnTaskCompletion(ctx, vmCtx)

	// Before going further, we need the VM's managed object reference.
	vmRef, err := findVM(ctx, vmCtx)
	if err != nil {
		if !isNotFound(err) {
			return vm, err
		}

		// If the machine was not found by BIOS UUID, it could mean that the machine got deleted from vcenter directly,
		// but sometimes this error is transient, for instance, if the storage was temporarily disconnected but
		// later recovered, the machine will recover from this error.
		if wasNotFoundByBIOSUUID(err) {
			conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.NotFoundByBIOSUUIDReason, clusterv1.ConditionSeverityWarning, err.Error())
			vm.State = infrav1.VirtualMachineStateNotFound
			return vm, err
		}

		// Otherwise, this is a new machine and the VM should be created.
		// NOTE: We are setting this condition only in case it does not exist, so we avoid to get flickering LastConditionTime
		// in case of cloning errors or powering on errors.
		if !conditions.Has(vmCtx.VSphereVM, infrav1.VMProvisionedCondition) {
			conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.CloningReason, clusterv1.ConditionSeverityInfo, "")
		}

		// Get the bootstrap data.
		bootstrapData, format, err := vms.getBootstrapData(ctx, vmCtx)
		if err != nil {
			conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.CloningFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return vm, err
		}

		// Create the VM.
		err = createVM(ctx, vmCtx, bootstrapData, format)
		if err != nil {
			conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.CloningFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return vm, err
		}
		return vm, nil
	}

	//
	// At this point we know the VM exists, so it needs to be updated.
	//

	// Create a new virtualMachineContext to reconcile the VM.
	virtualMachineCtx := &virtualMachineContext{
		VMContext: *vmCtx,
		Obj:       object.NewVirtualMachine(vmCtx.Session.Client.Client, vmRef),
		Ref:       vmRef,
		State:     &vm,
	}
	vm.VMRef = vmRef.String()

	vms.reconcileUUID(ctx, virtualMachineCtx)

	if ok, err := vms.reconcileHardwareVersion(ctx, virtualMachineCtx); err != nil || !ok {
		return vm, err
	}

	if err := vms.reconcilePCIDevices(ctx, virtualMachineCtx); err != nil {
		return vm, err
	}

	if err := vms.reconcileNetworkStatus(ctx, virtualMachineCtx); err != nil {
		return vm, err
	}

	if ok, err := vms.reconcileIPAddresses(ctx, virtualMachineCtx); err != nil || !ok {
		return vm, err
	}

	if ok, err := vms.reconcileMetadata(ctx, virtualMachineCtx); err != nil || !ok {
		return vm, err
	}

	if err := vms.reconcileStoragePolicy(ctx, virtualMachineCtx); err != nil {
		return vm, err
	}

	if ok, err := vms.reconcileVMGroupInfo(ctx, virtualMachineCtx); err != nil || !ok {
		return vm, err
	}

	if err := vms.reconcileClusterModuleMembership(ctx, virtualMachineCtx); err != nil {
		return vm, err
	}

	if ok, err := vms.reconcilePowerState(ctx, virtualMachineCtx); err != nil || !ok {
		return vm, err
	}

	if err := vms.reconcileHostInfo(ctx, virtualMachineCtx); err != nil {
		return vm, err
	}

	if err := vms.reconcileTags(ctx, virtualMachineCtx); err != nil {
		conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.TagsAttachmentFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return vm, err
	}

	vm.State = infrav1.VirtualMachineStateReady
	return vm, nil
}

// DestroyVM powers off and destroys a virtual machine.
func (vms *VMService) DestroyVM(ctx context.Context, vmCtx *capvcontext.VMContext) (reconcile.Result, infrav1.VirtualMachine, error) {
	log := ctrl.LoggerFrom(ctx)

	vm := infrav1.VirtualMachine{
		Name:  vmCtx.VSphereVM.Name,
		State: infrav1.VirtualMachineStatePending,
	}

	// If there is an in-flight task associated with this VM then do not
	// reconcile the VM until the task is completed.
	if inFlight, err := reconcileInFlightTask(ctx, vmCtx); err != nil || inFlight {
		return reconcile.Result{}, vm, err
	}

	// This deferred function will trigger a reconcile event for the
	// VSphereVM resource once its associated task completes. If
	// there is no task for the VSphereVM resource then no reconcile
	// event is triggered.
	defer reconcileVSphereVMOnTaskCompletion(ctx, vmCtx)

	// Before going further, we need the VM's managed object reference.
	vmRef, err := findVM(ctx, vmCtx)
	if err != nil {
		// If the VM's MoRef could not be found then the VM no longer exists. This
		// is the desired state.
		if isNotFound(err) || isFolderNotFound(err) {
			vm.State = infrav1.VirtualMachineStateNotFound
			return reconcile.Result{}, vm, nil
		}
		return reconcile.Result{}, vm, err
	}

	//
	// At this point we know the VM exists, so it needs to be destroyed.
	//

	// Create a new virtualMachineContext to reconcile the VM.
	virtualMachineCtx := &virtualMachineContext{
		VMContext: *vmCtx,
		Obj:       object.NewVirtualMachine(vmCtx.Session.Client.Client, vmRef),
		Ref:       vmRef,
		State:     &vm,
	}

	// Shut down the VM
	powerState, err := vms.getPowerState(ctx, virtualMachineCtx)
	if err != nil {
		return reconcile.Result{}, vm, err
	}

	if powerState == infrav1.VirtualMachinePowerStatePoweredOn {
		// Trigger the soft power off and set the condition.
		softPowerOffPending, err := vms.triggerSoftPowerOff(ctx, virtualMachineCtx)
		if err != nil {
			return reconcile.Result{}, vm, err
		}

		if softPowerOffPending {
			// Return to reconcile later when we have a pending soft power off operation.
			return reconcile.Result{RequeueAfter: time.Second * 20}, vm, nil
		}

		// Hard shut off VM.
		task, err := virtualMachineCtx.Obj.PowerOff(ctx)
		if err != nil {
			return reconcile.Result{}, vm, err
		}

		virtualMachineCtx.VSphereVM.Status.TaskRef = task.Reference().Value
		if err = virtualMachineCtx.Patch(ctx); err != nil {
			return reconcile.Result{}, vm, errors.Wrapf(err, "failed to patch VSphereVM")
		}

		log.Info("Wait for VM to be powered off")
		return reconcile.Result{}, vm, nil
	}

	// Only set the GuestPowerOffCondition to true when the guest shutdown has been initiated.
	if conditions.Has(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededCondition) {
		conditions.MarkTrue(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededCondition)
	}

	log.Info("VM is powered off")
	if vmCtx.ClusterModuleInfo != nil {
		log := log.WithValues("moduleUUID", *vmCtx.ClusterModuleInfo)
		ctx := ctrl.LoggerInto(ctx, log)

		provider := clustermodules.NewProvider(vmCtx.Session.TagManager.Client)
		err := provider.RemoveMoRefFromModule(ctx, *vmCtx.ClusterModuleInfo, virtualMachineCtx.Ref)
		if err != nil && !rest.IsStatusError(err, http.StatusNotFound) {
			return reconcile.Result{}, vm, err
		}
		vmCtx.VSphereVM.Status.ModuleUUID = nil
	}

	// At this point the VM is not powered on and can be destroyed. Store the
	// destroy task's reference and return a requeue error.
	log.Info("Destroying vm")
	task, err := virtualMachineCtx.Obj.Destroy(ctx)
	if err != nil {
		return reconcile.Result{}, vm, err
	}
	vmCtx.VSphereVM.Status.TaskRef = task.Reference().Value
	log.Info("Wait for VM to be destroyed")
	return reconcile.Result{}, vm, nil
}

func (vms *VMService) reconcileNetworkStatus(ctx context.Context, virtualMachineCtx *virtualMachineContext) error {
	netStatus, err := vms.getNetworkStatus(ctx, virtualMachineCtx)
	if err != nil {
		return err
	}
	virtualMachineCtx.State.Network = netStatus
	return nil
}

// reconcileIPAddresses works to check that all the IPAddressClaim objects for the
// VSphereVM object have been bound.
// This function is a no-op if the VSphereVM has no associated IPAddressClaims.
// A discovered IPAddress is expected to contain a valid IP, Prefix and Gateway.
func (vms *VMService) reconcileIPAddresses(ctx context.Context, virtualMachineCtx *virtualMachineContext) (bool, error) {
	ipamState, err := ipam.BuildState(ctx, virtualMachineCtx.VMContext, virtualMachineCtx.State.Network)
	if err != nil && !errors.Is(err, ipam.ErrWaitingForIPAddr) {
		return false, err
	}
	if errors.Is(err, ipam.ErrWaitingForIPAddr) {
		conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.WaitingForIPAddressReason, clusterv1.ConditionSeverityInfo, err.Error())
		return false, nil
	}
	virtualMachineCtx.IPAMState = ipamState
	return true, nil
}

func (vms *VMService) reconcileMetadata(ctx context.Context, virtualMachineCtx *virtualMachineContext) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	existingMetadata, err := vms.getMetadata(ctx, virtualMachineCtx)
	if err != nil {
		return false, err
	}

	newMetadata, err := util.GetMachineMetadata(virtualMachineCtx.VSphereVM.Name, *virtualMachineCtx.VSphereVM, virtualMachineCtx.IPAMState, virtualMachineCtx.State.Network...)
	if err != nil {
		return false, err
	}

	// If the metadata is the same then return early.
	if string(newMetadata) == existingMetadata {
		return true, nil
	}

	log.Info("Updating VM metadata")
	taskRef, err := vms.setMetadata(ctx, virtualMachineCtx, newMetadata)
	if err != nil {
		return false, errors.Wrapf(err, "unable to set metadata on vm %s", ctx)
	}

	virtualMachineCtx.VSphereVM.Status.TaskRef = taskRef
	log.Info("Wait for VM metadata to be updated")
	return false, nil
}

func (vms *VMService) reconcilePowerState(ctx context.Context, virtualMachineCtx *virtualMachineContext) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	powerState, err := vms.getPowerState(ctx, virtualMachineCtx)
	if err != nil {
		return false, err
	}
	switch powerState {
	case infrav1.VirtualMachinePowerStatePoweredOff:
		log.Info("Powering on VM")
		task, err := virtualMachineCtx.Obj.PowerOn(ctx)
		if err != nil {
			conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.PoweringOnFailedReason, clusterv1.ConditionSeverityWarning, err.Error())
			return false, errors.Wrapf(err, "failed to trigger power on op for vm %s", ctx)
		}
		conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.PoweringOnReason, clusterv1.ConditionSeverityInfo, "")

		// Update the VSphereVM.Status.TaskRef to track the power-on task.
		virtualMachineCtx.VSphereVM.Status.TaskRef = task.Reference().Value
		if err = virtualMachineCtx.Patch(ctx); err != nil {
			return false, errors.Wrapf(err, "failed to patch VSphereVM")
		}

		// Once the VM is successfully powered on, a reconcile request should be
		// triggered once the VM reports IP addresses are available.
		reconcileVSphereVMWhenNetworkIsReady(ctx, virtualMachineCtx, task)

		log.Info("Wait for VM to be powered on")
		return false, nil
	case infrav1.VirtualMachinePowerStatePoweredOn:
		log.Info("VM is powered on")
		return true, nil
	default:
		return false, errors.Errorf("unexpected power state %q for vm %s", powerState, ctx)
	}
}

func (vms *VMService) reconcileStoragePolicy(ctx context.Context, virtualMachineCtx *virtualMachineContext) error {
	log := ctrl.LoggerFrom(ctx)

	if virtualMachineCtx.VSphereVM.Spec.StoragePolicyName == "" {
		log.V(5).Info("Storage policy not defined. skipping reconcile storage policy")
		return nil
	}

	// return early if the VM is already powered on
	powerState, err := vms.getPowerState(ctx, virtualMachineCtx)
	if err != nil {
		return err
	}
	if powerState == infrav1.VirtualMachinePowerStatePoweredOn {
		log.Info("VM powered on. Skipping reconcile storage policy")
		return nil
	}

	pbmClient, err := pbm.NewClient(ctx, virtualMachineCtx.Session.Client.Client)
	if err != nil {
		return errors.Wrap(err, "unable to create pbm client")
	}
	storageProfileID, err := pbmClient.ProfileIDByName(ctx, virtualMachineCtx.VSphereVM.Spec.StoragePolicyName)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve storage profile ID")
	}

	var changes []types.BaseVirtualDeviceConfigSpec
	devices, err := virtualMachineCtx.Obj.Device(ctx)
	if err != nil {
		return err
	}

	disksRefs := make([]pbmTypes.PbmServerObjectRef, 0)
	// diskMap is just an auxiliar map so we don't need to iterate over and over disks to get their configs
	// if we realize they are not on the right storage policy
	diskMap := make(map[string]*types.VirtualDisk)

	disks := devices.SelectByType((*types.VirtualDisk)(nil))

	// We iterate over disks and create an array of disks refs, so we just need to make a single call
	// against vCenter, instead of one call per disk
	// the diskMap is an auxiliar way of, besides the disksRefs, we have a "searchable" disk configuration
	// in case we need to reconfigure a disk, to get its config
	for _, d := range disks {
		disk := d.(*types.VirtualDisk)
		// entities associated with storage policy has key in the form <vm-ID>:<disk>
		diskID := fmt.Sprintf("%s:%d", virtualMachineCtx.Obj.Reference().Value, disk.Key)
		diskMap[diskID] = disk

		disksRefs = append(disksRefs, pbmTypes.PbmServerObjectRef{
			ObjectType: string(pbmTypes.PbmObjectTypeVirtualDiskId),
			Key:        diskID,
		})
	}

	diskObjects, err := pbmClient.QueryAssociatedProfiles(ctx, disksRefs)
	if err != nil {
		return errors.Wrap(err, "unable to query disks associated profiles")
	}

	// Ensure storage policy is set correctly for all disks of the VM
	for k := range diskObjects {
		if !isStoragePolicyIDPresent(storageProfileID, diskObjects[k]) {
			log.V(5).Info("Storage policy not found on disk, adding for reconciliation", "disk", diskObjects[k].Object.Key)
			config := &types.VirtualDeviceConfigSpec{
				Operation: types.VirtualDeviceConfigSpecOperationEdit,
				Device:    diskMap[diskObjects[k].Object.Key],
				Profile: []types.BaseVirtualMachineProfileSpec{
					&types.VirtualMachineDefinedProfileSpec{ProfileId: storageProfileID},
				},
			}
			changes = append(changes, config)
		}
	}

	// If there are pending changes for Storage Policies, do it before moving next
	if len(changes) > 0 {
		task, err := virtualMachineCtx.Obj.Reconfigure(ctx, types.VirtualMachineConfigSpec{
			VmProfile: []types.BaseVirtualMachineProfileSpec{
				&types.VirtualMachineDefinedProfileSpec{ProfileId: storageProfileID},
			},
			DeviceChange: changes,
		})
		if err != nil {
			return errors.Wrapf(err, "unable to set storagePolicy on vm %s", ctx)
		}
		virtualMachineCtx.VSphereVM.Status.TaskRef = task.Reference().Value
	}
	return nil
}

func (vms *VMService) reconcileUUID(ctx context.Context, virtualMachineCtx *virtualMachineContext) {
	virtualMachineCtx.State.BiosUUID = virtualMachineCtx.Obj.UUID(ctx)
}

func (vms *VMService) reconcileHardwareVersion(ctx context.Context, virtualMachineCtx *virtualMachineContext) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	if virtualMachineCtx.VSphereVM.Spec.HardwareVersion != "" {
		var virtualMachine mo.VirtualMachine
		if err := virtualMachineCtx.Obj.Properties(ctx, virtualMachineCtx.Obj.Reference(), []string{"config.version"}, &virtualMachine); err != nil {
			return false, errors.Wrapf(err, "error getting guestInfo version information from VM %s", virtualMachineCtx.VSphereVM.Name)
		}
		toUpgrade, err := util.LessThan(virtualMachine.Config.Version, virtualMachineCtx.VSphereVM.Spec.HardwareVersion)
		if err != nil {
			return false, errors.Wrapf(err, "failed to parse hardware version")
		}
		if toUpgrade {
			log.Info("Upgrading hardware version", "fromVersion", virtualMachine.Config.Version, "toVersion", virtualMachineCtx.VSphereVM.Spec.HardwareVersion)
			task, err := virtualMachineCtx.Obj.UpgradeVM(ctx, virtualMachineCtx.VSphereVM.Spec.HardwareVersion)
			if err != nil {
				return false, errors.Wrapf(err, "error trigging upgrade op for machine %s", ctx)
			}
			virtualMachineCtx.VSphereVM.Status.TaskRef = task.Reference().Value
			return false, nil
		}
	}
	return true, nil
}

func (vms *VMService) reconcilePCIDevices(ctx context.Context, virtualMachineCtx *virtualMachineContext) error {
	log := ctrl.LoggerFrom(ctx)

	if expectedPciDevices := virtualMachineCtx.VSphereVM.Spec.VirtualMachineCloneSpec.PciDevices; len(expectedPciDevices) != 0 {
		specsToBeAdded, err := pci.CalculateDevicesToBeAdded(ctx, virtualMachineCtx.Obj, expectedPciDevices)
		if err != nil {
			return err
		}

		if len(specsToBeAdded) == 0 {
			if conditions.Has(virtualMachineCtx.VSphereVM, infrav1.PCIDevicesDetachedCondition) {
				conditions.Delete(virtualMachineCtx.VSphereVM, infrav1.PCIDevicesDetachedCondition)
			}
			log.V(5).Info("No new PCI devices to be added")
			return nil
		}

		powerState, err := virtualMachineCtx.Obj.PowerState(ctx)
		if err != nil {
			return err
		}
		if powerState == types.VirtualMachinePowerStatePoweredOn {
			// This would arise only when the PCI device is manually removed from
			// the VM post creation.
			log.Info("PCI device cannot be attached in powered on state")
			conditions.MarkFalse(virtualMachineCtx.VSphereVM,
				infrav1.PCIDevicesDetachedCondition,
				infrav1.NotFoundReason,
				clusterv1.ConditionSeverityWarning,
				"PCI devices removed after VM was powered on")
			return errors.Errorf("missing PCI devices")
		}
		log.Info("PCI devices to be added", "number", len(specsToBeAdded))
		if err := virtualMachineCtx.Obj.AddDevice(ctx, pci.ConstructDeviceSpecs(specsToBeAdded)...); err != nil {
			return errors.Wrapf(err, "error adding pci devices for %q", ctx)
		}
	}
	return nil
}

func (vms *VMService) getMetadata(ctx context.Context, virtualMachineCtx *virtualMachineContext) (string, error) {
	var (
		obj mo.VirtualMachine

		pc    = property.DefaultCollector(virtualMachineCtx.Session.Client.Client)
		props = []string{"config.extraConfig"}
	)

	if err := pc.RetrieveOne(ctx, virtualMachineCtx.Ref, props, &obj); err != nil {
		return "", errors.Wrapf(err, "unable to fetch props %v for vm %s", props, ctx)
	}
	if obj.Config == nil {
		return "", nil
	}

	var metadataBase64 string
	for _, ec := range obj.Config.ExtraConfig {
		if optVal := ec.GetOptionValue(); optVal != nil && optVal.Key == guestInfoKeyMetadata {
			if v, ok := optVal.Value.(string); ok {
				metadataBase64 = v
			}
		}
	}

	if metadataBase64 == "" {
		return "", nil
	}

	metadataBuf, err := base64.StdEncoding.DecodeString(metadataBase64)
	if err != nil {
		return "", errors.Wrapf(err, "unable to decode metadata for %s", ctx)
	}

	return string(metadataBuf), nil
}

func (vms *VMService) reconcileHostInfo(ctx context.Context, virtualMachineCtx *virtualMachineContext) error {
	host, err := virtualMachineCtx.Obj.HostSystem(ctx)
	if err != nil {
		return err
	}
	name, err := host.ObjectName(ctx)
	if err != nil {
		return err
	}
	virtualMachineCtx.VSphereVM.Status.Host = name
	return nil
}

func (vms *VMService) setMetadata(ctx context.Context, virtualMachineCtx *virtualMachineContext, metadata []byte) (string, error) {
	var extraConfig extra.Config

	extraConfig.SetCloudInitMetadata(metadata)

	task, err := virtualMachineCtx.Obj.Reconfigure(ctx, types.VirtualMachineConfigSpec{
		ExtraConfig: extraConfig,
	})
	if err != nil {
		return "", errors.Wrapf(err, "unable to set metadata on vm %s", ctx)
	}

	return task.Reference().Value, nil
}

func (vms *VMService) getNetworkStatus(ctx context.Context, virtualMachineCtx *virtualMachineContext) ([]infrav1.NetworkStatus, error) {
	log := ctrl.LoggerFrom(ctx)

	allNetStatus, err := govmominet.GetNetworkStatus(ctx, virtualMachineCtx.Session.Client.Client, virtualMachineCtx.Ref)
	if err != nil {
		return nil, err
	}
	log.V(4).Info("Got allNetStatus", "status", allNetStatus)
	apiNetStatus := []infrav1.NetworkStatus{}
	for _, s := range allNetStatus {
		apiNetStatus = append(apiNetStatus, infrav1.NetworkStatus{
			Connected:   s.Connected,
			IPAddrs:     sanitizeIPAddrs(ctx, s.IPAddrs),
			MACAddr:     s.MACAddr,
			NetworkName: s.NetworkName,
		})
	}
	return apiNetStatus, nil
}

// getBootstrapData obtains a machine's bootstrap data from the relevant k8s secret and returns the
// data and its format.
func (vms *VMService) getBootstrapData(ctx context.Context, vmCtx *capvcontext.VMContext) ([]byte, bootstrapv1.Format, error) {
	log := ctrl.LoggerFrom(ctx)

	if vmCtx.VSphereVM.Spec.BootstrapRef == nil {
		log.Info("VM has no bootstrap data")
		return nil, "", nil
	}

	secret := &corev1.Secret{}
	secretKey := apitypes.NamespacedName{
		Namespace: vmCtx.VSphereVM.Spec.BootstrapRef.Namespace,
		Name:      vmCtx.VSphereVM.Spec.BootstrapRef.Name,
	}
	if err := vmCtx.Client.Get(ctx, secretKey, secret); err != nil {
		return nil, "", errors.Wrapf(err, "failed to get bootstrap data secret for %s", ctx)
	}

	format, ok := secret.Data["format"]
	if !ok || len(format) == 0 {
		// Bootstrap data format is missing or empty - assume cloud-config.
		format = []byte(bootstrapv1.CloudConfig)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, bootstrapv1.Format(format), nil
}

func (vms *VMService) reconcileVMGroupInfo(ctx context.Context, virtualMachineCtx *virtualMachineContext) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	if virtualMachineCtx.VSphereFailureDomain == nil || virtualMachineCtx.VSphereFailureDomain.Spec.Topology.Hosts == nil {
		log.V(5).Info("Hosts topology in failure domain not defined. skipping reconcile VM group")
		return true, nil
	}

	topology := virtualMachineCtx.VSphereFailureDomain.Spec.Topology
	vmGroup, err := cluster.FindVMGroup(ctx, virtualMachineCtx, *topology.ComputeCluster, topology.Hosts.VMGroupName)
	if err != nil {
		return false, errors.Wrapf(err, "unable to find VM Group %s", topology.Hosts.VMGroupName)
	}

	hasVM, err := vmGroup.HasVM(virtualMachineCtx.Ref)
	if err != nil {
		return false, errors.Wrapf(err, "unable to find VM Group %s membership", topology.Hosts.VMGroupName)
	}

	if !hasVM {
		task, err := vmGroup.Add(ctx, virtualMachineCtx.Ref)
		if err != nil {
			return false, errors.Wrapf(err, "failed to add VM %s to VM group", virtualMachineCtx.VSphereVM.Name)
		}
		virtualMachineCtx.VSphereVM.Status.TaskRef = task.Reference().Value
		log.Info("Wait for VM to be added to group")
		return false, nil
	}
	return true, nil
}

func (vms *VMService) reconcileTags(ctx context.Context, virtualMachineCtx *virtualMachineContext) error {
	log := ctrl.LoggerFrom(ctx)

	if len(virtualMachineCtx.VSphereVM.Spec.TagIDs) == 0 {
		log.V(5).Info("No tags defined. skipping tags reconciliation")
		return nil
	}

	err := virtualMachineCtx.Session.TagManager.AttachMultipleTagsToObject(ctx, virtualMachineCtx.VSphereVM.Spec.TagIDs, virtualMachineCtx.Ref)
	if err != nil {
		return errors.Wrapf(err, "failed to attach tags %v to VM %s", virtualMachineCtx.VSphereVM.Spec.TagIDs, virtualMachineCtx.VSphereVM.Name)
	}

	return nil
}

func (vms *VMService) reconcileClusterModuleMembership(ctx context.Context, virtualMachineCtx *virtualMachineContext) error {
	log := ctrl.LoggerFrom(ctx)

	if virtualMachineCtx.ClusterModuleInfo != nil {
		log := log.WithValues("moduleUUID", *virtualMachineCtx.ClusterModuleInfo)
		ctx := ctrl.LoggerInto(ctx, log)

		provider := clustermodules.NewProvider(virtualMachineCtx.Session.TagManager.Client)

		if err := provider.AddMoRefToModule(ctx, *virtualMachineCtx.ClusterModuleInfo, virtualMachineCtx.Ref); err != nil {
			return err
		}
		virtualMachineCtx.VSphereVM.Status.ModuleUUID = virtualMachineCtx.ClusterModuleInfo
	}
	return nil
}
