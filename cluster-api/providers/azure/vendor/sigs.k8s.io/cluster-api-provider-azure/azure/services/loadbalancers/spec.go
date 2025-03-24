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

package loadbalancers

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// LBSpec defines the specification for a Load Balancer.
type LBSpec struct {
	Name                 string
	ResourceGroup        string
	SubscriptionID       string
	ClusterName          string
	Location             string
	ExtendedLocation     *infrav1.ExtendedLocationSpec
	Role                 string
	Type                 infrav1.LBType
	SKU                  infrav1.SKU
	VNetName             string
	VNetResourceGroup    string
	SubnetName           string
	BackendPoolName      string
	FrontendIPConfigs    []infrav1.FrontendIP
	APIServerPort        int32
	IdleTimeoutInMinutes *int32
	AdditionalTags       map[string]string
}

// ResourceName returns the name of the load balancer.
func (s *LBSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *LBSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for load balancers.
func (s *LBSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the load balancer.
func (s *LBSpec) Parameters(_ context.Context, existing interface{}) (parameters interface{}, err error) {
	var (
		etag                *string
		frontendIDs         []*armnetwork.SubResource
		frontendIPConfigs   []*armnetwork.FrontendIPConfiguration
		loadBalancingRules  []*armnetwork.LoadBalancingRule
		backendAddressPools []*armnetwork.BackendAddressPool
		outboundRules       []*armnetwork.OutboundRule
		probes              []*armnetwork.Probe
	)

	if existing != nil {
		existingLB, ok := existing.(armnetwork.LoadBalancer)
		if !ok {
			return nil, errors.Errorf("%T is not an armnetwork.LoadBalancer", existing)
		}
		// LB already exists
		// We append the existing LB etag to the header to ensure we only apply the updates if the LB has not been modified.
		etag = existingLB.Etag
		update := false

		// merge existing LB properties with desired properties
		frontendIPConfigs = existingLB.Properties.FrontendIPConfigurations
		wantedIPs, wantedFrontendIDs := getFrontendIPConfigs(*s)
		for _, ip := range wantedIPs {
			if !ipExists(frontendIPConfigs, *ip) {
				update = true
				frontendIPConfigs = append(frontendIPConfigs, ip)
			}
		}

		loadBalancingRules = existingLB.Properties.LoadBalancingRules
		for _, rule := range getLoadBalancingRules(*s, wantedFrontendIDs) {
			if !lbRuleExists(loadBalancingRules, *rule) {
				update = true
				loadBalancingRules = append(loadBalancingRules, rule)
			}
		}

		backendAddressPools = existingLB.Properties.BackendAddressPools
		for _, pool := range getBackendAddressPools(*s) {
			if !poolExists(backendAddressPools, *pool) {
				update = true
				backendAddressPools = append(backendAddressPools, pool)
			}
		}

		outboundRules = existingLB.Properties.OutboundRules
		for _, rule := range getOutboundRules(*s, wantedFrontendIDs) {
			if !outboundRuleExists(outboundRules, *rule) {
				update = true
				outboundRules = append(outboundRules, rule)
			}
		}

		probes = existingLB.Properties.Probes
		for _, probe := range getProbes(*s) {
			if !probeExists(probes, *probe) {
				update = true
				probes = append(probes, probe)
			}
		}

		if !update {
			// load balancer already exists with all required defaults
			return nil, nil
		}
	} else {
		frontendIPConfigs, frontendIDs = getFrontendIPConfigs(*s)
		loadBalancingRules = getLoadBalancingRules(*s, frontendIDs)
		backendAddressPools = getBackendAddressPools(*s)
		outboundRules = getOutboundRules(*s, frontendIDs)
		probes = getProbes(*s)
	}

	lb := armnetwork.LoadBalancer{
		Etag:             etag,
		SKU:              &armnetwork.LoadBalancerSKU{Name: ptr.To(converters.SKUtoSDK(s.SKU))},
		Location:         ptr.To(s.Location),
		ExtendedLocation: converters.ExtendedLocationToNetworkSDK(s.ExtendedLocation),
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Role:        ptr.To(s.Role),
			Additional:  s.AdditionalTags,
		})),
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: frontendIPConfigs,
			BackendAddressPools:      backendAddressPools,
			OutboundRules:            outboundRules,
			Probes:                   probes,
			LoadBalancingRules:       loadBalancingRules,
		},
	}

	return lb, nil
}

