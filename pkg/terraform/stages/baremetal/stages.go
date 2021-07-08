package baremetal

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in
// Bare Metal.
var PlatformStages = []terraform.Stage{
	stages.NewStage("baremetal", "bootstrap", stages.WithNormalDestroy()),
	stages.NewStage("baremetal", "masters"),
}
