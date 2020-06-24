package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/packet"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *packet.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	return allErrs
}
