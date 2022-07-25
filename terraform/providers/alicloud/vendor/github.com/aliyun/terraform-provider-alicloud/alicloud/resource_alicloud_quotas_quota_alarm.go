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

func resourceAlicloudQuotasQuotaAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudQuotasQuotaAlarmCreate,
		Read:   resourceAlicloudQuotasQuotaAlarmRead,
		Update: resourceAlicloudQuotasQuotaAlarmUpdate,
		Delete: resourceAlicloudQuotasQuotaAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_alarm_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"quota_dimensions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				ForceNew: true,
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"threshold_percent": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"web_hook": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudQuotasQuotaAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateQuotaAlarm"
	request := make(map[string]interface{})
	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}
	request["SourceIp"] = client.SourceIp
	request["ProductCode"] = d.Get("product_code")
	request["QuotaActionCode"] = d.Get("quota_action_code")
	request["AlarmName"] = d.Get("quota_alarm_name")
	if v, ok := d.GetOk("quota_dimensions"); ok {
		quotaDimensionsMaps := make([]map[string]interface{}, 0)
		for _, quotaDimensions := range v.(*schema.Set).List() {
			quotaDimensionsMap := make(map[string]interface{})
			quotaDimensionsArg := quotaDimensions.(map[string]interface{})
			quotaDimensionsMap["Key"] = quotaDimensionsArg["key"]
			quotaDimensionsMap["Value"] = quotaDimensionsArg["value"]
			quotaDimensionsMaps = append(quotaDimensionsMaps, quotaDimensionsMap)
		}
		request["QuotaDimensions"] = quotaDimensionsMaps

	}

	if v, ok := d.GetOk("threshold"); ok {
		request["Threshold"] = v
	}

	if v, ok := d.GetOk("threshold_percent"); ok {
		request["ThresholdPercent"] = v
	}

	if v, ok := d.GetOk("web_hook"); ok {
		request["WebHook"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_quotas_quota_alarm", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AlarmId"]))

	return resourceAlicloudQuotasQuotaAlarmRead(d, meta)
}
func resourceAlicloudQuotasQuotaAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasService := QuotasService{client}
	object, err := quotasService.DescribeQuotasQuotaAlarm(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_quota_alarm quotasService.DescribeQuotasQuotaAlarm Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("product_code", object["ProductCode"])
	d.Set("quota_action_code", object["QuotaActionCode"])
	d.Set("quota_alarm_name", object["AlarmName"])

	quotaDimensionList := make([]map[string]interface{}, 0)
	if quotaDimension, ok := object["QuotaDimension"]; ok {
		for k, v := range quotaDimension.(map[string]interface{}) {
			quotaDimensionMap := make(map[string]interface{})
			quotaDimensionMap["key"] = k
			quotaDimensionMap["value"] = v
			quotaDimensionList = append(quotaDimensionList, quotaDimensionMap)
		}
	}

	if err := d.Set("quota_dimensions", quotaDimensionList); err != nil {
		return WrapError(err)
	}
	d.Set("threshold", object["Threshold"])
	d.Set("threshold_percent", object["ThresholdPercent"])
	return nil
}
func resourceAlicloudQuotasQuotaAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"AlarmId":  d.Id(),
		"SourceIp": client.SourceIp,
	}
	if d.HasChange("quota_alarm_name") {
		update = true
	}
	request["AlarmName"] = d.Get("quota_alarm_name")
	if d.HasChange("threshold") {
		update = true
		request["Threshold"] = d.Get("threshold")
	}
	if d.HasChange("threshold_percent") {
		update = true
		request["ThresholdPercent"] = d.Get("threshold_percent")
	}
	if d.HasChange("web_hook") {
		update = true
		request["WebHook"] = d.Get("web_hook")
	}
	if update {
		action := "UpdateQuotaAlarm"
		conn, err := client.NewQuotasClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
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
	return resourceAlicloudQuotasQuotaAlarmRead(d, meta)
}
func resourceAlicloudQuotasQuotaAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteQuotaAlarm"
	var response map[string]interface{}
	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AlarmId":  d.Id(),
		"SourceIp": client.SourceIp,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
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
	return nil
}
