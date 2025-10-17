package manifests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	agentconfig "github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types/agent"
)

func TestNMStateConfig_Generate(t *testing.T) {
	cases := []struct {
		name               string
		dependencies       []asset.Asset
		requiresNmstatectl bool
		expectedConfig     []*aiv1beta1.NMStateConfig
		expectedError      string
	}{
		{
			name: "add-nodes workflow",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{},
				getAgentHostsNoHosts(),
				&agentconfig.OptionalInstallConfig{},
			},
			requiresNmstatectl: false,
			expectedConfig:     nil,
			expectedError:      "",
		},
		{
			name: "add-nodes workflow - agentHosts with some hosts without networkconfig",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{
					Namespace:   "cluster0",
					ClusterName: "ostest",
					Nodes:       &v1.NodeList{},
				},
				getAgentHostsWithSomeHostsWithoutNetworkConfig(),
				&agentconfig.OptionalInstallConfig{},
			},
			requiresNmstatectl: true,
			expectedConfig: []*aiv1beta1.NMStateConfig{
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ostest-0",
						Namespace: "cluster0",
						Labels:    getNMStateConfigLabels("ostest"),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "enp2t0",
								MacAddress: "98:af:65:a5:8d:02",
							},
						},
						NetConfig: aiv1beta1.NetConfig{
							Raw: unmarshalJSON([]byte(rawNMStateConfigNoIP)),
						},
					},
				},
			},
			expectedError: "",
		},
		{
			name: "add-nodes workflow - invalid ip",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{
					Namespace:   "cluster0",
					ClusterName: "ostest",
					Nodes: &v1.NodeList{
						Items: []v1.Node{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "master-0",
								},
								Status: v1.NodeStatus{
									Addresses: []v1.NodeAddress{
										{
											Address: "192.168.122.21", // configured by getValidAgentHostsConfig()
										},
									},
								},
							},
						},
					},
				},
				getValidAgentHostsConfig(),
				&agentconfig.OptionalInstallConfig{},
			},
			requiresNmstatectl: false,
			expectedError:      "address conflict found. The configured address 192.168.122.21 is already used by the cluster node master-0",
		},
		{
			name: "add-nodes workflow - invalid hostname",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{
					Namespace:   "cluster0",
					ClusterName: "ostest",
					Nodes: &v1.NodeList{
						Items: []v1.Node{
							{
								ObjectMeta: metav1.ObjectMeta{
									Name: "control-0.example.org",
								},
								Status: v1.NodeStatus{
									Addresses: []v1.NodeAddress{
										{
											Address: "control-0.example.org", // configured by getValidAgentHostsConfig()
										},
									},
								},
							},
						},
					},
				},
				getValidAgentHostsConfig(),
				&agentconfig.OptionalInstallConfig{},
			},
			requiresNmstatectl: false,
			expectedError:      "hostname conflict found. The configured hostname control-0.example.org is already used in the cluster",
		},
		{
			name: "agentHosts does not contain networkConfig",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				getAgentHostsNoHosts(),
				getValidOptionalInstallConfig(),
			},
			requiresNmstatectl: false,
			expectedConfig:     nil,
			expectedError:      "",
		},
		{
			name: "agentHosts with some hosts without networkconfig",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				getAgentHostsWithSomeHostsWithoutNetworkConfig(),
				getValidOptionalInstallConfig(),
			},
			requiresNmstatectl: true,
			expectedConfig: []*aiv1beta1.NMStateConfig{
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprint(getValidOptionalInstallConfig().ClusterName(), "-0"),
						Namespace: getValidOptionalInstallConfig().ClusterNamespace(),
						Labels:    getNMStateConfigLabels(getValidOptionalInstallConfig().ClusterName()),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "enp2t0",
								MacAddress: "98:af:65:a5:8d:02",
							},
						},
						NetConfig: aiv1beta1.NetConfig{
							Raw: unmarshalJSON([]byte(rawNMStateConfigNoIP)),
						},
					},
				},
			},
			expectedError: "",
		},
		{
			name: "valid config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				getValidAgentHostsConfig(),
				getValidOptionalInstallConfig(),
			},
			requiresNmstatectl: true,
			expectedConfig: []*aiv1beta1.NMStateConfig{
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprint(getValidOptionalInstallConfig().ClusterName(), "-0"),
						Namespace: getValidOptionalInstallConfig().ClusterNamespace(),
						Labels:    getNMStateConfigLabels(getValidOptionalInstallConfig().ClusterName()),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "enp2s0",
								MacAddress: "98:af:65:a5:8d:01",
							},
							{
								Name:       "enp3s1",
								MacAddress: "28:d2:44:d2:b2:1a",
							},
						},
						NetConfig: aiv1beta1.NetConfig{
							Raw: unmarshalJSON([]byte(rawNMStateConfig)),
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprint(getValidOptionalInstallConfig().ClusterName(), "-1"),
						Namespace: getValidOptionalInstallConfig().ClusterNamespace(),
						Labels:    getNMStateConfigLabels(getValidOptionalInstallConfig().ClusterName()),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "enp2t0",
								MacAddress: "98:af:65:a5:8d:02",
							},
						},
						NetConfig: aiv1beta1.NetConfig{
							Raw: unmarshalJSON([]byte(rawNMStateConfig)),
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprint(getValidOptionalInstallConfig().ClusterName(), "-2"),
						Namespace: getValidOptionalInstallConfig().ClusterNamespace(),
						Labels:    getNMStateConfigLabels(getValidOptionalInstallConfig().ClusterName()),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "enp2u0",
								MacAddress: "98:af:65:a5:8d:03",
							},
						},
						NetConfig: aiv1beta1.NetConfig{
							Raw: unmarshalJSON([]byte(rawNMStateConfig)),
						},
					},
				},
			},
			expectedError: "",
		},
		{
			name: "invalid networkConfig",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				getInValidAgentHostsConfig(),
				getValidOptionalInstallConfig(),
			},
			requiresNmstatectl: true,
			expectedConfig:     nil,
			expectedError:      "failed to validate network yaml",
		},
		{
			name: "invalid networkConfig, no interfaces in interface table",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				getAgentHostsConfigNoInterfaces(),
				getValidOptionalInstallConfig(),
			},
			requiresNmstatectl: true,
			expectedConfig:     nil,
			expectedError:      "at least one interface for host 0 must be provided",
		},
		{
			name: "invalid networkConfig, invalid mac in interface table",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				getAgentHostsConfigInvalidMac(),
				getValidOptionalInstallConfig(),
			},
			requiresNmstatectl: true,
			expectedConfig:     nil,
			expectedError:      "MAC address 98-af-65-a5-8d-02 for host 0 is incorrectly formatted, use XX:XX:XX:XX:XX:XX format",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &NMStateConfig{}
			err := asset.Generate(context.Background(), parents)

			// Check if the test failed because nmstatectl is not available in CI
			if tc.requiresNmstatectl {
				_, execErr := exec.LookPath("nmstatectl")
				if execErr != nil {
					assert.ErrorContains(t, err, "executable file not found")
					t.Skip("No nmstatectl binary available")
				}
			}

			switch {
			case tc.expectedError != "":
				assert.ErrorContains(t, err, tc.expectedError)
			case len(tc.expectedConfig) == 0:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-manifests/nmstateconfig.yaml", configFile.Filename)

				// Split up the file into multiple YAMLs if it contains NMStateConfig for more than one node
				yamlList, err := GetMultipleYamls[aiv1beta1.NMStateConfig](configFile.Data)

				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedConfig), len(yamlList))

				for i := range tc.expectedConfig {
					assert.Equal(t, *tc.expectedConfig[i], yamlList[i])
				}
				assert.Equal(t, len(tc.expectedConfig), len(asset.StaticNetworkConfig))
			}
		})
	}
}

