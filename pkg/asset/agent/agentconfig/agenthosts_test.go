package agentconfig

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	agentAsset "github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/agent"
	"github.com/openshift/installer/pkg/types/baremetal"
)

const (
	agentNetworkConfigOne = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.80
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:d2:b2:1a
  name: eth0
  state: up
  type: ethernet
`
	agentNetworkConfigTwo = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.81
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:d2:b2:1b
  name: eth0
  state: up
  type: ethernet
`
	installNetworkConfigOne = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.80
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:b0:bf:01
  name: eth0
  state: up
  type: ethernet
`
	installNetworkConfigTwo = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.81
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:b0:bf:02
  name: eth0
  state: up
  type: ethernet
`
	agentNetworkConfigEmbeddedRendezvousIPOne = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.1
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:d2:b2:1b
  name: eth0
  state: up
  type: ethernet
`
	agentNetworkConfigEmbeddedRendezvousIPTwo = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.2
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:d2:b2:1b
  name: eth0
  state: up
  type: ethernet
routes:
  config:
  - destination: 0.0.0.0/0
    next-hop-address: 192.168.111.126
    next-hop-interface: eth0
    table-id: 254
`
	agentNetworkConfigBond = `interfaces:
- name: eth0
  type: ethernet
  state: up
  mac-address: 28:d2:44:d2:b2:1a
- name: eth1
  type: ethernet
  state: up
  mac-address: 28:d2:44:d2:b2:1b
- name: bond0
  type: bond
  state: up
  link-aggregation:
    mode: active-backup
    port:
    - eth0
    - eth1
  ipv4:
    enabled: true
    dhcp: false
    address:
    - ip: 192.168.111.80
      prefix-length: 24
`
)

