//go:build !altinfra
// +build !altinfra

package platform

import (
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure"
	awscapi "github.com/openshift/installer/pkg/infrastructure/aws/clusterapi"
	azureinfra "github.com/openshift/installer/pkg/infrastructure/azure"
	baremetalinfra "github.com/openshift/installer/pkg/infrastructure/baremetal"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcpcapi "github.com/openshift/installer/pkg/infrastructure/gcp/clusterapi"
	ibmcloudcapi "github.com/openshift/installer/pkg/infrastructure/ibmcloud/clusterapi"
	nutanixcapi "github.com/openshift/installer/pkg/infrastructure/nutanix/clusterapi"
	openstackcapi "github.com/openshift/installer/pkg/infrastructure/openstack/clusterapi"
	powervscapi "github.com/openshift/installer/pkg/infrastructure/powervs/clusterapi"
	vspherecapi "github.com/openshift/installer/pkg/infrastructure/vsphere/clusterapi"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages/azure"
	"github.com/openshift/installer/pkg/terraform/stages/ibmcloud"
	"github.com/openshift/installer/pkg/terraform/stages/ovirt"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/featuregates"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
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
		return clusterapi.InitializeProvider(&awscapi.Provider{}), nil
	case azuretypes.Name:
		return clusterapi.InitializeProvider(&azureinfra.Provider{}), nil
	case azuretypes.StackTerraformName:
		return terraform.InitializeProvider(azure.StackPlatformStages), nil
	case baremetaltypes.Name:
		return baremetalinfra.InitializeProvider(), nil
	case gcptypes.Name:
		return clusterapi.InitializeProvider(gcpcapi.Provider{}), nil
	case ibmcloudtypes.Name:
		if types.ClusterAPIFeatureGateEnabled(platform, fg) {
			return clusterapi.InitializeProvider(ibmcloudcapi.Provider{}), nil
		}
		return terraform.InitializeProvider(ibmcloud.PlatformStages), nil
	case nutanixtypes.Name:
		return clusterapi.InitializeProvider(nutanixcapi.Provider{}), nil
	case powervstypes.Name:
		return clusterapi.InitializeProvider(&powervscapi.Provider{}), nil
	case openstacktypes.Name:
		return clusterapi.InitializeProvider(openstackcapi.Provider{}), nil
	case ovirttypes.Name:
		return terraform.InitializeProvider(ovirt.PlatformStages), nil
	case vspheretypes.Name:
		return clusterapi.InitializeProvider(vspherecapi.Provider{}), nil
	case nonetypes.Name:
		// terraform is not used when the platform is "none"
		return terraform.InitializeProvider([]terraform.Stage{}), nil
	case externaltypes.Name:
		// terraform is not used when the platform is "external"
		return terraform.InitializeProvider([]terraform.Stage{}), nil
	}
	return nil, fmt.Errorf("unsupported platform %q", platform)
}
