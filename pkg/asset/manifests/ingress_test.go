package manifests

import (
	"testing"

	"github.com/ghodss/yaml"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

// installConfigFromTopologies generates an install config that would yield the
// given topologies when determineTopologies is called on it
func installConfigFromTopologies(t *testing.T, options []icOption,
	controlPlaneTopology configv1.TopologyMode, infrastructureTopology configv1.TopologyMode) *types.InstallConfig {
	installConfig := icBuild.build(options...)

	highlyAvailable := int64(3)
	singleReplica := int64(1)

	switch controlPlaneTopology {
	case configv1.HighlyAvailableTopologyMode:
		installConfig.ControlPlane = &types.MachinePool{
			Replicas: &highlyAvailable,
		}
	case configv1.SingleReplicaTopologyMode:
		installConfig.ControlPlane = &types.MachinePool{
			Replicas: &singleReplica,
		}
	}

	switch infrastructureTopology {
	case configv1.HighlyAvailableTopologyMode:
		installConfig.Compute = []types.MachinePool{
			{Replicas: &highlyAvailable},
		}
	case configv1.SingleReplicaTopologyMode:
		installConfig.Compute = []types.MachinePool{
			{Replicas: &singleReplica},
		}
	}

	// Assert that this function actually works
	generatedControlPlaneTopology, generatedInfrastructureTopology := determineTopologies(installConfig)
	assert.Equal(t, generatedControlPlaneTopology, controlPlaneTopology)
	assert.Equal(t, generatedInfrastructureTopology, infrastructureTopology)

	return installConfig
}

func TestGenerateIngerssDefaultPlacement(t *testing.T) {
	cases := []struct {
		name                      string
		installConfigBuildOptions []icOption
		controlPlaneTopology      configv1.TopologyMode
		infrastructureTopology    configv1.TopologyMode
		expectedIngressPlacement  configv1.DefaultPlacement
	}{
		{
			// AWS currently uses a load balancer even on single-node, so the
			// default placement should be workers
			name:                      "aws single node with 0 or 1 day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                      "aws multi-node with 1 day-1 worker",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			// AWS currently uses a load balancer even on single-node, so the
			// default placement should be workers
			name:                      "aws single-node with multiple day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                      "vanilla aws",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                      "none-platform single node with 0 or 1 day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forNone()},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementControlPlane,
		},
		{
			name:                      "none-platform multi-node with 1 day-1 worker",
			installConfigBuildOptions: []icOption{icBuild.forNone()},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			// For the sake of consistency, we want ingress traffic to go
			// through the single control plane node even when there are
			// workers on day 1. This is even though it would make sense
			// for the user to want to set up a day-1 load balancer in this
			// situation for highly available ingress.
			name:                      "none-platform single-node with multiple day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forNone()},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementControlPlane,
		},
		{
			name:                      "vanilla none-platform",
			installConfigBuildOptions: []icOption{icBuild.forNone()},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(
				&installconfig.ClusterID{
					UUID:    "test-uuid",
					InfraID: "test-infra-id",
				},
				&installconfig.InstallConfig{
					Config: installConfigFromTopologies(t, tc.installConfigBuildOptions,
						tc.controlPlaneTopology, tc.infrastructureTopology),
				},
			)
			ingressAsset := &Ingress{}
			err := ingressAsset.Generate(parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}
			if !assert.Len(t, ingressAsset.FileList, 1, "expected only one file to be generated") {
				return
			}
			assert.Equal(t, ingressAsset.FileList[0].Filename, "manifests/cluster-ingress-02-config.yml")
			var actualIngress configv1.Ingress
			err = yaml.Unmarshal(ingressAsset.FileList[0].Data, &actualIngress)
			if !assert.NoError(t, err, "failed to unmarshal infra manifest") {
				return
			}
			assert.Equal(t, tc.expectedIngressPlacement, actualIngress.Status.DefaultPlacement)
		})
	}
}
