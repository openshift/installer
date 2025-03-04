package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/common"
)

// ValidateFencingCredentials checks that the provided fencing credentials are valid.
func ValidateFencingCredentials(fencingCredentials []*common.FencingCredential, fldPath *field.Path) (errors field.ErrorList) {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, common.ValidateUniqueAndRequiredFields(fencingCredentials, fldPath, func([]byte) bool { return false }, "fencingCredentials")...)
	allErrs = append(allErrs, common.ValidateTwoFencingCredentials(fencingCredentials, fldPath.Child("FencingCredentials"))...)

	return allErrs
}
