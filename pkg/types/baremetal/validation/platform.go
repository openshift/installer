package validation

import (
	"fmt"
	"net"
	"net/url"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/validate"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"strings"
)

// dynamicValidator is a function that validates certain fields in the platform.
type dynamicValidator func(*baremetal.Platform, *field.Path) field.ErrorList

// dynamicValidators is an array of dynamicValidator functions. This array can be added to by an init function, and
// is intended to be used for validations that require dependencies not built with the default tags, e.g. libvirt
// libraries.
var dynamicValidators []dynamicValidator

func validateIPinMachineCIDR(vip string, n *types.Networking) error {
	for _, network := range n.MachineNetwork {
		if network.CIDR.Contains(net.ParseIP(vip)) {
			return nil
		}
	}
	return fmt.Errorf("the virtual IP is expected to be in one of the machine networks")
}

func validateIPNotinMachineCIDR(ip string, n *types.Networking) error {
	for _, network := range n.MachineNetwork {
		if network.CIDR.Contains(net.ParseIP(ip)) {
			return fmt.Errorf("the IP must not be in one of the machine networks")
		}
	}
	return nil
}

func validateOSImageURI(uri string) error {
	// Check for valid URI and sha256 checksum part of the URL
	parsedURL, err := url.ParseRequestURI(uri)
	if err != nil {
		return fmt.Errorf("the URI provided: %s is invalid", uri)
	}
	if parsedURL.Scheme == "http" || parsedURL.Scheme == "https" {
		var sha256Checksum string
		if sha256Checksums, ok := parsedURL.Query()["sha256"]; ok {
			sha256Checksum = sha256Checksums[0]
		}
		if sha256Checksum == "" {
			return fmt.Errorf("the sha256 parameter in the %s URI is missing", uri)
		}
		if len(sha256Checksum) != 64 {
			return fmt.Errorf("the sha256 parameter in the %s URI is invalid", uri)
		}
	} else {
		return fmt.Errorf("the URI provided: %s must begin with http/https", uri)
	}
	return nil
}

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *baremetal.Platform, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if err := validate.URI(p.LibvirtURI); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("libvirtURI"), p.LibvirtURI, err.Error()))
	}

	if err := validate.IP(p.ClusterProvisioningIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningHostIP"), p.ClusterProvisioningIP, err.Error()))
	}

	if p.ProvisioningNetworkCIDR != nil && !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(p.ClusterProvisioningIP)) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterProvisioningIP"), p.ClusterProvisioningIP, fmt.Sprintf("%q is not in the provisioning network", p.ClusterProvisioningIP)))
	}

	if err := validate.IP(p.BootstrapProvisioningIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, err.Error()))
	}

	if p.ProvisioningNetworkCIDR != nil && !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(p.BootstrapProvisioningIP)) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapProvisioningIP"), p.BootstrapProvisioningIP, fmt.Sprintf("%q is not in the provisioning network", p.BootstrapProvisioningIP)))
	}

	if p.ProvisioningDHCPRange != "" {
		dhcpRange := strings.Split(p.ProvisioningDHCPRange, ",")
		if len(dhcpRange) != 2 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, "provisioning dhcp range should be in format: start_ip,end_ip"))
		} else {
			for _, ip := range dhcpRange {
				// Ensure IP is valid
				if err := validate.IP(ip); err != nil {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, fmt.Sprintf("%s: %s", ip, err.Error())))
				}

				// Validate IP is in the provisioning network
				if p.ProvisioningNetworkCIDR != nil && !p.ProvisioningNetworkCIDR.Contains(net.ParseIP(ip)) {
					allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningDHCPRange"), p.ProvisioningDHCPRange, fmt.Sprintf("%q is not in the provisioning network", ip)))
				}
			}
		}
	}

	// Make sure the provisioning interface is set.  Very little we can do to validate this  as it's not on this machine.
	if p.ProvisioningNetworkInterface == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningNetworkInterface"), p.ProvisioningNetworkInterface, "no provisioning network interface is configured, please set this value to be the interface on the provisioning network on your cluster's baremetal hosts"))
	}

	if p.Hosts == nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("hosts"), p.Hosts, "bare metal hosts are missing"))
	}

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if err := validateIPinMachineCIDR(p.APIVIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	if err := validateIPinMachineCIDR(p.IngressVIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	if err := validate.IP(p.DNSVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("dnsVIP"), p.DNSVIP, err.Error()))
	}

	if err := validateIPinMachineCIDR(p.DNSVIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("dnsVIP"), p.DNSVIP, err.Error()))
	}
	if err := validateIPNotinMachineCIDR(p.ClusterProvisioningIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("provisioningHostIP"), p.ClusterProvisioningIP, err.Error()))
	}
	if err := validateIPNotinMachineCIDR(p.BootstrapProvisioningIP, n); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapHostIP"), p.BootstrapProvisioningIP, err.Error()))
	}
	if p.BootstrapOSImage != "" {
		if err := validateOSImageURI(p.BootstrapOSImage); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("bootstrapOSImage"), p.BootstrapOSImage, err.Error()))
		}
	}
	if p.ClusterOSImage != "" {
		if err := validateOSImageURI(p.ClusterOSImage); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), p.ClusterOSImage, err.Error()))
		}
		if !strings.Contains(p.ClusterOSImage, "-openstack") {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), p.ClusterOSImage, "cluster baremetal image must contain '-openstack'"))
		}
		// The URL will contain a checksum after the '?' so we can just check for the suffix followed by a '?'.
		if !strings.Contains(p.ClusterOSImage, "qcow2.gz?") && !strings.Contains(p.ClusterOSImage, "qcow2.xz?") {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), p.ClusterOSImage, "cluster baremetal image must end in .qcow2.gz or .qcow2.xz"))
		}
	}

	for _, validator := range dynamicValidators {
		allErrs = append(allErrs, validator(p, fldPath)...)
	}

	return allErrs
}
