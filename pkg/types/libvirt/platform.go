package libvirt

import (
	"net"
)

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// URI is the identifier for the libvirtd connection.  It must be
	// reachable from both the host (where the installer is run) and the
	// cluster (where the cluster-API controller pod will be running).
	URI string `json:"URI"`

	// Image is the URL to the OS image.
	// E.g. "http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz"
	Image string `json:"image"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on libvirt for machine pools which do not define their
	// own platform configuration.
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network
	Network Network `json:"network"`

	// MasterIPs
	MasterIPs []net.IP `json:"masterIPs,omitempty"`
}
