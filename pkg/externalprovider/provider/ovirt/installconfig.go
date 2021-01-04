package ovirt

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
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

func (ovirt *ovirtProvider) ValidateInstallConfig(
	Config *types.InstallConfig,
	_ *asset.File,
	_ *aws.Metadata,
	_ *icazure.Metadata,
) error {
	return ovirt2.Validate(Config)
}
