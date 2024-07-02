package image

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

func TestRegistriesConf_Generate(t *testing.T) {
	cases := []struct {
		name         string
		dependencies []asset.Asset

		expectedError string
		expectedData  []byte
	}{
		{
			name: "default",
			dependencies: []asset.Asset{
				&ImageBasedInstallationConfig{
					Config: ibiConfig().
						imageDigestSources([]types.ImageDigestSource{
							{
								Source:  "quay.io",
								Mirrors: []string{"mirror-quay.io"},
							},
						}).
						build(),
				},
			},

			expectedData: []byte("credential-helpers = []\nshort-name-mode = \"\"\nunqualified-search-registries = []\n\n[[registry]]\n  location = \"quay.io\"\n  mirror-by-digest-only = true\n  prefix = \"\"\n\n  [[registry.mirror]]\n    location = \"mirror-quay.io\"\n"),
		},
		{
			name: "deprecated image content sources",
			dependencies: []asset.Asset{
				&ImageBasedInstallationConfig{
					Config: ibiConfig().
						deprecatedImageContentSources([]types.ImageContentSource{
							{
								Source:  "quay.io",
								Mirrors: []string{"mirror-quay.io"},
							},
						}).
						build(),
				},
			},

			expectedData: []byte("credential-helpers = []\nshort-name-mode = \"\"\nunqualified-search-registries = []\n\n[[registry]]\n  location = \"quay.io\"\n  mirror-by-digest-only = true\n  prefix = \"\"\n\n  [[registry.mirror]]\n    location = \"mirror-quay.io\"\n"),
		},
		{
			name: "both image digest sources and deprecated image contet sources",
			dependencies: []asset.Asset{
				&ImageBasedInstallationConfig{
					Config: ibiConfig().
						imageDigestSources([]types.ImageDigestSource{
							{
								Source:  "quay.io",
								Mirrors: []string{"mirror-quay.io"},
							},
						}).
						deprecatedImageContentSources([]types.ImageContentSource{
							{
								Source:  "quay.io",
								Mirrors: []string{"mirror-quay.io"},
							},
						}).
						build(),
				},
			},

			expectedError: "invalid imagebased-installation-config.yaml, cannot set imageContentSources and imageDigestSources at the same time",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			registriesConf := &RegistriesConf{}
			err := registriesConf.Generate(parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, registriesConf.Data)
			}
		})
	}
}
