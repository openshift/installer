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

package vnetpeerings

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// VnetPeeringSpec defines the specification for a virtual network peering.
type VnetPeeringSpec struct {
	SourceResourceGroup       string
	SourceVnetName            string
	RemoteResourceGroup       string
	RemoteVnetName            string
	PeeringName               string
	SubscriptionID            string
	AllowForwardedTraffic     *bool
	AllowGatewayTransit       *bool
	AllowVirtualNetworkAccess *bool
	UseRemoteGateways         *bool
}

// ResourceName returns the name of the virtual network peering.
func (s *VnetPeeringSpec) ResourceName() string {
	return s.PeeringName
}

// ResourceGroupName returns the name of the resource group.
func (s *VnetPeeringSpec) ResourceGroupName() string {
	return s.SourceResourceGroup
}

// OwnerResourceName is a no-op for virtual network peerings.
func (s *VnetPeeringSpec) OwnerResourceName() string {
	return s.SourceVnetName
}

// Parameters returns the parameters for the virtual network peering.
func (s *VnetPeeringSpec) Parameters(ctx context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(armnetwork.VirtualNetworkPeering); !ok {
			return nil, errors.Errorf("%T is not an armnetwork.VnetPeering", existing)
		}
		// virtual network peering already exists
		return nil, nil
	}
	vnetID := azure.VNetID(s.SubscriptionID, s.RemoteResourceGroup, s.RemoteVnetName)
	peeringProperties := armnetwork.VirtualNetworkPeeringPropertiesFormat{
		RemoteVirtualNetwork: &armnetwork.SubResource{
			ID: ptr.To(vnetID),
		},
		AllowForwardedTraffic:     s.AllowForwardedTraffic,
		AllowGatewayTransit:       s.AllowGatewayTransit,
		AllowVirtualNetworkAccess: s.AllowVirtualNetworkAccess,
		UseRemoteGateways:         s.UseRemoteGateways,
	}
	return armnetwork.VirtualNetworkPeering{
		Name:       ptr.To(s.PeeringName),
		Properties: &peeringProperties,
	}, nil
}
