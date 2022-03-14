package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudNasDataFlows() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasDataFlowsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Starting", "Running", "Stopped", "Updating", "Deleting", "Stopping", "Misconfigured"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_flow_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fset_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fset_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_security_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"throughput": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNasDataFlowsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDataFlows"
	request := make(map[string]interface{})
	request["FileSystemId"] = d.Get("file_system_id")
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_data_flows", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DataFlowInfo.DataFlow", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DataFlowInfo.DataFlow", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["FileSystemId"], ":", item["DataFlowId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
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
		mapping := map[string]interface{}{
			"create_time":          object["CreateTime"],
			"id":                   fmt.Sprint(object["FileSystemId"], ":", object["DataFlowId"]),
			"data_flow_id":         fmt.Sprint(object["DataFlowId"]),
			"description":          object["Description"],
			"error_message":        object["ErrorMessage"],
			"file_system_id":       object["FileSystemId"],
			"file_system_path":     object["FileSystemPath"],
			"fset_description":     object["FsetDescription"],
			"fset_id":              object["FsetId"],
			"source_security_type": object["SourceSecurityType"],
			"source_storage":       object["SourceStorage"],
			"status":               object["Status"],
			"throughput":           formatInt(object["Throughput"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("flows", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
