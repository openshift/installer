// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package metricsrouter

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/metricsrouterv3"
)

func ResourceIBMMetricsRouterRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMMetricsRouterRouteCreate,
		ReadContext:   resourceIBMMetricsRouterRouteRead,
		UpdateContext: resourceIBMMetricsRouterRouteUpdate,
		DeleteContext: resourceIBMMetricsRouterRouteDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_metrics_router_route", "name"),
				Description:  "The name of the route. The name must be 1000 characters or less and cannot include any special characters other than `(space) - . _ :`. Do not include any personal identifying information (PII) in any resource names.",
			},
			"rules": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Routing rules that will be evaluated in their order of the array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "The action if the inclusion_filters matches, default is `send` action.",
							ValidateFunc: validate.InvokeValidator("ibm_metrics_router_route", "action"),
						},
						"targets": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A collection of targets with ID in the request.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The target uuid for a pre-defined metrics router target.",
										ValidateFunc: validate.InvokeValidator("ibm_metrics_router_route", "id"),
									},
								},
							},
						},
						"inclusion_filters": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of conditions to be satisfied for routing metrics to pre-defined target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operand": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										Description:  "Part of CRN that can be compared with values.",
										ValidateFunc: validate.InvokeValidator("ibm_metrics_router_route", "operand"),
									},
									"operator": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The operation to be performed between operand and the provided values. 'is' to be used with one value and 'in' can support upto 20 values in the array.",
										ValidateFunc: validate.InvokeValidator("ibm_metrics_router_route", "operator"),
									},
									"values": &schema.Schema{
										Type:        schema.TypeList,
										Required:    true,
										Description: "The provided string values of the operand to be compared with.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the route resource.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the route creation time.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the route last updated time.",
			},
		},
	}
}

func ResourceIBMMetricsRouterRouteValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 \-._:]+$`,
			MinValueLength:             1,
			MaxValueLength:             1000,
		},
		validate.ValidateSchema{
			Identifier:                 "id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 \-._:]+$`,
			MinValueLength:             3,
			MaxValueLength:             1000,
		},
		validate.ValidateSchema{
			Identifier:                 "operand",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "location, service_name, service_instance, resource_type, resource",
		},
		validate.ValidateSchema{
			Identifier:                 "operator",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "is, in",
		},
		validate.ValidateSchema{
			Identifier:                 "action",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "send, drop",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_metrics_router_route", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMMetricsRouterRouteCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	createRouteOptions := &metricsrouterv3.CreateRouteOptions{}

	createRouteOptions.SetName(d.Get("name").(string))
	var rules []metricsrouterv3.RulePrototype
	for _, e := range d.Get("rules").([]interface{}) {
		value := e.(map[string]interface{})
		rulesItem, err := resourceIBMMetricsRouterRouteMapToRulePrototype(value)
		if err != nil {
			return diag.FromErr(err)
		}
		rules = append(rules, *rulesItem)
	}
	createRouteOptions.SetRules(rules)

	route, response, err := metricsRouterClient.CreateRouteWithContext(context, createRouteOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateRouteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateRouteWithContext failed %s\n%s", err, response))
	}

	d.SetId(*route.ID)

	return resourceIBMMetricsRouterRouteRead(context, d, meta)
}

func resourceIBMMetricsRouterRouteRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getRouteOptions := &metricsrouterv3.GetRouteOptions{}

	getRouteOptions.SetID(d.Id())

	route, response, err := metricsRouterClient.GetRouteWithContext(context, getRouteOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetRouteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRouteWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("name", route.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	rules := []map[string]interface{}{}
	for _, rulesItem := range route.Rules {
		rulesItemMap, err := resourceIBMMetricsRouterRouteRulePrototypeToMap(&rulesItem)
		if err != nil {
			return diag.FromErr(err)
		}
		rules = append(rules, rulesItemMap)
	}
	if err = d.Set("rules", rules); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rules: %s", err))
	}
	if err = d.Set("crn", route.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(route.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(route.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	return nil
}

func resourceIBMMetricsRouterRouteUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	updateRouteOptions := &metricsrouterv3.UpdateRouteOptions{}

	updateRouteOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("name") || d.HasChange("rules") {
		updateRouteOptions.SetName(d.Get("name").(string))
		// TODO: handle Rules of type TypeList -- not primitive, not model

		var rules []metricsrouterv3.RulePrototype
		for _, e := range d.Get("rules").([]interface{}) {
			value := e.(map[string]interface{})
			rulesItem, err := resourceIBMMetricsRouterRouteMapToRulePrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			rules = append(rules, *rulesItem)
		}
		updateRouteOptions.SetRules(rules)
		hasChange = true
	}

	if hasChange {
		_, response, err := metricsRouterClient.UpdateRouteWithContext(context, updateRouteOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateRouteWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateRouteWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMMetricsRouterRouteRead(context, d, meta)
}

func resourceIBMMetricsRouterRouteDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	metricsRouterClient, err := meta.(conns.ClientSession).MetricsRouterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRouteOptions := &metricsrouterv3.DeleteRouteOptions{}

	deleteRouteOptions.SetID(d.Id())

	response, err := metricsRouterClient.DeleteRouteWithContext(context, deleteRouteOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteRouteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteRouteWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMMetricsRouterRouteMapToRulePrototype(modelMap map[string]interface{}) (*metricsrouterv3.RulePrototype, error) {
	model := &metricsrouterv3.RulePrototype{}
	if modelMap["action"] != nil && modelMap["action"].(string) != "" {
		model.Action = core.StringPtr(modelMap["action"].(string))
	}
	targets := []metricsrouterv3.TargetIdentity{}
	for _, targetsItem := range modelMap["targets"].([]interface{}) {
		targetsItemModel, err := resourceIBMMetricsRouterRouteMapToTargetIdentity(targetsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		targets = append(targets, *targetsItemModel)
	}
	model.Targets = targets
	inclusionFilters := []metricsrouterv3.InclusionFilterPrototype{}
	for _, inclusionFiltersItem := range modelMap["inclusion_filters"].([]interface{}) {
		inclusionFiltersItemModel, err := resourceIBMMetricsRouterRouteMapToInclusionFilterPrototype(inclusionFiltersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		inclusionFilters = append(inclusionFilters, *inclusionFiltersItemModel)
	}
	model.InclusionFilters = inclusionFilters
	return model, nil
}

func resourceIBMMetricsRouterRouteMapToTargetIdentity(modelMap map[string]interface{}) (*metricsrouterv3.TargetIdentity, error) {
	model := &metricsrouterv3.TargetIdentity{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMMetricsRouterRouteMapToInclusionFilterPrototype(modelMap map[string]interface{}) (*metricsrouterv3.InclusionFilterPrototype, error) {
	model := &metricsrouterv3.InclusionFilterPrototype{}
	model.Operand = core.StringPtr(modelMap["operand"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	values := []string{}
	for _, valuesItem := range modelMap["values"].([]interface{}) {
		values = append(values, valuesItem.(string))
	}
	model.Values = values
	return model, nil
}

func resourceIBMMetricsRouterRouteRulePrototypeToMap(model *metricsrouterv3.Rule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Action != nil {
		modelMap["action"] = model.Action
	}
	targets := []map[string]interface{}{}
	for _, targetsItem := range model.Targets {
		targetsItemMap, err := resourceIBMMetricsRouterRouteTargetIdentityToMap(&targetsItem)
		if err != nil {
			return modelMap, err
		}
		targets = append(targets, targetsItemMap)
	}
	modelMap["targets"] = targets
	inclusionFilters := []map[string]interface{}{}
	for _, inclusionFiltersItem := range model.InclusionFilters {
		inclusionFiltersItemMap, err := resourceIBMMetricsRouterRouteInclusionFilterPrototypeToMap(&inclusionFiltersItem)
		if err != nil {
			return modelMap, err
		}
		inclusionFilters = append(inclusionFilters, inclusionFiltersItemMap)
	}
	modelMap["inclusion_filters"] = inclusionFilters
	return modelMap, nil
}

func resourceIBMMetricsRouterRouteTargetIdentityToMap(model *metricsrouterv3.TargetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMMetricsRouterRouteInclusionFilterPrototypeToMap(model *metricsrouterv3.InclusionFilter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["operand"] = model.Operand
	modelMap["operator"] = model.Operator
	modelMap["values"] = model.Values
	return modelMap, nil
}
