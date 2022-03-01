package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEhpcJobTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEhpcJobTemplatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"array_request": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clock_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_line": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mem": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"package_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"re_runable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"runas_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stderr_redirect_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stdout_redirect_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"thread": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"variables": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEhpcJobTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListJobTemplates"
	request := make(map[string]interface{})

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
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
	conn, err := client.NewEhpcClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-04-12"), StringPointer("AK"), request, nil, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ehpc_job_templates", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Templates.JobTemplates", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Templates.JobTemplates", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"array_request":        object["ArrayRequest"],
			"clock_time":           object["ClockTime"],
			"command_line":         object["CommandLine"],
			"gpu":                  formatInt(object["Gpu"]),
			"id":                   fmt.Sprint(object["Id"]),
			"job_template_id":      fmt.Sprint(object["Id"]),
			"job_template_name":    object["Name"],
			"mem":                  object["Mem"],
			"node":                 formatInt(object["Node"]),
			"package_path":         object["PackagePath"],
			"priority":             formatInt(object["Priority"]),
			"queue":                object["Queue"],
			"re_runable":           object["ReRunable"],
			"runas_user":           object["RunasUser"],
			"stderr_redirect_path": object["StderrRedirectPath"],
			"stdout_redirect_path": object["StdoutRedirectPath"],
			"task":                 formatInt(object["Task"]),
			"thread":               formatInt(object["Thread"]),
			"variables":            object["Variables"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("templates", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
