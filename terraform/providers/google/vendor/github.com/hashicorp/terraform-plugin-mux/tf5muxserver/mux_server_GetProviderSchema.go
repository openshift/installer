package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// GetProviderSchema merges the schemas returned by the
// tfprotov5.ProviderServers associated with muxServer into a single schema.
// Resources and data sources must be returned from only one server. Provider
// and ProviderMeta schemas must be identical between all servers.
func (s muxServer) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	rpc := "GetProviderSchema"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	logging.MuxTrace(ctx, "serving cached schema information")

	resp := &tfprotov5.GetProviderSchemaResponse{
		Provider:          s.providerSchema,
		ResourceSchemas:   s.resourceSchemas,
		DataSourceSchemas: s.dataSourceSchemas,
		ProviderMeta:      s.providerMetaSchema,
	}

	for _, diff := range s.serverProviderSchemaDifferences {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Summary: "Invalid Provider Server Combination",
			Detail: "The combined provider has differing provider schema implementations across providers. " +
				"Provider schemas must be identical across providers. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Provider schema difference: " + diff,
		})
	}

	for _, diff := range s.serverProviderMetaSchemaDifferences {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Summary: "Invalid Provider Server Combination",
			Detail: "The combined provider has differing provider meta schema implementations across providers. " +
				"Provider meta schemas must be identical across providers. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Provider meta schema difference: " + diff,
		})
	}

	for _, dataSourceType := range s.serverDataSourceSchemaDuplicates {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Summary: "Invalid Provider Server Combination",
			Detail: "The combined provider has multiple implementations of the same data source type across providers. " +
				"Data source types must be implemented by only one provider. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Duplicate data source type: " + dataSourceType,
		})
	}

	for _, resourceType := range s.serverResourceSchemaDuplicates {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Summary: "Invalid Provider Server Combination",
			Detail: "The combined provider has multiple implementations of the same resource type across providers. " +
				"Resource types must be implemented by only one provider. " +
				"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
				"Duplicate resource type: " + resourceType,
		})
	}

	return resp, nil
}
