package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// Validate validates the given installconfig for OpenStack platform
func Validate(ic *types.InstallConfig) error {
	ci, err := newCloudInfo(ic)
	if err != nil {
		return err
	}

	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validatePlatform(ic, ci)...)
	allErrs = append(allErrs, validateMachinePool(ic.ControlPlane.Platform.OpenStack, field.NewPath("controlPlane", "platform", "openstack"))...)
	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.OpenStack != nil {
			allErrs = append(
				allErrs,
				validateMachinePool(compute.Platform.OpenStack, fldPath.Child("platform", "openstack"))...)
		}
	}

	return allErrs.ToAggregate()
}
