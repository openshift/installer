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

func dataSourceAlicloudNasAccessGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasAccessGroupsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"type": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'type' has been deprecated from provider version 1.95.0. New field 'access_group_type' replaces it.",
			},
			"access_group_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"access_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"extreme", "standard"}, false),
			},
			"useutc_date_time": {
				Type:     schema.TypeBool,
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
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNasAccessGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeAccessGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("access_group_name"); ok {
		request["AccessGroupName"] = v
	}
	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v
	}
	if v, ok := d.GetOkExists("useutc_date_time"); ok {
		request["UseUTCDateTime"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var accessGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		accessGroupNameRegex = r
	}
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_access_groups", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.AccessGroups.AccessGroup", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AccessGroups.AccessGroup", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if accessGroupNameRegex != nil {
				if !accessGroupNameRegex.MatchString(fmt.Sprint(item["AccessGroupName"])) {
					continue
				}
			}
			if v, ok := d.GetOk("type"); ok && v.(string) != "" && item["AccessGroupType"].(string) != v.(string) {
				continue
			}
			if v, ok := d.GetOk("access_group_type"); ok && v.(string) != "" && item["AccessGroupType"].(string) != v.(string) {
				continue
			}
			if v, ok := d.GetOk("description"); ok && v.(string) != "" && item["Description"].(string) != v.(string) {
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
			"id":                 fmt.Sprint(object["AccessGroupName"], ":", request["FileSystemType"]),
			"access_group_name":  object["AccessGroupName"],
			"access_group_type":  object["AccessGroupType"],
			"type":               object["AccessGroupType"],
			"description":        object["Description"],
			"mount_target_count": formatInt(object["MountTargetCount"]),
			"rule_count":         formatInt(object["RuleCount"]),
		}
		ids = append(ids, fmt.Sprint(object["AccessGroupName"], ":", request["FileSystemType"]))
		names = append(names, object["AccessGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
