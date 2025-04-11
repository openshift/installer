package joiner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestAddNodesConfig_Load(t *testing.T) {
	cases := []struct {
		name               string
		addNodesParamsData string
		nodesConfigData    string
		expectedError      string
	}{
		{
			name: "default",
			nodesConfigData: `hosts:
- hostname: master-0
  interfaces:
  - name: eth0
    macAddress: 00:ef:29:72:b9:771`,
		},
		{
			name:          "empty nodes-config.yaml",
			expectedError: "hosts: Required value: at least one host must be defined",
		},
		{
			name: "ssh key",
			nodesConfigData: `hosts:
- hostname: master-0
  interfaces:
  - name: eth0
    macAddress: 00:ef:29:72:b9:771
sshKey: "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q=="`,
		},
		{
			name: "invalid ssh key",
			nodesConfigData: `hosts:
- hostname: master-0
  interfaces:
  - name: eth0
    macAddress: 00:ef:29:72:b9:771
sshKey: "not a valid ssh key"`,
			expectedError: "sshKey: Invalid value: \"not a valid ssh key\": ssh: no key found",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			if tc.addNodesParamsData == "" {
				tc.addNodesParamsData = "{}"
			}

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(addNodesParamsFile).
				Return(
					&asset.File{
						Filename: addNodesParamsFile,
						Data:     []byte(tc.addNodesParamsData)},
					nil,
				)
			fileFetcher.EXPECT().FetchByName(nodesConfigFilename).
				Return(
					&asset.File{
						Filename: nodesConfigFilename,
						Data:     []byte(tc.nodesConfigData)},
					nil,
				)

			addNodesConfig := &AddNodesConfig{}
			_, err := addNodesConfig.Load(fileFetcher)
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
