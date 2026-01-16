package manifests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
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

// assetClusterID returns a clusterID asset for tests.
func assetClusterID() *installconfig.ClusterID {
	return &installconfig.ClusterID{
		UUID:    "test-uuid",
		InfraID: "test-infra-id",
	}
}

// awsBYOSubnets returns subnets with roles assigned
// for generating an AWS byo-subnet install config.
func awsBYOSubnets(publish types.PublishingStrategy) []awstypes.Subnet {
	if publish == types.ExternalPublishingStrategy {
		return []awstypes.Subnet{
			{
				ID: "subnet-valid-private-a",
				Roles: []awstypes.SubnetRole{
					{Type: awstypes.ClusterNodeSubnetRole},
					{Type: awstypes.ControlPlaneInternalLBSubnetRole},
				},
			},
			{
				ID: "subnet-valid-private-b",
				Roles: []awstypes.SubnetRole{
					{Type: awstypes.ClusterNodeSubnetRole},
					{Type: awstypes.ControlPlaneInternalLBSubnetRole},
				},
			},
			{
				ID: "subnet-valid-public-a",
				Roles: []awstypes.SubnetRole{
					{Type: awstypes.ControlPlaneExternalLBSubnetRole},
					{Type: awstypes.IngressControllerLBSubnetRole},
				},
			},
			{
				ID: "subnet-valid-public-b",
				Roles: []awstypes.SubnetRole{
					{Type: awstypes.ControlPlaneExternalLBSubnetRole},
					{Type: awstypes.IngressControllerLBSubnetRole},
					{Type: awstypes.BootstrapNodeSubnetRole},
				},
			},
		}
	}
	return []awstypes.Subnet{
		{
			ID: "subnet-valid-private-a",
			Roles: []awstypes.SubnetRole{
				{Type: awstypes.ClusterNodeSubnetRole},
				{Type: awstypes.ControlPlaneInternalLBSubnetRole},
				{Type: awstypes.IngressControllerLBSubnetRole},
				{Type: awstypes.BootstrapNodeSubnetRole},
			},
		},
		{
			ID: "subnet-valid-private-b",
			Roles: []awstypes.SubnetRole{
				{Type: awstypes.ClusterNodeSubnetRole},
				{Type: awstypes.ControlPlaneInternalLBSubnetRole},
				{Type: awstypes.IngressControllerLBSubnetRole},
			},
		},
	}
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
	}{
		{
			// AWS currently uses a load balancer even on single-node, so the
			// default placement should be workers
			name:                      "aws single node with 0 or 1 day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forAWS(), icBuild.withMachineNetwork(defaultMachineNetwork)},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                      "aws multi-node with 1 day-1 worker",
			installConfigBuildOptions: []icOption{icBuild.forAWS(), icBuild.withMachineNetwork(defaultMachineNetwork)},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.SingleReplicaTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			// AWS currently uses a load balancer even on single-node, so the
			// default placement should be workers
			name:                      "aws single-node with multiple day-1 workers",
			installConfigBuildOptions: []icOption{icBuild.forAWS(), icBuild.withMachineNetwork(defaultMachineNetwork)},
			controlPlaneTopology:      configv1.SingleReplicaTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                      "vanilla aws",
			installConfigBuildOptions: []icOption{icBuild.forAWS(), icBuild.withMachineNetwork(defaultMachineNetwork)},
			controlPlaneTopology:      configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:    configv1.HighlyAvailableTopologyMode,
			expectedIngressAWSLBType:  configv1.Classic,
			expectedIngressPlacement:  configv1.DefaultPlacementWorkers,
		},
		{
			name:                        "test setting of aws lb type to NLB",
			installConfigBuildOptions:   []icOption{icBuild.withLBType(configv1.NLB), icBuild.withMachineNetwork(defaultMachineNetwork)},
			controlPlaneTopology:        configv1.HighlyAvailableTopologyMode,
			infrastructureTopology:      configv1.HighlyAvailableTopologyMode,
			expectedIngressPlacement:    configv1.DefaultPlacementWorkers,
			expectedIngressAWSLBType:    configv1.NLB,
			expectedIngressPlatformType: configv1.AWSPlatformType,
		},
		{
			name:                        "test setting of aws lb type to Classic",
			installConfigBuildOptions:   []icOption{icBuild.withLBType(configv1.Classic), icBuild.withMachineNetwork(defaultMachineNetwork)},
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
				assetClusterID(),
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
			if len(tc.expectedIngressPlatformType) != 0 {
				assert.Equal(t, tc.expectedIngressAWSLBType, actualIngress.Spec.LoadBalancer.Platform.AWS.Type)
				assert.Equal(t, tc.expectedIngressPlatformType, actualIngress.Spec.LoadBalancer.Platform.Type)
			}
		})
	}
}

