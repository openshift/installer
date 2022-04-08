package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
// TODO(nutanix): Revisit for further expanding the validation logic
func ValidatePlatform(p *nutanix.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(p.PrismCentral) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral"), "must specify the Prism Central"))
	}
	if len(p.PrismElementUUID) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElement"), "must specify the Prism Element"))
	}
	if len(p.Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("username"), "must specify the username"))
	}
	if len(p.Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("password"), "must specify the password"))
	}
	if len(p.Port) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("port"), "must specify the port"))
	}
	if len(p.SubnetUUID) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("subnet"), "must specify the subnet"))
	}
	if len(p.PrismCentral) != 0 {
		if err := validate.Host(p.PrismCentral); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("prismCentral"), p.PrismCentral, "must be the domain name or IP address of the Prism Central"))
		}
	}

	// If all VIPs are empty, skip IP validation.  All VIPs are required to be defined together.
	if p.APIVIP != "" || p.IngressVIP != "" {
		allErrs = append(allErrs, validateVIPs(p, fldPath)...)
	}

	return allErrs
}

// validateVIPs checks that all required VIPs are provided and are valid IP addresses.
func validateVIPs(p *nutanix.Platform, fldPath *field.Path) field.ErrorList {
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

	if p.APIVIP == p.IngressVIP {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, "IPs for both API and Ingress should not be the same"))
	}

	return allErrs
}
