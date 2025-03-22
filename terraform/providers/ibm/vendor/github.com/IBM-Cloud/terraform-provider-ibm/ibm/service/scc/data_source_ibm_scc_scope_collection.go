// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccScopeCollection() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccScopeCollectionRead,

		Schema: map[string]*schema.Schema{
			"scopes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The array of scopes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the scopes.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Security and Compliance instance owning the scope.",
						},
						"account_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the account associated with the scope.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the attachment was created.",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who created the attachment.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the attachment was updated.",
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who updated the attachment.",
						},
						"environment": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The environment that relates to this scope.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the scope",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the scope",
						},
						"properties": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "The properties supported for scoping by this environment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the property.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The value of the property.",
									},
								},
							},
						},
						"attachment_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of attachments linked to that scope.",
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccScopeCollectionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	scopesOptions := &securityandcompliancecenterapiv3.ListScopesOptions{}
	scopesOptions.SetInstanceID(d.Get("instance_id").(string))

	pager, err := securityAndComplianceCenterApIsClient.NewScopesPager(scopesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListScopesWithContext failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ListScopesWithContext failed %s", err))
	}

	scopesCollection, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] ListScopesWithContext failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ListScopesWithContext failed %s", err))
	}

	d.SetId(dataSourceIbmSccScopeCollectionID(d))

	scopes := []map[string]interface{}{}
	if scopesCollection != nil {
		for _, modelItem := range scopesCollection {
			modelMap, err := dataSourceIbmSccScopeCollectionScopeItemToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			scopes = append(scopes, modelMap)
		}
	}
	if err = d.Set("scopes", scopes); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting scope %s", err))
	}

	return nil
}

// dataSourceIbmSccProviderTypeCollectionID returns a reasonable ID for the list.
func dataSourceIbmSccScopeCollectionID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccScopeCollectionScopeItemToMap(model *securityandcompliancecenterapiv3.Scope) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	if model.AccountID != nil {
		modelMap["account_id"] = model.AccountID
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = model.InstanceID
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.CreatedOn != nil {
		modelMap["created_on"] = model.CreatedOn.String()
	}
	if model.UpdatedBy != nil {
		modelMap["updated_by"] = model.UpdatedBy
	}
	if model.UpdatedOn != nil {
		modelMap["updated_on"] = model.UpdatedOn.String()
	}
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := dataSourceIbmSccReportScopePropertyToMap(propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	return modelMap, nil
}
