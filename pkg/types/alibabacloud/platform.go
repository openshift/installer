package alibabacloud

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// Region specifies the Alibaba Cloud region where the cluster will be created.
	Region string `json:"region"`

	// ResourceGroupName is the name of an already existing resource group where the
	// cluster should be installed.
	//
	// This resource group must be empty with no other resources when trying to
	// use it for creating a cluster.
	// +optional
	ResourceGroupName string `json:"resourceGroupName"`

	// DefaultMachinePlatform is the default configuration used when installing
	// on Alibaba Cloud for machine pools which do not define their own platform
	// configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}

// ClusterResourceGroupName returns the ID of the resource group for the cluster.
func (p *Platform) ClusterResourceGroupName(infraID string) string {
	if len(p.ResourceGroupName) > 0 {
		return p.ResourceGroupName
	}
	return infraID
}