func TestGenerateDefaultIngressController(t *testing.T) {
	cases := []struct {
		name                 string
		installConfig        *types.InstallConfig
		expectedScope        operatorv1.LoadBalancerScope
		expectedAWSLbType    operatorv1.AWSLoadBalancerType
		expectedAWSSubnetIDs []operatorv1.AWSSubnetID
		expectIngressCtr     bool
	}{
		{
			name: "aws platform, managed vpc and public ingress",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withPublish(types.ExternalPublishingStrategy),
				icBuild.withMachineNetwork(defaultMachineNetwork),
			),
		},
		{
			name: "aws platform, managed vpc and private ingress",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withMachineNetwork(defaultMachineNetwork),
			),
			expectedScope:    operatorv1.InternalLoadBalancer,
			expectIngressCtr: true,
		},
		{
			name: "aws platform, byo subnets without roles and public ingress",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withPublish(types.ExternalPublishingStrategy),
				icBuild.withAWSBYOSubnets(
					awstypes.Subnet{ID: "subnet-valid-private-a"},
					awstypes.Subnet{ID: "subnet-valid-private-b"},
				),
			),
		},
		{
			name: "aws platform, byo subnets without roles and private ingress",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withAWSBYOSubnets(
					awstypes.Subnet{ID: "subnet-valid-private-a"},
					awstypes.Subnet{ID: "subnet-valid-private-b"},
				),
			),
			expectedScope:    operatorv1.InternalLoadBalancer,
			expectIngressCtr: true,
		},
		{
			name: "aws platform, byo subnets with roles, public ingress and network lb",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withPublish(types.ExternalPublishingStrategy),
				icBuild.withLBType(configv1.NLB),
				icBuild.withAWSBYOSubnets(awsBYOSubnets(types.ExternalPublishingStrategy)...),
			),
			expectedScope:     operatorv1.ExternalLoadBalancer,
			expectedAWSLbType: operatorv1.AWSNetworkLoadBalancer,
			expectedAWSSubnetIDs: []operatorv1.AWSSubnetID{
				"subnet-valid-public-a",
				"subnet-valid-public-b",
			},
			expectIngressCtr: true,
		},
		{
			name: "aws platform, byo subnets with roles, public ingress and classic lb",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withPublish(types.ExternalPublishingStrategy),
				icBuild.withLBType(configv1.Classic),
				icBuild.withAWSBYOSubnets(awsBYOSubnets(types.ExternalPublishingStrategy)...),
			),
			expectedScope:     operatorv1.ExternalLoadBalancer,
			expectedAWSLbType: operatorv1.AWSClassicLoadBalancer,
			expectedAWSSubnetIDs: []operatorv1.AWSSubnetID{
				"subnet-valid-public-a",
				"subnet-valid-public-b",
			},
			expectIngressCtr: true,
		},
		{
			name: "aws platform, byo subnets with roles, private ingress and network lb",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withLBType(configv1.NLB),
				icBuild.withAWSBYOSubnets(awsBYOSubnets(types.InternalPublishingStrategy)...),
			),
			expectedScope:     operatorv1.InternalLoadBalancer,
			expectedAWSLbType: operatorv1.AWSNetworkLoadBalancer,
			expectedAWSSubnetIDs: []operatorv1.AWSSubnetID{
				"subnet-valid-private-a",
				"subnet-valid-private-b",
			},
			expectIngressCtr: true,
		},
		{
			name: "aws platform, byo subnets with roles, private ingress and classic lb",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withLBType(configv1.Classic),
				icBuild.withAWSBYOSubnets(awsBYOSubnets(types.InternalPublishingStrategy)...),
			),
			expectedScope:     operatorv1.InternalLoadBalancer,
			expectedAWSLbType: operatorv1.AWSClassicLoadBalancer,
			expectedAWSSubnetIDs: []operatorv1.AWSSubnetID{
				"subnet-valid-private-a",
				"subnet-valid-private-b",
			},
			expectIngressCtr: true,
		},
		{
			name: "other platforms, and public ingress",
			installConfig: icBuild.build(
				icBuild.forGCP(),
				icBuild.withPublish(types.ExternalPublishingStrategy),
			),
		},
		{
			name: "other platforms, and private ingress",
			installConfig: icBuild.build(
				icBuild.forGCP(),
				icBuild.withPublish(types.InternalPublishingStrategy),
			),
			expectedScope:    operatorv1.InternalLoadBalancer,
			expectIngressCtr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(assetClusterID(), installconfig.MakeAsset(tc.installConfig))

			ingressAsset := &Ingress{}

			err := ingressAsset.Generate(context.Background(), parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}

			numOfFiles := 1 // ingress manifest
			if tc.expectIngressCtr {
				numOfFiles++ // +1 default ingress controller manifest
			}
			if !assert.Lenf(t, ingressAsset.FileList, numOfFiles, "expected %d file(s) to be generated", numOfFiles) {
				return
			}

			assert.Equal(t, ingressAsset.FileList[0].Filename, "manifests/cluster-ingress-02-config.yml")
			if tc.expectIngressCtr {
				assert.Equal(t, ingressAsset.FileList[1].Filename, "manifests/cluster-ingress-default-ingresscontroller.yaml")

				var actualIngressCtrl operatorv1.IngressController

				err = yaml.Unmarshal(ingressAsset.FileList[1].Data, &actualIngressCtrl)
				if !assert.NoError(t, err, "failed to unmarshal default ingress controller manifest") {
					return
				}

				assert.Equal(t, "default", actualIngressCtrl.Name)
				assert.Equal(t, "openshift-ingress-operator", actualIngressCtrl.Namespace)
				assert.Equal(t, operatorv1.LoadBalancerServiceStrategyType, actualIngressCtrl.Spec.EndpointPublishingStrategy.Type)
				assert.Equal(t, tc.expectedScope, actualIngressCtrl.Spec.EndpointPublishingStrategy.LoadBalancer.Scope)

				// Install case: AWS byo subnets with roles.
				if len(tc.expectedAWSSubnetIDs) > 0 {
					assert.Equal(t, operatorv1.AWSLoadBalancerProvider, actualIngressCtrl.Spec.EndpointPublishingStrategy.LoadBalancer.ProviderParameters.Type)
					assert.Equal(t, tc.expectedAWSLbType, actualIngressCtrl.Spec.EndpointPublishingStrategy.LoadBalancer.ProviderParameters.AWS.Type)

					if tc.expectedAWSLbType == operatorv1.AWSNetworkLoadBalancer {
						assert.Equal(t, tc.expectedAWSSubnetIDs, actualIngressCtrl.Spec.EndpointPublishingStrategy.LoadBalancer.ProviderParameters.AWS.NetworkLoadBalancerParameters.Subnets.IDs)
						assert.Nil(t, actualIngressCtrl.Spec.EndpointPublishingStrategy.LoadBalancer.ProviderParameters.AWS.ClassicLoadBalancerParameters)
					} else {
						assert.Equal(t, tc.expectedAWSSubnetIDs, actualIngressCtrl.Spec.EndpointPublishingStrategy.LoadBalancer.ProviderParameters.AWS.ClassicLoadBalancerParameters.Subnets.IDs)
						assert.Nil(t, actualIngressCtrl.Spec.EndpointPublishingStrategy.LoadBalancer.ProviderParameters.AWS.NetworkLoadBalancerParameters)
					}
				}
			}
		})
	}
}
