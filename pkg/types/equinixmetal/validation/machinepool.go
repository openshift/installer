package validation

import (
	"github.com/openshift/installer/pkg/types/equinixmetal"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *equinixmetal.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	return allErrs
}
