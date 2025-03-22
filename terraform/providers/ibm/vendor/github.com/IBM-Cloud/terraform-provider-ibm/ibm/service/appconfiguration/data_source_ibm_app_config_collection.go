// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0
package appconfiguration

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppConfigCollection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigCollectionRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"collection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection Id of the collection.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Collection name.",
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
	}
}

func dataSourceIbmAppConfigCollectionRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.GetCollectionOptions{}

	options.SetCollectionID(d.Get("collection_id").(string))

	if _, ok := GetFieldExists(d, "include"); ok {
		includes := []string{}
		for _, includeItem := range d.Get("include").([]interface{}) {
			includes = append(includes, includeItem.(string))
		}
		options.SetInclude(includes)
	}
	if _, ok := GetFieldExists(d, "expand"); ok {
		options.SetExpand(d.Get("expand").(bool))
	}

	result, response, err := appconfigClient.GetCollection(options)
	if err != nil {
		return flex.FmtErrorf("GetCollection failed %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", guid, *result.CollectionID))
	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting name: %s", err)
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
	if result.Features != nil {
		err = d.Set("features", dataSourceCollectionFlattenFeatures(result.Features))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting features %s", err)
		}
	}
	if result.Properties != nil {
		err = d.Set("properties", dataSourceCollectionFlattenProperties(result.Properties))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting properties %s", err)
		}
	}

	return nil
}

func dataSourceCollectionFlattenFeatures(result []appconfigurationv1.FeatureOutput) (features []map[string]interface{}) {
	for _, featuresItem := range result {
		features = append(features, dataSourceCollectionFeaturesToMap(featuresItem))
	}

	return features
}

func dataSourceCollectionFlattenProperties(result []appconfigurationv1.PropertyOutput) (properties []map[string]interface{}) {
	for _, propertiesItem := range result {
		properties = append(properties, dataSourceCollectionPropertiesToMap(propertiesItem))
	}

	return properties
}

func dataSourceCollectionFeaturesToMap(featuresItem appconfigurationv1.FeatureOutput) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.FeatureID != nil {
		featuresMap["feature_id"] = featuresItem.FeatureID
	}
	if featuresItem.Name != nil {
		featuresMap["name"] = featuresItem.Name
	}
	return featuresMap
}

func dataSourceCollectionPropertiesToMap(propertiesItem appconfigurationv1.PropertyOutput) (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}

	if propertiesItem.PropertyID != nil {
		propertiesMap["property_id"] = propertiesItem.PropertyID
	}
	if propertiesItem.Name != nil {
		propertiesMap["name"] = propertiesItem.Name
	}
	return propertiesMap
}
