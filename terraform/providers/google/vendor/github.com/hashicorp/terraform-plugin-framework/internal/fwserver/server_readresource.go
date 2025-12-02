package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ReadResourceRequest is the framework server request for the
// ReadResource RPC.
type ReadResourceRequest struct {
	CurrentState *tfsdk.State
	Resource     resource.Resource
	Private      *privatestate.Data
	ProviderMeta *tfsdk.Config
}

// ReadResourceResponse is the framework server response for the
// ReadResource RPC.
type ReadResourceResponse struct {
	Diagnostics diag.Diagnostics
	NewState    *tfsdk.State
	Private     *privatestate.Data
}

// ReadResource implements the framework server ReadResource RPC.
func (s *Server) ReadResource(ctx context.Context, req *ReadResourceRequest, resp *ReadResourceResponse) {
	if req == nil {
		return
	}

	if req.CurrentState == nil {
		resp.Diagnostics.AddError(
			"Unexpected Read Request",
			"An unexpected error was encountered when reading the resource. The current state was missing.\n\n"+
				"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.",
		)

		return
	}

	if _, ok := req.Resource.(resource.ResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "Resource implements ResourceWithConfigure")

		configureReq := resource.ConfigureRequest{
			ProviderData: s.ResourceConfigureData,
		}
		configureResp := resource.ConfigureResponse{}

		logging.FrameworkDebug(ctx, "Calling provider defined Resource Configure")
		req.Resource.(resource.ResourceWithConfigure).Configure(ctx, configureReq, &configureResp)
		logging.FrameworkDebug(ctx, "Called provider defined Resource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	readReq := resource.ReadRequest{
		State: tfsdk.State{
			Schema: req.CurrentState.Schema,
			Raw:    req.CurrentState.Raw.Copy(),
		},
	}
	readResp := resource.ReadResponse{
		State: tfsdk.State{
			Schema: req.CurrentState.Schema,
			Raw:    req.CurrentState.Raw.Copy(),
		},
	}

	if req.ProviderMeta != nil {
		readReq.ProviderMeta = *req.ProviderMeta
	}

	privateProviderData := privatestate.EmptyProviderData(ctx)

	readReq.Private = privateProviderData
	readResp.Private = privateProviderData

	if req.Private != nil {
		if req.Private.Provider != nil {
			readReq.Private = req.Private.Provider
			readResp.Private = req.Private.Provider
		}

		resp.Private = req.Private
	}

	logging.FrameworkDebug(ctx, "Calling provider defined Resource Read")
	req.Resource.Read(ctx, readReq, &readResp)
	logging.FrameworkDebug(ctx, "Called provider defined Resource Read")

	resp.Diagnostics = readResp.Diagnostics
	resp.NewState = &readResp.State

	if readResp.Private != nil {
		if resp.Private == nil {
			resp.Private = &privatestate.Data{}
		}

		resp.Private.Provider = readResp.Private
	}
}
