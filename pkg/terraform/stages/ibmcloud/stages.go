package ibmcloud

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in IBM Cloud.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"ibmcloud",
		"network",
		[]providers.Provider{providers.IBM},
	),
	stages.NewStage(
		"ibmcloud",
		"bootstrap",
		[]providers.Provider{providers.IBM},
		stages.WithNormalBootstrapDestroy(),
	),
	stages.NewStage(
		"ibmcloud",
		"master",
		[]providers.Provider{providers.IBM},
	),
}
