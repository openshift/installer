package manifests

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterImageSet_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  bool
		expectedConfig *hivev1.ClusterImageSet
	}{
		{
			name: "valid-config-file",
			data: `
metadata:
  name: openshift-v4.10.0
spec:
  releaseImage: quay.io:443/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64`,
			expectedFound: true,
			expectedConfig: &hivev1.ClusterImageSet{
				ObjectMeta: v1.ObjectMeta{
					Name: "openshift-v4.10.0",
				},
				Spec: hivev1.ClusterImageSetSpec{
					ReleaseImage: "quay.io:443/openshift-release-dev/ocp-release:4.10.0-rc.1-x86_64",
				},
			},
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: true,
		},
		{
			name:           "empty",
			data:           "",
			expectedFound:  true,
			expectedConfig: &hivev1.ClusterImageSet{},
			expectedError:  false,
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
		{
			name: "unknown-field",
			data: `
metadata:
  name: test-cluster-image-set
  namespace: cluster0
spec:
  wrongField: wrongValue`,
			expectedError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(clusterImageSetFilename).
				Return(
					&asset.File{
						Filename: clusterImageSetFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &ClusterImageSet{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError {
				assert.Error(t, err, "expected error from Load")
			} else {
				assert.NoError(t, err, "unexpected error from Load")
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in ClusterImageSet")
			}
		})
	}

}
