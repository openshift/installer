// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetMetadataResponse returns the *tfprotov5.GetMetadataResponse
// equivalent of a *fwserver.GetMetadataResponse.
func GetMetadataResponse(ctx context.Context, fw *fwserver.GetMetadataResponse) *tfprotov5.GetMetadataResponse {
	if fw == nil {
		return nil
	}

	protov6 := &tfprotov5.GetMetadataResponse{
		DataSources:        []tfprotov5.DataSourceMetadata{},
		Diagnostics:        Diagnostics(ctx, fw.Diagnostics),
		Resources:          []tfprotov5.ResourceMetadata{},
		ServerCapabilities: ServerCapabilities(ctx, fw.ServerCapabilities),
	}

	for _, datasource := range fw.DataSources {
		protov6.DataSources = append(protov6.DataSources, DataSourceMetadata(ctx, datasource))
	}

	for _, resource := range fw.Resources {
		protov6.Resources = append(protov6.Resources, ResourceMetadata(ctx, resource))
	}

	return protov6
}
