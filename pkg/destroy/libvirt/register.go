// +build libvirt_destroy

package libvirt

import (
	"github.com/openshift/installer/pkg/destroy"
)

func init() {
	destroy.Registry["libvirt"] = New
}
