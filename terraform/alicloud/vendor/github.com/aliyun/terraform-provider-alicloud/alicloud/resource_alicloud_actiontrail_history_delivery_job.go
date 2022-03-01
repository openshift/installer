package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudActiontrailHistoryDeliveryJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudActiontrailHistoryDeliveryJobCreate,
		Read:   resourceAlicloudActiontrailHistoryDeliveryJobRead,
		Delete: resourceAlicloudActiontrailHistoryDeliveryJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"trail_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudActiontrailHistoryDeliveryJobCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDeliveryHistoryJob"
	request := make(map[string]interface{})
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	request["TrailName"] = d.Get("trail_name")
	request["ClientToken"] = buildClientToken("CreateDeliveryHistoryJob")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrail_history_delivery_job", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(formatInt(response["JobId"])))
	actiontrailService := ActiontrailService{client}
	stateConf := BuildStateConf([]string{}, []string{"2", "3"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, actiontrailService.ActiontrailHistoryDeliveryJobStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudActiontrailHistoryDeliveryJobRead(d, meta)
}
func resourceAlicloudActiontrailHistoryDeliveryJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}
	object, err := actiontrailService.DescribeActiontrailHistoryDeliveryJob(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_history_delivery_job actiontrailService.DescribeActiontrailHistoryDeliveryJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", formatInt(object["JobStatus"]))
	d.Set("trail_name", object["TrailName"])
	return nil
}
func resourceAlicloudActiontrailHistoryDeliveryJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDeliveryHistoryJob"
	var response map[string]interface{}
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"JobId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
