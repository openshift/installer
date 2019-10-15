package installconfig

import (
	"context"
	"fmt"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azureconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
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
	ic := &InstallConfig{}
	dependencies.Get(ic)

	var err error
	platform := ic.Config.Platform.Name()
	switch platform {
	case aws.Name:
		ssn, err := awsconfig.GetSession()
		if err != nil {
			return errors.Wrap(err, "creating AWS session")
		}
		err = awsconfig.ValidateCreds(ssn, ic.Config.Platform.AWS.Region)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case gcp.Name:
		_, err = gcpconfig.GetSession(context.TODO())
		if err != nil {
			return errors.Wrap(err, "creating GCP session")
		}
	case openstack.Name:
		opts := new(clientconfig.ClientOpts)
		opts.Cloud = ic.Config.Platform.OpenStack.Cloud
		_, err = clientconfig.GetCloudFromYAML(opts)
	case baremetal.Name, libvirt.Name, none.Name, vsphere.Name:
		// no creds to check
	case azure.Name:
		_, err = azureconfig.GetSession()
		if err != nil {
			return errors.Wrap(err, "creating Azure session")
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
