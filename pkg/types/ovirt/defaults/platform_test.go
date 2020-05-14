package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func defaultPlatform() *ovirt.Platform {
	return &ovirt.Platform{
		NetworkName: DefaultNetworkName,
	}
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *ovirt.Platform
		expected *ovirt.Platform
	}{
		{
			name:     "empty",
			platform: &ovirt.Platform{},
			expected: defaultPlatform(),
		},
		{
			name:     "URL present",
			platform: &ovirt.Platform{},
			expected: func() *ovirt.Platform {
				p := defaultPlatform()
				return p
			}(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform)
			assert.Equal(t, tc.expected, tc.platform, "unexpected platform")
		})
	}
}
