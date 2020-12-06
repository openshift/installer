package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/openstack"
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
		{
			name: "deprecated OpenShiftSDN spelling",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					NetworkType: "OpenshiftSDN",
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					NetworkType: "OpenShiftSDN",
				},
			},
		},
		{
			name: "deprecated OpenStack LbFloatingIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
						APIFloatingIP:          "10.0.109.1",
					},
				},
			},
		},
		{
			name: "deprecated OpenStack LbFloatingIP with APIFloatingIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
						APIFloatingIP:          "10.0.109.1",
					},
				},
			},
			expectedError: "cannot specify lbFloatingIP and apiFloatingIP together",
		},

		// BareMetal platform conversions
		{
			name: "baremetal external DHCP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedProvisioningDHCPExternal: true,
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedProvisioningDHCPExternal: true,
						ProvisioningNetwork:                "Unmanaged",
					},
				},
			},
		},
		{
			name: "baremetal provisioningHostIP -> clusterProvisioningIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedProvisioningHostIP: "172.22.0.3",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						ClusterProvisioningIP:        "172.22.0.3",
						DeprecatedProvisioningHostIP: "172.22.0.3",
					},
				},
			},
		},
		{
			name: "baremetal provisioningHostIP mismatch clusterProvisioningIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						ClusterProvisioningIP:        "172.22.0.4",
						DeprecatedProvisioningHostIP: "172.22.0.3",
					},
				},
			},
			expectedError: "provisioningHostIP is deprecated; only clusterProvisioningIP needs to be specified",
		},
		{
			name: "deprecated OpenStack computeFlavor",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedFlavorName: "big-flavor",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedFlavorName: "big-flavor",
						DefaultMachinePlatform: &openstack.MachinePool{
							FlavorName: "big-flavor",
						},
					},
				},
			},
		},
		{
			name: "deprecated OpenStack computeFlavor with type in defaultMachinePlatform",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedFlavorName: "flavor1",
						DefaultMachinePlatform: &openstack.MachinePool{
							FlavorName: "flavor2",
						},
					},
				},
			},
			expectedError: "cannot specify computeFlavor and type in defaultMachinePlatform together",
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
