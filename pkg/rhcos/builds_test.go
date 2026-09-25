package rhcos

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatURLWithCompressedIntegrity(t *testing.T) {
	cases := []struct {
		name          string
		artifact      stream.Artifact
		expected      string
		expectedError string
	}{
		{
			name: "carries the compressed digest",
			artifact: stream.Artifact{
				Location:           "https://example.com/rhcos.tar.gz",
				Sha256:             "compressed",
				UncompressedSha256: "uncompressed",
			},
			expected: "https://example.com/rhcos.tar.gz?" + CompressedSHA256Param + "=compressed",
		},
		{
			// The GCP artifacts carry no uncompressed digest, which is how an
			// unverifiable URL slipped through before: the caller reads an empty
			// parameter and has no way to tell it apart from a missing one.
			name: "rejects an artifact with no compressed digest",
			artifact: stream.Artifact{
				Location:           "https://example.com/rhcos.tar.gz",
				UncompressedSha256: "uncompressed",
			},
			expectedError: "artifact https://example.com/rhcos.tar.gz has no sha256 checksum",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FormatURLWithCompressedIntegrity(&tc.artifact)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestFindCompressedArtifactURL(t *testing.T) {
	disk := &stream.Artifact{Location: "https://example.com/rhcos.tar.gz", Sha256: "compressed"}

	cases := []struct {
		name          string
		artifacts     stream.PlatformArtifacts
		expected      string
		expectedError string
	}{
		{
			name: "single disk artifact",
			artifacts: stream.PlatformArtifacts{
				Formats: map[string]stream.ImageFormat{"tar.gz": {Disk: disk}},
			},
			expected: "https://example.com/rhcos.tar.gz?" + CompressedSHA256Param + "=compressed",
		},
		{
			name: "multiple disk artifacts",
			artifacts: stream.PlatformArtifacts{
				Formats: map[string]stream.ImageFormat{"tar.gz": {Disk: disk}, "tar.xz": {Disk: disk}},
			},
			expectedError: `multiple "disk" artifacts found`,
		},
		{
			name: "no disk artifact",
			artifacts: stream.PlatformArtifacts{
				Formats: map[string]stream.ImageFormat{"iso": {}},
			},
			expectedError: `no "disk" artifact found`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FindCompressedArtifactURL(tc.artifacts)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

// The GCP sovereign cloud flow uploads the artifact in the compressed form the
// mirror serves it in, so the compressed digest is the only one it can verify.
// Every embedded stream must carry one, or the boot image is built from an
// unverified download.
//
// The stream files are read straight from disk rather than through
// FetchCoreOSBuild, which resolves them relative to the working directory of the
// installer binary.
func TestEmbeddedGCPArtifactsCarryCompressedChecksum(t *testing.T) {
	streamFiles, err := filepath.Glob("../../data/data/coreos/*.json")
	require.NoError(t, err)
	require.NotEmpty(t, streamFiles)

	for _, streamFile := range streamFiles {
		t.Run(filepath.Base(streamFile), func(t *testing.T) {
			body, err := os.ReadFile(streamFile) //nolint:gosec // path comes from the glob above
			require.NoError(t, err)

			var st stream.Stream
			require.NoError(t, json.Unmarshal(body, &st))

			for archName, arch := range st.Architectures {
				artifacts, ok := arch.Artifacts["gcp"]
				if !ok {
					continue
				}
				rawURL, err := FindCompressedArtifactURL(artifacts)
				require.NoError(t, err, "architecture %s", archName)

				u, err := url.Parse(rawURL)
				require.NoError(t, err, "architecture %s", archName)
				assert.NotEmpty(t, u.Query().Get(CompressedSHA256Param), "architecture %s", archName)
			}
		})
	}
}
