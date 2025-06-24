package validation

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervc"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *powervc.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// @TODO
	logrus.Debugf("HAMZY: ValidatePlatform")

	return allErrs
}
