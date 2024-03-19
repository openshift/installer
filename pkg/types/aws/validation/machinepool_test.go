package validation

import (
	"fmt"
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
			name: "valid zone io instance",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-east-1b"},
				EC2RootVolume: aws.EC2RootVolume{
					Type: "io2",
					Size: 128,
					IOPS: 10000,
				},
			},
		},
		{
			name: "valid zone gp3 instance no iops",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-east-1b"},
				EC2RootVolume: aws.EC2RootVolume{
					Type: "gp3",
					Size: 128,
				},
			},
		},
		{
			name: "valid zone gp3 instance with iops",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-east-1b"},
				EC2RootVolume: aws.EC2RootVolume{
					Type: "gp3",
					Size: 128,
					IOPS: 10000,
				},
			},
		},
		{
			name: "valid zone gp2 instance no iops",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-east-1b"},
				EC2RootVolume: aws.EC2RootVolume{
					Type: "gp2",
					Size: 128,
				},
			},
		},
		{
			name: "invalid zone gp2 instance with iops",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-east-1b"},
				EC2RootVolume: aws.EC2RootVolume{
					Type: "gp2",
					Size: 128,
					IOPS: 10000,
				},
			},
			expected: fmt.Sprintf("test-path.iops: Invalid value: 10000: iops not supported for type gp2"),
		},
		{
			name: "invalid zone",
			pool: &aws.MachinePool{
				Zones: []string{"us-east-1a", "us-west-1a"},
				EC2RootVolume: aws.EC2RootVolume{
					Type: "io2",
					Size: 128,
					IOPS: 10000,
				},
			},
			expected: `^test-path\.zones\[1]: Invalid value: "us-west-1a": Zone not in configured region \(us-east-1\)$`,
		},
		{
			name: "invalid volume type",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: "bad-volume-type",
					Size: 128,
				},
			},
			expected: fmt.Sprintf("test-path.type: Invalid value: \"bad-volume-type\": failed to find volume type bad-volume-type"),
		},
		{
			name: "invalid volume size using zero",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: "io2",
					Size: 0,
					IOPS: 10000,
				},
			},
			expected: fmt.Sprintf("test-path.size: Invalid value: 0: volume size value must be a positive number"),
		},
		{
			name: "invalid volume size using negative",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: "gp3",
					Size: -1,
					IOPS: 10000,
				},
			},
			expected: fmt.Sprintf("test-path.size: Invalid value: -1: volume size value must be a positive number"),
		},
		{
			name: "invalid metadata auth option",
			pool: &aws.MachinePool{
				EC2Metadata: aws.EC2Metadata{
					Authentication: "foobarbaz",
				},
			},
			expected: `^test-path\.authentication: Invalid value: \"foobarbaz\": must be either Required or Optional$`,
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

func Test_ValdidateSecurityGroups(t *testing.T) {
	tooManySecurityGroups := make([]string, 0, maxUserSecurityGroupsCount+1)
	for i := 0; i < maxUserSecurityGroupsCount+1; i++ {
		tooManySecurityGroups = append(tooManySecurityGroups, fmt.Sprintf("sg-valid-%d", i))
	}
	cases := []struct {
		name     string
		platform *aws.Platform
		pool     *aws.MachinePool
		err      string
	}{
		{
			name: "valid security group config",
			platform: &aws.Platform{
				Region:  "us-east-1",
				Subnets: []string{"valid-subnet-1", "valid-subnet-2"},
			},
			pool: &aws.MachinePool{
				AdditionalSecurityGroupIDs: []string{
					"sg-valid-security-group",
				},
			},
		},
		{
			name:     "invalid security group config",
			platform: &aws.Platform{Region: "us-east-1"},
			pool:     &aws.MachinePool{AdditionalSecurityGroupIDs: []string{"sg-valid-security-group"}},
			err:      "test-path.platform.subnets: Required value: subnets must be provided when additional security groups are present",
		},
		{
			name: "invalid security group config exceeds maximum",
			platform: &aws.Platform{
				Region:  "us-east-1",
				Subnets: []string{"valid-subnet-1", "valid-subnet-2"},
			},
			pool: &aws.MachinePool{
				AdditionalSecurityGroupIDs: tooManySecurityGroups,
			},
			err: "test-path: Too many: 11: must have at most 10 items",
		},
		{
			name: "valid maximum security group config",
			platform: &aws.Platform{
				Region:  "us-east-1",
				Subnets: []string{"valid-subnet-1", "valid-subnet-2"},
			},
			pool: &aws.MachinePool{
				AdditionalSecurityGroupIDs: tooManySecurityGroups[:maxUserSecurityGroupsCount],
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSecurityGroups(tc.platform, tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.err, err)
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
