package powervc

import (
	"github.com/openshift/installer/pkg/types/powervc"
)

// Platform collects powervc-specific configuration.
func Platform() (*powervc.Platform, error) {
	var p powervc.Platform

	return &p, nil
}
