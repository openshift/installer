package capiutils

import (
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
)

var (
	defaultCIDR = ipnet.MustParseCIDR("10.0.0.0/16")
	// AnyIPv4CidrBlock is the CIDR block to match all IPv4 addresses.
	AnyIPv4CidrBlock = ipnet.MustParseCIDR("0.0.0.0/0")
	// AnyIPv6CidrBlock is the CIDR block to match all IPv6 addresses.
	AnyIPv6CidrBlock = ipnet.MustParseCIDR("::/0")
)

// CIDRFromInstallConfig generates the CIDR from the install config,
// or returns the default CIDR if none is found.
func CIDRFromInstallConfig(installConfig *installconfig.InstallConfig) *ipnet.IPNet {
	if len(installConfig.Config.MachineNetwork) > 0 {
		return &installConfig.Config.MachineNetwork[0].CIDR
	}
	return defaultCIDR
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

// MachineCIDRsFromInstallConfig returns the machine network CIDRs from the install config.
func MachineCIDRsFromInstallConfig(ic *installconfig.InstallConfig) []ipnet.IPNet {
	cidrs := make([]ipnet.IPNet, 0, len(ic.Config.MachineNetwork))
	for _, cidr := range ic.Config.MachineNetwork {
		cidrs = append(cidrs, cidr.CIDR)
	}
	return cidrs
}

// CIDRsToString returns the string representation of network CIDRs.
func CIDRsToString(cidrs []ipnet.IPNet) []string {
	cidrStrings := make([]string, 0, len(cidrs))
	for _, cidr := range cidrs {
		cidrStrings = append(cidrStrings, cidr.String())
	}
	return cidrStrings
}
