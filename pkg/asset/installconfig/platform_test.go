package installconfig

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
)

func TestPlatformDependencies(t *testing.T) {
	platform := &Platform{}
	act := platform.Dependencies()
	assert.Empty(t, act, "expected no dependencies")
}

func TestPlatformGenerate(t *testing.T) {
	cases := []struct {
		name             string
		input            string
		expectedContents []string
	}{
		{
			name: "aws",
			input: `aws
test_region`,
			expectedContents: []string{
				"aws",
				"test_region",
			},
		},
		{
			name: "libvirt",
			input: `libvirt
test_uri`,
			expectedContents: []string{
				"libvirt",
				"test_uri",
			},
		},
		{
			name: "case insensitive platform",
			input: `AWS
test_region`,
			expectedContents: []string{
				"aws",
				"test_region",
			},
		},
		{
			name: "invalid platform",
			input: `bad-platform
aws
test_region`,
			expectedContents: []string{
				"aws",
				"test_region",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			platform := &Platform{
				InputReader: bufio.NewReader(strings.NewReader(tc.input)),
			}
			deps := map[asset.Asset]*asset.State{}
			state, err := platform.Generate(deps)
			assert.NoError(t, err, "unexpected error generating platform")
			assert.NotNil(t, state, "expected non-nil asset state")
			act := make([]string, len(state.Contents))
			for i, c := range state.Contents {
				act[i] = string(c.Data)
			}
			assert.Equal(t, tc.expectedContents, act, "unexpected contents in asset state")
		})
	}
}
