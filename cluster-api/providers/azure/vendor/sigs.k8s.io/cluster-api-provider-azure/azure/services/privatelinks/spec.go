/*
Copyright The Kubernetes Authors.

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

package privatelinks

import (
	"context"
	"fmt"

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// PrivateLinkSpec defines the specification for a private link service.
type PrivateLinkSpec struct {
	Name                      string
	ResourceGroup             string
	SubscriptionID            string
	Location                  string
	VNetResourceGroup         string
	VNet                      string
	NATIPConfiguration        []NATIPConfiguration
	LoadBalancerName          string
	LBFrontendIPConfigNames   []string
	AllowedSubscriptions      []string
	AutoApprovedSubscriptions []string
	EnableProxyProtocol       *bool
	ClusterName               string
	AdditionalTags            infrav1.Tags
}

// NATIPConfiguration defines the NAT IP configuration for the private link service.
type NATIPConfiguration struct {
	// AllocationMethod can be Static or Dynamic.
	AllocationMethod infrav1.PrivateLinkNATIPAllocationMethod

	// Subnet from the VNet from which the IP is allocated.
	Subnet string

	// PrivateIPAddress is the optional static private IP address from the specified Subnet.
	PrivateIPAddress string
}

// ResourceRef implements [azure.ASOResourceSpecGetter.ResourceRef].
func (s *PrivateLinkSpec) ResourceRef() *asonetworkv1.PrivateLinkService {
	return &asonetworkv1.PrivateLinkService{
		ObjectMeta: metav1.ObjectMeta{
			Name: azure.GetNormalizedKubernetesName(s.Name),
		},
	}
}

// Parameters implements [azure.ASOResourceSpecGetter.Parameters].
func (s *PrivateLinkSpec) Parameters(_ context.Context, existingPrivateLink *asonetworkv1.PrivateLinkService) (*asonetworkv1.PrivateLinkService, error) {
	privateLink := existingPrivateLink
	if privateLink == nil {
		privateLink = new(asonetworkv1.PrivateLinkService)
	}

	privateLink.Spec.AzureName = s.Name
	privateLink.Spec.Location = ptr.To(s.Location)
	privateLink.Spec.EnableProxyProtocol = s.EnableProxyProtocol

	privateLink.Spec.IpConfigurations = make([]asonetworkv1.PrivateLinkServiceIpConfiguration, len(s.NATIPConfiguration))
	for i, natIPConfiguration := range s.NATIPConfiguration {
		privateLink.Spec.IpConfigurations[i] = asonetworkv1.PrivateLinkServiceIpConfiguration{
			Name: ptr.To(fmt.Sprintf("%s-natipconfig-%d", natIPConfiguration.Subnet, i+1)),
			Subnet: &asonetworkv1.Subnet_PrivateLinkService_SubResourceEmbedded{
				Reference: &genruntime.ResourceReference{
					ARMID: azure.SubnetID(s.SubscriptionID, s.VNetResourceGroup, s.VNet, natIPConfiguration.Subnet),
				},
			},
			PrivateIPAllocationMethod: ptr.To(asonetworkv1.IPAllocationMethod(natIPConfiguration.AllocationMethod)),
		}

		if *privateLink.Spec.IpConfigurations[i].PrivateIPAllocationMethod == asonetworkv1.IPAllocationMethod_Static {
			privateLink.Spec.IpConfigurations[i].PrivateIPAddress = ptr.To(natIPConfiguration.PrivateIPAddress)
		}
	}
	if len(privateLink.Spec.IpConfigurations) == 0 {
		return nil, fmt.Errorf("private link %s has no NAT IP configurations", s.Name)
	}
	privateLink.Spec.IpConfigurations[0].Primary = ptr.To(true)

	privateLink.Spec.LoadBalancerFrontendIpConfigurations = make([]asonetworkv1.FrontendIPConfiguration_PrivateLinkService_SubResourceEmbedded, len(s.LBFrontendIPConfigNames))
	for i, frontendIPConfigName := range s.LBFrontendIPConfigNames {
		privateLink.Spec.LoadBalancerFrontendIpConfigurations[i] = asonetworkv1.FrontendIPConfiguration_PrivateLinkService_SubResourceEmbedded{
			Reference: &genruntime.ResourceReference{
				ARMID: azure.FrontendIPConfigID(s.SubscriptionID, s.ResourceGroup, s.LoadBalancerName, frontendIPConfigName),
			},
		}
	}

	privateLink.Spec.Visibility = nil
	if len(s.AllowedSubscriptions) > 0 {
		privateLink.Spec.Visibility = &asonetworkv1.ResourceSet{
			Subscriptions: s.AllowedSubscriptions,
		}
	}

	privateLink.Spec.AutoApproval = nil
	if len(s.AutoApprovedSubscriptions) > 0 {
		privateLink.Spec.AutoApproval = &asonetworkv1.ResourceSet{
			Subscriptions: s.AutoApprovedSubscriptions,
		}
	}

	privateLink.Spec.Owner = &genruntime.KnownResourceReference{
		Name: azure.GetNormalizedKubernetesName(s.ResourceGroup),
	}

	privateLink.Spec.Tags = infrav1.Build(infrav1.BuildParams{
		ClusterName: s.ClusterName,
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        ptr.To(s.Name),
		Additional:  s.AdditionalTags,
	})

	return privateLink, nil
}

// WasManaged implements [azure.ASOResourceSpecGetter.WasManaged].
// It always returns true since CAPZ doesn't support BYO private links.
func (s *PrivateLinkSpec) WasManaged(_ *asonetworkv1.PrivateLinkService) bool {
	return true
}
