package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm"
)

func init() {
	ibmProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ibm.Provider,
		})
	}
	KnownPlugins["terraform-provider-ibm"] = ibmProvider
}
