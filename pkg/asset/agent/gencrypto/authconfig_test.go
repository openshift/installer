package gencrypto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
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
			parents.Add(agentWorkflow)

			authConfigAsset := &AuthConfig{}
			err := authConfigAsset.Generate(context.Background(), parents)
			assert.NoError(t, err)
			assert.NotEmpty(t, authConfigAsset.PublicKey)
			assert.NotEmpty(t, authConfigAsset.AgentAuthToken)

			// All the 3 tokens should be unique
			assert.NotEqual(t, authConfigAsset.AgentAuthToken, authConfigAsset.UserAuthToken)
			assert.NotEqual(t, authConfigAsset.AgentAuthToken, authConfigAsset.WatcherAuthToken)
			assert.NotEqual(t, authConfigAsset.UserAuthToken, authConfigAsset.WatcherAuthToken)

			// verify each token is signed with correct persona
			claims , err := ParseToken(authConfigAsset.AgentAuthToken)
			assert.NoError(t, err)
			persona, _ := claims["sub"].(string)
			assert.Equal(t, persona, agentPersona)

			claims , err = ParseToken(authConfigAsset.UserAuthToken)
			assert.NoError(t, err)
			persona, _ = claims["sub"].(string)
			assert.Equal(t, persona, userPersona)

			claims , err = ParseToken(authConfigAsset.WatcherAuthToken)
			assert.NoError(t, err)
			persona, _ = claims["sub"].(string)
			assert.Equal(t, persona, watcherPersona)
			
		})
	}
}
