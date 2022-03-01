package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudHbrServerBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrServerBackupPlanCreate,
		Read:   resourceAlicloudHbrServerBackupPlanRead,
		Update: resourceAlicloudHbrServerBackupPlanUpdate,
		Delete: resourceAlicloudHbrServerBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"detail": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_consistent": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"snapshot_group": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"enable_fs_freeze": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"pre_script_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"post_script_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout_in_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
						"disk_id_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"do_copy": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"destination_region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_retention": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
					},
				},
			},
			"ecs_server_backup_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudHbrServerBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request["BackupType"] = "COMPLETE"
	if v, err := hbrService.convertDetailToString(d.Get("detail").(*schema.Set).List()); err == nil {
		request["Detail"] = v
	} else {
		return WrapError(err)
	}

	request["PlanName"] = d.Get("ecs_server_backup_plan_name")
	request["InstanceId"] = d.Get("instance_id")
	request["UdmRegionId"] = client.RegionId
	request["Retention"] = d.Get("retention")

	request["Schedule"] = d.Get("schedule")
	request["SourceType"] = "UDM_ECS"

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_server_backup_plan", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["PlanId"]))

	return resourceAlicloudHbrServerBackupPlanUpdate(d, meta)
}
func resourceAlicloudHbrServerBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrServerBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_server_backup_plan hbrService.DescribeHbrServerBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
		if err := d.Set("detail", details); err != nil {
			return WrapError(err)
		}
	}
	d.Set("ecs_server_backup_plan_name", object["PlanName"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("retention", formatInt(object["Retention"]))
	d.Set("schedule", object["Schedule"])
	d.Set("disabled", object["Disabled"])
	return nil
}
func resourceAlicloudHbrServerBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"PlanId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("retention") {
		update = true
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}
	if !d.IsNewResource() && d.HasChange("detail") {
		update = true
		if v, err := hbrService.convertDetailToString(d.Get("detail").(*schema.Set).List()); err == nil {
			request["Detail"] = v
		} else {
			return WrapError(err)
		}
	}
	if !d.IsNewResource() && d.HasChange("ecs_server_backup_plan_name") {
		update = true
		request["PlanName"] = d.Get("ecs_server_backup_plan_name")
	}
	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
		request["Schedule"] = d.Get("schedule")
	}
	request["SourceType"] = "UDM_ECS"
	if update {
		action := "UpdateBackupPlan"
		conn, err := client.NewHbrClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("retention")
		d.SetPartial("detail")
		d.SetPartial("ecs_server_backup_plan_name")
		d.SetPartial("schedule")
	}
	if d.HasChange("disabled") {
		object, err := hbrService.DescribeHbrServerBackupPlan(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strconv.FormatBool(d.Get("disabled").(bool))
		if strconv.FormatBool(object["Disabled"].(bool)) != target {
			if target == "false" {
				request := map[string]interface{}{
					"PlanId": d.Id(),
				}
				request["SourceType"] = "UDM_ECS"
			}
			action := "EnableBackupPlan"
			if target == "false" {
				action = "EnableBackupPlan"
			} else {
				action = "DisableBackupPlan"
			}
			conn, err := client.NewHbrClient()
			if err != nil {
				return WrapError(err)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			if fmt.Sprint(response["Success"]) == "false" {
				return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
			}
		}
		d.SetPartial("disabled")
	}

	d.Partial(false)
	return resourceAlicloudHbrServerBackupPlanRead(d, meta)
}
func resourceAlicloudHbrServerBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBackupPlan"
	var response map[string]interface{}
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PlanId": d.Id(),
	}

	request["SourceType"] = "UDM_ECS"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
