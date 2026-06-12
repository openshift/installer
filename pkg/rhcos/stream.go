//go:build !(okd || scos)

package rhcos

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
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
	return types.OSImageStreamRHCOS9
}

func getStreamFileName(stream types.OSImageStream) string {
	return fmt.Sprintf("coreos/coreos-%v.json", stream)
}

func getMarketplaceStreamFileName(stream types.OSImageStream) string {
	return fmt.Sprintf("coreos/marketplace/coreos-%v.json", stream)
}
