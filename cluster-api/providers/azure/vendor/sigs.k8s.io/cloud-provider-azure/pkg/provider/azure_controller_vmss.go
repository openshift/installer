/*
Copyright 2020 The Kubernetes Authors.

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

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	azcache "sigs.k8s.io/cloud-provider-azure/pkg/cache"
	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
	"sigs.k8s.io/cloud-provider-azure/pkg/util/errutils"
)

// AttachDisk attaches a disk to vm
func (ss *ScaleSet) AttachDisk(ctx context.Context, nodeName types.NodeName, diskMap map[string]*AttachDiskOptions) error {
	vmName := mapNodeNameToVMName(nodeName)
	vm, err := ss.getVmssVM(ctx, vmName, azcache.CacheReadTypeDefault)
	if err != nil {
		return err
	}

	nodeResourceGroup, err := ss.GetNodeResourceGroup(vmName)
	if err != nil {
		return err
	}

	var disks []*armcompute.DataDisk

	storageProfile := vm.AsVirtualMachineScaleSetVM().Properties.StorageProfile

	if storageProfile != nil && storageProfile.DataDisks != nil {
		disks = make([]*armcompute.DataDisk, len(storageProfile.DataDisks))
		copy(disks, storageProfile.DataDisks)
	}

	for k, v := range diskMap {
		diSKURI := k
		opt := v
		attached := false
		for _, disk := range storageProfile.DataDisks {
			if disk.ManagedDisk != nil && strings.EqualFold(*disk.ManagedDisk.ID, diSKURI) && disk.Lun != nil {
				if *disk.Lun == opt.Lun {
					attached = true
					break
				}
				return fmt.Errorf("disk(%s) already attached to node(%s) on LUN(%d), but target LUN is %d", diSKURI, nodeName, *disk.Lun, opt.Lun)

			}
		}
		if attached {
			klog.V(2).Infof("azureDisk - disk(%s) already attached to node(%s) on LUN(%d)", diSKURI, nodeName, opt.Lun)
			continue
		}

		managedDisk := &armcompute.ManagedDiskParameters{ID: &diSKURI}
		if opt.DiskEncryptionSetID == "" {
			if storageProfile.OSDisk != nil &&
				storageProfile.OSDisk.ManagedDisk != nil &&
				storageProfile.OSDisk.ManagedDisk.DiskEncryptionSet != nil &&
				storageProfile.OSDisk.ManagedDisk.DiskEncryptionSet.ID != nil {
				// set diskEncryptionSet as value of os disk by default
				opt.DiskEncryptionSetID = *storageProfile.OSDisk.ManagedDisk.DiskEncryptionSet.ID
			}
		}
		if opt.DiskEncryptionSetID != "" {
			managedDisk.DiskEncryptionSet = &armcompute.DiskEncryptionSetParameters{ID: &opt.DiskEncryptionSetID}
		}
		disks = append(disks,
			&armcompute.DataDisk{
				Name:                    &opt.DiskName,
				Lun:                     &opt.Lun,
				Caching:                 to.Ptr(opt.CachingMode),
				CreateOption:            to.Ptr(armcompute.DiskCreateOptionTypesAttach),
				ManagedDisk:             managedDisk,
				WriteAcceleratorEnabled: ptr.To(opt.WriteAcceleratorEnabled),
			})
	}

	newVM := &armcompute.VirtualMachineScaleSetVM{
		Properties: &armcompute.VirtualMachineScaleSetVMProperties{
			StorageProfile: &armcompute.StorageProfile{
				DataDisks: disks,
			},
		},
	}

	klog.V(2).Infof("azureDisk - update: rg(%s) vm(%s) - attach disk list(%+v)", nodeResourceGroup, nodeName, diskMap)
	result, rerr := ss.ComputeClientFactory.GetVirtualMachineScaleSetVMClient().Update(ctx, nodeResourceGroup, vm.VMSSName, vm.InstanceID, *newVM)
	if rerr != nil {
		klog.Errorf("azureDisk - attach disk list(%+v) on rg(%s) vm(%s) failed, err: %v", diskMap, nodeResourceGroup, nodeName, rerr)
		if exists, err := errutils.CheckResourceExistsFromAzcoreError(rerr); exists && err == nil {
			klog.Errorf("azureDisk - begin to filterNonExistingDisks(%v) on rg(%s) vm(%s)", diskMap, nodeResourceGroup, nodeName)
			disks := FilterNonExistingDisks(ctx, ss.ComputeClientFactory, newVM.Properties.StorageProfile.DataDisks)
			newVM.Properties.StorageProfile.DataDisks = disks
			result, rerr = ss.ComputeClientFactory.GetVirtualMachineScaleSetVMClient().Update(ctx, nodeResourceGroup, vm.VMSSName, vm.InstanceID, *newVM)
		}
	}

	klog.V(2).Infof("azureDisk - update: rg(%s) vm(%s) - attach disk list(%+v) returned with %v", nodeResourceGroup, nodeName, diskMap, rerr)
	if rerr != nil {
		return rerr
	}

	// clean node cache first and then update cache
	_ = ss.DeleteCacheForNode(ctx, vmName)

	var vmssName, instanceID string
	if result != nil && result.Properties != nil {
		// get vmssName, instanceID from cache first
		vm, err := ss.getVmssVM(ctx, vmName, azcache.CacheReadTypeDefault)
		if err == nil && vm != nil {
			vmssName = vm.VMSSName
			instanceID = vm.InstanceID
		} else {
			klog.Errorf("getVmssVM failed with error(%v) or nil pointer", err)
		}

		if err := ss.updateCache(ctx, vmName, nodeResourceGroup, vmssName, instanceID, result); err != nil {
			klog.Errorf("updateCache(%s, %s, %s, %s) failed with error: %v", vmName, nodeResourceGroup, vmssName, instanceID, err)
		}
	}
	return rerr
}

// DetachDisk detaches a disk from VM
func (ss *ScaleSet) DetachDisk(ctx context.Context, nodeName types.NodeName, diskMap map[string]string, forceDetach bool) error {
	vmName := mapNodeNameToVMName(nodeName)
	vm, err := ss.getVmssVM(ctx, vmName, azcache.CacheReadTypeDefault)
	if err != nil {
		return err
	}

	nodeResourceGroup, err := ss.GetNodeResourceGroup(vmName)
	if err != nil {
		return err
	}

	var disks []*armcompute.DataDisk

	if vm != nil && vm.VirtualMachineScaleSetVMProperties != nil {
		storageProfile := vm.VirtualMachineScaleSetVMProperties.StorageProfile
		if storageProfile != nil && storageProfile.DataDisks != nil {
			disks = make([]*armcompute.DataDisk, len(storageProfile.DataDisks))
			copy(disks, storageProfile.DataDisks)
		}
	}
	bFoundDisk := false
	for i, disk := range disks {
		for diSKURI, diskName := range diskMap {
			if disk.Lun != nil && (disk.Name != nil && diskName != "" && strings.EqualFold(*disk.Name, diskName)) ||
				(disk.Vhd != nil && disk.Vhd.URI != nil && diSKURI != "" && strings.EqualFold(*disk.Vhd.URI, diSKURI)) ||
				(disk.ManagedDisk != nil && diSKURI != "" && strings.EqualFold(*disk.ManagedDisk.ID, diSKURI)) {
				// found the disk
				klog.V(2).Infof("azureDisk - detach disk: name %s uri %s", diskName, diSKURI)
				disks[i].ToBeDetached = ptr.To(true)
				if forceDetach {
					disks[i].DetachOption = to.Ptr(armcompute.DiskDetachOptionTypesForceDetach)
				}
				bFoundDisk = true
			}
		}
	}

	if !bFoundDisk {
		// only log here, next action is to update VM status with original meta data
		klog.Warningf("detach azure disk on node(%s): disk list(%s) not found", nodeName, diskMap)
	} else {
		if strings.EqualFold(ss.Environment.Name, consts.AzureStackCloudName) && !ss.Config.DisableAzureStackCloud {
			// Azure stack does not support ToBeDetached flag, use original way to detach disk
			var newDisks []*armcompute.DataDisk
			for _, disk := range disks {
				if !ptr.Deref(disk.ToBeDetached, false) {
					newDisks = append(newDisks, disk)
				}
			}
			disks = newDisks
		}
	}

	newVM := &armcompute.VirtualMachineScaleSetVM{
		Properties: &armcompute.VirtualMachineScaleSetVMProperties{
			StorageProfile: &armcompute.StorageProfile{
				DataDisks: disks,
			},
		},
	}

	var result *armcompute.VirtualMachineScaleSetVM
	var rerr error

	defer func() {
		_ = ss.DeleteCacheForNode(ctx, vmName)
		if rerr == nil {
			if err := ss.updateCache(ctx, vmName, nodeResourceGroup, vm.VMSSName, vm.InstanceID, result); err != nil {
				klog.Errorf("updateCache(%s, %s, %s, %s) failed with error: %v", vmName, nodeResourceGroup, vm.VMSSName, vm.InstanceID, err)
			}
		}
	}()

	klog.V(2).Infof("azureDisk - update(%s): vm(%s) - detach disk list(%s)", nodeResourceGroup, nodeName, diskMap)
	result, rerr = ss.ComputeClientFactory.GetVirtualMachineScaleSetVMClient().Update(ctx, nodeResourceGroup, vm.VMSSName, vm.InstanceID, *newVM)
	if rerr != nil {
		klog.Errorf("azureDisk - detach disk list(%s) on rg(%s) vm(%s) failed, err: %v", diskMap, nodeResourceGroup, nodeName, rerr)
		if exists, err := errutils.CheckResourceExistsFromAzcoreError(rerr); exists && err == nil {
			klog.Errorf("azureDisk - begin to filterNonExistingDisks(%v) on rg(%s) vm(%s)", diskMap, nodeResourceGroup, nodeName)
			disks := FilterNonExistingDisks(ctx, ss.ComputeClientFactory, newVM.Properties.StorageProfile.DataDisks)
			newVM.Properties.StorageProfile.DataDisks = disks
			result, rerr = ss.ComputeClientFactory.GetVirtualMachineScaleSetVMClient().Update(ctx, nodeResourceGroup, vm.VMSSName, vm.InstanceID, *newVM)
		}
	}

	klog.V(2).Infof("azureDisk - update(%s): vm(%s) - detach disk(%v) returned with %v", nodeResourceGroup, nodeName, diskMap, rerr)
	return rerr
}

// UpdateVM updates a vm
func (ss *ScaleSet) UpdateVM(ctx context.Context, nodeName types.NodeName) error {
	vmName := mapNodeNameToVMName(nodeName)
	vm, err := ss.getVmssVM(ctx, vmName, azcache.CacheReadTypeDefault)
	if err != nil {
		return err
	}

	nodeResourceGroup, err := ss.GetNodeResourceGroup(vmName)
	if err != nil {
		return err
	}

	_, rerr := ss.ComputeClientFactory.GetVirtualMachineScaleSetVMClient().Update(ctx, nodeResourceGroup, vm.VMSSName, vm.InstanceID, armcompute.VirtualMachineScaleSetVM{})
	if rerr != nil {
		return rerr
	}
	return nil
}

// GetDataDisks gets a list of data disks attached to the node.
func (ss *ScaleSet) GetDataDisks(ctx context.Context, nodeName types.NodeName, crt azcache.AzureCacheReadType) ([]*armcompute.DataDisk, *string, error) {
	vm, err := ss.getVmssVM(ctx, string(nodeName), crt)
	if err != nil {
		return nil, nil, err
	}

	if vm != nil && vm.AsVirtualMachineScaleSetVM() != nil && vm.AsVirtualMachineScaleSetVM().Properties != nil {
		storageProfile := vm.AsVirtualMachineScaleSetVM().Properties.StorageProfile

		if storageProfile == nil || storageProfile.DataDisks == nil {
			return nil, nil, nil
		}
		return storageProfile.DataDisks, vm.AsVirtualMachineScaleSetVM().Properties.ProvisioningState, nil
	}

	return nil, nil, nil
}
