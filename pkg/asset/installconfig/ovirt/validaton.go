package ovirt

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/ovirt/validation"
	"github.com/openshift/installer/pkg/validate"
)

// Validate executes ovirt specific validation
func Validate(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	ovirtPlatformPath := field.NewPath("platform", "ovirt")

	if ic.Platform.Ovirt == nil {
		return errors.New(field.Required(
			ovirtPlatformPath,
			"oVirt validation requires a oVirt platform configuration").Error())
	}

	allErrs = append(
		allErrs,
		validation.ValidatePlatform(ic.Platform.Ovirt, ovirtPlatformPath)...)

	con, err := NewConnection()
	if err != nil {
		return err
	}
	defer con.Close()

	if err := validateVNICProfile(*ic.Ovirt, con); err != nil {
		allErrs = append(
			allErrs,
			field.Invalid(ovirtPlatformPath.Child("vnicProfileID"), ic.Ovirt.VNICProfileID, err.Error()))
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.Ovirt != nil {
		allErrs = append(
			allErrs,
			validateMachinePool(con, field.NewPath("controlPlane", "platform", "ovirt"), ic.ControlPlane.Platform.Ovirt)...)
	}
	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.Ovirt != nil {
			allErrs = append(
				allErrs,
				validateMachinePool(con, fldPath.Child("platform", "ovirt"), compute.Platform.Ovirt)...)
		}
	}

	return allErrs.ToAggregate()
}

func validateMachinePool(con *ovirtsdk.Connection, child *field.Path, pool *ovirt.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateInstanceTypeID(con, child, pool)...)
	return allErrs
}

func validateInstanceTypeID(con *ovirtsdk.Connection, child *field.Path, machinePool *ovirt.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}
	if machinePool.InstanceTypeID != "" {
		_, err := con.SystemService().InstanceTypesService().InstanceTypeService(machinePool.InstanceTypeID).Get().Send()
		if err != nil {
			allErrs = append(allErrs, field.NotFound(child.Child("instanceTypeID"), machinePool.InstanceTypeID))
		}
	}
	return allErrs
}

// authenticated takes an ovirt platform and validates
// its connection to the API by establishing
// the connection and authenticating successfully.
// The API connection is closed in the end and must leak
// or be reused in any way.
func authenticated(c *Config) survey.Validator {
	return func(val interface{}) error {
		connection, err := ovirtsdk.NewConnectionBuilder().
			URL(c.URL).
			Username(c.Username).
			Password(fmt.Sprint(val)).
			CAFile(c.CAFile).
			Insecure(c.Insecure).
			Build()

		if err != nil {
			return errors.Errorf("failed to construct connection to oVirt platform %s", err)
		}

		defer connection.Close()

		err = connection.Test()
		if err != nil {
			return errors.Errorf("failed to connect to oVirt platform %s", err)
		}
		return nil
	}

}

// validate the provided vnic profile exists and belongs the the cluster network
func validateVNICProfile(platform ovirt.Platform, con *ovirtsdk.Connection) error {
	if platform.VNICProfileID != "" {
		profiles, err := FetchVNICProfileByClusterNetwork(con, platform.ClusterID, platform.NetworkName)
		if err != nil {
			return err
		}

		for _, p := range profiles {
			if platform.VNICProfileID == p.MustId() {
				return nil
			}
		}

		return fmt.Errorf(
			"vNic profile ID %s does not belong to cluster network %s",
			platform.VNICProfileID,
			platform.NetworkName)
	}
	return nil
}

func validURL(val interface{}) error {
	uri, ok := val.(string)
	if !ok {
		return fmt.Errorf("cannot check url validity on type %T", val)
	}

	return validate.URIWithProtocol(uri, "https")
}
