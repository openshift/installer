// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

const (
	accountSettings         = "ibm_iam_account_settings"
	restrictCreateServiceId = "restrict_create_service_id"
	restrictCreateApiKey    = "restrict_create_platform_apikey"
	mfa                     = "mfa"
)

func ResourceIBMIAMAccountSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIamAccountSettingsCreate,
		ReadContext:   resourceIbmIamAccountSettingsRead,
		UpdateContext: resourceIbmIamAccountSettingsUpdate,
		DeleteContext: resourceIbmIamAccountSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"include_history": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"restrict_create_service_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator(accountSettings, "restrict_create_service_id"),
				Description:  "Defines whether or not creating a Service Id is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"restrict_create_platform_apikey": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator(accountSettings, "restrict_create_platform_apikey"),
				Description:  "Defines whether or not creating platform API keys is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"allowed_ip_addresses": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the IP addresses and subnets from which IAM tokens can be created for the account.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Version of the account settings.",
			},
			"mfa": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator(accountSettings, "mfa"),
				Description:  "Defines the MFA trait for the account. Valid values:  * NONE - No MFA trait set  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
			},
			"if_match": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
				Description: "Version of the account settings to be updated. Specify the version that you retrieved as entity_tag (ETag header) when reading the account. This value helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might result in stale updates.",
			},
			"user_mfa": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of users that are exempted from the MFA requirement of the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The iam_id of the user.",
						},
						"mfa": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
						},
					},
				},
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the Account Settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action of the history entry.",
						},
						"params": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Params of the history entry.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"message": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
			"session_expiration_in_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the session expiration in seconds for the account. Valid values:  * Any whole number between between '900' and '86400'  * NOT_SET - To unset account setting and use service default.",
			},
			"session_invalidation_in_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:  * Any whole number between '900' and '7200'  * NOT_SET - To unset account setting and use service default.",
			},
			"max_sessions_per_identity": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the max allowed sessions per identity required by the account. Value values:  * Any whole number greater than 0  * NOT_SET - To unset account setting and use service default.",
			},
			"system_access_token_expiration_in_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the access token expiration in seconds. Valid values:  * Any whole number between '900' and '3600'  * NOT_SET - To unset account setting and use service default.",
			},
			"system_refresh_token_expiration_in_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the refresh token expiration in seconds. Valid values:  * Any whole number between '900' and '2592000'  * NOT_SET - To unset account setting and use service default.",
			},
		},
	}
}

func ResourceIBMIAMAccountSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)

	restrict_values := "RESTRICTED, NOT_RESTRICTED, NOT_SET"
	mfa_values := "NONE, TOTP, TOTP4ALL, LEVEL1, LEVEL2, LEVEL3"

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 restrictCreateServiceId,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              restrict_values})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 restrictCreateApiKey,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              restrict_values})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 mfa,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              mfa_values})

	ibmIAMAccountSettingsValidator := validate.ResourceValidator{ResourceName: "ibm_iam_account_settings", Schema: validateSchema}
	return &ibmIAMAccountSettingsValidator
}

func resourceIbmIamAccountSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions := &iamidentityv1.GetAccountSettingsOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	getAccountSettingsOptions.SetAccountID(userDetails.UserAccount)
	if _, ok := d.GetOk("include_history"); ok {
		getAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}

	accountSettingsResponse, response, err := iamIdentityClient.GetAccountSettings(getAccountSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetAccountSettings failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*accountSettingsResponse.AccountID)

	return resourceIbmIamAccountSettingsUpdate(context, d, meta)
}

func resourceIbmIamAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	if err = d.Set("restrict_create_service_id", accountSettingsResponse.RestrictCreateServiceID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting restrict_create_service_id: %s", err))
	}
	if err = d.Set("restrict_create_platform_apikey", accountSettingsResponse.RestrictCreatePlatformApikey); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting restrict_create_platform_apikey: %s", err))
	}
	if err = d.Set("allowed_ip_addresses", accountSettingsResponse.AllowedIPAddresses); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting allowed_ip_addresses: %s", err))
	}
	if err = d.Set("entity_tag", accountSettingsResponse.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting entity_tag: %s", err))
	}
	if err = d.Set("mfa", accountSettingsResponse.Mfa); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting mfa: %s", err))
	}
	if accountSettingsResponse.History != nil {
		history := []map[string]interface{}{}
		for _, historyItem := range accountSettingsResponse.History {
			historyItemMap := resourceIbmIamAccountSettingsEnityHistoryRecordToMap(historyItem)
			history = append(history, historyItemMap)
		}
		if err = d.Set("history", history); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting history: %s", err))
		}
	}
	userMfa := []map[string]interface{}{}
	if accountSettingsResponse.UserMfa != nil {
		for _, userMfaItem := range accountSettingsResponse.UserMfa {
			userMfaItemMap, err := resourceIBMIamAccountSettingsAccountSettingsUserMfaToMap(&userMfaItem)
			if err != nil {
				return diag.FromErr(err)
			}
			userMfa = append(userMfa, userMfaItemMap)
		}
	}
	if err = d.Set("user_mfa", userMfa); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting user_mfa: %s", err))
	}
	if err = d.Set("session_expiration_in_seconds", accountSettingsResponse.SessionExpirationInSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting session_expiration_in_seconds: %s", err))
	}
	if err = d.Set("session_invalidation_in_seconds", accountSettingsResponse.SessionInvalidationInSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting session_invalidation_in_seconds: %s", err))
	}
	if err = d.Set("max_sessions_per_identity", accountSettingsResponse.MaxSessionsPerIdentity); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting max_sessions_per_identity: %s", err))
	}
	if err = d.Set("system_access_token_expiration_in_seconds", accountSettingsResponse.SystemAccessTokenExpirationInSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting system_access_token_expiration_in_seconds: %s", err))
	}
	if err = d.Set("system_refresh_token_expiration_in_seconds", accountSettingsResponse.SystemRefreshTokenExpirationInSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting system_refresh_token_expiration_in_seconds: %s", err))
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

func resourceIbmIamAccountSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
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
	var user_mfa []iamidentityv1.AccountSettingsUserMfa
	if d.HasChange("user_mfa") {
		for _, e := range d.Get("user_mfa").([]interface{}) {
			value := e.(map[string]interface{})
			userMfaItem := resourceIBMIamAccountSettingsMapToAccountSettingsUserMfa(value)
			user_mfa = append(user_mfa, userMfaItem)
		}
		updateAccountSettingsOptions.SetUserMfa(user_mfa)
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
	if d.HasChange("system_access_token_expiration_in_seconds") {
		updateAccountSettingsOptions.SetSystemAccessTokenExpirationInSeconds(d.Get("system_access_token_expiration_in_seconds").(string))
		hasChange = true
	}
	if d.HasChange("system_refresh_token_expiration_in_seconds") {
		updateAccountSettingsOptions.SetSystemRefreshTokenExpirationInSeconds(d.Get("system_refresh_token_expiration_in_seconds").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := iamIdentityClient.UpdateAccountSettings(updateAccountSettingsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountSettings failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIamAccountSettingsRead(context, d, meta)
}

func resourceIBMIamAccountSettingsMapToAccountSettingsUserMfa(userMfaMap map[string]interface{}) iamidentityv1.AccountSettingsUserMfa {
	userMfa := iamidentityv1.AccountSettingsUserMfa{}
	userMfa.IamID = core.StringPtr(userMfaMap["iam_id"].(string))
	userMfa.Mfa = core.StringPtr(userMfaMap["mfa"].(string))
	return userMfa
}

func resourceIBMIamAccountSettingsAccountSettingsUserMfaToMap(model *iamidentityv1.AccountSettingsUserMfa) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["iam_id"] = model.IamID
	modelMap["mfa"] = model.Mfa
	return modelMap, nil
}

func resourceIbmIamAccountSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// DELETE NOT SUPPORTED
	d.SetId("")

	return nil
}
