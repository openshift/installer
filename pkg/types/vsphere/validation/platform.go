package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	return field.ErrorList{}
}
