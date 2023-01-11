package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateDataSourceConfigRequest is the framework server request for the
// ValidateDataSourceConfig RPC.
type ValidateDataSourceConfigRequest struct {
	Config     *tfsdk.Config
	DataSource datasource.DataSource
}

// ValidateDataSourceConfigResponse is the framework server response for the
// ValidateDataSourceConfig RPC.
type ValidateDataSourceConfigResponse struct {
	Diagnostics diag.Diagnostics
}

// ValidateDataSourceConfig implements the framework server ValidateDataSourceConfig RPC.
func (s *Server) ValidateDataSourceConfig(ctx context.Context, req *ValidateDataSourceConfigRequest, resp *ValidateDataSourceConfigResponse) {
	if req == nil || req.Config == nil {
		return
	}

	if _, ok := req.DataSource.(datasource.DataSourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "DataSource implements DataSourceWithConfigure")

		configureReq := datasource.ConfigureRequest{
			ProviderData: s.DataSourceConfigureData,
		}
		configureResp := datasource.ConfigureResponse{}

		logging.FrameworkDebug(ctx, "Calling provider defined DataSource Configure")
		req.DataSource.(datasource.DataSourceWithConfigure).Configure(ctx, configureReq, &configureResp)
		logging.FrameworkDebug(ctx, "Called provider defined DataSource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	vdscReq := datasource.ValidateConfigRequest{
		Config: *req.Config,
	}

	if dataSource, ok := req.DataSource.(datasource.DataSourceWithConfigValidators); ok {
		logging.FrameworkTrace(ctx, "DataSource implements DataSourceWithConfigValidators")

		for _, configValidator := range dataSource.ConfigValidators(ctx) {
			vdscResp := &datasource.ValidateConfigResponse{
				Diagnostics: resp.Diagnostics,
			}

			logging.FrameworkDebug(
				ctx,
				"Calling provider defined ConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)
			configValidator.ValidateDataSource(ctx, vdscReq, vdscResp)
			logging.FrameworkDebug(
				ctx,
				"Called provider defined ConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)

			resp.Diagnostics = vdscResp.Diagnostics
		}
	}

	if dataSource, ok := req.DataSource.(datasource.DataSourceWithValidateConfig); ok {
		logging.FrameworkTrace(ctx, "DataSource implements DataSourceWithValidateConfig")

		vdscResp := &datasource.ValidateConfigResponse{
			Diagnostics: resp.Diagnostics,
		}

		logging.FrameworkDebug(ctx, "Calling provider defined DataSource ValidateConfig")
		dataSource.ValidateConfig(ctx, vdscReq, vdscResp)
		logging.FrameworkDebug(ctx, "Called provider defined DataSource ValidateConfig")

		resp.Diagnostics = vdscResp.Diagnostics
	}

	validateSchemaReq := ValidateSchemaRequest{
		Config: *req.Config,
	}
	validateSchemaResp := ValidateSchemaResponse{
		Diagnostics: resp.Diagnostics,
	}

	SchemaValidate(ctx, req.Config.Schema, validateSchemaReq, &validateSchemaResp)

	resp.Diagnostics = validateSchemaResp.Diagnostics
}
