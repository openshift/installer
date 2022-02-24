package main

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ibm.Provider,
	})
}
