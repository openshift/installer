// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func resourceIbmIbmAppConfigFeature() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmIbmAppConfigFeatureCreate,
		Read:     resourceIbmIbmAppConfigFeatureRead,
		Update:   resourceIbmIbmAppConfigFeatureUpdate,
		Delete:   resourceIbmIbmAppConfigFeatureDelete,
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
				Description: "Feature name.",
			},
			"feature_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Feature id.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_app_config_feature", "type"),
				Description:  "Type of the feature (BOOLEAN, STRING, NUMERIC).",
			},
			"enabled_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value of the feature when it is enabled. The value can be BOOLEAN, STRING or a NUMERIC value as per the `type` attribute.",
			},
			"disabled_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value of the feature when it is disabled. The value can be BOOLEAN, STRING or a NUMERIC value as per the `type` attribute.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Feature description.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags associated with the feature.",
			},
			"segment_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify the targeting rules that is used to set different feature flag values for different segments.",
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
				Description: "List of collection id representing the collections that are associated with the specified feature flag.",
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
				Description: "Denotes if the targeting rules are specified for the feature flag.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The state of the feature flag.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the feature flag.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of the feature flag data.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Feature flag URL.",
			},
		},
	}
}

func resourceIbmIbmAppConfigFeatureCreate(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)
	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return err
	}
	options := &appconfigurationv1.CreateFeatureOptions{}
	options.SetType(d.Get("type").(string))
	options.SetName(d.Get("name").(string))
	options.SetFeatureID(d.Get("feature_id").(string))
	options.SetEnabledValue(d.Get("enabled_value").(string))
	options.SetEnvironmentID(d.Get("environment_id").(string))
	options.SetDisabledValue(d.Get("disabled_value").(string))

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}

	if _, ok := d.GetOk("segment_rules"); ok {
		var segmentRules []appconfigurationv1.SegmentRule
		for _, e := range d.Get("segment_rules").([]interface{}) {
			value := e.(map[string]interface{})
			segmentRulesItem, err := resourceIbmAppConfigFeatureMapToSegmentRule(d, value)
			if err != nil {
				return err
			}
			segmentRules = append(segmentRules, segmentRulesItem)
		}
		options.SetSegmentRules(segmentRules)
	}
	if _, ok := d.GetOk("collections"); ok {
		var collections []appconfigurationv1.CollectionRef
		for _, e := range d.Get("collections").([]interface{}) {
			value := e.(map[string]interface{})
			collectionsItem := resourceIbmAppConfigFeatureMapToCollectionRef(value)
			collections = append(collections, collectionsItem)
		}
		options.SetCollections(collections)
	}

	feature, response, err := appconfigClient.CreateFeature(options)

	if err != nil {
		log.Printf("CreateFeature failed %s\n%s", err, response)
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", guid, *options.EnvironmentID, *feature.FeatureID))
	return resourceIbmIbmAppConfigFeatureRead(d, meta)
}

func resourceIbmIbmAppConfigFeatureUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := idParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return err
	}

	options := &appconfigurationv1.UpdateFeatureOptions{}
	options.SetEnvironmentID(parts[1])
	options.SetFeatureID(parts[2])

	if ok := d.HasChanges("name", "enabled_value", "disabled_value", "description", "tags", "segment_rules", "collections"); ok {
		options.SetName(d.Get("name").(string))
		options.SetEnabledValue(d.Get("enabled_value").(string))
		options.SetDisabledValue(d.Get("disabled_value").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}
		if _, ok := d.GetOk("tags"); ok {
			options.SetTags(d.Get("tags").(string))
		}
		if _, ok := d.GetOk("segment_rules"); ok {
			var segmentRules []appconfigurationv1.SegmentRule
			for _, e := range d.Get("segment_rules").([]interface{}) {
				value := e.(map[string]interface{})
				segmentRulesItem, err := resourceIbmAppConfigFeatureMapToSegmentRule(d, value)
				if err != nil {
					return err
				}
				segmentRules = append(segmentRules, segmentRulesItem)
			}
			options.SetSegmentRules(segmentRules)
		}
		if _, ok := d.GetOk("collections"); ok {
			var collections []appconfigurationv1.CollectionRef
			for _, e := range d.Get("collections").([]interface{}) {
				value := e.(map[string]interface{})
				collectionsItem := resourceIbmAppConfigFeatureMapToCollectionRef(value)
				collections = append(collections, collectionsItem)
			}
			options.SetCollections(collections)
		}

		_, response, err := appconfigClient.UpdateFeature(options)
		if err != nil {
			log.Printf("[DEBUG] UpdateFeature %s\n%s", err, response)
			return err
		}
		return resourceIbmIbmAppConfigFeatureRead(d, meta)
	}
	return nil
}

func resourceIbmIbmAppConfigFeatureRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := idParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return err
	}

	options := &appconfigurationv1.GetFeatureOptions{}
	options.SetEnvironmentID(parts[1])
	options.SetFeatureID(parts[2])

	result, response, err := appconfigClient.GetFeature(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] GetFeature failed %s\n%s", err, response)
	}

	d.Set("guid", parts[0])
	d.Set("environment_id", parts[1])
	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return fmt.Errorf("error setting name: %s", err)
		}
	}
	if result.FeatureID != nil {
		if err = d.Set("feature_id", result.FeatureID); err != nil {
			return fmt.Errorf("error setting feature_id: %s", err)
		}
	}
	if result.Type != nil {
		if err = d.Set("type", result.Type); err != nil {
			return fmt.Errorf("error setting type: %s", err)
		}
	}
	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return fmt.Errorf("error setting description: %s", err)
		}

	}
	if result.Tags != nil {
		if err = d.Set("tags", result.Tags); err != nil {
			return fmt.Errorf("error setting tags: %s", err)
		}
	}

	if result.SegmentRules != nil {
		segmentRules := []map[string]interface{}{}
		for _, segmentRulesItem := range result.SegmentRules {
			segmentRulesItemMap := resourceIbmAppConfigFeatureSegmentRuleToMap(segmentRulesItem)
			segmentRules = append(segmentRules, segmentRulesItemMap)
		}
		if err = d.Set("segment_rules", segmentRules); err != nil {
			return fmt.Errorf("error setting segment_rules: %s", err)
		}
	}
	if result.Collections != nil {
		collections := []map[string]interface{}{}
		for _, collectionsItem := range result.Collections {
			collectionsItemMap := resourceIbmAppConfigFeatureCollectionRefToMap(collectionsItem)
			collections = append(collections, collectionsItemMap)
		}
		if err = d.Set("collections", collections); err != nil {
			return fmt.Errorf("error setting collections: %s", err)
		}
	}
	if result.SegmentExists != nil {
		if err = d.Set("segment_exists", result.SegmentExists); err != nil {
			return fmt.Errorf("error setting segment_exists: %s", err)
		}
	}
	if result.CreatedTime != nil {
		if err = d.Set("created_time", result.CreatedTime.String()); err != nil {
			return fmt.Errorf("error setting created_time: %s", err)
		}
	}
	if result.UpdatedTime != nil {
		if err = d.Set("updated_time", result.UpdatedTime.String()); err != nil {
			return fmt.Errorf("error setting updated_time: %s", err)
		}
	}
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return fmt.Errorf("error setting href: %s", err)
		}
	}
	if result.Enabled != nil {
		if err = d.Set("enabled", result.Enabled); err != nil {
			return fmt.Errorf("error setting enabled: %s", err)
		}
	}

	if result.EnabledValue != nil {
		enabledValue := result.EnabledValue

		switch enabledValue.(interface{}).(type) {
		case string:
			d.Set("enabled_value", enabledValue.(string))
		case float64:
			d.Set("enabled_value", fmt.Sprintf("%v", enabledValue))
		case bool:
			d.Set("enabled_value", strconv.FormatBool(enabledValue.(bool)))
		}
	}

	if result.DisabledValue != nil {
		disabledValue := result.DisabledValue

		switch disabledValue.(interface{}).(type) {
		case string:
			d.Set("disabled_value", disabledValue.(string))
		case float64:
			d.Set("disabled_value", fmt.Sprintf("%v", disabledValue))
		case bool:
			d.Set("disabled_value", strconv.FormatBool(disabledValue.(bool)))
		}
	}
	return nil
}

func resourceIbmIbmAppConfigFeatureDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := idParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return err
	}

	options := &appconfigurationv1.DeleteFeatureOptions{}
	options.SetEnvironmentID(parts[1])
	options.SetFeatureID(parts[2])

	response, err := appconfigClient.DeleteFeature(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[DEBUG] DeleteFeature failed %s\n%s", err, response)
	}

	d.SetId("")

	return nil
}

func resourceIbmAppConfigFeatureValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "BOOLEAN, NUMERIC, STRING",
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_app_config_feature", Schema: validateSchema}
	return &resourceValidator
}

// output
func resourceIbmAppConfigFeatureSegmentRuleToMap(segmentRule appconfigurationv1.SegmentRule) map[string]interface{} {
	segmentRuleMap := map[string]interface{}{}

	rules := []map[string]interface{}{}
	for _, rulesItem := range segmentRule.Rules {
		rulesItemMap := resourceIbmAppConfigFeatureRuleToMap(rulesItem)
		rules = append(rules, rulesItemMap)
	}

	segmentRuleMap["rules"] = rules
	segmentRuleMap["order"] = intValue(segmentRule.Order)

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

func resourceIbmAppConfigFeatureRuleToMap(rule appconfigurationv1.TargetSegments) map[string]interface{} {
	ruleMap := map[string]interface{}{}
	ruleMap["segments"] = rule.Segments
	return ruleMap
}

func resourceIbmAppConfigFeatureCollectionRefToMap(collectionRef appconfigurationv1.CollectionRef) map[string]interface{} {
	collectionRefMap := map[string]interface{}{}
	collectionRefMap["collection_id"] = collectionRef.CollectionID
	collectionRefMap["name"] = collectionRef.Name
	return collectionRefMap
}

// input
func resourceIbmAppConfigFeatureMapToSegmentRule(d *schema.ResourceData, segmentRuleMap map[string]interface{}) (appconfigurationv1.SegmentRule, error) {
	segmentRule := appconfigurationv1.SegmentRule{}

	rules := []appconfigurationv1.TargetSegments{}
	for _, rulesItem := range segmentRuleMap["rules"].([]interface{}) {
		rulesItemModel := resourceIbmAppConfigFeatureMapToRule(rulesItem.(map[string]interface{}))
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
			return segmentRule, fmt.Errorf("'value' parameter in 'segment_rules' has wrong value: %s", err)
		}
		segmentRule.Value = v
	case "BOOLEAN":
		if ruleValue == "false" {
			segmentRule.Value = false
		} else if ruleValue == "true" {
			segmentRule.Value = true
		} else {
			return segmentRule, fmt.Errorf("'value' parameter in 'segment_rules' has wrong value")
		}
	}

	return segmentRule, nil
}

func resourceIbmAppConfigFeatureMapToRule(ruleMap map[string]interface{}) appconfigurationv1.TargetSegments {
	rule := appconfigurationv1.TargetSegments{}

	segments := []string{}
	for _, segmentsItem := range ruleMap["segments"].([]interface{}) {
		segments = append(segments, segmentsItem.(string))
	}
	rule.Segments = segments

	return rule
}

func resourceIbmAppConfigFeatureMapToCollectionRef(collectionRefMap map[string]interface{}) appconfigurationv1.CollectionRef {
	collectionRef := appconfigurationv1.CollectionRef{}
	collectionRef.CollectionID = core.StringPtr(collectionRefMap["collection_id"].(string))

	return collectionRef
}
