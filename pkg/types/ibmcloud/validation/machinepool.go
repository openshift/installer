package validation

import (
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateMachinePool validates the MachinePool.
func ValidateMachinePool(mp *ibmcloud.MachinePool, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if mp.BootVolume != nil {
		allErrs = append(allErrs, validateBootVolume(mp.BootVolume, path.Child("bootVolume"))...)
	}
	return allErrs
}

func validateBootVolume(bv *ibmcloud.BootVolume, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if bv.EncryptionKey != "" {
		_, parseErr := crn.Parse(bv.EncryptionKey)
		if parseErr != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("encryptionKey"), bv.EncryptionKey, "encryptionKey is not a valid IBM CRN"))
		}
	}
	return allErrs
}
