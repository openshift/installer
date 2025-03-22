package appconfiguration

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMAppConfigCollection() *schema.Resource {
	return &schema.Resource{
		Read:     resourceIbmIbmAppConfigCollectiontRead,
		Create:   resourceIbmIbmAppConfigCollectiontCreate,
		Update:   resourceIbmIbmAppConfigCollectionUpdate,
		Delete:   resourceIbmIbmAppConfigCollectionDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection name.",
			},
			"collection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection Id.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Collection description",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tags associated with the collection",
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Number of features associated with the collection.",
			},
			"properties_count": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Number of properties associated with the collection.",
			},
		},
	}
}

func resourceIbmIbmAppConfigCollectiontCreate(d *schema.ResourceData, meta interface{}) error {

	guid := d.Get("guid").(string)
	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}
	options := &appconfigurationv1.CreateCollectionOptions{}
	options.SetName(d.Get("name").(string))
	options.SetCollectionID(d.Get("collection_id").(string))

	if _, ok := GetFieldExists(d, "description"); ok {
		options.SetDescription(d.Get("description").(string))
	}
	if _, ok := GetFieldExists(d, "tags"); ok {
		options.SetTags(d.Get("tags").(string))
	}

	collection, response, err := appconfigClient.CreateCollection(options)

	if err != nil {
		return flex.FmtErrorf("CreateCollection failed %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", guid, *collection.CollectionID))

	return resourceIbmIbmAppConfigCollectiontRead(d, meta)
}

func resourceIbmIbmAppConfigCollectiontRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	if len(parts) != 2 {
		return flex.FmtErrorf("Kindly check the id")
	}

	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.GetCollectionOptions{}
	options.SetCollectionID(parts[1])

	result, response, err := appconfigClient.GetCollection(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
		}
		return flex.FmtErrorf("[ERROR] GetCollection failed %s\n%s", err, response)
	}

	d.Set("guid", parts[0])

	if result.Name != nil {
		if err = d.Set("name", result.Name); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting name: %s", err)
		}
	}
	if result.CollectionID != nil {
		if err = d.Set("collection_id", result.CollectionID); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting collection_id: %s", err)
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
	if result.FeaturesCount != nil {
		if err = d.Set("features_count", result.FeaturesCount); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting features_count: %s", err)
		}
	}
	if result.FeaturesCount != nil {
		if err = d.Set("features_count", result.FeaturesCount); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting features_count: %s", err)
		}
	}
	if result.PropertiesCount != nil {
		if err = d.Set("properties_count", result.PropertiesCount); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting properties_count: %s", err)
		}
	}
	if result.Features != nil {
		err = d.Set("features", resourceIbmAppConfigCollectionFeatureToMap(result.Features))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting features %s", err)
		}
	}
	if result.Properties != nil {
		err = d.Set("properties", resourceIbmAppConfigCollectionPropertiesToMap(result.Properties))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting properties %s", err)
		}
	}

	return nil
}

func resourceIbmAppConfigCollectionFeatureToMap(result []appconfigurationv1.FeatureOutput) (features []map[string]interface{}) {
	for _, featuresItem := range result {
		features = append(features, resourceCollectionFeaturesToMap(featuresItem))
	}
	return features
}

func resourceIbmAppConfigCollectionPropertiesToMap(result []appconfigurationv1.PropertyOutput) (properties []map[string]interface{}) {
	for _, propertiesItem := range result {
		properties = append(properties, resourceCollectionPropertiesToMap(propertiesItem))
	}
	return properties
}

func resourceCollectionFeaturesToMap(featuresItem appconfigurationv1.FeatureOutput) (featuresMap map[string]interface{}) {
	featuresMap = map[string]interface{}{}

	if featuresItem.FeatureID != nil {
		featuresMap["feature_id"] = featuresItem.FeatureID
	}
	if featuresItem.Name != nil {
		featuresMap["name"] = featuresItem.Name
	}
	return featuresMap
}

func resourceCollectionPropertiesToMap(propertiesItem appconfigurationv1.PropertyOutput) (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}

	if propertiesItem.PropertyID != nil {
		propertiesMap["property_id"] = propertiesItem.PropertyID
	}
	if propertiesItem.Name != nil {
		propertiesMap["name"] = propertiesItem.Name
	}
	return propertiesMap
}

func resourceIbmIbmAppConfigCollectionUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}
	options := &appconfigurationv1.UpdateCollectionOptions{}

	options.SetCollectionID(parts[1])

	if ok := d.HasChanges("name", "description", "tags"); ok {

		if _, ok := GetFieldExists(d, "name"); ok {
			options.SetName(d.Get("name").(string))
		}
		if _, ok := GetFieldExists(d, "description"); ok {
			options.SetDescription(d.Get("description").(string))
		}
		if _, ok := GetFieldExists(d, "tags"); ok {
			options.SetTags(d.Get("tags").(string))
		}

		_, response, err := appconfigClient.UpdateCollection(options)
		if err != nil {
			return flex.FmtErrorf("UpdateCollection failed %s\n%s", err, response)
		}
		return resourceIbmIbmAppConfigCollectiontRead(d, meta)
	}
	return nil
}

func resourceIbmIbmAppConfigCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.DeleteCollectionOptions{}
	options.SetCollectionID(parts[1])

	response, err := appconfigClient.DeleteCollection(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return flex.FmtErrorf("[ERROR] DeleteCollection failed %s\n%s", err, response)
	}

	d.SetId("")

	return nil
}
