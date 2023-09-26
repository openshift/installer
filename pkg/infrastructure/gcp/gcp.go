package gcp

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/gcp"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return gcp.InitializeProvider(installDir)
}
