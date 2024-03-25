package ibi

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift-kni/lifecycle-agent/api/ibiconfig"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types/ibi"
)

func TestImageBasedInstallConfig_LoadedFromDisk(t *testing.T) {
	cases := []struct {
		name       string
		data       string
		fetchError error

		expectedError  string
		expectedFound  bool
		expectedConfig *ImageBasedInstallConfigBuilder
	}{
		{
			name: "valid-config-single-node",
			data: `
apiVersion: v1beta1
metadata:
  name: ibi-config
  namespace: cluster0
seedImage: quay.io/openshift-kni/seed-image:4.16.0
seedVersion: 4.16.0
pullSecretFile: pull-secret.json
authFile: auth-file.json
sshPublicKeyFile: id_rsa.pub
rhcosLiveIso: https://mirror.openshift.com/pub/openshift-v4/amd64/dependencies/rhcos/latest/rhcos-live.x86_64.iso
installationDisk: vda
useContainersFolder: false
precacheBestEffort: false
precacheDisabled: false`,
			expectedFound:  true,
			expectedConfig: ibiConfig(),
		},
		{
			name: "not-yaml",
			data: `This is not a yaml file`,

			expectedFound: false,
			expectedError: "failed to unmarshal ibi-config.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type ibi.Config",
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
			expectedError: "failed to load ibi-config.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
apiVersion: v1beta1
metadata:
  name: ibi-config-unknown-field
wrongField: wrongValue`,

			expectedFound: false,
			expectedError: "failed to unmarshal ibi-config.yaml: error unmarshaling JSON: while decoding JSON: json: unknown field \"wrongField\"",
		},
		{
			name: "empty-authFile",
			data: `
apiVersion: v1beta1
metadata:
  name: ibi-config`,

			expectedFound: false,
			expectedError: "invalid Image-based Install Config configuration: []: Invalid value: \"\": authFile is required",
		},
		{
			name: "empty-seedImage",
			data: `
apiVersion: v1beta1
metadata:
  name: ibi-config
authFile: "auth-file.json"`,

			expectedFound: false,
			expectedError: "invalid Image-based Install Config configuration: []: Invalid value: \"\": seedImage is required",
		},
		{
			name: "empty-seedVersion",
			data: `
apiVersion: v1beta1
metadata:
  name: ibi-config
authFile: "auth-file.json"
seedImage: "quay.io/openshift-kni/seed-image:4.16.0"`,

			expectedFound: false,
			expectedError: "invalid Image-based Install Config configuration: []: Invalid value: \"\": seedVersion is required",
		},
		{
			name: "empty-installationDisk",
			data: `
apiVersion: v1beta1
metadata:
  name: ibi-config
authFile: "auth-file.json"
seedImage: "quay.io/openshift-kni/seed-image:4.16.0"
seedVersion: "4.16.0"`,

			expectedFound: false,
			expectedError: "invalid Image-based Install Config configuration: []: Invalid value: \"\": installationDisk is required",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(ibiConfigFilename).
				Return(
					&asset.File{
						Filename: ibiConfigFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ImageBasedInstallConfig{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found)
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				if tc.expectedConfig != nil {
					assert.Equal(t, tc.expectedConfig.build(), asset.Config, "unexpected Config in ImageBasedInstallConfig")
				}
			}
		})
	}
}

// ImageBasedInstallConfigBuilder it's a builder class to make it easier
// creating ibi.Config instance used in the test cases.
type ImageBasedInstallConfigBuilder struct {
	ibi.Config
}

func ibiConfig() *ImageBasedInstallConfigBuilder {
	return &ImageBasedInstallConfigBuilder{
		Config: ibi.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ibi-config",
				Namespace: "cluster0",
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: ibi.ImageBasedInstallConfigVersion,
			},
			IBIPrepareConfig: ibiconfig.IBIPrepareConfig{
				SSHPublicKeyFile:     "id_rsa.pub",
				RHCOSLiveISO:         "https://mirror.openshift.com/pub/openshift-v4/amd64/dependencies/rhcos/latest/rhcos-live.x86_64.iso",
				SeedImage:            "quay.io/openshift-kni/seed-image:4.16.0",
				SeedVersion:          "4.16.0",
				AuthFile:             "auth-file.json",
				PullSecretFile:       "pull-secret.json",
				InstallationDisk:     "vda",
				PrecacheBestEffort:   false,
				PrecacheDisabled:     false,
				ExtraPartitionStart:  "",
				ExtraPartitionLabel:  "",
				ExtraPartitionNumber: 0,
				Shutdown:             false,
				UseContainersFolder:  false,
			},
		},
	}
}

func (icb *ImageBasedInstallConfigBuilder) build() *ibi.Config {
	return &icb.Config
}
