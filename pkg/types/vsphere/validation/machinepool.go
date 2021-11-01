package validation

import (
	"github.com/openshift/installer/pkg/validate"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *vsphere.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.DiskSizeGB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), p.DiskSizeGB, "storage disk size must be positive"))
	}
	if p.MemoryMiB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("memoryMB"), p.MemoryMiB, "memory size must be positive"))
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
	if len(p.Zones) > 0 {
		for _, zone := range p.Zones {
			err := validate.ClusterName1035(zone)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("zones"), p.Zones, err.Error()))
			}
		}
	}
	return allErrs
}
