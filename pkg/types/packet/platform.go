package packet

type Platform struct {
	// FacilityCode represents the Packet region and datacenter where your devices will be provisioned (https://www.packet.com/developers/docs/getting-started/facilities/)
	FacilityCode string `json:"facility_code,omitempty"`

	// ProjectID represents the Packet project used for logical grouping and invoicing (https://www.packet.com/developers/docs/API/getting-started/)
	ProjectID string `json:"project_id,omitempty"`

	// APIVIP is the static IP on the nodes subnet that the api port for
	// openshift will be assigned
	// Default: will be set to the 5 on the first entry in the machineNetwork
	// CIDR
	// +optional
	// +kubebuilder:validation:Format=ip
	APIVIP string `json:"apivip,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Packet for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network specifies an existing VPC where the cluster should be created
	// rather than provisioning a new one.
	// +optional
	Network string `json:"network,omitempty"`

	// ControlPlaneSubnet is an existing subnet where the control plane will be deployed.
	// The value should be the name of the subnet.
	// +optional
	ControlPlaneSubnet string `json:"controlPlaneSubnet,omitempty"`

	// ComputeSubnet is an existing subnet where the compute nodes will be deployed.
	// The value should be the name of the subnet.
	// +optional
	ComputeSubnet string `json:"computeSubnet,omitempty"`

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
