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

package privateendpoints

import (
	"context"
	"sort"

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// PrivateLinkServiceConnection defines the specification for a private link service connection associated with a private endpoint.
type PrivateLinkServiceConnection struct {
	Name                 string
	PrivateLinkServiceID string
	GroupIDs             []string
	RequestMessage       string
}

// PrivateEndpointSpec defines the specification for a private endpoint.
type PrivateEndpointSpec struct {
	Name                          string
	ResourceGroup                 string
	Location                      string
	CustomNetworkInterfaceName    string
	PrivateIPAddresses            []string
	SubnetID                      string
	ApplicationSecurityGroups     []string
	ManualApproval                bool
	PrivateLinkServiceConnections []PrivateLinkServiceConnection
	AdditionalTags                infrav1.Tags
	ClusterName                   string
}

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *PrivateEndpointSpec) ResourceRef() *asonetworkv1.PrivateEndpoint {
	return &asonetworkv1.PrivateEndpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name: azure.GetNormalizedKubernetesName(s.Name),
		},
	}
}

// Parameters implements azure.ASOResourceSpecGetter.
func (s *PrivateEndpointSpec) Parameters(_ context.Context, existingPrivateEndpoint *asonetworkv1.PrivateEndpoint) (*asonetworkv1.PrivateEndpoint, error) {
	privateEndpoint := &asonetworkv1.PrivateEndpoint{}
	if existingPrivateEndpoint != nil {
		privateEndpoint = existingPrivateEndpoint
	}

	if len(s.ApplicationSecurityGroups) > 0 {
		applicationSecurityGroups := make([]asonetworkv1.ApplicationSecurityGroupSpec_PrivateEndpoint_SubResourceEmbedded, 0, len(s.ApplicationSecurityGroups))
		for _, applicationSecurityGroup := range s.ApplicationSecurityGroups {
			applicationSecurityGroups = append(applicationSecurityGroups, asonetworkv1.ApplicationSecurityGroupSpec_PrivateEndpoint_SubResourceEmbedded{
				Reference: &genruntime.ResourceReference{
					ARMID: applicationSecurityGroup,
				},
			})
		}
		// Sort the slices in order to get the same order of elements for both new and existing application Security Groups.
		sort.SliceStable(applicationSecurityGroups, func(i, j int) bool {
			return applicationSecurityGroups[i].Reference.ARMID < applicationSecurityGroups[j].Reference.ARMID
		})
		privateEndpoint.Spec.ApplicationSecurityGroups = applicationSecurityGroups
	}

	privateEndpoint.Spec.AzureName = s.Name
	privateEndpoint.Spec.CustomNetworkInterfaceName = ptr.To(s.CustomNetworkInterfaceName)
	privateEndpoint.Spec.Location = ptr.To(s.Location)

	if len(s.PrivateIPAddresses) > 0 {
		privateIPAddresses := make([]asonetworkv1.PrivateEndpointIPConfiguration, 0, len(s.PrivateIPAddresses))
		for _, address := range s.PrivateIPAddresses {
			ipConfig := asonetworkv1.PrivateEndpointIPConfiguration{PrivateIPAddress: ptr.To(address)}
			privateIPAddresses = append(privateIPAddresses, ipConfig)
		}
		sort.SliceStable(privateIPAddresses, func(i, j int) bool {
			return *privateIPAddresses[i].PrivateIPAddress < *privateIPAddresses[j].PrivateIPAddress
		})
		privateEndpoint.Spec.IpConfigurations = privateIPAddresses
	}

	if len(s.PrivateLinkServiceConnections) > 0 {
		privateLinkServiceConnections := make([]asonetworkv1.PrivateLinkServiceConnection, 0, len(s.PrivateLinkServiceConnections))
		for _, privateLinkServiceConnection := range s.PrivateLinkServiceConnections {
			linkServiceConnection := asonetworkv1.PrivateLinkServiceConnection{
				Name: ptr.To(privateLinkServiceConnection.Name),
				PrivateLinkServiceReference: &genruntime.ResourceReference{
					ARMID: privateLinkServiceConnection.PrivateLinkServiceID,
				},
			}

			if len(privateLinkServiceConnection.GroupIDs) > 0 {
				linkServiceConnection.GroupIds = privateLinkServiceConnection.GroupIDs
			}

			if privateLinkServiceConnection.RequestMessage != "" {
				linkServiceConnection.RequestMessage = ptr.To(privateLinkServiceConnection.RequestMessage)
			}
			privateLinkServiceConnections = append(privateLinkServiceConnections, linkServiceConnection)
		}
		sort.SliceStable(privateLinkServiceConnections, func(i, j int) bool {
			return *privateLinkServiceConnections[i].Name < *privateLinkServiceConnections[j].Name
		})

		if s.ManualApproval {
			privateEndpoint.Spec.ManualPrivateLinkServiceConnections = privateLinkServiceConnections
		} else {
			privateEndpoint.Spec.PrivateLinkServiceConnections = privateLinkServiceConnections
		}
	}

	privateEndpoint.Spec.Owner = &genruntime.KnownResourceReference{
		Name: azure.GetNormalizedKubernetesName(s.ResourceGroup),
	}

	privateEndpoint.Spec.Subnet = &asonetworkv1.Subnet_PrivateEndpoint_SubResourceEmbedded{
		Reference: &genruntime.ResourceReference{
			ARMID: s.SubnetID,
		},
	}

	privateEndpoint.Spec.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.ClusterName,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        ptr.To(s.Name),
		Additional:  s.AdditionalTags,
	})

	return privateEndpoint, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
// It always returns true since CAPZ doesn't support BYO private endpoints.
func (s *PrivateEndpointSpec) WasManaged(_ *asonetworkv1.PrivateEndpoint) bool {
	return true
}