func TestAgentHosts_Generate(t *testing.T) {
	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *AgentHostsBuilder
	}{
		{
			name: "host-from-add-nodes-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.AddNodesConfig{
					Config: joiner.Config{
						Hosts: []agent.Host{
							{
								Hostname: "extra-worker-0",
								Role:     "worker",
								Interfaces: []*aiv1beta1.Interface{
									{
										Name:       "enp3s1",
										MacAddress: "28:d2:44:d2:b2:1a",
									},
								},
							},
						},
					},
				},
				getNoHostsInstallConfig(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: agentHosts().hosts(agentHost().name("extra-worker-0").role("worker").interfaces(iface("enp3s1", "28:d2:44:d2:b2:1a"))),
		},
		{
			name: "no-hosts",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getNoHostsInstallConfig(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: nil,
		},
		{
			name: "host-from-agent-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigSingleHost(),
			},
			expectedConfig: agentHosts().hosts(agentHost().name("test").role("master").interfaces(iface("enp3s1", "28:d2:44:d2:b2:1a")).deviceHint()),
		},
		{
			name: "host-from-install-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: agentHosts().hosts(agentHost().name("test").role("master").interfaces(iface("boot", "28:d2:44:b0:bf:01")).deviceHint()),
		},
		{
			name: "multi-host-from-agent-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigMultiHost("worker"),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:d2:b2:1a")).deviceHint().networkConfig(agentNetworkConfigOne),
				agentHost().name("test-2").role("worker").interfaces(iface("eth0", "28:d2:44:d2:b2:1b")).networkConfig(agentNetworkConfigTwo)),
		},
		{
			name: "multi-host-from-install-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigMultiHost(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:b0:bf:01")).deviceHint().networkConfig(installNetworkConfigOne),
				agentHost().name("test-2").role("worker").interfaces(iface("eth0", "28:d2:44:b0:bf:02")).networkConfig(installNetworkConfigTwo)),
		},
		{
			name: "unsupported-device-name-agent-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigUnsupportedDeviceName(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].rootDeviceHints.deviceName: Invalid value: \"/dev/disk/by-id/wwn-0x600508e000000000ce506dc50ab0ad05\": Device Name of root device hint must be path in /dev/ or /dev/disk/by-path/",
			expectedConfig: nil,
		},
		{
			name: "unsupported-wwn-extension-install-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigUnsupportedWWNExtension(),
				getNoHostsAgentConfig(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].rootDeviceHints.wwnWithExtension: Forbidden: WWN extensions are not supported in root device hints",
			expectedConfig: nil,
		},
		{
			name: "unsupported-www-vendor-extension-agent-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigUnsupportedWWNVendorExtension(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].rootDeviceHints.wwnVendorExtension: Forbidden: WWN vendor extensions are not supported in root device hints",
			expectedConfig: nil,
		},
		{
			name: "node-hostname-and-role-are-not-required",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigNoHostnameOrRole(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: agentHosts().hosts(agentHost().interfaces(iface("boot", "28:d2:44:b0:bf:01")).deviceHint()),
		},
		{
			name: "host-roles-have-incorrect-values",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigInvalidHostRole(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].role: Unsupported value: \"invalid-role\": supported values: \"master\", \"worker\", \"arbiter\"",
			expectedConfig: nil,
		},
		{
			name: "different-hosts-cannot-have-same-mac",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSameMac(),
				getNoHostsAgentConfig(),
			},
			expectedError:  "invalid Hosts configuration: hosts[1].interfaces[0].macAddress: Invalid value: \"28:d2:44:b0:bf:01\": duplicate MAC address found",
			expectedConfig: nil,
		},
		{
			name: "invalid-mac",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigInvalidMac(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].interfaces[0].macAddress: Invalid value: \"000000\": address 000000: invalid MAC address",
			expectedConfig: nil,
		},
		{
			name: "duplicate-mac-same-host-agent-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigInvalidInterfaces(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].interfaces[1].macAddress: Invalid value: \"28:d2:44:d2:b2:1a\": duplicate MAC address found",
			expectedConfig: nil,
		},
		{
			name: "duplicate-mac-same-host-install-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigInvalidInterfaces(),
				getNoHostsAgentConfig(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].interfaces[1].macAddress: Invalid value: \"28:d2:44:b0:bf:01\": duplicate MAC address found",
			expectedConfig: nil,
		},
		{
			name: "invalid-rendezvous-agent-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigInvalidRendezvousIP(),
			},
			expectedError:  "invalid Hosts configuration: hosts[1].role: Forbidden: Host test-2 has role 'worker' and has the rendezvousIP assigned to it. The rendezvousIP must be assigned to a control plane host.",
			expectedConfig: nil,
		},
		{
			name: "invalid-rendezvous-install-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigInvalidRendezvousIP(),
				getNoHostsAgentConfig(),
			},
			expectedError:  "invalid Hosts configuration: hosts[0].role: Forbidden: Host test has role 'worker' and has the rendezvousIP assigned to it. The rendezvousIP must be assigned to a control plane host.",
			expectedConfig: nil,
		},
		{
			name: "host-missing-interface-error",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigMissingInterfaces(),
			},
			expectedError:  "invalid Hosts configuration: [hosts[0].interfaces: Required value: at least one interface must be defined for each node, hosts[1].interfaces: Required value: at least one interface must be defined for each node, hosts[2].interfaces: Required value: at least one interface must be defined for each node]",
			expectedConfig: nil,
		},
		{
			name: "rendezvousip-in-worker-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getNoHostsInstallConfig(),
				getAgentConfigMultiHostEmbeddedRendezvousIP(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:d2:b2:1a")).deviceHint().networkConfig(agentNetworkConfigEmbeddedRendezvousIPOne),
				agentHost().name("test-2").role("worker").interfaces(iface("eth0", "28:d2:44:d2:b2:1b")).networkConfig(agentNetworkConfigEmbeddedRendezvousIPTwo)),
		},
		{
			name: "multi-host-from-agent-config-with-arbiter",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigMultiHost("arbiter"),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:d2:b2:1a")).deviceHint().networkConfig(agentNetworkConfigOne),
				agentHost().name("test-2").role("arbiter").interfaces(iface("eth0", "28:d2:44:d2:b2:1b")).networkConfig(agentNetworkConfigTwo)),
		},
		{
			name: "interface-name-mismatch-with-networkconfig-warns-only",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigMismatchedInterfaceName(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("enp3s0", "28:d2:44:d2:b2:1a")).deviceHint().networkConfig(agentNetworkConfigOne)),
		},
		{
			name: "interface-name-matches-networkconfig",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigMatchingInterfaceName(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:d2:b2:1a")).deviceHint().networkConfig(agentNetworkConfigOne)),
		},
		{
			name: "bond-networkconfig-with-matching-interfaces",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigBondMatching(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:d2:b2:1a"), iface("eth1", "28:d2:44:d2:b2:1b")).deviceHint().networkConfig(agentNetworkConfigBond)),
		},
		{
			name: "bond-networkconfig-with-mismatched-interfaces-warns-only",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigBondMismatched(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("nic1", "28:d2:44:d2:b2:1a"), iface("nic2", "28:d2:44:d2:b2:1b")).deviceHint().networkConfig(agentNetworkConfigBond)),
		},
		{
			name: "install-config-with-networkconfig-no-warning-inert-path",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithMismatchedNetworkConfig(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:b0:bf:01")).deviceHint().networkConfig(installNetworkConfigOne)),
		},
		{
			name: "add-nodes-interface-mismatch-warns-only",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.AddNodesConfig{
					Config: joiner.Config{
						Hosts: []agent.Host{
							{
								Hostname: "extra-worker-0",
								Role:     "worker",
								Interfaces: []*aiv1beta1.Interface{
									{
										Name:       "nic0",
										MacAddress: "28:d2:44:d2:b2:1a",
									},
								},
								NetworkConfig: aiv1beta1.NetConfig{
									Raw: []byte(agentNetworkConfigOne),
								},
							},
						},
					},
				},
				getNoHostsInstallConfig(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("extra-worker-0").role("worker").interfaces(iface("nic0", "28:d2:44:d2:b2:1a")).networkConfig(agentNetworkConfigOne)),
		},
		{
			name: "add-nodes-interface-matches-no-warning",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.AddNodesConfig{
					Config: joiner.Config{
						Hosts: []agent.Host{
							{
								Hostname: "extra-worker-0",
								Role:     "worker",
								Interfaces: []*aiv1beta1.Interface{
									{
										Name:       "eth0",
										MacAddress: "28:d2:44:d2:b2:1a",
									},
								},
								NetworkConfig: aiv1beta1.NetConfig{
									Raw: []byte(agentNetworkConfigOne),
								},
							},
						},
					},
				},
				getNoHostsInstallConfig(),
				getNoHostsAgentConfig(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("extra-worker-0").role("worker").interfaces(iface("eth0", "28:d2:44:d2:b2:1a")).networkConfig(agentNetworkConfigOne)),
		},
		{
			name: "agent-config-empty-interface-name-skipped",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigEmptyInterfaceName(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("", "28:d2:44:d2:b2:1a")).deviceHint().networkConfig(agentNetworkConfigOne)),
		},
		{
			name: "agent-config-malformed-networkconfig-no-panic",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigMalformedNetworkConfig(),
			},
			expectedConfig: agentHosts().hosts(
				agentHost().name("test").role("master").interfaces(iface("eth0", "28:d2:44:d2:b2:1a")).deviceHint().rawNetworkConfig("not: valid: yaml: [[")),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &AgentHosts{}
			err := asset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				if tc.expectedConfig != nil {
					assert.Equal(t, tc.expectedConfig.build().Hosts, asset.Hosts, "unexpected Config in AgentHosts")
				} else {
					assert.Nil(t, asset.Hosts)
				}
			}
		})
	}
}

