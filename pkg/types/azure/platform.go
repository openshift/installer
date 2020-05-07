package azure

import "strings"

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// Region specifies the Azure region where the cluster will be created.
	Region string `json:"region"`

	// BaseDomainResourceGroupName specifies the resource group where the Azure DNS zone for the base domain is found.
	BaseDomainResourceGroupName string `json:"baseDomainResourceGroupName,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Azure for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// NetworkResourceGroupName specifies the network resource group that contains an existing VNet
	//
	// +optional
	NetworkResourceGroupName string `json:"networkResourceGroupName,omitempty"`

	// VirtualNetwork specifies the name of an existing VNet for the installer to use
	//
	// +optional
	VirtualNetwork string `json:"virtualNetwork,omitempty"`

	// ControlPlaneSubnet specifies an existing subnet for use by the control plane nodes
	//
	// +optional
	ControlPlaneSubnet string `json:"controlPlaneSubnet,omitempty"`

	// ComputeSubnet specifies an existing subnet for use by compute nodes
	//
	// +optional
	ComputeSubnet string `json:"computeSubnet,omitempty"`
}

//SetBaseDomain parses the baseDomainID and sets the related fields on azure.Platform
func (p *Platform) SetBaseDomain(baseDomainID string) error {
	parts := strings.Split(baseDomainID, "/")
	p.BaseDomainResourceGroupName = parts[4]
	return nil
}
