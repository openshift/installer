package powervs

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/powervs"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return powervs.InitializeProvider(installDir)
}
