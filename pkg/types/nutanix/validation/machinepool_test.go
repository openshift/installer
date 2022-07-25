package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/nutanix"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name           string
		pool           *nutanix.MachinePool
		expectedErrMsg string
	}{
		{
			name:           "empty",
			pool:           &nutanix.MachinePool{},
			expectedErrMsg: "",
		}, {
			name: "negative disk size",
			pool: &nutanix.MachinePool{
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: -1,
				},
			},
			expectedErrMsg: `^test-path\.diskSizeGiB: Invalid value: -1: storage disk size must be positive$`,
		}, {
			name: "negative CPUs",
			pool: &nutanix.MachinePool{
				NumCPUs: -1,
			},
			expectedErrMsg: `^test-path\.cpus: Invalid value: -1: number of CPUs must be positive$`,
		}, {
			name: "negative cores",
			pool: &nutanix.MachinePool{
				NumCoresPerSocket: -1,
			},
			expectedErrMsg: `^test-path\.coresPerSocket: Invalid value: -1: cores per socket must be positive$`,
		}, {
			name: "negative memory",
			pool: &nutanix.MachinePool{
				MemoryMiB: -1,
			},
			expectedErrMsg: `^test-path\.memoryMiB: Invalid value: -1: memory size must be positive$`,
		}, {
			name: "less CPUs than cores per socket",
			pool: &nutanix.MachinePool{
				NumCPUs:           1,
				NumCoresPerSocket: 8,
			},
			expectedErrMsg: `^test-path\.coresPerSocket: Invalid value: 8: cores per socket must be less than number of CPUs$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expectedErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedErrMsg, err)
			}
		})
	}
}
