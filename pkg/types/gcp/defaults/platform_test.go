package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/network"
)

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *gcp.Platform
		expected *gcp.Platform
	}{
		{
			name:     "empty platform should default IPFamily to IPv4",
			platform: &gcp.Platform{},
			expected: &gcp.Platform{
				IPFamily: network.IPv4,
			},
		},
		{
			name: "IPFamily already set should not be overridden",
			platform: &gcp.Platform{
				IPFamily: network.DualStackIPv4Primary,
			},
			expected: &gcp.Platform{
				IPFamily: network.DualStackIPv4Primary,
			},
		},
		{
			name: "DualStackIPv6Primary should not be overridden",
			platform: &gcp.Platform{
				IPFamily: network.DualStackIPv6Primary,
			},
			expected: &gcp.Platform{
				IPFamily: network.DualStackIPv6Primary,
			},
		},
		{
			name:     "nil platform should not panic",
			platform: nil,
			expected: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform)
			assert.Equal(t, tc.expected, tc.platform)
		})
	}
}
