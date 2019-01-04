package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.IOPS < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("iops"), p.IOPS, "Storage IOPS must be positive"))
	}
	if p.Size < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("size"), p.IOPS, "Storage size must be positive"))
	}
	return allErrs
}
