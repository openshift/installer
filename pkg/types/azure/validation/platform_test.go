package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

func validPlatform() *azure.Platform {
	return &azure.Platform{
		Region:                      "eastus",
		BaseDomainResourceGroupName: "group",
		OutboundType:                azure.LoadbalancerOutboundType,
		CloudName:                   azure.PublicCloud,
	}
}

func validNetworkPlatform() *azure.Platform {
	p := validPlatform()
	p.NetworkResourceGroupName = "networkresourcegroup"
	p.VirtualNetwork = "virtualnetwork"
	p.ComputeSubnet = "computesubnet"
	p.ControlPlaneSubnet = "controlplanesubnet"
	return p
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *azure.Platform
		expected string
	}{
		{
			name: "invalid region",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.Region = ""
				return p
			}(),
			expected: `^test-path\.region: Required value: region should be set to one of the supported Azure regions$`,
		},
		{
			name: "invalid baseDomainResourceGroupName",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.BaseDomainResourceGroupName = ""
				return p
			}(),
			expected: `^test-path\.baseDomainResourceGroupName: Required value: baseDomainResourceGroupName is the resource group name where the azure dns zone is deployed$`,
		},
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "valid machine pool",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &azure.MachinePool{}
				return p
			}(),
		},
		{
			name:     "valid subnets & virtual network",
			platform: validNetworkPlatform(),
		},
		{
			name: "missing subnets",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.ControlPlaneSubnet = ""
				return p
			}(),
			expected: `^test-path\.controlPlaneSubnet: Required value: must provide a control plane subnet when a virtual network is specified$`,
		},
		{
			name: "subnets missing virtual network",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.ControlPlaneSubnet = ""
				p.VirtualNetwork = ""
				return p
			}(),
			expected: `^test-path\.virtualNetwork: Required value: must provide a virtual network when supplying subnets$`,
		},
		{
			name: "missing network resource group",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.NetworkResourceGroupName = ""
				return p
			}(),
			expected: `^\[test-path\.networkResourceGroupName: Required value: must provide a network resource group when a virtual network is specified, test-path\.networkResourceGroupName: Required value: must provide a network resource group when supplying subnets\]$`,
		},
		{
			name: "missing cloud name",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.CloudName = ""
				return p
			}(),
			expected: `^test-path\.cloudName: Unsupported value: "": supported values:`,
		},
		{
			name: "invalid cloud name",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.CloudName = azure.CloudEnvironment("AzureOtherCloud")
				return p
			}(),
			expected: `^test-path\.cloudName: Unsupported value: "AzureOtherCloud": supported values:`,
		},
		{
			name: "invalid outbound type",
			platform: func() *azure.Platform {
				p := validNetworkPlatform()
				p.OutboundType = "random-egress"
				return p
			}(),
			expected: `^test-path\.outboundType: Unsupported value: "random-egress": supported values: "Loadbalancer", "UserDefinedRouting"$`,
		},
		{
			name: "invalid user defined type",
			platform: func() *azure.Platform {
				p := validPlatform()
				p.OutboundType = azure.UserDefinedRoutingOutboundType
				return p
			}(),
			expected: `^test-path\.outboundType: Invalid value: "UserDefinedRouting": UserDefinedRouting is only allowed when installing to pre-existing network$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, types.ExternalPublishingStrategy, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
