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

// Package vcenter has tools for cloning virtual machines in vcenter.
package vcenter

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/pbm"
	pbmTypes "github.com/vmware/govmomi/pbm/types"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"k8s.io/utils/ptr"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/extra"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/template"
)

const (
	fullCloneDiskMoveType = types.VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndConsolidate
	linkCloneDiskMoveType = types.VirtualMachineRelocateDiskMoveOptionsCreateNewChildDiskBacking

	// maxUnitNumber constant is used to define the maximum number of devices that can be assigned to a virtual machine's controller.
	// Not all controllers support up to 30, but the maximum is 30.
	// xref: https://docs.vmware.com/en/VMware-vSphere/8.0/vsphere-vm-administration/GUID-5872D173-A076-42FE-8D0B-9DB0EB0E7362.html#:~:text=If%20you%20add%20a%20hard,values%20from%200%20to%2014.
	maxUnitNumber = 30
)

// Clone kicks off a clone operation on vCenter to create a new virtual machine. This function does not wait for
// the virtual machine to be created on the vCenter, which can be resolved by waiting on the task reference stored
// in VMContext.VSphereVM.Status.TaskRef.
func Clone(ctx context.Context, vmCtx *capvcontext.VMContext, bootstrapData []byte, format bootstrapv1.Format) error {
	log := ctrl.LoggerFrom(ctx)

	vmCtx = &capvcontext.VMContext{
		ControllerManagerContext: vmCtx.ControllerManagerContext,
		VSphereVM:                vmCtx.VSphereVM,
		Session:                  vmCtx.Session,
		PatchHelper:              vmCtx.PatchHelper,
	}
	log.Info("Starting clone process")

	var extraConfig extra.Config
	if len(bootstrapData) > 0 {
		log.Info("Applied bootstrap data to VM clone spec")
		switch format {
		case bootstrapv1.CloudConfig:
			extraConfig.SetCloudInitUserData(bootstrapData)
		case bootstrapv1.Ignition:
			extraConfig.SetIgnitionUserData(bootstrapData)
		}
	}
	if vmCtx.VSphereVM.Spec.CustomVMXKeys != nil {
		log.Info("Applied custom VMX keys to VM clone spec")
		if err := extraConfig.SetCustomVMXKeys(vmCtx.VSphereVM.Spec.CustomVMXKeys); err != nil {
			return err
		}
	}
	tpl, err := template.FindTemplate(ctx, vmCtx.GetSession(), vmCtx.VSphereVM.Spec.Template)
	if err != nil {
		return err
	}

	// If a linked clone is requested then a MoRef for a snapshot must be
	// found with which to perform the linked clone.
	var snapshotRef *types.ManagedObjectReference

	if vmCtx.VSphereVM.Spec.CloneMode == "" || vmCtx.VSphereVM.Spec.CloneMode == infrav1.LinkedClone {
		log.Info("Linked clone requested")
		// If the name of a snapshot was not provided then find the template's
		// current snapshot.
		if snapshotName := vmCtx.VSphereVM.Spec.Snapshot; snapshotName == "" {
			log.Info("Searching for current snapshot")
			var vm mo.VirtualMachine
			if err := tpl.Properties(ctx, tpl.Reference(), []string{"snapshot"}, &vm); err != nil {
				return errors.Wrapf(err, "error getting snapshot information for template %s", vmCtx.VSphereVM.Spec.Template)
			}
			if vm.Snapshot != nil {
				snapshotRef = vm.Snapshot.CurrentSnapshot
			}
		} else {
			log.Info("Searching for snapshot by name", "snapshotName", snapshotName)
			var err error
			snapshotRef, err = tpl.FindSnapshot(ctx, snapshotName)
			if err != nil {
				log.Info("Failed to find snapshot", "snapshotName", snapshotName)
			}
		}
	}

	// The type of clone operation depends on whether there is a snapshot
	// from which to do a linked clone.
	diskMoveType := fullCloneDiskMoveType
	vmCtx.VSphereVM.Status.CloneMode = infrav1.FullClone
	if snapshotRef != nil {
		// Record the actual type of clone mode used as well as the name of
		// the snapshot (if not the current snapshot).
		vmCtx.VSphereVM.Status.CloneMode = infrav1.LinkedClone
		vmCtx.VSphereVM.Status.Snapshot = snapshotRef.Value
		diskMoveType = linkCloneDiskMoveType
	}

	folder, err := vmCtx.Session.Finder.FolderOrDefault(ctx, vmCtx.VSphereVM.Spec.Folder)
	if err != nil {
		return errors.Wrapf(err, "unable to get folder for %q", vmCtx)
	}

	pool, err := vmCtx.Session.Finder.ResourcePoolOrDefault(ctx, vmCtx.VSphereVM.Spec.ResourcePool)
	if err != nil {
		return errors.Wrapf(err, "unable to get resource pool for %q", vmCtx)
	}

	devices, err := tpl.Device(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting devices for %q", vmCtx)
	}

	// Create a new list of device specs for cloning the VM.
	var deviceSpecs []types.BaseVirtualDeviceConfigSpec

	// Only non-linked clones may expand the size of the template's disk.
	if snapshotRef == nil {
		diskSpecs, err := getDiskSpec(vmCtx, devices)
		if err != nil {
			return errors.Wrapf(err, "error getting disk spec for %q", vmCtx)
		}
		deviceSpecs = append(deviceSpecs, diskSpecs...)
	}

	// Process all DataDisks definitions to dynamically create and add disks to the VM
	if len(vmCtx.VSphereVM.Spec.DataDisks) > 0 {
		dataDisks, err := createDataDisks(ctx, vmCtx.VSphereVM.Spec.DataDisks, devices)
		if err != nil {
			return errors.Wrapf(err, "error getting data disks")
		}
		log.V(4).Info("Adding the following data disks", "disks", dataDisks)
		deviceSpecs = append(deviceSpecs, dataDisks...)
	}

	networkSpecs, err := getNetworkSpecs(ctx, vmCtx, devices)
	if err != nil {
		return errors.Wrapf(err, "error getting network specs for %q", vmCtx)
	}

	deviceSpecs = append(deviceSpecs, networkSpecs...)

	numCPUs := vmCtx.VSphereVM.Spec.NumCPUs
	if numCPUs < 2 {
		numCPUs = 2
	}
	numCoresPerSocket := vmCtx.VSphereVM.Spec.NumCoresPerSocket
	if numCoresPerSocket == 0 {
		numCoresPerSocket = numCPUs
	}
	memMiB := vmCtx.VSphereVM.Spec.MemoryMiB
	if memMiB == 0 {
		memMiB = 2048
	}

	// Disable the vAppConfig during VM creation to ensure Cloud-Init inside of the guest does not
	// activate and prefer the OVF datasource over the VMware datasource.
	vappConfigRemoved := true

	spec := types.VirtualMachineCloneSpec{
		Config: &types.VirtualMachineConfigSpec{
			// Assign the clone's InstanceUUID the value of the Kubernetes Machine
			// object's UID. This allows lookup of the cloned VM prior to knowing
			// the VM's UUID.
			InstanceUuid:      string(vmCtx.VSphereVM.UID),
			Flags:             newVMFlagInfo(),
			DeviceChange:      deviceSpecs,
			ExtraConfig:       extraConfig,
			NumCPUs:           numCPUs,
			NumCoresPerSocket: numCoresPerSocket,
			MemoryMB:          memMiB,
			VAppConfigRemoved: &vappConfigRemoved,
		},
		Location: types.VirtualMachineRelocateSpec{
			DiskMoveType: string(diskMoveType),
			Folder:       types.NewReference(folder.Reference()),
			Pool:         types.NewReference(pool.Reference()),
		},
		// This is implicit, but making it explicit as it is important to not
		// power the VM on before its virtual hardware is created and the MAC
		// address(es) used to build and inject the VM with cloud-init metadata
		// are generated.
		PowerOn:  false,
		Snapshot: snapshotRef,
	}

	// For PCI devices, the memory for the VM needs to be reserved
	// We can replace this once we have another way of reserving memory option
	// exposed via the API types.
	if len(vmCtx.VSphereVM.Spec.PciDevices) > 0 {
		spec.Config.MemoryReservationLockedToMax = ptr.To(true)
	}

	var datastoreRef *types.ManagedObjectReference
	if vmCtx.VSphereVM.Spec.Datastore != "" {
		datastore, err := vmCtx.Session.Finder.Datastore(ctx, vmCtx.VSphereVM.Spec.Datastore)
		if err != nil {
			return errors.Wrapf(err, "unable to get datastore %s for %q", vmCtx.VSphereVM.Spec.Datastore, vmCtx)
		}
		datastoreRef = types.NewReference(datastore.Reference())
		spec.Location.Datastore = datastoreRef
	}

	var storageProfileID string
	if vmCtx.VSphereVM.Spec.StoragePolicyName != "" {
		pbmClient, err := pbm.NewClient(ctx, vmCtx.Session.Client.Client)
		if err != nil {
			return errors.Wrapf(err, "unable to create pbm client for %q", vmCtx)
		}

		storageProfileID, err = pbmClient.ProfileIDByName(ctx, vmCtx.VSphereVM.Spec.StoragePolicyName)
		if err != nil {
			return errors.Wrapf(err, "unable to get storageProfileID from name %s for %q", vmCtx.VSphereVM.Spec.StoragePolicyName, vmCtx)
		}

		var hubs []pbmTypes.PbmPlacementHub

		// If there's a Datastore configured, it should be the only one for which we check if it matches the requirements of the Storage Policy
		if datastoreRef != nil {
			hubs = append(hubs, pbmTypes.PbmPlacementHub{
				HubType: datastoreRef.Type,
				HubId:   datastoreRef.Value,
			})
		} else {
			// Otherwise we should get just the Datastores connected to our pool
			cluster, err := pool.Owner(ctx)
			if err != nil {
				return errors.Wrapf(err, "failed to get owning cluster of resourcepool %q to calculate datastore based on storage policy", pool)
			}

			dsList, err := object.NewComputeResource(vmCtx.Session.Client.Client, cluster.Reference()).Datastores(ctx)
			if err != nil {
				return errors.Wrapf(err, "unable to list datastores from owning cluster of requested resourcepool")
			}

			var refs []types.ManagedObjectReference
			for i := range dsList {
				refs = append(refs, dsList[i].Reference())
			}

			var datastores []mo.Datastore
			if err := property.DefaultCollector(vmCtx.Session.Client.Client).Retrieve(ctx, refs, []string{"summary"}, &datastores); err != nil {
				return errors.Wrapf(err, "unable to collect datastore properties to validate maintenance mode")
			}

			for _, ds := range datastores {
				if ds.Summary.MaintenanceMode != string(types.DatastoreSummaryMaintenanceModeStateNormal) {
					log.V(4).Info("datastore is in maintenance mode, skipping", "datastore", ds.Summary.Name)
					continue
				}

				hubs = append(hubs, pbmTypes.PbmPlacementHub{
					HubType: ds.Reference().Type,
					HubId:   ds.Reference().Value,
				})
			}
		}

		var constraints []pbmTypes.BasePbmPlacementRequirement
		constraints = append(constraints, &pbmTypes.PbmPlacementCapabilityProfileRequirement{ProfileId: pbmTypes.PbmProfileId{UniqueId: storageProfileID}})
		result, err := pbmClient.CheckRequirements(ctx, hubs, nil, constraints)
		if err != nil {
			return errors.Wrapf(err, "unable to check requirements for storage policy")
		}

		if len(result.CompatibleDatastores()) == 0 {
			return fmt.Errorf("no compatible datastores found for storage policy: %s", vmCtx.VSphereVM.Spec.StoragePolicyName)
		}

		// If datastoreRef is nil here it means that the user didn't specify a Datastore. So we should
		// select one of the datastores of the owning cluster of the resource pool that matched the
		// requirements of the storage policy.
		if datastoreRef == nil {
			r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // We won't need cryptographically secure randomness here.
			ds := result.CompatibleDatastores()[r.Intn(len(result.CompatibleDatastores()))]
			datastoreRef = &types.ManagedObjectReference{Type: ds.HubType, Value: ds.HubId}
		}
	}

	// if datastoreRef is nil here, means that user didn't specified a datastore NOR a
	// storagepolicy, so we should select the default
	if datastoreRef == nil {
		// if no datastore defined through VM spec or storage policy, use default
		datastore, err := vmCtx.Session.Finder.DefaultDatastore(ctx)
		if err != nil {
			return errors.Wrapf(err, "unable to get default datastore for %q", vmCtx)
		}
		datastoreRef = types.NewReference(datastore.Reference())
	}

	disks := devices.SelectByType((*types.VirtualDisk)(nil))
	isLinkedClone := snapshotRef != nil
	spec.Location.Disk = getDiskLocators(disks, *datastoreRef, isLinkedClone)
	spec.Location.Datastore = datastoreRef

	log.Info(fmt.Sprintf("Cloning Machine with clone mode %s", vmCtx.VSphereVM.Status.CloneMode))
	task, err := tpl.Clone(ctx, folder, vmCtx.VSphereVM.Name, spec)
	if err != nil {
		return errors.Wrapf(err, "error trigging clone op for machine %s", vmCtx)
	}

	vmCtx.VSphereVM.Status.TaskRef = task.Reference().Value

	// patch the vsphereVM early to ensure that the task is
	// reflected in the status right away, this avoids situations
	// of concurrent clones
	if err := vmCtx.Patch(ctx); err != nil {
		log.Error(err, "Failed to patch VSphereVM (best-effort)")
	}
	return nil
}

