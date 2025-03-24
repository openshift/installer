package validation

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/azure/defaults"
)

const (
	enabled = "Enabled"
)

var (
	validSecurityEncryptionTypes = map[azure.SecurityEncryptionTypes]bool{
		azure.SecurityEncryptionTypesVMGuestStateOnly:     true,
		azure.SecurityEncryptionTypesDiskWithVMGuestState: true,
	}

	validSecurityEncryptionTypeValues = func() []string {
		v := make([]string, 0, len(validSecurityEncryptionTypes))
		for n := range validSecurityEncryptionTypes {
			v = append(v, string(n))
		}
		sort.Strings(v)
		return v
	}()
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

	allErrs = append(allErrs, validateSecurityProfile(p, platform.CloudName, fldPath.Child("defaultMachinePlatform"))...)

	if p.VMNetworkingType != "" {
		acceleratedNetworkingOptions := sets.NewString(string(azure.VMnetworkingTypeAccelerated), string(azure.VMNetworkingTypeBasic))
		if !acceleratedNetworkingOptions.Has(p.VMNetworkingType) {
			allErrs = append(allErrs,
				field.NotSupported(fldPath.Child("acceleratedNetworking"),
					p.VMNetworkingType, acceleratedNetworkingOptions.List()))
		}
	}

	allErrs = append(allErrs, validateOSImage(p, fldPath)...)
	allErrs = append(allErrs, validateIdentity(poolName, p, fldPath.Child("identity"))...)

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

func validateSecurityProfile(p *azure.MachinePool, cloudName azure.CloudEnvironment, fieldPath *field.Path) field.ErrorList {
	var errs field.ErrorList

	if p.Settings == nil && p.OSDisk.SecurityProfile == nil {
		return errs
	}
	if p.Settings == nil && p.OSDisk.SecurityProfile.SecurityEncryptionType != "" {
		return append(errs, field.Required(fieldPath.Child("settings"),
			"settings should be set when osDisk.securityProfile.securityEncryptionType is defined."))
	}

	if cloudName == azure.StackCloud {
		return append(errs, field.Invalid(fieldPath.Child("settings").Child("securityType"),
			p.Settings.SecurityType,
			fmt.Sprintf("the securityType field is not supported on %s.", azure.StackCloud)))
	}

	switch p.Settings.SecurityType {
	case azure.SecurityTypesConfidentialVM:
		if p.OSDisk.SecurityProfile == nil || p.OSDisk.SecurityProfile.SecurityEncryptionType == "" {
			securityProfileFieldPath := fieldPath.Child("osDisk").Child("securityProfile")
			return append(errs, field.Required(securityProfileFieldPath.Child("securityEncryptionType"),
				fmt.Sprintf("securityEncryptionType should be set when securityType is set to %s.",
					azure.SecurityTypesConfidentialVM)))
		}

		if !validSecurityEncryptionTypes[p.OSDisk.SecurityProfile.SecurityEncryptionType] {
			securityProfileFieldPath := fieldPath.Child("osDisk").Child("securityProfile")
			return append(errs, field.NotSupported(securityProfileFieldPath.Child("securityEncryptionType"),
				p.OSDisk.SecurityProfile.SecurityEncryptionType, validSecurityEncryptionTypeValues))
		}

		if p.Settings.ConfidentialVM == nil {
			return append(errs, field.Required(fieldPath.Child("settings").Child("confidentialVM"),
				fmt.Sprintf("confidentialVM should be set when securityType is set to %s.",
					azure.SecurityTypesConfidentialVM)))
		}

		if p.Settings.ConfidentialVM.UEFISettings == nil {
			return append(errs, field.Required(fieldPath.Child("settings").Child("confidentialVM").Child("uefiSettings"),
				fmt.Sprintf("uefiSettings should be set when securityType is set to %s.",
					azure.SecurityTypesConfidentialVM)))
		}

		if p.Settings.ConfidentialVM.UEFISettings.VirtualizedTrustedPlatformModule != nil &&
			*p.Settings.ConfidentialVM.UEFISettings.VirtualizedTrustedPlatformModule != enabled {
			uefiSettingsFieldPath := fieldPath.Child("settings").Child("confidentialVM").Child("uefiSettings")
			return append(errs, field.Invalid(uefiSettingsFieldPath.Child("virtualizedTrustedPlatformModule"),
				*p.Settings.ConfidentialVM.UEFISettings.VirtualizedTrustedPlatformModule,
				fmt.Sprintf("virtualizedTrustedPlatformModule should be enabled when securityType is set to %s.",
					azure.SecurityTypesConfidentialVM)))
		}

		if p.OSDisk.SecurityProfile.SecurityEncryptionType == azure.SecurityEncryptionTypesDiskWithVMGuestState {
			if p.EncryptionAtHost {
				return append(errs, field.Invalid(fieldPath.Child("encryptionAtHost"), p.EncryptionAtHost,
					fmt.Sprintf("encryptionAtHost cannot be set to true when securityEncryptionType is set to %s.",
						azure.SecurityEncryptionTypesDiskWithVMGuestState)))
			}

			if p.Settings.ConfidentialVM.UEFISettings.SecureBoot != nil &&
				*p.Settings.ConfidentialVM.UEFISettings.SecureBoot != enabled {
				uefiSettingsFieldPath := fieldPath.Child("settings").Child("confidentialVM").Child("uefiSettings")
				return append(errs, field.Invalid(uefiSettingsFieldPath.Child("secureBoot"),
					*p.Settings.ConfidentialVM.UEFISettings.SecureBoot,
					fmt.Sprintf("secureBoot should be enabled when securityEncryptionType is set to %s.",
						azure.SecurityEncryptionTypesDiskWithVMGuestState)))
			}
		}
	case azure.SecurityTypesTrustedLaunch:
		if p.Settings.TrustedLaunch == nil {
			return append(errs, field.Required(fieldPath.Child("settings").Child("trustedLaunch"),
				fmt.Sprintf("trustedLaunch should be set when securityType is set to %s.",
					azure.SecurityTypesTrustedLaunch)))
		}
	default:
		if p.OSDisk.SecurityProfile != nil && p.OSDisk.SecurityProfile.SecurityEncryptionType != "" {
			return append(errs, field.Invalid(fieldPath.Child("settings").Child("securityType"),
				p.Settings.SecurityType,
				fmt.Sprintf("securityType should be set to %s when securityEncryptionType is defined.",
					azure.SecurityTypesConfidentialVM)))
		}

		if p.Settings.TrustedLaunch != nil && p.Settings.TrustedLaunch.UEFISettings != nil &&
			((p.Settings.TrustedLaunch.UEFISettings.SecureBoot != nil && *p.Settings.TrustedLaunch.UEFISettings.SecureBoot == enabled) ||
				(p.Settings.TrustedLaunch.UEFISettings.VirtualizedTrustedPlatformModule != nil && *p.Settings.TrustedLaunch.UEFISettings.VirtualizedTrustedPlatformModule == enabled)) {
			return append(errs, field.Invalid(fieldPath.Child("settings").Child("securityType"),
				p.Settings.SecurityType,
				fmt.Sprintf("securityType should be set to %s when uefiSettings are enabled.",
					azure.SecurityTypesTrustedLaunch)))
		}
	}

	return errs
}

func validateIdentity(poolName string, p *azure.MachinePool, fldPath *field.Path) field.ErrorList {
	id := p.Identity
	if id == nil {
		return nil
	}

	var errs field.ErrorList
	if poolName == "worker" && id.Type != capz.VMIdentityUserAssigned {
		return append(errs, field.Invalid(fldPath.Child("type"), id.Type, "only user-assigned identities are supported for compute nodes"))
	}

	if id.Type == "" {
		return append(errs, field.Required(fldPath.Child("type"), "type must be specified if using identity"))
	}

	if id.Type != capz.VMIdentityNone && id.Type != capz.VMIdentitySystemAssigned && id.Type != capz.VMIdentityUserAssigned {
		supportedValues := []capz.VMIdentity{capz.VMIdentityNone, capz.VMIdentitySystemAssigned, capz.VMIdentityUserAssigned}
		return append(errs, field.NotSupported(fldPath.Child("type"), id.Type, supportedValues))
	}

	if id.SystemAssignedIdentityRole != nil && id.Type != capz.VMIdentitySystemAssigned {
		errs = append(errs, field.Invalid(fldPath.Child("type"), id.Type, "systemAssignedIdentityRole may only be used with type: SystemAssigned"))
	}

	if id.Type == capz.VMIdentityUserAssigned && len(id.UserAssignedIdentities) == 0 {
		errs = append(errs, field.Required(fldPath.Child("userAssignedIdentities"), "userAssignedIdentities must be specified when using type: UserAssigned"))
	}

	if id.UserAssignedIdentities != nil && id.Type != capz.VMIdentityUserAssigned {
		errs = append(errs, field.Invalid(fldPath.Child("type"), id.Type, "userAssignedIdentities may only be used with type: UserAssigned"))
	}

	if s := id.SystemAssignedIdentityRole; s != nil && uuid.Validate(s.Name) != nil {
		errMsg := "name must be a valid UUID, please provide a valid UUID or leave name black to have one generated for you"
		errs = append(errs, field.Invalid(fldPath.Child("systemAssignedIdentityRole", "name"), s.Name, errMsg))
	}

	return errs
}
