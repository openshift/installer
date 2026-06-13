//go:build scos

package rhcos

import "github.com/openshift/installer/pkg/types"

// defaultOSImageStream is the SCOS stream name.
const defaultOSImageStream types.OSImageStream = types.OSImageStreamCentos10

// GetDefaultOSImageStream returns the default OS image stream.
// SCOS only has a single stream so the install-config is ignored.
func GetDefaultOSImageStream(_ *types.InstallConfig) types.OSImageStream {
	return defaultOSImageStream
}

// BuildDefaultOSImageStream returns the default OS image stream for the
// build version's default FeatureSet, for callers that don't have an
// install-config to read the FeatureSet from.
func BuildDefaultOSImageStream() types.OSImageStream {
	return defaultOSImageStream
}

func getStreamFileName(_ types.OSImageStream) string {
	return "coreos/scos.json"
}

func getMarketplaceStreamFileName(_ types.OSImageStream) string {
	// There is no current need for scos marketplace images,
	// so this file does not currently exist. The calling
	// functions will gracefully handle the missing file.
	return "coreos/marketplace/marketplace-scos.json"
}
