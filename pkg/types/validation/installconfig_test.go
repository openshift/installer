package validation

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/openstack/validation/mock"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		BaseDomain: "test-domain",
		Networking: &types.Networking{
			NetworkType:    "OpenShiftSDN",
			MachineCIDR:    ipnet.MustParseCIDR("10.0.0.0/16"),
			ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("172.30.0.0/16")},
			ClusterNetwork: []types.ClusterNetworkEntry{
				{
					CIDR:       *ipnet.MustParseCIDR("192.168.1.0/24"),
					HostPrefix: 28,
				},
			},
		},
		ControlPlane: &types.MachinePool{
			Name:     "master",
			Replicas: pointer.Int64Ptr(3),
		},
		Compute: []types.MachinePool{
			{
				Name:     "worker",
				Replicas: pointer.Int64Ptr(3),
			},
		},
		Platform: types.Platform{
			AWS: validAWSPlatform(),
		},
		PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
	}
}

func validAWSPlatform() *aws.Platform {
	return &aws.Platform{
		Region: "us-east-1",
	}
}

func validLibvirtPlatform() *libvirt.Platform {
	return &libvirt.Platform{
		URI: "qemu+tcp://192.168.122.1/system",
		Network: &libvirt.Network{
			IfName: "tt0",
		},
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
			name: "invalid version",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.APIVersion = "bad-version"
				return c
			}(),
			expectedError: fmt.Sprintf(`^apiVersion: Invalid value: "bad-version": install-config version must be %q`, types.InstallConfigVersion),
		},
		{
			name: "invalid name",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ObjectMeta.Name = "bad-name-"
				return c
			}(),
			expectedError: `^metadata.name: Invalid value: "bad-name-": a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '\.', and must start and end with an alphanumeric character \(e\.g\. 'example\.com', regex used for validation is '\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\(\\\.\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\)\*'\)$`,
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
			name: "overly long cluster domain",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ObjectMeta.Name = fmt.Sprintf("test-cluster%050d", 0)
				c.BaseDomain = fmt.Sprintf("test-domain%050d.a%060d.b%060d.c%060d", 0, 0, 0, 0)
				return c
			}(),
			expectedError: `^baseDomain: Invalid value: "` + fmt.Sprintf("test-cluster%050d.test-domain%050d.a%060d.b%060d.c%060d", 0, 0, 0, 0, 0) + `": must be no more than 253 characters$`,
		},
		{
			name: "missing networking",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = nil
				return c
			}(),
			expectedError: `^networking: Required value: networking is required$`,
		},
		{
			name: "invalid network type",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = ""
				return c
			}(),
			expectedError: `^networking.networkType: Required value: network provider type required$`,
		},
		{
			name: "missing service network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceNetwork = nil
				return c
			}(),
			expectedError: `^networking\.serviceNetwork: Required value: a service network is required$`,
		},
		{
			name: "invalid service network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceNetwork[0] = *ipnet.MustParseCIDR("13.0.128.0/16")
				return c
			}(),
			expectedError: `^networking\.serviceNetwork\[0\]: Invalid value: "13\.0\.128\.0/16": invalid network address. got 13\.0\.128\.0/16, expecting 13\.0\.0\.0/16$`,
		},
		{
			name: "overlapping service network and machine cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceNetwork[0] = *ipnet.MustParseCIDR("10.0.2.0/24")
				return c
			}(),
			expectedError: `^networking\.serviceNetwork\[0\]: Invalid value: "10\.0\.2\.0/24": service network must not overlap with machineCIDR$`,
		},
		{
			name: "overlapping service network and service network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceNetwork = []ipnet.IPNet{
					*ipnet.MustParseCIDR("13.0.0.0/16"),
					*ipnet.MustParseCIDR("13.0.2.0/24"),
				}

				return c
			}(),
			// also triggers the only-one-service-network validation
			expectedError: `^\[networking\.serviceNetwork\[1\]: Invalid value: "13\.0\.2\.0/24": service network must not overlap with service network 0, networking\.serviceNetwork: Invalid value: "13\.0\.0\.0/16, 13\.0\.2\.0/24": only one service network can be specified]$`,
		},
		{
			name: "missing machine cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineCIDR = nil
				return c
			}(),
			expectedError: `^networking\.machineCIDR: Required value: a machine CIDR is required$`,
		},
		{
			name: "invalid machine cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineCIDR = ipnet.MustParseCIDR("11.0.128.0/16")
				return c
			}(),
			expectedError: `^networking\.machineCIDR: Invalid value: "11\.0\.128\.0/16": invalid network address. got 11\.0\.128\.0/16, expecting 11\.0\.0\.0/16$`,
		},
		{
			name: "invalid cluster network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetwork = []types.ClusterNetworkEntry{{CIDR: *ipnet.MustParseCIDR("12.0.128.0/16"), HostPrefix: 23}}
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[0]\.cidr: Invalid value: "12\.0\.128\.0/16": invalid network address. got 12\.0\.128\.0/16, expecting 12\.0\.0\.0/16$`,
		},
		{
			name: "overlapping cluster network and machine cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetwork[0].CIDR = *ipnet.MustParseCIDR("10.0.3.0/24")
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[0]\.cidr: Invalid value: "10\.0\.3\.0/24": cluster network must not overlap with machine CIDR$`,
		},
		{
			name: "overlapping cluster network and service network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetwork[0].CIDR = *ipnet.MustParseCIDR("172.30.2.0/24")
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[0]\.cidr: Invalid value: "172\.30\.2\.0/24": cluster network must not overlap with service network 0$`,
		},
		{
			name: "overlapping cluster network and cluster network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetwork = []types.ClusterNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("12.0.0.0/16"), HostPrefix: 23},
					{CIDR: *ipnet.MustParseCIDR("12.0.3.0/24"), HostPrefix: 25},
				}
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[1]\.cidr: Invalid value: "12\.0\.3\.0/24": cluster network must not overlap with cluster network 0$`,
		},
		{
			name: "cluster network host prefix too large",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetwork[0].CIDR = *ipnet.MustParseCIDR("192.168.1.0/24")
				c.Networking.ClusterNetwork[0].HostPrefix = 23
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[0]\.hostPrefix: Invalid value: 23: cluster network host subnetwork prefix must not be larger size than CIDR 192.168.1.0/24$`,
		},
		{
			name: "missing control plane",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ControlPlane = nil
				return c
			}(),
			expectedError: `^controlPlane: Required value: controlPlane is required$`,
		},
		{
			name: "control plane with 0 replicas",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ControlPlane.Replicas = pointer.Int64Ptr(0)
				return c
			}(),
			expectedError: `^controlPlane.replicas: Invalid value: 0: number of control plane replicas must be positive$`,
		},
		{
			name: "invalid control plane",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ControlPlane.Replicas = nil
				return c
			}(),
			expectedError: `^controlPlane.replicas: Required value: replicas is required$`,
		},
		{
			name: "missing compute",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute = nil
				return c
			}(),
		},
		{
			name: "empty compute",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute = []types.MachinePool{}
				return c
			}(),
		},
		{
			name: "duplicate compute",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute = []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64Ptr(1),
					},
					{
						Name:     "worker",
						Replicas: pointer.Int64Ptr(2),
					},
				}
				return c
			}(),
			expectedError: `^compute\[1\]\.name: Duplicate value: "worker"$`,
		},
		{
			name: "no compute replicas",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute = []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64Ptr(0),
					},
				}
				return c
			}(),
		},
		{
			name: "invalid compute",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute = []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64Ptr(3),
						Platform: types.MachinePoolPlatform{
							OpenStack: &openstack.MachinePool{},
						},
					},
				}
				return c
			}(),
			expectedError: `^compute\[0\]\.platform.openstack: Invalid value: openstack.MachinePool{FlavorName:""}: cannot specify "openstack" for machine pool when cluster is using "aws"$`,
		},
		{
			name: "missing platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{}
				return c
			}(),
			expectedError: `^platform: Invalid value: types\.Platform{((, )?\w+:\(\*\w+\.Platform\)\(nil\))+}: must specify one of the platforms \(aws, azure, none, openstack, vsphere\)$`,
		},
		{
			name: "multiple platforms",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.Libvirt = validLibvirtPlatform()
				return c
			}(),
			expectedError: `^platform: Invalid value: types\.Platform{((, )?\w+:\(\*\w+\.Platform\)\(\w+\))+}: must only specify a single type of platform; cannot use both "aws" and "libvirt"$`,
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
			expectedError: `^platform\.aws\.region: Unsupported value: "": supported values: "ap-northeast-1", "ap-northeast-2", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-west-3", "sa-east-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2"$`,
		},
		{
			name: "valid libvirt platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Libvirt: validLibvirtPlatform(),
				}
				return c
			}(),
			expectedError: `^platform: Invalid value: types\.Platform{((, )?(\w+:\(\*\w+\.Platform\)\(nil\)|Libvirt:\(\*libvirt\.Platform\)\(0x[0-9a-f]*\)))+}: must specify one of the platforms \(aws, azure, none, openstack, vsphere\)$`,
		},
		{
			name: "invalid libvirt platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Libvirt: validLibvirtPlatform(),
				}
				c.Platform.Libvirt.URI = ""
				return c
			}(),
			expectedError: `^\[platform: Invalid value: types\.Platform{((, )?(\w+:\(\*\w+\.Platform\)\(nil\)|Libvirt:\(\*libvirt\.Platform\)\(0x[0-9a-f]*\)))+}: must specify one of the platforms \(aws, azure, none, openstack, vsphere\), platform\.libvirt\.uri: Invalid value: "": invalid URI "" \(no scheme\)]$`,
		},
		{
			name: "valid none platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					None: &none.Platform{},
				}
				return c
			}(),
		},
		{
			name: "valid openstack platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					OpenStack: &openstack.Platform{
						Region:          "test-region",
						Cloud:           "test-cloud",
						ExternalNetwork: "test-network",
						FlavorName:      "test-flavor",
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
			expectedError: `^platform\.openstack\.cloud: Unsupported value: "": supported values: "test-cloud"$`,
		},
		{
			name: "valid vsphere platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: &vsphere.Platform{},
				}
				return c
			}(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fetcher := mock.NewMockValidValuesFetcher(mockCtrl)
			fetcher.EXPECT().GetCloudNames().Return([]string{"test-cloud"}, nil).AnyTimes()
			fetcher.EXPECT().GetRegionNames(gomock.Any()).Return([]string{"test-region"}, nil).AnyTimes()
			fetcher.EXPECT().GetNetworkNames(gomock.Any()).Return([]string{"test-network"}, nil).AnyTimes()
			fetcher.EXPECT().GetFlavorNames(gomock.Any()).Return([]string{"test-flavor"}, nil).AnyTimes()
			fetcher.EXPECT().GetNetworkExtensionsAliases(gomock.Any()).Return([]string{"trunk"}, nil).AnyTimes()

			err := ValidateInstallConfig(tc.installConfig, fetcher).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
