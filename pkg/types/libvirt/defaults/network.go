package defaults

import (
	"github.com/openshift/installer/pkg/ipnet"

	"github.com/openshift/installer/pkg/types/libvirt"
)

const (
	defaultIfName = "tt0"
)

var (
	// DefaultMachineCIDR is the libvirt default IP address space from
	// which to assign machine IPs.
	DefaultMachineCIDR = ipnet.MustParseCIDR("192.168.126.0/24")
)

// SetNetworkDefaults sets the defaults for the network.
func SetNetworkDefaults(n *libvirt.Network) {
	if n.IfName == "" {
		n.IfName = defaultIfName
	}
}
