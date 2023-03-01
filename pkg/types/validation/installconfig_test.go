package validation

import (
	"fmt"
	"net"
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	utilsslice "k8s.io/utils/strings/slices"
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
			NoProxy:    "valid-proxy.com,172.30.0.0/16",
		},
	}
}

func validAlibabaCloudCloudPlatform() *alibabacloud.Platform {
	return &alibabacloud.Platform{
		Region:          "cn-hangzhou",
		ResourceGroupID: "test-resource-group",
	}
}

func validAWSPlatform() *aws.Platform {
	return &aws.Platform{
		Region: "us-east-1",
	}
}

func validAzureStackPlatform() *azure.Platform {
	return &azure.Platform{
		Region:                      "test-region",
		ARMEndpoint:                 "http://test-endpoint.com",
		BaseDomainResourceGroupName: "test-basedomain-rg",
		CloudName:                   azure.StackCloud,
		OutboundType:                "Loadbalancer",
	}
}

func validGCPPlatform() *gcp.Platform {
	return &gcp.Platform{
		ProjectID: "myProject",
		Region:    "us-east1",
	}
}

func validIBMCloudPlatform() *ibmcloud.Platform {
	return &ibmcloud.Platform{
		Region: "us-south",
	}
}

func validPowerVSPlatform() *powervs.Platform {
	return &powervs.Platform{
		Region: "dal",
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
		ProvisioningNetwork:          baremetal.ManagedProvisioningNetwork,
		Hosts: []*baremetal.Host{
			{
				Name:           "host1",
				Role:           "master",
				BootMACAddress: "CA:FE:CA:FE:00:00",
				BMC: baremetal.BMC{
					Username: "root",
					Password: "password",
					Address:  "ipmi://192.168.111.1",
				},
			},
			{
				Name:           "host2",
				Role:           "worker",
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
		APIVIPs:                []string{"10.0.0.5"},
		IngressVIPs:            []string{"10.0.0.4"},
	}
}

func validOpenStackPlatform() *openstack.Platform {
	return &openstack.Platform{
		Cloud:           "test-cloud",
		ExternalNetwork: "test-network",
		DefaultMachinePlatform: &openstack.MachinePool{
			FlavorName: "test-flavor",
		},
		APIVIPs:     []string{"10.0.0.5"},
		IngressVIPs: []string{"10.0.0.4"},
	}
}

func validNutanixPlatform() *nutanix.Platform {
	return &nutanix.Platform{
		PrismCentral: nutanix.PrismCentral{
			Endpoint: nutanix.PrismEndpoint{Address: "test-pc", Port: 8080},
			Username: "test-username-pc",
			Password: "test-password-pc",
		},
		PrismElements: []nutanix.PrismElement{{
			UUID:     "test-pe-uuid",
			Endpoint: nutanix.PrismEndpoint{Address: "test-pe", Port: 8081},
		}},
		SubnetUUIDs: []string{"test-subnet"},
	}
}

func validIPv4NetworkingConfig() *types.Networking {
	return &types.Networking{
		NetworkType: "OVNKubernetes",
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
			*ipnet.MustParseCIDR("ffd1::/112"),
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
				CIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
			},
			{
				CIDR: *ipnet.MustParseCIDR("ffd0::/48"),
			},
		},
		ServiceNetwork: []ipnet.IPNet{
			*ipnet.MustParseCIDR("172.30.0.0/16"),
			*ipnet.MustParseCIDR("ffd1::/112"),
		},
		ClusterNetwork: []types.ClusterNetworkEntry{
			{
				CIDR:       *ipnet.MustParseCIDR("192.168.1.0/24"),
				HostPrefix: 28,
			},
			{
				CIDR:       *ipnet.MustParseCIDR("ffd2::/48"),
				HostPrefix: 64,
			},
		},
	}
}

