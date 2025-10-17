package configimage

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types/imagebased"
)

func TestImageBasedConfig_LoadedFromDisk(t *testing.T) {
	skipTestIfnmstatectlIsMissing(t)

	cases := []struct {
		name       string
		data       string
		fetchError error

		expectedError  string
		expectedFound  bool
		expectedConfig *imagebased.Config
	}{
		{
			name: "valid-config",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-config
hostname: somehostname
releaseRegistry: quay.io
networkConfig:
  interfaces:
    - name: eth0
      type: ethernet
      state: up
      mac-address: 00:00:00:00:00:00
      ipv4:
        enabled: true
        address:
          - ip: 192.168.122.2
            prefix-length: 23
        dhcp: false`,
			expectedFound:  true,
			expectedConfig: imageBasedConfig().Config,
		},
		{
			name: "not-yaml",
			data: `This is not a yaml file`,

			expectedFound: false,
			expectedError: "failed to unmarshal image-based-config.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type imagebased.Config",
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
			expectedError: "failed to load image-based-config.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-config-wrong
wrongField: wrongValue`,

			expectedFound: false,
			expectedError: "failed to unmarshal image-based-config.yaml: error unmarshaling JSON: while decoding JSON: json: unknown field \"wrongField\"",
		},
		{
			name: "invalid-nmstate-config",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-config
hostname: somehostname
releaseRegistry: quay.io
networkConfig:
  interfaces:
    - name: eth0
      type: ethernet
      state: invalid`,

			expectedFound: false,
			expectedError: "interfaces: unknown variant `invalid`",
		},
		{
			name: "invalid-additional-ntp-sources",
			data: `
apiVersion: v1beta1
metadata:
  name: image-based-config
hostname: somehostname
releaseRegistry: quay.io
additionalNTPSources:
  - invalid-domain@
networkConfig:
  interfaces:
    - name: eth0
      type: ethernet
      state: up`,

			expectedFound: false,
			expectedError: "invalid Image-based Config configuration: AdditionalNTPSources: Invalid value: \"invalid-domain@\": NTP server is not a valid IP or domain name",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(imageBasedConfigFilename).
				Return(
					&asset.File{
						Filename: imageBasedConfigFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ImageBasedConfig{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				if tc.expectedConfig != nil {
					assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in ImageBasedConfig")
				}
			}
		})
	}
}

func skipTestIfnmstatectlIsMissing(t *testing.T) {
	t.Helper()

	_, execErr := exec.LookPath("nmstatectl")
	if execErr != nil {
		t.Skip("No nmstatectl binary available")
	}
}