func getNoHostsInstallConfig() *agentAsset.OptionalInstallConfig {
	_, newCidr, err1 := net.ParseCIDR("192.168.111.0/24")
	_, machineNetCidr, err2 := net.ParseCIDR("10.10.11.0/24")
	if err1 != nil || err2 != nil {
		return nil
	}

	return &agentAsset.OptionalInstallConfig{
		AssetBase: installconfig.AssetBase{
			Config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				BaseDomain: "test-domain",
				PullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}",
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64(3),
					Platform: types.MachinePoolPlatform{},
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

func getNoHostsAgentConfig() *AgentConfig {
	return &AgentConfig{
		Config: &agent.Config{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AgentConfig",
				APIVersion: "v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ocp-edge-cluster-0",
				Namespace: "cluster-0",
			},
			RendezvousIP: "192.168.111.80",
		},
	}
}

func getAgentConfigSingleHost() *AgentConfig {
	a := getNoHostsAgentConfig()
	a.Config.Hosts = []agent.Host{
		{
			Hostname: "test",
			Role:     "master",
			Interfaces: []*aiv1beta1.Interface{
				{
					Name:       "enp3s1",
					MacAddress: "28:d2:44:d2:b2:1a",
				},
			},
			RootDeviceHints: baremetal.RootDeviceHints{
				DeviceName: "/dev/sda",
			},
		},
	}
	return a
}

func getAgentConfigMultiHost(role string) *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces[0].Name = "eth0"
	a.Config.Hosts[0].NetworkConfig.Raw = []byte(agentNetworkConfigOne)
	host := agent.Host{
		Hostname: "test-2",
		Role:     role,
		Interfaces: []*aiv1beta1.Interface{
			{
				Name:       "eth0",
				MacAddress: "28:d2:44:d2:b2:1b",
			},
		},
		NetworkConfig: aiv1beta1.NetConfig{
			Raw: []byte(agentNetworkConfigTwo),
		},
	}
	a.Config.Hosts = append(a.Config.Hosts, host)
	return a
}

