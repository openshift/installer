package libvirt

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// URI is the identifier for the libvirtd connection.  It must be
	// reachable from both the host (where the installer is run) and the
	// cluster (where the cluster-API controller pod will be running).
	// +optional
	// Default is qemu+tcp://192.168.122.1/system
	URI string `json:"URI,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on libvirt for machine pools which do not define their
	// own platform configuration.
	// +optional
	// Default will set the image field to the latest RHCOS image.
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network
	// +optional
	Network *Network `json:"network,omitempty"`
}
