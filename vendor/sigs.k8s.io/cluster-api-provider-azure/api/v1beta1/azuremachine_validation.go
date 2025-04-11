/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
)

// ValidateAzureMachineSpec checks an AzureMachineSpec and returns any validation errors.
func ValidateAzureMachineSpec(spec AzureMachineSpec) field.ErrorList {
	var allErrs field.ErrorList

	if errs := ValidateImage(spec.Image, field.NewPath("image")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateOSDisk(spec.OSDisk, field.NewPath("osDisk")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateConfidentialCompute(spec.OSDisk.ManagedDisk, spec.SecurityProfile, field.NewPath("securityProfile")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateSSHKey(spec.SSHPublicKey, field.NewPath("sshPublicKey")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateUserAssignedIdentity(spec.Identity, spec.UserAssignedIdentities, field.NewPath("userAssignedIdentities")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateDataDisks(spec.DataDisks, field.NewPath("dataDisks")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateDiagnostics(spec.Diagnostics, field.NewPath("diagnostics")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateNetwork(spec.SubnetName, spec.AcceleratedNetworking, spec.NetworkInterfaces, field.NewPath("networkInterfaces")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateSystemAssignedIdentityRole(spec.Identity, spec.RoleAssignmentName, spec.SystemAssignedIdentityRole, field.NewPath("systemAssignedIdentityRole")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateCapacityReservationGroupID(spec.CapacityReservationGroupID, field.NewPath("capacityReservationGroupID")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := ValidateVMExtensions(spec.DisableExtensionOperations, spec.VMExtensions, field.NewPath("vmExtensions")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	return allErrs
}

// ValidateNetwork validates the network configuration.
func ValidateNetwork(subnetName string, acceleratedNetworking *bool, networkInterfaces []NetworkInterface, fldPath *field.Path) field.ErrorList {
	if (networkInterfaces != nil) && len(networkInterfaces) > 0 && subnetName != "" {
		return field.ErrorList{field.Invalid(fldPath, networkInterfaces, "cannot set both networkInterfaces and machine subnetName")}
	}

	if (networkInterfaces != nil) && len(networkInterfaces) > 0 && acceleratedNetworking != nil {
		return field.ErrorList{field.Invalid(fldPath, networkInterfaces, "cannot set both networkInterfaces and machine acceleratedNetworking")}
	}

	for _, nic := range networkInterfaces {
		if nic.PrivateIPConfigs < 1 {
			return field.ErrorList{field.Invalid(fldPath, networkInterfaces, "number of privateIPConfigs per interface must be at least 1")}
		}
	}

	return field.ErrorList{}
}

// ValidateSSHKey validates an SSHKey.
func ValidateSSHKey(sshKey string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	decoded, err := base64.StdEncoding.DecodeString(sshKey)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, sshKey, "the SSH public key is not properly base64 encoded"))
		return allErrs
	}

	if _, _, _, _, err := ssh.ParseAuthorizedKey(decoded); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, sshKey, "the SSH public key is not valid"))
		return allErrs
	}

	return allErrs
}

// ValidateSystemAssignedIdentity validates the system-assigned identities list.
func ValidateSystemAssignedIdentity(identityType VMIdentity, oldIdentity, newIdentity string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if identityType == VMIdentitySystemAssigned {
		if _, err := uuid.Parse(newIdentity); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, newIdentity, "Role assignment name must be a valid GUID. It is optional and will be auto-generated when not specified."))
		}
		if oldIdentity != "" && oldIdentity != newIdentity {
			allErrs = append(allErrs, field.Invalid(fldPath, newIdentity, "Role assignment name should not be modified after AzureMachine creation."))
		}
	} else if newIdentity != "" {
		allErrs = append(allErrs, field.Forbidden(fldPath, "Role assignment name should only be set when using system assigned identity."))
	}

	return allErrs
}

// ValidateUserAssignedIdentity validates the user-assigned identities list.
func ValidateUserAssignedIdentity(identityType VMIdentity, userAssignedIdentities []UserAssignedIdentity, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(userAssignedIdentities) > 0 && identityType != VMIdentityUserAssigned {
		allErrs = append(allErrs, field.Invalid(fldPath, identityType, "must be set to 'UserAssigned' when assigning any user identity to the machine"))
	}

	if identityType == VMIdentityUserAssigned {
		if len(userAssignedIdentities) == 0 {
			allErrs = append(allErrs, field.Required(fldPath, "must be specified for the 'UserAssigned' identity type"))
		}
		for _, identity := range userAssignedIdentities {
			if identity.ProviderID != "" {
				if _, err := azureutil.ParseResourceID(identity.ProviderID); err != nil {
					allErrs = append(allErrs, field.Invalid(fldPath, identity.ProviderID, "must be a valid Azure resource ID"))
				}
			}
		}
	}

	return allErrs
}

// ValidateSystemAssignedIdentityRole validates the system-assigned identity role.
func ValidateSystemAssignedIdentityRole(identityType VMIdentity, roleAssignmentName string, role *SystemAssignedIdentityRole, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if roleAssignmentName != "" && role != nil && role.Name != "" {
		allErrs = append(allErrs, field.Invalid(fldPath, role.Name, "cannot set both roleAssignmentName and systemAssignedIdentityRole.name"))
	}
	if identityType == VMIdentitySystemAssigned && role != nil {
		if role.DefinitionID == "" {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "systemAssignedIdentityRole", "definitionID"), role.DefinitionID, "the definitionID field cannot be empty"))
		}
		if role.Scope == "" {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "systemAssignedIdentityRole", "scope"), role.Scope, "the scope field cannot be empty"))
		}
	}
	if identityType != VMIdentitySystemAssigned && role != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "role"), "systemAssignedIdentityRole can only be set when identity is set to SystemAssigned"))
	}
	return allErrs
}

