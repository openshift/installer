package alibabacloud

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// Region specifies the Alibaba Cloud region where the cluster will be created.
	Region string `json:"region"`

	// ResourceGroupID is the ID of an already existing resource group where the cluster should be installed.
	// If empty, the installer will create a new resource group for the cluster.
	// +optional
	ResourceGroupID string `json:"resourceGroupID,omitempty"`

	// Tags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	// +optional
	Tags map[string]string `json:"tags,omitempty"`

	// DefaultMachinePlatform is the default configuration used when installing
	// on Alibaba Cloud for machine pools which do not define their own platform
	// configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
