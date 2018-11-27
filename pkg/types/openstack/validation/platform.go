package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	if err := validate.SubnetCIDR(&p.NetworkCIDRBlock.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("NetworkCIDRBlock"), p.NetworkCIDRBlock, err.Error()))
	}
	if err := validate.URI(p.BaseImage); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("baseImage"), p.BaseImage, err.Error()))
	}
	return allErrs
}
