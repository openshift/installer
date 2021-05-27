// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func dataSourceIbmIamApiKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmIamApiKeyRead,

		Schema: map[string]*schema.Schema{
			"apikey_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of the API key.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the API Key details object. You need to specify this value when updating the API key to avoid stale updates.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud Resource Name of the item. Example Cloud Resource Name: 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::apikey:1234-9012-5678'.",
			},
			"locked": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The API key cannot be changed if set to true.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which created the API key.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the API key. The name is not checked for uniqueness. Therefore multiple names with the same value can exist. Access is done via the UUID of the API key.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The optional description of the API key. The 'description' property is only available if a description was provided during a create of an API key.",
			},
			"iam_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The iam_id that this API key authenticates.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the account that this API key authenticates for.",
			},
			"apikey": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API key value. This property only contains the API key value for the following cases: create an API key, update a service ID API key that stores the API key value as retrievable, or get a service ID API key that stores the API key value as retrievable. All other operations don't return the API key value, for example all user API key related operations, except for create, don't contain the API key value.",
			},
		},
	}
}

func dataSourceIbmIamApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	getApiKeyOptions := &iamidentityv1.GetAPIKeyOptions{}

	getApiKeyOptions.SetID(d.Get("apikey_id").(string))

	apiKey, response, err := iamIdentityClient.GetAPIKey(getApiKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetApiKey failed %s\n%s", err, response)
		return err
	}

	d.SetId(*apiKey.ID)

	if err = d.Set("entity_tag", apiKey.EntityTag); err != nil {
		return fmt.Errorf("Error setting entity_tag: %s", err)
	}
	if err = d.Set("crn", apiKey.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("locked", apiKey.Locked); err != nil {
		return fmt.Errorf("Error setting locked: %s", err)
	}
	if err = d.Set("created_at", apiKey.CreatedAt.String()); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}
	if err = d.Set("created_by", apiKey.CreatedBy); err != nil {
		return fmt.Errorf("Error setting created_by: %s", err)
	}
	if err = d.Set("modified_at", apiKey.ModifiedAt.String()); err != nil {
		return fmt.Errorf("Error setting modified_at: %s", err)
	}
	if err = d.Set("name", apiKey.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("description", apiKey.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("iam_id", apiKey.IamID); err != nil {
		return fmt.Errorf("Error setting iam_id: %s", err)
	}
	if err = d.Set("account_id", apiKey.AccountID); err != nil {
		return fmt.Errorf("Error setting account_id: %s", err)
	}
	if err = d.Set("apikey", apiKey.Apikey); err != nil {
		return fmt.Errorf("Error setting apikey: %s", err)
	}

	return nil
}