func getFrontendIPConfigs(lbSpec LBSpec) ([]*armnetwork.FrontendIPConfiguration, []*armnetwork.SubResource) {
	frontendIPConfigurations := make([]*armnetwork.FrontendIPConfiguration, 0)
	frontendIDs := make([]*armnetwork.SubResource, 0)
	for _, ipConfig := range lbSpec.FrontendIPConfigs {
		var properties armnetwork.FrontendIPConfigurationPropertiesFormat
		if lbSpec.Type == infrav1.Internal {
			properties = armnetwork.FrontendIPConfigurationPropertiesFormat{
				PrivateIPAllocationMethod: ptr.To(armnetwork.IPAllocationMethodStatic),
				Subnet: &armnetwork.Subnet{
					ID: ptr.To(azure.SubnetID(lbSpec.SubscriptionID, lbSpec.VNetResourceGroup, lbSpec.VNetName, lbSpec.SubnetName)),
				},
				PrivateIPAddress: ptr.To(ipConfig.PrivateIPAddress),
			}
		} else {
			properties = armnetwork.FrontendIPConfigurationPropertiesFormat{
				PublicIPAddress: &armnetwork.PublicIPAddress{
					ID: ptr.To(azure.PublicIPID(lbSpec.SubscriptionID, lbSpec.ResourceGroup, ipConfig.PublicIP.Name)),
				},
			}
		}
		frontendIPConfigurations = append(frontendIPConfigurations, &armnetwork.FrontendIPConfiguration{
			Properties: &properties,
			Name:       ptr.To(ipConfig.Name),
		})
		frontendIDs = append(frontendIDs, &armnetwork.SubResource{
			ID: ptr.To(azure.FrontendIPConfigID(lbSpec.SubscriptionID, lbSpec.ResourceGroup, lbSpec.Name, ipConfig.Name)),
		})
	}
	return frontendIPConfigurations, frontendIDs
}

func getOutboundRules(lbSpec LBSpec, frontendIDs []*armnetwork.SubResource) []*armnetwork.OutboundRule {
	if lbSpec.Type == infrav1.Internal {
		return []*armnetwork.OutboundRule{}
	}
	return []*armnetwork.OutboundRule{
		{
			Name: ptr.To(outboundNAT),
			Properties: &armnetwork.OutboundRulePropertiesFormat{
				Protocol:                 ptr.To(armnetwork.LoadBalancerOutboundRuleProtocolAll),
				IdleTimeoutInMinutes:     lbSpec.IdleTimeoutInMinutes,
				FrontendIPConfigurations: frontendIDs,
				BackendAddressPool: &armnetwork.SubResource{
					ID: ptr.To(azure.AddressPoolID(lbSpec.SubscriptionID, lbSpec.ResourceGroup, lbSpec.Name, lbSpec.BackendPoolName)),
				},
			},
		},
	}
}

