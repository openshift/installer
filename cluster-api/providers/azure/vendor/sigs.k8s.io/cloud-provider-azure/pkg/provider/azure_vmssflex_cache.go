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
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	azcache "sigs.k8s.io/cloud-provider-azure/pkg/cache"
	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
	"sigs.k8s.io/cloud-provider-azure/pkg/util/errutils"
)

func (fs *FlexScaleSet) newVmssFlexCache() (azcache.Resource, error) {
	getter := func(ctx context.Context, _ string) (interface{}, error) {
		localCache := &sync.Map{}

		allResourceGroups, err := fs.GetResourceGroups()
		if err != nil {
			return nil, err
		}

		for _, resourceGroup := range allResourceGroups.UnsortedList() {
			allScaleSets, rerr := fs.ComputeClientFactory.GetVirtualMachineScaleSetClient().List(ctx, resourceGroup)
			if rerr != nil {
				if exists, err := errutils.CheckResourceExistsFromAzcoreError(rerr); !exists && err == nil {
					klog.Warningf("Skip caching vmss for resource group %s due to error: %v", resourceGroup, rerr.Error())
					continue
				}
				klog.Errorf("VirtualMachineScaleSetsClient.List failed: %v", rerr)
				return nil, rerr
			}

			for i := range allScaleSets {
				scaleSet := allScaleSets[i]
				if scaleSet.ID == nil || *scaleSet.ID == "" {
					klog.Warning("failed to get the ID of VMSS Flex")
					continue
				}

				if *scaleSet.Properties.OrchestrationMode == armcompute.OrchestrationModeFlexible {
					localCache.Store(*scaleSet.ID, scaleSet)
				}
			}
		}

		return localCache, nil
	}

	if fs.Config.VmssFlexCacheTTLInSeconds == 0 {
		fs.Config.VmssFlexCacheTTLInSeconds = consts.VmssFlexCacheTTLDefaultInSeconds
	}
	return azcache.NewTimedCache(time.Duration(fs.Config.VmssFlexCacheTTLInSeconds)*time.Second, getter, fs.Cloud.Config.DisableAPICallCache)
}

func (fs *FlexScaleSet) newVmssFlexVMCache() (azcache.Resource, error) {
	getter := func(ctx context.Context, key string) (interface{}, error) {
		localCache := &sync.Map{}
		armResource, err := arm.ParseResourceID("/" + key)
		if err != nil {
			return nil, err
		}
		vms, rerr := fs.ComputeClientFactory.GetVirtualMachineClient().ListVmssFlexVMsWithOutInstanceView(ctx, armResource.ResourceGroupName, key)
		if rerr != nil {
			klog.Errorf("List failed: %v", rerr)
			return nil, rerr
		}

		for i := range vms {
			vm := vms[i]
			if vm.Properties.OSProfile != nil && vm.Properties.OSProfile.ComputerName != nil {
				localCache.Store(strings.ToLower(*vm.Properties.OSProfile.ComputerName), vm)
				fs.vmssFlexVMNameToVmssID.Store(strings.ToLower(*vm.Properties.OSProfile.ComputerName), key)
				fs.vmssFlexVMNameToNodeName.Store(*vm.Name, strings.ToLower(*vm.Properties.OSProfile.ComputerName))
			}
		}

		vms, rerr = fs.ComputeClientFactory.GetVirtualMachineClient().ListVmssFlexVMsWithOnlyInstanceView(ctx, armResource.ResourceGroupName, key)
		if rerr != nil {
			klog.Errorf("ListVMInstanceView failed: %v", rerr)
			return nil, rerr
		}

		for i := range vms {
			vm := vms[i]
			if vm.Name != nil {
				nodeName, ok := fs.vmssFlexVMNameToNodeName.Load(*vm.Name)
				if !ok {
					continue
				}

				cached, ok := localCache.Load(nodeName)
				if ok {
					cachedVM := cached.(*armcompute.VirtualMachine)
					cachedVM.Properties.InstanceView = vm.Properties.InstanceView
				}
			}
		}

		return localCache, nil
	}

	if fs.Config.VmssFlexVMCacheTTLInSeconds == 0 {
		fs.Config.VmssFlexVMCacheTTLInSeconds = consts.VmssFlexVMCacheTTLDefaultInSeconds
	}
	return azcache.NewTimedCache(time.Duration(fs.Config.VmssFlexVMCacheTTLInSeconds)*time.Second, getter, fs.Cloud.Config.DisableAPICallCache)
}

