package openstack

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in
// OpenStack.
var PlatformStages = []terraform.Stage{
	stages.NewStage("openstack", "masters"),
	stages.NewStage("openstack", "bootstrap", stages.WithNormalDestroy()),
}
