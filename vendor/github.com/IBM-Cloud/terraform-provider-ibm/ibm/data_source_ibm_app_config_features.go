// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
)

func dataSourceIbmAppConfigFeatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigFeaturesRead,

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
			"includes": {
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
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of Features.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Feature name.",
						},
						"feature_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Feature id.",
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
						"segment_exists": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Denotes if the targeting rules are specified for the feature flag.",
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
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of records returned in the current response.",
			},
			"next": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the next list of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
			"first": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the first page of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
			"previous": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the previous list of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
			"last": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "URL to navigate to the last page of records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL of the response.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmAppConfigFeaturesRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return err
	}

	options := &appconfigurationv1.ListFeaturesOptions{}
	options.SetEnvironmentID(d.Get("environment_id").(string))
	if _, ok := d.GetOk("expand"); ok {
		options.SetExpand(d.Get("expand").(bool))
	}
	if _, ok := d.GetOk("tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}
	if _, ok := d.GetOk("collections"); ok {
		collections := []string{}
		for _, segmentsItem := range d.Get("collections").([]interface{}) {
			collections = append(collections, segmentsItem.(string))
		}
		options.SetCollections(collections)
	}
	if _, ok := d.GetOk("segments"); ok {
		segments := []string{}
		for _, segmentsItem := range d.Get("segments").([]interface{}) {
			segments = append(segments, segmentsItem.(string))
		}
		options.SetSegments(segments)
	}
	if _, ok := d.GetOk("includes"); ok {
		includes := []string{}
		for _, segmentsItem := range d.Get("includes").([]interface{}) {
			includes = append(includes, segmentsItem.(string))
		}
		options.SetInclude(includes)
	}

	var featuresList *appconfigurationv1.FeaturesList
	var offset int64
	var limit int64 = 10
	var isLimit bool
	finalList := []appconfigurationv1.Feature{}

	if _, ok := d.GetOk("limit"); ok {
		isLimit = true
		limit = int64(d.Get("limit").(int))
	}
	options.SetLimit(limit)
	if _, ok := d.GetOk("offset"); ok {
		offset = int64(d.Get("offset").(int))
	}
	for {
		options.Offset = &offset
		result, response, err := appconfigClient.ListFeatures(options)
		featuresList = result
		if err != nil {
			log.Printf("[DEBUG] ListFeatures failed %s\n%s", err, response)
			return err
		}
		if isLimit {
			offset = 0
		} else {
			offset = dataSourceFeaturesListGetNext(result.Next)
		}
		finalList = append(finalList, result.Features...)
		if offset == 0 {
			break
		}
	}

	featuresList.Features = finalList

	d.SetId(fmt.Sprintf("%s/%s", guid, *options.EnvironmentID))

	if featuresList.Features != nil {
		err = d.Set("features", dataSourceFeaturesListFlattenFeatures(featuresList.Features))
		if err != nil {
			return fmt.Errorf("error setting features %s", err)
		}
	}
	if featuresList.TotalCount != nil {
		if err = d.Set("total_count", featuresList.TotalCount); err != nil {
			return fmt.Errorf("error setting total_count: %s", err)
		}
	}
	if featuresList.Limit != nil {
		if err = d.Set("limit", featuresList.Limit); err != nil {
			return fmt.Errorf("error setting limit: %s", err)
		}
	}
	if featuresList.Offset != nil {
		if err = d.Set("offset", featuresList.Offset); err != nil {
			return fmt.Errorf("error setting offset: %s", err)
		}
	}
	if featuresList.First != nil {
		err = d.Set("first", dataSourceFeatureListFlattenPagination(*featuresList.First))
		if err != nil {
			return fmt.Errorf("error setting first %s", err)
		}
	}

	if featuresList.Previous != nil {
		err = d.Set("previous", dataSourceFeatureListFlattenPagination(*featuresList.Previous))
		if err != nil {
			return fmt.Errorf("error setting previous %s", err)
		}
	}

	if featuresList.Last != nil {
		err = d.Set("last", dataSourceFeatureListFlattenPagination(*featuresList.Last))
		if err != nil {
			return fmt.Errorf("error setting last %s", err)
		}
	}
	if featuresList.Next != nil {
		err = d.Set("next", dataSourceFeatureListFlattenPagination(*featuresList.Next))
		if err != nil {
			return fmt.Errorf("error setting next %s", err)
		}
	}

	return nil
}

func dataSourceFeaturesListFlattenFeatures(result []appconfigurationv1.Feature) (features []map[string]interface{}) {
	for _, featuresItem := range result {
		features = append(features, dataSourceFeaturesListFeaturesToMap(featuresItem))
	}

	return features
}

