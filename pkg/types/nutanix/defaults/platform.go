package defaults

import (
	"github.com/openshift/installer/pkg/types/nutanix"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *nutanix.Platform) {
	if p.PrismAPICallTimeout == nil {
		timeout := nutanix.DefaultPrismAPICallTimeout
		p.PrismAPICallTimeout = &timeout
	}
}
