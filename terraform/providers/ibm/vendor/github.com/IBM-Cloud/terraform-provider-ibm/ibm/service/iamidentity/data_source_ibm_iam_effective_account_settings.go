// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.93.0-c40121e6-20240729-182103
 */

package iamidentity

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamEffectiveAccountSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamEffectiveAccountSettingsRead,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of the account.",
			},
			"include_history": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"resolve_user_mfa": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enrich MFA exemptions with user information.",
			},
			"context": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Context with key properties for problem determination.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transaction_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The transaction ID of the inbound REST request.",
						},
						"operation": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operation of the inbound REST request.",
						},
						"user_agent": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user agent of the inbound REST request.",
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of that cluster.",
						},
						"instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID of the server instance processing the request.",
						},
						"thread_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The thread ID of the server instance processing the request.",
						},
						"host": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host of the server instance processing the request.",
						},
						"start_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start time of the request.",
						},
						"end_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The finish time of the request.",
						},
						"elapsed_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The elapsed time in msec.",
						},
						"cluster_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster name.",
						},
					},
				},
			},
			"effective": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restrict_create_service_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines whether or not creating a service ID is access controlled. Valid values:  * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service IDs, including the account owner  * NOT_RESTRICTED - all members of an account can create service IDs  * NOT_SET - to 'unset' a previous set value.",
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
						"mfa": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
						},
						"user_mfa": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of users that are exempted from the MFA requirement of the account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iam_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The iam_id of the user.",
									},
									"mfa": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the user account.",
									},
									"user_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "userName of the user.",
									},
									"email": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "email of the user.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "optional description.",
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
							Description: "Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:  * Any whole number between '900' and '7200'  * NOT_SET - To unset account setting and use service default.",
						},
						"max_sessions_per_identity": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the max allowed sessions per identity required by the account. Valid values:  * Any whole number greater than 0  * NOT_SET - To unset account setting and use service default.",
						},
						"system_access_token_expiration_in_seconds": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the access token expiration in seconds. Valid values:  * Any whole number between '900' and '3600'  * NOT_SET - To unset account setting and use service default.",
						},
						"system_refresh_token_expiration_in_seconds": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the refresh token expiration in seconds. Valid values:  * Any whole number between '900' and '259200'  * NOT_SET - To unset account setting and use service default.",
						},
					},
				},
			},
			"account": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of the account.",
						},
						"restrict_create_service_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines whether or not creating a service ID is access controlled. Valid values:  * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service IDs, including the account owner  * NOT_RESTRICTED - all members of an account can create service IDs  * NOT_SET - to 'unset' a previous set value.",
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
						"mfa": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
						},
						"user_mfa": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of users that are exempted from the MFA requirement of the account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iam_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The iam_id of the user.",
									},
									"mfa": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the user account.",
									},
									"user_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "userName of the user.",
									},
									"email": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "email of the user.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "optional description.",
									},
								},
							},
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
							Description: "Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:  * Any whole number between '900' and '7200'  * NOT_SET - To unset account setting and use service default.",
						},
						"max_sessions_per_identity": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the max allowed sessions per identity required by the account. Valid values:  * Any whole number greater than 0  * NOT_SET - To unset account setting and use service default.",
						},
						"system_access_token_expiration_in_seconds": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the access token expiration in seconds. Valid values:  * Any whole number between '900' and '3600'  * NOT_SET - To unset account setting and use service default.",
						},
						"system_refresh_token_expiration_in_seconds": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the refresh token expiration in seconds. Valid values:  * Any whole number between '900' and '259200'  * NOT_SET - To unset account setting and use service default.",
						},
					},
				},
			},
			"assigned_templates": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "assigned template section.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template Id.",
						},
						"template_version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Template version.",
						},
						"template_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template name.",
						},
						"restrict_create_service_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines whether or not creating a service ID is access controlled. Valid values:  * RESTRICTED - only users assigned the 'Service ID creator' role on the IAM Identity Service can create service IDs, including the account owner  * NOT_RESTRICTED - all members of an account can create service IDs  * NOT_SET - to 'unset' a previous set value.",
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
						"mfa": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
						},
						"user_mfa": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of users that are exempted from the MFA requirement of the account.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iam_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The iam_id of the user.",
									},
									"mfa": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Defines the MFA requirement for the user. Valid values:  * NONE - No MFA trait set  * NONE_NO_ROPC- No MFA, disable CLI logins with only a password  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the user account.",
									},
									"user_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "userName of the user.",
									},
									"email": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "email of the user.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "optional description.",
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
							Description: "Defines the period of time in seconds in which a session will be invalidated due to inactivity. Valid values:  * Any whole number between '900' and '7200'  * NOT_SET - To unset account setting and use service default.",
						},
						"max_sessions_per_identity": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the max allowed sessions per identity required by the account. Valid values:  * Any whole number greater than 0  * NOT_SET - To unset account setting and use service default.",
						},
						"system_access_token_expiration_in_seconds": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the access token expiration in seconds. Valid values:  * Any whole number between '900' and '3600'  * NOT_SET - To unset account setting and use service default.",
						},
						"system_refresh_token_expiration_in_seconds": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the refresh token expiration in seconds. Valid values:  * Any whole number between '900' and '259200'  * NOT_SET - To unset account setting and use service default.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIamEffectiveAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_effective_account_settings", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getEffectiveAccountSettingsOptions := &iamidentityv1.GetEffectiveAccountSettingsOptions{}

	getEffectiveAccountSettingsOptions.SetAccountID(d.Get("account_id").(string))
	if _, ok := d.GetOk("include_history"); ok {
		getEffectiveAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}
	if _, ok := d.GetOk("resolve_user_mfa"); ok {
		getEffectiveAccountSettingsOptions.SetResolveUserMfa(d.Get("resolve_user_mfa").(bool))
	}

	effectiveAccountSettingsResponse, _, err := iamIdentityClient.GetEffectiveAccountSettingsWithContext(context, getEffectiveAccountSettingsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetEffectiveAccountSettingsWithContext failed: %s", err.Error()), "(Data) ibm_iam_effective_account_settings", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIamEffectiveAccountSettingsID(d))

	if !core.IsNil(effectiveAccountSettingsResponse.Context) {
		context := []map[string]interface{}{}
		contextMap, err := DataSourceIBMIamEffectiveAccountSettingsResponseContextToMap(effectiveAccountSettingsResponse.Context)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_effective_account_settings", "read", "context-to-map").GetDiag()
		}
		context = append(context, contextMap)
		if err = d.Set("context", context); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting context: %s", err), "(Data) ibm_iam_effective_account_settings", "read", "set-context").GetDiag()
		}
	}

	effective := []map[string]interface{}{}
	effectiveMap, err := DataSourceIBMIamEffectiveAccountSettingsAccountSettingsEffectiveSectionToMap(effectiveAccountSettingsResponse.Effective)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_effective_account_settings", "read", "effective-to-map").GetDiag()
	}
	effective = append(effective, effectiveMap)
	if err = d.Set("effective", effective); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting effective: %s", err), "(Data) ibm_iam_effective_account_settings", "read", "set-effective").GetDiag()
	}

	account := []map[string]interface{}{}
	accountMap, err := DataSourceIBMIamEffectiveAccountSettingsAccountSettingsAccountSectionToMap(effectiveAccountSettingsResponse.Account)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_effective_account_settings", "read", "account-to-map").GetDiag()
	}
	account = append(account, accountMap)
	if err = d.Set("account", account); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting account: %s", err), "(Data) ibm_iam_effective_account_settings", "read", "set-account").GetDiag()
	}

	if !core.IsNil(effectiveAccountSettingsResponse.AssignedTemplates) {
		assignedTemplates := []map[string]interface{}{}
		for _, assignedTemplatesItem := range effectiveAccountSettingsResponse.AssignedTemplates {
			assignedTemplatesItemMap, err := DataSourceIBMIamEffectiveAccountSettingsAccountSettingsAssignedTemplatesSectionToMap(&assignedTemplatesItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_effective_account_settings", "read", "assigned_templates-to-map").GetDiag()
			}
			assignedTemplates = append(assignedTemplates, assignedTemplatesItemMap)
		}
		if err = d.Set("assigned_templates", assignedTemplates); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting assigned_templates: %s", err), "(Data) ibm_iam_effective_account_settings", "read", "set-assigned_templates").GetDiag()
		}
	}

	return nil
}

