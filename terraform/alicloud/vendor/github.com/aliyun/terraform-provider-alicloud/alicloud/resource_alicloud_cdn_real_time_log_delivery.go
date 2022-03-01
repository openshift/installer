package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCdnRealTimeLogDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCdnRealTimeLogDeliveryCreate,
		Read:   resourceAlicloudCdnRealTimeLogDeliveryRead,
		Delete: resourceAlicloudCdnRealTimeLogDeliveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logstore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sls_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCdnRealTimeLogDeliveryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateRealTimeLogDelivery"
	request := make(map[string]interface{})
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request["Domain"] = d.Get("domain")
	request["Logstore"] = d.Get("logstore")
	request["Project"] = d.Get("project")
	request["Region"] = d.Get("sls_region")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cdn_real_time_log_delivery", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Domain"]))

	return resourceAlicloudCdnRealTimeLogDeliveryRead(d, meta)
}
func resourceAlicloudCdnRealTimeLogDeliveryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnService := CdnService{client}
	object, err := cdnService.DescribeCdnRealTimeLogDelivery(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cdn_real_time_log_delivery cdnService.DescribeCdnRealTimeLogDelivery Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain", d.Id())
	d.Set("logstore", object["Logstore"])
	d.Set("project", object["Project"])
	d.Set("sls_region", object["Region"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudCdnRealTimeLogDeliveryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteRealtimeLogDelivery"
	var response map[string]interface{}
	conn, err := client.NewCdnClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Domain": d.Id(),
	}

	request["Logstore"] = d.Get("logstore")
	request["Project"] = d.Get("project")
	request["Region"] = d.Get("sls_region")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2018-05-10"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
