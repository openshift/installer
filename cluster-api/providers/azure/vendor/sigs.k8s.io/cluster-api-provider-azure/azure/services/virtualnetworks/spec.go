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

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
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

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *VNetSpec) ResourceRef() *asonetworkv1.VirtualNetwork {
	return &asonetworkv1.VirtualNetwork{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Name,
		},
	}
}

// Parameters implements azure.ASOResourceSpecGetter.
func (s *VNetSpec) Parameters(ctx context.Context, existing *asonetworkv1.VirtualNetwork) (*asonetworkv1.VirtualNetwork, error) {
	vnet := existing
	if existing == nil {
		vnet = &asonetworkv1.VirtualNetwork{
			Spec: asonetworkv1.VirtualNetwork_Spec{
				Tags: infrav1.Build(infrav1.BuildParams{
					ClusterName: s.ClusterName,
					Lifecycle:   infrav1.ResourceLifecycleOwned,
					Name:        ptr.To(s.Name),
					Role:        ptr.To(infrav1.CommonRole),
					Additional:  s.AdditionalTags,
				}),
			},
		}
	}

	vnet.Spec.AzureName = s.Name
	vnet.Spec.Owner = &genruntime.KnownResourceReference{
		Name: s.ResourceGroup,
	}
	vnet.Spec.Location = ptr.To(s.Location)
	vnet.Spec.ExtendedLocation = converters.ExtendedLocationToNetworkASO(s.ExtendedLocation)
	vnet.Spec.AddressSpace = &asonetworkv1.AddressSpace{
		AddressPrefixes: s.CIDRs,
	}

	return vnet, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *VNetSpec) WasManaged(resource *asonetworkv1.VirtualNetwork) bool {
	return infrav1.Tags(resource.Status.Tags).HasOwned(s.ClusterName)
}
