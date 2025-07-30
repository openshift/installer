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
	g := c.GCP
	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateGCPClusterHostedDNSInstall,
			Condition:       g.UserProvisionedDNS == dns.UserProvisionedDNSEnabled,
			Field:           field.NewPath("platform", "gcp", "userProvisionedDNS"),
		},
		{
			FeatureGateName: features.FeatureGateGCPCustomAPIEndpointsInstall,
			Condition:       len(g.ServiceEndpoints) > 0,
			Field:           field.NewPath("platform", "gcp", "serviceEndpoints"),
		},
	}
}
