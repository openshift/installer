package alibabacloud

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/terraform/stages"
)

// PlatformStages are the stages to run to provision the infrastructure in Alibaba Cloud.
var PlatformStages = []terraform.Stage{
	stages.NewStage("alibabacloud", "cluster"),
	stages.NewStage("alibabacloud", "bootstrap", stages.WithNormalDestroy()),
}
