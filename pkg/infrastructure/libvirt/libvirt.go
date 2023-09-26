package libvirt

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/libvirt"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return libvirt.InitializeProvider(installDir)
}
