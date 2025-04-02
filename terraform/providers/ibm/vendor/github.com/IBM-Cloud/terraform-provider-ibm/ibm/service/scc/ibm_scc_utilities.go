// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	INSTANCE_ID               = "instance_id"
	MAX_REQUIRED_CONFIG_DEPTH = 5
)

// AddSchemaData will add the Schemas 'instance_id' and 'region' to the resource
func AddSchemaData(resource *schema.Resource) *schema.Resource {
	resource.Schema["instance_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the Security and Compliance Center instance.",
	}
	return resource
}

// getRegionData will check if the field region is defined
func getRegionData(client securityandcompliancecenterapiv3.SecurityAndComplianceCenterApiV3, d *schema.ResourceData) string {
	val, ok := d.GetOk("region")
	if ok {
		return val.(string)
	} else {
		url := client.Service.GetServiceURL()
		return strings.Split(url, ".")[1]
	}
}

// setRegionData will set the field "region" field if the field was previously defined
func setRegionData(d *schema.ResourceData, region string) error {
	if val, ok := d.GetOk("region"); ok {
		return d.Set("region", val.(string))
	}
	return nil
}

// getRequiredConfigSchema will return the schema for a scc rule required_config. This schema is recursive.
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
			Description: "The value that you want to apply to `value` field. Options differ depending on the rule or template that you configure. For more information, refer to the service documentation.",
		},
		"operator": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "The way in which the `name` field is compared to its value.There are three types of operators: string, numeric, and boolean.",
			ValidateFunc: validate.InvokeValidator("ibm_scc_rule", "operator"),
		},
	}
	if currentDepth > MAX_REQUIRED_CONFIG_DEPTH {
		return baseMap
	}

	baseMap["and"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A list of property conditions where all items need to be satisfied",
		Elem: &schema.Resource{
			Schema: getRequiredConfigSchema(currentDepth + 1),
		},
	}

	baseMap["or"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A list of property conditions where any item needs to be satisfied",
		Elem: &schema.Resource{
			Schema: getRequiredConfigSchema(currentDepth + 1),
		},
	}

	baseMap["all"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A condition with the SubRule all logical operator.",
		Elem: &schema.Resource{
			Schema: getSubRuleSchema(currentDepth + 1),
		},
	}

	baseMap["all_if"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A condition with the SubRule all_ifexists logical operator.",
		Elem: &schema.Resource{
			Schema: getSubRuleSchema(currentDepth + 1),
		},
	}

	baseMap["any"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A condition with the SubRule any logical operator.",
		Elem: &schema.Resource{
			Schema: getSubRuleSchema(currentDepth + 1),
		},
	}

	baseMap["any_if"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "A condition with the SubRule any_ifexists logical operator.",
		Elem: &schema.Resource{
			Schema: getSubRuleSchema(currentDepth + 1),
		},
	}
	return baseMap
}

// getTargetSchema returns a terraform Schema defining the attributes of a Target object
func getTargetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The target service name.",
		},
		"service_display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The display name of the target service.",
			// Manual Intervention
			DiffSuppressFunc: func(_, oldVal, newVal string, d *schema.ResourceData) bool {
				if newVal == "" {
					return true
				}
				if strings.ToLower(oldVal) == strings.ToLower(newVal) {
					return true
				}
				return false
			},
			// End Manual Intervention
		},
		"reference_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The target reference name",
		},
		"resource_kind": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The target resource kind.",
		},
		"additional_target_attributes": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The list of targets supported properties.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The additional target attribute name.",
					},
					"operator": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The operator.",
					},
					"value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The value.",
					},
				},
			},
		},
	}
}

// getSubRuleSchema returns a terraform Schema that define attributes of a subRule
func getSubRuleSchema(currentDepth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"required_config": {
			Description: "The requirements that must be met to determine the resource's level of compliance in accordance with the rule. Use logical operators (and/or) to define multiple property checks and conditions. To define requirements for a rule, list one or more property check objects in the and array. To add conditions to a property check, use or.",
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Resource{
				Schema: getRequiredConfigSchema(currentDepth + 1),
			},
			MaxItems: 1,
		},
		"target": {
			Description: "The requirements that must be met to determine the resource's level of compliance in accordance with the rule. Use logical operators (and/or) to define multiple property checks and conditions. To define requirements for a rule, list one or more property check objects in the and array. To add conditions to a property check, use or.",
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Resource{
				Schema: getTargetSchema(),
			},
			MaxItems: 1,
		},
	}
}

// requiredConfigItemsToListMap will dicipher the list of conditions and return a []map[string]interface{} to abide to the
// terraform type TypeList
func requiredConfigItemsToListMap(items []securityandcompliancecenterapiv3.RequiredConfigIntf) ([]map[string]interface{}, error) {

	rcItems := []map[string]interface{}{}
	for _, rcItem := range items {
		rcMap, err := requiredConfigToModelMap(rcItem)
		if err != nil {
			return []map[string]interface{}{}, err
		}
		rcItems = append(rcItems, rcMap)
	}
	return rcItems, nil
}

