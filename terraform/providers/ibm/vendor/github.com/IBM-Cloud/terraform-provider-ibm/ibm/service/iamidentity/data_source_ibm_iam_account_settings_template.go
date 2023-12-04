// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMAccountSettingsTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAccountSettingsTemplateRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the account settings template.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version of the account settings template.",
			},
			"include_history": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the the template.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the account where the template resides.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the trusted profile template. This is visible only in the enterprise account.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the trusted profile template. Describe the template for enterprise account users.",
			},
			"committed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Committed flag determines if the template is ready for assignment.",
			},
			"account_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restrict_create_service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines whether or not creating a service ID is access controlled. Valid values:  * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service IDs, including the account owner  * NOT_RESTRICTED - all members of an account can create service IDs  * NOT_SET - to 'unset' a previous set value.",
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
										Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
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
							Description: "Defines the refresh token expiration in seconds. Valid values:  * Any whole number between '900' and '259200'  * NOT_SET - To unset account setting and use service default.",
						},
					},
				},
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the Template.",
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
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Entity tag for this templateId-version combination.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud resource name.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template Created At.",
			},
			"created_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the creator.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template last modified at.",
			},
			"last_modified_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the identity that made the latest modification.",
			},
		},
	}
}

func dataSourceIBMAccountSettingsTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsTemplateVersionOptions := &iamidentityv1.GetAccountSettingsTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Get("template_id").(string))
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateRead failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateRead failed %s", err))
	}
	if version == "" {
		version = d.Get("version").(string)
	}

	getAccountSettingsTemplateVersionOptions.SetTemplateID(id)
	getAccountSettingsTemplateVersionOptions.SetVersion(version)

	if _, ok := d.GetOk("include_history"); ok {
		getAccountSettingsTemplateVersionOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}

	accountSettingsTemplateResponse, response, err := iamIdentityClient.GetAccountSettingsTemplateVersionWithContext(context, getAccountSettingsTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] GetAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response))
	}

	d.SetId(buildResourceIdFromTemplateVersion(*accountSettingsTemplateResponse.ID, *accountSettingsTemplateResponse.Version))

	if err = d.Set("id", accountSettingsTemplateResponse.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	if !core.IsNil(accountSettingsTemplateResponse.Version) {
		versionStr := strconv.Itoa(int(*accountSettingsTemplateResponse.Version))
		if err = d.Set("version", versionStr); err != nil {
			return diag.FromErr(fmt.Errorf("error setting version: %s", err))
		}
	}

	if err = d.Set("account_id", accountSettingsTemplateResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}

	if err = d.Set("name", accountSettingsTemplateResponse.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("description", accountSettingsTemplateResponse.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("committed", accountSettingsTemplateResponse.Committed); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting committed: %s", err))
	}

	var accountSettings []map[string]interface{}
	if accountSettingsTemplateResponse.AccountSettings != nil {
		modelMap, err := dataSourceIBMAccountSettingsTemplateAccountSettingsComponentToMap(accountSettingsTemplateResponse.AccountSettings)
		if err != nil {
			return diag.FromErr(err)
		}
		accountSettings = append(accountSettings, modelMap)
	}
	if err = d.Set("account_settings", accountSettings); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_settings %s", err))
	}

	var history []map[string]interface{}
	if accountSettingsTemplateResponse.History != nil {
		for _, modelItem := range accountSettingsTemplateResponse.History {
			modelMap, err := dataSourceIBMAccountSettingsTemplateEnityHistoryRecordToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			history = append(history, modelMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting history %s", err))
	}

	if err = d.Set("entity_tag", accountSettingsTemplateResponse.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}

	if err = d.Set("crn", accountSettingsTemplateResponse.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("created_at", accountSettingsTemplateResponse.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("created_by_id", accountSettingsTemplateResponse.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}

	if err = d.Set("last_modified_at", accountSettingsTemplateResponse.LastModifiedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_at: %s", err))
	}

	if err = d.Set("last_modified_by_id", accountSettingsTemplateResponse.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_by_id: %s", err))
	}

	return nil
}

func dataSourceIBMAccountSettingsTemplateAccountSettingsComponentToMap(model *iamidentityv1.AccountSettingsComponent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RestrictCreateServiceID != nil {
		modelMap["restrict_create_service_id"] = model.RestrictCreateServiceID
	}
	if model.RestrictCreatePlatformApikey != nil {
		modelMap["restrict_create_platform_apikey"] = model.RestrictCreatePlatformApikey
	}
	if model.AllowedIPAddresses != nil {
		modelMap["allowed_ip_addresses"] = model.AllowedIPAddresses
	}
	if model.Mfa != nil {
		modelMap["mfa"] = model.Mfa
	}
	if model.UserMfa != nil {
		var userMfa []map[string]interface{}
		for _, userMfaItem := range model.UserMfa {
			userMfaItemMap, err := dataSourceIBMAccountSettingsTemplateAccountSettingsUserMfaToMap(&userMfaItem)
			if err != nil {
				return modelMap, err
			}
			userMfa = append(userMfa, userMfaItemMap)
		}
		modelMap["user_mfa"] = userMfa
	}
	if model.SessionExpirationInSeconds != nil {
		modelMap["session_expiration_in_seconds"] = model.SessionExpirationInSeconds
	}
	if model.SessionInvalidationInSeconds != nil {
		modelMap["session_invalidation_in_seconds"] = model.SessionInvalidationInSeconds
	}
	if model.MaxSessionsPerIdentity != nil {
		modelMap["max_sessions_per_identity"] = model.MaxSessionsPerIdentity
	}
	if model.SystemAccessTokenExpirationInSeconds != nil {
		modelMap["system_access_token_expiration_in_seconds"] = model.SystemAccessTokenExpirationInSeconds
	}
	if model.SystemRefreshTokenExpirationInSeconds != nil {
		modelMap["system_refresh_token_expiration_in_seconds"] = model.SystemRefreshTokenExpirationInSeconds
	}
	return modelMap, nil
}

func dataSourceIBMAccountSettingsTemplateAccountSettingsUserMfaToMap(model *iamidentityv1.AccountSettingsUserMfa) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["iam_id"] = model.IamID
	modelMap["mfa"] = model.Mfa
	return modelMap, nil
}

func dataSourceIBMAccountSettingsTemplateEnityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = model.Timestamp
	modelMap["iam_id"] = model.IamID
	modelMap["iam_id_account"] = model.IamIDAccount
	modelMap["action"] = model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = model.Message
	return modelMap, nil
}
