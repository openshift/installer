package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRdsAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsAccountsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Unavailable"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_name": {
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
									"account_privilege_detail": {
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
						"priv_exceeded": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRdsAccountsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeAccounts"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var accountNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		accountNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_accounts", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Accounts.DBInstanceAccount", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Accounts.DBInstanceAccount", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if accountNameRegex != nil {
				if !accountNameRegex.MatchString(fmt.Sprint(item["AccountName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AccountName"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["AccountStatus"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_description": object["AccountDescription"],
			"id":                  fmt.Sprint(object["AccountName"]),
			"account_name":        fmt.Sprint(object["AccountName"]),
			"account_type":        object["AccountType"],
			"priv_exceeded":       object["PrivExceeded"],
			"status":              object["AccountStatus"],
		}

		databasePrivilege := make([]map[string]interface{}, 0)
		if databasePrivilegeList, ok := object["DatabasePrivileges"].(map[string]interface{})["DatabasePrivilege"].([]interface{}); ok {
			for _, v := range databasePrivilegeList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"account_privilege":        m1["AccountPrivilege"],
						"account_privilege_detail": m1["AccountPrivilegeDetail"],
						"db_name":                  m1["DBName"],
					}
					databasePrivilege = append(databasePrivilege, temp1)
				}
			}
		}
		mapping["database_privileges"] = databasePrivilege
		ids = append(ids, fmt.Sprint(object["AccountName"]))
		names = append(names, object["AccountName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("accounts", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
