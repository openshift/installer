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

	errs = append(errs, validateSubnets(platform, machinePool, fldPath)...)

	return errs
}

func validateSubnets(platform *openstack.Platform, machinePool *openstack.MachinePool, fldPath *field.Path) (errs field.ErrorList) {
	subnetsNumber := len(machinePool.Subnets)
	if subnetsNumber > 0 && platform.MachinesSubnet != "" {
		errs = append(errs, field.Invalid(fldPath.Child("subnets"), machinePool.Subnets, "can't be used together with machinesSubnet"))
	}
	zonesNumber := len(machinePool.Zones)
	if subnetsNumber > 0 && zonesNumber > 0 && subnetsNumber != zonesNumber {
		errs = append(errs, field.Invalid(fldPath.Child("subnets"), subnetsNumber, "must match the number of zones"))
	}
	return errs
}
