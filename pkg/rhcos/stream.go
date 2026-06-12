//go:build !(okd || scos)

package rhcos

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version/versioninfo"
)

// GetDefaultOSImageStream returns the default OS image stream for the
// build version's default FeatureSet.
func GetDefaultOSImageStream(_ *types.InstallConfig) types.OSImageStream {
	// Note: This function is in place to allow stream
	// selection based on FeatureGates in the future
	return BuildDefaultOSImageStream()
}

// BuildDefaultOSImageStream returns the default OS image stream for the
// build version's default FeatureSet, for callers that don't have an
// install-config to read the FeatureSet from.
func BuildDefaultOSImageStream() types.OSImageStream {
	versionInfo := versioninfo.GetInfo()
	if versionInfo.Major < 5 {
		return types.OSImageStreamRHCOS9
	}
	return types.OSImageStreamRHCOS10
}

func getStreamFileName(stream types.OSImageStream) string {
	return fmt.Sprintf("coreos/coreos-%v.json", stream)
}

func getMarketplaceStreamFileName(stream types.OSImageStream) string {
	return fmt.Sprintf("coreos/marketplace/coreos-%v.json", stream)
}
