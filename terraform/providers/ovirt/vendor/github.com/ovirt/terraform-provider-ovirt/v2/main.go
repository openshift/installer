package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/ovirt/terraform-provider-ovirt/v2/internal/ovirt"
)

func main() {
	var debugMode bool

	flag.BoolVar(
		&debugMode,
		"debug",
		false,
		"set to true to run the provider with support for debuggers like delve",
	)
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderAddr: "registry.terraform.io/oVirt/ovirt",
		ProviderFunc: ovirt.New(),
	}

	if debugMode {
		opts.Debug = true
	}

	plugin.Serve(opts)
}
