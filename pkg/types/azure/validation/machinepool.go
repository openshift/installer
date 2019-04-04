package validation

import (
	"github.com/openshift/installer/pkg/types/azure"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *azure.MachinePool, fldPath *field.Path) field.ErrorList {
	//TODO: implement machine pool validation
	allErrs := field.ErrorList{}
	return allErrs
}
