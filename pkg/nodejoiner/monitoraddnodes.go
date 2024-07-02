package nodejoiner

import (
	"context"

	agentpkg "github.com/openshift/installer/pkg/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

// NewMonitorAddNodesCommand creates a new command for monitor add nodes.
func NewMonitorAddNodesCommand(directory, kubeconfigPath string, ips []string) error {
	err := saveParams(directory, kubeconfigPath)
	if err != nil {
		return err
	}

	// sshKey is not required parameter for monitor-add-nodes
	sshKey := ""

	clusters := []*agentpkg.Cluster{}
	for _, ip := range ips {
		cluster, err := agentpkg.NewCluster(context.Background(), directory, ip, kubeconfigPath, sshKey, workflow.AgentWorkflowTypeAddNodes)
		if err != nil {
			return err
		}
		clusters = append(clusters, cluster)
	}
	agentpkg.MonitorAddNodes(clusters, ips)

	return nil
}
