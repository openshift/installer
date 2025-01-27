package manifests

import (
	"context"
	"testing"

	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/stretchr/testify/assert"
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
		name                        string
		installConfigBuildOptions   []icOption
		controlPlaneTopology        configv1.TopologyMode
		infrastructureTopology      configv1.TopologyMode
		expectedIngressPlacement    configv1.DefaultPlacement
		expectedIngressAWSLBType    configv1.AWSLBType
		expectedIngressPlatformType configv1.PlatformType
		expectedAllocatedEIP        []operatorv1.EIPAllocation
	}{
		{
			// AWS currently uses a load balancer even on single-node, so the
			// default placement should be workers
			name:                      "aws single node with 0 or 1 day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                      "aws multi-node with 1 day-1 worker",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			// AWS currently uses a load balancer even on single-node, so the
			// default placement should be workers
			name:                      "aws single-node with multiple day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                      "vanilla aws",
			installConfigBuildOptions: []icOption{icBuild.forAWS()},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                        "test setting of aws lb type to NLB",
			installConfigBuildOptions:   []icOption{icBuild.withLBType(configv1.NLB)},
			controlPlaneTopology:        configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:      configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:    configv1.DefaultPlacementWorkers,
			expectedIngressAWSLBType:    configv1.NLB,
			expectedIngressPlatformType: configv1.AWSPlatformType,
		},
		{
			name:                        "test setting of aws lb type to NLB and allocate EIP",
			installConfigBuildOptions:   []icOption{icBuild.withLBTypeAndEIPAllocations(configv1.NLB, []aws.EIPAllocation{"eipalloc-<redacted>,eipalloc-<redacted>,eipalloc-<redacted>,eipalloc-<redacted>,eipalloc-<redacted>"})},
			controlPlaneTopology:        configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:      configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:    configv1.DefaultPlacementWorkers,
			expectedIngressAWSLBType:    configv1.NLB,
			expectedIngressPlatformType: configv1.AWSPlatformType,
			expectedAllocatedEIP:        []operatorv1.EIPAllocation{"eipalloc-<redacted>,eipalloc-<redacted>,eipalloc-<redacted>,eipalloc-<redacted>,eipalloc-<redacted>"},
		},
		{
			name:                        "test setting of aws lb type to Classic",
			installConfigBuildOptions:   []icOption{icBuild.withLBType(configv1.Classic)},
			controlPlaneTopology:        configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:      configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:    configv1.DefaultPlacementWorkers,
			expectedIngressAWSLBType:    configv1.Classic,
			expectedIngressPlatformType: configv1.AWSPlatformType,
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
				installconfig.MakeAsset(
					installConfigFromTopologies(t, tc.installConfigBuildOptions,
						tc.controlPlaneTopology, tc.infrastructureTopology),
				),
			)
			ingressAsset := &Ingress{}
			err := ingressAsset.Generate(context.Background(), parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}

			assert.Equal(t, ingressAsset.FileList[0].Filename, "manifests/cluster-ingress-02-config.yml")
			if len(tc.expectedAllocatedEIP) > 0 {
				if !assert.Len(t, ingressAsset.FileList, 2, "expected only two files to be generated") {
					return
				} else {
					assert.Equal(t, ingressAsset.FileList[0].Filename, "manifests/cluster-ingress-02-config.yml")
					assert.Equal(t, ingressAsset.FileList[1].Filename, "manifests/cluster-ingress-default-ingresscontroller.yaml")
				}

			} else if !assert.Len(t, ingressAsset.FileList, 1, "expected only one file to be generated") {
				return
			} else {
				assert.Equal(t, ingressAsset.FileList[0].Filename, "manifests/cluster-ingress-02-config.yml")
			}

			var actualIngress configv1.Ingress
			err = yaml.Unmarshal(ingressAsset.FileList[0].Data, &actualIngress)
			if !assert.NoError(t, err, "failed to unmarshal infra manifest") {
				return
			}
			assert.Equal(t, tc.expectedIngressPlacement, actualIngress.Status.DefaultPlacement)
			if len(tc.expectedIngressPlatformType) != 0 {
				assert.Equal(t, tc.expectedIngressAWSLBType, actualIngress.Spec.LoadBalancer.Platform.AWS.Type)
				assert.Equal(t, tc.expectedIngressPlatformType, actualIngress.Spec.LoadBalancer.Platform.Type)
			}
			if len(tc.expectedAllocatedEIP) != 0 {
				var actualIngressController operatorv1.IngressController
				err = yaml.Unmarshal(ingressAsset.FileList[1].Data, &actualIngressController)
				assert.Equal(t, tc.expectedAllocatedEIP, actualIngressController.Spec.EndpointPublishingStrategy.LoadBalancer.ProviderParameters.AWS.NetworkLoadBalancerParameters.EIPAllocations)
			}
		})
	}
}
