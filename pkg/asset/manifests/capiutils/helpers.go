package capiutils

import (
	"encoding/binary"
	"math"
	"net"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/pkg/errors"
)

var (
	defaultCIDR         = ipnet.MustParseCIDR("10.0.0.0/16")
	maxAddressByNetwork = 1024
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
	return installConfig.Config.EnabledFeatureGates().Enabled(v1.FeatureGateClusterAPIInstall)
}

// GenerateBoostrapMachineName generates the Cluster API Machine used for bootstrapping
// from the cluster ID and machine type.
func GenerateBoostrapMachineName(infraID string) string {
	return infraID + "-bootstrap"
}

// GenerateSubnetCIDRs generates the CIDR blocks for the subnets from the zones.
// Modified from CAPA https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/baf8d59ff495f86cded519cf81bc0ea00e1b0126/pkg/internal/cidr/cidr.go#L29-L67
func GenerateSubnetCIDRs(CIDR string, numSubnets int) ([]string, error) {

	subnetCIDRs := make([]string, 0, numSubnets)

	_, parent, err := net.ParseCIDR(CIDR)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse CIDR")
	}

	subnetBits := math.Ceil(math.Log2(float64(numSubnets)))

	networkLen, addrLen := parent.Mask.Size()
	modifiedNetworkLen := networkLen + int(subnetBits)

	if modifiedNetworkLen > addrLen {
		return nil, errors.Errorf("cidr %s cannot accommodate %d subnets", CIDR, numSubnets)
	}

	var subnets []*net.IPNet
	for i := 0; i < numSubnets; i++ {
		ip4 := parent.IP.To4()
		if ip4 == nil {
			return nil, errors.Errorf("unexpected IP address type: %s", parent)
		}

		n := binary.BigEndian.Uint32(ip4)
		n += uint32(i) << uint(32-modifiedNetworkLen)
		subnetIP := make(net.IP, len(ip4))
		binary.BigEndian.PutUint32(subnetIP, n)

		ip := &net.IPNet{
			IP:   subnetIP,
			Mask: net.CIDRMask(modifiedNetworkLen, 32),
		}
		subnets = append(subnets, ip)
		subnetCIDRs = append(subnetCIDRs, ip.String())
	}

	return subnetCIDRs, nil

	// return subnetCIDRs
}
