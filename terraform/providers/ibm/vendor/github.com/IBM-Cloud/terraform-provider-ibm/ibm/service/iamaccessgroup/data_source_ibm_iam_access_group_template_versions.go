// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamaccessgroup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
)

func DataSourceIBMIAMAccessGroupTemplateVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIAMAccessGroupTemplateVersionRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the template that you want to list all versions of.",
			},
			"first": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A string containing the link’s URL.",
						},
					},
				},
			},
			"previous": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A string containing the link’s URL.",
						},
					},
				},
			},
			"last": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A string containing the link’s URL.",
						},
					},
				},
			},
			"group_template_versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of access group template versions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the template.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the template.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the account associated with the template.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version number of the template.",
						},
						"committed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "A boolean indicating whether the template is committed or not.",
						},
						"group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Access Group Component.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Give the access group a unique name that doesn't conflict with other templates access group name in the given account. This is shown in child accounts.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Access group description. This is shown in child accounts.",
									},
									"members": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Array of enterprise users to add to the template. All enterprise users that you add to the template must be invited to the child accounts where the template is assigned.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"users": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Array of enterprise users to add to the template. All enterprise users that you add to the template must be invited to the child accounts where the template is assigned.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"services": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Array of service IDs to add to the template.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"action_controls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Control whether or not access group administrators in child accounts can add and remove members from the enterprise-managed access group in their account.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"add": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Action control for adding child account members to an enterprise-managed access group. If an access group administrator in a child account adds a member, they can always remove them.",
															},
															"remove": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Action control for removing enterprise-managed members from an enterprise-managed access group.",
															},
														},
													},
												},
											},
										},
									},
									"assertions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Assertions Input Component.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Dynamic rules to automatically add federated users to access groups based on specific identity attributes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Dynamic rule name.",
															},
															"expiration": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Session duration in hours. Access group membership is revoked after this time period expires. Users must log back in to refresh their access group membership.",
															},
															"realm_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The identity provider (IdP) URL.",
															},
															"conditions": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Conditions of membership. You can think of this as a key:value pair.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"claim": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The key in the key:value pair.",
																		},
																		"operator": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Compares the claim and the value.",
																		},
																		"value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The value in the key:value pair.",
																		},
																	},
																},
															},
															"action_controls": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Control whether or not access group administrators in child accounts can update and remove this dynamic rule in the enterprise-managed access group in their account.This overrides outer level AssertionsActionControls.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"remove": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Action control for removing this enterprise-managed dynamic rule.",
																		},
																	},
																},
															},
														},
													},
												},
												"action_controls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Control whether or not access group administrators in child accounts can add and remove dynamic rules for the enterprise-managed access group in their account. The inner level RuleActionControls override these action controls.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"add": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Action control for adding dynamic rules to an enterprise-managed access group. If an access group administrator in a child account adds a dynamic rule, they can always update or remove it.",
															},
															"remove": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Action control for removing enterprise-managed dynamic rules in an enterprise-managed access group.",
															},
														},
													},
												},
											},
										},
									},
									"action_controls": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Access group action controls component.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"access": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Control whether or not access group administrators in child accounts can add access policies to the enterprise-managed access group in their account.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"add": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Action control for adding access policies to an enterprise-managed access group in a child account. If an access group administrator in a child account adds a policy, they can always update or remove it.",
															},
														},
													},
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
							Description: "A list of policy templates associated with the template.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy template ID.",
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy template version.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to the template resource.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time the template was created.",
						},
						"created_by_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user who created the template.",
						},
						"last_modified_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time the template was last modified.",
						},
						"last_modified_by_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user who last modified the template.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIAMAccessGroupTemplateVersionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listTemplateVersionsOptions := &iamaccessgroupsv2.ListTemplateVersionsOptions{}

	listTemplateVersionsOptions.SetTemplateID(d.Get("template_id").(string))

	var pager *iamaccessgroupsv2.TemplateVersionsPager
	pager, err = iamAccessGroupsClient.NewTemplateVersionsPager(listTemplateVersionsOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] TemplateVersionsPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("TemplateVersionsPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIBMIAMAccessGroupTemplateVersionID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIAMAccessGroupTemplateVersionListTemplateVersionResponseToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("group_template_versions", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting group_template_versions %s", err))
	}

	return nil
}

