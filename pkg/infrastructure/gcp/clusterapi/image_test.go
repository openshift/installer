package clusterapi

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/infrastructure/gcp/clusterapi/mock"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	testInfraID   = "test-infra-id"
	testProjectID = "eu0:sovereign-project"
	testRegion    = "u-northeast1"
	testImageName = testInfraID + "-rhcos"
	testBucket    = testInfraID + "-rhcos-image"
	testImageRef  = "projects/" + testProjectID + "/global/images/" + testImageName
)

// testPreProvisionInput builds a sovereign-cloud input whose control plane and
// compute pools resolve to the same artifact, which is the supported case.
func testPreProvisionInput(imageURL string) clusterapi.PreProvisionInput {
	return clusterapi.PreProvisionInput{
		InfraID: testInfraID,
		InstallConfig: installconfig.MakeAsset(&types.InstallConfig{
			Platform: types.Platform{
				GCP: &gcptypes.Platform{
					ProjectID: testProjectID,
					Region:    testRegion,
				},
			},
		}),
		RhcosImage: &rhcos.Image{ControlPlane: imageURL, Compute: imageURL},
	}
}

// stubDownloader returns a downloader that writes a placeholder artifact into a
// throwaway directory, standing in for the multi-gigabyte real download. It
// records the arguments it was called with so that URL parsing can be asserted.
func stubDownloader(t *testing.T, gotURL, gotChecksum *string) imageDownloader {
	t.Helper()
	return func(imageURL, sha256Checksum string) (string, error) {
		*gotURL = imageURL
		*gotChecksum = sha256Checksum
		path := filepath.Join(t.TempDir(), "rhcos-gcp.tar.gz")
		require.NoError(t, os.WriteFile(path, []byte("rhcos"), 0o600))
		return path, nil
	}
}

// failingDownloader fails the test if the flow reaches the download step.
func failingDownloader(t *testing.T) imageDownloader {
	t.Helper()
	return func(string, string) (string, error) {
		t.Error("download must not be attempted")
		return "", errors.New("unexpected download")
	}
}

func TestPublishRHCOSImageRejectsHeterogeneousArchitectures(t *testing.T) {
	in := testPreProvisionInput("https://example.com/rhcos-aarch64.tar.gz")
	in.RhcosImage.Compute = "https://example.com/rhcos-x86_64.tar.gz"

	// No EXPECT calls: the guard must reject before touching GCP at all.
	client := mock.NewMockImageClient(gomock.NewController(t))

	_, err := publishRHCOSImage(context.Background(), in, client, failingDownloader(t))
	assert.ErrorContains(t, err, "heterogeneous architectures are not supported on sovereign clouds")
	assert.ErrorContains(t, err, "rhcos-aarch64.tar.gz")
	assert.ErrorContains(t, err, "rhcos-x86_64.tar.gz")
}

func TestPublishRHCOSImageReusesExistingImage(t *testing.T) {
	client := mock.NewMockImageClient(gomock.NewController(t))
	client.EXPECT().GetImage(gomock.Any(), testProjectID, testImageName).
		Return(&compute.Image{Name: testImageName, Status: imageStatusReady}, nil)

	// A retried PreProvision must short circuit: no download, no staging, no
	// second Images.Insert.
	ref, err := publishRHCOSImage(context.Background(),
		testPreProvisionInput("https://example.com/rhcos.tar.gz"), client, failingDownloader(t))

	assert.NoError(t, err)
	assert.Equal(t, testImageRef, ref)
}

func TestPublishRHCOSImageRejectsUnusableExistingImage(t *testing.T) {
	client := mock.NewMockImageClient(gomock.NewController(t))
	client.EXPECT().GetImage(gomock.Any(), testProjectID, testImageName).
		Return(&compute.Image{Name: testImageName, Status: "FAILED"}, nil)

	_, err := publishRHCOSImage(context.Background(),
		testPreProvisionInput("https://example.com/rhcos.tar.gz"), client, failingDownloader(t))

	assert.ErrorContains(t, err, "is not usable (status FAILED)")
	assert.ErrorContains(t, err, "delete it and retry the install")
}

func TestPublishRHCOSImageSurfacesLookupFailure(t *testing.T) {
	// A transient lookup failure must not be mistaken for a missing image and
	// send the install into a multi-gigabyte download.
	client := mock.NewMockImageClient(gomock.NewController(t))
	client.EXPECT().GetImage(gomock.Any(), testProjectID, testImageName).
		Return(nil, errors.New("failed to check for existing image: quota exceeded"))

	_, err := publishRHCOSImage(context.Background(),
		testPreProvisionInput("https://example.com/rhcos.tar.gz"), client, failingDownloader(t))

	assert.ErrorContains(t, err, "quota exceeded")
}

