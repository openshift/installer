// Package azure generates Machine objects for azure.
package azure

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	capi "sigs.k8s.io/cluster-api/api/core/v1beta1"

	"github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

const (
	genV2Suffix string = "-gen2"
)

// MachineInput defines the inputs needed to generate a machine asset.
type MachineInput struct {
	Subnet          string
	Role            string
	UserDataSecret  string
	HyperVGen       string
	StorageSuffix   string
	UseImageGallery bool
	Private         bool
	UserTags        map[string]string
	Platform        *aztypes.Platform
	Pool            *types.MachinePool
}

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID, resourceGroup, subscriptionID string, session *icazure.Session, in *MachineInput) ([]*asset.RuntimeFile, error) {
	if poolPlatform := in.Pool.Platform.Name(); poolPlatform != aztypes.Name {
		return nil, fmt.Errorf("non-Azure machine-pool: %q", poolPlatform)
	}
	mpool := in.Pool.Platform.Azure

	total := int64(1)
	if in.Pool.Replicas != nil {
		total = *in.Pool.Replicas
	}

	if len(mpool.Zones) == 0 {
		// if no azs are given we set to []string{""} for convenience over later operations.
		// It means no-zoned for the machine API
		mpool.Zones = []string{""}
	}
	tags, err := CapzTagsFromUserTags(clusterID, in.UserTags)
	if err != nil {
		return nil, fmt.Errorf("failed to create machineapi.TagSpecifications from UserTags: %w", err)
	}

	var image *capz.Image
	osImage := mpool.OSImage
	galleryName := strings.ReplaceAll(clusterID, "-", "_")

	switch {
	case osImage.Publisher != "":
		image = &capz.Image{
			Marketplace: &capz.AzureMarketplaceImage{
				ImagePlan: capz.ImagePlan{
					Publisher: osImage.Publisher,
					Offer:     osImage.Offer,
					SKU:       osImage.SKU,
				},
				Version:         osImage.Version,
				ThirdPartyImage: osImage.Plan != aztypes.ImageNoPurchasePlan,
			},
		}
	case in.UseImageGallery:
		// image gallery names cannot have dashes
		imageID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/gallery_%s/images/%s", subscriptionID, resourceGroup, galleryName, clusterID)
		if in.HyperVGen == "V2" && in.Platform.CloudName != aztypes.StackCloud {
			imageID += genV2Suffix
		}
		image = &capz.Image{ID: &imageID}
	default:
		// AzureStack is the only use for managed images & supports only Gen1 VMs:
		// https://learn.microsoft.com/en-us/azure-stack/user/azure-stack-vm-considerations?view=azs-2501&tabs=az1%2Caz2#vm-differences
		imageID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/images/%s", subscriptionID, resourceGroup, clusterID)
		image = &capz.Image{ID: &imageID}
	}

	// Set up OSDisk
	osDisk := capz.OSDisk{
		OSType:     "Linux",
		DiskSizeGB: &mpool.DiskSizeGB,
		ManagedDisk: &capz.ManagedDiskParameters{
			StorageAccountType: mpool.DiskType,
		},
		CachingType: "ReadWrite",
	}
	if in.Pool.Platform.Azure.DiskEncryptionSet != nil {
		osDisk.ManagedDisk.DiskEncryptionSet = &capz.DiskEncryptionSetParameters{
			ID: mpool.OSDisk.DiskEncryptionSet.ToID(),
		}
	}

	var diskSecurityProfile capz.VMDiskSecurityProfile
	if mpool.OSDisk.SecurityProfile != nil && mpool.OSDisk.SecurityProfile.SecurityEncryptionType != "" {
		diskSecurityProfile = capz.VMDiskSecurityProfile{
			SecurityEncryptionType: capz.SecurityEncryptionType(mpool.OSDisk.SecurityProfile.SecurityEncryptionType),
		}

		if mpool.OSDisk.SecurityProfile.DiskEncryptionSet != nil {
			diskSecurityProfile.DiskEncryptionSet = &capz.DiskEncryptionSetParameters{
				ID: mpool.OSDisk.SecurityProfile.DiskEncryptionSet.ToID(),
			}
		}
		osDisk.ManagedDisk.SecurityProfile = &diskSecurityProfile
	}

	ultrassd := mpool.UltraSSDCapability == "Enabled"
	additionalCapabilities := &capz.AdditionalCapabilities{
		UltraSSDEnabled: &ultrassd,
	}

	machineProfile := generateSecurityProfile(mpool)
	securityProfile := &capz.SecurityProfile{
		EncryptionAtHost: machineProfile.EncryptionAtHost,
		SecurityType:     capz.SecurityTypes(machineProfile.Settings.SecurityType),
	}
	if machineProfile.Settings.ConfidentialVM != nil {
		securityProfile.UefiSettings = &capz.UefiSettings{
			VTpmEnabled:       ptr.To[bool](machineProfile.Settings.ConfidentialVM.UEFISettings.VirtualizedTrustedPlatformModule == v1beta1.VirtualizedTrustedPlatformModulePolicyEnabled),
			SecureBootEnabled: ptr.To[bool](machineProfile.Settings.ConfidentialVM.UEFISettings.SecureBoot == v1beta1.SecureBootPolicyEnabled),
		}
	} else if machineProfile.Settings.TrustedLaunch != nil {
		securityProfile.UefiSettings = &capz.UefiSettings{
			VTpmEnabled:       ptr.To(machineProfile.Settings.TrustedLaunch.UEFISettings.VirtualizedTrustedPlatformModule == v1beta1.VirtualizedTrustedPlatformModulePolicyEnabled),
			SecureBootEnabled: ptr.To(machineProfile.Settings.TrustedLaunch.UEFISettings.SecureBoot == v1beta1.SecureBootPolicyEnabled),
		}
	}

	userAssignedIdentities := make([]capz.UserAssignedIdentity, len(mpool.Identity.UserAssignedIdentities))
	for i, id := range mpool.Identity.UserAssignedIdentities {
		userAssignedIdentities[i] = capz.UserAssignedIdentity{ProviderID: id.ProviderID()}
	}

	// If identity type is UserAssigned, but no identities are provided, the installer
	// will create one. Populate the manifest with a reference to that identity.
	if mpool.Identity.Type == capz.VMIdentityUserAssigned && len(userAssignedIdentities) == 0 {
		userAssignedIdentities = []capz.UserAssignedIdentity{
			{
				ProviderID: fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s-identity", subscriptionID, resourceGroup, clusterID),
			},
		}
	}

	storageAccountName := aztypes.GetStorageAccountName(clusterID)

	defaultDiag := &capz.Diagnostics{
		Boot: &capz.BootDiagnostics{
			StorageAccountType: capz.ManagedDiagnosticsStorage,
		},
	}

	if in.Platform.DefaultMachinePlatform != nil && in.Platform.DefaultMachinePlatform.BootDiagnostics != nil {
		defaultDiag.Boot.StorageAccountType = in.Platform.DefaultMachinePlatform.BootDiagnostics.Type
		if saURI := bootDiagStorageURIBuilder(in.Platform.DefaultMachinePlatform.BootDiagnostics, session.Environment.StorageEndpointSuffix); saURI != "" {
			defaultDiag.Boot.UserManaged = &capz.UserManagedBootDiagnostics{
				StorageAccountURI: saURI,
			}
		}
	}

	var controlPlaneDiag *capz.Diagnostics
	if mpool.BootDiagnostics != nil {
		controlPlaneDiag = &capz.Diagnostics{
			Boot: &capz.BootDiagnostics{
				StorageAccountType: mpool.BootDiagnostics.Type,
			},
		}
		controlPlaneDiag.Boot.StorageAccountType = mpool.BootDiagnostics.Type
		if saURI := bootDiagStorageURIBuilder(mpool.BootDiagnostics, session.Environment.StorageEndpointSuffix); saURI != "" {
			controlPlaneDiag.Boot.UserManaged = &capz.UserManagedBootDiagnostics{
				StorageAccountURI: saURI,
			}
		}
	}
	if controlPlaneDiag == nil {
		controlPlaneDiag = defaultDiag
	}

	var result []*asset.RuntimeFile
	for idx := int64(0); idx < total; idx++ {
		zone := mpool.Zones[int(idx)%len(mpool.Zones)]
		azureMachine := &capz.AzureMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", clusterID, in.Pool.Name, idx),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
					"cluster.x-k8s.io/cluster-name":  clusterID,
				},
			},
			Spec: capz.AzureMachineSpec{
				VMSize:                     mpool.InstanceType,
				FailureDomain:              ptr.To(zone),
				Image:                      image,
				OSDisk:                     osDisk, // required
				AdditionalTags:             tags,
				AdditionalCapabilities:     additionalCapabilities,
				DisableExtensionOperations: ptr.To(true),
				AllocatePublicIP:           false,
				EnableIPForwarding:         false,
				SecurityProfile:            securityProfile,
				NetworkInterfaces: []capz.NetworkInterface{
					{
						SubnetName:            in.Subnet,
						AcceleratedNetworking: ptr.To(mpool.VMNetworkingType == string(aztypes.VMnetworkingTypeAccelerated) || mpool.VMNetworkingType == string(aztypes.AcceleratedNetworkingEnabled)),
					},
				},
				Identity:               mpool.Identity.Type,
				UserAssignedIdentities: userAssignedIdentities,
				Diagnostics:            controlPlaneDiag,
				DataDisks:              mpool.DataDisks,
			},
		}

		if len(zone) == 0 {
			// FailureDomain must be nil (not empty) to trigger availability set.
			azureMachine.Spec.FailureDomain = nil
		}

		if in.Platform.CloudName == aztypes.StackCloud {
			azureMachine.Spec.Diagnostics = &capz.Diagnostics{
				Boot: &capz.BootDiagnostics{
					StorageAccountType: capz.UserManagedDiagnosticsStorage,
					UserManaged: &capz.UserManagedBootDiagnostics{
						StorageAccountURI: fmt.Sprintf("https://%s.blob.%s", storageAccountName, in.StorageSuffix),
					},
				},
			}
		}

		azureMachine.SetGroupVersionKind(capz.GroupVersion.WithKind("AzureMachine"))
		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", azureMachine.Name)},
			Object: azureMachine,
		})

		controlPlaneMachine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: azureMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, in.Role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: capz.GroupVersion.String(),
					Kind:       "AzureMachine",
					Name:       azureMachine.Name,
				},
			},
		}
		controlPlaneMachine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", azureMachine.Name)},
			Object: controlPlaneMachine,
		})
	}

	bootstrapAzureMachine := &capz.AzureMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name: capiutils.GenerateBoostrapMachineName(clusterID),
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
				"install.openshift.io/bootstrap": "",
				"cluster.x-k8s.io/cluster-name":  clusterID,
			},
		},
		Spec: capz.AzureMachineSpec{
			VMSize:                     mpool.InstanceType,
			Image:                      image,
			FailureDomain:              ptr.To(mpool.Zones[0]),
			OSDisk:                     osDisk,
			AdditionalTags:             tags,
			DisableExtensionOperations: ptr.To(true),
			AllocatePublicIP:           !in.Private,
			AdditionalCapabilities:     additionalCapabilities,
			SecurityProfile:            securityProfile,
			Identity:                   mpool.Identity.Type,
			Diagnostics:                controlPlaneDiag,
			UserAssignedIdentities:     userAssignedIdentities,
		},
	}

	if len(mpool.Zones[0]) == 0 {
		// FailureDomain must be nil (not empty) to trigger availability set.
		bootstrapAzureMachine.Spec.FailureDomain = nil
	}

	if in.Platform.CloudName == aztypes.StackCloud {
		bootstrapAzureMachine.Spec.Diagnostics = &capz.Diagnostics{
			Boot: &capz.BootDiagnostics{
				StorageAccountType: capz.UserManagedDiagnosticsStorage,
				UserManaged: &capz.UserManagedBootDiagnostics{
					StorageAccountURI: fmt.Sprintf("https://%s.blob.%s", storageAccountName, in.StorageSuffix),
				},
			},
		}
	}

	bootstrapAzureMachine.SetGroupVersionKind(capz.GroupVersion.WithKind("AzureMachine"))

	result = append(result, &asset.RuntimeFile{
		File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", bootstrapAzureMachine.Name)},
		Object: bootstrapAzureMachine,
	})

	bootstrapMachine := &capi.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name: bootstrapAzureMachine.Name,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capi.MachineSpec{
			ClusterName: clusterID,
			Bootstrap: capi.Bootstrap{
				DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, "bootstrap")),
			},
			InfrastructureRef: v1.ObjectReference{
				APIVersion: capz.GroupVersion.String(),
				Kind:       "AzureMachine",
				Name:       bootstrapAzureMachine.Name,
			},
		},
	}
	bootstrapMachine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

	result = append(result, &asset.RuntimeFile{
		File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", bootstrapMachine.Name)},
		Object: bootstrapMachine,
	})

	return result, nil
}

