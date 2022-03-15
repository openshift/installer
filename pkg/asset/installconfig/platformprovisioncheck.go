package installconfig

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	alibabacloudconfig "github.com/openshift/installer/pkg/asset/installconfig/alibabacloud"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	bmconfig "github.com/openshift/installer/pkg/asset/installconfig/baremetal"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ibmcloudconfig "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	osconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	vsconfig "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
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
	case gcp.Name:
		client, err2 := gcpconfig.NewClient(context.TODO())
		if err2 != nil {
			return err2
		}
		err = gcpconfig.ValidatePreExitingPublicDNS(client, ic.Config)
	case ibmcloud.Name:
		client, err2 := ibmcloudconfig.NewClient()
		if err2 != nil {
			return err2
		}
		err = ibmcloudconfig.ValidatePreExitingPublicDNS(client, ic.Config, ic.IBMCloud)
	case openstack.Name:
		err = osconfig.ValidateForProvisioning(ic.Config)
	case vsphere.Name:
		err = vsconfig.ValidateForProvisioning(ic.Config)
	case ovirt.Name:
		err = ovirtconfig.ValidateForProvisioning(ic.Config)
	case alibabacloud.Name:
		client, err2 := ic.AlibabaCloud.Client()
		if err2 != nil {
			return err2
		}
		err = alibabacloudconfig.ValidateForProvisioning(client, ic.Config, ic.AlibabaCloud)
	case powervs.Name:
		err = powervsconfig.ValidateForProvisioning()
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
