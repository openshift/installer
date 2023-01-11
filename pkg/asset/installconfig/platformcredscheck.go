package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ibmcloudconfig "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	openstackconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// PlatformCredsCheck is an asset that checks the platform credentials, asks for them or errors out if invalid
// the cluster.
type PlatformCredsCheck struct {
}

var _ asset.Asset = (*PlatformCredsCheck)(nil)

// Dependencies returns the dependencies for PlatformCredsCheck
func (a *PlatformCredsCheck) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate queries for input from the user.
func (a *PlatformCredsCheck) Generate(dependencies asset.Parents) error {
	ctx := context.TODO()
	ic := &InstallConfig{}
	dependencies.Get(ic)

	var err error
	platform := ic.Config.Platform.Name()
	switch platform {
	case alibabacloud.Name:
		_, err = ic.AlibabaCloud.Client()
		if err != nil {
			return errors.Wrap(err, "creating AlibabaCloud Cloud session")
		}

	case aws.Name:
		_, err := ic.AWS.Session(ctx)
		if err != nil {
			return err
		}
	case gcp.Name:
		client, err := gcpconfig.NewClient(context.TODO())
		if err != nil {
			return err
		}

		err = gcpconfig.ValidateCredentialMode(client, ic.Config)
		if err != nil {
			return errors.Wrap(err, "validating credentials")
		}
	case ibmcloud.Name:
		_, err = ibmcloudconfig.NewClient()
		if err != nil {
			return errors.Wrap(err, "creating IBM Cloud session")
		}
	case powervs.Name:
		bxCli, err := powervsconfig.NewBxClient()
		if err != nil {
			return err
		}
		err = bxCli.NewPISession()
		if err != nil {
			return errors.Wrap(err, "creating IBM Cloud session")
		}
	case openstack.Name:
		_, err = openstackconfig.GetSession(ic.Config.Platform.OpenStack.Cloud)
		if err != nil {
			return errors.Wrap(err, "creating OpenStack session")
		}
	case baremetal.Name, libvirt.Name, none.Name, vsphere.Name, nutanix.Name:
		// no creds to check
	case azure.Name:
		azureSession, err := ic.Azure.Session()
		if err != nil {
			return errors.Wrap(err, "creating Azure session")
		}
		if azureSession.Credentials.ClientCertificatePath != "" && ic.Config.CredentialsMode != types.ManualCredentialsMode {
			return fmt.Errorf("authentication with client certificates is only supported in manual credentials mode")
		}
	case ovirt.Name:
		con, err := ovirtconfig.NewConnection()
		if err != nil {
			return errors.Wrap(err, "creating Engine connection")
		}
		err = con.Test()
		if err != nil {
			return errors.Wrap(err, "testing Engine connection")
		}
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}

	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformCredsCheck) Name() string {
	return "Platform Credentials Check"
}
