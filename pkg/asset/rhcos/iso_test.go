package rhcos

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/stretchr/testify/assert"
)

func TestBaseIso(t *testing.T) {
	ocReleaseImage := "416.94.202402130130-0"
	ocBaseIsoFilename := "openshift-4.16"

	cases := []struct {
		name                       string
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
			name:                    "default",
			expectedBaseIsoFilename: ocBaseIsoFilename,
		},
		{
			name:                    "direct download if oc is not available",
			getIsoError:             &exec.Error{},
			expectedBaseIsoFilename: ocReleaseImage,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup a fake http server, to serve the future download request.
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Answer with a fixed size randomly filled buffer
				buffer := make([]byte, 1024)
				_, err := rand.Read(buffer)
				assert.NoError(t, err)
				_, err = w.Write(buffer)
				assert.NoError(t, err)
			}))
			defer svr.Close()
			// Creates a tmp folder to store the .cache downloaded images.
			tmpPath, err := os.MkdirTemp("", "agent-baseiso-test")
			assert.NoError(t, err)
			previousXdgCacheHomeValue := os.Getenv("XDG_CACHE_HOME")
			t.Setenv("XDG_CACHE_HOME", tmpPath)
			// Set the image override if defined
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

			fetcher := NewBaseISOFetcher(
				&mockRelease{
					isoBaseVersion:  ocReleaseImage,
					baseIsoFileName: ocBaseIsoFilename,
					baseIsoError:    tc.getIsoError,
				},
				func(ctx context.Context) (*stream.Stream, error) {
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
						},
					}, nil
				})
			filename, err := fetcher.GetBaseISOFilename(context.Background(), "")

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Regexp(t, tc.expectedBaseIsoFilename, filename)
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

func (m *mockRelease) GetBaseIso(architecture string, streamGetter CoreOSBuildFetcher) (string, error) {
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
