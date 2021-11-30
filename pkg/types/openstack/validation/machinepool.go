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
func ValidateMachinePool(platform *openstack.Platform, machinePool *openstack.MachinePool, poolName string, fldPath *field.Path) field.ErrorList {
	if poolName == "master" {
		return validateMasterMachinePool(machinePool, fldPath)
	}
	return validateWorkerMachinePool(machinePool, fldPath)
}

func validateMasterMachinePool(pool *openstack.MachinePool, fldPath *field.Path) field.ErrorList {
	var errs field.ErrorList
	switch pool.ServerGroupPolicy {
	case openstack.SGPolicyUnset, openstack.SGPolicyAffinity, openstack.SGPolicyAntiAffinity, openstack.SGPolicySoftAffinity, openstack.SGPolicySoftAntiAffinity:
	default:
		errs = append(errs, field.NotSupported(fldPath.Child("serverGroupPolicy"), pool.ServerGroupPolicy, validServerGroupPolicies))
	}
	return errs
}

func validateWorkerMachinePool(pool *openstack.MachinePool, fldPath *field.Path) field.ErrorList {
	var errs field.ErrorList
	if pool.ServerGroupPolicy != openstack.SGPolicyUnset {
		errs = append(errs, field.Invalid(fldPath.Child("serverGroupPolicy"), pool.ServerGroupPolicy, "server group policy cannot be set for compute machines"))
	}
	return errs
}

func validateDefaultMachinePool(pool *openstack.MachinePool, fldPath *field.Path) field.ErrorList {
	if pool == nil {
		return nil
	}
	var errs field.ErrorList
	if pool.ServerGroupPolicy != openstack.SGPolicyUnset {
		errs = append(errs, field.Invalid(fldPath.Child("serverGroupPolicy"), pool.ServerGroupPolicy, "server group policy cannot be set as default because compute machines do not support it"))
	}
	return errs
}
