package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCmsMetricRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsMetricRuleTemplateCreate,
		Read:   resourceAlicloudCmsMetricRuleTemplateRead,
		Update: resourceAlicloudCmsMetricRuleTemplateUpdate,
		Delete: resourceAlicloudCmsMetricRuleTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_templates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ecs", "rds", "ads", "slb", "vpc", "apigateway", "cdn", "cs", "dcdn", "ddos", "eip", "elasticsearch", "emr", "ess", "hbase", "iot_edge", "kvstore_sharding", "kvstore_splitrw", "kvstore_standard", "memcache", "mns", "mongodb", "mongodb_cluster", "mongodb_sharding", "mq_topic", "ocs", "opensearch", "oss", "polardb", "petadata", "scdn", "sharebandwidthpackages", "sls", "vpn"}, false),
						},
						"escalations": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"critical": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"info": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"warn": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"webhook": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"apply_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_rule_template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notify_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rest_version": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCmsMetricRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMetricRuleTemplate"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("alert_templates"); ok {
		alertTemplatesMaps := make([]map[string]interface{}, 0)
		for _, alertTemplates := range v.(*schema.Set).List() {
			alertTemplatesArg := alertTemplates.(map[string]interface{})
			alertTemplatesMap := map[string]interface{}{}
			alertTemplatesMap["Category"] = alertTemplatesArg["category"]
			if escalationsMaps, ok := alertTemplatesArg["escalations"]; ok {
				escalationsMap := map[string]interface{}{}
				for _, escalationsArg := range escalationsMaps.(*schema.Set).List() {
					if criticalMaps, ok := escalationsArg.(map[string]interface{})["critical"]; ok {
						requestCriticalArg := map[string]interface{}{}
						for _, criticalMap := range criticalMaps.(*schema.Set).List() {
							criticalArg := criticalMap.(map[string]interface{})
							requestCriticalArg["ComparisonOperator"] = criticalArg["comparison_operator"]
							requestCriticalArg["Statistics"] = criticalArg["statistics"]
							requestCriticalArg["Threshold"] = criticalArg["threshold"]
							requestCriticalArg["Times"] = criticalArg["times"]
						}
						escalationsMap["Critical"] = requestCriticalArg
					}
					if infoMaps, ok := escalationsArg.(map[string]interface{})["info"]; ok {
						requestInfoArg := map[string]interface{}{}
						for _, infoMap := range infoMaps.(*schema.Set).List() {
							infoArg := infoMap.(map[string]interface{})
							requestInfoArg["ComparisonOperator"] = infoArg["comparison_operator"]
							requestInfoArg["Statistics"] = infoArg["statistics"]
							requestInfoArg["Threshold"] = infoArg["threshold"]
							requestInfoArg["Times"] = infoArg["times"]
						}
						escalationsMap["Info"] = requestInfoArg
					}
					if warnMaps, ok := escalationsArg.(map[string]interface{})["warn"]; ok {
						requestWarnArg := map[string]interface{}{}
						for _, warnMap := range warnMaps.(*schema.Set).List() {
							warnArg := warnMap.(map[string]interface{})
							requestWarnArg["ComparisonOperator"] = warnArg["comparison_operator"]
							requestWarnArg["Statistics"] = warnArg["statistics"]
							requestWarnArg["Threshold"] = warnArg["threshold"]
							requestWarnArg["Times"] = warnArg["times"]
						}
						escalationsMap["Warn"] = requestWarnArg
					}
				}
				alertTemplatesMap["Escalations"] = escalationsMap
			}
			alertTemplatesMap["MetricName"] = alertTemplatesArg["metric_name"]
			alertTemplatesMap["Namespace"] = alertTemplatesArg["namespace"]
			alertTemplatesMap["RuleName"] = alertTemplatesArg["rule_name"]
			alertTemplatesMap["Webhook"] = alertTemplatesArg["webhook"]
			alertTemplatesMaps = append(alertTemplatesMaps, alertTemplatesMap)
		}

		request["AlertTemplates"] = alertTemplatesMaps
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Name"] = d.Get("metric_rule_template_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_metric_rule_template", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAlicloudCmsMetricRuleTemplateUpdate(d, meta)
}
func resourceAlicloudCmsMetricRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsMetricRuleTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_metric_rule_template cmsService.DescribeCmsMetricRuleTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if alertTemplatesMap, ok := object["AlertTemplates"].(map[string]interface{}); ok && alertTemplatesMap != nil {
		if alertTemplateList, ok := alertTemplatesMap["AlertTemplate"]; ok && alertTemplateList != nil {
			alertTemplatesMaps := make([]map[string]interface{}, 0)
			for _, alertTemplateListItem := range alertTemplateList.([]interface{}) {
				if alertTemplateListItemMap, ok := alertTemplateListItem.(map[string]interface{}); ok {
					alertTempArg := make(map[string]interface{}, 0)
					alertTempArg["category"] = alertTemplateListItemMap["Category"]
					alertTempArg["metric_name"] = alertTemplateListItemMap["MetricName"]
					alertTempArg["namespace"] = alertTemplateListItemMap["Namespace"]
					alertTempArg["rule_name"] = alertTemplateListItemMap["RuleName"]
					alertTempArg["webhook"] = alertTemplateListItemMap["Webhook"]
					escalationsMaps := make([]map[string]interface{}, 0)
					escalationsMap := map[string]interface{}{}
					if EscalationsMap, ok := alertTemplateListItemMap["Escalations"].(map[string]interface{}); ok && len(EscalationsMap) > 0 {
						EscalationsArg := EscalationsMap

						if criticalMap, ok := EscalationsArg["Critical"].(map[string]interface{}); ok && len(criticalMap) > 0 {
							criticalMaps := make([]map[string]interface{}, 0)
							criticalArg := map[string]interface{}{}
							criticalArg["comparison_operator"] = criticalMap["ComparisonOperator"]
							criticalArg["statistics"] = criticalMap["Statistics"]
							criticalArg["threshold"] = criticalMap["Threshold"]
							criticalArg["times"] = criticalMap["Times"]
							criticalMaps = append(criticalMaps, criticalArg)
							escalationsMap["critical"] = criticalMaps
						}

						if infoMap, ok := EscalationsArg["Info"].(map[string]interface{}); ok && len(infoMap) > 0 {
							infoMaps := make([]map[string]interface{}, 0)
							infoArg := map[string]interface{}{}
							infoArg["comparison_operator"] = infoMap["ComparisonOperator"]
							infoArg["statistics"] = infoMap["Statistics"]
							infoArg["threshold"] = infoMap["Threshold"]
							infoArg["times"] = infoMap["Times"]
							infoMaps = append(infoMaps, infoArg)
							escalationsMap["info"] = infoMaps
						}

						if warnMap, ok := EscalationsArg["Warn"].(map[string]interface{}); ok && len(warnMap) > 0 {
							warnMaps := make([]map[string]interface{}, 0)
							warnArg := make(map[string]interface{}, 0)
							warnArg["comparison_operator"] = warnMap["ComparisonOperator"]
							warnArg["statistics"] = warnMap["Statistics"]
							warnArg["threshold"] = warnMap["Threshold"]
							warnArg["times"] = warnMap["Times"]
							warnMaps = append(warnMaps, warnArg)
							escalationsMap["warn"] = warnMaps
						}
					}
					escalationsMaps = append(escalationsMaps, escalationsMap)

					alertTempArg["escalations"] = escalationsMaps
					alertTemplatesMaps = append(alertTemplatesMaps, alertTempArg)
				}
			}
			d.Set("alert_templates", alertTemplatesMaps)
		}
	}
	d.Set("description", object["Description"])
	d.Set("metric_rule_template_name", object["Name"])
	d.Set("rest_version", object["RestVersion"])
	return nil
}
func resourceAlicloudCmsMetricRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"TemplateIds": d.Id(),
	}
	if d.HasChange("group_id") {
		update = true
	}
	if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	}
	if update {
		if v, ok := d.GetOk("apply_mode"); ok {
			request["ApplyMode"] = v
		}
		if v, ok := d.GetOk("enable_end_time"); ok {
			request["EnableEndTime"] = v
		}
		if v, ok := d.GetOk("enable_start_time"); ok {
			request["EnableStartTime"] = v
		}
		if v, ok := d.GetOk("notify_level"); ok {
			request["NotifyLevel"] = v
		}
		if v, ok := d.GetOk("silence_time"); ok {
			request["SilenceTime"] = v
		}
		if v, ok := d.GetOk("webhook"); ok {
			request["Webhook"] = v
		}
		action := "ApplyMetricRuleTemplate"
		conn, err := client.NewCmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("group_id")
	}
	update = false
	modifyMetricRuleTemplateReq := map[string]interface{}{
		"TemplateId": d.Id(),
	}

	if v, ok := d.GetOk("rest_version"); ok {
		modifyMetricRuleTemplateReq["RestVersion"] = v
	}
	if !d.IsNewResource() && d.HasChange("alert_templates") {
		update = true
		if v, ok := d.GetOk("alert_templates"); ok {
			alertTemplatesMaps := make([]map[string]interface{}, 0)
			for _, alertTemplates := range v.(*schema.Set).List() {
				alertTemplatesArg := alertTemplates.(map[string]interface{})
				alertTemplatesMap := map[string]interface{}{}
				alertTemplatesMap["Category"] = alertTemplatesArg["category"]
				if escalationsMaps, ok := alertTemplatesArg["escalations"]; ok {
					escalationsMap := map[string]interface{}{}
					for _, escalationsArg := range escalationsMaps.(*schema.Set).List() {
						if criticalMaps, ok := escalationsArg.(map[string]interface{})["critical"]; ok {
							requestCriticalArg := map[string]interface{}{}
							for _, criticalMap := range criticalMaps.(*schema.Set).List() {
								criticalArg := criticalMap.(map[string]interface{})
								requestCriticalArg["ComparisonOperator"] = criticalArg["comparison_operator"]
								requestCriticalArg["Statistics"] = criticalArg["statistics"]
								requestCriticalArg["Threshold"] = criticalArg["threshold"]
								requestCriticalArg["Times"] = criticalArg["times"]
							}
							escalationsMap["Critical"] = requestCriticalArg
						}
						if infoMaps, ok := escalationsArg.(map[string]interface{})["info"]; ok {
							requestInfoArg := map[string]interface{}{}
							for _, infoMap := range infoMaps.(*schema.Set).List() {
								infoArg := infoMap.(map[string]interface{})
								requestInfoArg["ComparisonOperator"] = infoArg["comparison_operator"]
								requestInfoArg["Statistics"] = infoArg["statistics"]
								requestInfoArg["Threshold"] = infoArg["threshold"]
								requestInfoArg["Times"] = infoArg["times"]
							}
							escalationsMap["Info"] = requestInfoArg
						}
						if warnMaps, ok := escalationsArg.(map[string]interface{})["warn"]; ok {
							requestWarnArg := map[string]interface{}{}
							for _, warnMap := range warnMaps.(*schema.Set).List() {
								warnArg := warnMap.(map[string]interface{})
								requestWarnArg["ComparisonOperator"] = warnArg["comparison_operator"]
								requestWarnArg["Statistics"] = warnArg["statistics"]
								requestWarnArg["Threshold"] = warnArg["threshold"]
								requestWarnArg["Times"] = warnArg["times"]
							}
							escalationsMap["Warn"] = requestWarnArg
						}
					}
					alertTemplatesMap["Escalations"] = escalationsMap
				}
				alertTemplatesMap["MetricName"] = alertTemplatesArg["metric_name"]
				alertTemplatesMap["Namespace"] = alertTemplatesArg["namespace"]
				alertTemplatesMap["RuleName"] = alertTemplatesArg["rule_name"]
				alertTemplatesMap["Webhook"] = alertTemplatesArg["webhook"]
				alertTemplatesMaps = append(alertTemplatesMaps, alertTemplatesMap)
			}
			modifyMetricRuleTemplateReq["AlertTemplates"] = alertTemplatesMaps
		}
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			modifyMetricRuleTemplateReq["Description"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("metric_rule_template_name") {
		update = true
		modifyMetricRuleTemplateReq["Name"] = d.Get("metric_rule_template_name")
	}
	if update {
		action := "ModifyMetricRuleTemplate"
		conn, err := client.NewCmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, modifyMetricRuleTemplateReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyMetricRuleTemplateReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("rest_version")
		d.SetPartial("alert_templates")
		d.SetPartial("description")
		d.SetPartial("metric_rule_template_name")
	}
	d.Partial(false)
	return resourceAlicloudCmsMetricRuleTemplateRead(d, meta)
}
func resourceAlicloudCmsMetricRuleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMetricRuleTemplate"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
