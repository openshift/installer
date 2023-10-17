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

package subnets

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// SubnetSpec defines the specification for a Subnet.
type SubnetSpec struct {
	Name              string
	ResourceGroup     string
	SubscriptionID    string
	CIDRs             []string
	VNetName          string
	VNetResourceGroup string
	IsVNetManaged     bool
	RouteTableName    string
	SecurityGroupName string
	Role              infrav1.SubnetRole
	NatGatewayName    string
	ServiceEndpoints  infrav1.ServiceEndpoints
}

// ResourceName returns the name of the subnet.
func (s *SubnetSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group of the VNet that owns this subnet.
func (s *SubnetSpec) ResourceGroupName() string {
	return s.VNetResourceGroup
}

// OwnerResourceName returns the name of the VNet that owns this subnet.
func (s *SubnetSpec) OwnerResourceName() string {
	return s.VNetName
}

// Parameters returns the parameters for the subnet.
func (s *SubnetSpec) Parameters(ctx context.Context, existing interface{}) (parameters interface{}, err error) {
	if existing != nil {
		existingSubnet, ok := existing.(armnetwork.Subnet)
		if !ok {
			return nil, errors.Errorf("%T is not an armnetwork.Subnet", existing)
		}

		if !s.shouldUpdate(existingSubnet) {
			return nil, nil
		}
	}

	if !s.IsVNetManaged {
		// TODO: change this to terminal error once we add support for handling them
		return nil, errors.Errorf("custom vnet was provided but subnet %s is missing", s.Name)
	}
	subnetProperties := armnetwork.SubnetPropertiesFormat{
		AddressPrefixes: azure.PtrSlice(&s.CIDRs),
	}

	// workaround needed to avoid SubscriptionNotRegisteredForFeature for feature Microsoft.Network/AllowMultipleAddressPrefixesOnSubnet.
	if len(s.CIDRs) == 1 {
		subnetProperties = armnetwork.SubnetPropertiesFormat{
			AddressPrefix: &s.CIDRs[0],
		}
	}

	if s.RouteTableName != "" {
		subnetProperties.RouteTable = &armnetwork.RouteTable{
			ID: ptr.To(azure.RouteTableID(s.SubscriptionID, s.ResourceGroup, s.RouteTableName)),
		}
	}

	if s.NatGatewayName != "" {
		subnetProperties.NatGateway = &armnetwork.SubResource{
			ID: ptr.To(azure.NatGatewayID(s.SubscriptionID, s.ResourceGroup, s.NatGatewayName)),
		}
	}

	if s.SecurityGroupName != "" {
		subnetProperties.NetworkSecurityGroup = &armnetwork.SecurityGroup{
			ID: ptr.To(azure.SecurityGroupID(s.SubscriptionID, s.ResourceGroup, s.SecurityGroupName)),
		}
	}

	serviceEndpoints := make([]armnetwork.ServiceEndpointPropertiesFormat, 0, len(s.ServiceEndpoints))
	for _, se := range s.ServiceEndpoints {
		se := se
		serviceEndpoints = append(serviceEndpoints, armnetwork.ServiceEndpointPropertiesFormat{Service: ptr.To(se.Service), Locations: azure.PtrSlice(&se.Locations)})
	}
	subnetProperties.ServiceEndpoints = azure.PtrSlice(&serviceEndpoints)

	return armnetwork.Subnet{
		Properties: &subnetProperties,
	}, nil
}

// shouldUpdate returns true if an existing subnet should be updated.
func (s *SubnetSpec) shouldUpdate(existingSubnet armnetwork.Subnet) bool {
	// No modifications for non-managed subnets
	if !s.IsVNetManaged {
		return false
	}

	// Update the subnet a NAT Gateway was added for backwards compatibility.
	if s.NatGatewayName != "" && existingSubnet.Properties.NatGateway == nil {
		return true
	}

	// Update the subnet if the service endpoints changed.
	if existingSubnet.Properties.ServiceEndpoints != nil || len(s.ServiceEndpoints) > 0 {
		var existingServiceEndpoints []armnetwork.ServiceEndpointPropertiesFormat
		if existingSubnet.Properties.ServiceEndpoints != nil {
			for _, se := range existingSubnet.Properties.ServiceEndpoints {
				existingServiceEndpoints = append(existingServiceEndpoints, armnetwork.ServiceEndpointPropertiesFormat{Service: se.Service, Locations: se.Locations})
			}
		}
		newServiceEndpoints := make([]armnetwork.ServiceEndpointPropertiesFormat, len(s.ServiceEndpoints))
		for _, se := range s.ServiceEndpoints {
			se := se
			newServiceEndpoints = append(newServiceEndpoints, armnetwork.ServiceEndpointPropertiesFormat{Service: ptr.To(se.Service), Locations: azure.PtrSlice(&se.Locations)})
		}

		diff := cmp.Diff(newServiceEndpoints, existingServiceEndpoints)
		return diff != ""
	}
	return false
}
