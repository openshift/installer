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

package natgateways

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
)

// NatGatewaySpec defines the specification for a NAT gateway.
type NatGatewaySpec struct {
	Name           string
	ResourceGroup  string
	SubscriptionID string
	Location       string
	NatGatewayIP   infrav1.PublicIPSpec
	ClusterName    string
	AdditionalTags infrav1.Tags
}

// ResourceName returns the name of the NAT gateway.
func (s *NatGatewaySpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *NatGatewaySpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for NAT gateways.
func (s *NatGatewaySpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the NAT gateway.
func (s *NatGatewaySpec) Parameters(ctx context.Context, existing interface{}) (params interface{}, err error) {
	if existing != nil {
		existingNatGateway, ok := existing.(armnetwork.NatGateway)
		if !ok {
			return nil, errors.Errorf("%T is not an armnetwork.NatGateway", existing)
		}

		if hasPublicIP(existingNatGateway, s.NatGatewayIP.Name) {
			// Skip update for NAT gateway as it exists with expected values
			return nil, nil
		}
	}

	natGatewayToCreate := armnetwork.NatGateway{
		Name:     ptr.To(s.Name),
		Location: ptr.To(s.Location),
		SKU:      &armnetwork.NatGatewaySKU{Name: ptr.To(armnetwork.NatGatewaySKUNameStandard)},
		Properties: &armnetwork.NatGatewayPropertiesFormat{
			PublicIPAddresses: []*armnetwork.SubResource{
				{
					ID: ptr.To(azure.PublicIPID(s.SubscriptionID, s.ResourceGroupName(), s.NatGatewayIP.Name)),
				},
			},
		},
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Additional:  s.AdditionalTags,
		})),
	}

	return natGatewayToCreate, nil
}

func hasPublicIP(natGateway armnetwork.NatGateway, publicIPName string) bool {
	for _, publicIP := range natGateway.Properties.PublicIPAddresses {
		if publicIP != nil && publicIP.ID != nil {
			resource, err := azureutil.ParseResourceID(*publicIP.ID)
			if err != nil {
				continue
			}
			if resource.Name == publicIPName {
				return true
			}
		}
	}
	return false
}
