package clusterapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckSourceRanges(t *testing.T) {
	cases := []struct {
		name     string
		region   string
		external bool
		expected []string
	}{
		{
			name:     "public region internal",
			region:   "us-central1",
			external: false,
			expected: []string{"35.191.0.0/16", "130.211.0.0/22"},
		},
		{
			name:     "public region external appends google ranges",
			region:   "us-central1",
			external: true,
			expected: []string{"35.191.0.0/16", "130.211.0.0/22", "209.85.152.0/22", "209.85.204.0/22"},
		},
		{
			name:     "gcd berlin uses sovereign prober ranges",
			region:   "u-germany-northeast1",
			external: false,
			expected: []string{"34.3.144.0/23", "34.3.151.0/26", "34.3.151.64/26", "136.124.104.0/22", "136.124.108.0/22"},
		},
		{
			name:     "gcd berlin ignores external (private-only)",
			region:   "u-germany-northeast1",
			external: true,
			expected: []string{"34.3.144.0/23", "34.3.151.0/26", "34.3.151.64/26", "136.124.104.0/22", "136.124.108.0/22"},
		},
		{
			name:     "gcd france uses sovereign prober ranges",
			region:   "u-france-east1",
			external: false,
			expected: []string{"177.222.80.0/23", "177.222.87.0/26", "177.222.87.64/26", "136.124.104.0/22", "136.124.108.0/22"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, healthCheckSourceRanges(tc.region, tc.external))
		})
	}
}
