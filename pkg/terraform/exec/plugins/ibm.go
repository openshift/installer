package plugins

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func init() {
	ibmProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: ibm.Provider,
		})
	}
	KnownPlugins["terraform-provider-ibm"] = ibmProvider
}
