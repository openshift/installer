package manifests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestExtraManifests_Load(t *testing.T) {
	cases := []struct {
		name       string
		files      []string
		fetchError error

		expectedFound bool
		expectedFiles []string
		expectedError string
	}{
		{
			name:  "no-extras",
			files: []string{},

			expectedFound: false,
			expectedFiles: []string{},
		},
		{
			name:  "just-yaml",
			files: []string{"/openshift/test-configmap.yaml"},

			expectedFound: true,
			expectedFiles: []string{"/openshift/test-configmap.yaml"},
		},
		{
			name:  "just-yml",
			files: []string{"/openshift/another-test-configmap.yml"},

			expectedFound: true,
			expectedFiles: []string{"/openshift/another-test-configmap.yml"},
		},
		{
			name: "mixed",
			files: []string{
				"/openshift/test-configmap.yaml",
				"/openshift/another-test-configmap.yml",
			},

			expectedFound: true,
			expectedFiles: []string{
				"/openshift/test-configmap.yaml",
				"/openshift/another-test-configmap.yml",
			},
		},
		{
			name:       "error",
			fetchError: os.ErrNotExist,

			expectedError: "failed to load *.yaml files: file does not exist",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			yamlFiles := []*asset.File{}
			ymlFiles := []*asset.File{}
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
				default:
					t.Error("extension not valid")
				}
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByPattern("openshift/*.yaml").Return(
				yamlFiles,
				tc.fetchError,
			)
			if tc.fetchError == nil {
				fileFetcher.EXPECT().FetchByPattern("openshift/*.yml").Return(
					ymlFiles,
					nil,
				)
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
