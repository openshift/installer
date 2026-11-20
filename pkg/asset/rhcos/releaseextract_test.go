package rhcos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

// TestCoreOsFileNames verifies the machine-os-images file paths for each
// OS image stream. They must match the layout produced by machine-os-images
// (fetch_image.sh): version 9 uses the coreos- prefix and version 10 uses
// the coreos10- prefix, for both RHCOS (rhel-*) and SCOS (centos-*) streams.
func TestCoreOsFileNames(t *testing.T) {
	cases := []struct {
		name           string
		stream         types.OSImageStream
		arch           string
		expectedISO    string
		expectedSha256 string
		expectedStream string
	}{
		{
			name:           "rhel-9 x86_64",
			stream:         types.OSImageStreamRHCOS9,
			arch:           "x86_64",
			expectedISO:    "/coreos/coreos-x86_64.iso",
			expectedSha256: "/coreos/coreos-x86_64.iso.sha256",
			expectedStream: "/coreos/coreos-stream.json",
		},
		{
			name:           "rhel-9 aarch64",
			stream:         types.OSImageStreamRHCOS9,
			arch:           "aarch64",
			expectedISO:    "/coreos/coreos-aarch64.iso",
			expectedSha256: "/coreos/coreos-aarch64.iso.sha256",
			expectedStream: "/coreos/coreos-stream.json",
		},
		{
			name:           "rhel-10 x86_64",
			stream:         types.OSImageStreamRHCOS10,
			arch:           "x86_64",
			expectedISO:    "/coreos/coreos10-x86_64.iso",
			expectedSha256: "/coreos/coreos10-x86_64.iso.sha256",
			expectedStream: "/coreos/coreos10-stream.json",
		},
		{
			name:           "rhel-10 aarch64",
			stream:         types.OSImageStreamRHCOS10,
			arch:           "aarch64",
			expectedISO:    "/coreos/coreos10-aarch64.iso",
			expectedSha256: "/coreos/coreos10-aarch64.iso.sha256",
			expectedStream: "/coreos/coreos10-stream.json",
		},
		{
			name:           "centos-10 x86_64",
			stream:         types.OSImageStreamCentos10,
			arch:           "x86_64",
			expectedISO:    "/coreos/coreos10-x86_64.iso",
			expectedSha256: "/coreos/coreos10-x86_64.iso.sha256",
			expectedStream: "/coreos/coreos10-stream.json",
		},
		{
			name:           "centos-10 aarch64",
			stream:         types.OSImageStreamCentos10,
			arch:           "aarch64",
			expectedISO:    "/coreos/coreos10-aarch64.iso",
			expectedSha256: "/coreos/coreos10-aarch64.iso.sha256",
			expectedStream: "/coreos/coreos10-stream.json",
		},
		{
			name:           "empty stream defaults to coreos- layout",
			stream:         "",
			arch:           "x86_64",
			expectedISO:    "/coreos/coreos-x86_64.iso",
			expectedSha256: "/coreos/coreos-x86_64.iso.sha256",
			expectedStream: "/coreos/coreos-stream.json",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedISO, getCoreOsFileName(tc.stream, tc.arch))
			assert.Equal(t, tc.expectedSha256, getCoreOsSha256FileName(tc.stream, tc.arch))
			assert.Equal(t, tc.expectedStream, getCoreOsStreamFileName(tc.stream))
		})
	}
}
