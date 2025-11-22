package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

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
		{
			name: "host placement any available",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{
					Affinity: ptr.To(aws.HostAffinityAnyAvailable),
				},
			},
		},
		{
			name: "valid dedicated hosts",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{
					Affinity: ptr.To(aws.HostAffinityDedicatedHost),
					DedicatedHost: []aws.DedicatedHost{
						{
							ID: "h-09dcf61cb388b0149",
						},
					},
				},
			},
		},
		{
			name: "invalid dedicated hosts - missing hostID",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{
					Affinity: ptr.To(aws.HostAffinityDedicatedHost),
					DedicatedHost: []aws.DedicatedHost{
						{},
					},
				},
			},
			expected: `^test-path.hostPlacement.dedicatedHost\[0].id: Required value: a hostID must be specified when configuring 'dedicatedHost'$`,
		},
		{
			name: "invalid - hostPlacement without affinity",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{},
			},
			expected: `^test-path.hostPlacement.affinity: Required value: affinity is required when hostPlacement is configured$`,
		},
		{
			name: "invalid unknown affinity",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{
					Affinity: ptr.To(aws.HostAffinity("Unknown")),
				},
			},
			expected: `^test-path.hostPlacement.affinity: Unsupported value: "Unknown": supported values: "AnyAvailable", "DedicatedHost"$`,
		},
		{
			name: "any available with dedicated host set",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{
					Affinity:      ptr.To(aws.HostAffinityAnyAvailable),
					DedicatedHost: []aws.DedicatedHost{{ID: "h-09dcf61cb388b0149"}},
				},
			},
			expected: `^test-path.hostPlacement.dedicatedHost: Required value: dedicatedHost is required when 'affinity' is set to DedicatedHost, and forbidden otherwise$`,
		},
		{
			name: "invalid - DedicatedHost affinity without dedicatedHost",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{
					Affinity: ptr.To(aws.HostAffinityDedicatedHost),
				},
			},
			expected: `^test-path.hostPlacement.dedicatedHost: Required value: dedicatedHost is required when 'affinity' is set to DedicatedHost, and forbidden otherwise$`,
		},
		{
			name: "invalid dedicated host - bad hostID",
			pool: &aws.MachinePool{
				HostPlacement: &aws.HostPlacement{
					Affinity: ptr.To(aws.HostAffinityDedicatedHost),
					DedicatedHost: []aws.DedicatedHost{
						{ID: "h-09DCFABC"},
					},
				},
			},
			expected: `^test-path.hostPlacement.dedicatedHost\[0\].id: Invalid value: "h-09DCFABC": id must start with 'h-' followed by 17 lowercase hexadecimal characters \(0-9 and a-f\)$`,
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
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "valid-subnet-1"},
						{ID: "valid-subnet-2"},
					},
				},
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
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "valid-subnet-1"},
						{ID: "valid-subnet-2"},
					},
				},
			},
			pool: &aws.MachinePool{
				AdditionalSecurityGroupIDs: tooManySecurityGroups,
			},
			err: "test-path: Too many: 11: must have at most 10 items",
		},
		{
			name: "valid maximum security group config",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "valid-subnet-1"},
						{ID: "valid-subnet-2"},
					},
				},
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

func Test_validateCPUOptions(t *testing.T) {
	cases := []struct {
		name string
		pool *aws.MachinePool
		err  string
	}{{
		name: "confidential compute policy set to AMD SEV-SNP",
		pool: &aws.MachinePool{
			CPUOptions: &aws.CPUOptions{
				ConfidentialCompute: ptr.To(aws.ConfidentialComputePolicySEVSNP),
			},
		},
	}, {
		name: "confidential compute disabled",
		pool: &aws.MachinePool{
			CPUOptions: &aws.CPUOptions{
				ConfidentialCompute: ptr.To(aws.ConfidentialComputePolicyDisabled),
			},
		},
	}, {
		name: "empty confidential compute policy",
		pool: &aws.MachinePool{
			CPUOptions: &aws.CPUOptions{
				ConfidentialCompute: ptr.To(aws.ConfidentialComputePolicy("")),
			},
		},
		err: `^test-path.confidentialCompute: Unsupported value: "": supported values: "Disabled", "AMDEncryptedVirtualizationNestedPaging"$`,
	}, {
		name: "invalid confidential compute policy",
		pool: &aws.MachinePool{
			CPUOptions: &aws.CPUOptions{
				ConfidentialCompute: ptr.To(aws.ConfidentialComputePolicy("invalid")),
			},
		},
		err: `^test-path.confidentialCompute: Unsupported value: "invalid": supported values: "Disabled", "AMDEncryptedVirtualizationNestedPaging"$`,
	}, {
		name: "empty cpu options",
		pool: &aws.MachinePool{
			CPUOptions: &aws.CPUOptions{},
		},
		err: `^test-path.cpuOptions: Invalid value: "{}": At least one field must be set if cpuOptions is provided$`,
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCPUOptions(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.err, err)
			}
		})
	}
}
