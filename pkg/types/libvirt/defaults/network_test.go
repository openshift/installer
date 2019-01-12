package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/libvirt"
)

func defaultNetwork() *libvirt.Network {
	return &libvirt.Network{
		IfName: defaultIfName,
	}
}

func TestSetNetworkDefaults(t *testing.T) {
	cases := []struct {
		name     string
		network  *libvirt.Network
		expected *libvirt.Network
	}{
		{
			name:     "empty",
			network:  &libvirt.Network{},
			expected: defaultNetwork(),
		},
		{
			name: "IfName present",
			network: &libvirt.Network{
				IfName: "test-if",
			},
			expected: func() *libvirt.Network {
				n := defaultNetwork()
				n.IfName = "test-if"
				return n
			}(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetNetworkDefaults(tc.network)
			assert.Equal(t, tc.expected, tc.network, "unexpected network")
		})
	}
}
