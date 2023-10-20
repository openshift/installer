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

package virtualnetworks

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// VNetSpec defines the specification for a Virtual Network.
type VNetSpec struct {
	ResourceGroup    string
	Name             string
	CIDRs            []string
	Location         string
	ExtendedLocation *infrav1.ExtendedLocationSpec
	ClusterName      string
	AdditionalTags   infrav1.Tags
}

// ResourceName returns the name of the vnet.
func (s *VNetSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *VNetSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for vnets.
func (s *VNetSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the vnet.
func (s *VNetSpec) Parameters(ctx context.Context, existing interface{}) (interface{}, error) {
	if existing != nil {
		// vnet already exists, nothing to update.
		return nil, nil
	}
	return armnetwork.VirtualNetwork{
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Role:        ptr.To(infrav1.CommonRole),
			Additional:  s.AdditionalTags,
		})),
		Location:         ptr.To(s.Location),
		ExtendedLocation: converters.ExtendedLocationToNetworkSDK(s.ExtendedLocation),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: azure.PtrSlice(&s.CIDRs),
			},
		},
	}, nil
}
