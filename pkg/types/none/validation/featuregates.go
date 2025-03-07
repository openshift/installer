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
			// TODO mshitrit this is just a place holder for features.FeatureGateDualReplica until FeatureGateDualReplica PR is merged https://github.com/openshift/api/pull/2196
			FeatureGateName: features.FeatureGateGCPCustomAPIEndpoints,
			Condition:       c.None.FencingCredentials != nil,
			Field:           field.NewPath("platform", "none", "fencingCredentials"),
		},
	}
}
