package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azureconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	openstackconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
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
	case aws.Name:
		permissionGroups := []awsconfig.PermissionGroup{awsconfig.PermissionCreateBase, awsconfig.PermissionDeleteBase}
		// If subnets are not provided in install-config.yaml, include network permissions
		if len(ic.Config.AWS.Subnets) == 0 {
			permissionGroups = append(permissionGroups, awsconfig.PermissionCreateNetworking, awsconfig.PermissionDeleteNetworking)
		}

		ssn, err := ic.AWS.Session(ctx)
		if err != nil {
			return err
		}
		err = awsconfig.ValidateCreds(ssn, permissionGroups)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case gcp.Name:
		_, err = gcpconfig.GetSession(ctx)
		if err != nil {
			return errors.Wrap(err, "creating GCP session")
		}
	case openstack.Name:
		_, err = openstackconfig.GetSession(ic.Config.Platform.OpenStack.Cloud)
	case baremetal.Name, libvirt.Name, none.Name, vsphere.Name:
		// no creds to check
	case azure.Name:
		_, err = azureconfig.GetSession()
		if err != nil {
			return errors.Wrap(err, "creating Azure session")
		}
	case ovirt.Name:
		ovirtConfig, err := ovirtconfig.NewConfig()
		if err != nil {
			return errors.Wrap(err, "getting ovirt configuration")
		}
		con, err := ovirtconfig.GetConnection(ovirtConfig)
		if err != nil {
			return errors.Wrap(err, "establishing ovirt connection")
		}
		err = con.Test()
		if err != nil {
			return errors.Wrap(err, "testing ovirt connection")
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
