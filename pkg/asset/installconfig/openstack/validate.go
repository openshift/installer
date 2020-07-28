package openstack

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/types"
)

// Validate validates the given installconfig for OpenStack platform
func Validate(ic *types.InstallConfig) error {
	ci, err := validation.GetCloudInfo(ic)
	if err != nil {
		return err
	}

	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validation.ValidatePlatform(ic.Platform.OpenStack, ic.Networking, ci)...)
	if ic.ControlPlane.Platform.OpenStack != nil {
		allErrs = append(allErrs, validation.ValidateMachinePool(ic.ControlPlane.Platform.OpenStack, ci, field.NewPath("controlPlane", "platform", "openstack"))...)
	}
	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.OpenStack != nil {
			allErrs = append(
				allErrs,
				validation.ValidateMachinePool(compute.Platform.OpenStack, ci, fldPath.Child("platform", "openstack"))...)
		}
	}

	return allErrs.ToAggregate()
}
