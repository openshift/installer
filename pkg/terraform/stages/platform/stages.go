package platform

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages/alibabacloud"
	"github.com/openshift/installer/pkg/terraform/stages/aws"
	"github.com/openshift/installer/pkg/terraform/stages/azure"
	"github.com/openshift/installer/pkg/terraform/stages/baremetal"
	"github.com/openshift/installer/pkg/terraform/stages/compat"
	"github.com/openshift/installer/pkg/terraform/stages/gcp"
	"github.com/openshift/installer/pkg/terraform/stages/ibmcloud"
	"github.com/openshift/installer/pkg/terraform/stages/libvirt"
	"github.com/openshift/installer/pkg/terraform/stages/vsphere"
	alibabacloudtypes "github.com/openshift/installer/pkg/types/alibabacloud"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	libvirttypes "github.com/openshift/installer/pkg/types/libvirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// StagesForPlatform returns the terraform stages to run to provision the infrastructure for the specified platform.
func StagesForPlatform(platform string) []terraform.Stage {
	switch platform {
	case alibabacloudtypes.Name:
		return alibabacloud.PlatformStages
	case awstypes.Name:
		return aws.PlatformStages
	case azuretypes.Name:
		return azure.PlatformStages
	case azuretypes.StackTerraformName:
		return azure.StackPlatformStages
	case baremetaltypes.Name:
		return baremetal.PlatformStages
	case gcptypes.Name:
		return gcp.PlatformStages
	case ibmcloudtypes.Name:
		return ibmcloud.PlatformStages
	case libvirttypes.Name:
		return libvirt.PlatformStages
	case vspheretypes.Name:
		return vsphere.PlatformStages
	default:
		return compat.PlatformStages(platform)
	}
}
