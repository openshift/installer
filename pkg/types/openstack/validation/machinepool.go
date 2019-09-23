package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *openstack.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate Root Volumes
	if p.RootVolume != nil {
		if p.RootVolume.Type == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("rootVolume").Child("type"), p.RootVolume.Type, "Volume type must be specified to use root volumes"))
		}
		if p.RootVolume.Size <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("rootVolume").Child("size"), p.RootVolume.Size, "Volume size must be greater than zero to use root volumes"))
		}
	}

	return allErrs
}
