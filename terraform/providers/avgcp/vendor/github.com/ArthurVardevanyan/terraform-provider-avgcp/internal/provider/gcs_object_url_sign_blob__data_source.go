// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &GcsObjectUrlSignBlobDataSource{}

func NewGcsObjectUrlSignBlobDataSource() datasource.DataSource {
	return &GcsObjectUrlSignBlobDataSource{}
}

// GcsObjectUrlSignBlobDataSource defines the data source implementation.
type GcsObjectUrlSignBlobDataSource struct {
	client *http.Client
}

// GcsObjectUrlSignBlobDataSourceModel describes the data source data model.
type GcsObjectUrlSignBlobDataSourceModel struct {
	GoogleAccessID types.String `tfsdk:"google_access_id"`
	Bucket         types.String `tfsdk:"bucket"`
	Path           types.String `tfsdk:"path"`
	SignedUrl      types.String `tfsdk:"signed_url"`
}

func (d *GcsObjectUrlSignBlobDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcs_object_url_sign_blob"
}

func (d *GcsObjectUrlSignBlobDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"google_access_id": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
			},
			"bucket": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
			},
			"signed_url": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
		},
	}
}

func (d *GcsObjectUrlSignBlobDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *GcsObjectUrlSignBlobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GcsObjectUrlSignBlobDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.SignedUrl = types.StringValue("gcs_object_url_sign_blob-signed_url")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
