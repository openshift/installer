// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

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
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
)

func ResourceIBMIAMPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIAMPolicyTemplateCreate,
		ReadContext:   resourceIBMIAMPolicyTemplateVersionRead,
		UpdateContext: resourceIBMIAMPolicyTemplateVersionUpdate,
		DeleteContext: resourceIBMIAMPolicyTemplateVersionDelete,
		Exists:        resourceIBMIAMPolicyTemplateVersionExists,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_policy_template", "name"),
				Description:  "name of template.",
			},
			"policy": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The core set of properties associated with the template's policy objet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The policy type; either 'access' or 'authorization'.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Allows the customer to use their own words to record the purpose/context related to a policy.",
						},
						"resource": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The resource attributes to which the policy grants access.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attributes": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of resource attributes to which the policy grants access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of a resource attribute.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator of an attribute.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.",
												},
											},
										},
									},
									"tags": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Optional list of resource tags to which the policy grants access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of an access management tag.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The value of an access management tag.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator of an access management tag.",
												},
											},
										},
									},
								},
							},
						},
						"pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates pattern of rule, either 'time-based-conditions:once', 'time-based-conditions:weekly:all-day', or 'time-based-conditions:weekly:custom-hours'.",
						},
						"rule_conditions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Rule conditions enforced by the policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key of the condition",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Operator of the condition",
									},
									"value": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Value of the condition",
									},
									"conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Additional Rule conditions enforced by the policy",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Key of the condition",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Operator of the condition",
												},
												"value": {
													Type:        schema.TypeList,
													Required:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Value of the condition",
												},
											},
										},
									},
								},
							},
						},

						"rule_operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator that multiple rule conditions are evaluated over",
						},
						"roles": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Role names of the policy definition",
						},
						"subject": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "The subject attributes for authorization type templates",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attributes": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of resource attributes to which the policy grants access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of a resource attribute.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator of an attribute.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.",
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
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_policy_template", "description"),
				Description:  "description of template purpose.",
			},
			"committed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "committed status for the template.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template ID.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template Version.",
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIAMPolicyTemplateValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Required:                   true,
			MinValueLength:             1,
			MaxValueLength:             300,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^.*$`,
			MinValueLength:             0,
			MaxValueLength:             300,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_iam_policy_template", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIAMPolicyTemplateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch BluemixUserDetails %s", err))
	}

	accountID := userDetails.UserAccount

	createPolicyTemplateOptions := &iampolicymanagementv1.CreatePolicyTemplateOptions{}

	createPolicyTemplateOptions.SetName(d.Get("name").(string))
	createPolicyTemplateOptions.SetAccountID(accountID)

	policyModel, err := generateTemplatePolicy(d, iamPolicyManagementClient)
	if err != nil {
		return diag.FromErr(err)
	}
	createPolicyTemplateOptions.SetPolicy(policyModel)
	if _, ok := d.GetOk("description"); ok {
		createPolicyTemplateOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("committed"); ok {
		createPolicyTemplateOptions.SetCommitted(d.Get("committed").(bool))
	}

	policyTemplate, response, err := iamPolicyManagementClient.CreatePolicyTemplateWithContext(context, createPolicyTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] CreatePolicyTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreatePolicyTemplateWithContext failed %s\n%s", err, response))
	}

	version, _ := strconv.Atoi(*policyTemplate.Version)
	d.SetId(fmt.Sprintf("%s/%d", *policyTemplate.ID, version))
	return resourceIBMIAMPolicyTemplateVersionRead(context, d, meta)
}

func generateTemplatePolicy(d *schema.ResourceData, iamPolicyManagementClient *iampolicymanagementv1.IamPolicyManagementV1) (*iampolicymanagementv1.TemplatePolicy, error) {
	var sourceServiceName, targetServiceName, serviceGroupID string
	model := &iampolicymanagementv1.TemplatePolicy{}
	modelMap := d.Get("policy.0").(map[string]interface{})
	model.Type = core.StringPtr(modelMap["type"].(string))

	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	var roleList *iampolicymanagementv1.RoleCollection
	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{}
	var err error
	if _, ok := d.GetOk("policy.0.resource"); ok {
		modelMap := modelMap["resource"].([]interface{})[0].(map[string]interface{})
		modelResource := &iampolicymanagementv1.V2PolicyResource{}
		attributes := []iampolicymanagementv1.V2PolicyResourceAttribute{}
		for _, attributesItem := range modelMap["attributes"].([]interface{}) {
			attributesItemModel := &iampolicymanagementv1.V2PolicyResourceAttribute{}
			attributesItemModel.Key = core.StringPtr(attributesItem.(map[string]interface{})["key"].(string))
			attributesItemModel.Operator = core.StringPtr(attributesItem.(map[string]interface{})["operator"].(string))
			attributesItemModel.Value = attributesItem.(map[string]interface{})["value"].(string)

			if *attributesItemModel.Key == "serviceName" &&
				(*attributesItemModel.Operator == "stringMatch" ||
					*attributesItemModel.Operator == "stringEquals") {
				targetServiceName = fmt.Sprintf("%v", attributesItemModel.Value)
			}

			if *attributesItemModel.Key == "service_group_id" && (*attributesItemModel.Operator == "stringMatch" ||
				*attributesItemModel.Operator == "stringEquals") {
				serviceGroupID = fmt.Sprintf("%v", attributesItemModel.Value)
			}

			if *attributesItemModel.Key == "serviceType" && attributesItemModel.Value.(string) == "service" && (*attributesItemModel.Operator == "stringMatch" ||
				*attributesItemModel.Operator == "stringEquals") {
				listRoleOptions.ServiceName = core.StringPtr("alliamserviceroles")
			}

			if *model.Type == "authorization" && *attributesItemModel.Key == "resourceType" && targetServiceName == "" {
				targetServiceName = "resource-controller"
			}

			if *attributesItemModel.Operator == "stringExists" {
				if attributesItemModel.Value == "true" {
					attributesItemModel.Value = true
				} else if attributesItemModel.Value == "false" {
					attributesItemModel.Value = false
				} else {
					return model, fmt.Errorf("[ERROR] When operator equals stringExists, value should be either \"true\" or \"false\", instead of %s",
						attributesItemModel.Value)
				}
			}

			attributes = append(attributes, *attributesItemModel)
		}
		modelResource.Attributes = attributes
		if modelMap["tags"] != nil {
			tags := []iampolicymanagementv1.V2PolicyResourceTag{}
			for _, tagsItem := range modelMap["tags"].([]interface{}) {
				tagsItemModel, err := generateTemplatePolicyTag(tagsItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				tags = append(tags, *tagsItemModel)
			}
			modelResource.Tags = tags
		}
		model.Resource = modelResource
	}

	// check subject only for authorization type
	if _, ok := d.GetOk("policy.0.subject"); ok {
		modelMap := (modelMap["subject"]).(*schema.Set).List()
		modelSubject := &iampolicymanagementv1.V2PolicySubject{}
		attributes := []iampolicymanagementv1.V2PolicySubjectAttribute{}
		for _, attributesItem := range modelMap {
			attribute := (attributesItem.(map[string]interface{}))["attributes"]
			for _, item := range (attribute).([]interface{}) {
				attributesItemModel := &iampolicymanagementv1.V2PolicySubjectAttribute{}
				attributesItemModel.Key = core.StringPtr((item.((map[string]interface{}))["key"].(string)))
				attributesItemModel.Operator = core.StringPtr(item.(map[string]interface{})["operator"].(string))
				attributesItemModel.Value = core.StringPtr(item.((map[string]interface{}))["value"].(string))
				if *attributesItemModel.Key == "serviceName" {
					sourceServiceName = item.((map[string]interface{}))["value"].(string)
				}
				if *attributesItemModel.Operator == "stringExists" {
					if attributesItemModel.Value == "true" {
						attributesItemModel.Value = true
					} else if attributesItemModel.Value == "false" {
						attributesItemModel.Value = false
					} else {
						return model, fmt.Errorf("[ERROR] Only values \"true\" and \"false\" are allowed when operator is \"stringExists\". Received %s.", attributesItemModel.Value)
					}
				}
				if *model.Type == "authorization" && *attributesItemModel.Operator == "" && attributesItemModel.Value == "*" && *attributesItemModel.Key == "resourceGroupId" {
					attributesItemModel.Value = true
					*attributesItemModel.Operator = "stringExists"
				}
				attributes = append(attributes, *attributesItemModel)
			}
		}
		modelSubject.Attributes = attributes
		model.Subject = modelSubject
	}

	if targetServiceName != "" {
		listRoleOptions.ServiceName = &targetServiceName
	}
	if serviceGroupID != "" {
		listRoleOptions.ServiceGroupID = &serviceGroupID
	}
	if sourceServiceName != "" {
		listRoleOptions.SourceServiceName = &sourceServiceName
		listRoleOptions.PolicyType = core.StringPtr("authorization")
	}
	roles, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
	roleList = roles
	if err != nil {
		return model, err
	}

	if _, ok := d.GetOk("policy.0.roles"); ok && roleList != nil {
		controlModel, err := generateTemplatePolicyControl(modelMap["roles"].([]interface{}), roleList)
		if err != nil {
			return nil, err
		}

		model.Control = controlModel
	}

	if modelMap["pattern"] != nil && modelMap["pattern"].(string) != "" {
		model.Pattern = core.StringPtr(modelMap["pattern"].(string))
	}

	if modelMap["rule_conditions"] != nil && len(modelMap["rule_conditions"].(*schema.Set).List()) > 0 {
		conditions := []iampolicymanagementv1.NestedConditionIntf{}
		for _, condition := range modelMap["rule_conditions"].(*schema.Set).List() {
			c := condition.(map[string]interface{})
			key := c["key"].(string)
			operator := c["operator"].(string)
			r := &iampolicymanagementv1.NestedCondition{
				Key:      &key,
				Operator: &operator,
			}

			interfaceValues := c["value"].([]interface{})
			values := make([]string, len(interfaceValues))
			for i, v := range interfaceValues {
				values[i] = fmt.Sprint(v)
			}

			if len(values) > 1 {
				r.Value = &values
			} else if operator == "stringExists" && values[0] == "true" {
				r.Value = true
			} else {
				r.Value = &values[0]
			}

			conditions = append(conditions, r)
		}
		rule := new(iampolicymanagementv1.V2PolicyRule)
		if len(conditions) == 1 {
			ruleCondition := conditions[0].(*iampolicymanagementv1.NestedCondition)
			rule.Key = ruleCondition.Key
			rule.Operator = ruleCondition.Operator
			rule.Value = ruleCondition.Value
		} else {
			ruleOperator := modelMap["rule_operator"].(string)
			rule.Operator = &ruleOperator
			rule.Conditions = conditions
		}
		model.Rule = rule
	}
	return model, nil
}

func generateTemplatePolicyTag(modelMap map[string]interface{}) (*iampolicymanagementv1.V2PolicyResourceTag, error) {
	model := &iampolicymanagementv1.V2PolicyResourceTag{}
	model.Key = core.StringPtr(modelMap["key"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	return model, nil
}

func generateTemplatePolicyControl(roles []interface{}, roleList *iampolicymanagementv1.RoleCollection) (*iampolicymanagementv1.Control, error) {
	policyRoles := flex.MapRoleListToPolicyRoles(*roleList)

	policyRoles, err := flex.GetRolesFromRoleNames(flex.ExpandStringList(roles), policyRoles)
	if err != nil {
		return &iampolicymanagementv1.Control{}, err
	}
	policyGrant := &iampolicymanagementv1.Grant{
		Roles: flex.MapPolicyRolesToRoles(policyRoles),
	}
	policyControl := &iampolicymanagementv1.Control{
		Grant: policyGrant,
	}
	return policyControl, nil
}

func flattenTemplatePolicy(model *iampolicymanagementv1.TemplatePolicy, iamPolicyManagementClient *iampolicymanagementv1.IamPolicyManagementV1) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	var subjectMap map[string]interface{}
	var err error
	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{}
	// Check subject details exists
	if model.Subject != nil {
		subjectMap, listRoleOptions, err = flattenTemplatePolicySubject(model.Subject, listRoleOptions)

		if err != nil {
			return nil, err
		}
		modelMap["subject"] = []map[string]interface{}{subjectMap}
	}
	if *model.Type == "authorization" {
		listRoleOptions.SetPolicyType("authorization")
	}
	if model.Resource != nil {
		resourceMap, roleList, err := flattenTemplatePolicyResource(model.Resource, listRoleOptions, iamPolicyManagementClient)
		if err != nil {
			return nil, err
		}
		controlResponse := model.Control
		policyRoles := flex.MapRolesToPolicyRoles(controlResponse.Grant.Roles)

		rolesWithCrn := flex.MapRoleListToPolicyRoles(*roleList)
		roleNames := []string{}
		for _, role := range policyRoles {
			role, err := flex.FindRoleByCRN(rolesWithCrn, *role.RoleID)
			if err != nil {
				return nil, err
			}
			roleNames = append(roleNames, *role.DisplayName)
		}
		modelMap["resource"] = []map[string]interface{}{resourceMap}
		modelMap["roles"] = roleNames
	}

	if model.Pattern != nil {
		modelMap["pattern"] = model.Pattern
	}
	if model.Rule != nil {
		modelMap["rule_conditions"] = flex.FlattenRuleConditions(*model.Rule.(*iampolicymanagementv1.V2PolicyRule))
		if len(model.Rule.(*iampolicymanagementv1.V2PolicyRule).Conditions) > 0 {
			modelMap["rule_operator"] = model.Rule.(*iampolicymanagementv1.V2PolicyRule).Operator
		}
	}
	return modelMap, nil
}

func flattenTemplatePolicyResource(model *iampolicymanagementv1.V2PolicyResource, listRoleOptions *iampolicymanagementv1.ListRolesOptions, iamPolicyManagementClient *iampolicymanagementv1.IamPolicyManagementV1) (map[string]interface{}, *iampolicymanagementv1.RoleCollection, error) {
	modelMap := make(map[string]interface{})
	attributes := []map[string]interface{}{}

	for _, attributesItem := range model.Attributes {
		if *attributesItem.Key == "serviceName" &&
			(*attributesItem.Operator == "stringMatch" ||
				*attributesItem.Operator == "stringEquals") {
			listRoleOptions.ServiceName = core.StringPtr(attributesItem.Value.(string))
		}

		if *attributesItem.Key == "service_group_id" && (*attributesItem.Operator == "stringMatch" ||
			*attributesItem.Operator == "stringEquals") {
			listRoleOptions.ServiceGroupID = core.StringPtr(attributesItem.Value.(string))
		}

		if *attributesItem.Key == "serviceType" && attributesItem.Value.(string) == "service" && (*attributesItem.Operator == "stringMatch" ||
			*attributesItem.Operator == "stringEquals") {
			listRoleOptions.ServiceName = core.StringPtr("alliamserviceroles")
		}
		attributesItemMap := make(map[string]interface{})
		attributesItemMap["key"] = *attributesItem.Key
		attributesItemMap["operator"] = *attributesItem.Operator
		attributesItemMap["value"] = fmt.Sprintf("%v", attributesItem.Value)
		attributes = append(attributes, attributesItemMap)
	}
	modelMap["attributes"] = attributes
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap := make(map[string]interface{})
			tagsItemMap["key"] = *tagsItem.Key
			tagsItemMap["operator"] = *tagsItem.Operator
			tagsItemMap["value"] = *tagsItem.Value
			tags = append(tags, tagsItemMap)
		}
		modelMap["tags"] = tags
	}
	roleList, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
	if err != nil {
		return nil, nil, err
	}
	return modelMap, roleList, nil
}

func flattenTemplatePolicySubject(model *iampolicymanagementv1.V2PolicySubject, listRoleOptions *iampolicymanagementv1.ListRolesOptions) (map[string]interface{}, *iampolicymanagementv1.ListRolesOptions, error) {
	modelMap := make(map[string]interface{})
	attributes := []map[string]interface{}{}
	for _, attributesItem := range model.Attributes {
		attributesItemMap := make(map[string]interface{})
		attributesItemMap["key"] = *attributesItem.Key
		attributesItemMap["operator"] = *attributesItem.Operator
		attributesItemMap["value"] = *&attributesItem.Value
		if *attributesItem.Key == "serviceName" {
			listRoleOptions.SourceServiceName = core.StringPtr(fmt.Sprintf("%v", attributesItem.Value))
		}
		attributes = append(attributes, attributesItemMap)
	}
	modelMap["attributes"] = attributes

	return modelMap, listRoleOptions, nil
}
