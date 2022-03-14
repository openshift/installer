package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudHbrNasBackupPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbrNasBackupPlanCreate,
		Read:   resourceAlicloudHbrNasBackupPlanRead,
		Update: resourceAlicloudHbrNasBackupPlanUpdate,
		Delete: resourceAlicloudHbrNasBackupPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"backup_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPLETE"}, false),
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"create_time": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'create_time' has been deprecated from provider version 1.153.0.",
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nas_backup_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"options": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudHbrNasBackupPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBackupPlan"
	request := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request["BackupType"] = d.Get("backup_type")
	if v, ok := d.GetOk("create_time"); ok {
		request["CreateTime"] = ConvertNasFileSystemStringToUnix(v.(string))
	}
	if v, ok := d.GetOk("file_system_id"); ok {
		request["FileSystemId"] = v
	}
	request["PlanName"] = d.Get("nas_backup_plan_name")
	if v, ok := d.GetOk("options"); ok {
		request["Options"] = v
	}
	if v, ok := d.GetOk("retention"); ok {
		request["Retention"] = v
	}
	request["Schedule"] = d.Get("schedule")
	request["SourceType"] = "NAS"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_nas_backup_plan", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["PlanId"]))

	return resourceAlicloudHbrNasBackupPlanUpdate(d, meta)
}
func resourceAlicloudHbrNasBackupPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrService := HbrService{client}
	object, err := hbrService.DescribeHbrNasBackupPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_nas_backup_plan hbrService.DescribeHbrNasBackupPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("backup_type", object["BackupType"])
	t := int64(formatInt(object["CreateTime"]))
	d.Set("create_time", ConvertNasFileSystemUnixToString(d.Get("create_time").(string), t))
	d.Set("file_system_id", object["FileSystemId"])
	d.Set("nas_backup_plan_name", object["PlanName"])
	d.Set("options", object["Options"])
	if object["Paths"] != nil {
		d.Set("path", object["Paths"].(map[string]interface{})["Path"])
	}
	d.Set("retention", fmt.Sprint(formatInt(object["Retention"])))
	d.Set("schedule", object["Schedule"])
	d.Set("vault_id", object["VaultId"])
	d.Set("disabled", object["Disabled"])
	return nil
}
func resourceAlicloudHbrNasBackupPlanUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if v, ok := d.GetOk("vault_id"); ok {
		request["VaultId"] = v
	}
	if !d.IsNewResource() && d.HasChange("nas_backup_plan_name") {
		update = true
		request["PlanName"] = d.Get("nas_backup_plan_name")
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
	if !d.IsNewResource() && d.HasChange("schedule") {
		update = true
		request["Schedule"] = d.Get("schedule")
	}
	request["SourceType"] = "NAS"
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
		d.SetPartial("vault_id")
		d.SetPartial("nas_backup_plan_name")
		d.SetPartial("options")
		d.SetPartial("path")
		d.SetPartial("schedule")
	}
	if d.HasChange("disabled") {
		object, err := hbrService.DescribeHbrNasBackupPlan(d.Id())
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
			request["SourceType"] = "NAS"
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
	return resourceAlicloudHbrNasBackupPlanRead(d, meta)
}
func resourceAlicloudHbrNasBackupPlanDelete(d *schema.ResourceData, meta interface{}) error {
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

	request["SourceType"] = "NAS"
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

	{
		action := "DetachNasFileSystem"
		request := map[string]interface{}{
			"PlanId": d.Id(),
		}
		request["VaultId"] = d.Get("vault_id")
		if v, ok := d.GetOk("file_system_id"); ok {
			request["FileSystemId"] = v
		}
		if v, ok := d.GetOk("create_time"); ok {
			request["CreateTime"] = ConvertNasFileSystemStringToUnix(v.(string))
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
			if IsExpectedErrors(err, []string{"NasFileSystemNotExist"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return nil
}

func ConvertNasFileSystemUnixToString(v string, t int64) string {
	var s string
	c := time.Unix(t, 0)
	if strings.HasSuffix(v, "CST") {
		s = fmt.Sprint(c.Format("2006-01-02T15:04:05CST"))
	} else if strings.HasSuffix(v, "Z") {
		s = fmt.Sprint(c.Format("2006-01-02T15:04:05Z"))
	} else {
		s = fmt.Sprint(c.Format("2006-01-02T15:04:05CST"))
	}
	return s
}

func ConvertNasFileSystemStringToUnix(v string) int64 {
	var t time.Time
	if strings.HasSuffix(v, "CST") {
		t, _ = time.ParseInLocation("2006-01-02T15:04:05CST", v, time.Local)
	} else if strings.HasSuffix(v, "Z") {
		t, _ = time.ParseInLocation("2006-01-02T15:04:05Z", v, time.FixedZone("CST", 8*3600))
	} else {
		t, _ = time.ParseInLocation("2006-01-02T15:04:05CST", v, time.Local)
	}
	return t.Unix()
}
