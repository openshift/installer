package powervs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVPCZonesForVPCRegion(t *testing.T) {
	cases := []struct {
		name      string
		vpcRegion string
		expected  []string
		expectErr bool
	}{
		{
			name:      "default VPC region of a PowerVS region",
			vpcRegion: "us-south",
			expected:  []string{"us-south-1", "us-south-2", "us-south-3"},
		},
		{
			// A PowerVS region can be paired with a VPC region other than its own, e.g. dal
			// with us-east. The zones must follow the VPC region, not the PowerVS region.
			name:      "non-default VPC region",
			vpcRegion: "us-east",
			expected:  []string{"us-east-1", "us-east-2", "us-east-3"},
		},
		{
			name:      "region with only two zones",
			vpcRegion: "eu-es",
			expected:  []string{"eu-es-1", "eu-es-2"},
		},
		{
			name:      "unknown VPC region",
			vpcRegion: "nowhere",
			expectErr: true,
		},
		{
			name:      "empty VPC region",
			vpcRegion: "",
			expectErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := VPCZonesForVPCRegion(tc.vpcRegion)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}
