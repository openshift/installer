package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

func TestConvertInstallConfig(t *testing.T) {
	cases := []struct {
		name     string
		config   *types.InstallConfig
		expected *types.InstallConfig
	}{
		{
			name:     "empty",
			config:   &types.InstallConfig{},
			expected: &types.InstallConfig{},
		},
		{
			name: "empty networking",
			config: &types.InstallConfig{
				Networking: &types.Networking{},
			},
			expected: &types.InstallConfig{
				Networking: &types.Networking{},
			},
		},
		{
			// all deprecated fields
			name: "old networking",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					MachineCIDR:           ipnet.MustParseCIDR("1.1.1.1/24"),
					DeprecatedType:        "foo",
					DeprecatedServiceCIDR: ipnet.MustParseCIDR("1.2.3.4/32"),
					DeprecatedClusterNetworks: []types.ClusterNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("1.2.3.5/32"),
							DeprecatedHostSubnetLength: 8,
						},
					},
				},
			},
			expected: &types.InstallConfig{
				Networking: &types.Networking{
					MachineCIDR:    ipnet.MustParseCIDR("1.1.1.1/24"),
					NetworkType:    "foo",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/32")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("1.2.3.5/32"),

							HostPrefix:                 24,
							DeprecatedHostSubnetLength: 8,
						},
					},

					// deprecated fields are preserved
					DeprecatedType:        "foo",
					DeprecatedServiceCIDR: ipnet.MustParseCIDR("1.2.3.4/32"),
					DeprecatedClusterNetworks: []types.ClusterNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("1.2.3.5/32"),

							HostPrefix:                 24,
							DeprecatedHostSubnetLength: 8,
						},
					},
				},
			},
		},
		{
			name: "new networking",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					MachineCIDR:    ipnet.MustParseCIDR("1.1.1.1/24"),
					NetworkType:    "foo",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/32")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("1.2.3.5/32"),
							HostPrefix: 24,
						},
					},
				},
			},
			expected: &types.InstallConfig{
				Networking: &types.Networking{
					MachineCIDR:    ipnet.MustParseCIDR("1.1.1.1/24"),
					NetworkType:    "foo",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/32")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("1.2.3.5/32"),
							HostPrefix: 24,
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ConvertInstallConfig(tc.config)
			assert.Equal(t, tc.expected, tc.config, "unexpected install config")
		})
	}
}
