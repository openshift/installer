package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/ovirt/terraform-provider-ovirt/ovirt"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ovirt.ProviderContext(),
	})
}
