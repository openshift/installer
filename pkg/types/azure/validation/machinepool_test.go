package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types"
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
		{
			name: "valid disk",
			pool: &azure.MachinePool{
				OSDisk: azure.OSDisk{
					DiskType: "Premium_LRS",
				},
			},
		},
		{
			name: "invalid disk",
			pool: &azure.MachinePool{
				OSDisk: azure.OSDisk{
					DiskType: "LRS",
				},
			},
			expected: `^test-path\.diskType: Unsupported value: "LRS": supported values: "Premium_LRS", "StandardSSD_LRS", "Standard_LRS"$`,
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

func TestValidateMasterDiskType(t *testing.T) {
	cases := []struct {
		name     string
		pool     *types.MachinePool
		expected string
	}{
		{
			name: "empty",
			pool: &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				},
			},
		},
		{
			name: "unsupported disk master",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskType: "Standard_LRS",
						},
					},
				},
			},
			expected: `^test-path\.diskType: Invalid value: "Standard_LRS": Standard_LRS not compatible with control planes.$`,
		},
		{
			name: "supported disk worker",
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskType: "Standard_LRS",
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMasterDiskType(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
