package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudHbrBackupJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrBackupJobsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"sort_direction": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ASCEND", "DESCEND"}, false),
			},
			"source_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ECS_FILE", "NAS", "OSS"}, false),
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"EQUAL", "NOT_EQUAL", "GREATER_THAN", "GREATER_THAN_OR_EQUAL", "LESS_THAN", "LESS_THAN_OR_EQUAL", "BETWEEN", "IN"}, false),
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE", "FAILED", "PARTIAL_COMPLETE", "UNAVAILABLE"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actual_bytes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"actual_items": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"back_job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bytes_done": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bytes_total": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"complete_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"include": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"items_done": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"items_total": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nas_create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"paths": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vault_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrBackupJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBackupJobs2"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("sort_direction"); ok {
		request["SortDirection"] = v
	}
	request["SourceType"] = d.Get("source_type")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("filter"); ok {
		for filterPtr, filter := range v.([]interface{}) {
			filterArg := filter.(map[string]interface{})
			request[fmt.Sprintf("Filters.%d.Key", filterPtr+1)] = filterArg["key"]
			request[fmt.Sprintf("Filters.%d.Operator", filterPtr+1)] = filterArg["operator"]
			if filterArg["key"] == "CompleteTime" {
				for valuePtr, value := range filterArg["values"].([]interface{}) {
					request[fmt.Sprintf("Filters.%d.Values.%d", filterPtr+1, valuePtr+1)] = strconv.FormatInt(ConvertNasFileSystemStringToUnix(value.(string)), 10)
				}
			} else {
				for valuePtr, value := range filterArg["values"].([]interface{}) {
					request[fmt.Sprintf("Filters.%d.Values.%d", filterPtr+1, valuePtr+1)] = value
				}
			}
		}
	}

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
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_backup_jobs", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.BackupJobs.BackupJob", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BackupJobs.BackupJob", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["JobId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"actual_bytes":    fmt.Sprint(object["ActualBytes"]),
			"actual_items":    fmt.Sprint(object["ActualItems"]),
			"back_job_name":   object["JobName"],
			"id":              fmt.Sprint(object["JobId"]),
			"backup_job_id":   fmt.Sprint(object["JobId"]),
			"backup_type":     object["BackupType"],
			"bucket":          object["Bucket"],
			"bytes_done":      fmt.Sprint(object["BytesDone"]),
			"bytes_total":     fmt.Sprint(object["BytesTotal"]),
			"complete_time":   fmt.Sprint(object["CompleteTime"]),
			"create_time":     fmt.Sprint(object["CreatedTime"]),
			"exclude":         object["Exclude"],
			"file_system_id":  object["FileSystemId"],
			"include":         object["Include"],
			"instance_id":     object["InstanceId"],
			"items_done":      fmt.Sprint(object["ItemsDone"]),
			"items_total":     fmt.Sprint(object["ItemsTotal"]),
			"progress":        fmt.Sprint(object["Progress"]),
			"nas_create_time": fmt.Sprint(object["CreateTime"]),
			"plan_id":         object["PlanId"],
			"prefix":          object["Prefix"],
			"source_type":     object["SourceType"],
			"error_message":   object["ErrorMessage"],
			"start_time":      fmt.Sprint(object["StartTime"]),
			"status":          object["Status"],
			"updated_time":    fmt.Sprint(object["UpdatedTime"]),
			"vault_id":        object["VaultId"],
		}
		if object["Path"] != nil {
			mapping["path"] = object["Paths"].(map[string]interface{})["Path"]
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("jobs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
