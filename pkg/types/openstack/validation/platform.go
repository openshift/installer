package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
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

	allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, "default", fldPath.Child("defaultMachinePlatform"))...)

	if c.OpenStack.LoadBalancer != nil {
		if !validateLoadBalancer(c.OpenStack.LoadBalancer.Type) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("loadBalancer", "type"), c.OpenStack.LoadBalancer.Type, "invalid load balancer type"))
		}
	}

	if c.OpenStack.ControlPlanePort != nil {
		allErrs = append(allErrs, validateControlPlanePort(c, fldPath)...)
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

// validateControlPlanePort returns all the errors found when the control plane port is not valid.
func validateControlPlanePort(c *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	controlPlanePort := c.OpenStack.ControlPlanePort
	var allErrs field.ErrorList
	if len(controlPlanePort.FixedIPs) <= 2 {
		for _, fixedIP := range controlPlanePort.FixedIPs {
			if fixedIP.Subnet.ID != "" && !validation.ValidUUIDv4(fixedIP.Subnet.ID) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlanePort").Child("fixedIPs"), fixedIP.Subnet.ID, "invalid subnet ID"))
			}
		}
		if controlPlanePort.Network.ID != "" && !validation.ValidUUIDv4(controlPlanePort.Network.ID) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("controlPlanePort").Child("network"), controlPlanePort.Network.ID, "invalid network ID"))
		}
	} else {
		allErrs = append(allErrs, field.TooMany(fldPath.Child("fixedIPs"), len(controlPlanePort.FixedIPs), 2))
	}
	return allErrs
}
