package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDtsJobMonitorRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDtsJobMonitorRuleCreate,
		Read:   resourceAlicloudDtsJobMonitorRuleRead,
		Update: resourceAlicloudDtsJobMonitorRuleUpdate,
		Delete: resourceAlicloudDtsJobMonitorRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"delay_rule_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"dts_job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Y", "N"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"delay", "error"}, false),
			},
		},
	}
}

func resourceAlicloudDtsJobMonitorRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateJobMonitorRule"
	request := make(map[string]interface{})
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("delay_rule_time"); ok {
		request["DelayRuleTime"] = v
	}
	request["DtsJobId"] = d.Get("dts_job_id")
	if v, ok := d.GetOk("phone"); ok {
		request["Phone"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("state"); ok {
		request["State"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dts_job_monitor_rule", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["DtsJobId"]))

	return resourceAlicloudDtsJobMonitorRuleRead(d, meta)
}
func resourceAlicloudDtsJobMonitorRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := DtsService{client}
	object, err := dtsService.DescribeDtsJobMonitorRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_job_monitor_rule dtsService.DescribeDtsJobMonitorRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("dts_job_id", d.Id())
	if v := object["MonitorRules"].([]interface{}); len(v) > 0 {
		d.Set("delay_rule_time", v[0].(map[string]interface{})["DelayRuleTime"])
		d.Set("phone", v[0].(map[string]interface{})["Phone"])
		d.Set("state", v[0].(map[string]interface{})["State"])
		d.Set("type", v[0].(map[string]interface{})["Type"])
	}
	return nil
}
func resourceAlicloudDtsJobMonitorRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	if d.HasChange("delay_rule_time") {
		update = true
		if v, ok := d.GetOk("delay_rule_time"); ok {
			request["DelayRuleTime"] = v
		}
	}
	if d.HasChange("phone") {
		update = true
		if v, ok := d.GetOk("phone"); ok {
			request["Phone"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("state") {
		update = true
		if v, ok := d.GetOk("state"); ok {
			request["State"] = v
		}
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if update {
		action := "CreateJobMonitorRule"
		conn, err := client.NewDtsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudDtsJobMonitorRuleRead(d, meta)
}
func resourceAlicloudDtsJobMonitorRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudDtsJobMonitorRule. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
