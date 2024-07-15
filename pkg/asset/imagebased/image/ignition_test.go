package image

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/containers/image/v5/pkg/sysregistriesv2"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

var (
	testCert = `-----BEGIN CERTIFICATE-----
MIICYTCCAcqgAwIBAgIJAI2kA+uXAbhOMA0GCSqGSIb3DQEBCwUAMEgxCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJDQTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzEUMBIG
A1UECgwLUmVkIEhhdCBJbmMwHhcNMTkwMjEyMTkzMjUzWhcNMTkwMjEzMTkzMjUz
WjBIMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFu
Y2lzY28xFDASBgNVBAoMC1JlZCBIYXQgSW5jMIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQC+HOC0mKig/oINAKPo88LqxDJ4l7lozdLtp5oGeqWrLUXSfkvXAkQY
2QYdvPAjpRfH7Ii7G0Asx+HTKdvula7B5fXDjc6NYKuEpTJZRV1ugntI97bozF/E
C2BBmxxEnJN3+Xe8RYXMjz5Q4aqPw9vZhlWN+0hrREl1Ea/zHuWFIQIDAQABo1Mw
UTAdBgNVHQ4EFgQUvTS1XjlvOdsufSyWxukyQu3LriEwHwYDVR0jBBgwFoAUvTS1
XjlvOdsufSyWxukyQu3LriEwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsF
AAOBgQB9gFcOXnzJrM65QqxeCB9Z5l5JMjp45UFC9Bj2cgwDHP80Zvi4omlaacC6
aavmnLd67zm9PbYDWRaOIWAMeB916Iwaw/v6I0jwhAk/VxX5Fl6cGlZu9jZ3zbFE
2sDqkwzIuSjCG2A23s6d4M1S3IXCCydoCSLMu+WhLkbboK6jEg==
-----END CERTIFICATE-----
`
)

func TestIgnition_Generate(t *testing.T) {
	// This patch currently allows testing the Ignition asset using the embedded resources.
	// TODO: Replace it by mocking the filesystem in bootstrap.AddStorageFiles()
	workingDirectory, err := os.Getwd()
	assert.NoError(t, err)
	err = os.Chdir(path.Join(workingDirectory, "../../../../data"))
	assert.NoError(t, err)

	registriesConf := &RegistriesConf{
		Config: &sysregistriesv2.V2RegistriesConf{
			Registries: []sysregistriesv2.Registry{
				{
					Endpoint:           sysregistriesv2.Endpoint{Location: "quay.io"},
					MirrorByDigestOnly: true,
					Mirrors:            []sysregistriesv2.Endpoint{{Location: "mirror-quay.io"}},
				},
			},
		},
	}
	registriesConfData, err := toml.Marshal(registriesConf.Config)
	assert.NoError(t, err)

	testNetworkConfig := "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.21\n      prefix-length: 24\n    enabled: true\n  mac-address: 00:01:02:03:04:05\n  name: eth0\n  state: up\n  type: ethernet\n\n"

	ignitionConfigOverride := `{"ignition": {"version": "3.2.0"}, "storage": {"files": [{"path": "/tmp/example", "contents": {"source": "data:text/plain;charset=utf-8;base64,aGVscGltdHJhcHBlZGluYXN3YWdnZXJzcGVj"}}]}}`

	postDeploymentScript := "#!/bin/bash\necho \"Post image-based installation script\""

	cases := []struct {
		name         string
		dependencies []asset.Asset

		expectedError        string
		expectedFiles        map[string]string
		expectedSystemdUnits map[string]bool
	}{
		{
			name: "comprehensive configuration",
			dependencies: []asset.Asset{
				&ImageBasedInstallationConfig{
					Config: ibiConfig().
						imageDigestSources([]types.ImageDigestSource{
							{
								Source:  "quay.io",
								Mirrors: []string{"mirror-quay.io"},
							},
						}).
						additionalTrustBundle(testCert).
						networkConfig(aiv1beta1.NetConfig{Raw: []byte(testNetworkConfig)}).
						ignitionConfigOverride(ignitionConfigOverride).
						build(),
				},
				&RegistriesConf{Data: registriesConfData},
				&PostDeployment{File: &asset.File{
					Filename: "post.sh",
					Data:     []byte(postDeploymentScript),
				}},
			},

			expectedFiles: map[string]string{
				"/var/tmp/pull-secret.json": "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}\n",

				"/usr/local/bin/install-rhcos-and-restore-seed.sh": `#!/bin/bash

set -e # Halt on error

seed_image=${1:-$SEED_IMAGE}
authfile=${PULL_SECRET_FILE:-"/var/tmp/pull-secret.json"}
ibi_config=${IBI_CONFIGURATION_FILE:-"/var/tmp/ibi-configuration.json"}

# Copy the lca-cli binary to the host, pulling the seed image can sometimes fail
until podman create --authfile "${authfile}" --name lca-cli "${seed_image}" lca-cli ; do
    sleep 10
done
podman cp lca-cli:lca-cli /usr/local/bin/lca-cli
podman rm lca-cli

/usr/local/bin/lca-cli ibi -f "${ibi_config}"
`,

				"/var/tmp/ibi-configuration.json": "{\"extraPartitionLabel\":\"varlibcontainers\",\"extraPartitionNumber\":5,\"extraPartitionStart\":\"-40G\",\"installationDisk\":\"/dev/vda\",\"seedImage\":\"quay.io/openshift-kni/seed-image:4.16.0\",\"seedVersion\":\"4.16.0\"}\n",

				"/etc/containers/registries.conf": "credential-helpers = []\nshort-name-mode = \"\"\nunqualified-search-registries = []\n\n[[registry]]\n  location = \"quay.io\"\n  mirror-by-digest-only = true\n  prefix = \"\"\n\n  [[registry.mirror]]\n    location = \"mirror-quay.io\"\n",

				"/etc/pki/ca-trust/source/anchors/additional-trust-bundle.pem": testCert,

				"/var/tmp/network-config.yaml": testNetworkConfig + "\n",

				"/tmp/example": "helpimtrappedinaswaggerspec",

				"/var/tmp/post.sh": "#!/bin/bash\necho \"Post image-based installation script\"",
			},
			expectedSystemdUnits: map[string]bool{
				"install-rhcos-and-restore-seed.service": true,

				"network-config.service": true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			ignitionAsset := &Ignition{}
			err := ignitionAsset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)

				assertFiles(t, ignitionAsset.Config, tc.expectedFiles)

				assertSystemdUnits(t, ignitionAsset.Config, tc.expectedSystemdUnits)
			}
		})
	}
}

func assertFiles(t *testing.T, config *igntypes.Config, files map[string]string) {
	t.Helper()

	if len(files) > 0 {
		assert.Equal(t, len(files), len(config.Storage.Files))

		for name, content := range files {
			found := false
			for _, i := range config.Storage.Files {
				if i.Node.Path == name {
					actualData, err := dataurl.DecodeString(*i.FileEmbedded1.Contents.Source)
					assert.NoError(t, err)
					assert.Equal(t, content, string(actualData.Data))

					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("expected file %s not found", name))
		}
	}
}

func assertSystemdUnits(t *testing.T, config *igntypes.Config, units map[string]bool) {
	t.Helper()

	for name, enabled := range units {
		for _, unit := range config.Systemd.Units {
			if unit.Name != name {
				continue
			}

			if unit.Enabled == nil {
				assert.Equal(t, enabled, false)
			} else {
				assert.Equal(t, enabled, *unit.Enabled)
			}
		}
	}
}
