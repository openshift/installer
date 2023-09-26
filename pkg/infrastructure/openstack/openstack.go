package openstack

import (
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform/stages/openstack"
)

func InitializeProvider(installDir string) ([]infrastructure.Stage, func() error, error) {
	return openstack.InitializeProvider(installDir)
}
