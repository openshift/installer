package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/featuregates"
)

// GatedFeatures determines all of the vSphere install config fields that should
// be validated to ensure that the proper featuregate is enabled when the field is used.
func GatedFeatures(c *types.InstallConfig) []featuregates.GatedInstallConfigFeature {
	v := c.VSphere

	multiNetworksFound := false
	nodeNetworkingDefined := v.NodeNetworking != nil

	for _, fd := range v.FailureDomains {
		if len(fd.Topology.Networks) > 1 {
			multiNetworksFound = true
		}
	}

	cpDef := c.ControlPlane.Platform.VSphere
	computeDefs := c.Compute

	return []featuregates.GatedInstallConfigFeature{
		{
			FeatureGateName: features.FeatureGateVSphereStaticIPs,
			Condition:       len(v.Hosts) > 0,
			Field:           field.NewPath("platform", "vsphere", "hosts"),
		},
		{
			FeatureGateName: features.FeatureGateVSphereMultiVCenters,
			Condition:       len(v.VCenters) > 1,
			Field:           field.NewPath("platform", "vsphere", "vcenters"),
		},
		{
			FeatureGateName: features.FeatureGateVSphereMultiNetworks,
			Condition:       multiNetworksFound,
			Field:           field.NewPath("platform", "vsphere", "failureDomains", "topology", "networks"),
		},
		{
			FeatureGateName: features.FeatureGateVSphereMultiNetworks,
			Condition:       nodeNetworkingDefined,
			Field:           field.NewPath("platform", "vsphere", "nodeNetworking"),
		},
		{
			FeatureGateName: features.FeatureGateVSphereMultiDisk,
			Condition:       cpDef != nil && len(cpDef.AdditionalDisks) > 0, // Here we need to check disk count
			Field:           field.NewPath("controlPlane", "platform", "vsphere", "additionalDisks"),
		},
		{
			FeatureGateName: features.FeatureGateVSphereMultiDisk,
			Condition:       hasAdditionalDisks(computeDefs), // Here we need to check disk count
			Field:           field.NewPath("compute", "platform", "vsphere", "additionalDisks"),
		},
	}
}

func hasAdditionalDisks(pool []types.MachinePool) bool {
	foundAdditionalDisks := false
	for _, machine := range pool {
		if machine.Platform.VSphere != nil && len(machine.Platform.VSphere.AdditionalDisks) > 0 {
			foundAdditionalDisks = true
			break
		}
	}
	return foundAdditionalDisks
}
