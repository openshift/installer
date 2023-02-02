package defaults

import (
	"github.com/openshift/installer/pkg/ipnet"
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

func TestGetMachineCIDRDefault(t *testing.T) {
	cidr := GetMachineCIDR()
	assert.Equal(t, *defaultMachineCIDR, *cidr, "unexpected machine CIDR")
}

func TestGetMachineCIDROverride(t *testing.T) {
	t.Setenv("NUTANIX_MACHINE_CIDR_OVERRIDE", "10.40.0.0/16")
	cidr := GetMachineCIDR()
	assert.NotEqual(t, *defaultMachineCIDR, *cidr, "unexpected machine CIDR")
	expectedCIDR := ipnet.MustParseCIDR("10.40.0.0/16")
	assert.Equal(t, *expectedCIDR, *cidr, "unexpected machine CIDR")
}
