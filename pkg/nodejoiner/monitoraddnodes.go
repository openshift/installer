package nodejoiner

import (
	"context"

	"github.com/sirupsen/logrus"

	agentpkg "github.com/openshift/installer/pkg/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

// NewMonitorAddNodesCommand creates a new command for monitor add nodes.
func NewMonitorAddNodesCommand(directory, kubeconfigPath string, ips []string) error {
	err := saveParams(directory, kubeconfigPath)
	if err != nil {
		return err
	}

	var token string
	cluster, err := agentpkg.NewCluster(context.Background(), "", ips[0], kubeconfigPath, "", token, workflow.AgentWorkflowTypeAddNodes)
	if err != nil {
		// TODO exit code enumerate
		logrus.Exit(1)
	}

	if err != nil {
		return err
	}

	return agentpkg.MonitorAddNodes(cluster, ips[0])
}
