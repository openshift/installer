package manifests

import (
	"net"

	"github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	agenttypes "github.com/openshift/installer/pkg/types/agent"
	"github.com/openshift/installer/pkg/types/baremetal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

var (
	// TestSSHKey provides a ssh key for unit tests
	TestSSHKey = `|
	ssh-rsa AAAAB3NzaC1y1LJe3zew1ghc= root@localhost.localdomain`
	// TestSecret provides a ssh key for unit tests
	TestSecret = `'{"auths":{"cloud.openshift.com":{"auth":"b3BlUTA=","email":"test@redhat.com"}}}`
)

// GetValidOptionalInstallConfig returns a valid optional install config
func getValidOptionalInstallConfig() *agent.OptionalInstallConfig {
	_, newCidr, _ := net.ParseCIDR("192.168.111.0/24")
	_, machineNetCidr, _ := net.ParseCIDR("10.10.11.0/24")

	return &agent.OptionalInstallConfig{
		InstallConfig: installconfig.InstallConfig{
			Config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ocp-edge-cluster-0",
					Namespace: "cluster-0",
				},
				BaseDomain: "testing.com",
				PullSecret: TestSecret,
				SSHKey:     TestSSHKey,
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64Ptr(3),
					Platform: types.MachinePoolPlatform{},
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker-machine-pool-1",
						Replicas: pointer.Int64Ptr(2),
					},
					{
						Name:     "worker-machine-pool-2",
						Replicas: pointer.Int64Ptr(3),
					},
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{
							CIDR: ipnet.IPNet{IPNet: *machineNetCidr},
						},
					},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       ipnet.IPNet{IPNet: *newCidr},
							HostPrefix: 23,
						},
					},
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
					NetworkType: "OVNKubernetes",
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						APIVIPs:     []string{"192.168.122.10"},
						IngressVIPs: []string{"192.168.122.11"},
					},
				},
			},
		},
		Supplied: true,
	}
}

// GetValidOptionalInstallConfigDualStack returns a valid optional install config for dual stack
func getValidOptionalInstallConfigDualStack() *agent.OptionalInstallConfig {
	_, newCidr, _ := net.ParseCIDR("192.168.111.0/24")
	_, newCidrIPv6, _ := net.ParseCIDR("2001:db8:1111:2222::/64")
	_, machineNetCidr, _ := net.ParseCIDR("10.10.11.0/24")
	_, machineNetCidrIPv6, _ := net.ParseCIDR("2001:db8:5dd8:c956::/64")

	return &agent.OptionalInstallConfig{
		InstallConfig: installconfig.InstallConfig{
			Config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ocp-edge-cluster-0",
					Namespace: "cluster-0",
				},
				BaseDomain: "testing.com",
				PullSecret: TestSecret,
				SSHKey:     TestSSHKey,
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64Ptr(3),
					Platform: types.MachinePoolPlatform{},
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker-machine-pool-1",
						Replicas: pointer.Int64Ptr(2),
					},
					{
						Name:     "worker-machine-pool-2",
						Replicas: pointer.Int64Ptr(3),
					},
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{
							CIDR: ipnet.IPNet{IPNet: *machineNetCidr},
						},
						{
							CIDR: ipnet.IPNet{IPNet: *machineNetCidrIPv6},
						},
					},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       ipnet.IPNet{IPNet: *newCidr},
							HostPrefix: 23,
						},
						{
							CIDR:       ipnet.IPNet{IPNet: *newCidrIPv6},
							HostPrefix: 64,
						},
					},
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"), *ipnet.MustParseCIDR("fd02::/112"),
					},
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						APIVIPs:     []string{"192.168.122.10"},
						IngressVIPs: []string{"192.168.122.11"},
					},
				},
			},
		},
		Supplied: true,
	}
}

