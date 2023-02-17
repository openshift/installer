package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

var validServerGroupPolicies = []string{
	string(openstack.SGPolicyAffinity),
	string(openstack.SGPolicyAntiAffinity),
	string(openstack.SGPolicySoftAffinity),
	string(openstack.SGPolicySoftAntiAffinity),
}

// ValidateMachinePool validates Control plane and Compute MachinePools
func ValidateMachinePool(installConfig *types.InstallConfig, machinePool *openstack.MachinePool, typ string, fldPath *field.Path) field.ErrorList {
	if machinePool == nil {
		return nil
	}
	var errs field.ErrorList
	switch machinePool.ServerGroupPolicy {
	case openstack.SGPolicyUnset, openstack.SGPolicyAffinity, openstack.SGPolicyAntiAffinity, openstack.SGPolicySoftAffinity, openstack.SGPolicySoftAntiAffinity:
	default:
		errs = append(errs, field.NotSupported(fldPath.Child("serverGroupPolicy"), machinePool.ServerGroupPolicy, validServerGroupPolicies))
	}

	if len(machinePool.FailureDomains) > 0 {
		errs = append(errs, validateFailureDomains(installConfig, machinePool, typ, fldPath)...)
	}

	return errs
}

func validateFailureDomains(installConfig *types.InstallConfig, machinePool *openstack.MachinePool, typ string, fldPath *field.Path) field.ErrorList {
	var errs field.ErrorList
	// We don't get installConfig when validating the default machine-pool.
	if installConfig != nil {
		if installConfig.FeatureSet != configv1.TechPreviewNoUpgrade {
			return append(errs, field.Forbidden(fldPath.Child("failureDomain"), "failureDomain is only available within the TechPreviewNoUpgrade featureSet"))
		}
	}

	if typ == "default" {
		return append(errs, field.Forbidden(fldPath.Child("failureDomain"), "failureDomain are not accepted in the default machine-pool. Use Control plane and Compute machine-pools."))
	}

	if len(machinePool.Zones) > 1 || len(machinePool.Zones) == 1 && machinePool.Zones[0] != "" {
		errs = append(errs, field.Forbidden(fldPath.Child("zones"), "set compute zones in failureDomains when using failure domains"))
	}
	if machinePool.RootVolume != nil && (len(machinePool.RootVolume.Zones) > 1 || len(machinePool.RootVolume.Zones) == 1 && machinePool.RootVolume.Zones[0] != "") {
		errs = append(errs, field.Forbidden(fldPath.Child("rootVolume", "zones"), "set storage zones in failureDomains when using failure domains"))
	}
	return errs
}
