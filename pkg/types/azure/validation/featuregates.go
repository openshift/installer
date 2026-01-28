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
			FeatureGateName: features.FeatureGateAzureMultiDisk,
			Condition:       defMp != nil && len(defMp.DataDisks) != 0,
			Field:           field.NewPath("platform", "azure", "defaultMachinePlatform", "dataDisks"),
		},
		{
			FeatureGateName: features.FeatureGateAzureMultiDisk,
			Condition:       cp.Azure != nil && len(cp.Azure.DataDisks) != 0,
			Field:           field.NewPath("controlPlane", "azure", "dataDisks"),
		},
		{
			FeatureGateName: features.FeatureGateAzureMultiDisk,
			Condition: func() bool {
				computeMachinePool := c.Compute
				for _, compute := range computeMachinePool {
					if compute.Platform.Azure != nil {
						if len(compute.Platform.Azure.DataDisks) != 0 {
							return true
						}
					}
				}
				return false
			}(),
			Field: field.NewPath("compute", "azure", "dataDisks"),
		},
		{
			FeatureGateName: features.FeatureGateAzureClusterHostedDNSInstall,
			Condition:       azure.UserProvisionedDNS == dns.UserProvisionedDNSEnabled,
			Field:           field.NewPath("platform", "azure", "userProvisionedDNS"),
		},
		{
			FeatureGateName: features.FeatureGateAzureDualStackInstall,
			Condition:       azure.IPFamily.DualStackEnabled(),
			Field:           field.NewPath("platform", "azure", "ipFamily"),
		},
	}
}
