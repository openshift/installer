//go:build scos

package rhcos

import "github.com/openshift/installer/pkg/types"

// DefaultOSImageStream Not used in SCOS
const DefaultOSImageStream types.OSImageStream = ""

func getStreamFileName(_ types.OSImageStream) string {
	return "coreos/scos.json"
}

func getMarketplaceStreamFileName(_ types.OSImageStream) string {
	// There is no current need for scos marketplace images,
	// so this file does not currently exist. The calling
	// functions will gracefully handle the missing file.
	return "coreos/marketplace/marketplace-scos.json"
}

// GetPayloadImageStreamTag returns the payload image stream tag corresponding
// to the given OS image stream. For SCOS, this always returns "stream-coreos".
func GetPayloadImageStreamTag(_ types.OSImageStream) string {
	return "stream-coreos"
}
