package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

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

// azurePoolForDiskSetup returns an Azure machine pool carrying one data disk per
// diskSetup entry, paired by index and name the way Azure expects. The cases
// below exercise the platform-agnostic disk setup rules, but Azure now rejects a
// diskSetup entry that has no data disk to configure, so the disks have to be
// present for a pool to reach the rule under test.
func azurePoolForDiskSetup(diskSetup []types.Disk) *azure.MachinePool {
	p := &azure.MachinePool{}
	for i, ds := range diskSetup {
		// A malformed entry has no disk ID; an empty name suffix still pairs,
		// because the match is only checked for well-formed entries.
		nameSuffix := ""
		switch ds.Type {
		case types.Etcd:
			if ds.Etcd != nil {
				nameSuffix = ds.Etcd.PlatformDiskID
			}
		case types.Swap:
			if ds.Swap != nil {
				nameSuffix = ds.Swap.PlatformDiskID
			}
		case types.UserDefined:
			if ds.UserDefined != nil {
				nameSuffix = ds.UserDefined.PlatformDiskID
			}
		}
		p.DataDisks = append(p.DataDisks, capz.DataDisk{
			NameSuffix: nameSuffix,
			DiskSizeGB: 32,
			Lun:        ptr.To(int32(i)),
		})
	}
	return p
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
			name:     "valid CAPI management on AWS worker",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool(types.MachinePoolComputeRoleName)
				p.Management = types.ClusterAPI
				return p
			}(),
			valid: true,
		},
		{
			name:     "valid CAPI management on AWS edge",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool(types.MachinePoolEdgeRoleName)
				p.Management = types.ClusterAPI
				return p
			}(),
			valid: true,
		},
		{
			name:     "invalid CAPI management on AWS master",
			platform: &types.Platform{AWS: &aws.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool(types.MachinePoolControlPlaneRoleName)
				p.Management = types.ClusterAPI
				return p
			}(),
			expectedError: `master machines cannot be managed by Cluster API`,
		},
		{
			name:     "invalid CAPI management on GCP",
			platform: &types.Platform{GCP: &gcp.Platform{Region: "us-east-1"}},
			pool: func() *types.MachinePool {
				p := validMachinePool(types.MachinePoolComputeRoleName)
				p.Management = types.ClusterAPI
				return p
			}(),
			expectedError: `machines cannot be managed by Cluster API for platform gcp`,
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.etcd: Invalid value: "type: etcd\\n": etcd configuration must be created$`,
		},
		{
			name:     "invalid swap disk on master machine pool",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{
					Type: "swap",
					Swap: &types.DiskSwap{PlatformDiskID: "swap"},
				})
				p.Platform = types.MachinePoolPlatform{
					Azure: azurePoolForDiskSetup(p.DiskSetup),
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.swap: Invalid value: "swap:\\n  platformDiskID: swap\\ntype: swap\\n": swap is unsupported on control plane nodes$`,
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.userDefined: Invalid value: "type: user-defined\\n": userDefined configuration must be created$`,
		},
		{
			// The cases from here to "valid distinct user-defined disks" cover rules
			// that hold for every platform, so they leave the platform machine pool
			// unset and reach validateDiskSetup without any platform validation
			// running on top.
			name:     "invalid duplicate platformDiskID across disk setup types",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				// Both entries claim the same data disk, but resolve to different
				// labels ("etcd" and "shared"), so only the ID check should fire.
				p.DiskSetup = append(p.DiskSetup,
					types.Disk{Type: types.Etcd, Etcd: &types.DiskEtcd{PlatformDiskID: "shared"}},
					types.Disk{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{PlatformDiskID: "shared", MountPath: "/mnt/data"}},
				)
				return p
			}(),
			expectedError: `^test-path\.diskSetup\[1\]\.userDefined\.platformDiskID: Duplicate value: "shared"$`,
		},
		{
			name:     "invalid user-defined label normalizing onto the etcd label",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				// "et-cd" is stripped of non-alphanumerics to "etcd", colliding with
				// the etcd MachineConfig. The platform disk IDs differ, so only the
				// label check should fire.
				p.DiskSetup = append(p.DiskSetup,
					types.Disk{Type: types.Etcd, Etcd: &types.DiskEtcd{PlatformDiskID: "etcddisk"}},
					types.Disk{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{PlatformDiskID: "et-cd", MountPath: "/mnt/data"}},
				)
				return p
			}(),
			expectedError: `^test-path\.diskSetup\[1\]: Invalid value: .*: resolves to disk label "etcd", which is already used by an earlier diskSetup entry in this pool`,
		},
		{
			name:     "invalid duplicate user-defined labels differing only in case",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				// The MachineConfig name is lower-cased, so these name one object.
				p.DiskSetup = append(p.DiskSetup,
					types.Disk{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{PlatformDiskID: "Data", MountPath: "/mnt/a"}},
					types.Disk{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{PlatformDiskID: "data", MountPath: "/mnt/b"}},
				)
				return p
			}(),
			expectedError: `^test-path\.diskSetup\[1\]: Invalid value: .*: resolves to disk label "data", which is already used by an earlier diskSetup entry in this pool`,
		},
		{
			name:     "invalid unrecognized disk setup type",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{Type: "bogus"})
				return p
			}(),
			expectedError: `^test-path\.diskSetup\[0\]\.type: Unsupported value: "bogus": supported values: "etcd", "swap", "user-defined"$`,
		},
		{
			name:     "invalid empty disk setup type",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("master")
				p.DiskSetup = append(p.DiskSetup, types.Disk{})
				return p
			}(),
			expectedError: `^test-path\.diskSetup\[0\]\.type: Unsupported value: "": supported values: "etcd", "swap", "user-defined"$`,
		},
		{
			name:     "valid distinct user-defined disks",
			platform: &types.Platform{Azure: &azure.Platform{Region: "eastus"}},
			pool: func() *types.MachinePool {
				p := validMachinePool("worker")
				p.DiskSetup = append(p.DiskSetup,
					types.Disk{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{PlatformDiskID: "containers", MountPath: "/var/lib/containers"}},
					types.Disk{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{PlatformDiskID: "kubelet", MountPath: "/var/lib/kubelet"}},
				)
				return p
			}(),
			valid: true,
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.etcd: Too many: 2: must have at most 1 item$`,
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
				}
				return p
			}(),
			expectedError: `^test-path\.diskSetup\.swap: Too many: 2: must have at most 1 item$`,
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
					Azure: azurePoolForDiskSetup(p.DiskSetup),
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
