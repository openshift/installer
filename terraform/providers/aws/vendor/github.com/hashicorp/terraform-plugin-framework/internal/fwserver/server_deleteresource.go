// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// DeleteResourceRequest is the framework server request for a delete request
// with the ApplyResourceChange RPC.
type DeleteResourceRequest struct {
	PlannedPrivate *privatestate.Data
	PriorState     *tfsdk.State
	ProviderMeta   *tfsdk.Config
	ResourceSchema fwschema.Schema
	Resource       resource.Resource
}

// DeleteResourceResponse is the framework server response for a delete request
// with the ApplyResourceChange RPC.
type DeleteResourceResponse struct {
	Diagnostics diag.Diagnostics
	NewState    *tfsdk.State
	Private     *privatestate.Data
}

// DeleteResource implements the framework server delete request logic for the
// ApplyResourceChange RPC.
func (s *Server) DeleteResource(ctx context.Context, req *DeleteResourceRequest, resp *DeleteResourceResponse) {
	if req == nil {
		return
	}

	if resourceWithConfigure, ok := req.Resource.(resource.ResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "Resource implements ResourceWithConfigure")

		configureReq := resource.ConfigureRequest{
			ProviderData: s.ResourceConfigureData,
		}
		configureResp := resource.ConfigureResponse{}

		logging.FrameworkDebug(ctx, "Calling provider defined Resource Configure")
		resourceWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkDebug(ctx, "Called provider defined Resource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	deleteReq := resource.DeleteRequest{
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    tftypes.NewValue(req.ResourceSchema.Type().TerraformType(ctx), nil),
		},
	}
	deleteResp := resource.DeleteResponse{
		State: tfsdk.State{
			Schema: req.ResourceSchema,
			Raw:    tftypes.NewValue(req.ResourceSchema.Type().TerraformType(ctx), nil),
		},
	}

	if req.PriorState != nil {
		deleteReq.State = *req.PriorState
		deleteResp.State = *req.PriorState
	}

	if req.ProviderMeta != nil {
		deleteReq.ProviderMeta = *req.ProviderMeta
	}

	if req.PlannedPrivate != nil {
		deleteReq.Private = req.PlannedPrivate.Provider
	}

	logging.FrameworkDebug(ctx, "Calling provider defined Resource Delete")
	req.Resource.Delete(ctx, deleteReq, &deleteResp)
	logging.FrameworkDebug(ctx, "Called provider defined Resource Delete")

	if !deleteResp.Diagnostics.HasError() {
		logging.FrameworkTrace(ctx, "No provider defined Delete errors detected, ensuring State is cleared")
		deleteResp.State.RemoveResource(ctx)
	}

	resp.Diagnostics = deleteResp.Diagnostics
	resp.NewState = &deleteResp.State
}
