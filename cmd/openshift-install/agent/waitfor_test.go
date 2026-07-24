package agent

import (
	"testing"

	"github.com/openshift/installer/pkg/types"
)

func int64Ptr(i int64) *int64 {
	return &i
}

func TestExtractExpectedNodeCounts(t *testing.T) {
	tests := []struct {
		name            string
		config          *types.InstallConfig
		expectedMasters int
		expectedWorkers int
	}{
		{
			name:            "nil config",
			config:          nil,
			expectedMasters: 0,
			expectedWorkers: 0,
		},
		{
			name: "3 masters and 2 workers",
			config: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: int64Ptr(3),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: int64Ptr(2),
					},
				},
			},
			expectedMasters: 3,
			expectedWorkers: 2,
		},
		{
			name: "compact cluster with 3 masters and 0 workers",
			config: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: int64Ptr(3),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: int64Ptr(0),
					},
				},
			},
			expectedMasters: 3,
			expectedWorkers: 0,
		},
		{
			name: "single node openshift",
			config: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: int64Ptr(1),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: int64Ptr(0),
					},
				},
			},
			expectedMasters: 1,
			expectedWorkers: 0,
		},
		{
			name: "multiple compute pools",
			config: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: int64Ptr(3),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: int64Ptr(2),
					},
					{
						Name:     "worker",
						Replicas: int64Ptr(3),
					},
				},
			},
			expectedMasters: 3,
			expectedWorkers: 5,
		},
		{
			name: "nil replicas fields",
			config: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Name: "master",
				},
				Compute: []types.MachinePool{
					{
						Name: "worker",
					},
				},
			},
			expectedMasters: 0,
			expectedWorkers: 0,
		},
		{
			name: "nil control plane with workers",
			config: &types.InstallConfig{
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: int64Ptr(2),
					},
				},
			},
			expectedMasters: 0,
			expectedWorkers: 2,
		},
		{
			name:            "empty config",
			config:          &types.InstallConfig{},
			expectedMasters: 0,
			expectedWorkers: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			masters, workers := extractExpectedNodeCounts(tt.config)
			if masters != tt.expectedMasters {
				t.Errorf("expected %d masters but got %d", tt.expectedMasters, masters)
			}
			if workers != tt.expectedWorkers {
				t.Errorf("expected %d workers but got %d", tt.expectedWorkers, workers)
			}
		})
	}
}
