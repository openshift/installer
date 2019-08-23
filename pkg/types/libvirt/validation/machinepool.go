package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/libvirt"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *libvirt.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.DomainMemoryMiB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("DomainMemory"), p.DomainMemoryMiB, "Domain Memory (MiB) must be non-negative"))
	}
	if p.DomainVcpuCount < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("DomainVcpu"), p.DomainVcpuCount, "Domain VCPU count must be positive"))
	}
	return allErrs

}
