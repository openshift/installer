package baremetal

import (
	"github.com/openshift/installer/pkg/types/baremetal"
)

// TemplateData holds data specific to templates used for the baremetal platform.
type TemplateData struct {
	// ProvisioningIP holds the IP the bootstrap node will use to service Ironic, TFTP, etc.
	ProvisioningIP string

	// ProvisioningIPv6 determines if we are using IPv6 or not.
	ProvisioningIPv6 bool

	// ProvisioningCIDR has the integer CIDR notation, e.g. 255.255.255.0 should be "24"
	ProvisioningCIDR int

	// ProvisioningDHCPRange has the DHCP range, if DHCP is not external. Otherwise it
	// should be blank.
	ProvisioningDHCPRange string
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *baremetal.Platform) *TemplateData {
	var templateData TemplateData

	templateData.ProvisioningIP = config.BootstrapProvisioningIP

	cidr, _ := config.ProvisioningNetworkCIDR.Mask.Size()
	templateData.ProvisioningCIDR = cidr

	templateData.ProvisioningIPv6 = config.ProvisioningNetworkCIDR.IP.To4() == nil

	if !config.ProvisioningDHCPExternal {
		templateData.ProvisioningDHCPRange = config.ProvisioningDHCPRange
	}

	return &templateData
}
