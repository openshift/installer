package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const testClusterName = "test-cluster"

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
			ic := &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: testClusterName,
				},
			}
			SetPlatformDefaults(tc.platform, ic)
			assert.Equal(t, tc.expected, tc.platform, "unexpected platform")
		})
	}
}
