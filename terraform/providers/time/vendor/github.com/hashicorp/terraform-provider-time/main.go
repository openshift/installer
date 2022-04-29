package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/hashicorp/terraform-provider-time/internal/tftime"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tftime.Provider,
	})
}
