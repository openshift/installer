// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamApiKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIamApiKeyRead,

		Schema: map[string]*schema.Schema{
			"apikey_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of the API key.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the API Key details object. You need to specify this value when updating the API key to avoid stale updates.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud Resource Name of the item. Example Cloud Resource Name: 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::apikey:1234-9012-5678'.",
			},
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The API key cannot be changed if set to true.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which created the API key.",
			},
			"modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the API key. The name is not checked for uniqueness. Therefore multiple names with the same value can exist. Access is done via the UUID of the API key.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The optional description of the API key. The 'description' property is only available if a description was provided during a create of an API key.",
			},
			"iam_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The iam_id that this API key authenticates.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the account that this API key authenticates for.",
			},
		},
	}
}

func dataSourceIbmIamApiKeyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getApiKeyOptions := &iamidentityv1.GetAPIKeyOptions{}

	getApiKeyOptions.SetID(d.Get("apikey_id").(string))

	apiKey, response, err := iamIdentityClient.GetAPIKey(getApiKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetApiKey failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*apiKey.ID)

	if err = d.Set("entity_tag", apiKey.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting entity_tag: %s", err))
	}
	if err = d.Set("crn", apiKey.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("locked", apiKey.Locked); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting locked: %s", err))
	}
	if err = d.Set("created_at", apiKey.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", apiKey.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by: %s", err))
	}
	if err = d.Set("modified_at", apiKey.ModifiedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting modified_at: %s", err))
	}
	if err = d.Set("name", apiKey.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("description", apiKey.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}
	if err = d.Set("iam_id", apiKey.IamID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting iam_id: %s", err))
	}
	if err = d.Set("account_id", apiKey.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting account_id: %s", err))
	}

	return nil
}
