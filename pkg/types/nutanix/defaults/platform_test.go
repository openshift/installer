package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/nutanix"
)

const testClusterName = "test-cluster"

func defaultPlatform() *nutanix.Platform {
	return &nutanix.Platform{}
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *nutanix.Platform
		expected *nutanix.Platform
	}{
		{
			name:     "empty",
			platform: &nutanix.Platform{},
			expected: defaultPlatform(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform)
			assert.Equal(t, tc.expected, tc.platform, "unexpected platform")
		})
	}
}
