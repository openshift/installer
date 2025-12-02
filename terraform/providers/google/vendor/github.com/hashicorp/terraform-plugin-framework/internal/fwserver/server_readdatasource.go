package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ReadDataSourceRequest is the framework server request for the
// ReadDataSource RPC.
type ReadDataSourceRequest struct {
	Config           *tfsdk.Config
	DataSourceSchema fwschema.Schema
	DataSource       datasource.DataSource
	ProviderMeta     *tfsdk.Config
}

// ReadDataSourceResponse is the framework server response for the
// ReadDataSource RPC.
type ReadDataSourceResponse struct {
	Diagnostics diag.Diagnostics
	State       *tfsdk.State
}

// ReadDataSource implements the framework server ReadDataSource RPC.
func (s *Server) ReadDataSource(ctx context.Context, req *ReadDataSourceRequest, resp *ReadDataSourceResponse) {
	if req == nil {
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

	readReq := datasource.ReadRequest{
		Config: tfsdk.Config{
			Schema: req.DataSourceSchema,
		},
	}
	readResp := datasource.ReadResponse{
		State: tfsdk.State{
			Schema: req.DataSourceSchema,
		},
	}

	if req.Config != nil {
		readReq.Config = *req.Config
		readResp.State.Raw = req.Config.Raw.Copy()
	}

	if req.ProviderMeta != nil {
		readReq.ProviderMeta = *req.ProviderMeta
	}

	logging.FrameworkDebug(ctx, "Calling provider defined DataSource Read")
	req.DataSource.Read(ctx, readReq, &readResp)
	logging.FrameworkDebug(ctx, "Called provider defined DataSource Read")

	resp.Diagnostics = readResp.Diagnostics
	resp.State = &readResp.State
}
