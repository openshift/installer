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

	// DnsmasqOptions is the dnsmasq options to be used when installing on
	// libvirt.
	//
	// +optional
	DnsmasqOptions []DnsmasqOption `json:"dnsmasqOptions,omitempty"`
}

// DnsmasqOption contains the name and value of the option to configure in the
// libvirt network.
type DnsmasqOption struct {
	// The dnsmasq option name. A full list of options and an explanation for
	// each can be found in /etc/dnsmasq.conf
	//
	// +optional
	Name string `json:"name,omitempty"`

	// The value that is being set for the particular option.
	//
	// +optional
	Value string `json:"value,omitempty"`
}
