/*
Copyright 2022 The Kubernetes Authors.

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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	azcache "sigs.k8s.io/cloud-provider-azure/pkg/cache"
	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
)

// AttachDisk attaches a disk to vm
func (fs *FlexScaleSet) AttachDisk(ctx context.Context, nodeName types.NodeName, diskMap map[string]*AttachDiskOptions) error {
	vmName := mapNodeNameToVMName(nodeName)
	vm, err := fs.getVmssFlexVM(ctx, vmName, azcache.CacheReadTypeDefault)
	if err != nil {
		return err
	}

	nodeResourceGroup, err := fs.GetNodeResourceGroup(vmName)
	if err != nil {
		return err
	}

	disks := make([]*armcompute.DataDisk, len(vm.Properties.StorageProfile.DataDisks))
	copy(disks, vm.Properties.StorageProfile.DataDisks)

	for k, v := range diskMap {
		diSKURI := k
		opt := v
		attached := false
		for _, disk := range vm.Properties.StorageProfile.DataDisks {
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
			if vm.Properties.StorageProfile.OSDisk != nil &&
				vm.Properties.StorageProfile.OSDisk.ManagedDisk != nil &&
				vm.Properties.StorageProfile.OSDisk.ManagedDisk.DiskEncryptionSet != nil &&
				vm.Properties.StorageProfile.OSDisk.ManagedDisk.DiskEncryptionSet.ID != nil {
				// set diskEncryptionSet as value of os disk by default
				opt.DiskEncryptionSetID = *vm.Properties.StorageProfile.OSDisk.ManagedDisk.DiskEncryptionSet.ID
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

	newVM := armcompute.VirtualMachine{
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				DataDisks: disks,
			},
		},
		Location: vm.Location,
	}

	klog.V(2).Infof("azureDisk - update: rg(%s) vm(%s) - attach disk list(%+v)", nodeResourceGroup, vmName, diskMap)
	result, err := fs.ComputeClientFactory.GetVirtualMachineClient().CreateOrUpdate(ctx, nodeResourceGroup, *vm.Name, newVM)
	var rerr *azcore.ResponseError
	if err != nil && errors.As(err, &rerr) {
		klog.Errorf("azureDisk - attach disk list(%+v) on rg(%s) vm(%s) failed, err: %v", diskMap, nodeResourceGroup, vmName, rerr)
		if rerr.StatusCode == http.StatusNotFound {
			klog.Errorf("azureDisk - begin to filterNonExistingDisks(%v) on rg(%s) vm(%s)", diskMap, nodeResourceGroup, vmName)
			disks := FilterNonExistingDisks(ctx, fs.ComputeClientFactory, newVM.Properties.StorageProfile.DataDisks)
			newVM.Properties.StorageProfile.DataDisks = disks
			result, err = fs.ComputeClientFactory.GetVirtualMachineClient().CreateOrUpdate(ctx, nodeResourceGroup, *vm.Name, newVM)
		}
	}

	klog.V(2).Infof("azureDisk - update(%s): vm(%s) - attach disk list(%+v) returned with %v", nodeResourceGroup, vmName, diskMap, rerr)
	if err != nil {
		return err
	}
	_ = fs.DeleteCacheForNode(ctx, vmName)
	if err := fs.updateCache(ctx, vmName, result); err != nil {
		klog.Errorf("updateCache(%s) failed with error: %v", vmName, err)
	}
	return nil
}

// DetachDisk detaches a disk from VM
func (fs *FlexScaleSet) DetachDisk(ctx context.Context, nodeName types.NodeName, diskMap map[string]string, forceDetach bool) error {
	vmName := mapNodeNameToVMName(nodeName)
	vm, err := fs.getVmssFlexVM(ctx, vmName, azcache.CacheReadTypeDefault)
	if err != nil {
		// if host doesn't exist, no need to detach
		klog.Warningf("azureDisk - cannot find node %s, skip detaching disk list(%s)", nodeName, diskMap)
		return nil
	}

	nodeResourceGroup, err := fs.GetNodeResourceGroup(vmName)
	if err != nil {
		return err
	}

	disks := make([]*armcompute.DataDisk, len(vm.Properties.StorageProfile.DataDisks))
	copy(disks, vm.Properties.StorageProfile.DataDisks)

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
		if strings.EqualFold(fs.Environment.Name, consts.AzureStackCloudName) && !fs.Config.DisableAzureStackCloud {
			// Azure stack does not support ToBeDetached flag, use original way to detach disk
			newDisks := []*armcompute.DataDisk{}
			for _, disk := range disks {
				if !ptr.Deref(disk.ToBeDetached, false) {
					newDisks = append(newDisks, disk)
				}
			}
			disks = newDisks
		}
	}

	newVM := armcompute.VirtualMachine{
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				DataDisks: disks,
			},
		},
		Location: vm.Location,
	}

	var result *armcompute.VirtualMachine
	defer func() {
		_ = fs.DeleteCacheForNode(ctx, vmName)
		if err == nil {
			if err := fs.updateCache(ctx, vmName, result); err != nil {
				klog.Errorf("updateCache(%s) failed with error: %v", vmName, err)
			}
		}
	}()

	klog.V(2).Infof("azureDisk - update(%s): vm(%s) node(%s)- detach disk list(%s)", nodeResourceGroup, vmName, nodeName, diskMap)

	result, err = fs.ComputeClientFactory.GetVirtualMachineClient().CreateOrUpdate(ctx, nodeResourceGroup, *vm.Name, newVM)
	if err != nil {
		klog.Errorf("azureDisk - detach disk list(%s) on rg(%s) vm(%s) failed, err: %v", diskMap, nodeResourceGroup, vmName, err)
		var rerr *azcore.ResponseError
		if errors.As(err, &rerr) {
			if rerr.StatusCode == http.StatusNotFound {
				klog.Errorf("azureDisk - begin to filterNonExistingDisks(%v) on rg(%s) vm(%s)", diskMap, nodeResourceGroup, vmName)
				disks := FilterNonExistingDisks(ctx, fs.ComputeClientFactory, vm.Properties.StorageProfile.DataDisks)
				newVM.Properties.StorageProfile.DataDisks = disks
				result, err = fs.ComputeClientFactory.GetVirtualMachineClient().CreateOrUpdate(ctx, nodeResourceGroup, *vm.Name, newVM)
			}
		}
	}

	klog.V(2).Infof("azureDisk - update(%s): vm(%s) - detach disk list(%s) returned with %v", nodeResourceGroup, vmName, diskMap, err)
	if err != nil {
		return err
	}
	// clean node cache first and then update cache
	_ = fs.DeleteCacheForNode(ctx, vmName)
	if err := fs.updateCache(ctx, vmName, result); err != nil {
		klog.Errorf("updateCache(%s) failed with error: %v", vmName, err)
	}
	return nil
}

// UpdateVM updates a vm
func (fs *FlexScaleSet) UpdateVM(ctx context.Context, nodeName types.NodeName) error {
	vmName := mapNodeNameToVMName(nodeName)
	vm, err := fs.getVmssFlexVM(ctx, vmName, azcache.CacheReadTypeDefault)
	if err != nil {
		// if host doesn't exist, no need to update
		klog.Warningf("azureDisk - cannot find node %s, skip updating vm", nodeName)
		return nil
	}
	nodeResourceGroup, err := fs.GetNodeResourceGroup(vmName)
	if err != nil {
		return err
	}

	_, rerr := fs.ComputeClientFactory.GetVirtualMachineClient().CreateOrUpdate(ctx, nodeResourceGroup, *vm.Name, armcompute.VirtualMachine{})
	if rerr != nil {
		return rerr
	}
	return nil
}

func (fs *FlexScaleSet) updateCache(ctx context.Context, nodeName string, vm *armcompute.VirtualMachine) error {
	if nodeName == "" {
		return fmt.Errorf("nodeName is empty")
	}
	if vm == nil {
		return fmt.Errorf("vm is nil")
	}
	if vm.Name == nil {
		return fmt.Errorf("vm.Name is nil")
	}
	if vm.Properties == nil {
		return fmt.Errorf("vm.Properties is nil")
	}
	if vm.Properties.OSProfile == nil || vm.Properties.OSProfile.ComputerName == nil {
		return fmt.Errorf("vm.Properties.OSProfile.ComputerName is nil")
	}

	vmssFlexID, err := fs.getNodeVmssFlexID(ctx, nodeName)
	if err != nil {
		return err
	}

	fs.lockMap.LockEntry(vmssFlexID)
	defer fs.lockMap.UnlockEntry(vmssFlexID)
	cached, err := fs.vmssFlexVMCache.Get(ctx, vmssFlexID, azcache.CacheReadTypeDefault)
	if err != nil {
		return err
	}
	vmMap := cached.(*sync.Map)
	vmMap.Store(nodeName, vm)

	fs.vmssFlexVMNameToVmssID.Store(strings.ToLower(*vm.Properties.OSProfile.ComputerName), vmssFlexID)
	fs.vmssFlexVMNameToNodeName.Store(*vm.Name, strings.ToLower(*vm.Properties.OSProfile.ComputerName))
	klog.V(2).Infof("updateCache(%s) for vmssFlexID(%s) successfully", nodeName, vmssFlexID)
	return nil
}

// GetDataDisks gets a list of data disks attached to the node.
func (fs *FlexScaleSet) GetDataDisks(ctx context.Context, nodeName types.NodeName, crt azcache.AzureCacheReadType) ([]*armcompute.DataDisk, *string, error) {
	vm, err := fs.getVmssFlexVM(ctx, string(nodeName), crt)
	if err != nil {
		return nil, nil, err
	}

	if vm.Properties.StorageProfile.DataDisks == nil {
		return nil, nil, nil
	}
	return vm.Properties.StorageProfile.DataDisks, vm.Properties.ProvisioningState, nil
}
