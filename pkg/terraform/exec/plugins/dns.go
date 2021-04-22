package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	dns "github.com/terraform-providers/terraform-provider-dns/dns"
)

func init() {
	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: dns.Provider,
		})
	}
	// TODO(displague) update to equinix-metal when TF 0.13+ sdk can be used
	KnownPlugins["terraform-provider-dns"] = exec
}
