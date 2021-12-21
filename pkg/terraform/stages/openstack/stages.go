package openstack

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in
// OpenStack.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"openstack",
		"masters",
		[]providers.Provider{providers.OpenStack, providers.Ignition},
	),
	stages.NewStage(
		"openstack",
		"bootstrap",
		[]providers.Provider{providers.OpenStack},
		stages.WithNormalBootstrapDestroy(),
	),
}