// CapzTagsFromUserTags converts a map of user tags to a map of capz.Tags.
func CapzTagsFromUserTags(clusterID string, usertags map[string]string) (capz.Tags, error) {
	tags := capz.Tags{}
	tags[fmt.Sprintf("kubernetes.io_cluster.%s", clusterID)] = "owned"

	forbiddenTags := sets.New[string]()
	for key := range tags {
		forbiddenTags.Insert(key)
	}

	userTagKeys := sets.New[string]()
	for key := range usertags {
		userTagKeys.Insert(key)
	}
	if clobberedTags := userTagKeys.Intersection(forbiddenTags); clobberedTags.Len() > 0 {
		return nil, fmt.Errorf("user tag keys %v are not allowed", sets.List(clobberedTags))
	}
	for _, k := range sets.List(userTagKeys) {
		tags[k] = usertags[k]
	}
	return tags, nil
}

func bootDiagStorageURIBuilder(diag *aztypes.BootDiagnostics, storageEndpointSuffix string) string {
	storageAccountURI := "https://%s.blob.%s"
	if diag.Type == capz.UserManagedDiagnosticsStorage && diag.StorageAccountName != "" {
		return fmt.Sprintf(storageAccountURI, diag.StorageAccountName, storageEndpointSuffix)
	}
	return ""
}
