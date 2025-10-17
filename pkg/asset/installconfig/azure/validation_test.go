package azure

import (
	"fmt"
	"net"
	"testing"

	azres "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/resources"
	azsubs "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/resources/mgmt/subscriptions"
	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	azenc "github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig/azure/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

type editFunctions []func(ic *types.InstallConfig)

var (
	validVirtualNetwork            = "valid-virtual-network"
	validNetworkResourceGroup      = "valid-network-resource-group"
	validRegion                    = "centralus"
	validRegionsList               = []string{"centralus", "westus", "australiacentral2"}
	resourcesCapableRegionsList    = []string{"centralus", "westus"}
	validComputeSubnet             = "valid-compute-subnet"
	validControlPlaneSubnet        = "valid-controlplane-subnet"
	validCIDR                      = "10.0.0.0/16"
	validComputeSubnetCIDR         = "10.0.0.0/24"
	validControlPlaneSubnetCIDR    = "10.0.32.0/24"
	validResourceGroupNamespace    = "Microsoft.Resources"
	validResourceGroupResourceType = "resourceGroups"
	validResourceSkuRegions        = "southeastasia"

	vmCapabilities = map[string]map[string]string{
		"Standard_D8s_v3":    {"vCPUsAvailable": "4", "MemoryGB": "16", "PremiumIO": "True", "HyperVGenerations": "V1,V2", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "x64"},
		"Standard_D4s_v3":    {"vCPUsAvailable": "4", "MemoryGB": "32", "PremiumIO": "True", "HyperVGenerations": "V1", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "x64"},
		"Standard_A1_v2":     {"vCPUsAvailable": "1", "MemoryGB": "2", "PremiumIO": "True", "HyperVGenerations": "V1,V2", "AcceleratedNetworkingEnabled": "False", "CpuArchitectureType": "x64"},
		"Standard_D2_v4":     {"vCPUsAvailable": "2", "MemoryGB": "8", "PremiumIO": "True", "HyperVGenerations": "V1,V2", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "x64"},
		"Standard_D4_v4":     {"vCPUsAvailable": "4", "MemoryGB": "16", "PremiumIO": "False", "HyperVGenerations": "V1,V2", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "x64"},
		"Standard_D2s_v3":    {"vCPUsAvailable": "4", "MemoryGB": "16", "PremiumIO": "True", "HyperVGenerations": "V1,V2", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "x64"},
		"Standard_Dc4_v4":    {"vCPUsAvailable": "4", "MemoryGB": "16", "PremiumIO": "True", "HyperVGenerations": "V2", "CpuArchitectureType": "x64"},
		"Standard_B4ms":      {"vCPUsAvailable": "4", "MemoryGB": "16", "PremiumIO": "True", "HyperVGenerations": "V1,V2", "AcceleratedNetworkingEnabled": "False", "CpuArchitectureType": "x64"},
		"Standard_D8ps_v5":   {"vCPUsAvailable": "8", "MemoryGB": "32", "PremiumIO": "True", "HyperVGenerations": "V2", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "Arm64", "TrustedLaunchDisabled": "True"},
		"Standard_D4ps_v5":   {"vCPUsAvailable": "4", "MemoryGB": "16", "PremiumIO": "True", "HyperVGenerations": "V2", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "Arm64", "TrustedLaunchDisabled": "True"},
		"Standard_DC8ads_v5": {"vCPUsAvailable": "8", "MemoryGB": "32", "PremiumIO": "True", "HyperVGenerations": "V2", "AcceleratedNetworkingEnabled": "False", "CpuArchitectureType": "x64", "ConfidentialComputingType": "SNP"},
		"Standard_DC8eds_v5": {"vCPUsAvailable": "8", "MemoryGB": "32", "PremiumIO": "True", "HyperVGenerations": "V2", "AcceleratedNetworkingEnabled": "False", "CpuArchitectureType": "x64", "ConfidentialComputingType": "TDX"},
		"Standard_DC8s_v3":   {"vCPUsAvailable": "8", "MemoryGB": "32", "PremiumIO": "True", "HyperVGenerations": "V2", "AcceleratedNetworkingEnabled": "True", "CpuArchitectureType": "x64", "ConfidentialComputingType": "SGX"},
	}

	instanceTypeSku = func() []*azenc.ResourceSku {
		instances := make([]*azenc.ResourceSku, 0, len(vmCapabilities))
		for typeName, capsMap := range vmCapabilities {
			capabilities := make([]azenc.ResourceSkuCapabilities, 0, len(capsMap))
			for name, value := range capsMap {
				capabilities = append(capabilities, azenc.ResourceSkuCapabilities{
					Name: to.StringPtr(name), Value: to.StringPtr(value),
				})
			}
			instances = append(instances, &azenc.ResourceSku{
				Name: to.StringPtr(typeName), Capabilities: &capabilities,
			})
		}
		return instances
	}()

	validInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D4s_v3"
		ic.ControlPlane.Platform.Azure.InstanceType = "Standard_D8s_v3"
		ic.Compute[0].Platform.Azure.InstanceType = "Standard_D4s_v3"
	}

	validArm64InstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D4ps_v5"
		ic.ControlPlane.Architecture = types.ArchitectureARM64
		ic.ControlPlane.Platform.Azure.InstanceType = "Standard_D8ps_v5"
		ic.Compute[0].Architecture = types.ArchitectureARM64
		ic.Compute[0].Platform.Azure.InstanceType = "Standard_D4ps_v5"
	}

	invalidArchInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D4ps_v5"
		ic.ControlPlane.Platform.Azure.InstanceType = "Standard_D8ps_v5"
		ic.Compute[0].Platform.Azure.InstanceType = "Standard_D4ps_v5"
	}

	invalidateDefaultInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_A1_v2"
	}

	invalidateControlPlaneInstanceTypes = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.InstanceType = "Standard_A1_v2"
	}

	invalidateComputeInstanceTypes = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.InstanceType = "Standard_A1_v2"
	}

	undefinedDefaultInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Dne_D2_v4"
	}

	ultraSSDAvailableInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D8s_v3"
	}

	validVMNetworkingInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D8s_v3"
	}

	invalidVMNetworkingIstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_B4ms"
	}

	validConfidentialVMInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_DC8ads_v5"
	}

	validConfidentialVMSNPInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_DC8ads_v5"
	}

	validConfidentialVMTDXInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_DC8eds_v5"
	}

	invalidConfidentialVMInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_B4ms"
	}

	invalidConfidentialVMSGXInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_DC8s_v3"
	}

	validTrustedLaunchInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D8s_v3"
	}

	invalidTrustedLaunchInstanceTypes = func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D8ps_v5"
	}

	invalidateMachineCIDR = func(ic *types.InstallConfig) {
		_, newCidr, _ := net.ParseCIDR("192.168.111.0/24")
		ic.MachineNetwork = []types.MachineNetworkEntry{
			{CIDR: ipnet.IPNet{IPNet: *newCidr}},
		}
	}
	invalidResourceSkuRegion = "centralus"

	invalidateVirtualNetwork               = func(ic *types.InstallConfig) { ic.Azure.VirtualNetwork = "invalid-virtual-network" }
	invalidateComputeSubnet                = func(ic *types.InstallConfig) { ic.Azure.ComputeSubnet = "invalid-compute-subnet" }
	invalidateControlPlaneSubnet           = func(ic *types.InstallConfig) { ic.Azure.ControlPlaneSubnet = "invalid-controlplane-subnet" }
	invalidateRegion                       = func(ic *types.InstallConfig) { ic.Azure.Region = "neverland" }
	invalidateRegionCapabilities           = func(ic *types.InstallConfig) { ic.Azure.Region = "australiacentral2" }
	invalidateRegionLetterCase             = func(ic *types.InstallConfig) { ic.Azure.Region = "Central US" }
	removeVirtualNetwork                   = func(ic *types.InstallConfig) { ic.Azure.VirtualNetwork = "" }
	removeSubnets                          = func(ic *types.InstallConfig) { ic.Azure.ComputeSubnet, ic.Azure.ControlPlaneSubnet = "", "" }
	premiumDiskCompute                     = func(ic *types.InstallConfig) { ic.Compute[0].Platform.Azure.OSDisk.DiskType = "Premium_LRS" }
	nonpremiumInstanceTypeDiskCompute      = func(ic *types.InstallConfig) { ic.Compute[0].Platform.Azure.InstanceType = "Standard_D4_v4" }
	premiumDiskControlPlane                = func(ic *types.InstallConfig) { ic.ControlPlane.Platform.Azure.OSDisk.DiskType = "Premium_LRS" }
	nonpremiumInstanceTypeDiskControlPlane = func(ic *types.InstallConfig) { ic.ControlPlane.Platform.Azure.InstanceType = "Standard_D4_v4" }
	// premiumDiskDefault                      = func(ic *types.InstallConfig) { ic.Azure.DefaultMachinePlatform.OSDisk.DiskType = "Premium_LRS" }
	nonpremiumInstanceTypeDiskDefault       = func(ic *types.InstallConfig) { ic.Azure.DefaultMachinePlatform.InstanceType = "Standard_D4_v4" }
	enabledSSDCapabilityControlPlane        = func(ic *types.InstallConfig) { ic.ControlPlane.Platform.Azure.UltraSSDCapability = "Enabled" }
	enabledSSDCapabilityCompute             = func(ic *types.InstallConfig) { ic.Compute[0].Platform.Azure.UltraSSDCapability = "Enabled" }
	enabledSSDCapabilityDefault             = func(ic *types.InstallConfig) { ic.Azure.DefaultMachinePlatform.UltraSSDCapability = "Enabled" }
	vmNetworkingTypeAcceleratedControlPlane = func(ic *types.InstallConfig) { ic.ControlPlane.Platform.Azure.VMNetworkingType = "Accelerated" }
	vmNetworkingTypeAcceleratedCompute      = func(ic *types.InstallConfig) { ic.Compute[0].Platform.Azure.VMNetworkingType = "Accelerated" }
	vmNetworkingTypeAcceleratedDefault      = func(ic *types.InstallConfig) { ic.Azure.DefaultMachinePlatform.VMNetworkingType = "Accelerated" }

	securityTypeConfidentialVMDefaultMachinePlatform = func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.Settings = &azure.SecuritySettings{SecurityType: "ConfidentialVM"}
	}
	securityTypeConfidentialVMControlPlane = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.Settings = &azure.SecuritySettings{SecurityType: "ConfidentialVM"}
	}
	securityTypeConfidentialVMCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.Settings = &azure.SecuritySettings{SecurityType: "ConfidentialVM"}
	}
	securityTypeTrustedLaunchDefaultMachinePlatform = func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.Settings = &azure.SecuritySettings{SecurityType: "TrustedLaunch"}
	}
	securityTypeTrustedLaunchControlPlane = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.Settings = &azure.SecuritySettings{SecurityType: "TrustedLaunch"}
	}
	securityTypeTrustedLaunchCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.Settings = &azure.SecuritySettings{SecurityType: "TrustedLaunch"}
	}

	virtualNetworkAPIResult = &aznetwork.VirtualNetwork{
		Name: &validVirtualNetwork,
	}
	computeSubnetAPIResult = &aznetwork.Subnet{
		Name: &validComputeSubnet,
		SubnetPropertiesFormat: &aznetwork.SubnetPropertiesFormat{
			AddressPrefix: &validComputeSubnetCIDR,
		},
	}
	controlPlaneSubnetAPIResult = &aznetwork.Subnet{
		Name: &validControlPlaneSubnet,
		SubnetPropertiesFormat: &aznetwork.SubnetPropertiesFormat{
			AddressPrefix: &validControlPlaneSubnetCIDR,
		},
	}
	locationsAPIResult = func() *[]azsubs.Location {
		r := []azsubs.Location{}
		for i := 0; i < len(validRegionsList); i++ {
			r = append(r, azsubs.Location{
				Name:        &validRegionsList[i],
				DisplayName: &validRegionsList[i],
			})
		}
		return &r
	}()

	marketplaceImageAPIResult = azenc.VirtualMachineImage{
		Name: to.StringPtr("VMImage"),
		VirtualMachineImageProperties: &azenc.VirtualMachineImageProperties{
			HyperVGeneration: azenc.HyperVGenerationTypesV1,
			Plan:             &azenc.PurchasePlan{},
		},
	}

	marketplaceImageAPIResultNoPlan = azenc.VirtualMachineImage{
		Name: to.StringPtr("VMImage"),
		VirtualMachineImageProperties: &azenc.VirtualMachineImageProperties{
			HyperVGeneration: azenc.HyperVGenerationTypesV1,
		},
	}

	resourcesProviderAPIResult = &azres.Provider{
		Namespace: &validResourceGroupNamespace,
		ResourceTypes: &[]azres.ProviderResourceType{
			{
				ResourceType: &validResourceGroupResourceType,
				Locations:    &resourcesCapableRegionsList,
			},
		},
	}

	diskEncryptionSetID          = "test-encryption-set-id"
	diskEncryptionSetName        = "test-encryption-set-name"
	diskEncryptionSetType        = "test-encryption-set-type"
	diskEncryptionSetLocation    = "disk-encryption-set-location"
	validDiskEncryptionSetResult = &azenc.DiskEncryptionSet{
		ID:       to.StringPtr(diskEncryptionSetID),
		Name:     to.StringPtr(diskEncryptionSetName),
		Type:     to.StringPtr(diskEncryptionSetType),
		Location: to.StringPtr(diskEncryptionSetLocation),
	}
	validConfidentialVMDiskEncryptionSetResult = &azenc.DiskEncryptionSet{
		ID:                      to.StringPtr(diskEncryptionSetID),
		Name:                    to.StringPtr(diskEncryptionSetName),
		Type:                    to.StringPtr(diskEncryptionSetType),
		Location:                to.StringPtr(diskEncryptionSetLocation),
		EncryptionSetProperties: &azenc.EncryptionSetProperties{EncryptionType: azenc.ConfidentialVMEncryptedWithCustomerKey},
	}

	validDiskEncryptionSetSubscriptionID = "test-encryption-set-subscription-id"
	validDiskEncryptionSetResourceGroup  = "test-encryption-set-resource-group"
	validDiskEncryptionSetName           = "test-encryption-set-name"
	validDiskEncryptionSetConfig         = func() *azure.DiskEncryptionSet {
		return &azure.DiskEncryptionSet{
			SubscriptionID: validDiskEncryptionSetSubscriptionID,
			ResourceGroup:  validDiskEncryptionSetResourceGroup,
			Name:           validDiskEncryptionSetName,
		}
	}
	validConfidentialVMDiskEncryptionSetName   = "test-confidential-vm-encryption-set-name"
	validConfidentialVMDiskEncryptionSetConfig = func() *azure.DiskEncryptionSet {
		return &azure.DiskEncryptionSet{
			SubscriptionID: validDiskEncryptionSetSubscriptionID,
			ResourceGroup:  validDiskEncryptionSetResourceGroup,
			Name:           validConfidentialVMDiskEncryptionSetName,
		}
	}
	invalidDiskEncryptionSetName   = "test-encryption-set-invalid-name"
	invalidDiskEncryptionSetConfig = func() *azure.DiskEncryptionSet {
		return &azure.DiskEncryptionSet{
			SubscriptionID: validDiskEncryptionSetSubscriptionID,
			ResourceGroup:  validDiskEncryptionSetResourceGroup,
			Name:           invalidDiskEncryptionSetName,
		}
	}

	validOSImagePublisher            = "test-publisher"
	validOSImageOffer                = "test-offer"
	validOSImageSKU                  = "test-sku"
	validOSImageVersion              = "test-version"
	noPlanOSImageSKU                 = "no-plan-sku"
	invalidOSImageSKU                = "bad-sku"
	erroringOSImageSKU               = "test-sku-gen2"
	erroringLicenseTermsOSImageSKU   = "erroring-license-terms"
	unacceptedLicenseTermsOSImageSKU = "unaccepted-license-terms"
	validOSImage                     = azure.OSImage{
		Publisher: validOSImagePublisher,
		Offer:     validOSImageOffer,
		SKU:       validOSImageSKU,
		Version:   validOSImageVersion,
	}

	validDiskEncryptionSetDefaultMachinePlatform = func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.OSDisk.DiskEncryptionSet = validDiskEncryptionSetConfig()
	}
	validDiskEncryptionSetControlPlane = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.OSDisk.DiskEncryptionSet = validDiskEncryptionSetConfig()
	}
	validDiskEncryptionSetCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.OSDisk.DiskEncryptionSet = validDiskEncryptionSetConfig()
	}
	invalidDiskEncryptionSetDefaultMachinePlatform = func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.OSDisk.DiskEncryptionSet = invalidDiskEncryptionSetConfig()
	}
	invalidDiskEncryptionSetControlPlane = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.OSDisk.DiskEncryptionSet = invalidDiskEncryptionSetConfig()
	}
	invalidDiskEncryptionSetCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.OSDisk.DiskEncryptionSet = invalidDiskEncryptionSetConfig()
	}

	validConfidentialVMDiskEncryptionSetDefaultMachinePlatform = func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: validConfidentialVMDiskEncryptionSetConfig()}
	}
	validConfidentialVMDiskEncryptionSetControlPlane = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: validConfidentialVMDiskEncryptionSetConfig()}
	}
	validConfidentialVMDiskEncryptionSetCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: validConfidentialVMDiskEncryptionSetConfig()}
	}
	invalidConfidentialVMDiskEncryptionSetDefaultMachinePlatform = func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: invalidDiskEncryptionSetConfig()}
	}
	invalidConfidentialVMDiskEncryptionSetControlPlane = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: invalidDiskEncryptionSetConfig()}
	}
	invalidConfidentialVMDiskEncryptionSetCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: invalidDiskEncryptionSetConfig()}
	}
	invalidTypeConfidentialVMDiskEncryptionSetDefaultMachinePlatform = func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: validDiskEncryptionSetConfig()}
	}
	invalidTypeConfidentialVMDiskEncryptionSetControlPlane = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: validDiskEncryptionSetConfig()}
	}
	invalidTypeConfidentialVMDiskEncryptionSetCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{DiskEncryptionSet: validDiskEncryptionSetConfig()}
	}

	validOSImageCompute = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.OSImage = validOSImage
	}
	invalidOSImageCompute = func(ic *types.InstallConfig) {
		validOSImageCompute(ic)
		ic.Compute[0].Platform.Azure.OSImage.SKU = invalidOSImageSKU
	}
	erroringLicenseTermsOSImageCompute = func(ic *types.InstallConfig) {
		validOSImageCompute(ic)
		ic.Compute[0].Platform.Azure.OSImage.SKU = erroringLicenseTermsOSImageSKU
	}
	unacceptedLicenseTermsOSImageCompute = func(ic *types.InstallConfig) {
		validOSImageCompute(ic)
		ic.Compute[0].Platform.Azure.OSImage.SKU = unacceptedLicenseTermsOSImageSKU
	}
	erroringGenerationOsImageCompute = func(ic *types.InstallConfig) {
		validOSImageCompute(ic)
		ic.Compute[0].Platform.Azure.OSImage.SKU = erroringOSImageSKU
	}
	validBootDiagnosticsStorageAccount = "validstorageaccount"
	validBootDiagnosticsResourceGroup  = "valid-resource-group"
	validStorageAccountValues          = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.BootDiagnostics = &azure.BootDiagnostics{
			Type:               v1beta1.UserManagedDiagnosticsStorage,
			ResourceGroup:      validBootDiagnosticsResourceGroup,
			StorageAccountName: validBootDiagnosticsStorageAccount,
		}
	}
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Platform: types.Platform{
			Azure: &azure.Platform{
				Region:                   validRegion,
				NetworkResourceGroupName: validNetworkResourceGroup,
				VirtualNetwork:           validVirtualNetwork,
				ComputeSubnet:            validComputeSubnet,
				ControlPlaneSubnet:       validControlPlaneSubnet,
				DefaultMachinePlatform:   &azure.MachinePool{},
			},
		},
		ControlPlane: &types.MachinePool{
			Architecture: types.ArchitectureAMD64,
			Platform: types.MachinePoolPlatform{
				Azure: &azure.MachinePool{},
			},
		},
		Compute: []types.MachinePool{{
			Architecture: types.ArchitectureAMD64,
			Platform: types.MachinePoolPlatform{
				Azure: &azure.MachinePool{},
			},
		}},
	}
}

