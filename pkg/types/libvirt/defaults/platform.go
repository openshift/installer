package defaults

import (
	"github.com/openshift/installer/pkg/types/libvirt"
)

const (
	// DefaultURI is the default URI of the libvirtd connection.
	DefaultURI = "qemu+tcp://192.168.122.1/system"
	// DefaultPoolPath directory path for the libvirt storage pool
	DefaultPoolPath = "/var/lib/libvirt/openshift-images"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *libvirt.Platform) {
	if p.URI == "" {
		p.URI = DefaultURI
	}
	if p.Network == nil {
		p.Network = &libvirt.Network{}
	}
	SetNetworkDefaults(p.Network)
	if p.PoolPath == "" {
		p.PoolPath = DefaultPoolPath
	}
}
