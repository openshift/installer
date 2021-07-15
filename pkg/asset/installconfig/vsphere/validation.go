package vsphere

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/types/vsphere/validation"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	p := ic.Platform.VSphere
	vim25Client, _, err := vspheretypes.CreateVSphereClients(context.TODO(),
		p.VCenter,
		p.Username,
		p.Password)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "vsphere"), errors.Wrapf(err, "unable to connect to vCenter %s.", p.VCenter)).Error())
	}

	finder := NewFinder(vim25Client)
	return validateResources(finder, ic)
}

func validateResources(finder Finder, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	p := ic.Platform.VSphere
	allErrs = append(allErrs, validation.ValidatePlatform(p, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, validateNetwork(finder, p, field.NewPath("platform").Child("vsphere").Child("network"))...)
	return allErrs.ToAggregate()
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for vSphere clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	p := ic.Platform.VSphere
	vim25Client, _, err := vspheretypes.CreateVSphereClients(context.TODO(),
		p.VCenter,
		p.Username,
		p.Password)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "vsphere"), errors.Wrapf(err, "unable to connect to vCenter %s.", p.VCenter)).Error())
	}

	finder := NewFinder(vim25Client)
	return validateProvisioning(finder, ic)
}

func validateProvisioning(finder Finder, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateForProvisioning(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, folderExists(finder, ic, field.NewPath("platform").Child("vsphere").Child("folder"))...)

	return allErrs.ToAggregate()
}

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(finder Finder, ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	cfg := ic.VSphere

	// If no folder is specified, skip this check as the folder will be created.
	if cfg.Folder == "" {
		return allErrs
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if _, err := finder.Folder(ctx, cfg.Folder); err != nil {
		return append(allErrs, field.Invalid(fldPath, cfg.Folder, err.Error()))
	}
	return nil
}

func validateNetwork(finder Finder, p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	dcName := p.Datacenter
	if !strings.HasPrefix(dcName, "/") && !strings.HasPrefix(dcName, "./") {
		dcName = "./" + dcName
	}

	dataCenter, err := finder.Datacenter(ctx, dcName)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, p.Datacenter, err.Error())}
	}
	networkPath := fmt.Sprintf("%s/network/%s", dataCenter.InventoryPath, p.Network)
	_, err = finder.Network(ctx, networkPath)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, p.Network, "unable to find network provided")}
	}
	return nil
}
