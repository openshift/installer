package alibabacloud

import "fmt"

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// Region specifies the Alibaba Cloud region where the cluster will be created.
	Region string `json:"region"`

	// ResourceGroupID is the ID of an already existing resource group where the cluster should be installed.
	// If empty, the installer will create a new resource group for the cluster.
	// +optional
	ResourceGroupID string `json:"resourceGroupID,omitempty"`

	// VpcID is the ID of an already existing VPC where the cluster should be installed.
	// If empty, the installer will create a new VPC for the cluster.
	// +optional
	VpcID string `json:"vpcID,omitempty"`

	// VSwitchIDs is the ID list of already existing VSwitches where cluster resources will be created.
	// The existing VSwitches can only be used when also using existing VPC.
	// If empty, the installer will create new VSwitches for the cluster.
	// +optional
	VSwitchIDs []string `json:"vswitchIDs,omitempty"`

	// PrivateZoneID is the ID of an existing private zone into which to add DNS
	// records for the cluster's internal API. An existing private zone can
	// only be used when also using existing VPC. The private zone must be
	// associated with the VPC containing the subnets.
	// Leave the private zone unset to have the installer create the private zone
	// on your behalf.
	// +optional
	PrivateZoneID string `json:"privateZoneID,omitempty"`

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

// ClusterResourceGroupName returns the name of the resource group for the cluster.
func (p *Platform) ClusterResourceGroupName(clusterID string) string {
	return fmt.Sprintf("%s-rg", clusterID)
}
