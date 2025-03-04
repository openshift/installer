package installconfig

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	bmconfig "github.com/openshift/installer/pkg/asset/installconfig/baremetal"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	ibmcloudconfig "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	nutanixconfig "github.com/openshift/installer/pkg/asset/installconfig/nutanix"
	osconfig "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	ovirtconfig "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	vsconfig "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	nonevalidation "github.com/openshift/installer/pkg/types/none/validation"
	"github.com/openshift/installer/pkg/types/nutanix"
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
//
//nolint:gocyclo
func (a *PlatformProvisionCheck) Generate(ctx context.Context, dependencies asset.Parents) error {
	ic := &InstallConfig{}
	dependencies.Get(ic)
	platform := ic.Config.Platform.Name()

	// IPI requires MachineAPI capability
	enabledCaps := ic.Config.GetEnabledCapabilities()
	if !enabledCaps.Has(configv1.ClusterVersionCapabilityMachineAPI) {
		return errors.New("IPI requires MachineAPI capability")
	}

	switch platform {
	case aws.Name:
		session, err := ic.AWS.Session(ctx)
		if err != nil {
			return err
		}
		client := awsconfig.NewClient(session)
		return awsconfig.ValidateForProvisioning(client, ic.Config, ic.AWS)
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
		err := bmconfig.ValidateBaremetalPlatformSet(ic.Config)
		if err != nil {
			return err
		}
		err = bmconfig.ValidateProvisioning(ic.Config)
		if err != nil {
			return err
		}
		err = bmconfig.ValidateStaticBootstrapNetworking(ic.Config)
		if err != nil {
			return err
		}
		err = baremetalvalidation.ValidateHosts(ic.Config.BareMetal, field.NewPath("platform"), ic.Config).ToAggregate()
		if err != nil {
			return err
		}

	case gcp.Name:
		err := gcpconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	case ibmcloud.Name:
		client, err := ibmcloudconfig.NewClient(ic.Config.Platform.IBMCloud.ServiceEndpoints)
		if err != nil {
			return err
		}
		metadata := ibmcloudconfig.NewMetadata(ic.Config)
		err = ibmcloudconfig.ValidatePreExistingPublicDNS(client, ic.Config, metadata)
		if err != nil {
			return err
		}
	case openstack.Name:
		err := osconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	case vsphere.Name:
		if err := vsconfig.ValidateForProvisioning(ic.Config); err != nil {
			return err
		}
	case ovirt.Name:
		err := ovirtconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	case powervs.Name:
		client, err := powervsconfig.NewClient()
		if err != nil {
			return err
		}

		err = powervsconfig.ValidatePERAvailability(client, ic.Config)
		if err != nil {
			return err
		}

		err = powervsconfig.ValidatePreExistingDNS(client, ic.Config, ic.PowerVS)
		if err != nil {
			return err
		}

		err = powervsconfig.ValidateCustomVPCSetup(client, ic.Config)
		if err != nil {
			return err
		}

		err = powervsconfig.ValidateResourceGroup(client, ic.Config)
		if err != nil {
			return err
		}

		err = powervsconfig.ValidateSystemTypeForZone(client, ic.Config)
		if err != nil {
			return err
		}

		err = powervsconfig.ValidateServiceInstance(client, ic.Config)
		if err != nil {
			return err
		}

		powervsconfig.ValidateClusterOSImageNotSet(client, ic.Config)

		err = powervsconfig.ValidateTransitGateway(client, ic.Config)
		if err != nil {
			return err
		}
	case none.Name:
		err := nonevalidation.ValidateFencingCredentials(ic.Config.None.FencingCredentials, field.NewPath("Platform", "None")).ToAggregate()
		if err != nil {
			return err
		}
	case external.Name:
		// no special provisioning requirements to check
	case nutanix.Name:
		err := nutanixconfig.ValidateForProvisioning(ic.Config)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown platform type %q", platform)
	}
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *PlatformProvisionCheck) Name() string {
	return "Platform Provisioning Check"
}
