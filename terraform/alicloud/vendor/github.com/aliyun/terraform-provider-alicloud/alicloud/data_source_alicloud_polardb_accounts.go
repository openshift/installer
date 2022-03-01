package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPolarDBAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPolarDBAccountsRead,

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
			"accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_lock_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_privileges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_privilege": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_name": {
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

func dataSourceAlicloudPolarDBAccountsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := polardb.CreateDescribeAccountsRequest()

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
		return client.DescribeAccounts(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_accounts", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*polardb.DescribeAccountsResponse)
	var ids []string
	var accounts []map[string]interface{}
	if len(response.Accounts) > 0 {
		for _, item := range response.Accounts {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.AccountName) {
					continue
				}
			}
			var nodes []map[string]interface{}
			for _, node := range item.DatabasePrivileges {
				nodeMap := map[string]interface{}{
					"account_privilege": node.AccountPrivilege,
					"db_name":           node.DBName,
				}
				nodes = append(nodes, nodeMap)
			}
			mapping := map[string]interface{}{
				"account_description": item.AccountDescription,
				"account_lock_state":  item.AccountLockState,
				"account_name":        item.AccountName,
				"account_status":      item.AccountStatus,
				"account_type":        item.AccountType,
				"database_privileges": nodes,
			}
			ids = append(ids, item.AccountName)
			accounts = append(accounts, mapping)
		}
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("accounts", accounts); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), accounts)
	}
	return nil
}
