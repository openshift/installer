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

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// LibvirtURI is the identifier for the libvirtd connection.  It must be
	// reachable from the host where the installer is run.
	// +optional
	// Default is qemu:///system
	LibvirtURI string `json:"libvirtURI,omitempty"`

	// ClusterProvisioningIP is the IP on the dedicated provisioning network
	// where the baremetal-operator pod runs provisioning services,
	// and an http server to cache some downloaded content e.g RHCOS/IPA images
	// +optional
	ClusterProvisioningIP string `json:"provisioningHostIP,omitempty"`

	// BootstrapProvisioningIP is the IP used on the bootstrap VM to
	// bring up provisioning services that are used to create the
	// control-plane machines
	// +optional
	BootstrapProvisioningIP string `json:"bootstrapProvisioningIP,omitempty"`

	// External bridge is used for external communication.
	// +optional
	ExternalBridge string `json:"externalBridge,omitempty"`

	// Provisioning bridge is used for provisioning nodes.
	// +optional
	ProvisioningBridge string `json:"provisioningBridge,omitempty"`

	// Hosts is the information needed to create the objects in Ironic.
	Hosts []*Host `json:"hosts"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on bare metal for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// APIVIP is the VIP to use for internal API communication
	APIVIP string `json:"apiVIP"`

	// IngressVIP is the VIP to use for ingress traffic
	IngressVIP string `json:"ingressVIP"`

	// DNSVIP is the VIP to use for internal DNS communication
	DNSVIP string `json:"dnsVIP"`
}
