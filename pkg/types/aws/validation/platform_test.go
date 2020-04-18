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
				Region: "",
			},
			expected: `^test-path\.region: Required value: region must be specified$`,
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
						IOPS: -10,
					},
				},
			},
			expected: `^test-path\.defaultMachinePlatform\.iops: Invalid value: -10: Storage IOPS must be positive$`,
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
