package image

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

// Kargs is an Asset that generates the additional kernel args.
type Kargs struct {
	consoleArgs string
	fips        bool
}

// Dependencies returns the assets on which the Kargs asset depends.
func (a *Kargs) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&joiner.ClusterInfo{},
		&manifests.AgentClusterInstall{},
	}
}

// Generate generates the kernel args configurations for the agent ISO image and PXE assets.
func (a *Kargs) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	clusterInfo := &joiner.ClusterInfo{}
	agentClusterInstall := &manifests.AgentClusterInstall{}
	dependencies.Get(agentClusterInstall, agentWorkflow, clusterInfo)

	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		a.fips = agentClusterInstall.FIPSEnabled()
		// Add kernel args for external oci platform
		if agentClusterInstall.GetExternalPlatformName() == agent.ExternalPlatformNameOci {
			logrus.Debugf("Added kernel args to enable serial console for %s %s platform", hiveext.ExternalPlatformType, agent.ExternalPlatformNameOci)
			a.consoleArgs = " console=ttyS0"
		}

	case workflow.AgentWorkflowTypeAddNodes:
		a.fips = clusterInfo.FIPS

	default:
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Kargs) Name() string {
	return "Agent ISO/PXE files Kernel Arguments"
}

// KernelCmdLine returns the data to be appended to the kernel arguments.
func (a *Kargs) KernelCmdLine() string {
	cmdLine := a.consoleArgs
	if a.fips {
		cmdLine += " fips=1"
	}
	return cmdLine
}
