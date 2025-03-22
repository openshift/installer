// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package cdtoolchain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtoolchainv2"
)

func DataSourceIBMCdToolchain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdToolchainRead,

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the toolchain.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toolchain name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Describes the toolchain.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account ID where toolchain can be found.",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toolchain region.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group where the toolchain is located.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toolchain CRN.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI that can be used to retrieve toolchain.",
			},
			"ui_href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of a user-facing user interface for this toolchain.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Toolchain creation timestamp.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest toolchain update timestamp.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity that created the toolchain.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Toolchain tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMCdToolchainRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_toolchain", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getToolchainByIDOptions := &cdtoolchainv2.GetToolchainByIDOptions{}

	getToolchainByIDOptions.SetToolchainID(d.Get("toolchain_id").(string))

	toolchain, _, err := cdToolchainClient.GetToolchainByIDWithContext(context, getToolchainByIDOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetToolchainByIDWithContext failed: %s", err.Error()), "(Data) ibm_cd_toolchain", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getToolchainByIDOptions.ToolchainID)

	tags, err := flex.GetTagsUsingCRN(meta, *toolchain.CRN)
	if err != nil {
		log.Printf(
			"Error on get of toolchain (%s) tags: %s", d.Id(), err)
	}
	d.Set("tags", tags)

	if err = d.Set("name", toolchain.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_cd_toolchain", "read", "set-name").GetDiag()
	}

	if err = d.Set("description", toolchain.Description); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_cd_toolchain", "read", "set-description").GetDiag()
	}

	if err = d.Set("account_id", toolchain.AccountID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting account_id: %s", err), "(Data) ibm_cd_toolchain", "read", "set-account_id").GetDiag()
	}

	if err = d.Set("location", toolchain.Location); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting location: %s", err), "(Data) ibm_cd_toolchain", "read", "set-location").GetDiag()
	}

	if err = d.Set("resource_group_id", toolchain.ResourceGroupID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_group_id: %s", err), "(Data) ibm_cd_toolchain", "read", "set-resource_group_id").GetDiag()
	}

	if err = d.Set("crn", toolchain.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cd_toolchain", "read", "set-crn").GetDiag()
	}

	if err = d.Set("href", toolchain.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_cd_toolchain", "read", "set-href").GetDiag()
	}

	if err = d.Set("ui_href", toolchain.UIHref); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ui_href: %s", err), "(Data) ibm_cd_toolchain", "read", "set-ui_href").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(toolchain.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_cd_toolchain", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("updated_at", flex.DateTimeToString(toolchain.UpdatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_cd_toolchain", "read", "set-updated_at").GetDiag()
	}

	if err = d.Set("created_by", toolchain.CreatedBy); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_by: %s", err), "(Data) ibm_cd_toolchain", "read", "set-created_by").GetDiag()
	}

	return nil
}
