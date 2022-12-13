package defaults

import (
	"fmt"
	"math/rand"
	"net"
	"sort"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/openshift/installer/pkg/ipnet"
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
	BootMode                = baremetal.UEFI
	ExternalMACAddress      = ""
	ProvisioningMACAddress  = ""
)

// Wrapper for net.LookupHost so we can override in the test
var lookupHost = func(host string) (addrs []string, err error) {
	return net.LookupHost(host)
}

// GenerateMAC a randomized MAC address with the libvirt prefix
func GenerateMAC() string {
	buf := make([]byte, 3)
	rand.Seed(time.Now().UnixNano())
	rand.Read(buf)

	// set local bit and unicast
	buf[0] = (buf[0] | 2) & 0xfe

	// avoid libvirt-reserved addresses
	if buf[0] == 0xfe {
		buf[0] = 0xee
	}

	return fmt.Sprintf("52:54:00:%02x:%02x:%02x", buf[0], buf[1], buf[2])
}

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *baremetal.Platform, c *types.InstallConfig) {
	if p.LibvirtURI == "" {
		p.LibvirtURI = LibvirtURI
	}

	if p.ExternalMACAddress == "" {
		p.ExternalMACAddress = GenerateMAC()
	}
	if p.ProvisioningMACAddress == "" {
		p.ProvisioningMACAddress = GenerateMAC()
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

	if len(p.APIVIPs) == 0 && p.DeprecatedAPIVIP == "" {
		// This name should resolve to exactly one address
		if vip, err := lookupHost("api." + c.ClusterDomain()); err == nil {
			p.APIVIPs = []string{vip[0]}
		}
	}

	if len(p.IngressVIPs) == 0 && p.DeprecatedIngressVIP == "" {
		// This name should resolve to exactly one address
		if vip, err := lookupHost("test.apps." + c.ClusterDomain()); err == nil {
			p.IngressVIPs = []string{vip[0]}
		}
	}

	if p.Hosts != nil {
		sort.SliceStable(p.Hosts, func(i, j int) bool {
			return p.Hosts[i].CompareByRole(p.Hosts[j])
		})
	}
}
