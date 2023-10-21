package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAlarmCreate,
		Read:   resourceAlicloudCmsAlarmRead,
		Update: resourceAlicloudCmsAlarmUpdate,
		Delete: resourceAlicloudCmsAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     schema.TypeString,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"escalations_critical": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientCriticalSuppressFunc,
				MaxItems:         1,
			},
			"escalations_warn": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientWarnSuppressFunc,
				MaxItems:         1,
			},
			"escalations_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientInfoSuppressFunc,
				MaxItems:         1,
			},
			"statistics": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
				Deprecated:   "Field 'statistics' has been deprecated from provider version 1.94.0. New field 'escalations_critical.statistics' instead.",
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, Equal, NotEqual,
				}, false),
				Deprecated: "Field 'operator' has been deprecated from provider version 1.94.0. New field 'escalations_critical.comparison_operator' instead.",
			},
			"threshold": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'threshold' has been deprecated from provider version 1.94.0. New field 'escalations_critical.threshold' instead.",
			},
			"triggered_count": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'triggered_count' has been deprecated from provider version 1.94.0. New field 'escalations_critical.times' instead.",
			},
			"contact_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:      0,
				//ValidateFunc: validation.IntBetween(0, 24),
				Deprecated: "Field 'start_time' has been deprecated from provider version 1.50.0. New field 'effective_interval' instead.",
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:      24,
				//ValidateFunc: validation.IntBetween(0, 24),
				Deprecated: "Field 'end_time' has been deprecated from provider version 1.50.0. New field 'effective_interval' instead.",
			},
			"effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "00:00-23:59",
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validation.IntBetween(300, 86400),
			},

			"notify_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Removed:      "Field 'notify_type' has been removed from provider version 1.50.0.",
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

func resourceAlicloudCmsAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(resource.UniqueId())
	return resourceAlicloudCmsAlarmUpdate(d, meta)
}

func resourceAlicloudCmsAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	alarm, err := cmsService.DescribeAlarm(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", alarm.RuleName)
	d.Set("project", alarm.Namespace)
	d.Set("metric", alarm.MetricName)
	if period, err := strconv.Atoi(alarm.Period); err != nil {
		return WrapError(err)
	} else {
		d.Set("period", period)
	}

	d.Set("statistics", alarm.Escalations.Critical.Statistics)
	oper := convertOperator(alarm.Escalations.Critical.ComparisonOperator)
	if oper == MoreThan && d.Get("operator").(string) == Equal {
		oper = Equal
	}
	d.Set("operator", oper)
	d.Set("threshold", alarm.Escalations.Critical.Threshold)
	d.Set("triggered_count", alarm.Escalations.Critical.Times)

	escalationsCritical := make([]map[string]interface{}, 1)

	mapping := map[string]interface{}{
		"statistics":          alarm.Escalations.Critical.Statistics,
		"comparison_operator": convertOperator(alarm.Escalations.Critical.ComparisonOperator),
		"threshold":           alarm.Escalations.Critical.Threshold,
		"times":               alarm.Escalations.Critical.Times,
	}
	escalationsCritical[0] = mapping
	d.Set("escalations_critical", escalationsCritical)

	escalationsWarn := make([]map[string]interface{}, 1)
	if alarm.Escalations.Warn.Times != "" {
		if count, err := strconv.Atoi(alarm.Escalations.Warn.Times); err != nil {
			return WrapError(err)
		} else {
			mappingWarn := map[string]interface{}{
				"statistics":          alarm.Escalations.Warn.Statistics,
				"comparison_operator": convertOperator(alarm.Escalations.Warn.ComparisonOperator),
				"threshold":           alarm.Escalations.Warn.Threshold,
				"times":               count,
			}
			escalationsWarn[0] = mappingWarn
			d.Set("escalations_warn", escalationsWarn)
		}
	}

	escalationsInfo := make([]map[string]interface{}, 1)
	if alarm.Escalations.Info.Times != "" {
		if count, err := strconv.Atoi(alarm.Escalations.Info.Times); err != nil {
			return WrapError(err)
		} else {
			mappingInfo := map[string]interface{}{
				"statistics":          alarm.Escalations.Info.Statistics,
				"comparison_operator": convertOperator(alarm.Escalations.Info.ComparisonOperator),
				"threshold":           alarm.Escalations.Info.Threshold,
				"times":               count,
			}
			escalationsInfo[0] = mappingInfo
			d.Set("escalations_info", escalationsInfo)
		}
	}

	d.Set("effective_interval", alarm.EffectiveInterval)
	//d.Set("start_time", parts[0])
	//d.Set("end_time", parts[1])

	d.Set("silence_time", alarm.SilenceTime)

	d.Set("status", alarm.AlertState)
	d.Set("enabled", alarm.EnableState)
	d.Set("webhook", alarm.Webhook)
	d.Set("contact_groups", strings.Split(alarm.ContactGroups, ","))

	var dims []string
	if alarm.Dimensions != "" {
		if err := json.Unmarshal([]byte(alarm.Dimensions), &dims); err != nil {
			return fmt.Errorf("Unmarshaling Dimensions got an error: %#v.", err)
		}
	}
	d.Set("dimensions", dims)

	return nil
}

func resourceAlicloudCmsAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	d.Partial(true)

	request := cms.CreatePutResourceMetricRuleRequest()
	request.RuleId = d.Id()
	request.RuleName = d.Get("name").(string)
	request.Namespace = d.Get("project").(string)
	request.MetricName = d.Get("metric").(string)
	request.Period = strconv.Itoa(d.Get("period").(int))
	request.ContactGroups = strings.Join(expandStringList(d.Get("contact_groups").([]interface{})), ",")

	// 兼容弃用参数
	request.EscalationsCriticalStatistics = d.Get("statistics").(string)
	request.EscalationsCriticalComparisonOperator = convertOperator(d.Get("operator").(string))
	if v, ok := d.GetOk("threshold"); ok && v.(string) != "" {
		request.EscalationsCriticalThreshold = v.(string)
	}
	request.EscalationsCriticalThreshold = d.Get("threshold").(string)
	request.EscalationsCriticalTimes = requests.NewInteger(d.Get("triggered_count").(int))

	// Critical
	if v, ok := d.GetOk("escalations_critical"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request.EscalationsCriticalStatistics = val["statistics"].(string)
			request.EscalationsCriticalComparisonOperator = convertOperator(val["comparison_operator"].(string))
			request.EscalationsCriticalThreshold = val["threshold"].(string)
			request.EscalationsCriticalTimes = requests.NewInteger(val["times"].(int))
		}
	}
	// Warn
	if v, ok := d.GetOk("escalations_warn"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request.EscalationsWarnStatistics = val["statistics"].(string)
			request.EscalationsWarnComparisonOperator = convertOperator(val["comparison_operator"].(string))
			request.EscalationsWarnThreshold = val["threshold"].(string)
			request.EscalationsWarnTimes = requests.NewInteger(val["times"].(int))
		}
	}
	// Info
	if v, ok := d.GetOk("escalations_info"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request.EscalationsInfoStatistics = val["statistics"].(string)
			request.EscalationsInfoComparisonOperator = convertOperator(val["comparison_operator"].(string))
			request.EscalationsInfoThreshold = val["threshold"].(string)
			request.EscalationsInfoTimes = requests.NewInteger(val["times"].(int))

		}
	}

	if v, ok := d.GetOk("effective_interval"); ok && v.(string) != "" {
		request.EffectiveInterval = v.(string)
	} else {
		start, startOk := d.GetOk("start_time")
		end, endOk := d.GetOk("end_time")
		if startOk && endOk && end.(int) > 0 {
			// The EffectiveInterval valid value between 00:00 and 23:59
			request.EffectiveInterval = fmt.Sprintf("%d:00-%d:59", start.(int), end.(int)-1)
		}
	}
	request.SilenceTime = requests.NewInteger(d.Get("silence_time").(int))

	if webhook, ok := d.GetOk("webhook"); ok && webhook.(string) != "" {
		request.Webhook = webhook.(string)
	}

	var dimList []map[string]string
	if dimensions, ok := d.GetOk("dimensions"); ok {
		for k, v := range dimensions.(map[string]interface{}) {
			values := strings.Split(v.(string), COMMA_SEPARATED)
			if len(values) > 0 {
				for _, vv := range values {
					dimList = append(dimList, map[string]string{k: Trim(vv)})
				}
			} else {
				dimList = append(dimList, map[string]string{k: Trim(v.(string))})
			}

		}
	}
	if len(dimList) > 0 {
		if bytes, err := json.Marshal(dimList); err != nil {
			return fmt.Errorf("Marshaling dimensions to json string got an error: %#v.", err)
		} else {
			request.Resources = string(bytes[:])
		}
	}
	_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.PutResourceMetricRule(request)
	})
	if err != nil {
		return fmt.Errorf("Putting alarm got an error: %#v", err)
	}
	d.SetPartial("name")
	d.SetPartial("period")
	d.SetPartial("statistics")
	d.SetPartial("operator")
	d.SetPartial("threshold")
	d.SetPartial("triggered_count")
	d.SetPartial("contact_groups")
	d.SetPartial("effective_interval")
	d.SetPartial("start_time")
	d.SetPartial("end_time")
	d.SetPartial("silence_time")
	d.SetPartial("notify_type")
	d.SetPartial("webhook")

	if d.Get("enabled").(bool) {
		request := cms.CreateEnableMetricRulesRequest()
		request.RuleId = &[]string{d.Id()}

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.EnableMetricRules(request)
			})

			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(fmt.Errorf("Enabling alarm got an error: %#v", err))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Enabling alarm got an error: %#v", err)
		}
	} else {
		request := cms.CreateDisableMetricRulesRequest()
		request.RuleId = &[]string{d.Id()}

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.DisableMetricRules(request)
			})

			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(fmt.Errorf("Disableing alarm got an error: %#v", err))
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Disableing alarm got an error: %#v", err)
		}
	}
	if err := cmsService.WaitForCmsAlarm(d.Id(), d.Get("enabled").(bool), 102); err != nil {
		return err
	}

	d.Partial(false)

	return resourceAlicloudCmsAlarmRead(d, meta)
}

func resourceAlicloudCmsAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	request := cms.CreateDeleteMetricRulesRequest()

	request.Id = &[]string{d.Id()}

	wait := incrementalWait(1*time.Second, 2*time.Second)
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteMetricRules(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
		}

		_, err = cmsService.DescribeAlarm(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe alarm rule got an error: %#v", err))
		}

		return resource.RetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
	})
}

func convertOperator(operator string) string {
	switch operator {
	case MoreThan:
		return "GreaterThanThreshold"
	case MoreThanOrEqual:
		return "GreaterThanOrEqualToThreshold"
	case LessThan:
		return "LessThanThreshold"
	case LessThanOrEqual:
		return "LessThanOrEqualToThreshold"
	case NotEqual:
		return "NotEqualToThreshold"
	case Equal:
		return "GreaterThanThreshold"
	case "GreaterThanThreshold":
		return MoreThan
	case "GreaterThanOrEqualToThreshold":
		return MoreThanOrEqual
	case "LessThanThreshold":
		return LessThan
	case "LessThanOrEqualToThreshold":
		return LessThanOrEqual
	case "NotEqualToThreshold":
		return NotEqual
	default:
		return ""
	}
}