func validOvirtPlatform() *ovirt.Platform {
	return &ovirt.Platform{
		ClusterID:       uuid.NewRandom().String(),
		StorageDomainID: uuid.NewRandom().String(),
		APIVIPs:         []string{"10.0.1.1"},
		IngressVIPs:     []string{"10.0.1.3"},
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
			expectedError: `^metadata.name: Invalid value: "bad-name-": a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '\.', and must start and end with an alphanumeric character \(e\.g\. 'example\.com', regex used for validation is '\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\(\\\.\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\)\*'\)$`,
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
			expectedError: `^baseDomain: Invalid value: "\.bad-domain\.": a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '\.', and must start and end with an alphanumeric character \(e\.g\. 'example\.com', regex used for validation is '\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\(\\\.\[a-z0-9]\(\[-a-z0-9]\*\[a-z0-9]\)\?\)\*'\)$`,
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
			name: "cluster network host prefix unset",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = "OVNKubernetes"
				c.Networking.ClusterNetwork[0].CIDR = *ipnet.MustParseCIDR("192.168.1.0/24")
				c.Networking.ClusterNetwork[0].HostPrefix = 0
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[0]\.hostPrefix: Invalid value: 0: cluster network host subnetwork prefix must not be larger size than CIDR 192.168.1.0/24$`,
		},
		{
			name: "cluster network host prefix unset ignored",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = "HostPrefixNotRequiredPlugin"
				c.Networking.ClusterNetwork[0].CIDR = *ipnet.MustParseCIDR("192.168.1.0/24")
				return c
			}(),
			expectedError: ``,
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
			name: "missing platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{}
				return c
			}(),
			expectedError: `^platform: Invalid value: "": must specify one of the platforms \(alibabacloud, aws, azure, baremetal, gcp, ibmcloud, none, nutanix, openstack, ovirt, powervs, vsphere\)$`,
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
			expectedError: `^platform: Invalid value: "libvirt": must specify one of the platforms \(alibabacloud, aws, azure, baremetal, gcp, ibmcloud, none, nutanix, openstack, ovirt, powervs, vsphere\)$`,
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
			expectedError: `^\[platform: Invalid value: "libvirt": must specify one of the platforms \(alibabacloud, aws, azure, baremetal, gcp, ibmcloud, none, nutanix, openstack, ovirt, powervs, vsphere\), platform\.libvirt\.uri: Invalid value: "": invalid URI "" \(no scheme\)]$`,
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
			expectedError: `^platform\.vsphere\.folder: Invalid value: \"my-folder\": folder must be absolute path: expected prefix /test-datacenter/vm/$`,
		},
		{
			name: "invalid vsphere resource pool",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.ResourcePool = "my-resource-pool"
				return c
			}(),
			expectedError: `^platform\.vsphere\.resourcePool: Invalid value: \"my-resource-pool\": resourcePool must be absolute path: expected prefix /test-datacenter/host/<cluster>/Resources/$`,
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
			expectedError: `^proxy.httpProxy: Invalid value: "http://bad%20uri": parse "http://bad%20uri": invalid URL escape "%20"$`,
		},
		{
			name: "invalid HTTPProxy Schema missing",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http//baduri"
				return c
			}(),
			expectedError: `^proxy.httpProxy: Invalid value: "http//baduri": parse "http//baduri": invalid URI for request$`,
		},
		{
			name: "HTTPProxy with port overlapping with Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://192.168.1.25:3030"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpProxy: Invalid value: "http://192.168.1.25:3030": proxy value is part of the cluster networks$`,
		},
		{
			name: "overlapping HTTPProxy and Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://192.168.1.25"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpProxy: Invalid value: "http://192.168.1.25": proxy value is part of the cluster networks$`,
		},
		{
			name: "non-overlapping HTTPProxy and Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://192.169.1.25"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
		},
		{
			name: "overlapping HTTPProxy and more than one Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://192.168.1.25"
				c.Networking = validIPv4NetworkingConfig()
				c.ClusterNetwork = append(c.ClusterNetwork, []types.ClusterNetworkEntry{
					{
						CIDR:       *ipnet.MustParseCIDR("192.168.0.0/16"),
						HostPrefix: 28,
					},
				}...,
				)
				return c
			}(),
			expectedError: `^\Q[networking.clusterNetwork[1].cidr: Invalid value: "192.168.0.0/16": cluster network must not overlap with cluster network 0, proxy.httpProxy: Invalid value: "http://192.168.1.25": proxy value is part of the cluster networks]\E$`,
		},
		{
			name: "non-overlapping HTTPProxy and Service Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://172.31.0.25"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
		},
		{
			name: "HTTPProxy with port overlapping with Service Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://172.30.0.25:3030"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpProxy: Invalid value: "http://172.30.0.25:3030": proxy value is part of the service networks$`,
		},
		{
			name: "overlapping HTTPProxy and Service Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://172.30.0.25"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpProxy: Invalid value: "http://172.30.0.25": proxy value is part of the service networks$`,
		},
		{
			name: "overlapping HTTPProxy and more than one Service Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://172.30.0.25"
				c.Networking = validIPv4NetworkingConfig()
				c.ServiceNetwork = append(c.ServiceNetwork, []ipnet.IPNet{
					*ipnet.MustParseCIDR("172.30.1.0/24"),
				}...,
				)
				return c
			}(),
			expectedError: `^\Q[networking.serviceNetwork[1]: Invalid value: "172.30.1.0/24": service network must not overlap with service network 0, networking.serviceNetwork: Invalid value: "172.30.0.0/16, 172.30.1.0/24": only one service network can be specified, proxy.httpProxy: Invalid value: "http://172.30.0.25": proxy value is part of the service networks]\E$`,
		},
		{
			name: "non-overlapping HTTPSProxy and Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http://192.168.2.25"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
		},
		{
			name: "HTTPSProxy with port overlapping with Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http://192.168.1.25:3030"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpsProxy: Invalid value: "http://192.168.1.25:3030": proxy value is part of the cluster networks$`,
		},
		{
			name: "overlapping HTTPSProxy and Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http://192.168.1.25"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpsProxy: Invalid value: "http://192.168.1.25": proxy value is part of the cluster networks$`,
		},
		{
			name: "overlapping HTTPSProxy and more than one Cluster Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http://192.168.1.25"
				c.Networking = validIPv4NetworkingConfig()
				c.ClusterNetwork = append(c.ClusterNetwork, []types.ClusterNetworkEntry{
					{
						CIDR:       *ipnet.MustParseCIDR("192.168.0.0/16"),
						HostPrefix: 28,
					},
				}...,
				)
				return c
			}(),
			expectedError: `^\Q[networking.clusterNetwork[1].cidr: Invalid value: "192.168.0.0/16": cluster network must not overlap with cluster network 0, proxy.httpsProxy: Invalid value: "http://192.168.1.25": proxy value is part of the cluster networks]\E$`,
		},
		{
			name: "overlapping HTTPSProxy and Service Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http://172.30.0.25"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpsProxy: Invalid value: "http://172.30.0.25": proxy value is part of the service networks$`,
		},
		{
			name: "HTTPSProxy with port overlapping with Service Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http://172.30.0.25:3030"
				c.Networking = validIPv4NetworkingConfig()
				return c
			}(),
			expectedError: `^proxy.httpsProxy: Invalid value: "http://172.30.0.25:3030": proxy value is part of the service networks$`,
		},
		{
			name: "overlapping HTTPSProxy and more than one Service Networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http://172.30.0.25"
				c.Networking = validIPv4NetworkingConfig()
				c.ServiceNetwork = append(c.ServiceNetwork, []ipnet.IPNet{
					*ipnet.MustParseCIDR("172.30.1.0/24"),
				}...,
				)
				return c
			}(),
			expectedError: `^\Q[networking.serviceNetwork[1]: Invalid value: "172.30.1.0/24": service network must not overlap with service network 0, networking.serviceNetwork: Invalid value: "172.30.0.0/16, 172.30.1.0/24": only one service network can be specified, proxy.httpsProxy: Invalid value: "http://172.30.0.25": proxy value is part of the service networks]\E$`,
		},
		{
			name: "invalid HTTPProxy Schema different schema",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "ftp://baduri"
				return c
			}(),
			expectedError: `^proxy.httpProxy: Unsupported value: "ftp": supported values: "http"$`,
		},
		{
			name: "invalid HTTPSProxy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "https://bad%20uri"
				return c
			}(),
			expectedError: `^proxy.httpsProxy: Invalid value: "https://bad%20uri": parse "https://bad%20uri": invalid URL escape "%20"$`,
		},
		{
			name: "invalid HTTPSProxy Schema missing",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "http//baduri"
				return c
			}(),
			expectedError: `^proxy.httpsProxy: Invalid value: "http//baduri": parse "http//baduri": invalid URI for request$`,
		},
		{
			name: "invalid HTTPSProxy Schema different schema",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPSProxy = "ftp://baduri"
				return c
			}(),
			expectedError: `^proxy.httpsProxy: Unsupported value: "ftp": supported values: "http", "https"$`,
		},
		{
			name: "invalid NoProxy domain",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com,*.bad-proxy"
				return c
			}(),
			expectedError: `^\Qproxy.noProxy: Invalid value: "good-no-proxy.com,*.bad-proxy": each element of noProxy must be a CIDR or domain without wildcard characters, which is violated by element 1 "*.bad-proxy"\E$`,
		},
		{
			name: "invalid NoProxy spaces",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com, *.bad-proxy"
				return c
			}(),
			expectedError: `^\Q[proxy.noProxy: Invalid value: "good-no-proxy.com, *.bad-proxy": noProxy must not have spaces, proxy.noProxy: Invalid value: "good-no-proxy.com, *.bad-proxy": each element of noProxy must be a CIDR or domain without wildcard characters, which is violated by element 1 "*.bad-proxy"]\E$`,
		},
		{
			name: "invalid NoProxy CIDR",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com,172.bad.CIDR.0/16"
				return c
			}(),
			expectedError: `^\Qproxy.noProxy: Invalid value: "good-no-proxy.com,172.bad.CIDR.0/16": each element of noProxy must be a CIDR or domain without wildcard characters, which is violated by element 1 "172.bad.CIDR.0/16"\E$`,
		},
		{
			name: "invalid NoProxy domain & CIDR",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com,a-good-one,*.bad-proxy.,another,172.bad.CIDR.0/16,good-end"
				return c
			}(),
			expectedError: `^\Q[proxy.noProxy: Invalid value: "good-no-proxy.com,a-good-one,*.bad-proxy.,another,172.bad.CIDR.0/16,good-end": each element of noProxy must be a CIDR or domain without wildcard characters, which is violated by element 2 "*.bad-proxy.", proxy.noProxy: Invalid value: "good-no-proxy.com,a-good-one,*.bad-proxy.,another,172.bad.CIDR.0/16,good-end": each element of noProxy must be a CIDR or domain without wildcard characters, which is violated by element 4 "172.bad.CIDR.0/16"]\E$`,
		},
		{
			name: "valid * NoProxy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "*"
				return c
			}(),
		},
		{
			name: "valid alibabacloud platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					AlibabaCloud: validAlibabaCloudCloudPlatform(),
				}
				return c
			}(),
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
			name: "valid ibmcloud platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					IBMCloud: validIBMCloudPlatform(),
				}
				return c
			}(),
		},
		{
			name: "invalid ibmcloud platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					IBMCloud: &ibmcloud.Platform{},
				}
				return c
			}(),
			expectedError: `^\Qplatform.ibmcloud.region: Required value: region must be specified\E$`,
		},
		{
			name: "valid powervs platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					PowerVS: validPowerVSPlatform(),
				}
				return c
			}(),
		},
		{
			name: "valid powervs platform manual credential mod",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					PowerVS: validPowerVSPlatform(),
				}
				c.CredentialsMode = types.ManualCredentialsMode
				return c
			}(),
		},
		{
			name: "invalid powervs platform mint credential mod",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					PowerVS: validPowerVSPlatform(),
				}
				c.CredentialsMode = types.MintCredentialsMode
				return c
			}(),
			expectedError: `^credentialsMode: Unsupported value: "Mint": supported values: "Manual"$`,
		},
		{
			name: "invalid powervs platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					PowerVS: &powervs.Platform{},
				}
				return c
			}(),
			expectedError: `^\Qplatform.powervs.region: Required value: region must be specified\E$`,
		},
		{
			name: "valid azurestack platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Azure: validAzureStackPlatform(),
				}
				return c
			}(),
		},
		{
			name: "invalid azurestack platform mint credentials mod",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Azure: validAzureStackPlatform(),
				}
				c.CredentialsMode = types.MintCredentialsMode
				return c
			}(),
			expectedError: `^credentialsMode: Unsupported value: "Mint": supported values: "Manual"$`,
		},
		{
			name: "release image source is not valid",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source: "ocp/release-x.y",
				}}
				return c
			}(),
			expectedError: `^imageContentSources\[0\]\.source: Invalid value: "ocp/release-x\.y": the repository provided is invalid: a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, \'\-\' or \'\.\', and must start and end with an alphanumeric character \(e.g. \'example\.com\', regex used for validation is \'\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\(\\\.\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\)\*\'\)`,
		},
		{
			name: "release image source's mirror is not valid",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source:  "q.io/ocp/release-x.y",
					Mirrors: []string{"ocp/openshift-x.y"},
				}}
				return c
			}(),
			expectedError: `^imageContentSources\[0\]\.mirrors\[0\]: Invalid value: "ocp/openshift-x\.y": the repository provided is invalid: a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, \'\-\' or \'\.\', and must start and end with an alphanumeric character \(e.g. \'example\.com\', regex used for validation is \'\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\(\\\.\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\)\*\'\)`,
		},
		{
			name: "release image source's mirror is valid",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageContentSources = []types.ImageContentSource{{
					Source:  "q.io/ocp/release-x.y",
					Mirrors: []string{"mirror.example.com:5000"},
				}}
				return c
			}(),
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
				c.Networking.MachineNetwork = c.Networking.MachineNetwork[:1]
				return c
			}(),
			expectedError: `Invalid value: "10.0.0.0/16": dual-stack IPv4/IPv6 requires an IPv6 network in this list`,
		},
		{
			name: "invalid dual-stack configuration, IPv6-primary",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.ServiceNetwork = []ipnet.IPNet{
					c.Networking.ServiceNetwork[1],
					c.Networking.ServiceNetwork[0],
				}
				return c
			}(),
			expectedError: `Invalid value: "ffd1::/112, 172.30.0.0/16": IPv4 addresses must be listed before IPv6 addresses`,
		},
		{
			name: "valid dual-stack configuration with mixed-order clusterNetworks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.ClusterNetwork = append(c.Networking.ClusterNetwork,
					types.ClusterNetworkEntry{
						CIDR:       *ipnet.MustParseCIDR("192.168.2.0/24"),
						HostPrefix: 28,
					},
				)
				// ClusterNetwork is now "IPv4, IPv6, IPv4", which is allowed
				return c
			}(),
		},
		{
			name: "invalid IPv6 hostprefix",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validIPv6NetworkingConfig()
				c.Networking.ClusterNetwork[0].HostPrefix = 72
				return c
			}(),
			expectedError: `Invalid value: 72: cluster network host subnetwork prefix must be 64 for IPv6 networks`,
		},
		{
			name: "invalid IPv6 service network size",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validIPv6NetworkingConfig()
				c.Networking.ServiceNetwork[0] = *ipnet.MustParseCIDR("ffd1::/48")
				return c
			}(),
			expectedError: `Invalid value: "ffd1::/48": subnet size for IPv6 service network should be /112`,
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
			name: "architecture is not supported",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute[0].Architecture = types.ArchitectureS390X
				c.ControlPlane.Architecture = types.ArchitectureS390X
				return c
			}(),
			expectedError: `[controlPlane.architecture: Unsupported value: "s390x": supported values: "amd64", "arm64", compute\[0\].architecture: Unsupported value: "s390x": supported values: "amd64", "arm64"]`,
		},
		{
			name: "architecture is not supported",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute[0].Architecture = types.ArchitecturePPC64LE
				c.ControlPlane.Architecture = types.ArchitecturePPC64LE
				return c
			}(),
			expectedError: `[controlPlane.architecture: Unsupported value: "ppc64le": supported values: "amd64", "arm64", compute\[0\].architecture: Unsupported value: "ppc64le": supported values: "amd64", "arm64"]`,
		},
		{
			name: "cluster is not heteregenous",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute[0].Architecture = types.ArchitectureARM64
				return c
			}(),
			expectedError: `^compute\[0\].architecture: Invalid value: "arm64": heteregeneous multi-arch is not supported; compute pool architecture must match control plane$`,
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
		{
			name: "allowed docker bridge with non-libvirt",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{{CIDR: *ipnet.MustParseCIDR("172.17.64.0/18")}}
				return c
			}(),
			expectedError: ``,
		},
		{
			name: "docker bridge not allowed with libvirt",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{Libvirt: validLibvirtPlatform()}
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{{CIDR: *ipnet.MustParseCIDR("172.17.64.0/18")}}
				return c
			}(),
			expectedError: `\Q[networking.machineNewtork[0]: Invalid value: "172.17.64.0/18": overlaps with default Docker Bridge subnet, platform: Invalid value: "libvirt": must specify one of the platforms (\E.*\Q)]\E`,
		},
		{
			name: "publish internal for non-cloud platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{VSphere: validVSpherePlatform()}
				c.Publish = types.InternalPublishingStrategy
				return c
			}(),
			expectedError: `publish: Invalid value: "Internal": Internal publish strategy is not supported on "vsphere" platform`,
		},
		{
			name: "publish internal for cloud platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{GCP: validGCPPlatform()}
				c.Publish = types.InternalPublishingStrategy
				return c
			}(),
		}, {
			name: "valid nutanix platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Nutanix: validNutanixPlatform(),
				}
				return c
			}(),
		}, {
			name: "invalid nutanix platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Nutanix: validNutanixPlatform(),
				}
				c.Platform.Nutanix.PrismCentral.Endpoint.Address = ""
				return c
			}(),
			expectedError: `^platform\.nutanix\.prismCentral\.endpoint\.address: Required value: must specify the Prism Central endpoint address$`,
		},
		{
			name: "invalid credentials mode for nutanix",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Nutanix: validNutanixPlatform(),
				}
				c.CredentialsMode = types.PassthroughCredentialsMode
				return c
			}(),
			expectedError: `credentialsMode: Unsupported value: "Passthrough": supported values: "Manual"$`,
		},
		{
			name: "valid credentials mode for nutanix",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Nutanix: validNutanixPlatform(),
				}
				c.CredentialsMode = types.ManualCredentialsMode
				return c
			}(),
		},
		{
			name: "valid baseline capability set",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "v4.11"}
				return c
			}(),
		},
		{
			name: "invalid empty string baseline capability set",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: ""}
				return c
			}(),
			expectedError: `capabilities.baselineCapabilitySet: Unsupported value: "": supported values: .*`,
		},
		{
			name: "invalid baseline capability set specified",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "vNotValid"}
				return c
			}(),
			expectedError: `capabilities.baselineCapabilitySet: Unsupported value: "vNotValid": supported values: .*`,
		},
		{
			name: "valid additional enabled capability specified",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "v4.11",
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{"openshift-samples"}}
				return c
			}(),
		},
		{
			name: "invalid empty additional enabled capability specified",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "v4.11",
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{""}}
				return c
			}(),
			expectedError: `capabilities.additionalEnabledCapabilities\[0\]: Unsupported value: "": supported values: .*`,
		},
		{
			name: "invalid additional enabled capability specified",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "v4.11",
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{"not-valid"}}
				return c
			}(),
			expectedError: `capabilities.additionalEnabledCapabilities\[0\]: Unsupported value: "not-valid": supported values: .*`,
		},
		//VIP tests
		{
			name: "apivip_v4_not_in_machinenetwork_cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("fe80::/10")},
				}
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"192.168.222.1"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"192.168.222.1\": IP expected to be in one of the machine networks: 10.0.0.0/16,fe80::/10",
		},
		{
			name: "apivip_v6_not_in_machinenetwork_cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("fe80::/10")},
				}
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"2001::1"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"2001::1\": IP expected to be in one of the machine networks: 10.0.0.0/16,fe80::/10",
		},
		{
			name: "apivips_v6_on_openshiftsdn",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = validIPv6NetworkingConfig()
				c.Networking.NetworkType = string(operv1.NetworkTypeOpenShiftSDN)

				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"ffd0::1"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"ffd0::1\": IPv6 is not supported on OpenShiftSDN",
		},
		{
			name: "ingressvips_v6_on_openshiftsdn",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = validIPv6NetworkingConfig()
				c.Networking.NetworkType = string(operv1.NetworkTypeOpenShiftSDN)

				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"ffd0::1"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \"ffd0::1\": IPv6 is not supported on OpenShiftSDN",
		},
		{
			name: "too_many_apivips",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"fe80::1", "fe80::2", "fe80::3"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Too many: 3: must have at most 2 items",
		},
		{
			name: "invalid_apivip",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{""}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"\": \"\" is not a valid IP",
		},
		{
			name: "invalid_apivip_2",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"123.456.789.000"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"123.456.789.000\": \"123.456.789.000\" is not a valid IP",
		},
		{
			name: "invalid_apivip_format",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"foobar"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "invalid_apivip_format_one_of_many",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"192.168.1.0", "foobar"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "invalid_apivips_both_ipv4",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"192.168.111.1", "192.168.111.2"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \\[\\]string\\{\"192.168.111.1\", \"192.168.111.2\"\\}: If two API VIPs are given, one must be an IPv4 address, the other an IPv6",
		},
		{
			name: "invalid_apis_both_ipv6",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"fe80::1", "fe80::2"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \\[\\]string\\{\"fe80::1\", \"fe80::2\"\\}: If two API VIPs are given, one must be an IPv4 address, the other an IPv6",
		},
		{
			name: "ingressvip_v4_not_in_machinenetwork_cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("fe80::/10")},
				}
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"192.168.222.4"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \"192.168.222.4\": IP expected to be in one of the machine networks: 10.0.0.0/16,fe80::/10",
		},
		{
			name: "ingressvip_v6_not_in_machinenetwork_cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("fe80::/10")},
				}
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"2001::1"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \"2001::1\": IP expected to be in one of the machine networks: 10.0.0.0/16,fe80::/10",
		},
		{
			name: "vsphere_ingressvip_v4_not_in_machinenetwork_cidr",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("fe80::/10")},
				}
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.IngressVIPs = []string{"192.168.222.4"}
				c.Platform.VSphere.APIVIPs = []string{"192.168.1.0"}

				return c
			}(),
		},
		{
			name: "too_many_ingressvips",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"fe80::1", "fe80::2", "fe80::3"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Too many: 3: must have at most 2 items",
		},
		{
			name: "invalid_ingressvip",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{""}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \"\": \"\" is not a valid IP",
		},
		{
			name: "invalid_ingressvip_format",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"foobar"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "invalid_ingressvip_format_one_of_many",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"192.1.1.1", "foobar"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "invalid_ingressvips_both_ipv4",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"192.168.111.4", "192.168.111.5"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \\[\\]string\\{\"192.168.111.4\", \"192.168.111.5\"\\}: If two Ingress VIPs are given, one must be an IPv4 address, the other an IPv6",
		},
		{
			name: "invalid_ingressvips_both_ipv6",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"fe80::1", "fe80::2"}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \\[\\]string\\{\"fe80::1\", \"fe80::2\"\\}: If two Ingress VIPs are given, one must be an IPv4 address, the other an IPv6",
		},
		{
			name: "identical_apivip_ingressvip",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"fe80::1"}
				c.Platform.BareMetal.IngressVIPs = []string{"fe80::1"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"fe80::1\": VIP for API must not be one of the Ingress VIPs",
		},
		{
			name: "identical_apivips_ingressvips_multiple_ips",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"fe80::1", "192.1.2.3"}
				c.Platform.BareMetal.IngressVIPs = []string{"fe80::1", "192.1.2.4"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"fe80::1\": VIP for API must not be one of the Ingress VIPs",
		},
		{
			name: "apivip_ingressvip_are_synonyms",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"2001:db8::5"}
				c.Platform.BareMetal.IngressVIPs = []string{"2001:db8:0::5"}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \"2001:db8::5\": VIP for API must not be one of the Ingress VIPs",
		},
		{
			name: "empty_api_vip_fields",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.DeprecatedAPIVIP = ""
				c.Platform.BareMetal.APIVIPs = []string{}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Required value: must specify at least one VIP for the API",
		},
		{
			name: "empty_ingress_vip_fields",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.DeprecatedIngressVIP = ""
				c.Platform.BareMetal.IngressVIPs = []string{}

				return c
			}(),
			expectedError: "platform.baremetal.ingressVIPs: Required value: must specify at least one VIP for the Ingress",
		},
		{
			name: "baremetal API VIP set to an incorrect IP Family",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = validDualStackNetworkingConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"ffd0::"}
				return c
			}(),
			expectedError: `platform.baremetal.apiVIPs: Invalid value: "ffd0::": VIP for the API must be of the same IP family with machine network's primary IP Family for dual-stack IPv4/IPv6`,
		},
		{
			name: "baremetal Ingress VIP set to an incorrect IP Family",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = validDualStackNetworkingConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"ffd0::"}
				return c
			}(),
			expectedError: `platform.baremetal.ingressVIPs: Invalid value: "ffd0::": VIP for the Ingress must be of the same IP family with machine network's primary IP Family for dual-stack IPv4/IPv6`,
		},
		{
			name: "should validate vips on baremetal (required)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.DeprecatedAPIVIP = ""
				c.Platform.BareMetal.APIVIPs = []string{}

				return c
			}(),
			expectedError: "platform.baremetal.apiVIPs: Required value: must specify at least one VIP for the API",
		},
		{
			name: "should validate vips on OpenStack (vips are required on openstack)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					OpenStack: validOpenStackPlatform(),
				}
				c.Platform.OpenStack.DeprecatedAPIVIP = ""
				c.Platform.OpenStack.APIVIPs = []string{}

				return c
			}(),
			expectedError: "platform.openstack.apiVIPs: Required value: must specify at least one VIP for the API",
		},
		// {
		// 	name: "should not validate vips on OpenStack if not set (vips are not required on openstack)",
		// 	installConfig: func() *types.InstallConfig {
		// 		c := validInstallConfig()
		// 		c.Platform = types.Platform{
		// 			OpenStack: validOpenStackPlatform(),
		// 		}
		// 		c.Platform.OpenStack.DeprecatedAPIVIP = ""
		// 		c.Platform.OpenStack.APIVIPs = []string{}

		// 		return c
		// 	}(),
		// },
		{
			name: "should validate vips on OpenStack if set (vips are not required on openstack)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					OpenStack: validOpenStackPlatform(),
				}
				c.Platform.OpenStack.APIVIPs = []string{"foobar"}

				return c
			}(),
			expectedError: "platform.openstack.apiVIPs: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "should not validate vips on VSphere if not set (vips are not required on VSphere)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.DeprecatedAPIVIP = ""
				c.Platform.VSphere.APIVIPs = []string{}

				return c
			}(),
		},
		{
			name: "should validate vips on VSphere if set (vips are not required on VSphere)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.APIVIPs = []string{"foobar"}

				return c
			}(),
			expectedError: "platform.vsphere.apiVIPs: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "should not validate vips on Nutanix if not set (vips are not required on Nutanix)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Nutanix: validNutanixPlatform(),
				}
				c.Platform.Nutanix.DeprecatedAPIVIP = ""
				c.Platform.Nutanix.APIVIPs = []string{}

				return c
			}(),
		},
		{
			name: "should validate vips on Nutanix if set (vips are not required on Nutanix)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Nutanix: validNutanixPlatform(),
				}
				c.Platform.Nutanix.APIVIPs = []string{"foobar"}

				return c
			}(),
			expectedError: "platform.nutanix.apiVIPs: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "should return error on missing vips on Ovirt if not set (vips are required on Ovirt)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Ovirt: validOvirtPlatform(),
				}
				c.Platform.Ovirt.DeprecatedAPIVIP = ""
				c.Platform.Ovirt.APIVIPs = []string{}

				return c
			}(),
			expectedError: "platform.ovirt.api_vips: Required value: must specify at least one VIP for the API",
		},
		{
			name: "should validate vips on Ovirt (vips are required on Ovirt)",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					Ovirt: validOvirtPlatform(),
				}
				c.Platform.Ovirt.APIVIPs = []string{"foobar"}

				return c
			}(),
			expectedError: "platform.ovirt.api_vips: Invalid value: \"foobar\": \"foobar\" is not a valid IP",
		},
		{
			name: "should return error if only API VIP is set",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.APIVIPs = []string{"10.0.0.1"}
				c.Platform.VSphere.IngressVIPs = []string{}

				return c
			}(),
			expectedError: "platform.vsphere.ingressVIPs: Required value: must specify VIP for ingress, when VIP for API is set",
		},
		{
			name: "should return error if only Ingress VIP is set",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.APIVIPs = []string{}
				c.Platform.VSphere.IngressVIPs = []string{"10.0.0.1"}

				return c
			}(),
			expectedError: "platform.vsphere.apiVIPs: Required value: must specify VIP for API, when VIP for ingress is set",
		},
		{
			name: "GCP Create Firewall Rules should return error if used WITHOUT tech preview when not enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					GCP: validGCPPlatform(),
				}
				c.GCP.CreateFirewallRules = gcp.CreateFirewallRulesDisabled

				return c
			}(),
			expectedError: "platform.gcp.createFirewallRules: Forbidden: the TechPreviewNoUpgrade feature set must be enabled to use this field",
		},
		{
			name: "GCP BYO PUBLIC DNS SHOULD return error if used WITHOUT tech preview",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					GCP: validGCPPlatform(),
				}
				c.Platform.GCP.PublicDNSZone = &gcp.DNSZone{
					ID:        "myZone",
					ProjectID: "myProject",
				}

				return c
			}(),
			expectedError: "platform.gcp.publicDNSZone.projectID: Forbidden: the TechPreviewNoUpgrade feature set must be enabled to use this field",
		},
		{
			name: "GCP BYO PRIVATE DNS SHOULD return error if used WITHOUT tech preview",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					GCP: validGCPPlatform(),
				}
				c.Platform.GCP.PrivateDNSZone = &gcp.DNSZone{
					ID:        "myZone",
					ProjectID: "myProject",
				}

				return c
			}(),
			expectedError: "platform.gcp.privateDNSZone.projectID: Forbidden: the TechPreviewNoUpgrade feature set must be enabled to use this field",
		},
		{
			name: "GCP BYO DNS should NOT return error if used WITH tech preview",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					GCP: validGCPPlatform(),
				}
				c.Platform.GCP.PublicDNSZone = &gcp.DNSZone{
					ID:        "myZone",
					ProjectID: "myProject",
				}
				c.Platform.GCP.PrivateDNSZone = &gcp.DNSZone{
					ID:        "myZone",
					ProjectID: "myProject",
				}
				c.FeatureSet = "TechPreviewNoUpgrade"

				return c
			}(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateInstallConfig(tc.installConfig).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}

