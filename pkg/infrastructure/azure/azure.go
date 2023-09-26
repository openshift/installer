package azure

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/azure"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return azure.InitializeProvider(installDir)
}

func InitializeStackProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return azure.InitializeStackProvider(installDir)
}
