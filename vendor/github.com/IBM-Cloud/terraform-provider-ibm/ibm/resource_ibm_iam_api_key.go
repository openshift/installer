// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func resourceIbmIamApiKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmIamApiKeyCreate,
		Read:     resourceIbmIamApiKeyRead,
		Update:   resourceIbmIamApiKeyUpdate,
		Delete:   resourceIbmIamApiKeyDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the API key. The name is not checked for uniqueness. Therefore multiple names with the same value can exist. Access is done via the UUID of the API key.",
			},
			"iam_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The iam_id that this API key authenticates.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional description of the API key. The 'description' property is only available if a description was provided during a create of an API key.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID of the API key.",
			},
			"apikey": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "You can optionally passthrough the API key value for this API key. If passed, NO validation of that apiKey value is done, i.e. the value can be non-URL safe. If omitted, the API key management will create an URL safe opaque API key value. The value of the API key is checked for uniqueness. Please ensure enough variations when passing in this value.",
			},
			"store_value": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Send true or false to set whether the API key value is retrievable in the future by using the Get details of an API key request. If you create an API key for a user, you must specify `false` or omit the value. We don't allow storing of API keys for users.",
			},
			"entity_lock": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "false",
				Description: "Indicates if the API key is locked for further write operations. False by default.",
			},
			"apikey_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of this API Key.",
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
		},
	}
}

func resourceIbmIamApiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	createApiKeyOptions := &iamidentityv1.CreateAPIKeyOptions{}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	iamID := userDetails.userID
	accountID := userDetails.userAccount

	createApiKeyOptions.SetName(d.Get("name").(string))
	createApiKeyOptions.SetIamID(iamID)
	createApiKeyOptions.SetAccountID(accountID)

	if _, ok := d.GetOk("description"); ok {
		createApiKeyOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("apikey"); ok {
		createApiKeyOptions.SetApikey(d.Get("apikey").(string))
	}
	if _, ok := d.GetOk("store_value"); ok {
		createApiKeyOptions.SetStoreValue(d.Get("store_value").(bool))
	}
	if _, ok := d.GetOk("locked"); ok {
		createApiKeyOptions.SetEntityLock(d.Get("locked").(string))
	}

	apiKey, response, err := iamIdentityClient.CreateAPIKey(createApiKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateApiKey failed %s\n%s", err, response)
		return err
	}

	d.SetId(*apiKey.ID)
	d.Set("apikey", *apiKey.Apikey)

	if keyfile, ok := d.GetOk("file"); ok {
		if err := saveToFile(apiKey, keyfile.(string)); err != nil {
			log.Printf("Error writing API Key Details to file: %s", err)
		}
	}

	return resourceIbmIamApiKeyRead(d, meta)
}

func resourceIbmIamApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	getApiKeyOptions := &iamidentityv1.GetAPIKeyOptions{}

	getApiKeyOptions.SetID(d.Id())

	apiKey, response, err := iamIdentityClient.GetAPIKey(getApiKeyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetApiKey failed %s\n%s", err, response)
		return err
	}

	if err = d.Set("name", apiKey.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("iam_id", apiKey.IamID); err != nil {
		return fmt.Errorf("Error setting iam_id: %s", err)
	}
	if err = d.Set("description", apiKey.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("account_id", apiKey.AccountID); err != nil {
		return fmt.Errorf("Error setting account_id: %s", err)
	}
	if err = d.Set("apikey", apiKey.Apikey); err != nil {
		return fmt.Errorf("Error setting apikey: %s", err)
	}
	if err = d.Set("locked", apiKey.Locked); err != nil {
		return fmt.Errorf("Error setting entity_lock: %s", err)
	}
	if err = d.Set("apikey_id", apiKey.ID); err != nil {
		return fmt.Errorf("Error setting id: %s", err)
	}
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

	return nil
}

func resourceIbmIamApiKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	updateApiKeyOptions := &iamidentityv1.UpdateAPIKeyOptions{}

	updateApiKeyOptions.SetIfMatch("*")
	updateApiKeyOptions.SetID(d.Id())
	updateApiKeyOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("description"); ok {
		updateApiKeyOptions.SetDescription(d.Get("description").(string))
	}
	_, response, err := iamIdentityClient.UpdateAPIKey(updateApiKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateApiKey failed %s\n%s", err, response)
		return err
	}

	return resourceIbmIamApiKeyRead(d, meta)
}

func resourceIbmIamApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	deleteApiKeyOptions := &iamidentityv1.DeleteAPIKeyOptions{}

	deleteApiKeyOptions.SetID(d.Id())

	response, err := iamIdentityClient.DeleteAPIKey(deleteApiKeyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteApiKey failed %s\n%s", err, response)
		return err
	}

	d.SetId("")

	return nil
}
