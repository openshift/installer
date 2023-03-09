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

func withZones() func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) { mp.Zones = []string{"one", "two", "three"} }
}

func withVolumeZones() func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) {
		mp.RootVolume = &openstack.RootVolume{Zones: []string{"one", "two", "three"}}
	}
}

func withFailureDomain(fd openstack.FailureDomain) func(*openstack.MachinePool) {
	return func(mp *openstack.MachinePool) {
		mp.FailureDomains = append(mp.FailureDomains, fd)
	}
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
			"failure domains",
			testMachinePool(withFailureDomain(openstack.FailureDomain{
				PortTargets: []openstack.NamedPortTarget{{ID: "one"}, {ID: "two"}},
			})),
			"master",
			check(noError),
		},
		{
			"failure domains on worker",
			testMachinePool(withFailureDomain(openstack.FailureDomain{
				PortTargets: []openstack.NamedPortTarget{{ID: "one"}, {ID: "two"}},
			})),
			"worker",
			check(
				someErrorType(field.ErrorTypeForbidden),
				exactlyNErrors(1),
			),
		},
		{
			"failure domains, duplicate portTarget ID",
			testMachinePool(withFailureDomain(openstack.FailureDomain{
				PortTargets: []openstack.NamedPortTarget{{ID: "one"}, {ID: "one"}},
			})),
			"master",
			check(
				someErrorType(field.ErrorTypeDuplicate),
				exactlyNErrors(1),
			),
		},
		{
			"failure domains together with zones",
			testMachinePool(
				withZones(),
				withFailureDomain(openstack.FailureDomain{
					PortTargets: []openstack.NamedPortTarget{{ID: "one"}, {ID: "two"}},
				}),
			),
			"master",
			check(
				someErrorType(field.ErrorTypeForbidden),
				exactlyNErrors(1),
			),
		},
		{
			"failure domains together with root volume zones",
			testMachinePool(
				withVolumeZones(),
				withFailureDomain(openstack.FailureDomain{
					PortTargets: []openstack.NamedPortTarget{{ID: "one"}, {ID: "two"}},
				}),
			),
			"master",
			check(
				someErrorType(field.ErrorTypeForbidden),
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
