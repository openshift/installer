package validation

import (
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/openstack/validation/mock"
	"github.com/openshift/installer/pkg/types/ovirt"
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
		BaseDomain:   "test-domain",
		Networking:   validIPv4NetworkingConfig(),
		ControlPlane: validMachinePool("master"),
		Compute:      []types.MachinePool{*validMachinePool("worker")},
		Platform: types.Platform{
			AWS: validAWSPlatform(),
		},
		PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
		Publish:    types.ExternalPublishingStrategy,
		Proxy: &types.Proxy{
			HTTPProxy:  "http://user:password@127.0.0.1:8080",
			HTTPSProxy: "https://user:password@127.0.0.1:8080",
			NoProxy:    "valid-proxy.com, 172.30.0.0/16",
		},
	}
}

func validAWSPlatform() *aws.Platform {
	return &aws.Platform{
		Region: "us-east-1",
	}
}

func validAzurePlatform() *azure.Platform {
	return &azure.Platform{
		Region:                      "us-east-1",
		BaseDomainResourceGroupName: "my-resource-group",
	}
}

func validGCPPlatform() *gcp.Platform {
	return &gcp.Platform{
		ProjectID: "myProject",
		Region:    "us-east1",
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

func validVSpherePlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenter:          "test-server",
		Username:         "test-username",
		Password:         "test-password",
		Datacenter:       "test-datacenter",
		DefaultDatastore: "test-datastore",
	}
}

func validBareMetalPlatform() *baremetal.Platform {
	iface, _ := net.Interfaces()
	return &baremetal.Platform{
		LibvirtURI:                   "qemu+tcp://192.168.122.1/system",
		ProvisioningNetworkInterface: "ens3",
		ProvisioningNetworkCIDR:      ipnet.MustParseCIDR("192.168.111.0/24"),
		BootstrapProvisioningIP:      "192.168.111.1",
		ClusterProvisioningIP:        "192.168.111.2",
		Hosts: []*baremetal.Host{
			{
				Name:           "host1",
				BootMACAddress: "CA:FE:CA:FE:00:00",
				BMC: baremetal.BMC{
					Username: "root",
					Password: "password",
					Address:  "ipmi://192.168.111.1",
				},
			},
			{
				Name:           "host2",
				BootMACAddress: "CA:FE:CA:FE:00:01",
				BMC: baremetal.BMC{
					Username: "root",
					Password: "password",
					Address:  "ipmi://192.168.111.2",
				},
			},
		},
		ExternalBridge:         iface[0].Name,
		ProvisioningBridge:     iface[0].Name,
		DefaultMachinePlatform: &baremetal.MachinePool{},
		APIVIP:                 "10.0.0.5",
		IngressVIP:             "10.0.0.4",
	}
}

func validOpenStackPlatform() *openstack.Platform {
	return &openstack.Platform{
		Cloud:           "test-cloud",
		ExternalNetwork: "test-network",
		FlavorName:      "test-flavor",
	}
}

func validIPv4NetworkingConfig() *types.Networking {
	return &types.Networking{
		NetworkType: "OpenShiftSDN",
		MachineNetwork: []types.MachineNetworkEntry{
			{
				CIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
			},
		},
		ServiceNetwork: []ipnet.IPNet{
			*ipnet.MustParseCIDR("172.30.0.0/16"),
		},
		ClusterNetwork: []types.ClusterNetworkEntry{
			{
				CIDR:       *ipnet.MustParseCIDR("192.168.1.0/24"),
				HostPrefix: 28,
			},
		},
	}
}

func validIPv6NetworkingConfig() *types.Networking {
	return &types.Networking{
		NetworkType: "OVNKubernetes",
		MachineNetwork: []types.MachineNetworkEntry{
			{
				CIDR: *ipnet.MustParseCIDR("ffd0::/48"),
			},
		},
		ServiceNetwork: []ipnet.IPNet{
			*ipnet.MustParseCIDR("ffd1::/48"),
		},
		ClusterNetwork: []types.ClusterNetworkEntry{
			{
				CIDR:       *ipnet.MustParseCIDR("ffd2::/48"),
				HostPrefix: 64,
			},
		},
	}
}

