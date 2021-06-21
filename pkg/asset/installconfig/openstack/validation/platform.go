package validation

import (
	"errors"
	"fmt"
	"net/url"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, ci *CloudInfo) field.ErrorList {
	var allErrs field.ErrorList
	fldPath := field.NewPath("platform", "openstack")

	// validate BYO machinesSubnet usage
	allErrs = append(allErrs, validateMachinesSubnet(p, n, ci, fldPath)...)

	// validate the externalNetwork
	allErrs = append(allErrs, validateExternalNetwork(p, ci, fldPath)...)

	// validate floating ips
	allErrs = append(allErrs, validateFloatingIPs(p, ci, fldPath)...)

	// validate vips
	allErrs = append(allErrs, validateVIPs(p, ci, fldPath)...)

	// validate custom cluster os image
	allErrs = append(allErrs, validateClusterOSImage(p, ci, fldPath)...)

	return allErrs
}

// validateMachinesSubnet validates the machines subnet and enforces proper byo subnet usage and returns a list of all validation errors
func validateMachinesSubnet(p *openstack.Platform, n *types.Networking, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if p.MachinesSubnet != "" {
		if len(p.ExternalDNS) > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, "externalDNS is set, externalDNS is not supported when machinesSubnet is set"))
		}
		if ci.MachinesSubnet == nil {
			allErrs = append(allErrs, field.NotFound(fldPath.Child("machinesSubnet"), p.MachinesSubnet))
		} else if !validUUIDv4(p.MachinesSubnet) {
			allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), errors.New("invalid subnet ID")))
		} else {
			if n.MachineNetwork[0].CIDR.String() != ci.MachinesSubnet.CIDR {
				allErrs = append(allErrs, field.InternalError(fldPath.Child("machinesSubnet"), fmt.Errorf("the first CIDR in machineNetwork, %s, doesn't match the CIDR of the machineSubnet, %s", n.MachineNetwork[0].CIDR.String(), ci.MachinesSubnet.CIDR)))
			}
		}
	}

	return allErrs
}

// validateExternalNetwork validates the user's input for the externalNetwork and returns a list of all validation errors
func validateExternalNetwork(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	// Return an error if external network was specified in the install config, but hasn't been found
	if p.ExternalNetwork != "" && ci.ExternalNetwork == nil {
		allErrs = append(allErrs, field.NotFound(fldPath.Child("externalNetwork"), p.ExternalNetwork))
	}
	return allErrs
}

func validateFloatingIPs(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if p.APIFloatingIP != "" {
		if ci.APIFIP == nil {
			allErrs = append(allErrs, field.NotFound(fldPath.Child("apiFloatingIP"), p.APIFloatingIP))
		} else if ci.APIFIP.Status != "DOWN" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("apiFloatingIP"), p.APIFloatingIP, "Floating IP already in use"))
		} else if p.ExternalNetwork == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("apiFloatingIP"), p.APIFloatingIP, "Cannot set floating ips when external network not specified"))
		}
	}

	if p.IngressFloatingIP != "" {
		if ci.IngressFIP == nil {
			allErrs = append(allErrs, field.NotFound(fldPath.Child("ingressFloatingIP"), p.IngressFloatingIP))
		} else if ci.IngressFIP.Status != "DOWN" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressFloatingIP"), p.IngressFloatingIP, "Floating IP already in use"))
		} else if p.ExternalNetwork == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressFloatingIP"), p.IngressFloatingIP, "Cannot set floating ips when external network not specified"))
		}
		if p.APIFloatingIP != "" && p.APIFloatingIP == p.IngressFloatingIP {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressFloatingIP"), p.IngressFloatingIP, "ingressFloatingIP can not be the same as apiFloatingIP"))
		}
	}
	return allErrs
}

func validateVIPs(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if p.APIVIP != "" && p.IngressVIP != "" {
		if p.APIVIP == p.IngressVIP {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, "ingressVIP can not be the same as apiVIP"))
		}
	}
	return allErrs
}

// validateExternalNetwork validates the user's input for the clusterOSImage and returns a list of all validation errors
func validateClusterOSImage(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	if p.ClusterOSImage == "" {
		return
	}

	// For URLs we support only 'http(s)' and 'file' schemes
	if uri, err := url.ParseRequestURI(p.ClusterOSImage); err == nil {
		switch uri.Scheme {
		case "http", "https", "file":
		default:
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), p.ClusterOSImage, fmt.Sprintf("URL scheme should be either http(s) or file but it is '%v'", uri.Scheme)))
		}
		return
	}

	// Image should exist in OpenStack Glance
	if ci.OSImage == nil {
		allErrs = append(allErrs, field.NotFound(fldPath.Child("clusterOSImage"), p.ClusterOSImage))
		return allErrs
	}

	// Image should have "active" status
	if ci.OSImage.Status != images.ImageStatusActive {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), p.ClusterOSImage, fmt.Sprintf("OS image must be active but its status is '%s'", ci.OSImage.Status)))
	}

	return allErrs
}