func TestNMStateConfig_LoadedFromDisk(t *testing.T) {
	cases := []struct {
		name               string
		data               string
		fetchError         error
		expectedFound      bool
		expectedError      string
		requiresNmstatectl bool
		expectedConfig     []*models.HostStaticNetworkConfig
	}{
		{
			name: "valid-config-file",
			data: `
metadata:
  name: mynmstateconfig
  namespace: spoke-cluster
  labels:
    cluster0-nmstate-label-name: cluster0-nmstate-label-value
spec:
  config:
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
          table-id: 254
  interfaces:
    - name: "eth0"
      macAddress: "52:54:01:aa:aa:a1"
    - name: "eth1"
      macAddress: "52:54:01:bb:bb:b1"`,
			requiresNmstatectl: true,
			expectedFound:      true,
			expectedConfig: []*models.HostStaticNetworkConfig{
				{
					MacInterfaceMap: models.MacInterfaceMap{
						{LogicalNicName: "eth0", MacAddress: "52:54:01:aa:aa:a1"},
						{LogicalNicName: "eth1", MacAddress: "52:54:01:bb:bb:b1"},
					},
					NetworkYaml: "dns-resolver:\n  config:\n    server:\n    - 192.168.122.1\ninterfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.21\n      prefix-length: 24\n    dhcp: false\n    enabled: true\n  mac-address: 52:54:01:aa:aa:a1\n  name: eth0\n  state: up\n  type: ethernet\nroutes:\n  config:\n  - destination: 0.0.0.0/0\n    next-hop-address: 192.168.122.1\n    next-hop-interface: eth0\n    table-id: 254\n",
				},
			},
		},

		{
			name: "valid-config-multiple-yamls",
			data: `
metadata:
  name: mynmstateconfig
  namespace: spoke-cluster
  labels:
    cluster0-nmstate-label-name: cluster0-nmstate-label-value
spec:
  config:
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
  interfaces:
    - name: "eth0"
      macAddress: "52:54:01:aa:aa:a1"
---
metadata:
  name: mynmstateconfig-2
  namespace: spoke-cluster
  labels:
    cluster0-nmstate-label-name: cluster0-nmstate-label-value
spec:
  config:
    interfaces:
      - name: eth0
        type: ethernet
        state: up
        mac-address: 52:54:01:cc:cc:c1
        ipv4:
          enabled: true
          address:
            - ip: 192.168.122.22
              prefix-length: 24
  interfaces:
    - name: "eth0"
      macAddress: "52:54:01:cc:cc:c1"`,
			requiresNmstatectl: true,
			expectedFound:      true,
			expectedConfig: []*models.HostStaticNetworkConfig{
				{
					MacInterfaceMap: models.MacInterfaceMap{
						{LogicalNicName: "eth0", MacAddress: "52:54:01:aa:aa:a1"},
					},
					NetworkYaml: "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.21\n      prefix-length: 24\n    enabled: true\n  mac-address: 52:54:01:aa:aa:a1\n  name: eth0\n  state: up\n  type: ethernet\n",
				},
				{
					MacInterfaceMap: models.MacInterfaceMap{
						{LogicalNicName: "eth0", MacAddress: "52:54:01:cc:cc:c1"},
					},
					NetworkYaml: "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.22\n      prefix-length: 24\n    enabled: true\n  mac-address: 52:54:01:cc:cc:c1\n  name: eth0\n  state: up\n  type: ethernet\n",
				},
			},
		},

		{
			name: "invalid-interfaces",
			data: `
metadata:
  name: mynmstateconfig
  namespace: spoke-cluster
  labels:
    cluster0-nmstate-label-name: cluster0-nmstate-label-value
spec:
  interfaces:
    - name: "eth0"
      macAddress: "52:54:01:aa:aa:a1"
    - name: "eth0"
      macAddress: "52:54:01:bb:bb:b1"`,
			requiresNmstatectl: true,
			expectedError:      "staticNetwork configuration is not valid",
		},

		// This test case currently does not work for libnmstate 2.2.9,
		// due a regression that will be fixed in https://github.com/nmstate/nmstate/issues/2311
		// 		{
		// 			name: "invalid-address-for-type",
		// 			data: `
		// metadata:
		//   name: mynmstateconfig
		//   namespace: spoke-cluster
		//   labels:
		//     cluster0-nmstate-label-name: cluster0-nmstate-label-value
		// spec:
		//   config:
		//     interfaces:
		//       - name: eth0
		//         type: ethernet
		//         state: up
		//         mac-address: 52:54:01:aa:aa:a1
		//         ipv6:
		//           enabled: true
		//           address:
		//             - ip: 192.168.122.21
		//               prefix-length: 24
		//   interfaces:
		//     - name: "eth0"
		//       macAddress: "52:54:01:aa:aa:a1"`,
		// 			requiresNmstatectl: true,
		// 			expectedError:      "staticNetwork configuration is not valid",
		// 		},

		{
			name: "missing-label",
			data: `
metadata:
  name: mynmstateconfig
  namespace: spoke-cluster
spec:
  config:
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
  interfaces:
    - name: "eth0"
      macAddress: "52:54:01:aa:aa:a1"`,
			requiresNmstatectl: true,
			expectedError:      "invalid NMStateConfig configuration: labels: Required value: mynmstateconfig does not have any label set",
		},

		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: "could not decode YAML for cluster-manifests/nmstateconfig.yaml: Error reading multiple YAMLs: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1beta1.NMStateConfig",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load file cluster-manifests/nmstateconfig.yaml: fetch failed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// nmstate may not be installed yet in CI so skip this test if not
			if tc.requiresNmstatectl {
				_, execErr := exec.LookPath("nmstatectl")
				if execErr != nil {
					t.Skip("No nmstatectl binary available")
				}
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(nmStateConfigFilename).
				Return(
					&asset.File{
						Filename: nmStateConfigFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &NMStateConfig{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.StaticNetworkConfig, "unexpected Config in NMStateConfig")
				assert.Equal(t, len(tc.expectedConfig), len(asset.Config))
				for i := 0; i < len(tc.expectedConfig); i++ {

					staticNetworkConfig := asset.StaticNetworkConfig[i]
					nmStateConfig := asset.Config[i]

					for n := 0; n < len(staticNetworkConfig.MacInterfaceMap); n++ {
						macInterfaceMap := staticNetworkConfig.MacInterfaceMap[n]
						iface := nmStateConfig.Spec.Interfaces[n]

						assert.Equal(t, macInterfaceMap.LogicalNicName, iface.Name)
						assert.Equal(t, macInterfaceMap.MacAddress, iface.MacAddress)
					}
					assert.YAMLEq(t, staticNetworkConfig.NetworkYaml, string(nmStateConfig.Spec.NetConfig.Raw))
				}

			}
		})
	}
}

