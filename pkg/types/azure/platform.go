package azure

import "strings"

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// Region specifies the Azure region where the cluster will be created.
	Region string `json:"region"`

	// BaseDomainResourceGroupName specifies the resource group where the azure DNS zone for the base domain is found
	BaseDomainResourceGroupName string `json:"baseDomainResourceGroupName,omitempty"`
	// DefaultMachinePlatform is the default configuration used when
	// installing on Azure for machine pools which do not define their own
	// platform configuration.
	// +optional

	// UserTags specifies additional tags for Azure resources created for the cluster.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}

//SetBaseDomain parses the baseDomainID and sets the related fields on azure.Platform
func (p *Platform) SetBaseDomain(baseDomainID string) error {
	parts := strings.Split(baseDomainID, "/")
	p.BaseDomainResourceGroupName = parts[4]
	return nil
}
