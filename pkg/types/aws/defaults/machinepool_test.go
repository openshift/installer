package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

func TestSetMachinePoolDefaults(t *testing.T) {
	cases := []struct {
		name     string
		pool     *aws.MachinePool
		role     string
		expected *aws.MachinePool
	}{
		{
			name:     "nil pool",
			pool:     nil,
			role:     types.MachinePoolComputeRoleName,
			expected: nil,
		},
		{
			name: "empty pool - worker pool sets gp3",
			pool: &aws.MachinePool{},
			role: types.MachinePoolComputeRoleName,
			expected: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: aws.VolumeTypeGp3,
					Size: defaultRootVolumeSize,
				},
			},
		},
		{
			name: "empty pool - control plane pool sets gp3",
			pool: &aws.MachinePool{},
			role: types.MachinePoolControlPlaneRoleName,
			expected: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: aws.VolumeTypeGp3,
					Size: defaultRootVolumeSize,
				},
			},
		},
		{
			name: "empty pool - arbiter pool sets gp3",
			pool: &aws.MachinePool{},
			role: types.MachinePoolArbiterRoleName,
			expected: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: aws.VolumeTypeGp3,
					Size: defaultRootVolumeSize,
				},
			},
		},
		{
			name: "empty pool - edge pool sets gp2",
			pool: &aws.MachinePool{},
			role: types.MachinePoolEdgeRoleName,
			expected: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: aws.VolumeTypeGp2,
					Size: defaultRootVolumeSize,
				},
			},
		},
		{
			name: "pool with existing volume type - should not override",
			pool: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: aws.VolumeTypeGp2,
					Size: defaultRootVolumeSize,
				},
			},
			role: types.MachinePoolComputeRoleName,
			expected: &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					Type: aws.VolumeTypeGp2,
					Size: defaultRootVolumeSize,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetMachinePoolDefaults(tc.pool, tc.role)
			assert.Equal(t, tc.expected, tc.pool, "unexpected machine pool defaults")
		})
	}
}
