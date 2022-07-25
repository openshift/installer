package main

import (
	"github.com/community-terraform-providers/terraform-provider-ignition/v2/ignition"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ignition.Provider})
}
