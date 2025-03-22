// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.1-71478489-20240820-161623
 */

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmCodeEngineAllowedOutboundDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineAllowedOutboundDestinationRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your allowed outbound destination.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the allowed outbound destination, which is used to achieve optimistic locking.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.",
			},
			"cidr_block": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 address range.",
			},
		},
	}
}

func dataSourceIbmCodeEngineAllowedOutboundDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_allowed_outbound_destination", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getAllowedOutboundDestinationOptions := &codeenginev2.GetAllowedOutboundDestinationOptions{}

	getAllowedOutboundDestinationOptions.SetProjectID(d.Get("project_id").(string))
	getAllowedOutboundDestinationOptions.SetName(d.Get("name").(string))

	allowedOutboundDestinationIntf, _, err := codeEngineClient.GetAllowedOutboundDestinationWithContext(context, getAllowedOutboundDestinationOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetAllowedOutboundDestinationWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_allowed_outbound_destination", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	allowedOutboundDestination := allowedOutboundDestinationIntf.(*codeenginev2.AllowedOutboundDestination)

	d.SetId(fmt.Sprintf("%s/%s", *getAllowedOutboundDestinationOptions.ProjectID, *getAllowedOutboundDestinationOptions.Name))

	if err = d.Set("entity_tag", allowedOutboundDestination.EntityTag); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_code_engine_allowed_outbound_destination", "read", "set-entity_tag").GetDiag()
	}

	if err = d.Set("type", allowedOutboundDestination.Type); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_code_engine_allowed_outbound_destination", "read", "set-type").GetDiag()
	}

	if !core.IsNil(allowedOutboundDestination.CidrBlock) {
		if err = d.Set("cidr_block", allowedOutboundDestination.CidrBlock); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cidr_block: %s", err), "(Data) ibm_code_engine_allowed_outbound_destination", "read", "set-cidr_block").GetDiag()
		}
	}

	return nil
}
