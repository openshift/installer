package agent

import (
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// SupportedPlatforms lists the supported platforms for agent installer
var SupportedPlatforms = []string{baremetal.Name, vsphere.Name, none.Name}

// IsSupportedPlatform returns true if provided platform is baremeral, vsphere or none.
// Otherwise, returns false
func IsSupportedPlatform(platform string) bool {
	for _, p := range SupportedPlatforms {
		if p == platform {
			return true
		}
	}
	return false
}
