// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0
package appconfiguration

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
)

func DataSourceIBMAppConfigProperties() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigPropertiesRead,

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
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort the feature details based on the specified attribute.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated tags. Returns resources associated with any of the specified tags.",
			},
			"collections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter features by a list of comma separated collections.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"segments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter features by a list of comma separated segments.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"expand": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to `true`, returns expanded view of the resource details.",
			},
			"include": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Include the associated collections or targeting rules details in the response.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different set of records, use `limit` with `offset` to page through the available records.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset` value. Use `offset` with `limit` to page through the available records.",
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Property name.",
						},
						"property_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Property id.",
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
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of records returned in the current response.",
			},
		},
	}
}

func dataSourceIbmAppConfigPropertiesRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.ListPropertiesOptions{}

	options.SetEnvironmentID(d.Get("environment_id").(string))

	if _, ok := GetFieldExists(d, "expand"); ok {
		options.SetExpand(d.Get("expand").(bool))
	}
	if _, ok := GetFieldExists(d, "sort"); ok {
		options.SetSort(d.Get("sort").(string))
	}
	if _, ok := GetFieldExists(d, "tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}
	if _, ok := GetFieldExists(d, "collections"); ok {
		collections := []string{}
		for _, item := range d.Get("collections").([]interface{}) {
			collections = append(collections, item.(string))
		}
		options.SetCollections(collections)
	}
	if _, ok := GetFieldExists(d, "segments"); ok {
		segments := []string{}
		for _, item := range d.Get("segments").([]interface{}) {
			segments = append(segments, item.(string))
		}
		options.SetSegments(segments)
	}
	if _, ok := GetFieldExists(d, "include"); ok {
		includes := []string{}
		for _, item := range d.Get("include").([]interface{}) {
			includes = append(includes, item.(string))
		}
		options.SetInclude(includes)
	}

	var propertiesList *appconfigurationv1.PropertiesList
	var offset int64
	var limit int64 = 10
	var isLimit bool

	finalList := []appconfigurationv1.Property{}
	if _, ok := GetFieldExists(d, "limit"); ok {
		isLimit = true
		limit = int64(d.Get("limit").(int))
	}
	options.SetLimit(limit)
	if _, ok := GetFieldExists(d, "offset"); ok {
		offset = int64(d.Get("offset").(int))
	}
	for {
		options.Offset = &offset
		result, response, err := appconfigClient.ListProperties(options)
		propertiesList = result
		if err != nil {
			return flex.FmtErrorf("ListProperties failed %s\n%s", err, response)
		}
		if isLimit {
			offset = 0
		} else {
			offset = dataSourcePropertiesListGetNext(result.Next)
		}
		finalList = append(finalList, result.Properties...)
		if offset == 0 {
			break
		}
	}

	propertiesList.Properties = finalList

	d.SetId(fmt.Sprintf("%s/%s", guid, *options.EnvironmentID))

	if propertiesList.Properties != nil {
		err = d.Set("properties", dataSourcePropertiesListFlattenProperties(propertiesList.Properties))
		if err != nil {
			return flex.FmtErrorf("error setting properties %s", err)
		}
	}
	if propertiesList.TotalCount != nil {
		if err = d.Set("total_count", propertiesList.TotalCount); err != nil {
			return flex.FmtErrorf("error setting total_count: %s", err)
		}
	}
	if propertiesList.Limit != nil {
		if err = d.Set("limit", propertiesList.Limit); err != nil {
			return flex.FmtErrorf("error setting limit: %s", err)
		}
	}
	if propertiesList.Offset != nil {
		if err = d.Set("offset", propertiesList.Offset); err != nil {
			return flex.FmtErrorf("error setting offset: %s", err)
		}
	}

	return nil
}

func dataSourcePropertiesListGetNext(next interface{}) int64 {
	if reflect.ValueOf(next).IsNil() {
		return 0
	}

	u, err := url.Parse(reflect.ValueOf(next).Elem().FieldByName("Href").Elem().String())
	if err != nil {
		return 0
	}

	q := u.Query()
	var page string

	if q.Get("start") != "" {
		page = q.Get("start")
	} else if q.Get("offset") != "" {
		page = q.Get("offset")
	}

	convertedVal, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return 0
	}
	return convertedVal
}

func dataSourcePropertiesListFlattenProperties(result []appconfigurationv1.Property) (properties []map[string]interface{}) {
	for _, propertyItem := range result {
		properties = append(properties, dataSourcePropertiesListPropertiesToMap(propertyItem))
	}
	return properties
}

func dataSourcePropertiesListPropertiesToMap(property appconfigurationv1.Property) (propertyMap map[string]interface{}) {
	propertyMap = map[string]interface{}{}

	if property.Name != nil {
		propertyMap["name"] = property.Name
	}
	if property.PropertyID != nil {
		propertyMap["property_id"] = property.PropertyID
	}
	if property.Description != nil {
		propertyMap["description"] = property.Description
	}
	if property.Type != nil {
		propertyMap["type"] = property.Type
	}
	if property.Value != nil {
		value := property.Value
		switch value.(interface{}).(type) {
		case string:
			propertyMap["value"] = value.(string)
		case float64:
			propertyMap["value"] = fmt.Sprintf("%v", value)
		case bool:
			propertyMap["value"] = strconv.FormatBool(value.(bool))
		}
	}
	if property.Tags != nil {
		propertyMap["tags"] = property.Tags
	}
	if property.SegmentExists != nil {
		propertyMap["segment_exists"] = property.SegmentExists
	}
	if property.UpdatedTime != nil {
		propertyMap["updated_time"] = property.UpdatedTime.String()
	}
	if property.CreatedTime != nil {
		propertyMap["created_time"] = property.CreatedTime.String()
	}
	if property.Href != nil {
		propertyMap["href"] = property.Href
	}
	if property.Format != nil {
		propertyMap["format"] = property.Format
	}
	if property.Collections != nil {
		collectionsList := []map[string]interface{}{}
		for _, collectionsItem := range property.Collections {
			collectionsList = append(collectionsList, dataSourcePropertiesListPropertiesCollectionsToMap(collectionsItem))
		}
		propertyMap["collections"] = collectionsList
	}
	if property.SegmentRules != nil {
		segmentRulesList := []map[string]interface{}{}
		for _, segmentRulesItem := range property.SegmentRules {
			segmentRulesList = append(segmentRulesList, dataSourcePropertyListPropertiesSegmentRulesToMap(segmentRulesItem))
		}
		propertyMap["segment_rules"] = segmentRulesList
	}
	return propertyMap
}

func dataSourcePropertiesListPropertiesCollectionsToMap(collectionsItem appconfigurationv1.CollectionRef) (collectionsMap map[string]interface{}) {
	collectionsMap = map[string]interface{}{}

	if collectionsItem.CollectionID != nil {
		collectionsMap["collection_id"] = collectionsItem.CollectionID
	}
	if collectionsItem.Name != nil {
		collectionsMap["name"] = collectionsItem.Name
	}

	return collectionsMap
}

func dataSourcePropertyListPropertiesSegmentRulesToMap(segmentRulesItem appconfigurationv1.SegmentRule) (segmentRulesMap map[string]interface{}) {
	segmentRulesMap = map[string]interface{}{}

	if segmentRulesItem.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, rulesItem := range segmentRulesItem.Rules {
			rulesList = append(rulesList, dataSourceListPropertiesSegmentRulesToMap(rulesItem))
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

func dataSourceListPropertiesSegmentRulesToMap(rule appconfigurationv1.TargetSegments) map[string]interface{} {
	ruleMap := map[string]interface{}{}

	ruleMap["segments"] = rule.Segments

	return ruleMap
}
