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

func dataSourceAlicloudHbrRestoreJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrRestoreJobsRead,
		Schema: map[string]*schema.Schema{
			"restore_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ECS_FILE", "NAS", "OSS"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CANCELED", "CANCELING", "COMPLETE", "CREATED", "EXPIRED", "FAILED", "PARTIAL_COMPLETE", "QUEUED", "RUNNING"}, false),
			},
			"restore_id": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_type": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_bucket": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_file_system_id": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_instance_id": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vault_id": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"actual_bytes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"actual_items": {
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
						"error_file": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
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
						"options": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"restore_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restore_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
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
						"target_bucket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_data_source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_file_system_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_prefix": {
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrRestoreJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeRestoreJobs2"
	request := make(map[string]interface{})
	request["RestoreType"] = d.Get("restore_type")

	filtersMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("vault_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "vaultId",
			"Values": v,
		})
	}
	if v, ok := d.GetOk("restore_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "restoreId",
			"Values": v,
		})
	}
	if v, ok := d.GetOk("source_type"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "sourceType",
			"Values": v,
		})
	}
	if v, ok := d.GetOk("target_bucket"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "targetBucket",
			"Values": v,
		})
	}
	if v, ok := d.GetOk("target_file_system_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "targetFileSystemId",
			"Values": v,
		})
	}
	if v, ok := d.GetOk("target_instance_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"Key":    "targetInstanceId",
			"Values": v,
		})
	}
	request["Filters"] = filtersMapList
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_restore_jobs", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.RestoreJobs.RestoreJob", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.RestoreJobs.RestoreJob", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
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
			"id":                    fmt.Sprint(object["RestoreId"], ":", object["RestoreType"]),
			"actual_bytes":          fmt.Sprint(object["ActualBytes"]),
			"actual_items":          fmt.Sprint(object["ActualItems"]),
			"bytes_done":            fmt.Sprint(object["BytesDone"]),
			"bytes_total":           fmt.Sprint(object["BytesTotal"]),
			"complete_time":         fmt.Sprint(object["CompleteTime"]),
			"create_time":           fmt.Sprint(object["CreatedTime"]),
			"error_file":            object["ErrorFile"],
			"error_message":         object["ErrorMessage"],
			"expire_time":           fmt.Sprint(object["ExpireTime"]),
			"items_done":            fmt.Sprint(object["ItemsDone"]),
			"items_total":           fmt.Sprint(object["ItemsTotal"]),
			"options":               object["Options"],
			"parent_id":             object["ParentId"],
			"progress":              formatInt(object["Progress"]),
			"restore_job_id":        object["RestoreId"],
			"restore_type":          object["RestoreType"],
			"snapshot_hash":         object["SnapshotHash"],
			"snapshot_id":           object["SnapshotId"],
			"source_type":           object["SourceType"],
			"start_time":            fmt.Sprint(object["StartTime"]),
			"status":                object["Status"],
			"target_bucket":         object["TargetBucket"],
			"target_client_id":      object["TargetClientId"],
			"target_create_time":    fmt.Sprint(object["TargetCreateTime"]),
			"target_data_source_id": object["TargetDataSourceId"],
			"target_file_system_id": object["TargetFileSystemId"],
			"target_instance_id":    object["TargetInstanceId"],
			"target_path":           object["TargetPath"],
			"target_prefix":         object["TargetPrefix"],
			"updated_time":          fmt.Sprint(object["UpdatedTime"]),
			"vault_id":              object["VaultId"],
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
