package image

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

func TestBaseIso_Generate(t *testing.T) {
	ocReleaseImage := "417.94.202407011030-0"

	cases := []struct {
		name                       string
		envVarOsImageOverrideValue string
		expectedBaseIsoFilename    string
		architecture               string
	}{
		{
			name:                       "os image override",
			envVarOsImageOverrideValue: "openshift-4.18",
			expectedBaseIsoFilename:    "openshift-4.18",
		},
		{
			name:                    "default",
			expectedBaseIsoFilename: ocReleaseImage,
		},
		{
			name:                    "arm64 architecture",
			expectedBaseIsoFilename: fmt.Sprintf("%s-arm64", ocReleaseImage),
			architecture:            types.ArchitectureARM64,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup an HTTP server, to serve the dummy base ISO payload.
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Respond with a fixed size, randomly filled buffer.
				buffer := make([]byte, 1024)
				_, err := rand.Read(buffer)
				assert.NoError(t, err)
				_, err = w.Write(buffer)
				assert.NoError(t, err)
			}))
			defer svr.Close()
			// Creates a temp folder to store the .cache downloaded ISO images.
			tmpPath, err := os.MkdirTemp("", "imagebased-baseiso-test")
			assert.NoError(t, err)
			previousXdgCacheHomeValue := os.Getenv("XDG_CACHE_HOME")
			t.Setenv("XDG_CACHE_HOME", tmpPath)
			// Set the OS image override if defined.
			previousOpenshiftInstallOsImageOverrideValue := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE")
			if tc.envVarOsImageOverrideValue != "" {
				newOsImageOverride := fmt.Sprintf("%s/%s", svr.URL, tc.envVarOsImageOverrideValue)
				t.Setenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE", newOsImageOverride)
			}
			// Cleanup on exit.
			defer func() {
				t.Setenv("XDG_CACHE_HOME", previousXdgCacheHomeValue)
				t.Setenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE", previousOpenshiftInstallOsImageOverrideValue)
				err = os.RemoveAll(tmpPath)
				assert.NoError(t, err)
			}()

			baseIso := &BaseIso{
				streamGetter: func(ctx context.Context) (*stream.Stream, error) {
					return &stream.Stream{
						Architectures: map[string]stream.Arch{
							"x86_64": {
								Artifacts: map[string]stream.PlatformArtifacts{
									"metal": {
										Release: ocReleaseImage,
										Formats: map[string]stream.ImageFormat{
											"iso": {
												Disk: &stream.Artifact{
													Location: fmt.Sprintf("%s/%s", svr.URL, ocReleaseImage),
												},
											},
										},
									},
								},
							},
							"aarch64": {
								Artifacts: map[string]stream.PlatformArtifacts{
									"metal": {
										Release: ocReleaseImage,
										Formats: map[string]stream.ImageFormat{
											"iso": {
												Disk: &stream.Artifact{
													Location: fmt.Sprintf("%s/%s-arm64", svr.URL, ocReleaseImage),
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
			parents := asset.Parents{}
			parents.Add(&ImageBasedInstallationConfig{
				Config: ibiConfig().
					architecture(tc.architecture).
					build(),
			})
			err = baseIso.Generate(context.TODO(), parents)

			assert.NoError(t, err)
			assert.Regexp(t, tc.expectedBaseIsoFilename, baseIso.File.Filename)
		})
	}
}
