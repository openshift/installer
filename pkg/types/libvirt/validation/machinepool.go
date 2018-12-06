package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/libvirt"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *libvirt.MachinePool, fldPath *field.Path) field.ErrorList {
	return field.ErrorList{}
}