func TestGetNodeZeroIP(t *testing.T) {
	cases := []struct {
		name          string
		expectedIP    string
		expectedError string
		configs       []string
		hosts         []agent.Host
	}{
		{
			name:          "no interfaces",
			expectedError: "no interface IPs set",
		},
		{
			name:       "first interface",
			expectedIP: "192.168.122.21",
			configs: []string{
				`
interfaces:
  - name: eth0
    type: ethernet
    ipv4:
      address:
        - ip: 192.168.122.21
  - name: eth1
    type: ethernet
    ipv4:
      address:
        - ip: 192.168.122.22
`,
			},
		},
		{
			name:       "second interface",
			expectedIP: "192.168.122.22",
			configs: []string{
				`
interfaces:
  - name: eth0
    type: ethernet
  - name: eth1
    type: ethernet
    ipv4:
      address:
        - ip: 192.168.122.22
`,
			},
		},
		{
			name:       "second host",
			expectedIP: "192.168.122.22",
			configs: []string{
				`
interfaces:
  - name: eth0
    type: ethernet
  - name: eth1
    type: ethernet
`,
				`
interfaces:
  - name: eth0
    type: ethernet
  - name: eth1
    type: ethernet
    ipv4:
      address:
        - ip: 192.168.122.22
`,
			},
		},
		{
			name:       "ipv4 first",
			expectedIP: "192.168.122.22",
			configs: []string{
				`
interfaces:
  - name: eth0
    type: ethernet
    ipv6:
      address:
        - ip: "2001:0db8::0001"
    ipv4:
      address:
        - ip: 192.168.122.22
`,
			},
		},
		{
			name:       "ipv6 host first",
			expectedIP: "2001:0db8::0001",
			configs: []string{
				`
interfaces:
  - name: eth0
    type: ethernet
    ipv6:
      address:
        - ip: "2001:0db8::0001"
`,
				`
interfaces:
  - name: eth0
    type: ethernet
    ipv4:
      address:
        - ip: 192.168.122.31
`,
			},
		},
		{
			name:       "ipv6 first",
			expectedIP: "2001:0db8::0001",
			configs: []string{
				`
interfaces:
  - name: eth0
    type: ethernet
    ipv6:
      address:
        - ip: "2001:0db8::0001"
  - name: eth1
    type: ethernet
    ipv4:
      address:
        - ip: 192.168.122.22
`,
			},
		},
		{
			name:       "ipv6",
			expectedIP: "2001:0db8::0001",
			configs: []string{
				`
interfaces:
  - name: eth0
    type: ethernet
    ipv6:
      address:
        - ip: "2001:0db8::0001"
`,
			},
		},
		{
			name:       "skip workers/nodes without role",
			expectedIP: "192.168.122.22",
			hosts: []agent.Host{
				{
					Role: "worker",
					NetworkConfig: aiv1beta1.NetConfig{Raw: []byte(`
interfaces:
- name: eth0
  type: ethernet
  ipv4:
    address:
      - ip: 192.168.122.31`)},
				},
				{
					Role: "",
					NetworkConfig: aiv1beta1.NetConfig{Raw: []byte(`
interfaces:
- name: eth0
  type: ethernet
  ipv4:
    address:
      - ip: 192.168.122.32`)},
				},
				{
					Role: "master",
					NetworkConfig: aiv1beta1.NetConfig{Raw: []byte(`
interfaces:
- name: eth0
  type: ethernet
  ipv4:
    address:
      - ip: 192.168.122.22`)},
				},
			},
		},
		{
			name:          "fail if only workers",
			expectedError: "invalid NMState configurations provided, no interface IPs set",
			hosts: []agent.Host{
				{
					Role: "worker",
					NetworkConfig: aiv1beta1.NetConfig{Raw: []byte(`
interfaces:
- name: eth0
  type: ethernet
  ipv4:
    address:
      - ip: 192.168.122.31`)},
				},
			},
		},
		{
			name:          "fail if only master without static configuration",
			expectedError: "invalid NMState configurations provided, no interface IPs set",
			hosts: []agent.Host{
				{
					Role: "master",
				},
			},
		},
		{
			name:       "fallback on configs if missing host definition",
			expectedIP: "192.168.122.22",
			hosts: []agent.Host{
				{
					Role: "master",
				},
			},
			configs: []string{`
interfaces:
  - name: eth0
    type: ethernet
    ipv4:
      address:
        - ip: 192.168.122.22`,
			},
		},
		{
			name:       "implicit masters",
			expectedIP: "192.168.122.32",
			hosts: []agent.Host{
				{
					Role: "worker",
					NetworkConfig: aiv1beta1.NetConfig{Raw: []byte(`
interfaces:
- name: eth0
  type: ethernet
  ipv4:
    address:
      - ip: 192.168.122.31`)},
				},
				{
					Role: "",
					NetworkConfig: aiv1beta1.NetConfig{Raw: []byte(`
interfaces:
- name: eth0
  type: ethernet
  ipv4:
    address:
      - ip: 192.168.122.32`)},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var configs []*aiv1beta1.NMStateConfig
			for _, hostRaw := range tc.configs {
				configs = append(configs, &aiv1beta1.NMStateConfig{
					Spec: aiv1beta1.NMStateConfigSpec{
						NetConfig: aiv1beta1.NetConfig{
							Raw: aiv1beta1.RawNetConfig(hostRaw),
						},
					},
				})
			}

			ip, err := GetNodeZeroIP(tc.hosts, configs)
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedIP, ip)
			} else {
				assert.ErrorContains(t, err, tc.expectedError)
			}
		})
	}
}