func validDualStackNetworkingConfig() *types.Networking {
	return &types.Networking{
		NetworkType: "OVNKubernetes",
		MachineNetwork: []types.MachineNetworkEntry{
			{
				CIDR: *ipnet.MustParseCIDR("ffd0::/48"),
			},
			{
				CIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
			},
		},
		ServiceNetwork: []ipnet.IPNet{
			*ipnet.MustParseCIDR("ffd1::/48"),
			*ipnet.MustParseCIDR("172.30.0.0/16"),
		},
		ClusterNetwork: []types.ClusterNetworkEntry{
			{
				CIDR:       *ipnet.MustParseCIDR("ffd2::/48"),
				HostPrefix: 64,
			},
			{
				CIDR:       *ipnet.MustParseCIDR("192.168.1.0/24"),
				HostPrefix: 28,
			},
		},
	}
}

func validOvirtPlatform() *ovirt.Platform {
	return &ovirt.Platform{
		ClusterID:       uuid.NewRandom().String(),
		StorageDomainID: uuid.NewRandom().String(),
		APIVIP:          "1.1.1.1",
		DNSVIP:          "1.1.1.2",
		IngressVIP:      "1.1.1.3",
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
				c.ObjectMeta.Name = fmt.Sprintf("test-cluster%042d", 0)
				c.BaseDomain = fmt.Sprintf("test-domain%056d.a%060d.b%060d.c%060d", 0, 0, 0, 0)
				return c
			}(),
			expectedError: `^baseDomain: Invalid value: "` + fmt.Sprintf("test-cluster%042d.test-domain%056d.a%060d.b%060d.c%060d", 0, 0, 0, 0, 0) + `": must be no more than 253 characters$`,
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
			expectedError: `^networking\.serviceNetwork\[0\]: Invalid value: "10\.0\.2\.0/24": service network must not overlap with any of the machine networks$`,
		},
		{
			name: "overlapping machine network and machine network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("13.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("13.0.2.0/24")},
				}

				return c
			}(),
			// also triggers the only-one-machine-network validation
			expectedError: `^networking\.machineNetwork\[1\]: Invalid value: "13\.0\.2\.0/24": machine network must not overlap with machine network 0$`,
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
			name: "missing machine networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = nil
				return c
			}(),
			expectedError: `^networking\.machineNetwork: Required value: at least one machine network is required$`,
		},
		{
			name: "invalid machine cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{{CIDR: *ipnet.MustParseCIDR("11.0.128.0/16")}}
				return c
			}(),
			expectedError: `^networking\.machineNetwork\[0\]: Invalid value: "11\.0\.128\.0/16": invalid network address. got 11\.0\.128\.0/16, expecting 11\.0\.0\.0/16$`,
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
			expectedError: `^networking\.clusterNetwork\[0]\.cidr: Invalid value: "10\.0\.3\.0/24": cluster network must not overlap with any of the machine networks$`,
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
					*validMachinePool("worker"),
					*validMachinePool("worker"),
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
					func() types.MachinePool {
						p := *validMachinePool("worker")
						p.Replicas = pointer.Int64Ptr(0)
						return p
					}(),
				}
				return c
			}(),
		},
		{
			name: "invalid compute",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute = []types.MachinePool{
					func() types.MachinePool {
						p := *validMachinePool("worker")
						p.Platform = types.MachinePoolPlatform{
							OpenStack: &openstack.MachinePool{},
						}
						return p
					}(),
				}
				return c
			}(),
			expectedError: `^compute\[0\]\.platform\.openstack: Invalid value: openstack\.MachinePool{.*}: cannot specify "openstack" for machine pool when cluster is using "aws"$`,
		},
		{
			name: "missing platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{}
				return c
			}(),
			expectedError: `^platform: Invalid value: "": must specify one of the platforms \(aws, azure, baremetal, gcp, none, openstack, ovirt, vsphere\)$`,
		},
		{
			name: "multiple platforms",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.Libvirt = validLibvirtPlatform()
				return c
			}(),
			expectedError: `^platform: Invalid value: "aws": must only specify a single type of platform; cannot use both "aws" and "libvirt"$`,
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
			expectedError: `^platform\.aws\.region: Required value: region must be specified$`,
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
			expectedError: `^platform: Invalid value: "libvirt": must specify one of the platforms \(aws, azure, baremetal, gcp, none, openstack, ovirt, vsphere\)$`,
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
			expectedError: `^\[platform: Invalid value: "libvirt": must specify one of the platforms \(aws, azure, baremetal, gcp, none, openstack, ovirt, vsphere\), platform\.libvirt\.uri: Invalid value: "": invalid URI "" \(no scheme\)]$`,
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
					OpenStack: validOpenStackPlatform(),
				}
				return c
			}(),
		},
		{
			name: "invalid openstack platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					OpenStack: validOpenStackPlatform(),
				}
				c.Platform.OpenStack.Cloud = ""
				return c
			}(),
			expectedError: `^platform\.openstack\.cloud: Unsupported value: "": supported values: "test-cloud"$`,
		},
		{
			name: "valid baremetal platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				return c
			}(),
		},
		{
			name: "invalid baremetal platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIP = ""
				return c
			}(),
			expectedError: `^\[platform\.baremetal\.apiVIP: Invalid value: "": "" is not a valid IP, platform\.baremetal\.apiVIP: Invalid value: "": the virtual IP is expected to be in one of the machine networks]$`,
		},
		{
			name: "baremetal API VIP not an IP",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIP = "test"
				return c
			}(),
			expectedError: `^\[platform\.baremetal\.apiVIP: Invalid value: "test": "test" is not a valid IP, platform\.baremetal\.apiVIP: Invalid value: "test": the virtual IP is expected to be in one of the machine networks]$`,
		},
		{
			name: "baremetal API VIP set to an incorrect value",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIP = "10.1.0.5"
				return c
			}(),
			expectedError: `^platform\.baremetal\.apiVIP: Invalid value: "10\.1\.0\.5": the virtual IP is expected to be in one of the machine networks$`,
		},
		{
			name: "baremetal Ingress VIP not an IP",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIP = "test"
				return c
			}(),
			expectedError: `^\[platform\.baremetal\.ingressVIP: Invalid value: "test": "test" is not a valid IP, platform\.baremetal\.ingressVIP: Invalid value: "test": the virtual IP is expected to be in one of the machine networks]$`,
		},
		{
			name: "baremetal Ingress VIP set to an incorrect value",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIP = "10.1.0.7"
				return c
			}(),
			expectedError: `^platform\.baremetal\.ingressVIP: Invalid value: "10\.1\.0\.7": the virtual IP is expected to be in one of the machine networks$`,
		}, {
			name: "valid vsphere platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				return c
			}(),
		},
		{
			name: "invalid vsphere platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.VCenter = ""
				return c
			}(),
			expectedError: `^platform\.vsphere.vCenter: Required value: must specify the name of the vCenter$`,
		},
		{
			name: "invalid vsphere folder",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.Folder = "my-folder"
				return c
			}(),
			expectedError: `^platform\.vsphere.folder: Invalid value: \"my-folder\": folder must be absolute path: expected prefix /test-datacenter/vm/$`,
		},
		{
			name: "empty proxy settings",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = ""
				c.Proxy.HTTPSProxy = ""
				c.Proxy.NoProxy = ""
				return c
			}(),
			expectedError: `^proxy: Required value: must include httpProxy or httpsProxy$`,
		},
		{
			name: "invalid HTTPProxy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://bad%20uri"
				return c
			}(),
			expectedError: `^\QHTTPProxy: Invalid value: "http://bad%20uri": parse http://bad%20uri: invalid URL escape "%20"\E$`,
		},
		{
			name: "invalid HTTPSProxy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "https://bad%20uri"
				return c
			}(),
			expectedError: `^\QHTTPSProxy: Invalid value: "https://bad%20uri": parse https://bad%20uri: invalid URL escape "%20"\E$`,
		},
		{
			name: "invalid NoProxy domain",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com, *.bad-proxy"
				return c
			}(),
			expectedError: `^\QNoProxy: Invalid value: "*.bad-proxy": must be a CIDR or domain, without wildcard characters\E$`,
		},
		{
			name: "invalid NoProxy CIDR",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com, 172.bad.CIDR.0/16"
				return c
			}(),
			expectedError: `^\QNoProxy: Invalid value: "172.bad.CIDR.0/16": must be a CIDR or domain, without wildcard characters\E$`,
		},
		{
			name: "invalid NoProxy domain & CIDR",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com, a-good-one, *.bad-proxy., another,   172.bad.CIDR.0/16, good-end"
				return c
			}(),
			expectedError: `^\Q[NoProxy: Invalid value: "*.bad-proxy.": must be a CIDR or domain, without wildcard characters, NoProxy: Invalid value: "172.bad.CIDR.0/16": must be a CIDR or domain, without wildcard characters]\E$`,
		},
		{
			name: "valid GCP platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					GCP: validGCPPlatform(),
				}
				return c
			}(),
		},
		{
			name: "invalid GCP cluster name",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					GCP: validGCPPlatform(),
				}
				c.ObjectMeta.Name = "1-invalid-cluster"
				return c
			}(),
			expectedError: `^metadata\.name: Invalid value: "1-invalid-cluster": cluster name must begin with a lower-case letter$`,
		},
		{
			name: "release image source is not canonical",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source: "ocp/release-x.y",
				}}
				return c
			}(),
			expectedError: `^imageContentSources\[0\]\.source: Invalid value: "ocp/release-x\.y": failed to parse: repository name must be canonical$`,
		},
		{
			name: "release image source's mirror is not canonical",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source:  "q.io/ocp/release-x.y",
					Mirrors: []string{"ocp/openshift-x.y"},
				}}
				return c
			}(),
			expectedError: `^imageContentSources\[0\]\.mirrors\[0\]: Invalid value: "ocp/openshift-x\.y": failed to parse: repository name must be canonical$`,
		},
		{
			name: "release image source is not repository but reference by digest",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source: "quay.io/ocp/release-x.y@sha256:397c867cc10bcc90cf05ae9b71dd3de6000535e27cb6c704d9f503879202582c",
				}}
				return c
			}(),
			expectedError: `^imageContentSources\[0\]\.source: Invalid value: "quay\.io/ocp/release-x\.y@sha256:397c867cc10bcc90cf05ae9b71dd3de6000535e27cb6c704d9f503879202582c": must be repository--not reference$`,
		},
		{
			name: "release image source is not repository but reference by tag",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source: "quay.io/ocp/release-x.y:latest",
				}}
				return c
			}(),
			expectedError: `^imageContentSources\[0\]\.source: Invalid value: "quay\.io/ocp/release-x\.y:latest": must be repository--not reference$`,
		},
		{
			name: "valid release image source",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source: "quay.io/ocp/release-x.y",
				}}
				return c
			}(),
		},
		{
			name: "invalid publishing strategy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Publish = types.PublishingStrategy("ExternalInternalDoNotCare")
				return c
			}(),
			expectedError: `^publish: Unsupported value: \"ExternalInternalDoNotCare\": supported values: \"External\", \"Internal\"`,
		},

		{
			name: "valid dual-stack configuration",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				return c
			}(),
		},
		{
			name: "valid single-stack IPv6 configuration",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validIPv6NetworkingConfig()
				return c
			}(),
		},
		{
			name: "invalid dual-stack configuration, bad platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{GCP: validGCPPlatform()}
				c.Networking = validDualStackNetworkingConfig()
				return c
			}(),
			expectedError: `Invalid value: "DualStack": dual-stack IPv4/IPv6 is not supported for this platform, specify only one type of address`,
		},
		{
			name: "invalid single-stack IPv6 configuration, bad platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{GCP: validGCPPlatform()}
				c.Networking = validIPv6NetworkingConfig()
				return c
			}(),
			expectedError: `Invalid value: "IPv6": single-stack IPv6 is not supported for this platform`,
		},
		{
			name: "invalid dual-stack configuration, bad plugin",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.NetworkType = "OpenShiftSDN"
				return c
			}(),
			expectedError: `IPv6 is not supported for this networking plugin`,
		},
		{
			name: "invalid single-stack IPv6 configuration, bad plugin",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validIPv6NetworkingConfig()
				c.Networking.NetworkType = "OpenShiftSDN"
				return c
			}(),
			expectedError: `IPv6 is not supported for this networking plugin`,
		},
		{
			name: "invalid dual-stack configuration, machine has no IPv6",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.MachineNetwork = c.Networking.MachineNetwork[1:]
				return c
			}(),
			expectedError: `Invalid value: "10.0.0.0": dual-stack IPv4/IPv6 requires an IPv6 address in this list`,
		},
		{
			name: "valid dual-stack configuration, machine has no IPv6 but is on AWS",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.MachineNetwork = c.Networking.MachineNetwork[1:]
				return c
			}(),
			expectedError: `Invalid value: "DualStack": dual-stack IPv4/IPv6 is not supported for this platform, specify only one type of address`,
		},

		{
			name: "valid ovirt platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Ovirt: validOvirtPlatform(),
				}
				return c
			}(),
		},
		{
			name: "cluster is not heteregenous",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute[0].Architecture = types.ArchitectureS390X
				return c
			}(),
			expectedError: `^compute\[0\].architecture: Invalid value: "s390x": heteregeneous multi-arch is not supported; compute pool architecture must match control plane$`,
		},
		{
			name: "cluster is not heteregenous",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute[0].Architecture = types.ArchitecturePPC64LE
				return c
			}(),
			expectedError: `^compute\[0\].architecture: Invalid value: "ppc64le": heteregeneous multi-arch is not supported; compute pool architecture must match control plane$`,
		},
		{
			name: "valid cloud credentials mode",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.CredentialsMode = types.PassthroughCredentialsMode
				return c
			}(),
		},
		{
			name: "unsupported manual cloud credentials mode",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{GCP: validGCPPlatform()}
				c.CredentialsMode = types.ManualCredentialsMode
				return c
			}(),
			expectedError: `^credentialsMode: Unsupported value: "Manual": supported values: "Mint", "Passthrough"$`,
		},
		{
			name: "invalidly set cloud credentials mode",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{BareMetal: validBareMetalPlatform()}
				c.CredentialsMode = types.PassthroughCredentialsMode
				return c
			}(),
			expectedError: `^credentialsMode: Invalid value: "Passthrough": cannot be set when using the "baremetal" platform$`,
		},
		{
			name: "bad cloud credentials mode",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.CredentialsMode = "bad-mode"
				return c
			}(),
			expectedError: `^credentialsMode: Unsupported value: "bad-mode": supported values: "Manual", "Mint", "Passthrough"$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fetcher := mock.NewMockValidValuesFetcher(mockCtrl)
			fetcher.EXPECT().GetCloudNames().Return([]string{"test-cloud"}, nil).AnyTimes()
			fetcher.EXPECT().GetNetworkNames(gomock.Any()).Return([]string{"test-network"}, nil).AnyTimes()
			fetcher.EXPECT().GetFlavorNames(gomock.Any()).Return([]string{"test-flavor"}, nil).AnyTimes()
			fetcher.EXPECT().GetNetworkExtensionsAliases(gomock.Any()).Return([]string{"trunk"}, nil).AnyTimes()
			fetcher.EXPECT().GetServiceCatalog(gomock.Any()).Return([]string{"octavia"}, nil).AnyTimes()

			err := ValidateInstallConfig(tc.installConfig, fetcher).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}
