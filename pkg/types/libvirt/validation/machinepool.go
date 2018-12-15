package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/validate"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *libvirt.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Image != "" {
		if err := validate.URI(p.Image); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("image"), p.Image, err.Error()))
		}
	}
	return allErrs
}
