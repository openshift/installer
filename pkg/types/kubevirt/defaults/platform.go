package defaults

import (
	"github.com/openshift/installer/pkg/types/kubevirt"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *kubevirt.Platform) {
	if p.PersistentVolumeAccessMode == "" {
		p.PersistentVolumeAccessMode = "ReadWriteMany"
	}
}
