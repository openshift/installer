package vsphere

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/vsphere"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return vsphere.InitializeProvider(installDir)
}