// requiredConfigSubRuleToMap will dicipher the sub rule brought in and return a []map[string]interface{} for the terraform
// state file
func requiredConfigSubRuleToMap(subRule *securityandcompliancecenterapiv3.SubRule) ([]map[string]interface{}, error) {
	srMap := make(map[string]interface{})
	subRuleTarget, err := targetToModelMap(subRule.Target)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	srMap["target"] = []interface{}{subRuleTarget}
	subRuleConfig, err := requiredConfigToModelMap(subRule.RequiredConfig)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	srMap["required_config"] = []interface{}{subRuleConfig}
	return []map[string]interface{}{srMap}, nil
}

// ibmSccRuleRequiredConfigToMap will dicipher a scc rule required_config and return back to a map[string]interface{}
// for the terraform state file
func requiredConfigToModelMap(model securityandcompliancecenterapiv3.RequiredConfigIntf) (map[string]interface{}, error) {
	if rc, ok := model.(*securityandcompliancecenterapiv3.RequiredConfig); ok {
		modelMap := make(map[string]interface{})
		if rc.Description != nil {
			modelMap["description"] = rc.Description
		}
		if rc.And != nil {
			if rcItems, err := requiredConfigItemsToListMap(rc.And); err != nil {
				return map[string]interface{}{}, err
			} else {
				modelMap["and"] = rcItems
			}
		}
		if rc.Or != nil {
			if rcItems, err := requiredConfigItemsToListMap(rc.Or); err != nil {
				return map[string]interface{}{}, err
			} else {
				modelMap["or"] = rcItems
			}
		}
		// sub rules
		if rc.All != nil {
			if subRule, err := requiredConfigSubRuleToMap(rc.All); err != nil {
				return map[string]interface{}{}, err
			} else {
				modelMap["all"] = subRule
			}
		}
		if rc.AllIfexists != nil {
			if subRule, err := requiredConfigSubRuleToMap(rc.All); err != nil {
				return map[string]interface{}{}, err
			} else {
				modelMap["all_if"] = subRule
			}
		}
		if rc.Any != nil {
			if subRule, err := requiredConfigSubRuleToMap(rc.Any); err != nil {
				return map[string]interface{}{}, err
			} else {
				modelMap["any"] = subRule
			}
		}
		if rc.AnyIfexists != nil {
			if subRule, err := requiredConfigSubRuleToMap(rc.AnyIfexists); err != nil {
				return map[string]interface{}{}, err
			} else {
				modelMap["any_if"] = subRule
			}
		}
		// base condition handling
		if rc.Property != nil {
			modelMap["property"] = rc.Property
		}
		if rc.Operator != nil {
			modelMap["operator"] = rc.Operator
		}
		if rc.Value != nil {
			// rc.Value can be implicitly cast as a []interface, it needs to be stringified
			if valList, ok := rc.Value.([]interface{}); ok {
				s := make([]string, len(valList))
				for i, v := range valList {
					s[i] = fmt.Sprint(v)
				}
				modelMap["value"] = fmt.Sprintf("[%s]", strings.Join(s, ","))
			} else {
				modelMap["value"] = rc.Value
			}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized securityandcompliancecenterapiv3.RequiredConfigIntf subtype encountered %#v", model)
	}
}

// listMapToSccSubRule is a helper function that converts a map into a SubRule
func listMapToSccSubRule(subRuleModel []interface{}) (*securityandcompliancecenterapiv3.SubRule, error) {
	subRule := securityandcompliancecenterapiv3.SubRule{}
	subMap := subRuleModel[0].(map[string]interface{})
	target, err := modelMapToTarget(subMap["target"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return &subRule, err
	}
	subRule.Target = target
	rc, err := modelMapToRequiredConfig(subMap["required_config"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return &subRule, err
	}
	subRule.RequiredConfig = rc.(*securityandcompliancecenterapiv3.RequiredConfig)
	return &subRule, nil
}

// modelMapToRequiredConfig converts the map to a RequiredConfig
func modelMapToRequiredConfig(modelMap map[string]interface{}) (securityandcompliancecenterapiv3.RequiredConfigIntf, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfig{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["or"] != nil {
		or := []securityandcompliancecenterapiv3.RequiredConfigIntf{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			orItemModel, err := modelMapToRequiredConfig(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	if modelMap["and"] != nil {
		and := []securityandcompliancecenterapiv3.RequiredConfigIntf{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			andItemModel, err := modelMapToRequiredConfig(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, andItemModel)
		}
		model.And = and
	}
	if anySM, ok := modelMap["any"].([]interface{}); ok && len(anySM) > 0 {
		anySubRule, err := listMapToSccSubRule(anySM)
		if err != nil {
			return model, err
		}
		model.Any = anySubRule
	}
	if anyIfSM, ok := modelMap["any_if"].([]interface{}); ok && len(anyIfSM) > 0 {
		anyIfSubRule, err := listMapToSccSubRule(anyIfSM)
		if err != nil {
			return model, err
		}
		model.AnyIfexists = anyIfSubRule
	}
	if allSM, ok := modelMap["all"].([]interface{}); ok && len(allSM) > 0 {
		allSubRule, err := listMapToSccSubRule(allSM)
		if err != nil {
			return model, err
		}
		model.All = allSubRule
	}
	if allIfSM, ok := modelMap["all_if"].([]interface{}); ok && len(allIfSM) > 0 {
		allIfSubRule, err := listMapToSccSubRule(allIfSM)
		if err != nil {
			return model, err
		}
		model.AllIfexists = allIfSubRule
	}
	if modelMap["property"] != nil && modelMap["property"].(string) != "" {
		model.Property = core.StringPtr(modelMap["property"].(string))
	}
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	if modelMap["value"] != nil && len(modelMap["value"].(string)) > 0 {
		// model.Value = modelMap["value"].(string)
		sLit := strings.Trim(modelMap["value"].(string), "[]")
		sList := strings.Split(sLit, ",")
		if len(sList) == 1 {
			model.Value = modelMap["value"].(string)
		} else {
			model.Value = sList
		}
	}

	return model, nil
}

// modelMapToTarget transforms a map and transforms into a Target object
func modelMapToTarget(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RuleTarget, error) {
	model := &securityandcompliancecenterapiv3.RuleTarget{}
	model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	if modelMap["reference_name"] != nil && modelMap["reference_name"].(string) != "" {
		model.Ref = core.StringPtr(modelMap["reference_name"].(string))
	}
	model.ResourceKind = core.StringPtr(modelMap["resource_kind"].(string))
	if modelMap["additional_target_attributes"] != nil {
		additionalTargetAttributes := []securityandcompliancecenterapiv3.AdditionalTargetAttribute{}
		for _, additionalTargetAttributesItem := range modelMap["additional_target_attributes"].([]interface{}) {
			additionalTargetAttributesItemModel, err := ruleMapToAdditionalTargetAttribute(additionalTargetAttributesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, *additionalTargetAttributesItemModel)
		}
		model.AdditionalTargetAttributes = additionalTargetAttributes
	}
	return model, nil
}

// modelMapToTargetPrototype transforms a map and transforms into a TargetPrototype object
func modelMapToTargetPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RuleTargetPrototype, error) {
	model := &securityandcompliancecenterapiv3.RuleTargetPrototype{}
	model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	if modelMap["reference_name"] != nil && modelMap["reference_name"].(string) != "" {
		model.Ref = core.StringPtr(modelMap["reference_name"].(string))
	}
	model.ResourceKind = core.StringPtr(modelMap["resource_kind"].(string))
	if modelMap["additional_target_attributes"] != nil {
		additionalTargetAttributes := []securityandcompliancecenterapiv3.AdditionalTargetAttribute{}
		for _, additionalTargetAttributesItem := range modelMap["additional_target_attributes"].([]interface{}) {
			additionalTargetAttributesItemModel, err := ruleMapToAdditionalTargetAttribute(additionalTargetAttributesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, *additionalTargetAttributesItemModel)
		}
		model.AdditionalTargetAttributes = additionalTargetAttributes
	}
	return model, nil
}

// ruleMapToAdditionalTargetAttribute will convert a given map to a AdditionalTargetAttribute object
func ruleMapToAdditionalTargetAttribute(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.AdditionalTargetAttribute, error) {
	model := &securityandcompliancecenterapiv3.AdditionalTargetAttribute{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

// targetToModelMap will convert a Target object to a map for the terraform state file.
func targetToModelMap(model *securityandcompliancecenterapiv3.RuleTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})

	modelMap["service_name"] = model.ServiceName

	modelMap["resource_kind"] = model.ResourceKind

	if model.Ref != nil {
		modelMap["reference_name"] = model.Ref
	}

	if model.ServiceDisplayName != nil {
		modelMap["service_display_name"] = model.ServiceDisplayName
	}

	if model.AdditionalTargetAttributes != nil {
		additionalTargetAttributes := []map[string]interface{}{}
		for _, additionalTargetAttributesItem := range model.AdditionalTargetAttributes {
			additionalTargetAttributesItemMap, err := ruleAdditionalTargetAttributeToMap(&additionalTargetAttributesItem)
			if err != nil {
				return modelMap, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, additionalTargetAttributesItemMap)
		}
		modelMap["additional_target_attributes"] = additionalTargetAttributes
	}
	return modelMap, nil
}

// converts a AdditionalTargetAttribute object into a map for the terraform state file.
func ruleAdditionalTargetAttributeToMap(model *securityandcompliancecenterapiv3.AdditionalTargetAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

// scopePropertiesToMap returns a map[string]interface{}
//
// This function is used for any scc resource/datasource that
// need to read scope.properties
func scopePropertiesToMap(model securityandcompliancecenterapiv3.ScopePropertyIntf) (map[string]interface{}, error) {
	if prop, ok := model.(*securityandcompliancecenterapiv3.ScopeProperty); ok && prop.Name != nil && prop.Value != nil {
		modelMap := make(map[string]interface{})
		modelMap["name"] = prop.Name
		if val, ok := prop.Value.(string); !ok {
			modelMap["value"] = fmt.Sprintf("%v", val)
		} else {
			modelMap["value"] = prop.Value
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized securityandcompliancecenterv3.ScopePropertyIntf subtype encountered")
	}
}
