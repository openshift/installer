package image

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
)

func TestPostDeployment_Load(t *testing.T) {
	cases := []struct {
		name       string
		fetchError error

		expectedFound bool
		expectedFile  string
		expectedError string
	}{
		{
			name:       "no-post-deployment-script",
			fetchError: os.ErrNotExist,

			expectedFound: false,
		},
		{
			name: "post-deployment-script",

			expectedFound: true,
			expectedFile:  "post.sh",
		},
		{
			name:       "error",
			fetchError: os.ErrPermission,

			expectedError: "failed to load the post.sh file: permission denied",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assetFile := &asset.File{
				Filename: postDeploymentFilename,
				Data:     []byte("test-bash-script"),
			}

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)

			if tc.fetchError != nil {
				fileFetcher.EXPECT().FetchByName("post.sh").Return(nil, tc.fetchError)
			} else {
				fileFetcher.EXPECT().FetchByName("post.sh").Return(assetFile, nil)
			}

			postDeploymentAsset := &PostDeployment{}
			found, err := postDeploymentAsset.Load(fileFetcher)

			assert.Equal(t, tc.expectedFound, found)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)

				if tc.expectedFile != "" {
					assert.True(t,
						tc.expectedFile == postDeploymentAsset.File.Filename,
						fmt.Sprintf("Expected file %s not found", tc.expectedFile))
				}
			}
		})
	}
}
