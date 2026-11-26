package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"

	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func TestInternalLoadBalancerFromPlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *gcptypes.Platform
		expected *capg.LoadBalancer
	}{
		{
			name:     "nil platform",
			platform: nil,
		},
		{
			name:     "unset load balancer",
			platform: &gcptypes.Platform{Region: "us-east1"},
		},
		{
			name: "empty client access",
			platform: &gcptypes.Platform{
				LoadBalancer: &gcptypes.LoadBalancer{},
			},
		},
		{
			name: "global client access",
			platform: &gcptypes.Platform{
				LoadBalancer: &gcptypes.LoadBalancer{
					ClientAccess: gcptypes.ClientAccessGlobal,
				},
			},
			expected: &capg.LoadBalancer{InternalAccess: capg.InternalAccessGlobal},
		},
		{
			name: "local client access",
			platform: &gcptypes.Platform{
				LoadBalancer: &gcptypes.LoadBalancer{
					ClientAccess: gcptypes.ClientAccessLocal,
				},
			},
			expected: &capg.LoadBalancer{InternalAccess: capg.InternalAccessRegional},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, internalLoadBalancerFromPlatform(tc.platform))
		})
	}
}
