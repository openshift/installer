package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
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

func TestSetMahcinePoolDefaults(t *testing.T) {
	defaultEdgeReplicaCount := int64(0)
	cases := []struct {
		name     string
		pool     *types.MachinePool
		platform *types.Platform
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

func TestSetMachinePoolDefaultsGCPDefaultMachinePlatform(t *testing.T) {
	cases := []struct {
		name     string
		pool     *types.MachinePool
		platform *types.Platform
		expected *gcp.MachinePool
	}{
		{
			name: "diskType from defaultMachinePlatform applied to empty pool",
			pool: &types.MachinePool{Name: "worker"},
			platform: &types.Platform{
				GCP: &gcp.Platform{
					DefaultMachinePlatform: &gcp.MachinePool{
						OSDisk: gcp.OSDisk{
							DiskType: "pd-ssd",
						},
					},
				},
			},
			expected: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskType: "pd-ssd",
				},
			},
		},
		{
			name: "pool diskType takes precedence over defaultMachinePlatform",
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{
						OSDisk: gcp.OSDisk{
							DiskType: "pd-balanced",
						},
					},
				},
			},
			platform: &types.Platform{
				GCP: &gcp.Platform{
					DefaultMachinePlatform: &gcp.MachinePool{
						OSDisk: gcp.OSDisk{
							DiskType: "pd-ssd",
						},
					},
				},
			},
			expected: &gcp.MachinePool{
				OSDisk: gcp.OSDisk{
					DiskType: "pd-balanced",
				},
			},
		},
		{
			name: "multiple fields from defaultMachinePlatform",
			pool: &types.MachinePool{Name: "master"},
			platform: &types.Platform{
				GCP: &gcp.Platform{
					DefaultMachinePlatform: &gcp.MachinePool{
						InstanceType: "n2-standard-4",
						OSDisk: gcp.OSDisk{
							DiskType:   "pd-ssd",
							DiskSizeGB: 256,
						},
					},
				},
			},
			expected: &gcp.MachinePool{
				InstanceType: "n2-standard-4",
				OSDisk: gcp.OSDisk{
					DiskType:   "pd-ssd",
					DiskSizeGB: 256,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetMachinePoolDefaults(tc.pool, tc.platform)
			assert.Equal(t, tc.expected, tc.pool.Platform.GCP)
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
