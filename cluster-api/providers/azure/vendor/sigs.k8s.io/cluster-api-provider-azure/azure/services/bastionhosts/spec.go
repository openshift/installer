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

package bastionhosts

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// AzureBastionSpec defines the specification for azure bastion feature.
type AzureBastionSpec struct {
	Name            string
	ResourceGroup   string
	Location        string
	ClusterName     string
	SubnetID        string
	PublicIPID      string
	Sku             infrav1.BastionHostSkuName
	EnableTunneling bool
}

// AzureBastionSpecInput defines the required inputs to construct an azure bastion spec.
type AzureBastionSpecInput struct {
	SubnetName   string
	PublicIPName string
	VNetName     string
}

// ResourceName returns the name of the bastion host.
func (s *AzureBastionSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *AzureBastionSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for bastion hosts.
func (s *AzureBastionSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the bastion host.
func (s *AzureBastionSpec) Parameters(ctx context.Context, existing interface{}) (parameters interface{}, err error) {
	if existing != nil {
		if _, ok := existing.(armnetwork.BastionHost); !ok {
			return nil, errors.Errorf("%T is not an armnetwork.BastionHost", existing)
		}
		// bastion host already exists
		return nil, nil
	}

	bastionHostIPConfigName := fmt.Sprintf("%s-%s", s.Name, "bastionIP")

	return armnetwork.BastionHost{
		Name:     ptr.To(s.Name),
		Location: ptr.To(s.Location),
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Role:        ptr.To("Bastion"),
		})),
		SKU: &armnetwork.SKU{
			Name: ptr.To(armnetwork.BastionHostSKUName(s.Sku)),
		},
		Properties: &armnetwork.BastionHostPropertiesFormat{
			EnableTunneling: ptr.To(s.EnableTunneling),
			DNSName:         ptr.To(fmt.Sprintf("%s-bastion", strings.ToLower(s.Name))),
			IPConfigurations: []*armnetwork.BastionHostIPConfiguration{
				{
					Name: ptr.To(bastionHostIPConfigName),
					Properties: &armnetwork.BastionHostIPConfigurationPropertiesFormat{
						Subnet: &armnetwork.SubResource{
							ID: &s.SubnetID,
						},
						PublicIPAddress: &armnetwork.SubResource{
							ID: &s.PublicIPID,
						},
						PrivateIPAllocationMethod: ptr.To(armnetwork.IPAllocationMethodDynamic),
					},
				},
			},
		},
	}, nil
}
