package ibmcloud

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/ibmcloud"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return ibmcloud.InitializeProvider(installDir)
}
