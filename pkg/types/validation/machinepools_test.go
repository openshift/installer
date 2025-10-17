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
	"github.com/openshift/installer/pkg/types/ibmcloud"
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

// Cursor generated disk Setup tests

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name          string
		platform      *types.Platform
		pool          *types.MachinePool
		valid         bool
		expectedError string
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
							Type: "io1",
							Size: 128,
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
			platform: &types.Platform{IBMCloud: &ibmcloud.Platform{}},
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
					AWS:      &aws.MachinePool{},
					IBMCloud: &ibmcloud.MachinePool{},
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
		{
			name:     "valid GCP disk type master",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{},
				}
				p.Platform.GCP.OSDisk.DiskSizeGB = 100
				p.Platform.GCP.OSDisk.DiskType = "hyperdisk-balanced"
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid GCP disk type master",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{},
				}
				p.Platform.GCP.OSDisk.DiskSizeGB = 100
				p.Platform.GCP.OSDisk.DiskType = "pd-standard"
				return p
			}(),
			valid: false,
		},
		{
			name:     "valid GCP service account use",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1", NetworkProjectID: "ExampleNetworkProject"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{
						ServiceAccount: "ExampleServiceAccount@ExampleServiceAccount.com",
					},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid GCP service account on machine pool type",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{
						ServiceAccount: "ExampleServiceAccount@ExampleServiceAccount.com",
					},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid GCP service account non xpn install",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.Platform = types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{
						ServiceAccount: "ExampleServiceAccount@ExampleServiceAccount.com",
					},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "valid multiple disks",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "etcd",
					UserDefined: nil,
					Etcd:        &types.DiskEtcd{PlatformDiskID: "etcd"},
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid etcd disk type on worker machine pool",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "etcd",
					UserDefined: nil,
					Etcd:        &types.DiskEtcd{PlatformDiskID: "etcd"},
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.etcd: Invalid value: "etcd:\\n  platformDiskID: etcd\\ntype: etcd\\n": cannot specify etcd on worker machine pools$`,
		},
		{
			name:     "valid etcd disk on master machine pool",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "etcd",
					UserDefined: nil,
					Etcd:        &types.DiskEtcd{PlatformDiskID: "etcd"},
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid etcd disk with nil Etcd field",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "etcd",
					UserDefined: nil,
					Etcd:        nil,
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.etcd: Invalid value: "type: etcd\\n": etcd configuration must be created$`,
		},
		{
			name:     "valid swap disk",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "swap",
					UserDefined: nil,
					Etcd:        nil,
					Swap:        &types.DiskSwap{PlatformDiskID: "swap"},
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid swap disk with nil Swap field",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "swap",
					UserDefined: nil,
					Etcd:        nil,
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.swap: Invalid value: "type: swap\\n": swap configuration must be created$`,
		},
		{
			name:     "valid user-defined disk",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "user-defined",
					UserDefined: &types.DiskUserDefined{PlatformDiskID: "userdisk", MountPath: "/mnt/data"},
					Etcd:        nil,
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid user-defined disk platformDiskId too long",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "user-defined",
					UserDefined: &types.DiskUserDefined{PlatformDiskID: "userdiskuserdisk", MountPath: "/mnt/data"},
					Etcd:        nil,
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.userDefined\.platformDiskId: Invalid value: \"type: user-defined\\nuserDefined:\\n  mountPath: /mnt/data\\n  platformDiskID: userdiskuserdisk\\n": cannot be longer than 12 characters$`,
		},
		{
			name:     "invalid user-defined disk with nil UserDefined field",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "user-defined",
					UserDefined: nil,
					Etcd:        nil,
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.userDefined: Invalid value: "type: user-defined\\n": userDefined configuration must be created$`,
		},
		{
			name:     "invalid multiple etcd disks",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "etcd",
					UserDefined: nil,
					Etcd:        &types.DiskEtcd{PlatformDiskID: "etcd1"},
					Swap:        nil,
				})
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "etcd",
					UserDefined: nil,
					Etcd:        &types.DiskEtcd{PlatformDiskID: "etcd2"},
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.etcd: Too many: 2: must have at most 1 items$`,
		},
		{
			name:     "invalid multiple swap disks",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "swap",
					UserDefined: nil,
					Etcd:        nil,
					Swap:        &types.DiskSwap{PlatformDiskID: "swap1"},
				})
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "swap",
					UserDefined: nil,
					Etcd:        nil,
					Swap:        &types.DiskSwap{PlatformDiskID: "swap2"},
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.swap: Too many: 2: must have at most 1 items$`,
		},
		{
			name:     "valid mixed disk types",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "etcd",
					UserDefined: nil,
					Etcd:        &types.DiskEtcd{PlatformDiskID: "etcd"},
					Swap:        nil,
				})
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type:        "user-defined",
					UserDefined: &types.DiskUserDefined{PlatformDiskID: "userdisk", MountPath: "/mnt/data"},
					Etcd:        nil,
					Swap:        nil,
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				}
				return p
			}(),
			valid: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.platform, tc.pool, field.NewPath("test-path")).ToAggregate()

			switch {
			case tc.expectedError != "":
				assert.Regexp(t, tc.expectedError, err)
			case tc.valid:
				assert.NoError(t, err)
			case !tc.valid:
				assert.Error(t, err)
			}
		})
	}
}
