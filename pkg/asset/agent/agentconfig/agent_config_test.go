package agentconfig

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types/agent"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAgentConfig_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name       string
		data       string
		fetchError error

		expectedError  string
		expectedFound  bool
		expectedConfig *AgentConfigBuilder
	}{
		{
			name: "valid-config-single-node",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - hostname: control-0.example.org
    role: master
    rootDeviceHints:
      deviceName: "/dev/sda"
      hctl: "hctl-value"
      model: "model-value"
      vendor: "vendor-value"
      serialNumber: "serial-number-value"
      minSizeGigabytes: 20
      wwn: "wwn-value"
      rotational: false
    interfaces:
      - name: enp2s0
        macAddress: 98:af:65:a5:8d:01
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a
    networkConfig:
      interfaces:`,

			expectedFound:  true,
			expectedConfig: agentConfig().hosts(defaultAgentHost("control-0.example.org")),
		},
		{
			name: "valid-config-multiple-nodes",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - hostname: control-0.example.org
    role: master
    rootDeviceHints:
        deviceName: "/dev/sda"
        hctl: "hctl-value"
        model: "model-value"
        vendor: "vendor-value"
        serialNumber: "serial-number-value"
        minSizeGigabytes: 20
        wwn: "wwn-value"
        rotational: false
    interfaces:
      - name: enp2s0
        macAddress: 98:af:65:a5:8d:01
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a
    networkConfig:
      interfaces:
  - hostname: control-1.example.org
    role: master
    interfaces:
      - name: enp2s0
        macAddress: 98:af:65:a5:8d:02
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1b`,

			expectedFound: true,
			expectedConfig: agentConfig().
				hosts(
					defaultAgentHost("control-0.example.org"),
					agentHost().
						name("control-1.example.org").
						role("master").
						interfaces(
							iface("enp2s0", "98:af:65:a5:8d:02"),
							iface("enp3s1", "28:d2:44:d2:b2:1b"),
						),
				),
		},
		{
			name: "not-yaml",
			data: `This is not a yaml file`,

			expectedFound: false,
			expectedError: "failed to unmarshal agent-config.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type agent.Config",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},

			expectedFound: false,
		},
		{
			name:       "error-fetching-file",
			fetchError: errors.New("fetch failed"),

			expectedFound: false,
			expectedError: "failed to load agent-config.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-wrong
wrongField: wrongValue`,

			expectedFound: false,
			expectedError: "failed to unmarshal agent-config.yaml: error unmarshaling JSON: while decoding JSON: json: unknown field \"wrongField\"",
		},
		{
			name: "interface-missing-mac-address-error",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - hostname: control-0.example.org
    interfaces:
      - name: enp2s0
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: Hosts[0].Interfaces[0].macAddress: Required value: each interface must have a MAC address defined",
		},
		{
			name: "unsupported wwn extension root device hint",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - hostname: control-0.example.org
    interfaces:
      - name: enp2s0
        macAddress: 98:af:65:a5:8d:01
    rootDeviceHints:
      wwnWithExtension: "wwn-with-extension-value"`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: Hosts[0].RootDeviceHints.WWNWithExtension: Forbidden: WWN extensions are not supported in root device hints",
		},
		{
			name: "unsupported wwn vendor extension root device hint",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - hostname: control-0.example.org
    interfaces:
      - name: enp2s0
        macAddress: 98:af:65:a5:8d:01
    rootDeviceHints:
      wwnVendorExtension: "wwn-with-vendor-extension-value"`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: Hosts[0].RootDeviceHints.WWNVendorExtension: Forbidden: WWN vendor extensions are not supported in root device hints",
		},
		{
			name: "node-hostname-and-role-are-not-required",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - interfaces:
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a`,

			expectedFound: true,
			expectedConfig: agentConfig().hosts(
				agentHost().interfaces(
					iface("enp3s1", "28:d2:44:d2:b2:1a"),
				)),
		},
		{
			name: "host-roles-have-correct-values",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - role: master
    interfaces:
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a
  - role: worker
    interfaces:
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1b`,

			expectedFound: true,
			expectedConfig: agentConfig().hosts(
				agentHost().role("master").interfaces(iface("enp3s1", "28:d2:44:d2:b2:1a")),
				agentHost().role("worker").interfaces(iface("enp3s1", "28:d2:44:d2:b2:1b")),
			),
		},
		{
			name: "host-roles-have-incorrect-values",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - role: invalid-role
    interfaces:
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: Hosts[0].Host: Forbidden: host role has incorrect value. Role must either be 'master' or 'worker'",
		},
		{
			name: "different-ifaces-same-host-cannot-have-same-mac",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - interfaces:
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a
      - name: enp3s2
        macAddress: 28:d2:44:d2:b2:1a`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: Hosts[0].Interfaces[1].macAddress: Invalid value: \"28:d2:44:d2:b2:1a\": duplicate MAC address found",
		},
		{
			name: "different-hosts-cannot-have-same-mac",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - interfaces:
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a
  - interfaces:
      - name: enp3s1
        macAddress: 28:d2:44:d2:b2:1a`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: Hosts[1].Interfaces[0].macAddress: Invalid value: \"28:d2:44:d2:b2:1a\": duplicate MAC address found",
		},
		{
			name: "invalid-mac",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: 192.168.111.80
hosts:
  - interfaces:
      - name: enp3s1
        macAddress: "000000"`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: Hosts[0].Interfaces[0].macAddress: Invalid value: \"000000\": address 000000: invalid MAC address",
		},
		{
			name: "empty-rendezvousIP",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0`,

			expectedFound:  true,
			expectedConfig: agentConfig().rendezvousIP(""),
		},
		{
			name: "invalid-rendezvousIP",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
rendezvousIP: not-a-valid-ip`,

			expectedFound: false,
			expectedError: "invalid Agent Config configuration: rendezvousIP: Invalid value: \"not-a-valid-ip\": \"not-a-valid-ip\" is not a valid IP",
		},
		{
			name: "invalid-additionalNTPSourceDomain",
			data: `
apiVersion: v1beta1
metadata:
  name: agent-config-cluster0
additionalNTPSources:
  - 0.fedora.pool.ntp.org
  - 1.fedora.pool.ntp.org
  - 192.168.111.14
  - fd10:39:192:1::1337
  - invalid_pool.ntp.org
rendezvousIP: 192.168.111.80`,
			expectedFound: false,
			expectedError: "invalid Agent Config configuration: AdditionalNTPSources[4]: Invalid value: \"invalid_pool.ntp.org\": NTP source is not a valid domain name nor a valid IP",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(agentConfigFilename).
				Return(
					&asset.File{
						Filename: agentConfigFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &AgentConfig{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				if tc.expectedConfig != nil {
					assert.Equal(t, tc.expectedConfig.build(), asset.Config, "unexpected Config in AgentConfig")
				}
			}
		})
	}

}

