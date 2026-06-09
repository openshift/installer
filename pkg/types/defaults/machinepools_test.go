package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
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
			platform: &types.Platform{},
			expected: defaultMachinePool(""),
		},
		{
			name:     "empty",
			pool:     &types.MachinePool{Replicas: &defaultEdgeReplicaCount},
			platform: &types.Platform{},
			expected: defaultEdgeMachinePool(""),
		},
		{
			name:     "edge",
			pool:     &types.MachinePool{Name: "edge"},
			platform: &types.Platform{},
			expected: defaultEdgeMachinePool("edge"),
		},
		{
			name:     "arbiter",
			pool:     &types.MachinePool{Name: "arbiter"},
			platform: &types.Platform{},
			expected: defaultMachinePoolWithReplicaCount("arbiter", 0),
		},
		{
			name:     "default",
			pool:     defaultMachinePool("test-name"),
			platform: &types.Platform{},
			expected: defaultMachinePool("test-name"),
		},
		{
			name:     "default",
			pool:     defaultEdgeMachinePool("test-name"),
			platform: &types.Platform{},
			expected: defaultEdgeMachinePool("test-name"),
		},
		{
			name:     "non-default replicas",
			platform: &types.Platform{},
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
			name:     "non-default replicas",
			platform: &types.Platform{},
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
			name:     "non-default hyperthreading",
			platform: &types.Platform{},
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
			name:     "non-default hyperthreading",
			platform: &types.Platform{},
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
			// Use default feature set (no special features enabled)
			config := &types.InstallConfig{}
			SetMachinePoolDefaults(tc.pool, tc.platform, config.EnabledFeatureGates())
			assert.Equal(t, tc.expected, tc.pool, "unexpected machine pool")
		})
	}
}

func TestSetMachinePoolDefaultsWithFeatureGates(t *testing.T) {
	cases := []struct {
		name               string
		pool               *types.MachinePool
		platform           *types.Platform
		featureSet         configv1.FeatureSet
		expectedManagement types.MachineManagementAPI
	}{
		{
			name:               "control plane with DevPreviewNoUpgrade feature set",
			pool:               &types.MachinePool{Name: types.MachinePoolControlPlaneRoleName},
			platform:           &types.Platform{},
			featureSet:         configv1.DevPreviewNoUpgrade,
			expectedManagement: types.ClusterAPI,
		},
		{
			name:               "control plane with default feature set",
			pool:               &types.MachinePool{Name: types.MachinePoolControlPlaneRoleName},
			platform:           &types.Platform{},
			featureSet:         configv1.Default,
			expectedManagement: "",
		},
		{
			name:               "compute with DevPreviewNoUpgrade feature set",
			pool:               &types.MachinePool{Name: types.MachinePoolComputeRoleName},
			platform:           &types.Platform{},
			featureSet:         configv1.DevPreviewNoUpgrade,
			expectedManagement: types.ClusterAPI,
		},
		{
			name:               "compute with default feature set",
			pool:               &types.MachinePool{Name: types.MachinePoolComputeRoleName},
			platform:           &types.Platform{},
			featureSet:         configv1.Default,
			expectedManagement: "",
		},
		{
			name:               "control plane with management already set",
			pool:               &types.MachinePool{Name: types.MachinePoolControlPlaneRoleName, Management: types.MachineAPI},
			platform:           &types.Platform{},
			featureSet:         configv1.DevPreviewNoUpgrade,
			expectedManagement: types.MachineAPI,
		},
		{
			name:               "compute with management already set",
			pool:               &types.MachinePool{Name: types.MachinePoolComputeRoleName, Management: types.MachineAPI},
			platform:           &types.Platform{},
			featureSet:         configv1.DevPreviewNoUpgrade,
			expectedManagement: types.MachineAPI,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.InstallConfig{
				FeatureSet: tc.featureSet,
			}
			SetMachinePoolDefaults(tc.pool, tc.platform, config.EnabledFeatureGates())
			assert.Equal(t, tc.expectedManagement, tc.pool.Management, "unexpected management API")
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
			config := &types.InstallConfig{}
			SetMachinePoolDefaults(tc.pool, tc.platform, config.EnabledFeatureGates())
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
