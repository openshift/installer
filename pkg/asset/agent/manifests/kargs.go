package manifests

import (
	"github.com/sirupsen/logrus"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset/agent"
)

// GetAdditionalKernelArgs returns the additional kernel arguments that are
// appended to the agent ISO and PXE artifacts. It is shared by the Kargs asset
// and the InfraEnv asset so that the InfraEnv reports the same kernel args
// that are actually used to boot.
func GetAdditionalKernelArgs(fips bool, externalPlatformName string) []string {
	var kargs []string
	// Add kernel args for external oci platform
	if externalPlatformName == agent.ExternalPlatformNameOci {
		logrus.Debugf("Added kernel args to enable serial console for %s %s platform", hiveext.ExternalPlatformType, agent.ExternalPlatformNameOci)
		kargs = append(kargs, "console=ttyS0")
	}
	if fips {
		kargs = append(kargs, "fips=1")
	}
	return kargs
}
