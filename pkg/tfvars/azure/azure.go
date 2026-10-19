package azure

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

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
	UseMSI                    bool   `json:"azure_use_msi,omitempty"`
}

// OSImage is the marketplace image to be used for RHCOS.
type OSImage struct {
	Publisher string `json:"azure_marketplace_image_publisher,omitempty"`
	Offer     string `json:"azure_marketplace_image_offer,omitempty"`
	SKU       string `json:"azure_marketplace_image_sku,omitempty"`
	Version   string `json:"azure_marketplace_image_version,omitempty"`
}

type config struct {
	Auth                                    `json:",inline"`
	Environment                             string            `json:"azure_environment"`
	ARMEndpoint                             string            `json:"azure_arm_endpoint"`
	ExtraTags                               map[string]string `json:"azure_extra_tags,omitempty"`
	MasterInstanceType                      string            `json:"azure_master_vm_type,omitempty"`
	MasterAvailabilityZones                 []string          `json:"azure_master_availability_zones"`
	MasterEncryptionAtHostEnabled           bool              `json:"azure_master_encryption_at_host_enabled"`
	MasterDiskEncryptionSetID               string            `json:"azure_master_disk_encryption_set_id,omitempty"`
	ControlPlaneUltraSSDEnabled             bool              `json:"azure_control_plane_ultra_ssd_enabled"`
	VolumeType                              string            `json:"azure_master_root_volume_type"`
	VolumeSize                              int32             `json:"azure_master_root_volume_size"`
	ImageURL                                string            `json:"azure_image_url,omitempty"`
	ImageRelease                            string            `json:"azure_image_release,omitempty"`
	Region                                  string            `json:"azure_region,omitempty"`
	BaseDomainResourceGroupName             string            `json:"azure_base_domain_resource_group_name,omitempty"`
	ResourceGroupName                       string            `json:"azure_resource_group_name"`
	NetworkResourceGroupName                string            `json:"azure_network_resource_group_name"`
	VirtualNetwork                          string            `json:"azure_virtual_network"`
	ControlPlaneSubnet                      string            `json:"azure_control_plane_subnet"`
	ComputeSubnet                           string            `json:"azure_compute_subnet"`
	PreexistingNetwork                      bool              `json:"azure_preexisting_network"`
	Private                                 bool              `json:"azure_private"`
	OutboundType                            string            `json:"azure_outbound_routing_type"`
	BootstrapIgnitionStub                   string            `json:"azure_bootstrap_ignition_stub"`
	BootstrapIgnitionURLPlaceholder         string            `json:"azure_bootstrap_ignition_url_placeholder"`
	HyperVGeneration                        string            `json:"azure_hypervgeneration_version"`
	VMNetworkingType                        bool              `json:"azure_control_plane_vm_networking_type"`
	RandomStringPrefix                      string            `json:"random_storage_account_suffix"`
	VMArchitecture                          string            `json:"azure_vm_architecture"`
	UseMarketplaceImage                     bool              `json:"azure_use_marketplace_image"`
	MarketplaceImageHasPlan                 bool              `json:"azure_marketplace_image_has_plan"`
	OSImage                                 `json:",inline"`
	SecurityEncryptionType                  string            `json:"azure_master_security_encryption_type,omitempty"`
	SecureVirtualMachineDiskEncryptionSetID string            `json:"azure_master_secure_vm_disk_encryption_set_id,omitempty"`
	SecureBoot                              string            `json:"azure_master_secure_boot,omitempty"`
	VirtualizedTrustedPlatformModule        string            `json:"azure_master_virtualized_trusted_platform_module,omitempty"`
	KeyVaultResourceGroup                   string            `json:"azure_keyvault_resource_group,omitempty"`
	KeyVaultName                            string            `json:"azure_keyvault_name,omitempty"`
	KeyVaultKeyName                         string            `json:"azure_keyvault_key_name,omitempty"`
	UserAssignedIdentity                    string            `json:"azure_user_assigned_identity_key,omitempty"`
	ResourceGroupMetadataTags               map[string]string `json:"azure_resource_group_metadata_tags"`
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
	ImageRelease                    string
	PreexistingNetwork              bool
	Publish                         types.PublishingStrategy
	OutboundType                    azure.OutboundType
	BootstrapIgnStub                string
	BootstrapIgnitionURLPlaceholder string
	HyperVGeneration                string
	VMArchitecture                  types.Architecture
	InfrastructureName              string
	KeyVault                        azure.KeyVault
	UserAssignedIdentityKey         string
	LBPrivate                       bool
}

