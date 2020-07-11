package baremetal

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal/validation"
)

// ValidateProvisioning performs platform validation specifically for any optional requirement
// to be called when the cluster creation takes place
func ValidateProvisioning(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	if ic.Platform.BareMetal == nil {
		return errors.New(field.Required(field.NewPath("platform", "baremetal"), "Baremetal validation requires a baremetal platform configuration").Error())
	}

	allErrs = append(allErrs, validation.ValidateProvisioning(ic.Platform.BareMetal, ic.Networking, field.NewPath("platform").Child("baremetal"))...)

	return allErrs.ToAggregate()
}
