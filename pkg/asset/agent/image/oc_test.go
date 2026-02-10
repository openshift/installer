package image

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

func TestGetIcspContents(t *testing.T) {

	cases := []struct {
		name           string
		mirrorConfig   types.MirrorConfig
		expectedError  string
		expectedConfig string
	}{
		{
			name: "valid-config",
			mirrorConfig: types.MirrorConfig{
				{
					Location: "registry.ci.openshift.org/ocp/release",
					Mirrors:  []string{"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"},
				},
				{
					Location: "quay.io/openshift-release-dev/ocp-v4.0-art-dev",
					Mirrors:  []string{"virthost.ostest.test.metalkube.org:5000/localimages/local-release-image"},
				},
			},
			expectedConfig: "apiVersion: operator.openshift.io/v1alpha1\nkind: ImageContentSourcePolicy\nmetadata:\n  name: image-policy\nspec:\n  repositoryDigestMirrors:\n  - mirrors:\n    - virthost.ostest.test.metalkube.org:5000/localimages/local-release-image\n    source: registry.ci.openshift.org/ocp/release\n  - mirrors:\n    - virthost.ostest.test.metalkube.org:5000/localimages/local-release-image\n    source: quay.io/openshift-release-dev/ocp-v4.0-art-dev\n",
			expectedError:  "",
		},
		{
			name:           "empty-config",
			mirrorConfig:   types.MirrorConfig{},
			expectedConfig: "apiVersion: operator.openshift.io/v1alpha1\nkind: ImageContentSourcePolicy\nmetadata:\n  name: image-policy\nspec:\n  repositoryDigestMirrors: []\n",
			expectedError:  "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			contents, err := tc.mirrorConfig.GetICSPContents()
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedConfig, string(contents))
		})
	}
}
