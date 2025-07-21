package defaults

import (
	"github.com/openshift/installer/pkg/types/libvirt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

func defaultMachinePoolWithReplicaCount(name string, replicaCount int) *types.MachinePool {
	repCount := int64(replicaCount)
	return &types.MachinePool{
		Name:           name,
		Replicas:       &repCount,
		Hyperthreading: types.HyperthreadingEnabled,
		Architecture:   types.ArchitectureAMD64,
	}
}

func defaultMachinePool(name string) *types.MachinePool {
	return defaultMachinePoolWithReplicaCount(name, 3)
}

func defaultEdgeMachinePool(name string) *types.MachinePool {
	return defaultMachinePoolWithReplicaCount(name, 0)
}

func TestSetMachinePoolDefaults(t *testing.T) {
	defaultEdgeReplicaCount := int64(0)
	cases := []struct {
		name     string
		pool     *types.MachinePool
		config   *types.InstallConfig
		expected *types.MachinePool
	}{
		{
			name:     "empty",
			pool:     &types.MachinePool{},
			config:   &types.InstallConfig{},
			expected: defaultMachinePool(""),
		},
		{
			name:     "empty",
			pool:     &types.MachinePool{Replicas: &defaultEdgeReplicaCount},
			config:   &types.InstallConfig{},
			expected: defaultEdgeMachinePool(""),
		},
		{
			name:     "edge",
			pool:     &types.MachinePool{Name: "edge"},
			expected: defaultEdgeMachinePool("edge"),
		},
		{
			name:     "arbiter",
			pool:     &types.MachinePool{Name: "arbiter"},
			expected: defaultMachinePoolWithReplicaCount("arbiter", 0),
		},
		{
			name:     "default",
			pool:     defaultMachinePool("test-name"),
			config:   &types.InstallConfig{},
			expected: defaultMachinePool("test-name"),
		},
		{
			name:     "default",
			pool:     defaultEdgeMachinePool("test-name"),
			config:   &types.InstallConfig{},
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
			config: &types.InstallConfig{},
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
			config: &types.InstallConfig{},
			expected: func() *types.MachinePool {
				p := defaultEdgeMachinePool("test-name")
				repCount := int64(5)
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
			config: &types.InstallConfig{},
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
			config: &types.InstallConfig{},
			expected: func() *types.MachinePool {
				p := defaultEdgeMachinePool("test-name")
				p.Hyperthreading = types.HyperthreadingMode("test-hyperthreading")
				return p
			}(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetMachinePoolDefaults(tc.pool, tc.config)
			assert.Equal(t, tc.expected, tc.pool, "unexpected machine pool")
		})
	}
}

func TestHasEdgePoolConfig(t *testing.T) {
	cases := []struct {
		name     string
		pool     []types.MachinePool
		expected bool
	}{
		{
			name:     "empty",
			pool:     []types.MachinePool{*defaultMachinePool("non-edge")},
			expected: false,
		}, {
			name:     "worker",
			pool:     []types.MachinePool{*defaultMachinePool("worker")},
			expected: false,
		}, {
			name:     "edge",
			pool:     []types.MachinePool{*defaultEdgeMachinePool("edge")},
			expected: true,
		}, {
			name:     "edge",
			pool:     []types.MachinePool{*defaultEdgeMachinePool("edge"), *defaultMachinePool("non-edge")},
			expected: true,
		}, {
			name:     "edge",
			pool:     []types.MachinePool{*defaultEdgeMachinePool("edge"), *defaultMachinePool("worker")},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := hasEdgePoolConfig(tc.pool)
			assert.Equal(t, tc.expected, res, "unexpected machine pool")
		})
	}
}
