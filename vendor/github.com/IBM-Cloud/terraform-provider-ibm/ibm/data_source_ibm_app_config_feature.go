// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
)

func dataSourceIbmAppConfigFeature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigFeatureRead,

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
			"feature_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Feature Id.",
			},
			"includes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Include the associated collections in the response.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Feature name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Feature description.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the feature (BOOLEAN, STRING, NUMERIC).",
			},
			"enabled_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the feature when it is enabled. The value can be Boolean, String or a Numeric value as per the `type` attribute.",
			},
			"disabled_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the feature when it is disabled. The value can be Boolean, String or a Numeric value as per the `type` attribute.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The state of the feature flag.",
			},
			"tags": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tags associated with the feature.",
			},
			"segment_exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes if the targeting rules are specified for the feature flag.",
			},
			"segment_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specify the targeting rules that is used to set different feature flag values for different segments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rules array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"segments": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of segment ids that are used for targeting using the rule.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value to be used for evaluation for this rule. The value can be Boolean, String or a Numeric value as per the `type` attribute.",
						},
						"order": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Order of the rule, used during evaluation. The evaluation is performed in the order defined and the value associated with the first matching rule is used for evaluation.",
						},
					},
				},
			},
			"collections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of collection id representing the collections that are associated with the specified feature flag.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Collection id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the collection.",
						},
					},
				},
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

func dataSourceIbmAppConfigFeatureRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return err
	}

	options := &appconfigurationv1.GetFeatureOptions{}

	options.SetEnvironmentID(d.Get("environment_id").(string))
	options.SetFeatureID(d.Get("feature_id").(string))

	if _, ok := d.GetOk("includes"); ok {
		options.SetInclude(d.Get("includes").(string))
	}

	result, response, err := appconfigClient.GetFeature(options)
	if err != nil {
		log.Printf("[DEBUG] GetFeature failed %s\n%s", err, response)
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", guid, *options.EnvironmentID, *result.FeatureID))
	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return fmt.Errorf("error setting name: %s", err)
		}
	}
	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return fmt.Errorf("error setting description: %s", err)
		}
	}
	if result.Type != nil {
		if err = d.Set("type", result.Type); err != nil {
			return fmt.Errorf("error setting type: %s", err)
		}
	}
	if result.Enabled != nil {
		if err = d.Set("enabled", result.Enabled); err != nil {
			return fmt.Errorf("error setting enabled: %s", err)
		}
	}
	if result.Tags != nil {
		if err = d.Set("tags", result.Tags); err != nil {
			return fmt.Errorf("error setting tags: %s", err)
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

	if result.SegmentRules != nil {
		err = d.Set("segment_rules", dataSourceFeatureFlattenSegmentRules(result.SegmentRules))
		if err != nil {
			return fmt.Errorf("error setting segment_rules %s", err)
		}
	}

	if result.Collections != nil {
		err = d.Set("collections", dataSourceFeatureFlattenCollections(result.Collections))
		if err != nil {
			return fmt.Errorf("error setting collections %s", err)
		}
	}
	return nil
}

func dataSourceFeatureFlattenSegmentRules(result []appconfigurationv1.SegmentRule) (segmentRules []map[string]interface{}) {
	for _, segmentRulesItem := range result {
		segmentRules = append(segmentRules, dataSourceFeatureSegmentRulesToMap(segmentRulesItem))
	}

	return segmentRules
}

func dataSourceFeatureSegmentRulesToMap(segmentRulesItem appconfigurationv1.SegmentRule) (segmentRulesMap map[string]interface{}) {
	segmentRulesMap = map[string]interface{}{}

	if segmentRulesItem.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, rulesItem := range segmentRulesItem.Rules {
			rulesList = append(rulesList, dataSourceFeatureSegmentRulesRulesToMap(rulesItem))
		}
		segmentRulesMap["rules"] = rulesList
	}
	if segmentRulesItem.Value != nil {
		segmentValue := segmentRulesItem.Value
		switch segmentValue.(interface{}).(type) {
		case string:
			segmentRulesMap["value"] = segmentValue.(string)
		case float64:
			segmentRulesMap["value"] = fmt.Sprintf("%v", segmentValue)
		case bool:
			segmentRulesMap["value"] = strconv.FormatBool(segmentValue.(bool))
		}
	}
	if segmentRulesItem.Order != nil {
		segmentRulesMap["order"] = segmentRulesItem.Order
	}

	return segmentRulesMap
}

func dataSourceFeatureSegmentRulesRulesToMap(rulesItem appconfigurationv1.TargetSegments) (rulesMap map[string]interface{}) {
	rulesMap = map[string]interface{}{}

	if rulesItem.Segments != nil {
		rulesMap["segments"] = rulesItem.Segments
	}

	return rulesMap
}

func dataSourceFeatureFlattenCollections(result []appconfigurationv1.CollectionRef) (collections []map[string]interface{}) {
	for _, collectionsItem := range result {
		collections = append(collections, dataSourceFeatureCollectionsToMap(collectionsItem))
	}

	return collections
}

func dataSourceFeatureCollectionsToMap(collectionsItem appconfigurationv1.CollectionRef) (collectionsMap map[string]interface{}) {
	collectionsMap = map[string]interface{}{}

	if collectionsItem.CollectionID != nil {
		collectionsMap["collection_id"] = collectionsItem.CollectionID
	}
	if collectionsItem.Name != nil {
		collectionsMap["name"] = collectionsItem.Name
	}

	return collectionsMap
}
