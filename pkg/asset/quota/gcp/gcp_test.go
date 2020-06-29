package gcp

import (
	"fmt"
	"testing"

	"github.com/openshift/installer/pkg/quota"
	"github.com/stretchr/testify/assert"
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

func Test_machineTypeToQuota(t *testing.T) {
	cases := []struct {
		machineType string
		expected    quota.Constraint
	}{{
		machineType: "e2-standard-2",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		machineType: "e2-standard-8-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 8},
	}, {
		machineType: "e2-highmem-2",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		machineType: "e2-highmem-8-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 8},
	}, {
		machineType: "e2-highcpu-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 4},
	}, {
		machineType: "e2-highcpu-16-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 16},
	}, {
		machineType: "n2-standard-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 4},
	}, {
		machineType: "n2-standard-16-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 16},
	}, {
		machineType: "n2-standard-16.5-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 0},
	}, {
		machineType: "n2-highmem-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 4},
	}, {
		machineType: "n2-highmem-32-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 32},
	}, {
		machineType: "n2-highcpu-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 4},
	}, {
		machineType: "n2-highcpu-16",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2_cpus", Count: 16},
	}, {
		machineType: "n2d-standard-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2d_cpus", Count: 4},
	}, {
		machineType: "n2d-standard-16",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2d_cpus", Count: 16},
	}, {
		machineType: "n2d-highmem-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2d_cpus", Count: 4},
	}, {
		machineType: "n2d-highmem-16",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2d_cpus", Count: 16},
	}, {
		machineType: "n2d-highcpu-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2d_cpus", Count: 4},
	}, {
		machineType: "n2d-highcpu-16",
		expected:    quota.Constraint{Name: "compute.googleapis.com/n2d_cpus", Count: 16},
	}, {
		machineType: "n1-standard-2",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		machineType: "n1-standard-8-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 8},
	}, {
		machineType: "c2-standard-4",
		expected:    quota.Constraint{Name: "compute.googleapis.com/c2_cpus", Count: 4},
	}, {
		machineType: "c2-standard-16-4096",
		expected:    quota.Constraint{Name: "compute.googleapis.com/c2_cpus", Count: 16},
	}, {
		machineType: "e2-micro",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		machineType: "e2-medium",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 2},
	}, {
		machineType: "f1-micro",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 1},
	}, {
		machineType: "g1-small",
		expected:    quota.Constraint{Name: "compute.googleapis.com/cpus", Count: 1},
	}}
	for idx, test := range cases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			got := machineTypeToQuota(test.machineType)
			assert.EqualValues(t, test.expected, got)
		})
	}
}
