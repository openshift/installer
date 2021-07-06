// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func dataSourceIBMIAMAccountSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmIamAccountSettingsRead,

		Schema: map[string]*schema.Schema{
			"include_history": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the account.",
			},
			"restrict_create_service_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines whether or not creating a Service Id is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"restrict_create_platform_apikey": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines whether or not creating platform API keys is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"allowed_ip_addresses": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the IP addresses and subnets from which IAM tokens can be created for the account.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the account settings.",
			},
			"mfa": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the MFA trait for the account. Valid values:  * NONE - No MFA trait set  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
			},
			"history": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the Account Settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action of the history entry.",
						},
						"params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Params of the history entry.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
			"session_expiration_in_seconds": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the session expiration in seconds for the account. Valid values:  * Any whole number between between '900' and '86400'  * NOT_SET - To unset account setting and use service default.",
			},
			"session_invalidation_in_seconds": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the period of time in seconds in which a session will be invalidated due  to inactivity. Valid values:   * Any whole number between '900' and '7200'   * NOT_SET - To unset account setting and use service default.",
			},
			"max_sessions_per_identity": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the max allowed sessions per identity required by the account. Value values: * Any whole number greater than '0'   * NOT_SET - To unset account setting and use service default.",
			},
		},
	}
}

func dataSourceIbmIamAccountSettingsRead(d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	getAccountSettingsOptions := &iamidentityv1.GetAccountSettingsOptions{}
	getAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	getAccountSettingsOptions.SetAccountID(userDetails.userAccount)

	accountSettingsResponse, response, err := iamIdentityClient.GetAccountSettings(getAccountSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetAccountSettings failed %s\n%s", err, response)
		return err
	}

	d.SetId(userDetails.userAccount)

	if err = d.Set("account_id", accountSettingsResponse.AccountID); err != nil {
		return fmt.Errorf("Error setting account_id: %s", err)
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
		err = d.Set("history", dataSourceAccountSettingsResponseFlattenHistory(accountSettingsResponse.History))
		if err != nil {
			return fmt.Errorf("Error setting history %s", err)
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

func dataSourceAccountSettingsResponseFlattenHistory(result []iamidentityv1.EnityHistoryRecord) (history []map[string]interface{}) {
	for _, historyItem := range result {
		history = append(history, dataSourceAccountSettingsResponseHistoryToMap(historyItem))
	}

	return history
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
