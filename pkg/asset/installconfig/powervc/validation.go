package powervc

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// Validate executes platform specific validation.
func Validate(config *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if config.Platform.PowerVC == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("platform", "powervc"), "Power VC Validation requires a Power VC platform configuration."))
		return allErrs.ToAggregate()
	}
	if config.Platform.OpenStack == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("platform", "powervc"), "PowerVC should have also initialized OpenStack"))
		return allErrs.ToAggregate()
	}

	return allErrs.ToAggregate()
}
