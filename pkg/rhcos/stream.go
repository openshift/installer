//go:build !(okd || scos)

package rhcos

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
)

// DefaultOSImageStream is the OS image stream used when the install-config
// does not specify one.
const DefaultOSImageStream = types.OSImageStreamRHCOS9

func getStreamFileName(stream types.OSImageStream) string {
	if stream == "" {
		stream = DefaultOSImageStream
	}
	return fmt.Sprintf("coreos/coreos-%v.json", stream)
}

func getMarketplaceStreamFileName(stream types.OSImageStream) string {
	if stream == "" {
		stream = DefaultOSImageStream
	}
	return fmt.Sprintf("coreos/marketplace/coreos-%v.json", stream)
}