func TestAzureInstallConfigValidation(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "Valid virtual network & subnets",
			edits:    editFunctions{},
			errorMsg: "",
		},
		{
			name:     "Valid install config without virtual network & subnets",
			edits:    editFunctions{removeVirtualNetwork, removeSubnets},
			errorMsg: "",
		},
		{
			name:     "Invalid subnet range",
			edits:    editFunctions{invalidateMachineCIDR},
			errorMsg: "subnet .+ address prefix is outside of the specified machine networks",
		},
		{
			name:     "Invalid virtual network",
			edits:    editFunctions{invalidateVirtualNetwork},
			errorMsg: "invalid virtual network",
		},
		{
			name:     "Invalid compute subnet",
			edits:    editFunctions{invalidateComputeSubnet},
			errorMsg: "failed to retrieve compute subnet",
		},
		{
			name:     "Invalid control plane subnet",
			edits:    editFunctions{invalidateControlPlaneSubnet},
			errorMsg: "failed to retrieve control plane subnet",
		},
		{
			name:     "Invalid both subnets",
			edits:    editFunctions{invalidateControlPlaneSubnet, invalidateComputeSubnet},
			errorMsg: "failed to retrieve compute subnet",
		},
		{
			name:     "Valid instance types",
			edits:    editFunctions{validInstanceTypes},
			errorMsg: "",
		},
		{
			name:     "Valid Arm64 instance types",
			edits:    editFunctions{validArm64InstanceTypes},
			errorMsg: "",
		},
		{
			name:     "Invalid arch instance types",
			edits:    editFunctions{invalidArchInstanceTypes},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D8ps_v5": instance type architecture 'Arm64' does not match install config architecture amd64, compute\[0\].platform.azure.type: Invalid value: "Standard_D4ps_v5": instance type architecture 'Arm64' does not match install config architecture amd64, platform.azure.defaultMachinePlatform.type: Invalid value: "Standard_D4ps_v5": instance type architecture 'Arm64' does not match install config architecture amd64\]`,
		},
		{
			name:     "Invalid default machine type",
			edits:    editFunctions{invalidateDefaultInstanceTypes},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_A1_v2": instance type does not meet minimum resource requirements of 4 vCPUsAvailable, controlPlane.platform.azure.type: Invalid value: "Standard_A1_v2": instance type does not meet minimum resource requirements of 16 GB Memory, compute\[0\].platform.azure.type: Invalid value: "Standard_A1_v2": instance type does not meet minimum resource requirements of 2 vCPUsAvailable, compute\[0\].platform.azure.type: Invalid value: "Standard_A1_v2": instance type does not meet minimum resource requirements of 8 GB Memory\]`,
		},
		{
			name:     "Invalid control plane instance types",
			edits:    editFunctions{invalidateControlPlaneInstanceTypes},
			errorMsg: `[controlPlane.platform.azure.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 4 vCPUs, controlPlane.platform.azure.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 16 GB Memory]`,
		},
		{
			name:     "Undefined default instance types",
			edits:    editFunctions{undefinedDefaultInstanceTypes},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Dne_D2_v4": not found in region centralus, compute\[0\].platform.azure.type: Invalid value: "Dne_D2_v4": not found in region centralus, controlPlane.platform.azure.type: Invalid value: "Dne_D2_v4": unable to determine HyperVGeneration version\]`,
		},
		{
			name:     "Invalid compute instance types",
			edits:    editFunctions{invalidateComputeInstanceTypes},
			errorMsg: `\[compute\[0\].platform.azure.type: Invalid value: "Standard_A1_v2": instance type does not meet minimum resource requirements of 2 vCPUsAvailable, compute\[0\].platform.azure.type: Invalid value: "Standard_A1_v2": instance type does not meet minimum resource requirements of 8 GB Memory\]`,
		},
		{
			name:     "Invalid region",
			edits:    editFunctions{invalidateRegion},
			errorMsg: "region \"neverland\" is not valid or not available for this account$",
		},
		{
			name:     "Invalid region uncapable",
			edits:    editFunctions{invalidateRegionCapabilities},
			errorMsg: "region \"australiacentral2\" does not support resource creation$",
		},
		{
			name:     "Invalid region letter case",
			edits:    editFunctions{invalidateRegionLetterCase},
			errorMsg: "region \"Central US\" is not valid or not available for this account, did you mean \"centralus\"\\?$",
		},
		{
			name:     "Non-premium instance disk type for compute",
			edits:    editFunctions{premiumDiskCompute, nonpremiumInstanceTypeDiskCompute},
			errorMsg: `compute\[0\].platform.azure.osDisk.diskType: Invalid value: "Premium_LRS": PremiumIO not supported for instance type Standard_D4_v4`,
		},
		{
			name:     "Non-premium instance disk type for control-plane",
			edits:    editFunctions{premiumDiskControlPlane, nonpremiumInstanceTypeDiskControlPlane},
			errorMsg: `controlPlane.platform.azure.osDisk.diskType: Invalid value: "Premium_LRS": PremiumIO not supported for instance type Standard_D4_v4$`,
		},
		{
			name:     "Supported AcceleratedNetworking as default",
			edits:    editFunctions{validVMNetworkingInstanceTypes, vmNetworkingTypeAcceleratedDefault},
			errorMsg: "",
		},
		{
			name:     "Unsupported VMNetworkingType in Control Plane",
			edits:    editFunctions{invalidVMNetworkingIstanceTypes, vmNetworkingTypeAcceleratedControlPlane},
			errorMsg: `controlPlane.platform.azure.vmNetworkingType: Invalid value: "Accelerated": vm networking type is not supported for instance type Standard_B4ms`,
		},
		{
			name:     "Unsupported VMNetworkingType in Compute",
			edits:    editFunctions{invalidVMNetworkingIstanceTypes, vmNetworkingTypeAcceleratedCompute},
			errorMsg: `compute\[0\].platform.azure.vmNetworkingType: Invalid value: "Accelerated": vm networking type is not supported for instance type Standard_B4ms`,
		},
		{
			name:     "Supported ConfidentialVM security type SNP",
			edits:    editFunctions{validConfidentialVMSNPInstanceTypes, securityTypeConfidentialVMControlPlane},
			errorMsg: "",
		},
		{
			name:     "Supported ConfidentialVM security type TDX",
			edits:    editFunctions{validConfidentialVMTDXInstanceTypes, securityTypeConfidentialVMControlPlane},
			errorMsg: "",
		},
		{
			name:     "Unsupported ConfidentialVM security type in control plane",
			edits:    editFunctions{invalidConfidentialVMInstanceTypes, securityTypeConfidentialVMControlPlane},
			errorMsg: `controlPlane.platform.azure.settings.securityType: Invalid value: "ConfidentialVM": this security type is not supported for instance type Standard_B4ms`,
		},
		{
			name:     "Unsupported ConfidentialVM security type in control plane",
			edits:    editFunctions{invalidConfidentialVMSGXInstanceTypes, securityTypeConfidentialVMControlPlane},
			errorMsg: `controlPlane.platform.azure.settings.securityType: Invalid value: "ConfidentialVM": this security type is not supported for instance type Standard_DC8s_v3`,
		},
		{
			name:     "Unsupported ConfidentialVM security type in compute",
			edits:    editFunctions{invalidConfidentialVMInstanceTypes, securityTypeConfidentialVMCompute},
			errorMsg: `compute\[0\].platform.azure.settings.securityType: Invalid value: "ConfidentialVM": this security type is not supported for instance type Standard_B4ms`,
		},
		{
			name:     "Unsupported ConfidentialVM security type in compute for SGX instance type",
			edits:    editFunctions{invalidConfidentialVMSGXInstanceTypes, securityTypeConfidentialVMCompute},
			errorMsg: `compute\[0\].platform.azure.settings.securityType: Invalid value: "ConfidentialVM": this security type is not supported for instance type Standard_DC8s_v3`,
		},
		{
			name:     "Unsupported ConfidentialVM security type in default machine platform",
			edits:    editFunctions{invalidConfidentialVMInstanceTypes, securityTypeConfidentialVMDefaultMachinePlatform},
			errorMsg: `[compute\[0\].platform.azure.settings.securityType: Invalid value: "ConfidentialVM": this security type is not supported for instance type Standard_B4ms,controlPlane.platform.azure.settings.securityType: Invalid valud: "ConfidentialVM": this security type is not supported for instance type Standard_B4ms]`,
		},
		{
			name:     "Supported TrustedLaunch security type",
			edits:    editFunctions{validTrustedLaunchInstanceTypes, securityTypeTrustedLaunchControlPlane},
			errorMsg: "",
		},
		{
			name:     "Unsupported TrustedLaunch security type in control plane",
			edits:    editFunctions{validArm64InstanceTypes, invalidTrustedLaunchInstanceTypes, securityTypeTrustedLaunchControlPlane},
			errorMsg: `controlPlane.platform.azure.settings.securityType: Invalid value: "TrustedLaunch": this security type is not supported for instance type Standard_D8ps_v5`,
		},
		{
			name:     "Unsupported TrustedLaunch security type in compute",
			edits:    editFunctions{validArm64InstanceTypes, invalidTrustedLaunchInstanceTypes, securityTypeTrustedLaunchCompute},
			errorMsg: `compute\[0\].platform.azure.settings.securityType: Invalid value: "TrustedLaunch": this security type is not supported for instance type Standard_D4ps_v5`,
		},
		{
			name:     "Unsupported TrustedLaunch security type for Confindential VM size in compute",
			edits:    editFunctions{validConfidentialVMInstanceTypes, securityTypeTrustedLaunchCompute},
			errorMsg: `compute\[0\].platform.azure.settings.securityType: Invalid value: "TrustedLaunch": this security type is not supported for instance type Standard_DC8ads_v5`,
		},
		{
			name:     "Unsupported TrustedLaunch security type in default machine platform",
			edits:    editFunctions{validArm64InstanceTypes, invalidTrustedLaunchInstanceTypes, securityTypeTrustedLaunchDefaultMachinePlatform},
			errorMsg: `[compute\[0\].platform.azure.settings.securityType: Invalid value: "TrustedLaunch": this security type is not supported for instance type Standard_D4ps_v5,controlPlane.platform.azure.settings.securityType: Invalid valud: "ConfidentialVM": this security type is not supported for instance type Standard_D4ps_v5]`,
		},
		{
			name:     "BootDiagnostics type user managed and valid storage account values",
			edits:    editFunctions{validStorageAccountValues},
			errorMsg: "",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	// InstanceType
	for _, value := range instanceTypeSku {
		azureClient.EXPECT().GetVirtualMachineSku(gomock.Any(), to.String(value.Name), gomock.Any()).Return(value, nil).AnyTimes()
	}
	azureClient.EXPECT().GetVirtualMachineSku(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	for key, value := range vmCapabilities {
		azureClient.EXPECT().GetVMCapabilities(gomock.Any(), key, validRegion).Return(value, nil).AnyTimes()
	}
	azureClient.EXPECT().GetVMCapabilities(gomock.Any(), "Dne_D2_v4", validRegion).Return(nil, fmt.Errorf("not found in region centralus")).AnyTimes()
	azureClient.EXPECT().GetVMCapabilities(gomock.Any(), gomock.Any(), gomock.Any()).Return(vmCapabilities["Standard_D8s_v3"], nil).AnyTimes()

	// VirtualNetwork
	azureClient.EXPECT().GetVirtualNetwork(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork).Return(virtualNetworkAPIResult, nil).AnyTimes()
	azureClient.EXPECT().GetVirtualNetwork(gomock.Any(), gomock.Not(validNetworkResourceGroup), gomock.Not(validVirtualNetwork)).Return(&aznetwork.VirtualNetwork{}, fmt.Errorf("invalid network resource group")).AnyTimes()
	azureClient.EXPECT().GetVirtualNetwork(gomock.Any(), validNetworkResourceGroup, gomock.Not(validVirtualNetwork)).Return(&aznetwork.VirtualNetwork{}, fmt.Errorf("invalid virtual network")).AnyTimes()

	// ComputeSubnet
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, validComputeSubnet).Return(computeSubnetAPIResult, nil).AnyTimes()
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), gomock.Not(validNetworkResourceGroup), validVirtualNetwork, validComputeSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid network resource group")).AnyTimes()
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), validNetworkResourceGroup, gomock.Not(validVirtualNetwork), validComputeSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid virtual network")).AnyTimes()
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, gomock.Not(validComputeSubnet)).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid compute subnet")).AnyTimes()

	// ControlPlaneSubnet
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, validControlPlaneSubnet).Return(controlPlaneSubnetAPIResult, nil).AnyTimes()
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), gomock.Not(validNetworkResourceGroup), validVirtualNetwork, validControlPlaneSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid network resource group")).AnyTimes()
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), validNetworkResourceGroup, gomock.Not(validVirtualNetwork), validControlPlaneSubnet).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid virtual network")).AnyTimes()
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, gomock.Not(validControlPlaneSubnet)).Return(&aznetwork.Subnet{}, fmt.Errorf("invalid control plane subnet")).AnyTimes()

	// Location
	azureClient.EXPECT().ListLocations(gomock.Any()).Return(locationsAPIResult, nil).AnyTimes()

	// ResourceProvider
	azureClient.EXPECT().GetResourcesProvider(gomock.Any(), validResourceGroupNamespace).Return(resourcesProviderAPIResult, nil).AnyTimes()

	// Resource SKUs
	azureClient.EXPECT().GetDiskSkus(gomock.Any(), validResourceSkuRegions).Return(nil, fmt.Errorf("invalid disk type")).AnyTimes()
	azureClient.EXPECT().GetDiskSkus(gomock.Any(), invalidResourceSkuRegion).Return(nil, fmt.Errorf("invalid region")).AnyTimes()

	azureClient.EXPECT().GetAvailabilityZones(gomock.Any(), gomock.Any(), gomock.Any()).Return([]string{"1", "2", "3"}, nil).AnyTimes()

	azureClient.EXPECT().GetVirtualMachineFamily(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	azureClient.EXPECT().CheckIfExistsStorageAccount(gomock.Any(), validBootDiagnosticsResourceGroup, validBootDiagnosticsStorageAccount, validRegion).Return(nil)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := Validate(azureClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

var validGroupResult = &azres.Group{
	ID:       to.StringPtr("valid-resource-group"),
	Location: to.StringPtr("centralus"),
}

var invalidGroupOutsideRegionResult = &azres.Group{
	ID:       to.StringPtr("invalid-resource-group-useast2"),
	Location: to.StringPtr("useast2"),
}

var validGroupWithTagsResult = &azres.Group{
	ID:       to.StringPtr("valid-resource-group-tags"),
	Location: to.StringPtr("centralus"),
	Tags: map[string]*string{
		"key": to.StringPtr("value"),
	},
}

var validGroupWithConflictinsTagsResult = &azres.Group{
	ID:       to.StringPtr("valid-resource-group-conf-tags"),
	Location: to.StringPtr("centralus"),
	Tags: map[string]*string{
		"kubernetes.io_cluster.test-cluster-12345": to.StringPtr("owned"),
	},
}

func Test_validateResourceGroup(t *testing.T) {
	cases := []struct {
		groupName string
		wantSkip  bool
		err       string
	}{{
		groupName: "non-existent-group",
		err:       `^\Qplatform.azure.resourceGroupName: Internal error: failed to get resource group: resource group /resourceGroups/non-existent-group was not found\E$`,
	}, {
		groupName: "valid-resource-group",
	}, {
		groupName: "invalid-resource-group-useast2",
		err:       `^\Qplatform.azure.resourceGroupName: Invalid value: "invalid-resource-group-useast2": expected to in region centralus, but found it to be in useast2\E$`,
	}, {
		groupName: "valid-resource-group-tags",
	}, {
		groupName: "valid-resource-group-conf-tags",
		err:       `^\Qplatform.azure.resourceGroupName: Invalid value: "valid-resource-group-conf-tags": resource group has conflicting tags kubernetes.io_cluster.test-cluster-12345\E$`,
	}, {
		groupName: "valid-resource-group-with-resources",
		err:       `^\Qplatform.azure.resourceGroupName: Invalid value: "valid-resource-group-with-resources": resource group must be empty but it has 3 resources like id1, id2 ...\E$`,
	}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)
	azureClient.EXPECT().GetGroup(gomock.Any(), "non-existent-group").Return(nil, fmt.Errorf("resource group /resourceGroups/non-existent-group was not found")).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group").Return(validGroupResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "invalid-resource-group-useast2").Return(invalidGroupOutsideRegionResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group-with-resources").Return(validGroupResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group-tags").Return(validGroupWithTagsResult, nil).AnyTimes()
	azureClient.EXPECT().GetGroup(gomock.Any(), "valid-resource-group-conf-tags").Return(validGroupWithConflictinsTagsResult, nil).AnyTimes()
	azureClient.EXPECT().ListResourceIDsByGroup(gomock.Any(), gomock.Not("valid-resource-group-with-resources")).Return(nil, nil).AnyTimes()
	azureClient.EXPECT().ListResourceIDsByGroup(gomock.Any(), "valid-resource-group-with-resources").Return([]string{"id1", "id2", "id3"}, nil).AnyTimes()

	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			if test.wantSkip {
				t.Skip()
			}
			err := validateResourceGroup(azureClient, field.NewPath("platform").Child("azure"), &azure.Platform{ResourceGroupName: test.groupName, Region: "centralus"})
			if test.err != "" {
				assert.Regexp(t, test.err, err.ToAggregate())
			} else {
				assert.NoError(t, err.ToAggregate())
			}
		})
	}
}

