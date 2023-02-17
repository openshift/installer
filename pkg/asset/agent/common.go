package agent

import (
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// SupportedInstallerPlatforms lists the supported platforms for agent installer.
func SupportedInstallerPlatforms() []string {
	return []string{baremetal.Name, vsphere.Name, none.Name}
}

var supportedHivePlatforms = []hiveext.PlatformType{
	hiveext.BareMetalPlatformType,
	hiveext.VSpherePlatformType,
	hiveext.NonePlatformType,
}

// SupportedHivePlatforms lists the supported platforms for AgentClusterInstall.
func SupportedHivePlatforms() []string {
	platforms := []string{}
	for _, p := range supportedHivePlatforms {
		platforms = append(platforms, string(p))
	}
	return platforms
}

func HivePlatformType(platform types.Platform) hiveext.PlatformType {
	switch platform.Name() {
	case baremetal.Name:
		return hiveext.BareMetalPlatformType
	case none.Name:
		return hiveext.NonePlatformType
	case vsphere.Name:
		return hiveext.VSpherePlatformType
	}
	return ""
}

// IsSupportedPlatform returns true if provided platform is supported.
// Otherwise, returns false.
func IsSupportedPlatform(platform hiveext.PlatformType) bool {
	for _, p := range supportedHivePlatforms {
		if p == platform {
			return true
		}
	}
	return false
}
