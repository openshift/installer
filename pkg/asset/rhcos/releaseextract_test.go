package rhcos

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

// mirrorConfig returns a MirrorConfig with a single mirror entry.
func mirrorConfig() types.MirrorConfig {
	return types.MirrorConfig{
		{
			Location: "registry.example.com/ocp/release",
			Mirrors:  []string{"mirror.example.com/ocp/release"},
		},
	}
}

// hasArgWithPrefix reports whether any argument in cmd starts with prefix.
func hasArgWithPrefix(cmd []string, prefix string) bool {
	for _, arg := range cmd {
		if strings.HasPrefix(arg, prefix) {
			return true
		}
	}
	return false
}

func TestAppendMirrorArgs(t *testing.T) {
	cases := []struct {
		name           string
		mirrorConfig   types.MirrorConfig
		expectInsecure bool
		expectICSP     bool
	}{
		{
			name:           "no mirror omits --insecure",
			mirrorConfig:   nil,
			expectInsecure: false,
			expectICSP:     false,
		},
		{
			name:           "mirror includes --insecure and icsp file",
			mirrorConfig:   mirrorConfig(),
			expectInsecure: true,
			expectICSP:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &releasePayload{mirrorConfig: tc.mirrorConfig}
			base := []string{"oc", "adm", "release", "info"}

			cmd, cleanup, err := r.appendMirrorArgs(base)
			assert.NoError(t, err)
			assert.NotNil(t, cleanup)
			defer cleanup()

			assert.Equal(t, tc.expectInsecure, assertContains(cmd, "--insecure=true"),
				"unexpected --insecure=true presence in %v", cmd)
			assert.Equal(t, tc.expectICSP, hasArgWithPrefix(cmd, "--icsp-file="),
				"unexpected --icsp-file presence in %v", cmd)
		})
	}
}

// The release-info command must only carry --insecure when mirrors are set.
func TestGetImageFromReleaseCmd(t *testing.T) {
	cases := []struct {
		name           string
		mirrorConfig   types.MirrorConfig
		expectInsecure bool
	}{
		{name: "connected install", mirrorConfig: nil, expectInsecure: false},
		{name: "mirrored install", mirrorConfig: mirrorConfig(), expectInsecure: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &releasePayload{
				releaseImage: "registry.example.com/ocp/release:latest",
				mirrorConfig: tc.mirrorConfig,
			}

			cmd := []string{"oc", "adm", "release", "info", "--image-for=machine-os-images", "--filter-by-os=linux/amd64"}
			cmd, cleanup, err := r.appendMirrorArgs(cmd)
			assert.NoError(t, err)
			defer cleanup()

			assert.Equal(t, tc.expectInsecure, assertContains(cmd, "--insecure=true"))
		})
	}
}

// The image-extract command must only carry --insecure when mirrors are set.
func TestExtractFileFromImageCmd(t *testing.T) {
	cases := []struct {
		name           string
		mirrorConfig   types.MirrorConfig
		expectInsecure bool
	}{
		{name: "connected install", mirrorConfig: nil, expectInsecure: false},
		{name: "mirrored install", mirrorConfig: mirrorConfig(), expectInsecure: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &releasePayload{mirrorConfig: tc.mirrorConfig}

			cmd := []string{"oc", "image", "extract", "--path=/coreos/coreos-x86_64.iso:/tmp", "--filter-by-os=linux/amd64", "--confirm"}
			cmd, cleanup, err := r.appendMirrorArgs(cmd)
			assert.NoError(t, err)
			defer cleanup()

			assert.Equal(t, tc.expectInsecure, assertContains(cmd, "--insecure=true"))
		})
	}
}

// assertContains reports whether cmd contains an exact match for arg.
func assertContains(cmd []string, arg string) bool {
	for _, a := range cmd {
		if a == arg {
			return true
		}
	}
	return false
}
