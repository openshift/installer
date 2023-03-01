package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/openstack"
)

var validServerGroupPolicies = []string{
	string(openstack.SGPolicyAffinity),
	string(openstack.SGPolicyAntiAffinity),
	string(openstack.SGPolicySoftAffinity),
	string(openstack.SGPolicySoftAntiAffinity),
}

// ValidateMachinePool validates Control plane and Compute MachinePools
func ValidateMachinePool(_ configv1.FeatureSet, machinePool *openstack.MachinePool, _ string, fldPath *field.Path) field.ErrorList {
	if machinePool == nil {
		return nil
	}
	var errs field.ErrorList
	switch machinePool.ServerGroupPolicy {
	case openstack.SGPolicyUnset, openstack.SGPolicyAffinity, openstack.SGPolicyAntiAffinity, openstack.SGPolicySoftAffinity, openstack.SGPolicySoftAntiAffinity:
	default:
		errs = append(errs, field.NotSupported(fldPath.Child("serverGroupPolicy"), machinePool.ServerGroupPolicy, validServerGroupPolicies))
	}
	return errs
}
