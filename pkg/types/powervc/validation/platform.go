package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervc"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *powervc.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// In the future, we will check for PowerVC specific install-config.yaml entries here.
	// Currently, we check for OpenStack configurations which we don't support.
	if p.ExternalNetwork != "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("externalNetwork"), p.ExternalNetwork, "Cannot set external network with powervc"))
	}

	// Without controlPlanePort the installer creates its own network and
	// later finds the machine subnet by Neutron tag. PowerVC does not let a
	// non-admin user tag Neutron resources, so that lookup cannot work and
	// the install fails late. Require an existing network up front. The
	// contents are validated by the OpenStack validation, since the PowerVC
	// platform is converted into the OpenStack platform.
	if p.ControlPlanePort == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("controlPlanePort"), "controlPlanePort is required with powervc"))
	}

	return allErrs
}
