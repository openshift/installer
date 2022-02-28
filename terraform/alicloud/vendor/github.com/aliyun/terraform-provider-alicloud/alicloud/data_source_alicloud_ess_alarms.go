package alicloud

import (
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEssAlarms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssAlarmsRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "custom"}, false),
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"alarm_actions": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeInt,
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
						"comparison_operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"evaluation_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cloud_monitor_group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssAlarmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDescribeAlarmsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request.ScalingGroupId = scalingGroupId.(string)
	}
	if metric_type, ok := d.GetOk("metric_type"); ok && metric_type.(string) != "" {
		request.MetricType = metric_type.(string)
	}
	var allAlarms []ess.Alarm
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeAlarms(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_alarms", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response := raw.(*ess.DescribeAlarmsResponse)
		if len(response.AlarmList.Alarm) < 1 {
			break
		}
		allAlarms = append(allAlarms, response.AlarmList.Alarm...)
		if len(response.AlarmList.Alarm) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	var filteredAlarms = make([]ess.Alarm, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, alarm := range allAlarms {
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if r != nil && !r.MatchString(alarm.Name) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[alarm.AlarmTaskId]; !ok {
					continue
				}
			}
			filteredAlarms = append(filteredAlarms, alarm)
		}
	} else {
		filteredAlarms = allAlarms
	}

	return alarmsDescriptionAttribute(d, filteredAlarms, meta)
}

func alarmsDescriptionAttribute(d *schema.ResourceData, alarms []ess.Alarm, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, alarm := range alarms {
		dims := make([]ess.Dimension, 0, len(alarm.Dimensions.Dimension))
		for _, dimension := range alarm.Dimensions.Dimension {
			dims = append(dims, dimension)
		}
		client := meta.(*connectivity.AliyunClient)
		essService := EssService{client}
		mapping := map[string]interface{}{
			"state":               alarm.State,
			"id":                  alarm.AlarmTaskId,
			"name":                alarm.Name,
			"scaling_group_id":    alarm.ScalingGroupId,
			"metric_name":         alarm.MetricName,
			"description":         alarm.Description,
			"enable":              alarm.Enable,
			"alarm_actions":       alarm.AlarmActions.AlarmAction,
			"metric_type":         alarm.MetricType,
			"period":              alarm.Period,
			"statistics":          alarm.Statistics,
			"threshold":           strconv.FormatFloat(alarm.Threshold, 'f', -1, 32),
			"comparison_operator": alarm.ComparisonOperator,
			"evaluation_count":    alarm.EvaluationCount,
			"dimensions":          essService.flattenDimensionsToMap(dims),
		}
		ids = append(ids, alarm.AlarmTaskId)
		names = append(names, alarm.AlarmTaskName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("alarms", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
