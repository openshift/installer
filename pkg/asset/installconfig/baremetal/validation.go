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

// ValidateStaticBootstrapNetworking ensures that both or neither of  BootstrapExternalStaticIP and BootstrapExternalStaticGateway are set
func ValidateStaticBootstrapNetworking(ic *types.InstallConfig) error {
	if ic.Platform.BareMetal == nil {
		return errors.New(field.Required(field.NewPath("platform", "baremetal"), "Baremetal validation requires a baremetal platform configuration").Error())
	}

	if ic.Platform.BareMetal.BootstrapExternalStaticIP != "" && ic.Platform.BareMetal.BootstrapExternalStaticGateway == "" {
		return errors.New(field.Required(field.NewPath("platform", "baremetal"), "You must specify a value for BootstrapExternalStaticGateway when BootstrapExternalStaticIP is set.").Error())
	}

	if ic.Platform.BareMetal.BootstrapExternalStaticGateway != "" && ic.Platform.BareMetal.BootstrapExternalStaticIP == "" {
		return errors.New(field.Required(field.NewPath("platform", "baremetal"), "You must specify a value for BootstrapExternalStaticIP when BootstrapExternalStaticGateway is set.").Error())
	}

	return nil
}