func getAgentConfigMultiHostEmbeddedRendezvousIP() *AgentConfig {
	a := getAgentConfigMultiHost("worker")
	a.Config.RendezvousIP = "192.168.111.1"
	a.Config.Hosts[0].NetworkConfig.Raw = []byte(agentNetworkConfigEmbeddedRendezvousIPOne)
	a.Config.Hosts[1].NetworkConfig.Raw = []byte(agentNetworkConfigEmbeddedRendezvousIPTwo)
	return a
}

func getAgentConfigUnsupportedDeviceName() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].RootDeviceHints = baremetal.RootDeviceHints{
		DeviceName: "/dev/disk/by-id/wwn-0x600508e000000000ce506dc50ab0ad05",
	}
	return a
}

func getAgentConfigUnsupportedWWNVendorExtension() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].RootDeviceHints = baremetal.RootDeviceHints{
		WWNVendorExtension: "wwn-with-vendor-extension-value",
	}
	return a
}

func getAgentConfigInvalidHostRole() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Role = "invalid-role"
	return a
}

func getAgentConfigInvalidMac() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces[0].MacAddress = "000000"
	return a
}

func getAgentConfigInvalidInterfaces() *AgentConfig {
	a := getNoHostsAgentConfig()
	a.Config.Hosts = []agent.Host{
		{
			Hostname: "test",
			Role:     "master",
			Interfaces: []*aiv1beta1.Interface{
				{
					Name:       "enp3s1",
					MacAddress: "28:d2:44:d2:b2:1a",
				},
				{
					Name:       "enp3s2",
					MacAddress: "28:d2:44:d2:b2:1a",
				},
			},
		},
	}
	return a
}

func getAgentConfigBondMatching() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces = []*aiv1beta1.Interface{
		{Name: "eth0", MacAddress: "28:d2:44:d2:b2:1a"},
		{Name: "eth1", MacAddress: "28:d2:44:d2:b2:1b"},
	}
	a.Config.Hosts[0].NetworkConfig.Raw = unmarshalJSON([]byte(agentNetworkConfigBond))
	return a
}

func getAgentConfigBondMismatched() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces = []*aiv1beta1.Interface{
		{Name: "nic1", MacAddress: "28:d2:44:d2:b2:1a"},
		{Name: "nic2", MacAddress: "28:d2:44:d2:b2:1b"},
	}
	a.Config.Hosts[0].NetworkConfig.Raw = unmarshalJSON([]byte(agentNetworkConfigBond))
	return a
}

func getAgentConfigMissingInterfaces() *AgentConfig {
	a := getNoHostsAgentConfig()
	a.Config.Hosts = []agent.Host{
		{
			Hostname: "control-0.example.org",
			Role:     "master",
		},
		{
			Hostname: "control-1.example.org",
			Role:     "master",
		},
		{
			Hostname: "control-2.example.org",
			Role:     "master",
		},
	}
	return a
}

func getAgentConfigMismatchedInterfaceName() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces[0].Name = "enp3s0"
	a.Config.Hosts[0].NetworkConfig.Raw = []byte(agentNetworkConfigOne)
	return a
}

func getAgentConfigMatchingInterfaceName() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces[0].Name = "eth0"
	a.Config.Hosts[0].NetworkConfig.Raw = []byte(agentNetworkConfigOne)
	return a
}

