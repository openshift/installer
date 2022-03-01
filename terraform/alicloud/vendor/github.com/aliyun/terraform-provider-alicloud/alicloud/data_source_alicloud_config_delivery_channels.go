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

func dataSourceAlicloudConfigDeliveryChannels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigDeliveryChannelsRead,
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
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delivery_channel_assume_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_target_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudConfigDeliveryChannelsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDeliveryChannels"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var deliveryChannelNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		deliveryChannelNameRegex = r
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
	status, statusOk := d.GetOkExists("status")
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-01-08"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_delivery_channels", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.DeliveryChannels", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DeliveryChannels", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if deliveryChannelNameRegex != nil {
			if !deliveryChannelNameRegex.MatchString(item["DeliveryChannelName"].(string)) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DeliveryChannelId"])]; !ok {
				continue
			}
		}
		if statusOk && status.(int) != formatInt(item["Status"]) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"delivery_channel_assume_role_arn": object["DeliveryChannelAssumeRoleArn"],
			"delivery_channel_condition":       object["DeliveryChannelCondition"],
			"id":                               fmt.Sprint(object["DeliveryChannelId"]),
			"delivery_channel_id":              fmt.Sprint(object["DeliveryChannelId"]),
			"delivery_channel_name":            object["DeliveryChannelName"],
			"delivery_channel_target_arn":      object["DeliveryChannelTargetArn"],
			"delivery_channel_type":            object["DeliveryChannelType"],
			"description":                      object["Description"],
			"status":                           formatInt(object["Status"]),
		}
		ids = append(ids, fmt.Sprint(object["DeliveryChannelId"]))
		names = append(names, object["DeliveryChannelName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("channels", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