func newVMFlagInfo() *types.VirtualMachineFlagInfo {
	diskUUIDEnabled := true
	return &types.VirtualMachineFlagInfo{
		DiskUuidEnabled: &diskUUIDEnabled,
	}
}

func getDiskLocators(disks object.VirtualDeviceList, datastoreRef types.ManagedObjectReference, isLinkedClone bool) []types.VirtualMachineRelocateSpecDiskLocator {
	diskLocators := make([]types.VirtualMachineRelocateSpecDiskLocator, 0, len(disks))
	for _, disk := range disks {
		dl := types.VirtualMachineRelocateSpecDiskLocator{
			DiskId:       disk.GetVirtualDevice().Key,
			DiskMoveType: string(types.VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndDisallowSharing),
			Datastore:    datastoreRef,
		}

		if isLinkedClone {
			dl.DiskMoveType = string(linkCloneDiskMoveType)
		}
		if vmDiskBacking, ok := disk.(*types.VirtualDisk).Backing.(*types.VirtualDiskFlatVer2BackingInfo); ok {
			dl.DiskBackingInfo = vmDiskBacking
		}
		diskLocators = append(diskLocators, dl)
	}

	return diskLocators
}

func getDiskSpec(vmCtx *capvcontext.VMContext, devices object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	disks := devices.SelectByType((*types.VirtualDisk)(nil))
	if len(disks) == 0 {
		return nil, errors.Errorf("Invalid disk count: %d", len(disks))
	}

	// There is at least one disk
	var diskSpecs []types.BaseVirtualDeviceConfigSpec
	primaryDisk := disks[0].(*types.VirtualDisk)
	primaryCloneCapacityKB := int64(vmCtx.VSphereVM.Spec.DiskGiB) * 1024 * 1024
	primaryDiskConfigSpec, err := getDiskConfigSpec(primaryDisk, primaryCloneCapacityKB)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting disk config spec for primary disk")
	}
	diskSpecs = append(diskSpecs, primaryDiskConfigSpec)

	// Check for additional disks
	// CAPV will not spin up additional extra disks provided in the conf but not available in the template
	if len(disks) > 1 {
		// Disk range starts from 1 to avoid primary disk
		for i, disk := range disks[1:] {
			var diskCloneCapacityKB int64
			// Check if additional Disks have been provided
			if len(vmCtx.VSphereVM.Spec.AdditionalDisksGiB) > i {
				diskCloneCapacityKB = int64(vmCtx.VSphereVM.Spec.AdditionalDisksGiB[i]) * 1024 * 1024
			} else {
				diskCloneCapacityKB = disk.(*types.VirtualDisk).CapacityInKB
			}
			additionalDiskConfigSpec, err := getDiskConfigSpec(disk.(*types.VirtualDisk), diskCloneCapacityKB)
			if err != nil {
				return nil, errors.Wrap(err, "Error getting disk config spec for additional disk")
			}
			diskSpecs = append(diskSpecs, additionalDiskConfigSpec)
		}
	}
	return diskSpecs, nil
}

