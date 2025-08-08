package defaults

import "github.com/openshift/installer/pkg/types/gcp"

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *gcp.Platform) {
	if gcpDmp := p.DefaultMachinePlatform; gcpDmp != nil {
		if ek := gcpDmp.EncryptionKey; ek != nil {
			if ek.KMSKey.ProjectID == "" {
				ek.KMSKey.ProjectID = p.ProjectID
			}

			if ek.KMSKey.Location == "" {
				ek.KMSKey.Location = p.Region
			}
		}
	}
}
