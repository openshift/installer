package validation

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
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
			expected: `^\Qtest-path.userTags[Name]: Invalid value: "test-cluster": Name key is not allowed for user defined tags\E$`,
		},
		{
			name: "invalid userTags, key with kubernetes.io/cluster/",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"kubernetes.io/cluster/test-cluster": "shared",
				},
			},
			expected: `^\Qtest-path.userTags[kubernetes.io/cluster/test-cluster]: Invalid value: "shared": Keys with prefix 'kubernetes.io/cluster/' are not allowed for user defined tags\E$`,
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
			expected: `^test-path\.hostedZoneRole: Invalid value: "test-hosted-zone-role": may not specify a role to assume for hosted zone operations without also specifying a hosted zone$`,
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
			expected: `^test-path\.hostedZoneRole: Forbidden: when specifying a hostedZoneRole, either Passthrough or Manual credential mode must be specified$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, tc.credMode, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.credMode, err)
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
