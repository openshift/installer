// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package appconfiguration

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIbmAppConfigProperty() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmIbmAppConfigPropertyCreate,
		Read:     resourceIbmIbmAppConfigPropertyRead,
		Update:   resourceIbmIbmAppConfigPropertyUpdate,
		Delete:   resourceIbmIbmAppConfigPropertyDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Environment Id.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Property name.",
			},
			"property_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Property id.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the Property  (BOOLEAN, STRING, NUMERIC).",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value of the Property. The value can be Boolean, String or a Numeric value as per the `type` attribute.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Property description.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags associated with the property.",
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Format of the feature (TEXT, JSON, YAML).",
			},
			"segment_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the targeting rules that is used to set different property values for different segments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rules": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Rules array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segments": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of segment ids that are used for targeting using the rule.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value to be used for evaluation for this rule. The value can be Boolean, String or a Numeric value as per the `type` attribute.",
						},
						"order": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Order of the rule, used during evaluation. The evaluation is performed in the order defined and the value associated with the first matching rule is used for evaluation.",
						},
					},
				},
			},
			"collections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of collection id representing the collections that are associated with the specified property.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Collection id.",
						},
					},
				},
			},
			"segment_exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes if the targeting rules are specified for the property.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the property.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of the property data.",
			},
			"evaluation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last occurrence of the property value evaluation.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property URL.",
			},
		},
	}
}

func resourceIbmIbmAppConfigPropertyCreate(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)
	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.CreatePropertyOptions{}

	options.SetName(d.Get("name").(string))
	options.SetType(d.Get("type").(string))
	options.SetEnvironmentID(d.Get("environment_id").(string))
	options.SetPropertyID(d.Get("property_id").(string))
	options.SetValue(d.Get("value").(string))

	if _, ok := GetFieldExists(d, "description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := GetFieldExists(d, "tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}
	if _, ok := GetFieldExists(d, "format"); ok {
		options.SetFormat(d.Get("format").(string))
	}
	if _, ok := GetFieldExists(d, "collections"); ok {
		var collections []appconfigurationv1.CollectionRef
		for _, e := range d.Get("collections").([]interface{}) {
			value := e.(map[string]interface{})
			collectionsItem := resourceIbmAppConfigPropertyMapToCollectionRef(value)
			collections = append(collections, collectionsItem)
		}
		options.SetCollections(collections)
	}
	if _, ok := GetFieldExists(d, "segment_rules"); ok {
		var segmentRules []appconfigurationv1.SegmentRule
		for _, e := range d.Get("segment_rules").([]interface{}) {
			value := e.(map[string]interface{})
			segmentRulesItem, err := resourceIbmAppConfigPropertyMapToSegmentRule(d, value)
			if err != nil {
				return flex.FmtErrorf(fmt.Sprintf("%s", err))
			}
			segmentRules = append(segmentRules, segmentRulesItem)
		}
		options.SetSegmentRules(segmentRules)
	}

	result, response, err := appconfigClient.CreateProperty(options)
	if err != nil {
		return flex.FmtErrorf("CreateProperty failed %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", guid, *options.EnvironmentID, *result.PropertyID))

	return resourceIbmIbmAppConfigPropertyRead(d, meta)
}

func resourceIbmIbmAppConfigPropertyRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	if len(parts) != 3 {
		return flex.FmtErrorf("Kindly check the id")
	}

	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.GetPropertyOptions{}

	options.SetEnvironmentID(parts[1])
	options.SetPropertyID(parts[2])

	result, response, err := appconfigClient.GetProperty(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
		}
		return flex.FmtErrorf("[ERROR] GetProperty failed %s\n%s", err, response)
	}

	d.Set("guid", parts[0])
	d.Set("environment_id", parts[1])

	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return flex.FmtErrorf("error setting name: %s", err)
		}
	}
	if result.PropertyID != nil {
		if err = d.Set("property_id", result.PropertyID); err != nil {
			return flex.FmtErrorf("error setting property_id: %s", err)
		}
	}
	if result.Type != nil {
		if err = d.Set("type", result.Type); err != nil {
			return flex.FmtErrorf("error setting type: %s", err)
		}
	}
	if result.Value != nil {
		value := result.Value
		switch value.(interface{}).(type) {
		case string:
			d.Set("value", value.(string))
		case float64:
			d.Set("value", fmt.Sprintf("%v", value))
		case bool:
			d.Set("value", strconv.FormatBool(value.(bool)))
		}
	}
	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return flex.FmtErrorf("error setting description: %s", err)
		}
	}
	if result.Tags != nil {
		if err = d.Set("tags", result.Tags); err != nil {
			return flex.FmtErrorf("error setting tags: %s", err)
		}
	}
	if result.SegmentExists != nil {
		if err = d.Set("segment_exists", result.SegmentExists); err != nil {
			return flex.FmtErrorf("error setting segment_exists: %s", err)
		}
	}
	if result.CreatedTime != nil {
		if err = d.Set("created_time", result.CreatedTime.String()); err != nil {
			return flex.FmtErrorf("error setting created_time: %s", err)
		}
	}
	if result.UpdatedTime != nil {
		if err = d.Set("updated_time", result.UpdatedTime.String()); err != nil {
			return flex.FmtErrorf("error setting updated_time: %s", err)
		}
	}
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return flex.FmtErrorf("error setting href: %s", err)
		}
	}

	if result.SegmentRules != nil {
		segmentRules := []map[string]interface{}{}
		for _, segmentRulesItem := range result.SegmentRules {
			segmentRulesItemMap := resourceIbmAppConfigPropertySegmentRuleToMap(segmentRulesItem)
			segmentRules = append(segmentRules, segmentRulesItemMap)
		}
		if err = d.Set("segment_rules", segmentRules); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting segment_rules: %s", err)
		}
	}
	if result.Collections != nil {
		collections := []map[string]interface{}{}
		for _, collectionsItem := range result.Collections {
			collectionsItemMap := resourceIbmAppConfigPropertyCollectionRefToMap(collectionsItem)
			collections = append(collections, collectionsItemMap)
		}
		if err = d.Set("collections", collections); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting collections: %s", err)
		}
	}
	return nil
}

func resourceIbmIbmAppConfigPropertyUpdate(d *schema.ResourceData, meta interface{}) error {
	if ok := d.HasChanges("name", "value", "description", "tags", "segment_rules", "collections"); ok {
		parts, err := flex.IdParts(d.Id())
		if err != nil {
			return nil
		}
		appconfigClient, err := getAppConfigClient(meta, parts[0])
		if err != nil {
			return flex.FmtErrorf("getAppConfigClient failed %s", err)
		}
		options := &appconfigurationv1.UpdatePropertyOptions{}

		options.SetEnvironmentID(parts[1])
		options.SetPropertyID(parts[2])

		options.SetName(d.Get("name").(string))
		options.SetValue(d.Get("value").(string))

		if _, ok := GetFieldExists(d, "description"); ok {
			options.SetDescription(d.Get("description").(string))
		}
		if _, ok := GetFieldExists(d, "tags"); ok {
			options.SetTags(d.Get("tags").(string))
		}
		if _, ok := GetFieldExists(d, "collections"); ok {
			var collections []appconfigurationv1.CollectionRef
			for _, e := range d.Get("collections").([]interface{}) {
				value := e.(map[string]interface{})
				collectionsItem := resourceIbmAppConfigPropertyMapToCollectionRef(value)
				collections = append(collections, collectionsItem)
			}
			options.SetCollections(collections)
		}
		if _, ok := GetFieldExists(d, "segment_rules"); ok {
			var segmentRules []appconfigurationv1.SegmentRule
			for _, e := range d.Get("segment_rules").([]interface{}) {
				value := e.(map[string]interface{})
				segmentRulesItem, err := resourceIbmAppConfigPropertyMapToSegmentRule(d, value)
				if err != nil {
					return flex.FmtErrorf(fmt.Sprintf("%s", err))
				}
				segmentRules = append(segmentRules, segmentRulesItem)
			}
			options.SetSegmentRules(segmentRules)
		}
		_, response, err := appconfigClient.UpdateProperty(options)
		if err != nil {
			return flex.FmtErrorf("UpdateProperty failed %s\n%s", err, response)
		}

		return resourceIbmIbmAppConfigPropertyRead(d, meta)
	}
	return nil
}

func resourceIbmIbmAppConfigPropertyDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.DeletePropertyOptions{}

	options.SetEnvironmentID(parts[1])
	options.SetPropertyID(parts[2])

	response, err := appconfigClient.DeleteProperty(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return flex.FmtErrorf("[ERROR] DeleteProperty failed %s\n%s", err, response)
	}

	d.SetId("")

	return nil
}

func resourceIbmAppConfigPropertyMapToCollectionRef(collectionRefMap map[string]interface{}) appconfigurationv1.CollectionRef {
	collectionRef := appconfigurationv1.CollectionRef{}
	collectionRef.CollectionID = core.StringPtr(collectionRefMap["collection_id"].(string))
	return collectionRef
}

func resourceIbmAppConfigPropertyMapToSegmentRule(d *schema.ResourceData, segmentRuleMap map[string]interface{}) (appconfigurationv1.SegmentRule, error) {
	segmentRule := appconfigurationv1.SegmentRule{}

	rules := []appconfigurationv1.TargetSegments{}
	for _, rulesItem := range segmentRuleMap["rules"].([]interface{}) {
		rulesItemModel := resourceIbmAppConfigPropertyMapToRule(rulesItem.(map[string]interface{}))
		rules = append(rules, rulesItemModel)
	}
	segmentRule.Rules = rules

	segmentRule.Order = core.Int64Ptr(int64(segmentRuleMap["order"].(int)))
	ruleValue := segmentRuleMap["value"].(string)
	switch d.Get("type").(string) {
	case "STRING":
		segmentRule.Value = ruleValue
	case "NUMERIC":
		v, err := strconv.ParseFloat(ruleValue, 64)
		if err != nil {
			return segmentRule, flex.FmtErrorf("'value' parameter in 'segment_rules' has wrong value: %s", err)
		}
		segmentRule.Value = v
	case "BOOLEAN":
		if ruleValue == "false" {
			segmentRule.Value = false
		} else if ruleValue == "true" {
			segmentRule.Value = true
		} else {
			return segmentRule, flex.FmtErrorf("'value' parameter in 'segment_rules' has wrong value")
		}
	}

	return segmentRule, nil
}

func resourceIbmAppConfigPropertyMapToRule(ruleMap map[string]interface{}) appconfigurationv1.TargetSegments {
	rule := appconfigurationv1.TargetSegments{}

	segments := []string{}
	for _, segmentsItem := range ruleMap["segments"].([]interface{}) {
		segments = append(segments, segmentsItem.(string))
	}
	rule.Segments = segments

	return rule
}

// output
func resourceIbmAppConfigPropertySegmentRuleToMap(segmentRule appconfigurationv1.SegmentRule) map[string]interface{} {
	segmentRuleMap := map[string]interface{}{}

	rules := []map[string]interface{}{}
	for _, rulesItem := range segmentRule.Rules {
		rulesItemMap := resourceIbmAppConfigPropertyRuleToMap(rulesItem)
		rules = append(rules, rulesItemMap)
	}

	segmentRuleMap["rules"] = rules
	segmentRuleMap["order"] = flex.IntValue(segmentRule.Order)
	segmentValue := segmentRule.Value
	switch segmentValue.(interface{}).(type) {
	case string:
		segmentRuleMap["value"] = segmentValue.(string)
	case float64:
		segmentRuleMap["value"] = fmt.Sprintf("%v", segmentValue)
	case bool:
		segmentRuleMap["value"] = strconv.FormatBool(segmentValue.(bool))
	}

	return segmentRuleMap
}

func resourceIbmAppConfigPropertyRuleToMap(rule appconfigurationv1.TargetSegments) map[string]interface{} {
	ruleMap := map[string]interface{}{}
	ruleMap["segments"] = rule.Segments
	return ruleMap
}

func resourceIbmAppConfigPropertyCollectionRefToMap(collectionRef appconfigurationv1.CollectionRef) map[string]interface{} {
	collectionRefMap := map[string]interface{}{}
	collectionRefMap["collection_id"] = collectionRef.CollectionID
	collectionRefMap["name"] = collectionRef.Name
	return collectionRefMap
}
