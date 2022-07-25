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

func dataSourceAlicloudQuotasQuotaAlarms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudQuotasQuotaAlarmsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"quota_alarm_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"quota_dimensions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_action_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_alarm_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_dimensions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"threshold_percent": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"web_hook": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudQuotasQuotaAlarmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListQuotaAlarms"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	if v, ok := d.GetOk("product_code"); ok {
		request["ProductCode"] = v
	}
	if v, ok := d.GetOk("quota_action_code"); ok {
		request["QuotaActionCode"] = v
	}
	if v, ok := d.GetOk("quota_alarm_name"); ok {
		request["AlarmName"] = v
	}
	if v, ok := d.GetOk("quota_dimensions"); ok {
		quotaDimensionsMaps := make([]map[string]interface{}, 0)
		for _, quotaDimensions := range v.([]interface{}) {
			quotaDimensionsMap := make(map[string]interface{})
			quotaDimensionsArg := quotaDimensions.(map[string]interface{})
			quotaDimensionsMap["Key"] = quotaDimensionsArg["key"]
			quotaDimensionsMap["Value"] = quotaDimensionsArg["value"]
			quotaDimensionsMaps = append(quotaDimensionsMaps, quotaDimensionsMap)
		}
		request["QuotaDimensions"] = quotaDimensionsMaps

	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var quotaAlarmNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		quotaAlarmNameRegex = r
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
	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_quotas_quota_alarms", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.QuotaAlarms", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.QuotaAlarms", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if quotaAlarmNameRegex != nil {
				if !quotaAlarmNameRegex.MatchString(fmt.Sprint(item["AlarmName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AlarmId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["AlarmId"]),
			"alarm_id":          fmt.Sprint(object["AlarmId"]),
			"product_code":      object["ProductCode"],
			"quota_action_code": object["QuotaActionCode"],
			"quota_alarm_name":  object["AlarmName"],
			"threshold":         object["Threshold"],
			"threshold_percent": object["ThresholdPercent"],
			"web_hook":          object["WebHook"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["AlarmId"]))
			names = append(names, object["AlarmName"])
			s = append(s, mapping)
			continue
		}

		quotasService := QuotasService{client}
		id := fmt.Sprint(object["AlarmId"])
		getResp, err := quotasService.DescribeQuotasQuotaAlarm(id)
		if err != nil {
			return WrapError(err)
		}

		quotaDimensionList := make([]map[string]interface{}, 0)
		if quotaDimension, ok := getResp["QuotaDimension"]; ok {
			for k, v := range quotaDimension.(map[string]interface{}) {
				quotaDimensionMap := make(map[string]interface{})
				quotaDimensionMap["key"] = k
				quotaDimensionMap["value"] = v
				quotaDimensionList = append(quotaDimensionList, quotaDimensionMap)
			}
		}

		mapping["quota_dimensions"] = quotaDimensionList
		ids = append(ids, fmt.Sprint(object["AlarmId"]))
		names = append(names, object["AlarmName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("alarms", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
