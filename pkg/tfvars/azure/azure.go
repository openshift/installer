package azure

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/types/azure/defaults"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1alpha1"
)

type config struct {
	ExtraTags                   map[string]string `json:"azure_extra_tags,omitempty"`
	BootstrapInstanceType       string            `json:"azure_bootstrap_vm_type,omitempty"`
	MasterInstanceType          string            `json:"azure_master_vm_type,omitempty"`
	VolumeSize                  int32             `json:"azure_master_root_volume_size,omitempty"`
	VMImageID                   string            `json:"azure_image_id,omitempty"`
	Region                      string            `json:"azure_region,omitempty"`
	BaseDomainResourceGroupName string            `json:"azure_base_domain_resource_group_name,omitempty"`
}

// TFVars generates Azure-specific Terraform variables launching the cluster.
func TFVars(baseDomainResourceGroupName string, masterConfigs []*azureprovider.AzureMachineProviderSpec) ([]byte, error) {
	masterConfig := masterConfigs[0]
	region := masterConfig.Location
	cfg := &config{
		Region: region,
		BaseDomainResourceGroupName: baseDomainResourceGroupName,
		BootstrapInstanceType:       defaults.InstanceClass(region),
		MasterInstanceType:          masterConfig.VMSize,
		VolumeSize:                  masterConfig.OSDisk.DiskSizeGB,
		VMImageID:                   masterConfig.Image.ResourceID,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
