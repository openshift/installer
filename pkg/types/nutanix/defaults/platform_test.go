package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/nutanix"
)

func defaultPlatform() *nutanix.Platform {
	return &nutanix.Platform{
		DNSRecordsType: configv1.DNSRecordsTypeInternal,
	}
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
