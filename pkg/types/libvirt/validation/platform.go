package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *libvirt.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if err := validate.URI(p.URI); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("uri"), p.URI, err.Error()))
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	if p.Network.IfName == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("network").Child("if"), p.Network.IfName))
	}
	if err := validate.SubnetCIDR(&p.Network.IPRange.IPNet); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("network").Child("ipRange"), p.Network.IPRange, err.Error()))
	}
	return allErrs
}