func (fs *FlexScaleSet) getNodeNameByVMName(ctx context.Context, vmName string) (string, error) {
	fs.lockMap.LockEntry(consts.GetNodeVmssFlexIDLockKey)
	defer fs.lockMap.UnlockEntry(consts.GetNodeVmssFlexIDLockKey)
	cachedNodeName, isCached := fs.vmssFlexVMNameToNodeName.Load(vmName)
	if isCached {
		return fmt.Sprintf("%v", cachedNodeName), nil
	}

	getter := func(ctx context.Context, vmName string, crt azcache.AzureCacheReadType) (string, error) {
		cached, err := fs.vmssFlexCache.Get(ctx, consts.VmssFlexKey, crt)
		if err != nil {
			return "", err
		}
		vmssFlexes := cached.(*sync.Map)

		vmssFlexes.Range(func(key, _ interface{}) bool {
			vmssFlexID := key.(string)
			_, err := fs.vmssFlexVMCache.Get(ctx, vmssFlexID, azcache.CacheReadTypeForceRefresh)
			if err != nil {
				klog.Errorf("failed to refresh vmss flex VM cache for vmssFlexID %s", vmssFlexID)
			}
			return true
		})

		cachedNodeName, isCached = fs.vmssFlexVMNameToNodeName.Load(vmName)
		if isCached {
			return fmt.Sprintf("%v", cachedNodeName), nil
		}
		return "", cloudprovider.InstanceNotFound
	}

	nodeName, err := getter(ctx, vmName, azcache.CacheReadTypeDefault)
	if errors.Is(err, cloudprovider.InstanceNotFound) {
		klog.V(2).Infof("Could not find node (%s) in the existing cache. Forcely freshing the cache to check again...", nodeName)
		return getter(ctx, vmName, azcache.CacheReadTypeForceRefresh)
	}
	return nodeName, err

}

func (fs *FlexScaleSet) getNodeVmssFlexID(ctx context.Context, nodeName string) (string, error) {
	fs.lockMap.LockEntry(consts.GetNodeVmssFlexIDLockKey)
	defer fs.lockMap.UnlockEntry(consts.GetNodeVmssFlexIDLockKey)
	cachedVmssFlexID, isCached := fs.vmssFlexVMNameToVmssID.Load(nodeName)

	if isCached {
		return fmt.Sprintf("%v", cachedVmssFlexID), nil
	}

	getter := func(ctx context.Context, nodeName string, crt azcache.AzureCacheReadType) (string, error) {
		cached, err := fs.vmssFlexCache.Get(ctx, consts.VmssFlexKey, crt)
		if err != nil {
			return "", err
		}
		vmssFlexes := cached.(*sync.Map)

		var vmssFlexIDs []string
		vmssFlexes.Range(func(key, value interface{}) bool {
			vmssFlexID := key.(string)
			vmssFlex := value.(*armcompute.VirtualMachineScaleSet)
			vmssPrefix := ptr.Deref(vmssFlex.Name, "")
			if vmssFlex.Properties.VirtualMachineProfile != nil &&
				vmssFlex.Properties.VirtualMachineProfile.OSProfile != nil &&
				vmssFlex.Properties.VirtualMachineProfile.OSProfile.ComputerNamePrefix != nil {
				vmssPrefix = ptr.Deref(vmssFlex.Properties.VirtualMachineProfile.OSProfile.ComputerNamePrefix, "")
			}
			if strings.EqualFold(vmssPrefix, nodeName[:len(nodeName)-6]) {
				// we should check this vmss first since nodeName and vmssFlex.Name or
				// ComputerNamePrefix belongs to same vmss, so prepend here
				vmssFlexIDs = append([]string{vmssFlexID}, vmssFlexIDs...)
			} else {
				vmssFlexIDs = append(vmssFlexIDs, vmssFlexID)
			}
			return true
		})

		for _, vmssID := range vmssFlexIDs {
			if _, err := fs.vmssFlexVMCache.Get(ctx, vmssID, azcache.CacheReadTypeForceRefresh); err != nil {
				klog.Errorf("failed to refresh vmss flex VM cache for vmssFlexID %s", vmssID)
			}
			// if the vm is cached stop refreshing
			cachedVmssFlexID, isCached = fs.vmssFlexVMNameToVmssID.Load(nodeName)
			if isCached {
				return fmt.Sprintf("%v", cachedVmssFlexID), nil
			}
		}
		return "", cloudprovider.InstanceNotFound
	}

	vmssFlexID, err := getter(ctx, nodeName, azcache.CacheReadTypeDefault)
	if errors.Is(err, cloudprovider.InstanceNotFound) {
		klog.V(2).Infof("Could not find node (%s) in the existing cache. Forcely freshing the cache to check again...", nodeName)
		return getter(ctx, nodeName, azcache.CacheReadTypeForceRefresh)
	}
	return vmssFlexID, err

}

func (fs *FlexScaleSet) getVmssFlexVM(ctx context.Context, nodeName string, crt azcache.AzureCacheReadType) (vm *armcompute.VirtualMachine, err error) {
	vmssFlexID, err := fs.getNodeVmssFlexID(ctx, nodeName)
	if err != nil {
		return vm, err
	}

	cached, err := fs.vmssFlexVMCache.Get(ctx, vmssFlexID, crt)
	if err != nil {
		return vm, err
	}
	vmMap := cached.(*sync.Map)
	cachedVM, ok := vmMap.Load(nodeName)
	if !ok {
		klog.V(2).Infof("did not find node (%s) in the existing cache, which means it is deleted...", nodeName)
		return vm, cloudprovider.InstanceNotFound
	}

	return (cachedVM.(*armcompute.VirtualMachine)), nil
}

