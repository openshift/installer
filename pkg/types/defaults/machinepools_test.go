package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
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
		featureGates       []string
		expectedManagement types.MachineManagementAPI
	}{
		{
			name:               "AWS compute sets ClusterAPI management when ClusterAPIMachineManagementAWS enabled",
			pool:               &types.MachinePool{Name: types.MachinePoolComputeRoleName},
			platform:           &types.Platform{AWS: &aws.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=True"},
			expectedManagement: types.ClusterAPI,
		},
		{
			name:               "AWS compute management is unchanged when ClusterAPIMachineManagementAWS disabled",
			pool:               &types.MachinePool{Name: types.MachinePoolComputeRoleName},
			platform:           &types.Platform{AWS: &aws.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=False"},
			expectedManagement: "",
		},
		{
			name:               "AWS edge compute sets ClusterAPI management when ClusterAPIMachineManagementAWS enabled",
			pool:               &types.MachinePool{Name: types.MachinePoolEdgeRoleName},
			platform:           &types.Platform{AWS: &aws.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=True"},
			expectedManagement: types.ClusterAPI,
		},
		{
			name:               "AWS edge compute management is unchanged when ClusterAPIMachineManagementAWS disabled",
			pool:               &types.MachinePool{Name: types.MachinePoolEdgeRoleName},
			platform:           &types.Platform{AWS: &aws.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=False"},
			expectedManagement: "",
		},
		{
			name:               "AWS control plane is unaffected by ClusterAPIMachineManagementAWS",
			pool:               &types.MachinePool{Name: types.MachinePoolControlPlaneRoleName},
			platform:           &types.Platform{AWS: &aws.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=True"},
			expectedManagement: "",
		},
		{
			name:               "non-AWS compute is unaffected by ClusterAPIMachineManagementAWS",
			pool:               &types.MachinePool{Name: types.MachinePoolComputeRoleName},
			platform:           &types.Platform{GCP: &gcp.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=True"},
			expectedManagement: "",
		},
		{
			name:               "AWS compute management is not overwritten when already set",
			pool:               &types.MachinePool{Name: types.MachinePoolComputeRoleName, Management: types.MachineAPI},
			platform:           &types.Platform{AWS: &aws.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=True"},
			expectedManagement: types.MachineAPI,
		},
		{
			name:               "AWS edge compute management is not overwritten when already set",
			pool:               &types.MachinePool{Name: types.MachinePoolEdgeRoleName, Management: types.MachineAPI},
			platform:           &types.Platform{AWS: &aws.Platform{}},
			featureGates:       []string{"ClusterAPIMachineManagementAWS=True"},
			expectedManagement: types.MachineAPI,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.InstallConfig{
				FeatureSet:   configv1.CustomNoUpgrade,
				FeatureGates: tc.featureGates,
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

// TestSetMachinePoolDefaultsAzureDataDiskCaching checks that the Azure data disk
// defaults are applied through the generic machine pool defaulting, after the
// defaultMachinePlatform has been merged into the pool.
func TestSetMachinePoolDefaultsAzureDataDiskCaching(t *testing.T) {
	cases := []struct {
		name     string
		pool     *types.MachinePool
		platform *types.Platform
		expected []capz.DataDisk
	}{
		{
			name: "omitted cachingType defaults to None",
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{NameSuffix: "etcd", DiskSizeGB: 256, Lun: ptr.To(int32(0))}},
					},
				},
			},
			platform: &types.Platform{Azure: &azure.Platform{}},
			expected: []capz.DataDisk{{NameSuffix: "etcd", DiskSizeGB: 256, Lun: ptr.To(int32(0)), CachingType: "None"}},
		},
		{
			name: "explicit cachingType is preserved",
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{NameSuffix: "data", DiskSizeGB: 64, Lun: ptr.To(int32(0)), CachingType: "ReadOnly"}},
					},
				},
			},
			platform: &types.Platform{Azure: &azure.Platform{}},
			expected: []capz.DataDisk{{NameSuffix: "data", DiskSizeGB: 64, Lun: ptr.To(int32(0)), CachingType: "ReadOnly"}},
		},
		{
			name: "data disks inherited from defaultMachinePlatform are defaulted",
			pool: &types.MachinePool{Name: "worker"},
			platform: &types.Platform{
				Azure: &azure.Platform{
					DefaultMachinePlatform: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{NameSuffix: "data", DiskSizeGB: 64, Lun: ptr.To(int32(0))}},
					},
				},
			},
			expected: []capz.DataDisk{{NameSuffix: "data", DiskSizeGB: 64, Lun: ptr.To(int32(0)), CachingType: "None"}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.InstallConfig{}
			SetMachinePoolDefaults(tc.pool, tc.platform, config.EnabledFeatureGates())
			assert.Equal(t, tc.expected, tc.pool.Platform.Azure.DataDisks)
			if dmp := tc.platform.Azure.DefaultMachinePlatform; dmp != nil {
				for _, dataDisk := range dmp.DataDisks {
					assert.Empty(t, dataDisk.CachingType, "defaulting must not leak into defaultMachinePlatform")
				}
			}
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
