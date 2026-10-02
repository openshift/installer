package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	cp := c.ControlPlane.Platform
	defMp := c.Platform.Azure.DefaultMachinePlatform
	azure := c.Azure

	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateMachineAPIMigration,
			Condition:       cp.Azure != nil && cp.Azure.Identity != nil && cp.Azure.Identity.UserAssignedIdentities != nil && len(cp.Azure.Identity.UserAssignedIdentities) > 1,
			Field:           field.NewPath("controlPlane", "azure", "identity", "userAssignedIdentities"),
		},
		{
			FeatureGateName: features.FeatureGateMachineAPIMigration,
			Condition:       defMp != nil && defMp.Identity != nil && defMp.Identity.UserAssignedIdentities != nil && len(defMp.Identity.UserAssignedIdentities) > 1,
			Field:           field.NewPath("platform", "azure", "defaultMachinePlatform", "identity", "userAssignedIdentities"),
		},
		{
			FeatureGateName: features.FeatureGateAzureDualStackInstall,
			Condition:       azure.IPFamily.DualStackEnabled(),
			Field:           field.NewPath("platform", "azure", "ipFamily"),
		},
	}
}
