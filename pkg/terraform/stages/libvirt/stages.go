package libvirt

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in libvirt.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"libvirt",
		"cluster",
		[]providers.Provider{providers.Libvirt},
	),
	stages.NewStage(
		"libvirt",
		"bootstrap",
		[]providers.Provider{providers.Libvirt},
		stages.WithNormalBootstrapDestroy(),
	),
}
