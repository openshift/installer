package validation

import (
	"github.com/openshift/installer/pkg/types/azure"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *azure.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.OSDisk.DiskSizeGB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), p.OSDisk.DiskSizeGB, "Storage DiskSizeGB must be positive"))
	}
	return allErrs
}
