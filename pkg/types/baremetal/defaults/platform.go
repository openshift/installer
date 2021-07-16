package defaults

import (
	"fmt"
	"github.com/openshift/installer/pkg/ipnet"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// Defaults for the baremetal platform.
const (
	LibvirtURI              = "qemu:///system"
	ProvisioningNetworkCIDR = "172.22.0.0/24"
	ExternalBridge          = "baremetal"
	ProvisioningBridge      = "provisioning"
	HardwareProfile         = "default"
	APIVIP                  = ""
	IngressVIP              = ""
	BootMode                = baremetal.UEFI
)

// Wrapper for net.LookupHost so we can override in the test
var lookupHost = func(host string) (addrs []string, err error) {
	return net.LookupHost(host)
}

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *baremetal.Platform, c *types.InstallConfig) {
	if p.LibvirtURI == "" {
		p.LibvirtURI = LibvirtURI
	}

	if p.ProvisioningNetwork == "" {
		p.ProvisioningNetwork = baremetal.ManagedProvisioningNetwork
	}

	switch p.ProvisioningNetwork {
	case baremetal.DisabledProvisioningNetwork:
		if p.ClusterProvisioningIP != "" {
			for _, network := range c.MachineNetwork {
				if network.CIDR.Contains(net.ParseIP(p.ClusterProvisioningIP)) {
					p.ProvisioningNetworkCIDR = &network.CIDR
				}
			}
		}
	default:
		if p.ProvisioningNetworkCIDR == nil {
			p.ProvisioningNetworkCIDR = ipnet.MustParseCIDR(ProvisioningNetworkCIDR)
		}

		if p.BootstrapProvisioningIP == "" {
			// Default to the second address in provisioning network, e.g 172.22.0.2
			ip, err := cidr.Host(&p.ProvisioningNetworkCIDR.IPNet, 2)
			if err == nil {
				p.BootstrapProvisioningIP = ip.String()
			}
		}

		if p.ClusterProvisioningIP == "" {
			// Default to the third address in provisioning network, e.g 172.22.0.3
			ip, err := cidr.Host(&p.ProvisioningNetworkCIDR.IPNet, 3)
			if err == nil {
				p.ClusterProvisioningIP = ip.String()
			}
		}

		if p.ProvisioningBridge == "" {
			p.ProvisioningBridge = ProvisioningBridge
		}
	}

	// If network is managed, and user didn't specify a range, let's use the rest of the subnet range from
	// the 10th address until the end.
	if p.ProvisioningNetwork == baremetal.ManagedProvisioningNetwork && p.ProvisioningDHCPRange == "" {
		startIP, _ := cidr.Host(&p.ProvisioningNetworkCIDR.IPNet, 10)
		_, broadcastIP := cidr.AddressRange(&p.ProvisioningNetworkCIDR.IPNet)
		endIP := cidr.Dec(broadcastIP)
		p.ProvisioningDHCPRange = fmt.Sprintf("%s,%s", startIP, endIP)
	}

	if p.ExternalBridge == "" {
		p.ExternalBridge = ExternalBridge
	}

	for _, host := range p.Hosts {
		if host.HardwareProfile == "" {
			host.HardwareProfile = HardwareProfile
		}

		if host.BootMode == "" {
			host.BootMode = BootMode
		}
	}

	if p.APIVIP == APIVIP {
		// This name should resolve to exactly one address
		vip, err := lookupHost("api." + c.ClusterDomain())
		if err != nil {
			// This will fail validation and abort the install
			p.APIVIP = fmt.Sprintf("DNS lookup failure: %s", err.Error())
		} else {
			p.APIVIP = vip[0]
		}
	}

	if p.IngressVIP == IngressVIP {
		// This name should resolve to exactly one address
		vip, err := lookupHost("test.apps." + c.ClusterDomain())
		if err != nil {
			// This will fail validation and abort the install
			p.IngressVIP = fmt.Sprintf("DNS lookup failure: %s", err.Error())
		} else {
			p.IngressVIP = vip[0]
		}
	}
}
