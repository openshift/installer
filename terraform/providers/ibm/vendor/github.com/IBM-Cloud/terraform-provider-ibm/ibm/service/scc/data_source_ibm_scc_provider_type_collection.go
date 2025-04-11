// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccProviderTypeCollection() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccProviderTypeCollectionRead,

		Schema: map[string]*schema.Schema{
			"provider_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The array of provder type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the provider type.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the provider type.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the provider type.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provider type description.",
						},
						"s2s_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "A boolean that indicates whether the provider type is s2s-enabled.",
						},
						"instance_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of instances that can be created for the provider type.",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mode that is used to get results from provider (`PUSH` or `PULL`).",
						},
						"data_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The format of the results that a provider supports.",
						},
						"icon": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The icon of a provider in .svg format that is encoded as a base64 string.",
						},
						"label": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The label that is associated with the provider type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"text": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The text of the label.",
									},
									"tip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The text to be shown when user hover overs the label.",
									},
								},
							},
						},
						"attributes": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The attributes that are required when you're creating an instance of a provider type. The attributes field can have multiple  keys in its value. Each of those keys has a value  object that includes the type, and display name as keys. For example, `{type:\"\", display_name:\"\"}`. **NOTE;** If the provider type is s2s-enabled, which means that if the `s2s_enabled` field is set to `true`, then a CRN field of type text is required in the attributes value object.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time at which resource was created.",
						},
						"updated_at": {
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

func dataSourceIbmSccProviderTypeCollectionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	listProviderTypesOptions := &securityandcompliancecenterapiv3.ListProviderTypesOptions{}
	listProviderTypesOptions.SetInstanceID(d.Get("instance_id").(string))

	providerTypesCollection, response, err := securityAndComplianceCenterApIsClient.ListProviderTypesWithContext(context, listProviderTypesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListProviderTypesWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("ListProviderTypesWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmSccProviderTypeCollectionID(d))

	providerTypes := []map[string]interface{}{}
	if providerTypesCollection.ProviderTypes != nil {
		for _, modelItem := range providerTypesCollection.ProviderTypes {
			modelMap, err := dataSourceIbmSccProviderTypeCollectionProviderTypeItemToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			providerTypes = append(providerTypes, modelMap)
		}
	}
	if err = d.Set("provider_types", providerTypes); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting provider_types %s", err))
	}

	return nil
}

// dataSourceIbmSccProviderTypeCollectionID returns a reasonable ID for the list.
func dataSourceIbmSccProviderTypeCollectionID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccProviderTypeCollectionProviderTypeItemToMap(model *securityandcompliancecenterapiv3.ProviderType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["description"] = model.Description
	modelMap["s2s_enabled"] = model.S2sEnabled
	modelMap["instance_limit"] = flex.IntValue(model.InstanceLimit)
	modelMap["mode"] = model.Mode
	modelMap["data_type"] = model.DataType
	modelMap["icon"] = model.Icon
	if model.Label != nil {
		labelMap, err := dataSourceIbmSccProviderTypeCollectionLabelTypeToMap(model.Label)
		if err != nil {
			return modelMap, err
		}
		modelMap["label"] = []map[string]interface{}{labelMap}
	}
	attributes := make(map[string]interface{})
	for k, v := range model.Attributes {
		bytes, err := json.Marshal(v)
		if err != nil {
			return modelMap, err
		}
		attributes[k] = string(bytes)
	}
	modelMap["attributes"] = attributes
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSccProviderTypeCollectionLabelTypeToMap(model *securityandcompliancecenterapiv3.LabelType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Text != nil {
		modelMap["text"] = model.Text
	}
	if model.Tip != nil {
		modelMap["tip"] = model.Tip
	}
	return modelMap, nil
}

func dataSourceIbmSccProviderTypeCollectionAdditionalPropertyToMap(model *securityandcompliancecenterapiv3.AdditionalProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["display_name"] = model.DisplayName
	return modelMap, nil
}
