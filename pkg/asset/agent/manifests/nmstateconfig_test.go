package manifests

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestNMStateConfig_Generate(t *testing.T) {

	// TODO

}

func TestNMStateConfig_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name               string
		data               string
		fetchError         error
		expectedFound      bool
		expectedError      bool
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
			expectedError:      true,
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
			expectedError:      true,
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
			expectedError: true,
		},

		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: true,
		},
		{
			name:          "empty",
			data:          "",
			expectedFound: true,
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: true,
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
			if tc.expectedError {
				assert.Error(t, err, "expected error from Load")
			} else {
				assert.NoError(t, err, "unexpected error from Load")
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.StaticNetworkConfig, "unexpected Config in NMStateConfig")
			}
		})
	}
}
