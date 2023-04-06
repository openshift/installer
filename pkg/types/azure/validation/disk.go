package validation

import (
	"regexp"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/azure"
)

var (
	// RxDiskEncryptionSetID is a regular expression that validates a disk encryption set ID.
	RxDiskEncryptionSetID = regexp.MustCompile(`(?i)^/subscriptions/([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})/resourceGroups/([-a-zA-Z0-9_().]{0,89}[-a-zA-Z0-9_()])/providers/Microsoft\.Compute/diskEncryptionSets/([-a-zA-Z0-9_]{1,80})$`)

	// RxSubscriptionID is a regular expression that validates a subscription ID.
	RxSubscriptionID = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	// RxResourceGroup is a regular expression that validates a resource group.
	RxResourceGroup = regexp.MustCompile(`^[-a-zA-Z0-9_().]{0,89}[-a-zA-Z0-9_()]$`)

	// RxDiskEncryptionSetName is a regular expression that validates a disk encryption set name
	RxDiskEncryptionSetName = regexp.MustCompile(`^[-a-zA-Z0-9_]{1,80}$`)
)

// ValidateDiskEncryption checks that the specified disk encryption configuration is valid.
func ValidateDiskEncryption(p *azure.MachinePool, cloudName azure.CloudEnvironment, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	childFldPath := fldPath.Child("osDisk", "diskEncryptionSet")

	diskEncryptionSet := p.OSDisk.DiskEncryptionSet
	if diskEncryptionSet != nil && cloudName == azure.StackCloud {
		return append(allErrs, field.Invalid(childFldPath.Child("diskEncryptionSet"), diskEncryptionSet, "disk encryption sets are not supported on this platform"))
	}
	if diskEncryptionSet.SubscriptionID == "" {
		return append(allErrs, field.Required(childFldPath.Child("subscriptionID"), "subscription ID is required"))
	}
	if !RxSubscriptionID.MatchString(diskEncryptionSet.SubscriptionID) {
		return append(allErrs, field.Invalid(childFldPath.Child("subscriptionID"), diskEncryptionSet.SubscriptionID, "invalid subscription ID format"))
	}
	if !RxResourceGroup.MatchString(diskEncryptionSet.ResourceGroup) {
		return append(allErrs, field.Invalid(childFldPath.Child("resourceGroup"), diskEncryptionSet.ResourceGroup, "invalid resource group format"))
	}
	if !RxDiskEncryptionSetName.MatchString(diskEncryptionSet.Name) {
		return append(allErrs, field.Invalid(childFldPath.Child("diskEncryptionSetName"), diskEncryptionSet.Name, "invalid name format"))
	}

	return allErrs
}

// ValidateEncryptionAtHost checks that the encryption at host configuration is valid.
func ValidateEncryptionAtHost(p *azure.MachinePool, cloudName azure.CloudEnvironment, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	encryptionAtHost := p.EncryptionAtHost
	if encryptionAtHost == true && cloudName == azure.StackCloud {
		return append(allErrs, field.Invalid(fldPath.Child("encryptionAtHost"), encryptionAtHost, "encryption at host is not supported on this platform"))
	}

	return allErrs
}
