package baremetal

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/baremetal"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return baremetal.InitializeProvider(installDir)
}
