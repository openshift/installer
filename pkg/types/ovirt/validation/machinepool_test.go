package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/ovirt"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name  string
		pool  *ovirt.MachinePool
		valid bool
	}{
		{
			name:  "empty",
			pool:  &ovirt.MachinePool{},
			valid: true,
		},
		{
			name: "invalid CPU cores",
			pool: &ovirt.MachinePool{
				CPU: &ovirt.CPU{
					Cores:   0,
					Sockets: 1,
				},
			},
			valid: false,
		},
		{
			name: "invalid CPU sockets",
			pool: &ovirt.MachinePool{
				CPU: &ovirt.CPU{
					Cores:   1,
					Sockets: 0,
				},
			},
			valid: false,
		},
		{
			name: "invalid OSDisk",
			pool: &ovirt.MachinePool{
				OSDisk: &ovirt.Disk{
					SizeGB: 0,
				},
			},
			valid: false,
		},
		{
			name: "invalid vmType",
			pool: &ovirt.MachinePool{
				VMType: "strong",
			},
			valid: false,
		},
		{
			name: "invalid instance type id",
			pool: &ovirt.MachinePool{
				InstanceTypeID: "aaaaa-123",
			},
			valid: false,
		},
		{
			name: "invalid mix of cpu and instance type id",
			pool: &ovirt.MachinePool{
				CPU: &ovirt.CPU{
					Sockets: 1,
					Cores:   4,
				},
				InstanceTypeID: "85c65199-2df1-43bf-94f6-7e1567e6b238",
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
