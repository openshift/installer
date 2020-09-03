package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

func TestValidateMachinePool(t *testing.T) {
	platform := &aws.Platform{Region: "us-east-1"}
	cases := []struct {
		name     string
		pool     *aws.MachinePool
		expected string
	}{
		{
			name: "empty",
			pool: &aws.MachinePool{},
		},
		{
			name: "valid zone",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-east-1b"},
			},
		},
		{
			name: "invalid zone",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-west-1a"},
			},
			expected: `^test-path\.zones\[1]: Invalid value: "us-west-1a": Zone not in configured region \(us-east-1\)$`,
		},
		{
			name: "valid iops",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					IOPS: 10,
				},
			},
		},
		{
			name: "invalid iops",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					IOPS: -10,
				},
			},
			expected: `^test-path\.iops: Invalid value: -10: Storage IOPS must be positive$`,
		},
		{
			name: "valid size",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Size: 10,
				},
			},
		},
		{
			name: "invalid size",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Size: -10,
				},
			},
			expected: `^test-path\.size: Invalid value: -10: Storage size must be positive$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(platform, tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

func Test_validateAMIID(t *testing.T) {
	cases := []struct {
		platform *aws.Platform
		pool     *aws.MachinePool

		err string
	}{{
		platform: &aws.Platform{Region: "us-east-1"},
	}, {
		platform: &aws.Platform{Region: "us-gov-east-1"},
		err:      `^test-path: Required value: AMI ID must be provided for regions .*$`,
	}, {
		platform: &aws.Platform{Region: "cn-north-1"},
		err:      `^test-path: Required value: AMI ID must be provided for regions .*$`,
	}, {
		platform: &aws.Platform{Region: "us-gov-east-1", AMIID: "ami"},
	}, {
		platform: &aws.Platform{Region: "us-gov-east-1", DefaultMachinePlatform: &aws.MachinePool{AMIID: "ami"}},
	}, {
		platform: &aws.Platform{Region: "us-gov-east-1"},
		pool:     &aws.MachinePool{AMIID: "ami"},
	}}
	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			err := ValidateAMIID(test.platform, test.pool, field.NewPath("test-path")).ToAggregate()
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.err, err)
			}
		})
	}
}
