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

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// NatGatewaySpec defines the specification for a NAT gateway.
type NatGatewaySpec struct {
	Name           string
	Namespace      string
	ResourceGroup  string
	SubscriptionID string
	Location       string
	NatGatewayIP   infrav1.PublicIPSpec
	ClusterName    string
	AdditionalTags infrav1.Tags
	IsVnetManaged  bool
}

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *NatGatewaySpec) ResourceRef() *asonetworkv1.NatGateway {
	return &asonetworkv1.NatGateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
	}
}

// Parameters implements azure.ASOResourceSpecGetter.
func (s *NatGatewaySpec) Parameters(ctx context.Context, existingNatGateway *asonetworkv1.NatGateway) (params *asonetworkv1.NatGateway, err error) {
	natGateway := &asonetworkv1.NatGateway{}
	natGateway.Spec = asonetworkv1.NatGateway_Spec{}

	if existingNatGateway != nil {
		natGateway = existingNatGateway
	}

	natGateway.Spec.AzureName = s.Name
	natGateway.Spec.Owner = &genruntime.KnownResourceReference{
		Name: s.ResourceGroup,
	}
	natGateway.Spec.Location = ptr.To(s.Location)
	natGateway.Spec.Sku = &asonetworkv1.NatGatewaySku{
		Name: ptr.To(asonetworkv1.NatGatewaySku_Name_Standard),
	}
	natGateway.Spec.PublicIpAddresses = []asonetworkv1.ApplicationGatewaySubResource{
		{
			Reference: &genruntime.ResourceReference{
				ARMID: azure.PublicIPID(s.SubscriptionID, s.ResourceGroup, s.NatGatewayIP.Name),
			},
		},
	}
	natGateway.Spec.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.ClusterName,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        ptr.To(s.Name),
		Additional:  s.AdditionalTags,
	})

	return natGateway, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *NatGatewaySpec) WasManaged(resource *asonetworkv1.NatGateway) bool {
	return s.IsVnetManaged
}
