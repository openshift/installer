//go:build !(okd || scos)

package rhcos

import "github.com/openshift/installer/pkg/types"

// DefaultOSImageStream is the OS image stream used when the install-config
// does not specify one.
const DefaultOSImageStream = types.OSImageStreamRHCOS9

func getStreamFileName() string {
	return "coreos/rhcos.json"
}

func getMarketplaceStreamFileName() string {
	return "coreos/marketplace-rhcos.json"
}
