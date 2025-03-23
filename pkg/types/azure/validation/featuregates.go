package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	cp := c.ControlPlane.Platform
	defMp := c.Platform.Azure.DefaultMachinePlatform
	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateMachineAPIMigration,
			Condition:       cp.Azure != nil && cp.Azure.Identity != nil && cp.Azure.Identity.Type == capz.VMIdentitySystemAssigned,
			Field:           field.NewPath("controlPlane", "azure", "identity", "systemAssignedIdentityRole"),
		},
		{
			FeatureGateName: features.FeatureGateMachineAPIMigration,
			Condition:       defMp != nil && defMp.Identity != nil && defMp.Identity.Type == capz.VMIdentitySystemAssigned,
			Field:           field.NewPath("platform", "azure", "defaultMachinePlatform", "identity", "systemAssignedIdentityRole"),
		},
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
	}
}
