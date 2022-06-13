package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/nutanix"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *nutanix.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.DiskSizeGiB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGiB"), p.DiskSizeGiB, "storage disk size must be positive"))
	}
	if p.MemoryMiB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("memoryMiB"), p.MemoryMiB, "memory size must be positive"))
	}
	if p.NumCPUs < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cpus"), p.NumCPUs, "number of CPUs must be positive"))
	}
	if p.NumCoresPerSocket < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("coresPerSocket"), p.NumCoresPerSocket, "cores per socket must be positive"))
	}
	if p.NumCoresPerSocket >= 0 && p.NumCPUs >= 0 && p.NumCoresPerSocket > p.NumCPUs {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("coresPerSocket"), p.NumCoresPerSocket, "cores per socket must be less than number of CPUs"))
	}
	return allErrs
}
