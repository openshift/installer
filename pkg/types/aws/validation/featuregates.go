package validation

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	gatedFeatures := []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateAWSClusterHostedDNSInstall,
			Condition:       c.AWS.UserProvisionedDNS == dns.UserProvisionedDNSEnabled,
			Field:           field.NewPath("platform", "aws", "userProvisionedDNS"),
		},
		{
			FeatureGateName: features.FeatureGateAWSDualStackInstall,
			Condition:       c.AWS.IPFamily.DualStackEnabled(),
			Field:           field.NewPath("platform", "aws", "ipFamily"),
		},
		{
			FeatureGateName: features.FeatureGateAWSEuropeanSovereignCloudInstall,
			Condition:       strings.HasPrefix(c.AWS.Region, "eusc-"),
			Field:           field.NewPath("platform", "aws", "region"),
		},
	}

	// Check if any compute pool has dedicated hosts configured
	for idx, compute := range c.Compute {
		if compute.Platform.AWS != nil && compute.Platform.AWS.HostPlacement != nil {
			gatedFeatures = append(gatedFeatures, featuregates.GatedInstallConfigFeature{
				FeatureGateName: features.FeatureGateAWSDedicatedHosts,
				Condition:       true,
				Field:           field.NewPath("compute").Index(idx).Child("platform", "aws", "hostPlacement"),
			})
			break // Only need to add the feature gate once
		}
	}

	return gatedFeatures
}
