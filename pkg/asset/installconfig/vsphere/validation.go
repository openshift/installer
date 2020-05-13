package vsphere

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
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

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for vSphere clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	allErrs = append(allErrs, validation.ValidateForProvisioning(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, folderExists(ic, field.NewPath("platform").Child("vsphere").Child("folder"))...)

	return allErrs.ToAggregate()
}

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	cfg := ic.VSphere

	// If no folder is specified, skip this check as the folder will be created.
	if cfg.Folder == "" {
		return allErrs
	}

	vim25Client, _, err := vspheretypes.CreateVSphereClients(context.TODO(), cfg.VCenter, cfg.Username, cfg.Password)
	if err != nil {
		err = errors.Wrap(err, "unable to connect to vCenter API")
		return append(allErrs, field.InternalError(fldPath, err))
	}

	finder := find.NewFinder(vim25Client)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if _, err = finder.Folder(ctx, cfg.Folder); err != nil {
		return append(allErrs, field.Invalid(fldPath, cfg.Folder, err.Error()))
	}
	return nil
}
