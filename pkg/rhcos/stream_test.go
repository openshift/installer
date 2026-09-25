//go:build !(okd || scos)

package rhcos

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

func TestBuildDefaultOSImageStreamUsesResolvedVersion(t *testing.T) {
	originalRaw := version.Raw
	t.Cleanup(func() { version.Raw = originalRaw })

	tests := []struct {
		name     string
		raw      string
		expected types.OSImageStream
	}{
		{name: "OpenShift 4", raw: "v4.18.3", expected: types.OSImageStreamRHCOS9},
		{name: "OpenShift 5", raw: "v5.0.0", expected: types.OSImageStreamRHCOS10},
		{name: "invalid build version", raw: "was not built correctly", expected: types.OSImageStreamRHCOS10},
		{name: "zero build version", raw: "0.0", expected: types.OSImageStreamRHCOS10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version.Raw = tt.raw
			assert.Equal(t, tt.expected, BuildDefaultOSImageStream())
		})
	}
}
