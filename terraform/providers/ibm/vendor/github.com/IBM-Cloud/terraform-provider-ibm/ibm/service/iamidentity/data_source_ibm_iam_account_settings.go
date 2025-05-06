// Copyright IBM Corp. 2022 All Rights Reserved.
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

func DataSourceIBMIAMAccountSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIamAccountSettingsRead,

		Schema: map[string]*schema.Schema{
			"include_history": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the account.",
			},
			"restrict_create_service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines whether or not creating a Service Id is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"restrict_create_platform_apikey": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines whether or not creating platform API keys is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"allowed_ip_addresses": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the IP addresses and subnets from which IAM tokens can be created for the account.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the account settings.",
			},
			"mfa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the MFA trait for the account. Valid values:  * NONE - No MFA trait set  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
			},
			"user_mfa": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of users that are exempted from the MFA requirement of the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The iam_id of the user.",
						},
						"mfa": {
							Type:        schema.TypeString,
							Computed:    true,
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
							Computed:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action of the history entry.",
						},
						"params": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Params of the history entry.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
			"session_expiration_in_seconds": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the session expiration in seconds for the account. Valid values:  * Any whole number between between '900' and '86400'  * NOT_SET - To unset account setting and use service default.",
			},
			"session_invalidation_in_seconds": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:  * Any whole number between '900' and '7200'  * NOT_SET - To unset account setting and use service default.",
			},
			"max_sessions_per_identity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the max allowed sessions per identity required by the account. Valid values:  * Any whole number greater than 0  * NOT_SET - To unset account setting and use service default.",
			},
			"system_access_token_expiration_in_seconds": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the access token expiration in seconds. Valid values:  * Any whole number between '900' and '3600'  * NOT_SET - To unset account setting and use service default.",
			},
			"system_refresh_token_expiration_in_seconds": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the refresh token expiration in seconds. Valid values:  * Any whole number between '900' and '2592000'  * NOT_SET - To unset account setting and use service default.",
			},
		},
	}
}

func dataSourceIbmIamAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions := &iamidentityv1.GetAccountSettingsOptions{}
	getAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions.SetAccountID(userDetails.UserAccount)

	accountSettingsResponse, response, err := iamIdentityClient.GetAccountSettings(getAccountSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetAccountSettings failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(userDetails.UserAccount)

	if err = d.Set("account_id", accountSettingsResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting account_id: %s", err))
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
	userMfa := []map[string]interface{}{}
	if accountSettingsResponse.UserMfa != nil {
		for _, modelItem := range accountSettingsResponse.UserMfa {
			modelMap, err := dataSourceIBMIamAccountSettingsAccountSettingsUserMfaToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			userMfa = append(userMfa, modelMap)
		}
	}
	if err = d.Set("user_mfa", userMfa); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting user_mfa %s", err))
	}
	if accountSettingsResponse.History != nil {
		err = d.Set("history", dataSourceAccountSettingsResponseFlattenHistory(accountSettingsResponse.History))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting history %s", err))
		}
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

func dataSourceAccountSettingsResponseFlattenHistory(result []iamidentityv1.EnityHistoryRecord) (history []map[string]interface{}) {
	for _, historyItem := range result {
		history = append(history, dataSourceAccountSettingsResponseHistoryToMap(historyItem))
	}

	return history
}

func dataSourceIBMIamAccountSettingsAccountSettingsUserMfaToMap(model *iamidentityv1.AccountSettingsUserMfa) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IamID != nil {
		modelMap["iam_id"] = *model.IamID
	}
	if model.Mfa != nil {
		modelMap["mfa"] = *model.Mfa
	}
	return modelMap, nil
}

func dataSourceAccountSettingsResponseHistoryToMap(historyItem iamidentityv1.EnityHistoryRecord) (historyMap map[string]interface{}) {
	historyMap = map[string]interface{}{}

	if historyItem.Timestamp != nil {
		historyMap["timestamp"] = historyItem.Timestamp
	}
	if historyItem.IamID != nil {
		historyMap["iam_id"] = historyItem.IamID
	}
	if historyItem.IamIDAccount != nil {
		historyMap["iam_id_account"] = historyItem.IamIDAccount
	}
	if historyItem.Action != nil {
		historyMap["action"] = historyItem.Action
	}
	if historyItem.Params != nil {
		historyMap["params"] = historyItem.Params
	}
	if historyItem.Message != nil {
		historyMap["message"] = historyItem.Message
	}

	return historyMap
}
