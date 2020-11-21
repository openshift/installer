package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/libvirt"
)

func defaultPlatform() *libvirt.Platform {
	n := &libvirt.Network{}
	SetNetworkDefaults(n)
	return &libvirt.Platform{
		URI:      DefaultURI,
		Network:  n,
		PoolPath: DefaultPoolPath,
	}
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *libvirt.Platform
		expected *libvirt.Platform
	}{
		{
			name:     "empty",
			platform: &libvirt.Platform{},
			expected: defaultPlatform(),
		},
		{
			name: "URI present",
			platform: &libvirt.Platform{
				URI: "test-uri",
			},
			expected: func() *libvirt.Platform {
				p := defaultPlatform()
				p.URI = "test-uri"
				return p
			}(),
		},
		{
			name: "Network present",
			platform: &libvirt.Platform{
				Network: func() *libvirt.Network {
					n := &libvirt.Network{}
					SetNetworkDefaults(n)
					n.IfName = "test-if"
					return n
				}(),
			},
			expected: func() *libvirt.Platform {
				p := defaultPlatform()
				p.Network = &libvirt.Network{}
				SetNetworkDefaults(p.Network)
				p.Network.IfName = "test-if"
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
