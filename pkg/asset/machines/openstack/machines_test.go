package openstack

import (
	"fmt"
	"strings"
	"testing"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types/openstack"
)

func mpWithZones(zones ...string) func(*openstack.MachinePool) {
	return func(mpool *openstack.MachinePool) {
		mpool.Zones = zones
	}
}

func mpWithRootVolumeZones(zones ...string) func(*openstack.MachinePool) {
	return func(mpool *openstack.MachinePool) {
		if mpool.RootVolume != nil {
			mpool.RootVolume.Zones = zones
		} else {
			mpool.RootVolume = &openstack.RootVolume{Zones: zones}
		}
	}
}

func mpWithRootVolumeTypes(types ...string) func(*openstack.MachinePool) {
	return func(mpool *openstack.MachinePool) {
		if mpool.RootVolume != nil {
			mpool.RootVolume.Types = types
		} else {
			mpool.RootVolume = &openstack.RootVolume{Types: types}
		}
	}
}

func generateMachinePool(options ...func(*openstack.MachinePool)) openstack.MachinePool {
	mpool := openstack.MachinePool{}
	for _, apply := range options {
		apply(&mpool)
	}
	return mpool
}

func TestFailureDomains(t *testing.T) {
	type checkFunc func([]machinev1.OpenStackFailureDomain, error) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasNFailureDomains := func(want int) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			if have := len(fds); want != have {
				return fmt.Errorf("expected %d failure domains, got %d", want, have)
			}
			return nil
		}
	}

	hasComputeZones := func(wantZones ...string) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			haveZones := make([]string, len(fds))
			for i := range fds {
				haveZones[i] = fds[i].AvailabilityZone
			}

			if wantLen, haveLen := len(wantZones), len(haveZones); wantLen != haveLen {
				return fmt.Errorf("expected compute zones %v (len %d), got %v (len %d)", wantZones, wantLen, haveZones, haveLen)
			}

			for i := range fds {
				if want, have := wantZones[i], haveZones[i]; want != have {
					return fmt.Errorf("expected compute zones %v, got %v", wantZones, haveZones)
				}
			}

			return nil
		}
	}

	hasNilRootVolume := func(fds []machinev1.OpenStackFailureDomain, _ error) error {
		for i := range fds {
			if fds[i].RootVolume != nil {
				return fmt.Errorf("failure domain %d has unexpectedly non-nil RootVolume", i)
			}
		}
		return nil
	}

	hasRootVolumeZones := func(wantZones ...string) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			haveZones := make([]string, len(fds))
			for i := range fds {
				if fds[i].RootVolume == nil {
					return fmt.Errorf("failure domain %d has unexpectedly nil RootVolume", i)
				}
				haveZones[i] = fds[i].RootVolume.AvailabilityZone
			}

			if wantLen, haveLen := len(wantZones), len(haveZones); wantLen != haveLen {
				return fmt.Errorf("expected root volume zones %v, got %v", wantZones, haveZones)
			}

			for i := range fds {
				if want, have := wantZones[i], haveZones[i]; want != have {
					return fmt.Errorf("expected root volume zones %v, got %v", wantZones, haveZones)
				}
			}

			return nil
		}
	}

	hasRootVolumeTypes := func(wantTypes ...string) checkFunc {
		return func(fds []machinev1.OpenStackFailureDomain, _ error) error {
			haveTypes := make([]string, len(fds))
			for i := range fds {
				if fds[i].RootVolume == nil {
					return fmt.Errorf("failure domain %d has unexpectedly nil RootVolume", i)
				}
				haveTypes[i] = fds[i].RootVolume.VolumeType
			}

			if wantLen, haveLen := len(wantTypes), len(haveTypes); wantLen != haveLen {
				return fmt.Errorf("expected root volume types %v, got %v", wantTypes, haveTypes)
			}

			for i := range fds {
				if want, have := wantTypes[i], haveTypes[i]; want != have {
					return fmt.Errorf("expected root volume types %v, got %v", wantTypes, haveTypes)
				}
			}

			return nil
		}
	}

	doesNotPanic := func(_ []machinev1.OpenStackFailureDomain, have error) error {
		if have != nil {
			return fmt.Errorf("unexpected panic: %w", have)
		}
		return nil
	}

	panicsWith := func(want string) checkFunc {
		return func(_ []machinev1.OpenStackFailureDomain, have error) error {
			if have == nil {
				return fmt.Errorf("unexpectedly, didn't panic")
			}
			if have := fmt.Sprintf("%v", have); !strings.Contains(have, want) {
				return fmt.Errorf("expected panic with %q, got %q", want, have)
			}
			return nil
		}
	}

	for _, tc := range [...]struct {
		name   string
		mpool  openstack.MachinePool
		checks []checkFunc
	}{
		{
			"no_zones",
			generateMachinePool(),
			check(
				hasNFailureDomains(1),
				hasComputeZones(""),
				hasNilRootVolume,
				doesNotPanic,
			),
		},
		{
			"one_compute_zone",
			generateMachinePool(
				mpWithZones("one"),
			),
			check(
				hasNFailureDomains(1),
				hasComputeZones("one"),
				hasNilRootVolume,
				doesNotPanic,
			),
		},
		{
			"three_compute_zones",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "two", "three"),
				hasNilRootVolume,
				doesNotPanic,
			),
		},
		{
			"three_compute_zones_one_root_volume_zone",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
				mpWithRootVolumeZones("volume_one"),
				mpWithRootVolumeTypes("type-1"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "two", "three"),
				hasRootVolumeZones("volume_one", "volume_one", "volume_one"),
				hasRootVolumeTypes("type-1", "type-1", "type-1"),
				doesNotPanic,
			),
		},
		{
			"one_compute_zone_three_root_volume_zones",
			generateMachinePool(
				mpWithZones("one"),
				mpWithRootVolumeZones("volume_one", "volume_two", "volume_three"),
				mpWithRootVolumeTypes("type-1"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "one", "one"),
				hasRootVolumeZones("volume_one", "volume_two", "volume_three"),
				hasRootVolumeTypes("type-1", "type-1", "type-1"),
				doesNotPanic,
			),
		},
		{
			"three_compute_zone_two_root_volume_zones_panics",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
				mpWithRootVolumeZones("volume_one", "volume_two"),
			),
			check(
				// We have to check for a partial result here, because the mapping
				// of compute zones to root volume zones is handled in a map therefore
				// the order is not deterministic.
				panicsWith("availability zones should have equal length"),
			),
		},
		{
			"three_compute_zones_three_root_volume_types",
			generateMachinePool(
				mpWithZones("one", "two", "three"),
				mpWithRootVolumeZones("volume_one", "volume_two", "volume_three"),
				mpWithRootVolumeTypes("type-1", "type-2", "type-3"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("one", "two", "three"),
				hasRootVolumeZones("volume_one", "volume_two", "volume_three"),
				hasRootVolumeTypes("type-1", "type-2", "type-3"),
				doesNotPanic,
			),
		},
		{
			"three_root_volume_types",
			generateMachinePool(
				mpWithRootVolumeTypes("type-1", "type-2", "type-3"),
			),
			check(
				hasNFailureDomains(3),
				hasComputeZones("", "", ""),
				hasRootVolumeZones("", "", ""),
				hasRootVolumeTypes("type-1", "type-2", "type-3"),
				doesNotPanic,
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			failureDomains, recoveredPanic := func() (fds []machinev1.OpenStackFailureDomain, recoveredPanic error) {
				defer func() {
					if r := recover(); r != nil {
						recoveredPanic = fmt.Errorf("%v", r)
					}
				}()

				fds = failureDomainsFromSpec(tc.mpool)
				return
			}()

			for _, check := range tc.checks {
				if err := check(failureDomains, recoveredPanic); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestPruneFailureDomains(t *testing.T) {
}
