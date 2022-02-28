package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudHbrSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"query": {
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
			"source_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ECS_FILE", "OSS", "NAS"}, false),
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"complete_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"complete_time_checker": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MATCH_TERM", "GREATER_THAN", "GREATER_THAN_OR_EQUAL", "LESS_THAN", "LESS_THAN_OR_EQUAL", "BETWEEN", "IN", "NOT_IN"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE", "PARTIAL_COMPLETE"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshots": {
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
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_file": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_id": {
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
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_snapshot_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
						"updated_time": {
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

type QueryField struct {
	Field     string `json:"field"`
	Value     string `json:"value"`
	Operation string `json:"operation"`
}

func dataSourceAlicloudHbrSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "SearchHistoricalSnapshots"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("limit"); ok {
		request["Limit"] = v
	}
	request["SourceType"] = d.Get("source_type")

	filtersMapList := make([]map[string]interface{}, 0)

	if v, ok := d.GetOk("vault_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"field":     "VaultId",
			"value":     v.(string),
			"operation": "MATCH_TERM",
		})
	}
	if v, ok := d.GetOk("instance_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"field":     "InstanceId",
			"value":     v.(string),
			"operation": "MATCH_TERM",
		})
	}
	if v, ok := d.GetOk("bucket"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"field":     "Bucket",
			"value":     v.(string),
			"operation": "MATCH_TERM",
		})
	}
	if v, ok := d.GetOk("file_system_id"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"field":     "FileSystemId",
			"value":     v.(string),
			"operation": "MATCH_TERM",
		})
	}
	if v, ok := d.GetOk("create_time"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"field":     "CreateTime",
			"value":     strconv.FormatInt(ConvertNasFileSystemStringToUnix(v.(string)), 10),
			"operation": "MATCH_TERM",
		})
	}
	if v, ok := d.GetOk("complete_time"); ok {
		completeChecker := make(map[string]interface{})
		completeChecker["field"] = "CompleteTime"
		completeChecker["operation"] = "MATCH_TERM"
		if vv, ok := d.GetOk("complete_time_checker"); ok {
			completeChecker["operation"] = vv.(string)
		}
		if completeChecker["operation"] == "BETWEEN" {
			timeSection := strings.Split(v.(string), ",")
			completeChecker["value"] = ConvertNasFileSystemStringToUnix(timeSection[0])
			completeChecker["value2"] = ConvertNasFileSystemStringToUnix(timeSection[1])
		} else {
			completeChecker["value"] = ConvertNasFileSystemStringToUnix(v.(string))
		}
		filtersMapList = append(filtersMapList, completeChecker)
	}
	if v, ok := d.GetOk("status"); ok {
		filtersMapList = append(filtersMapList, map[string]interface{}{
			"field":     "Status",
			"value":     v.(string),
			"operation": "MATCH_TERM",
		})
	}
	request["Query"], _ = convertListMapToJsonString(filtersMapList)

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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_snapshots", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Snapshots.Snapshot", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Snapshots.Snapshot", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SnapshotId"])]; !ok {
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
			"actual_bytes":         fmt.Sprint(object["ActualBytes"]),
			"actual_items":         fmt.Sprint(object["ActualItems"]),
			"backup_type":          object["BackupType"],
			"bucket":               object["Bucket"],
			"bytes_done":           fmt.Sprint(object["BytesDone"]),
			"bytes_total":          fmt.Sprint(object["BytesTotal"]),
			"complete_time":        fmt.Sprint(object["CompleteTime"]),
			"create_time":          fmt.Sprint(object["CreateTime"]),
			"created_time":         fmt.Sprint(object["CreatedTime"]),
			"error_file":           object["ErrorFile"],
			"file_system_id":       object["FileSystemId"],
			"instance_id":          object["InstanceId"],
			"items_done":           fmt.Sprint(object["ItemsDone"]),
			"items_total":          fmt.Sprint(object["ItemsTotal"]),
			"job_id":               object["JobId"],
			"parent_snapshot_hash": object["ParentSnapshotHash"],
			"prefix":               object["Prefix"],
			"client_id":            object["ClientId"],
			"retention":            fmt.Sprint(object["Retention"]),
			"snapshot_hash":        object["SnapshotHash"],
			"id":                   fmt.Sprint(object["SnapshotId"]),
			"snapshot_id":          fmt.Sprint(object["SnapshotId"]),
			"source_type":          object["SourceType"],
			"start_time":           fmt.Sprint(object["StartTime"]),
			"status":               object["Status"],
			"updated_time":         fmt.Sprint(object["UpdatedTime"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("snapshots", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