// dataSourceIBMIamEffectiveAccountSettingsID returns a reasonable ID for the list.
func dataSourceIBMIamEffectiveAccountSettingsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIamEffectiveAccountSettingsResponseContextToMap(model *iamidentityv1.ResponseContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TransactionID != nil {
		modelMap["transaction_id"] = *model.TransactionID
	}
	if model.Operation != nil {
		modelMap["operation"] = *model.Operation
	}
	if model.UserAgent != nil {
		modelMap["user_agent"] = *model.UserAgent
	}
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = *model.InstanceID
	}
	if model.ThreadID != nil {
		modelMap["thread_id"] = *model.ThreadID
	}
	if model.Host != nil {
		modelMap["host"] = *model.Host
	}
	if model.StartTime != nil {
		modelMap["start_time"] = *model.StartTime
	}
	if model.EndTime != nil {
		modelMap["end_time"] = *model.EndTime
	}
	if model.ElapsedTime != nil {
		modelMap["elapsed_time"] = *model.ElapsedTime
	}
	if model.ClusterName != nil {
		modelMap["cluster_name"] = *model.ClusterName
	}
	return modelMap, nil
}

func DataSourceIBMIamEffectiveAccountSettingsAccountSettingsEffectiveSectionToMap(model *iamidentityv1.AccountSettingsEffectiveSection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RestrictCreateServiceID != nil {
		modelMap["restrict_create_service_id"] = *model.RestrictCreateServiceID
	}
	if model.RestrictCreatePlatformApikey != nil {
		modelMap["restrict_create_platform_apikey"] = *model.RestrictCreatePlatformApikey
	}
	if model.AllowedIPAddresses != nil {
		modelMap["allowed_ip_addresses"] = *model.AllowedIPAddresses
	}
	if model.Mfa != nil {
		modelMap["mfa"] = *model.Mfa
	}
	if model.UserMfa != nil {
		userMfa := []map[string]interface{}{}
		for _, userMfaItem := range model.UserMfa {
			userMfaItemMap, err := DataSourceIBMIamEffectiveAccountSettingsEffectiveAccountSettingsUserMfaToMap(&userMfaItem)
			if err != nil {
				return modelMap, err
			}
			userMfa = append(userMfa, userMfaItemMap)
		}
		modelMap["user_mfa"] = userMfa
	}
	if model.SessionExpirationInSeconds != nil {
		modelMap["session_expiration_in_seconds"] = *model.SessionExpirationInSeconds
	}
	if model.SessionInvalidationInSeconds != nil {
		modelMap["session_invalidation_in_seconds"] = *model.SessionInvalidationInSeconds
	}
	if model.MaxSessionsPerIdentity != nil {
		modelMap["max_sessions_per_identity"] = *model.MaxSessionsPerIdentity
	}
	if model.SystemAccessTokenExpirationInSeconds != nil {
		modelMap["system_access_token_expiration_in_seconds"] = *model.SystemAccessTokenExpirationInSeconds
	}
	if model.SystemRefreshTokenExpirationInSeconds != nil {
		modelMap["system_refresh_token_expiration_in_seconds"] = *model.SystemRefreshTokenExpirationInSeconds
	}
	return modelMap, nil
}

