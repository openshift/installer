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
				return fmt.Errorf("expected %d errors, got %d", want, have)
			}
			return nil
		}
	}
	noError := exactlyNErrors(0)

	for _, tc := range [...]struct {
		name        string
		machinePool *openstack.MachinePool
		checks      []checkFunc
	}{
		{
			"empty",
			testMachinePool(),
			check(noError),
		},
		{
			"with valid server group policy",
			testMachinePool(withServerGroupPolicy("anti-affinity")),
			check(noError),
		},
		{
			"with invalid server group policy",
			testMachinePool(withServerGroupPolicy("anti-gravity")),
			check(
				someErrorType(field.ErrorTypeNotSupported),
				exactlyNErrors(1),
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateMachinePool(nil, tc.machinePool, "", nil)
			for _, check := range tc.checks {
				if e := check(errs); e != nil {
					t.Error(e)
				}
			}
		})
	}
}
