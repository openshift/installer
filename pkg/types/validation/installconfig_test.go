package validation

import (
	"fmt"
	"net"
	"testing"

	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	utilsslice "k8s.io/utils/strings/slices"

	configv1 "github.com/openshift/api/config/v1"
	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const TechPreviewNoUpgrade = "TechPreviewNoUpgrade"

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

func validSSHKey() string {
	return "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD1+D0ns3LYRPeFK2nqOtVKBGueBQGdBLre5A+afvjaIj/QgtJuwv3rb6Uso8GMPbFlj693/b9BcV0TGxa5lC8cAGKrpxKUPvZ0WLRFLMP5HKBFf6+N4SQR9NKi7Liw8Km1GW9l+s/gMFz/ypANTg8PqvR4yglW+6jJEuKdCy/q14s9kEn4czifBzqiBw60gUiDdWbawl8yF+TxiqeKTCfw4HTeY6j1vui0ROuN2XAWgdH999rNAr1QY8BPMTjQJ5X7jeFgagq7u+snXgWycoDsn4fZP1XL91nQXLdZZgJ3T/qtjUbQt4wUuiqCu4cyN8KRoFQBtX9X7TKU8aH/Kkf+t67zS/SE0ZgvCkNr+iaqYVyHpmBoLh3AaWUYJ2bQ7fx9FvEGLcDYNkwqBED6VwuqB7nw+zGYVouGLs+2UKjfc+A1BOP0Q/2ACEkt1u5iLA+dfEC5nMMThIMNgXpjpsYLsGDKV+e9fEzrTphYtYs/XKaYlG634kGMk7wdgsHoTL0= localhost"
}

func validPowerVSPlatform() *powervs.Platform {
	return &powervs.Platform{
		Zone: "dal10",
	}
}

func validVSpherePlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenters: []vsphere.VCenter{
			{
				Server:   "test-vcenter",
				Port:     443,
				Username: "test-username",
				Password: "test-password",
				Datacenters: []string{
					"test-datacenter",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "test-east-1a",
				Region: "test-east",
				Zone:   "test-east-1a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					ResourcePool:   "/test-datacenter/host/test-cluster/Resources/test-resourcepool",
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
			{
				Name:   "test-east-2a",
				Region: "test-east",
				Zone:   "test-east-2a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
		},
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
		ClusterNetworkMTU: 0,
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
func InvalidPrimaryV6DualStackNetworkingConfig() *types.Networking {
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
			*ipnet.MustParseCIDR("172.30.0.0/16"),
			*ipnet.MustParseCIDR("ffd1::/112"),
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
			name: "invalid empty network type",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = ""
				return c
			}(),
			expectedError: `^networking.networkType: Required value: network provider type required$`,
		},
		{
			name: "invalid Kuryr network type",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = "Kuryr"
				return c
			}(),
			expectedError: `^networking\.networkType: Invalid value: "Kuryr": networkType Kuryr is not supported on OpenShift later than 4\.14$`,
		},
		{
			name: "invalid OpenShiftSDN network type",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = "OpenShiftSDN"
				return c
			}(),
			expectedError: `^networking\.networkType: Invalid value: "OpenShiftSDN": networkType OpenShiftSDN is not supported, please use OVNKubernetes$`,
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
			name: "overlapping service network and machine network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceNetwork[0] = *ipnet.MustParseCIDR("10.0.2.0/24")
				return c
			}(),
			expectedError: `^networking\.serviceNetwork\[0\]: Invalid value: "10\.0\.2\.0/24": service network must not overlap with any of the machine networks$`,
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
			name: "overlapping service network and default OVNKubernetes join networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.ServiceNetwork = []ipnet.IPNet{
					*ipnet.MustParseCIDR("100.64.2.0/24"),
					*ipnet.MustParseCIDR("fd98::/112"),
				}
				return c
			}(),
			expectedError: `^\[networking\.serviceNetwork\[0\]: Invalid value: "100\.64\.2\.0/24": must not overlap with OVNKubernetes default internal subnet 100\.64\.0\.0/16, networking\.serviceNetwork\[1\]: Invalid value: "fd98::/112": must not overlap with OVNKubernetes default internal subnet fd98::/64\]$`,
		},
		{
			name: "overlapping service network and default OVNKubernetes switch networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.ServiceNetwork = []ipnet.IPNet{
					*ipnet.MustParseCIDR("100.88.2.0/24"),
					*ipnet.MustParseCIDR("fd97::/112"),
				}
				return c
			}(),
			expectedError: `^\[networking\.serviceNetwork\[0\]: Invalid value: "100\.88\.2\.0/24": must not overlap with OVNKubernetes default transit subnet 100\.88\.0\.0/16, networking\.serviceNetwork\[1\]: Invalid value: "fd97::/112": must not overlap with OVNKubernetes default transit subnet fd97::/64\]$`,
		},
		{
			name: "overlapping service network and default OVNKubernetes masquerade networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.ServiceNetwork = []ipnet.IPNet{
					*ipnet.MustParseCIDR("169.254.2.0/24"),
					*ipnet.MustParseCIDR("fd69::/112"),
				}
				return c
			}(),
			expectedError: `^\[networking\.serviceNetwork\[0\]: Invalid value: "169\.254\.2\.0/24": must not overlap with OVNKubernetes default masquerade subnet 169\.254\.0\.0/17, networking\.serviceNetwork\[1\]: Invalid value: "fd69::/112": must not overlap with OVNKubernetes default masquerade subnet fd69::/112\]$`,
		},
		{
			name: "overlapping service network and user-provided IPv4 InternalJoinSubnet",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ServiceNetwork[0] = *ipnet.MustParseCIDR("13.0.2.0/24")
				c.Networking.OVNKubernetesConfig = &types.OVNKubernetesConfig{IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("13.0.0.0/16")}}
				return c
			}(),
			expectedError: `^networking\.ovnKubernetesConfig.ipv4.internalJoinSubnet: Invalid value: "13\.0\.0\.0/16": must not overlap with serviceNetwork 13\.0\.2\.0/24$`,
		},
		{
			name: "overlapping machine network and user-provided IPv4 InternalJoinSubnet",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.OVNKubernetesConfig = &types.OVNKubernetesConfig{IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("10.0.2.0/24")}}
				return c
			}(),
			expectedError: `^networking\.ovnKubernetesConfig.ipv4.internalJoinSubnet: Invalid value: "10\.0\.2\.0/24": must not overlap with machineNetwork 10\.0\.0\.0/16$`,
		},
		{
			name: "overlapping cluster network and user-provided IPv4 InternalJoinSubnet",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.OVNKubernetesConfig = &types.OVNKubernetesConfig{IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("192.168.0.0/16")}}
				return c
			}(),
			expectedError: `^networking\.ovnKubernetesConfig\.ipv4\.internalJoinSubnet: Invalid value: "192\.168\.0\.0/16": must not overlap with clusterNetwork 192\.168\.1\.0/24$`,
		},
		{
			name: "HTTPProxy overlapping with user-provided IPv4 Internal Join Subnet",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.HTTPProxy = "http://100.64.1.2:3030"
				c.Networking = validIPv4NetworkingConfig()
				c.Networking.OVNKubernetesConfig = &types.OVNKubernetesConfig{IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("100.64.0.0/16")}}
				return c
			}(),
			expectedError: `^proxy.httpProxy: Invalid value: "http://100.64.1.2:3030": proxy value is part of the ovn-kubernetes IPv4 InternalJoinSubnet$`,
		},
		{
			name: "invalid user-provided IPv4 InternalJoinSubnet",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.OVNKubernetesConfig = &types.OVNKubernetesConfig{IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("192.168.2.0/16")}}
				return c
			}(),
			expectedError: `^networking\.ovnKubernetesConfig\.ipv4\.internalJoinSubnet: Invalid value: "192\.168\.2\.0/16": invalid network address. got 192\.168\.2\.0/16, expecting 192\.168\.0\.0/16$`,
		},
		{
			name: "user-provided IPv4 InternalJoinSubnet too small",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.OVNKubernetesConfig = &types.OVNKubernetesConfig{IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("100.64.0.0/24")}}
				c.Networking.ClusterNetwork = []types.ClusterNetworkEntry{
					{
						CIDR:       *ipnet.MustParseCIDR("10.128.0.0/14"),
						HostPrefix: int32(23),
					},
				}
				return c
			}(),
			expectedError: `^networking\.ovnKubernetesConfig\.ipv4\.internalJoinSubnet: Invalid value: "100\.64\.0\.0/24": ipv4InternalJoinSubnet is not large enough for the maximum number of nodes which can be supported by ClusterNetwork$`,
		},
		{
			name: "valid user-provided IPv4 InternalJoinSubnet but invalid hostPrefix",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.OVNKubernetesConfig = &types.OVNKubernetesConfig{IPv4: &types.IPv4OVNKubernetesConfig{InternalJoinSubnet: ipnet.MustParseCIDR("100.64.0.0/24")}}
				c.Networking.ClusterNetwork = []types.ClusterNetworkEntry{
					{
						CIDR:       *ipnet.MustParseCIDR("10.128.0.0/24"),
						HostPrefix: int32(23),
					},
				}
				return c
			}(),
			expectedError: `^\[networking\.clusterNetwork\[0\]\.hostPrefix: Invalid value: 23: cluster network host subnetwork prefix must not be larger size than CIDR 10\.128\.0\.0/24, networking\.ovnKubernetesConfig\.ipv4\.internalJoinSubnet: Internal error: cannot determine the number of nodes supported by cluster network 0 due to invalid hostPrefix\]$`,
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
			name: "overlapping machine network and default OVNKubernetes networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("100.64.2.0/24")},
					{CIDR: *ipnet.MustParseCIDR("100.88.2.0/24")},
					{CIDR: *ipnet.MustParseCIDR("169.254.2.0/24")},
					{CIDR: *ipnet.MustParseCIDR("fd98::/48")},
					{CIDR: *ipnet.MustParseCIDR("fd97::/48")},
					{CIDR: *ipnet.MustParseCIDR("fd69::/48")},
				}
				return c
			}(),
			expectedError: `\[networking\.machineNetwork\[0\]: Invalid value: "100\.64\.2\.0/24": must not overlap with OVNKubernetes default internal subnet 100\.64\.0\.0/16, networking\.machineNetwork\[1\]: Invalid value: "100\.88\.2\.0/24": must not overlap with OVNKubernetes default transit subnet 100\.88\.0\.0/16, networking\.machineNetwork\[2\]: Invalid value: "169\.254\.2\.0/24": must not overlap with OVNKubernetes default masquerade subnet 169\.254\.0\.0/17, networking\.machineNetwork\[3\]: Invalid value: "fd98::/48": must not overlap with OVNKubernetes default internal subnet fd98::/64, networking\.machineNetwork\[4\]: Invalid value: "fd97::/48": must not overlap with OVNKubernetes default transit subnet fd97::/64, networking\.machineNetwork\[5\]: Invalid value: "fd69::/48": must not overlap with OVNKubernetes default masquerade subnet fd69::/112\]`,
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
			name: "overlapping cluster network and machine network",
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
					{CIDR: *ipnet.MustParseCIDR("12.0.0.0/16"), HostPrefix: 28},
					{CIDR: *ipnet.MustParseCIDR("12.0.3.0/24"), HostPrefix: 28},
				}
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[1]\.cidr: Invalid value: "12\.0\.3\.0/24": cluster network must not overlap with cluster network 0$`,
		},
		{
			name: "overlapping cluster network and default OVNKubernetes networks",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.ClusterNetwork = []types.ClusterNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("100.64.2.0/24"), HostPrefix: 28},
					{CIDR: *ipnet.MustParseCIDR("100.88.2.0/24"), HostPrefix: 28},
					{CIDR: *ipnet.MustParseCIDR("169.254.2.0/24"), HostPrefix: 28},
					{CIDR: *ipnet.MustParseCIDR("fd98::/48"), HostPrefix: 64},
					{CIDR: *ipnet.MustParseCIDR("fd97::/48"), HostPrefix: 64},
					{CIDR: *ipnet.MustParseCIDR("fd69::/48"), HostPrefix: 64},
				}
				return c
			}(),
			expectedError: `^\[networking\.clusterNetwork\[0\]\.cidr: Invalid value: "100\.64\.2\.0/24": must not overlap with OVNKubernetes default internal subnet 100\.64\.0\.0/16, networking\.clusterNetwork\[1\]\.cidr: Invalid value: "100\.88\.2\.0/24": must not overlap with OVNKubernetes default transit subnet 100\.88\.0\.0/16, networking\.clusterNetwork\[2\]\.cidr: Invalid value: "169\.254\.2\.0/24": must not overlap with OVNKubernetes default masquerade subnet 169\.254\.0\.0/17, networking\.clusterNetwork\[3\]\.cidr: Invalid value: "fd98::/48": must not overlap with OVNKubernetes default internal subnet fd98::/64, networking\.clusterNetwork\[4\]\.cidr: Invalid value: "fd97::/48": must not overlap with OVNKubernetes default transit subnet fd97::/64, networking\.clusterNetwork\[5\]\.cidr: Invalid value: "fd69::/48": must not overlap with OVNKubernetes default masquerade subnet fd69::/112\]$`,
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
			name: "cluster network host prefix negative",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetwork[0].HostPrefix = -23
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[0]\.hostPrefix: Invalid value: -23: hostPrefix must be positive$`,
		},
		{
			name: "multiple cluster network host prefix different",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				// Use dual-stack config to ensure the validation only applies to IPv4 CIDRs
				c.Platform = types.Platform{None: &none.Platform{}}
				c.Networking = validDualStackNetworkingConfig()
				c.Networking.ClusterNetwork = append(c.Networking.ClusterNetwork,
					types.ClusterNetworkEntry{
						CIDR:       *ipnet.MustParseCIDR("192.168.2.0/24"),
						HostPrefix: 30,
					},
				)
				return c
			}(),
			expectedError: `^networking\.clusterNetwork\[2]\.hostPrefix: Invalid value: 30: cluster network host subnetwork prefix must be the same value for IPv4 networks$`,
		},
		{
			name: "networking clusterNetworkMTU - valid high limit ovn",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = string(operv1.NetworkTypeOVNKubernetes)
				c.Networking.ClusterNetworkMTU = 8901
				fmt.Println(c.Platform.Name())
				return c
			}(),
		},
		{
			name: "networking clusterNetworkMTU - valid low limit",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = string(operv1.NetworkTypeOVNKubernetes)
				c.Networking.ClusterNetworkMTU = 1000
				return c
			}(),
		},
		{
			name: "networking clusterNetworkMTU - invalid value lower",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworkMTU = 999
				return c
			}(),
			expectedError: `^networking\.clusterNetworkMTU: Invalid value: 999: cluster network MTU is lower than the minimum value of 1000$`,
		},
		{
			name: "networking clusterNetworkMTU - invalid value ovn",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = string(operv1.NetworkTypeOVNKubernetes)
				c.Networking.ClusterNetworkMTU = 8951
				return c
			}(),
			expectedError: `^networking\.clusterNetworkMTU: Invalid value: 8951: cluster network MTU exceeds the maximum value with the network plugin OVNKubernetes of 8901$`,
		},
		{
			name: "networking clusterNetworkMTU - invalid jumbo value",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.ClusterNetworkMTU = 9002
				return c
			}(),
			expectedError: `^networking\.clusterNetworkMTU: Invalid value: 9002: cluster network MTU exceeds the maximum value of 9001$`,
		},
		{
			name: "networking clusterNetworkMTU - invalid for non-aws",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = string(operv1.NetworkTypeOVNKubernetes)
				c.Networking.ClusterNetworkMTU = 8901
				c.Platform = types.Platform{
					None: &none.Platform{},
				}
				return c
			}(),
			expectedError: `^networking\.clusterNetworkMTU: Invalid value: 8901: cluster network MTU is allowed only in AWS deployments`,
		},
		{
			name: "networking clusterNetworkMTU - unsupported network type",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking.NetworkType = string(operv1.NetworkTypeOpenShiftSDN)
				c.Networking.ClusterNetworkMTU = 8000
				return c
			}(),
			expectedError: `networking.networkType: Invalid value: "OpenShiftSDN": networkType OpenShiftSDN is not supported, please use OVNKubernetes, networking.clusterNetworkMTU: Invalid value: 8000: cluster network MTU is not valid with network plugin OpenShiftSDN`,
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
			expectedError: `^platform: Invalid value: "": must specify one of the platforms \(aws, azure, baremetal, external, gcp, ibmcloud, none, nutanix, openstack, powervs, vsphere\)$`,
		},
		{
			name: "multiple platforms",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.IBMCloud = validIBMCloudPlatform()
				return c
			}(),
			expectedError: `^platform: Invalid value: "aws": must only specify a single type of platform; cannot use both "aws" and "ibmcloud"$`,
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
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "v4.11"}
				c.Capabilities.AdditionalEnabledCapabilities = append(c.Capabilities.AdditionalEnabledCapabilities, configv1.ClusterVersionCapabilityIngress, configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityCloudControllerManager, configv1.ClusterVersionCapabilityOperatorLifecycleManager)
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
				c.Platform.VSphere.VCenters[0].Server = ""
				return c
			}(),
			expectedError: `platform\.vsphere\.vcenters\[0]\.server: Required value: must be the domain name or IP address of the vCenter(.*)`,
		},
		{
			name: "invalid vsphere folder",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.FailureDomains[0].Topology.Folder = "my-folder"
				return c
			}(),
			expectedError: `^platform\.vsphere\.failureDomains\.topology.folder: Invalid value: "my-folder": full path of folder must be provided in format /<datacenter>/vm/<folder>$`,
		},
		{
			name: "invalid vsphere resource pool",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					VSphere: validVSpherePlatform(),
				}
				c.Platform.VSphere.FailureDomains[0].Topology.ResourcePool = "my-resource-pool"
				return c
			}(),
			expectedError: `^platform\.vsphere\.failureDomains\.topology\.resourcePool: Invalid value: "my-resource-pool": full path of resource pool must be provided in format /<datacenter>/host/<cluster>/\.\.\.$`,
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
			expectedError: `^\Qproxy.noProxy: Invalid value: "good-no-proxy.com,*.bad-proxy": each element of noProxy must be a IP, CIDR or domain without wildcard characters, which is violated by element 1 "*.bad-proxy"\E$`,
		},
		{
			name: "invalid NoProxy spaces",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com, *.bad-proxy"
				return c
			}(),
			expectedError: `^\Q[proxy.noProxy: Invalid value: "good-no-proxy.com, *.bad-proxy": noProxy must not have spaces, proxy.noProxy: Invalid value: "good-no-proxy.com, *.bad-proxy": each element of noProxy must be a IP, CIDR or domain without wildcard characters, which is violated by element 1 "*.bad-proxy"]\E$`,
		},
		{
			name: "invalid NoProxy CIDR",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com,172.bad.CIDR.0/16"
				return c
			}(),
			expectedError: `^\Qproxy.noProxy: Invalid value: "good-no-proxy.com,172.bad.CIDR.0/16": each element of noProxy must be a IP, CIDR or domain without wildcard characters, which is violated by element 1 "172.bad.CIDR.0/16"\E$`,
		},
		{
			name: "invalid NoProxy domain & CIDR",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Proxy.NoProxy = "good-no-proxy.com,a-good-one,*.bad-proxy.,another,172.bad.CIDR.0/16,good-end"
				return c
			}(),
			expectedError: `^\Q[proxy.noProxy: Invalid value: "good-no-proxy.com,a-good-one,*.bad-proxy.,another,172.bad.CIDR.0/16,good-end": each element of noProxy must be a IP, CIDR or domain without wildcard characters, which is violated by element 2 "*.bad-proxy.", proxy.noProxy: Invalid value: "good-no-proxy.com,a-good-one,*.bad-proxy.,another,172.bad.CIDR.0/16,good-end": each element of noProxy must be a IP, CIDR or domain without wildcard characters, which is violated by element 4 "172.bad.CIDR.0/16"]\E$`,
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
				c.SSHKey = validSSHKey()
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
				c.SSHKey = validSSHKey()
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
				c.SSHKey = validSSHKey()
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
				c.SSHKey = validSSHKey()
				c.Platform = types.Platform{
					PowerVS: &powervs.Platform{},
				}
				return c
			}(),
			expectedError: `^\Qplatform.powervs.zone: Required value: zone must be specified\E$`,
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
				c.DeprecatedImageContentSources = []types.ImageContentSource{{
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
				c.DeprecatedImageContentSources = []types.ImageContentSource{{
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
				c.DeprecatedImageContentSources = []types.ImageContentSource{{
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
				c.DeprecatedImageContentSources = []types.ImageContentSource{{
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
				c.DeprecatedImageContentSources = []types.ImageContentSource{{
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
				c.DeprecatedImageContentSources = []types.ImageContentSource{{
					Source: "quay.io/ocp/release-x.y",
				}}
				return c
			}(),
		},
		{
			name: "release image source is not valid ImageDigestSource",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source: "ocp/release-x.y",
				}}
				return c
			}(),
			expectedError: `^imageDigestSources\[0\]\.source: Invalid value: "ocp/release-x\.y": the repository provided is invalid: a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, \'\-\' or \'\.\', and must start and end with an alphanumeric character \(e.g. \'example\.com\', regex used for validation is \'\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\(\\\.\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\)\*\'\)`,
		},
		{
			name: "release image source's mirror is not valid ImageDigestSource",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source:  "q.io/ocp/release-x.y",
					Mirrors: []string{"ocp/openshift-x.y"},
				}}
				return c
			}(),
			expectedError: `^imageDigestSources\[0\]\.mirrors\[0\]: Invalid value: "ocp/openshift-x\.y": the repository provided is invalid: a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, \'\-\' or \'\.\', and must start and end with an alphanumeric character \(e.g. \'example\.com\', regex used for validation is \'\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\(\\\.\[a\-z0\-9\]\(\[\-a\-z0\-9\]\*\[a\-z0\-9\]\)\?\)\*\'\)`,
		},
		{
			name: "release image source's mirror is valid ImageDigestSource",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source:  "q.io/ocp/release-x.y",
					Mirrors: []string{"mirror.example.com:5000"},
				}}
				return c
			}(),
		},
		{
			name: "release image source is not repository but reference by digest ImageDigestSource",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source: "quay.io/ocp/release-x.y@sha256:397c867cc10bcc90cf05ae9b71dd3de6000535e27cb6c704d9f503879202582c",
				}}
				return c
			}(),
			expectedError: `^imageDigestSources\[0\]\.source: Invalid value: "quay\.io/ocp/release-x\.y@sha256:397c867cc10bcc90cf05ae9b71dd3de6000535e27cb6c704d9f503879202582c": must be repository--not reference$`,
		},
		{
			name: "release image source is not repository but reference by tag ImageDigestSource",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source: "quay.io/ocp/release-x.y:latest",
				}}
				return c
			}(),
			expectedError: `^imageDigestSources\[0\]\.source: Invalid value: "quay\.io/ocp/release-x\.y:latest": must be repository--not reference$`,
		},
		{
			name: "valid release image source ImageDigstSourrce",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source: "quay.io/ocp/release-x.y",
				}}
				return c
			}(),
		},
		{
			name: "valid release image source ImageDigstSource with valid mirror and sourcePolicy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source:       "quay.io/ocp/release-x.y",
					Mirrors:      []string{"mirror.example.com:5000"},
					SourcePolicy: "NeverContactSource",
				}}
				return c
			}(),
		},
		{
			name: "valid release image source ImageDigstSource with no mirror and valid sourcePolicy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source:       "quay.io/ocp/release-x.y",
					SourcePolicy: "NeverContactSource",
				}}
				return c
			}(),
			expectedError: `^imageDigestSources\[0\]\.sourcePolicy: Invalid value: "NeverContactSource": sourcePolicy cannot be configured without a mirror$`,
		},
		{
			name: "valid release image source ImageDigstSource with invalid sourcePolicy",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source:       "quay.io/ocp/release-x.y",
					Mirrors:      []string{"mirror.example.com:5000"},
					SourcePolicy: "InvalidPolicy",
				}}
				return c
			}(),
			expectedError: `^imageDigestSources\[0\]\.sourcePolicy: Invalid value: "InvalidPolicy": supported values are "NeverContactSource" and "AllowContactingSource"$`,
		},
		{
			name: "error out ImageContentSources and ImageDigestSources and are set at the same time",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.DeprecatedImageContentSources = []types.ImageContentSource{{
					Source:  "q.io/ocp/source",
					Mirrors: []string{"ocp/openshift/mirror"},
				}}
				c.ImageDigestSources = []types.ImageDigestSource{{
					Source:  "q.io/ocp/source",
					Mirrors: []string{"ocp-digest/openshift/mirror"}}}
				return c
			}(),
			expectedError: `cannot set imageContentSources and imageDigestSources at the same time`,
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
				c.Platform = types.Platform{GCP: &gcp.Platform{}}
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
				c.Platform = types.Platform{Azure: validAzureStackPlatform()}
				c.Compute[0].Architecture = types.ArchitectureARM64
				return c
			}(),
			expectedError: `^compute\[0\].architecture: Invalid value: "arm64": heteregeneous multi-arch is not supported; compute pool architecture must match control plane$`,
		},
		{
			name: "aws cluster is heteregeneous",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Compute[0].Architecture = types.ArchitectureARM64
				return c
			}(),
		},
		{
			name: "gcp cluster is heteregeneous",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{GCP: validGCPPlatform()}
				c.Compute[0].Architecture = types.ArchitectureARM64
				return c
			}(),
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
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetCurrent}
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
				c.Capabilities.AdditionalEnabledCapabilities = append(c.Capabilities.AdditionalEnabledCapabilities, configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityCloudControllerManager, configv1.ClusterVersionCapabilityIngress, configv1.ClusterVersionCapabilityOperatorLifecycleManager)
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
			name: "invalid capability marketplace specified without OperatorLifecycleManager",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "None",
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{"marketplace"}}
				return c
			}(),
			expectedError: `additionalEnabledCapabilities: Invalid value: \["marketplace"\]: the marketplace capability requires the OperatorLifecycleManager capability`,
		},
		{
			name: "valid additional enabled capability specified",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "v4.11"}
				c.Capabilities.AdditionalEnabledCapabilities = append(c.Capabilities.AdditionalEnabledCapabilities, configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityOpenShiftSamples, configv1.ClusterVersionCapabilityCloudControllerManager, configv1.ClusterVersionCapabilityIngress, configv1.ClusterVersionCapabilityOperatorLifecycleManager)
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
		{
			name: "baremetal platform requires the baremetal capability",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Capabilities = &types.Capabilities{BaselineCapabilitySet: "None", AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{"marketplace"}}
				return c
			}(),
			expectedError: `additionalEnabledCapabilities: Invalid value: \["marketplace"\]: platform baremetal requires the baremetal capability`,
		},
		// VIP tests
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
			name: "apivip_v4_not_in_machinenetwork_cidr_usermanaged_loadbalancer",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = configv1.TechPreviewNoUpgrade
				c.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					{CIDR: *ipnet.MustParseCIDR("fe80::/10")},
				}
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.LoadBalancer = &configv1.BareMetalPlatformLoadBalancer{Type: configv1.LoadBalancerTypeUserManaged}
				c.Platform.BareMetal.APIVIPs = []string{"192.168.222.1"}

				return c
			}(),
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
			expectedError: "[networking.networkType: Invalid value: \"OpenShiftSDN\": networkType OpenShiftSDN is not supported, please use OVNKubernetes, platform.baremetal.ingressVIPs: Invalid value: \"10.0.0.4\": IP expected to be in one of the machine networks: ffd0::/48]",
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
			expectedError: "[networking.networkType: Invalid value: \"OpenShiftSDN\": networkType OpenShiftSDN is not supported, please use OVNKubernetes, platform.baremetal.apiVIPs: Invalid value: \"10.0.0.5\": IP expected to be in one of the machine networks: ffd0::/48]",
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
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \\[\"192.168.111.1\",\"192.168.111.2\"\\]: If two API VIPs are given, one must be an IPv4 address, the other an IPv6",
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
			expectedError: "platform.baremetal.apiVIPs: Invalid value: \\[\"fe80::1\",\"fe80::2\"\\]: If two API VIPs are given, one must be an IPv4 address, the other an IPv6",
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
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \\[\"192.168.111.4\",\"192.168.111.5\"\\]: If two Ingress VIPs are given, one must be an IPv4 address, the other an IPv6",
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
			expectedError: "platform.baremetal.ingressVIPs: Invalid value: \\[\"fe80::1\",\"fe80::2\"\\]: If two Ingress VIPs are given, one must be an IPv4 address, the other an IPv6",
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
			name: "identical_apivip_ingressvip_usermanaged_loadbalancer",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = configv1.TechPreviewNoUpgrade
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.LoadBalancer = &configv1.BareMetalPlatformLoadBalancer{Type: configv1.LoadBalancerTypeUserManaged}
				c.Platform.BareMetal.APIVIPs = []string{"fe80::1"}
				c.Platform.BareMetal.IngressVIPs = []string{"fe80::1"}

				return c
			}(),
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
			expectedError: `[platform.baremetal.apiVIPs: Invalid value: "ffd0::": clusterNetwork primary IP Family and primary IP family for the API VIP should match, platform.baremetal.apiVIPs: Invalid value: "ffd0::": machineNetwork primary IP Family and primary IP family for the API VIP should match, platform.baremetal.apiVIPs: Invalid value: "ffd0::": serviceNetwork primary IP Family and primary IP family for the API VIP should match]`,
		},
		{
			name: "baremetal API VIP set to an incorrect IP Family with invalid primary IPv6 network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = InvalidPrimaryV6DualStackNetworkingConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.APIVIPs = []string{"ffd0::"}
				return c
			}(),
			expectedError: `platform.baremetal.apiVIPs: Invalid value: "ffd0::": serviceNetwork primary IP Family and primary IP family for the API VIP should match`,
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
			expectedError: `[platform.baremetal.ingressVIPs: Invalid value: "ffd0::": clusterNetwork primary IP Family and primary IP family for the Ingress VIP should match, platform.baremetal.ingressVIPs: Invalid value: "ffd0::": machineNetwork primary IP Family and primary IP family for the Ingress VIP should match, platform.baremetal.ingressVIPs: Invalid value: "ffd0::": serviceNetwork primary IP Family and primary IP family for the Ingress VIP should match]`,
		},
		{
			name: "baremetal Ingress VIP set to an incorrect IP Family with invalid primary IPv6 network",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Networking = InvalidPrimaryV6DualStackNetworkingConfig()
				c.Platform = types.Platform{
					BareMetal: validBareMetalPlatform(),
				}
				c.Platform.BareMetal.IngressVIPs = []string{"ffd0::"}
				return c
			}(),
			expectedError: `platform.baremetal.ingressVIPs: Invalid value: "ffd0::": serviceNetwork primary IP Family and primary IP family for the Ingress VIP should match`,
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
			name: "valid custom features",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = configv1.CustomNoUpgrade
				c.FeatureGates = []string{
					"CustomFeature1=True",
					"CustomFeature2=False",
				}
				return c
			}(),
		},
		{
			name: "invalid custom features",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = configv1.CustomNoUpgrade
				c.FeatureGates = []string{
					"CustomFeature1=True",
					"CustomFeature2",
				}
				return c
			}(),
			expectedError: `featureGates\[1\]: Invalid value: "CustomFeature2": must match the format <feature-name>=<bool>`,
		},
		{
			name: "invalid custom features bool",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = configv1.CustomNoUpgrade
				c.FeatureGates = []string{
					"CustomFeature1=foo",
					"CustomFeature2=False",
				}
				return c
			}(),
			expectedError: `featureGates\[0\]: Invalid value: "CustomFeature1=foo": must match the format <feature-name>=<bool>, could not parse boolean value`,
		},
		{
			name: "custom features supplied with non-custom featureset",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = configv1.TechPreviewNoUpgrade
				c.FeatureGates = []string{
					"CustomFeature1=True",
					"CustomFeature2=False",
				}
				return c
			}(),
			expectedError: "featureGates: Forbidden: featureGates can only be used with the CustomNoUpgrade feature set",
		},
		{
			name: "valid disabled MAPI with baseline none and baremetal enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityBaremetal, configv1.ClusterVersionCapabilityIngress, configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityCloudControllerManager},
				}
				return c
			}(),
		},
		{
			name: "valid disabled MAPI capability configuration",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
				}
				c.Capabilities.AdditionalEnabledCapabilities = append(c.Capabilities.AdditionalEnabledCapabilities, configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityCloudControllerManager, configv1.ClusterVersionCapabilityIngress)
				return c
			}(),
		},
		{
			name: "valid enabled MAPI capability configuration",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityBaremetal, configv1.ClusterVersionCapabilityMachineAPI, configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityCloudControllerManager, configv1.ClusterVersionCapabilityIngress},
				}
				return c
			}(),
		},
		{
			name: "valid enabled MAPI capability configuration 2",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityMachineAPI, configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityCloudControllerManager, configv1.ClusterVersionCapabilityIngress},
				}
				return c
			}(),
		},
		{
			name: "CloudCredential is enabled in cloud",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetCurrent,
				}
				return c
			}(),
		},
		{
			name: "CloudCredential is disabled in cloud aws",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
				}
				return c
			}(),
			expectedError: "disabling CloudCredential capability available only for baremetal platforms",
		},
		{
			name: "CloudCredential is disabled in cloud gcp",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.GCP = validGCPPlatform()
				c.AWS = nil
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
				}
				return c
			}(),
			expectedError: "disabling CloudCredential capability available only for baremetal platforms",
		},
		{
			name: "CloudCredential is enabled in baremetal",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.BareMetal = validBareMetalPlatform()
				c.AWS = nil
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetCurrent,
				}
				return c
			}(),
		},
		{
			name: "CloudCredential is disabled in baremetal",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.BareMetal = validBareMetalPlatform()
				c.AWS = nil
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityBaremetal, configv1.ClusterVersionCapabilityMachineAPI, configv1.ClusterVersionCapabilityIngress},
				}
				return c
			}(),
		},
		{
			name: "CloudController can't be disabled on cloud",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
				}
				return c
			}(),
			expectedError: "disabling CloudControllerManager is only supported on the Baremetal, None, or External platform with cloudControllerManager value none",
		},
		{
			name: "valid disabled CloudController configuration none platform",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.AWS = nil
				c.Platform.None = &none.Platform{}
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityIngress},
				}
				return c
			}(),
		},
		{
			name: "valid disabled CloudController configuration platform baremetal",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.AWS = nil
				c.Platform.BareMetal = validBareMetalPlatform()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityBaremetal, configv1.ClusterVersionCapabilityMachineAPI, configv1.ClusterVersionCapabilityIngress},
				}
				return c
			}(),
		},
		{
			name: "valid disabled CloudController configuration platform External",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.AWS = nil
				c.Platform.External = &external.Platform{}
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityIngress},
				}
				return c
			}(),
		},
		{
			name: "valid disabled CloudController configuration platform External 2",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.AWS = nil
				c.Platform.External = &external.Platform{
					CloudControllerManager: external.CloudControllerManagerTypeNone,
				}
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityCloudCredential, configv1.ClusterVersionCapabilityIngress},
				}
				return c
			}(),
		},
		{
			name: "invalid disabled CloudController configuration platform External 2",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Platform.AWS = nil
				c.Platform.External = &external.Platform{
					CloudControllerManager: external.CloudControllerManagerTypeExternal,
				}
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet:         configv1.ClusterVersionCapabilitySetNone,
					AdditionalEnabledCapabilities: []configv1.ClusterVersionCapability{configv1.ClusterVersionCapabilityCloudCredential},
				}
				return c
			}(),
			expectedError: "disabling CloudControllerManager on External platform supported only with cloudControllerManager value none",
		},
		{
			name: "Ingress can't be disabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.Capabilities = &types.Capabilities{
					BaselineCapabilitySet: configv1.ClusterVersionCapabilitySetNone,
				}
				return c
			}(),
			expectedError: "the Ingress capability is required",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateInstallConfig(tc.installConfig, false).ToAggregate()
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

