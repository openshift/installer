// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamaccessgroup

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
)

func ResourceIBMIAMAccessGroupTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIAMAccessGroupTemplateCreate,
		ReadContext:   resourceIBMIAMAccessGroupTemplateVersionRead,
		UpdateContext: resourceIBMIAMAccessGroupTemplateVersionUpdate,
		DeleteContext: resourceIBMIAMAccessGroupTemplateVersionDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"transaction_id": {
				Type:     schema.TypeString,
				Optional: true,

				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template", "transaction_id"),
				Description:  "An optional transaction id for the request.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template", "name"),
				Description:  "The name of the access group template.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,

				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template", "description"),
				Description:  "The description of the access group template.",
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,

				Description: "The ID of the account to which the access group template is assigned.",
			},
			"group": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access Group Component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Give the access group a unique name that doesn't conflict with other templates access group name in the given account. This is shown in child accounts.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access group description. This is shown in child accounts.",
						},
						"members": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Array of enterprise users to add to the template. All enterprise users that you add to the template must be invited to the child accounts where the template is assigned.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"users": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Array of enterprise users to add to the template. All enterprise users that you add to the template must be invited to the child accounts where the template is assigned.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"services": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Array of service IDs to add to the template.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"action_controls": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Control whether or not access group administrators in child accounts can add and remove members from the enterprise-managed access group in their account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"add": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Action control for adding child account members to an enterprise-managed access group. If an access group administrator in a child account adds a member, they can always remove them.",
												},
												"remove": {
													Type:        schema.TypeBool,
													Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Assertions Input Component.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Dynamic rules to automatically add federated users to access groups based on specific identity attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Dynamic rule name.",
												},
												"expiration": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Session duration in hours. Access group membership is revoked after this time period expires. Users must log back in to refresh their access group membership.",
												},
												"realm_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The identity provider (IdP) URL.",
												},
												"conditions": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Conditions of membership. You can think of this as a key:value pair.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"claim": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The key in the key:value pair.",
															},
															"operator": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Compares the claim and the value.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The value in the key:value pair.",
															},
														},
													},
												},
												"action_controls": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Control whether or not access group administrators in child accounts can update and remove this dynamic rule in the enterprise-managed access group in their account.This overrides outer level AssertionsActionControls.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"remove": {
																Type:        schema.TypeBool,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "Control whether or not access group administrators in child accounts can add, remove, and update dynamic rules for the enterprise-managed access group in their account. The inner level RuleActionControls override these action controls.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"add": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Action control for adding dynamic rules to an enterprise-managed access group. If an access group administrator in a child account adds a dynamic rule, they can always update or remove it.",
												},
												"remove": {
													Type:        schema.TypeBool,
													Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Access group action controls component.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Control whether or not access group administrators in child accounts can add access policies to the enterprise-managed access group in their account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"add": {
													Type:        schema.TypeBool,
													Optional:    true,
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
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "References to policy templates assigned to the access group template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Policy template ID.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Policy template version.",
						},
					},
				},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the access group template.",
			},
			"committed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A boolean indicating whether the access group template is committed. You must commit a template before you can assign it to child accounts.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the access group template resource.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the access group template was created.",
			},
			"created_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user who created the access group template.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the access group template was last modified.",
			},
			"last_modified_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user who last modified the access group template.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template ID.",
			},
		},
	}
}

