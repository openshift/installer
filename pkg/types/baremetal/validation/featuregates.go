package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	features "github.com/openshift/api/features"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateOnPremInternalDNSRecords,
			Condition:       c.BareMetal.InternalDNSRecords == configv1.InternalDNSRecordsDisabled,
			Field:           field.NewPath("platform", "baremetal", "internalDNSRecords"),
		},
	}
}
