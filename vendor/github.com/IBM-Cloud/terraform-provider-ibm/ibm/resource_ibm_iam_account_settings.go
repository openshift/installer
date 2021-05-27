// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

const (
	accountSettings         = "ibm_iam_account_settings"
	restrictCreateServiceId = "restrict_create_service_id"
	restrictCreateApiKey    = "restrict_create_platform_apikey"
	mfa                     = "mfa"
)

func resourceIbmIamAccountSettings() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmIamAccountSettingsCreate,
		Read:     resourceIbmIamAccountSettingsRead,
		Update:   resourceIbmIamAccountSettingsUpdate,
		Delete:   resourceIbmIamAccountSettingsDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"include_history": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"restrict_create_service_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator(accountSettings, restrictCreateServiceId),
				Description:  "Defines whether or not creating a Service Id is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"restrict_create_platform_apikey": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator(accountSettings, restrictCreateApiKey),
				Description:  "Defines whether or not creating platform API keys is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"allowed_ip_addresses": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the IP addresses and subnets from which IAM tokens can be created for the account.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Version of the account settings.",
			},
			"mfa": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator(accountSettings, mfa),
				Description:  "Defines the MFA trait for the account. Valid values:  * NONE - No MFA trait set  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
			},
			"if_match": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
				Description: "Version of the account settings to be updated. Specify the version that you retrieved as entity_tag (ETag header) when reading the account. This value helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might result in stale updates.",
			},
			"history": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the Account Settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action of the history entry.",
						},
						"params": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Params of the history entry.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
			"session_expiration_in_seconds": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the session expiration in seconds for the account. Valid values:  * Any whole number between between '900' and '86400'  * NOT_SET - To unset account setting and use service default.",
			},
			"session_invalidation_in_seconds": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the period of time in seconds in which a session will be invalidated due  to inactivity. Valid values:   * Any whole number between '900' and '7200'   * NOT_SET - To unset account setting and use service default.",
			},
			"max_sessions_per_identity": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the max allowed sessions per identity required by the account. Value values: * Any whole number greater than '0'   * NOT_SET - To unset account setting and use service default.",
			},
		},
	}
}

func resourceIBMIAMAccountSettingsValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)

	restrict_values := "RESTRICTED, NOT_RESTRICTED, NOT_SET"
	mfa_values := "NONE, TOTP, TOTP4ALL, LEVEL1, LEVEL2, LEVEL3"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 restrictCreateServiceId,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              restrict_values})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 restrictCreateApiKey,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              restrict_values})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 mfa,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              mfa_values})

	ibmIAMAccountSettingsValidator := ResourceValidator{ResourceName: "ibm_iam_account_settings", Schema: validateSchema}
	return &ibmIAMAccountSettingsValidator
}

func resourceIbmIamAccountSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	getAccountSettingsOptions := &iamidentityv1.GetAccountSettingsOptions{}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	getAccountSettingsOptions.SetAccountID(userDetails.userAccount)
	if _, ok := d.GetOk("include_history"); ok {
		getAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}

	accountSettingsResponse, response, err := iamIdentityClient.GetAccountSettings(getAccountSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetAccountSettings failed %s\n%s", err, response)
		return err
	}

	d.SetId(fmt.Sprintf("%s", *accountSettingsResponse.AccountID))

	return resourceIbmIamAccountSettingsUpdate(d, meta)
}

func resourceIbmIamAccountSettingsRead(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	getAccountSettingsOptions := &iamidentityv1.GetAccountSettingsOptions{}

	getAccountSettingsOptions.SetAccountID(d.Id())
	getAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))

	accountSettingsResponse, response, err := iamIdentityClient.GetAccountSettings(getAccountSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountSettings failed %s\n%s", err, response)
		return err
	}

	if err = d.Set("restrict_create_service_id", accountSettingsResponse.RestrictCreateServiceID); err != nil {
		return fmt.Errorf("Error setting restrict_create_service_id: %s", err)
	}
	if err = d.Set("restrict_create_platform_apikey", accountSettingsResponse.RestrictCreatePlatformApikey); err != nil {
		return fmt.Errorf("Error setting restrict_create_platform_apikey: %s", err)
	}
	if err = d.Set("allowed_ip_addresses", accountSettingsResponse.AllowedIPAddresses); err != nil {
		return fmt.Errorf("Error setting allowed_ip_addresses: %s", err)
	}
	if err = d.Set("entity_tag", accountSettingsResponse.EntityTag); err != nil {
		return fmt.Errorf("Error setting entity_tag: %s", err)
	}
	if err = d.Set("mfa", accountSettingsResponse.Mfa); err != nil {
		return fmt.Errorf("Error setting mfa: %s", err)
	}
	if accountSettingsResponse.History != nil {
		history := []map[string]interface{}{}
		for _, historyItem := range accountSettingsResponse.History {
			historyItemMap := resourceIbmIamAccountSettingsEnityHistoryRecordToMap(historyItem)
			history = append(history, historyItemMap)
		}
		if err = d.Set("history", history); err != nil {
			return fmt.Errorf("Error setting history: %s", err)
		}
	}
	if err = d.Set("session_expiration_in_seconds", accountSettingsResponse.SessionExpirationInSeconds); err != nil {
		return fmt.Errorf("Error setting session_expiration_in_seconds: %s", err)
	}
	if err = d.Set("session_invalidation_in_seconds", accountSettingsResponse.SessionInvalidationInSeconds); err != nil {
		return fmt.Errorf("Error setting session_invalidation_in_seconds: %s", err)
	}
	if err = d.Set("max_sessions_per_identity", accountSettingsResponse.MaxSessionsPerIdentity); err != nil {
		return fmt.Errorf("Error setting max_sessions_per_identity: %s", err)
	}

	return nil
}

