package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPolarDBDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBDatabasesRead,

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"character_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"account_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"privilege_status": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataSourceAlicloudPolarDBDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := polardb.CreateDescribeDatabasesRequest()

	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("db_cluster_id").(string)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

	raw, err := client.WithPolarDBClient(func(client *polardb.Client) (interface{}, error) {
		return client.DescribeDatabases(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_databases", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeDatabasesResponse)
	var ids []string
	var databases []map[string]interface{}
	if len(response.Databases.Database) > 0 {
		for _, item := range response.Databases.Database {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBName) {
					continue
				}
			}
			var nodes []map[string]interface{}
			for _, node := range item.Accounts.Account {
				nodeMap := map[string]interface{}{
					"account_name":     node.AccountName,
					"account_status":   node.AccountStatus,
					"privilege_status": node.PrivilegeStatus,
				}
				nodes = append(nodes, nodeMap)
			}
			mapping := map[string]interface{}{
				"character_set_name": item.CharacterSetName,
				"db_description":     item.DBDescription,
				"db_name":            item.DBName,
				"db_status":          item.DBStatus,
				"engine":             item.Engine,
				"accounts":           nodes,
			}
			ids = append(ids, item.DBName)
			databases = append(databases, mapping)
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("databases", databases); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), databases)
	}
	return nil
}
