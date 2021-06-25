package ovirt

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
)

// PlatformStages are the stages to run to provision the infrastructure in oVirt.
var PlatformStages = []terraform.Stage{
	stages.NewStage(ovirttypes.Name, "image", stages.WithNormalDestroy()),
	stages.NewStage(ovirttypes.Name, "cluster"),
	stages.NewStage(ovirttypes.Name, "bootstrap", stages.WithNormalDestroy()),
}

// Note that the bootstrap IP and control-plane IP extraction is not supported on oVirt
