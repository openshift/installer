// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package atracker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/atrackerv1"
	"github.com/IBM/platform-services-go-sdk/atrackerv2"
)

func ResourceIBMAtrackerRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAtrackerRouteCreate,
		ReadContext:   resourceIBMAtrackerRouteRead,
		UpdateContext: resourceIBMAtrackerRouteUpdate,
		DeleteContext: resourceIBMAtrackerRouteDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_atracker_route", "name"),
				Description:  "The name of the route. The name must be 1000 characters or less and cannot include any special characters other than `(space) - . _ :`.",
			},
			"receive_global_events": {
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "use rules.locations instead",
				Description: "Indicates whether or not all global events should be forwarded to this region.",
			},
			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Routing rules that will be evaluated in their order of the array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_ids": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "The target ID List. All the events will be send to all targets listed in the rule. You can include targets from other regions.  ",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"locations": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Logs from these locations will be sent to the targets specified. Locations is a superset of regions including global and *.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the route resource.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the route.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "use created_at instead",
				Description: "The timestamp of the route creation time.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the route creation time.",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "use updated_at instead",
				Description: "The timestamp of the route last updated time.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the route last updated time.",
			},
			"api_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The API version of the route.",
			},
		},
	}
}

func ResourceIBMAtrackerRouteValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 -._:]+$`,
			MinValueLength:             1,
			MaxValueLength:             1000,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_atracker_route", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMAtrackerRouteCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createRouteOptions := &atrackerv2.CreateRouteOptions{}

	createRouteOptions.SetName(d.Get("name").(string))
	var rules []atrackerv2.RulePrototype
	for _, e := range d.Get("rules").([]interface{}) {
		value := e.(map[string]interface{})
		rulesItem := resourceIBMAtrackerRouteMapToRule(value, d.Get("receive_global_events").(bool))
		rules = append(rules, rulesItem)
	}

	createRouteOptions.SetRules(rules)

	route, response, err := atrackerClient.CreateRouteWithContext(context, createRouteOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateRouteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateRouteWithContext failed %s\n%s", err, response))
	}

	d.SetId(*route.ID)
	d.Set("api_version", 2)

	return resourceIBMAtrackerRouteRead(context, d, meta)
}

func resourceIBMAtrackerRouteRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClientv1, atrackerClient, err := getAtrackerClients(meta)
	// We need both route methods to ensure backwards compatibility and the ability to migrate
	if err != nil {
		return diag.FromErr(err)
	}

	getRouteOptions := &atrackerv2.GetRouteOptions{}
	getRouteOptions.SetID(d.Id())
	apiVersion := d.Get("api_version")

	// Try v2 first, otherwise try v1
	route, response, err := atrackerClient.GetRouteWithContext(context, getRouteOptions)

	if err != nil && response != nil && response.StatusCode != 404 {
		log.Printf("[DEBUG] GetRouteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRouteWithContext failed %s\n%s", err, response))
	}
	if err == nil && response != nil {
		if err = d.Set("name", route.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
		rules := make([]map[string]interface{}, len(route.Rules), len(route.Rules))
		for i, rulesItem := range route.Rules {
			rulesItemMap, _, _ := resourceIBMAtrackerRouteRulePrototypeToMap(&rulesItem)
			rules[i] = rulesItemMap
		}
		if err = d.Set("rules", rules); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting rules: %s", err))
		}
		if err = d.Set("crn", route.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
		}
		if err = d.Set("version", flex.IntValue(route.Version)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
		}
		if err = d.Set("created_at", flex.DateTimeToString(route.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
		if err = d.Set("updated_at", flex.DateTimeToString(route.UpdatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
		}
		if err = d.Set("api_version", flex.IntValue(route.APIVersion)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting api_version: %s", err))
		}
		d.Set("receive_global_events", false)
	} else if apiVersion != 2 {
		getRouteV1Options := &atrackerv1.GetRouteOptions{}
		getRouteV1Options.SetID(d.Id())
		routeV1, responseV1, err := atrackerClientv1.GetRouteWithContext(context, getRouteV1Options)
		if err != nil {
			if response != nil && responseV1.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			log.Printf("[DEBUG] GetRouteWithContext failed %s\n%s", err, responseV1)
			return diag.FromErr(fmt.Errorf("GetRouteWithContext failed %s\n%s", err, responseV1))
		}
		if err = d.Set("name", routeV1.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
		rules := []map[string]interface{}{}
		for _, rulesItem := range routeV1.Rules {
			rulesItemMap, _ := resourceIBMAtrackerRouteRulePrototypeToMapV1(&rulesItem)
			rules = append(rules, rulesItemMap)
		}
		if err = d.Set("rules", rules); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting rules: %s", err))
		}
		if err = d.Set("crn", routeV1.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
		}
		if err = d.Set("receive_global_events", routeV1.ReceiveGlobalEvents); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
		}
		if err = d.Set("version", flex.IntValue(routeV1.Version)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
		}
		if err = d.Set("created_at", flex.DateTimeToString(routeV1.Created)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
		if err = d.Set("updated_at", flex.DateTimeToString(routeV1.Updated)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
		}
		if err = d.Set("api_version", 1); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting api_version: %s", err))
		}
	}

	return nil
}

func resourceIBMAtrackerRouteUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClientV1, atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	apiVersion := d.Get("api_version").(int)

	if apiVersion > 1 {
		replaceRouteOptions := &atrackerv2.ReplaceRouteOptions{}

		replaceRouteOptions.SetID(d.Id())
		replaceRouteOptions.SetName(d.Get("name").(string))

		var rules []atrackerv2.RulePrototype = make([]atrackerv2.RulePrototype, 0)
		for _, e := range d.Get("rules").([]interface{}) {
			value := e.(map[string]interface{})
			rulesItem := resourceIBMAtrackerRouteMapToRule(value, d.Get("receive_global_events").(bool))
			rules = append(rules, rulesItem)
		}
		replaceRouteOptions.SetRules(rules)

		_, response, err := atrackerClient.ReplaceRouteWithContext(context, replaceRouteOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceRouteWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplaceRouteWithContext failed %s\n%s", err, response))
		}
		return resourceIBMAtrackerRouteRead(context, d, meta)
	}
	// TODO: to remove once version 1 is fully deprecated
	replaceRouteOptionsV1 := &atrackerv1.ReplaceRouteOptions{}
	replaceRouteOptionsV1.SetID(d.Id())
	replaceRouteOptionsV1.SetName(d.Get("name").(string))
	replaceRouteOptionsV1.SetReceiveGlobalEvents(d.Get("receive_global_events").(bool))

	var rules []atrackerv1.Rule = make([]atrackerv1.Rule, 0)
	for _, e := range d.Get("rules").([]interface{}) {
		value := e.(map[string]interface{})
		rulesItem := resourceIBMAtrackerRouteMapToRuleV1(value)
		rules = append(rules, rulesItem)
	}
	replaceRouteOptionsV1.SetRules(rules)

	_, response, err := atrackerClientV1.ReplaceRouteWithContext(context, replaceRouteOptionsV1)
	if err != nil {
		log.Printf("[DEBUG] ReplaceRouteWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ReplaceRouteWithContext failed %s\n%s", err, response))
	}

	return resourceIBMAtrackerRouteRead(context, d, meta)
}

func resourceIBMAtrackerRouteDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClientV1, atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	apiVersion := d.Get("api_version").(int)

	if apiVersion > 1 {
		deleteRouteOptions := &atrackerv2.DeleteRouteOptions{}

		deleteRouteOptions.SetID(d.Id())

		response, err := atrackerClient.DeleteRouteWithContext(context, deleteRouteOptions)
		if err != nil {
			log.Printf("[DEBUG] DeleteRouteWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("DeleteRouteWithContext failed %s\n%s", err, response))
		}
	} else {
		deleteRouteOptions := &atrackerv1.DeleteRouteOptions{}

		deleteRouteOptions.SetID(d.Id())

		response, err := atrackerClientV1.DeleteRouteWithContext(context, deleteRouteOptions)
		if err != nil {
			log.Printf("[DEBUG] DeleteRouteWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("DeleteRouteWithContext failed %s\n%s", err, response))
		}
	}

	d.SetId("")

	return nil
}

func resourceIBMAtrackerRouteRulePrototypeToMap(ruleModel *atrackerv2.Rule) (map[string]interface{}, bool, error) {
	receives_global_events := false
	ruleMap := make(map[string]interface{})
	if ruleModel != nil {
		ruleMap["target_ids"] = make([]string, len(ruleModel.TargetIds))
		if ruleModel.TargetIds != nil {
			for i, target_id := range ruleModel.TargetIds {
				ruleMap["target_ids"].([]string)[i] = target_id
			}
		}

		ruleMap["locations"] = make([]string, len(ruleModel.Locations))
		if ruleModel.Locations != nil {
			for i, location := range ruleModel.Locations {
				ruleMap["locations"].([]string)[i] = location
				if strings.Contains(location, "*") || strings.Contains(location, "global") {
					receives_global_events = true
				}
			}
		}
		return ruleMap, receives_global_events, nil
	}
	return ruleMap, false, nil
}

func resourceIBMAtrackerRouteRulePrototypeToMapV1(model *atrackerv1.Rule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_ids"] = model.TargetIds
	return modelMap, nil
}

func resourceIBMAtrackerRouteMapToRule(ruleMap map[string]interface{}, addGlobalFlag bool) atrackerv2.RulePrototype {
	rule := atrackerv2.RulePrototype{}

	targetIds := make([]string, 0)
	for _, targetIdsItem := range ruleMap["target_ids"].(*schema.Set).List() {
		if targetIdsItem != nil {
			targetIds = append(targetIds, targetIdsItem.(string))
		}
	}
	rule.TargetIds = targetIds

	locations := make([]string, 0)
	globalDetected := false
	for _, locationsItem := range ruleMap["locations"].(*schema.Set).List() {
		if strings.Contains(locationsItem.(string), "*") || strings.Contains(locationsItem.(string), "global") {
			globalDetected = true
		}
		locations = append(locations, locationsItem.(string))
	}

	if addGlobalFlag && !globalDetected {
		locations = append(locations, "global")
	}
	rule.Locations = locations

	return rule
}

func resourceIBMAtrackerRouteMapToRuleV1(ruleMap map[string]interface{}) atrackerv1.Rule {
	rule := atrackerv1.Rule{}

	targetIds := []string{}
	for _, targetIdsItem := range ruleMap["target_ids"].(*schema.Set).List() {
		targetIds = append(targetIds, targetIdsItem.(string))
	}
	rule.TargetIds = targetIds

	return rule
}
