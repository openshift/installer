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

func DataSourceIBMAppConfigSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"collection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the response based on the specified collection_id.",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the response based on the specified environment_id.",
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
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of records.",
			},
			"git_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of git_config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_config_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Git config name.",
						},
						"git_config_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Git config id",
						},
						"git_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Git url which will be used to connect to the github account",
						},
						"git_branch": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Git url which will be used to connect to the github account",
						},
						"git_file_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Git file path, this is a path where your configuration file will be written.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the git config.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time of the git config data.",
						},
						"last_sync_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest time when the snapshot was synced to git.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Git config URL.",
						},
						"collection": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"collection_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Collection name.",
									},
									"collection_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Collection id.",
									},
								},
							},
						},
						"environment": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Environment object",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment name.",
									},
									"environment_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment id.",
									},
									"color_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment color code.",
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

func dataSourceIbmAppConfigSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.ListSnapshotsOptions{}

	if _, ok := GetFieldExists(d, "collection_id"); ok {
		options.SetCollectionID(d.Get("collection_id").(string))
	}
	if _, ok := GetFieldExists(d, "environment_id"); ok {
		options.SetEnvironmentID(d.Get("environment_id").(string))
	}
	if _, ok := GetFieldExists(d, "search"); ok {
		options.SetSearch(d.Get("search").(string))
	}
	if _, ok := GetFieldExists(d, "sort"); ok {
		options.SetSort(d.Get("sort").(string))
	}

	var shapshotsList *appconfigurationv1.GitConfigList
	var offset int64
	var limit int64 = 10
	var isLimit bool

	finalList := []appconfigurationv1.GitConfig{}

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
		result, response, err := appconfigClient.ListSnapshots(options)
		shapshotsList = result
		if err != nil {
			return flex.FmtErrorf("ListSnapshots failed %s\n%s", err, response)
		}
		if isLimit {
			offset = 0
		} else {
			offset = dataSourceSnapshotsListGetNext(result.Next)
		}
		finalList = append(finalList, result.GitConfig...)
		if offset == 0 {
			break
		}
	}

	shapshotsList.GitConfig = finalList

	d.SetId(fmt.Sprintf("%s", guid))

	if shapshotsList.GitConfig != nil {
		err = d.Set("git_config", dataSourceFeaturesListFlattenSnapshots(shapshotsList.GitConfig))
		if err != nil {
			return flex.FmtErrorf("[ERROR] Error setting git_config %s", err)
		}
	}
	if shapshotsList.TotalCount != nil {
		if err = d.Set("total_count", shapshotsList.TotalCount); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting total_count: %s", err)
		}
	}
	if shapshotsList.Limit != nil {
		if err = d.Set("limit", shapshotsList.Limit); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting limit: %s", err)
		}
	}
	if shapshotsList.Offset != nil {
		if err = d.Set("offset", shapshotsList.Offset); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting offset: %s", err)
		}
	}

	return nil
}

func dataSourceSnapshotsListGetNext(next interface{}) int64 {
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

func dataSourceFeaturesListFlattenSnapshots(result []appconfigurationv1.GitConfig) (snapshots []map[string]interface{}) {
	for _, snapshotsItem := range result {
		snapshots = append(snapshots, dataSourceSnapshotsListSnapshotsToMap(snapshotsItem))
	}
	return snapshots
}

func dataSourceSnapshotsListSnapshotsToMap(snapshotsItem appconfigurationv1.GitConfig) (snapshotsMap map[string]interface{}) {
	snapshotsMap = map[string]interface{}{}

	if snapshotsItem.GitConfigName != nil {
		snapshotsMap["git_config_name"] = snapshotsItem.GitConfigName
	}
	if snapshotsItem.GitConfigID != nil {
		snapshotsMap["git_config_id"] = snapshotsItem.GitConfigID
	}
	if snapshotsItem.GitURL != nil {
		snapshotsMap["git_url"] = snapshotsItem.GitURL
	}
	if snapshotsItem.GitBranch != nil {
		snapshotsMap["git_branch"] = snapshotsItem.GitBranch
	}
	if snapshotsItem.GitFilePath != nil {
		snapshotsMap["git_file_path"] = snapshotsItem.GitFilePath
	}
	if snapshotsItem.CreatedTime != nil {
		snapshotsMap["created_time"] = snapshotsItem.CreatedTime.String()
	}
	if snapshotsItem.UpdatedTime != nil {
		snapshotsMap["updated_time"] = snapshotsItem.UpdatedTime.String()
	}
	if snapshotsItem.LastSyncTime != nil {
		snapshotsMap["last_sync_time"] = snapshotsItem.LastSyncTime.String()
	}
	if snapshotsItem.Href != nil {
		snapshotsMap["href"] = snapshotsItem.Href
	}

	if snapshotsItem.Collection != nil {
		collectionItemMap := resourceIbmAppConfigSnapshotsCollectionRefToMap(snapshotsItem.Collection)
		snapshotsMap["collection"] = collectionItemMap
	}
	if snapshotsItem.Environment != nil {
		environmentItemMap := resourceIbmAppConfigSnapshotsEnvironmentRefToMap(snapshotsItem.Environment)
		snapshotsMap["environment"] = environmentItemMap
	}
	return snapshotsMap
}

func resourceIbmAppConfigSnapshotsCollectionRefToMap(collectionRef interface{}) []CollectionRef {
	collections := getSnapshotsCollection(collectionRef)
	collectionRefMap := CollectionRef{}
	var collectionMap []CollectionRef
	collectionRefMap["collection_id"] = collections.Collection_id
	collectionRefMap["collection_name"] = collections.Collection_name
	collectionMap = append(collectionMap, collectionRefMap)
	return collectionMap
}

func getSnapshotsCollection(data interface{}) Collections {
	m := data.(map[string]interface{})
	collection := Collections{}
	if name, ok := m["name"].(string); ok {
		collection.Collection_name = name
	}
	if id, ok := m["collection_id"].(string); ok {
		collection.Collection_id = id
	}
	return collection
}

func resourceIbmAppConfigSnapshotsEnvironmentRefToMap(environmentRef interface{}) []EnvironmentRef {
	environments := getSnapshotsEnvironment(environmentRef)
	environmentRefMap := EnvironmentRef{}
	var environmentMap []EnvironmentRef
	environmentRefMap["environment_id"] = environments.Environment_id
	environmentRefMap["environment_name"] = environments.Environment_name
	environmentRefMap["color_code"] = environments.Color_code
	environmentMap = append(environmentMap, environmentRefMap)
	return environmentMap
}

func getSnapshotsEnvironment(data interface{}) Environments {
	m := data.(map[string]interface{})
	environment := Environments{}
	if name, ok := m["name"].(string); ok {
		environment.Environment_name = name
	}
	if id, ok := m["environment_id"].(string); ok {
		environment.Environment_id = id
	}
	if color, ok := m["color_code"].(string); ok {
		environment.Color_code = color
	}
	return environment
}
