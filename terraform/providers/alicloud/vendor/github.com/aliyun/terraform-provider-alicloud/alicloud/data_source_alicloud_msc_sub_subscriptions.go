package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMscSubSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMscSubSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"item_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"item_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pmsg_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sms_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tts_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"webhook_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"webhook_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMscSubSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSubscriptionItems"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("filter"); ok {
		request["Filter"] = v
	}
	request["Locale"] = "en"
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewMscopensubscriptionClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2021-07-13"), StringPointer("AK"), request, nil, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_msc_sub_subscriptions", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.SubscriptionItems", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SubscriptionItems", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"channel":        object["Channel"],
			"contact_ids":    object["ContactIds"],
			"description":    object["Description"],
			"email_status":   formatInt(object["EmailStatus"]),
			"id":             fmt.Sprint(object["ItemId"]),
			"item_id":        fmt.Sprint(object["ItemId"]),
			"item_name":      object["ItemName"],
			"pmsg_status":    formatInt(object["PmsgStatus"]),
			"sms_status":     formatInt(object["SmsStatus"]),
			"tts_status":     formatInt(object["TtsStatus"]),
			"webhook_ids":    object["WebhookIds"],
			"webhook_status": formatInt(object["WebhookStatus"]),
		}
		s = append(s, mapping)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("subscriptions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
