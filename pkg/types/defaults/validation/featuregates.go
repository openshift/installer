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
		{
			FeatureGateName: features.FeatureGateOSStreams,
			Condition:       len(c.OSImageStream) != 0,
			Field:           field.NewPath("osImageStream"),
		},
		{
			FeatureGateName: features.FeatureGateOSStreams,
			Condition:       c.ControlPlane != nil && len(c.ControlPlane.OSImageStream) != 0,
			Field:           field.NewPath("controlPlane", "osImageStream"),
		},
		{
			FeatureGateName: features.FeatureGateOSStreams,
			Condition:       c.Arbiter != nil && len(c.Arbiter.OSImageStream) != 0,
			Field:           field.NewPath("arbiter", "osImageStream"),
		},
		{
			FeatureGateName: features.FeatureGateOSStreams,
			Condition: func() bool {
				for _, compute := range c.Compute {
					if len(compute.OSImageStream) != 0 {
						return true
					}
				}
				return false
			}(),
			Field: field.NewPath("compute", "osImageStream"),
		},
	}
}
