package main

import (
	"github.com/hashicorp/terraform/plugin"

	"github.com/coreos/terraform-provider-matchbox/matchbox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: matchbox.Provider,
	})
}
