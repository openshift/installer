package manifests

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestNMStateConfig_Generate(t *testing.T) {

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig []*aiv1beta1.NMStateConfig
	}{
		{
			name: "missing-config",
			dependencies: []asset.Asset{
				&agentconfig.AgentConfig{},
			},
			expectedError: "missing configuration or manifest file",
		},
		{
			name: "valid config",
			dependencies: []asset.Asset{
				getValidAgentConfig(),
			},
			expectedConfig: []*aiv1beta1.NMStateConfig{
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprint(getNMStateConfigName(getValidAgentConfig()), "-0"),
						Namespace: getNMStateConfigNamespace(getValidAgentConfig()),
						Labels:    getNMStateConfigLabelsFromAgentConfig(getValidAgentConfig()),
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
							Raw: unmarshalJSON([]byte("interfaces:")),
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprint(getNMStateConfigName(getValidAgentConfig()), "-1"),
						Namespace: getNMStateConfigNamespace(getValidAgentConfig()),
						Labels:    getNMStateConfigLabelsFromAgentConfig(getValidAgentConfig()),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "enp2t0",
								MacAddress: "98:af:65:a5:8d:02",
							},
						},
						NetConfig: aiv1beta1.NetConfig{
							Raw: unmarshalJSON([]byte("interfaces:")),
						},
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind:       "NMStateConfig",
						APIVersion: "agent-install.openshift.io/v1beta1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprint(getNMStateConfigName(getValidAgentConfig()), "-2"),
						Namespace: getNMStateConfigNamespace(getValidAgentConfig()),
						Labels:    getNMStateConfigLabelsFromAgentConfig(getValidAgentConfig()),
					},
					Spec: aiv1beta1.NMStateConfigSpec{
						Interfaces: []*aiv1beta1.Interface{
							{
								Name:       "enp2u0",
								MacAddress: "98:af:65:a5:8d:03",
							},
						},
						NetConfig: aiv1beta1.NetConfig{
							Raw: unmarshalJSON([]byte("interfaces:")),
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &NMStateConfig{}
			err := asset.Generate(parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-manifests/nmstateconfig.yaml", configFile.Filename)

				// Split up the file into multiple YAMLs if it contains NMStateConfig for more than one node
				var decoder nmStateConfigYamlDecoder
				yamlList, err := getMultipleYamls(configFile.Data, &decoder)

				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedConfig), len(yamlList))

				for i := range tc.expectedConfig {
					assert.Equal(t, tc.expectedConfig[i], yamlList[i])

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

		{
			name: "invalid-address-for-type",
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
        ipv6:
          enabled: true
          address:
            - ip: 192.168.122.21
              prefix-length: 24
  interfaces:
    - name: "eth0"
      macAddress: "52:54:01:aa:aa:a1"`,
			requiresNmstatectl: true,
			expectedError:      "staticNetwork configuration is not valid",
		},

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
			expectedError: "invalid NMStateConfig configuration: ObjectMeta.Labels: Required value: mynmstateconfig does not have any label set",
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
		// nmstate may not be installed yet in CI so skip this test if not
		if tc.requiresNmstatectl {
			_, execErr := exec.LookPath("nmstatectl")
			if execErr != nil {
				continue
			}
		}

		t.Run(tc.name, func(t *testing.T) {

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
