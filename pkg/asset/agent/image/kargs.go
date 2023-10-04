package image

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
)

// Kargs is an Asset the generates the additional kernel args.
type Kargs struct {
	fips bool
}

// Dependencies returns the assets on which the AgentArtifacts asset depends.
func (a *Kargs) Dependencies() []asset.Asset {
	return []asset.Asset{
		&manifests.AgentClusterInstall{},
	}
}

// Generate generates the configurations for the agent ISO image and PXE assets.
func (a *Kargs) Generate(dependencies asset.Parents) error {
	agentClusterInstall := &manifests.AgentClusterInstall{}
	dependencies.Get(agentClusterInstall)

	a.fips = agentClusterInstall.FIPSEnabled()

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Kargs) Name() string {
	return "Agent ISO Kernel Arguments"
}

// KernelCmdLine returns the data to be appended to the kernel arguments.
func (a *Kargs) KernelCmdLine() []byte {
	if a.fips {
		return []byte(" fips=1")
	}
	return nil
}
