package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name     string
		pool     *powervs.MachinePool
		expected string
	}{
		{
			name: "empty",
			pool: &powervs.MachinePool{},
		},
		{
			name: "valid volumeIDs",
			pool: &powervs.MachinePool{
				VolumeIDs: []string{"c8b709c4-93f1-499e-915e-0820bcc51406", "587c5788-107f-4351-aabc-1652c54c4491"},
			},
		},
		{
			name: "invalid volumeIDs",
			pool: &powervs.MachinePool{
				VolumeIDs: []string{"c8b709c4-93f1-499e-915e-0820bcc51406", "abc123"},
			},
			expected: `^test-path\.volumeIDs\[1]: Invalid value: "abc123": volume ID must be a valid UUID$`,
		},
		{
			name: "valid memory",
			pool: &powervs.MachinePool{
				Memory: "5",
			},
		},
		{
			name: "invalid memory under",
			pool: &powervs.MachinePool{
				Memory: "1",
			},
			expected: `^test-path\.memory: Invalid value: "1": memory must be an integer number of GB that is at least 2 and no more than 64$`,
		},
		{
			name: "invalid memory over",
			pool: &powervs.MachinePool{
				Memory: "65",
			},
			expected: `^test-path\.memory: Invalid value: "65": memory must be an integer number of GB that is at least 2 and no more than 64$`,
		},
		{
			name: "invalid memory string",
			pool: &powervs.MachinePool{
				Memory: "all",
			},
			expected: `^test-path\.memory: Invalid value: "all": memory must be an integer number of GB that is at least 2 and no more than 64$`,
		},
		{
			name: "valid processors",
			pool: &powervs.MachinePool{
				Processors: "1.25",
			},
		},
		{
			name: "invalid processors under",
			pool: &powervs.MachinePool{
				Processors: "0",
			},
			expected: `^test-path\.processors: Invalid value: "0": number of processors must be from \.25 to 32 cores$`,
		},
		{
			name: "invalid processors over",
			pool: &powervs.MachinePool{
				Processors: "33",
			},
			expected: `^test-path\.processors: Invalid value: "33": number of processors must be from \.25 to 32 cores$`,
		},
		{
			name: "invalid processors string",
			pool: &powervs.MachinePool{
				Processors: "all",
			},
			expected: `^test-path\.processors: Invalid value: "all": processors must be a valid floating point number$`,
		},
		{
			name: "invalid processors increment",
			pool: &powervs.MachinePool{
				Processors: "1.33",
			},
			expected: `^test-path\.processors: Invalid value: "1\.33": processors must be in increments of \.25$`,
		},
		{
			name: "valid procType",
			pool: &powervs.MachinePool{
				ProcType: "shared",
			},
		},
		{
			name: "invalid procType",
			pool: &powervs.MachinePool{
				ProcType: "none",
			},
			expected: `^test-path\.procType: Unsupported value: "none": supported values: "capped", "dedicated", "shared"$`,
		},
		{
			name: "valid sysType",
			pool: &powervs.MachinePool{
				SysType: "s922",
			},
		},
		{
			name: "invalid sysType",
			pool: &powervs.MachinePool{
				SysType: "p922",
			},
			expected: `^test-path\.sysType: Invalid value: "p922": system type must be one of {e980,s922}$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
