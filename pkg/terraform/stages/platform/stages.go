package platform

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages/compat"
)

// StagesForPlatform returns the terraform stages to run to provision the infrastructure for the specified platform.
func StagesForPlatform(platform string) []terraform.Stage {
	return compat.PlatformStages(platform)
}