func getDiskConfigSpec(disk *types.VirtualDisk, diskCloneCapacityKB int64) (types.BaseVirtualDeviceConfigSpec, error) {
	switch {
	case diskCloneCapacityKB == 0:
		// No disk size specified for the clone. Default to the template disk capacity.
	case diskCloneCapacityKB > 0 && diskCloneCapacityKB >= disk.CapacityInKB:
		disk.CapacityInKB = diskCloneCapacityKB
	case diskCloneCapacityKB > 0 && diskCloneCapacityKB < disk.CapacityInKB:
		return nil, errors.Errorf(
			"can't resize template disk down, initial capacity is larger: %dKiB > %dKiB",
			disk.CapacityInKB, diskCloneCapacityKB)
	}

	return &types.VirtualDeviceConfigSpec{
		Operation: types.VirtualDeviceConfigSpecOperationEdit,
		Device:    disk,
	}, nil
}

// createDataDisks parses through the list of VSphereDisk objects and generates the VirtualDeviceConfigSpec for each one.
func createDataDisks(ctx context.Context, dataDiskDefs []infrav1.VSphereDisk, devices object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log := ctrl.LoggerFrom(ctx)
	additionalDisks := []types.BaseVirtualDeviceConfigSpec{}

	disks := devices.SelectByType((*types.VirtualDisk)(nil))
	if len(disks) == 0 {
		return nil, errors.Errorf("Invalid disk count: %d", len(disks))
	}

	// There is at least one disk
	primaryDisk := disks[0].(*types.VirtualDisk)

	// Get the controller of the primary disk.
	controller, ok := devices.FindByKey(primaryDisk.ControllerKey).(types.BaseVirtualController)
	if !ok {
		return nil, errors.Errorf("unable to find controller with key=%v", primaryDisk.ControllerKey)
	}

	controllerKey := controller.GetVirtualController().Key
	unitNumberAssigner, err := newUnitNumberAssigner(controller, devices)
	if err != nil {
		return nil, err
	}

	for i, dataDisk := range dataDiskDefs {
		log.V(2).Info("Adding disk", "name", dataDisk.Name, "spec", dataDisk)

		backing := &types.VirtualDiskFlatVer2BackingInfo{
			DiskMode: string(types.VirtualDiskModePersistent),
			VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
				FileName: "",
			},
		}

		// Set provisioning type for the new data disk.
		// Currently, if ThinProvisioned is not set, GOVC will set default to false.  We may want to change this behavior
		// to match what template image OS disk has configured to make them match if not set.
		switch dataDisk.ProvisioningMode {
		case infrav1.ThinProvisioningMode:
			backing.ThinProvisioned = types.NewBool(true)
		case infrav1.ThickProvisioningMode:
			backing.ThinProvisioned = types.NewBool(false)
		case infrav1.EagerlyZeroedProvisioningMode:
			backing.ThinProvisioned = types.NewBool(false)
			backing.EagerlyScrub = types.NewBool(true)
		default:
			log.V(2).Info("No provisioning type detected.  Leaving configuration empty.")
		}

		dev := &types.VirtualDisk{
			VirtualDevice: types.VirtualDevice{
				// Key needs to be unique and cannot match another new disk being added.  So we'll use the index as an
				// input to NewKey.  NewKey() will always return same value since our new devices are not part of devices yet.
				Key:           devices.NewKey() - int32(i),
				Backing:       backing,
				ControllerKey: controller.GetVirtualController().Key,
			},
			CapacityInKB: int64(dataDisk.SizeGiB) * 1024 * 1024,
		}

		vd := dev.GetVirtualDevice()
		vd.ControllerKey = controllerKey

		// Assign unit number to the new disk.  Should be next available slot on the controller.
		unitNumber, err := unitNumberAssigner.assign()
		if err != nil {
			return nil, err
		}
		vd.UnitNumber = &unitNumber

		log.V(4).Info("Created device for data disk device", "name", dataDisk.Name, "spec", dataDisk, "device", dev)

		additionalDisks = append(additionalDisks, &types.VirtualDeviceConfigSpec{
			Device:        dev,
			Operation:     types.VirtualDeviceConfigSpecOperationAdd,
			FileOperation: types.VirtualDeviceConfigSpecFileOperationCreate,
		})
	}

	return additionalDisks, nil
}

