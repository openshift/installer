package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-google/google"

	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func init() {
	googleProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: google.Provider,
		})
	}
	KnownPlugins["terraform-provider-google"] = &TFPlugin{
		Name: gcptypes.Name,
		Exec: googleProvider,
	}
}
