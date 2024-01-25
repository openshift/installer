// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func ResourceIBMAccountSettingsTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAccountSettingsTemplateCreate,
		ReadContext:   resourceIBMAccountSettingsTemplateRead,
		UpdateContext: resourceIBMAccountSettingsTemplateUpdate,
		DeleteContext: resourceIBMAccountSettingsTemplateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the account where the template resides.",
			},
			"name": {
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"name", "description", "account_settings"},
				Optional:     true,
				Description:  "The name of the trusted profile template. This is visible only in the enterprise account.",
			},
			"description": {
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"name", "description", "account_settings"},
				Optional:     true,
				Description:  "The description of the trusted profile template. Describe the template for enterprise account users.",
			},
			"account_settings": {
				Type:         schema.TypeList,
				AtLeastOneOf: []string{"name", "description", "account_settings"},
				MaxItems:     1,
				Optional:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restrict_create_service_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines whether or not creating a service ID is access controlled. Valid values:  * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service IDs, including the account owner  * NOT_RESTRICTED - all members of an account can create service IDs  * NOT_SET - to 'unset' a previous set value.",
						},
						"restrict_create_platform_apikey": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines whether or not creating platform API keys is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
						},
						"allowed_ip_addresses": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the IP addresses and subnets from which IAM tokens can be created for the account.",
						},
						"mfa": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the MFA trait for the account. Valid values:  * NONE - No MFA trait set  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
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
										Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
									},
								},
							},
						},
						"session_expiration_in_seconds": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the session expiration in seconds for the account. Valid values:  * Any whole number between between '900' and '86400'  * NOT_SET - To unset account setting and use service default.",
						},
						"session_invalidation_in_seconds": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:  * Any whole number between '900' and '7200'  * NOT_SET - To unset account setting and use service default.",
						},
						"max_sessions_per_identity": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the max allowed sessions per identity required by the account. Valid values:  * Any whole number greater than 0  * NOT_SET - To unset account setting and use service default.",
						},
						"system_access_token_expiration_in_seconds": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the access token expiration in seconds. Valid values:  * Any whole number between '900' and '3600'  * NOT_SET - To unset account setting and use service default.",
						},
						"system_refresh_token_expiration_in_seconds": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Defines the refresh token expiration in seconds. Valid values:  * Any whole number between '900' and '259200'  * NOT_SET - To unset account setting and use service default.",
						},
					},
				},
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the the template.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the the template resource.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of the the template.",
			},
			"committed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Committed flag determines if the template is ready for assignment.",
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
							Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceIBMAccountSettingsTemplateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk("template_id"); ok { // if template_id is present then we need to create a new version of this template instead
		return resourceIBMAccountSettingsTemplateCreateVersion(context, d, meta)
	}

	createAccountSettingsTemplateOptions := &iamidentityv1.CreateAccountSettingsTemplateOptions{}

	if _, ok := d.GetOk("name"); ok {
		createAccountSettingsTemplateOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createAccountSettingsTemplateOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("account_settings"); ok {
		accountSettingsModel, err := resourceIBMAccountSettingsTemplateMapToAccountSettingsComponent(d.Get("account_settings.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createAccountSettingsTemplateOptions.SetAccountSettings(accountSettingsModel)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	accountID := userDetails.UserAccount
	createAccountSettingsTemplateOptions.SetAccountID(accountID)

	accountSettingsTemplateResponse, response, err := iamIdentityClient.CreateAccountSettingsTemplateWithContext(context, createAccountSettingsTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateAccountSettingsTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateAccountSettingsTemplateWithContext failed %s\n%s", err, response))
	}

	d.SetId(buildResourceIdFromTemplateVersion(*accountSettingsTemplateResponse.ID, *accountSettingsTemplateResponse.Version))

	if d.Get("committed").(bool) {
		err := resourceIBMAccountSettingsTemplateCommit(context, d, meta)
		if err != nil {
			log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateCommit failed %s", err)
			return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateCommit failed %s", err))
		}
	}

	return resourceIBMAccountSettingsTemplateRead(context, d, meta)
}

func resourceIBMAccountSettingsTemplateCreateVersion(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createAccountSettingsTemplateVersionOptions := &iamidentityv1.CreateAccountSettingsTemplateVersionOptions{}

	id, _, err := parseResourceId(d.Get("template_id").(string))
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateRead failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateRead failed %s", err))
	}

	createAccountSettingsTemplateVersionOptions.SetTemplateID(id)

	if _, ok := d.GetOk("name"); ok {
		createAccountSettingsTemplateVersionOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createAccountSettingsTemplateVersionOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("account_settings"); ok {
		accountSettingsModel, err := resourceIBMAccountSettingsTemplateMapToAccountSettingsComponent(d.Get("account_settings.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createAccountSettingsTemplateVersionOptions.SetAccountSettings(accountSettingsModel)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	accountID := userDetails.UserAccount
	createAccountSettingsTemplateVersionOptions.SetAccountID(accountID)

	accountSettingsTemplateResponse, response, err := iamIdentityClient.CreateAccountSettingsTemplateVersionWithContext(context, createAccountSettingsTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response))
	}

	d.SetId(buildResourceIdFromTemplateVersion(*accountSettingsTemplateResponse.ID, *accountSettingsTemplateResponse.Version))

	if d.Get("committed").(bool) {
		err := resourceIBMAccountSettingsTemplateCommit(context, d, meta)
		if err != nil {
			log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateCommit failed %s", err)
			return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateCommit failed %s", err))
		}
	}

	return resourceIBMAccountSettingsTemplateRead(context, d, meta)
}

func resourceIBMAccountSettingsTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsTemplateVersionOptions := &iamidentityv1.GetAccountSettingsTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateRead failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateRead failed %s", err))
	}

	getAccountSettingsTemplateVersionOptions.SetTemplateID(id)
	getAccountSettingsTemplateVersionOptions.SetVersion(version)

	accountSettingsTemplateResponse, response, err := iamIdentityClient.GetAccountSettingsTemplateVersionWithContext(context, getAccountSettingsTemplateVersionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response))
	}

	if !core.IsNil(accountSettingsTemplateResponse.Version) {
		if err = d.Set("version", accountSettingsTemplateResponse.Version); err != nil {
			return diag.FromErr(fmt.Errorf("error setting version: %s", err))
		}
	}
	if !core.IsNil(accountSettingsTemplateResponse.AccountID) {
		if err = d.Set("account_id", accountSettingsTemplateResponse.AccountID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting account_id: %s", err))
		}
	}
	if !core.IsNil(accountSettingsTemplateResponse.Name) {
		if err = d.Set("name", accountSettingsTemplateResponse.Name); err != nil {
			return diag.FromErr(fmt.Errorf("error setting name: %s", err))
		}
	}
	if !core.IsNil(accountSettingsTemplateResponse.Description) {
		if err = d.Set("description", accountSettingsTemplateResponse.Description); err != nil {
			return diag.FromErr(fmt.Errorf("error setting description: %s", err))
		}
	}
	if !core.IsNil(accountSettingsTemplateResponse.AccountSettings) {
		accountSettingsMap, err := resourceIBMAccountSettingsTemplateAccountSettingsComponentToMap(accountSettingsTemplateResponse.AccountSettings)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("account_settings", []map[string]interface{}{accountSettingsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting account_settings: %s", err))
		}
	}
	if err = d.Set("id", accountSettingsTemplateResponse.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %s", err))
	}
	if err = d.Set("committed", accountSettingsTemplateResponse.Committed); err != nil {
		return diag.FromErr(fmt.Errorf("error setting committed: %s", err))
	}
	var history []map[string]interface{}
	if !core.IsNil(accountSettingsTemplateResponse.History) {
		for _, historyItem := range accountSettingsTemplateResponse.History {
			historyItemMap, err := resourceIBMAccountSettingsTemplateEntityHistoryRecordToMap(&historyItem)
			if err != nil {
				return diag.FromErr(err)
			}
			history = append(history, historyItemMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		return diag.FromErr(fmt.Errorf("error setting history: %s", err))
	}

	if err = d.Set("entity_tag", accountSettingsTemplateResponse.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("error setting entity_tag: %s", err))
	}
	if err = d.Set("crn", accountSettingsTemplateResponse.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("error setting crn: %s", err))
	}
	if !core.IsNil(accountSettingsTemplateResponse.CreatedAt) {
		if err = d.Set("created_at", accountSettingsTemplateResponse.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("error setting created_at: %s", err))
		}
	}
	if !core.IsNil(accountSettingsTemplateResponse.CreatedByID) {
		if err = d.Set("created_by_id", accountSettingsTemplateResponse.CreatedByID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting created_by_id: %s", err))
		}
	}
	if !core.IsNil(accountSettingsTemplateResponse.LastModifiedAt) {
		if err = d.Set("last_modified_at", accountSettingsTemplateResponse.LastModifiedAt); err != nil {
			return diag.FromErr(fmt.Errorf("error setting last_modified_at: %s", err))
		}
	}
	if !core.IsNil(accountSettingsTemplateResponse.LastModifiedByID) {
		if err = d.Set("last_modified_by_id", accountSettingsTemplateResponse.LastModifiedByID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting last_modified_by_id: %s", err))
		}
	}

	return nil
}

func resourceIBMAccountSettingsTemplateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAccountSettingsTemplateVersionOptions := &iamidentityv1.UpdateAccountSettingsTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateUpdate failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateUpdate failed %s", err))
	}

	updateAccountSettingsTemplateVersionOptions.SetTemplateID(id)
	updateAccountSettingsTemplateVersionOptions.SetVersion(version)
	updateAccountSettingsTemplateVersionOptions.SetIfMatch(d.Get("entity_tag").(string))

	hasChange := false

	if d.HasChange("name") {
		updateAccountSettingsTemplateVersionOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateAccountSettingsTemplateVersionOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("account_settings") {
		accountSettings, err := resourceIBMAccountSettingsTemplateMapToAccountSettingsComponent(d.Get("account_settings.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateAccountSettingsTemplateVersionOptions.SetAccountSettings(accountSettings)
		hasChange = true
	}

	if hasChange {
		_, response, err := iamIdentityClient.UpdateAccountSettingsTemplateVersionWithContext(context, updateAccountSettingsTemplateVersionOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response))
		}
	}

	if d.HasChange("committed") {
		if d.Get("committed").(bool) {
			err := resourceIBMAccountSettingsTemplateCommit(context, d, meta)
			if err != nil {
				log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateCommit failed %s", err)
				return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateCommit failed %s", err))
			}
		} else {
			return diag.FromErr(fmt.Errorf("A committed template cannot be uncommitted"))
		}
	}

	return resourceIBMAccountSettingsTemplateRead(context, d, meta)
}

func resourceIBMAccountSettingsTemplateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteAccountSettingsTemplateVersionOptions := &iamidentityv1.DeleteAccountSettingsTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateDelete failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateDelete failed %s", err))
	}

	deleteAccountSettingsTemplateVersionOptions.SetTemplateID(id)
	deleteAccountSettingsTemplateVersionOptions.SetVersion(version)

	response, err := iamIdentityClient.DeleteAccountSettingsTemplateVersionWithContext(context, deleteAccountSettingsTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteAccountSettingsTemplateVersionWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMAccountSettingsTemplateCommit(context context.Context, d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		return err
	}

	commitAccountSettingsTemplateVersionOptions := iamIdentityClient.NewCommitAccountSettingsTemplateOptions(id, version)
	_, err = iamIdentityClient.CommitAccountSettingsTemplateWithContext(context, commitAccountSettingsTemplateVersionOptions)
	if err != nil {
		return err
	}

	return nil
}

func resourceIBMAccountSettingsTemplateMapToAccountSettingsComponent(modelMap map[string]interface{}) (*iamidentityv1.AccountSettingsComponent, error) {
	model := &iamidentityv1.AccountSettingsComponent{}
	if modelMap["restrict_create_service_id"] != nil && modelMap["restrict_create_service_id"].(string) != "" {
		model.RestrictCreateServiceID = core.StringPtr(modelMap["restrict_create_service_id"].(string))
	}
	if modelMap["restrict_create_platform_apikey"] != nil && modelMap["restrict_create_platform_apikey"].(string) != "" {
		model.RestrictCreatePlatformApikey = core.StringPtr(modelMap["restrict_create_platform_apikey"].(string))
	}
	if modelMap["allowed_ip_addresses"] != nil && modelMap["allowed_ip_addresses"].(string) != "" {
		model.AllowedIPAddresses = core.StringPtr(modelMap["allowed_ip_addresses"].(string))
	}
	if modelMap["mfa"] != nil && modelMap["mfa"].(string) != "" {
		model.Mfa = core.StringPtr(modelMap["mfa"].(string))
	}
	if modelMap["user_mfa"] != nil {
		var userMfa []iamidentityv1.AccountSettingsUserMfa
		for _, userMfaItem := range modelMap["user_mfa"].([]interface{}) {
			userMfaItemModel, err := resourceIBMAccountSettingsTemplateMapToAccountSettingsUserMfa(userMfaItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			userMfa = append(userMfa, *userMfaItemModel)
		}
		model.UserMfa = userMfa
	}
	if modelMap["session_expiration_in_seconds"] != nil && modelMap["session_expiration_in_seconds"].(string) != "" {
		model.SessionExpirationInSeconds = core.StringPtr(modelMap["session_expiration_in_seconds"].(string))
	}
	if modelMap["session_invalidation_in_seconds"] != nil && modelMap["session_invalidation_in_seconds"].(string) != "" {
		model.SessionInvalidationInSeconds = core.StringPtr(modelMap["session_invalidation_in_seconds"].(string))
	}
	if modelMap["max_sessions_per_identity"] != nil && modelMap["max_sessions_per_identity"].(string) != "" {
		model.MaxSessionsPerIdentity = core.StringPtr(modelMap["max_sessions_per_identity"].(string))
	}
	if modelMap["system_access_token_expiration_in_seconds"] != nil && modelMap["system_access_token_expiration_in_seconds"].(string) != "" {
		model.SystemAccessTokenExpirationInSeconds = core.StringPtr(modelMap["system_access_token_expiration_in_seconds"].(string))
	}
	if modelMap["system_refresh_token_expiration_in_seconds"] != nil && modelMap["system_refresh_token_expiration_in_seconds"].(string) != "" {
		model.SystemRefreshTokenExpirationInSeconds = core.StringPtr(modelMap["system_refresh_token_expiration_in_seconds"].(string))
	}
	return model, nil
}

func resourceIBMAccountSettingsTemplateMapToAccountSettingsUserMfa(modelMap map[string]interface{}) (*iamidentityv1.AccountSettingsUserMfa, error) {
	model := &iamidentityv1.AccountSettingsUserMfa{}
	model.IamID = core.StringPtr(modelMap["iam_id"].(string))
	model.Mfa = core.StringPtr(modelMap["mfa"].(string))
	return model, nil
}

func resourceIBMAccountSettingsTemplateAccountSettingsComponentToMap(model *iamidentityv1.AccountSettingsComponent) (map[string]interface{}, error) {
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
			userMfaItemMap, err := resourceIBMAccountSettingsTemplateAccountSettingsUserMfaToMap(&userMfaItem)
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

func resourceIBMAccountSettingsTemplateAccountSettingsUserMfaToMap(model *iamidentityv1.AccountSettingsUserMfa) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["iam_id"] = model.IamID
	modelMap["mfa"] = model.Mfa
	return modelMap, nil
}

func resourceIBMAccountSettingsTemplateEntityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = model.Timestamp
	modelMap["iam_id"] = model.IamID
	modelMap["iam_id_account"] = model.IamIDAccount
	modelMap["action"] = model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = model.Message
	return modelMap, nil
}

func buildResourceIdFromTemplateVersion(id string, version int64) string {
	versionStr := strconv.Itoa(int(version))
	return fmt.Sprintf("%s/%s", id, versionStr)
}

func parseResourceId(ID string) (templateId, templateVersion string, err error) {
	if !core.IsNil(ID) {
		resourceIdParts := strings.Split(ID, "/")

		if len(resourceIdParts) == 1 {
			return resourceIdParts[0], "", nil
		}

		return resourceIdParts[0], resourceIdParts[1], nil
	}

	return "", "", errors.New("resource ID is null")
}
