package defaults

import (
	"fmt"
	"net"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"

	"github.com/apparentlymart/go-cidr/cidr"
)

// Defaults for the baremetal platform.
const (
	LibvirtURI              = "qemu:///system"
	ExternalBridge          = "baremetal"
	ProvisioningBridge      = "provisioning"
	HardwareProfile         = "default"
	APIVIP                  = ""
	IngressVIP              = ""
	ProvisioningInterface   = "ens3"
	ProvisioningNetworkCIDR = "172.22.0.0/24"
	CachedImageURL          = "http://192.168.111.1/images"
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

	if p.ProvisioningNetworkCIDR == "" {
		p.ProvisioningNetworkCIDR = ProvisioningNetworkCIDR
	}

	_, provNet, _ := net.ParseCIDR(p.ProvisioningNetworkCIDR)

	if p.BootstrapProvisioningIP == "" {
		// Default to .2 address for CIDR e.g 172.22.0.2
		ip, err := cidr.Host(provNet, 2)
		if err == nil {
			p.BootstrapProvisioningIP = ip.String()
		}
	}

	if p.ClusterProvisioningIP == "" {
		// Default to .3 address for CIDR e.g 172.22.0.3
		ip, err := cidr.Host(provNet, 3)
		if err == nil {
			p.ClusterProvisioningIP = ip.String()
		}
	}

	if p.ExternalBridge == "" {
		p.ExternalBridge = ExternalBridge
	}

	if p.ProvisioningBridge == "" {
		p.ProvisioningBridge = ProvisioningBridge
	}

	for _, host := range p.Hosts {
		if host.HardwareProfile == "" {
			host.HardwareProfile = HardwareProfile
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

	if p.ProvisioningInterface == "" {
		p.ProvisioningInterface = ProvisioningInterface
	}

	if p.ProvisioningDHCPStart == "" {
		// Default to .20 address for CIDR e.g 172.22.0.20
		ip, err := cidr.Host(provNet, 20)
		if err == nil {
			p.ProvisioningDHCPStart = ip.String()
		}
	}

	if p.ProvisioningDHCPEnd == "" {
		// Default to .200 address for CIDR e.g 172.22.0.200
		ip, err := cidr.Host(provNet, 200)
		if err == nil {
			p.ProvisioningDHCPEnd = ip.String()
		}
	}

	if p.CachedImageURL == "" {
		p.CachedImageURL = CachedImageURL
	}
}
