package manifests

import (
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
)

func TestGetInfrastructureTopology(t *testing.T) {
	cases := []struct {
		name                 string
		installConfig        *types.InstallConfig
		expectedTopologyMode configv1.TopologyMode
	}{{
		name:                 "no workers, no control-plane replicas",
		installConfig:        icBuild.build(),
		expectedTopologyMode: configv1.HighlyAvailableTopologyMode,
	}, {
		name: "no workers, 1 control-plane replica",
		installConfig: icBuild.build(
			icBuild.withControlPlaneReplicas(1),
		),
		expectedTopologyMode: configv1.SingleReplicaTopologyMode,
	}, {
		name: "no workers, 3 control-plane replicas",
		installConfig: icBuild.build(
			icBuild.withControlPlaneReplicas(3),
		),
		expectedTopologyMode: configv1.HighlyAvailableTopologyMode,
	}, {
		name: "1 worker",
		installConfig: icBuild.build(
			icBuild.withWorkerReplicas(1),
		),
		expectedTopologyMode: configv1.SingleReplicaTopologyMode,
	}, {
		name: "2 workers",
		installConfig: icBuild.build(
			icBuild.withWorkerReplicas(2),
		),
		expectedTopologyMode: configv1.HighlyAvailableTopologyMode,
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedTopologyMode, getInfrastructureTopology(tc.installConfig))
		})
	}
}

func (b icBuildNamespace) withControlPlaneReplicas(replicas int) icOption {
	return func(ic *types.InstallConfig) {
		i := int64(replicas)
		ic.ControlPlane = &types.MachinePool{
			Replicas: &i,
		}
	}
}

func (b icBuildNamespace) withWorkerReplicas(replicas int) icOption {
	return func(ic *types.InstallConfig) {
		i := int64(replicas)
		ic.Compute = []types.MachinePool{{
			Replicas: &i,
		}}
	}
}
