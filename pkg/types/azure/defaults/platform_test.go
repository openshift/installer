package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/azure"
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
