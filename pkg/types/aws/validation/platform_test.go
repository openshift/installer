package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *aws.Platform
		expected string
	}{
		{
			name: "minimal",
			platform: &aws.Platform{
				Region: "us-east-1",
			},
		},
		{
			name: "invalid region",
			platform: &aws.Platform{
				Region: "bad-region",
			},
			expected: "^test-path.region: Unsupported value: \"bad-region\": supported values: .*\"us-east-1\"",
		},
		{
			name: "valid region without RHCOS",
			platform: &aws.Platform{
				Region: "us-gov-west-1",
			},
			expected: "^test-path.region: Invalid value: \"us-gov-west-1\": no RHCOS AMIs found in \"us-gov-west-1\" (.*us-east-1.*)$",
		},
		{
			name: "valid region without RHCOS, but with an explicit AMI",
			platform: &aws.Platform{
				AMIID:  "ami-123",
				Region: "us-gov-west-1",
			},
		},
		{
			name: "valid machine pool",
			platform: &aws.Platform{
				Region:                 "us-east-1",
				DefaultMachinePlatform: &aws.MachinePool{},
			},
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
			expected: "^test-path.defaultMachinePlatform.iops: Invalid value: -10: Storage IOPS must be positive$",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
