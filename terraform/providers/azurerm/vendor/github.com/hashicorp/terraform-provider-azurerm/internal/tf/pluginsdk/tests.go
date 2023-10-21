package pluginsdk

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type TestCheckFunc = resource.TestCheckFunc

type InstanceState = terraform.InstanceState

type State = terraform.State
