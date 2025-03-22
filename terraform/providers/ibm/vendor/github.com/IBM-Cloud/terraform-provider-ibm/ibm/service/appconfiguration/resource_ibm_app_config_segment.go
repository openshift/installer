package appconfiguration

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIbmAppConfigSegment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmIbmAppConfigSegmentCreate,
		Read:     resourceIbmIbmAppConfigSegmentRead,
		Update:   resourceIbmIbmAppConfigSegmentUpdate,
		Delete:   resourceIbmIbmAppConfigSegmentDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Segment name.",
			},
			"segment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Segment id.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Segment description.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags associated with the segments.",
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
				Description: "Segment URL.",
			},
			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of rules that determine if the entity belongs to the segment during feature / property evaluation. An entity is identified by an unique identifier and the attributes that it defines.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute name.",
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Operator to be used for the evaluation if the entity belongs to the segment.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of values. Entities matching any of the given values will be considered to belong to the segment.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceIbmIbmAppConfigSegmentCreate(d *schema.ResourceData, meta interface{}) error {

	guid := d.Get("guid").(string)
	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}
	options := &appconfigurationv1.CreateSegmentOptions{}
	options.SetName(d.Get("name").(string))
	options.SetSegmentID(d.Get("segment_id").(string))

	if _, ok := GetFieldExists(d, "description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := GetFieldExists(d, "tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}

	if _, ok := GetFieldExists(d, "rules"); ok {
		var segmentRules []appconfigurationv1.Rule
		for _, e := range d.Get("rules").([]interface{}) {
			value := e.(map[string]interface{})
			segmentRulesItem, err := resourceIbmAppConfigMapToSegmentRule(value)
			if err != nil {
				return flex.FmtErrorf(fmt.Sprintf("%s", err))
			}
			segmentRules = append(segmentRules, segmentRulesItem)
		}
		options.SetRules(segmentRules)
	}

	segment, response, err := appconfigClient.CreateSegment(options)

	if err != nil {
		return flex.FmtErrorf("CreateSegment failed %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", guid, *segment.SegmentID))
	return resourceIbmIbmAppConfigSegmentRead(d, meta)
}

func resourceIbmAppConfigMapToSegmentRule(segmentRuleMap map[string]interface{}) (appconfigurationv1.Rule, error) {
	segmentRule := appconfigurationv1.Rule{}

	segmentRule.AttributeName = core.StringPtr(segmentRuleMap["attribute_name"].(string))
	segmentRule.Operator = core.StringPtr(segmentRuleMap["operator"].(string))
	var values []string
	for _, rulesItem := range segmentRuleMap["values"].([]interface{}) {
		values = append(values, rulesItem.(string))
	}
	segmentRule.Values = values
	return segmentRule, nil
}

func resourceIbmIbmAppConfigSegmentRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	if len(parts) != 2 {
		return flex.FmtErrorf("Kindly check the id")
	}

	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}

	options := &appconfigurationv1.GetSegmentOptions{}
	options.SetSegmentID(parts[1])

	result, response, err := appconfigClient.GetSegment(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
		}
		return flex.FmtErrorf("[ERROR] GetSegment failed %s\n%s", err, response)
	}

	d.Set("guid", parts[0])

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
			return flex.FmtErrorf("[ERROR] Error setting createdTime: %s", err)
		}
	}
	if result.UpdatedTime != nil {
		if err = d.Set("updated_time", result.UpdatedTime.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting updatedTime: %s", err)
		}
	}
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting href: %s", err)
		}
	}
	if result.Rules != nil {
		segmentRules := []map[string]interface{}{}
		for _, ruleItem := range result.Rules {
			segmentRulesItemMap := resourceIbmAppConfigSegmentRuleToMap(ruleItem)
			segmentRules = append(segmentRules, segmentRulesItemMap)
		}
		if err = d.Set("rules", segmentRules); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting segment_rules: %s", err)
		}
	}
	if result.Features != nil {
		err = d.Set("features", resourceIbmAppConfigSegmentFeatureToMap(result.Features))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting features %s", err)
		}
	}
	if result.Properties != nil {
		err = d.Set("properties", resourceIbmAppConfigSegmentPropertiesToMap(result.Properties))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting properties %s", err)
		}
	}
	return nil
}

func resourceIbmIbmAppConfigSegmentUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}
	options := &appconfigurationv1.UpdateSegmentOptions{}

	options.SetSegmentID(parts[1])

	if ok := d.HasChanges("name", "description", "tags", "rules", "attribute_name", "operator", "values"); ok {

		if _, ok := GetFieldExists(d, "name"); ok {
			options.SetName(d.Get("name").(string))
		}
		if _, ok := GetFieldExists(d, "description"); ok {
			options.SetDescription(d.Get("description").(string))
		}
		if _, ok := GetFieldExists(d, "tags"); ok {
			options.SetTags(d.Get("tags").(string))
		}
		if _, ok := GetFieldExists(d, "rules"); ok {
			var segmentRules []appconfigurationv1.Rule
			for _, e := range d.Get("rules").([]interface{}) {
				value := e.(map[string]interface{})
				segmentRulesItem, err := resourceIbmAppConfigMapToSegmentRule(value)
				if err != nil {
					return flex.FmtErrorf(fmt.Sprintf("%s", err))
				}
				segmentRules = append(segmentRules, segmentRulesItem)
			}
			options.SetRules(segmentRules)
		}

		_, response, err := appconfigClient.UpdateSegment(options)
		if err != nil {
			return flex.FmtErrorf("[ERROR] UpdateSegment %s\n%s", err, response)
		}
		return resourceIbmIbmAppConfigSegmentRead(d, meta)
	}
	return nil
}

func resourceIbmIbmAppConfigSegmentDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}

	options := &appconfigurationv1.DeleteSegmentOptions{}
	options.SetSegmentID(parts[1])

	response, err := appconfigClient.DeleteSegment(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return flex.FmtErrorf("[ERROR] DeleteSegment failed %s\n%s", err, response)
	}

	d.SetId("")

	return nil
}

func resourceIbmAppConfigSegmentRuleToMap(segmentRulesItem appconfigurationv1.Rule) (segmentRulesMap map[string]interface{}) {
	segmentRulesMap = map[string]interface{}{}

	segmentRulesMap["values"] = segmentRulesItem.Values
	segmentRulesMap["attribute_name"] = segmentRulesItem.AttributeName
	segmentRulesMap["operator"] = segmentRulesItem.Operator

	return segmentRulesMap
}

func resourceIbmAppConfigSegmentFeatureToMap(result []appconfigurationv1.FeatureOutput) (features []map[string]interface{}) {
	for _, featuresItem := range result {
		features = append(features, resourceSegmentFeaturesToMap(featuresItem))
	}
	return features
}

func resourceIbmAppConfigSegmentPropertiesToMap(result []appconfigurationv1.PropertyOutput) (properties []map[string]interface{}) {
	for _, propertiesItem := range result {
		properties = append(properties, resourceSegmentPropertiesToMap(propertiesItem))
	}
	return properties
}

func resourceSegmentFeaturesToMap(featuresItem appconfigurationv1.FeatureOutput) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.FeatureID != nil {
		featuresMap["feature_id"] = featuresItem.FeatureID
	}
	if featuresItem.Name != nil {
		featuresMap["name"] = featuresItem.Name
	}
	return featuresMap
}

func resourceSegmentPropertiesToMap(propertiesItem appconfigurationv1.PropertyOutput) (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}

	if propertiesItem.PropertyID != nil {
		propertiesMap["property_id"] = propertiesItem.PropertyID
	}
	if propertiesItem.Name != nil {
		propertiesMap["name"] = propertiesItem.Name
	}

	return propertiesMap
}
