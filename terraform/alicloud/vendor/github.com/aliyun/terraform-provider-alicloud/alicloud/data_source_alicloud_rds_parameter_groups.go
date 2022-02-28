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

func dataSourceAlicloudRdsParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsParameterGroupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"force_restart": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"param_counts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"param_detail": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"param_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"parameter_group_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudRdsParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeParameterGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var parameterGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		parameterGroupNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_parameter_groups", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.ParameterGroups.ParameterGroup", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ParameterGroups.ParameterGroup", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if parameterGroupNameRegex != nil {
			if !parameterGroupNameRegex.MatchString(fmt.Sprint(item["ParameterGroupName"])) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ParameterGroupId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"engine":               object["Engine"],
			"engine_version":       object["EngineVersion"],
			"force_restart":        formatInt(object["ForceRestart"]),
			"param_counts":         formatInt(object["ParamCounts"]),
			"parameter_group_desc": object["ParameterGroupDesc"],
			"id":                   fmt.Sprint(object["ParameterGroupId"]),
			"parameter_group_id":   fmt.Sprint(object["ParameterGroupId"]),
			"parameter_group_name": object["ParameterGroupName"],
			"parameter_group_type": formatInt(object["ParameterGroupType"]),
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["ParameterGroupId"]))
			names = append(names, object["ParameterGroupName"])
			s = append(s, mapping)
			continue
		}

		rdsService := RdsService{client}
		id := fmt.Sprint(object["ParameterGroupId"])
		getResp, err := rdsService.DescribeRdsParameterGroup(id)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := getResp["ParamDetail"].(map[string]interface{})["ParameterDetail"].([]interface{}); ok {
			parameterDetail := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				parameterDetail = append(parameterDetail, map[string]interface{}{
					"param_name":  item["ParamName"],
					"param_value": item["ParamValue"],
				})
			}
			mapping["param_detail"] = parameterDetail
		}
		ids = append(ids, fmt.Sprint(object["ParameterGroupId"]))
		names = append(names, object["ParameterGroupName"])
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