func resourceIbmIamAccountSettingsEnityHistoryRecordToMap(enityHistoryRecord iamidentityv1.EnityHistoryRecord) map[string]interface{} {
	enityHistoryRecordMap := map[string]interface{}{}

	enityHistoryRecordMap["timestamp"] = enityHistoryRecord.Timestamp
	enityHistoryRecordMap["iam_id"] = enityHistoryRecord.IamID
	enityHistoryRecordMap["iam_id_account"] = enityHistoryRecord.IamIDAccount
	enityHistoryRecordMap["action"] = enityHistoryRecord.Action
	enityHistoryRecordMap["params"] = enityHistoryRecord.Params
	enityHistoryRecordMap["message"] = enityHistoryRecord.Message

	return enityHistoryRecordMap
}

func resourceIbmIamAccountSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	updateAccountSettingsOptions := &iamidentityv1.UpdateAccountSettingsOptions{}

	updateAccountSettingsOptions.SetAccountID(d.Id())
	updateAccountSettingsOptions.SetIfMatch(d.Get("if_match").(string))

	hasChange := false

	if d.HasChange("allowed_ip_addresses") {
		allowed_ip_addresses_str := d.Get("allowed_ip_addresses").(string)
		updateAccountSettingsOptions.SetAllowedIPAddresses(allowed_ip_addresses_str)
		hasChange = true
	}

	if d.HasChange("restrict_create_service_id") {
		restrict_create_service_id_str := d.Get("restrict_create_service_id").(string)
		updateAccountSettingsOptions.SetRestrictCreateServiceID(restrict_create_service_id_str)
		hasChange = true
	}

	if d.HasChange("restrict_create_platform_apikey") {
		restrict_create_platform_apikey_str := d.Get("restrict_create_platform_apikey").(string)
		updateAccountSettingsOptions.SetRestrictCreatePlatformApikey(restrict_create_platform_apikey_str)
		hasChange = true
	}

	if d.HasChange("mfa") {
		mfa_str := d.Get("mfa").(string)
		updateAccountSettingsOptions.SetMfa(mfa_str)
		hasChange = true
	}

	if d.HasChange("session_expiration_in_seconds") {
		session_expiration_in_seconds_str := d.Get("session_expiration_in_seconds").(string)
		updateAccountSettingsOptions.SetSessionExpirationInSeconds(session_expiration_in_seconds_str)
		hasChange = true
	}

	if d.HasChange("session_invalidation_in_seconds") {
		session_invalidation_in_seconds_str := d.Get("session_invalidation_in_seconds").(string)
		updateAccountSettingsOptions.SetSessionInvalidationInSeconds(session_invalidation_in_seconds_str)
		hasChange = true
	}

	if d.HasChange("max_sessions_per_identity") {
		max_sessions_per_identity_str := d.Get("max_sessions_per_identity").(string)
		updateAccountSettingsOptions.SetMaxSessionsPerIdentity(max_sessions_per_identity_str)
		hasChange = true
	}

	if hasChange {
		_, response, err := iamIdentityClient.UpdateAccountSettings(updateAccountSettingsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountSettings failed %s\n%s", err, response)
			return err
		}
	}

	return resourceIbmIamAccountSettingsRead(d, meta)
}

func resourceIbmIamAccountSettingsDelete(d *schema.ResourceData, meta interface{}) error {

	// DELETE NOT SUPPORTED
	d.SetId("")

	return nil
}