func TestCheckAzureStackClusterOSImageSet(t *testing.T) {
	cases := []struct {
		ClusterOSImage string
		err            string
	}{{
		ClusterOSImage: "https://storage.test-endpoint.com/rhcos-image",
		err:            "",
	}, {
		ClusterOSImage: "",
		err:            "^platform.azure.clusterOSImage: Required value: clusterOSImage must be set when installing on Azure Stack$",
	}}
	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			err := checkAzureStackClusterOSImageSet(test.ClusterOSImage, field.NewPath("platform").Child("azure"))
			if test.err != "" {
				assert.Regexp(t, test.err, err.ToAggregate())
			} else {
				assert.NoError(t, err.ToAggregate())
			}
		})
	}
}

func TestValidateAzureStackClusterOSImage(t *testing.T) {
	cases := []struct {
		StorageEndpointSuffix string
		ClusterOSImage        string
		err                   string
	}{{
		StorageEndpointSuffix: "storage.test-endpoint.com",
		ClusterOSImage:        "https://storage.test-endpoint.com/rhcos-image",
		err:                   "",
	}, {
		StorageEndpointSuffix: "storage.test-endpoint.com",
		ClusterOSImage:        "https://storage.not-in-the-cluster.com/rhcos-image",
		err:                   `^platform.azure.clusterOSImage: Invalid value: "https://storage.not-in-the-cluster.com/rhcos-image": clusterOSImage must be in the Azure Stack environment$`,
	}}
	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			err := validateAzureStackClusterOSImage(test.StorageEndpointSuffix, test.ClusterOSImage, field.NewPath("platform").Child("azure"))
			if test.err != "" {
				assert.Regexp(t, test.err, err.ToAggregate())
			} else {
				assert.NoError(t, err.ToAggregate())
			}
		})
	}
}

