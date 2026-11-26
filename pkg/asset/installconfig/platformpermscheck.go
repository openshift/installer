package installconfig

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ibmcloudconfig "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
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
	"github.com/openshift/installer/pkg/types/powervc"
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

	platform := ic.Config.Platform.Name()
	// IBM Cloud IPI requires credentialsMode: Manual. That only means CCO will
	// not mint in-cluster credentials; the installer API key still provisions
	// infrastructure and must be checked. Other platforms skip when the mode is set.
	if skipPermsCheckForCredentialsMode(platform, string(ic.Config.CredentialsMode)) {
		logrus.Debug("CredentialsMode is set. Skipping platform permissions checks before attempting installation.")
		return nil
	}
	logrus.Debug("Performing platform permissions checks before attempting installation.")

	var err error
	switch platform {
	case aws.Name:
		permissionGroups := awsconfig.RequiredPermissionGroups(ic.Config)
		region := ic.Config.Platform.AWS.Region
		endpoints := ic.Config.AWS.ServiceEndpoints

		awsConfig, err := awsconfig.GetConfigWithOptions(ctx, config.WithRegion(region))
		if err != nil {
			return err
		}

		endpointResolver := awsconfig.NewServiceEndpointResolver(awsconfig.EndpointOptions{
			Region:    region,
			Endpoints: endpoints,
		})

		var iamEndpoint string
		if ep, found := endpointResolver.GetCustomEndpoint(iam.ServiceID); found {
			iamEndpoint = ep.URL
		}

		err = awsconfig.ValidateCreds(ctx, awsConfig, permissionGroups, region, iamEndpoint)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case gcp.Name:
		client, err := gcpconfig.NewClient(ctx, ic.Config.GCP.Endpoint)
		if err != nil {
			return err
		}

		if err = gcpconfig.ValidateEnabledServices(ctx, client, ic.Config.GCP.ProjectID); err != nil {
			return errors.Wrap(err, "failed to validate services in this project")
		}
	case ibmcloud.Name:
		if ic.Config.Platform.IBMCloud == nil {
			return fmt.Errorf("ibmcloud platform configuration is required")
		}
		client, clientErr := ibmcloudconfig.NewClient(ic.Config.Platform.IBMCloud.ServiceEndpoints)
		if clientErr != nil {
			return errors.Wrap(clientErr, "creating IBM Cloud session")
		}
		if err = ibmcloudconfig.ValidatePerms(ctx, client, ic.Config); err != nil {
			return errors.Wrap(err, "validate IBM Cloud permissions")
		}
	case powervs.Name:
		// Nothing needs to be done here
	case azure.Name, baremetal.Name, external.Name, none.Name, openstack.Name, powervc.Name, ovirt.Name, vsphere.Name, nutanix.Name:
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

// skipPermsCheckForCredentialsMode reports whether PlatformPermsCheck should
// skip because credentialsMode is set. IBM Cloud is excluded: Manual is
// required there, but the installer API key still needs IAM probes.
func skipPermsCheckForCredentialsMode(platform, credentialsMode string) bool {
	if credentialsMode == "" {
		return false
	}
	return platform != ibmcloud.Name
}
