package vsphere

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in vsphere.
var PlatformStages = []terraform.Stage{
	stages.NewStage("vsphere", "pre-bootstrap"),
	stages.NewStage("vsphere", "bootstrap", stages.WithNormalDestroy()),
	stages.NewStage("vsphere", "master"),
}
