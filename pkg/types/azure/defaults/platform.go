package defaults

import (
	"github.com/openshift/installer/pkg/types/azure"
	dnstypes "github.com/openshift/installer/pkg/types/dns"
)

var (
	// Overrides
	defaultMachineClass = map[string]string{}

	// AzurestackMinimumDiskSize is the minimum disk size value for azurestack.
	AzurestackMinimumDiskSize int32 = 128
	// AzurestackMaximumDiskSize is the maximum disk size value for azurestack.
	AzurestackMaximumDiskSize int32 = 1023
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *azure.Platform) {
	if p.CloudName == "" {
		p.CloudName = azure.PublicCloud
	}
	if p.OutboundType == "" {
		p.OutboundType = azure.LoadbalancerOutboundType
	}
	if p.UserProvisionedDNS == "" {
		p.UserProvisionedDNS = dnstypes.UserProvisionedDNSDisabled
	}
}

// getInstanceClass returns the instance "class" we should use for a given region.
func getInstanceClass(region string) string {
	if class, ok := defaultMachineClass[region]; ok {
		return class
	}

	return "Standard"
}
