package image

import (
	"github.com/sirupsen/logrus"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
)

// Kargs is an Asset that generates the additional kernel args.
type Kargs struct {
	fips        bool
	consoleArgs string
}

// Dependencies returns the assets on which the Kargs asset depends.
func (a *Kargs) Dependencies() []asset.Asset {
	return []asset.Asset{
		&manifests.AgentClusterInstall{},
	}
}

// Generate generates the kernel args configurations for the agent ISO image and PXE assets.
func (a *Kargs) Generate(dependencies asset.Parents) error {
	agentClusterInstall := &manifests.AgentClusterInstall{}
	dependencies.Get(agentClusterInstall)

	a.fips = agentClusterInstall.FIPSEnabled()

	// Add kernel args for external oci platform
	if agentClusterInstall.GetExternalPlatformName() == string(models.PlatformTypeOci) {
		logrus.Debugf("Added kernel args to enable serial console for %s %s platform", hiveext.ExternalPlatformType, string(models.PlatformTypeOci))
		a.consoleArgs = " console=ttyS0"
	}
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Kargs) Name() string {
	return "Agent ISO/PXE files Kernel Arguments"
}

// KernelCmdLine returns the data to be appended to the kernel arguments.
func (a *Kargs) KernelCmdLine() []byte {
	if a.fips {
		return []byte(" fips=1")
	}
	return []byte(a.consoleArgs)
}
