package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name     string
		pool     *azure.MachinePool
		expected string
	}{
		{
			name: "empty",
			pool: &azure.MachinePool{},
		},
		{
			name: "valid iops",
			pool: &azure.MachinePool{
				OSDisk: azure.OSDisk{
					DiskSizeGB: 120,
				},
			},
		},
		{
			name: "invalid iops",
			pool: &azure.MachinePool{
				OSDisk: azure.OSDisk{
					DiskSizeGB: -120,
				},
			},
			expected: `^test-path\.diskSizeGB: Invalid value: -120: Storage DiskSizeGB must be positive$`,
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