func TestValidateReleaseArchitecture(t *testing.T) {
	t.Run("multi arch payload is always valid", func(t *testing.T) {
		releaseArch := types.Architecture("multi")
		t.Run("for default single arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = ""
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = ""
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("for amd64 single arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("for arm64 single arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = types.ArchitectureARM64
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("for mixed arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
	})
	t.Run("unknown arch payload is always valid", func(t *testing.T) {
		releaseArch := types.Architecture("unknown")
		t.Run("for default single arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = ""
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = ""
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("for amd64 single arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("for arm64 single arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = types.ArchitectureARM64
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("for mixed arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
	})
	t.Run("amd64 arch payload", func(t *testing.T) {
		releaseArch := types.Architecture("amd64")
		t.Run("for default single arch amd64 cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = ""
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = ""
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("is valid for amd64 cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("is not valid for arm64 cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = types.ArchitectureARM64
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 2)
			assert.Regexp(t, `^\[controlPlane\.architecture: Invalid value: \"arm64\": .*compute\[0\]\.architecture: Invalid value: \"arm64\": .*\]$`, errs.ToAggregate())
		})
		t.Run("is not valid for mixed arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 1)
			assert.Regexp(t, `^compute\[0\]\.architecture: Invalid value: \"arm64\": .*$`, errs.ToAggregate())
		})
	})
	t.Run("arm64 arch payload", func(t *testing.T) {
		releaseArch := types.Architecture("arm64")
		t.Run("is valid for arm64 cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = types.ArchitectureARM64
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 0, "expected no errors")
		})
		t.Run("is not valid for default single arch amd64 cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			controlplanePool.Architecture = ""
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = ""
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Regexp(t, `^controlPlane\.architecture: Invalid value: \"amd64\": .*$`, errs.ToAggregate())
		})
		t.Run("is not valid for amd64 cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 2)
			assert.Regexp(t, `^\[controlPlane\.architecture: Invalid value: \"amd64\": .*compute\[0\]\.architecture: Invalid value: \"amd64\": .*\]$`, errs.ToAggregate())
		})
		t.Run("is not valid for mixed arch cluster", func(t *testing.T) {
			controlplanePool := validMachinePool("master")
			computePool := []types.MachinePool{*validMachinePool("worker")}
			computePool[0].Architecture = types.ArchitectureARM64
			errs := validateReleaseArchitecture(controlplanePool, computePool, releaseArch)
			assert.Len(t, errs, 1)
			assert.Regexp(t, `^controlPlane\.architecture: Invalid value: \"amd64\": .*$`, errs.ToAggregate())
		})
	})
}

