package capiutils

import (
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
)

var (
	defaultIPv4CIDR = ipnet.MustParseCIDR("10.0.0.0/16")
	defaultIPv6CIDR = ipnet.MustParseCIDR("fd00::/8")
)

// CIDRFromInstallConfig generates the CIDR from the install config,
// or returns the default CIDR if none is found.
func CIDRFromInstallConfig(installConfig *installconfig.InstallConfig) *ipnet.IPNet {
	if len(installConfig.Config.MachineNetwork) > 0 {
		return &installConfig.Config.MachineNetwork[0].CIDR
	}
	return defaultIPv4CIDR
}

// CIDRsFromInstallConfig generates multiple CIDRs from the install config,
// or returns the default IPv4 CIDR if none is found.
func CIDRsFromInstallConfig(installConfig *installconfig.InstallConfig) []ipnet.IPNet {
	var cidrs []ipnet.IPNet
	for _, machineNetwork := range installConfig.Config.MachineNetwork {
		cidrs = append(cidrs, machineNetwork.CIDR)
	}
	if len(cidrs) == 0 {
		// XXX: Do we even support single stack IPv6?
		cidrs = append(cidrs, *defaultIPv4CIDR)
	}
	return cidrs
}

// IsEnabled returns true if the feature gate is enabled.
func IsEnabled(installConfig *installconfig.InstallConfig) bool {
	// TODO(padillon): refactor to remove IsEnabled function.
	return true
}

// GenerateBoostrapMachineName generates the Cluster API Machine used for bootstrapping
// from the cluster ID and machine type.
func GenerateBoostrapMachineName(infraID string) string {
	return infraID + "-bootstrap"
}