func (fs *FlexScaleSet) getVmssFlexByVmssFlexID(ctx context.Context, vmssFlexID string, crt azcache.AzureCacheReadType) (*armcompute.VirtualMachineScaleSet, error) {
	cached, err := fs.vmssFlexCache.Get(ctx, consts.VmssFlexKey, crt)
	if err != nil {
		return nil, err
	}
	vmssFlexes := cached.(*sync.Map)
	if vmssFlex, ok := vmssFlexes.Load(vmssFlexID); ok {
		result := vmssFlex.(*armcompute.VirtualMachineScaleSet)
		return result, nil
	}

	klog.V(2).Infof("Couldn't find VMSS Flex with ID %s, refreshing the cache", vmssFlexID)
	cached, err = fs.vmssFlexCache.Get(ctx, consts.VmssFlexKey, azcache.CacheReadTypeForceRefresh)
	if err != nil {
		return nil, err
	}
	vmssFlexes = cached.(*sync.Map)
	if vmssFlex, ok := vmssFlexes.Load(vmssFlexID); ok {
		result := vmssFlex.(*armcompute.VirtualMachineScaleSet)
		return result, nil
	}
	return nil, cloudprovider.InstanceNotFound
}

func (fs *FlexScaleSet) getVmssFlexByNodeName(ctx context.Context, nodeName string, crt azcache.AzureCacheReadType) (*armcompute.VirtualMachineScaleSet, error) {
	vmssFlexID, err := fs.getNodeVmssFlexID(ctx, nodeName)
	if err != nil {
		return nil, err
	}
	vmssFlex, err := fs.getVmssFlexByVmssFlexID(ctx, vmssFlexID, crt)
	if err != nil {
		return nil, err
	}
	return vmssFlex, nil
}

func (fs *FlexScaleSet) getVmssFlexIDByName(ctx context.Context, vmssFlexName string) (string, error) {
	cached, err := fs.vmssFlexCache.Get(ctx, consts.VmssFlexKey, azcache.CacheReadTypeDefault)
	if err != nil {
		return "", err
	}
	var targetVmssFlexID string
	vmssFlexes := cached.(*sync.Map)
	vmssFlexes.Range(func(key, _ interface{}) bool {
		vmssFlexID := key.(string)
		name, err := getLastSegment(vmssFlexID, "/")
		if err != nil {
			return true
		}
		if strings.EqualFold(name, vmssFlexName) {
			targetVmssFlexID = vmssFlexID
			return false
		}
		return true
	})
	if targetVmssFlexID != "" {
		return targetVmssFlexID, nil
	}
	return "", cloudprovider.InstanceNotFound
}

func (fs *FlexScaleSet) getVmssFlexByName(ctx context.Context, vmssFlexName string) (*armcompute.VirtualMachineScaleSet, error) {
	cached, err := fs.vmssFlexCache.Get(ctx, consts.VmssFlexKey, azcache.CacheReadTypeDefault)
	if err != nil {
		return nil, err
	}

	var targetVmssFlex *armcompute.VirtualMachineScaleSet
	vmssFlexes := cached.(*sync.Map)
	vmssFlexes.Range(func(key, value interface{}) bool {
		vmssFlexID := key.(string)
		vmssFlex := value.(*armcompute.VirtualMachineScaleSet)
		name, err := getLastSegment(vmssFlexID, "/")
		if err != nil {
			return true
		}
		if strings.EqualFold(name, vmssFlexName) {
			targetVmssFlex = vmssFlex
			return false
		}
		return true
	})
	if targetVmssFlex != nil {
		return targetVmssFlex, nil
	}
	return nil, cloudprovider.InstanceNotFound
}

func (fs *FlexScaleSet) DeleteCacheForNode(ctx context.Context, nodeName string) error {
	if fs.Config.DisableAPICallCache {
		return nil
	}
	vmssFlexID, err := fs.getNodeVmssFlexID(ctx, nodeName)
	if err != nil {
		klog.Errorf("getNodeVmssFlexID(%s) failed with %v", nodeName, err)
		return err
	}

	fs.lockMap.LockEntry(vmssFlexID)
	defer fs.lockMap.UnlockEntry(vmssFlexID)
	cached, err := fs.vmssFlexVMCache.Get(ctx, vmssFlexID, azcache.CacheReadTypeDefault)
	if err != nil {
		klog.Errorf("vmssFlexVMCache.Get(%s, %s) failed with %v", vmssFlexID, nodeName, err)
		return err
	}
	if cached == nil {
		err := fmt.Errorf("nil cache returned from %s", vmssFlexID)
		klog.Errorf("DeleteCacheForNode(%s, %s) failed with %v", vmssFlexID, nodeName, err)
		return err
	}
	vmMap := cached.(*sync.Map)
	vmMap.Delete(nodeName)

	fs.vmssFlexVMCache.Update(vmssFlexID, vmMap)
	fs.vmssFlexVMNameToVmssID.Delete(nodeName)

	klog.V(2).Infof("DeleteCacheForNode(%s, %s) successfully", vmssFlexID, nodeName)
	return nil
}
