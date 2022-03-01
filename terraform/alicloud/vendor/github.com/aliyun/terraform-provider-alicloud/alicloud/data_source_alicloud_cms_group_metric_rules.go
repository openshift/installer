package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCmsGroupMetricRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsGroupMetricRulesRead,
		Schema: map[string]*schema.Schema{
			"dimensions": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_state": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"group_metric_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALARM", "INSUFFICIENT_DATA", "OK"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact_groups": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effective_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_state": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"escalations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"critical": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"statistics": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"times": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"info": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"statistics": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"times": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"warn": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"statistics": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"times": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_metric_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"no_effective_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"silence_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"webhook": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCmsGroupMetricRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeMetricRuleList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("dimensions"); ok {
		request["Dimensions"] = v
	}
	if v, ok := d.GetOkExists("enable_state"); ok {
		request["EnableState"] = v
	}
	if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	}
	if v, ok := d.GetOk("group_metric_rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("metric_name"); ok {
		request["MetricName"] = v
	}
	if v, ok := d.GetOk("namespace"); ok {
		request["Namespace"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["AlertState"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["Page"] = 1
	var objects []map[string]interface{}
	var groupMetricRuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		groupMetricRuleNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_group_metric_rules", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Alarms.Alarm", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Alarms.Alarm", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if groupMetricRuleNameRegex != nil {
				if !groupMetricRuleNameRegex.MatchString(item["RuleName"].(string)) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RuleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["Page"] = request["Page"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"contact_groups":         object["ContactGroups"],
			"dimensions":             object["Dimensions"],
			"effective_interval":     object["EffectiveInterval"],
			"email_subject":          object["MailSubject"],
			"enable_state":           object["EnableState"],
			"group_id":               object["GroupId"],
			"group_metric_rule_name": object["RuleName"],
			"metric_name":            object["MetricName"],
			"namespace":              object["Namespace"],
			"no_effective_interval":  object["NoEffectiveInterval"],
			"period":                 formatInt(object["Period"]),
			"resources":              object["Resources"],
			"id":                     fmt.Sprint(object["RuleId"]),
			"rule_id":                fmt.Sprint(object["RuleId"]),
			"silence_time":           formatInt(object["SilenceTime"]),
			"source_type":            object["SourceType"],
			"status":                 object["AlertState"],
			"webhook":                object["Webhook"],
		}

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
		mapping["escalations"] = escalationsSli
		ids = append(ids, fmt.Sprint(object["RuleId"]))
		names = append(names, object["RuleName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
