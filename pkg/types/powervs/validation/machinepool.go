package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

    "github.com/openshift/installer/pkg/types/powervs"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *powervs.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
    /*
	for i, zone := range p.Zones {
		if !strings.HasPrefix(zone, platform.Region) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("zones").Index(i), zone, fmt.Sprintf("Zone not in configured region (%s)", platform.Region)))
		}
	}

	if p.OSDisk.DiskType != "" {
		diskTypes := sets.NewString("pd-standard", "pd-ssd")
		if !diskTypes.Has(p.OSDisk.DiskType) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskType"), p.OSDisk.DiskType, diskTypes.List()))
		}
	}*/

	return allErrs
}
