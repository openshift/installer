package kubevirt

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in
// Kubevirt.
var PlatformStages = []terraform.Stage{
	stages.NewStage("kubevirt", "masters"),
	stages.NewStage("kubevirt", "bootstrap", stages.WithNormalDestroy()),
}
