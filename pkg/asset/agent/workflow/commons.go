package workflow

// AgentWorkflowType defines the supported
// agent workflows.
type AgentWorkflowType string

const (
	// AgentWorkflowTypeInstall identifies the install workflow.
	AgentWorkflowTypeInstall AgentWorkflowType = "install"
	// AgentWorkflowTypeAddNodes identifies the add nodes workflow.
	AgentWorkflowTypeAddNodes AgentWorkflowType = "addnodes"

	agentWorkflowFilename = ".agentworkflow"
)
