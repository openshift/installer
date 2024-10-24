package joiner

import (
	"context"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

func TestImportClusterConfig_Generate(t *testing.T) {
	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedConfig ClusterConfig
	}{
		{
			name: "skip",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&ClusterInfo{},
			},
		},
		{
			name: "default",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&ClusterInfo{
					PlatformType: v1beta1.BareMetalPlatformType,
				},
			},
			expectedConfig: ClusterConfig{
				Networking: Networking{
					UserManagedNetworking: swag.Bool(false),
				},
			},
		},
		{
			name: "set UserManagedNetworking for platform None",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&ClusterInfo{
					PlatformType: v1beta1.NonePlatformType,
				},
			},
			expectedConfig: ClusterConfig{
				Networking: Networking{
					UserManagedNetworking: swag.Bool(true),
				},
			},
		},
		{
			name: "set UserManagedNetworking for platform External",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&ClusterInfo{
					PlatformType: v1beta1.ExternalPlatformType,
				},
			},
			expectedConfig: ClusterConfig{
				Networking: Networking{
					UserManagedNetworking: swag.Bool(true),
				},
			},
		},
	}
	for _, tc := range cases {
		icc := &ImportClusterConfig{}

		parents := asset.Parents{}
		parents.Add(tc.dependencies...)
		err := icc.Generate(context.Background(), parents)

		assert.NoError(t, err)
		assert.Equal(t, tc.expectedConfig, icc.Config)
	}
}
