package libvirt

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in libvirt.
var PlatformStages = []terraform.Stage{
	stages.NewStage("libvirt", "cluster"),
	stages.NewStage("libvirt", "bootstrap", stages.WithNormalDestroy()),
}
