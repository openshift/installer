package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsGroupMetricRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsGroupMetricRuleCreate,
		Read:   resourceAlicloudCmsGroupMetricRuleRead,
		Update: resourceAlicloudCmsGroupMetricRuleUpdate,
		Delete: resourceAlicloudCmsGroupMetricRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contact_groups": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dimensions": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_subject": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"escalations": {
				Type:     schema.TypeSet,
				Required: true,
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
										Type:     schema.TypeString,
										Optional: true,
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
										Type:     schema.TypeInt,
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
										Type:     schema.TypeString,
										Optional: true,
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
										Type:     schema.TypeInt,
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
										Type:     schema.TypeString,
										Optional: true,
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
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_metric_rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interval": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"no_effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"silence_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCmsGroupMetricRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutGroupMetricRule"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request["Category"] = d.Get("category")
	if v, ok := d.GetOk("contact_groups"); ok {
		request["ContactGroups"] = v
	}

	if v, ok := d.GetOk("dimensions"); ok {
		request["Dimensions"] = v
	}

	if v, ok := d.GetOk("effective_interval"); ok {
		request["EffectiveInterval"] = v
	}

	if v, ok := d.GetOk("email_subject"); ok {
		request["EmailSubject"] = v
	}

	if v, ok := d.GetOk("escalations"); ok {
		if v != nil {
			escalationsMap := make(map[string]interface{})
			for _, escalations := range v.(*schema.Set).List() {
				escalationsArg := escalations.(map[string]interface{})
				if escalationsArg["critical"] != nil {
					criticalMap := make(map[string]interface{})
					for _, critical := range escalationsArg["critical"].(*schema.Set).List() {
						criticalArg := critical.(map[string]interface{})
						criticalMap["ComparisonOperator"] = criticalArg["comparison_operator"].(string)
						criticalMap["Statistics"] = criticalArg["statistics"].(string)
						criticalMap["Threshold"] = criticalArg["threshold"].(string)
						criticalMap["Times"] = requests.NewInteger(criticalArg["times"].(int))
					}
					escalationsMap["Critical"] = criticalMap
				}
				if escalationsArg["info"] != nil {
					infoMap := make(map[string]interface{})
					for _, info := range escalationsArg["info"].(*schema.Set).List() {
						infoArg := info.(map[string]interface{})
						infoMap["ComparisonOperator"] = infoArg["comparison_operator"].(string)
						infoMap["Statistics"] = infoArg["statistics"].(string)
						infoMap["Threshold"] = infoArg["threshold"].(string)
						infoMap["Times"] = requests.NewInteger(infoArg["times"].(int))
					}
					escalationsMap["Info"] = infoMap
				}
				if escalationsArg["warn"] != nil {
					warnMap := make(map[string]interface{})
					for _, warn := range escalationsArg["warn"].(*schema.Set).List() {
						warnArg := warn.(map[string]interface{})
						warnMap["ComparisonOperator"] = warnArg["comparison_operator"].(string)
						warnMap["Statistics"] = warnArg["statistics"].(string)
						warnMap["Threshold"] = warnArg["threshold"].(string)
						warnMap["Times"] = requests.NewInteger(warnArg["times"].(int))
					}
					escalationsMap["Warn"] = warnMap
				}
			}
			request["Escalations"] = escalationsMap
		}
	}
	request["GroupId"] = d.Get("group_id")
	request["RuleName"] = d.Get("group_metric_rule_name")
	if v, ok := d.GetOk("interval"); ok {
		request["Interval"] = v
	}

	request["MetricName"] = d.Get("metric_name")
	request["Namespace"] = d.Get("namespace")
	if v, ok := d.GetOk("no_effective_interval"); ok {
		request["NoEffectiveInterval"] = v
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}

	request["RuleId"] = d.Get("rule_id")
	if v, ok := d.GetOk("silence_time"); ok {
		request["SilenceTime"] = v
	}

	if v, ok := d.GetOk("webhook"); ok {
		request["Webhook"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ExceedingQuota", "Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_group_metric_rule", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("PutGroupMetricRule failed for " + response["Message"].(string)))
	}

	d.SetId(fmt.Sprint(request["RuleId"]))

	return resourceAlicloudCmsGroupMetricRuleRead(d, meta)
}
func resourceAlicloudCmsGroupMetricRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsGroupMetricRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_group_metric_rule cmsService.DescribeCmsGroupMetricRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule_id", d.Id())
	d.Set("contact_groups", object["ContactGroups"])
	d.Set("dimensions", object["Dimensions"])
	d.Set("effective_interval", object["EffectiveInterval"])
	d.Set("email_subject", object["MailSubject"])

	escalationsSli := make([]map[string]interface{}, 0)
	if len(object["Escalations"].(map[string]interface{})) > 0 {
		escalations := object["Escalations"]
		escalationsMap := make(map[string]interface{})

		criticalSli := make([]map[string]interface{}, 0)
		if len(escalations.(map[string]interface{})["Critical"].(map[string]interface{})) > 0 {
			critical := escalations.(map[string]interface{})["Critical"]
			criticalMap := make(map[string]interface{})
			criticalMap["comparison_operator"] = critical.(map[string]interface{})["ComparisonOperator"]
			criticalMap["statistics"] = critical.(map[string]interface{})["Statistics"]
			criticalMap["threshold"] = critical.(map[string]interface{})["Threshold"]
			criticalMap["times"] = critical.(map[string]interface{})["Times"]
			criticalSli = append(criticalSli, criticalMap)
		}
		escalationsMap["critical"] = criticalSli

		infoSli := make([]map[string]interface{}, 0)
		if len(escalations.(map[string]interface{})["Info"].(map[string]interface{})) > 0 {
			info := escalations.(map[string]interface{})["Info"]
			infoMap := make(map[string]interface{})
			infoMap["comparison_operator"] = info.(map[string]interface{})["ComparisonOperator"]
			infoMap["statistics"] = info.(map[string]interface{})["Statistics"]
			infoMap["threshold"] = info.(map[string]interface{})["Threshold"]
			infoMap["times"] = info.(map[string]interface{})["Times"]
			infoSli = append(infoSli, infoMap)
		}
		escalationsMap["info"] = infoSli

		warnSli := make([]map[string]interface{}, 0)
		if len(escalations.(map[string]interface{})["Warn"].(map[string]interface{})) > 0 {
			warn := escalations.(map[string]interface{})["Warn"]
			warnMap := make(map[string]interface{})
			warnMap["comparison_operator"] = warn.(map[string]interface{})["ComparisonOperator"]
			warnMap["statistics"] = warn.(map[string]interface{})["Statistics"]
			warnMap["threshold"] = warn.(map[string]interface{})["Threshold"]
			warnMap["times"] = warn.(map[string]interface{})["Times"]
			warnSli = append(warnSli, warnMap)
		}
		escalationsMap["warn"] = warnSli
		escalationsSli = append(escalationsSli, escalationsMap)
	}
	d.Set("escalations", escalationsSli)
	d.Set("group_id", object["GroupId"])
	d.Set("group_metric_rule_name", object["RuleName"])
	d.Set("metric_name", object["MetricName"])
	d.Set("namespace", object["Namespace"])
	d.Set("no_effective_interval", object["NoEffectiveInterval"])
	d.Set("period", formatInt(object["Period"]))
	d.Set("silence_time", formatInt(object["SilenceTime"]))
	d.Set("status", object["AlertState"])
	d.Set("webhook", object["Webhook"])
	return nil
}
func resourceAlicloudCmsGroupMetricRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}
	if d.HasChange("group_id") {
		update = true
	}
	request["GroupId"] = d.Get("group_id")
	if d.HasChange("group_metric_rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("group_metric_rule_name")
	if d.HasChange("metric_name") {
		update = true
	}
	request["MetricName"] = d.Get("metric_name")
	request["Namespace"] = d.Get("namespace")
	if d.HasChange("contact_groups") {
		update = true
		request["ContactGroups"] = d.Get("contact_groups")
	}
	if d.HasChange("dimensions") {
		update = true
		request["Dimensions"] = d.Get("dimensions")
	}
	if d.HasChange("effective_interval") {
		update = true
		request["EffectiveInterval"] = d.Get("effective_interval")
	}
	if d.HasChange("email_subject") {
		update = true
		request["EmailSubject"] = d.Get("email_subject")
	}
	if d.HasChange("escalations") {
		update = true
		if d.Get("escalations") != nil {
			escalationsMap := make(map[string]interface{})
			for _, escalations := range d.Get("escalations").(*schema.Set).List() {
				escalationsArg := escalations.(map[string]interface{})
				if escalationsArg["critical"] != nil {
					criticalMap := make(map[string]interface{})
					for _, critical := range escalationsArg["critical"].(*schema.Set).List() {
						criticalArg := critical.(map[string]interface{})
						criticalMap["ComparisonOperator"] = criticalArg["comparison_operator"].(string)
						criticalMap["Statistics"] = criticalArg["statistics"].(string)
						criticalMap["Threshold"] = criticalArg["threshold"].(string)
						criticalMap["Times"] = requests.NewInteger(criticalArg["times"].(int))
					}
					escalationsMap["Critical"] = criticalMap
				}
				if escalationsArg["info"] != nil {
					infoMap := make(map[string]interface{})
					for _, info := range escalationsArg["info"].(*schema.Set).List() {
						infoArg := info.(map[string]interface{})
						infoMap["ComparisonOperator"] = infoArg["comparison_operator"].(string)
						infoMap["Statistics"] = infoArg["statistics"].(string)
						infoMap["Threshold"] = infoArg["threshold"].(string)
						infoMap["Times"] = requests.NewInteger(infoArg["times"].(int))
					}
					escalationsMap["Info"] = infoMap
				}
				if escalationsArg["warn"] != nil {
					warnMap := make(map[string]interface{})
					for _, warn := range escalationsArg["warn"].(*schema.Set).List() {
						warnArg := warn.(map[string]interface{})
						warnMap["ComparisonOperator"] = warnArg["comparison_operator"].(string)
						warnMap["Statistics"] = warnArg["statistics"].(string)
						warnMap["Threshold"] = warnArg["threshold"].(string)
						warnMap["Times"] = requests.NewInteger(warnArg["times"].(int))
					}
					escalationsMap["Warn"] = warnMap
				}
			}
			request["Escalations"] = escalationsMap
		}
	}
	if d.HasChange("no_effective_interval") {
		update = true
		request["NoEffectiveInterval"] = d.Get("no_effective_interval")
	}
	if d.HasChange("period") {
		update = true
		request["Period"] = d.Get("period")
	}
	if d.HasChange("silence_time") {
		update = true
		request["SilenceTime"] = d.Get("silence_time")
	}
	if d.HasChange("webhook") {
		update = true
		request["Webhook"] = d.Get("webhook")
	}
	if update {
		request["Category"] = d.Get("category")
		if _, ok := d.GetOk("interval"); ok {
			request["Interval"] = d.Get("interval")
		}
		action := "PutGroupMetricRule"
		conn, err := client.NewCmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ExceedingQuota", "Throttling.User"}) || NeedRetry(err) {
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
		if fmt.Sprintf(`%v`, response["Code"]) != "200" {
			return WrapError(Error("PutGroupMetricRule failed for " + response["Message"].(string)))
		}
	}
	return resourceAlicloudCmsGroupMetricRuleRead(d, meta)
}
func resourceAlicloudCmsGroupMetricRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMetricRules"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id": []string{d.Id()},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"ExceedingQuota", "Throttling.User"}) || NeedRetry(err) {
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
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"400", "403", "404", "ResourceNotFound"}) {
		return nil
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("DeleteMetricRules failed for " + response["Message"].(string)))
	}
	return nil
}
