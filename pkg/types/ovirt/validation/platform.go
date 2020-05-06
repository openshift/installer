package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *ovirt.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if err := validate.UUID(p.ClusterID); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ovirt_cluster_id"), p.ClusterID, err.Error()))
	}
	if err := validate.UUID(p.StorageDomainID); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ovirt_storage_domain_id"), p.StorageDomainID, err.Error()))
	}
	if err := validate.IP(p.APIVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("api_vip"), p.APIVIP, err.Error()))
	}
	if err := validate.IP(p.DNSVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("dns_vip"), p.DNSVIP, err.Error()))
	}
	if err := validate.IP(p.IngressVIP); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingress_vip"), p.IngressVIP, err.Error()))
	}
	if p.VNICProfileID != "" {
		if err := validate.UUID(p.VNICProfileID); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("vnicProfileID"), p.IngressVIP, err.Error()))
		}
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}
