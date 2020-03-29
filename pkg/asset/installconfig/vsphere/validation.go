package vsphere

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere/validation"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	allErrs = append(allErrs, validation.ValidatePlatform(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)

	return allErrs.ToAggregate()
}
