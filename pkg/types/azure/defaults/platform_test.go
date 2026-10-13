package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/network"
)

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *azure.Platform
		expected *azure.Platform
	}{
		{
			name:     "empty",
			platform: &azure.Platform{},
			expected: &azure.Platform{
				CloudName:    azure.PublicCloud,
				OutboundType: azure.LoadbalancerOutboundType,
				IPFamily:     network.IPv4,
			},
		},
		{
			name: "default cloud",
			platform: &azure.Platform{
				CloudName: azure.PublicCloud,
			},
			expected: &azure.Platform{
				CloudName:    azure.PublicCloud,
				OutboundType: azure.LoadbalancerOutboundType,
				IPFamily:     network.IPv4,
			},
		},
		{
			name: "non-default cloud name",
			platform: &azure.Platform{
				CloudName: azure.USGovernmentCloud,
			},
			expected: &azure.Platform{
				CloudName:    azure.USGovernmentCloud,
				OutboundType: azure.LoadbalancerOutboundType,
				IPFamily:     network.IPv4,
			},
		},
		{
			name: "default outbound",
			platform: &azure.Platform{
				CloudName:    azure.PublicCloud,
				OutboundType: azure.LoadbalancerOutboundType,
			},
			expected: &azure.Platform{
				CloudName:    azure.PublicCloud,
				OutboundType: azure.LoadbalancerOutboundType,
				IPFamily:     network.IPv4,
			},
		},
		{
			name: "non-default cloud name",
			platform: &azure.Platform{
				CloudName:    azure.USGovernmentCloud,
				OutboundType: azure.UserDefinedRoutingOutboundType,
			},
			expected: &azure.Platform{
				CloudName:    azure.USGovernmentCloud,
				OutboundType: azure.UserDefinedRoutingOutboundType,
				IPFamily:     network.IPv4,
			},
		},
		{
			name: "non-default IPFamily",
			platform: &azure.Platform{
				IPFamily: network.DualStackIPv4Primary,
			},
			expected: &azure.Platform{
				CloudName:    azure.PublicCloud,
				OutboundType: azure.LoadbalancerOutboundType,
				IPFamily:     network.DualStackIPv4Primary,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform)
			assert.Equal(t, tc.expected, tc.platform, "unexpected platform")
		})
	}
}
