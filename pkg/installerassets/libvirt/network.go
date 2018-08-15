package libvirt

import (
	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["libvirt/network/interface-name"] = installerassets.ConstantDefault([]byte("tt0"))
	installerassets.Defaults["libvirt/network/node-cidr"] = installerassets.ConstantDefault([]byte("192.168.126.0/24"))
}
