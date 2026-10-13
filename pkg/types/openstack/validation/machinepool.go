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

	if machinePool.RootVolume != nil {
		if len(machinePool.Zones) > 0 && len(machinePool.RootVolume.Zones) == 0 {
			errs = append(errs, field.Required(fldPath.Child("rootVolume").Child("zones"), "root volume availability zones must be specified when compute availability zones are specified"))
		}

		rootVolumeType := machinePool.RootVolume.DeprecatedType
		rootVolumeTypes := machinePool.RootVolume.Types
		typePath := fldPath.Child("rootVolume").Child("type")
		typesPath := fldPath.Child("rootVolume").Child("types")

		if rootVolumeType != "" && len(rootVolumeTypes) > 0 {
			errs = append(errs, field.Invalid(typePath, rootVolumeType, "Only one of type or types can be specified"))
			errs = append(errs, field.Invalid(typesPath, rootVolumeTypes, "Only one of type or types can be specified"))
		}

		if rootVolumeType == "" && len(rootVolumeTypes) == 0 {
			errs = append(errs, field.Invalid(typePath, rootVolumeType, "Either type or types must be specified"))
			errs = append(errs, field.Invalid(typesPath, rootVolumeTypes, "Either type or types must be specified"))
		}

		// When distributing the Root volumes across multiple failure domains, we suggest using multiple Storage types so they use a different backend.
		// Storage availability zones are purely cosmetic for now and can also be used to define where a volume should be created, however
		// we don't want to force a user to define a Storage availability zone when using multiple Storage types.
		// Therefore we decided to require as many Storage types as there are Compute availability zones, if there are multiple Compute availability zones
		// and more than one Storage type is defined.
		// Even if we support a single Storage type across multiple failure domains, we still allow doing it.
		// e.g. it would not make sense to have 3 Compute availability zones and 2 Storage types, because one of the Storage types would be used twice and
		// therefore the number of failure domains would not be 3 anymore.
		if machinePool.RootVolume.Types != nil {
			if computes, volumes := len(machinePool.Zones), len(machinePool.RootVolume.Types); computes > 1 && volumes > 1 && volumes != computes {
				errs = append(errs, field.Invalid(typesPath, rootVolumeTypes, "Compute and Storage availability zones in a MachinePool should have been validated to have equal length when more than one Storage type is defined"))
			}
		}
	}

	return errs
}
