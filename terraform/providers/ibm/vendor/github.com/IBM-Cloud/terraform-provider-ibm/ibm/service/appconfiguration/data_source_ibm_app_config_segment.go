package appconfiguration

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppConfigSegment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigSegmentRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"segment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Segment id.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Segment name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Segment description.",
			},
			"tags": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tags associated with the segment.",
			},
			"includes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Include feature and property details in the response.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the segment.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of the segment data.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Segment flag URL.",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of rules that determine if the entity belongs to the segment during feature / property evaluation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of values. Entities matching any of the given values will be considered to belong to the segment",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"attribute_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Attribute name.",
						},
						"operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operator to be used for the evaluation if the entity belongs to the segment.",
						},
					},
				},
			},
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Features associated with the segment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feature_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "feature id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Feature.",
						},
					},
				},
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Features associated with the segment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "property id.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Property.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmAppConfigSegmentRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}

	options := &appconfigurationv1.GetSegmentOptions{}
	options.SetSegmentID(d.Get("segment_id").(string))

	if _, ok := GetFieldExists(d, "includes"); ok {
		includes := []string{}
		for _, segmentsItem := range d.Get("includes").([]interface{}) {
			includes = append(includes, segmentsItem.(string))
		}
		options.SetInclude(includes)
	}

	result, response, err := appconfigClient.GetSegment(options)
	if err != nil {
		return flex.FmtErrorf("[ERROR] GetSegment failed %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", guid, *result.SegmentID))

	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting name: %s", err)
		}
	}
	if result.SegmentID != nil {
		if err = d.Set("segment_id", result.SegmentID); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting segment_id: %s", err)
		}
	}
	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting description: %s", err)
		}
	}
	if result.Tags != nil {
		if err = d.Set("tags", result.Tags); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting tags: %s", err)
		}
	}
	if result.CreatedTime != nil {
		if err = d.Set("created_time", result.CreatedTime.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting created_time: %s", err)
		}
	}
	if result.UpdatedTime != nil {
		if err = d.Set("updated_time", result.UpdatedTime.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting updated_time: %s", err)
		}
	}
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting href: %s", err)
		}
	}
	if result.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, ruleItem := range result.Rules {
			rulesList = append(rulesList, dataSourceSegmentListSegmentRulesToMap(ruleItem))
			if err = d.Set("rules", rulesList); err != nil {
				return flex.FmtErrorf("[ERROR] Error setting rules %s", err)
			}
		}
	}
	if result.Features != nil {
		err = d.Set("features", dataSourceSegmentFlattenFeatures(result.Features))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting features %s", err)
		}
	}
	if result.Properties != nil {
		err = d.Set("properties", dataSourceSegmentFlattenProperties(result.Properties))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting properties %s", err)
		}
	}
	return nil
}

func dataSourceSegmentListSegmentRulesToMap(segmentRulesItem appconfigurationv1.Rule) (segmentRulesMap map[string]interface{}) {
	segmentRulesMap = map[string]interface{}{}

	segmentRulesMap["values"] = segmentRulesItem.Values
	segmentRulesMap["attribute_name"] = segmentRulesItem.AttributeName
	segmentRulesMap["operator"] = segmentRulesItem.Operator

	return segmentRulesMap
}

func dataSourceSegmentFlattenFeatures(result []appconfigurationv1.FeatureOutput) (features []map[string]interface{}) {
	for _, featuresItem := range result {
		features = append(features, dataSourceSegmentFeaturesToMap(featuresItem))
	}

	return features
}

func dataSourceSegmentFeaturesToMap(featuresItem appconfigurationv1.FeatureOutput) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.FeatureID != nil {
		featuresMap["feature_id"] = featuresItem.FeatureID
	}
	if featuresItem.Name != nil {
		featuresMap["name"] = featuresItem.Name
	}

	return featuresMap
}

func dataSourceSegmentFlattenProperties(result []appconfigurationv1.PropertyOutput) (properties []map[string]interface{}) {
	for _, propertiesItem := range result {
		properties = append(properties, dataSourceSegmentPropertiesToMap(propertiesItem))
	}
	return properties
}

func dataSourceSegmentPropertiesToMap(propertiesItem appconfigurationv1.PropertyOutput) (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}

	if propertiesItem.PropertyID != nil {
		propertiesMap["property_id"] = propertiesItem.PropertyID
	}
	if propertiesItem.Name != nil {
		propertiesMap["name"] = propertiesItem.Name
	}

	return propertiesMap
}
