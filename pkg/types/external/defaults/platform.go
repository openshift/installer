package defaults

import "github.com/openshift/installer/pkg/types/external"

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *external.Platform) {
	p.PlatformName = "Unknown"
	p.CloudControllerManager = external.CloudControllerManagerTypeNone
}
