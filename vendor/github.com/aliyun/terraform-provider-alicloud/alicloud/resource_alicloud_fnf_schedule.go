package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudFnfSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFnfScheduleCreate,
		Read:   resourceAlicloudFnfScheduleRead,
		Update: resourceAlicloudFnfScheduleUpdate,
		Delete: resourceAlicloudFnfScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cron_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"flow_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payload": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudFnfScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSchedule"
	request := make(map[string]interface{})
	conn, err := client.NewFnfClient()
	if err != nil {
		return WrapError(err)
	}
	request["CronExpression"] = d.Get("cron_expression")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("enable"); ok {
		request["Enable"] = v
	}

	request["FlowName"] = d.Get("flow_name")
	if v, ok := d.GetOk("payload"); ok {
		request["Payload"] = v
	}

	request["ScheduleName"] = d.Get("schedule_name")
	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fnf_schedule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FlowName"], ":", response["ScheduleName"]))

	return resourceAlicloudFnfScheduleRead(d, meta)
}
func resourceAlicloudFnfScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fnfService := FnfService{client}
	object, err := fnfService.DescribeFnfSchedule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fnf_schedule fnfService.DescribeFnfSchedule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("flow_name", parts[0])
	d.Set("schedule_name", parts[1])
	d.Set("cron_expression", object["CronExpression"])
	d.Set("description", object["Description"])
	d.Set("enable", object["Enable"])
	d.Set("last_modified_time", object["LastModifiedTime"])
	d.Set("payload", object["Payload"])
	d.Set("schedule_id", object["ScheduleId"])
	return nil
}
func resourceAlicloudFnfScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"FlowName":     parts[0],
		"ScheduleName": parts[1],
	}
	if d.HasChange("cron_expression") {
		update = true
		request["CronExpression"] = d.Get("cron_expression")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("enable") {
		update = true
		request["Enable"] = d.Get("enable")
	}
	if d.HasChange("payload") {
		update = true
		request["Payload"] = d.Get("payload")
	}
	if update {
		action := "UpdateSchedule"
		conn, err := client.NewFnfClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 1*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudFnfScheduleRead(d, meta)
}
func resourceAlicloudFnfScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteSchedule"
	var response map[string]interface{}
	conn, err := client.NewFnfClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FlowName":     parts[0],
		"ScheduleName": parts[1],
	}

	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-03-15"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ExecutionNotExists", "FlowNotExists", "ScheduleNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