func DataSourceIBMIamEffectiveAccountSettingsEffectiveAccountSettingsUserMfaToMap(model *iamidentityv1.EffectiveAccountSettingsUserMfa) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["iam_id"] = *model.IamID
	modelMap["mfa"] = *model.Mfa
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.UserName != nil {
		modelMap["user_name"] = *model.UserName
	}
	if model.Email != nil {
		modelMap["email"] = *model.Email
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	return modelMap, nil
}

func DataSourceIBMIamEffectiveAccountSettingsAccountSettingsAccountSectionToMap(model *iamidentityv1.AccountSettingsAccountSection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AccountID != nil {
		modelMap["account_id"] = *model.AccountID
	}
	if model.RestrictCreateServiceID != nil {
		modelMap["restrict_create_service_id"] = *model.RestrictCreateServiceID
	}
	if model.RestrictCreatePlatformApikey != nil {
		modelMap["restrict_create_platform_apikey"] = *model.RestrictCreatePlatformApikey
	}
	if model.AllowedIPAddresses != nil {
		modelMap["allowed_ip_addresses"] = *model.AllowedIPAddresses
	}
	if model.Mfa != nil {
		modelMap["mfa"] = *model.Mfa
	}
	if model.UserMfa != nil {
		userMfa := []map[string]interface{}{}
		for _, userMfaItem := range model.UserMfa {
			userMfaItemMap, err := DataSourceIBMIamEffectiveAccountSettingsEffectiveAccountSettingsUserMfaToMap(&userMfaItem)
			if err != nil {
				return modelMap, err
			}
			userMfa = append(userMfa, userMfaItemMap)
		}
		modelMap["user_mfa"] = userMfa
	}
	if model.History != nil {
		history := []map[string]interface{}{}
		for _, historyItem := range model.History {
			historyItemMap, err := DataSourceIBMIamEffectiveAccountSettingsEnityHistoryRecordToMap(&historyItem)
			if err != nil {
				return modelMap, err
			}
			history = append(history, historyItemMap)
		}
		modelMap["history"] = history
	}
	if model.SessionExpirationInSeconds != nil {
		modelMap["session_expiration_in_seconds"] = *model.SessionExpirationInSeconds
	}
	if model.SessionInvalidationInSeconds != nil {
		modelMap["session_invalidation_in_seconds"] = *model.SessionInvalidationInSeconds
	}
	if model.MaxSessionsPerIdentity != nil {
		modelMap["max_sessions_per_identity"] = *model.MaxSessionsPerIdentity
	}
	if model.SystemAccessTokenExpirationInSeconds != nil {
		modelMap["system_access_token_expiration_in_seconds"] = *model.SystemAccessTokenExpirationInSeconds
	}
	if model.SystemRefreshTokenExpirationInSeconds != nil {
		modelMap["system_refresh_token_expiration_in_seconds"] = *model.SystemRefreshTokenExpirationInSeconds
	}
	return modelMap, nil
}

