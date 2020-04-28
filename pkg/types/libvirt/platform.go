package libvirt

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// URI is the identifier for the libvirtd connection.  It must be
	// reachable from both the host (where the installer is run) and the
	// cluster (where the cluster-API controller pod will be running).
	// Default is qemu+tcp://192.168.122.1/system
	//
	// +kubebuilder:default="qemu+tcp://192.168.122.1/system"
	// +optional
	URI string `json:"URI,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on libvirt for machine pools which do not define their
	// own platform configuration.
	// Default will set the image field to the latest RHCOS image.
	//
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network
	// +optional
	Network *Network `json:"network,omitempty"`
}

// Network is the configuration of the libvirt network.
type Network struct {
	// The interface make used for the network.
	// Default is tt0.
	//
	// +kubebuilder:default="tt0"
	// +optional
	IfName string `json:"if,omitempty"`
}
