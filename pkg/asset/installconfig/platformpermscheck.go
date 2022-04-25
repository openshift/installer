package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"

	gcpclient "github.com/openshift/installer/pkg/client/gcp"
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
		permissionGroups := []awsconfig.PermissionGroup{awsconfig.PermissionCreateBase}
		usingExistingVPC := len(ic.Config.AWS.Subnets) != 0

		if !usingExistingVPC {
			permissionGroups = append(permissionGroups, awsconfig.PermissionCreateNetworking)
		}

		// Add delete permissions for non-C2S installs.
		if !aws.IsSecretRegion(ic.Config.AWS.Region) {
			permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteBase)
			if usingExistingVPC {
				permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteSharedNetworking)
			} else {
				permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteNetworking)
			}
		}

		ssn, err := ic.AWS.Session(ctx)
		if err != nil {
			return err
		}

		err = awsconfig.ValidateCreds(ssn, permissionGroups, ic.Config.Platform.AWS.Region)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case gcp.Name:
		client, err := gcpclient.NewClient(context.TODO())
		if err != nil {
			return err
		}

		if err = gcpconfig.ValidateEnabledServices(ctx, client, ic.Config.GCP.ProjectID); err != nil {
			return errors.Wrap(err, "failed to validate services in this project")
		}
	case ibmcloud.Name:
		// TODO: IBM[#90]: platformpermscheck
	case powervs.Name:
		//@TODO add check that the account plan is anything but "lite"
	case azure.Name, baremetal.Name, libvirt.Name, none.Name, openstack.Name, ovirt.Name, vsphere.Name, alibabacloud.Name, nutanix.Name:
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
