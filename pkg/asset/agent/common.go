package agent

import (
	"fmt"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/yaml"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/vsphere"
)

type nmStateConfig struct {
	Interfaces []struct {
		IPV4 struct {
			Address []struct {
				IP string `yaml:"ip,omitempty"`
			} `yaml:"address,omitempty"`
		} `yaml:"ipv4,omitempty"`
		IPV6 struct {
			Address []struct {
				IP string `yaml:"ip,omitempty"`
			} `yaml:"address,omitempty"`
		} `yaml:"ipv6,omitempty"`
	} `yaml:"interfaces,omitempty"`
}

const (
	// ExternalPlatformNameOci is the name of the external platform for OCP.
	ExternalPlatformNameOci = "oci"
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
		if strings.Contains(err.Error(), "unable to read image") {
			logrus.Debugf("Could not get release image to check architecture in a disconnected setup")
		}
		return "", err
	}
	logrus.Debugf("Release Image arch is: %s", releaseArch)
	return releaseArch, nil
}

// GetUserManagedNetworkingByPlatformType returns the expected value for userManagedNetworking
// based on the current platform type.
func GetUserManagedNetworkingByPlatformType(platformType hiveext.PlatformType) *bool {
	switch platformType {
	case hiveext.NonePlatformType, hiveext.ExternalPlatformType:
		logrus.Debugf("Setting UserManagedNetworking to true for %s platform", platformType)
		return swag.Bool(true)
	default:
		return swag.Bool(false)
	}
}

// GetFirstIP returns the firt IP found in the nmstate configuration for this host.
func GetFirstIP(nmstateRaw []byte) (string, error) {
	var nmStateConfig nmStateConfig
	err := yaml.Unmarshal(nmstateRaw, &nmStateConfig)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling NMStateConfig: %w", err)
	}

	for _, intf := range nmStateConfig.Interfaces {
		for _, addr4 := range intf.IPV4.Address {
			if addr4.IP != "" {
				return addr4.IP, nil
			}
		}
		for _, addr6 := range intf.IPV6.Address {
			if addr6.IP != "" {
				return addr6.IP, nil
			}
		}
	}

	return "", nil
}

// GetAllHostIPs returns a map of host IPs from the nmstate configuration for this host.
func GetAllHostIPs(config aiv1beta1.NetConfig) (map[string]struct{}, error) {
	var nmStateConfig nmStateConfig
	hostIPs := make(map[string]struct{})

	err := yaml.Unmarshal(config.Raw, &nmStateConfig)
	if err != nil {
		return hostIPs, fmt.Errorf("error unmarshalling NMStateConfig: %w", err)
	}

	for _, intf := range nmStateConfig.Interfaces {
		for _, addr4 := range intf.IPV4.Address {
			if addr4.IP != "" {
				hostIPs[addr4.IP] = struct{}{}
			}
		}
		for _, addr6 := range intf.IPV6.Address {
			if addr6.IP != "" {
				hostIPs[addr6.IP] = struct{}{}
			}
		}
	}
	return hostIPs, nil
}
