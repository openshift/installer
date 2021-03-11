package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: aws.Provider,
		})
	}
	KnownPlugins["terraform-provider-aws"] = &TFPlugin{
		Name:    "aws",
		Exec:    exec,
		Version: GetAwsVersion(),
	}
}
