/*
Copyright 2023 The Kubernetes Authors.

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
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	azcache "sigs.k8s.io/cloud-provider-azure/pkg/cache"
	"sigs.k8s.io/cloud-provider-azure/pkg/consts"
)

// GetVirtualMachineWithRetry invokes az.getVirtualMachine with exponential backoff retry
func (az *Cloud) GetVirtualMachineWithRetry(ctx context.Context, name types.NodeName, crt azcache.AzureCacheReadType) (*armcompute.VirtualMachine, error) {
	var machine *armcompute.VirtualMachine
	var retryErr error
	err := wait.ExponentialBackoff(az.RequestBackoff(), func() (bool, error) {
		machine, retryErr = az.getVirtualMachine(ctx, name, crt)
		if errors.Is(retryErr, cloudprovider.InstanceNotFound) {
			return true, cloudprovider.InstanceNotFound
		}
		if retryErr != nil {
			klog.Errorf("GetVirtualMachineWithRetry(%s): backoff failure, will retry, err=%v", name, retryErr)
			return false, nil
		}
		klog.V(2).Infof("GetVirtualMachineWithRetry(%s): backoff success", name)
		return true, nil
	})
	if errors.Is(err, wait.ErrWaitTimeout) {
		err = retryErr
	}
	return machine, err
}

// ListVirtualMachines invokes az.ComputeClientFactory.GetVirtualMachineScaleSetClient().List with exponential backoff retry
func (az *Cloud) ListVirtualMachines(ctx context.Context, resourceGroup string) ([]*armcompute.VirtualMachine, error) {
	allNodes, err := az.ComputeClientFactory.GetVirtualMachineClient().List(ctx, resourceGroup)
	if err != nil {
		klog.Errorf("ComputeClientFactory.GetVirtualMachineScaleSetClient().List(%v) failure with err=%v", resourceGroup, err)
		return nil, err
	}
	klog.V(6).Infof("ComputeClientFactory.GetVirtualMachineScaleSetClient().List(%v) success", resourceGroup)
	return allNodes, nil
}

// getPrivateIPsForMachine is wrapper for optional backoff getting private ips
// list of a node by name
func (az *Cloud) getPrivateIPsForMachine(ctx context.Context, nodeName types.NodeName) ([]string, error) {
	return az.getPrivateIPsForMachineWithRetry(ctx, nodeName)
}

func (az *Cloud) getPrivateIPsForMachineWithRetry(ctx context.Context, nodeName types.NodeName) ([]string, error) {
	var privateIPs []string
	err := wait.ExponentialBackoff(az.RequestBackoff(), func() (bool, error) {
		var retryErr error
		privateIPs, retryErr = az.VMSet.GetPrivateIPsByNodeName(ctx, string(nodeName))
		if retryErr != nil {
			// won't retry since the instance doesn't exist on Azure.
			if errors.Is(retryErr, cloudprovider.InstanceNotFound) {
				return true, retryErr
			}
			klog.Errorf("GetPrivateIPsByNodeName(%s): backoff failure, will retry,err=%v", nodeName, retryErr)
			return false, nil
		}
		klog.V(3).Infof("GetPrivateIPsByNodeName(%s): backoff success", nodeName)
		return true, nil
	})
	return privateIPs, err
}

func (az *Cloud) getIPForMachine(ctx context.Context, nodeName types.NodeName) (string, string, error) {
	return az.GetIPForMachineWithRetry(ctx, nodeName)
}

// GetIPForMachineWithRetry invokes az.getIPForMachine with exponential backoff retry
func (az *Cloud) GetIPForMachineWithRetry(ctx context.Context, name types.NodeName) (string, string, error) {
	var ip, publicIP string
	err := wait.ExponentialBackoffWithContext(ctx, az.RequestBackoff(), func(ctx context.Context) (bool, error) {
		var retryErr error
		ip, publicIP, retryErr = az.VMSet.GetIPByNodeName(ctx, string(name))
		if retryErr != nil {
			klog.Errorf("GetIPForMachineWithRetry(%s): backoff failure, will retry,err=%v", name, retryErr)
			return false, nil
		}
		klog.V(3).Infof("GetIPForMachineWithRetry(%s): backoff success", name)
		return true, nil
	})
	return ip, publicIP, err
}

func (az *Cloud) newVMCache() (azcache.Resource, error) {
	getter := func(ctx context.Context, key string) (interface{}, error) {
		// Currently InstanceView request are used by azure_zones, while the calls come after non-InstanceView
		// request. If we first send an InstanceView request and then a non InstanceView request, the second
		// request will still hit throttling. This is what happens now for cloud controller manager: In this
		// case we do get instance view every time to fulfill the azure_zones requirement without hitting
		// throttling.
		// Consider adding separate parameter for controlling 'InstanceView' once node update issue #56276 is fixed

		resourceGroup, err := az.GetNodeResourceGroup(key)
		if err != nil {
			return nil, err
		}

		vm, verr := az.ComputeClientFactory.GetVirtualMachineClient().Get(ctx, resourceGroup, key, nil)
		exists, rerr := checkResourceExistsFromError(verr)
		if rerr != nil {
			return nil, rerr
		}

		if !exists {
			klog.V(2).Infof("Virtual machine %q not found", key)
			return nil, nil
		}

		if vm != nil && vm.Properties != nil &&
			strings.EqualFold(ptr.Deref(vm.Properties.ProvisioningState, ""), string(consts.ProvisioningStateDeleting)) {
			klog.V(2).Infof("Virtual machine %q is under deleting", key)
			return nil, nil
		}

		return vm, nil
	}

	if az.VMCacheTTLInSeconds == 0 {
		az.VMCacheTTLInSeconds = vmCacheTTLDefaultInSeconds
	}
	return azcache.NewTimedCache(time.Duration(az.VMCacheTTLInSeconds)*time.Second, getter, az.Config.DisableAPICallCache)
}

// getVirtualMachine calls 'ComputeClientFactory.GetVirtualMachineScaleSetClient().Get' with a timed cache
// The service side has throttling control that delays responses if there are multiple requests onto certain vm
// resource request in short period.
func (az *Cloud) getVirtualMachine(ctx context.Context, nodeName types.NodeName, crt azcache.AzureCacheReadType) (vm *armcompute.VirtualMachine, err error) {
	vmName := string(nodeName)
	cachedVM, err := az.vmCache.Get(ctx, vmName, crt)
	if err != nil {
		return vm, err
	}

	if cachedVM == nil {
		klog.Warningf("Unable to find node %s: %v", nodeName, cloudprovider.InstanceNotFound)
		return vm, cloudprovider.InstanceNotFound
	}

	return (cachedVM.(*armcompute.VirtualMachine)), nil
}
