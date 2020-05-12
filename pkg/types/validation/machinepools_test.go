package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
)

func validMachinePool(name string) *types.MachinePool {
	return &types.MachinePool{
		Name:           name,
		Replicas:       pointer.Int64Ptr(1),
		Hyperthreading: types.HyperthreadingDisabled,
		Architecture:   types.ArchitectureAMD64,
	}
}

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name     string
		platform *types.Platform
		pool     *types.MachinePool
		valid    bool
	}{
		{
			name:     "minimal",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool:     validMachinePool("test-name"),
			valid:    true,
		},
		{
			name:     "missing replicas",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Replicas = nil
				return p
			}(),
			valid: false,
		},
		{
			name:     "invalid replicas",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Replicas = pointer.Int64Ptr(-1)
				return p
			}(),
			valid: false,
		},
		{
			name:     "valid aws",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid aws",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						EC2RootVolume: aws.EC2RootVolume{
							IOPS: -10,
						},
					},
				}
				return p
			}(),
			valid: false,
		},
		{
			name:     "valid azure",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "valid libvirt",
			platform: &types.Platform{Libvirt: &libvirt.Platform{}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					Libvirt: &libvirt.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "valid openstack",
			platform: &types.Platform{OpenStack: &openstack.Platform{}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					OpenStack: &openstack.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "mis-matched platform",
			platform: &types.Platform{Libvirt: &libvirt.Platform{}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				}
				return p
			}(),
			valid: false,
		},
		{
			name:     "multiple platforms",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					AWS:     &aws.MachinePool{},
					Libvirt: &libvirt.MachinePool{},
				}
				return p
			}(),
			valid: false,
		},
		{
			name:     "valid GCP",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{},
				}
				p.Platform.GCP.OSDisk.DiskSizeGB = 100
				p.Platform.GCP.OSDisk.DiskType = "pd-standard"
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid GCP disk size",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{},
				}
				p.Platform.GCP.OSDisk.DiskSizeGB = -100
				p.Platform.GCP.OSDisk.DiskType = "pd-standard"
				return p
			}(),
			valid: false,
		},
		{
			name:     "invalid GCP disk type",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("test-name")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{},
				}
				p.Platform.GCP.OSDisk.DiskSizeGB = 100
				p.Platform.GCP.OSDisk.DiskType = "pd-"
				return p
			}(),
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.platform, tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
