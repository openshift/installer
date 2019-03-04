package conversion

import (
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

// ConvertInstallConfig is modeled after the k8s conversion schemes, which is
// how deprecated values are upconverted.
func ConvertInstallConfig(config *types.InstallConfig) {
	ConvertNetworking(config)
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

	if len(netconf.ServiceNetwork) == 0 && netconf.DeprecatedServiceCIDR != nil {
		netconf.ServiceNetwork = []ipnet.IPNet{*netconf.DeprecatedServiceCIDR}
	}

	// Convert type to networkType if the latter is missing
	if netconf.NetworkType == "" {
		netconf.NetworkType = netconf.DeprecatedType
	}

	// Convert hostSubnetLength to hostPrefix
	for i, entry := range netconf.ClusterNetwork {
		if entry.HostPrefix == 0 && entry.DeprecatedHostSubnetLength != 0 {
			_, size := entry.CIDR.Mask.Size()
			netconf.ClusterNetwork[i].HostPrefix = int32(size) - entry.DeprecatedHostSubnetLength
		}
	}
}
