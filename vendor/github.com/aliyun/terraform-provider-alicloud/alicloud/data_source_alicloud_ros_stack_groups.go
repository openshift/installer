package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRosStackGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRosStackGroupsRead,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "DELETED"}, false),
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
						"administration_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execution_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parameter_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"stack_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_body": {
							Type:     schema.TypeString,
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

func dataSourceAlicloudRosStackGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListStackGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var stackGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		stackGroupNameRegex = r
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
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ros_stack_groups", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.StackGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.StackGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if stackGroupNameRegex != nil {
				if !stackGroupNameRegex.MatchString(item["StackGroupName"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["StackGroupName"])]; !ok {
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
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"description":      object["Description"],
			"stack_group_id":   object["StackGroupId"],
			"id":               fmt.Sprint(object["StackGroupName"]),
			"stack_group_name": fmt.Sprint(object["StackGroupName"]),
			"status":           object["Status"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["StackGroupName"]))
			names = append(names, object["StackGroupName"].(string))
			s = append(s, mapping)
			continue
		}

		rosService := RosService{client}
		id := fmt.Sprint(object["StackGroupName"])
		getResp, err := rosService.DescribeRosStackGroup(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["administration_role_name"] = getResp["AdministrationRoleName"]
		mapping["execution_role_name"] = getResp["ExecutionRoleName"]

		parameters := make([]map[string]interface{}, 0)
		if parametersList, ok := getResp["Parameters"].([]interface{}); ok {
			for _, v := range parametersList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"parameter_key":   m1["ParameterKey"],
						"parameter_value": m1["ParameterValue"],
					}
					parameters = append(parameters, temp1)
				}
			}
		}
		mapping["parameters"] = parameters

		b, err := json.Marshal(getResp["TemplateBody"])
		if err != nil {
			return WrapError(err)
		}
		mapping["template_body"] = string(b)

		ids = append(ids, fmt.Sprint(object["StackGroupName"]))
		names = append(names, object["StackGroupName"].(string))
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
