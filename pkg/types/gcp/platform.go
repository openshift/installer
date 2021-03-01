package gcp

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// ProjectID is the the project that will be used for the cluster.
	ProjectID string `json:"projectID"`

	// Region specifies the GCP region where the cluster will be created.
	Region string `json:"region"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on GCP for machine pools which do not define their own
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

	// Licenses is no longer used and should be either empty or contain the nested virt license.
	// +optional
	Licenses []string `json:"licenses,omitempty"`
}
