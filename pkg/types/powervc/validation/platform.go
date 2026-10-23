package validation

import (
	"net/url"

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

	// The installer does not upload the RHCOS image on PowerVC, but the
	// machines still reference the Glance image name it would have
	// generated for a URL. Require the name of an existing image, such as
	// one created by the powervc-image tool. A value that parses as a URL is
	// rejected using the same test as rhcos.GenerateOpenStackImageName.
	if p.ClusterOSImage == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterOSImage"), "clusterOSImage must be the name of an existing image with powervc"))
	} else if _, err := url.ParseRequestURI(p.ClusterOSImage); err == nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), p.ClusterOSImage, "clusterOSImage must be the name of an existing image, not a URL, with powervc"))
	}

	return allErrs
}
