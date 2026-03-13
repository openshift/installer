//go:build scos

package rhcos

import "github.com/openshift/installer/pkg/types"

// DefaultOSImageStream Not used in SCOS
const DefaultOSImageStream types.OSImageStream = ""

func getStreamFileName() string {
	return "coreos/scos.json"
}

func getMarketplaceStreamFileName() string {
	// There is no current need for scos marketplace images,
	// so this file does not currently exist. The calling
	// functions will gracefully handle the missing file.
	return "coreos/marketplace-scos.json"
}