type unitNumberAssigner struct {
	used   []bool
	offset int32
}

func newUnitNumberAssigner(controller types.BaseVirtualController, existingDevices object.VirtualDeviceList) (*unitNumberAssigner, error) {
	if controller == nil {
		return nil, errors.New("controller parameter cannot be nil")
	}
	used := make([]bool, maxUnitNumber)

	// SCSIControllers also use a unit.
	if scsiController, ok := controller.(types.BaseVirtualSCSIController); ok {
		used[scsiController.GetVirtualSCSIController().ScsiCtlrUnitNumber] = true
	}
	controllerKey := controller.GetVirtualController().Key

	// Mark all unit numbers of existing devices as used
	for _, device := range existingDevices {
		d := device.GetVirtualDevice()
		if d.ControllerKey == controllerKey && d.UnitNumber != nil {
			used[*d.UnitNumber] = true
		}
	}

	// Set offset to 0, it will auto-increment on the first assignment.
	return &unitNumberAssigner{used: used, offset: 0}, nil
}

func (a *unitNumberAssigner) assign() (int32, error) {
	if int(a.offset) > len(a.used) {
		return -1, fmt.Errorf("all unit numbers are already in-use")
	}
	for i, isInUse := range a.used[a.offset:] {
		unit := int32(i) + a.offset
		if !isInUse {
			a.used[unit] = true
			a.offset++
			return unit, nil
		}
	}
	return -1, fmt.Errorf("all unit numbers are already in-use")
}

