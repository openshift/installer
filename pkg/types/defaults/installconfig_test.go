package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	"github.com/openshift/installer/pkg/types/none"
	nonedefaults "github.com/openshift/installer/pkg/types/none/defaults"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	"github.com/openshift/installer/pkg/types/ovirt"
	ovirtdefaults "github.com/openshift/installer/pkg/types/ovirt/defaults"
)

func defaultInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		AdditionalTrustBundlePolicy: defaultAdditionalTrustBundlePolicy(),
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *DefaultMachineCIDR},
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

func defaultInstallConfigWithEdge() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Compute = append(c.Compute, *defaultMachinePool("edge"))
	return c
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

func defaultOpenStackInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.OpenStack = &openstack.Platform{}
	openstackdefaults.SetPlatformDefaults(c.Platform.OpenStack, c.Networking)
	return c
}

func defaultOvirtInstallConfig() *types.InstallConfig {
	c := defaultInstallConfig()
	c.Platform.Ovirt = &ovirt.Platform{}
	ovirtdefaults.SetPlatformDefaults(c.Platform.Ovirt)
	ovirtdefaults.SetControlPlaneDefaults(c.Platform.Ovirt, c.ControlPlane)
	for i := range c.Compute {
		ovirtdefaults.SetComputeDefaults(c.Platform.Ovirt, &c.Compute[i])
	}
	return c
}

func defaultAdditionalTrustBundlePolicy() types.PolicyType {
	return types.PolicyProxyOnly
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
			name: "empty OpenStack",
			config: &types.InstallConfig{
				Platform: types.Platform{
					OpenStack: &openstack.Platform{},
				},
			},
			expected: defaultOpenStackInstallConfig(),
		},
		{
			name: "empty oVirt",
			config: &types.InstallConfig{
				Platform: types.Platform{
					Ovirt: &ovirt.Platform{},
				},
			},
			expected: defaultOvirtInstallConfig(),
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
			name: "arbiter present",
			config: &types.InstallConfig{
				Arbiter: &types.MachinePool{},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Arbiter = defaultMachinePoolWithReplicaCount("arbiter", 0)
				return c
			}(),
		},
		{
			name: "Compute present",
			config: &types.InstallConfig{
				Compute: []types.MachinePool{{Name: "worker"}},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfig()
				c.Compute = []types.MachinePool{*defaultMachinePool("worker")}
				return c
			}(),
		},
		{
			name: "Edge Compute present",
			config: &types.InstallConfig{
				Compute: []types.MachinePool{{Name: "worker"}, {Name: "edge"}},
			},
			expected: func() *types.InstallConfig {
				c := defaultInstallConfigWithEdge()
				c.Compute = []types.MachinePool{
					*defaultMachinePool("worker"),
					*defaultEdgeMachinePool("edge"),
				}
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
