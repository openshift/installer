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

func dataSourceAlicloudHbrServerBackupPlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudHbrServerBackupPlansRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"filters": {
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
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detail": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_consistent": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"snapshot_group": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"enable_fs_freeze": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"pre_script_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"post_script_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timeout_in_seconds": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"disk_id_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"do_copy": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"destination_region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_retention": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_server_backup_plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_server_backup_plan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudHbrServerBackupPlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBackupPlans"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		for filtersPtr, filters := range v.([]interface{}) {
			filtersArg := filters.(map[string]interface{})
			request[fmt.Sprintf("Filters.%d.Key", filtersPtr+1)] = filtersArg["key"]
			for valuesPtr, values := range filtersArg["values"].([]interface{}) {
				request[fmt.Sprintf("Filters.%d.Values.%d", filtersPtr+1, valuesPtr+1)] = values
			}
		}
	}
	request["SourceType"] = "UDM_ECS"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_hbr_server_backup_plans", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.BackupPlans.BackupPlan", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BackupPlans.BackupPlan", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PlanId"])]; !ok {
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
			"create_time":                 fmt.Sprint(object["CreatedTime"]),
			"disabled":                    object["Disabled"],
			"id":                          fmt.Sprint(object["PlanId"]),
			"ecs_server_backup_plan_id":   fmt.Sprint(object["PlanId"]),
			"ecs_server_backup_plan_name": object["PlanName"],
			"instance_id":                 object["InstanceId"],
			"retention":                   fmt.Sprint(object["Retention"]),
			"schedule":                    object["Schedule"],
		}
		if v, ok := object["Detail"].(map[string]interface{}); ok {
			details := make([]map[string]interface{}, 0)
			detail := make(map[string]interface{})
			detail["app_consistent"] = v["appConsistent"]
			detail["snapshot_group"] = v["snapshotGroup"]
			detail["enable_fs_freeze"] = v["enableFsFreeze"]
			detail["pre_script_path"] = v["preScriptPath"]
			detail["post_script_path"] = v["postScriptPath"]
			detail["timeout_in_seconds"] = formatInt(v["timeoutInSeconds"])
			detail["do_copy"] = v["doCopy"]
			detail["destination_region_id"] = v["destinationRegionId"]
			detail["destination_retention"] = v["destinationRetention"]
			if v["diskIdList"] != nil {
				detail["disk_id_list"] = expandStringList(v["diskIdList"].([]interface{}))
			}
			details = append(details, detail)
			mapping["detail"] = details
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("plans", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
