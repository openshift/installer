package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name     string
		pool     *types.MachinePool
		platform string
		valid    bool
	}{
		{
			name: "minimal",
			pool: &types.MachinePool{
				Name: "master",
			},
			platform: "aws",
			valid:    true,
		},
		{
			name: "invalid name",
			pool: &types.MachinePool{
				Name: "bad-name",
			},
			platform: "aws",
			valid:    false,
		},
		{
			name: "invalid replicas",
			pool: &types.MachinePool{
				Name:     "master",
				Replicas: func(x int64) *int64 { return &x }(-1),
			},
			platform: "aws",
			valid:    false,
		},
		{
			name: "valid aws",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				},
			},
			platform: "aws",
			valid:    true,
		},
		{
			name: "invalid aws",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						EC2RootVolume: aws.EC2RootVolume{
							IOPS: -10,
						},
					},
				},
			},
			platform: "aws",
			valid:    false,
		},
		{
			name: "valid libvirt",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Libvirt: &libvirt.MachinePool{},
				},
			},
			platform: "libvirt",
			valid:    true,
		},
		{
			name: "valid openstack",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					OpenStack: &openstack.MachinePool{},
				},
			},
			platform: "openstack",
			valid:    true,
		},
		{
			name: "mis-matched platform",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				},
			},
			platform: "libvirt",
			valid:    false,
		},
		{
			name: "multiple platforms",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					AWS:     &aws.MachinePool{},
					Libvirt: &libvirt.MachinePool{},
				},
			},
			platform: "aws",
			valid:    false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path"), tc.platform).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