const ethCardType = "vmxnet3"

func getNetworkSpecs(ctx context.Context, vmCtx *capvcontext.VMContext, devices object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	log := ctrl.LoggerFrom(ctx)

	deviceSpecs := []types.BaseVirtualDeviceConfigSpec{}

	// Remove any existing NICs
	for _, dev := range devices.SelectByType((*types.VirtualEthernetCard)(nil)) {
		deviceSpecs = append(deviceSpecs, &types.VirtualDeviceConfigSpec{
			Device:    dev,
			Operation: types.VirtualDeviceConfigSpecOperationRemove,
		})
	}

	// Add new NICs based on the machine config.
	key := int32(-100)
	for i := range vmCtx.VSphereVM.Spec.Network.Devices {
		netSpec := &vmCtx.VSphereVM.Spec.Network.Devices[i]
		ref, err := vmCtx.Session.Finder.Network(ctx, netSpec.NetworkName)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to find network %q", netSpec.NetworkName)
		}
		backing, err := ref.EthernetCardBackingInfo(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to create new ethernet card backing info for network %q on %q", netSpec.NetworkName, vmCtx)
		}
		dev, err := object.EthernetCardTypes().CreateEthernetCard(ethCardType, backing)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to create new ethernet card %q for network %q on %q", ethCardType, netSpec.NetworkName, vmCtx)
		}

		// Get the actual NIC object. This is safe to assert without a check
		// because "object.EthernetCardTypes().CreateEthernetCard" returns a
		// "types.BaseVirtualEthernetCard" as a "types.BaseVirtualDevice".
		nic := dev.(types.BaseVirtualEthernetCard).GetVirtualEthernetCard()

		if netSpec.MACAddr != "" {
			nic.MacAddress = netSpec.MACAddr
			// Please see https://www.vmware.com/support/developer/converter-sdk/conv60_apireference/vim.vm.device.VirtualEthernetCard.html#addressType
			// for the valid values for this field.
			nic.AddressType = string(types.VirtualEthernetCardMacTypeManual)
			log.V(4).Info("Configured manual MAC address", "macAddress", nic.MacAddress)
		}

		// Assign a temporary device key to ensure that a unique one will be
		// generated when the device is created.
		nic.Key = key

		deviceSpecs = append(deviceSpecs, &types.VirtualDeviceConfigSpec{
			Device:    dev,
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
		})
		log.V(4).Info("Created network device", "ethCardType", ethCardType, "networkSpec", netSpec)
		key--
	}

	return deviceSpecs, nil
}
