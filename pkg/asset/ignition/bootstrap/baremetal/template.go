package baremetal

import (
	"fmt"
	"net"
	"strings"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// TemplateData holds data specific to templates used for the baremetal platform.
type TemplateData struct {
	// ProvisioningInterface holds the interface the bootstrap node will use to host the ProvisioningIP below.
	// When the provisioning network is disabled, this is the external baremetal network interface.
	ProvisioningInterface string

	// ProvisioningIP holds the IP the bootstrap node will use to service Ironic, TFTP, etc.
	ProvisioningIP string

	// ProvisioningIPv6 determines if we are using IPv6 or not.
	ProvisioningIPv6 bool

	// ProvisioningCIDR has the integer CIDR notation, e.g. 255.255.255.0 should be "24"
	ProvisioningCIDR int

	// ProvisioningDNSMasq determines if we start the dnsmasq service on the bootstrap node.
	ProvisioningDNSMasq bool

	// ProvisioningDHCPRange has the DHCP range, if DHCP is not external. Otherwise it
	// should be blank.
	ProvisioningDHCPRange string

	// ProvisioningDHCPAllowList contains a space-separated list of all of the control plane's boot
	// MAC addresses. Requests to bootstrap DHCP from other hosts will be ignored.
	ProvisioningDHCPAllowList string

	// IronicUsername contains the username for authentication to Ironic
	IronicUsername string

	// IronicUsername contains the password for authentication to Ironic
	IronicPassword string

	// BaremetalEndpointOverride contains the url for the baremetal endpoint
	BaremetalEndpointOverride string

	// BaremetalIntrospectionEndpointOverride contains the url for the baremetal introspection endpoint
	BaremetalIntrospectionEndpointOverride string

	// ClusterOSImage contains 4 URLs to download RHCOS live iso, kernel, rootfs and initramfs
	ClusterOSImage string
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *baremetal.Platform, networks []types.MachineNetworkEntry, ironicUsername, ironicPassword string) *TemplateData {
	var templateData TemplateData

	templateData.ProvisioningIP = config.BootstrapProvisioningIP
	templateData.BaremetalEndpointOverride = fmt.Sprintf("http://%s/v1", net.JoinHostPort(config.APIVIP, "6385"))
	templateData.BaremetalIntrospectionEndpointOverride = fmt.Sprintf("http://%s/v1", net.JoinHostPort(config.APIVIP, "5050"))

	if config.ProvisioningNetwork != baremetal.DisabledProvisioningNetwork {
		cidr, _ := config.ProvisioningNetworkCIDR.Mask.Size()
		templateData.ProvisioningCIDR = cidr
		templateData.ProvisioningIPv6 = config.ProvisioningNetworkCIDR.IP.To4() == nil
		templateData.ProvisioningInterface = "ens4"
		templateData.ProvisioningDNSMasq = true
	}

	switch config.ProvisioningNetwork {
	case baremetal.ManagedProvisioningNetwork:
		cidr, _ := config.ProvisioningNetworkCIDR.Mask.Size()

		// When provisioning network is managed, we set a DHCP range including
		// netmask for dnsmasq.
		templateData.ProvisioningDHCPRange = fmt.Sprintf("%s,%d", config.ProvisioningDHCPRange, cidr)

		var dhcpAllowList []string
		for _, host := range config.Hosts {
			if host.Role == "master" {
				dhcpAllowList = append(dhcpAllowList, host.BootMACAddress)
			}
		}
		templateData.ProvisioningDHCPAllowList = strings.Join(dhcpAllowList, " ")
	case baremetal.DisabledProvisioningNetwork:
		templateData.ProvisioningInterface = "ens3"
		templateData.ProvisioningDNSMasq = false

		if templateData.ProvisioningIP != "" {
			for _, network := range networks {
				if network.CIDR.Contains(net.ParseIP(templateData.ProvisioningIP)) {
					templateData.ProvisioningIPv6 = network.CIDR.IP.To4() == nil

					cidr, _ := network.CIDR.Mask.Size()
					templateData.ProvisioningCIDR = cidr
				}
			}
		}
	}

	templateData.IronicUsername = ironicUsername
	templateData.IronicPassword = ironicPassword
	templateData.ClusterOSImage = config.ClusterOSImage

	return &templateData
}
