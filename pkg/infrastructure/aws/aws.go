package aws

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/aws"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return aws.InitializeProvider(installDir)
}
