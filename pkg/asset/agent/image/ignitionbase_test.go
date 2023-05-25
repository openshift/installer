package image

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
	v1 "k8s.io/api/core/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
)

func TestIgnitionBase_Generate(t *testing.T) {
	setupIgnitionGenerateTest(t)

	cases := []struct {
		name          string
		expectedError string
		expectedFiles []string
	}{
		{
			name:          "default",
			expectedFiles: generatedFilesIgnitionBase(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			deps := buildIgnitionBaseAssetDefaultDependencies()

			parents := asset.Parents{}
			parents.Add(deps...)

			ignitionBaseAsset := &IgnitionBase{}
			err := ignitionBaseAsset.Generate(parents)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)

				assertExpectedFiles(t, &ignitionBaseAsset.Config, tc.expectedFiles, nil)
			}
		})
	}
}

func TestIgnitionGenerateDoesNotChangeIgnitionBaseAsset(t *testing.T) {
	setupIgnitionGenerateTest(t)

	cases := []struct {
		name                      string
		expectedIgnitionBaseFiles []string
		expectedIgnitionFiles     []string
		expectedFileContent       map[string]string
	}{
		{
			name:                      "Ignition.Generate should not change IgnitionBase content",
			expectedIgnitionBaseFiles: generatedFilesIgnitionBase(),
			expectedIgnitionFiles:     generatedFiles(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			depsBase := buildIgnitionBaseAssetDefaultDependencies()

			parentsBase := asset.Parents{}
			parentsBase.Add(depsBase...)

			ignitionBaseAsset := &IgnitionBase{}
			errBase := ignitionBaseAsset.Generate(parentsBase)

			assert.NoError(t, errBase)
			assertExpectedFiles(t, &ignitionBaseAsset.Config, tc.expectedIgnitionBaseFiles, nil)

			deps := buildIgnitionAssetDefaultDependencies()
			parents := asset.Parents{}
			parents.Add(append(deps, ignitionBaseAsset)...)

			ignitionAsset := &Ignition{}
			err := ignitionAsset.Generate(parents)

			assert.NoError(t, err)
			assertExpectedFiles(t, ignitionAsset.Config, tc.expectedIgnitionFiles, nil)

			// The contents of IgnitionBase should not be changed by Ignition.Generate
			assertExpectedFiles(t, &ignitionBaseAsset.Config, tc.expectedIgnitionBaseFiles, nil)
			assert.NotEqual(t, len(ignitionBaseAsset.Config.Storage.Files), len(ignitionAsset.Config.Storage.Files))
			assert.NotEqual(t,
				&ignitionBaseAsset.Config.Passwd.Users[0].PasswordHash,
				ignitionAsset.Config.Passwd.Users[0].PasswordHash)
		})
	}
}

func setupIgnitionGenerateTest(t *testing.T) {
	// Generate calls addStaticNetworkConfig which calls nmstatectl
	_, execErr := exec.LookPath("nmstatectl")
	if execErr != nil {
		t.Skip("No nmstatectl binary available")
	}

	// This patch currently allows testing the Ignition asset using the embedded resources.
	// TODO: Replace it by mocking the filesystem in bootstrap.AddStorageFiles()
	workingDirectory, _ := os.Getwd()
	os.Chdir(path.Join(workingDirectory, "../../../../data"))
}

func assertExpectedFiles(t *testing.T, config *igntypes.Config, expectedFiles []string, expectedFileContent map[string]string) {
	if len(expectedFiles) > 0 {
		assert.Equal(t, len(expectedFiles), len(config.Storage.Files))

		for _, f := range expectedFiles {
			found := false
			for _, i := range config.Storage.Files {
				if i.Node.Path == f {
					if expectedData, ok := expectedFileContent[i.Node.Path]; ok {
						actualData, err := dataurl.DecodeString(*i.FileEmbedded1.Contents.Source)
						assert.NoError(t, err)
						assert.Regexp(t, expectedData, string(actualData.Data))
					}

					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("Expected file %s not found", f))
		}
	}
}

// This test util create the minimum valid set of dependencies for the
// IgnitionBase asset
func buildIgnitionBaseAssetDefaultDependencies() []asset.Asset {
	secretDataBytes, _ := base64.StdEncoding.DecodeString("super-secret")

	return []asset.Asset{
		&manifests.InfraEnv{
			Config: &aiv1beta1.InfraEnv{
				Spec: aiv1beta1.InfraEnvSpec{
					SSHAuthorizedKey: "my-ssh-key",
				},
			},
			File: &asset.File{
				Filename: "infraenv.yaml",
				Data:     []byte("infraenv"),
			},
		},
		&manifests.AgentPullSecret{
			Config: &v1.Secret{
				Data: map[string][]byte{
					".dockerconfigjson": secretDataBytes,
				},
			},
			File: &asset.File{
				Filename: "pull-secret.yaml",
				Data:     []byte("pull-secret"),
			},
		},
		&manifests.ClusterImageSet{
			Config: &hivev1.ClusterImageSet{
				Spec: hivev1.ClusterImageSetSpec{
					ReleaseImage: "registry.ci.openshift.org/origin/release:4.11",
				},
			},
			File: &asset.File{
				Filename: "cluster-image-set.yaml",
				Data:     []byte("cluster-image-set"),
			},
		},
	}
}

func generatedFilesIgnitionBase(otherFiles ...string) []string {
	defaultFiles := []string{
		"/etc/issue",
		"/etc/multipath.conf",
		"/etc/containers/containers.conf",
		"/root/.docker/config.json",
		"/root/assisted.te",
		"/usr/local/bin/common.sh",
		"/usr/local/bin/agent-gather",
		"/usr/local/bin/extract-agent.sh",
		"/usr/local/bin/get-container-images.sh",
		"/usr/local/bin/set-hostname.sh",
		"/usr/local/bin/start-agent.sh",
		"/usr/local/bin/start-cluster-installation.sh",
		"/usr/local/bin/wait-for-assisted-service.sh",
		"/usr/local/bin/set-node-zero.sh",
		"/usr/local/share/assisted-service/assisted-db.env",
		"/usr/local/share/assisted-service/assisted-service.env",
		"/usr/local/share/assisted-service/images.env",
		"/usr/local/bin/bootstrap-service-record.sh",
		"/usr/local/bin/release-image.sh",
		"/usr/local/bin/release-image-download.sh",
		"/etc/assisted/manifests/pull-secret.yaml",
		"/etc/assisted/manifests/cluster-image-set.yaml",
		"/etc/assisted/manifests/infraenv.yaml",
		"/etc/assisted/agent-installer.env",
		"/etc/motd.d/10-agent-installer",
		"/etc/systemd/system.conf.d/10-default-env.conf",
		"/usr/local/bin/install-status.sh",
		"/usr/local/bin/issue_status.sh",
	}
	return append(defaultFiles, otherFiles...)
}
