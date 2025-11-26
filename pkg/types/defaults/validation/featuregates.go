package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateDualReplica,
			Condition:       c.ControlPlane != nil && c.ControlPlane.Fencing != nil,
			Field:           field.NewPath("platform", "none", "fencingCredentials"),
		},
		{
			FeatureGateName: features.FeatureGateMultiDiskSetup,
			Condition:       c.ControlPlane != nil && len(c.ControlPlane.DiskSetup) != 0,
			Field:           field.NewPath("controlPlane", "diskSetup"),
		},
		{
			FeatureGateName: features.FeatureGateMultiDiskSetup,
			Condition: func() bool {
				computeMachinePool := c.Compute
				for _, compute := range computeMachinePool {
					if len(compute.DiskSetup) != 0 {
						return true
					}
				}
				return false
			}(),
			Field: field.NewPath("compute", "diskSetup"),
		},
	}
}
