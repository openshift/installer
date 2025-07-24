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

	if c.OpenStack.DNSRecordsType == configv1.DNSRecordsTypeExternal && (c.OpenStack.LoadBalancer == nil || c.OpenStack.LoadBalancer.Type != configv1.LoadBalancerTypeUserManaged) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("dnsRecordsType"), c.OpenStack.DNSRecordsType, "external DNS records can only be configured with user-managed loadbalancers"))
	}

	if controlPlanePort := c.OpenStack.ControlPlanePort; controlPlanePort != nil {
		allErrs = append(allErrs, validateControlPlanePort(controlPlanePort, fldPath.Child("controlPlanePort"))...)
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
func validateControlPlanePort(controlPlanePort *openstack.PortTarget, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if controlPlanePort.Network.ID != "" && !validation.ValidUUIDv4(controlPlanePort.Network.ID) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("network"), controlPlanePort.Network.ID, "invalid network ID: must be a UUIDv4"))
	}

	fixedIPsField := fldPath.Child("fixedIPs")

	switch l := len(controlPlanePort.FixedIPs); l {
	case 0:
		allErrs = append(allErrs, field.Required(fixedIPsField, "it is required to set a subnet filter to the controlPlanePort"))
	case 1, 2:
		for i, fixedIP := range controlPlanePort.FixedIPs {
			subnetField := fixedIPsField.Index(i).Child("subnet")
			if fixedIP.Subnet.ID != "" && !validation.ValidUUIDv4(fixedIP.Subnet.ID) {
				allErrs = append(allErrs, field.Invalid(subnetField.Child("id"), fixedIP.Subnet.ID, "invalid subnet ID: must be a UUIDv4"))
			}
			if fixedIP.Subnet.ID == "" && fixedIP.Subnet.Name == "" {
				allErrs = append(allErrs, field.Required(subnetField, "either ID or Name must be set on the subnet filter"))
			}
		}
	default:
		allErrs = append(allErrs, field.TooMany(fixedIPsField, l, 2))
	}

	return allErrs
}