func TestPublishRHCOSImage(t *testing.T) {
	client := mock.NewMockImageClient(gomock.NewController(t))

	wantLabels := map[string]string{"kubernetes-io-cluster-" + testInfraID: "owned"}
	var uploadedPath string

	gomock.InOrder(
		client.EXPECT().GetImage(gomock.Any(), testProjectID, testImageName).Return(nil, nil),
		client.EXPECT().CreateStagingBucket(gomock.Any(), testBucket, testProjectID, testRegion, wantLabels).Return(nil),
		client.EXPECT().UploadObject(gomock.Any(), testBucket, rhcosImageObjectName, gomock.Any()).
			DoAndReturn(func(_ context.Context, _, _, path string) error {
				uploadedPath = path
				return nil
			}),
		client.EXPECT().UniverseDomain(gomock.Any()).Return("sovereign.example.com", nil),
		// The GCS source must use the universe domain, otherwise a sovereign
		// cloud resolves the wrong storage host.
		client.EXPECT().CreateImage(gomock.Any(), testProjectID, testImageName,
			"https://storage.sovereign.example.com/"+testBucket+"/"+rhcosImageObjectName, wantLabels).Return(nil),
		client.EXPECT().CleanupStaging(gomock.Any(), testBucket, rhcosImageObjectName).Return(nil),
	)

	var gotURL, gotChecksum string
	in := testPreProvisionInput("https://example.com/rhcos.tar.gz?sha256=abc123")
	ref, err := publishRHCOSImage(context.Background(), in, client, stubDownloader(t, &gotURL, &gotChecksum))

	assert.NoError(t, err)
	assert.Equal(t, testImageRef, ref)
	// The checksum is carried in the URL query and must be split out rather
	// than passed to the mirror as part of the download URL.
	assert.Equal(t, "https://example.com/rhcos.tar.gz", gotURL)
	assert.Equal(t, "abc123", gotChecksum)
	// The temp directory holding the artifact is removed on the way out.
	assert.NoFileExists(t, uploadedPath)
}

func TestPublishRHCOSImageCleansUpStagingOnFailure(t *testing.T) {
	cases := []struct {
		name          string
		uploadErr     error
		createErr     error
		expectedError string
	}{
		{
			name:          "upload failure",
			uploadErr:     errors.New("connection reset"),
			expectedError: "failed to upload RHCOS image to GCS: connection reset",
		},
		{
			name:          "image creation failure",
			createErr:     errors.New("operation test-op failed"),
			expectedError: "failed to create compute image: operation test-op failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := mock.NewMockImageClient(gomock.NewController(t))
			client.EXPECT().GetImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			client.EXPECT().CreateStagingBucket(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			client.EXPECT().UploadObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.uploadErr)
			if tc.uploadErr == nil {
				client.EXPECT().UniverseDomain(gomock.Any()).Return("sovereign.example.com", nil)
				client.EXPECT().CreateImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tc.createErr)
			}
			// The bucket is scratch space; leaking it on the error path leaves
			// gigabytes of billable storage behind until destroy runs.
			client.EXPECT().CleanupStaging(gomock.Any(), testBucket, rhcosImageObjectName).Return(nil)

			var url, checksum string
			_, err := publishRHCOSImage(context.Background(),
				testPreProvisionInput("https://example.com/rhcos.tar.gz"), client, stubDownloader(t, &url, &checksum))

			assert.EqualError(t, err, tc.expectedError)
		})
	}
}

func TestPublishRHCOSImageSkipsCleanupWhenBucketWasNeverCreated(t *testing.T) {
	client := mock.NewMockImageClient(gomock.NewController(t))
	client.EXPECT().GetImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	client.EXPECT().CreateStagingBucket(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("permission denied"))
	// No CleanupStaging expectation: deleting a bucket we failed to create
	// would only produce a misleading warning.

	var url, checksum string
	_, err := publishRHCOSImage(context.Background(),
		testPreProvisionInput("https://example.com/rhcos.tar.gz"), client, stubDownloader(t, &url, &checksum))

	assert.EqualError(t, err, "failed to create staging bucket: permission denied")
}

