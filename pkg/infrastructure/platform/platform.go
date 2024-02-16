//go:build !altinfra
// +build !altinfra

package platform

import (
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/infrastructure"
	awsinfra "github.com/openshift/installer/pkg/infrastructure/aws"
	awscapi "github.com/openshift/installer/pkg/infrastructure/aws/clusterapi"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcpcapi "github.com/openshift/installer/pkg/infrastructure/gcp/clusterapi"
	openstackcapi "github.com/openshift/installer/pkg/infrastructure/openstack/clusterapi"
	vspherecapi "github.com/openshift/installer/pkg/infrastructure/vsphere/clusterapi"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages/aws"
	"github.com/openshift/installer/pkg/terraform/stages/azure"
	"github.com/openshift/installer/pkg/terraform/stages/baremetal"
	"github.com/openshift/installer/pkg/terraform/stages/gcp"
	"github.com/openshift/installer/pkg/terraform/stages/ibmcloud"
	"github.com/openshift/installer/pkg/terraform/stages/libvirt"
	"github.com/openshift/installer/pkg/terraform/stages/nutanix"
	"github.com/openshift/installer/pkg/terraform/stages/openstack"
	"github.com/openshift/installer/pkg/terraform/stages/ovirt"
	"github.com/openshift/installer/pkg/terraform/stages/powervs"
	"github.com/openshift/installer/pkg/terraform/stages/vsphere"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/featuregates"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// ProviderForPlatform returns the stages to run to provision the infrastructure for the specified platform.
func ProviderForPlatform(platform string, fg featuregates.FeatureGate) (infrastructure.Provider, error) {
	switch platform {
	case awstypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(&awscapi.Provider{}), nil
		}
		if fg.Enabled(configv1.FeatureGateInstallAlternateInfrastructureAWS) {
			return awsinfra.InitializeProvider(), nil
		}
		return terraform.InitializeProvider(aws.PlatformStages), nil
	case azuretypes.Name:
		return terraform.InitializeProvider(azure.PlatformStages), nil
	case azuretypes.StackTerraformName:
		return terraform.InitializeProvider(azure.StackPlatformStages), nil
	case baremetaltypes.Name:
		return terraform.InitializeProvider(baremetal.PlatformStages), nil
	case gcptypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(gcpcapi.Provider{}), nil
		}
		return terraform.InitializeProvider(gcp.PlatformStages), nil
	case ibmcloudtypes.Name:
		return terraform.InitializeProvider(ibmcloud.PlatformStages), nil
	case libvirttypes.Name:
		return terraform.InitializeProvider(libvirt.PlatformStages), nil
	case nutanixtypes.Name:
		return terraform.InitializeProvider(nutanix.PlatformStages), nil
	case powervstypes.Name:
		return terraform.InitializeProvider(powervs.PlatformStages), nil
	case openstacktypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(openstackcapi.Provider{}), nil
		}
		return terraform.InitializeProvider(openstack.PlatformStages), nil
	case ovirttypes.Name:
		return terraform.InitializeProvider(ovirt.PlatformStages), nil
	case vspheretypes.Name:
		if fg.Enabled(configv1.FeatureGateClusterAPIInstall) {
			return clusterapi.InitializeProvider(vspherecapi.Provider{}), nil
		}
		return terraform.InitializeProvider(vsphere.PlatformStages), nil
	case nonetypes.Name:
		// terraform is not used when the platform is "none"
		return terraform.InitializeProvider([]terraform.Stage{}), nil
	case externaltypes.Name:
		// terraform is not used when the platform is "external"
		return terraform.InitializeProvider([]terraform.Stage{}), nil
	}
	return nil, fmt.Errorf("unsupported platform %q", platform)
}