// ValidateDataDisks validates a list of data disks.
func ValidateDataDisks(dataDisks []DataDisk, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	lunSet := make(map[int32]struct{})
	nameSet := make(map[string]struct{})
	for _, disk := range dataDisks {
		// validate that the disk size is between 4 and 32767.
		if disk.DiskSizeGB < 4 || disk.DiskSizeGB > 32767 {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("DiskSizeGB"), "", "the disk size should be a value between 4 and 32767"))
		}

		// validate that all names are unique
		if disk.NameSuffix == "" {
			allErrs = append(allErrs, field.Required(fieldPath.Child("NameSuffix"), "the name suffix cannot be empty"))
		}
		if _, ok := nameSet[disk.NameSuffix]; ok {
			allErrs = append(allErrs, field.Duplicate(fieldPath, disk.NameSuffix))
		} else {
			nameSet[disk.NameSuffix] = struct{}{}
		}

		// validate optional managed disk option
		if disk.ManagedDisk != nil {
			if errs := validateManagedDisk(disk.ManagedDisk, fieldPath.Child("managedDisk"), false); len(errs) > 0 {
				allErrs = append(allErrs, errs...)
			}
		}

		// validate that all LUNs are unique and between 0 and 63.
		if disk.Lun == nil {
			allErrs = append(allErrs, field.Required(fieldPath, "LUN should not be nil"))
		} else if *disk.Lun < 0 || *disk.Lun > 63 {
			allErrs = append(allErrs, field.Invalid(fieldPath, disk.Lun, "logical unit number must be between 0 and 63"))
		} else if _, ok := lunSet[*disk.Lun]; ok {
			allErrs = append(allErrs, field.Duplicate(fieldPath, disk.Lun))
		} else {
			lunSet[*disk.Lun] = struct{}{}
		}

		// validate cachingType
		allErrs = append(allErrs, validateCachingType(disk.CachingType, fieldPath, disk.ManagedDisk)...)
	}
	return allErrs
}

