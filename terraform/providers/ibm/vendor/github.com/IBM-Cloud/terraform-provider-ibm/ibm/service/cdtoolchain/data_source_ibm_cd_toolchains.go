// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func DataSourceIBMCdToolchains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdToolchainsRead,

		Schema: map[string]*schema.Schema{
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource group ID where the toolchains exist.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of toolchain to look up.",
			},
			"toolchains": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Toolchain results returned from the collection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Toolchain ID.",
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
					},
				},
			},
		},
	}
}

func dataSourceIBMCdToolchainsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listToolchainsOptions := &cdtoolchainv2.ListToolchainsOptions{}

	listToolchainsOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	if _, ok := d.GetOk("name"); ok {
		listToolchainsOptions.SetName(d.Get("name").(string))
	}

	var pager *cdtoolchainv2.ToolchainsPager
	pager, err = cdToolchainClient.NewToolchainsPager(listToolchainsOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] ToolchainsPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("ToolchainsPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIBMCdToolchainsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMCdToolchainsToolchainModelToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("toolchains", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchains %s", err))
	}

	return nil
}

// dataSourceIBMCdToolchainsID returns a reasonable ID for the list.
func dataSourceIBMCdToolchainsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMCdToolchainsToolchainModelToMap(model *cdtoolchainv2.ToolchainModel) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["description"] = model.Description
	modelMap["account_id"] = model.AccountID
	modelMap["location"] = model.Location
	modelMap["resource_group_id"] = model.ResourceGroupID
	modelMap["crn"] = model.CRN
	modelMap["href"] = model.Href
	modelMap["ui_href"] = model.UIHref
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["updated_at"] = model.UpdatedAt.String()
	modelMap["created_by"] = model.CreatedBy
	return modelMap, nil
}
