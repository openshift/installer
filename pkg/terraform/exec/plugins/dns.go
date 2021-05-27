package plugins

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-dns/dns"
)

func init() {
	dnsProvider := func() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: dns.Provider,
		})
	}
	KnownPlugins["terraform-provider-dns"] = dnsProvider
}