func ResourceIBMIAMAccessGroupTemplateValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "transaction_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9_-]+$`,
			MinValueLength:             1,
			MaxValueLength:             50,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9!@#$%^&*()_+{}:;"'<>,.?\/|\-\s]+$`,
			MinValueLength:             1,
			MaxValueLength:             100,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9!@#$%^&*()_+{}:;"'<>,.?\/|\-\s]+$`,
			MinValueLength:             0,
			MaxValueLength:             250,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_iam_access_group_template", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIAMAccessGroupTemplateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTemplateOptions := &iamaccessgroupsv2.CreateTemplateOptions{}

	createTemplateOptions.SetName(d.Get("name").(string))

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	createTemplateOptions.SetAccountID(userDetails.UserAccount)

	if _, ok := d.GetOk("description"); ok {
		createTemplateOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("group"); ok {
		groupModel, err := resourceIBMIAMAccessGroupTemplateMapToAccessGroupRequest(d.Get("group.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTemplateOptions.SetGroup(groupModel)
	}
	if _, ok := d.GetOk("policy_template_references"); ok {
		var policyTemplateReferences []iamaccessgroupsv2.PolicyTemplates
		for _, v := range d.Get("policy_template_references").(*schema.Set).List() {
			value := v.(map[string]interface{})
			policyTemplateReferencesItem, err := resourceIBMIAMAccessGroupTemplateMapToPolicyTemplates(value)
			if err != nil {
				return diag.FromErr(err)
			}
			policyTemplateReferences = append(policyTemplateReferences, *policyTemplateReferencesItem)
		}
		createTemplateOptions.SetPolicyTemplateReferences(policyTemplateReferences)
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		createTemplateOptions.SetTransactionID(d.Get("transaction_id").(string))
	}

	templateResponse, response, err := iamAccessGroupsClient.CreateTemplateWithContext(context, createTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTemplateWithContext failed %s\n%s", err, response))
	}
	version, _ := strconv.Atoi(*templateResponse.Version)

	if d.Get("committed").(bool) {
		commitTemplateOptions := &iamaccessgroupsv2.CommitTemplateOptions{}
		commitTemplateOptions.SetTemplateID(*templateResponse.ID)
		commitTemplateOptions.SetVersionNum(*templateResponse.Version)
		commitTemplateOptions.SetIfMatch(response.Headers.Get("ETag"))
		response, err = iamAccessGroupsClient.CommitTemplateWithContext(context, commitTemplateOptions)
		if err != nil {
			log.Printf("[DEBUG] CommitTemplateWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("CommitTemplateWithContext failed %s\n%s", err, response))
		}
	}
	d.SetId(fmt.Sprintf("%s/%d", *templateResponse.ID, version))

	return resourceIBMIAMAccessGroupTemplateVersionRead(context, d, meta)
}

func resourceIBMIAMAccessGroupTemplateMapToAccessGroupRequest(modelMap map[string]interface{}) (*iamaccessgroupsv2.AccessGroupRequest, error) {
	model := &iamaccessgroupsv2.AccessGroupRequest{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["members"] != nil && len(modelMap["members"].([]interface{})) > 0 {
		MembersModel, err := resourceIBMIAMAccessGroupTemplateMapToMembers(modelMap["members"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Members = MembersModel
	}
	if modelMap["assertions"] != nil && len(modelMap["assertions"].([]interface{})) > 0 {
		AssertionsModel, err := resourceIBMIAMAccessGroupTemplateMapToAssertions(modelMap["assertions"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Assertions = AssertionsModel
	}
	if modelMap["action_controls"] != nil && len(modelMap["action_controls"].([]interface{})) > 0 {
		ActionControlsModel, err := resourceIBMIAMAccessGroupTemplateMapToGroupActionControls(modelMap["action_controls"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ActionControls = ActionControlsModel
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToMembers(modelMap map[string]interface{}) (*iamaccessgroupsv2.Members, error) {
	model := &iamaccessgroupsv2.Members{}
	if modelMap["users"] != nil {
		users := []string{}
		for _, usersItem := range modelMap["users"].([]interface{}) {
			users = append(users, usersItem.(string))
		}
		model.Users = users
	}
	if modelMap["services"] != nil {
		services := []string{}
		for _, servicesItem := range modelMap["services"].([]interface{}) {
			services = append(services, servicesItem.(string))
		}
		model.Services = services
	}
	if modelMap["action_controls"] != nil && len(modelMap["action_controls"].([]interface{})) > 0 {
		ActionControlsModel, err := resourceIBMIAMAccessGroupTemplateMapToMembersActionControls(modelMap["action_controls"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ActionControls = ActionControlsModel
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToMembersActionControls(modelMap map[string]interface{}) (*iamaccessgroupsv2.MembersActionControls, error) {
	model := &iamaccessgroupsv2.MembersActionControls{}
	if modelMap["add"] != nil {
		model.Add = core.BoolPtr(modelMap["add"].(bool))
	}
	if modelMap["remove"] != nil {
		model.Remove = core.BoolPtr(modelMap["remove"].(bool))
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToAssertions(modelMap map[string]interface{}) (*iamaccessgroupsv2.Assertions, error) {
	model := &iamaccessgroupsv2.Assertions{}
	if modelMap["rules"] != nil {
		rules := []iamaccessgroupsv2.AssertionsRule{}
		for _, rulesItem := range modelMap["rules"].([]interface{}) {
			rulesItemModel, err := resourceIBMIAMAccessGroupTemplateMapToAssertionsRule(rulesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			rules = append(rules, *rulesItemModel)
		}
		model.Rules = rules
	}
	if modelMap["action_controls"] != nil && len(modelMap["action_controls"].([]interface{})) > 0 {
		ActionControlsModel, err := resourceIBMIAMAccessGroupTemplateMapToAssertionsActionControls(modelMap["action_controls"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ActionControls = ActionControlsModel
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToAssertionsRule(modelMap map[string]interface{}) (*iamaccessgroupsv2.AssertionsRule, error) {
	model := &iamaccessgroupsv2.AssertionsRule{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["expiration"] != nil {
		model.Expiration = core.Int64Ptr(int64(modelMap["expiration"].(int)))
	}
	if modelMap["realm_name"] != nil && modelMap["realm_name"].(string) != "" {
		model.RealmName = core.StringPtr(modelMap["realm_name"].(string))
	}
	if modelMap["conditions"] != nil {
		conditions := []iamaccessgroupsv2.Conditions{}
		for _, conditionsItem := range modelMap["conditions"].([]interface{}) {
			conditionsItemModel, err := resourceIBMIAMAccessGroupTemplateMapToConditions(conditionsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			conditions = append(conditions, *conditionsItemModel)
		}
		model.Conditions = conditions
	}
	if modelMap["action_controls"] != nil && len(modelMap["action_controls"].([]interface{})) > 0 {
		ActionControlsModel, err := resourceIBMIAMAccessGroupTemplateMapToRuleActionControls(modelMap["action_controls"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ActionControls = ActionControlsModel
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToConditions(modelMap map[string]interface{}) (*iamaccessgroupsv2.Conditions, error) {
	model := &iamaccessgroupsv2.Conditions{}
	if modelMap["claim"] != nil && modelMap["claim"].(string) != "" {
		model.Claim = core.StringPtr(modelMap["claim"].(string))
	}
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToRuleActionControls(modelMap map[string]interface{}) (*iamaccessgroupsv2.RuleActionControls, error) {
	model := &iamaccessgroupsv2.RuleActionControls{}
	if modelMap["remove"] != nil {
		model.Remove = core.BoolPtr(modelMap["remove"].(bool))
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToAssertionsActionControls(modelMap map[string]interface{}) (*iamaccessgroupsv2.AssertionsActionControls, error) {
	model := &iamaccessgroupsv2.AssertionsActionControls{}
	if modelMap["add"] != nil {
		model.Add = core.BoolPtr(modelMap["add"].(bool))
	}
	if modelMap["remove"] != nil {
		model.Remove = core.BoolPtr(modelMap["remove"].(bool))
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToGroupActionControls(modelMap map[string]interface{}) (*iamaccessgroupsv2.GroupActionControls, error) {
	model := &iamaccessgroupsv2.GroupActionControls{}
	if modelMap["access"] != nil && len(modelMap["access"].([]interface{})) > 0 {
		AccessModel, err := resourceIBMIAMAccessGroupTemplateMapToAccessActionControls(modelMap["access"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Access = AccessModel
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToAccessActionControls(modelMap map[string]interface{}) (*iamaccessgroupsv2.AccessActionControls, error) {
	model := &iamaccessgroupsv2.AccessActionControls{}
	if modelMap["add"] != nil {
		model.Add = core.BoolPtr(modelMap["add"].(bool))
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateMapToPolicyTemplates(modelMap map[string]interface{}) (*iamaccessgroupsv2.PolicyTemplates, error) {
	model := &iamaccessgroupsv2.PolicyTemplates{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	return model, nil
}

func resourceIBMIAMAccessGroupTemplateAccessGroupResponseToMap(model *iamaccessgroupsv2.AccessGroupResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Members != nil {
		membersMap, err := resourceIBMIAMAccessGroupTemplateMembersToMap(model.Members)
		if err != nil {
			return modelMap, err
		}
		modelMap["members"] = []map[string]interface{}{membersMap}
	}
	if model.Assertions != nil {
		assertionsMap, err := resourceIBMIAMAccessGroupTemplateAssertionsToMap(model.Assertions)
		if err != nil {
			return modelMap, err
		}
		modelMap["assertions"] = []map[string]interface{}{assertionsMap}
	}
	if model.ActionControls != nil {
		actionControlsMap, err := resourceIBMIAMAccessGroupTemplateGroupActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateMembersToMap(model *iamaccessgroupsv2.Members) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Users != nil {
		modelMap["users"] = model.Users
	}
	if model.Services != nil {
		modelMap["services"] = model.Services
	}
	if model.ActionControls != nil {
		actionControlsMap, err := resourceIBMIAMAccessGroupTemplateMembersActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateMembersActionControlsToMap(model *iamaccessgroupsv2.MembersActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Add != nil {
		modelMap["add"] = model.Add
	}
	if model.Remove != nil {
		modelMap["remove"] = model.Remove
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateAssertionsToMap(model *iamaccessgroupsv2.Assertions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Rules != nil {
		rules := []map[string]interface{}{}
		for _, rulesItem := range model.Rules {
			rulesItemMap, err := resourceIBMIAMAccessGroupTemplateAssertionsRuleToMap(&rulesItem)
			if err != nil {
				return modelMap, err
			}
			rules = append(rules, rulesItemMap)
		}
		modelMap["rules"] = rules
	}
	if model.ActionControls != nil {
		actionControlsMap, err := resourceIBMIAMAccessGroupTemplateAssertionsActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateAssertionsRuleToMap(model *iamaccessgroupsv2.AssertionsRule) (map[string]interface{}, error) {
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
			conditionsItemMap, err := resourceIBMIAMAccessGroupTemplateConditionsToMap(&conditionsItem)
			if err != nil {
				return modelMap, err
			}
			conditions = append(conditions, conditionsItemMap)
		}
		modelMap["conditions"] = conditions
	}
	if model.ActionControls != nil {
		actionControlsMap, err := resourceIBMIAMAccessGroupTemplateRuleActionControlsToMap(model.ActionControls)
		if err != nil {
			return modelMap, err
		}
		modelMap["action_controls"] = []map[string]interface{}{actionControlsMap}
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateConditionsToMap(model *iamaccessgroupsv2.Conditions) (map[string]interface{}, error) {
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

func resourceIBMIAMAccessGroupTemplateRuleActionControlsToMap(model *iamaccessgroupsv2.RuleActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Remove != nil {
		modelMap["remove"] = model.Remove
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateAssertionsActionControlsToMap(model *iamaccessgroupsv2.AssertionsActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Add != nil {
		modelMap["add"] = model.Add
	}
	if model.Remove != nil {
		modelMap["remove"] = model.Remove
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateGroupActionControlsToMap(model *iamaccessgroupsv2.GroupActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Access != nil {
		accessMap, err := resourceIBMIAMAccessGroupTemplateAccessActionControlsToMap(model.Access)
		if err != nil {
			return modelMap, err
		}
		modelMap["access"] = []map[string]interface{}{accessMap}
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplateAccessActionControlsToMap(model *iamaccessgroupsv2.AccessActionControls) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Add != nil {
		modelMap["add"] = model.Add
	}
	return modelMap, nil
}

func resourceIBMIAMAccessGroupTemplatePolicyTemplatesToMap(model *iamaccessgroupsv2.PolicyTemplates) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	return modelMap, nil
}
