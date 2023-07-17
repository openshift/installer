package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
)

var validServerGroupPolicies = []string{
	string(openstack.SGPolicyAffinity),
	string(openstack.SGPolicyAntiAffinity),
	string(openstack.SGPolicySoftAffinity),
	string(openstack.SGPolicySoftAntiAffinity),
}

// ValidateMachinePool validates Control plane and Compute MachinePools
func ValidateMachinePool(_ *openstack.Platform, machinePool *openstack.MachinePool, role string, fldPath *field.Path) field.ErrorList {
	if machinePool == nil {
		return nil
	}
	var errs field.ErrorList
	switch machinePool.ServerGroupPolicy {
	case openstack.SGPolicyUnset, openstack.SGPolicyAffinity, openstack.SGPolicyAntiAffinity, openstack.SGPolicySoftAffinity, openstack.SGPolicySoftAntiAffinity:
	default:
		errs = append(errs, field.NotSupported(fldPath.Child("serverGroupPolicy"), machinePool.ServerGroupPolicy, validServerGroupPolicies))
	}

	if len(machinePool.Zones) > 0 && machinePool.RootVolume != nil && len(machinePool.RootVolume.Zones) == 0 {
		errs = append(errs, field.Required(fldPath.Child("rootVolume").Child("zones"), "root volume availability zones must be specified when compute availability zones are specified"))
	}

	return errs
}
