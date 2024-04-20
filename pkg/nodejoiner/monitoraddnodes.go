package nodejoiner

import (
	"context"

	"github.com/sirupsen/logrus"

	agentpkg "github.com/openshift/installer/pkg/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

// NewMonitorAddNodesCommand creates a new command for monitor add nodes.
func NewMonitorAddNodesCommand(directory, kubeconfigPath string, autoApproveCSRs bool, ips []string) error {
	cluster, err := agentpkg.NewCluster(context.Background(), "", ips[0], kubeconfigPath, workflow.AgentWorkflowTypeAddNodes)
	if err != nil {
		// TODO exit code enumerate
		logrus.Exit(1)
	}

	return agentpkg.MonitorAddNodes(cluster, autoApproveCSRs, ips[0])
}
