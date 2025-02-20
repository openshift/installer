package workflow

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
)

// AgentWorkflowInstallInteractiveDisconnected is meant just to define
// the add nodes workflow.
type AgentWorkflowInstallInteractiveDisconnected struct {
	AgentWorkflow
}

var _ asset.WritableAsset = (*AgentWorkflowInstallInteractiveDisconnected)(nil)

// Name returns a human friendly name for the asset.
func (*AgentWorkflowInstallInteractiveDisconnected) Name() string {
	return "Agent Workflow Install Interactive Disconnected"
}

// Generate generates the AgentWorkflow asset.
func (a *AgentWorkflowInstallInteractiveDisconnected) Generate(_ context.Context, dependencies asset.Parents) error {
	a.Workflow = AgentWorkflowTypeInstallInteractiveDisconnected
	a.File = &asset.File{
		Filename: agentWorkflowFilename,
		Data:     []byte(a.Workflow),
	}

	return nil
}
