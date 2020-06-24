package packet

import (
	packngo "github.com/packethost/packngo"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/packet"
	"github.com/openshift/installer/pkg/types/packet/validation"
)

// Validate executes packet specific validation
func Validate(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	packetPlatformPath := field.NewPath("platform", "packet")

	if ic.Platform.Packet == nil {
		return errors.New(field.Required(
			packetPlatformPath,
			"validation requires a Engine platform configuration").Error())
	}

	allErrs = append(
		allErrs,
		validation.ValidatePlatform(ic.Platform.Packet, packetPlatformPath)...)

	con, err := packngo.NewClient()
	if err != nil {
		return err
	}

	// TODO(displague) validate networks

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.Packet != nil {
		allErrs = append(
			allErrs,
			validateMachinePool(con, field.NewPath("controlPlane", "platform", "packet"), ic.ControlPlane.Platform.Packet)...)
	}
	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.Packet != nil {
			allErrs = append(
				allErrs,
				validateMachinePool(con, fldPath.Child("platform", "packet"), compute.Platform.Packet)...)
		}
	}

	return allErrs.ToAggregate()
}

func validateMachinePool(con *packngo.Client, child *field.Path, pool *packet.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}

// authenticated takes an packet platform and validates
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