func Test_ensureIPv4IsFirstInDualStackSlice(t *testing.T) {
	tests := []struct {
		name    string
		vips    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "should switch VIPs",
			vips:    []string{"fe80::0", "192.168.1.1"},
			want:    []string{"192.168.1.1", "fe80::0"},
			wantErr: false,
		},
		{
			name:    "should do nothing on single stack",
			vips:    []string{"192.168.1.1"},
			want:    []string{"192.168.1.1"},
			wantErr: false,
		},
		{
			name:    "should do nothing on correct order",
			vips:    []string{"192.168.1.1", "fe80::0"},
			want:    []string{"192.168.1.1", "fe80::0"},
			wantErr: false,
		},
		{
			name:    "return error on invalid number of vips",
			vips:    []string{"192.168.1.1", "fe80::0", "192.168.1.1"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureIPv4IsFirstInDualStackSlice(&tt.vips, field.NewPath("test")); (len(err) > 0) != tt.wantErr {
				t.Errorf("ensureIPv4IsFirstInDualStackSlice() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !utilsslice.Equal(tt.vips, tt.want) && len(tt.vips) == 2 {
				t.Errorf("ensureIPv4IsFirstInDualStackSlice() changed to %v, expected %v", tt.vips, tt.want)
			}
		})
	}
}
