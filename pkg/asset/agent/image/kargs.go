package image

import (
	"github.com/openshift/installer/pkg/asset"
)

// Kargs is an Asset the generates the additional kernel args.
type Kargs struct {
}

// Dependencies returns the assets on which the AgentArtifacts asset depends.
func (a *Kargs) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the configurations for the agent ISO image and PXE assets.
func (a *Kargs) Generate(dependencies asset.Parents) error {
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Kargs) Name() string {
	return "Agent ISO Kernel Arguments"
}

// KernelCmdLine returns the data to be appended to the kernel arguments.
func (a *Kargs) KernelCmdLine() []byte {
	return nil
}
