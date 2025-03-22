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

func DataSourceIBMAppConfigSegments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigSegmentsRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter the resources to be returned based on the associated tags.",
			},
			"sort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort the segment details based on the specified attribute.",
			},
			"include": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Segment details to include the associated rules in the response",
			},
			"expand": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to `true`, returns expanded view of the resource details.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different set of records, use `limit` with `offset` to page through the available records.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Total number of records.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset` value. Use `offset` with `limit` to page through the available records.",
			},
			"segments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of Segments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Segment name.",
						},
						"segment_id": {
							Type:        schema.TypeString,
							Computed:    true,
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
							Description: "Segment URL..",
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
							Description: "List of properties associated with the segment.",
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
					},
				},
			},
		},
	}
}

func dataSourceIbmAppConfigSegmentsRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.ListSegmentsOptions{}

	if _, ok := GetFieldExists(d, "expand"); ok {
		options.SetExpand(d.Get("expand").(bool))
	}
	if _, ok := GetFieldExists(d, "tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}
	if _, ok := GetFieldExists(d, "sort"); ok {
		options.SetTags(d.Get("sort").(string))
	}

	if _, ok := GetFieldExists(d, "include"); ok {
		options.SetInclude(d.Get("include").(string))
	}

	var segmentsList *appconfigurationv1.SegmentsList
	var offset int64
	var limit int64 = 10
	var isLimit bool

	finalList := []appconfigurationv1.Segment{}

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
		result, response, err := appconfigClient.ListSegments(options)
		segmentsList = result
		if err != nil {
			return flex.FmtErrorf("[ERROR] ListSegments failed %s\n%s", err, response)
		}
		if isLimit {
			offset = 0
		} else {
			offset = dataSourceSegmentsListGetNext(result.Next)
		}
		finalList = append(finalList, result.Segments...)
		if offset == 0 {
			break
		}
	}

	segmentsList.Segments = finalList

	d.SetId(fmt.Sprintf("%s", guid))

	if segmentsList.Segments != nil {
		err = d.Set("segments", dataSourceSegmentsListFlattenSegments(segmentsList.Segments))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting segments %s", err)
		}
	}
	if segmentsList.Limit != nil {
		if err = d.Set("limit", segmentsList.Limit); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting limit: %s", err)
		}
	}
	if segmentsList.Offset != nil {
		if err = d.Set("offset", segmentsList.Offset); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting offset: %s", err)
		}
	}
	if segmentsList.TotalCount != nil {
		if err = d.Set("total_count", segmentsList.TotalCount); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting total_count: %s", err)
		}
	}
	return nil
}

func dataSourceSegmentsListGetNext(next interface{}) int64 {
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

func dataSourceSegmentsListFlattenSegments(result []appconfigurationv1.Segment) (segments []map[string]interface{}) {
	for _, segmentsItem := range result {
		segments = append(segments, dataSourceSegmentsListSegmentsToMap(segmentsItem))
	}
	return segments
}

func dataSourceSegmentsListSegmentsToMap(segmentsItem appconfigurationv1.Segment) (segmentsMap map[string]interface{}) {
	segmentsMap = map[string]interface{}{}

	if segmentsItem.Name != nil {
		segmentsMap["name"] = segmentsItem.Name
	}
	if segmentsItem.SegmentID != nil {
		segmentsMap["segment_id"] = segmentsItem.SegmentID
	}
	if segmentsItem.Description != nil {
		segmentsMap["description"] = segmentsItem.Description
	}
	if segmentsItem.Tags != nil {
		segmentsMap["tags"] = segmentsItem.Tags
	}
	if segmentsItem.CreatedTime != nil {
		segmentsMap["created_time"] = segmentsItem.CreatedTime.String()
	}
	if segmentsItem.UpdatedTime != nil {
		segmentsMap["updated_time"] = segmentsItem.UpdatedTime.String()
	}
	if segmentsItem.Href != nil {
		segmentsMap["href"] = segmentsItem.Href
	}
	if segmentsItem.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, segmentRulesItem := range segmentsItem.Rules {
			rulesList = append(rulesList, dataSourceSegmentsListSegmentRulesToMap(segmentRulesItem))
		}
		segmentsMap["rules"] = rulesList
	}
	if segmentsItem.Features != nil {
		featuresList := []map[string]interface{}{}
		for _, segmentsItem := range segmentsItem.Features {
			featuresList = append(featuresList, dataSourceSegmentsListSegmentsFeaturesToMap(segmentsItem))
		}
		segmentsMap["features"] = featuresList
	}
	if segmentsItem.Properties != nil {
		propertiesList := []map[string]interface{}{}
		for _, segmentsItem := range segmentsItem.Properties {
			propertiesList = append(propertiesList, dataSourcePropertiesListSegmentsPropertiesToMap(segmentsItem))
		}
		segmentsMap["properties"] = propertiesList
	}

	return segmentsMap
}

func dataSourceSegmentsListSegmentRulesToMap(segmentRulesItem appconfigurationv1.Rule) (segmentRulesMap map[string]interface{}) {
	segmentRulesMap = map[string]interface{}{}

	segmentRulesMap["values"] = segmentRulesItem.Values
	segmentRulesMap["attribute_name"] = segmentRulesItem.AttributeName
	segmentRulesMap["operator"] = segmentRulesItem.Operator

	return segmentRulesMap
}

func dataSourceSegmentsListSegmentsFeaturesToMap(featuresItem appconfigurationv1.FeatureOutput) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.FeatureID != nil {
		featuresMap["feature_id"] = featuresItem.FeatureID
	}
	if featuresItem.Name != nil {
		featuresMap["name"] = featuresItem.Name
	}

	return featuresMap
}

func dataSourcePropertiesListSegmentsPropertiesToMap(propertiesItem appconfigurationv1.PropertyOutput) (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}

	if propertiesItem.PropertyID != nil {
		propertiesMap["property_id"] = propertiesItem.PropertyID
	}
	if propertiesItem.Name != nil {
		propertiesMap["name"] = propertiesItem.Name
	}

	return propertiesMap
}