func TestValidateTNF(t *testing.T) {
	cases := []struct {
		name         string
		config       *types.InstallConfig
		machinePool  *types.MachinePool
		checkCompute bool
		expected     string
	}{
		{
			config: installConfig().
				PlatformBMWithHosts().
				CpReplicas(3).
				build(),
			name:     "valid_empty_credentials_for_non_tnf",
			expected: "",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1(), c2())).
				CpReplicas(2).
				build(),
			name:     "valid_two_credentials",
			expected: "",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1().CertificateVerification(types.CertificateVerificationDisabled), c2())).
				CpReplicas(2).
				build(),
			name:     "valid_disabled_cert_verification",
			expected: "",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1().CertificateVerification("Invalid"), c2())).
				CpReplicas(2).
				build(),
			name:     "invalid_cert_verification",
			expected: "controlPlane.fencing.credentials\\[0\\]\\[CertificateVerification\\]: Invalid value: \"Invalid\": invalid certificate verification; \"Invalid\" should set to one of the following: \\['Enabled' \\(default\\), 'Disabled'\\]",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1(), c2(), c3())).
				CpReplicas(2).build(),
			name:     "invalid_number_of_credentials",
			expected: "controlPlane.fencing.credentials: Forbidden: there should be exactly two fencing credentials to support the two node cluster, instead 3 credentials were found",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolArbiter(machinePool()).
				MachinePoolCP(machinePool().Credential(c1(), c2())).
				ArbiterReplicas(1).
				CpReplicas(2).build(),
			name:     "skip_number_of_credentials_validation_for_arbiter_deployment",
			expected: "",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolArbiter(machinePool()).
				MachinePoolCP(machinePool().Credential(c1(), c2(), c3())).
				ArbiterReplicas(1).
				CpReplicas(2).build(),
			name:     "skip_number_of_credentials_validation_for_arbiter_deployment_invalid_credentials_count",
			expected: "",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1(), c2())).
				CpReplicas(3).build(),
			name:     "invalid_number_of_credentials_for_non_tnf",
			expected: "controlPlane.fencing.credentials: Forbidden: there should not be any fencing credentials configured for a non dual replica control plane \\(Two Nodes Fencing\\) cluster, instead 2 credentials were found",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(
						c1().FencingCredentialAddress("ipmi://192.168.111.1"),
						c2().FencingCredentialAddress("ipmi://192.168.111.1"))).
				CpReplicas(2).build(),
			name:     "fencing_credential_address_not_unique",
			expected: "controlPlane.fencing.credentials\\[1\\].address: Duplicate value: \"ipmi://192.168.111.1\"",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1().FencingCredentialAddress(""), c2())).
				CpReplicas(2).build(),
			name:     "fencing_credential_address_required",
			expected: "controlPlane.fencing.credentials\\[0\\].address: Required value: missing Address",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1(), c2().FencingCredentialUsername(""))).
				CpReplicas(2).build(),
			name:     "fencing_credential_username_required",
			expected: "controlPlane.fencing.credentials\\[1\\].username: Required value: missing Username",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1().FencingCredentialPassword(""), c2())).
				CpReplicas(2).build(),
			name:     "fencing_credential_password_required",
			expected: "controlPlane.fencing.credentials\\[0\\].password: Required value: missing Password",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1().HostName(""), c2())).
				CpReplicas(2).build(),
			name:     "fencing_credential_host_name_required",
			expected: "controlPlane.fencing.credentials\\[0\\].hostName: Required value: missing HostName",
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Architecture(types.ArchitectureAMD64).
					Credential(c1(), c2())).
				MachinePoolCompute(
					machinePool().Name("worker").
						Hyperthreading(types.HyperthreadingDisabled).
						Architecture(types.ArchitectureAMD64).
						Replicas(ptr.Int64(3)).
						Credential(c1())).
				CpReplicas(2).build(),
			name:         "fencing_only_valid_for_control_plane",
			checkCompute: true,
			expected:     `compute\[\d+\]\.fencing: Invalid value: \{"credentials":\[\{"hostName":"host1","username":"root","password":"password","address":"ipmi://192.168.111.1"\}\]\}: fencing is only valid for control plane`,
		},
		{
			config: installConfig().
				PlatformBMWithHosts().
				MachinePoolCP(machinePool().
					Credential(c1(), c2())).
				CpReplicas(2).
				build(),
			name:     "supported_platform_bm",
			expected: "",
		},
		{
			config: installConfig().
				PlatformExternal().
				MachinePoolCP(machinePool().
					Credential(c1(), c2())).
				CpReplicas(2).
				build(),
			name:     "supported_platform_ext",
			expected: "",
		},
		{
			config: installConfig().
				PlatformNone().
				MachinePoolCP(machinePool().
					Credential(c1(), c2())).
				CpReplicas(2).
				build(),
			name:     "supported_platform_none",
			expected: "",
		},
		{
			config: installConfig().
				PlatformAWS().
				MachinePoolCP(machinePool().
					Credential(c1(), c2())).
				CpReplicas(2).
				build(),
			name:     "unsupported_platform",
			expected: "controlPlane.fencing: Forbidden: fencing is only supported on baremetal, external or none platforms, instead aws platform was found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			var err error
			if tc.checkCompute {
				err = validateCompute(&tc.config.Platform, tc.config.ControlPlane, tc.config.Compute, field.NewPath("compute")).ToAggregate()
			} else {
				err = validateFencingCredentials(tc.config).ToAggregate()
			}

			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

type credentialBuilder struct {
	types.Credential
}

func (hb *credentialBuilder) build() *types.Credential {
	return &hb.Credential
}

func c1() *credentialBuilder {
	return &credentialBuilder{
		types.Credential{
			HostName: "host1",
			Username: "root",
			Password: "password",
			Address:  "ipmi://192.168.111.1",
		},
	}
}

func c2() *credentialBuilder {
	return &credentialBuilder{
		types.Credential{
			HostName: "host2",
			Username: "root",
			Password: "password",
			Address:  "ipmi://192.168.111.2",
		},
	}
}

func c3() *credentialBuilder {
	return &credentialBuilder{
		types.Credential{
			HostName: "host3",
			Username: "root",
			Password: "password",
			Address:  "ipmi://192.168.111.3",
		},
	}
}

func (hb *credentialBuilder) HostName(value string) *credentialBuilder {
	hb.Credential.HostName = value
	return hb
}

func (hb *credentialBuilder) FencingCredentialAddress(value string) *credentialBuilder {
	hb.Credential.Address = value
	return hb
}

func (hb *credentialBuilder) FencingCredentialUsername(value string) *credentialBuilder {
	hb.Credential.Username = value
	return hb
}

func (hb *credentialBuilder) FencingCredentialPassword(value string) *credentialBuilder {
	hb.Credential.Password = value
	return hb
}

func (hb *credentialBuilder) CertificateVerification(value types.CertificateVerificationPolicy) *credentialBuilder {
	hb.Credential.CertificateVerification = value
	return hb
}

type machinePoolBuilder struct {
	types.MachinePool
}

func (pb *machinePoolBuilder) build() *types.MachinePool {
	return &pb.MachinePool
}

func machinePool() *machinePoolBuilder {
	return &machinePoolBuilder{
		types.MachinePool{}}
}

func (pb *machinePoolBuilder) Name(name string) *machinePoolBuilder {
	pb.MachinePool.Name = name
	return pb
}

func (pb *machinePoolBuilder) Replicas(replicas *int64) *machinePoolBuilder {
	pb.MachinePool.Replicas = replicas
	return pb
}

func (pb *machinePoolBuilder) Architecture(architecture types.Architecture) *machinePoolBuilder {
	pb.MachinePool.Architecture = architecture
	return pb
}

func (pb *machinePoolBuilder) Hyperthreading(hyperthreading types.HyperthreadingMode) *machinePoolBuilder {
	pb.MachinePool.Hyperthreading = hyperthreading
	return pb
}

func (pb *machinePoolBuilder) Credential(builders ...*credentialBuilder) *machinePoolBuilder {
	pb.MachinePool.Fencing = &types.Fencing{}
	for _, builder := range builders {
		pb.MachinePool.Fencing.Credentials = append(pb.MachinePool.Fencing.Credentials, builder.build())
	}
	return pb
}

type installConfigBuilder struct {
	types.InstallConfig
}

func installConfig() *installConfigBuilder {
	return &installConfigBuilder{
		InstallConfig: types.InstallConfig{},
	}
}

func (icb *installConfigBuilder) PlatformAWS() *installConfigBuilder {
	icb.InstallConfig.Platform = types.Platform{AWS: validAWSPlatform()}
	return icb
}

func (icb *installConfigBuilder) PlatformBMWithHosts() *installConfigBuilder {
	icb.InstallConfig.Platform = types.Platform{BareMetal: validBareMetalPlatform()}
	return icb
}

func (icb *installConfigBuilder) PlatformNone() *installConfigBuilder {
	icb.InstallConfig.Platform = types.Platform{None: &none.Platform{}}
	return icb
}

func (icb *installConfigBuilder) PlatformExternal() *installConfigBuilder {
	icb.InstallConfig.Platform = types.Platform{External: &external.Platform{}}
	return icb
}

func (icb *installConfigBuilder) CpReplicas(numOfCpReplicas int64) *installConfigBuilder {
	if icb.InstallConfig.ControlPlane == nil {
		icb.InstallConfig.ControlPlane = &types.MachinePool{}
	}
	icb.InstallConfig.ControlPlane.Replicas = &numOfCpReplicas
	return icb
}

func (icb *installConfigBuilder) MachinePoolCP(builder *machinePoolBuilder) *installConfigBuilder {
	icb.InstallConfig.ControlPlane = builder.build()
	return icb
}

func (icb *installConfigBuilder) ArbiterReplicas(numOfCpReplicas int64) *installConfigBuilder {
	if icb.InstallConfig.Arbiter == nil {
		icb.InstallConfig.Arbiter = &types.MachinePool{}
	}
	icb.InstallConfig.Arbiter.Replicas = &numOfCpReplicas
	return icb
}

func (icb *installConfigBuilder) MachinePoolArbiter(builder *machinePoolBuilder) *installConfigBuilder {
	icb.InstallConfig.Arbiter = builder.build()
	return icb
}

func (icb *installConfigBuilder) MachinePoolCompute(builders ...*machinePoolBuilder) *installConfigBuilder {
	for _, builder := range builders {
		icb.InstallConfig.Compute = append(icb.InstallConfig.Compute, *builder.build())
	}
	return icb
}

func (icb *installConfigBuilder) build() *types.InstallConfig {
	return &icb.InstallConfig
}
