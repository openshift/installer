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
	"strings"

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	NatGatewayName    string
	ServiceEndpoints  infrav1.ServiceEndpoints
}

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *SubnetSpec) ResourceRef() *asonetworkv1.VirtualNetworksSubnet {
	return &asonetworkv1.VirtualNetworksSubnet{
		ObjectMeta: metav1.ObjectMeta{
			// s.Name isn't unique per-cluster, so combine with vnet name to avoid collisions.
			// ToLower makes the name compatible with standard Kubernetes name requirements.
			Name: s.VNetName + "-" + strings.ToLower(s.Name),
		},
	}
}

// Parameters implements azure.ASOResourceSpecGetter.
func (s *SubnetSpec) Parameters(ctx context.Context, existing *asonetworkv1.VirtualNetworksSubnet) (parameters *asonetworkv1.VirtualNetworksSubnet, err error) {
	subnet := existing
	if subnet == nil {
		subnet = &asonetworkv1.VirtualNetworksSubnet{}
	}

	subnet.Spec = asonetworkv1.VirtualNetworks_Subnet_Spec{
		AzureName: s.Name,
		Owner: &genruntime.KnownResourceReference{
			Name: s.VNetName,
		},
		AddressPrefixes: s.CIDRs,
	}
	// workaround needed to avoid SubscriptionNotRegisteredForFeature for feature Microsoft.Network/AllowMultipleAddressPrefixesOnSubnet.
	if len(s.CIDRs) == 1 {
		subnet.Spec.AddressPrefix = &s.CIDRs[0]
	}

	if s.RouteTableName != "" {
		subnet.Spec.RouteTable = &asonetworkv1.RouteTableSpec_VirtualNetworks_Subnet_SubResourceEmbedded{
			Reference: &genruntime.ResourceReference{
				ARMID: azure.RouteTableID(s.SubscriptionID, s.VNetResourceGroup, s.RouteTableName),
			},
		}
	}

	if s.NatGatewayName != "" {
		subnet.Spec.NatGateway = &asonetworkv1.SubResource{
			Reference: &genruntime.ResourceReference{
				ARMID: azure.NatGatewayID(s.SubscriptionID, s.ResourceGroup, s.NatGatewayName),
			},
		}
	}

	if s.SecurityGroupName != "" {
		subnet.Spec.NetworkSecurityGroup = &asonetworkv1.NetworkSecurityGroupSpec_VirtualNetworks_Subnet_SubResourceEmbedded{
			Reference: &genruntime.ResourceReference{
				ARMID: azure.SecurityGroupID(s.SubscriptionID, s.VNetResourceGroup, s.SecurityGroupName),
			},
		}
	}

	//nolint:prealloc // pre-allocating this slice isn't going to make any meaningful performance difference
	// and makes it harder to keep this value nil when s.ServiceEndpoints is empty as is necessary.
	var serviceEndpoints []asonetworkv1.ServiceEndpointPropertiesFormat
	for _, se := range s.ServiceEndpoints {
		serviceEndpoints = append(serviceEndpoints, asonetworkv1.ServiceEndpointPropertiesFormat{Service: ptr.To(se.Service), Locations: se.Locations})
	}
	subnet.Spec.ServiceEndpoints = serviceEndpoints

	return subnet, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *SubnetSpec) WasManaged(resource *asonetworkv1.VirtualNetworksSubnet) bool {
	return s.IsVNetManaged
}
