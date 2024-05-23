package configimage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
)

func TestImageDigestSources_Generate(t *testing.T) {
	installConfigWithImageDigestSources := defaultInstallConfig()
	installConfigWithImageDigestSources.Config.ImageDigestSources = imageDigestSources()

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig []types.ImageDigestSource
	}{
		{
			name: "missing install config should not generate image-digest-sources.json",
			dependencies: []asset.Asset{
				&InstallConfig{},
			},
		},
		{
			name: "missing ImageDigestSources in install config should not generate image-digest-sources.json",
			dependencies: []asset.Asset{
				defaultInstallConfig(),
			},
		},
		{
			name: "valid configuration",
			dependencies: []asset.Asset{
				installConfigWithImageDigestSources,
			},
			expectedConfig: imageDigestSources(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &ImageDigestSources{}
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
				assert.Equal(t, "cluster-configuration/manifests/image-digest-sources.json", configFile.Filename)

				var actualConfig []types.ImageDigestSource
				err = json.Unmarshal(configFile.Data, &actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, actualConfig)
			}
		})
	}
}

func TestImageDigestSources_LoadedFromDisk(t *testing.T) {
	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig []types.ImageDigestSource
	}{
		{
			name:           "valid-image-digest-sources-file",
			data:           `[{"source":"test.registry.io","mirrors":["another.registry"]}]`,
			expectedFound:  true,
			expectedConfig: imageDigestSources(),
		},
		{
			name:          "not-json",
			data:          `This is not a JSON file`,
			expectedError: "failed to unmarshal cluster-configuration/manifests/image-digest-sources.json: invalid JSON syntax",
		},
		{
			name:          "empty",
			data:          "",
			expectedError: "failed to unmarshal cluster-configuration/manifests/image-digest-sources.json: invalid JSON syntax",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load cluster-configuration/manifests/image-digest-sources.json file: fetch failed",
		},
		{
			name:          "unknown-field",
			data:          `[{"wrongField":"wrongValue"}]`,
			expectedError: "failed to unmarshal cluster-configuration/manifests/image-digest-sources.json: unknown field \"[0].wrongField\"",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(imageDigestSourcesFilename).
				Return(
					&asset.File{
						Filename: imageDigestSourcesFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ImageDigestSources{}
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