// ValidateOSDisk validates the OSDisk spec.
func ValidateOSDisk(osDisk OSDisk, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if osDisk.DiskSizeGB != nil {
		if *osDisk.DiskSizeGB <= 0 || *osDisk.DiskSizeGB > 2048 {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("DiskSizeGB"), "", "the Disk size should be a value between 1 and 2048"))
		}
	}

	if osDisk.OSType == "" {
		allErrs = append(allErrs, field.Required(fieldPath.Child("OSType"), "the OS type cannot be empty"))
	}

	allErrs = append(allErrs, validateCachingType(osDisk.CachingType, fieldPath, osDisk.ManagedDisk)...)

	if osDisk.ManagedDisk != nil {
		if errs := validateManagedDisk(osDisk.ManagedDisk, fieldPath.Child("managedDisk"), true); len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
	}

	if osDisk.DiffDiskSettings != nil && osDisk.ManagedDisk != nil && osDisk.ManagedDisk.DiskEncryptionSet != nil {
		allErrs = append(allErrs, field.Invalid(
			fieldPath.Child("managedDisks").Child("diskEncryptionSet"),
			osDisk.ManagedDisk.DiskEncryptionSet.ID,
			"diskEncryptionSet is not supported when diffDiskSettings.option is 'Local'",
		))
	}
	if osDisk.DiffDiskSettings != nil && osDisk.DiffDiskSettings.Placement != nil {
		if osDisk.DiffDiskSettings.Option != string(armcompute.DiffDiskOptionsLocal) {
			allErrs = append(allErrs, field.Invalid(
				fieldPath.Child("diffDiskSettings"),
				osDisk.DiffDiskSettings,
				"placement is only supported when diffDiskSettings.option is 'Local'",
			))
		}
	}

	return allErrs
}

// validateManagedDisk validates updates to the ManagedDiskParameters field.
func validateManagedDisk(m *ManagedDiskParameters, fieldPath *field.Path, isOSDisk bool) field.ErrorList {
	allErrs := field.ErrorList{}

	if m != nil {
		allErrs = append(allErrs, validateStorageAccountType(m.StorageAccountType, fieldPath.Child("StorageAccountType"), isOSDisk)...)

		// DiskEncryptionSet can only be set when SecurityEncryptionType is set to DiskWithVMGuestState
		// https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/create-or-update?tabs=HTTP#securityencryptiontypes
		if isOSDisk && m.SecurityProfile != nil && m.SecurityProfile.DiskEncryptionSet != nil {
			if m.SecurityProfile.SecurityEncryptionType != SecurityEncryptionTypeDiskWithVMGuestState {
				allErrs = append(allErrs, field.Invalid(
					fieldPath.Child("securityProfile").Child("diskEncryptionSet"),
					m.SecurityProfile.DiskEncryptionSet.ID,
					"diskEncryptionSet is only supported when securityEncryptionType is set to DiskWithVMGuestState",
				))
			}
		}
	}

	return allErrs
}

// ValidateDataDisksUpdate validates updates to Data disks.
func ValidateDataDisksUpdate(oldDataDisks, newDataDisks []DataDisk, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	diskErrMsg := "adding/removing data disks after machine creation is not allowed"
	fieldErrMsg := "modifying data disk's fields after machine creation is not allowed"

	if len(oldDataDisks) != len(newDataDisks) {
		allErrs = append(allErrs, field.Invalid(fieldPath, newDataDisks, diskErrMsg))
		return allErrs
	}

	oldDisks := make(map[string]DataDisk)

	for _, disk := range oldDataDisks {
		oldDisks[disk.NameSuffix] = disk
	}

	for i, newDisk := range newDataDisks {
		if oldDisk, ok := oldDisks[newDisk.NameSuffix]; ok {
			if newDisk.DiskSizeGB != oldDisk.DiskSizeGB {
				allErrs = append(allErrs, field.Invalid(fieldPath.Index(i).Child("diskSizeGB"), newDataDisks, fieldErrMsg))
			}

			allErrs = append(allErrs, validateManagedDisksUpdate(oldDisk.ManagedDisk, newDisk.ManagedDisk, fieldPath.Index(i).Child("managedDisk"))...)

			if (newDisk.Lun != nil && oldDisk.Lun != nil) && (*newDisk.Lun != *oldDisk.Lun) {
				allErrs = append(allErrs, field.Invalid(fieldPath.Index(i).Child("lun"), newDataDisks, fieldErrMsg))
			} else if (newDisk.Lun != nil && oldDisk.Lun == nil) || (newDisk.Lun == nil && oldDisk.Lun != nil) {
				allErrs = append(allErrs, field.Invalid(fieldPath.Index(i).Child("lun"), newDataDisks, fieldErrMsg))
			}

			if newDisk.CachingType != oldDisk.CachingType {
				allErrs = append(allErrs, field.Invalid(fieldPath.Index(i).Child("cachingType"), newDataDisks, fieldErrMsg))
			}
		} else {
			allErrs = append(allErrs, field.Invalid(fieldPath.Index(i).Child("nameSuffix"), newDataDisks, diskErrMsg))
		}
	}

	return allErrs
}

