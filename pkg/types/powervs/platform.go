package powervs

// Platform stores all the global configuration that all machinesets
// use.
// Note: The subsequent mentions of future-TF support refer to work that is
// undergoing and should be available to test well in time for 4.10 feature-
// freeze. We do not plan to GA with these as required inputs.
type Platform struct {

	// ServiceInstanceID is the ID of the Power IAAS instance created from the IBM Cloud Catalog
	ServiceInstanceID string `json:"serviceInstance"`

	// PowerVSResourceGroup is the resource group for creating Power VS resources.
	PowerVSResourceGroup string `json:"powervsResourceGroup"`

	// Region specifies the IBM Cloud region where the cluster will be created.
	Region string `json:"region"`

	// Zone specifies the IBM Cloud colo region where the cluster will be created.
	// Required for multi-zone regions.
	Zone string `json:"zone"`

	// Zone in the region used to create VPC resources. Leave unset
	// to allow installer to randomly select a zone.
	//
	// +optional
	VPCZone string `json:"vpcRegion,omitempty"`

	// UserID is the login for the user's IBM Cloud account.
	UserID string `json:"userID"`

	// APIKey is the API key for the user's IBM Cloud account.
	//
	// +optional
	APIKey string `json:"APIKey,omitempty"`

	// SSHKeyName is the name of an SSH key stored in the Service Instance.
	SSHKeyName string `json:"SSHKeyName,omitempty"`

	// VPC is a VPC inside IBM Cloud. Needed in order to create VPC Load Balancers.
	//
	// @TODO: make this +optional when we have TF support
	VPC string `json:"vpc,omitempty"`

	// Subnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	// @TODO: Rename VPCSubnetID & make into string
	//
	// @TODO: make this +optional when we have TF support
	Subnets []string `json:"subnets,omitempty"`

	// PVSNetworkName specifies an existing network within the Power VS Service Instance.
	// @TODO: make this +optional when we have TF support
	PVSNetworkName string `json:"pvsNetworkName,omitempty"`

	// UserTags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	// BootstrapOSImage is a pre-created boot image to override the default
	// OS image for the bootstrap node.
	//
	// +optional
	BootstrapOSImage string `json:"bootstrapOSImage,omitempty"`

	// ClusterOSImage is a pre-created Power VS boot image that overrides the
	// default image for cluster nodes.
	//
	// @TODO: make this +optional, add omitempty when we have TF support
	ClusterOSImage string `json:"clusterOSImage"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Power VS for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
