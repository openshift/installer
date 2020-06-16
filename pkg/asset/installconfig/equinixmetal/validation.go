package equinixmetal

import (
	packngo "github.com/packethost/packngo"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/equinixmetal"
	"github.com/openshift/installer/pkg/types/equinixmetal/validation"
)

// Validate executes Equinix Metal specific validation
func Validate(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	equinixmetalPlatformPath := field.NewPath("platform", "equinixmetal")

	if ic.Platform.EquinixMetal == nil {
		return errors.New(field.Required(
			equinixmetalPlatformPath,
			"validation requires a Engine platform configuration").Error())
	}

	allErrs = append(
		allErrs,
		validation.ValidatePlatform(ic.Platform.EquinixMetal, equinixmetalPlatformPath)...)

	con, err := packngo.NewClient()
	if err != nil {
		return err
	}

	// TODO(displague) validate networks

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.EquinixMetal != nil {
		allErrs = append(
			allErrs,
			validateMachinePool(con, field.NewPath("controlPlane", "platform", "equinixmetal"), ic.ControlPlane.Platform.EquinixMetal)...)
	}
	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.EquinixMetal != nil {
			allErrs = append(
				allErrs,
				validateMachinePool(con, fldPath.Child("platform", "equinixmetal"), compute.Platform.EquinixMetal)...)
		}
	}

	return allErrs.ToAggregate()
}

func validateMachinePool(con *packngo.Client, child *field.Path, pool *equinixmetal.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}

// authenticated takes an equinixmetal platform and validates
// its connection to the API by establishing
// the connection and authenticating successfully.
// The API connection is closed in the end and must leak
// or be reused in any way.
func authenticated(c *Config) survey.Validator {
	return func(val interface{}) error {
		_, err := packngo.NewClient()

		if err != nil {
			return errors.Errorf("failed to construct connection to Engine platform %s", err)
		}

		return nil
	}

}
