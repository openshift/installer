package vsphere

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25"
	"k8s.io/apimachinery/pkg/util/validation/field"

	vsphereclient "github.com/openshift/installer/pkg/client/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/types/vsphere/validation"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	return validation.ValidatePlatform(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere")).ToAggregate()
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for vSphere clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	p := ic.Platform.VSphere
	vim25Client, _, cleanup, err := vsphereclient.CreateVSphereClients(context.TODO(),
		p.VCenter,
		p.Username,
		p.Password)

	if err != nil {
		return errors.New(field.InternalError(field.NewPath("platform", "vsphere"), errors.Wrapf(err, "unable to connect to vCenter %s.", p.VCenter)).Error())
	}
	defer cleanup()

	finder := vsphereclient.NewFinder(vim25Client)
	return validateProvisioning(vim25Client, finder, ic)
}

func validateProvisioning(client *vim25.Client, finder vsphereclient.Finder, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validation.ValidateForProvisioning(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, folderExists(finder, ic, field.NewPath("platform").Child("vsphere").Child("folder"))...)
	allErrs = append(allErrs, resourcePoolExists(finder, ic, field.NewPath("platform").Child("vsphere").Child("resourcePool"))...)
	if p := ic.Platform.VSphere; p.Network != "" {
		allErrs = append(allErrs, validateNetwork(client, finder, p, field.NewPath("platform").Child("vsphere").Child("network"))...)
	}

	return allErrs.ToAggregate()
}

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(finder vsphereclient.Finder, ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {
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

func validateNetwork(client *vim25.Client, finder vsphereclient.Finder, p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	// It's not possible to validate a network if datacenter or cluster are empty strings
	if p.Datacenter == "" || p.Cluster == "" {
		return nil
	}
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

	_, err = vsphereclient.GetNetworkMoID(ctx, client, finder, dataCenter.Name(), p.Cluster, p.Network)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, p.Network, err.Error())}
	}
	return nil
}

// resourcePoolExists returns an error if a resourcePool is specified in the vSphere platform but a resourcePool with that name is not found in the datacenter.
func resourcePoolExists(finder vsphereclient.Finder, ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {
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