func dataSourceFeaturesListFeaturesToMap(featuresItem appconfigurationv1.Feature) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.Name != nil {
		featuresMap["name"] = featuresItem.Name
	}
	if featuresItem.FeatureID != nil {
		featuresMap["feature_id"] = featuresItem.FeatureID
	}
	if featuresItem.Description != nil {
		featuresMap["description"] = featuresItem.Description
	}
	if featuresItem.Type != nil {
		featuresMap["type"] = featuresItem.Type
	}
	if featuresItem.Enabled != nil {
		featuresMap["enabled"] = featuresItem.Enabled
	}
	if featuresItem.Tags != nil {
		featuresMap["tags"] = featuresItem.Tags
	}
	if featuresItem.SegmentRules != nil {
		segmentRulesList := []map[string]interface{}{}
		for _, segmentRulesItem := range featuresItem.SegmentRules {
			segmentRulesList = append(segmentRulesList, dataSourceFeaturesListFeaturesSegmentRulesToMap(segmentRulesItem))
		}
		featuresMap["segment_rules"] = segmentRulesList
	}
	if featuresItem.SegmentExists != nil {
		featuresMap["segment_exists"] = featuresItem.SegmentExists
	}
	if featuresItem.Collections != nil {
		collectionsList := []map[string]interface{}{}
		for _, collectionsItem := range featuresItem.Collections {
			collectionsList = append(collectionsList, dataSourceFeaturesListFeaturesCollectionsToMap(collectionsItem))
		}
		featuresMap["collections"] = collectionsList
	}
	if featuresItem.CreatedTime != nil {
		featuresMap["created_time"] = featuresItem.CreatedTime.String()
	}
	if featuresItem.UpdatedTime != nil {
		featuresMap["updated_time"] = featuresItem.UpdatedTime.String()
	}
	if featuresItem.Href != nil {
		featuresMap["href"] = featuresItem.Href
	}
	if featuresItem.EnabledValue != nil {
		enabledValue := featuresItem.EnabledValue

		switch enabledValue.(interface{}).(type) {
		case string:
			featuresMap["enabled_value"] = enabledValue.(string)
		case float64:
			featuresMap["enabled_value"] = fmt.Sprintf("%v", enabledValue)
		case bool:
			featuresMap["enabled_value"] = strconv.FormatBool(enabledValue.(bool))
		}
	}

	if featuresItem.DisabledValue != nil {
		disabledValue := featuresItem.DisabledValue

		switch disabledValue.(interface{}).(type) {
		case string:
			featuresMap["disabled_value"] = disabledValue.(string)
		case float64:
			featuresMap["disabled_value"] = fmt.Sprintf("%v", disabledValue)
		case bool:
			featuresMap["disabled_value"] = strconv.FormatBool(disabledValue.(bool))
		}
	}
	return featuresMap
}

func dataSourceFeaturesListFeaturesSegmentRulesToMap(segmentRulesItem appconfigurationv1.SegmentRule) (segmentRulesMap map[string]interface{}) {
	segmentRulesMap = map[string]interface{}{}

	if segmentRulesItem.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, rulesItem := range segmentRulesItem.Rules {
			rulesList = append(rulesList, dataSourceListFeaturesSegmentRulesRulesToMap(rulesItem))
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

func dataSourceFeaturesListFeaturesCollectionsToMap(collectionsItem appconfigurationv1.CollectionRef) (collectionsMap map[string]interface{}) {
	collectionsMap = map[string]interface{}{}

	if collectionsItem.CollectionID != nil {
		collectionsMap["collection_id"] = collectionsItem.CollectionID
	}
	if collectionsItem.Name != nil {
		collectionsMap["name"] = collectionsItem.Name
	}

	return collectionsMap
}

func dataSourceFeaturesListGetNext(next interface{}) int64 {
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

func dataSourceListFeaturesSegmentRulesRulesToMap(rule appconfigurationv1.TargetSegments) map[string]interface{} {
	ruleMap := map[string]interface{}{}

	ruleMap["segments"] = rule.Segments

	return ruleMap
}

func dataSourceFeatureListFlattenPagination(result appconfigurationv1.PageHrefResponse) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceFeatureListURLToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceFeatureListURLToMap(urlItem appconfigurationv1.PageHrefResponse) (urlMap map[string]interface{}) {
	urlMap = map[string]interface{}{}

	if urlItem.Href != nil {
		urlMap["href"] = urlItem.Href
	}

	return urlMap
}
