//go:build !(okd || scos)

package rhcos

import (
	"fmt"

	"github.com/openshift/installer/pkg/types"
)

const (
	// DefaultOSImageStream is the OS image stream used when the install-config
	// does not specify one.
	DefaultOSImageStream = types.OSImageStreamRHCOS9

	payloadImageStreamTagRHCOS9  = "rhel-coreos"
	payloadImageStreamTagRHCOS10 = "rhel-coreos-10"
)

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

// GetPayloadImageStreamTag returns the payload image stream tag corresponding
// to the given OS image stream.
func GetPayloadImageStreamTag(stream types.OSImageStream) string {
	if stream == "" || stream == types.OSImageStreamRHCOS9 {
		return payloadImageStreamTagRHCOS9
	}
	return payloadImageStreamTagRHCOS10
}
