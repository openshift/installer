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

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

// CreateOrUpdateInterface invokes az.ComputeClientFactory.GetInterfaceClient().CreateOrUpdate with exponential backoff retry
func (az *Cloud) CreateOrUpdateInterface(ctx context.Context, service *v1.Service, nic *armnetwork.Interface) error {
	_, rerr := az.ComputeClientFactory.GetInterfaceClient().CreateOrUpdate(ctx, az.ResourceGroup, *nic.Name, *nic)
	klog.V(10).Infof("InterfacesClient.CreateOrUpdate(%s): end", *nic.Name)
	if rerr != nil {
		klog.Errorf("InterfacesClient.CreateOrUpdate(%s) failed: %s", *nic.Name, rerr.Error())
		az.Event(service, v1.EventTypeWarning, "CreateOrUpdateInterface", rerr.Error())
		return rerr
	}

	return nil
}
