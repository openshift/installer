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
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
)

// LBSpec defines the specification for a Load Balancer.
type LBSpec struct {
	Name                  string
	ResourceGroup         string
	SubscriptionID        string
	ClusterName           string
	Location              string
	ExtendedLocation      *infrav1.ExtendedLocationSpec
	Role                  string
	Type                  infrav1.LBType
	SKU                   infrav1.SKU
	VNetName              string
	VNetResourceGroup     string
	SubnetName            string
	BackendPoolName       string
	FrontendIPConfigs     []infrav1.FrontendIP
	APIServerPort         int32
	APIServerFrontendPort int32
	IdleTimeoutInMinutes  *int32
	AdditionalTags        map[string]string
	AdditionalPorts       []infrav1.LoadBalancerPort
	AvailabilityZones     []string
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
func (s *LBSpec) Parameters(_ context.Context, existing any) (parameters any, err error) {
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
		if mergedRules, changed := mergeLoadBalancingRules(loadBalancingRules, getLoadBalancingRules(*s, wantedFrontendIDs)); changed {
			update = true
			loadBalancingRules = mergedRules
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
		if mergedProbes, changed := mergeProbes(probes, getProbes(*s)); changed {
			update = true
			probes = mergedProbes
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

	// Convert availability zones to []*string for Azure SDK.
	// IMPORTANT: Zones can only be set on frontend IP configurations for internal load balancers
	// (where the frontend references a subnet). For public load balancers, zone-redundancy is
	// achieved by setting zones on the associated public IP address resource, NOT on the load
	// balancer's frontend IP configuration.
	//
	// Azure returns error "LoadBalancerFrontendIPConfigCannotHaveZoneWhenReferencingPublicIPAddress"
	// if zones are specified on a frontend that references a public IP.
	//
	// See: https://learn.microsoft.com/en-us/azure/reliability/reliability-load-balancer#zone-redundant-load-balancer
	// Section: "Zone-redundant load balancer" - "For public load balancers, if the public IP in the
	// Load balancer's frontend is zone redundant then the load balancer is also zone-redundant."
	var zones []*string
	if len(lbSpec.AvailabilityZones) > 0 && lbSpec.Type == infrav1.Internal {
		zones = make([]*string, len(lbSpec.AvailabilityZones))
		for i, zone := range lbSpec.AvailabilityZones {
			zones[i] = ptr.To(zone)
		}
	}

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
			Zones:      zones,
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
		rules := []*armnetwork.LoadBalancingRule{
			{
				Name: ptr.To(lbRuleHTTPS),
				Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
					DisableOutboundSnat:     ptr.To(true),
					Protocol:                ptr.To(armnetwork.TransportProtocolTCP),
					FrontendPort:            ptr.To(lbSpec.apiServerFrontendPort()),
					BackendPort:             ptr.To(lbSpec.APIServerPort),
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

		for _, port := range lbSpec.AdditionalPorts {
			rules = append(rules, &armnetwork.LoadBalancingRule{
				Name: ptr.To(port.Name),
				Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
					DisableOutboundSnat:     ptr.To(true),
					Protocol:                ptr.To(armnetwork.TransportProtocolTCP),
					FrontendPort:            ptr.To(port.Port),
					BackendPort:             ptr.To(port.Port),
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
			})
		}

		return rules
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
					ProbeThreshold:    ptr.To[int32](1),
				},
			},
		}
	}
	return []*armnetwork.Probe{}
}

func (s LBSpec) apiServerFrontendPort() int32 {
	if s.APIServerFrontendPort != 0 {
		return s.APIServerFrontendPort
	}
	return s.APIServerPort
}

func mergeProbes(existing, desired []*armnetwork.Probe) ([]*armnetwork.Probe, bool) {
	merged := append([]*armnetwork.Probe(nil), existing...)
	changed := false
	for _, desiredProbe := range desired {
		found := false
		for i, existingProbe := range merged {
			if !strings.EqualFold(ptr.Deref(existingProbe.Name, ""), ptr.Deref(desiredProbe.Name, "")) {
				continue
			}
			found = true
			if !probeMatches(existingProbe, desiredProbe) {
				merged[i] = desiredProbe
				changed = true
			}
			break
		}
		if !found {
			merged = append(merged, desiredProbe)
			changed = true
		}
	}
	return merged, changed
}

func probeMatches(existing, desired *armnetwork.Probe) bool {
	if existing == nil || desired == nil || existing.Properties == nil || desired.Properties == nil {
		return existing == desired
	}
	return strings.EqualFold(ptr.Deref(existing.Name, ""), ptr.Deref(desired.Name, "")) &&
		pointerMatches(existing.Properties.Protocol, desired.Properties.Protocol) &&
		pointerMatches(existing.Properties.Port, desired.Properties.Port) &&
		pointerMatches(existing.Properties.RequestPath, desired.Properties.RequestPath) &&
		pointerMatches(existing.Properties.IntervalInSeconds, desired.Properties.IntervalInSeconds) &&
		pointerMatches(existing.Properties.ProbeThreshold, desired.Properties.ProbeThreshold)
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

func mergeLoadBalancingRules(existing, desired []*armnetwork.LoadBalancingRule) ([]*armnetwork.LoadBalancingRule, bool) {
	merged := append([]*armnetwork.LoadBalancingRule(nil), existing...)
	changed := false
	for _, desiredRule := range desired {
		found := false
		for i, existingRule := range merged {
			if !strings.EqualFold(ptr.Deref(existingRule.Name, ""), ptr.Deref(desiredRule.Name, "")) {
				continue
			}
			found = true
			if !loadBalancingRuleMatches(existingRule, desiredRule) {
				merged[i] = desiredRule
				changed = true
			}
			break
		}
		if !found {
			merged = append(merged, desiredRule)
			changed = true
		}
	}
	return merged, changed
}

func loadBalancingRuleMatches(existing, desired *armnetwork.LoadBalancingRule) bool {
	if existing == nil || desired == nil || existing.Properties == nil || desired.Properties == nil {
		return existing == desired
	}
	return strings.EqualFold(ptr.Deref(existing.Name, ""), ptr.Deref(desired.Name, "")) &&
		pointerMatches(existing.Properties.DisableOutboundSnat, desired.Properties.DisableOutboundSnat) &&
		pointerMatches(existing.Properties.Protocol, desired.Properties.Protocol) &&
		pointerMatches(existing.Properties.FrontendPort, desired.Properties.FrontendPort) &&
		pointerMatches(existing.Properties.BackendPort, desired.Properties.BackendPort) &&
		pointerMatches(existing.Properties.IdleTimeoutInMinutes, desired.Properties.IdleTimeoutInMinutes) &&
		pointerMatches(existing.Properties.EnableFloatingIP, desired.Properties.EnableFloatingIP) &&
		pointerMatches(existing.Properties.LoadDistribution, desired.Properties.LoadDistribution) &&
		subResourceMatches(existing.Properties.FrontendIPConfiguration, desired.Properties.FrontendIPConfiguration) &&
		subResourceMatches(existing.Properties.BackendAddressPool, desired.Properties.BackendAddressPool) &&
		subResourceMatches(existing.Properties.Probe, desired.Properties.Probe)
}

func pointerMatches[T comparable](existing, desired *T) bool {
	return desired == nil || existing != nil && *existing == *desired
}

func subResourceMatches(existing, desired *armnetwork.SubResource) bool {
	return desired == nil || existing != nil && strings.EqualFold(ptr.Deref(existing.ID, ""), ptr.Deref(desired.ID, ""))
}

func ipExists(configs []*armnetwork.FrontendIPConfiguration, config armnetwork.FrontendIPConfiguration) bool {
	for _, ip := range configs {
		if ptr.Deref(ip.Name, "") == ptr.Deref(config.Name, "") {
			return true
		}
	}
	return false
}
