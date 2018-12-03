package validation

import (
	"testing"

	"github.com/golang/mock/gomock"
	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/openstack/validation/mock"
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		ClusterID:  "test-cluster-id",
		BaseDomain: "test-domain",
		Networking: types.Networking{
			Type:        "OpenshiftSDN",
			ServiceCIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
			ClusterNetworks: []netopv1.ClusterNetwork{
				{
					CIDR:             "192.168.1.0/24",
					HostSubnetLength: 4,
				},
			},
		},
		Machines: []types.MachinePool{
			{
				Name: "master",
			},
			{
				Name: "worker",
			},
		},
		Platform: types.Platform{
			AWS: &aws.Platform{
				Region: "us-east-1",
			},
		},
		PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
	}
}

func TestValidateInstallConfig(t *testing.T) {
	cases := []struct {
		name          string
		installConfig *types.InstallConfig
		valid         bool
	}{
		{
			name:          "minimal",
			installConfig: validInstallConfig(),
			valid:         true,
		},
		{
			name: "missing name",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ObjectMeta.Name = ""
				return c
			}(),
			valid: false,
		},
		{
			name: "missing cluster ID",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ClusterID = ""
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid ssh key",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.SSHKey = "bad-ssh-key"
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid base domain",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.BaseDomain = ".bad-domain."
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid network type",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.Type = "bad-type"
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid service cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceCIDR = ipnet.IPNet{}
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid cluster network cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworks[0].CIDR = "bad-cidr"
				return c
			}(),
			valid: false,
		},
		{
			name: "overlapping cluster network cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworks[0].CIDR = "10.0.0.0/24"
				return c
			}(),
			valid: false,
		},
		{
			name: "cluster network host subnet length too large",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworks[0].CIDR = "192.168.1.0/24"
				c.Networking.ClusterNetworks[0].HostSubnetLength = 9
				return c
			}(),
			valid: false,
		},
		{
			name: "missing master machine pool",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Machines = []types.MachinePool{
					{
						Name: "worker",
					},
				}
				return c
			}(),
			valid: false,
		},
		{
			name: "missing worker machine pool",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Machines = []types.MachinePool{
					{
						Name: "master",
					},
				}
				return c
			}(),
			valid: false,
		},
		{
			name: "duplicate machine pool",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Machines = append(c.Machines, types.MachinePool{
					Name: "master",
				})
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid machine pool",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Machines = append(c.Machines, types.MachinePool{
					Name: "other",
				})
				return c
			}(),
			valid: false,
		},
		{
			name: "missing platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{}
				return c
			}(),
			valid: false,
		},
		{
			name: "multiple platforms",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.Libvirt = &libvirt.Platform{}
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid aws platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					AWS: &aws.Platform{},
				}
				return c
			}(),
			valid: false,
		},
		{
			name: "valid libvirt platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Libvirt: &libvirt.Platform{
						URI: "qemu+tcp://192.168.122.1/system",
						Network: libvirt.Network{
							IfName:  "tt0",
							IPRange: *ipnet.MustParseCIDR("10.0.0.0/16"),
						},
					},
				}
				return c
			}(),
			valid: true,
		},
		{
			name: "invalid libvirt platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Libvirt: &libvirt.Platform{},
				}
				return c
			}(),
			valid: false,
		},
		{
			name: "valid openstack platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					OpenStack: &openstack.Platform{
						Region:           "test-region",
						NetworkCIDRBlock: *ipnet.MustParseCIDR("10.0.0.0/16"),
						BaseImage:        "test-image",
						Cloud:            "test-cloud",
						ExternalNetwork:  "test-network",
						FlavorName:       "test-flavor",
					},
				}
				return c
			}(),
			valid: true,
		},
		{
			name: "invalid openstack platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					OpenStack: &openstack.Platform{},
				}
				return c
			}(),
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fetcher := mock.NewMockValidValuesFetcher(mockCtrl)
			fetcher.EXPECT().GetCloudNames().Return([]string{"test-cloud"}, nil).AnyTimes()
			fetcher.EXPECT().GetRegionNames(gomock.Any()).Return([]string{"test-region"}, nil).AnyTimes()
			fetcher.EXPECT().GetImageNames(gomock.Any()).Return([]string{"test-image"}, nil).AnyTimes()
			fetcher.EXPECT().GetNetworkNames(gomock.Any()).Return([]string{"test-network"}, nil).AnyTimes()
			fetcher.EXPECT().GetFlavorNames(gomock.Any()).Return([]string{"test-flavor"}, nil).AnyTimes()

			err := ValidateInstallConfig(tc.installConfig, fetcher).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
