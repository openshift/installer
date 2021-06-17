package alibabacloud

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// ImageID is the ID of image that should be used to boot machines for the cluster.
	// If set, the image should belong to the same region as the cluster.
	//
	// +optional
	ImageID string `json:"imageID,omitempty"`

	// Region specifies the Alibaba Cloud region where the cluster will be created.
	Region string `json:"region"`

	// ResourceGroupID is the ID of an already existing resource group where the
	// cluster should be installed.
	ResourceGroupID string `json:"resourceGroupID"`

	// Tags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	Tags map[string]string `json:"tags"`

	// DefaultMachinePlatform is the default configuration used when installing
	// on Alibaba Cloud for machine pools which do not define their own platform
	// configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
