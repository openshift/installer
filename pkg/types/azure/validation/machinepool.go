package validation

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/azure/defaults"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *azure.MachinePool, poolName string, platform *azure.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.OSDisk.DiskSizeGB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), p.OSDisk.DiskSizeGB, "Storage DiskSizeGB must be positive"))
	} else if platform.CloudName == azure.StackCloud && p.OSDisk.DiskSizeGB != 0 && (p.OSDisk.DiskSizeGB < defaults.AzurestackMinimumDiskSize || p.OSDisk.DiskSizeGB > defaults.AzurestackMaximumDiskSize) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), p.OSDisk.DiskSizeGB, "Storage DiskSizeGB must be between 128 and 1023 inclusive for Azure Stack"))
	}

	if p.OSDisk.DiskType != "" {
		diskTypes := sets.NewString(
			"StandardSSD_LRS",
			// "UltraSSD_LRS" needs azure terraform version 2.0
			"Premium_LRS")
		// The control plane cannot use Standard_LRS. Don't let the default machine pool specify "Standard_LRS" either.
		if poolName != "" && poolName != "master" {
			diskTypes.Insert("Standard_LRS")
		}

		if !diskTypes.Has(p.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.OSDisk.DiskType, diskTypes.List()))
		}
	}

	if p.UltraSSDCapability != "" {
		ultraSSDCapabilities := sets.NewString("Enabled", "Disabled")
		if !ultraSSDCapabilities.Has(p.UltraSSDCapability) {
			allErrs = append(allErrs,
				field.NotSupported(fldPath.Child("ultraSSDCapability"),
					p.UltraSSDCapability, ultraSSDCapabilities.List()))
		}
	}

	allErrs = append(allErrs, ValidateEncryptionAtHost(p, platform.CloudName, fldPath.Child("defaultMachinePlatform"))...)
	if p.OSDisk.DiskEncryptionSet != nil {
		allErrs = append(allErrs, ValidateDiskEncryption(p, platform.CloudName, fldPath.Child("defaultMachinePlatform"))...)
	}

	if p.VMNetworkingType != "" {
		acceleratedNetworkingOptions := sets.NewString(string(azure.VMnetworkingTypeAccelerated), string(azure.VMNetworkingTypeBasic))
		if !acceleratedNetworkingOptions.Has(p.VMNetworkingType) {
			allErrs = append(allErrs,
				field.NotSupported(fldPath.Child("acceleratedNetworking"),
					p.VMNetworkingType, acceleratedNetworkingOptions.List()))
		}
	}

	allErrs = append(allErrs, validateOSImage(p, fldPath)...)

	return allErrs
}

func validateOSImage(p *azure.MachinePool, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	osImageFldPath := fldPath.Child("osImage")

	emptyOSImage := azure.OSImage{}
	if p.OSImage != emptyOSImage {
		if p.OSImage.Plan != "" {
			planOptions := sets.NewString(
				string(azure.ImageNoPurchasePlan),
				string(azure.ImageWithPurchasePlan),
			)
			if !planOptions.Has(string(p.OSImage.Plan)) {
				allErrs = append(allErrs, field.NotSupported(osImageFldPath.Child("plan"), p.OSImage.Plan, planOptions.List()))
			}
		}
		if p.OSImage.Publisher == "" {
			allErrs = append(allErrs, field.Required(osImageFldPath.Child("publisher"), "must specify publisher for the OS image"))
		}
		if p.OSImage.Offer == "" {
			allErrs = append(allErrs, field.Required(osImageFldPath.Child("offer"), "must specify offer for the OS image"))
		}
		if p.OSImage.SKU == "" {
			allErrs = append(allErrs, field.Required(osImageFldPath.Child("sku"), "must specify SKU for the OS image"))
		}
		if p.OSImage.Version == "" {
			allErrs = append(allErrs, field.Required(osImageFldPath.Child("version"), "must specify version for the OS image"))
		}
	}

	return allErrs
}
