package image

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

// Kargs is an Asset that generates the additional kernel args.
type Kargs struct {
	kargs []string
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
		fips := false
		if agentClusterInstall.Config != nil {
			fips = agentClusterInstall.FIPSEnabled()
		}
		a.kargs = manifests.GetAdditionalKernelArgs(fips, agentClusterInstall.GetExternalPlatformName())

	case workflow.AgentWorkflowTypeAddNodes:
		a.kargs = manifests.GetAdditionalKernelArgs(clusterInfo.FIPS, "")

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
	cmdLine := ""
	for _, karg := range a.kargs {
		cmdLine += " " + karg
	}
	return cmdLine
}
