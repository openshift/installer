package powervs

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in PowerVS.
var PlatformStages = []terraform.Stage{
	stages.NewStage("powervs",
		"cluster",
		[]providers.Provider{providers.IBM, providers.Ignition}),
	stages.NewStage("powervs",
		"post-install",
		[]providers.Provider{providers.IBM}),
}
