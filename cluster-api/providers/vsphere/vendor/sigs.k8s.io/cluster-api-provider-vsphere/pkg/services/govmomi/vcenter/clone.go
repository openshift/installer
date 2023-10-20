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

package vcenter

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/pbm"
	pbmTypes "github.com/vmware/govmomi/pbm/types"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"k8s.io/utils/pointer"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/extra"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/template"
)

const (
	fullCloneDiskMoveType = types.VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndConsolidate
	linkCloneDiskMoveType = types.VirtualMachineRelocateDiskMoveOptionsCreateNewChildDiskBacking
)

// Clone kicks off a clone operation on vCenter to create a new virtual machine. This function does not wait for
// the virtual machine to be created on the vCenter, which can be resolved by waiting on the task reference stored
// in VMContext.VSphereVM.Status.TaskRef.
//
//nolint:gocognit,gocyclo
func Clone(ctx *context.VMContext, bootstrapData []byte, format bootstrapv1.Format) error {
	ctx = &context.VMContext{
		ControllerContext: ctx.ControllerContext,
		VSphereVM:         ctx.VSphereVM,
		Session:           ctx.Session,
		Logger:            ctx.Logger.WithName("vcenter"),
		PatchHelper:       ctx.PatchHelper,
	}
	ctx.Logger.Info("starting clone process")

	var extraConfig extra.Config
	if len(bootstrapData) > 0 {
		ctx.Logger.Info("applied bootstrap data to VM clone spec")
		switch format {
		case bootstrapv1.CloudConfig:
			extraConfig.SetCloudInitUserData(bootstrapData)
		case bootstrapv1.Ignition:
			extraConfig.SetIgnitionUserData(bootstrapData)
		}
	}
	if ctx.VSphereVM.Spec.CustomVMXKeys != nil {
		ctx.Logger.Info("applied custom vmx keys o VM clone spec")
		if err := extraConfig.SetCustomVMXKeys(ctx.VSphereVM.Spec.CustomVMXKeys); err != nil {
			return err
		}
	}
	tpl, err := template.FindTemplate(ctx, ctx.VSphereVM.Spec.Template)
	if err != nil {
		return err
	}

	// If a linked clone is requested then a MoRef for a snapshot must be
	// found with which to perform the linked clone.
	var snapshotRef *types.ManagedObjectReference
	//nolint:nestif
	if ctx.VSphereVM.Spec.CloneMode == "" || ctx.VSphereVM.Spec.CloneMode == infrav1.LinkedClone {
		ctx.Logger.Info("linked clone requested")
		// If the name of a snapshot was not provided then find the template's
		// current snapshot.
		if snapshotName := ctx.VSphereVM.Spec.Snapshot; snapshotName == "" {
			ctx.Logger.Info("searching for current snapshot")
			var vm mo.VirtualMachine
			if err := tpl.Properties(ctx, tpl.Reference(), []string{"snapshot"}, &vm); err != nil {
				return errors.Wrapf(err, "error getting snapshot information for template %s", ctx.VSphereVM.Spec.Template)
			}
			if vm.Snapshot != nil {
				snapshotRef = vm.Snapshot.CurrentSnapshot
			}
		} else {
			ctx.Logger.Info("searching for snapshot by name", "snapshotName", snapshotName)
			var err error
			snapshotRef, err = tpl.FindSnapshot(ctx, snapshotName)
			if err != nil {
				ctx.Logger.Info("failed to find snapshot", "snapshotName", snapshotName)
			}
		}
	}

	// The type of clone operation depends on whether or not there is a snapshot
	// from which to do a linked clone.
	diskMoveType := fullCloneDiskMoveType
	ctx.VSphereVM.Status.CloneMode = infrav1.FullClone
	if snapshotRef != nil {
		// Record the actual type of clone mode used as well as the name of
		// the snapshot (if not the current snapshot).
		ctx.VSphereVM.Status.CloneMode = infrav1.LinkedClone
		ctx.VSphereVM.Status.Snapshot = snapshotRef.Value
		diskMoveType = linkCloneDiskMoveType
	}

	folder, err := ctx.Session.Finder.FolderOrDefault(ctx, ctx.VSphereVM.Spec.Folder)
	if err != nil {
		return errors.Wrapf(err, "unable to get folder for %q", ctx)
	}

	pool, err := ctx.Session.Finder.ResourcePoolOrDefault(ctx, ctx.VSphereVM.Spec.ResourcePool)
	if err != nil {
		return errors.Wrapf(err, "unable to get resource pool for %q", ctx)
	}

	devices, err := tpl.Device(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting devices for %q", ctx)
	}

	// Create a new list of device specs for cloning the VM.
	var deviceSpecs []types.BaseVirtualDeviceConfigSpec

	// Only non-linked clones may expand the size of the template's disk.
	if snapshotRef == nil {
		diskSpecs, err := getDiskSpec(ctx, devices)
		if err != nil {
			return errors.Wrapf(err, "error getting disk spec for %q", ctx)
		}
		deviceSpecs = append(deviceSpecs, diskSpecs...)
	}

	networkSpecs, err := getNetworkSpecs(ctx, devices)
	if err != nil {
		return errors.Wrapf(err, "error getting network specs for %q", ctx)
	}

	deviceSpecs = append(deviceSpecs, networkSpecs...)

	if err != nil {
		return errors.Wrapf(err, "error getting network specs for %q", ctx)
	}

	numCPUs := ctx.VSphereVM.Spec.NumCPUs
	if numCPUs < 2 {
		numCPUs = 2
	}
	numCoresPerSocket := ctx.VSphereVM.Spec.NumCoresPerSocket
	if numCoresPerSocket == 0 {
		numCoresPerSocket = numCPUs
	}
	memMiB := ctx.VSphereVM.Spec.MemoryMiB
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
			InstanceUuid:      string(ctx.VSphereVM.UID),
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
	if len(ctx.VSphereVM.Spec.PciDevices) > 0 {
		spec.Config.MemoryReservationLockedToMax = pointer.Bool(true)
	}

	var datastoreRef *types.ManagedObjectReference
	if ctx.VSphereVM.Spec.Datastore != "" {
		datastore, err := ctx.Session.Finder.Datastore(ctx, ctx.VSphereVM.Spec.Datastore)
		if err != nil {
			return errors.Wrapf(err, "unable to get datastore %s for %q", ctx.VSphereVM.Spec.Datastore, ctx)
		}
		datastoreRef = types.NewReference(datastore.Reference())
		spec.Location.Datastore = datastoreRef
	}

	var storageProfileID string
	//nolint:nestif
	if ctx.VSphereVM.Spec.StoragePolicyName != "" {
		pbmClient, err := pbm.NewClient(ctx, ctx.Session.Client.Client)
		if err != nil {
			return errors.Wrapf(err, "unable to create pbm client for %q", ctx)
		}

		storageProfileID, err = pbmClient.ProfileIDByName(ctx, ctx.VSphereVM.Spec.StoragePolicyName)
		if err != nil {
			return errors.Wrapf(err, "unable to get storageProfileID from name %s for %q", ctx.VSphereVM.Spec.StoragePolicyName, ctx)
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
			dsGetter := object.NewComputeResource(ctx.Session.Client.Client, cluster.Reference())
			datastores, err := dsGetter.Datastores(ctx)
			if err != nil {
				return errors.Wrapf(err, "unable to list datastores from owning cluster of requested resourcepool")
			}
			for _, ds := range datastores {
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
			return fmt.Errorf("no compatible datastores found for storage policy: %s", ctx.VSphereVM.Spec.StoragePolicyName)
		}

		// If datastoreRef is nil here it means that the user didn't specify a Datastore. So we should
		// select one of the datastores of the owning cluster of the resource pool that matched the
		// requirements of the storage policy.
		if datastoreRef == nil {
			r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
			ds := result.CompatibleDatastores()[r.Intn(len(result.CompatibleDatastores()))]
			datastoreRef = &types.ManagedObjectReference{Type: ds.HubType, Value: ds.HubId}
		}
	}

	// if datastoreRef is nil here, means that user didn't specified a datastore NOR a
	// storagepolicy, so we should select the default
	if datastoreRef == nil {
		// if no datastore defined through VM spec or storage policy, use default
		datastore, err := ctx.Session.Finder.DefaultDatastore(ctx)
		if err != nil {
			return errors.Wrapf(err, "unable to get default datastore for %q", ctx)
		}
		datastoreRef = types.NewReference(datastore.Reference())
	}

	disks := devices.SelectByType((*types.VirtualDisk)(nil))
	isLinkedClone := snapshotRef != nil
	spec.Location.Disk = getDiskLocators(disks, *datastoreRef, isLinkedClone)
	spec.Location.Datastore = datastoreRef

	ctx.Logger.Info("cloning machine", "namespace", ctx.VSphereVM.Namespace, "name", ctx.VSphereVM.Name, "cloneType", ctx.VSphereVM.Status.CloneMode)
	task, err := tpl.Clone(ctx, folder, ctx.VSphereVM.Name, spec)
	if err != nil {
		return errors.Wrapf(err, "error trigging clone op for machine %s", ctx)
	}

	ctx.VSphereVM.Status.TaskRef = task.Reference().Value

	// patch the vsphereVM early to ensure that the task is
	// reflected in the status right away, this avoid situations
	// of concurrent clones
	if err := ctx.Patch(); err != nil {
		ctx.Logger.Error(err, "patch failed", "vspherevm", ctx.VSphereVM)
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

func getDiskSpec(ctx *context.VMContext, devices object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
	disks := devices.SelectByType((*types.VirtualDisk)(nil))
	if len(disks) == 0 {
		return nil, errors.Errorf("Invalid disk count: %d", len(disks))
	}

	// There is at least one disk
	var diskSpecs []types.BaseVirtualDeviceConfigSpec
	primaryDisk := disks[0].(*types.VirtualDisk) //nolint:forcetypeassert
	primaryCloneCapacityKB := int64(ctx.VSphereVM.Spec.DiskGiB) * 1024 * 1024
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
			if len(ctx.VSphereVM.Spec.AdditionalDisksGiB) > i {
				diskCloneCapacityKB = int64(ctx.VSphereVM.Spec.AdditionalDisksGiB[i]) * 1024 * 1024
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

const ethCardType = "vmxnet3"

func getNetworkSpecs(ctx *context.VMContext, devices object.VirtualDeviceList) ([]types.BaseVirtualDeviceConfigSpec, error) {
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
	for i := range ctx.VSphereVM.Spec.Network.Devices {
		netSpec := &ctx.VSphereVM.Spec.Network.Devices[i]
		ref, err := ctx.Session.Finder.Network(ctx, netSpec.NetworkName)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to find network %q", netSpec.NetworkName)
		}
		backing, err := ref.EthernetCardBackingInfo(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to create new ethernet card backing info for network %q on %q", netSpec.NetworkName, ctx)
		}
		dev, err := object.EthernetCardTypes().CreateEthernetCard(ethCardType, backing)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to create new ethernet card %q for network %q on %q", ethCardType, netSpec.NetworkName, ctx)
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
			ctx.Logger.V(4).Info("configured manual mac address", "mac-addr", nic.MacAddress)
		}

		// Assign a temporary device key to ensure that a unique one will be
		// generated when the device is created.
		nic.Key = key

		deviceSpecs = append(deviceSpecs, &types.VirtualDeviceConfigSpec{
			Device:    dev,
			Operation: types.VirtualDeviceConfigSpecOperationAdd,
		})
		ctx.Logger.V(4).Info("created network device", "eth-card-type", ethCardType, "network-spec", netSpec)
		key--
	}

	return deviceSpecs, nil
}
