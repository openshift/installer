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
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

const (
	genV2Suffix string = "-gen2"
)

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(platform *azure.Platform, pool *types.MachinePool, userDataSecret string, clusterID string, role string, capabilities map[string]string, useImageGallery bool, userTags map[string]string, hyperVGen string, subnet string, resourceGroup string, subscriptionID string) ([]*asset.RuntimeFile, error) {
	if poolPlatform := pool.Platform.Name(); poolPlatform != azure.Name {
		return nil, fmt.Errorf("non-Azure machine-pool: %q", poolPlatform)
	}
	mpool := pool.Platform.Azure

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	if len(mpool.Zones) == 0 {
		// if no azs are given we set to []string{""} for convenience over later operations.
		// It means no-zoned for the machine API
		mpool.Zones = []string{""}
	}
	tags, err := CapzTagsFromUserTags(clusterID, userTags)
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
				Version: osImage.Version,
			},
		}
	case useImageGallery:
		// image gallery names cannot have dashes
		id := clusterID
		if hyperVGen == "V2" {
			id += genV2Suffix
		}
		imageID := fmt.Sprintf("/resourceGroups/%s/providers/Microsoft.Compute/galleries/gallery_%s/images/%s/versions/latest", resourceGroup, galleryName, id)
		image = &capz.Image{ID: &imageID}
	default:
		imageID := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/galleries/gallery_%s/images/%s", subscriptionID, resourceGroup, galleryName, clusterID)
		if hyperVGen == "V2" && platform.CloudName != azure.StackCloud {
			imageID += genV2Suffix
		}
		image = &capz.Image{ID: &imageID}
	}

	osDisk := capz.OSDisk{
		OSType:     "Linux",
		DiskSizeGB: &mpool.DiskSizeGB,
		ManagedDisk: &capz.ManagedDiskParameters{
			StorageAccountType: mpool.DiskType,
		},
		CachingType: "ReadWrite",
	}
	ultrassd := mpool.UltraSSDCapability == "Enabled"
	additionalCapabilities := &capz.AdditionalCapabilities{
		UltraSSDEnabled: &ultrassd,
	}
	if pool.Platform.Azure.DiskEncryptionSet != nil {
		osDisk.ManagedDisk.DiskEncryptionSet = &capz.DiskEncryptionSetParameters{
			ID: mpool.OSDisk.DiskEncryptionSet.ToID(),
		}
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

	var result []*asset.RuntimeFile
	for idx := int64(0); idx < total; idx++ {
		zone := mpool.Zones[int(idx)%len(mpool.Zones)]
		azureMachine := &capz.AzureMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
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
			},
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
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, role)),
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

	osDisk.ManagedDisk.DiskEncryptionSet = nil
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
			// Do not allocate a public IP since it isn't
			// accessible as we are using an outbound LB for the
			// control plane. This is temporary until we have a
			// workaround for accessing SSH	(Most likely port
			// forwarding SSH off the LB until the bootstrap node
			// is destroyed).
			AllocatePublicIP:       false,
			AdditionalCapabilities: additionalCapabilities,
			SecurityProfile:        securityProfile,
		},
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
