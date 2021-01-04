package ovirt

import (
	ovirt2 "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	"github.com/openshift/installer/pkg/types"
)

func (ovirt *ovirtProvider) AddToInstallConfigPlatform(
	p *types.Platform,
) error {
	platform, err := ovirt2.Platform()
	if err != nil {
		return err
	}
	p.Ovirt = platform
	return nil
}
