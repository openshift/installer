package baremetal

// BMC stores the information about a baremetal host's management controller.
type BMC struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

// Host stores all the configuration data for a baremetal host.
type Host struct {
	Name            string `json:"name,omitempty"`
	BMC             BMC    `json:"bmc"`
	Role            string `json:"role"`
	BootMACAddress  string `json:"bootMACAddress"`
	HardwareProfile string `json:"hardwareProfile"`
}

// Image stores details about the locations of various images needed for deployment.
// FIXME: This should be determined by the installer once Ironic and image downloading occurs in bootstrap VM.
type Image struct {
	Source        string `json:"source"`
	Checksum      string `json:"checksum"`
	DeployKernel  string `json:"deployKernel"`
	DeployRamdisk string `json:"deployRamdisk"`
}

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// LibvirtURI is the identifier for the libvirtd connection.  It must be
	// reachable from the host where the installer is run.
	// +optional
	// Default is qemu:///system
	LibvirtURI string `json:"libvirtURI,omitempty"`

	// IronicURI is the identifier for the Ironic connection.  It must be
	// reachable from the host where the installer is run.
	// +optional
	IronicURI string `json:"ironicURI,omitempty"`

	// External bridge is used for external communication.
	// +optional
	ExternalBridge string `json:"externalBridge,omitempty"`

	// Provisioning bridge is used for provisioning nodes.
	// +optional
	ProvisioningBridge string `json:"provisioningBridge,omitempty"`

	// Hosts is the information needed to create the objects in Ironic.
	Hosts []*Host `json:"hosts"`

	// Images contains the information needed to provision a host
	Image Image `json:"image"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on bare metal for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// APIVIP is the VIP to use for internal API communication
	APIVIP string `json:"apiVIP"`

	// IngressVIP is the VIP to use for ingress traffic
	IngressVIP string `json:"ingressVIP"`
}
