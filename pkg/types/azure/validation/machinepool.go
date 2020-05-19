package validation

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *azure.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.OSDisk.DiskSizeGB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskSizeGB"), p.OSDisk.DiskSizeGB, "Storage DiskSizeGB must be positive"))
	}

	if p.OSDisk.DiskType != "" {
		diskTypes := sets.NewString("Standard_LRS",
			"StandardSSD_LRS",
			// "UltraSSD_LRS" needs azure terraform version 2.0
			"Premium_LRS")

		if !diskTypes.Has(p.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.OSDisk.DiskType, diskTypes.List()))
		}
	}

	return allErrs
}

// ValidateMasterDiskType checks that the specified disk type is valid for control plane.
func ValidateMasterDiskType(p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Name == "master" && p.Platform.Azure.OSDisk.DiskType == "Standard_LRS" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("diskType"), p.Platform.Azure.OSDisk.DiskType, fmt.Sprintf("%s not compatible with control planes.", p.Platform.Azure.OSDisk.DiskType)))
	}

	return allErrs
}

// ValidateDefaultDiskType checks that the specified disk type is valid for default Azure Machine Platform.
func ValidateDefaultDiskType(p *azure.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.OSDisk.DiskType != "" {
		diskTypes := sets.NewString("StandardSSD_LRS", "Premium_LRS") // "UltraSSD_LRS"

		if !diskTypes.Has(p.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.OSDisk.DiskType, diskTypes.List()))
		}
	}

	return allErrs
}
