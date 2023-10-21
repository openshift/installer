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

func dataSourceAlicloudEcdCommands() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdCommandsRead,
		Schema: map[string]*schema.Schema{
			"command_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RunBatScript", "RunPowerShellScript"}, false),
			},
			"content_encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PlainText",
				ValidateFunc: validation.StringInSlice([]string{"Base64", "PlainText"}, false),
			},
			"desktop_id": {
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Pending", "Failed", "PartialFailed", "Running", "Stopped", "Stopping", "Finished", "Success"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"commands": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invoke_desktops": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"desktop_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dropped": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"error_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_info": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"exit_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"finish_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"invocation_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"output": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"repeats": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"stop_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invoke_id": {
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

func dataSourceAlicloudEcdCommandsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeInvocations"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("command_type"); ok {
		request["CommandType"] = v
	}
	if v, ok := d.GetOk("content_encoding"); ok {
		request["ContentEncoding"] = v
	}

	if v, ok := d.GetOk("desktop_id"); ok {
		request["DesktopId"] = v
	}

	request["IncludeOutput"] = "true"

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["InvokeStatus"] = v
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
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_commands", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Invocations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Invocations", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InvokeId"])]; !ok {
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
		mapping := map[string]interface{}{
			"command_content": object["CommandContent"],
			"command_type":    object["CommandType"],
			"create_time":     object["CreationTime"],
			"id":              fmt.Sprint(object["InvokeId"]),
			"invoke_id":       fmt.Sprint(object["InvokeId"]),
			"status":          object["InvocationStatus"],
		}

		invokeDesktops := make([]map[string]interface{}, 0)
		if invokeDesktopsList, ok := object["InvokeDesktops"].([]interface{}); ok {
			for _, v := range invokeDesktopsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"desktop_id":        m1["DesktopId"],
						"dropped":           formatInt(m1["Dropped"]),
						"error_code":        m1["ErrorCode"],
						"error_info":        m1["ErrorInfo"],
						"exit_code":         m1["ExitCode"],
						"finish_time":       m1["FinishTime"],
						"invocation_status": m1["InvocationStatus"],
						"output":            m1["Output"],
						"repeats":           formatInt(m1["Repeats"]),
						"start_time":        m1["StartTime"],
						"stop_time":         m1["StopTime"],
					}
					invokeDesktops = append(invokeDesktops, temp1)
				}
			}
		}
		mapping["invoke_desktops"] = invokeDesktops
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("commands", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
