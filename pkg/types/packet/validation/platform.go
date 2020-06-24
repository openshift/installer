package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/packet"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *packet.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}
