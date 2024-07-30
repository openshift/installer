package gencrypto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/common"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

func TestAuthConfig_Generate(t *testing.T) {
	cases := []struct {
		name     string
		workflow workflow.AgentWorkflowType
	}{
		{
			name:     "generate-public-key-and-token-install-workflow",
			workflow: workflow.AgentWorkflowTypeInstall,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			agentWorkflow := &workflow.AgentWorkflow{Workflow: tc.workflow}
			parents := asset.Parents{}
			parents.Add(&common.InfraEnvID{}, agentWorkflow)

			authConfigAsset := &AuthConfig{}
			err := authConfigAsset.Generate(context.Background(), parents)
			assert.NoError(t, err)
			assert.NotEmpty(t, authConfigAsset.PublicKey)
			assert.NotEmpty(t, authConfigAsset.AgentAuthToken)
		})
	}
}
