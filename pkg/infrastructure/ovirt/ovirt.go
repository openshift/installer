package ovirt

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/ovirt"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return ovirt.InitializeProvider(installDir)
}
