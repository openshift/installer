package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOosExecutions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOosExecutionsRead,
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AlarmTrigger", "EventTrigger", "Other", "TimerTrigger"}, false),
			},
			"end_date": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_date_after": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"executed_by": {
				Type:     schema.TypeString,
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
			"include_child_execution": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Automatic",
				ValidateFunc: validation.StringInSlice([]string{"Automatic", "Debug"}, false),
			},
			"parent_execution_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ram_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sort_field": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sort_order": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_date_after": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_date_before": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Cancelled", "Failed", "Queued", "Running", "Started", "Success", "Waiting"}, false),
			},
			"tags": tagsSchema(),
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"executions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"counters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"executed_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execution_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_parent": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"outputs": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_execution_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOosExecutionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListExecutions"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	if v, ok := d.GetOk("end_date"); ok {
		request["EndDateBefore"] = v
	}
	if v, ok := d.GetOk("end_date_after"); ok {
		request["EndDateAfter"] = v
	}
	if v, ok := d.GetOk("executed_by"); ok {
		request["ExecutedBy"] = v
	}
	if v, ok := d.GetOkExists("include_child_execution"); ok {
		request["IncludeChildExecution"] = v
	}
	if v, ok := d.GetOk("mode"); ok {
		request["Mode"] = v
	}
	if v, ok := d.GetOk("parent_execution_id"); ok {
		request["ParentExecutionId"] = v
	}
	if v, ok := d.GetOk("ram_role"); ok {
		request["RamRole"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("sort_field"); ok {
		request["SortField"] = v
	}
	if v, ok := d.GetOk("sort_order"); ok {
		request["SortOrder"] = v
	}
	if v, ok := d.GetOk("start_date_after"); ok {
		request["StartDateAfter"] = v
	}
	if v, ok := d.GetOk("start_date_before"); ok {
		request["StartDateBefore"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	if v, ok := d.GetOk("template_name"); ok {
		request["TemplateName"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}

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
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_oos_executions", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Executions", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Executions", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ExecutionId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		parameters, _ := convertMaptoJsonString(object["Parameters"].(map[string]interface{}))
		counters, _ := convertMapFloat64ToJsonString(object["Counters"].(map[string]interface{}))
		mapping := map[string]interface{}{
			"category":            object["Category"],
			"counters":            counters,
			"create_date":         object["CreateDate"],
			"end_date":            object["EndDate"],
			"executed_by":         object["ExecutedBy"],
			"id":                  fmt.Sprint(object["ExecutionId"]),
			"execution_id":        fmt.Sprint(object["ExecutionId"]),
			"is_parent":           object["IsParent"],
			"mode":                object["Mode"],
			"outputs":             object["Outputs"],
			"parameters":          parameters,
			"parent_execution_id": object["ParentExecutionId"],
			"ram_role":            object["RamRole"],
			"start_date":          object["StartDate"],
			"status":              object["Status"],
			"status_message":      object["StatusMessage"],
			"status_reason":       object["StatusReason"],
			"template_id":         object["TemplateId"],
			"template_name":       object["TemplateName"],
			"template_version":    object["TemplateVersion"],
			"update_date":         object["UpdateDate"],
		}
		ids = append(ids, fmt.Sprint(object["ExecutionId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("executions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
