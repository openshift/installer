package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCdnRealTimeLogDeliveries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCdnRealTimeLogDeliveriesRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"offline", "online"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deliveries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logstore": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_region": {
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

func dataSourceAlicloudCdnRealTimeLogDeliveriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDomainRealtimeLogDelivery"
	request := map[string]interface{}{}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}

	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), request, nil, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_real_time_log_deliveries", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$", response)
	}

	s := make([]map[string]interface{}, 0)
	if object, ok := resp.(map[string]interface{}); ok {
		if statusOk && status.(string) != "" && status.(string) == object["Status"].(string) {
			mapping := map[string]interface{}{
				"id":         fmt.Sprint(request["Domain"]),
				"domain":     fmt.Sprint(request["Domain"]),
				"logstore":   object["Logstore"],
				"project":    object["Project"],
				"sls_region": object["Region"],
				"status":     object["Status"],
			}
			s = append(s, mapping)
		}
	}

	d.SetId(fmt.Sprint(request["Domain"]))

	if err := d.Set("deliveries", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
