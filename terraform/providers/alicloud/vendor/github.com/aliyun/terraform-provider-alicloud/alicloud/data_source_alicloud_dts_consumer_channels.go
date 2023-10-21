package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDtsConsumerChannels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDtsConsumerChannelsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_group_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumption_checkpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_delay": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unconsumed_data": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDtsConsumerChannelsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeConsumerChannel"
	request := make(map[string]interface{})
	request["DtsInstanceId"] = d.Get("dts_instance_id")
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeMedium
	request["PageNumber"] = 1
	var objects []map[string]interface{}

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
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dts_consumer_channels", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ConsumerChannels", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ConsumerChannels", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["DtsInstanceId"], ":", item["ConsumerGroupId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                       fmt.Sprint(request["DtsInstanceId"], ":", object["ConsumerGroupId"]),
			"consumer_group_id":        fmt.Sprint(object["ConsumerGroupId"]),
			"consumer_group_name":      object["ConsumerGroupName"],
			"consumer_group_user_name": object["ConsumerGroupUserName"],
			"consumption_checkpoint":   fmt.Sprint(object["ConsumptionCheckpoint"]),
			"message_delay":            formatInt(object["MessageDelay"]),
			"unconsumed_data":          formatInt(object["UnconsumedData"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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