func TestAzureDiskEncryptionSet(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "Valid disk encryption set for default pool",
			edits:    editFunctions{validDiskEncryptionSetDefaultMachinePlatform},
			errorMsg: "",
		},
		{
			name:     "Invalid disk encryption set for default pool",
			edits:    editFunctions{invalidDiskEncryptionSetDefaultMachinePlatform},
			errorMsg: fmt.Sprintf(`^platform.azure.defaultMachinePlatform.osDisk.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: failed to get disk encryption set$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName),
		},
		{
			name:     "Valid disk encryption set for control-plane",
			edits:    editFunctions{validDiskEncryptionSetControlPlane},
			errorMsg: "",
		},
		{
			name:     "Invalid disk encryption set for control-plane",
			edits:    editFunctions{invalidDiskEncryptionSetControlPlane},
			errorMsg: fmt.Sprintf(`^platform.azure.osDisk.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: failed to get disk encryption set$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName),
		},
		{
			name:     "Valid disk encryption set for compute",
			edits:    editFunctions{validDiskEncryptionSetCompute},
			errorMsg: "",
		},
		{
			name:     "Invalid disk encryption set for compute",
			edits:    editFunctions{invalidDiskEncryptionSetCompute},
			errorMsg: fmt.Sprintf(`^compute\[0\].platform.azure.osDisk.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: failed to get disk encryption set$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName),
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	// DiskEncryptionSet
	azureClient.EXPECT().GetDiskEncryptionSet(gomock.Any(), validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, validDiskEncryptionSetName).Return(validDiskEncryptionSetResult, nil).AnyTimes()
	azureClient.EXPECT().GetDiskEncryptionSet(gomock.Any(), validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName).Return(nil, fmt.Errorf("failed to get disk encryption set")).AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			errors := ValidateDiskEncryptionSet(azureClient, editedInstallConfig)
			aggregatedErrors := errors.ToAggregate()
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestAzureSecurityProfileDiskEncryptionSet(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "Valid security profile disk encryption set for default pool",
			edits:    editFunctions{validConfidentialVMDiskEncryptionSetDefaultMachinePlatform},
			errorMsg: "",
		},
		{
			name:     "Invalid security profile disk encryption set not found for default pool",
			edits:    editFunctions{invalidConfidentialVMDiskEncryptionSetDefaultMachinePlatform},
			errorMsg: fmt.Sprintf(`^platform.azure.defaultMachinePlatform.osDisk.securityProfile.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: failed to get disk encryption set$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName),
		},
		{
			name:     "Invalid security profile disk encryption set with default encryption type for default pool",
			edits:    editFunctions{invalidTypeConfidentialVMDiskEncryptionSetDefaultMachinePlatform},
			errorMsg: fmt.Sprintf(`^platform.azure.defaultMachinePlatform.osDisk.securityProfile.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: the disk encryption set should be created with type %s$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, validDiskEncryptionSetName, azenc.ConfidentialVMEncryptedWithCustomerKey),
		},
		{
			name:     "Valid security profile disk encryption set for control-plane",
			edits:    editFunctions{validConfidentialVMDiskEncryptionSetControlPlane},
			errorMsg: "",
		},
		{
			name:     "Invalid security profile disk encryption set not found for control-plane",
			edits:    editFunctions{invalidConfidentialVMDiskEncryptionSetControlPlane},
			errorMsg: fmt.Sprintf(`^platform.azure.osDisk.securityProfile.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: failed to get disk encryption set$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName),
		},
		{
			name:     "Invalid security profile disk encryption set with default encryption type for control-plane",
			edits:    editFunctions{invalidTypeConfidentialVMDiskEncryptionSetControlPlane},
			errorMsg: fmt.Sprintf(`^platform.azure.osDisk.securityProfile.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: the disk encryption set should be created with type %s$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, validDiskEncryptionSetName, azenc.ConfidentialVMEncryptedWithCustomerKey),
		},
		{
			name:     "Valid security profile disk encryption set for compute",
			edits:    editFunctions{validConfidentialVMDiskEncryptionSetCompute},
			errorMsg: "",
		},
		{
			name:     "Invalid security profile disk encryption set not found for compute",
			edits:    editFunctions{invalidConfidentialVMDiskEncryptionSetCompute},
			errorMsg: fmt.Sprintf(`^compute\[0\].platform.azure.osDisk.securityProfile.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: failed to get disk encryption set$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName),
		},
		{
			name:     "Invalid security profile disk encryption set with default encryption type for compute",
			edits:    editFunctions{invalidTypeConfidentialVMDiskEncryptionSetCompute},
			errorMsg: fmt.Sprintf(`^compute\[0\].platform.azure.osDisk.securityProfile.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: the disk encryption set should be created with type %s$`, validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, validDiskEncryptionSetName, azenc.ConfidentialVMEncryptedWithCustomerKey),
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	azureClient.EXPECT().GetDiskEncryptionSet(gomock.Any(), validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, validConfidentialVMDiskEncryptionSetName).Return(validConfidentialVMDiskEncryptionSetResult, nil).AnyTimes()
	azureClient.EXPECT().GetDiskEncryptionSet(gomock.Any(), validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, validDiskEncryptionSetName).Return(validDiskEncryptionSetResult, nil).AnyTimes()
	azureClient.EXPECT().GetDiskEncryptionSet(gomock.Any(), validDiskEncryptionSetSubscriptionID, validDiskEncryptionSetResourceGroup, invalidDiskEncryptionSetName).Return(nil, fmt.Errorf("failed to get disk encryption set")).AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			errors := ValidateSecurityProfileDiskEncryptionSet(azureClient, editedInstallConfig)
			aggregatedErrors := errors.ToAggregate()
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestAzureUltraSSDCapability(t *testing.T) {
	locationInfoFull := &azenc.ResourceSkuLocationInfo{
		Location: to.StringPtr("centralus"),
		ZoneDetails: &[]azenc.ResourceSkuZoneDetails{
			{
				Name: to.StringSlicePtr([]string{"1", "3", "2"}),
				Capabilities: &[]azenc.ResourceSkuCapabilities{
					{Name: to.StringPtr("UltraSSDAvailable"), Value: to.StringPtr("True")},
				},
			},
		},
		Zones: to.StringSlicePtr([]string{"1", "2", "3"}),
	}
	locationInfoNoSSD := &azenc.ResourceSkuLocationInfo{
		Location: to.StringPtr("centralus"),
		ZoneDetails: &[]azenc.ResourceSkuZoneDetails{
			{
				Name:         to.StringSlicePtr([]string{"1", "3", "2"}),
				Capabilities: &[]azenc.ResourceSkuCapabilities{},
			},
		},
		Zones: to.StringSlicePtr([]string{"1", "2", "3"}),
	}
	locationInfoPartial := &azenc.ResourceSkuLocationInfo{
		Location: to.StringPtr("francecentral"),
		ZoneDetails: &[]azenc.ResourceSkuZoneDetails{
			{
				Name: to.StringSlicePtr([]string{"2", "3"}),
				Capabilities: &[]azenc.ResourceSkuCapabilities{
					{Name: to.StringPtr("UltraSSDAvailable"), Value: to.StringPtr("True")},
				},
			},
		},
		Zones: to.StringSlicePtr([]string{"1", "2", "3"}),
	}
	locationInfoSingle := &azenc.ResourceSkuLocationInfo{
		Location:    to.StringPtr("northcentralus"),
		ZoneDetails: &[]azenc.ResourceSkuZoneDetails{},
		Zones:       to.StringSlicePtr(nil),
	}
	locationInfoEmpty := &azenc.ResourceSkuLocationInfo{
		Location:    to.StringPtr("azurestack"),
		ZoneDetails: nil,
		Zones:       to.StringSlicePtr(nil),
	}

	ultraSSDSupportedInstanceTypes := func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D8s_v3"
	}
	ultraSSDUnsupportedInstanceTypes := func(ic *types.InstallConfig) {
		ic.Platform.Azure.DefaultMachinePlatform.InstanceType = "Standard_D2s_v3"
	}
	noZoneRegion := func(ic *types.InstallConfig) {
		ic.Platform.Azure.Region = "azurestack"
	}
	singleZoneRegion := func(ic *types.InstallConfig) {
		ic.Platform.Azure.Region = "northcentralus"
	}
	twoZoneRegion := func(ic *types.InstallConfig) {
		ic.Platform.Azure.Region = "francecentral"
	}

	// User provided availability zones restrictions
	setZones := func(where string, zones ...string) func(ic *types.InstallConfig) {
		switch where {
		case "controlplane", "master":
			return func(ic *types.InstallConfig) {
				ic.ControlPlane.Platform.Azure.Zones = zones
			}
		case "compute", "worker":
			return func(ic *types.InstallConfig) {
				ic.Compute[0].Platform.Azure.Zones = zones
			}
		default:
			return func(ic *types.InstallConfig) {
				ic.Platform.Azure.DefaultMachinePlatform.Zones = zones
			}
		}
	}

	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		// Tests that should fail
		{
			name:     "Unsupported LocationInfo",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, invalidateRegion},
			errorMsg: `\[platform.azure.region: Invalid value: "neverland": region "neverland" is not valid or not available for this account, controlPlane.platform.azure.type: Invalid value: "Standard_D8s_v3": could not determine Availability Zones support in the neverland region: error retrieving availability zones, compute\[0\].platform.azure.type: Invalid value: "Standard_D8s_v3": could not determine Availability Zones support in the neverland region: error retrieving availability zones\]$`,
		},
		{
			name:     "Unsupported UltraSSD in No Zone region when set in DefaultMachine",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDUnsupportedInstanceTypes, noZoneRegion},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones, compute\[0\].platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones\]$`,
		},
		{
			name:     "Unsupported UltraSSD in No Zone region when set in ControlPlane",
			edits:    editFunctions{enabledSSDCapabilityControlPlane, ultraSSDUnsupportedInstanceTypes, noZoneRegion},
			errorMsg: `controlPlane.platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones$`,
		},
		{
			name:     "Unsupported UltraSSD in No Zone region when set in Compute",
			edits:    editFunctions{enabledSSDCapabilityCompute, ultraSSDUnsupportedInstanceTypes, noZoneRegion},
			errorMsg: `compute\[0\].platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones$`,
		},
		{
			name:     "Unsupported UltraSSD in No Zone region when set in ControlPlane and Compute",
			edits:    editFunctions{enabledSSDCapabilityControlPlane, enabledSSDCapabilityCompute, ultraSSDUnsupportedInstanceTypes, noZoneRegion},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones, compute\[0\].platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones\]$`,
		},
		{
			name:     "Unsupported UltraSSD in No Zone region when set in ControlPlane and DefaultMachine",
			edits:    editFunctions{enabledSSDCapabilityDefault, enabledSSDCapabilityControlPlane, ultraSSDUnsupportedInstanceTypes, noZoneRegion},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones, compute\[0\].platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones\]$`,
		},
		{
			name:     "Unsupported UltraSSD in Single Zone region when set in DefaultMachine",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDUnsupportedInstanceTypes, singleZoneRegion},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region northcentralus does not support Availability Zones, compute\[0\].platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region northcentralus does not support Availability Zones\]`,
		},
		{
			name:     "Unsupported UltraSSD in Single Zone region when set in ControlPlane",
			edits:    editFunctions{enabledSSDCapabilityControlPlane, ultraSSDUnsupportedInstanceTypes, singleZoneRegion},
			errorMsg: `controlPlane.platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region northcentralus does not support Availability Zones$`,
		},
		{
			name:     "Unsupported UltraSSD in Single Zone region when set in Compute",
			edits:    editFunctions{enabledSSDCapabilityCompute, ultraSSDUnsupportedInstanceTypes, singleZoneRegion},
			errorMsg: `compute\[0\].platform.azure.type: Invalid value: "Standard_D2s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region northcentralus does not support Availability Zones$`,
		},
		{
			name:     "Unsupported UltraSSD in No Zone region because of Availability Sets",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, noZoneRegion},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones, compute\[0\].platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region azurestack does not support Availability Zones\]$`,
		},
		{
			name:     "Unsupported UltraSSD in Single Zone region because of Availability Sets",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, singleZoneRegion},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region northcentralus does not support Availability Zones, compute\[0\].platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability is not compatible with Availability Sets which are used because region northcentralus does not support Availability Zones\]$`,
		},
		{
			name:     "Unsupported UltraSSD in Two Zone region when set in DefaultMachine and zones not specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, twoZoneRegion},
			errorMsg: `compute\[0\].platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability only supported in zones \[2 3\] for this instance type in the francecentral region`,
		},
		{
			name:     "Unsupported UltraSSD in Two Zone region when set in DefaultMachine and single wrong zone specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, twoZoneRegion, setZones("default", "1")},
			errorMsg: `controlPlane.platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability only supported in zones \[2 3\] for this instance type in the francecentral region`,
		},
		{
			name:     "Unsupported UltraSSD in Two Zone region when set in DefaultMachine and one wrong zone specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, twoZoneRegion, setZones("default", "1", "2")},
			errorMsg: `\[controlPlane.platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability only supported in zones \[2 3\] for this instance type in the francecentral region, compute\[0\].platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability only supported in zones \[2 3\] for this instance type in the francecentral region\]`,
		},
		{
			name:     "Unsupported UltraSSD in Two Zone region when set in DefaultMachine and one wrong zone specified for ControlPlane",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, twoZoneRegion, setZones("master", "1", "2")},
			errorMsg: `controlPlane.platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability only supported in zones \[2 3\] for this instance type in the francecentral region`,
		},
		{
			name:     "Unsupported UltraSSD in Two Zone region when set in DefaultMachine and one wrong zone specified for Compute",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, twoZoneRegion, setZones("worker", "1", "2")},
			errorMsg: `compute\[0\].platform.azure.type: Invalid value: "Standard_D8s_v3": UltraSSD capability only supported in zones \[2 3\] for this instance type in the francecentral region`,
		},
		// Tests that should succeed
		{
			name:     "Unsupported UltraSSD in No Zone region when not set in config",
			edits:    editFunctions{ultraSSDUnsupportedInstanceTypes, noZoneRegion},
			errorMsg: "",
		},
		{
			name:     "Unsupported UltraSSD in Single Zone region when not set in config",
			edits:    editFunctions{ultraSSDUnsupportedInstanceTypes, singleZoneRegion},
			errorMsg: "",
		},
		{
			name:     "Unsupported UltraSSD in Multi Zone region when not set in config",
			edits:    editFunctions{ultraSSDUnsupportedInstanceTypes},
			errorMsg: "",
		},
		{
			name:     "Supported UltraSSD in Two Zone region when set in DefaultMachine and correct zones specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, twoZoneRegion, setZones("default", "2", "3")},
			errorMsg: "",
		},
		{
			name:     "Supported UltraSSD in Two Zone region when set in ControlPlane and one wrong zone specified for DefaultMachine but correct zones for ControlPlane",
			edits:    editFunctions{enabledSSDCapabilityControlPlane, ultraSSDSupportedInstanceTypes, twoZoneRegion, setZones("default", "1", "2"), setZones("master", "2", "3")},
			errorMsg: "",
		},
		{
			name:     "Supported UltraSSD in Two Zone region when set in Compute and one wrong zone specified for DefaultMachine but correct zones for Compute",
			edits:    editFunctions{enabledSSDCapabilityCompute, ultraSSDSupportedInstanceTypes, twoZoneRegion, setZones("default", "1", "2"), setZones("worker", "2", "3")},
			errorMsg: "",
		},
		{
			name:     "Supported UltraSSD in Multi Zone region when set in DefaultMachine and no zones specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes},
			errorMsg: "",
		},
		{
			name:     "Supported UltraSSD in Multi Zone region when set in DefaultMachine and single zone specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, setZones("default", "1")},
			errorMsg: "",
		},
		{
			name:     "Supported UltraSSD in Multi Zone region when set in DefaultMachine and two zones specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, setZones("default", "1", "2")},
			errorMsg: "",
		},
		{
			name:     "Supported UltraSSD in MultiZone region when set in DefaultMachine and all zones specified",
			edits:    editFunctions{enabledSSDCapabilityDefault, ultraSSDSupportedInstanceTypes, setZones("default", "1", "2", "3")},
			errorMsg: "",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	azureClient.EXPECT().GetVMCapabilities(gomock.Any(), "Standard_D8s_v3", gomock.Any()).Return(vmCapabilities["Standard_D8s_v3"], nil).AnyTimes()
	azureClient.EXPECT().GetVMCapabilities(gomock.Any(), "Standard_D2s_v3", gomock.Any()).Return(vmCapabilities["Standard_D2s_v3"], nil).AnyTimes()
	azureClient.EXPECT().GetVMCapabilities(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	azureClient.EXPECT().GetLocationInfo(gomock.Any(), "centralus", "Standard_D8s_v3").Return(locationInfoFull, nil).AnyTimes()
	azureClient.EXPECT().GetLocationInfo(gomock.Any(), "centralus", "Standard_D2s_v3").Return(locationInfoNoSSD, nil).AnyTimes()
	azureClient.EXPECT().GetLocationInfo(gomock.Any(), "francecentral", "Standard_D8s_v3").Return(locationInfoPartial, nil).AnyTimes()
	azureClient.EXPECT().GetLocationInfo(gomock.Any(), "francecentral", "Standard_D2s_v3").Return(locationInfoNoSSD, nil).AnyTimes()
	azureClient.EXPECT().GetLocationInfo(gomock.Any(), "northcentralus", gomock.Any()).Return(locationInfoSingle, nil).AnyTimes()
	azureClient.EXPECT().GetLocationInfo(gomock.Any(), "azurestack", gomock.Any()).Return(locationInfoEmpty, nil).AnyTimes()
	azureClient.EXPECT().GetLocationInfo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error retrieving availability zones")).AnyTimes()

	// VirtualNetwork
	azureClient.EXPECT().GetVirtualNetwork(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork).Return(virtualNetworkAPIResult, nil).AnyTimes()
	// ComputeSubnet
	azureClient.EXPECT().GetComputeSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, validComputeSubnet).Return(computeSubnetAPIResult, nil).AnyTimes()
	// ControlPlaneSubnet
	azureClient.EXPECT().GetControlPlaneSubnet(gomock.Any(), validNetworkResourceGroup, validVirtualNetwork, validControlPlaneSubnet).Return(controlPlaneSubnetAPIResult, nil).AnyTimes()

	validRegionList := []string{"centralus", "northcentralus", "francecentral", "azurestack"}
	locationsAPIResult = func() *[]azsubs.Location {
		r := []azsubs.Location{}
		for i := 0; i < len(validRegionList); i++ {
			r = append(r, azsubs.Location{Name: to.StringPtr(validRegionList[i]), DisplayName: to.StringPtr(validRegionList[i])})
		}
		return &r
	}()
	// Location
	azureClient.EXPECT().ListLocations(gomock.Any()).Return(locationsAPIResult, nil).AnyTimes()

	resourcesProviderAPIResult = &azres.Provider{
		Namespace: to.StringPtr(validResourceGroupNamespace),
		ResourceTypes: &[]azres.ProviderResourceType{
			{
				ResourceType: &validResourceGroupResourceType,
				Locations:    &validRegionList,
			},
		},
	}
	// ResourceProvider
	azureClient.EXPECT().GetResourcesProvider(gomock.Any(), validResourceGroupNamespace).Return(resourcesProviderAPIResult, nil).AnyTimes()

	azureClient.EXPECT().GetVirtualMachineFamily(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := Validate(azureClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestAzureMarketplaceImage(t *testing.T) {
	validOSImageNoPlan := azure.OSImage{
		Plan:      azure.ImageNoPurchasePlan,
		Publisher: validOSImagePublisher,
		SKU:       noPlanOSImageSKU,
		Version:   validOSImageVersion,
		Offer:     validOSImageOffer,
	}

	invalidOSImage := azure.OSImage{
		Publisher: validOSImagePublisher,
		SKU:       invalidOSImageSKU,
		Version:   validOSImageVersion,
		Offer:     validOSImageOffer,
	}

	allHyperVGens := sets.New("V1", "V2")

	cases := []struct {
		name       string
		osImage    *azure.OSImage
		hyperVGens sets.Set[string]
		errorMsg   string
	}{
		{
			name:       "Valid OS Image",
			osImage:    &validOSImage,
			hyperVGens: allHyperVGens,
			errorMsg:   "",
		},
		{
			name:       "Valid OS Image no purchase plan",
			osImage:    &validOSImageNoPlan,
			hyperVGens: allHyperVGens,
			errorMsg:   "",
		},
		{
			name:       "Invalid OS Image",
			osImage:    &invalidOSImage,
			hyperVGens: allHyperVGens,
			errorMsg:   `compute\[0\].platform.azure.osImage: Invalid value: .*: not found`,
		},
		{
			name: "OS Image causing error determining license terms",
			osImage: &azure.OSImage{
				Publisher: validOSImagePublisher,
				SKU:       erroringLicenseTermsOSImageSKU,
				Offer:     validOSImageOffer,
				Version:   validOSImageVersion,
			},
			hyperVGens: allHyperVGens,
			errorMsg:   `compute\[0\].platform.azure.osImage: Invalid value: .*: could not determine if the license terms for the marketplace image have been accepted: error`,
		},
		{
			name: "OS Image with unaccepted license terms",
			osImage: &azure.OSImage{
				Publisher: validOSImagePublisher,
				SKU:       unacceptedLicenseTermsOSImageSKU,
				Offer:     validOSImageOffer,
				Version:   validOSImageVersion,
			},
			hyperVGens: allHyperVGens,
			errorMsg:   `compute\[0\].platform.azure.osImage: Invalid value: .*: the license terms for the marketplace image have not been accepted`,
		},
		{
			name: "OS Image with wrong HyperV generation",
			osImage: &azure.OSImage{
				Publisher: validOSImagePublisher,
				SKU:       erroringOSImageSKU,
				Offer:     validOSImageOffer,
				Version:   validOSImageVersion,
			},
			hyperVGens: sets.New("V1"),
			errorMsg:   `compute\[0\].platform.azure.osImage: Invalid value: .*: instance type supports HyperVGenerations \[(V[12])\] but the specified image is for HyperVGeneration [^\\1].*$`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	// OS Images
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, validOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, validOSImageSKU).Return(true, nil).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, invalidOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, fmt.Errorf("not found")).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, erroringLicenseTermsOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, erroringLicenseTermsOSImageSKU).Return(false, fmt.Errorf("error")).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, unacceptedLicenseTermsOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, unacceptedLicenseTermsOSImageSKU).Return(false, nil).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, erroringOSImageSKU, validOSImageVersion).Return(azenc.VirtualMachineImage{
		VirtualMachineImageProperties: &azenc.VirtualMachineImageProperties{
			HyperVGeneration: azenc.HyperVGenerationTypesV2,
		},
	}, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, erroringOSImageSKU).Return(true, nil).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, noPlanOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResultNoPlan, nil).AnyTimes()
	// Should not check terms of images with no purchase plan
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, noPlanOSImageSKU).MaxTimes(0)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMarketplaceImage(azureClient, validRegion, tc.hyperVGens, tc.osImage, field.NewPath("compute").Index(0))
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAzureMarketplaceImages(t *testing.T) {
	defaultMachineInstanceTypeHyperVGen1 := func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.InstanceType = "Standard_D4s_v3"
	}

	validOSImageDefaultMachine := func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.OSImage = validOSImage
	}
	invalidOSImageDefaultMachine := func(ic *types.InstallConfig) {
		validOSImageDefaultMachine(ic)
		ic.Azure.DefaultMachinePlatform.OSImage.SKU = invalidOSImageSKU
	}
	erroringLicenseTermsOSImageDefaultMachine := func(ic *types.InstallConfig) {
		validOSImageDefaultMachine(ic)
		ic.Azure.DefaultMachinePlatform.OSImage.SKU = erroringLicenseTermsOSImageSKU
	}
	unacceptedLicenseTermsOSImageDefaultMachine := func(ic *types.InstallConfig) {
		validOSImageDefaultMachine(ic)
		ic.Azure.DefaultMachinePlatform.OSImage.SKU = unacceptedLicenseTermsOSImageSKU
	}
	erroringGenerationOsImageDefaultMachine := func(ic *types.InstallConfig) {
		validOSImageDefaultMachine(ic)
		ic.Azure.DefaultMachinePlatform.OSImage.SKU = erroringOSImageSKU
	}

	controlPlaneInstanceTypeHyperVGen1 := func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.InstanceType = "Standard_D4s_v3"
	}

	validOSImageControlPlane := func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.OSImage = validOSImage
	}
	invalidOSImageControlPlane := func(ic *types.InstallConfig) {
		validOSImageControlPlane(ic)
		ic.ControlPlane.Platform.Azure.OSImage.SKU = invalidOSImageSKU
	}
	erroringLicenseTermsOSImageControlPlane := func(ic *types.InstallConfig) {
		validOSImageControlPlane(ic)
		ic.ControlPlane.Platform.Azure.OSImage.SKU = erroringLicenseTermsOSImageSKU
	}
	unacceptedLicenseTermsOSImageControlPlane := func(ic *types.InstallConfig) {
		validOSImageControlPlane(ic)
		ic.ControlPlane.Platform.Azure.OSImage.SKU = unacceptedLicenseTermsOSImageSKU
	}
	erroringGenerationOsImageControlPlane := func(ic *types.InstallConfig) {
		validOSImageControlPlane(ic)
		ic.ControlPlane.Platform.Azure.OSImage.SKU = erroringOSImageSKU
	}

	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		// Marketplace definition on compute nodes
		{
			name:  "Valid OS Image compute",
			edits: editFunctions{validOSImageCompute, validInstanceTypes},
		},
		{
			name:     "Invalid OS Image compute",
			edits:    editFunctions{invalidOSImageCompute, validInstanceTypes},
			errorMsg: `compute\[0\].platform.azure.osImage: Invalid value: .*: not found`,
		},
		{
			name:     "OS Image causing error determining license terms for compute",
			edits:    editFunctions{erroringLicenseTermsOSImageCompute, validInstanceTypes},
			errorMsg: `compute\[0\].platform.azure.osImage: Invalid value: .*: could not determine if the license terms for the marketplace image have been accepted: error`,
		},
		{
			name:     "OS Image with unaccepted license terms for compute",
			edits:    editFunctions{unacceptedLicenseTermsOSImageCompute, validInstanceTypes},
			errorMsg: `compute\[0\].platform.azure.osImage: Invalid value: .*: the license terms for the marketplace image have not been accepted`,
		},
		{
			name:     "OS Image with wrong HyperV generation for compute",
			edits:    editFunctions{erroringGenerationOsImageCompute, validInstanceTypes},
			errorMsg: `compute\[0\].platform.azure.osImage: Invalid value: .*: instance type supports HyperVGenerations \[(V[12])\] but the specified image is for HyperVGeneration [^\\1].*`,
		},
		{
			name:     "OS Image with wrong HyperV generation for compute from default MachinePool",
			edits:    editFunctions{erroringGenerationOsImageDefaultMachine, validInstanceTypes},
			errorMsg: `compute\[0\].platform.azure.osImage: Invalid value: .*: instance type supports HyperVGenerations \[(V[12])\] but the specified image is for HyperVGeneration [^\\1].*`,
		},
		{
			name:     "OS Image with wrong HyperV generation for inherited instance type for compute",
			edits:    editFunctions{erroringGenerationOsImageCompute, defaultMachineInstanceTypeHyperVGen1},
			errorMsg: `compute\[0\].platform.azure.osImage: Invalid value: .*: instance type supports HyperVGenerations \[(V[12])\] but the specified image is for HyperVGeneration [^\\1].*`,
		},
		// Marketplace definition for ControlPlane nodes
		{
			name:  "Valid OS Image controlPlane",
			edits: editFunctions{validOSImageControlPlane, validInstanceTypes},
		},
		{
			name:     "Invalid OS Image controlPlane",
			edits:    editFunctions{invalidOSImageControlPlane, validOSImageCompute, validInstanceTypes},
			errorMsg: `controlPlane.platform.azure.osImage: Invalid value: .*: not found`,
		},
		{
			name:     "Invalid OS Image controlPlane from default MachinePool",
			edits:    editFunctions{invalidOSImageDefaultMachine, validOSImageCompute, validInstanceTypes},
			errorMsg: `controlPlane.platform.azure.osImage: Invalid value: .*: not found`,
		},
		{
			name:     "OS Image causing error determining license terms for controlPlane",
			edits:    editFunctions{erroringLicenseTermsOSImageControlPlane, validInstanceTypes},
			errorMsg: `controlPlane.platform.azure.osImage: Invalid value: .*: could not determine if the license terms for the marketplace image have been accepted: error`,
		},
		{
			name:     "OS Image with unaccepted license terms for controlPlane",
			edits:    editFunctions{unacceptedLicenseTermsOSImageControlPlane, validInstanceTypes},
			errorMsg: `controlPlane.platform.azure.osImage: Invalid value: .*: the license terms for the marketplace image have not been accepted`,
		},
		{
			name:     "OS Image with wrong HyperV generation for controlPlane",
			edits:    editFunctions{erroringGenerationOsImageControlPlane, validInstanceTypes, controlPlaneInstanceTypeHyperVGen1},
			errorMsg: `controlPlane.platform.azure.osImage: Invalid value: .*: instance type supports HyperVGenerations \[(V[12])\] but the specified image is for HyperVGeneration [^\\1].*`,
		},
		{
			name:     "OS Image with wrong HyperV generation for controlPlane from default MachinePool",
			edits:    editFunctions{erroringGenerationOsImageControlPlane, defaultMachineInstanceTypeHyperVGen1},
			errorMsg: `controlPlane.platform.azure.osImage: Invalid value: .*: instance type supports HyperVGenerations \[(V[12])\] but the specified image is for HyperVGeneration [^\\1].*`,
		},
		{
			name:     "OS Image with wrong HyperV generation for inherited instance type for controlPlane",
			edits:    editFunctions{erroringGenerationOsImageControlPlane, defaultMachineInstanceTypeHyperVGen1},
			errorMsg: `controlPlane.platform.azure.osImage: Invalid value: .*: instance type supports HyperVGenerations \[(V[12])\] but the specified image is for HyperVGeneration [^\\1].*`,
		},
		// Marketplace definition from Default Machine Pool
		{
			name:  "Valid OS Image from default MachinePool",
			edits: editFunctions{validOSImageDefaultMachine, validInstanceTypes},
		},
		{
			name:  "Valid OS Image when default MachinePool not valid",
			edits: editFunctions{invalidOSImageDefaultMachine, validOSImageControlPlane, validOSImageCompute, validInstanceTypes},
		},
		{
			name:     "Invalid OS Image from default MachinePool",
			edits:    editFunctions{invalidOSImageDefaultMachine, validInstanceTypes},
			errorMsg: `^\[controlPlane.platform.azure.osImage: Invalid value: .*: not found, compute\[0\].platform.azure.osImage: Invalid value: .*: not found\]$`,
		},
		{
			name:     "OS Image causing error determining license terms from default MachinePool",
			edits:    editFunctions{erroringLicenseTermsOSImageDefaultMachine, validInstanceTypes},
			errorMsg: `^\[controlPlane.platform.azure.osImage: Invalid value: .*: could not determine if the license terms for the marketplace image have been accepted: error, compute\[0\].platform.azure.osImage: Invalid value: .*: could not determine if the license terms for the marketplace image have been accepted: error\]$`,
		},
		{
			name:     "OS Image with unaccepted license terms from default MachinePool",
			edits:    editFunctions{unacceptedLicenseTermsOSImageDefaultMachine, validInstanceTypes},
			errorMsg: `^\[controlPlane.platform.azure.osImage: Invalid value: .*: the license terms for the marketplace image have not been accepted, compute\[0\].platform.azure.osImage: Invalid value: .*: the license terms for the marketplace image have not been accepted\]$`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	// Marketplace images
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, validOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, validOSImageSKU).Return(true, nil).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, invalidOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, fmt.Errorf("not found")).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, erroringLicenseTermsOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, erroringLicenseTermsOSImageSKU).Return(false, fmt.Errorf("error")).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, unacceptedLicenseTermsOSImageSKU, validOSImageVersion).Return(marketplaceImageAPIResult, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, unacceptedLicenseTermsOSImageSKU).Return(false, nil).AnyTimes()
	azureClient.EXPECT().GetMarketplaceImage(gomock.Any(), validRegion, validOSImagePublisher, validOSImageOffer, erroringOSImageSKU, validOSImageVersion).Return(azenc.VirtualMachineImage{
		VirtualMachineImageProperties: &azenc.VirtualMachineImageProperties{
			HyperVGeneration: azenc.HyperVGenerationTypesV2,
		},
	}, nil).AnyTimes()
	azureClient.EXPECT().AreMarketplaceImageTermsAccepted(gomock.Any(), validOSImagePublisher, validOSImageOffer, erroringOSImageSKU).Return(true, nil).AnyTimes()

	// VM Capabilities
	for key, value := range vmCapabilities {
		azureClient.EXPECT().GetVMCapabilities(gomock.Any(), key, validRegion).Return(value, nil).AnyTimes()
	}
	azureClient.EXPECT().GetVMCapabilities(gomock.Any(), "Dne_D2_v4", validRegion).Return(nil, fmt.Errorf("not found in region centralus")).AnyTimes()
	azureClient.EXPECT().GetVMCapabilities(gomock.Any(), gomock.Any(), gomock.Any()).Return(vmCapabilities["Standard_D8s_v3"], nil).AnyTimes()

	// HyperVGenerations
	azureClient.EXPECT().GetHyperVGenerationVersion(gomock.Any(), gomock.Any(), gomock.Any(), "V1").Return("", fmt.Errorf("instance type Standard_D8s_v3 supports HyperVGenerations [V2] but the specified image is for HyperVGeneration V1; to correct this issue either specify a compatible instance type or change the HyperVGeneration for the image by using a different SKU")).AnyTimes()
	azureClient.EXPECT().GetHyperVGenerationVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("V2", nil).AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}
			aggregatedErrors := validateMarketplaceImages(azureClient, editedInstallConfig)
			err := aggregatedErrors.ToAggregate()
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAzureStackDiskType(t *testing.T) {
	const unsupportedDiskType = "StandardSSD_LRS"

	validDefaultDiskType := func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.DiskType = "Premium_LRS"
	}
	invalidDefaultDiskType := func(ic *types.InstallConfig) {
		ic.Azure.DefaultMachinePlatform.DiskType = unsupportedDiskType
	}
	validControlPlaneDiskType := func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.DiskType = "Standard_LRS"
	}
	invalidControlPlaneDiskType := func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.Azure.DiskType = unsupportedDiskType
	}
	validComputeDiskType := func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.DiskType = "Standard_LRS"
	}
	invalidComputeDiskType := func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.Azure.DiskType = unsupportedDiskType
	}

	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "Valid defaultMachinePlatform DiskType",
			edits:    editFunctions{validDefaultDiskType},
			errorMsg: "",
		},
		{
			name:     "Invalid defaultMachinePlatform DiskType",
			edits:    editFunctions{invalidDefaultDiskType},
			errorMsg: `\[controlPlane.platform.azure.OSDisk.diskType: Invalid value: "StandardSSD_LRS": disk format not supported. Must be one of \[Premium_LRS Standard_LRS\] compute\[0\].platform.azure.OSDisk.diskType: Invalid value: "StandardSSD_LRS": disk format not supported. Must be one of \[Premium_LRS Standard_LRS\]\]`,
		},
		{
			name:     "Valid controlPlane DiskType",
			edits:    editFunctions{validControlPlaneDiskType},
			errorMsg: "",
		},
		{
			name:     "Invalid controlPlane DiskType",
			edits:    editFunctions{invalidControlPlaneDiskType},
			errorMsg: `\[controlPlane.platform.azure.OSDisk.diskType: Invalid value: "StandardSSD_LRS": disk format not supported. Must be one of \[Premium_LRS Standard_LRS\]\]`,
		},
		{
			name:     "Valid compute DiskType",
			edits:    editFunctions{validComputeDiskType},
			errorMsg: "",
		},
		{
			name:     "Invalid compute DiskType",
			edits:    editFunctions{invalidComputeDiskType},
			errorMsg: `\[compute\[0\].platform.azure.OSDisk.diskType: Invalid value: "StandardSSD_LRS": disk format not supported. Must be one of \[Premium_LRS Standard_LRS\]\]`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	azureClient := mock.NewMockAPI(mockCtrl)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := validateAzureStackDiskType(azureClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors.ToAggregate())
			}
		})
	}
}