// AgentConfigBuilder it's a builder class to make it easier creating agent.Config instance
// used in the test cases
type AgentConfigBuilder struct {
	agent.Config
}

func agentConfig() *AgentConfigBuilder {
	return &AgentConfigBuilder{
		Config: agent.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "agent-config-cluster0",
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: agent.AgentConfigVersion,
			},
			RendezvousIP: "192.168.111.80",
		},
	}
}

func (acb *AgentConfigBuilder) build() *agent.Config {
	return &acb.Config
}

func (acb *AgentConfigBuilder) hosts(builders ...*AgentHostBuilder) *AgentConfigBuilder {

	hosts := []agent.Host{}
	for _, b := range builders {
		hosts = append(hosts, *b.build())
	}
	acb.Config.Hosts = hosts

	return acb
}

func (acb *AgentConfigBuilder) rendezvousIP(ip string) *AgentConfigBuilder {
	acb.Config.RendezvousIP = ip
	return acb
}

// AgentHostBuilder it's a builder class to make it easier creating agent.Host instances
// used in the test cases, as part of the agent.Config type
type AgentHostBuilder struct {
	agent.Host
}

func agentHost() *AgentHostBuilder {
	return &AgentHostBuilder{}
}

func defaultAgentHost(name string) *AgentHostBuilder {
	return agentHost().
		name(name).
		role("master").
		interfaces(
			iface("enp2s0", "98:af:65:a5:8d:01"),
			iface("enp3s1", "28:d2:44:d2:b2:1a"),
		).
		defaultRootDeviceHints().
		networkConfig("interfaces:")
}

func (ahb *AgentHostBuilder) build() *agent.Host {
	return &ahb.Host
}

func (ahb *AgentHostBuilder) name(name string) *AgentHostBuilder {
	ahb.Host.Hostname = name
	return ahb
}

func (ahb *AgentHostBuilder) role(role string) *AgentHostBuilder {
	ahb.Host.Role = role
	return ahb
}

func (ahb *AgentHostBuilder) interfaces(builders ...*InterfacetBuilder) *AgentHostBuilder {
	ifaces := []*aiv1beta1.Interface{}
	for _, b := range builders {
		ifaces = append(ifaces, b.build())
	}
	ahb.Host.Interfaces = ifaces
	return ahb
}

func (ahb *AgentHostBuilder) networkConfig(raw string) *AgentHostBuilder {
	ahb.Host.NetworkConfig = aiv1beta1.NetConfig{
		Raw: unmarshalJSON([]byte(raw)),
	}
	return ahb
}

// TODO: Create BaremetalRootDeviceHintsBuilder, for the current tests not required
func (ahb *AgentHostBuilder) defaultRootDeviceHints() *AgentHostBuilder {
	falseBool := false
	ahb.Host.RootDeviceHints = baremetal.RootDeviceHints{
		DeviceName:       "/dev/sda",
		HCTL:             "hctl-value",
		Model:            "model-value",
		Vendor:           "vendor-value",
		SerialNumber:     "serial-number-value",
		MinSizeGigabytes: 20,
		WWN:              "wwn-value",
		Rotational:       &falseBool,
	}
	return ahb
}

// InterfacetBuilder it's a builder class to make it easier creating aiv1beta1.Interface instances
// used in the test cases, as part of the agent.Config type
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
