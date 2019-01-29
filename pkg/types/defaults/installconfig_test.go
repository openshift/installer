package defaults

import (
	"testing"

	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/openshift/installer/pkg/types/libvirt"
	libvirtdefaults "github.com/openshift/installer/pkg/types/libvirt/defaults"
	"github.com/openshift/installer/pkg/types/none"
	nonedefaults "github.com/openshift/installer/pkg/types/none/defaults"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	"k8s.io/utils/pointer"
)

func defaultInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineCIDR: defaultMachineCIDR,
			Type:        netopv1.NetworkTypeOpenshiftSDN,
			ServiceCIDR: defaultServiceCIDR,
			ClusterNetworks: []netopv1.ClusterNetwork{
				{
					CIDR:             defaultClusterCIDR,
					HostSubnetLength: uint32(defaultHostSubnetLength),
				},
			},
		},
		Machines: []types.MachinePool{
			{
				Name:     "master",
				Replicas: func(x int64) *int64 { return &x }(3),
			},
			{
				Name:     "worker",
				Replicas: func(x int64) *int64 { return &x }(3),
			},
		},
	}
}

func defaultAWSInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.AWS = &aws.Platform{}
	awsdefaults.SetPlatformDefaults(c.Platform.AWS)
	return c
}

func defaultLibvirtInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Networking.MachineCIDR = libvirtdefaults.DefaultMachineCIDR
	c.Platform.Libvirt = &libvirt.Platform{}
	libvirtdefaults.SetPlatformDefaults(c.Platform.Libvirt)
	for i, m := range c.Machines {
		m.Replicas = func(x int64) *int64 { return &x }(1)
		c.Machines[i] = m
	}
	return c
}

func defaultOpenStackInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.OpenStack = &openstack.Platform{}
	openstackdefaults.SetPlatformDefaults(c.Platform.OpenStack)
	return c
}

func defaultNoneInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.None = &none.Platform{}
	nonedefaults.SetPlatformDefaults(c.Platform.None)
	return c
}

func TestSetInstallConfigDefaults(t *testing.T) {
	cases := []struct {
		name     string
		config   *types.InstallConfig
		expected *types.InstallConfig
	}{
		{
			name:     "empty",
			config:   &types.InstallConfig{},
			expected: defaultInstallConfig(),
		},
		{
			name: "empty AWS",
			config: &types.InstallConfig{
				Platform: types.Platform{
					AWS: &aws.Platform{},
				},
			},
			expected: defaultAWSInstallConfig(),
		},
		{
			name: "empty Libvirt",
			config: &types.InstallConfig{
				Platform: types.Platform{
					Libvirt: &libvirt.Platform{},
				},
			},
			expected: defaultLibvirtInstallConfig(),
		},
		{
			name: "empty OpenStack",
			config: &types.InstallConfig{
				Platform: types.Platform{
					OpenStack: &openstack.Platform{},
				},
			},
			expected: defaultOpenStackInstallConfig(),
		},
		{
			name: "Networking present",
			config: &types.InstallConfig{
				Networking: &types.Networking{},
			},
			expected: defaultInstallConfig(),
		},
		{
			name: "Networking types present",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					Type: netopv1.NetworkType("test-networking-type"),
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Networking.Type = netopv1.NetworkType("test-networking-type")
				return c
			}(),
		},
		{
			name: "Service CIDR present",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					ServiceCIDR: ipnet.MustParseCIDR("1.2.3.4/8"),
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Networking.ServiceCIDR = ipnet.MustParseCIDR("1.2.3.4/8")
				return c
			}(),
		},
		{
			name: "Cluster Networks present",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetworks: []netopv1.ClusterNetwork{
						{
							CIDR:             "test-cidr",
							HostSubnetLength: 10,
						},
					},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Networking.ClusterNetworks = []netopv1.ClusterNetwork{
					{
						CIDR:             "test-cidr",
						HostSubnetLength: 10,
					},
				}
				return c
			}(),
		},
		{
			name: "Pod CIDR present",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					PodCIDR: ipnet.MustParseCIDR("1.2.3.4/8"),
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Networking.ClusterNetworks = nil
				c.Networking.PodCIDR = ipnet.MustParseCIDR("1.2.3.4/8")
				return c
			}(),
		},
		{
			name: "Machines present",
			config: &types.InstallConfig{
				Machines: []types.MachinePool{{Name: "test-machine"}},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Machines = []types.MachinePool{{
					Name:     "test-machine",
					Replicas: pointer.Int64Ptr(0),
				}}
				return c
			}(),
		},
		{
			name: "AWS platform present",
			config: &types.InstallConfig{
				Platform: types.Platform{
					AWS: &aws.Platform{},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultAWSInstallConfig()
				return c
			}(),
		},
		{
			name: "Libvirt platform present",
			config: &types.InstallConfig{
				Platform: types.Platform{
					Libvirt: &libvirt.Platform{
						URI: "test-uri",
					},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultLibvirtInstallConfig()
				c.Platform.Libvirt.URI = "test-uri"
				return c
			}(),
		},
		{
			name: "OpenStack platform present",
			config: &types.InstallConfig{
				Platform: types.Platform{
					OpenStack: &openstack.Platform{},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultOpenStackInstallConfig()
				return c
			}(),
		},
		{
			name: "None platform present",
			config: &types.InstallConfig{
				Platform: types.Platform{
					None: &none.Platform{},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultNoneInstallConfig()
				return c
			}(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetInstallConfigDefaults(tc.config)
			assert.Equal(t, tc.expected, tc.config, "unexpected install config")
		})
	}
}
