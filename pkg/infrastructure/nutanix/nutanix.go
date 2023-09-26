package nutanix

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/nutanix"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return nutanix.InitializeProvider(installDir)
}
