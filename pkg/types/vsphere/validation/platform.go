package validation

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
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

	// If all VIPs are empty, skip IP validation.  All VIPs are required to be defined together.
	if strings.Join([]string{p.APIVIP, p.IngressVIP}, "") != "" {
		allErrs = append(allErrs, validateVIPs(p, fldPath)...)
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

	allErrs = append(allErrs, validateVIPs(p, fldPath)...)
	return allErrs
}

// validateVIPs checks that all required VIPs are provided and are valid IP addresses.
func validateVIPs(p *vsphere.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(p.APIVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("apiVIP"), "must specify a VIP for the API"))
	} else if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	if len(p.IngressVIP) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("ingressVIP"), "must specify a VIP for Ingress"))
	} else if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
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
