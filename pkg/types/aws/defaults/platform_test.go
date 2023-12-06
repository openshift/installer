package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
)

func TestInstanceTypes(t *testing.T) {
	type testCase struct {
		name         string
		region       string
		architecture types.Architecture
		topology     configv1.TopologyMode
		expected     []string
		assert       func(*testCase)
	}
	cases := []testCase{
		{
			name:     "default instance types for AMD64",
			topology: configv1.HighlyAvailableTopologyMode,
			expected: []string{"m6i.xlarge", "m5.xlarge", "r5.xlarge", "c5.2xlarge", "m5.2xlarge", "c5d.2xlarge", "r5.2xlarge"},
			assert: func(tc *testCase) {
				instances := InstanceTypes(tc.region, tc.architecture, tc.topology)
				assert.Equal(t, tc.expected, instances, "unexepcted instance type for AMD64")
			},
		},
		{
			name:     "default instance types for AMD64",
			topology: configv1.SingleReplicaTopologyMode,
			expected: []string{"m6i.2xlarge", "m5.2xlarge", "r5.2xlarge", "c5.2xlarge", "m5.2xlarge", "c5d.2xlarge", "r5.2xlarge"},
			assert: func(tc *testCase) {
				instances := InstanceTypes(tc.region, tc.architecture, tc.topology)
				assert.Equal(t, tc.expected, instances, "unexepcted instance type for AMD64")
			},
		},
		{
			name:         "default instance types for ARM64",
			architecture: types.ArchitectureARM64,
			topology:     configv1.HighlyAvailableTopologyMode,
			expected:     []string{"m6g.xlarge"},
			assert: func(tc *testCase) {
				instances := InstanceTypes(tc.region, tc.architecture, tc.topology)
				assert.Equal(t, tc.expected, instances, "unexepcted instance type for ARM64")
			},
		},
		{
			name:         "default instance types for ARM64",
			architecture: types.ArchitectureARM64,
			topology:     configv1.SingleReplicaTopologyMode,
			expected:     []string{"m6g.2xlarge"},
			assert: func(tc *testCase) {
				instances := InstanceTypes(tc.region, tc.architecture, tc.topology)
				assert.Equal(t, tc.expected, instances, "unexepcted instance type for ARM64")
			},
		},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.assert(&tc)
		})
	}
}
