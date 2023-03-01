package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	for _, ip := range p.ExternalDNS {
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, err.Error()))
		}
	}

	allErrs = append(allErrs, ValidateMachinePool(c.FeatureSet, p.DefaultMachinePlatform, "default", fldPath.Child("defaultMachinePlatform"))...)

	// Platform fields only allowed in TechPreviewNoUpgrade
	if c.FeatureSet != configv1.TechPreviewNoUpgrade {
		if c.OpenStack.LoadBalancer != nil {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("loadBalancer"), "load balancer is not supported in this feature set"))
		}
	}

	if c.OpenStack.LoadBalancer != nil {
		if !validateLoadBalancer(c.OpenStack.LoadBalancer.Type) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("loadBalancer", "type"), c.OpenStack.LoadBalancer.Type, "invalid load balancer type"))
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
