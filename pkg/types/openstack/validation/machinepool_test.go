package validation

import (
	"fmt"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func withServerGroupPolicy(serverGroupPolicy string) func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) { mp.ServerGroupPolicy = openstack.ServerGroupPolicy(serverGroupPolicy) }
}

func withFailureDomain(failureDomain machinev1alpha1.FailureDomain) func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) { mp.FailureDomains = append(mp.FailureDomains, failureDomain) }
}

func withAValidFailureDomain() func(*openstack.MachinePool) {
	return withFailureDomain(machinev1alpha1.FailureDomain{})
}

func withAComputeZone() func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) { mp.Zones = []string{"thezone"} }
}

func withARootVolumeZone() func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) {
		mp.RootVolume = &openstack.RootVolume{
			Size:  1,
			Zones: []string{"thezone"},
			Type:  "thescentofadisk",
		}
	}
}

func testMachinePool(options ...func(*openstack.MachinePool)) *openstack.MachinePool {
	var mp openstack.MachinePool
	for _, apply := range options {
		apply(&mp)
	}
	return &mp
}

func withFeatureSet(featureSet configv1.FeatureSet) func(*types.InstallConfig) {
	return func(ic *types.InstallConfig) { ic.FeatureSet = featureSet }
}

func withFeatureSetTechPreviewNoUpgrade() func(*types.InstallConfig) {
	return withFeatureSet(configv1.TechPreviewNoUpgrade)
}

func testInstallConfig(options ...func(*types.InstallConfig)) *types.InstallConfig {
	var ic types.InstallConfig
	for _, apply := range options {
		apply(&ic)
	}
	return &ic
}

func TestValidateMachinePool(t *testing.T) {
	type checkFunc func(field.ErrorList) error
	check := func(fns ...checkFunc) []checkFunc { return fns }
	someErrorType := func(wantType field.ErrorType) checkFunc {
		return func(errs field.ErrorList) error {
			for _, err := range errs {
				if wantType == err.Type {
					return nil
				}
			}
			return fmt.Errorf("expected error type %q, not found", wantType)
		}
	}
	exactlyNErrors := func(want int) checkFunc {
		return func(errs field.ErrorList) error {
			if have := len(errs); want != have {
				errorMessages := make([]string, have)
				for i := range errs {
					errorMessages[i] = errs[i].Error()
				}
				return fmt.Errorf("expected %d errors, got %d: \n\t%s", want, have, strings.Join(errorMessages, "\n\t"))
			}
			return nil
		}
	}
	noError := exactlyNErrors(0)

	for _, tc := range [...]struct {
		name          string
		machinePool   *openstack.MachinePool
		installConfig *types.InstallConfig
		checks        []checkFunc
	}{
		{
			"empty",
			testMachinePool(),
			testInstallConfig(),
			check(noError),
		},
		{
			"with valid server group policy",
			testMachinePool(withServerGroupPolicy("anti-affinity")),
			testInstallConfig(),
			check(noError),
		},
		{
			"with invalid server group policy",
			testMachinePool(withServerGroupPolicy("anti-gravity")),
			testInstallConfig(),
			check(
				someErrorType(field.ErrorTypeNotSupported),
				exactlyNErrors(1),
			),
		},
		{
			"with valid server group policy",
			testMachinePool(withServerGroupPolicy("anti-affinity")),
			testInstallConfig(),
			check(noError),
		},
		{
			"with a failure domain",
			testMachinePool(withAValidFailureDomain()),
			testInstallConfig(withFeatureSetTechPreviewNoUpgrade()),
			check(noError),
		},
		{
			"with a failure domain, no TechPreview",
			testMachinePool(withAValidFailureDomain()),
			testInstallConfig(),
			check(
				someErrorType(field.ErrorTypeForbidden),
				exactlyNErrors(1),
			),
		},
		{
			"with a failure domain and a compute zone defined",
			testMachinePool(withAValidFailureDomain(), withAComputeZone()),
			testInstallConfig(withFeatureSetTechPreviewNoUpgrade()),
			check(
				someErrorType(field.ErrorTypeForbidden),
				exactlyNErrors(1),
			),
		},
		{
			"with a failure domain and a root volume zone defined",
			testMachinePool(withAValidFailureDomain(), withARootVolumeZone()),
			testInstallConfig(withFeatureSetTechPreviewNoUpgrade()),
			check(
				someErrorType(field.ErrorTypeForbidden),
				exactlyNErrors(1),
			),
		},
		{
			"with a failure domain, a compute zone and a root volume zone defined",
			testMachinePool(withAValidFailureDomain(), withAComputeZone(), withARootVolumeZone()),
			testInstallConfig(withFeatureSetTechPreviewNoUpgrade()),
			check(
				someErrorType(field.ErrorTypeForbidden),
				exactlyNErrors(2),
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateMachinePool(tc.installConfig, tc.machinePool, "", nil)
			for _, check := range tc.checks {
				if e := check(errs); e != nil {
					t.Error(e)
				}
			}
		})
	}
}
