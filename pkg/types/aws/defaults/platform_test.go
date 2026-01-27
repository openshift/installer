package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/network"
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

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name     string
		platform *aws.Platform
		expected *aws.Platform
	}{
		{
			name:     "empty platform should default IPFamily to IPv4",
			platform: &aws.Platform{},
			expected: &aws.Platform{
				IPFamily: network.IPv4,
			},
		},
		{
			name: "IPFamily already set should not be overridden",
			platform: &aws.Platform{
				IPFamily: network.DualStackIPv4Primary,
			},
			expected: &aws.Platform{
				IPFamily: network.DualStackIPv4Primary,
				LBType:   configv1.NLB, // LBType gets set to NLB when IPFamily is dual-stack
			},
		},
		{
			name: "LBType should default to NLB for DualStackIPv4Primary",
			platform: &aws.Platform{
				IPFamily: network.DualStackIPv4Primary,
			},
			expected: &aws.Platform{
				IPFamily: network.DualStackIPv4Primary,
				LBType:   configv1.NLB,
			},
		},
		{
			name: "LBType should default to NLB for DualStackIPv6Primary",
			platform: &aws.Platform{
				IPFamily: network.DualStackIPv6Primary,
			},
			expected: &aws.Platform{
				IPFamily: network.DualStackIPv6Primary,
				LBType:   configv1.NLB,
			},
		},
		{
			name: "LBType should not be set for IPv4",
			platform: &aws.Platform{
				IPFamily: network.IPv4,
			},
			expected: &aws.Platform{
				IPFamily: network.IPv4,
			},
		},
		{
			name: "LBType already set should not be overridden",
			platform: &aws.Platform{
				IPFamily: network.DualStackIPv4Primary,
				LBType:   configv1.Classic,
			},
			expected: &aws.Platform{
				IPFamily: network.DualStackIPv4Primary,
				LBType:   configv1.Classic,
			},
		},
		{
			name: "empty IPFamily should default to IPv4 and LBType remains empty",
			platform: &aws.Platform{
				LBType: "",
			},
			expected: &aws.Platform{
				IPFamily: network.IPv4,
				LBType:   "",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform)
			assert.Equal(t, tc.expected.IPFamily, tc.platform.IPFamily)
			assert.Equal(t, tc.expected.LBType, tc.platform.LBType)
		})
	}
}
