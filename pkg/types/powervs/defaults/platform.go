package defaults

import (
	"github.com/openshift/installer/pkg/ipnet"

	"github.com/openshift/installer/pkg/types/powervs"
)

var (
	// DefaultMachineCIDR is the PowerVS default IP address space from
	// which to assign machine IPs.
	DefaultMachineCIDR = ipnet.MustParseCIDR("192.168.0.0/16")
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *powervs.Platform) {
}
