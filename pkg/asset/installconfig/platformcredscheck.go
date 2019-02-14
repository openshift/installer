package installconfig

import (
	"fmt"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
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
		err = awsconfig.ValidateCreds(ssn)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case libvirt.Name:
	case none.Name:
	case openstack.Name:
		opts := new(clientconfig.ClientOpts)
		opts.Cloud = ic.Config.Platform.OpenStack.Cloud
		_, err = clientconfig.GetCloudFromYAML(opts)
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}

	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformCredsCheck) Name() string {
	return "Platform Credentials Check"
}
