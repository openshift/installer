package ibmcloud

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// Region specifies the IBM Cloud region where the cluster will be
	// created.
	Region string `json:"region"`

	// ResourceGroupName is the name of an already existing resource group where the
	// cluster should be installed. If empty, a new resource group will be created
	// for the cluster.
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`

	// VPCName is the name of an already existing VPC where the cluster should be
	// installed.
	// +optional
	VPCName string `json:"vpcName,omitempty"`

	// ControlPlaneSubnets are the names of already existing subnets where the
	// cluster control plane nodes should be created.
	// +optional
	ControlPlaneSubnets []string `json:"controlPlaneSubnets,omitempty"`

	// ComputeSubnets are the names of already existing subnets where the cluster
	// compute nodes should be created.
	// +optional
	ComputeSubnets []string `json:"computeSubnets,omitempty"`

	// DefaultMachinePlatform is the default configuration used when installing
	// on IBM Cloud for machine pools which do not define their own platform
	// configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}

// ClusterResourceGroupName returns the name of the resource group for the cluster.
func (p *Platform) ClusterResourceGroupName(infraID string) string {
	if len(p.ResourceGroupName) > 0 {
		return p.ResourceGroupName
	}
	return infraID
}

// GetVPCName returns the user provided name of the VPC for the cluster.
func (p *Platform) GetVPCName() string {
	if len(p.VPCName) > 0 {
		return p.VPCName
	}
	return ""
}
