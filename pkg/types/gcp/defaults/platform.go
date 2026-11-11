package defaults

import (
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/network"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *gcp.Platform) {
	if p == nil {
		return
	}

	if p.IPFamily == "" {
		p.IPFamily = network.IPv4
	}

	if gcpDmp := p.DefaultMachinePlatform; gcpDmp != nil {
		if ek := gcpDmp.EncryptionKey; ek != nil {
			if kms := ek.KMSKey; kms != nil {
				if kms.ProjectID == "" {
					kms.ProjectID = p.ProjectID
				}
				if kms.Location == "" {
					kms.Location = p.Region
				}
			}
		}
	}
}
