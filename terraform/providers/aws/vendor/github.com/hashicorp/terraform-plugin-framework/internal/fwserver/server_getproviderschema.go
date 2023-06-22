// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// GetProviderSchemaRequest is the framework server request for the
// GetProviderSchema RPC.
type GetProviderSchemaRequest struct{}

// GetProviderSchemaResponse is the framework server response for the
// GetProviderSchema RPC.
type GetProviderSchemaResponse struct {
	ServerCapabilities *ServerCapabilities
	Provider           fwschema.Schema
	ProviderMeta       fwschema.Schema
	ResourceSchemas    map[string]fwschema.Schema
	DataSourceSchemas  map[string]fwschema.Schema
	Diagnostics        diag.Diagnostics
}

// GetProviderSchema implements the framework server GetProviderSchema RPC.
func (s *Server) GetProviderSchema(ctx context.Context, req *GetProviderSchemaRequest, resp *GetProviderSchemaResponse) {
	resp.ServerCapabilities = &ServerCapabilities{
		PlanDestroy: true,
	}

	metadataReq := provider.MetadataRequest{}
	metadataResp := provider.MetadataResponse{}

	logging.FrameworkDebug(ctx, "Calling provider defined Provider Metadata")
	s.Provider.Metadata(ctx, metadataReq, &metadataResp)
	logging.FrameworkDebug(ctx, "Called provider defined Provider Metadata")

	s.providerTypeName = metadataResp.TypeName

	providerSchema, diags := s.ProviderSchema(ctx)

	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	resp.Provider = providerSchema

	providerMetaSchema, diags := s.ProviderMetaSchema(ctx)

	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	resp.ProviderMeta = providerMetaSchema

	resourceSchemas, diags := s.ResourceSchemas(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceSchemas = resourceSchemas

	dataSourceSchemas, diags := s.DataSourceSchemas(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceSchemas = dataSourceSchemas
}
