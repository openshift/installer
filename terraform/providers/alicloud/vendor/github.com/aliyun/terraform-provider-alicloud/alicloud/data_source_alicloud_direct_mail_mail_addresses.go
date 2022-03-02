package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDirectMailMailAddresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDirectMailMailAddressesRead,
		Schema: map[string]*schema.Schema{
			"key_word": {
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
			"sendtype": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"batch", "trigger"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"daily_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"daily_req_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mail_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"month_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"month_req_count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reply_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reply_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sendtype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDirectMailMailAddressesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "QueryMailAddressByParam"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("key_word"); ok {
		request["KeyWord"] = v
	}
	if v, ok := d.GetOk("sendtype"); ok {
		request["Sendtype"] = v
	}
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewDmClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-11-23"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_direct_mail_mail_addresses", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.data.mailAddress", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.data.mailAddress", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["MailAddressId"])]; !ok {
				continue
			}
		}
		if statusOk && status.(string) != "" && status.(string) != string(item["AccountStatus"].(json.Number)) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_name":    object["AccountName"],
			"create_time":     object["CreateTime"],
			"daily_count":     object["DailyCount"],
			"daily_req_count": object["DailyReqCount"],
			"domain_status":   object["DomainStatus"],
			"id":              fmt.Sprint(object["MailAddressId"]),
			"mail_address_id": fmt.Sprint(object["MailAddressId"]),
			"month_count":     object["MonthCount"],
			"month_req_count": object["MonthReqCount"],
			"reply_address":   object["ReplyAddress"],
			"reply_status":    object["ReplyStatus"],
			"sendtype":        object["Sendtype"],
			"status":          object["AccountStatus"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("addresses", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
