package validation

import (
	"github.com/openshift/installer/pkg/types/openstack"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var validServerGroupPolicies = []string{
	string(openstack.SGPolicyAffinity),
	string(openstack.SGPolicyAntiAffinity),
	string(openstack.SGPolicySoftAffinity),
	string(openstack.SGPolicySoftAntiAffinity),
}

// ValidateMachinePool validates Control plane and Compute MachinePools
func ValidateMachinePool(platform *openstack.Platform, machinePool *openstack.MachinePool, _ string, fldPath *field.Path) field.ErrorList {
	if machinePool == nil {
		return nil
	}
	var errs field.ErrorList

	switch machinePool.ServerGroupPolicy {
	case openstack.SGPolicyUnset, openstack.SGPolicyAffinity, openstack.SGPolicyAntiAffinity, openstack.SGPolicySoftAffinity, openstack.SGPolicySoftAntiAffinity:
	default:
		errs = append(errs, field.NotSupported(fldPath.Child("serverGroupPolicy"), machinePool.ServerGroupPolicy, validServerGroupPolicies))
	}

	errs = append(errs, validateFailureDomainsMachinePool(platform, machinePool, fldPath)...)

	return errs
}

func validateFailureDomainsMachinePool(platform *openstack.Platform, machinePool *openstack.MachinePool, fldPath *field.Path) (errs field.ErrorList) {
	if machinePool.FailureDomainNames != nil {
		var foundOne = false
		for _, name := range machinePool.FailureDomainNames {
			for _, domain := range platform.FailureDomains {
				if domain.Name == name {
					foundOne = true
					break
				}
			}
			if !foundOne {
				errs = append(errs, field.Invalid(fldPath.Child("failureDomainNames"), name, "failure domain not found"))
			}
		}

		if machinePool.Zones != nil {
			errs = append(errs, field.Invalid(fldPath.Child("zones"), machinePool.Zones, "zones cannot be specified when failureDomainNames is specified"))
		}
		if machinePool.RootVolume != nil && machinePool.RootVolume.Zones != nil {
			errs = append(errs, field.Invalid(fldPath.Child("zones"), machinePool.RootVolume.Zones, "zones cannot be specified when failureDomainNames is specified"))
		}
	}
	return errs
}
