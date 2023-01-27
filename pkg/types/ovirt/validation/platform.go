package validation

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *ovirt.Platform, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	if err := validate.UUID(p.ClusterID); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ovirt_cluster_id"), p.ClusterID, err.Error()))
	}
	if err := validate.UUID(p.StorageDomainID); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ovirt_storage_domain_id"), p.StorageDomainID, err.Error()))
	}
	if p.VNICProfileID != "" {
		if err := validate.UUID(p.VNICProfileID); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("vnicProfileID"), p.VNICProfileID, err.Error()))
		}
	}
	if p.AffinityGroups != nil {
		allErrs = append(allErrs, validateAffinityGroupFields(p, fldPath.Child("affinityGroups"))...)
		allErrs = append(allErrs, validateAffinityGroupDuplicate(p.AffinityGroups)...)
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	// Platform fields only allowed in TechPreviewNoUpgrade
	if c.FeatureSet != configv1.TechPreviewNoUpgrade {
		if c.Ovirt.LoadBalancer != nil {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("loadBalancer"), "load balancer is not supported in this feature set"))
		}
	}

	if c.Ovirt.LoadBalancer != nil {
		if !validateLoadBalancer(c.Ovirt.LoadBalancer.Type) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("loadBalancer", "type"), c.Ovirt.LoadBalancer.Type, "invalid load balancer type"))
		}
	}

	return allErrs
}

func validateAffinityGroupFields(platform *ovirt.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, ag := range platform.AffinityGroups {
		if ag.Name == "" {
			allErrs = append(
				allErrs,
				field.Required(fldPath,
					fmt.Sprintf("Invalid affinity group %v: name must be not empty", ag.Name)))
		}
		if ag.Priority < 0 || ag.Priority > 5 {
			allErrs = append(
				allErrs,
				field.Invalid(fldPath, ag,
					fmt.Sprintf(
						"Invalid affinity group %v: priority value must be between 0-5 found priority %v",
						ag.Name,
						ag.Priority)))
		}
	}
	return allErrs
}

// validateAffinityGroupDuplicate checks that there is no duplicated affinity group with different fields
func validateAffinityGroupDuplicate(agList []ovirt.AffinityGroup) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, ag1 := range agList {
		for _, ag2 := range agList[i+1:] {
			if ag1.Name == ag2.Name {
				if ag1.Priority != ag2.Priority ||
					ag1.Description != ag2.Description ||
					ag1.Enforcing != ag2.Enforcing {
					allErrs = append(
						allErrs,
						&field.Error{
							Type: field.ErrorTypeDuplicate,
							BadValue: errors.Errorf("Error validating affinity groups: found same "+
								"affinity group defined twice with different fields %v anf %v", ag1, ag2)})
				}
			}
		}
	}
	return allErrs
}

// validateLoadBalancer returns an error if the load balancer is not valid.
func validateLoadBalancer(lbType configv1.PlatformLoadBalancerType) bool {
	switch lbType {
	case configv1.LoadBalancerTypeOpenShiftManagedDefault, configv1.LoadBalancerTypeUserManaged:
		return true
	default:
		return false
	}
}
