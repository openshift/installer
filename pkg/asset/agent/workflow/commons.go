package workflow

// AgentWorkflowType defines the supported
// agent workflows.
type AgentWorkflowType string

const (
	AgentWorkflowTypeInstall  AgentWorkflowType = "install"
	AgentWorkflowTypeAddNodes AgentWorkflowType = "addnodes"

	agentWorkflowFilename = ".agentworkflow"
)
