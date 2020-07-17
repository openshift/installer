package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
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

// PlatformPermsCheck is an asset that checks platform credentials for the necessary permissions
// to create a cluster.
type PlatformPermsCheck struct {
}

var _ asset.Asset = (*PlatformPermsCheck)(nil)

// Dependencies returns the dependencies for PlatformPermsCheck
func (a *PlatformPermsCheck) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate queries for input from the user.
func (a *PlatformPermsCheck) Generate(dependencies asset.Parents) error {
	ctx := context.TODO()
	ic := &InstallConfig{}
	dependencies.Get(ic)

	if ic.Config.CredentialsMode != "" {
		return nil
	}

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

		err = awsconfig.ValidateCreds(ssn, permissionGroups, ic.Config.Platform.AWS.Region)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case azure.Name, baremetal.Name, gcp.Name, libvirt.Name, none.Name, openstack.Name, ovirt.Name, vsphere.Name:
		// no permissions to check
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}
	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformPermsCheck) Name() string {
	return "Platform Permissions Check"
}