func getLoadBalancingRules(lbSpec LBSpec, frontendIDs []*armnetwork.SubResource) []*armnetwork.LoadBalancingRule {
	if lbSpec.Role == infrav1.APIServerRole || lbSpec.Role == infrav1.APIServerRoleInternal {
		// We disable outbound SNAT explicitly in the HTTPS LB rule and enable TCP and UDP outbound NAT with an outbound rule.
		// For more information on Standard LB outbound connections see https://learn.microsoft.com/azure/load-balancer/load-balancer-outbound-connections.
		var frontendIPConfig *armnetwork.SubResource
		if len(frontendIDs) != 0 {
			frontendIPConfig = frontendIDs[0]
		}
		return []*armnetwork.LoadBalancingRule{
			{
				Name: ptr.To(lbRuleHTTPS),
				Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
					DisableOutboundSnat:     ptr.To(true),
					Protocol:                ptr.To(armnetwork.TransportProtocolTCP),
					FrontendPort:            ptr.To[int32](lbSpec.APIServerPort),
					BackendPort:             ptr.To[int32](lbSpec.APIServerPort),
					IdleTimeoutInMinutes:    lbSpec.IdleTimeoutInMinutes,
					EnableFloatingIP:        ptr.To(false),
					LoadDistribution:        ptr.To(armnetwork.LoadDistributionDefault),
					FrontendIPConfiguration: frontendIPConfig,
					BackendAddressPool: &armnetwork.SubResource{
						ID: ptr.To(azure.AddressPoolID(lbSpec.SubscriptionID, lbSpec.ResourceGroup, lbSpec.Name, lbSpec.BackendPoolName)),
					},
					Probe: &armnetwork.SubResource{
						ID: ptr.To(azure.ProbeID(lbSpec.SubscriptionID, lbSpec.ResourceGroup, lbSpec.Name, httpsProbe)),
					},
				},
			},
		}
	}
	return []*armnetwork.LoadBalancingRule{}
}

func getBackendAddressPools(lbSpec LBSpec) []*armnetwork.BackendAddressPool {
	return []*armnetwork.BackendAddressPool{
		{
			Name: ptr.To(lbSpec.BackendPoolName),
		},
	}
}

func getProbes(lbSpec LBSpec) []*armnetwork.Probe {
	if lbSpec.Role == infrav1.APIServerRole || lbSpec.Role == infrav1.APIServerRoleInternal {
		return []*armnetwork.Probe{
			{
				Name: ptr.To(httpsProbe),
				Properties: &armnetwork.ProbePropertiesFormat{
					Protocol:          ptr.To(armnetwork.ProbeProtocolHTTPS),
					Port:              ptr.To[int32](lbSpec.APIServerPort),
					RequestPath:       ptr.To(httpsProbeRequestPath),
					IntervalInSeconds: ptr.To[int32](15),
					NumberOfProbes:    ptr.To[int32](4),
				},
			},
		}
	}
	return []*armnetwork.Probe{}
}

func probeExists(probes []*armnetwork.Probe, probe armnetwork.Probe) bool {
	for _, p := range probes {
		if ptr.Deref(p.Name, "") == ptr.Deref(probe.Name, "") {
			return true
		}
	}
	return false
}

func outboundRuleExists(rules []*armnetwork.OutboundRule, rule armnetwork.OutboundRule) bool {
	for _, r := range rules {
		if ptr.Deref(r.Name, "") == ptr.Deref(rule.Name, "") {
			return true
		}
	}
	return false
}

func poolExists(pools []*armnetwork.BackendAddressPool, pool armnetwork.BackendAddressPool) bool {
	for _, p := range pools {
		if ptr.Deref(p.Name, "") == ptr.Deref(pool.Name, "") {
			return true
		}
	}
	return false
}

func lbRuleExists(rules []*armnetwork.LoadBalancingRule, rule armnetwork.LoadBalancingRule) bool {
	for _, r := range rules {
		if ptr.Deref(r.Name, "") == ptr.Deref(rule.Name, "") {
			return true
		}
	}
	return false
}

func ipExists(configs []*armnetwork.FrontendIPConfiguration, config armnetwork.FrontendIPConfiguration) bool {
	for _, ip := range configs {
		if ptr.Deref(ip.Name, "") == ptr.Deref(config.Name, "") {
			return true
		}
	}
	return false
}
