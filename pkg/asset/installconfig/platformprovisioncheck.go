package installconfig

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	bmconfig "github.com/openshift/installer/pkg/asset/installconfig/baremetal"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	kubevirtconfig "github.com/openshift/installer/pkg/asset/installconfig/kubevirt"
	osconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	vsconfig "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/kubevirt"
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
	case aws.Name:
		session, err := ic.AWS.Session(context.TODO())
		if err != nil {
			return err
		}
		return awsconfig.ValidateForProvisioning(session, ic.Config, ic.AWS)
	case azure.Name:
		dnsConfig, err := ic.Azure.DNSConfig()
		if err != nil {
			return err
		}
		err = azconfig.ValidatePublicDNS(ic.Config, dnsConfig)
		if err != nil {
			return err
		}
		client, err := ic.Azure.Client()
		if err != nil {
			return err
		}
		return azconfig.ValidateForProvisioning(client, ic.Config)
	case baremetal.Name:
		err = bmconfig.ValidateProvisioning(ic.Config)
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
	case ibmcloud.Name:
		// TODO: IBM[#91]: add validation for provisioning
	case openstack.Name:
		err = osconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	case vsphere.Name:
		err = vsconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	case kubevirt.Name:
		client, err := kubevirtconfig.NewClient()
		if err != nil {
			return err
		}
		err = kubevirtconfig.ValidateForProvisioning(client)
		if err != nil {
			return err
		}
	case ovirt.Name:
		err = ovirtconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	case libvirt.Name, none.Name:
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
