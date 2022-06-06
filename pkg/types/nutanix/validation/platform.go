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

	if len(p.PrismCentral.Endpoint.Address) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("endpoint").Child("address"),
			"must specify the Prism Central endpoint address"))
	} else {
		if err := validate.Host(p.PrismCentral.Endpoint.Address); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("prismCentral").Child("endpoint").Child("address"),
				p.PrismCentral.Endpoint.Address, "must be the domain name or IP address of the Prism Central"))
		}
	}

	if p.PrismCentral.Endpoint.Port < 1 || p.PrismCentral.Endpoint.Port > 65535 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("prismCentral").Child("endpoint").Child("port"),
			p.PrismCentral.Endpoint.Port, "The Prism Central endpoint port is invalid, must be in the range of 1 to 65535"))
	}

	if len(p.PrismCentral.Username) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("username"),
			"must specify the Prism Central username"))
	}

	if len(p.PrismCentral.Password) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismCentral").Child("password"),
			"must specify the Prism Central password"))
	}

	// Currently we only support one Prism Element for an OpenShift cluster
	if len(p.PrismElements) != 1 {
		allErrs = append(allErrs, field.Required(fldPath.Child("prismElements"), "must specify one Prism Element"))
	}

	for _, pe := range p.PrismElements {
		if len(pe.UUID) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("prismElements").Child("uuid"),
				"must specify the Prism Element UUID"))
		}

		if len(pe.Endpoint.Address) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("prismElements").Child("endpoint").Child("address"),
				"must specify the Prism Element endpoint address"))
		} else {
			if err := validate.Host(pe.Endpoint.Address); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("prismElements").Child("endpoint").Child("address"),
					pe.Endpoint.Address, "must be the domain name or IP address of the Prism Element (cluster)"))
			}
		}

		if pe.Endpoint.Port < 1 || pe.Endpoint.Port > 65535 {
			allErrs = append(allErrs, field.Required(fldPath.Child("prismElements").Child("endpoint").Child("port"),
				"The Prism Element endpoint port is invalid, must be in the range of 1 to 65535"))
		}
	}

	// Currently we only support one subnet for an OpenShift cluster
	if len(p.SubnetUUIDs) != 1 || len(p.SubnetUUIDs[0]) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("subnet"), "must specify the subnet"))
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
