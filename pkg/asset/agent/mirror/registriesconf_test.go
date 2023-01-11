package mirror

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/types"
)

func TestRegistriesConf_Generate(t *testing.T) {

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedConfig string
	}{
		{
			name: "missing-config",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{},
				&releaseimage.Image{},
			},
			expectedConfig: defaultRegistriesConf,
		},
		{
			name: "default",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{
					Supplied: true,
					InstallConfig: installconfig.InstallConfig{
						Config: &types.InstallConfig{
							ObjectMeta: v1.ObjectMeta{
								Namespace: "cluster-0",
							},
						},
					},
				},
				&releaseimage.Image{
					PullSpec: "registry.ci.openshift.org/ocp/release:4.11.0-0.ci-2022-05-16-202609",
				},
			},
			expectedConfig: defaultRegistriesConf,
		},
		{
			name: "invalid-config-image-content-source-does-not-match-releaseImage",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{
					Supplied: true,
					InstallConfig: installconfig.InstallConfig{
						Config: &types.InstallConfig{
							ObjectMeta: v1.ObjectMeta{
								Namespace: "cluster-0",
							},
							ImageContentSources: []types.ImageContentSource{
								{
									Source: "registry.ci.openshift.org/origin/release",
									Mirrors: []string{
										"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image",
									},
								},
								{
									Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
									Mirrors: []string{
										"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image",
									},
								},
							},
						},
					},
				},
				&releaseimage.Image{
					PullSpec: "registry.ci.openshift.org/ocp/release:4.11.0-0.ci-2022-05-16-202609",
				},
			},
			expectedConfig: "unqualified-search-registries = []\n\n[[registry]]\n  location = \"registry.ci.openshift.org/origin/release\"\n  mirror-by-digest-only = true\n  prefix = \"\"\n\n  [[registry.mirror]]\n    location = \"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image\"\n\n[[registry]]\n  location = \"quay.io/openshift-release-dev/ocp-v4.0-art-dev\"\n  mirror-by-digest-only = true\n  prefix = \"\"\n\n  [[registry.mirror]]\n    location = \"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image\"\n",
		},
		{
			name: "valid-image-content-sources",
			dependencies: []asset.Asset{
				&agent.OptionalInstallConfig{
					Supplied: true,
					InstallConfig: installconfig.InstallConfig{
						Config: &types.InstallConfig{
							ObjectMeta: v1.ObjectMeta{
								Namespace: "cluster-0",
							},
							ImageContentSources: []types.ImageContentSource{
								{
									Source: "registry.ci.openshift.org/ocp/release",
									Mirrors: []string{
										"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image",
									},
								},
								{
									Source: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
									Mirrors: []string{
										"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image",
									},
								},
							},
						},
					},
				},
				&releaseimage.Image{
					PullSpec: "registry.ci.openshift.org/ocp/release:4.11.0-0.ci-2022-05-16-202609",
				},
			},
			expectedConfig: `unqualified-search-registries = []

[[registry]]
  location = "registry.ci.openshift.org/ocp/release"
  mirror-by-digest-only = true
  prefix = ""

  [[registry.mirror]]
    location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"

[[registry]]
  location = "quay.io/openshift-release-dev/ocp-v4.0-art-dev"
  mirror-by-digest-only = true
  prefix = ""

  [[registry.mirror]]
    location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"
`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &RegistriesConf{}
			err := asset.Generate(parents)

			assert.NoError(t, err)

			files := asset.Files()
			assert.Len(t, files, 1)
			assert.Equal(t, tc.expectedConfig, string(files[0].Data))
		})
	}
}

func TestRegistries_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name          string
		data          string
		fetchError    error
		expectedFound bool
		expectedError string
	}{
		{
			name: "location-key-not-found-in-config-file",
			data: `
[[registry]]
wrongLocationkey = "registry.ci.openshift.org/origin/release" 
mirror-by-digest-only = false

[[registry.mirror]]
location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"

[[registry]]
locations = "quay.io/openshift-release-dev/ocp-v4.0-art-dev"
mirror-by-digest-only = false

[[registry.mirror]]
location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"`,
			expectedFound: true,
		},
		{
			name: "valid-config-file",
			data: `
[[registry]]
location = "registry.ci.openshift.org/origin/release" 
mirror-by-digest-only = false

[[registry.mirror]]
location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"

[[registry]]
location = "quay.io/openshift-release-dev/ocp-v4.0-art-dev"
mirror-by-digest-only = false

[[registry.mirror]]
location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"`,
			expectedFound: true,
			expectedError: "",
		},
		{
			name: "location-does-not-match-with-releaseImage",
			data: `
[[registry]]
location = "registry.ci.openshift.org/ocp/release" 
mirror-by-digest-only = false

[[registry.mirror]]
location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"

[[registry]]
location = "quay.io/openshift-release-dev/ocp-v4.0-art-dev"
mirror-by-digest-only = false

[[registry.mirror]]
location = "virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"`,
			expectedFound: true,
			expectedError: "",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load mirror/registries.conf file: fetch failed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(RegistriesConfFilename).
				Return(
					&asset.File{
						Filename: RegistriesConfFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &RegistriesConf{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