func DataSourceIBMIamEffectiveAccountSettingsEnityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = *model.Timestamp
	modelMap["iam_id"] = *model.IamID
	modelMap["iam_id_account"] = *model.IamIDAccount
	modelMap["action"] = *model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = *model.Message
	return modelMap, nil
}

func DataSourceIBMIamEffectiveAccountSettingsAccountSettingsAssignedTemplatesSectionToMap(model *iamidentityv1.AccountSettingsAssignedTemplatesSection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TemplateID != nil {
		modelMap["template_id"] = *model.TemplateID
	}
	if model.TemplateVersion != nil {
		modelMap["template_version"] = flex.IntValue(model.TemplateVersion)
	}
	if model.TemplateName != nil {
		modelMap["template_name"] = *model.TemplateName
	}
	if model.RestrictCreateServiceID != nil {
		modelMap["restrict_create_service_id"] = *model.RestrictCreateServiceID
	}
	if model.RestrictCreatePlatformApikey != nil {
		modelMap["restrict_create_platform_apikey"] = *model.RestrictCreatePlatformApikey
	}
	if model.AllowedIPAddresses != nil {
		modelMap["allowed_ip_addresses"] = *model.AllowedIPAddresses
	}
	if model.Mfa != nil {
		modelMap["mfa"] = *model.Mfa
	}
	if model.UserMfa != nil {
		userMfa := []map[string]interface{}{}
		for _, userMfaItem := range model.UserMfa {
			userMfaItemMap, err := DataSourceIBMIamEffectiveAccountSettingsEffectiveAccountSettingsUserMfaToMap(&userMfaItem)
			if err != nil {
				return modelMap, err
			}
			userMfa = append(userMfa, userMfaItemMap)
		}
		modelMap["user_mfa"] = userMfa
	}
	if model.SessionExpirationInSeconds != nil {
		modelMap["session_expiration_in_seconds"] = *model.SessionExpirationInSeconds
	}
	if model.SessionInvalidationInSeconds != nil {
		modelMap["session_invalidation_in_seconds"] = *model.SessionInvalidationInSeconds
	}
	if model.MaxSessionsPerIdentity != nil {
		modelMap["max_sessions_per_identity"] = *model.MaxSessionsPerIdentity
	}
	if model.SystemAccessTokenExpirationInSeconds != nil {
		modelMap["system_access_token_expiration_in_seconds"] = *model.SystemAccessTokenExpirationInSeconds
	}
	if model.SystemRefreshTokenExpirationInSeconds != nil {
		modelMap["system_refresh_token_expiration_in_seconds"] = *model.SystemRefreshTokenExpirationInSeconds
	}
	return modelMap, nil
}
