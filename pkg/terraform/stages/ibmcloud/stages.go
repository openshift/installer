package ibmcloud

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in IBM Cloud.
var PlatformStages = []terraform.Stage{
	stages.NewStage("ibmcloud", "network"),
	stages.NewStage("ibmcloud", "bootstrap", stages.WithNormalDestroy()),
	stages.NewStage("ibmcloud", "master"),
}
