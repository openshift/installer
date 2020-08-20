package baremetal

import (
	"github.com/openshift/installer/pkg/ipnet"
)

// BMC stores the information about a baremetal host's management controller.
type BMC struct {
	Username                       string `json:"username" validate:"required"`
	Password                       string `json:"password" validate:"required"`
	Address                        string `json:"address" validate:"required,uniqueField"`
	DisableCertificateVerification bool   `json:"disableCertificateVerification"`
}

// BootMode puts the server in legacy (BIOS) or UEFI mode for
// booting. The default is UEFI.
// +kubebuilder:validation:Enum=UEFI;legacy
type BootMode string

// Allowed boot mode from metal3
const (
	UEFI   BootMode = "UEFI"
	Legacy BootMode = "legacy"
)

// Host stores all the configuration data for a baremetal host.
type Host struct {
	Name            string           `json:"name,omitempty" validate:"required,uniqueField"`
	BMC             BMC              `json:"bmc"`
	Role            string           `json:"role"`
	BootMACAddress  string           `json:"bootMACAddress" validate:"required,uniqueField"`
	HardwareProfile string           `json:"hardwareProfile"`
	RootDeviceHints *RootDeviceHints `json:"rootDeviceHints,omitempty"`
	BootMode        BootMode         `json:"bootMode,omitempty"`
}

// ProvisioningNetwork determines how we will use the provisioning network.
// +kubebuilder:validation:Enum="";Managed;Unmanaged;Disabled
type ProvisioningNetwork string

const (
	// ManagedProvisioningNetwork indicates we should fully manage the provisioning network, including DHCP
	// services required for PXE-based provisioning.
	ManagedProvisioningNetwork ProvisioningNetwork = "Managed"

	// UnmanagedProvisioningNetwork indicates responsibility for managing the provisioning network is left to the
	// user. No DHCP server will be configured, however TFTP remains enabled if a user wants to use PXE-based provisioning.
	// However, they will need to configure external DHCP correctly with next-server definitions set to the relevant
	// provisioning IP's.
	UnmanagedProvisioningNetwork ProvisioningNetwork = "Unmanaged"

	// DisabledProvisioningNetwork indicates that no provisioning network will be used. Provisioning capabilities
	// will be limited to virtual media-based deployments only, and neither DHCP nor TFTP will be operated by the
	// cluster.
	DisabledProvisioningNetwork ProvisioningNetwork = "Disabled"
)

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// LibvirtURI is the identifier for the libvirtd connection.  It must be
	// reachable from the host where the installer is run.
	// Default is qemu:///system
	//
	// +kubebuilder:default="qemu:///system"
	// +optional
	LibvirtURI string `json:"libvirtURI,omitempty"`

	// ClusterProvisioningIP is the IP on the dedicated provisioning network
	// where the baremetal-operator pod runs provisioning services,
	// and an http server to cache some downloaded content e.g RHCOS/IPA images
	// +optional
	ClusterProvisioningIP string `json:"provisioningHostIP,omitempty"`

	// BootstrapProvisioningIP is the IP used on the bootstrap VM to
	// bring up provisioning services that are used to create the
	// control-plane machines
	//
	// +kubebuilder:validation:Format=ip
	// +optional
	BootstrapProvisioningIP string `json:"bootstrapProvisioningIP,omitempty"`

	// External bridge is used for external communication.
	// +optional
	ExternalBridge string `json:"externalBridge,omitempty"`

	// ProvisioningNetwork is used to indicate if we will have a provisioning network, and how it will be managed.
	// +kubebuilder:default=Managed
	// +optional
	ProvisioningNetwork ProvisioningNetwork `json:"provisioningNetwork,omitempty"`

	// Provisioning bridge is used for provisioning nodes, on the host that
	// will run the bootstrap VM.
	// +optional
	ProvisioningBridge string `json:"provisioningBridge,omitempty"`

	// ProvisioningNetworkInterface is the name of the network interface on a control plane
	// baremetal host that is connected to the provisioning network.
	ProvisioningNetworkInterface string `json:"provisioningNetworkInterface"`

	// ProvisioningNetworkCIDR defines the network to use for provisioning.
	// +optional
	ProvisioningNetworkCIDR *ipnet.IPNet `json:"provisioningNetworkCIDR,omitempty"`

	// DeprecatedProvisioningDHCPExternal indicates that DHCP is provided by an external service. This parameter is
	// replaced by ProvisioningNetwork being set to "Unmanaged".
	// +optional
	DeprecatedProvisioningDHCPExternal bool `json:"provisioningDHCPExternal,omitempty"`

	// ProvisioningDHCPRange is used to provide DHCP services to hosts
	// for provisioning.
	// +optional
	ProvisioningDHCPRange string `json:"provisioningDHCPRange,omitempty"`

	// Hosts is the information needed to create the objects in Ironic.
	Hosts []*Host `json:"hosts"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on bare metal for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// APIVIP is the VIP to use for internal API communication
	//
	// +kubebuilder:validation:Format=ip
	APIVIP string `json:"apiVIP"`

	// IngressVIP is the VIP to use for ingress traffic
	//
	// +kubebuilder:validation:Format=ip
	IngressVIP string `json:"ingressVIP"`

	// BootstrapOSImage is a URL to override the default OS image
	// for the bootstrap node. The URL must contain a sha256 hash of the image
	// e.g https://mirror.example.com/images/qemu.qcow2.gz?sha256=a07bd...
	//
	// +optional
	BootstrapOSImage string `json:"bootstrapOSImage,omitempty" validate:"omitempty,osimageuri,urlexist"`

	// ClusterOSImage is a URL to override the default OS image
	// for cluster nodes. The URL must contain a sha256 hash of the image
	// e.g https://mirror.example.com/images/metal.qcow2.gz?sha256=3b5a8...
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty" validate:"omitempty,osimageuri,urlexist"`
}