func validateManagedDisksUpdate(oldDiskParams, newDiskParams *ManagedDiskParameters, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	fieldErrMsg := "changing managed disk options after machine creation is not allowed"

	if newDiskParams != nil && oldDiskParams != nil {
		if newDiskParams.StorageAccountType != oldDiskParams.StorageAccountType {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("storageAccountType"), newDiskParams, fieldErrMsg))
		}
		if newDiskParams.DiskEncryptionSet != nil && oldDiskParams.DiskEncryptionSet != nil {
			if newDiskParams.DiskEncryptionSet.ID != oldDiskParams.DiskEncryptionSet.ID {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("diskEncryptionSet").Child("ID"), newDiskParams, fieldErrMsg))
			}
		} else if (newDiskParams.DiskEncryptionSet != nil && oldDiskParams.DiskEncryptionSet == nil) || (newDiskParams.DiskEncryptionSet == nil && oldDiskParams.DiskEncryptionSet != nil) {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("diskEncryptionSet"), newDiskParams, fieldErrMsg))
		}
	} else if (newDiskParams != nil && oldDiskParams == nil) || (newDiskParams == nil && oldDiskParams != nil) {
		allErrs = append(allErrs, field.Invalid(fieldPath, newDiskParams, fieldErrMsg))
	}

	return allErrs
}

func validateStorageAccountType(storageAccountType string, fieldPath *field.Path, isOSDisk bool) field.ErrorList {
	allErrs := field.ErrorList{}

	if isOSDisk && storageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) {
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("managedDisks").Child("storageAccountType"), storageAccountType, "UltraSSD_LRS can only be used with data disks, it cannot be used with OS Disks"))
	}

	if storageAccountType == "" {
		allErrs = append(allErrs, field.Required(fieldPath, "the Storage Account Type for Managed Disk cannot be empty"))
		return allErrs
	}

	for _, possibleStorageAccountType := range armcompute.PossibleDiskStorageAccountTypesValues() {
		if string(possibleStorageAccountType) == storageAccountType {
			return allErrs
		}
	}
	allErrs = append(allErrs, field.Invalid(fieldPath, "", fmt.Sprintf("allowed values are %v", armcompute.PossibleDiskStorageAccountTypesValues())))
	return allErrs
}

func validateCachingType(cachingType string, fieldPath *field.Path, managedDisk *ManagedDiskParameters) field.ErrorList {
	allErrs := field.ErrorList{}
	cachingTypeChildPath := fieldPath.Child("CachingType")

	if managedDisk != nil && managedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) {
		if cachingType != string(armcompute.CachingTypesNone) {
			allErrs = append(allErrs, field.Invalid(cachingTypeChildPath, cachingType, fmt.Sprintf("cachingType '%s' is not supported when storageAccountType is '%s'. Allowed values are: '%s'", cachingType, armcompute.StorageAccountTypesUltraSSDLRS, armcompute.CachingTypesNone)))
		}
	}

	for _, possibleCachingType := range armcompute.PossibleCachingTypesValues() {
		if string(possibleCachingType) == cachingType {
			return allErrs
		}
	}

	allErrs = append(allErrs, field.Invalid(cachingTypeChildPath, cachingType, fmt.Sprintf("allowed values are %v", armcompute.PossibleCachingTypesValues())))
	return allErrs
}

// ValidateDiagnostics validates the Diagnostic spec.
func ValidateDiagnostics(diagnostics *Diagnostics, fieldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if diagnostics != nil && diagnostics.Boot != nil {
		switch diagnostics.Boot.StorageAccountType {
		case UserManagedDiagnosticsStorage:
			if diagnostics.Boot.UserManaged == nil {
				allErrs = append(allErrs, field.Required(fieldPath.Child("UserManaged"),
					fmt.Sprintf("userManaged must be specified when storageAccountType is '%s'", UserManagedDiagnosticsStorage)))
			} else if diagnostics.Boot.UserManaged.StorageAccountURI == "" {
				allErrs = append(allErrs, field.Required(fieldPath.Child("StorageAccountURI"),
					fmt.Sprintf("StorageAccountURI cannot be empty when storageAccountType is '%s'", UserManagedDiagnosticsStorage)))
			}
		case ManagedDiagnosticsStorage:
			if diagnostics.Boot.UserManaged != nil &&
				diagnostics.Boot.UserManaged.StorageAccountURI != "" {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("StorageAccountURI"), diagnostics.Boot.UserManaged.StorageAccountURI,
					fmt.Sprintf("StorageAccountURI cannot be set when storageAccountType is '%s'",
						ManagedDiagnosticsStorage)))
			}
		case DisabledDiagnosticsStorage:
			if diagnostics.Boot.UserManaged != nil &&
				diagnostics.Boot.UserManaged.StorageAccountURI != "" {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("StorageAccountURI"), diagnostics.Boot.UserManaged.StorageAccountURI,
					fmt.Sprintf("StorageAccountURI cannot be set when storageAccountType is '%s'",
						ManagedDiagnosticsStorage)))
			}
		}
	}

	return allErrs
}

