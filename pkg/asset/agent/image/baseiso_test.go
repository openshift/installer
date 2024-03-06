package image

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/openshift/assisted-service/api/v1beta1"
	v1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

func setupEmbeddedResources(t *testing.T) {
	workingDirectory, err := os.Getwd()
	assert.NoError(t, err)
	err = os.Chdir(path.Join(workingDirectory, "../../../../data"))
	assert.NoError(t, err)
}

func TestBaseIso_Generate(t *testing.T) {
	setupEmbeddedResources(t)
	ocReleaseImage := "416.94.202402130130-0"
	ocBaseIsoFilename := "openshift-4.16"

	cases := []struct {
		name                       string
		dependencies               []asset.Asset
		envVarOsImageOverrideValue string
		getIsoError                error
		expectedBaseIsoFilename    string
		expectedError              string
	}{
		{
			name:                       "os image override",
			envVarOsImageOverrideValue: "openshift-4.15",
			expectedBaseIsoFilename:    "openshift-4.15",
		},
		{
			name: "default",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				&manifests.AgentManifests{
					InfraEnv: &v1beta1.InfraEnv{},
					ClusterImageSet: &v1.ClusterImageSet{
						Spec: v1.ClusterImageSetSpec{
							ReleaseImage: ocReleaseImage,
						},
					},
					PullSecret: &corev1.Secret{
						StringData: map[string]string{
							".dockerconfigjson": "supersecret",
						},
					},
				},
				&mirror.RegistriesConf{},
			},
			expectedBaseIsoFilename: ocBaseIsoFilename,
		},
		{
			name: "direct download if oc is not available",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				&manifests.AgentManifests{
					InfraEnv: &v1beta1.InfraEnv{},
					ClusterImageSet: &v1.ClusterImageSet{
						Spec: v1.ClusterImageSetSpec{
							ReleaseImage: ocReleaseImage,
						},
					},
					PullSecret: &corev1.Secret{
						StringData: map[string]string{
							".dockerconfigjson": "supersecret",
						},
					},
				},
				&mirror.RegistriesConf{},
			},
			getIsoError:             &exec.Error{"", exec.ErrNotFound},
			expectedBaseIsoFilename: ocReleaseImage,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dependencies := asset.Parents{}
			dependencies.Add(tc.dependencies...)

			// Setup a fake http server, to serve the future download request.
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Answer with a fixed size randomly filled buffer
				buffer := make([]byte, 1024)
				rand.Read(buffer)
				w.Write(buffer)
			}))
			defer svr.Close()
			// Creates a tmp folder to store the .cache downloaded images.
			tmpPath, err := os.MkdirTemp("", "agent-baseiso-test")
			assert.NoError(t, err)
			previousXdgCacheHomeValue := os.Getenv("XDG_CACHE_HOME")
			os.Setenv("XDG_CACHE_HOME", tmpPath)
			// Set the image override if defined
			previousOpenshiftInstallOsImageOverrideValue := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE")
			if tc.envVarOsImageOverrideValue != "" {
				newOsImageOverride := fmt.Sprintf("%s/%s", svr.URL, tc.envVarOsImageOverrideValue)
				os.Setenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE", newOsImageOverride)
			}
			// Cleanup on exit.
			defer func() {
				err := os.Setenv("XDG_CACHE_HOME", previousXdgCacheHomeValue)
				assert.NoError(t, err)
				err = os.Setenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE", previousOpenshiftInstallOsImageOverrideValue)
				assert.NoError(t, err)
				err = os.RemoveAll(tmpPath)
				assert.NoError(t, err)
			}()

			baseIso := &BaseIso{
				ocRelease: &mockRelease{
					isoBaseVersion:  ocReleaseImage,
					baseIsoFileName: ocBaseIsoFilename,
					baseIsoError:    tc.getIsoError,
				},
				streamGetter: func(ctx context.Context) (*stream.Stream, error) {
					return &stream.Stream{
						Architectures: map[string]stream.Arch{
							"x86_64": {
								Artifacts: map[string]stream.PlatformArtifacts{
									"metal": {
										Release: ocReleaseImage,
										Formats: map[string]stream.ImageFormat{
											"iso": stream.ImageFormat{
												Disk: &stream.Artifact{
													Location: fmt.Sprintf("%s/%s", svr.URL, ocReleaseImage),
												},
											},
										},
									},
								},
							},
						},
					}, nil
				},
			}
			err = baseIso.Generate(dependencies)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Regexp(t, tc.expectedBaseIsoFilename, baseIso.File.Filename)
			} else {
				assert.Equal(t, tc.expectedError, err.Error())
			}
		})
	}
}

type mockRelease struct {
	isoBaseVersion  string
	baseIsoFileName string
	baseIsoError    error
}

func (m *mockRelease) GetBaseIso(architecture string) (string, error) {
	if m.baseIsoError != nil {
		return "", m.baseIsoError
	}
	return m.baseIsoFileName, nil
}

func (m *mockRelease) GetBaseIsoVersion(architecture string) (string, error) {
	return m.isoBaseVersion, nil
}

func (m *mockRelease) ExtractFile(image string, filename string, architecture string) ([]string, error) {
	return []string{}, nil
}

func TestInfraBaseIso_GenerateOld(t *testing.T) {

	parents := asset.Parents{}
	manifests := &manifests.AgentManifests{}
	installConfig := &agent.OptionalInstallConfig{}
	parents.Add(manifests, installConfig)

	asset := &BaseIso{}
	err := asset.Generate(parents)
	assert.NoError(t, err)

	assert.NotEmpty(t, asset.Files())
	baseIso := asset.Files()[0]
	assert.Equal(t, baseIso.Filename, "some-openshift-release.iso")
}