func TestPublishRHCOSImageReportsCleanupFailureWithoutFailingInstall(t *testing.T) {
	client := mock.NewMockImageClient(gomock.NewController(t))
	client.EXPECT().GetImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	client.EXPECT().CreateStagingBucket(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	client.EXPECT().UploadObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	client.EXPECT().UniverseDomain(gomock.Any()).Return("sovereign.example.com", nil)
	client.EXPECT().CreateImage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	client.EXPECT().CleanupStaging(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("bucket not empty"))

	// The image exists at this point, so a failure to tear down scratch space
	// must not abort the install: destroy cleans the bucket up later.
	var url, checksum string
	ref, err := publishRHCOSImage(context.Background(),
		testPreProvisionInput("https://example.com/rhcos.tar.gz"), client, stubDownloader(t, &url, &checksum))

	assert.NoError(t, err)
	assert.Equal(t, testImageRef, ref)
}

func TestDownloadRHCOSImage(t *testing.T) {
	body := []byte("fake rhcos tarball")
	checksum := fmt.Sprintf("%x", sha256.Sum256(body))

	cases := []struct {
		name string
		// handler serves the artifact; nil serves body with a 200.
		handler       http.HandlerFunc
		checksum      string
		expectedError string
		expectedCalls int
	}{
		{
			name:          "downloads and verifies the checksum",
			checksum:      checksum,
			expectedCalls: 1,
		},
		{
			// A cluster built from a corrupted artifact fails in ways that are
			// very hard to diagnose, so a mismatch must stop the install.
			name:          "rejects a checksum mismatch",
			checksum:      "0000000000000000000000000000000000000000000000000000000000000000",
			expectedError: "checksum mismatch for RHCOS image",
			expectedCalls: 3,
		},
		{
			name:          "skips verification when no checksum is supplied",
			expectedCalls: 1,
		},
		{
			name: "retries a failing mirror",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
			},
			expectedError: "bad status downloading RHCOS image: 503",
			expectedCalls: 3,
		},
	}

	// Keep the retry cases from sleeping for the production interval.
	original := downloadRetryDelay
	downloadRetryDelay = time.Millisecond
	defer func() { downloadRetryDelay = original }()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			// A local server keeps the test off the network entirely.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if tc.handler != nil {
					tc.handler(w, r)
					return
				}
				if _, err := w.Write(body); err != nil {
					t.Errorf("failed to serve artifact: %v", err)
				}
			}))
			defer srv.Close()

			path, err := downloadRHCOSImage(srv.URL+"/rhcos.tar.gz", tc.checksum)
			assert.Equal(t, tc.expectedCalls, calls)

			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
				return
			}

			require.NoError(t, err)
			defer os.RemoveAll(filepath.Dir(path))
			got, err := os.ReadFile(path) //nolint:gosec // path is produced by the function under test
			require.NoError(t, err)
			assert.Equal(t, body, got)
		})
	}
}

func TestDownloadRHCOSImageRemovesTempDirOnFailure(t *testing.T) {
	original := downloadRetryDelay
	downloadRetryDelay = time.Millisecond
	defer func() { downloadRetryDelay = original }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	before, err := filepath.Glob(filepath.Join(os.TempDir(), "rhcos-gcp-*"))
	require.NoError(t, err)

	_, err = downloadRHCOSImage(srv.URL+"/rhcos.tar.gz", "")
	assert.Error(t, err)

	// A failed install attempt must not leave partial multi-gigabyte downloads
	// in the temp directory.
	after, err := filepath.Glob(filepath.Join(os.TempDir(), "rhcos-gcp-*"))
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestBuildImageLabels(t *testing.T) {
	labels := buildImageLabels(&gcptypes.Platform{
		UserLabels: []gcptypes.UserLabel{{Key: "team", Value: "cors"}},
	}, testInfraID)

	assert.Equal(t, map[string]string{
		"kubernetes-io-cluster-" + testInfraID: "owned",
		"team":                                 "cors",
	}, labels)
}

func TestIsNotFoundError(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "not found",
			err:      &googleapi.Error{Code: http.StatusNotFound},
			expected: true,
		},
		{
			// A wrapped error must still be recognised, since the lookup path
			// treats anything other than a 404 as a real failure.
			name:     "wrapped not found",
			err:      fmt.Errorf("looking up image: %w", &googleapi.Error{Code: http.StatusNotFound}),
			expected: true,
		},
		{
			name:     "forbidden is not a missing image",
			err:      &googleapi.Error{Code: http.StatusForbidden},
			expected: false,
		},
		{
			name:     "non-api error",
			err:      errors.New("connection reset"),
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, isNotFoundError(tc.err))
		})
	}
}
