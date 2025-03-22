package manifests

import (
	"testing"

	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
)

func Test_DetermineTopologies(t *testing.T) {
	testCases := []struct {
		desc                 string
		installConfig        *types.InstallConfig
		expectedControlPlane configv1.TopologyMode
		expectedInfra        configv1.TopologyMode
	}{
		{
			desc: "should default to HA for both infra and control plane when 3 control replicas",
			installConfig: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Replicas: ptr.To[int64](3),
				},
				Compute: []types.MachinePool{},
			},
			expectedControlPlane: configv1.HighlyAvailableTopologyMode,
			expectedInfra:        configv1.HighlyAvailableTopologyMode,
		},
		{
			desc: "should default to Single for both infra and control plane when 1 control replicas",
			installConfig: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Replicas: ptr.To[int64](1),
				},
				Compute: []types.MachinePool{},
			},
			expectedControlPlane: configv1.SingleReplicaTopologyMode,
			expectedInfra:        configv1.SingleReplicaTopologyMode,
		},
		{
			desc: "should default infra to HA and controlPlane to Single for 1 control replicas and 3 worker replicas",
			installConfig: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Replicas: ptr.To[int64](1),
				},
				Compute: []types.MachinePool{
					{
						Replicas: ptr.To[int64](3),
					},
				},
			},
			expectedControlPlane: configv1.SingleReplicaTopologyMode,
			expectedInfra:        configv1.HighlyAvailableTopologyMode,
		},
		{
			desc: "should default infra to Single and controlPlane to Single for 1 control replicas and 1 worker replicas",
			installConfig: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Replicas: ptr.To[int64](1),
				},
				Compute: []types.MachinePool{
					{
						Replicas: ptr.To[int64](1),
					},
				},
			},
			expectedControlPlane: configv1.SingleReplicaTopologyMode,
			expectedInfra:        configv1.SingleReplicaTopologyMode,
		},
		{
			desc: "should default infra to HA and controlPlane to DualReplica for 2 control replicas",
			installConfig: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Replicas: ptr.To[int64](2),
				},
				Compute: []types.MachinePool{},
			},
			expectedControlPlane: configv1.DualReplicaTopologyMode,
			expectedInfra:        configv1.HighlyAvailableTopologyMode,
		},
		{
			desc: "should default infra to HA and controlPlane to Arbiter for 2 control replicas and 1 arbiter",
			installConfig: &types.InstallConfig{
				Arbiter: &types.MachinePool{
					Replicas: ptr.To[int64](1),
				},
				ControlPlane: &types.MachinePool{
					Replicas: ptr.To[int64](2),
				},
				Compute: []types.MachinePool{},
			},
			expectedControlPlane: configv1.HighlyAvailableArbiterMode,
			expectedInfra:        configv1.HighlyAvailableTopologyMode,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			controlPlane, infra := determineTopologies(tc.installConfig)
			if controlPlane != tc.expectedControlPlane {
				t.Fatalf("expected control plane topology to be %s but got %s", tc.expectedControlPlane, controlPlane)
			}
			if infra != tc.expectedInfra {
				t.Fatalf("expected infra topology to be %s but got %s", tc.expectedInfra, infra)
			}
		})
	}
}

func Test_DetermineCPUPartitioning(t *testing.T) {
	testCases := []struct {
		desc                     string
		installConfig            *types.InstallConfig
		expectedPartitioningMode configv1.CPUPartitioningMode
	}{
		{
			desc: "should return AllNodes",
			installConfig: &types.InstallConfig{
				CPUPartitioning: types.CPUPartitioningAllNodes,
			},
			expectedPartitioningMode: configv1.CPUPartitioningAllNodes,
		},
		{
			desc: "should return None",
			installConfig: &types.InstallConfig{
				CPUPartitioning: types.CPUPartitioningNone,
			},
			expectedPartitioningMode: configv1.CPUPartitioningNone,
		},
		{
			desc:                     "should default to None",
			installConfig:            &types.InstallConfig{},
			expectedPartitioningMode: configv1.CPUPartitioningNone,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mode := determineCPUPartitioning(tc.installConfig)
			if mode != tc.expectedPartitioningMode {
				t.Fatalf("expected cpu partitioning mode to be %s but got %s", tc.expectedPartitioningMode, mode)
			}
		})
	}
}
