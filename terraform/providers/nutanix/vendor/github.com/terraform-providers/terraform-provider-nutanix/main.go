package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nutanix.Provider,
	})
}