func getValidAgentConfig() *agentconfig.AgentConfig {
	return &agentconfig.AgentConfig{
		Config: &agenttypes.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ocp-edge-cluster-0",
				Namespace: "cluster-0",
			},
			RendezvousIP: "192.168.122.2",
			Hosts: []agenttypes.Host{
				{
					Hostname: "control-0.example.org",
					Role:     "master",
					RootDeviceHints: baremetal.RootDeviceHints{
						DeviceName:         "/dev/sda",
						HCTL:               "hctl-value",
						Model:              "model-value",
						Vendor:             "vendor-value",
						SerialNumber:       "serial-number-value",
						MinSizeGigabytes:   20,
						WWN:                "wwn-value",
						WWNWithExtension:   "wwn-with-extension-value",
						WWNVendorExtension: "wwn-vendor-extension-value",
						Rotational:         new(bool),
					},
					Interfaces: []*v1beta1.Interface{
						{
							Name:       "enp2s0",
							MacAddress: "98:af:65:a5:8d:01",
						},
						{
							Name:       "enp3s1",
							MacAddress: "28:d2:44:d2:b2:1a",
						},
					},
					NetworkConfig: v1beta1.NetConfig{
						Raw: unmarshalJSON([]byte("interfaces:")),
					},
				},
				{
					Hostname: "control-1.example.org",
					Role:     "master",
					RootDeviceHints: baremetal.RootDeviceHints{
						DeviceName:         "/dev/sdb",
						HCTL:               "hctl-value",
						Model:              "model-value",
						Vendor:             "vendor-value",
						SerialNumber:       "serial-number-value",
						MinSizeGigabytes:   40,
						WWN:                "wwn-value",
						WWNWithExtension:   "wwn-with-extension-value",
						WWNVendorExtension: "wwn-vendor-extension-value",
						Rotational:         new(bool),
					},
					Interfaces: []*v1beta1.Interface{
						{
							Name:       "enp2t0",
							MacAddress: "98:af:65:a5:8d:02",
						},
					},
					NetworkConfig: v1beta1.NetConfig{
						Raw: unmarshalJSON([]byte("interfaces:")),
					},
				},
				{
					Hostname: "control-2.example.org",
					Role:     "master",
					RootDeviceHints: baremetal.RootDeviceHints{
						DeviceName:         "/dev/sdc",
						HCTL:               "hctl-value",
						Model:              "model-value",
						Vendor:             "vendor-value",
						SerialNumber:       "serial-number-value",
						MinSizeGigabytes:   60,
						WWN:                "wwn-value",
						WWNWithExtension:   "wwn-with-extension-value",
						WWNVendorExtension: "wwn-vendor-extension-value",
						Rotational:         new(bool),
					},
					Interfaces: []*v1beta1.Interface{
						{
							Name:       "enp2u0",
							MacAddress: "98:af:65:a5:8d:03",
						},
					},
					NetworkConfig: v1beta1.NetConfig{
						Raw: unmarshalJSON([]byte("interfaces:")),
					},
				},
			},
		},
	}
}

func getValidDHCPAgentConfigNoHosts() *agentconfig.AgentConfig {
	return &agentconfig.AgentConfig{
		Config: &agenttypes.Config{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AgentConfig",
				APIVersion: "v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ocp-edge-cluster-0",
				Namespace: "cluster-0",
			},
			RendezvousIP: "192.168.122.2",
		},
	}
}

func getValidDHCPAgentConfigWithSomeHostsWithoutNetworkConfig() *agentconfig.AgentConfig {
	return &agentconfig.AgentConfig{
		Config: &agenttypes.Config{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AgentConfig",
				APIVersion: "v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ocp-edge-cluster-0",
				Namespace: "cluster-0",
			},
			RendezvousIP: "192.168.122.2",
			Hosts: []agenttypes.Host{
				{
					Hostname: "control-0.example.org",
					Role:     "master",
					RootDeviceHints: baremetal.RootDeviceHints{
						DeviceName:         "/dev/sda",
						HCTL:               "hctl-value",
						Model:              "model-value",
						Vendor:             "vendor-value",
						SerialNumber:       "serial-number-value",
						MinSizeGigabytes:   20,
						WWN:                "wwn-value",
						WWNWithExtension:   "wwn-with-extension-value",
						WWNVendorExtension: "wwn-vendor-extension-value",
						Rotational:         new(bool),
					},
					Interfaces: []*v1beta1.Interface{
						{
							Name:       "enp2t0",
							MacAddress: "98:af:65:a5:8d:02",
						},
					},
					NetworkConfig: v1beta1.NetConfig{
						Raw: unmarshalJSON([]byte("interfaces:")),
					},
				},
				{
					Hostname: "control-1.example.org",
					Role:     "master",
					RootDeviceHints: baremetal.RootDeviceHints{
						DeviceName:         "/dev/sdb",
						HCTL:               "hctl-value",
						Model:              "model-value",
						Vendor:             "vendor-value",
						SerialNumber:       "serial-number-value",
						MinSizeGigabytes:   40,
						WWN:                "wwn-value",
						WWNWithExtension:   "wwn-with-extension-value",
						WWNVendorExtension: "wwn-vendor-extension-value",
						Rotational:         new(bool),
					},
					Interfaces: []*v1beta1.Interface{
						{
							Name:       "enp2t0",
							MacAddress: "98:af:65:a5:8d:03",
						},
					},
				},
			},
		},
	}
}

func unmarshalJSON(b []byte) []byte {
	output, _ := yaml.JSONToYAML(b)
	return output
}
