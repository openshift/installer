package manifests

import (
	"net"
	"strings"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/api/v1beta1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	agenttypes "github.com/openshift/installer/pkg/types/agent"
	"github.com/openshift/installer/pkg/types/baremetal"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"
)

var (
	testSSHKey = `|
	ssh-rsa AAAAB3NzaC1y1LJe3zew1ghc= root@localhost.localdomain`
	testSecret = `{"auths":{"cloud.openshift.com":{"auth":"b3BlUTA=","email":"test@redhat.com"}}}`

	rawNMStateConfig = `
    interfaces:
      - name: eth0
        type: ethernet
        state: up
        mac-address: 52:54:01:aa:aa:a1
        ipv4:
          enabled: true
          address:
            - ip: 192.168.122.21
              prefix-length: 24
          dhcp: false
    dns-resolver:
      config:
        server:
          - 192.168.122.1
    routes:
      config:
        - destination: 0.0.0.0/0
          next-hop-address: 192.168.122.1
          next-hop-interface: eth0
          table-id: 254`

	rawNMStateConfigNoIP = `
    interfaces:
      - name: eth0
        type: ethernet
        state: up
        mac-address: 52:54:01:aa:aa:a1`

	// config with route but no interface is invalid
	invalidRawNMStateConfig = `
    routes:
      config:
        - destination: 0.0.0.0/0
          next-hop-address: 192.168.122.1
          next-hop-interface: eth0
          table-id: 254`
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
				PullSecret: testSecret,
				SSHKey:     testSSHKey,
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
				PullSecret: testSecret,
				SSHKey:     testSSHKey,
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

func getValidOptionalInstallConfigDualStackDualVIPs() *agent.OptionalInstallConfig {
	installConfig := getValidOptionalInstallConfigDualStack()
	installConfig.Config.Platform.BareMetal.APIVIPs = append(installConfig.Config.Platform.BareMetal.APIVIPs, "2001:db8:1111:2222:ffff:ffff:ffff:cafe")
	installConfig.Config.Platform.BareMetal.IngressVIPs = append(installConfig.Config.Platform.BareMetal.IngressVIPs, "2001:db8:1111:2222:ffff:ffff:ffff:dead")
	return installConfig
}

// getProxyValidOptionalInstallConfig returns a valid optional install config for proxied installation
func getProxyValidOptionalInstallConfig() *agent.OptionalInstallConfig {
	validIC := getValidOptionalInstallConfig()
	validIC.InstallConfig.Config.Proxy = &types.Proxy{
		HTTPProxy:  "http://10.10.10.11:80",
		HTTPSProxy: "http://my-lab-proxy.org:443",
		NoProxy:    "internal.com",
	}
	return validIC
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
						Raw: unmarshalJSON([]byte(rawNMStateConfig)),
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
						Raw: unmarshalJSON([]byte(rawNMStateConfig)),
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
						Raw: unmarshalJSON([]byte(rawNMStateConfig)),
					},
				},
			},
		},
	}
}

func getValidAgentConfigWithAdditionalNTPSources() *agentconfig.AgentConfig {
	validAC := getValidAgentConfig()
	validAC.Config.AdditionalNTPSources = []string{
		"0.fedora.pool.ntp.org",
		"1.fedora.pool.ntp.org",
	}
	return validAC
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
						Raw: unmarshalJSON([]byte(rawNMStateConfigNoIP)),
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

func getGoodACI() *hiveext.AgentClusterInstall {
	goodACI := &hiveext.AgentClusterInstall{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getAgentClusterInstallName(getValidOptionalInstallConfig()),
			Namespace: getObjectMetaNamespace(getValidOptionalInstallConfig()),
		},
		Spec: hiveext.AgentClusterInstallSpec{
			ImageSetRef: &hivev1.ClusterImageSetReference{
				Name: getClusterImageSetReferenceName(),
			},
			ClusterDeploymentRef: corev1.LocalObjectReference{
				Name: getClusterDeploymentName(getValidOptionalInstallConfig()),
			},
			Networking: hiveext.Networking{
				MachineNetwork: []hiveext.MachineNetworkEntry{
					{
						CIDR: "10.10.11.0/24",
					},
				},
				ClusterNetwork: []hiveext.ClusterNetworkEntry{
					{
						CIDR:       "192.168.111.0/24",
						HostPrefix: 23,
					},
				},
				ServiceNetwork: []string{"172.30.0.0/16"},
				NetworkType:    "OVNKubernetes",
			},
			SSHPublicKey: strings.Trim(testSSHKey, "|\n\t"),
			ProvisionRequirements: hiveext.ProvisionRequirements{
				ControlPlaneAgents: 3,
				WorkerAgents:       5,
			},
			APIVIP:       "192.168.122.10",
			IngressVIP:   "192.168.122.11",
			PlatformType: hiveext.BareMetalPlatformType,
		},
	}
	return goodACI
}

func getInValidAgentConfig() *agentconfig.AgentConfig {
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
					Interfaces: []*v1beta1.Interface{
						{
							Name:       "enp2s0",
							MacAddress: "98:af:65:a5:8d:01",
						},
					},
					NetworkConfig: v1beta1.NetConfig{
						Raw: unmarshalJSON([]byte(invalidRawNMStateConfig)),
					},
				},
			},
		},
	}
}

func getGoodACIDualStack() *hiveext.AgentClusterInstall {
	goodACI := getGoodACI()
	goodACI.Spec.Networking.MachineNetwork = append(goodACI.Spec.Networking.MachineNetwork, hiveext.MachineNetworkEntry{
		CIDR: "2001:db8:5dd8:c956::/64",
	})
	goodACI.Spec.Networking.ClusterNetwork = append(goodACI.Spec.Networking.ClusterNetwork, hiveext.ClusterNetworkEntry{
		CIDR:       "2001:db8:1111:2222::/64",
		HostPrefix: 64,
	})
	goodACI.Spec.Networking.ServiceNetwork = append(goodACI.Spec.Networking.ServiceNetwork, "fd02::/112")

	return goodACI
}

func unmarshalJSON(b []byte) []byte {
	output, _ := yaml.JSONToYAML(b)
	return output
}
