package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types/azure"
)

func TestSetMachinePoolDefaults(t *testing.T) {
	cases := []struct {
		name     string
		pool     *azure.MachinePool
		expected []v1beta1.DataDisk
	}{
		{
			name: "both fields omitted",
			pool: &azure.MachinePool{
				DataDisks: []v1beta1.DataDisk{{
					NameSuffix: "etcd",
					DiskSizeGB: 256,
				}},
			},
			expected: []v1beta1.DataDisk{{
				NameSuffix:  "etcd",
				DiskSizeGB:  256,
				CachingType: CachingTypeNone,
			}},
		},
		{
			name: "cachingType omitted, managedDisk specified",
			pool: &azure.MachinePool{
				DataDisks: []v1beta1.DataDisk{{
					NameSuffix:  "etcd",
					DiskSizeGB:  256,
					ManagedDisk: &v1beta1.ManagedDiskParameters{StorageAccountType: "Premium_LRS"},
				}},
			},
			expected: []v1beta1.DataDisk{{
				NameSuffix:  "etcd",
				DiskSizeGB:  256,
				ManagedDisk: &v1beta1.ManagedDiskParameters{StorageAccountType: "Premium_LRS"},
				CachingType: CachingTypeNone,
			}},
		},
		{
			name: "cachingType specified, managedDisk omitted",
			pool: &azure.MachinePool{
				DataDisks: []v1beta1.DataDisk{{
					NameSuffix:  "etcd",
					DiskSizeGB:  256,
					CachingType: "ReadWrite",
				}},
			},
			expected: []v1beta1.DataDisk{{
				NameSuffix:  "etcd",
				DiskSizeGB:  256,
				CachingType: "ReadWrite",
			}},
		},
		{
			name: "both specified",
			pool: &azure.MachinePool{
				DataDisks: []v1beta1.DataDisk{{
					NameSuffix:  "etcd",
					DiskSizeGB:  256,
					ManagedDisk: &v1beta1.ManagedDiskParameters{StorageAccountType: "Premium_LRS"},
					CachingType: "ReadOnly",
				}},
			},
			expected: []v1beta1.DataDisk{{
				NameSuffix:  "etcd",
				DiskSizeGB:  256,
				ManagedDisk: &v1beta1.ManagedDiskParameters{StorageAccountType: "Premium_LRS"},
				CachingType: "ReadOnly",
			}},
		},
		{
			name:     "no data disks",
			pool:     &azure.MachinePool{},
			expected: nil,
		},
		{
			// None is the only caching mode Azure allows on an UltraSSD_LRS
			// disk, so defaulting to it is compatible with the UltraSSD
			// validation rules.
			name: "ultra SSD disk with cachingType omitted",
			pool: &azure.MachinePool{
				DataDisks: []v1beta1.DataDisk{{
					NameSuffix:  "etcd",
					DiskSizeGB:  256,
					ManagedDisk: &v1beta1.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
				}},
			},
			expected: []v1beta1.DataDisk{{
				NameSuffix:  "etcd",
				DiskSizeGB:  256,
				ManagedDisk: &v1beta1.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
				CachingType: CachingTypeNone,
			}},
		},
		{
			name: "mixed explicit and unset caching types",
			pool: &azure.MachinePool{
				DataDisks: []v1beta1.DataDisk{
					{
						NameSuffix: "etcd",
						DiskSizeGB: 256,
					},
					{
						NameSuffix:  "containers",
						DiskSizeGB:  100,
						CachingType: "ReadWrite",
					},
				},
			},
			expected: []v1beta1.DataDisk{
				{
					NameSuffix:  "etcd",
					DiskSizeGB:  256,
					CachingType: CachingTypeNone,
				},
				{
					NameSuffix:  "containers",
					DiskSizeGB:  100,
					CachingType: "ReadWrite",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Defaulting runs on manifests that may already have been
			// defaulted once (e.g. an install-config carried through
			// target transitions), so it must be safe to apply twice.
			SetMachinePoolDefaults(tc.pool)
			SetMachinePoolDefaults(tc.pool)
			assert.Equal(t, tc.expected, tc.pool.DataDisks, "unexpected dataDisks")
		})
	}
}

func TestSetMachinePoolDefaultsDoesNotMutateInheritedDisks(t *testing.T) {
	defaultMachinePlatform := &azure.MachinePool{
		DataDisks: []v1beta1.DataDisk{{
			NameSuffix: "etcd",
			DiskSizeGB: 256,
		}},
	}

	// A machine pool with no platform configuration inherits the default
	// machine platform's DataDisks wholesale, sharing the backing array.
	pool := &azure.MachinePool{}
	Apply(defaultMachinePlatform, pool)
	SetMachinePoolDefaults(pool)

	assert.Equal(t, CachingTypeNone, pool.DataDisks[0].CachingType, "inherited data disk should have been defaulted")
	// The pool's disks must be a private copy: mutating the shared backing
	// array would rewrite the default machine platform and with it every
	// other pool that inherits from it.
	assert.NotEqual(t, &defaultMachinePlatform.DataDisks[0], &pool.DataDisks[0],
		"defaulting must not write through to the default machine platform's slice")
	assert.Empty(t, defaultMachinePlatform.DataDisks[0].CachingType,
		"defaulting must not change the default machine platform's data disks")

	// A pool that overrides the inherited disks keeps its explicit values.
	override := &azure.MachinePool{
		DataDisks: []v1beta1.DataDisk{{
			NameSuffix:  "etcd",
			DiskSizeGB:  256,
			CachingType: "ReadWrite",
		}},
	}
	Apply(defaultMachinePlatform, override)
	SetMachinePoolDefaults(override)
	assert.Equal(t, "ReadWrite", override.DataDisks[0].CachingType, "explicit pool value should have precedence")
}
