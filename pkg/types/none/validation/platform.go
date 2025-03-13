package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/common"
)

// ValidateFencingCredentials checks that the provided fencing credentials are valid.
func ValidateFencingCredentials(installConfig *types.InstallConfig) (errors field.ErrorList) {
	fldPath := field.NewPath("platform", "none")
	fencingCredentials := installConfig.Platform.None.FencingCredentials
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, common.ValidateUniqueAndRequiredFields(fencingCredentials, fldPath, func([]byte) bool { return false }, "fencingCredentials")...)
	allErrs = append(allErrs, common.ValidateTwoFencingCredentials(*installConfig.ControlPlane.Replicas, fencingCredentials, fldPath.Child("fencingCredentials"))...)

	return allErrs
}
