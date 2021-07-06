// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
)

func resourceIBMIAMServiceAPIKey() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMServiceAPIkeyCreate,
		Read:     resourceIBMIAMServiceAPIKeyRead,
		Update:   resourceIBMIAMServiceAPIKeyUpdate,
		Delete:   resourceIBMIAMServiceAPIKeyDelete,
		Exists:   resourceIBMIAMServiceAPIKeyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Service API key",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: " description of the API key",
			},

			"iam_service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service iam_id that this API key authenticates",
			},

			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID of the API key",
			},

			"apikey": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "API key value for this API key",
			},

			"locked": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "The API key cannot be changed if set to true",
			},

			"store_value": {
				Type:             schema.TypeBool,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "Boolean value deciding whether API key value is retrievable in the future",
			},

			"file": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: applyOnce,
				Description:      "File where api key is to be stored",
			},

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "crn of the Service API Key",
			},

			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the API Key details object",
			},

			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the service which created the API key",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time Service API Key was created",
			},

			"modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time Service API Key was modified",
			},
		},
	}
}

type APIKey struct {
	Name        string
	Description string
	Apikey      string
	CreatedAt   string
	Locked      bool
}

func resourceIBMIAMServiceAPIkeyCreate(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	iamID := d.Get("iam_service_id").(string)

	createAPIKeyOptions := &iamidentityv1.CreateAPIKeyOptions{
		Name:  &name,
		IamID: &iamID,
	}

	if des, ok := d.GetOk("description"); ok {
		desString := des.(string)
		createAPIKeyOptions.Description = &desString
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	createAPIKeyOptions.AccountID = &userDetails.userAccount

	if key, ok := d.GetOk("apikey"); ok {
		apikeyString := key.(string)
		createAPIKeyOptions.Apikey = &apikeyString
	}

	if strvalue, ok := d.GetOk("store_value"); ok {
		value := strvalue.(bool)
		createAPIKeyOptions.StoreValue = &value
	}

	if lock, ok := d.GetOk("locked"); ok {
		elockstr := strconv.FormatBool(lock.(bool))
		createAPIKeyOptions.EntityLock = &elockstr
	}

	apiKey, response, err := iamIdentityClient.CreateAPIKey(createAPIKeyOptions)
	if err != nil || apiKey == nil {
		return fmt.Errorf("[DEBUG] Service API Key creation Error: %s\n%s", err, response)
	}

	d.SetId(*apiKey.ID)
	d.Set("apikey", *apiKey.Apikey)

	if keyfile, ok := d.GetOk("file"); ok {
		if err := saveToFile(apiKey, keyfile.(string)); err != nil {
			log.Printf("Error writing API Key Details to file: %s", err)
		}
	}

	return resourceIBMIAMServiceAPIKeyRead(d, meta)
}

func resourceIBMIAMServiceAPIKeyRead(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}
	apiKeyID := d.Id()

	getAPIKeyOptions := &iamidentityv1.GetAPIKeyOptions{
		ID: &apiKeyID,
	}

	apiKey, response, err := iamIdentityClient.GetAPIKey(getAPIKeyOptions)
	if err != nil || apiKey == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[DEBUG] Error retrieving Service API Key: %s\n%s", err, response)
	}
	if apiKey.Name != nil {
		d.Set("name", *apiKey.Name)
	}
	if apiKey.IamID != nil {
		d.Set("iam_service_id", *apiKey.IamID)
	}
	if apiKey.Description != nil {
		d.Set("description", *apiKey.Description)
	}
	if apiKey.AccountID != nil {
		d.Set("account_id", *apiKey.AccountID)
	}
	if *apiKey.Apikey != "" {
		d.Set("apikey", *apiKey.Apikey)
	}
	if apiKey.CRN != nil {
		d.Set("crn", *apiKey.CRN)
	}
	if apiKey.EntityTag != nil {
		d.Set("entity_tag", *apiKey.EntityTag)
	}
	if apiKey.Locked != nil {
		d.Set("locked", *apiKey.Locked)
	}
	if apiKey.CreatedBy != nil {
		d.Set("created_by", *apiKey.CreatedBy)
	}
	if apiKey.CreatedAt != nil {
		d.Set("created_at", apiKey.CreatedAt.String())
	}
	if apiKey.ModifiedAt != nil {
		d.Set("modified_at", apiKey.ModifiedAt.String())
	}

	return nil
}

func resourceIBMIAMServiceAPIKeyUpdate(d *schema.ResourceData, meta interface{}) error {

	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}
	apiKeyID := d.Id()

	getAPIKeyOptions := &iamidentityv1.GetAPIKeyOptions{
		ID: &apiKeyID,
	}

	apiKey, resp, err := iamIdentityClient.GetAPIKey(getAPIKeyOptions)
	if err != nil || apiKey == nil {
		return fmt.Errorf("[DEBUG] Error retrieving Service API Key: %s\n%s", err, resp)
	}

	updateAPIKeyOptions := &iamidentityv1.UpdateAPIKeyOptions{
		ID:      &apiKeyID,
		IfMatch: apiKey.EntityTag,
	}

	hasChange := false

	if d.HasChange("name") {
		namestr := d.Get("name").(string)
		updateAPIKeyOptions.Name = &namestr
		hasChange = true
	}

	if d.HasChange("description") {
		desc := d.Get("description").(string)
		updateAPIKeyOptions.Description = &desc
		hasChange = true
	}
	if hasChange {
		_, response, err := iamIdentityClient.UpdateAPIKey(updateAPIKeyOptions)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error updating Service API Key: %s\n%s", err, response)
		}
	}

	return resourceIBMIAMServiceAPIKeyRead(d, meta)

}

func resourceIBMIAMServiceAPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}
	apiKeyID := d.Id()

	getAPIKeyOptions := &iamidentityv1.GetAPIKeyOptions{
		ID: &apiKeyID,
	}

	_, response, err := iamIdentityClient.GetAPIKey(getAPIKeyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[DEBUG] Error retrieving Service API Key: %s\n%s", err, response)
	}

	deleteAPIKeyOptions := &iamidentityv1.DeleteAPIKeyOptions{
		ID: &apiKeyID,
	}

	resp, err := iamIdentityClient.DeleteAPIKey(deleteAPIKeyOptions)
	if err != nil {
		return fmt.Errorf("[DEBUG] Error deleting Service API Key: %s\n%s", err, resp)
	}
	d.SetId("")

	return nil
}

func resourceIBMIAMServiceAPIKeyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return false, err
	}
	apiKeyID := d.Id()

	getAPIKeyOptions := &iamidentityv1.GetAPIKeyOptions{
		ID: &apiKeyID,
	}

	apiKey, response, err := iamIdentityClient.GetAPIKey(getAPIKeyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving Service API Key: %s\n%s", err, response)
	}
	return *apiKey.ID == apiKeyID, nil
}

func saveToFile(apiKey *iamidentityv1.APIKey, filePath string) error {
	outputFilePath, err := homedir.Expand(filePath)
	if err != nil {
		return fmt.Errorf("Error generating API Key file path: %s", err)
	}

	key := &APIKey{
		Name:      *apiKey.Name,
		Apikey:    *apiKey.Apikey,
		CreatedAt: apiKey.CreatedAt.String(),
		Locked:    *apiKey.Locked,
	}
	if apiKey.Description != nil {
		key.Description = *apiKey.Description
	} else {
		key.Description = ""
	}

	out, err := json.MarshalIndent(key, "", "\t")

	err = ioutil.WriteFile(outputFilePath, out, 0666)
	if err == nil {
		log.Println("Successfully save API key information to ", outputFilePath)
	}

	return err
}
