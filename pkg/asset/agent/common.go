package agent

import (
	"github.com/sirupsen/logrus"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// SupportedInstallerPlatforms lists the supported platforms for agent installer.
func SupportedInstallerPlatforms() []string {
	return []string{baremetal.Name, vsphere.Name, none.Name, external.Name}
}

var supportedHivePlatforms = []hiveext.PlatformType{
	hiveext.BareMetalPlatformType,
	hiveext.VSpherePlatformType,
	hiveext.NonePlatformType,
	hiveext.ExternalPlatformType,
}

// SupportedHivePlatforms lists the supported platforms for AgentClusterInstall.
func SupportedHivePlatforms() []string {
	platforms := []string{}
	for _, p := range supportedHivePlatforms {
		platforms = append(platforms, string(p))
	}
	return platforms
}

// HivePlatformType returns the PlatformType for the ZTP Hive API corresponding
// to the given InstallConfig platform.
func HivePlatformType(platform types.Platform) hiveext.PlatformType {
	switch platform.Name() {
	case baremetal.Name:
		return hiveext.BareMetalPlatformType
	case external.Name:
		return hiveext.ExternalPlatformType
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

// DetermineReleaseImageArch returns the arch of the release image.
func DetermineReleaseImageArch(pullSecret, pullSpec string) (string, error) {
	templateFilter := "-o=go-template={{if and .metadata.metadata (index . \"metadata\" \"metadata\" \"release.openshift.io/architecture\")}}{{index . \"metadata\" \"metadata\" \"release.openshift.io/architecture\"}}{{else}}{{.config.architecture}}{{end}}"
	insecure := "--insecure=true"
	var getReleaseArch = []string{
		"oc",
		"adm",
		"release",
		"info",
		pullSpec,
		templateFilter,
		insecure,
	}

	releaseArch, err := ExecuteOC(pullSecret, getReleaseArch)
	if err != nil {
		logrus.Errorf("Release Image arch could not be found: %s", err)
		return "", err
	}
	logrus.Debugf("Release Image arch is: %s", releaseArch)
	return releaseArch, nil
}
