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

func resourceAlicloudArmsPrometheusAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudArmsPrometheusAlertRuleCreate,
		Read:   resourceAlicloudArmsPrometheusAlertRuleRead,
		Update: resourceAlicloudArmsPrometheusAlertRuleUpdate,
		Delete: resourceAlicloudArmsPrometheusAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dispatch_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("notify_type"); ok && v.(string) == "DISPATCH_RULE" {
						return false
					}
					return true
				},
			},
			"duration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_MANAGER", "DISPATCH_RULE"}, false),
			},
			"prometheus_alert_rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"prometheus_alert_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudArmsPrometheusAlertRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePrometheusAlertRule"
	request := make(map[string]interface{})
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("annotations"); ok {
		annotationsMaps := make([]map[string]interface{}, 0)
		for _, annotations := range v.(*schema.Set).List() {
			annotationsMap := annotations.(map[string]interface{})
			annotationsMaps = append(annotationsMaps, annotationsMap)
		}
		if v, err := convertArrayObjectToJsonString(annotationsMaps); err == nil {
			request["Annotations"] = v
		} else {
			return WrapError(err)
		}
	}
	request["ClusterId"] = d.Get("cluster_id")
	if v, ok := d.GetOk("dispatch_rule_id"); ok {
		request["DispatchRuleId"] = v
	}
	request["Duration"] = d.Get("duration")
	request["Expression"] = d.Get("expression")
	if v, ok := d.GetOk("labels"); ok {
		labelsMaps := make([]map[string]interface{}, 0)
		for _, labels := range v.(*schema.Set).List() {
			labelsMap := labels.(map[string]interface{})
			labelsMaps = append(labelsMaps, labelsMap)
		}
		if v, err := convertArrayObjectToJsonString(labelsMaps); err == nil {
			request["Labels"] = v
		} else {
			return WrapError(err)
		}
	}
	request["Message"] = d.Get("message")
	if v, ok := d.GetOk("notify_type"); ok {
		request["NotifyType"] = v
	}
	request["AlertName"] = d.Get("prometheus_alert_rule_name")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_arms_prometheus_alert_rule", action, AlibabaCloudSdkGoERROR)
	}
	responsePrometheusAlertRule := response["PrometheusAlertRule"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["ClusterId"], ":", responsePrometheusAlertRule["AlertId"]))

	return resourceAlicloudArmsPrometheusAlertRuleRead(d, meta)
}
func resourceAlicloudArmsPrometheusAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsPrometheusAlertRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_prometheus_alert_rule armsService.DescribeArmsPrometheusAlertRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("cluster_id", parts[0])
	d.Set("prometheus_alert_rule_id", parts[1])
	if v, ok := object["Annotations"].([]interface{}); ok {
		annotations := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			if item["Name"] == "message" {
				continue
			}

			temp := map[string]interface{}{
				"name":  item["Name"],
				"value": item["Value"],
			}

			annotations = append(annotations, temp)
		}
		if err := d.Set("annotations", annotations); err != nil {
			return WrapError(err)
		}
	}
	d.Set("dispatch_rule_id", fmt.Sprint(formatInt(object["DispatchRuleId"])))
	d.Set("duration", object["Duration"])
	d.Set("expression", object["Expression"])
	if v, ok := object["Labels"].([]interface{}); ok {
		labels := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"name":  item["Name"],
				"value": item["Value"],
			}

			labels = append(labels, temp)
		}
		if err := d.Set("labels", labels); err != nil {
			return WrapError(err)
		}
	}
	d.Set("message", object["Message"])
	d.Set("notify_type", object["NotifyType"])
	d.Set("prometheus_alert_rule_name", object["AlertName"])
	d.Set("status", fmt.Sprint(formatInt(object["Status"])))
	d.Set("type", object["Type"])
	return nil
}
func resourceAlicloudArmsPrometheusAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ClusterId": parts[0],
		"AlertId":   parts[1],
	}

	request["Duration"] = d.Get("duration")
	if d.HasChange("duration") {
		update = true
	}
	request["Expression"] = d.Get("expression")
	if d.HasChange("expression") {
		update = true
	}
	request["Message"] = d.Get("message")
	if d.HasChange("message") {
		update = true
	}
	request["AlertName"] = d.Get("prometheus_alert_rule_name")
	if d.HasChange("prometheus_alert_rule_name") {
		update = true
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("annotations") {
		update = true
		if v, ok := d.GetOk("annotations"); ok {
			annotationsMaps := make([]map[string]interface{}, 0)
			for _, annotations := range v.(*schema.Set).List() {
				annotationsMap := annotations.(map[string]interface{})
				annotationsMaps = append(annotationsMaps, annotationsMap)
			}
			if v, err := convertArrayObjectToJsonString(annotationsMaps); err == nil {
				request["Annotations"] = v
			} else {
				return WrapError(err)
			}
		}
	}
	if d.HasChange("dispatch_rule_id") {
		update = true
		if v, ok := d.GetOk("dispatch_rule_id"); ok {
			request["DispatchRuleId"] = v
		}
	}
	if d.HasChange("labels") {
		update = true
		if v, ok := d.GetOk("labels"); ok {
			labelsMaps := make([]map[string]interface{}, 0)
			for _, labels := range v.(*schema.Set).List() {
				labelsMap := labels.(map[string]interface{})
				labelsMaps = append(labelsMaps, labelsMap)
			}
			if v, err := convertArrayObjectToJsonString(labelsMaps); err == nil {
				request["Labels"] = v
			} else {
				return WrapError(err)
			}
		}
	}
	if d.HasChange("notify_type") {
		update = true
		if v, ok := d.GetOk("notify_type"); ok {
			request["NotifyType"] = v
		}
	}

	if update {
		action := "UpdatePrometheusAlertRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudArmsPrometheusAlertRuleRead(d, meta)
}
func resourceAlicloudArmsPrometheusAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeletePrometheusAlertRule"
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AlertId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
