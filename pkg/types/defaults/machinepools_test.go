package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

func defaultMachinePool(name string) *types.MachinePool {
	repCount := int64(3)
	return &types.MachinePool{
		Name:           name,
		Replicas:       &repCount,
		Hyperthreading: types.HyperthreadingEnabled,
		Architecture:   types.ArchitectureAMD64,
	}
}

func defaultEdgeMachinePool(name string) *types.MachinePool {
	pool := defaultMachinePool(name)
	defaultEdgeReplicaCount := int64(0)
	pool.Replicas = &defaultEdgeReplicaCount
	return pool
}

func TestSetMahcinePoolDefaults(t *testing.T) {
	defaultEdgeReplicaCount := int64(0)
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
			name:     "empty",
			pool:     &types.MachinePool{Replicas: &defaultEdgeReplicaCount},
			expected: defaultEdgeMachinePool(""),
		},
		{
			name:     "default",
			pool:     defaultMachinePool("test-name"),
			expected: defaultMachinePool("test-name"),
		},
		{
			name:     "default",
			pool:     defaultEdgeMachinePool("test-name"),
			expected: defaultEdgeMachinePool("test-name"),
		},
		{
			name: "non-default replicas",
			pool: func() *types.MachinePool {
				p := defaultMachinePool("test-name")
				repCount := int64(5)
				p.Replicas = &repCount
				return p
			}(),
			expected: func() *types.MachinePool {
				p := defaultMachinePool("test-name")
				repCount := int64(5)
				p.Replicas = &repCount
				return p
			}(),
		},
		{
			name: "non-default replicas",
			pool: func() *types.MachinePool {
				p := defaultEdgeMachinePool("test-name")
				repCount := int64(5)
				p.Replicas = &repCount
				return p
			}(),
			expected: func() *types.MachinePool {
				p := defaultEdgeMachinePool("test-name")
				repCount := int64(5)
				p.Replicas = &repCount
				return p
			}(),
		},
		{
			name:     "libvirt replicas",
			pool:     &types.MachinePool{},
			platform: "libvirt",
			expected: func() *types.MachinePool {
				p := defaultMachinePool("")
				repCount := int64(1)
				p.Replicas = &repCount
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
		{
			name: "non-default hyperthreading",
			pool: func() *types.MachinePool {
				p := defaultEdgeMachinePool("test-name")
				p.Hyperthreading = types.HyperthreadingMode("test-hyperthreading")
				return p
			}(),
			expected: func() *types.MachinePool {
				p := defaultEdgeMachinePool("test-name")
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
