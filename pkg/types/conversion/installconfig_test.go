package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

func TestConvertInstallConfig(t *testing.T) {
	cases := []struct {
		name          string
		config        *types.InstallConfig
		expected      *types.InstallConfig
		expectedError string
	}{
		{
			name: "empty",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
			},
		},
		{
			name: "empty networking",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{},
			},
		},
		{
			// all deprecated fields
			name: "old networking",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1beta3",
				},
				Networking: &types.Networking{
					DeprecatedMachineCIDR: ipnet.MustParseCIDR("1.1.1.1/24"),
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
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					NetworkType: "foo",
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("1.1.1.1/24")},
					},
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
					DeprecatedMachineCIDR: ipnet.MustParseCIDR("1.1.1.1/24"),
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
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("1.1.1.1/24")},
					},
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
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("1.1.1.1/24")},
					},
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
		{
			name: "empty APIVersion",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "",
				},
			},
			expectedError: "no version was provided",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ConvertInstallConfig(tc.config)
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, tc.config, "unexpected install config")
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
