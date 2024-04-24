package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/featuregates"
	"github.com/openshift/installer/pkg/types/gcp"
)

// GatedFeatures determines all of the install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	g := c.GCP
	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateGCPLabelsTags,
			Condition:       len(g.UserLabels) > 0,
			Field:           field.NewPath("platform", "gcp", "userLabels"),
		},
		{
			FeatureGateName: features.FeatureGateGCPLabelsTags,
			Condition:       len(g.UserTags) > 0,
			Field:           field.NewPath("platform", "gcp", "userTags"),
		},
		{
			FeatureGateName: features.FeatureGateGCPClusterHostedDNS,
			Condition:       g.UserProvisionedDNS == gcp.UserProvisionedDNSEnabled,
			Field:           field.NewPath("platform", "gcp", "userProvisionedDNS"),
		},
	}
}
