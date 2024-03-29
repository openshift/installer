package joiner

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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
			expectedError: "invalid nodes configuration: hosts: Required value: at least one host must be defined",
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
