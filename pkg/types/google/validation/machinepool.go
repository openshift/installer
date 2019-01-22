package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/google"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *google.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Size < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("size"), p.Size, "Storage size must be positive"))
	}
	return allErrs
}
