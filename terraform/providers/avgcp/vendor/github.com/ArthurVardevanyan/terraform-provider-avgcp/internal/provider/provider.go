// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure avgcpProvider satisfies various provider interfaces.
var _ provider.Provider = &avgcpProvider{}

// avgcpProvider defines the provider implementation.
type avgcpProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// avgcpProviderModel describes the provider data model.
type avgcpProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *avgcpProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "avgcp"
	resp.Version = p.version
}

func (p *avgcpProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *avgcpProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data avgcpProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *avgcpProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGcsObjectUrlSignBlobResource,
	}
}

func (p *avgcpProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewGcsObjectUrlSignBlobDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &avgcpProvider{
			version: version,
		}
	}
}
