// Copyright IBM Corp. 2021 All Rights Reserved.
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

func DataSourceIBMAppConfigCollections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigCollectionsRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of records to retrieve.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Skipped number of records.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of records.",
			},
			"include": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Include feature, property details in the response.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"expand": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, returns expanded view of the resource details.",
			},
			"collections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of Features.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Collection name.",
						},
						"collection_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Collection id.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Collection description.",
						},
						"tags": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tags associated with the collection.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the collection.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time of the collection data.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Collection URL.",
						},
						"features_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of features associated with the collection.",
						},
						"properties_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of properties associated with the collection.",
						},
						"features": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of Features associated with the collection.",
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
							Description: "List of properties associated with the collection.",
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
				},
			},
		},
	}
}

func dataSourceIbmAppConfigCollectionsRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.ListCollectionsOptions{}

	if _, ok := GetFieldExists(d, "expand"); ok {
		options.SetExpand(d.Get("expand").(bool))
	}
	if _, ok := GetFieldExists(d, "sort"); ok {
		options.SetExpand(d.Get("sort").(bool))
	}
	if _, ok := GetFieldExists(d, "tags"); ok {
		options.SetExpand(d.Get("tags").(bool))
	}
	if _, ok := GetFieldExists(d, "include"); ok {
		includes := []string{}
		for _, includeItem := range d.Get("include").([]interface{}) {
			includes = append(includes, includeItem.(string))
		}
		options.SetInclude(includes)
	}
	if _, ok := GetFieldExists(d, "features"); ok {
		features := []string{}
		for _, featureItem := range d.Get("features").([]interface{}) {
			features = append(features, featureItem.(string))
		}
		options.SetFeatures(features)
	}
	if _, ok := GetFieldExists(d, "properties"); ok {
		properties := []string{}
		for _, propertieItem := range d.Get("properties").([]interface{}) {
			properties = append(properties, propertieItem.(string))
		}
		options.SetProperties(properties)
	}

	var collectionsList *appconfigurationv1.CollectionList
	var offset int64
	var limit int64 = 10
	var isLimit bool
	finalList := []appconfigurationv1.Collection{}

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
		result, response, err := appconfigClient.ListCollections(options)
		collectionsList = result
		if err != nil {
			return flex.FmtErrorf("ListCollections failed %s\n%s", err, response)
		}
		if isLimit {
			offset = 0
		} else {
			offset = dataSourceCollectionsListGetNext(result.Next)
		}
		finalList = append(finalList, result.Collections...)
		if offset == 0 {
			break
		}
	}

	collectionsList.Collections = finalList

	d.SetId(fmt.Sprintf("%s", guid))

	if collectionsList.Collections != nil {
		err = d.Set("collections", dataSourceCollectionListFlattenCollections(collectionsList.Collections))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting collections %s", err)
		}
	}
	if collectionsList.Limit != nil {
		if err = d.Set("limit", collectionsList.Limit); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting limit: %s", err)
		}
	}
	if collectionsList.Offset != nil {
		if err = d.Set("offset", collectionsList.Offset); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting offset: %s", err)
		}
	}
	if collectionsList.TotalCount != nil {
		if err = d.Set("total_count", collectionsList.TotalCount); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting total_count: %s", err)
		}
	}

	return nil
}

func dataSourceCollectionListFlattenCollections(result []appconfigurationv1.Collection) (collections []map[string]interface{}) {
	for _, collectionsItem := range result {
		collections = append(collections, dataSourceCollectionsListCollectionsToMap(collectionsItem))
	}
	return collections
}

func dataSourceCollectionsListCollectionsToMap(collectionItem appconfigurationv1.Collection) (collectionsMap map[string]interface{}) {
	collectionsMap = map[string]interface{}{}

	if collectionItem.Name != nil {
		collectionsMap["name"] = collectionItem.Name
	}
	if collectionItem.CollectionID != nil {
		collectionsMap["collection_id"] = collectionItem.CollectionID
	}
	if collectionItem.Description != nil {
		collectionsMap["description"] = collectionItem.Description
	}
	if collectionItem.Tags != nil {
		collectionsMap["tags"] = collectionItem.Tags
	}
	if collectionItem.CreatedTime != nil {
		collectionsMap["created_time"] = collectionItem.CreatedTime.String()
	}
	if collectionItem.UpdatedTime != nil {
		collectionsMap["updated_time"] = collectionItem.UpdatedTime.String()
	}
	if collectionItem.Href != nil {
		collectionsMap["href"] = collectionItem.Href
	}
	if collectionItem.FeaturesCount != nil {
		collectionsMap["features_count"] = collectionItem.FeaturesCount
	}
	if collectionItem.PropertiesCount != nil {
		collectionsMap["properties_count"] = collectionItem.PropertiesCount
	}
	if collectionItem.Features != nil {
		featuresList := []map[string]interface{}{}
		for _, featuresItem := range collectionItem.Features {
			featuresList = append(featuresList, dataSourceCollectionsListCollectionsFeaturesToMap(featuresItem))
		}
		collectionsMap["features"] = featuresList
	}
	if collectionItem.Properties != nil {
		propertiesList := []map[string]interface{}{}
		for _, propertiesItem := range collectionItem.Properties {
			propertiesList = append(propertiesList, dataSourceCollectionsListCollectionsPropertiesToMap(propertiesItem))
		}
		collectionsMap["properties"] = propertiesList
	}

	return collectionsMap
}

func dataSourceCollectionsListCollectionsFeaturesToMap(featuresItem appconfigurationv1.FeatureOutput) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.FeatureID != nil {
		featuresMap["feature_id"] = featuresItem.FeatureID
	}
	if featuresItem.Name != nil {
		featuresMap["name"] = featuresItem.Name
	}
	return featuresMap
}

func dataSourceCollectionsListCollectionsPropertiesToMap(propertiesItem appconfigurationv1.PropertyOutput) (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}

	if propertiesItem.PropertyID != nil {
		propertiesMap["property_id"] = propertiesItem.PropertyID
	}
	if propertiesItem.Name != nil {
		propertiesMap["name"] = propertiesItem.Name
	}
	return propertiesMap
}

func dataSourceCollectionsListGetNext(next interface{}) int64 {
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

func dataSourceCollectionListFlattenPagination(result interface{}) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceCollectionsListURLToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceCollectionsListURLToMap(urlItem interface{}) (urlMap map[string]interface{}) {
	urlMap = map[string]interface{}{}

	var hrefUrl *string
	switch urlItem := urlItem.(type) {
	case appconfigurationv1.PaginatedListFirst:
		hrefUrl = urlItem.Href
	case appconfigurationv1.PaginatedListLast:
		hrefUrl = urlItem.Href
	case appconfigurationv1.PaginatedListNext:
		hrefUrl = urlItem.Href
	case appconfigurationv1.PaginatedListPrevious:
		hrefUrl = urlItem.Href
	}
	if hrefUrl != nil {
		urlMap["href"] = hrefUrl
	}

	return urlMap
}
