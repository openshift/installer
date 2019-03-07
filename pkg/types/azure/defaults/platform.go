package defaults

import (
	"github.com/openshift/installer/pkg/types/azure"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *azure.Platform) {
}

// InstanceClass returns the instance "class" we should use for a given
// region.
func InstanceClass(region string) string {
	return "Standard_DS2_v2"
}
