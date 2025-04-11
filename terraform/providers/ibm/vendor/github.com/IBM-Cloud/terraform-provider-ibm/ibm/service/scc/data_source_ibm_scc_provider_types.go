// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccProviderTypes() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccProviderTypesRead,

		Schema: map[string]*schema.Schema{
			"provider_types": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of provider_types found.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the provider type.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the provider type.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the provider type.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provider type description.",
						},
						"s2s_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "A boolean that indicates whether the provider type is s2s-enabled.",
						},
						"instance_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of instances that can be created for the provider type.",
						},
						"mode": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mode that is used to get results from provider (`PUSH` or `PULL`).",
						},
						"data_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The format of the results that a provider supports.",
						},
						"label": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The label that is associated with the provider type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"text": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The text of the label.",
									},
									"tip": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The text to be shown when user hover overs the label.",
									},
								},
							},
						},
						"attributes": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The attributes that are required when you're creating an instance of a provider type. The attributes field can have multiple  keys in its value. Each of those keys has a value  object that includes the type, and display name as keys. For example, `{type:\"\", display_name:\"\"}`. **NOTE;** If the provider type is s2s-enabled, which means that if the `s2s_enabled` field is set to `true`, then a CRN field of type text is required in the attributes value object.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time at which resource was created.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time at which resource was updated.",
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccProviderTypesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	listProviderTypesByIdOptions := &securityandcompliancecenterapiv3.ListProviderTypesOptions{}

	listProviderTypesByIdOptions.SetInstanceID(d.Get("instance_id").(string))

	providerTypeItems, response, err := securityAndComplianceCenterApIsClient.ListProviderTypesWithContext(context, listProviderTypesByIdOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProviderTypeByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProviderTypeByIDWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/provider_types", d.Get("instance_id").(string)))

	providerTypes := []map[string]interface{}{}
	for _, providerType := range providerTypeItems.ProviderTypes {
		modelMap, err := dataSourceIbmSccProviderToMap(&providerType)
		if err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting provider_type: %v\n%s", providerType, err))
		}
		providerTypes = append(providerTypes, modelMap)
	}
	if err = d.Set("provider_types", providerTypes); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting provider_types: %s", err))
	}

	return nil
}

func dataSourceIbmSccProviderToMap(model *securityandcompliancecenterapiv3.ProviderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}

	if model.Type != nil {
		modelMap["type"] = model.Type
	}

	if model.Name != nil {
		modelMap["name"] = model.Name
	}

	if model.Description != nil {
		modelMap["description"] = model.Description
	}

	if model.S2sEnabled != nil {
		modelMap["s2s_enabled"] = model.S2sEnabled
	}

	if model.InstanceLimit != nil {
		modelMap["instance_limit"] = model.InstanceLimit
	}

	if model.Mode != nil {
		modelMap["mode"] = model.Mode
	}

	if model.DataType != nil {
		modelMap["data_type"] = model.DataType
	}

	if model.Attributes != nil {
		convertedMap := make(map[string]interface{}, len(model.Attributes))
		for k, v := range model.Attributes {
			convertedMap[k] = v
		}
		modelMap["attributes"] = flex.Flatten(convertedMap)
	}

	if model.Label != nil {
		labelList := []map[string]interface{}{}
		convertedMap, err := dataSourceIbmSccProviderTypeLabelTypeToMap(model.Label)
		if err != nil {
			return modelMap, err
		}
		labelList = append(labelList, convertedMap)
		modelMap["label"] = labelList
	}

	if model.CreatedAt != nil {
		modelMap["created_at"] = flex.DateTimeToString(model.CreatedAt)
	}

	if model.UpdatedAt != nil {
		modelMap["updated_at"] = flex.DateTimeToString(model.UpdatedAt)
	}

	return modelMap, nil
}
