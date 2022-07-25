package nutanix

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in Nutanix.
var PlatformStages = []terraform.Stage{
	stages.NewStage("nutanix", "cluster", []providers.Provider{providers.Nutanix}),
	stages.NewStage("nutanix", "bootstrap", []providers.Provider{providers.Nutanix}, stages.WithNormalBootstrapDestroy()),
}
