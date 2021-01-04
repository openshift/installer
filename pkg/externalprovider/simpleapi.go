package externalprovider

import (
	"github.com/openshift/installer/pkg/externalprovider/provider/ovirt"
)

var providerRegistry = NewRegistry()

func init() {
	providerRegistry.Register(ovirt.NewOvirtProvider())
}
