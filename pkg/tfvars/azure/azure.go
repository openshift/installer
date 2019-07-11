package azure

import (
	"encoding/json"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/types/azure/defaults"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1alpha1"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	SubscriptionID string `json:"azure_subscription_id,omitempty"`
	ClientID       string `json:"azure_client_id,omitempty"`
	ClientSecret   string `json:"azure_client_secret,omitempty"`
	TenantID       string `json:"azure_tenant_id,omitempty"`
}

type config struct {
	Auth                        `json:",inline"`
	ExtraTags                   map[string]string `json:"azure_extra_tags,omitempty"`
	BootstrapInstanceType       string            `json:"azure_bootstrap_vm_type,omitempty"`
	MasterInstanceType          string            `json:"azure_master_vm_type,omitempty"`
	MasterAvailabilityZones     []string          `json:"azure_master_availability_zones"`
	VolumeSize                  int32             `json:"azure_master_root_volume_size,omitempty"`
	ImageURL                    string            `json:"azure_image_url,omitempty"`
	Region                      string            `json:"azure_region,omitempty"`
	BaseDomainResourceGroupName string            `json:"azure_base_domain_resource_group_name,omitempty"`
}

// TFVars generates Azure-specific Terraform variables launching the cluster.
func TFVars(auth Auth, baseDomainResourceGroupName string, imageURL string, masterConfigs []*azureprovider.AzureMachineProviderSpec) ([]byte, error) {
	masterConfig := masterConfigs[0]
	region := masterConfig.Location

	masterAvailabilityZones := make([]string, len(masterConfigs))
	for i, c := range masterConfigs {
		masterAvailabilityZones[i] = to.String(c.Zone)
	}

	cfg := &config{
		Auth:   auth,
		Region: region,
		BaseDomainResourceGroupName: baseDomainResourceGroupName,
		BootstrapInstanceType:       defaults.BootstrapInstanceType(region),
		MasterInstanceType:          masterConfig.VMSize,
		MasterAvailabilityZones:     masterAvailabilityZones,
		VolumeSize:                  masterConfig.OSDisk.DiskSizeGB,
		ImageURL:                    imageURL,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
