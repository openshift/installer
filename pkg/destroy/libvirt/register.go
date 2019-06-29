// +build libvirt

package libvirt

import (
	"github.com/openshift/installer/pkg/providers"
)

func init() {
	providers.Registry["libvirt"] = New
}
