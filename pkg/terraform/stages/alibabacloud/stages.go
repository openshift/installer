package alibabacloud

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/providers"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in Alibaba Cloud.
var PlatformStages = []terraform.Stage{
	stages.NewStage(
		"alibabacloud",
		"cluster",
		[]providers.Provider{providers.AliCloud},
	),
	stages.NewStage(
		"alibabacloud",
		"bootstrap",
		[]providers.Provider{providers.AliCloud},
		stages.WithNormalBootstrapDestroy(),
	),
}
