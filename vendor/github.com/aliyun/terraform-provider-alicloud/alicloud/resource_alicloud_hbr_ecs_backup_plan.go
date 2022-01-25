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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudHbrEcsBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrEcsBackupPlanCreate,
		Read:   resourceAlicloudHbrEcsBackupPlanRead,
		Update: resourceAlicloudHbrEcsBackupPlanUpdate,
		Delete: resourceAlicloudHbrEcsBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ecs_backup_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE"}, false),
			},
			"retention": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"speed_limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_paths": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Attribute update_paths has been deprecated in v1.139.0+ and you do not need to set it anymore.",
			},
		},
	}
}

func resourceAlicloudHbrEcsBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request["BackupType"] = d.Get("backup_type")
	if v, ok := d.GetOk("bucket"); ok {
		request["Bucket"] = v
	}
	if v, ok := d.GetOk("detail"); ok {
		request["Detail"] = v
	}
	request["PlanName"] = d.Get("ecs_backup_plan_name")
	if v, ok := d.GetOk("exclude"); ok {
		request["Exclude"] = v
	}
	if v, ok := d.GetOk("include"); ok {
		request["Include"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}
	request["Schedule"] = d.Get("schedule")
	request["SourceType"] = "ECS_FILE"
	if v, ok := d.GetOk("speed_limit"); ok {
		request["SpeedLimit"] = v
	}
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_ecs_backup_plan", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["PlanId"]))

	return resourceAlicloudHbrEcsBackupPlanUpdate(d, meta)
}
func resourceAlicloudHbrEcsBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrEcsBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_ecs_backup_plan hbrService.DescribeHbrEcsBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("backup_type", object["BackupType"])
	d.Set("detail", object["Detail"])
	d.Set("ecs_backup_plan_name", object["PlanName"])
	d.Set("exclude", object["Exclude"])
	d.Set("include", object["Include"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("options", object["Options"])
	if object["Paths"] != nil {
		d.Set("path", object["Paths"].(map[string]interface{})["Path"])
	}
	d.Set("retention", fmt.Sprint(formatInt(object["Retention"])))
	d.Set("schedule", object["Schedule"])
	d.Set("speed_limit", object["SpeedLimit"])
	d.Set("vault_id", object["VaultId"])
	d.Set("disabled", object["Disabled"])
	return nil
}
func resourceAlicloudHbrEcsBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"PlanId": d.Id(),
	}

	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	if !d.IsNewResource() && d.HasChange("detail") {
		update = true
		if v, ok := d.GetOk("detail"); ok {
			request["Detail"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("ecs_backup_plan_name") {
		update = true
		request["PlanName"] = d.Get("ecs_backup_plan_name")
	}
	if !d.IsNewResource() && d.HasChange("exclude") {
		update = true
		if v, ok := d.GetOk("exclude"); ok {
			request["Exclude"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("include") {
		update = true
		if v, ok := d.GetOk("include"); ok {
			request["Include"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("options") {
		update = true
		if v, ok := d.GetOk("options"); ok {
			request["Options"] = v
		}
	}
	if d.HasChange("path") {
		update = true
		request["UpdatePaths"] = true
		if v, ok := d.GetOk("path"); ok {
			request["Path"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("retention") {
		update = true
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}
	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
		request["Schedule"] = d.Get("schedule")
	}
	request["SourceType"] = "ECS_FILE"
	if !d.IsNewResource() && d.HasChange("speed_limit") {
		update = true
		if v, ok := d.GetOk("speed_limit"); ok {
			request["SpeedLimit"] = v
		}
	}
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
		d.SetPartial("vault_id")
		d.SetPartial("detail")
		d.SetPartial("ecs_backup_plan_name")
		d.SetPartial("exclude")
		d.SetPartial("include")
		d.SetPartial("options")
		d.SetPartial("path")
		d.SetPartial("retention")
		d.SetPartial("schedule")
		d.SetPartial("speed_limit")
	}
	if d.HasChange("disabled") {
		object, err := hbrService.DescribeHbrEcsBackupPlan(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strconv.FormatBool(d.Get("disabled").(bool))
		if strconv.FormatBool(object["Disabled"].(bool)) != target {
			action := "EnableBackupPlan"
			if target == "false" {
				action = "EnableBackupPlan"
			} else {
				action = "DisableBackupPlan"
			}
			request := map[string]interface{}{
				"PlanId": d.Id(),
			}
			request["VaultId"] = d.Get("vault_id")
			request["SourceType"] = "ECS_FILE"
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
	}
	d.SetPartial("disabled")

	d.Partial(false)
	return resourceAlicloudHbrEcsBackupPlanRead(d, meta)
}
func resourceAlicloudHbrEcsBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
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

	request["SourceType"] = "ECS_FILE"
	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
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
