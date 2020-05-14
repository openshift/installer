package installconfig

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	azconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	vsconfig "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// PlatformProvisionCheck is an asset that validates the install-config platform for
// any requirements specific for provisioning infrastructure.
type PlatformProvisionCheck struct {
}

var _ asset.Asset = (*PlatformProvisionCheck)(nil)

// Dependencies returns the dependencies for PlatformProvisionCheck
func (a *PlatformProvisionCheck) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate queries for input from the user.
func (a *PlatformProvisionCheck) Generate(dependencies asset.Parents) error {
	ic := &InstallConfig{}
	dependencies.Get(ic)

	var err error
	platform := ic.Config.Platform.Name()
	switch platform {
	case vsphere.Name:
		err = vsconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	case azure.Name:
		err = azconfig.ValidatePublicDNS(ic.Config)
		if err != nil {
			return err
		}
	case gcp.Name:
		client, err := gcpconfig.NewClient(context.TODO())
		if err != nil {
			return err
		}
		err = gcpconfig.ValidatePreExitingPublicDNS(client, ic.Config)
		if err != nil {
			return err
		}
	case aws.Name, baremetal.Name, libvirt.Name, none.Name, openstack.Name, ovirt.Name:
		// no special provisioning requirements to check
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}
	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformProvisionCheck) Name() string {
	return "Platform Provisioning Check"
}
