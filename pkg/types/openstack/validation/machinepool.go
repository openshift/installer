package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *openstack.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.RootVolume.IOPS < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("rootVolume").Child("iops"), p.RootVolume.IOPS, "Root volume IOPS must be positive"))
	}
	if p.RootVolume.Size < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("rootVolume").Child("size"), p.RootVolume.IOPS, "Root volume size must be positive"))
	}
	return allErrs
}
