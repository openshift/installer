package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/types"
)

func defaultMachinePool(name string) *types.MachinePool {
	return &types.MachinePool{
		Name:           name,
		Replicas:       pointer.Int64Ptr(3),
		Hyperthreading: types.HyperthreadingEnabled,
		Architecture:   types.ArchitectureAMD64,
	}
}

func TestSetMahcinePoolDefaults(t *testing.T) {
	cases := []struct {
		name     string
		pool     *types.MachinePool
		platform string
		expected *types.MachinePool
	}{
		{
			name:     "empty",
			pool:     &types.MachinePool{},
			expected: defaultMachinePool(""),
		},
		{
			name:     "default",
			pool:     defaultMachinePool("test-name"),
			expected: defaultMachinePool("test-name"),
		},
		{
			name: "non-default replicas",
			pool: func() *types.MachinePool {
				p := defaultMachinePool("test-name")
				p.Replicas = pointer.Int64Ptr(5)
				return p
			}(),
			expected: func() *types.MachinePool {
				p := defaultMachinePool("test-name")
				p.Replicas = pointer.Int64Ptr(5)
				return p
			}(),
		},
		{
			name:     "libvirt replicas",
			pool:     &types.MachinePool{},
			platform: "libvirt",
			expected: func() *types.MachinePool {
				p := defaultMachinePool("")
				p.Replicas = pointer.Int64Ptr(1)
				return p
			}(),
		},
		{
			name: "non-default hyperthreading",
			pool: func() *types.MachinePool {
				p := defaultMachinePool("test-name")
				p.Hyperthreading = types.HyperthreadingMode("test-hyperthreading")
				return p
			}(),
			expected: func() *types.MachinePool {
				p := defaultMachinePool("test-name")
				p.Hyperthreading = types.HyperthreadingMode("test-hyperthreading")
				return p
			}(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetMachinePoolDefaults(tc.pool, tc.platform)
			assert.Equal(t, tc.expected, tc.pool, "unexpected machine pool")
		})
	}
}
