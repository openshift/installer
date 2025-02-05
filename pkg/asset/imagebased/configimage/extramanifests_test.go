package configimage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestExtraManifests_Load(t *testing.T) {
	cases := []struct {
		name           string
		files          []string
		yamlFetchError error
		ymlFetchError  error
		jsonFetchError error

		expectedFound bool
		expectedFiles []string
		expectedError string
	}{
		{
			name:  "no-extra-manifests",
			files: []string{},

			expectedFound: false,
			expectedFiles: []string{},
		},
		{
			name:  "only-yaml",
			files: []string{"/extra-manifests/test.yaml"},

			expectedFound: true,
			expectedFiles: []string{"/extra-manifests/test.yaml"},
		},
		{
			name:  "only-yml",
			files: []string{"/extra-manifests/another-test.yml"},

			expectedFound: true,
			expectedFiles: []string{"/extra-manifests/another-test.yml"},
		},
		{
			name:  "only-json",
			files: []string{"/extra-manifests/test.json"},

			expectedFound: true,
			expectedFiles: []string{"/extra-manifests/test.json"},
		},
		{
			name: "all",
			files: []string{
				"/extra-manifests/test.yaml",
				"/extra-manifests/another-test.yml",
				"/extra-manifests/test.json",
			},

			expectedFound: true,
			expectedFiles: []string{
				"/extra-manifests/test.yaml",
				"/extra-manifests/another-test.yml",
				"/extra-manifests/test.json",
			},
		},
		{
			name:           "error",
			yamlFetchError: os.ErrNotExist,

			expectedError: "failed to load *.yaml files: file does not exist",
		},
		{
			name:          "error",
			ymlFetchError: os.ErrNotExist,

			expectedError: "failed to load *.yml files: file does not exist",
		},
		{
			name:           "error",
			jsonFetchError: os.ErrNotExist,

			expectedError: "failed to load *.json files: file does not exist",
		},
	}

	yamlPattern := "extra-manifests/*.yaml"
	ymlPattern := "extra-manifests/*.yml"
	jsonPattern := "extra-manifests/*.json"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			yamlFiles := []*asset.File{}
			ymlFiles := []*asset.File{}
			jsonFiles := []*asset.File{}
			for _, f := range tc.files {
				assetFile := &asset.File{
					Filename: f,
					Data:     []byte(f),
				}

				switch filepath.Ext(f) {
				case ".yaml":
					yamlFiles = append(yamlFiles, assetFile)
				case ".yml":
					ymlFiles = append(ymlFiles, assetFile)
				case ".json":
					jsonFiles = append(jsonFiles, assetFile)
				default:
					t.Error("invalid extension")
				}
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)

			if tc.yamlFetchError != nil {
				fileFetcher.EXPECT().FetchByPattern(yamlPattern).Return([]*asset.File{}, tc.yamlFetchError)
			} else {
				fileFetcher.EXPECT().FetchByPattern(yamlPattern).Return(yamlFiles, nil)

				if tc.ymlFetchError != nil {
					fileFetcher.EXPECT().FetchByPattern(ymlPattern).Return([]*asset.File{}, tc.ymlFetchError)
				} else {
					fileFetcher.EXPECT().FetchByPattern(ymlPattern).Return(ymlFiles, nil)

					if tc.jsonFetchError != nil {
						fileFetcher.EXPECT().FetchByPattern(jsonPattern).Return([]*asset.File{}, tc.jsonFetchError)
					} else {
						fileFetcher.EXPECT().FetchByPattern(jsonPattern).Return(jsonFiles, nil)
					}
				}
			}

			extraManifestsAsset := &ExtraManifests{}
			found, err := extraManifestsAsset.Load(fileFetcher)

			assert.Equal(t, tc.expectedFound, found)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedFiles), len(extraManifestsAsset.FileList))
				for _, f := range tc.expectedFiles {
					found := false
					for _, a := range extraManifestsAsset.FileList {
						if a.Filename == f {
							found = true
							break
						}
					}
					assert.True(t, found, fmt.Sprintf("Expected file %s not found", f))
				}
			}
		})
	}
}
