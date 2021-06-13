package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-aws/aws"

	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func init() {
	awsProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: aws.Provider,
		})
	}
	KnownPlugins["terraform-provider-aws"] = &TFPlugin{
		Name: awstypes.Name,
		Exec: awsProvider,
	}
}
