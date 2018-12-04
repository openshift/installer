package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/aws"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *aws.Platform
		valid    bool
	}{
		{
			name: "minimal",
			platform: &aws.Platform{
				Region: "us-east-1",
			},
			valid: true,
		},
		{
			name: "invalid region",
			platform: &aws.Platform{
				Region: "bad-region",
			},
			valid: false,
		},
		{
			name: "valid machine pool",
			platform: &aws.Platform{
				Region:                 "us-east-1",
				DefaultMachinePlatform: &aws.MachinePool{},
			},
			valid: true,
		},
		{
			name: "invalid machine pool",
			platform: &aws.Platform{
				Region: "us-east-1",
				DefaultMachinePlatform: &aws.MachinePool{
					EC2RootVolume: aws.EC2RootVolume{
						IOPS: -10,
					},
				},
			},
			valid: false,
		},
		{
			name: "valid CIDR",
			platform: &aws.Platform{
				Region:       "us-east-1",
				VPCCIDRBlock: ipnet.MustParseCIDR("192.168.0.0/16"),
			},
			valid: true,
		},
		{
			name: "invalid CIDR",
			platform: &aws.Platform{
				Region:       "us-east-1",
				VPCCIDRBlock: ipnet.MustParseCIDR("0.0.0.0/16"),
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
