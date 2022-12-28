package validation

import (
	"bytes"
	"fmt"
	"net"
	"net/url"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilsslice "k8s.io/utils/strings/slices"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, ci *CloudInfo) field.ErrorList {
	var allErrs field.ErrorList
	fldPath := field.NewPath("platform", "openstack")

	// validate BYO controlPlanePort usage
	if p.ControlPlanePort != nil {
		allErrs = append(allErrs, validateControlPlanePort(p, n, ci, fldPath)...)
	}

	// validate the externalNetwork
	allErrs = append(allErrs, validateExternalNetwork(p, ci, fldPath)...)

	// validate floating ips
	allErrs = append(allErrs, validateFloatingIPs(p, ci, fldPath)...)

	// validate vips (on OpenStack we need some additional checks)
	allErrs = append(allErrs, validateVIPs(p, ci, fldPath)...)

	// validate custom cluster os image
	allErrs = append(allErrs, validateClusterOSImage(p, ci, fldPath)...)

	return allErrs
}

// validateControlPlanePort validates the machines subnets and network, while enforcing proper byo subnets usage and returns a list of all validation errors.
func validateControlPlanePort(p *openstack.Platform, n *types.Networking, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	networkID := ""
	hasIPv4Subnet := false
	hasIPv6Subnet := false
	if len(p.ExternalDNS) > 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, "externalDNS is set, externalDNS is not supported when ControlPlanePort is set"))
		return allErrs
	}
	networkCIDRs := networksCIDRs(n.MachineNetwork)
	for _, fixedIP := range p.ControlPlanePort.FixedIPs {
		subnet := getSubnet(ci.ControlPlanePortSubnets, fixedIP.Subnet.ID, fixedIP.Subnet.Name)
		if subnet == nil {
			subnetDetail := fixedIP.Subnet.ID
			if subnetDetail == "" {
				subnetDetail = fixedIP.Subnet.Name
			}
			allErrs = append(allErrs, field.NotFound(fldPath.Child("controlPlanePort").Child("fixedIPs"), subnetDetail))
		} else {
			if subnet.IPVersion == 6 {
				hasIPv6Subnet = true
			} else {
				hasIPv4Subnet = true
			}
			if !utilsslice.Contains(networkCIDRs, subnet.CIDR) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlanePort").Child("fixedIPs"), subnet.CIDR, "controlPlanePort CIDR does not match machineNetwork"))
			}
			if networkID != "" && networkID != subnet.NetworkID {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlanePort").Child("fixedIPs"), subnet.NetworkID, "fixedIPs subnets must be on the same Network"))
			}
			networkID = subnet.NetworkID
		}
	}
	if !hasIPv4Subnet && hasIPv6Subnet {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("controlPlanePort").Child("fixedIPs"), fmt.Errorf("one IPv4 subnet must be specified")))
	} else if hasIPv4Subnet && !hasIPv6Subnet && len(p.ControlPlanePort.FixedIPs) == 2 {
		allErrs = append(allErrs, field.InternalError(fldPath.Child("controlPlanePort").Child("fixedIPs"), fmt.Errorf("multiple IPv4 subnets is not supported")))
	}
	controlPlaneNetwork := p.ControlPlanePort.Network
	if controlPlaneNetwork.ID != "" || controlPlaneNetwork.Name != "" {
		networkDetail := controlPlaneNetwork.ID
		if networkDetail == "" {
			networkDetail = controlPlaneNetwork.Name
		}
		// check if the networks does not exist. If it does, verifies if the network contains the subnets
		if ci.ControlPlanePortNetwork == nil {
			allErrs = append(allErrs, field.NotFound(fldPath.Child("controlPlanePort").Child("network"), networkDetail))
		} else if ci.ControlPlanePortNetwork.ID != networkID {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlanePort").Child("network"), networkDetail, "network must contain subnets"))
		}
	}
	return allErrs
}

func networksCIDRs(machineNetwork []types.MachineNetworkEntry) []string {
	networks := make([]string, 0, len(machineNetwork))
	for _, network := range machineNetwork {
		networks = append(networks, network.CIDR.String())
	}
	return networks
}

func getSubnet(controlPlaneSubnets []*subnets.Subnet, subnetID, subnetName string) *subnets.Subnet {
	for _, subnet := range controlPlaneSubnets {
		if subnet.ID == subnetID {
			return subnet
		} else if subnet.Name != "" && subnet.Name == subnetName {
			return subnet
		}
	}
	return nil
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

// validateVIPs adds some OpenStack specific VIP validation. The universal
// platform VIP validation is done in pkg/types/validation/installconfig.go,
// validateAPIAndIngressVIPs().
func validateVIPs(p *openstack.Platform, ci *CloudInfo, fldPath *field.Path) (allErrs field.ErrorList) {
	// If the subnet is not found in the CloudInfo object, abandon validation.
	// For dual-stack the user needs to pre-create the Port for API and Ingress, so no need for validation.
	if len(ci.ControlPlanePortSubnets) == 1 {
		for _, allocationPool := range ci.ControlPlanePortSubnets[0].AllocationPools {
			start := net.ParseIP(allocationPool.Start)
			end := net.ParseIP(allocationPool.End)

			// If the allocation pool is undefined, abandon validation
			if start == nil || end == nil {
				continue
			}

			for _, apiVIPString := range p.APIVIPs {
				apiVIP := net.ParseIP(apiVIPString)
				if bytes.Compare(start, apiVIP) <= 0 && bytes.Compare(end, apiVIP) >= 0 {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIPs"), apiVIPString, "apiVIP can not fall in a MachineNetwork allocation pool"))
				}

			}

			for _, ingressVIPString := range p.IngressVIPs {
				ingressVIP := net.ParseIP(ingressVIPString)
				if bytes.Compare(start, ingressVIP) <= 0 && bytes.Compare(end, ingressVIP) >= 0 {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIPs"), ingressVIPString, "ingressVIP can not fall in a MachineNetwork allocation pool"))
				}
			}
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
