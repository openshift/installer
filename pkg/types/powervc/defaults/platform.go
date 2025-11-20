package defaults

import (
	"os"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervc"
)

const (
	// DefaultCloudName is the default name of the cloud in clouds.yaml file.
	DefaultCloudName = "powervc"
	// DualStackVIPsPortTag is the identifier of VIPs Port with dual-stack addresses.
	DualStackVIPsPortTag = "-dual-stack-vips-port"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *powervc.Platform, n *types.Networking) {
	if p.Cloud == "" {
		p.Cloud = os.Getenv("OS_CLOUD")
		if p.Cloud == "" {
			p.Cloud = DefaultCloudName
		}
	}
}
