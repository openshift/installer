package privatednsresolver

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/private-dns-resolver"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Private DNS Resolver"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Private DNS Resolver",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		PrivateDNSResolverDnsForwardingRulesetDataSource{},
		PrivateDNSResolverDnsResolverDataSource{},
		PrivateDNSResolverForwardingRuleDataSource{},
		PrivateDNSResolverInboundEndpointDataSource{},
		PrivateDNSResolverOutboundEndpointDataSource{},
		PrivateDNSResolverVirtualNetworkLinkDataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		PrivateDNSResolverDnsForwardingRulesetResource{},
		PrivateDNSResolverDnsResolverResource{},
		PrivateDNSResolverForwardingRuleResource{},
		PrivateDNSResolverInboundEndpointResource{},
		PrivateDNSResolverOutboundEndpointResource{},
		PrivateDNSResolverVirtualNetworkLinkResource{},
	}
}