// dataSourceIBMIAMAccessGroupTemplateVersionID returns a reasonable ID for the list.
func dataSourceIBMIAMAccessGroupTemplateVersionID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIAMAccessGroupTemplateVersionHrefStructToMap(model *iamaccessgroupsv2.HrefStruct) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionListTemplateVersionResponseToMap(model *iamaccessgroupsv2.ListTemplateVersionResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["description"] = model.Description
	modelMap["account_id"] = model.AccountID
	modelMap["version"] = model.Version
	modelMap["committed"] = model.Committed
	groupMap, err := dataSourceIBMIAMAccessGroupTemplateVersionAccessGroupResponseToMap(model.Group)
	if err != nil {
		return modelMap, err
	}
	modelMap["group"] = []map[string]interface{}{groupMap}
	policyTemplateReferences := []map[string]interface{}{}
	for _, policyTemplateReferencesItem := range model.PolicyTemplateReferences {
		policyTemplateReferencesItemMap, err := dataSourceIBMIAMAccessGroupTemplateVersionPolicyTemplatesToMap(&policyTemplateReferencesItem)
		if err != nil {
			return modelMap, err
		}
		policyTemplateReferences = append(policyTemplateReferences, policyTemplateReferencesItemMap)
	}
	modelMap["policy_template_references"] = policyTemplateReferences
	modelMap["href"] = model.Href
	modelMap["created_at"] = model.CreatedAt
	modelMap["created_by_id"] = model.CreatedByID
	modelMap["last_modified_at"] = model.LastModifiedAt
	modelMap["last_modified_by_id"] = model.LastModifiedByID
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionAccessGroupResponseToMap(model *iamaccessgroupsv2.AccessGroupResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Members != nil {
		membersMap, err := dataSourceIBMIAMAccessGroupTemplateVersionMembersToMap(model.Members)
		if err != nil {
			return modelMap, err
		}
		modelMap["members"] = []map[string]interface{}{membersMap}
	}
	if model.Assertions != nil {
		assertionsMap, err := dataSourceIBMIAMAccessGroupTemplateVersionAssertionsToMap(model.Assertions)
		if err != nil {
			return modelMap, err
		}
		modelMap["assertions"] = []map[string]interface{}{assertionsMap}
	}
	if model.ActionControls != nil {
		actionControlsMap, err := dataSourceIBMIAMAccessGroupTemplateVersionGroupActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionMembersToMap(model *iamaccessgroupsv2.Members) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Users != nil {
		modelMap["users"] = model.Users
	}
	if model.Services != nil {
		modelMap["services"] = model.Services
	}
	if model.ActionControls != nil {
		actionControlsMap, err := dataSourceIBMIAMAccessGroupTemplateVersionMembersActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionMembersActionControlsToMap(model *iamaccessgroupsv2.MembersActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Add != nil {
		modelMap["add"] = model.Add
	}
	if model.Remove != nil {
		modelMap["remove"] = model.Remove
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionAssertionsToMap(model *iamaccessgroupsv2.Assertions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Rules != nil {
		rules := []map[string]interface{}{}
		for _, rulesItem := range model.Rules {
			rulesItemMap, err := dataSourceIBMIAMAccessGroupTemplateVersionAssertionsRuleToMap(&rulesItem)
			if err != nil {
				return modelMap, err
			}
			rules = append(rules, rulesItemMap)
		}
		modelMap["rules"] = rules
	}
	if model.ActionControls != nil {
		actionControlsMap, err := dataSourceIBMIAMAccessGroupTemplateVersionAssertionsActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionAssertionsRuleToMap(model *iamaccessgroupsv2.AssertionsRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Expiration != nil {
		modelMap["expiration"] = flex.IntValue(model.Expiration)
	}
	if model.RealmName != nil {
		modelMap["realm_name"] = model.RealmName
	}
	if model.Conditions != nil {
		conditions := []map[string]interface{}{}
		for _, conditionsItem := range model.Conditions {
			conditionsItemMap, err := dataSourceIBMIAMAccessGroupTemplateVersionConditionsToMap(&conditionsItem)
			if err != nil {
				return modelMap, err
			}
			conditions = append(conditions, conditionsItemMap)
		}
		modelMap["conditions"] = conditions
	}
	if model.ActionControls != nil {
		actionControlsMap, err := dataSourceIBMIAMAccessGroupTemplateVersionRuleActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionConditionsToMap(model *iamaccessgroupsv2.Conditions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Claim != nil {
		modelMap["claim"] = model.Claim
	}
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionRuleActionControlsToMap(model *iamaccessgroupsv2.RuleActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Remove != nil {
		modelMap["remove"] = model.Remove
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionAssertionsActionControlsToMap(model *iamaccessgroupsv2.AssertionsActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Add != nil {
		modelMap["add"] = model.Add
	}
	if model.Remove != nil {
		modelMap["remove"] = model.Remove
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionGroupActionControlsToMap(model *iamaccessgroupsv2.GroupActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Access != nil {
		accessMap, err := dataSourceIBMIAMAccessGroupTemplateVersionAccessActionControlsToMap(model.Access)
		if err != nil {
			return modelMap, err
		}
		modelMap["access"] = []map[string]interface{}{accessMap}
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionAccessActionControlsToMap(model *iamaccessgroupsv2.AccessActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Add != nil {
		modelMap["add"] = model.Add
	}
	return modelMap, nil
}

func dataSourceIBMIAMAccessGroupTemplateVersionPolicyTemplatesToMap(model *iamaccessgroupsv2.PolicyTemplates) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	return modelMap, nil
}
