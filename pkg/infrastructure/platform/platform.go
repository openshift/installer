package platform

import (
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages/alibabacloud"
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
	alibabacloudtypes "github.com/openshift/installer/pkg/types/alibabacloud"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
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
func ProviderForPlatform(platform string) infrastructure.Provider {
	switch platform {
	case alibabacloudtypes.Name:
		return terraform.InitializeProvider(alibabacloud.PlatformStages)
	case awstypes.Name:
		return terraform.InitializeProvider(aws.PlatformStages)
	case azuretypes.Name:
		return terraform.InitializeProvider(azure.PlatformStages)
	case azuretypes.StackTerraformName:
		return terraform.InitializeProvider(azure.StackPlatformStages)
	case baremetaltypes.Name:
		return terraform.InitializeProvider(baremetal.PlatformStages)
	case gcptypes.Name:
		return terraform.InitializeProvider(gcp.PlatformStages)
	case ibmcloudtypes.Name:
		return terraform.InitializeProvider(ibmcloud.PlatformStages)
	case libvirttypes.Name:
		return terraform.InitializeProvider(libvirt.PlatformStages)
	case nutanixtypes.Name:
		return terraform.InitializeProvider(nutanix.PlatformStages)
	case powervstypes.Name:
		return terraform.InitializeProvider(powervs.PlatformStages)
	case openstacktypes.Name:
		return terraform.InitializeProvider(openstack.PlatformStages)
	case ovirttypes.Name:
		return terraform.InitializeProvider(ovirt.PlatformStages)
	case vspheretypes.Name:
		return terraform.InitializeProvider(vsphere.PlatformStages)
	case nonetypes.Name:
		// terraform is not used when the platform is "none"
		return terraform.InitializeProvider([]terraform.Stage{})
	case externaltypes.Name:
		// terraform is not used when the platform is "external"
		return terraform.InitializeProvider([]terraform.Stage{})
	}
	panic(fmt.Sprintf("unsupported platform %q", platform))
}
