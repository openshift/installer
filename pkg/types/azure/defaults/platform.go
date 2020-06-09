package defaults

import (
	"github.com/openshift/installer/pkg/types/azure"
)

var (
	// Overrides
	defaultMachineClass = map[string]string{}
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *azure.Platform) {
	if p.CloudName == "" {
		p.CloudName = azure.PublicCloud
	}
	if p.OutboundType == "" {
		p.OutboundType = azure.LoadbalancerOutboundType
	}
}

// getInstanceClass returns the instance "class" we should use for a given region.
func getInstanceClass(region string) string {
	if class, ok := defaultMachineClass[region]; ok {
		return class
	}

	return "Standard"
}
