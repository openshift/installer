package powervs

import (
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types"
)

// OSImageNameFromStream returns the PowerVS catalog image name that corresponds
// to the given osImageStream value from the install-config.
// The stream value (e.g. "rhel-9", "rhel-10") is mapped to a catalog image
// name of the form "RHEL-CoreOS-<major>" (e.g. "RHEL-CoreOS-9").
// If the stream value is empty or does not contain a "-", it defaults to
// "RHEL-CoreOS-9".
func OSImageNameFromStream(stream types.OSImageStream) string {
	parts := strings.SplitN(string(stream), "-", 2)
	if len(parts) == 2 && parts[1] != "" {
		return fmt.Sprintf("RHEL-CoreOS-%s", parts[1])
	}
	return "RHEL-CoreOS-9"
}

// ImportedImageName returns the name of the IBMPowerVSImage resource (and the
// resulting workspace image) created when platform.powervs.osImage is set.
func ImportedImageName(clusterID string) string {
	return fmt.Sprintf("rhcos-%s", clusterID)
}

// ImageName returns the image name to be used by machines: the imported image
// when platform.powervs.osImage is set, otherwise the stock catalog image
// matching osImageStream.
func ImageName(clusterID string, config *types.InstallConfig) string {
	if config.PowerVS.OSImage != "" {
		return ImportedImageName(clusterID)
	}
	return OSImageNameFromStream(config.OSImageStream)
}
