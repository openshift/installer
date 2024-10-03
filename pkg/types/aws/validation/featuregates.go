package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateAWSClusterHostedDNS,
			Condition:       c.AWS.UserProvisionedDNS == dns.UserProvisionedDNSEnabled,
			Field:           field.NewPath("platform", "aws", "userProvisionedDNS"),
		},
	}
}
