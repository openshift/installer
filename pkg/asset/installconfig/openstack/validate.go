package openstack

import (
	"os"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
	"github.com/sirupsen/logrus"
)

// Validate validates the given installconfig for OpenStack platform
func Validate(ic *types.InstallConfig) error {
	if skip := os.Getenv("OPENSHIFT_INSTALL_SKIP_PREFLIGHT_VALIDATIONS"); skip == "1" {
		logrus.Warnf("OVERRIDE: pre-flight validation disabled.")
		return nil
	}

	ci, err := validation.GetCloudInfo(ic)
	if err != nil {
		return err
	}
	if ci == nil {
		logrus.Warnf("Empty OpenStack cloud info and therefore will skip pre-flight validation.")
		return nil
	}

	allErrs := field.ErrorList{}

	// Validate platform platform
	allErrs = append(allErrs, validation.ValidatePlatform(ic.Platform.OpenStack, ic.Networking, ci)...)

	// Validate control plane
	controlPlane := defaultOpenStackMachinePoolPlatform()
	controlPlane.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
	controlPlane.Set(ic.ControlPlane.Platform.OpenStack)
	if controlPlane.RootVolume != nil && controlPlane.RootVolume.Zones == nil {
		controlPlane.RootVolume.Zones = openstackdefaults.DefaultRootVolumeAZ()
	}
	allErrs = append(allErrs, validation.ValidateMachinePool(&controlPlane, ci, true, field.NewPath("controlPlane", "platform", "openstack"))...)

	// Validate computes
	for idx := range ic.Compute {
		compute := defaultOpenStackMachinePoolPlatform()
		compute.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		compute.Set(ic.Compute[idx].Platform.OpenStack)
		if compute.RootVolume != nil && compute.RootVolume.Zones == nil {
			compute.RootVolume.Zones = openstackdefaults.DefaultRootVolumeAZ()
		}
		fldPath := field.NewPath("compute").Index(idx)
		allErrs = append(allErrs, validation.ValidateMachinePool(&compute, ci, false, fldPath.Child("platform", "openstack"))...)
	}

	return allErrs.ToAggregate()
}

// ValidateForProvisioning validates that the install config is valid for provisioning the cluster.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	if ic.ControlPlane.Replicas != nil && *ic.ControlPlane.Replicas > 3 {
		return field.Invalid(field.NewPath("controlPlane", "replicas"), ic.ControlPlane.Replicas, "control plane cannot be more than three nodes when provisioning on OpenStack")
	}
	return nil
}

func defaultOpenStackMachinePoolPlatform() openstack.MachinePool {
	return openstack.MachinePool{
		Zones: []string{""},
	}
}
