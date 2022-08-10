// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v3/configurationgovernancev1"
)

// Functions that were changed for validation:
// - resourceIBMSccRuleMapToRuleCondition
// - resourceIBMSccRuleMapToRuleSingleProperty

const maxDepth = 1

func ResourceIBMSccRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccRuleCreate,
		ReadContext:   resourceIBMSccRuleRead,
		UpdateContext: resourceIBMSccRuleUpdate,
		DeleteContext: resourceIBMSccRuleDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Your IBM Cloud account ID.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-readable alias to assign to your rule.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "An extended description of your rule.",
			},
			"rule_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of rule. Rules that you create are `user_defined`.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Labels that you can use to group and search for similar rules, such as those that help you to meet a specific organization guideline.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"creation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the resource was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the user or application that created the resource.",
			},
			"modification_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// ForceNew:    true,      // Type 1 Fix
				Description: "The date the resource was last modified.",
			},
			"modified_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the user or application that last modified the resource.",
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enforcement_actions": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The actions that the service must run on your behalf when a request to create or modify the target resource does not comply with your conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "To block a request from completing, use `disallow`.",
						},
					},
				},
				MaxItems: 1,
			},
			"required_config": &schema.Schema{
				Description: "The requirements that must be met to determine the resource's level of compliance in accordance with the rule. Use logical operators (and/or) to define multiple property checks and conditions. To define requirements for a rule, list one or more property check objects in the and array. To add conditions to a property check, use or.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: getRequiredConfigSchema(0),
				},
				MaxItems: 1,
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The properties that describe the resource that you want to targetwith the rule or template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The programmatic name of the IBM Cloud service that you want to target with the rule or template.",
						},
						"resource_kind": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of resource that you want to target.",
						},
						"additional_target_attributes": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "An extra qualifier for the resource kind. When you include additional attributes, only the resources that match the definition are included in the rule or template.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the additional attribute that you want to use to further qualify the target.Options differ depending on the service or resource that you are targeting with a rule or template. For more information, refer to the service documentation.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The value that you want to apply to `name` field.Options differ depending on the rule or template that you configure. For more information, refer to the service documentation.",
									},
									"operator": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "The way in which the `name` field is compared to its value.There are three types of operators: string, numeric, and boolean.",
										ValidateFunc: validate.InvokeValidator("ibm_scc_rule", "operator"),
									},
								},
							},
						},
					},
				},
				MaxItems: 1,
			},
		},
		CustomizeDiff: customdiff.All(
			// update the version number via API GET if any of the fields are changed
			customdiff.ComputedIf("version", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.HasChange("name") || diff.HasChange("description") ||
					diff.HasChange("target") || diff.HasChange("labels") ||
					diff.HasChange("required_config") || diff.HasChange("enforcement_actions")
			}),
			// update the modification_date via API GET if any of the fields are changed
			customdiff.ComputedIf("modification_date", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.HasChange("name") || diff.HasChange("description") ||
					diff.HasChange("target") || diff.HasChange("labels") ||
					diff.HasChange("required_config") || diff.HasChange("enforcement_actions")
			}),
		),
	}
}

func getRequiredConfigSchema(currentDepth int) map[string]*schema.Schema {
	baseMap := map[string]*schema.Schema{
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The programmatic name of the IBM Cloud service that you want to target with the rule or template.",
		},
		"property": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the additional attribute that you want to use to further qualify the target.Options differ depending on the service or resource that you are targeting with a rule or template. For more information, refer to the service documentation.",
		},
		"value": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The value that you want to apply to `name` field.Options differ depending on the rule or template that you configure. For more information, refer to the service documentation.",
		},
		"operator": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "The way in which the `name` field is compared to its value.There are three types of operators: string, numeric, and boolean.",
			ValidateFunc: validate.InvokeValidator("ibm_scc_rule", "operator"),
		},
	}

	if currentDepth > maxDepth {
		return baseMap
	}
	baseMap["and"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A condition with the and logical operator.",
		Elem: &schema.Resource{
			Schema: getRequiredConfigSchema(currentDepth + 1),
		},
	}
	baseMap["or"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A condition with the or logical operator.",
		Elem: &schema.Resource{
			Schema: getRequiredConfigSchema(currentDepth + 1),
		},
	}
	return baseMap
}

func resourceIBMSccRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createRulesOptions := &configurationgovernancev1.CreateRulesOptions{}
	var rule []configurationgovernancev1.CreateRuleRequest
	ruleItem, err := resourceIBMSccRuleMapToCreateRuleRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	rule = append(rule, *ruleItem)
	createRulesOptions.SetRules(rule)

	createRulesResponse, response, err := configurationGovernanceClient.CreateRulesWithContext(context, createRulesOptions)
	if err != nil || response.GetStatusCode() == 207 || response.StatusCode > 300 {
		log.Printf("[DEBUG] CreateRulesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateRulesWithContext failed %s\n%s", err, response))
	}

	d.SetId(*createRulesResponse.Rules[0].Rule.RuleID)

	return resourceIBMSccRuleRead(context, d, meta)
}

func resourceIBMSccRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleOptions := &configurationgovernancev1.GetRuleOptions{}

	getRuleOptions.SetRuleID(d.Id())

	rule, response, err := configurationGovernanceClient.GetRuleWithContext(context, getRuleOptions)
	log.Println("[DEBUG] Grabbed a response from the Read Operation")

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRuleWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("account_id", rule.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("name", rule.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", rule.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("rule_type", rule.RuleType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rule_type: %s", err))
	}
	targetMap, e := resourceIBMSccRuleTargetResourceToMap(rule.Target)
	if e != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("target", []map[string]interface{}{targetMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target: %s", err))
	}
	requiredConfigMap, e := resourceIBMSccRuleRuleRequiredConfigToMap(rule.RequiredConfig)
	if e != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("required_config", []map[string]interface{}{requiredConfigMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting required_config: %s", err))
	}

	enforcementAction := []map[string]interface{}{}
	for _, enforcementActionItem := range rule.EnforcementActions {
		enforcementActionItemMap, err := resourceIBMSccRuleEnforcementActionToMap(&enforcementActionItem)
		if err != nil {
			return diag.FromErr(err)
		}
		enforcementAction = append(enforcementAction, enforcementActionItemMap)
	}

	if err = d.Set("enforcement_actions", enforcementAction); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enforcement_actions: %s", err))
	}
	if rule.Labels != nil {
		if err = d.Set("labels", rule.Labels); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting labels: %s", err))
		}
	}
	if err = d.Set("creation_date", flex.DateTimeToString(rule.CreationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting creation_date: %s", err))
	}
	if err = d.Set("created_by", rule.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("modification_date", flex.DateTimeToString(rule.ModificationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modification_date: %s", err))
	}
	if err = d.Set("modified_by", rule.ModifiedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_by: %s", err))
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	return nil
}

func resourceIBMSccRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateRuleOptions := &configurationgovernancev1.UpdateRuleOptions{}

	updateRuleOptions.SetRuleID(d.Id())

	hasChange := d.HasChange("name") || d.HasChange("description") ||
		d.HasChange("target") || d.HasChange("labels") ||
		d.HasChange("required_config") || d.HasChange("enforcement_actions")

	if hasChange {
		updateRuleOptions.SetName(d.Get("name").(string))
		updateRuleOptions.SetAccountID(d.Get("account_id").(string))
		updateRuleOptions.SetDescription(d.Get("description").(string))

		target, err := resourceIBMSccRuleMapToTargetResource(d.Get("target.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateRuleOptions.SetTarget(target)
		labels := []string{}
		if d.Get("labels") != nil {
			for _, labelsItem := range d.Get("labels").([]interface{}) {
				labels = append(labels, labelsItem.(string))
			}
		}
		updateRuleOptions.SetLabels(labels)

		required_config, err := resourceIBMSccRuleMapToRuleRequiredConfig(d.Get("required_config.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateRuleOptions.SetRequiredConfig(required_config)

		enforcementActions := []configurationgovernancev1.EnforcementAction{}
		for _, enforcementActionsItem := range d.Get("enforcement_actions").([]interface{}) {
			if enforcementActionsItem != nil {
				enforcementActionsItemModel, err := resourceIBMSccRuleMapToEnforcementAction(enforcementActionsItem.(map[string]interface{}))
				if err != nil {
					return diag.FromErr(err)
				}
				enforcementActions = append(enforcementActions, *enforcementActionsItemModel)
			}
		}
		updateRuleOptions.SetEnforcementActions(enforcementActions)

		updateRuleOptions.SetIfMatch(d.Get("version").(string))
		_, response, err := configurationGovernanceClient.UpdateRuleWithContext(context, updateRuleOptions)
		if err != nil || response.GetStatusCode() == 207 || response.StatusCode > 300 {
			log.Printf("[DEBUG] UpdateRuleWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateRuleWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMSccRuleRead(context, d, meta)
}

func resourceIBMSccRuleDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRuleOptions := &configurationgovernancev1.DeleteRuleOptions{}

	deleteRuleOptions.SetRuleID(d.Id())

	response, err := configurationGovernanceClient.DeleteRuleWithContext(context, deleteRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMSccRuleMapToCreateRuleRequest(d *schema.ResourceData) (*configurationgovernancev1.CreateRuleRequest, error) {
	model := &configurationgovernancev1.CreateRuleRequest{}
	if d.Get("request_id") != nil {
		model.RequestID = core.StringPtr(d.Get("request_id").(string))
	}
	RuleModel, err := resourceIBMSccRuleMapToRuleRequest(d)
	if err != nil {
		return model, err
	}
	model.Rule = RuleModel
	return model, nil
}

func resourceIBMSccRuleMapToRuleRequest(d *schema.ResourceData) (*configurationgovernancev1.RuleRequest, error) {
	model := &configurationgovernancev1.RuleRequest{}
	if d.Get("account_id") != nil {
		model.AccountID = core.StringPtr(d.Get("account_id").(string))
	}
	model.Name = core.StringPtr(d.Get("name").(string))
	model.Description = core.StringPtr(d.Get("description").(string))
	if d.Get("rule_type") != nil {
		model.RuleType = core.StringPtr(d.Get("rule_type").(string))
	}
	targetList := d.Get("target").([]interface{})
	TargetModel, err := resourceIBMSccRuleMapToTargetResource(targetList[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Target = TargetModel
	requiredConfigList := d.Get("required_config").([]interface{})
	RequiredConfigModel, err := resourceIBMSccRuleMapToRuleRequiredConfig(requiredConfigList[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.RequiredConfig = RequiredConfigModel
	enforcementActions := []configurationgovernancev1.EnforcementAction{}
	for _, enforcementActionsItem := range d.Get("enforcement_actions").([]interface{}) {
		if enforcementActionsItem != nil {
			enforcementActionsItemModel, err := resourceIBMSccRuleMapToEnforcementAction(enforcementActionsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			enforcementActions = append(enforcementActions, *enforcementActionsItemModel)
		}
	}
	model.EnforcementActions = enforcementActions
	if d.Get("labels") != nil {
		labels := []string{}
		for _, labelsItem := range d.Get("labels").([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	return model, nil
}

func resourceIBMSccRuleMapToTargetResource(modelMap map[string]interface{}) (*configurationgovernancev1.TargetResource, error) {
	model := &configurationgovernancev1.TargetResource{}
	model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	model.ResourceKind = core.StringPtr(modelMap["resource_kind"].(string))
	if modelMap["additional_target_attributes"] != nil {
		additionalTargetAttributes := []configurationgovernancev1.TargetResourceAdditionalTargetAttributesItem{}
		for _, additionalTargetAttributesItem := range modelMap["additional_target_attributes"].([]interface{}) {
			if additionalTargetAttributesItem != nil {
				additionalTargetAttributesItemModel, err := resourceIBMSccRuleMapToTargetResourceAdditionalTargetAttributesItem(additionalTargetAttributesItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				additionalTargetAttributes = append(additionalTargetAttributes, *additionalTargetAttributesItemModel)
			}
		}
		model.AdditionalTargetAttributes = additionalTargetAttributes
	}
	return model, nil
}

func resourceIBMSccRuleMapToTargetResourceAdditionalTargetAttributesItem(modelMap map[string]interface{}) (*configurationgovernancev1.TargetResourceAdditionalTargetAttributesItem, error) {
	model := &configurationgovernancev1.TargetResourceAdditionalTargetAttributesItem{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	return model, nil
}

func resourceIBMSccRuleMapToRuleRequiredConfig(modelMap map[string]interface{}) (configurationgovernancev1.RuleRequiredConfigIntf, error) {
	model := &configurationgovernancev1.RuleRequiredConfig{}

	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["property"] != nil {
		model.Property = core.StringPtr(modelMap["property"].(string))
	}
	if modelMap["operator"] != nil {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	// TODO: handle the usage of Lists/Arrays of strings(can't be done until the go-sdk is modified)
	if modelMap["value"] != nil {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	if modelMap["or"] != nil {
		or := []configurationgovernancev1.RuleConditionIntf{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			if orItem == nil {
				return model, errors.New("or block needs to be populated")
			}
			orItemModel, err := resourceIBMSccRuleMapToRuleCondition(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	if modelMap["and"] != nil {
		and := []configurationgovernancev1.RuleConditionIntf{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			if andItem == nil {
				return model, errors.New("and block needs to be populated")
			}
			andItemModel, err := resourceIBMSccRuleMapToRuleCondition(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, andItemModel)
		}
		model.And = and
	}
	// Error out if 'and' and 'or' are set at the same level
	if len(model.And) > 0 && len(model.Or) > 0 {
		return model, errors.New("attributes of required_config 'or' and 'and' cannot be set at the same level")
	}

	// Error out if the property, value, and operator are at the same level as 'and' and 'or'
	if (len(*model.Value) > 0 || len(*model.Property) > 0 || len(*model.Operator) > 0) &&
		(len(model.And) > 0 || len(model.Or) > 0) {
		return model, errors.New("'property','value','operator' should be nested inside 'and'/'or' or be by itself")
	}

	return model, nil
}

func resourceIBMSccRuleMapToRuleCondition(modelMap map[string]interface{}) (configurationgovernancev1.RuleConditionIntf, error) {
	model := &configurationgovernancev1.RuleCondition{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["property"] != nil {
		model.Property = core.StringPtr(modelMap["property"].(string))
	}
	if modelMap["operator"] != nil {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	if modelMap["value"] != nil {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	if modelMap["or"] != nil {
		or := []configurationgovernancev1.RuleSingleProperty{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			if orItem == nil {
				return model, errors.New("or block needs to be populated")
			}
			orItemModel, err := resourceIBMSccRuleMapToRuleSingleProperty(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, *orItemModel)
		}
		model.Or = or
	}
	if modelMap["and"] != nil {
		and := []configurationgovernancev1.RuleSingleProperty{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			if andItem == nil {
				return model, errors.New("and block needs to be populated")
			}
			andItemModel, err := resourceIBMSccRuleMapToRuleSingleProperty(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, *andItemModel)
		}
		model.And = and
	}
	// Error out if 'and' and 'or' are set at the same level
	if len(model.And) > 0 && len(model.Or) > 0 {
		return model, errors.New("attributes of required_config 'or' and 'and' cannot be set at the same level")
	}

	// Error out if the property, value, and operator are at the same level as 'and' and 'or'
	if (len(*model.Value) > 0 || len(*model.Property) > 0 || len(*model.Operator) > 0) &&
		(len(model.And) > 0 || len(model.Or) > 0) {
		return model, errors.New("'property','value','operator' should be nested inside 'and'/'or' or be by itself")
	}
	return model, nil
}

func resourceIBMSccRuleMapToRuleSingleProperty(modelMap map[string]interface{}) (*configurationgovernancev1.RuleSingleProperty, error) {
	model := &configurationgovernancev1.RuleSingleProperty{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Property = core.StringPtr(modelMap["property"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	if modelMap["value"] != nil {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func resourceIBMSccRuleMapToRuleConditionSingleProperty(modelMap map[string]interface{}) (*configurationgovernancev1.RuleConditionSingleProperty, error) {
	model := &configurationgovernancev1.RuleConditionSingleProperty{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Property = core.StringPtr(modelMap["property"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	if modelMap["value"] != nil {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func resourceIBMSccRuleMapToRuleConditionOrLvl2(modelMap map[string]interface{}) (*configurationgovernancev1.RuleConditionOrLvl2, error) {
	model := &configurationgovernancev1.RuleConditionOrLvl2{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	or := []configurationgovernancev1.RuleSingleProperty{}
	for _, orItem := range modelMap["or"].([]interface{}) {
		orItemModel, err := resourceIBMSccRuleMapToRuleSingleProperty(orItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		or = append(or, *orItemModel)
	}
	model.Or = or
	return model, nil
}

func resourceIBMSccRuleMapToRuleConditionAndLvl2(modelMap map[string]interface{}) (*configurationgovernancev1.RuleConditionAndLvl2, error) {
	model := &configurationgovernancev1.RuleConditionAndLvl2{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	and := []configurationgovernancev1.RuleSingleProperty{}
	for _, andItem := range modelMap["and"].([]interface{}) {
		andItemModel, err := resourceIBMSccRuleMapToRuleSingleProperty(andItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		and = append(and, *andItemModel)
	}
	model.And = and
	return model, nil
}

func resourceIBMSccRuleMapToRuleRequiredConfigSingleProperty(modelMap map[string]interface{}) (*configurationgovernancev1.RuleRequiredConfigSingleProperty, error) {
	model := &configurationgovernancev1.RuleRequiredConfigSingleProperty{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Property = core.StringPtr(modelMap["property"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	if modelMap["value"] != nil {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func resourceIBMSccRuleMapToRuleRequiredConfigMultipleProperties(modelMap map[string]interface{}) (configurationgovernancev1.RuleRequiredConfigMultiplePropertiesIntf, error) {
	model := &configurationgovernancev1.RuleRequiredConfigMultipleProperties{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["or"] != nil {
		or := []configurationgovernancev1.RuleConditionIntf{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			orItemModel, err := resourceIBMSccRuleMapToRuleCondition(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	if modelMap["and"] != nil {
		and := []configurationgovernancev1.RuleConditionIntf{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			andItemModel, err := resourceIBMSccRuleMapToRuleCondition(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, andItemModel)
		}
		model.And = and
	}
	return model, nil
}

func resourceIBMSccRuleMapToRuleRequiredConfigMultiplePropertiesConditionOr(modelMap map[string]interface{}) (*configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionOr, error) {
	model := &configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionOr{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	or := []configurationgovernancev1.RuleConditionIntf{}
	for _, orItem := range modelMap["or"].([]interface{}) {
		orItemModel, err := resourceIBMSccRuleMapToRuleCondition(orItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		or = append(or, orItemModel)
	}
	model.Or = or
	return model, nil
}

func resourceIBMSccRuleMapToRuleRequiredConfigMultiplePropertiesConditionAnd(modelMap map[string]interface{}) (*configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionAnd, error) {
	model := &configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionAnd{}
	if modelMap["description"] != nil {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	and := []configurationgovernancev1.RuleConditionIntf{}
	for _, andItem := range modelMap["and"].([]interface{}) {
		andItemModel, err := resourceIBMSccRuleMapToRuleCondition(andItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		and = append(and, andItemModel)
	}
	model.And = and
	return model, nil
}

func resourceIBMSccRuleMapToEnforcementAction(modelMap map[string]interface{}) (*configurationgovernancev1.EnforcementAction, error) {
	model := &configurationgovernancev1.EnforcementAction{}
	model.Action = core.StringPtr(modelMap["action"].(string))
	return model, nil
}

func resourceIBMSccRuleCreateRuleRequestToMap(model *configurationgovernancev1.CreateRuleRequest) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RequestID != nil {
		modelMap["request_id"] = model.RequestID
	}
	ruleMap, err := resourceIBMSccRuleRuleRequestToMap(model.Rule)
	if err != nil {
		return modelMap, err
	}
	modelMap["rule"] = []map[string]interface{}{ruleMap}
	return modelMap, nil
}

func resourceIBMSccRuleRuleRequestToMap(model *configurationgovernancev1.RuleRequest) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AccountID != nil {
		modelMap["account_id"] = model.AccountID
	}
	modelMap["name"] = model.Name
	modelMap["description"] = model.Description
	if model.RuleType != nil {
		modelMap["rule_type"] = model.RuleType
	}
	targetMap, err := resourceIBMSccRuleTargetResourceToMap(model.Target)
	if err != nil {
		return modelMap, err
	}
	modelMap["target"] = []map[string]interface{}{targetMap}
	requiredConfigMap, err := resourceIBMSccRuleRuleRequiredConfigToMap(model.RequiredConfig)
	if err != nil {
		return modelMap, err
	}
	modelMap["required_config"] = []map[string]interface{}{requiredConfigMap}
	enforcementActions := []map[string]interface{}{}
	for _, enforcementActionsItem := range model.EnforcementActions {
		enforcementActionsItemMap, err := resourceIBMSccRuleEnforcementActionToMap(&enforcementActionsItem)
		if err != nil {
			return modelMap, err
		}
		enforcementActions = append(enforcementActions, enforcementActionsItemMap)
	}
	modelMap["enforcement_actions"] = enforcementActions
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	return modelMap, nil
}

func resourceIBMSccRuleTargetResourceToMap(model *configurationgovernancev1.TargetResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["service_name"] = model.ServiceName
	modelMap["resource_kind"] = model.ResourceKind
	if model.AdditionalTargetAttributes != nil {
		additionalTargetAttributes := []map[string]interface{}{}
		for _, additionalTargetAttributesItem := range model.AdditionalTargetAttributes {
			additionalTargetAttributesItemMap, err := resourceIBMSccRuleTargetResourceAdditionalTargetAttributesItemToMap(&additionalTargetAttributesItem)
			if err != nil {
				return modelMap, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, additionalTargetAttributesItemMap)
		}
		modelMap["additional_target_attributes"] = additionalTargetAttributes
	}
	return modelMap, nil
}

func resourceIBMSccRuleTargetResourceAdditionalTargetAttributesItemToMap(model *configurationgovernancev1.TargetResourceAdditionalTargetAttributesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	modelMap["operator"] = model.Operator
	return modelMap, nil
}

func resourceIBMSccRuleRuleRequiredConfigToMap(model configurationgovernancev1.RuleRequiredConfigIntf) (map[string]interface{}, error) {
	if _, ok := model.(*configurationgovernancev1.RuleRequiredConfigSingleProperty); ok {
		return resourceIBMSccRuleRuleRequiredConfigSinglePropertyToMap(model.(*configurationgovernancev1.RuleRequiredConfigSingleProperty))
	} else if _, ok := model.(*configurationgovernancev1.RuleRequiredConfigMultipleProperties); ok {
		return resourceIBMSccRuleRuleRequiredConfigMultiplePropertiesToMap(model.(*configurationgovernancev1.RuleRequiredConfigMultipleProperties))
	} else if _, ok := model.(*configurationgovernancev1.RuleRequiredConfig); ok {
		modelMap := make(map[string]interface{})
		model := model.(*configurationgovernancev1.RuleRequiredConfig)
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Property != nil {
			modelMap["property"] = model.Property
		}
		if model.Operator != nil {
			modelMap["operator"] = model.Operator
		}
		if model.Value != nil {
			modelMap["value"] = model.Value
		}
		if model.Or != nil {
			or := []map[string]interface{}{}
			for _, orItem := range model.Or {
				orItemMap, err := resourceIBMSccRuleRuleConditionToMap(orItem)
				if err != nil {
					return modelMap, err
				}
				or = append(or, orItemMap)
			}
			modelMap["or"] = or
		}
		if model.And != nil {
			and := []map[string]interface{}{}
			for _, andItem := range model.And {
				andItemMap, err := resourceIBMSccRuleRuleConditionToMap(andItem)
				if err != nil {
					return modelMap, err
				}
				and = append(and, andItemMap)
			}
			modelMap["and"] = and
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized configurationgovernancev1.RuleRequiredConfigIntf subtype encountered")
	}
}

func resourceIBMSccRuleRuleConditionToMap(model configurationgovernancev1.RuleConditionIntf) (map[string]interface{}, error) {
	if _, ok := model.(*configurationgovernancev1.RuleConditionSingleProperty); ok {
		return resourceIBMSccRuleRuleConditionSinglePropertyToMap(model.(*configurationgovernancev1.RuleConditionSingleProperty))
	} else if _, ok := model.(*configurationgovernancev1.RuleConditionOrLvl2); ok {
		return resourceIBMSccRuleRuleConditionOrLvl2ToMap(model.(*configurationgovernancev1.RuleConditionOrLvl2))
	} else if _, ok := model.(*configurationgovernancev1.RuleConditionAndLvl2); ok {
		return resourceIBMSccRuleRuleConditionAndLvl2ToMap(model.(*configurationgovernancev1.RuleConditionAndLvl2))
	} else if _, ok := model.(*configurationgovernancev1.RuleCondition); ok {
		modelMap := make(map[string]interface{})
		model := model.(*configurationgovernancev1.RuleCondition)
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Property != nil {
			modelMap["property"] = model.Property
		}
		if model.Operator != nil {
			modelMap["operator"] = model.Operator
		}
		if model.Value != nil {
			modelMap["value"] = model.Value
		}
		if model.Or != nil {
			or := []map[string]interface{}{}
			for _, orItem := range model.Or {
				orItemMap, err := resourceIBMSccRuleRuleSinglePropertyToMap(&orItem)
				if err != nil {
					return modelMap, err
				}
				or = append(or, orItemMap)
			}
			modelMap["or"] = or
		}
		if model.And != nil {
			and := []map[string]interface{}{}
			for _, andItem := range model.And {
				andItemMap, err := resourceIBMSccRuleRuleSinglePropertyToMap(&andItem)
				if err != nil {
					return modelMap, err
				}
				and = append(and, andItemMap)
			}
			modelMap["and"] = and
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized configurationgovernancev1.RuleConditionIntf subtype encountered")
	}
}

func resourceIBMSccRuleRuleSinglePropertyToMap(model *configurationgovernancev1.RuleSingleProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["property"] = model.Property
	modelMap["operator"] = model.Operator
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIBMSccRuleRuleConditionSinglePropertyToMap(model *configurationgovernancev1.RuleConditionSingleProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["property"] = model.Property
	modelMap["operator"] = model.Operator
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIBMSccRuleRuleConditionOrLvl2ToMap(model *configurationgovernancev1.RuleConditionOrLvl2) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	or := []map[string]interface{}{}
	for _, orItem := range model.Or {
		orItemMap, err := resourceIBMSccRuleRuleSinglePropertyToMap(&orItem)
		if err != nil {
			return modelMap, err
		}
		or = append(or, orItemMap)
	}
	modelMap["or"] = or
	return modelMap, nil
}

func resourceIBMSccRuleRuleConditionAndLvl2ToMap(model *configurationgovernancev1.RuleConditionAndLvl2) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	and := []map[string]interface{}{}
	for _, andItem := range model.And {
		andItemMap, err := resourceIBMSccRuleRuleSinglePropertyToMap(&andItem)
		if err != nil {
			return modelMap, err
		}
		and = append(and, andItemMap)
	}
	modelMap["and"] = and
	return modelMap, nil
}

func resourceIBMSccRuleRuleRequiredConfigSinglePropertyToMap(model *configurationgovernancev1.RuleRequiredConfigSingleProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["property"] = model.Property
	modelMap["operator"] = model.Operator
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIBMSccRuleRuleRequiredConfigMultiplePropertiesToMap(model configurationgovernancev1.RuleRequiredConfigMultiplePropertiesIntf) (map[string]interface{}, error) {
	if _, ok := model.(*configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionOr); ok {
		return resourceIBMSccRuleRuleRequiredConfigMultiplePropertiesConditionOrToMap(model.(*configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionOr))
	} else if _, ok := model.(*configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionAnd); ok {
		return resourceIBMSccRuleRuleRequiredConfigMultiplePropertiesConditionAndToMap(model.(*configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionAnd))
	} else if _, ok := model.(*configurationgovernancev1.RuleRequiredConfigMultipleProperties); ok {
		modelMap := make(map[string]interface{})
		model := model.(*configurationgovernancev1.RuleRequiredConfigMultipleProperties)
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Or != nil {
			or := []map[string]interface{}{}
			for _, orItem := range model.Or {
				orItemMap, err := resourceIBMSccRuleRuleConditionToMap(orItem)
				if err != nil {
					return modelMap, err
				}
				or = append(or, orItemMap)
			}
			modelMap["or"] = or
		}
		if model.And != nil {
			and := []map[string]interface{}{}
			for _, andItem := range model.And {
				andItemMap, err := resourceIBMSccRuleRuleConditionToMap(andItem)
				if err != nil {
					return modelMap, err
				}
				and = append(and, andItemMap)
			}
			modelMap["and"] = and
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized configurationgovernancev1.RuleRequiredConfigMultiplePropertiesIntf subtype encountered")
	}
}

func resourceIBMSccRuleRuleRequiredConfigMultiplePropertiesConditionOrToMap(model *configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionOr) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	or := []map[string]interface{}{}
	for _, orItem := range model.Or {
		orItemMap, err := resourceIBMSccRuleRuleConditionToMap(orItem)
		if err != nil {
			return modelMap, err
		}
		or = append(or, orItemMap)
	}
	modelMap["or"] = or
	return modelMap, nil
}

func resourceIBMSccRuleRuleRequiredConfigMultiplePropertiesConditionAndToMap(model *configurationgovernancev1.RuleRequiredConfigMultiplePropertiesConditionAnd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	and := []map[string]interface{}{}
	for _, andItem := range model.And {
		andItemMap, err := resourceIBMSccRuleRuleConditionToMap(andItem)
		if err != nil {
			return modelMap, err
		}
		and = append(and, andItemMap)
	}
	modelMap["and"] = and
	return modelMap, nil
}

func resourceIBMSccRuleEnforcementActionToMap(model *configurationgovernancev1.EnforcementAction) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["action"] = model.Action
	return modelMap, nil
}
