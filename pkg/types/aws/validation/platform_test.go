package validation

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/stretchr/testify/assert"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *aws.Platform
		expected string
		credMode types.CredentialsMode
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
				Region: "",
			},
			expected: `^test-path\.region: Required value: region must be specified$`,
		},
		{
			name: "hosted zone with subnets",
			platform: &aws.Platform{
				Region:     "us-east-1",
				Subnets:    []string{"test-subnet"},
				HostedZone: "test-hosted-zone",
			},
		},
		{
			name: "hosted zone without subnets",
			platform: &aws.Platform{
				Region:     "us-east-1",
				HostedZone: "test-hosted-zone",
			},
			expected: `^test-path\.hostedZone: Invalid value: "test-hosted-zone": may not use an existing hosted zone when not using existing subnets$`,
		},
		{
			name: "invalid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "/path/some",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[0\]\.url: Invalid value: "(.*)": host cannot be empty, empty host provided$`,
		},
		{
			name: "invalid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "https://test-ec2.random.local/path/some",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[0\]\.url: Invalid value: "(.*)": no path or request parameters must be provided, "/path/some" was provided$`,
		},
		{
			name: "invalid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "https://test-ec2.random.local?foo=some",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[0\]\.url: Invalid value: "(.*)": no path or request parameters must be provided, "/\?foo=some" was provided$`,
		},
		{
			name: "valid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "test-ec2.random.local",
				}},
			},
		},
		{
			name: "valid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "https://test-ec2.random.local",
				}},
			},
		},
		{
			name: "duplicate service endpoints",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "test-ec2.random.local",
				}, {
					Name: "s3",
					URL:  "test-ec2.random.local",
				}, {
					Name: "ec2",
					URL:  "test-ec2.random.local",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[2\]\.name: Invalid value: "ec2": duplicate service endpoint not allowed for ec2, service endpoint already defined at test-path\.serviceEndpoints\[0\]$`,
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
						Type: "io1",
						IOPS: -10,
						Size: 128,
					},
				},
			},
			expected: `^test-path.*iops: Invalid value: -10: iops must be a positive number$`,
		},
		{
			name: "invalid userTags, Name key",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"Name": "test-cluster",
				},
			},
			expected: `^\Qtest-path.userTags[Name]: Invalid value: "test-cluster": "Name" key is not allowed for user defined tags\E$`,
		},
		{
			name: "invalid userTags, value with invalid characters",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"usage-user": "cloud-team-rebase-bot[bot]",
				},
			},
			expected: `^\Qtest-path.userTags[usage-user]: Invalid value: "cloud-team-rebase-bot[bot]": value contains invalid characters`,
		},
		{
			name: "valid userTags, value with spaces",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"test-key": "this test has spaces",
				},
				PropagateUserTag: true,
			},
		},
		{
			name: "invalid userTags, key with kubernetes.io/cluster/",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"kubernetes.io/cluster/test-cluster": "shared",
				},
			},
			expected: `^\Qtest-path.userTags[kubernetes.io/cluster/test-cluster]: Invalid value: "shared": key is in the kubernetes.io namespace\E$`,
		},
		{
			name: "valid userTags",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"app": "production",
				},
			},
		},
		{
			name: "too many userTags",
			platform: &aws.Platform{
				Region:   "us-east-1",
				UserTags: generateTooManyUserTags(),
			},
			expected: fmt.Sprintf(`^\Qtest-path.userTags: Too many: %d: must have at most %d items`, userTagLimit+1, userTagLimit),
		},
		{
			name: "hosted zone role without hosted zone should error",
			platform: &aws.Platform{
				Region:         "us-east-1",
				HostedZoneRole: "test-hosted-zone-role",
			},
			expected: `^\[test-path\.hostedZoneRole: Invalid value: "test-hosted-zone-role": may not specify a role to assume for hosted zone operations without also specifying a hosted zone, test-path\.credentialsMode: Forbidden: when specifying a hostedZoneRole, either Passthrough or Manual credential mode must be specified\]$`,
		},
		{
			name:     "valid hosted zone & role should not throw an error",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
			},
		},
		{
			name: "hosted zone role without credential mode should error",
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
			},
			expected: `^test-path\.credentialsMode: Forbidden: when specifying a hostedZoneRole, either Passthrough or Manual credential mode must be specified$`,
		},
		{
			name:     "valid when LBType is present while providing eipAllocations",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-1234567890abcdefa"},
				},
			},
		},
		{
			name:     "valid when LBType is present while providing eipAllocations",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         "",
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-1234567890abcdefa"},
				},
			},
			expected: `^test-path\.aws\.lbType: Required value: lbType NLB must be specified$`,
		},
		{
			name:     "Valid EIP Allocations",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-1234567890abcdefa"},
				},
			},
		},
		{
			name:     "Duplicate Allocation IDs",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet1", "test-subnet2"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-1234567890abcdefa", "eipalloc-1234567890abcdefa"},
				},
			},
			expected: `^test-path\.aws\.eipAllocations\.ingressNetworkLoadBalancer: Invalid value: "eipalloc-1234567890abcdefa": cannot have duplicate EIP Allocation IDs?`,
		},
		{
			name:     "Too Many Allocation IDs",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet1", "test-subnet2"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-0123456789abcdef0",
						"eipalloc-0abcdef1234567891",
						"eipalloc-0abcdef1234567892",
						"eipalloc-0abcdef1234567893",
						"eipalloc-0abcdef1234567894",
						"eipalloc-0abcdef1234567895",
						"eipalloc-0abcdef1234567896",
						"eipalloc-0abcdef1234567897",
						"eipalloc-0abcdef1234567898",
						"eipalloc-0abcdef1234567899",
						"eipalloc-0abcdef123456789a"},
				},
			},
			expected: `^test-path\.aws\.eipAllocations\.ingressNetworkLoadBalancer: Too many: 11: must have at most 10 items?`,
		},
		{
			name:     "Invalid Prefix",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"zipalloc-1234567890abcdefa"},
				},
			},
			expected: `^test-path\.aws\.eipAllocations\.ingressNetworkLoadBalancer: Invalid value: "zipalloc-1234567890abcdefa": eipAllocations should start with 'eipalloc-'?`,
		},
		{
			name:     "Invalid Characters",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-0123456789abcdezz"},
				},
			},
			expected: `^test-path\.aws\.eipAllocations\.ingressNetworkLoadBalancer: Invalid value: "eipalloc-0123456789abcdezz": eipAllocations must be 'eipalloc-' followed by exactly 17 hexadecimal characters (0-9, a-f, A-F)?`,
		},
		{
			name:     "Invalid Length (Too Short)",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-01234567"},
				},
			},
			expected: `test-path\.aws\.eipAllocations\.ingressNetworkLoadBalancer: Invalid value: "eipalloc-01234567": invalid EIP allocation ID length`,
		},
		{
			name:     "Invalid Length (Too Long)",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region:         "us-east-1",
				Subnets:        []string{"test-subnet"},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
				LBType:         configv1.NLB,
				EIPAllocations: &aws.EIPAllocations{
					IngressNetworkLoadBalancer: []aws.EIPAllocation{"eipalloc-0123456789abcdef012345"},
				},
			},
			expected: `test-path\.aws\.eipAllocations\.ingressNetworkLoadBalancer: Invalid value: "eipalloc-0123456789abcdef012345": invalid EIP allocation ID length`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, tc.credMode, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

func TestValidateTag(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		value     string
		expectErr bool
	}{{
		name:  "valid",
		key:   "test-key",
		value: "test-value",
	}, {
		name:      "invalid characters in key",
		key:       "bad-key***",
		value:     "test-value",
		expectErr: true,
	}, {
		name:      "invalid characters in value",
		key:       "test-key",
		value:     "bad-value***",
		expectErr: true,
	}, {
		name:      "empty key",
		key:       "",
		value:     "test-value",
		expectErr: true,
	}, {
		name:      "empty value",
		key:       "test-key",
		value:     "",
		expectErr: true,
	}, {
		name:      "key too long",
		key:       strings.Repeat("a", 129),
		value:     "test-value",
		expectErr: true,
	}, {
		name:      "value too long",
		key:       "test-key",
		value:     strings.Repeat("a", 257),
		expectErr: true,
	}, {
		name:      "key in kubernetes.io namespace",
		key:       "kubernetes.io/cluster/some-cluster",
		value:     "owned",
		expectErr: true,
	}, {
		name:      "key in openshift.io namespace",
		key:       "openshift.io/some-key",
		value:     "some-value",
		expectErr: true,
	}, {
		name:      "key in openshift.io subdomain namespace",
		key:       "other.openshift.io/some-key",
		value:     "some-value",
		expectErr: true,
	}, {
		name:  "key in namespace similar to openshift.io",
		key:   "otheropenshift.io/some-key",
		value: "some-value",
	}, {
		name:  "key with openshift.io in path",
		key:   "some-domain/openshift.io/some-key",
		value: "some-value",
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTag(tc.key, tc.value)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func generateTooManyUserTags() map[string]string {
	tags := map[string]string{}
	for i := 0; i <= userTagLimit; i++ {
		tags[strconv.Itoa(i)] = strconv.Itoa(i)
	}
	return tags
}