func getAgentConfigInvalidRendezvousIP() *AgentConfig {
	a := getAgentConfigMultiHost("worker")
	a.Config.RendezvousIP = "192.168.111.81"
	return a
}

func getInstallConfigSingleHost() *agentAsset.OptionalInstallConfig {
	a := getNoHostsInstallConfig()
	a.Config.Platform.BareMetal.Hosts = []*baremetal.Host{
		{
			Name:           "test",
			Role:           "master",
			BootMACAddress: "28:d2:44:b0:bf:01",
			RootDeviceHints: &baremetal.RootDeviceHints{
				DeviceName: "/dev/sda",
			},
		},
	}
	return a
}

func getInstallConfigMultiHost() *agentAsset.OptionalInstallConfig {
	a := getInstallConfigSingleHost()
	a.Config.Platform.BareMetal.Hosts[0].NetworkConfig = &apiextv1.JSON{
		Raw: []byte(installNetworkConfigOne),
	}
	host := &baremetal.Host{
		Name:           "test-2",
		Role:           "worker",
		BootMACAddress: "28:d2:44:b0:bf:02",
		NetworkConfig: &apiextv1.JSON{
			Raw: []byte(installNetworkConfigTwo),
		},
	}
	a.Config.Platform.BareMetal.Hosts = append(a.Config.Platform.BareMetal.Hosts, host)
	return a
}

