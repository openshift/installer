package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
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

	return []featuregates.GatedInstallConfigFeature{
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
			FeatureGateName: features.FeatureGateOnPremDNSRecords,
			Condition:       v.DNSRecordsType == configv1.DNSRecordsTypeExternal,
			Field:           field.NewPath("platform", "vsphere", "dnsRecordsType"),
		},
	}
}
