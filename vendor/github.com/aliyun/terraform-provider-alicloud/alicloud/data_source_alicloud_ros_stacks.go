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

func dataSourceAlicloudRosStacks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRosStacksRead,
		Schema: map[string]*schema.Schema{
			"parent_stack_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"show_nested_stack": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
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
			"stack_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CREATE_COMPLETE", "CREATE_FAILED", "CREATE_IN_PROGRESS", "DELETE_COMPLETE", "DELETE_FAILED", "DELETE_IN_PROGRESS", "ROLLBACK_COMPLETE", "ROLLBACK_FAILED", "ROLLBACK_IN_PROGRESS"}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stacks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deletion_protection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disable_rollback": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"drift_detection_time": {
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
						"parent_stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_drift_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_policy_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"template_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout_in_minutes": {
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

func dataSourceAlicloudRosStacksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListStacks"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("parent_stack_id"); ok {
		request["ParentStackId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("show_nested_stack"); ok {
		request["ShowNestedStack"] = v
	}
	if v, ok := d.GetOk("stack_name"); ok {
		request["StackName"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var stackNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		stackNameRegex = r
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
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ros_stacks", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Stacks", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Stacks", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if stackNameRegex != nil {
				if !stackNameRegex.MatchString(item["StackName"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["StackId"])]; !ok {
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
			"disable_rollback":     object["DisableRollback"],
			"drift_detection_time": object["DriftDetectionTime"],
			"parent_stack_id":      object["ParentStackId"],
			"stack_drift_status":   object["StackDriftStatus"],
			"id":                   fmt.Sprint(object["StackId"]),
			"stack_id":             fmt.Sprint(object["StackId"]),
			"stack_name":           object["StackName"],
			"status_reason":        object["StatusReason"],
			"timeout_in_minutes":   formatInt(object["TimeoutInMinutes"]),
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["StackId"]))
			names = append(names, object["StackName"].(string))
			s = append(s, mapping)
			continue
		}

		rosService := RosService{client}
		id := fmt.Sprint(object["StackId"])
		getResp, err := rosService.DescribeRosStack(id)
		if err != nil {
			return WrapError(err)
		}
		if statusOk && status != "" && status != getResp["Status"].(string) {
			continue
		}
		mapping["deletion_protection"] = getResp["DeletionProtection"]
		mapping["description"] = getResp["Description"]

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
		mapping["ram_role_name"] = getResp["RamRoleName"]
		mapping["root_stack_id"] = getResp["RootStackId"]
		mapping["status"] = getResp["Status"]
		mapping["template_description"] = getResp["TemplateDescription"]
		getResp1, err := rosService.GetStackPolicy(id)
		if err != nil {
			return WrapError(err)
		}
		b, err := json.Marshal(getResp1["StackPolicyBody"])
		mapping["stack_policy_body"] = string(b)

		getResp2, err := rosService.ListTagResources(id, "stack")
		if err != nil {
			return WrapError(err)
		}
		mapping["tags"] = tagsToMap(getResp2)

		ids = append(ids, fmt.Sprint(object["StackId"]))
		names = append(names, object["StackName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("stacks", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
