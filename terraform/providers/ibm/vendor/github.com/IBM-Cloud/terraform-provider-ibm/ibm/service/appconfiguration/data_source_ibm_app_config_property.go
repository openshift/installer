package appconfiguration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
)

func DataSourceIBMAppConfigProperty() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigPropertyRead,

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
			"property_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Property Id.",
			},
			"include": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Include the associated collections in the response.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property description.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the Property (BOOLEAN, STRING, NUMERIC).",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the Property. The value can be Boolean, String or a Numeric value as per the `type` attribute.",
			},
			"tags": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tags associated with the property.",
			},
			"format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Format of the feature (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required for `BOOLEAN` and `NUMERIC` types. This property is populated in the response body of `POST, PUT and GET` calls if the type `STRING` is used and not populated for `BOOLEAN` and `NUMERIC` types.",
			},
			"segment_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specify the targeting rules that is used to set different property values for different segments.",
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
			"segment_exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes if the targeting rules are specified for the property.",
			},
			"collections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of collection id representing the collections that are associated with the specified property.",
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
				Description: "Creation time of the property.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of the property data.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property URL.",
			},
		},
	}
}

func dataSourceIbmAppConfigPropertyRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.GetPropertyOptions{}

	options.SetEnvironmentID(d.Get("environment_id").(string))
	options.SetPropertyID(d.Get("property_id").(string))

	if _, ok := GetFieldExists(d, "include"); ok {
		dataString := d.Get("include").(string)
		options.SetInclude(strings.Split(dataString, ","))
	}

	property, response, err := appconfigClient.GetProperty(options)

	if err != nil {
		return flex.FmtErrorf("GetProperty failed %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", guid, *options.EnvironmentID, *property.PropertyID))

	if property.Name != nil {
		if err = d.Set("name", property.Name); err != nil {
			return flex.FmtErrorf("error setting name: %s", err)
		}
	}
	if property.Description != nil {
		if err = d.Set("description", property.Description); err != nil {
			return flex.FmtErrorf("error setting description: %s", err)
		}
	}
	if property.Type != nil {
		if err = d.Set("type", property.Type); err != nil {
			return flex.FmtErrorf("error setting type: %s", err)
		}
	}
	if property.Value != nil {
		value := property.Value
		switch value.(interface{}).(type) {
		case string:
			d.Set("value", value.(string))
		case float64:
			d.Set("value", fmt.Sprintf("%v", value))
		case bool:
			d.Set("value", strconv.FormatBool(value.(bool)))
		}
	}
	if property.Tags != nil {
		if err = d.Set("tags", property.Tags); err != nil {
			return flex.FmtErrorf("error setting tags: %s", err)
		}
	}
	if property.Format != nil {
		if err = d.Set("format", property.Format); err != nil {
			return flex.FmtErrorf("error setting format: %s", err)
		}
	}
	if property.SegmentRules != nil {
		err = d.Set("segment_rules", dataSourcePropertyFlattenSegmentRules(property.SegmentRules))
		if err != nil {
			return flex.FmtErrorf("error setting segment_rules %s", err)
		}
	}
	if property.SegmentExists != nil {
		if err = d.Set("segment_exists", property.SegmentExists); err != nil {
			return flex.FmtErrorf("error setting segment_exists: %s", err)
		}
	}
	if property.Collections != nil {
		err = d.Set("collections", dataSourcePropertyFlattenCollections(property.Collections))
		if err != nil {
			return flex.FmtErrorf("error setting collections %s", err)
		}
	}
	if property.CreatedTime != nil {
		if err = d.Set("created_time", property.CreatedTime.String()); err != nil {
			return flex.FmtErrorf("error setting created_time: %s", err)
		}
	}
	if property.UpdatedTime != nil {
		if err = d.Set("updated_time", property.UpdatedTime.String()); err != nil {
			return flex.FmtErrorf("error setting updated_time: %s", err)
		}
	}
	if property.Href != nil {
		if err = d.Set("href", property.Href); err != nil {
			return flex.FmtErrorf("error setting href: %s", err)
		}
	}
	return nil
}

func dataSourcePropertyFlattenCollections(result []appconfigurationv1.CollectionRef) (collections []map[string]interface{}) {
	for _, collectionsItem := range result {
		collections = append(collections, dataSourcePropertyCollectionsToMap(collectionsItem))
	}

	return collections
}

func dataSourcePropertyCollectionsToMap(collectionsItem appconfigurationv1.CollectionRef) (collectionsMap map[string]interface{}) {
	collectionsMap = map[string]interface{}{}

	if collectionsItem.CollectionID != nil {
		collectionsMap["collection_id"] = collectionsItem.CollectionID
	}
	if collectionsItem.Name != nil {
		collectionsMap["name"] = collectionsItem.Name
	}

	return collectionsMap
}

func dataSourcePropertyFlattenSegmentRules(result []appconfigurationv1.SegmentRule) (segmentRules []map[string]interface{}) {
	for _, segmentRulesItem := range result {
		segmentRules = append(segmentRules, dataSourcePropertySegmentRulesToMap(segmentRulesItem))
	}
	return segmentRules
}

func dataSourcePropertySegmentRulesToMap(segmentRulesItem appconfigurationv1.SegmentRule) (segmentRulesMap map[string]interface{}) {
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
