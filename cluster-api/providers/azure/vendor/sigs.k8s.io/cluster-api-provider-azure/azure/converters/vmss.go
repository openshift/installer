/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package converters

import (
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"k8s.io/utils/ptr"
	azprovider "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

const (
	// RegExpStrCommunityGalleryID is a regexp string used for matching community gallery IDs and capturing specific values.
	RegExpStrCommunityGalleryID = `/CommunityGalleries/(?P<gallery>.*)/Images/(?P<name>.*)/Versions/(?P<version>.*)`
	// RegExpStrComputeGalleryID is a regexp string used for matching compute gallery IDs and capturing specific values.
	RegExpStrComputeGalleryID = `/subscriptions/(?P<subID>.*)/resourceGroups/(?P<rg>.*)/providers/Microsoft.Compute/galleries/(?P<gallery>.*)/images/(?P<name>.*)/versions/(?P<version>.*)`
)

// SDKToVMSS converts an Azure SDK VirtualMachineScaleSet to the AzureMachinePool type.
func SDKToVMSS(sdkvmss armcompute.VirtualMachineScaleSet, sdkinstances []armcompute.VirtualMachineScaleSetVM) azure.VMSS {
	vmss := azure.VMSS{
		ID:    ptr.Deref(sdkvmss.ID, ""),
		Name:  ptr.Deref(sdkvmss.Name, ""),
		State: infrav1.ProvisioningState(ptr.Deref(sdkvmss.Properties.ProvisioningState, "")),
	}

	if sdkvmss.SKU != nil {
		vmss.Sku = ptr.Deref(sdkvmss.SKU.Name, "")
		vmss.Capacity = ptr.Deref[int64](sdkvmss.SKU.Capacity, 0)
	}

	for _, zone := range sdkvmss.Zones {
		vmss.Zones = append(vmss.Zones, *zone)
	}

	if len(sdkvmss.Tags) > 0 {
		vmss.Tags = MapToTags(sdkvmss.Tags)
	}

	if len(sdkinstances) > 0 {
		vmss.Instances = make([]azure.VMSSVM, len(sdkinstances))
		orchestrationMode := ptr.Deref(sdkvmss.Properties.OrchestrationMode, "")
		for i, vm := range sdkinstances {
			vmss.Instances[i] = *SDKToVMSSVM(vm)
			vmss.Instances[i].OrchestrationMode = infrav1.OrchestrationModeType(orchestrationMode)
		}
	}

	if sdkvmss.Properties.VirtualMachineProfile != nil &&
		sdkvmss.Properties.VirtualMachineProfile.StorageProfile != nil &&
		sdkvmss.Properties.VirtualMachineProfile.StorageProfile.ImageReference != nil {
		imageRef := sdkvmss.Properties.VirtualMachineProfile.StorageProfile.ImageReference
		vmss.Image = SDKImageToImage(imageRef, sdkvmss.Plan != nil)
	}

	return vmss
}

// SDKVMToVMSSVM converts an Azure SDK VM to a VMSS VM.
func SDKVMToVMSSVM(sdkInstance armcompute.VirtualMachine, mode infrav1.OrchestrationModeType) *azure.VMSSVM {
	instance := azure.VMSSVM{
		ID: ptr.Deref(sdkInstance.ID, ""),
	}

	if sdkInstance.Properties == nil {
		return &instance
	}

	instance.State = infrav1.Creating
	if sdkInstance.Properties.ProvisioningState != nil {
		instance.State = infrav1.ProvisioningState(ptr.Deref(sdkInstance.Properties.ProvisioningState, ""))
	}

	if sdkInstance.Properties.OSProfile != nil && sdkInstance.Properties.OSProfile.ComputerName != nil {
		instance.Name = *sdkInstance.Properties.OSProfile.ComputerName
	}

	if sdkInstance.Properties.StorageProfile != nil && sdkInstance.Properties.StorageProfile.ImageReference != nil {
		imageRef := sdkInstance.Properties.StorageProfile.ImageReference
		instance.Image = SDKImageToImage(imageRef, sdkInstance.Plan != nil)
	}

	if len(sdkInstance.Zones) > 0 {
		// An instance should have only 1 zone, so use the first item of the slice.
		instance.AvailabilityZone = *sdkInstance.Zones[0]
	}

	instance.OrchestrationMode = mode

	return &instance
}

// SDKToVMSSVM converts an Azure SDK VirtualMachineScaleSetVM into an infrav1exp.VMSSVM.
func SDKToVMSSVM(sdkInstance armcompute.VirtualMachineScaleSetVM) *azure.VMSSVM {
	// Convert resourceGroup Name ID ( ProviderID in capz objects )
	var convertedID string
	convertedID, err := azprovider.ConvertResourceGroupNameToLower(ptr.Deref(sdkInstance.ID, ""))
	if err != nil {
		convertedID = ptr.Deref(sdkInstance.ID, "")
	}

	instance := azure.VMSSVM{
		ID:         convertedID,
		InstanceID: ptr.Deref(sdkInstance.InstanceID, ""),
	}

	if sdkInstance.Properties == nil {
		return &instance
	}

	instance.State = infrav1.Creating
	if sdkInstance.Properties.ProvisioningState != nil {
		instance.State = infrav1.ProvisioningState(ptr.Deref(sdkInstance.Properties.ProvisioningState, ""))
	}

	if sdkInstance.Properties.OSProfile != nil && sdkInstance.Properties.OSProfile.ComputerName != nil {
		instance.Name = *sdkInstance.Properties.OSProfile.ComputerName
	}

	if sdkInstance.Resources != nil {
		for _, r := range sdkInstance.Resources {
			if r.Properties.ProvisioningState != nil && r.Name != nil &&
				(*r.Name == azure.BootstrappingExtensionLinux || *r.Name == azure.BootstrappingExtensionWindows) {
				instance.BootstrappingState = infrav1.ProvisioningState(ptr.Deref(r.Properties.ProvisioningState, ""))
				break
			}
		}
	}

	if sdkInstance.Properties.StorageProfile != nil && sdkInstance.Properties.StorageProfile.ImageReference != nil {
		imageRef := sdkInstance.Properties.StorageProfile.ImageReference
		instance.Image = SDKImageToImage(imageRef, sdkInstance.Plan != nil)
	}

	if len(sdkInstance.Zones) > 0 {
		// an instance should only have 1 zone, so we select the first item of the slice
		instance.AvailabilityZone = *sdkInstance.Zones[0]
	}

	return &instance
}

// SDKImageToImage converts a SDK image reference to infrav1.Image.
func SDKImageToImage(sdkImageRef *armcompute.ImageReference, isThirdPartyImage bool) infrav1.Image {
	if sdkImageRef.ID != nil {
		return IDImageRefToImage(*sdkImageRef.ID)
	}
	// community gallery image
	if sdkImageRef.CommunityGalleryImageID != nil {
		return cgImageRefToImage(*sdkImageRef.CommunityGalleryImageID)
	}
	// shared gallery image
	if sdkImageRef.SharedGalleryImageID != nil {
		return sgImageRefToImage(*sdkImageRef.SharedGalleryImageID)
	}
	// marketplace image
	return mpImageRefToImage(sdkImageRef, isThirdPartyImage)
}

// GetOrchestrationMode returns the compute.OrchestrationMode for the given infrav1.OrchestrationModeType.
func GetOrchestrationMode(modeType infrav1.OrchestrationModeType) armcompute.OrchestrationMode {
	if modeType == infrav1.FlexibleOrchestrationMode {
		return armcompute.OrchestrationModeFlexible
	}
	return armcompute.OrchestrationModeUniform
}

// IDImageRefToImage converts an ID to a infrav1.Image with ComputerGallery set or ID, depending on the structure of the ID.
func IDImageRefToImage(id string) infrav1.Image {
	// compute gallery image
	if ok, params := getParams(RegExpStrComputeGalleryID, id); ok {
		return infrav1.Image{
			ComputeGallery: &infrav1.AzureComputeGalleryImage{
				Gallery:        params["gallery"],
				Name:           params["name"],
				Version:        params["version"],
				SubscriptionID: ptr.To(params["subID"]),
				ResourceGroup:  ptr.To(params["rg"]),
			},
		}
	}

	// specific image
	return infrav1.Image{
		ID: &id,
	}
}

// mpImageRefToImage converts a marketplace gallery ImageReference to an infrav1.Image.
func mpImageRefToImage(sdkImageRef *armcompute.ImageReference, isThirdPartyImage bool) infrav1.Image {
	return infrav1.Image{
		Marketplace: &infrav1.AzureMarketplaceImage{
			ImagePlan: infrav1.ImagePlan{
				Publisher: ptr.Deref(sdkImageRef.Publisher, ""),
				Offer:     ptr.Deref(sdkImageRef.Offer, ""),
				SKU:       ptr.Deref(sdkImageRef.SKU, ""),
			},
			Version:         ptr.Deref(sdkImageRef.Version, ""),
			ThirdPartyImage: isThirdPartyImage,
		},
	}
}

// cgImageRefToImage converts a community gallery ImageReference to an infrav1.Image.
func cgImageRefToImage(id string) infrav1.Image {
	if ok, params := getParams(RegExpStrCommunityGalleryID, id); ok {
		return infrav1.Image{
			ComputeGallery: &infrav1.AzureComputeGalleryImage{
				Gallery: params["gallery"],
				Name:    params["name"],
				Version: params["version"],
			},
		}
	}
	return infrav1.Image{}
}

// sgImageRefToImage converts a shared gallery ImageReference to an infrav1.Image.
func sgImageRefToImage(id string) infrav1.Image {
	if ok, params := getParams(RegExpStrComputeGalleryID, id); ok {
		return infrav1.Image{
			SharedGallery: &infrav1.AzureSharedGalleryImage{
				SubscriptionID: params["subID"],
				ResourceGroup:  params["rg"],
				Gallery:        params["gallery"],
				Name:           params["name"],
				Version:        params["version"],
			},
		}
	}
	return infrav1.Image{}
}

func getParams(regStr, str string) (matched bool, params map[string]string) {
	re := regexp.MustCompile(regStr)
	match := re.FindAllStringSubmatch(str, -1)

	if len(match) == 1 {
		params = make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i > 0 && i <= len(match[0]) {
				params[name] = match[0][i]
			}
		}
		matched = true
	}

	return matched, params
}
