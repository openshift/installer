package validation

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *aws.Platform, p *aws.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, zone := range p.Zones {
		if !strings.HasPrefix(zone, platform.Region) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("zones").Index(i), zone, fmt.Sprintf("Zone not in configured region (%s)", platform.Region)))
		}
	}

	if p.IOPS < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("iops"), p.IOPS, "Storage IOPS must be positive"))
	}
	if p.Size < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("size"), p.IOPS, "Storage size must be positive"))
	}
	return allErrs
}