// ValidateConfidentialCompute validates the configuration options when the machine is a Confidential VM.
// https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/create-or-update?tabs=HTTP#vmdisksecurityprofile
// https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/create-or-update?tabs=HTTP#securityencryptiontypes
// https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/create-or-update?tabs=HTTP#uefisettings
func ValidateConfidentialCompute(managedDisk *ManagedDiskParameters, profile *SecurityProfile, fieldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	var securityEncryptionType SecurityEncryptionType

	if managedDisk != nil && managedDisk.SecurityProfile != nil {
		securityEncryptionType = managedDisk.SecurityProfile.SecurityEncryptionType
	}

	if profile != nil && securityEncryptionType != "" {
		// SecurityEncryptionType can only be set for Confindential VMs
		if profile.SecurityType != SecurityTypesConfidentialVM {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("SecurityType"), profile.SecurityType,
				fmt.Sprintf("SecurityType should be set to '%s' when securityEncryptionType is defined", SecurityTypesConfidentialVM)))
		}

		// Confidential VMs require vTPM to be enabled, irrespective of the SecurityEncryptionType used
		if profile.UefiSettings == nil {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("UefiSettings"), profile.UefiSettings,
				"UefiSettings should be set when securityEncryptionType is defined"))
		}

		if profile.UefiSettings != nil && (profile.UefiSettings.VTpmEnabled == nil || !*profile.UefiSettings.VTpmEnabled) {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("VTpmEnabled"), profile.UefiSettings.VTpmEnabled,
				"VTpmEnabled should be set to true when securityEncryptionType is defined"))
		}

		if securityEncryptionType == SecurityEncryptionTypeDiskWithVMGuestState {
			// DiskWithVMGuestState encryption type is not compatible with EncryptionAtHost
			if profile.EncryptionAtHost != nil && *profile.EncryptionAtHost {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("EncryptionAtHost"), profile.EncryptionAtHost,
					fmt.Sprintf("EncryptionAtHost cannot be set to 'true' when securityEncryptionType is set to '%s'", SecurityEncryptionTypeDiskWithVMGuestState)))
			}

			// DiskWithVMGuestState encryption type requires SecureBoot to be enabled
			if profile.UefiSettings != nil && (profile.UefiSettings.SecureBootEnabled == nil || !*profile.UefiSettings.SecureBootEnabled) {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("SecureBootEnabled"), profile.UefiSettings.SecureBootEnabled,
					fmt.Sprintf("SecureBootEnabled should be set to true when securityEncryptionType is set to '%s'", SecurityEncryptionTypeDiskWithVMGuestState)))
			}
		}
	}

	return allErrs
}

// ValidateCapacityReservationGroupID validates the capacity reservation group id.
func ValidateCapacityReservationGroupID(capacityReservationGroupID *string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if capacityReservationGroupID != nil {
		if _, err := azureutil.ParseResourceID(*capacityReservationGroupID); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, capacityReservationGroupID, "must be a valid Azure resource ID"))
		}
	}

	return allErrs
}

// ValidateVMExtensions validates the VMExtensions spec.
func ValidateVMExtensions(disableExtensionOperations *bool, vmExtensions []VMExtension, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ptr.Deref(disableExtensionOperations, false) && len(vmExtensions) > 0 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("AzureMachineTemplate", "spec", "template", "spec", "vmExtensions"), "VMExtensions must be empty when DisableExtensionOperations is true"))
	}

	return allErrs
}
