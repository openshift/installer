package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types"
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
		permissionGroups := []awsconfig.PermissionGroup{awsconfig.PermissionCreateBase}
		usingExistingVPC := len(ic.Config.AWS.Subnets) != 0
		usingExistingPrivateZone := len(ic.Config.AWS.HostedZone) != 0

		if !usingExistingVPC {
			permissionGroups = append(permissionGroups, awsconfig.PermissionCreateNetworking)
		}

		if !usingExistingPrivateZone {
			permissionGroups = append(permissionGroups, awsconfig.PermissionCreateHostedZone)
		}

		var ec2RootVolume = aws.EC2RootVolume{}
		var awsMachinePoolUsingKMS, masterMachinePoolUsingKMS bool
		if ic.Config.AWS.DefaultMachinePlatform != nil && ic.Config.AWS.DefaultMachinePlatform.EC2RootVolume != ec2RootVolume {
			awsMachinePoolUsingKMS = len(ic.Config.AWS.DefaultMachinePlatform.EC2RootVolume.KMSKeyARN) != 0
		}
		if ic.Config.ControlPlane != nil &&
			ic.Config.ControlPlane.Name == types.MachinePoolControlPlaneRoleName &&
			ic.Config.ControlPlane.Platform.AWS != nil &&
			ic.Config.ControlPlane.Platform.AWS.EC2RootVolume != ec2RootVolume {
			masterMachinePoolUsingKMS = len(ic.Config.ControlPlane.Platform.AWS.EC2RootVolume.KMSKeyARN) != 0
		}
		// Add KMS encryption keys, if provided.
		if awsMachinePoolUsingKMS || masterMachinePoolUsingKMS {
			logrus.Debugf("Adding %s to the group of permissions to validate", awsconfig.PermissionKMSEncryptionKeys)
			permissionGroups = append(permissionGroups, awsconfig.PermissionKMSEncryptionKeys)
		}

		// Add delete permissions for non-C2S installs.
		if !aws.IsSecretRegion(ic.Config.AWS.Region) {
			permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteBase)
			if usingExistingVPC {
				permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteSharedNetworking)
			} else {
				permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteNetworking)
			}
			if !usingExistingPrivateZone {
				permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteHostedZone)
			}
		}

		if ic.Config.AWS.PublicIpv4Pool != "" {
			permissionGroups = append(permissionGroups, awsconfig.PermissionPublicIpv4Pool)
		}

		if !ic.Config.AWS.BestEffortDeleteIgnition {
			permissionGroups = append(permissionGroups, awsconfig.PermissionDeleteIgnitionObjects)
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
		client, err := gcpconfig.NewClient(ctx)
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
