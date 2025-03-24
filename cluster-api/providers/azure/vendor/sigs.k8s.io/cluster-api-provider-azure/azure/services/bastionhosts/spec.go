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

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
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

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *AzureBastionSpec) ResourceRef() *asonetworkv1.BastionHost {
	return &asonetworkv1.BastionHost{
		ObjectMeta: metav1.ObjectMeta{
			Name: azure.GetNormalizedKubernetesName(s.Name),
		},
	}
}

// Parameters implements azure.ASOResourceSpecGetter.
func (s *AzureBastionSpec) Parameters(_ context.Context, existingBastionHost *asonetworkv1.BastionHost) (parameters *asonetworkv1.BastionHost, err error) {
	bastionHost := &asonetworkv1.BastionHost{}
	if existingBastionHost != nil {
		bastionHost = existingBastionHost
	}

	bastionHostIPConfigName := fmt.Sprintf("%s-%s", s.Name, "bastionIP")
	bastionHost.Spec.AzureName = s.Name
	bastionHost.Spec.Location = ptr.To(s.Location)
	bastionHost.Spec.Owner = &genruntime.KnownResourceReference{
		Name: azure.GetNormalizedKubernetesName(s.ResourceGroup),
	}
	bastionHost.Spec.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.ClusterName,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        ptr.To(s.Name),
		Role:        ptr.To("Bastion"),
	})
	bastionHost.Spec.Sku = &asonetworkv1.Sku{
		Name: ptr.To(asonetworkv1.Sku_Name(s.Sku)),
	}
	bastionHost.Spec.EnableTunneling = ptr.To(s.EnableTunneling)
	bastionHost.Spec.DnsName = ptr.To(fmt.Sprintf("%s-bastion", strings.ToLower(s.Name)))
	bastionHost.Spec.IpConfigurations = []asonetworkv1.BastionHostIPConfiguration{
		{
			Name: ptr.To(bastionHostIPConfigName),
			Subnet: &asonetworkv1.BastionHostSubResource{
				Reference: &genruntime.ResourceReference{
					ARMID: s.SubnetID,
				},
			},
			PublicIPAddress: &asonetworkv1.BastionHostSubResource{
				Reference: &genruntime.ResourceReference{
					ARMID: s.PublicIPID,
				},
			},
			PrivateIPAllocationMethod: ptr.To(asonetworkv1.IPAllocationMethod_Dynamic),
		},
	}

	return bastionHost, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *AzureBastionSpec) WasManaged(_ *asonetworkv1.BastionHost) bool {
	// returns always returns true as CAPZ does not support BYO bastion.
	return true
}
