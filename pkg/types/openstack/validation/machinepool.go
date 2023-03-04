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

	errs = append(errs, validateFailureDomains(machinePool, role, fldPath)...)

	return errs
}

func validateFailureDomains(machinePool *openstack.MachinePool, role string, fldPath *field.Path) (errs field.ErrorList) {
	if len(machinePool.FailureDomains) == 0 {
		return nil
	}

	fldPath = fldPath.Child("failureDomains")

	if role != "master" {
		return append(errs, field.Forbidden(fldPath, "failure domains can only be set on the master machine-pool"))
	}

	// No failure domains together with zones
	if len(machinePool.Zones) > 1 || len(machinePool.Zones) > 0 && machinePool.Zones[0] != "" {
		errs = append(errs, field.Forbidden(fldPath, "failure domains can not be set together with zones"))
	}

	// No failure domains together with root volume zones
	if machinePool.RootVolume != nil {
		if len(machinePool.RootVolume.Zones) > 1 || len(machinePool.RootVolume.Zones) > 0 && machinePool.RootVolume.Zones[0] != "" {
			errs = append(errs, field.Forbidden(fldPath, "failure domains can not be set together with rootVolume zones"))
		}
	}

	// portTarget IDs must be unique
	for i := range machinePool.FailureDomains {
		ids := make(map[string]struct{})
		for _, portTarget := range machinePool.FailureDomains[i].PortTargets {
			if _, ok := ids[portTarget.ID]; ok {
				errs = append(errs, field.Duplicate(fldPath.Index(i).Child("id"), portTarget.ID))
			}
			ids[portTarget.ID] = struct{}{}
		}
	}

	return errs
}
