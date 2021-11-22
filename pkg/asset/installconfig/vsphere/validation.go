package vsphere

import (
	"context"
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
	if errs := validation.ValidatePlatform(p, field.NewPath("platform").Child("vsphere")); len(errs) != 0 {
		return errs.ToAggregate()
	}

	vim25Client, _, err := vspheretypes.CreateVSphereClients(context.TODO(),
		p.VCenter,
		p.Username,
		p.Password)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "vsphere"), errors.Wrapf(err, "unable to connect to vCenter %s.", p.VCenter)).Error())
	}
	finder := vspheretypes.NewFinder(vim25Client)
	networkIDUtil := vspheretypes.NewNetworkUtil(vim25Client)
	return validateResources(finder, networkIDUtil, ic)
}

func validateResources(finder vspheretypes.Finder, networkIdentifier vspheretypes.NetworkIdentifier, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	p := ic.Platform.VSphere
	if p.Network != "" {
		allErrs = append(allErrs, validateNetwork(finder, networkIdentifier, p, field.NewPath("platform").Child("vsphere").Child("network"))...)
	}
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

	finder := vspheretypes.NewFinder(vim25Client)
	return validateProvisioning(finder, ic)
}

func validateProvisioning(finder vspheretypes.Finder, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateForProvisioning(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, folderExists(finder, ic, field.NewPath("platform").Child("vsphere").Child("folder"))...)
	allErrs = append(allErrs, resourcePoolExists(finder, ic, field.NewPath("platform").Child("vsphere").Child("resourcePool"))...)

	return allErrs.ToAggregate()
}

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(finder vspheretypes.Finder, ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {
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

func validateNetwork(finder vspheretypes.Finder, networkIdentifier vspheretypes.NetworkIdentifier, p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
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

	_, err = vspheretypes.GetNetworkMoID(ctx, networkIdentifier, finder, dataCenter.Name(), p.Cluster, p.Network)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, p.Network, err.Error())}
	}
	return nil
}

// resourcePoolExists returns an error if a resourcePool is specified in the vSphere platform but a resourcePool with that name is not found in the datacenter.
func resourcePoolExists(finder vspheretypes.Finder, ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	cfg := ic.VSphere

	// If no resourcePool is specified, skip this check as the root resourcePool will be used.
	if cfg.ResourcePool == "" {
		return allErrs
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if _, err := finder.ResourcePool(ctx, cfg.ResourcePool); err != nil {
		return append(allErrs, field.Invalid(fldPath, cfg.ResourcePool, err.Error()))
	}
	return nil
}
