package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *ovirt.Platform, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	return field.ErrorList{
		&field.Error{
			Type:     field.ErrorTypeForbidden,
			Field:    fldPath.String(),
			BadValue: "Unsupported platform",
			Detail:   "Platform oVirt is no longer supported",
		},
	}
}
