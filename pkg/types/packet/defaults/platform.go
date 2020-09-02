package defaults

import "github.com/openshift/installer/pkg/types/packet"

// Defaults for the packet platform.
const (
	// TODO(displague) what API? metadata?
	DefaultURI = "https://api.packet.com"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *packet.Platform) {
}
