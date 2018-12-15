package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/validate"
)

var (
	validRegionValues = func() []string {
		validValues := make([]string, len(aws.ValidRegions))
		i := 0
		for r := range aws.ValidRegions {
			validValues[i] = r
			i++
		}
		return validValues
	}()
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *aws.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if _, ok := aws.ValidRegions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, validRegionValues))
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	if p.VPCCIDRBlock != nil {
		if err := validate.SubnetCIDR(&p.VPCCIDRBlock.IPNet); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("vpcCIDRBlock"), p.VPCCIDRBlock, err.Error()))
		}
	}
	return allErrs
}
