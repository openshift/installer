package azure

import (
	"encoding/json"
	"os"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure/defaults"
	azureprovider "sigs.k8s.io/cluster-api-provider-azure/pkg/apis/azureprovider/v1beta1"
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
	VolumeType                  string            `json:"azure_master_root_volume_type"`
	VolumeSize                  int32             `json:"azure_master_root_volume_size"`
	ImageURL                    string            `json:"azure_image_url,omitempty"`
	Region                      string            `json:"azure_region,omitempty"`
	BaseDomainResourceGroupName string            `json:"azure_base_domain_resource_group_name,omitempty"`
	NetworkResourceGroupName    string            `json:"azure_network_resource_group_name"`
	VirtualNetwork              string            `json:"azure_virtual_network"`
	ControlPlaneSubnet          string            `json:"azure_control_plane_subnet"`
	ComputeSubnet               string            `json:"azure_compute_subnet"`
	PreexistingNetwork          bool              `json:"azure_preexisting_network"`
	Private                     bool              `json:"azure_private"`
	EmulateSingleStackIPv6      bool              `json:"azure_emulate_single_stack_ipv6"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                        Auth
	BaseDomainResourceGroupName string
	MasterConfigs               []*azureprovider.AzureMachineProviderSpec
	WorkerConfigs               []*azureprovider.AzureMachineProviderSpec
	ImageURL                    string
	PreexistingNetwork          bool
	Publish                     types.PublishingStrategy
}

// TFVars generates Azure-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]
	workerConfig := sources.WorkerConfigs[0]

	region := masterConfig.Location

	masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		masterAvailabilityZones[i] = to.String(c.Zone)
	}

	var emulateSingleStackIPv6 bool
	if os.Getenv("OPENSHIFT_INSTALL_AZURE_EMULATE_SINGLESTACK_IPV6") == "true" {
		emulateSingleStackIPv6 = true
	}

	cfg := &config{
		Auth:                        sources.Auth,
		Region:                      region,
		BootstrapInstanceType:       defaults.BootstrapInstanceType(region),
		MasterInstanceType:          masterConfig.VMSize,
		MasterAvailabilityZones:     masterAvailabilityZones,
		VolumeType:                  masterConfig.OSDisk.ManagedDisk.StorageAccountType,
		VolumeSize:                  masterConfig.OSDisk.DiskSizeGB,
		ImageURL:                    sources.ImageURL,
		Private:                     sources.Publish == types.InternalPublishingStrategy,
		BaseDomainResourceGroupName: sources.BaseDomainResourceGroupName,
		NetworkResourceGroupName:    masterConfig.NetworkResourceGroup,
		VirtualNetwork:              masterConfig.Vnet,
		ControlPlaneSubnet:          masterConfig.Subnet,
		ComputeSubnet:               workerConfig.Subnet,
		PreexistingNetwork:          sources.PreexistingNetwork,
		EmulateSingleStackIPv6:      emulateSingleStackIPv6,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
