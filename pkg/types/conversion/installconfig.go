package conversion

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	operv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// ConvertInstallConfig is modeled after the k8s conversion schemes, which is
// how deprecated values are upconverted.
// This updates the APIVersion to reflect the fact that we've internally
// upconverted.
func ConvertInstallConfig(config *types.InstallConfig) error {
	// check that the version is convertible
	switch config.APIVersion {
	case types.InstallConfigVersion, "v1beta3", "v1beta4":
		// works
	case "":
		return field.Required(field.NewPath("apiVersion"), "no version was provided")
	default:
		return field.Invalid(field.NewPath("apiVersion"), config.APIVersion, fmt.Sprintf("cannot upconvert from version %s", config.APIVersion))
	}
	ConvertNetworking(config)

	switch config.Platform.Name() {
	case baremetal.Name:
		ConvertBaremetal(config)
	}

	config.APIVersion = types.InstallConfigVersion
	return nil
}

// ConvertNetworking upconverts deprecated fields in networking
func ConvertNetworking(config *types.InstallConfig) {
	if config.Networking == nil {
		return
	}

	netconf := config.Networking

	if len(netconf.ClusterNetwork) == 0 {
		netconf.ClusterNetwork = netconf.DeprecatedClusterNetworks
	}

	if len(netconf.MachineNetwork) == 0 && netconf.DeprecatedMachineCIDR != nil {
		netconf.MachineNetwork = []types.MachineNetworkEntry{
			{CIDR: *netconf.DeprecatedMachineCIDR},
		}
	}

	if len(netconf.ServiceNetwork) == 0 && netconf.DeprecatedServiceCIDR != nil {
		netconf.ServiceNetwork = []ipnet.IPNet{*netconf.DeprecatedServiceCIDR}
	}

	// Convert type to networkType if the latter is missing
	if netconf.NetworkType == "" {
		netconf.NetworkType = netconf.DeprecatedType
	}

	// Recognize the default network plugin name regardless of capitalization, for
	// backward compatibility
	if strings.ToLower(netconf.NetworkType) == strings.ToLower(string(operv1.NetworkTypeOpenShiftSDN)) {
		netconf.NetworkType = string(operv1.NetworkTypeOpenShiftSDN)
	}

	// Convert hostSubnetLength to hostPrefix
	for i, entry := range netconf.ClusterNetwork {
		if entry.HostPrefix == 0 && entry.DeprecatedHostSubnetLength != 0 {
			_, size := entry.CIDR.Mask.Size()
			netconf.ClusterNetwork[i].HostPrefix = int32(size) - entry.DeprecatedHostSubnetLength
		}
	}
}

// ConvertBaremetal upconverts deprecated fields in the baremetal
// platform. ProvisioningDHCPExternal has been replaced by setting
// the ProvisioningNetwork field to "Unmanaged"
func ConvertBaremetal(config *types.InstallConfig) {
	if config.Platform.BareMetal.DeprecatedProvisioningDHCPExternal == true && config.Platform.BareMetal.ProvisioningNetwork == "" {
		config.Platform.BareMetal.ProvisioningNetwork = baremetal.UnmanagedProvisioningNetwork
	}
}