func getInstallConfigSameMac() *agentAsset.OptionalInstallConfig {
	var networkConfigSameMac = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.81
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:b0:bf:01
  name: eth0
  state: up
  type: ethernet
`
	a := getInstallConfigMultiHost()
	a.Config.Platform.BareMetal.Hosts[1].BootMACAddress = "28:d2:44:b0:bf:01"
	a.Config.Platform.BareMetal.Hosts[1].NetworkConfig = &apiextv1.JSON{
		Raw: []byte(networkConfigSameMac),
	}
	return a
}

func getInstallConfigInvalidInterfaces() *agentAsset.OptionalInstallConfig {
	var networkConfigSameMacSameHost = `interfaces:
- ipv4:
    address:
    - ip: 192.168.111.80
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:b0:bf:01
  name: eth0
  state: up
  type: ethernet
- ipv4:
    address:
    - ip: 192.168.111.81
      prefix-length: 24
    dhcp: false
    enabled: true
  mac-address: 28:d2:44:b0:bf:01
  name: eth0
  state: up
  type: ethernet
`
	a := getInstallConfigSingleHost()
	a.Config.Platform.BareMetal.Hosts[0].NetworkConfig = &apiextv1.JSON{
		Raw: []byte(networkConfigSameMacSameHost),
	}
	return a
}

func getInstallConfigUnsupportedWWNExtension() *agentAsset.OptionalInstallConfig {
	a := getInstallConfigSingleHost()
	a.Config.Platform.BareMetal.Hosts[0].RootDeviceHints = &baremetal.RootDeviceHints{
		WWNWithExtension: "wwn-with-extension-value",
	}
	return a
}

func getInstallConfigNoHostnameOrRole() *agentAsset.OptionalInstallConfig {
	a := getInstallConfigSingleHost()
	a.Config.Platform.BareMetal.Hosts[0].Name = ""
	a.Config.Platform.BareMetal.Hosts[0].Role = ""
	return a
}

func getInstallConfigInvalidRendezvousIP() *agentAsset.OptionalInstallConfig {
	a := getInstallConfigSingleHost()
	a.Config.Platform.BareMetal.Hosts[0].NetworkConfig = &apiextv1.JSON{
		Raw: []byte(installNetworkConfigOne),
	}
	a.Config.Platform.BareMetal.Hosts[0].Role = "worker"
	return a
}

// AgentConfigBuilder it's a builder class to make it easier creating agent.Config instance
// used in the test cases.
type AgentHostsBuilder struct {
	AgentHosts
}

func agentHosts() *AgentHostsBuilder {
	return &AgentHostsBuilder{
		AgentHosts: AgentHosts{
			Hosts: []agent.Host{{}},
		},
	}
}

func (ahb *AgentHostsBuilder) build() *AgentHosts {
	return &ahb.AgentHosts
}

func (ahb *AgentHostsBuilder) hosts(builders ...*HostBuilder) *AgentHostsBuilder {
	hosts := []agent.Host{}
	for _, b := range builders {
		hosts = append(hosts, *b.build())
	}
	ahb.Hosts = hosts

	return ahb
}

// HostBuilder it's a builder class to make it easier creating agent.Host instances
// used in the test cases, as part of the AgentHosts type.
type HostBuilder struct {
	agent.Host
}

func agentHost() *HostBuilder {
	return &HostBuilder{}
}

func (hb *HostBuilder) build() *agent.Host {
	return &hb.Host
}

func (hb *HostBuilder) name(name string) *HostBuilder {
	hb.Host.Hostname = name
	return hb
}

func (hb *HostBuilder) role(role string) *HostBuilder {
	hb.Host.Role = role
	return hb
}

func (hb *HostBuilder) interfaces(builders ...*InterfacetBuilder) *HostBuilder {
	ifaces := []*aiv1beta1.Interface{}
	for _, b := range builders {
		ifaces = append(ifaces, b.build())
	}
	hb.Host.Interfaces = ifaces
	return hb
}

func (hb *HostBuilder) networkConfig(raw string) *HostBuilder {
	hb.Host.NetworkConfig = aiv1beta1.NetConfig{
		Raw: unmarshalJSON([]byte(raw)),
	}
	return hb
}

func (hb *HostBuilder) rawNetworkConfig(raw string) *HostBuilder {
	hb.Host.NetworkConfig = aiv1beta1.NetConfig{
		Raw: []byte(raw),
	}
	return hb
}

func (hb *HostBuilder) deviceHint() *HostBuilder {
	hb.Host.RootDeviceHints = baremetal.RootDeviceHints{
		DeviceName: "/dev/sda",
	}
	return hb
}

// InterfacetBuilder it's a builder class to make it easier creating aiv1beta1.Interface instances
// used in the test cases, as part of the agent.Config type.
type InterfacetBuilder struct {
	aiv1beta1.Interface
}

func iface(name string, mac string) *InterfacetBuilder {
	return &InterfacetBuilder{
		Interface: aiv1beta1.Interface{
			Name:       name,
			MacAddress: mac,
		},
	}
}

func (ib *InterfacetBuilder) build() *aiv1beta1.Interface {
	return &ib.Interface
}

func TestAgentHosts_FencingCredentialsByHost(t *testing.T) {
	cases := []struct {
		name                  string
		dependencies          []asset.Asset
		expectedError         string
		expectedFencingByHost []FencingCredentialHost
	}{
		{
			name: "mac-keyed credentials populate FencingCredentialsByHost",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithFencing([]*types.Credential{
					{MACAddress: "28:d2:44:d2:b2:1a", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
				}),
				getAgentConfigMultiHost("worker"),
			},
			expectedFencingByHost: []FencingCredentialHost{
				{
					DirName: "test",
					Credentials: []*types.Credential{
						{MACAddress: "28:d2:44:d2:b2:1a", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
					},
				},
			},
		},
		{
			name: "hostname-keyed credentials are excluded from FencingCredentialsByHost",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithFencing([]*types.Credential{
					{HostName: "master-0", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
				}),
				getAgentConfigMultiHost("worker"),
			},
			expectedFencingByHost: nil,
		},
		{
			name: "mixed credentials only include mac-keyed in FencingCredentialsByHost",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithFencing([]*types.Credential{
					{HostName: "master-0", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
					{MACAddress: "28:d2:44:d2:b2:1b", Username: "admin2", Password: "pass2", Address: "redfish+https://10.0.0.2/redfish/v1/Systems/1"},
				}),
				getAgentConfigMultiHost("worker"),
			},
			expectedFencingByHost: []FencingCredentialHost{
				{
					DirName: "test-2",
					Credentials: []*types.Credential{
						{MACAddress: "28:d2:44:d2:b2:1b", Username: "admin2", Password: "pass2", Address: "redfish+https://10.0.0.2/redfish/v1/Systems/1"},
					},
				},
			},
		},
		{
			name: "credential with both hostname and mac is treated as hostname-keyed",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithFencing([]*types.Credential{
					{HostName: "master-0", MACAddress: "28:d2:44:d2:b2:1a", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
				}),
				getAgentConfigMultiHost("worker"),
			},
			expectedFencingByHost: nil,
		},
		{
			name: "mac with no matching host returns error",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithFencing([]*types.Credential{
					{MACAddress: "FF:FF:FF:FF:FF:FF", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
				}),
				getAgentConfigMultiHost("worker"),
			},
			expectedError: "does not match any configured host interface",
		},
		{
			name: "host without hostname uses host-i directory",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithFencing([]*types.Credential{
					{MACAddress: "28:d2:44:d2:b2:1a", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
				}),
				getAgentConfigNoHostname(),
			},
			expectedFencingByHost: []FencingCredentialHost{
				{
					DirName: "host-0",
					Credentials: []*types.Credential{
						{MACAddress: "28:d2:44:d2:b2:1a", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
					},
				},
			},
		},
		{
			name: "case-insensitive mac matching",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigWithFencing([]*types.Credential{
					{MACAddress: "28:D2:44:D2:B2:1A", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
				}),
				getAgentConfigMultiHost("worker"),
			},
			expectedFencingByHost: []FencingCredentialHost{
				{
					DirName: "test",
					Credentials: []*types.Credential{
						{MACAddress: "28:D2:44:D2:B2:1A", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
					},
				},
			},
		},
		{
			name: "no fencing config produces empty FencingCredentialsByHost",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.AddNodesConfig{},
				getInstallConfigSingleHost(),
				getAgentConfigSingleHost(),
			},
			expectedFencingByHost: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			ah := &AgentHosts{}
			err := ah.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedFencingByHost, ah.FencingCredentialsByHost)
		})
	}
}

func TestHostConfigFiles_FencingCredentials(t *testing.T) {
	ah := &AgentHosts{
		Hosts: []agent.Host{
			{
				Hostname: "master-0",
				Interfaces: []*aiv1beta1.Interface{
					{MacAddress: "aa:bb:cc:dd:ee:01"},
				},
			},
		},
		FencingCredentialsByHost: []FencingCredentialHost{
			{
				DirName: "master-0",
				Credentials: []*types.Credential{
					{MACAddress: "AA:BB:CC:DD:EE:01", Username: "admin", Password: "pass", Address: "redfish+https://10.0.0.1/redfish/v1/Systems/1"},
				},
			},
		},
	}

	files, err := ah.HostConfigFiles()
	assert.NoError(t, err)

	data, ok := files["master-0/fencing-credentials.yaml"]
	assert.True(t, ok, "expected fencing-credentials.yaml in master-0 directory")
	assert.Contains(t, string(data), "credentials:")
	assert.Contains(t, string(data), "username: admin")
}

func getInstallConfigWithFencing(credentials []*types.Credential) *agentAsset.OptionalInstallConfig {
	ic := getNoHostsInstallConfig()
	ic.Config.ControlPlane.Fencing = &types.Fencing{
		Credentials: credentials,
	}
	return ic
}

func getAgentConfigNoHostname() *AgentConfig {
	a := getNoHostsAgentConfig()
	a.Config.Hosts = []agent.Host{
		{
			Role: "master",
			Interfaces: []*aiv1beta1.Interface{
				{
					Name:       "eth0",
					MacAddress: "28:d2:44:d2:b2:1a",
				},
			},
		},
	}
	return a
}

func getInstallConfigWithMismatchedNetworkConfig() *agentAsset.OptionalInstallConfig {
	a := getInstallConfigSingleHost()
	a.Config.Platform.BareMetal.Hosts[0].NetworkConfig = &apiextv1.JSON{
		Raw: []byte(installNetworkConfigOne),
	}
	return a
}

func getAgentConfigEmptyInterfaceName() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces[0].Name = ""
	a.Config.Hosts[0].NetworkConfig.Raw = []byte(agentNetworkConfigOne)
	return a
}

func getAgentConfigMalformedNetworkConfig() *AgentConfig {
	a := getAgentConfigSingleHost()
	a.Config.Hosts[0].Interfaces[0].Name = "eth0"
	a.Config.Hosts[0].NetworkConfig.Raw = []byte("not: valid: yaml: [[")
	return a
}
