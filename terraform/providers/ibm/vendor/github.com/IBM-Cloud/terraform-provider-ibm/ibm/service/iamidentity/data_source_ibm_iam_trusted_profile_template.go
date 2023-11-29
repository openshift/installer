// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMTrustedProfileTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMTrustedProfileTemplateRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile template.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version of the Profile Template.",
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
			"profile": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Input body parameters for the TemplateProfileComponent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Profile.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the Profile.",
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rules for the Profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the claim rule to be created or updated.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the claim rule.",
									},
									"realm_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The realm name of the Idp this claim rule applies to. This field is required only if the type is specified as 'Profile-SAML'.",
									},
									"expiration": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Session expiration in seconds, only required if type is 'Profile-SAML'.",
									},
									"conditions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Conditions of this claim rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"claim": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The claim to evaluate against.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operation to perform on the claim. valid values are EQUALS, NOT_EQUALS, EQUALS_IGNORE_CASE, NOT_EQUALS_IGNORE_CASE, CONTAINS, IN.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The stringified JSON value that the claim is compared to using the operator.",
												},
											},
										},
									},
								},
							},
						},
						"identities": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Identities for the Profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iam_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IAM ID of the identity.",
									},
									"identifier": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifier of the identity that can assume the trusted profiles. This can be a user identifier (IAM id), serviceid or crn. Internally it uses account id of the service id for the identifier 'serviceid' and for the identifier 'crn' it uses account id contained in the CRN.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the identity.",
									},
									"accounts": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Only valid for the type user. Accounts from which a user can assume the trusted profile.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the identity that can assume the trusted profile. This is optional field for all the types of identities. When this field is not set for the identity type 'serviceid' then the description of the service id is used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.",
									},
								},
							},
						},
					},
				},
			},
			"policy_template_references": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Existing policy templates that you can reference to assign access in the trusted profile component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of Access Policy Template.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of Access Policy Template.",
						},
					},
				},
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the trusted profile template.",
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
				Description: "Timestamp of when the template was created.",
			},
			"created_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the creator.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of when the template was last modified.",
			},
			"last_modified_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the identity that made the latest modification.",
			},
		},
	}
}

func dataSourceIBMTrustedProfileTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileTemplateVersionOptions := &iamidentityv1.GetProfileTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Get("template_id").(string))
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateRead failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateRead failed %s", err))
	}
	if version == "" {
		version = d.Get("version").(string)
	}

	getProfileTemplateVersionOptions.SetTemplateID(id)
	getProfileTemplateVersionOptions.SetVersion(version)

	if _, ok := d.GetOk("include_history"); ok {
		getProfileTemplateVersionOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}

	trustedProfileTemplateResponse, response, err := iamIdentityClient.GetProfileTemplateVersionWithContext(context, getProfileTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProfileTemplateVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProfileTemplateVersionWithContext failed %s\n%s", err, response))
	}

	d.SetId(buildResourceIdFromTemplateVersion(*trustedProfileTemplateResponse.ID, *trustedProfileTemplateResponse.Version))

	if err = d.Set("id", trustedProfileTemplateResponse.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	if !core.IsNil(trustedProfileTemplateResponse.Version) {
		versionStr := strconv.Itoa(int(*trustedProfileTemplateResponse.Version))
		if err = d.Set("version", versionStr); err != nil {
			return diag.FromErr(fmt.Errorf("error setting version: %s", err))
		}
	}

	if err = d.Set("account_id", trustedProfileTemplateResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}

	if err = d.Set("name", trustedProfileTemplateResponse.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("description", trustedProfileTemplateResponse.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("committed", trustedProfileTemplateResponse.Committed); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting committed: %s", err))
	}

	var profile []map[string]interface{}
	if trustedProfileTemplateResponse.Profile != nil {
		modelMap, err := dataSourceIBMTrustedProfileTemplateTemplateProfileComponentResponseToMap(trustedProfileTemplateResponse.Profile)
		if err != nil {
			return diag.FromErr(err)
		}
		profile = append(profile, modelMap)
	}
	if err = d.Set("profile", profile); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting profile %s", err))
	}

	var policyTemplateReferences []map[string]interface{}
	if trustedProfileTemplateResponse.PolicyTemplateReferences != nil {
		for _, modelItem := range trustedProfileTemplateResponse.PolicyTemplateReferences {
			modelMap, err := dataSourceIBMTrustedProfileTemplatePolicyTemplateReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			policyTemplateReferences = append(policyTemplateReferences, modelMap)
		}
	}
	if err = d.Set("policy_template_references", policyTemplateReferences); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting policy_template_references %s", err))
	}

	var history []map[string]interface{}
	if trustedProfileTemplateResponse.History != nil {
		for _, modelItem := range trustedProfileTemplateResponse.History {
			modelMap, err := dataSourceIBMTrustedProfileTemplateEnityHistoryRecordToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			history = append(history, modelMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting history %s", err))
	}

	if err = d.Set("entity_tag", trustedProfileTemplateResponse.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}

	if err = d.Set("crn", trustedProfileTemplateResponse.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("created_at", trustedProfileTemplateResponse.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("created_by_id", trustedProfileTemplateResponse.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}

	if err = d.Set("last_modified_at", trustedProfileTemplateResponse.LastModifiedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_at: %s", err))
	}

	if err = d.Set("last_modified_by_id", trustedProfileTemplateResponse.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_by_id: %s", err))
	}

	return nil
}

func dataSourceIBMTrustedProfileTemplateTemplateProfileComponentResponseToMap(model *iamidentityv1.TemplateProfileComponentResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Rules != nil {
		var rules []map[string]interface{}
		for _, rulesItem := range model.Rules {
			rulesItemMap, err := dataSourceIBMTrustedProfileTemplateTrustedProfileTemplateClaimRuleToMap(&rulesItem)
			if err != nil {
				return modelMap, err
			}
			rules = append(rules, rulesItemMap)
		}
		modelMap["rules"] = rules
	}
	if model.Identities != nil {
		var identities []map[string]interface{}
		for _, identitiesItem := range model.Identities {
			identitiesItemMap, err := dataSourceIBMTrustedProfileTemplateProfileIdentityResponseToMap(&identitiesItem)
			if err != nil {
				return modelMap, err
			}
			identities = append(identities, identitiesItemMap)
		}
		modelMap["identities"] = identities
	}
	return modelMap, nil
}

func dataSourceIBMTrustedProfileTemplateTrustedProfileTemplateClaimRuleToMap(model *iamidentityv1.TrustedProfileTemplateClaimRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	modelMap["type"] = model.Type
	if model.RealmName != nil {
		modelMap["realm_name"] = model.RealmName
	}
	if model.Expiration != nil {
		modelMap["expiration"] = flex.IntValue(model.Expiration)
	}
	var conditions []map[string]interface{}
	for _, conditionsItem := range model.Conditions {
		conditionsItemMap, err := dataSourceIBMTrustedProfileTemplateProfileClaimRuleConditionsToMap(&conditionsItem)
		if err != nil {
			return modelMap, err
		}
		conditions = append(conditions, conditionsItemMap)
	}
	modelMap["conditions"] = conditions
	return modelMap, nil
}

func dataSourceIBMTrustedProfileTemplateProfileClaimRuleConditionsToMap(model *iamidentityv1.ProfileClaimRuleConditions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["claim"] = model.Claim
	modelMap["operator"] = model.Operator
	modelMap["value"] = model.Value
	return modelMap, nil
}

func dataSourceIBMTrustedProfileTemplateProfileIdentityResponseToMap(model *iamidentityv1.ProfileIdentityResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["iam_id"] = model.IamID
	modelMap["identifier"] = model.Identifier
	modelMap["type"] = model.Type
	if model.Accounts != nil {
		modelMap["accounts"] = model.Accounts
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func dataSourceIBMTrustedProfileTemplatePolicyTemplateReferenceToMap(model *iamidentityv1.PolicyTemplateReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["version"] = model.Version
	return modelMap, nil
}

func dataSourceIBMTrustedProfileTemplateEnityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = model.Timestamp
	modelMap["iam_id"] = model.IamID
	modelMap["iam_id_account"] = model.IamIDAccount
	modelMap["action"] = model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = model.Message
	return modelMap, nil
}
