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

func dataSourceAlicloudDmsEnterpriseInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDmsEnterpriseInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.ValidateRegexp,
				ForceNew:      true,
				ConflictsWith: []string{"instance_alias_regex"},
			},
			"instance_alias_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.ValidateRegexp,
				ForceNew:      true,
				ConflictsWith: []string{"name_regex"},
			},
			"env_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DELETED", "DISABLE", "NORMAL", "UNAVAILABLE"}, false),
			},
			"tid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_link_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dba_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dba_nick_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ddl_online": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ecs_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"env_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"export_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"query_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"safe_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_dsql": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDmsEnterpriseInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListInstances"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("env_type"); ok {
		request["EnvType"] = v
	}
	if v, ok := d.GetOk("instance_source"); ok {
		request["InstanceSource"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["DbType"] = v
	}
	if v, ok := d.GetOk("net_type"); ok {
		request["NetType"] = v
	}
	if v, ok := d.GetOk("search_key"); ok {
		request["SearchKey"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["InstanceState"] = v
	}
	if v, ok := d.GetOk("tid"); ok {
		request["Tid"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var instanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
	}
	if v, ok := d.GetOk("instance_alias_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		instanceNameRegex = r
	}
	var response map[string]interface{}
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dms_enterprise_instances", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.InstanceList.Instance", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.InstanceList.Instance", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if instanceNameRegex != nil {
				if !instanceNameRegex.MatchString(fmt.Sprint(item["InstanceAlias"])) {
					continue
				}
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
			"id":                fmt.Sprint(object["Host"], ":", formatInt(object["Port"])),
			"data_link_name":    object["DataLinkName"],
			"database_password": object["DatabasePassword"],
			"database_user":     object["DatabaseUser"],
			"dba_id":            object["DbaId"],
			"dba_nick_name":     object["DbaNickName"],
			"ddl_online":        formatInt(object["DdlOnline"]),
			"ecs_instance_id":   object["EcsInstanceId"],
			"ecs_region":        object["EcsRegion"],
			"env_type":          object["EnvType"],
			"export_timeout":    formatInt(object["ExportTimeout"]),
			"host":              object["Host"],
			"instance_id":       object["InstanceId"],
			"instance_name":     object["InstanceAlias"],
			"instance_alias":    object["InstanceAlias"],
			"instance_source":   object["InstanceSource"],
			"instance_type":     object["InstanceType"],
			"port":              formatInt(object["Port"]),
			"query_timeout":     formatInt(object["QueryTimeout"]),
			"safe_rule_id":      object["SafeRuleId"],
			"sid":               object["Sid"],
			"status":            object["State"],
			"use_dsql":          formatInt(object["UseDsql"]),
			"vpc_id":            object["VpcId"],
		}
		ids = append(ids, fmt.Sprint(object["Host"], ":", formatInt(object["Port"])))
		names = append(names, object["InstanceAlias"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
