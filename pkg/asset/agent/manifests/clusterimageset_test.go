package manifests

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestClusterImageSet_Generate(t *testing.T) {
	currentRelease, err := releaseimage.Default()
	assert.NoError(t, err)

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig *hivev1.ClusterImageSet
	}{
		{
			name: "missing install config",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{},
				&releaseimage.Image{},
			},
			expectedError: "missing configuration or manifest file",
		},
		{
			name: "invalid ClusterImageSet configuration",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfig(),
				&releaseimage.Image{},
			},
			expectedError: "invalid ClusterImageSet configuration: Spec.ReleaseImage: Forbidden: value must be equal to " + currentRelease,
		},
		{
			name: "valid configuration",
			dependencies: []asset.Asset{
				getValidOptionalInstallConfig(),
				&releaseimage.Image{
					PullSpec: currentRelease,
				},
			},
			expectedConfig: &hivev1.ClusterImageSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "openshift-was not built correctly",
					Namespace: getObjectMetaNamespace(getValidOptionalInstallConfig()),
				},
				Spec: hivev1.ClusterImageSetSpec{
					ReleaseImage: currentRelease,
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &ClusterImageSet{}
			err := asset.Generate(parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, asset.Config)
				assert.NotEmpty(t, asset.Files())

				configFile := asset.Files()[0]
				assert.Equal(t, "cluster-manifests/cluster-image-set.yaml", configFile.Filename)

				var actualConfig hivev1.ClusterImageSet
				err = yaml.Unmarshal(configFile.Data, &actualConfig)
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedConfig, actualConfig)
			}

		})
	}

}

func TestClusterImageSet_LoadedFromDisk(t *testing.T) {

	currentRelease, err := releaseimage.Default()
	assert.NoError(t, err)

	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig *hivev1.ClusterImageSet
	}{
		{
			name: "valid-config-file",
			data: `
metadata:
  name: openshift-v4.10.0
spec:
  releaseImage: ` + currentRelease,
			expectedFound: true,
			expectedConfig: &hivev1.ClusterImageSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "openshift-v4.10.0",
				},
				Spec: hivev1.ClusterImageSetSpec{
					ReleaseImage: currentRelease,
				},
			},
		},
		{
			name: "different-version-not-supported",
			data: `
metadata:
  name: openshift-v4.10.0
spec:
  releaseImage: 99.999`,
			expectedError: fmt.Sprintf("invalid ClusterImageSet configuration: Spec.ReleaseImage: Forbidden: value must be equal to %s", currentRelease),
		},
		{
			name:          "not-yaml",
			data:          `This is not a yaml file`,
			expectedError: "failed to unmarshal cluster-manifests/cluster-image-set.yaml: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type v1.ClusterImageSet",
		},
		{
			name:          "empty",
			data:          "",
			expectedError: fmt.Sprintf("invalid ClusterImageSet configuration: Spec.ReleaseImage: Forbidden: value must be equal to %s", currentRelease),
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load cluster-manifests/cluster-image-set.yaml file: fetch failed",
		},
		{
			name: "unknown-field",
			data: `
metadata:
  name: test-cluster-image-set
  namespace: cluster0
spec:
  wrongField: wrongValue`,
			expectedError: "failed to unmarshal cluster-manifests/cluster-image-set.yaml: error unmarshaling JSON: while decoding JSON: json: unknown field \"wrongField\"",
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
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in ClusterImageSet")
			}
		})
	}

}