// TFVars generates Azure-specific Terraform variables launching the cluster.
func TFVars(sources TFVarsSources) ([]byte, error) {
	masterConfig := sources.MasterConfigs[0]
	workerConfig := sources.WorkerConfigs[0]

	region := masterConfig.Location

	masterAvailabilityZones := make([]string, len(sources.MasterConfigs))
	for i, c := range sources.MasterConfigs {
		masterAvailabilityZones[i] = c.Zone
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

	var secureBoot string
	var virtualizedTrustedPlatformModule string
	if masterConfig.SecurityProfile != nil && masterConfig.SecurityProfile.Settings.SecurityType != "" {
		switch masterConfig.SecurityProfile.Settings.SecurityType {
		case machineapi.SecurityTypesConfidentialVM:
			secureBoot = string(masterConfig.SecurityProfile.Settings.ConfidentialVM.UEFISettings.SecureBoot)
			virtualizedTrustedPlatformModule = string(masterConfig.SecurityProfile.Settings.ConfidentialVM.UEFISettings.VirtualizedTrustedPlatformModule)
		case machineapi.SecurityTypesTrustedLaunch:
			secureBoot = string(masterConfig.SecurityProfile.Settings.TrustedLaunch.UEFISettings.SecureBoot)
			virtualizedTrustedPlatformModule = string(masterConfig.SecurityProfile.Settings.TrustedLaunch.UEFISettings.VirtualizedTrustedPlatformModule)
		}
	}

	vmarch := "x64"
	if sources.VMArchitecture == types.ArchitectureARM64 {
		vmarch = "Arm64"
	}

	tags := make(map[string]string, len(masterConfig.Tags)+1)
	// add OCP default tag
	tags[fmt.Sprintf("kubernetes.io_cluster.%s", sources.InfrastructureName)] = "owned"
	for k, v := range masterConfig.Tags {
		tags[k] = v
	}

	osImage := OSImage{
		Publisher: masterConfig.Image.Publisher,
		Offer:     masterConfig.Image.Offer,
		SKU:       masterConfig.Image.SKU,
		Version:   masterConfig.Image.Version,
	}

	// Metadata tags to be added to the resource group for the cluster destroy
	metadataTags := map[string]string{}
	metadataTags[azure.TagMetadataRegion] = region
	if len(sources.BaseDomainResourceGroupName) > 0 {
		metadataTags[azure.TagMetadataBaseDomainRG] = sources.BaseDomainResourceGroupName
	}
	if len(masterConfig.NetworkResourceGroup) > 0 {
		metadataTags[azure.TagMetadataNetworkRG] = masterConfig.NetworkResourceGroup
	}

	cfg := &config{
		Auth:                                    sources.Auth,
		Environment:                             environment,
		ARMEndpoint:                             sources.ARMEndpoint,
		Region:                                  region,
		MasterInstanceType:                      masterConfig.VMSize,
		MasterAvailabilityZones:                 masterAvailabilityZones,
		MasterEncryptionAtHostEnabled:           masterEncryptionAtHostEnabled,
		MasterDiskEncryptionSetID:               masterDiskEncryptionSetID,
		ControlPlaneUltraSSDEnabled:             masterConfig.UltraSSDCapability == machineapi.AzureUltraSSDCapabilityEnabled,
		VolumeType:                              masterConfig.OSDisk.ManagedDisk.StorageAccountType,
		VolumeSize:                              masterConfig.OSDisk.DiskSizeGB,
		ImageURL:                                sources.ImageURL,
		ImageRelease:                            sources.ImageRelease,
		Private:                                 sources.Publish == types.InternalPublishingStrategy || sources.LBPrivate,
		OutboundType:                            string(sources.OutboundType),
		ResourceGroupName:                       sources.ResourceGroupName,
		BaseDomainResourceGroupName:             sources.BaseDomainResourceGroupName,
		NetworkResourceGroupName:                masterConfig.NetworkResourceGroup,
		VirtualNetwork:                          masterConfig.Vnet,
		ControlPlaneSubnet:                      masterConfig.Subnet,
		ComputeSubnet:                           workerConfig.Subnet,
		PreexistingNetwork:                      sources.PreexistingNetwork,
		BootstrapIgnitionStub:                   sources.BootstrapIgnStub,
		BootstrapIgnitionURLPlaceholder:         sources.BootstrapIgnitionURLPlaceholder,
		HyperVGeneration:                        sources.HyperVGeneration,
		VMNetworkingType:                        masterConfig.AcceleratedNetworking,
		RandomStringPrefix:                      randomStringPrefixFunction(),
		VMArchitecture:                          vmarch,
		ExtraTags:                               tags,
		UseMarketplaceImage:                     osImage.Publisher != "",
		MarketplaceImageHasPlan:                 masterConfig.Image.Type != machineapi.AzureImageTypeMarketplaceNoPlan,
		OSImage:                                 osImage,
		SecurityEncryptionType:                  string(masterConfig.OSDisk.ManagedDisk.SecurityProfile.SecurityEncryptionType),
		SecureVirtualMachineDiskEncryptionSetID: masterConfig.OSDisk.ManagedDisk.SecurityProfile.DiskEncryptionSet.ID,
		SecureBoot:                              secureBoot,
		VirtualizedTrustedPlatformModule:        virtualizedTrustedPlatformModule,
		KeyVaultResourceGroup:                   sources.KeyVault.ResourceGroup,
		KeyVaultName:                            sources.KeyVault.Name,
		KeyVaultKeyName:                         sources.KeyVault.KeyName,
		UserAssignedIdentity:                    sources.UserAssignedIdentityKey,
		ResourceGroupMetadataTags:               metadataTags,
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

func randomStringPrefixFunction() string {
	length := 5
	rand.Seed(time.Now().UnixNano())
	suffix := make([]rune, length)
	for i := 0; i < length; i++ {
		suffix[i] = 97 + rand.Int31n(26)
	}
	return string(suffix)
}
