package image

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/common"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
)

func TestUnconfiguredIgnition_Generate(t *testing.T) {
	skipTestIfnmstatectlIsMissing(t)

	nmStateConfig := getTestNMStateConfig()

	cases := []struct {
		name              string
		overrideDeps      []asset.Asset
		expectedError     string
		expectedFiles     []string
		serviceEnabledMap map[string]bool
	}{
		{
			name:          "default-configs-and-no-nmstateconfigs",
			expectedFiles: generatedFilesUnconfiguredIgnition("/usr/local/bin/pre-network-manager-config.sh", "/usr/local/bin/oci-eval-user-data.sh"),
			serviceEnabledMap: map[string]bool{
				"pre-network-manager-config.service": false,
				"oci-eval-user-data.service":         true,
				"agent-check-config-image.service":   true},
		},
		{
			name: "with-mirror-configs",
			overrideDeps: []asset.Asset{
				&mirror.RegistriesConf{
					File: &asset.File{
						Filename: mirror.RegistriesConfFilename,
						Data:     []byte(""),
					},
					MirrorConfig: []mirror.RegistriesConfig{
						{
							Location: "some.registry.org/release",
							Mirror:   "some.mirror.org",
						},
					},
				},
				&mirror.CaBundle{
					File: &asset.File{
						Filename: "my.crt",
						Data:     []byte("my-certificate"),
					},
				},
			},
			expectedFiles: generatedFilesUnconfiguredIgnition(registriesConfPath,
				registryCABundlePath, "/usr/local/bin/pre-network-manager-config.sh", "/usr/local/bin/oci-eval-user-data.sh"),
			serviceEnabledMap: map[string]bool{
				"pre-network-manager-config.service": false,
				"oci-eval-user-data.service":         true,
				"agent-check-config-image.service":   true},
		},
		{
			name: "with-nmstateconfigs",
			overrideDeps: []asset.Asset{
				&nmStateConfig,
			},
			expectedFiles: generatedFilesUnconfiguredIgnition("/etc/assisted/network/host0/eth0.nmconnection",
				"/etc/assisted/network/host0/mac_interface.ini", "/usr/local/bin/pre-network-manager-config.sh", "/usr/local/bin/oci-eval-user-data.sh"),
			serviceEnabledMap: map[string]bool{
				"pre-network-manager-config.service": true,
				"oci-eval-user-data.service":         true,
				"agent-check-config-image.service":   true},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			deps := buildUnconfiguredIgnitionAssetDefaultDependencies(t)

			overrideDeps(deps, tc.overrideDeps)

			parents := asset.Parents{}
			parents.Add(deps...)

			unconfiguredIgnitionAsset := &UnconfiguredIgnition{}
			err := unconfiguredIgnitionAsset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)

				assertExpectedFiles(t, unconfiguredIgnitionAsset.Config, tc.expectedFiles, nil)

				assertServiceEnabled(t, unconfiguredIgnitionAsset.Config, tc.serviceEnabledMap)
			}
		})
	}
}

// This test util create the minimum valid set of dependencies for the
// UnconfiguredIgnition asset.
func buildUnconfiguredIgnitionAssetDefaultDependencies(t *testing.T) []asset.Asset {
	t.Helper()

	infraEnv := getTestInfraEnv()
	agentPullSecret := getTestAgentPullSecret(t)
	clusterImageSet := getTestClusterImageSet()

	return []asset.Asset{
		&infraEnv,
		&agentPullSecret,
		&clusterImageSet,
		&manifests.NMStateConfig{},
		&mirror.RegistriesConf{},
		&mirror.CaBundle{},
		&common.InfraEnvID{},
	}
}

func getTestInfraEnv() manifests.InfraEnv {
	return manifests.InfraEnv{
		Config: &aiv1beta1.InfraEnv{
			Spec: aiv1beta1.InfraEnvSpec{
				SSHAuthorizedKey: "my-ssh-key",
			},
		},
		File: &asset.File{
			Filename: "infraenv.yaml",
			Data:     []byte("infraenv"),
		},
	}
}

func getTestAgentPullSecret(t *testing.T) manifests.AgentPullSecret {
	t.Helper()
	secretDataBytes, err := base64.StdEncoding.DecodeString("c3VwZXItc2VjcmV0Cg==")
	assert.NoError(t, err)
	return manifests.AgentPullSecret{
		Config: &v1.Secret{
			Data: map[string][]byte{
				".dockerconfigjson": secretDataBytes,
			},
		},
		File: &asset.File{
			Filename: "pull-secret.yaml",
			Data:     []byte("pull-secret"),
		},
	}
}

func getTestClusterImageSet() manifests.ClusterImageSet {
	return manifests.ClusterImageSet{
		Config: &hivev1.ClusterImageSet{
			Spec: hivev1.ClusterImageSetSpec{
				ReleaseImage: "registry.ci.openshift.org/origin/release:4.11",
			},
		},
		File: &asset.File{
			Filename: "cluster-image-set.yaml",
			Data:     []byte("cluster-image-set"),
		},
	}
}

func getTestNMStateConfig() manifests.NMStateConfig {
	return manifests.NMStateConfig{
		Config: []*aiv1beta1.NMStateConfig{
			{
				Spec: aiv1beta1.NMStateConfigSpec{
					Interfaces: []*aiv1beta1.Interface{
						{
							Name:       "eth0",
							MacAddress: "00:01:02:03:04:05",
						},
					},
				},
			},
		},
		StaticNetworkConfig: []*models.HostStaticNetworkConfig{
			{
				MacInterfaceMap: models.MacInterfaceMap{
					{LogicalNicName: "eth0", MacAddress: "00:01:02:03:04:05"},
				},
				NetworkYaml: "interfaces:\n- ipv4:\n    address:\n    - ip: 192.168.122.21\n      prefix-length: 24\n    enabled: true\n  mac-address: 00:01:02:03:04:05\n  name: eth0\n  state: up\n  type: ethernet\n",
			},
		},
		File: &asset.File{
			Filename: "nmstateconfig.yaml",
			Data:     []byte("nmstateconfig"),
		},
	}
}

func generatedFilesUnconfiguredIgnition(otherFiles ...string) []string {
	unconfiguredIgnitionFiles := []string{
		"/etc/assisted/manifests/pull-secret.yaml",
		"/etc/assisted/manifests/cluster-image-set.yaml",
		"/etc/assisted/manifests/infraenv.yaml",
	}
	unconfiguredIgnitionFiles = append(unconfiguredIgnitionFiles, otherFiles...)
	return append(unconfiguredIgnitionFiles, commonFiles()...)
}
