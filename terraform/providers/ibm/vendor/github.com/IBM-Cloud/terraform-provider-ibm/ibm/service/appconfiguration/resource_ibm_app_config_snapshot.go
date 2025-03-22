package appconfiguration

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIbmAppConfigSnapshot() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIbmIbmAppConfigSnapshotCreate,
		Read:     resourceIbmIbmAppConfigSnapshotRead,
		Update:   resourceIbmIbmAppConfigSnapshotUpdate,
		Delete:   resourceIbmIbmAppConfigSnapshotDelete,
		Importer: &schema.ResourceImporter{},

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
				Required:    true,
				Description: "Git config name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only",
			},
			"git_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Git url which will be used to connect to the github account.",
			},
			"git_branch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Branch name to which you need to write or update the configuration.",
			},
			"git_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Git file path, this is a path where your configuration file will be written.",
			},
			"git_token": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Required:    true,
				Description: "Git token, this needs to be provided with enough permission to write and update the file.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the git config.",
			},
			"collection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Collection id.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "action promote",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Environment id.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of the git config data.",
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

func resourceIbmIbmAppConfigSnapshotCreate(d *schema.ResourceData, meta interface{}) error {

	guid := d.Get("guid").(string)
	appconfigClient, err := getAppConfigClient(meta, guid)
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}
	options := &appconfigurationv1.CreateGitconfigOptions{}

	options.SetGitConfigName(d.Get("git_config_name").(string))
	options.SetGitConfigID(d.Get("git_config_id").(string))

	options.SetCollectionID(d.Get("collection_id").(string))
	options.SetEnvironmentID(d.Get("environment_id").(string))
	options.SetGitURL(d.Get("git_url").(string))
	options.SetGitBranch(d.Get("git_branch").(string))
	options.SetGitFilePath(d.Get("git_file_path").(string))
	options.SetGitToken(d.Get("git_token").(string))

	snapshot, response, err := appconfigClient.CreateGitconfig(options)

	if err != nil {
		return flex.FmtErrorf("CreateGitconfig failed %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", guid, *snapshot.GitConfigID))
	return resourceIbmIbmAppConfigSnapshotRead(d, meta)
}

func resourceIbmIbmAppConfigSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}

	if ok := d.HasChanges("action"); ok {
		option := &appconfigurationv1.PromoteGitconfigOptions{}
		option.SetGitConfigID(parts[1])
		_, response, err := appconfigClient.PromoteGitconfig(option)
		if err != nil {
			return flex.FmtErrorf("[ERROR] PromoteGitconfig %s\n%s", err, response)
		}
		return resourceIbmIbmAppConfigSnapshotRead(d, meta)
	} else {
		if ok := d.HasChanges("git_config_name", "collection_id", "environment_id", "git_url", "git_branch", "git_file_path", "git_token"); ok {
			options := &appconfigurationv1.UpdateGitconfigOptions{}
			options.SetGitConfigID(parts[1])
			if _, ok := GetFieldExists(d, "git_config_name"); ok {
				options.SetGitConfigName(d.Get("git_config_name").(string))
			}
			if _, ok := GetFieldExists(d, "collection_id"); ok {
				options.SetCollectionID(d.Get("collection_id").(string))
			}
			if _, ok := GetFieldExists(d, "environment_id"); ok {
				options.SetEnvironmentID(d.Get("environment_id").(string))
			}
			if _, ok := GetFieldExists(d, "git_url"); ok {
				options.SetGitURL(d.Get("git_url").(string))
			}
			if _, ok := GetFieldExists(d, "git_branch"); ok {
				options.SetGitBranch(d.Get("git_branch").(string))
			}
			if _, ok := GetFieldExists(d, "git_file_path"); ok {
				options.SetGitFilePath(d.Get("git_file_path").(string))
			}
			if _, ok := GetFieldExists(d, "git_token"); ok {
				options.SetGitToken(d.Get("git_token").(string))
			}
			_, response, err := appconfigClient.UpdateGitconfig(options)
			if err != nil {
				return flex.FmtErrorf("[ERROR] UpdateGitconfig %s\n%s", err, response)
			}
			return resourceIbmIbmAppConfigSnapshotRead(d, meta)
		}
	}
	return nil
}

func resourceIbmIbmAppConfigSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}
	if len(parts) != 2 {
		return flex.FmtErrorf("Kindly check the id")
	}

	options := &appconfigurationv1.GetGitconfigOptions{}
	options.SetGitConfigID(parts[1])

	result, response, err := appconfigClient.GetGitconfig(options)
	if err != nil {
		return flex.FmtErrorf("[ERROR] GetGitconfigs failed %s\n%s", err, response)
	}

	d.Set("guid", parts[0])
	d.Set("git_config_id", parts[1])
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
	if result.Href != nil {
		if err = d.Set("href", result.Href); err != nil {
			return flex.FmtErrorf("[ERROR] Error setting href: %s", err)
		}
	}
	return nil
}

func resourceIbmIbmAppConfigSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil
	}
	appconfigClient, err := getAppConfigClient(meta, parts[0])
	if err != nil {
		return flex.FmtErrorf(fmt.Sprintf("%s", err))
	}

	options := &appconfigurationv1.DeleteGitconfigOptions{}
	options.SetGitConfigID(parts[1])

	response, err := appconfigClient.DeleteGitconfig(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return flex.FmtErrorf("[ERROR] DeleteGitconfig failed %s\n%s", err, response)
	}
	d.SetId("")

	return nil
}
