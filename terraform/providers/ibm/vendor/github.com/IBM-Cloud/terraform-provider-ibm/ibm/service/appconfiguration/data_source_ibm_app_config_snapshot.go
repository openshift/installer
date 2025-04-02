package appconfiguration

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppConfigSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmAppConfigSnapshotRead,

		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GUID of the App Configuration service. Get it from the service instance credentials section of the dashboard.",
			},
			"git_config_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Git config id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only",
			},
			"git_config_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Git config name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only",
			},
			"git_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Git url which will be used to connect to the github account.",
			},
			"git_branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Branch name to which you need to write or update the configuration.",
			},
			"git_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
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
	}
}

func dataSourceIbmAppConfigSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	guid := d.Get("guid").(string)

	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf("getAppConfigClient failed %s", err)
	}

	options := &appconfigurationv1.GetGitconfigOptions{}
	options.SetGitConfigID(d.Get("git_config_id").(string))

	result, response, err := appconfigClient.GetGitconfig(options)

	if err != nil {
		return flex.FmtErrorf("GetGitconfig failed %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s", guid, *result.GitConfigID))

	if result.GitConfigName != nil {
		if err = d.Set("git_config_name", result.GitConfigName); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting git_config_name: %s", err)
		}
	}
	if result.GitConfigID != nil {
		if err = d.Set("git_config_id", result.GitConfigID); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting git_config_id: %s", err)
		}
	}
	if result.GitURL != nil {
		if err = d.Set("git_url", result.GitURL); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting git_url: %s", err)
		}
	}
	if result.GitBranch != nil {
		if err = d.Set("git_branch", result.GitBranch); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting git_branch: %s", err)
		}
	}
	if result.GitFilePath != nil {
		if err = d.Set("git_file_path", result.GitFilePath); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting git_file_path: %s", err)
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
	if result.LastSyncTime != nil {
		if err = d.Set("last_sync_time", result.LastSyncTime.String()); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting last_sync_time: %s", err)
		}
	}
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting href: %s", err)
		}
	}
	if result.Collection != nil {
		collectionItemMap := resourceIbmAppConfigSnapshotCollectionRefToMap(result.Collection)
		if err = d.Set("collection", collectionItemMap); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting collection: %s", err)
		}
	}
	if result.Environment != nil {
		environmentItemMap := resourceIbmAppConfigSnapshotEnvironmentRefToMap(result.Environment)
		if err = d.Set("environment", environmentItemMap); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting environment: %s", err)
		}
	}
	return nil
}

type CollectionRef map[string]interface{}

type Collections struct {
	Collection_name string `json:"name"`
	Collection_id   string `json:"collection_id"`
}

func resourceIbmAppConfigSnapshotCollectionRefToMap(collectionRef interface{}) []CollectionRef {
	collections := getSnapshotCollection(collectionRef)
	collectionRefMap := CollectionRef{}
	var collectionMap []CollectionRef
	collectionRefMap["collection_id"] = collections.Collection_id
	collectionRefMap["collection_name"] = collections.Collection_name
	collectionMap = append(collectionMap, collectionRefMap)
	return collectionMap
}

func getSnapshotCollection(data interface{}) Collections {
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

type EnvironmentRef map[string]interface{}

type Environments struct {
	Environment_name string `json:"name"`
	Environment_id   string `json:"environment_id"`
	Color_code       string `json:"color_code"`
}

func resourceIbmAppConfigSnapshotEnvironmentRefToMap(environmentRef interface{}) []EnvironmentRef {
	environments := getSnapshotEnvironment(environmentRef)
	environmentRefMap := EnvironmentRef{}
	var environmentMap []EnvironmentRef
	environmentRefMap["environment_id"] = environments.Environment_id
	environmentRefMap["environment_name"] = environments.Environment_name
	environmentRefMap["color_code"] = environments.Color_code
	environmentMap = append(environmentMap, environmentRefMap)
	return environmentMap
}

func getSnapshotEnvironment(data interface{}) Environments {
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
