package alibabacloud

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/alibabacloud"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return alibabacloud.InitializeProvider(installDir)
}
