package aws

import (
	"testing"

	"github.com/openshift/installer/pkg/types/aws"
	"github.com/stretchr/testify/assert"
)

func TestCloudProviderConfig(t *testing.T) {
	customEndpoint := []aws.CustomEndpoint{
		aws.CustomEndpoint{
			Service: "ec2",
			URL:     "https://ec2.foo.bar",
		},
	}
	platform := &aws.Platform{
		Region:               "us-east-1",
		CustomRegionOverride: customEndpoint,
	}
	expected := `[Global]

[ServiceOverride "1"]
Service = ec2
Region  = us-east-1
URL     = https://ec2.foo.bar

`

	json, err := CloudProviderConfig(platform)
	assert.NoError(t, err, "failed to create cloud provider config")
	assert.Equal(t, expected, json, "unexpected cloud provider config")
}
