package validation

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/openstack"
)

func withServerGroupPolicy(serverGroupPolicy string) func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) { mp.ServerGroupPolicy = openstack.ServerGroupPolicy(serverGroupPolicy) }
}

func withRootVolume(rootVolume *openstack.RootVolume) func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) { mp.RootVolume = rootVolume }
}

func withAvailabilityZone(zones []string) func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) { mp.Zones = zones }
}

func testMachinePool(options ...func(*openstack.MachinePool)) *openstack.MachinePool {
	var mp openstack.MachinePool
	for _, apply := range options {
		apply(&mp)
	}
	return &mp
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
				return fmt.Errorf("expected %d errors, got %d: %v", want, have, errs)
			}
			return nil
		}
	}
	noError := exactlyNErrors(0)

	for _, tc := range [...]struct {
		name        string
		machinePool *openstack.MachinePool
		role        string
		checks      []checkFunc
	}{
		{
			"empty",
			testMachinePool(),
			"default",
			check(noError),
		},
		{
			"with valid server group policy",
			testMachinePool(withServerGroupPolicy("anti-affinity")),
			"default",
			check(noError),
		},
		{
			"with invalid server group policy",
			testMachinePool(withServerGroupPolicy("anti-gravity")),
			"default",
			check(
				someErrorType(field.ErrorTypeNotSupported),
				exactlyNErrors(1),
			),
		},
		{
			"with no rootVolume type nor types",
			testMachinePool(
				withRootVolume(&openstack.RootVolume{
					Size: 10,
				}),
			),
			"default",
			check(
				someErrorType(field.ErrorTypeInvalid),
				exactlyNErrors(2),
			),
		},
		{
			"with both rootVolume type and types",
			testMachinePool(
				withRootVolume(&openstack.RootVolume{
					Size:           10,
					DeprecatedType: "fast",
					Types:          []string{"fast"},
				}),
			),
			"default",
			check(
				someErrorType(field.ErrorTypeInvalid),
				exactlyNErrors(2),
			),
		},
		{
			"with three compute zones and one root volume type",
			testMachinePool(
				withRootVolume(&openstack.RootVolume{
					Size:  10,
					Types: []string{"fast"},
					Zones: []string{"az1", "az2", "az3"},
				}),
				withAvailabilityZone([]string{"az1", "az2", "az3"}),
			),
			"default",
			check(noError),
		},
		{
			"with three compute zones and two root volume types",
			testMachinePool(
				withRootVolume(&openstack.RootVolume{
					Size:  10,
					Types: []string{"fast", "slow"},
					Zones: []string{"az1", "az2", "az3"},
				}),
				withAvailabilityZone([]string{"az1", "az2", "az3"}),
			),
			"default",
			check(
				someErrorType(field.ErrorTypeInvalid),
				exactlyNErrors(1),
			),
		},
		{
			"with three compute zones and invalid root volume missing zones",
			testMachinePool(withAvailabilityZone([]string{"az0", "az1", "az2"}), withRootVolume(&openstack.RootVolume{Size: 100, Types: []string{"fast"}})),
			"default",
			check(
				someErrorType(field.ErrorTypeRequired),
				exactlyNErrors(1),
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateMachinePool(nil, tc.machinePool, tc.role, nil)
			for _, check := range tc.checks {
				if e := check(errs); e != nil {
					t.Error(e)
				}
			}
		})
	}
}
