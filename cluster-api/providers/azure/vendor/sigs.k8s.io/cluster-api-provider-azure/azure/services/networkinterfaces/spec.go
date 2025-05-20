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

package networkinterfaces

import (
	"context"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// NICSpec defines the specification for a Network Interface.
type NICSpec struct {
	Name                      string
	ResourceGroup             string
	Location                  string
	ExtendedLocation          *infrav1.ExtendedLocationSpec
	SubscriptionID            string
	MachineName               string
	SubnetName                string
	VNetName                  string
	VNetResourceGroup         string
	StaticIPAddress           string
	PublicLBName              string
	PublicLBAddressPoolName   string
	PublicLBNATRuleName       string
	InternalLBName            string
	InternalLBAddressPoolName string
	PublicIPName              string
	AcceleratedNetworking     *bool
	IPv6Enabled               bool
	EnableIPForwarding        bool
	SKU                       *resourceskus.SKU
	DNSServers                []string
	AdditionalTags            infrav1.Tags
	ClusterName               string
	IPConfigs                 []IPConfig
}

// IPConfig defines the specification for an IP address configuration.
type IPConfig struct {
	PrivateIP       *string
	PublicIPAddress *string
}

// ResourceName returns the name of the network interface.
func (s *NICSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *NICSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for network interfaces.
func (s *NICSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the network interface.
func (s *NICSpec) Parameters(ctx context.Context, existing interface{}) (parameters interface{}, err error) {
	_, log, done := tele.StartSpanWithLogger(ctx, "networkinterfaces.NICSpec.Parameters")
	defer done()

	if existing != nil {
		existingNIC, ok := existing.(armnetwork.Interface)
		if !ok {
			return nil, errors.Errorf("%T is not an armnetwork.Interface", existing)
		}

		// If the NIC already exists, return a nil parameters and nil errors to skip the create or update of the already existing NIC.
		// If the NIC is in ProvisioningFailed state, we will try to build the parameters of the existing NIC as the provisioning failed does not necessarily mean the NIC is in a bad state.
		// Reference: https://learn.microsoft.com/en-us/azure/networking/troubleshoot-failed-state#provisioning-states and
		// https://github.com/kubernetes-sigs/cluster-api-provider-azure/issues/5515
		if existingNIC.Properties != nil && existingNIC.Properties.ProvisioningState != nil && *existingNIC.Properties.ProvisioningState != armnetwork.ProvisioningStateFailed {
			// Return nil for both parameters and error as no changes are needed for the existing resource
			// otherwise rebuild the parameters of the existing NIC so that we can patch the ProvisioningState
			log.V(4).Info("existing NIC is not in ProvisioningFailed state, returning nil parameters and nil error", "ProvisioningState", *existingNIC.Properties.ProvisioningState)
			return nil, nil
		}
	}

	primaryIPConfig := &armnetwork.InterfaceIPConfigurationPropertiesFormat{
		Primary: ptr.To(true),
	}

	subnet := &armnetwork.Subnet{
		ID: ptr.To(azure.SubnetID(s.SubscriptionID, s.VNetResourceGroup, s.VNetName, s.SubnetName)),
	}
	primaryIPConfig.Subnet = subnet

	primaryIPConfig.PrivateIPAllocationMethod = ptr.To(armnetwork.IPAllocationMethodDynamic)
	if s.StaticIPAddress != "" {
		primaryIPConfig.PrivateIPAllocationMethod = ptr.To(armnetwork.IPAllocationMethodStatic)
		primaryIPConfig.PrivateIPAddress = ptr.To(s.StaticIPAddress)
	}

	backendAddressPools := []*armnetwork.BackendAddressPool{}
	if s.PublicLBName != "" {
		if s.PublicLBAddressPoolName != "" {
			backendAddressPools = append(backendAddressPools,
				&armnetwork.BackendAddressPool{
					ID: ptr.To(azure.AddressPoolID(s.SubscriptionID, s.ResourceGroup, s.PublicLBName, s.PublicLBAddressPoolName)),
				})
		}
		if s.PublicLBNATRuleName != "" {
			primaryIPConfig.LoadBalancerInboundNatRules = []*armnetwork.InboundNatRule{
				{
					ID: ptr.To(azure.NATRuleID(s.SubscriptionID, s.ResourceGroup, s.PublicLBName, s.PublicLBNATRuleName)),
				},
			}
		}
	}
	if s.InternalLBName != "" && s.InternalLBAddressPoolName != "" {
		backendAddressPools = append(backendAddressPools,
			&armnetwork.BackendAddressPool{
				ID: ptr.To(azure.AddressPoolID(s.SubscriptionID, s.ResourceGroup, s.InternalLBName, s.InternalLBAddressPoolName)),
			})
	}
	primaryIPConfig.LoadBalancerBackendAddressPools = backendAddressPools

	if s.PublicIPName != "" {
		primaryIPConfig.PublicIPAddress = &armnetwork.PublicIPAddress{
			ID: ptr.To(azure.PublicIPID(s.SubscriptionID, s.ResourceGroup, s.PublicIPName)),
		}
	}

	if s.AcceleratedNetworking == nil {
		// set accelerated networking to the capability of the VMSize
		if s.SKU == nil {
			return nil, errors.New("unable to get required network interface SKU from machine cache")
		}

		accelNet := s.SKU.HasCapability(resourceskus.AcceleratedNetworking)
		s.AcceleratedNetworking = &accelNet
	}

	dnsSettings := armnetwork.InterfaceDNSSettings{}
	if len(s.DNSServers) > 0 {
		dnsSettings.DNSServers = azure.PtrSlice(&s.DNSServers)
	}

	ipConfigurations := []*armnetwork.InterfaceIPConfiguration{
		{
			Name:       ptr.To("pipConfig"),
			Properties: primaryIPConfig,
		},
	}

	// Build additional IPConfigs if more than 1 is specified
	for i := 1; i < len(s.IPConfigs); i++ {
		c := s.IPConfigs[i]
		newIPConfigPropertiesFormat := &armnetwork.InterfaceIPConfigurationPropertiesFormat{}
		newIPConfigPropertiesFormat.Subnet = subnet
		config := &armnetwork.InterfaceIPConfiguration{
			Name:       ptr.To(s.Name + "-" + strconv.Itoa(i)),
			Properties: newIPConfigPropertiesFormat,
		}
		if c.PrivateIP != nil && *c.PrivateIP != "" {
			config.Properties.PrivateIPAllocationMethod = ptr.To(armnetwork.IPAllocationMethodStatic)
			config.Properties.PrivateIPAddress = c.PrivateIP
		} else {
			config.Properties.PrivateIPAllocationMethod = ptr.To(armnetwork.IPAllocationMethodDynamic)
		}

		if c.PublicIPAddress != nil && *c.PublicIPAddress != "" {
			config.Properties.PublicIPAddress = &armnetwork.PublicIPAddress{
				Properties: &armnetwork.PublicIPAddressPropertiesFormat{
					PublicIPAllocationMethod: ptr.To(armnetwork.IPAllocationMethodStatic),
					IPAddress:                c.PublicIPAddress,
				},
			}
		} else if c.PublicIPAddress != nil {
			config.Properties.PublicIPAddress = &armnetwork.PublicIPAddress{
				Properties: &armnetwork.PublicIPAddressPropertiesFormat{
					PublicIPAllocationMethod: ptr.To(armnetwork.IPAllocationMethodDynamic),
				},
			}
		}
		config.Properties.Primary = ptr.To(false)
		ipConfigurations = append(ipConfigurations, config)
	}
	if s.IPv6Enabled {
		ipv6Config := &armnetwork.InterfaceIPConfiguration{
			Name: ptr.To("ipConfigv6"),
			Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
				PrivateIPAddressVersion: ptr.To(armnetwork.IPVersionIPv6),
				Primary:                 ptr.To(false),
				Subnet:                  &armnetwork.Subnet{ID: subnet.ID},
			},
		}

		ipConfigurations = append(ipConfigurations, ipv6Config)
	}

	return armnetwork.Interface{
		Location:         ptr.To(s.Location),
		ExtendedLocation: converters.ExtendedLocationToNetworkSDK(s.ExtendedLocation),
		Properties: &armnetwork.InterfacePropertiesFormat{
			EnableAcceleratedNetworking: s.AcceleratedNetworking,
			IPConfigurations:            ipConfigurations,
			DNSSettings:                 &dnsSettings,
			EnableIPForwarding:          ptr.To(s.EnableIPForwarding),
		},
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Additional:  s.AdditionalTags,
		})),
	}, nil
}
