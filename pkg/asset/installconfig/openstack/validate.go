package openstack

import (
	"os"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/types"
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

	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validation.ValidatePlatform(ic.Platform.OpenStack, ic.Networking, ci)...)
	if ic.ControlPlane.Platform.OpenStack != nil {
		allErrs = append(allErrs, validation.ValidateMachinePool(ic.ControlPlane.Platform.OpenStack, ci, true, field.NewPath("controlPlane", "platform", "openstack"))...)
	}
	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.OpenStack != nil {
			allErrs = append(
				allErrs,
				validation.ValidateMachinePool(compute.Platform.OpenStack, ci, false, fldPath.Child("platform", "openstack"))...)
		}
	}

	return allErrs.ToAggregate()
}

// ValidateForProvisioning validates that the install config is valid for provisioning the cluster.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	if ic.ControlPlane.Replicas != nil && *ic.ControlPlane.Replicas != 3 {
		return field.Invalid(field.NewPath("controlPlane", "replicas"), ic.ControlPlane.Replicas, "control plane must be exactly three nodes when provisioning on OpenStack")
	}
	return nil
}
