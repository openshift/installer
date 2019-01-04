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
		expectedError string
	}{
		{
			name:          "minimal",
			installConfig: validInstallConfig(),
		},
		{
			name: "missing name",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ObjectMeta.Name = ""
				return c
			}(),
			expectedError: `^metadata.name: Required value: cluster name required$`,
		},
		{
			name: "missing cluster ID",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ClusterID = ""
				return c
			}(),
			expectedError: `^clusterID: Required value: cluster ID required$`,
		},
		{
			name: "invalid ssh key",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.SSHKey = "bad-ssh-key"
				return c
			}(),
			expectedError: `^sshKey: Invalid value: "bad-ssh-key": ssh: no key found$`,
		},
		{
			name: "invalid base domain",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.BaseDomain = ".bad-domain."
				return c
			}(),
			expectedError: `^baseDomain: Invalid value: "\.bad-domain\.": a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '\.', and must start and end with an alphanumeric character \(e\.g\. 'example\.com', regex used for validation is '\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\(\\\.\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\)\*'\)$`,
		},
		{
			name: "invalid network type",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.Type = "bad-type"
				return c
			}(),
			expectedError: `^networking.type: Unsupported value: "bad-type": supported values: "Calico", "Kuryr", "OVNKubernetes", "OpenshiftSDN"$`,
		},
		{
			name: "invalid service cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceCIDR = ipnet.IPNet{}
				return c
			}(),
			expectedError: `^networking\.serviceCIDR: Invalid value: ipnet\.IPNet{IPNet:net\.IPNet{IP:net\.IP\(nil\), Mask:net\.IPMask\(nil\)}}: must use IPv4$`,
		},
		{
			name: "invalid cluster network cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworks[0].CIDR = "bad-cidr"
				return c
			}(),
			expectedError: `^networking\.clusterNetworks\[0]\.cidr: Invalid value: "bad-cidr": invalid CIDR address: bad-cidr$`,
		},
		{
			name: "overlapping cluster network cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworks[0].CIDR = "10.0.0.0/24"
				return c
			}(),
			expectedError: `^networking\.clusterNetworks\[0]\.cidr: Invalid value: "10\.0\.0\.0/24": cluster network CIDR must not overlap with serviceCIDR$`,
		},
		{
			name: "cluster network host subnet length too large",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworks[0].CIDR = "192.168.1.0/24"
				c.Networking.ClusterNetworks[0].HostSubnetLength = 9
				return c
			}(),
			expectedError: `^networking\.clusterNetworks\[0]\.hostSubnetLength: Invalid value: 0x9: cluster network host subnet length must not be greater than CIDR length$`,
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
			expectedError: `^machines: Required value: must specify a machine pool with a name of 'master'$`,
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
			expectedError: `^machines: Required value: must specify a machine pool with a name of 'worker'$`,
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
			expectedError: `^machines\[2]: Duplicate value: types\.MachinePool{Name:"master", Replicas:\(\*int64\)\(nil\), Platform:types\.MachinePoolPlatform{AWS:\(\*aws\.MachinePool\)\(nil\), Libvirt:\(\*libvirt\.MachinePool\)\(nil\), OpenStack:\(\*openstack\.MachinePool\)\(nil\)}}$`,
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
			expectedError: `^machines\[2]\.name: Unsupported value: "other": supported values: "master", "worker"$`,
		},
		{
			name: "missing platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{}
				return c
			}(),
			expectedError: `^platform: Invalid value: types\.Platform{AWS:\(\*aws\.Platform\)\(nil\), Libvirt:\(\*libvirt\.Platform\)\(nil\), OpenStack:\(\*openstack\.Platform\)\(nil\)}: must specify one of the platforms \(aws, openstack\)$`,
		},
		{
			name: "multiple platforms",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.Libvirt = &libvirt.Platform{}
				return c
			}(),
			expectedError: `^\[platform: Invalid value: types\.Platform{AWS:\(\*aws\.Platform\)\(0x[0-9a-f]*\), Libvirt:\(\*libvirt\.Platform\)\(0x[0-9a-f]*\), OpenStack:\(\*openstack\.Platform\)\(nil\)}: must only specify a single type of platform; cannot use both "aws" and "libvirt", platform\.libvirt\.uri: Invalid value: "": invalid URI "" \(no scheme\), platform\.libvirt\.network\.if: Required value, platform\.libvirt\.network\.ipRange: Invalid value: ipnet\.IPNet{IPNet:net\.IPNet{IP:net\.IP\(nil\), Mask:net\.IPMask\(nil\)}}: must use IPv4]$`,
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
			expectedError: `^platform\.aws\.region: Unsupported value: "": supported values: "ap-northeast-1", "ap-northeast-2", "ap-northeast-3", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "cn-north-1", "cn-northwest-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-west-3", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2"$`,
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
			expectedError: `^platform: Invalid value: types\.Platform{AWS:\(\*aws\.Platform\)\(nil\), Libvirt:\(\*libvirt\.Platform\)\(0x[0-9a-f]*\), OpenStack:\(\*openstack\.Platform\)\(nil\)}: must specify one of the platforms \(aws, openstack\)$`,
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
			expectedError: `^\[platform: Invalid value: types\.Platform{AWS:\(\*aws\.Platform\)\(nil\), Libvirt:\(\*libvirt\.Platform\)\(0x[0-9a-f]*\), OpenStack:\(\*openstack\.Platform\)\(nil\)}: must specify one of the platforms \(aws, openstack\), platform\.libvirt\.uri: Invalid value: "": invalid URI "" \(no scheme\), platform\.libvirt\.network\.if: Required value, platform\.libvirt\.network\.ipRange: Invalid value: ipnet\.IPNet{IPNet:net\.IPNet{IP:net\.IP\(nil\), Mask:net\.IPMask\(nil\)}}: must use IPv4]$`,
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
			expectedError: `^\[platform\.openstack\.cloud: Unsupported value: "": supported values: "test-cloud", platform\.openstack\.NetworkCIDRBlock: Invalid value: ipnet\.IPNet{IPNet:net\.IPNet{IP:net\.IP\(nil\), Mask:net\.IPMask\(nil\)}}: must use IPv4]$`,
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
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
