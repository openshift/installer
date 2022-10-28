package powervs

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {

	// ServiceInstanceID is the ID of the Power IAAS instance created from the IBM Cloud Catalog
	ServiceInstanceID string `json:"serviceInstanceID"`

	// PowerVSResourceGroup is the resource group in which Power VS resources will be created.
	PowerVSResourceGroup string `json:"powervsResourceGroup"`

	// Region specifies the IBM Cloud colo region where the cluster will be created.
	Region string `json:"region"`

	// Zone specifies the IBM Cloud colo region where the cluster will be created.
	// At this time, only single-zone clusters are supported.
	Zone string `json:"zone"`

	// VPCRegion specifies the IBM Cloud region in which to create VPC resources.
	// Leave unset to allow installer to select the closest VPC region.
	//
	// +optional
	VPCRegion string `json:"vpcRegion,omitempty"`

	// VPCZone specifies the IBM Cloud zone in which to create VPC resources.
	// Leave unset to allow installer to select the closest VPC region.
	//
	// +optional
	VPCZone string `json:"vpcZone,omitempty"`

	// UserID is the login for the user's IBM Cloud account.
	UserID string `json:"userID"`

	// VPCName is the name of a pre-created VPC inside IBM Cloud.
	//
	// +optional
	VPCName string `json:"vpcName,omitempty"`

	// VPCSubnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	//
	// +optional
	VPCSubnets []string `json:"vpcSubnets,omitempty"`

	// PVSNetworkName specifies an existing network within the Power VS Service Instance.
	//
	// +optional
	PVSNetworkName string `json:"pvsNetworkName,omitempty"`

	// ClusterOSImage is a pre-created Power VS boot image that overrides the
	// default image for cluster nodes.
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Power VS for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// CloudConnctionName is the name of an existing Power VS Cloud connection.
	// If empty, one is created by the installer.
	// +optional
	CloudConnectionName string `json:"cloudConnectionName,omitempty"`
}
