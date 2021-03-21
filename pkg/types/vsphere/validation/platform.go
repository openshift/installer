package validation

import (
	"fmt"
	"net"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *vsphere.Platform, network *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(p.VCenter) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("vCenter"), "must specify the name of the vCenter"))
	}
	if len(p.Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("username"), "must specify the username"))
	}
	if len(p.Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("password"), "must specify the password"))
	}
	if len(p.Datacenter) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("datacenter"), "must specify the datacenter"))
	}
	if len(p.DefaultDatastore) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("defaultDatastore"), "must specify the default datastore"))
	}

	if len(p.VCenter) != 0 {
		if err := validate.Host(p.VCenter); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("vCenter"), p.VCenter, "must be the domain name or IP address of the vCenter"))
		}
	}

	// If all VIPs are empty, skip IP validation.  All VIPs are required to be defined together.
	if strings.Join([]string{p.APIVIP, p.IngressVIP}, "") != "" {
		allErrs = append(allErrs, validateVIPs(p, network, fldPath)...)
	}

	// folder is optional, but if provided should pass validation
	if len(p.Folder) != 0 {
		allErrs = append(allErrs, validateFolder(p, fldPath)...)
	}

	return allErrs
}

// ValidateForProvisioning checks that the specified platform is valid.
func ValidateForProvisioning(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.Cluster) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("cluster"), "must specify the cluster"))
	}

	if len(p.Network) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("network"), "must specify the network"))
	}

	if len(p.APIVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("apiVIP"), "must specify a VIP for API"))
	}

	if len(p.IngressVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("ingressVIP"), "must specify a VIP for Ingress"))
	}
	return allErrs
}

// ipInNetwork return true if the given ip is within one of the machine networks.
func ipInNetwork(vipIP net.IP, machineNetwork []types.MachineNetworkEntry) bool {
	for _, machine := range machineNetwork {
		if machine.CIDR.Contains(vipIP) {
			return true
		}
	}
	return false
}

// validateVIPs checks that all required VIPs are provided and are valid IP addresses.
func validateVIPs(p *vsphere.Platform, network *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if network == nil {
		return append(allErrs, field.Invalid(field.NewPath("networking"), network, "must specify the machine networks"))
	}

	if len(p.APIVIP) != 0 {
		ip := net.ParseIP(p.APIVIP)
		if ip != nil {
			if !ipInNetwork(ip, network.MachineNetwork) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, "must be contained within one of the machine networks"))
			}
		} else {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, "not a valid IP address"))
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("apiVIP"), "must specify a VIP for both API and Ingress VIPs when specifying either"))
	}

	if len(p.IngressVIP) != 0 {
		ip := net.ParseIP(p.IngressVIP)
		if ip != nil {
			if !ipInNetwork(ip, network.MachineNetwork) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, "must be contained within one of the machine networks"))
			}
		} else {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, "not a valid IP address"))
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("ingressVIP"), "must specify a VIP for both API and Ingress VIPs when specifying either"))
	}

	if len(p.APIVIP) != 0 && len(p.IngressVIP) != 0 && p.APIVIP == p.IngressVIP {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, "IPs for both API and Ingress should not be the same."))
	}
	return allErrs
}

// validateFolder checks that a provided folder is in absolute path in the correct datacenter.
func validateFolder(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	dc := p.Datacenter
	if len(dc) == 0 {
		dc = "<datacenter>"
	}
	expectedPrefix := fmt.Sprintf("/%s/vm/", dc)

	if !strings.HasPrefix(p.Folder, expectedPrefix) {
		errMsg := fmt.Sprintf("folder must be absolute path: expected prefix %s", expectedPrefix)
		allErrs = append(allErrs, field.Invalid(fldPath.Child("folder"), p.Folder, errMsg))
	}

	return allErrs
}
