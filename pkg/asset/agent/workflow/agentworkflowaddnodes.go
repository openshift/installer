package workflow

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
)

// AgentWorkflowAddNodes is meant just to define
// the add nodes workflow.
type AgentWorkflowAddNodes struct {
	AgentWorkflow
}

var _ asset.WritableAsset = (*AgentWorkflowAddNodes)(nil)

// Name returns a human friendly name for the asset.
func (*AgentWorkflowAddNodes) Name() string {
	return "Agent Workflow Add Nodes"
}

// Generate generates the AgentWorkflow asset.
func (a *AgentWorkflowAddNodes) Generate(_ context.Context, dependencies asset.Parents) error {
	a.Workflow = AgentWorkflowTypeAddNodes
	a.File = &asset.File{
		Filename: agentWorkflowFilename,
		Data:     []byte(a.Workflow),
	}

	return nil
}
