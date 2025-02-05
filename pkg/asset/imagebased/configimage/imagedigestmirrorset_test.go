package configimage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	apicfgv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
)

func TestImageDigestMirrorSet_Generate(t *testing.T) {
	installConfigWithImageDigestSources := defaultInstallConfig()
	installConfigWithImageDigestSources.Config.ImageDigestSources = imageDigestSources()

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *apicfgv1.ImageDigestMirrorSet
	}{
		{
			name: "missing install config should not generate image-digest-mirror-set.json",
			dependencies: []asset.Asset{
				&InstallConfig{},
			},
		},
		{
			name: "missing ImageDigestSources in install config should not generate image-digest-mirror-set.json",
			dependencies: []asset.Asset{
				defaultInstallConfig(),
			},
		},
		{
			name: "valid configuration",
			dependencies: []asset.Asset{
				installConfigWithImageDigestSources,
			},
			expectedConfig: imageDigestMirrorSet(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &ImageDigestMirrorSet{}
			err := asset.Generate(context.TODO(), parents)

			switch {
			case tc.expectedError != "":
				assert.Equal(t, tc.expectedError, err.Error())
			case tc.expectedConfig == nil:
				assert.NoError(t, err)
				assert.Empty(t, asset.Files())
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-configuration/manifests/image-digest-mirror-set.json", configFile.Filename)

				var actualConfig = new(apicfgv1.ImageDigestMirrorSet)
				err = json.Unmarshal(configFile.Data, actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, actualConfig)
			}
		})
	}
}

func TestImageDigestMirrorSet_LoadedFromDisk(t *testing.T) {
	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig *apicfgv1.ImageDigestMirrorSet
	}{
		{
			name:           "valid-image-digest-mirror-set-file",
			data:           `[{"source":"test.registry.io","mirrors":["another.registry"]}]`,
			expectedFound:  true,
			expectedConfig: imageDigestMirrorSet(),
		},
		{
			name:          "not-json",
			data:          `This is not a JSON file`,
			expectedError: "failed to unmarshal cluster-configuration/manifests/image-digest-mirror-set.json: invalid JSON syntax",
		},
		{
			name:          "empty",
			data:          "",
			expectedError: "failed to unmarshal cluster-configuration/manifests/image-digest-mirror-set.json: invalid JSON syntax",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load cluster-configuration/manifests/image-digest-mirror-set.json file: fetch failed",
		},
		{
			name:          "unknown-field",
			data:          `[{"wrongField":"wrongValue"}]`,
			expectedError: "failed to unmarshal cluster-configuration/manifests/image-digest-mirror-set.json: unknown field \"[0].wrongField\"",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(imageDigestMirrorSetFilename).
				Return(
					&asset.File{
						Filename: imageDigestMirrorSetFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ImageDigestMirrorSet{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in ImageDigestSources")
			}
		})
	}
}

func imageDigestSources() []types.ImageDigestSource {
	return []types.ImageDigestSource{
		{
			Source:  "test.registry.io",
			Mirrors: []string{"another.registry"},
		},
	}
}

func imageDigestMirrorSet() *apicfgv1.ImageDigestMirrorSet {
	return convertIDSToIDMS(imageDigestSources())
}
