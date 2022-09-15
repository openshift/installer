package azure

import (
	"encoding/json"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// Auth is the collection of credentials that will be used by terrform.
type Auth struct {
	SubscriptionID            string `json:"azure_subscription_id,omitempty"`
	ClientID                  string `json:"azure_client_id,omitempty"`
	ClientSecret              string `json:"azure_client_secret,omitempty"`
	TenantID                  string `json:"azure_tenant_id,omitempty"`
	ClientCertificatePath     string `json:"azure_certificate_path,omitempty"`
	ClientCertificatePassword string `json:"azure_certificate_password,omitempty"`
}

type config struct {
	Auth                            `json:",inline"`
	Environment                     string            `json:"azure_environment"`
	ARMEndpoint                     string            `json:"azure_arm_endpoint"`
	ExtraTags                       map[string]string `json:"azure_extra_tags,omitempty"`
	MasterInstanceType              string            `json:"azure_master_vm_type,omitempty"`
	MasterAvailabilityZones         []string          `json:"azure_master_availability_zones"`
	MasterEncryptionAtHostEnabled   bool              `json:"azure_master_encryption_at_host_enabled"`
	MasterDiskEncryptionSetID       string            `json:"azure_master_disk_encryption_set_id,omitempty"`
	ControlPlaneUltraSSDEnabled     bool              `json:"azure_control_plane_ultra_ssd_enabled"`
	VolumeType                      string            `json:"azure_master_root_volume_type"`
	VolumeSize                      int32             `json:"azure_master_root_volume_size"`
	ImageURL                        string            `json:"azure_image_url,omitempty"`
	Region                          string            `json:"azure_region,omitempty"`
	BaseDomainResourceGroupName     string            `json:"azure_base_domain_resource_group_name,omitempty"`
	ResourceGroupName               string            `json:"azure_resource_group_name"`
	NetworkResourceGroupName        string            `json:"azure_network_resource_group_name"`
	VirtualNetwork                  string            `json:"azure_virtual_network"`
	ControlPlaneSubnet              string            `json:"azure_control_plane_subnet"`
	ComputeSubnet                   string            `json:"azure_compute_subnet"`
	PreexistingNetwork              bool              `json:"azure_preexisting_network"`
	Private                         bool              `json:"azure_private"`
	OutboundUDR                     bool              `json:"azure_outbound_user_defined_routing"`
	BootstrapIgnitionStub           string            `json:"azure_bootstrap_ignition_stub"`
	BootstrapIgnitionURLPlaceholder string            `json:"azure_bootstrap_ignition_url_placeholder"`
	HyperVGeneration                string            `json:"azure_hypervgeneration_version"`
	VMNetworkingType                bool              `json:"azure_control_plane_vm_networking_type"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	Auth                            Auth
	CloudName                       azure.CloudEnvironment
	ARMEndpoint                     string
	ResourceGroupName               string
	BaseDomainResourceGroupName     string
	MasterConfigs                   []*machineapi.AzureMachineProviderSpec
	WorkerConfigs                   []*machineapi.AzureMachineProviderSpec
	ImageURL                        string
	PreexistingNetwork              bool
	Publish                         types.PublishingStrategy
	OutboundType                    azure.OutboundType
	BootstrapIgnStub                string
	BootstrapIgnitionURLPlaceholder string
	HyperVGeneration                string
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

	environment, err := environment(sources.CloudName)
	if err != nil {
		return nil, errors.Wrap(err, "could not determine Azure environment to use for Terraform")
	}

	masterEncryptionAtHostEnabled := masterConfig.SecurityProfile != nil &&
		(*masterConfig.SecurityProfile).EncryptionAtHost != nil &&
		*masterConfig.SecurityProfile.EncryptionAtHost

	var masterDiskEncryptionSetID string
	if masterConfig.OSDisk.ManagedDisk.DiskEncryptionSet != nil {
		masterDiskEncryptionSetID = masterConfig.OSDisk.ManagedDisk.DiskEncryptionSet.ID
	}

	cfg := &config{
		Auth:                            sources.Auth,
		Environment:                     environment,
		ARMEndpoint:                     sources.ARMEndpoint,
		Region:                          region,
		MasterInstanceType:              masterConfig.VMSize,
		MasterAvailabilityZones:         masterAvailabilityZones,
		MasterEncryptionAtHostEnabled:   masterEncryptionAtHostEnabled,
		MasterDiskEncryptionSetID:       masterDiskEncryptionSetID,
		ControlPlaneUltraSSDEnabled:     masterConfig.UltraSSDCapability == machineapi.AzureUltraSSDCapabilityEnabled,
		VolumeType:                      masterConfig.OSDisk.ManagedDisk.StorageAccountType,
		VolumeSize:                      masterConfig.OSDisk.DiskSizeGB,
		ImageURL:                        sources.ImageURL,
		Private:                         sources.Publish == types.InternalPublishingStrategy,
		OutboundUDR:                     sources.OutboundType == azure.UserDefinedRoutingOutboundType,
		ResourceGroupName:               sources.ResourceGroupName,
		BaseDomainResourceGroupName:     sources.BaseDomainResourceGroupName,
		NetworkResourceGroupName:        masterConfig.NetworkResourceGroup,
		VirtualNetwork:                  masterConfig.Vnet,
		ControlPlaneSubnet:              masterConfig.Subnet,
		ComputeSubnet:                   workerConfig.Subnet,
		PreexistingNetwork:              sources.PreexistingNetwork,
		BootstrapIgnitionStub:           sources.BootstrapIgnStub,
		BootstrapIgnitionURLPlaceholder: sources.BootstrapIgnitionURLPlaceholder,
		HyperVGeneration:                sources.HyperVGeneration,
		VMNetworkingType:                masterConfig.AcceleratedNetworking,
	}

	return json.MarshalIndent(cfg, "", "  ")
}

// environment returns the Azure environment to pass to Terraform
func environment(cloudName azure.CloudEnvironment) (string, error) {
	switch cloudName {
	case azure.PublicCloud:
		return "public", nil
	case azure.USGovernmentCloud:
		return "usgovernment", nil
	case azure.ChinaCloud:
		return "china", nil
	case azure.GermanCloud:
		return "german", nil
	case azure.StackCloud:
		// unused since stack uses its own provider
		return "", nil
	default:
		return "", errors.Errorf("unsupported cloud name %q", cloudName)
	}
}
