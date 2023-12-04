package defaults

import (
	dnstypes "github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *gcp.Platform) {
	if p.UserProvisionedDNS == "" {
		p.UserProvisionedDNS = dnstypes.UserProvisionedDNSDisabled
	}
}
