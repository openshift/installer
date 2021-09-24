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

func dataSourceAlicloudRosChangeSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRosChangeSetsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"change_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CREATE_COMPLETE", "CREATE_FAILED", "CREATE_IN_PROGRESS", "CREATE_PENDING", "DELETE_COMPLETE", "DELETE_FAILED"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"change_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"change_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"change_set_type": {
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
						"execution_status": {
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
						"stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_name": {
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

func dataSourceAlicloudRosChangeSetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListChangeSets"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("change_set_name"); ok {
		request["ChangeSetName"] = v
	}
	request["RegionId"] = client.RegionId
	request["StackId"] = d.Get("stack_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var changeSetNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		changeSetNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ros_change_sets", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.ChangeSets", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ChangeSets", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if changeSetNameRegex != nil {
				if !changeSetNameRegex.MatchString(item["ChangeSetName"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ChangeSetId"])]; !ok {
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
			"id":              fmt.Sprint(object["ChangeSetId"]),
			"change_set_id":   object["ChangeSetId"],
			"change_set_name": object["ChangeSetName"],
			"change_set_type": object["ChangeSetType"],
			"description":     object["Description"],
			"stack_id":        object["StackId"],
			"stack_name":      object["StackName"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["ChangeSetId"]))
			names = append(names, object["ChangeSetName"].(string))
			s = append(s, mapping)
			continue
		}

		rosService := RosService{client}
		id := fmt.Sprint(object["ChangeSetId"])
		getResp, err := rosService.DescribeRosChangeSet(id)
		if err != nil {
			return WrapError(err)
		}
		if statusOk && status != "" && status != getResp["Status"].(string) {
			continue
		}
		mapping["disable_rollback"] = getResp["DisableRollback"]
		mapping["execution_status"] = getResp["ExecutionStatus"]

		Parameters := getResp["Parameters"].([]interface{})
		parameters := make([]map[string]interface{}, len(Parameters))
		for i, v := range Parameters {
			mapping1 := map[string]interface{}{
				"parameter_key":   v.(map[string]interface{})["Parameters"],
				"parameter_value": v.(map[string]interface{})["Parameters"],
			}
			parameters[i] = mapping1
		}
		mapping["parameters"] = getResp["parameters"]
		mapping["status"] = getResp["Status"]
		mapping["template_body"] = getResp["TemplateBody"]
		mapping["timeout_in_minutes"] = getResp["TimeoutInMinutes"]
		ids = append(ids, fmt.Sprint(object["ChangeSetId"]))
		names = append(names, object["ChangeSetName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("sets", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
