package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
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
func (a *PlatformPermsCheck) Generate(ctx context.Context, dependencies asset.Parents) error {
	ic := &InstallConfig{}
	dependencies.Get(ic)

	if ic.Config.CredentialsMode != "" {
		logrus.Debug("CredentialsMode is set. Skipping platform permissions checks before attempting installation.")
		return nil
	}
	logrus.Debug("CredentialsMode is not set. Performing platform permissions checks before attempting installation.")

	var err error
	platform := ic.Config.Platform.Name()
	switch platform {
	case aws.Name:
		permissionGroups := awsconfig.RequiredPermissionGroups(ic.Config)

		ssn, err := ic.AWS.Session(ctx)
		if err != nil {
			return err
		}

		err = awsconfig.ValidateCreds(ssn, permissionGroups, ic.Config.Platform.AWS.Region)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case gcp.Name:
		client, err := gcpconfig.NewClient(ctx, ic.Config.GCP.ServiceEndpoints)
		if err != nil {
			return err
		}

		if err = gcpconfig.ValidateEnabledServices(ctx, client, ic.Config.GCP.ProjectID); err != nil {
			return errors.Wrap(err, "failed to validate services in this project")
		}
	case ibmcloud.Name:
		// TODO: IBM[#90]: platformpermscheck
	case powervs.Name:
		// Nothing needs to be done here
	case azure.Name, baremetal.Name, external.Name, none.Name, openstack.Name, ovirt.Name, vsphere.Name, nutanix.Name:
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
