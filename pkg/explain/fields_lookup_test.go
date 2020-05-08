package explain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lookup(t *testing.T) {
	schema, err := loadSchema(loadCRD(t))
	assert.NoError(t, err)

	cases := []struct {
		path []string

		desc string
		err  string
	}{{
		desc: `InstallConfig is the configuration for an OpenShift install.`,
	}, {
		path: []string{"publish"},
		desc: `Publish controls how the user facing endpoints of the cluster like the Kubernetes API, OpenShift routes etc. are exposed. When no strategy is specified, the strategy is "External".`,
	}, {
		path: []string{"publish", "unknown"},
		err:  `invalid field unknown, no such property found`,
	}, {
		path: []string{"platform"},
		desc: `Platform is the configuration for the specific platform upon which to perform the installation.`,
	}, {
		path: []string{"platform", "aws"},
		desc: `AWS is the configuration used when installing on AWS.`,
	}, {
		path: []string{"platform", "azure"},
		desc: `Azure is the configuration used when installing on Azure.`,
	}, {
		path: []string{"platform", "aws", "region"},
		desc: `Region specifies the AWS region where the cluster will be created.`,
	}, {
		path: []string{"platform", "aws", "subnets"},
		desc: `Subnets specifies existing subnets (by ID) where cluster resources will be created.  Leave unset to have the installer create subnets in a new VPC on your behalf.`,
	}, {
		path: []string{"platform", "aws", "userTags"},
		desc: `UserTags additional keys and values that the installer will add as tags to all resources that it creates. Resources created by the cluster itself may not include these tags.`,
	}, {
		path: []string{"platform", "aws", "serviceEndpoints"},
		desc: `ServiceEndpoints list contains custom endpoints which will override default service endpoint of AWS Services. There must be only one ServiceEndpoint for a service.`,
	}, {
		path: []string{"platform", "aws", "serviceEndpoints", "url"},
		desc: `URL is fully qualified URI with scheme https, that overrides the default generated endpoint for a client. This must be provided and cannot be empty.`,
	}}
	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			got, err := lookup(schema, test.path)
			if test.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.desc, got.Description)
			} else {
				assert.Regexp(t, test.err, err)
			}
		})
	}
}
