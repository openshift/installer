package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	"github.com/openshift/installer/pkg/types/libvirt"
	libvirtdefaults "github.com/openshift/installer/pkg/types/libvirt/defaults"
	"github.com/openshift/installer/pkg/types/none"
	nonedefaults "github.com/openshift/installer/pkg/types/none/defaults"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

func defaultInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *defaultMachineCIDR},
			},
			NetworkType:    defaultNetworkType,
			ServiceNetwork: []ipnet.IPNet{*defaultServiceNetwork},
			ClusterNetwork: []types.ClusterNetworkEntry{
				{
					CIDR:       *defaultClusterNetwork,
					HostPrefix: int32(defaultHostPrefix),
				},
			},
		},
		ControlPlane: defaultMachinePool("master"),
		Compute:      []types.MachinePool{*defaultMachinePool("worker")},
		Publish:      types.ExternalPublishingStrategy,
	}
}

func defaultAWSInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.AWS = &aws.Platform{}
	awsdefaults.SetPlatformDefaults(c.Platform.AWS)
	return c
}

func defaultAzureInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.Azure = &azure.Platform{}
	azuredefaults.SetPlatformDefaults(c.Platform.Azure)
	return c
}

func defaultLibvirtInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Networking.MachineNetwork[0].CIDR = *libvirtdefaults.DefaultMachineCIDR
	c.Platform.Libvirt = &libvirt.Platform{}
	libvirtdefaults.SetPlatformDefaults(c.Platform.Libvirt)
	c.ControlPlane.Replicas = pointer.Int64Ptr(1)
	c.Compute[0].Replicas = pointer.Int64Ptr(1)
	return c
}

func defaultOpenStackInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.OpenStack = &openstack.Platform{}
	openstackdefaults.SetPlatformDefaults(c.Platform.OpenStack, c.Networking)
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
			name: "empty Azure",
			config: &types.InstallConfig{
				Platform: types.Platform{
					Azure: &azure.Platform{},
				},
			},
			expected: defaultAzureInstallConfig(),
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
					NetworkType: "test-networking-type",
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Networking.NetworkType = "test-networking-type"
				return c
			}(),
		},
		{
			name: "Service network present",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/8")},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Networking.ServiceNetwork[0] = *ipnet.MustParseCIDR("1.2.3.4/8")
				return c
			}(),
		},
		{
			name: "Cluster network present",
			config: &types.InstallConfig{
				Networking: &types.Networking{
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("8.8.0.0/18"),
							HostPrefix: 22,
						},
					},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Networking.ClusterNetwork = []types.ClusterNetworkEntry{
					{
						CIDR:       *ipnet.MustParseCIDR("8.8.0.0/18"),
						HostPrefix: 22,
					},
				}
				return c
			}(),
		},
		{
			name: "control plane present",
			config: &types.InstallConfig{
				ControlPlane: &types.MachinePool{},
			},
			expected: defaultInstallConfig(),
		},
		{
			name: "Compute present",
			config: &types.InstallConfig{
				Compute: []types.MachinePool{{Name: "test-compute"}},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Compute = []types.MachinePool{*defaultMachinePool("test-compute")}
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
			name: "Azure platform present",
			config: &types.InstallConfig{
				Platform: types.Platform{
					Azure: &azure.Platform{},
				},
			},
			expected: func() *types.InstallConfig {
				c := defaultAzureInstallConfig()
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
