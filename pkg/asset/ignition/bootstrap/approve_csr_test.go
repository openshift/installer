package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func machineNetwork(cidrs ...string) []types.MachineNetworkEntry {
	entries := make([]types.MachineNetworkEntry, 0, len(cidrs))
	for _, cidr := range cidrs {
		entries = append(entries, types.MachineNetworkEntry{CIDR: *ipnet.MustParseCIDR(cidr)})
	}
	return entries
}

func TestApproveCSRNetworks(t *testing.T) {
	cases := []struct {
		name     string
		config   *types.InstallConfig
		expected []string
	}{{
		name:     "nil config",
		config:   nil,
		expected: nil,
	}, {
		name:     "no networking",
		config:   &types.InstallConfig{},
		expected: nil,
	}, {
		name: "single stack",
		config: &types.InstallConfig{
			Networking: &types.Networking{MachineNetwork: machineNetwork("10.0.0.0/16")},
		},
		expected: []string{"10.0.0.0/16"},
	}, {
		name: "dual stack keeps both families",
		config: &types.InstallConfig{
			Networking: &types.Networking{MachineNetwork: machineNetwork("10.0.0.0/16", "fd69::/64")},
		},
		expected: []string{"10.0.0.0/16", "fd69::/64"},
	}, {
		name: "baremetal adds the provisioning network",
		config: &types.InstallConfig{
			Networking: &types.Networking{MachineNetwork: machineNetwork("192.168.111.0/24")},
			Platform: types.Platform{BareMetal: &baremetal.Platform{
				ProvisioningNetwork:     baremetal.ManagedProvisioningNetwork,
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.22.0.0/24"),
			}},
		},
		expected: []string{"192.168.111.0/24", "172.22.0.0/24"},
	}, {
		name: "baremetal with an unmanaged provisioning network still adds it",
		config: &types.InstallConfig{
			Networking: &types.Networking{MachineNetwork: machineNetwork("192.168.111.0/24")},
			Platform: types.Platform{BareMetal: &baremetal.Platform{
				ProvisioningNetwork:     baremetal.UnmanagedProvisioningNetwork,
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.22.0.0/24"),
			}},
		},
		expected: []string{"192.168.111.0/24", "172.22.0.0/24"},
	}, {
		name: "baremetal with provisioning disabled adds nothing",
		config: &types.InstallConfig{
			Networking: &types.Networking{MachineNetwork: machineNetwork("192.168.111.0/24")},
			Platform: types.Platform{BareMetal: &baremetal.Platform{
				ProvisioningNetwork:     baremetal.DisabledProvisioningNetwork,
				ProvisioningNetworkCIDR: ipnet.MustParseCIDR("172.22.0.0/24"),
			}},
		},
		expected: []string{"192.168.111.0/24"},
	}, {
		name: "baremetal without a provisioning CIDR adds nothing",
		config: &types.InstallConfig{
			Networking: &types.Networking{MachineNetwork: machineNetwork("192.168.111.0/24")},
			Platform: types.Platform{BareMetal: &baremetal.Platform{
				ProvisioningNetwork: baremetal.ManagedProvisioningNetwork,
			}},
		},
		expected: []string{"192.168.111.0/24"},
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, approveCSRNetworks(tc.config))
		})
	}
}
