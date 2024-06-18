package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the OpenStack install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	v := c.VSphere
	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateVSphereStaticIPs,
			Condition:       len(v.Hosts) > 0,
			Field:           field.NewPath("platform", "vsphere", "hosts"),
		},
	}
}
