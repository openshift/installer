package gcp

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	computev1 "google.golang.org/api/compute/v1"

	"github.com/openshift/installer/pkg/quota"
)

func Test_aggregate(t *testing.T) {
	cases := []struct {
		input []quota.Constraint

		exp []quota.Constraint
	}{{
		input: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q1", Region: "g", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 2},
		},
	}, {
		input: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
			{Name: "q3", Region: "r1", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 2},
			{Name: "q3", Region: "r1", Count: 1},
		},
	}, {
		input: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
			{Name: "q3", Region: "r1", Count: 1},
			{Name: "q4", Region: "r1", Count: 1},
		},
		exp: []quota.Constraint{
			{Name: "q1", Region: "g", Count: 1},
			{Name: "q2", Region: "r1", Count: 1},
			{Name: "q3", Region: "r1", Count: 1},
			{Name: "q4", Region: "r1", Count: 1},
		},
	}}

	for idx, test := range cases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			got := aggregate(test.input)
			assert.EqualValues(t, test.exp, got)
		})
	}
}

func Test_guessMachineCPUCount(t *testing.T) {
	cases := []struct {
		machineType string
		expected    int64
	}{{
		machineType: "e2-standard-2",
		expected:    2,
	}, {
		machineType: "e2-standard-8-4096",
		expected:    8,
	}, {
		machineType: "e2-highmem-2",
		expected:    2,
	}, {
		machineType: "e2-highmem-8-4096",
		expected:    8,
	}, {
		machineType: "e2-highcpu-4",
		expected:    4,
	}, {
		machineType: "e2-highcpu-16-4096",
		expected:    16,
	}, {
		machineType: "n2-standard-4",
		expected:    4,
	}, {
		machineType: "n2-standard-16-4096",
		expected:    16,
	}, {
		machineType: "n2-standard-16.5-4096",
		expected:    0,
	}, {
		machineType: "n2-highmem-4",
		expected:    4,
	}, {
		machineType: "n2-highmem-32-4096",
		expected:    32,
	}, {
		machineType: "n2-highcpu-4",
		expected:    4,
	}, {
		machineType: "n2-highcpu-16",
		expected:    16,
	}, {
		machineType: "n2d-standard-4",
		expected:    4,
	}, {
		machineType: "n2d-standard-16",
		expected:    16,
	}, {
		machineType: "n2d-highmem-4",
		expected:    4,
	}, {
		machineType: "n2d-highmem-16",
		expected:    16,
	}, {
		machineType: "n2d-highcpu-4",
		expected:    4,
	}, {
		machineType: "n2d-highcpu-16",
		expected:    16,
	}, {
		machineType: "n1-standard-2",
		expected:    2,
	}, {
		machineType: "n1-standard-8-4096",
		expected:    8,
	}, {
		machineType: "c2-standard-4",
		expected:    4,
	}, {
		machineType: "c2-standard-16-4096",
		expected:    16,
	}, {
		machineType: "e2-micro",
		expected:    2,
	}, {
		machineType: "e2-medium",
		expected:    2,
	}, {
		machineType: "f1-micro",
		expected:    1,
	}, {
		machineType: "g1-small",
		expected:    1,
	}}
	for idx, test := range cases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			got := guessMachineCPUCount(test.machineType)
			assert.EqualValues(t, test.expected, got)
		})
	}
}

func Test_machineTypeToQuota(t *testing.T) {
	fake := newFakeMachineTypeGetter([]*computev1.MachineType{{
		Zone:      "a",
		Name:      "n1-standard-2",
		GuestCpus: 2,
	}, {
		Zone:      "a",
		Name:      "n1-custom-2-1024",
		GuestCpus: 2,
	}, {
		Zone:      "a",
		Name:      "custom-2-1024",
		GuestCpus: 2,
	}, {
		Zone:      "a",
		Name:      "n2-standard-4",
		GuestCpus: 4,
	}, {
		Zone:      "a",
		Name:      "n2-custom-4-1024",
		GuestCpus: 4,
	}})

	tests := []struct {
		zone        string
		machineType string
		expected    quota.Constraint
	}{{
		zone:        "a",
		machineType: "n1-standard-2",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		zone:        "a",
		machineType: "n1-custom-2-1024",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		zone:        "a",
		machineType: "n1-custom-2-2048",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		zone:        "b",
		machineType: "n1-standard-2",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		zone:        "a",
		machineType: "custom-2-1024",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		zone:        "a",
		machineType: "custom-2-2048",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 0},
	}, {
		zone:        "b",
		machineType: "custom-2-1024",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 0},
	}, {
		zone:        "a",
		machineType: "n2-standard-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 4},
	}, {
		zone:        "a",
		machineType: "n2-custom-4-1024",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 4},
	}, {
		zone:        "b",
		machineType: "n2-standard-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 4},
	}}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			got := machineTypeToQuota(fake, test.zone, test.machineType)
			assert.EqualValues(t, test.expected, got)
		})
	}
}

type fakeMachineTypeGetter struct {
	knownTypes map[string]*computev1.MachineType
}

func newFakeMachineTypeGetter(mtypes []*computev1.MachineType) *fakeMachineTypeGetter {
	fake := &fakeMachineTypeGetter{
		knownTypes: map[string]*computev1.MachineType{},
	}

	for idx, mtype := range mtypes {
		fake.knownTypes[fmt.Sprintf("%s__%s", mtype.Zone, mtype.Name)] = mtypes[idx]
	}

	return fake
}

func (fake *fakeMachineTypeGetter) GetMachineType(zone string, machineType string) (*computev1.MachineType, error) {
	mtype, ok := fake.knownTypes[fmt.Sprintf("%s__%s", zone, machineType)]
	if !ok {
		return nil, errors.New("unknwown")
	}
	return mtype, nil
}
