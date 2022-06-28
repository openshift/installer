package agentconfig

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types/agent"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAgentConfig_LoadedFromDisk(t *testing.T) {
	falseBool := false
	falsePtr := &falseBool

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig *agent.Config
	}{
		{
			name: "valid-config-single-node",
			data: `
metadata:
  name: agent-config-cluster0
spec:
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
          macAddress: 28:d2:44:d2:b2:1a`,
			expectedFound: true,
			expectedConfig: &agent.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: "agent-config-cluster0",
				},
				Spec: agent.Spec{
					RendezvousIP: "192.168.111.80",
					Hosts: []agent.Host{
						{
							Hostname: "control-0.example.org",
							Role:     "master",
							RootDeviceHints: baremetal.RootDeviceHints{
								DeviceName:       "/dev/sda",
								HCTL:             "hctl-value",
								Model:            "model-value",
								Vendor:           "vendor-value",
								SerialNumber:     "serial-number-value",
								MinSizeGigabytes: 20,
								WWN:              "wwn-value",
								Rotational:       falsePtr,
							},
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
						},
					},
				},
			},
		},
		{
			name: "valid-config-multiple-nodes",
			data: `
metadata:
  name: agent-config-cluster0
spec:
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
    - hostname: control-1.example.org
      role: master
      interfaces:
        - name: enp2s0
          macAddress: 98:af:65:a5:8d:02
        - name: enp3s1
          macAddress: 28:d2:44:d2:b2:1b`,
			expectedFound: true,
			expectedConfig: &agent.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: "agent-config-cluster0",
				},
				Spec: agent.Spec{
					RendezvousIP: "192.168.111.80",
					Hosts: []agent.Host{
						{
							Hostname: "control-0.example.org",
							Role:     "master",
							RootDeviceHints: baremetal.RootDeviceHints{
								DeviceName:       "/dev/sda",
								HCTL:             "hctl-value",
								Model:            "model-value",
								Vendor:           "vendor-value",
								SerialNumber:     "serial-number-value",
								MinSizeGigabytes: 20,
								WWN:              "wwn-value",
								Rotational:       falsePtr,
							},
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
						},
						{
							Hostname: "control-1.example.org",
							Role:     "master",
							Interfaces: []*aiv1beta1.Interface{
								{
									Name:       "enp2s0",
									MacAddress: "98:af:65:a5:8d:02",
								},
								{
									Name:       "enp3s1",
									MacAddress: "28:d2:44:d2:b2:1b",
								},
							},
						},
					},
				},
			},
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: "failed to unmarshal agent-config.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type agent.Config",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load agent-config.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
		metadata:
		  name: agent-config-wrong
		spec:
		  wrongField: wrongValue`,
			expectedError: "failed to unmarshal agent-config.yaml: error converting YAML to JSON: yaml: line 2: found character that cannot start any token",
		},
		{
			name: "interface-missing-mac-address-error",
			data: `
metadata:
  name: agent-config-cluster0
spec:
  rendezvousIP: 192.168.111.80
  hosts:
    - hostname: control-0.example.org
      interfaces:
        - name: enp2s0
        - name: enp3s1
          macAddress: 28:d2:44:d2:b2:1a`,
			expectedError: "invalid Agent Config configuration: Spec.Hosts[0].Interfaces[0].macAddress: Required value: each interface must have a MAC address defined",
		},
		{
			name: "unsupported wwn extension root device hint",
			data: `
metadata:
  name: agent-config-cluster0
spec:
  rendezvousIP: 192.168.111.80
  hosts:
    - hostname: control-0.example.org
      interfaces:
        - name: enp2s0
          macAddress: 98:af:65:a5:8d:01
      rootDeviceHints:
        wwnWithExtension: "wwn-with-extension-value"`,
			expectedError: "invalid Agent Config configuration: Spec.Hosts[0].RootDeviceHints.WWNWithExtension: Forbidden: WWN extensions are not supported in root device hints",
		},
		{
			name: "unsupported wwn vendor extension root device hint",
			data: `
metadata:
  name: agent-config-cluster0
spec:
  rendezvousIP: 192.168.111.80
  hosts:
    - hostname: control-0.example.org
      interfaces:
        - name: enp2s0
          macAddress: 98:af:65:a5:8d:01
      rootDeviceHints:
        wwnVendorExtension: "wwn-with-vendor-extension-value"`,
			expectedError: "invalid Agent Config configuration: Spec.Hosts[0].RootDeviceHints.WWNVendorExtension: Forbidden: WWN vendor extensions are not supported in root device hints",
		},
		{
			name: "node-hostname-and-role-are-not-required",
			data: `
metadata:
  name: agent-config-cluster0
spec:
  rendezvousIP: 192.168.111.80
  hosts:
    - interfaces:
        - name: enp3s1
          macAddress: 28:d2:44:d2:b2:1a`,
			expectedFound: true,
			expectedConfig: &agent.Config{
				ObjectMeta: metav1.ObjectMeta{
					Name: "agent-config-cluster0",
				},
				Spec: agent.Spec{
					RendezvousIP: "192.168.111.80",
					Hosts: []agent.Host{
						{
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

			asset := &Asset{}
			found, err := asset.Load(fileFetcher)
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in AgentConfig")
			}
		})
	}

}
