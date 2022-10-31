package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/azure"
)

// ValidateDiskEncryption checks that the specified disk encryption configuration is valid.
func ValidateDiskEncryption(p *azure.MachinePool, cloudName azure.CloudEnvironment, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	diskEncryptionSet := p.OSDisk.DiskEncryptionSet

	// diskEncryptionSet is optional, so it's valid if it's nil
	if diskEncryptionSet == nil {
		return allErrs
	}

	if cloudName == azure.StackCloud {
		return append(allErrs, field.Invalid(fldPath, diskEncryptionSet, "disk encryption sets are not supported on this platform"))
	}

	if len(diskEncryptionSet.Name) <= 0 {
		return append(allErrs, field.Required(fldPath.Child("name"), "name is required when specifying a diskEncryptionSet"))
	}

	if len(diskEncryptionSet.ResourceGroup) <= 0 {
		return append(allErrs, field.Required(fldPath.Child("resourceGroup"), "resourceGroup is required when specifying a diskEncryptionSet"))
	}

	return allErrs
}

// ValidateEncryptionAtHost checks that the encryption at host configuration is valid.
func ValidateEncryptionAtHost(p *azure.MachinePool, cloudName azure.CloudEnvironment, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	encryptionAtHost := p.EncryptionAtHost
	if encryptionAtHost && cloudName == azure.StackCloud {
		return append(allErrs, field.Invalid(fldPath, encryptionAtHost, "encryption at host is not supported on this platform"))
	}

	return allErrs
}
