package libvirt

import (
	"github.com/openshift/installer/pkg/destroy/providers"
)

func init() {
	providers.Registry["libvirt"] = New
}
