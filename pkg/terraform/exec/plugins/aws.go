package plugins

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/vrutkovs/terraform-provider-aws/v3/aws"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: aws.Provider,
		})
	}
	KnownPlugins["terraform-provider-aws"] = exec
}
