// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccProviderType() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccProviderTypeRead,

		Schema: map[string]*schema.Schema{
			"provider_type_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provider type ID.",
			},
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
			"icon": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The icon of a provider in .svg format that is encoded as a base64 string.",
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
	})
}

func dataSourceIbmSccProviderTypeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getProviderTypeByIdOptions := &securityandcompliancecenterapiv3.GetProviderTypeByIDOptions{}

	getProviderTypeByIdOptions.SetInstanceID(d.Get("instance_id").(string))
	getProviderTypeByIdOptions.SetProviderTypeID(d.Get("provider_type_id").(string))

	providerTypeItem, response, err := securityAndComplianceCenterApIsClient.GetProviderTypeByIDWithContext(context, getProviderTypeByIdOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProviderTypeByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProviderTypeByIDWithContext failed %s\n%s", err, response))
	}

	d.SetId(*providerTypeItem.ID)

	if err = d.Set("id", providerTypeItem.ID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting id: %s", err))
	}

	if err = d.Set("type", providerTypeItem.Type); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting type: %s", err))
	}

	if err = d.Set("name", providerTypeItem.Name); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting name: %s", err))
	}

	if err = d.Set("description", providerTypeItem.Description); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting description: %s", err))
	}

	if err = d.Set("s2s_enabled", providerTypeItem.S2sEnabled); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting s2s_enabled: %s", err))
	}

	if err = d.Set("instance_limit", flex.IntValue(providerTypeItem.InstanceLimit)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_limit: %s", err))
	}

	if err = d.Set("mode", providerTypeItem.Mode); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting mode: %s", err))
	}

	if err = d.Set("data_type", providerTypeItem.DataType); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting data_type: %s", err))
	}

	if err = d.Set("icon", providerTypeItem.Icon); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting icon: %s", err))
	}

	label := []map[string]interface{}{}
	if providerTypeItem.Label != nil {
		modelMap, err := dataSourceIbmSccProviderTypeLabelTypeToMap(providerTypeItem.Label)
		if err != nil {
			return diag.FromErr(err)
		}
		label = append(label, modelMap)
	}
	if err = d.Set("label", label); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting label %s", err))
	}

	if providerTypeItem.Attributes != nil {
		convertedMap := make(map[string]interface{}, len(providerTypeItem.Attributes))
		for k, v := range providerTypeItem.Attributes {
			convertedMap[k] = v
		}

		if err = d.Set("attributes", flex.Flatten(convertedMap)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting attributes: %s", err))
		}
		if err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting attributes %s", err))
		}
	}

	if err = d.Set("created_at", flex.DateTimeToString(providerTypeItem.CreatedAt)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(providerTypeItem.UpdatedAt)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_at: %s", err))
	}

	return nil
}

func dataSourceIbmSccProviderTypeLabelTypeToMap(model *securityandcompliancecenterapiv3.LabelType) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Text != nil {
		modelMap["text"] = model.Text
	}
	if model.Tip != nil {
		modelMap["tip"] = model.Tip
	}
	return modelMap, nil
}
