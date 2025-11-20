package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervc"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *powervc.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// In the future, we will check for PowerVC specific install-config.yaml entries here.
	// Currently, we check for OpenStack configurations which we don't support.
	if p.ExternalNetwork != "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("externalNetwork"), p.ExternalNetwork, "Cannot set external network with powervc"))
	}

	return allErrs
}
