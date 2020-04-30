package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

func TestAWSResolver(t *testing.T) {
	overrides := []typesaws.ServiceEndpoint{{
		Name: "ec2",
		URL:  "test-ec2.local",
	}, {
		Name: "s3",
		URL:  "https://test-s3.local",
	}}

	cases := []struct {
		iservice, iregion string
		overrides         []typesaws.ServiceEndpoint
		expected          string
	}{{
		iservice: "ec2",
		iregion:  "us-east-1",
		expected: "https://ec2.us-east-1.amazonaws.com",
	}, {
		iservice:  "ec2",
		iregion:   "us-east-1",
		overrides: overrides,
		expected:  "test-ec2.local",
	}, {
		iservice:  "s3",
		iregion:   "us-east-1",
		overrides: overrides,
		expected:  "https://test-s3.local",
	}, {
		iservice:  "elasticloadbalancing",
		iregion:   "us-east-1",
		overrides: overrides,
		expected:  "https://elasticloadbalancing.us-east-1.amazonaws.com",
	}}
	for idx, test := range cases {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			resolvers := newAWSResolver(test.iregion, test.overrides)
			endpoint, err := resolvers.EndpointFor(test.iservice, test.iregion)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, endpoint.URL)
		})
	}
}
