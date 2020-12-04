package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/kubevirt"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *kubevirt.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.CPU < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cpu"), p.CPU, "CPU must be positive"))
	}

	if p.StorageSize != "" {
		storageQuantity, err := resource.ParseQuantity(p.StorageSize)

		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("storageSize"), p.StorageSize, fmt.Sprintf("Storage size must be of Quantity type format, error: %s", err)))
		} else if storageQuantity.Sign() != 1 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("storageSize"), p.StorageSize, "Storage size must be positive value"))
		}
	}

	if p.Memory != "" {
		memoryQuantity, err := resource.ParseQuantity(p.Memory)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.Memory, fmt.Sprintf("Memory must be of Quantity type format, error: %s", err)))
		} else if memoryQuantity.Sign() != 1 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.Memory, "Memory must be positive value"))
		}
	}

	return allErrs
}
