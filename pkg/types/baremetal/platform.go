package baremetal

import (
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
)

// BMC stores the information about a baremetal host's management controller.
type BMC struct {
	Username                       string `json:"username" validate:"required"`
	Password                       string `json:"password" validate:"required"`
	Address                        string `json:"address" validate:"required,uniqueField"`
	DisableCertificateVerification bool   `json:"disableCertificateVerification"`
}

// BootMode puts the server in legacy (BIOS), UEFI secure boot or UEFI mode for
// booting. Secure boot is only enabled during the final instance boot.
// The default is UEFI.
// +kubebuilder:validation:Enum="";UEFI;UEFISecureBoot;legacy
type BootMode string

// Allowed boot mode from metal3
const (
	UEFI           BootMode = "UEFI"
	UEFISecureBoot BootMode = "UEFISecureBoot"
	Legacy         BootMode = "legacy"
)

const (
	masterRole string = "master"
	workerRole string = "worker"
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
	NetworkConfig   *apiextv1.JSON   `json:"networkConfig,omitempty"`
}

// IsMaster checks if the current host is a master
func (h *Host) IsMaster() bool {
	return h.Role == masterRole
}

// IsWorker checks if the current host is a worker
func (h *Host) IsWorker() bool {
	return h.Role == workerRole
}

var sortIndex = map[string]int{masterRole: -1, workerRole: 0, "": 1}

// CompareByRole allows to compare two hosts by the Role
func (h *Host) CompareByRole(k *Host) bool {
	return sortIndex[h.Role] < sortIndex[k.Role]
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
	ClusterProvisioningIP string `json:"clusterProvisioningIP,omitempty"`

	// DeprecatedProvisioningHostIP is the deprecated version of clusterProvisioningIP. When the
	// baremetal platform was initially added to the installer, the JSON field for ClusterProvisioningIP
	// was incorrectly set to "provisioningHostIP."  This field is here to allow backwards-compatibility.
	// +optional
	DeprecatedProvisioningHostIP string `json:"provisioningHostIP,omitempty"`

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

	// ExternalMACAddress is used to allow setting a static unicast MAC
	// address for the bootstrap host on the external network. Consider
	// using the QEMU vendor prefix `52:54:00`. If left blank, libvirt will
	// generate one for you.
	// +optional
	ExternalMACAddress string `json:"externalMACAddress,omitempty"`

	// ProvisioningNetwork is used to indicate if we will have a provisioning network, and how it will be managed.
	// +kubebuilder:default=Managed
	// +optional
	ProvisioningNetwork ProvisioningNetwork `json:"provisioningNetwork,omitempty"`

	// Provisioning bridge is used for provisioning nodes, on the host that
	// will run the bootstrap VM.
	// +optional
	ProvisioningBridge string `json:"provisioningBridge,omitempty"`

	// ProvisioningMACAddress is used to allow setting a static unicast MAC
	// address for the bootstrap host on the provisioning network. Consider
	// using the QEMU vendor prefix `52:54:00`. If left blank, libvirt will
	// generate one for you.
	// +optional
	ProvisioningMACAddress string `json:"provisioningMACAddress,omitempty"`

	// ProvisioningNetworkInterface is the name of the network interface on a control plane
	// baremetal host that is connected to the provisioning network.
	// +optional
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

	// DeprecatedAPIVIP is the VIP to use for internal API communication
	// Deprecated: Use APIVIPs
	//
	// +kubebuilder:validation:Format=ip
	// +optional
	DeprecatedAPIVIP string `json:"apiVIP,omitempty"`

	// APIVIPs contains the VIP(s) to use for internal API communication. In
	// dual stack clusters it contains an IPv4 and IPv6 address, otherwise only
	// one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:Format=ip
	// +optional
	APIVIPs []string `json:"apiVIPs,omitempty"`

	// DeprecatedIngressVIP is the VIP to use for ingress traffic
	// Deprecated: Use IngressVIPs
	//
	// +kubebuilder:validation:Format=ip
	// +optional
	DeprecatedIngressVIP string `json:"ingressVIP,omitempty"`

	// IngressVIPs contains the VIP(s) to use for ingress traffic. In dual stack
	// clusters it contains an IPv4 and IPv6 address, otherwise only one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:Format=ip
	// +optional
	IngressVIPs []string `json:"ingressVIPs,omitempty"`

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

	// BootstrapExternalStaticIP is the static IP address of the bootstrap node.
	// This can be useful in environments without a DHCP server.
	// +kubebuilder:validation:Format=ip
	// +optional
	BootstrapExternalStaticIP string `json:"bootstrapExternalStaticIP,omitempty"`

	// BootstrapExternalStaticGateway is the static network gateway of the bootstrap node.
	// This can be useful in environments without a DHCP server.
	// +kubebuilder:validation:Format=ip
	// +optional
	BootstrapExternalStaticGateway string `json:"bootstrapExternalStaticGateway,omitempty"`

	// LoadBalancer defines how the load balancer used by the cluster is configured.
	// LoadBalancer is available in TechPreview.
	// +optional
	LoadBalancer *configv1.BareMetalPlatformLoadBalancer `json:"loadBalancer,omitempty"`

	// BootstrapExternalStaticDNS is the static network DNS of the bootstrap node.
	// This can be useful in environments without a DHCP server.
	// +kubebuilder:validation:Format=ip
	// +optional
	BootstrapExternalStaticDNS string `json:"bootstrapExternalStaticDNS,omitempty"`
}
