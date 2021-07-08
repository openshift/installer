package platform

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages/aws"
	"github.com/openshift/installer/pkg/terraform/stages/compat"
	"github.com/openshift/installer/pkg/terraform/stages/gcp"
	"github.com/openshift/installer/pkg/terraform/stages/ovirt"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
)

// StagesForPlatform returns the terraform stages to run to provision the infrastructure for the specified platform.
func StagesForPlatform(platform string) []terraform.Stage {
	switch platform {
	case awstypes.Name:
		return aws.PlatformStages
	case gcptypes.Name:
		return gcp.PlatformStages
	case ovirttypes.Name:
		return ovirt.PlatformStages
	default:
		return compat.PlatformStages(platform)
	}
}
