package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func defaultPlatform() *vsphere.Platform {
	return &vsphere.Platform{}
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *vsphere.Platform
		expected *vsphere.Platform
	}{
		{
			name:     "empty",
			platform: &vsphere.Platform{},
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
