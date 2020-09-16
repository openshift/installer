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

	// Licenses is a list of licenses to apply to the compute images
	// The value should a list of strings (https URLs only) representing the license keys.
	// When set, this will cause the installer to copy the image into user's project.
	// This option is incompatible with any mechanism that makes use of pre-built images
	// such as the current env OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE
	// +optional
	Licenses []string `json:"licenses,omitempty"`

	// Labels is a map of key-value pairs to apply to resources that are created.
	// Not all GCP resources support labels.
	// Additionally, the following requirements apply:
	// Keys and values cannot be longer than be 63 characters each.
	// Keys and values can only contain lowercase letters, numeric characters, underscores,
	// and hyphens. International characters are allowed.
	// Keys must start with a lowercase letter and cannot be empty.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}
