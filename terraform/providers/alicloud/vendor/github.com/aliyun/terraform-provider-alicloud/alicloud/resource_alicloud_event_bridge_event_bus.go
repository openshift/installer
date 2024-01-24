package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEventBridgeEventBus() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEventBridgeEventBusCreate,
		Read:   resourceAlicloudEventBridgeEventBusRead,
		Update: resourceAlicloudEventBridgeEventBusUpdate,
		Delete: resourceAlicloudEventBridgeEventBusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_bus_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 127),
			},
		},
	}
}

func resourceAlicloudEventBridgeEventBusCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateEventBus"
	request := make(map[string]interface{})
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["EventBusName"] = d.Get("event_bus_name")
	request["ClientToken"] = buildClientToken("CreateEventBus")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_event_bus", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("CreateEventBus failed, response: %v", response))
	}

	d.SetId(fmt.Sprint(request["EventBusName"]))

	return resourceAlicloudEventBridgeEventBusRead(d, meta)
}
func resourceAlicloudEventBridgeEventBusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	object, err := eventbridgeService.DescribeEventBridgeEventBus(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_bus eventbridgeService.DescribeEventBridgeEventBus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("event_bus_name", d.Id())
	d.Set("description", object["Description"])
	return nil
}
func resourceAlicloudEventBridgeEventBusUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	if d.HasChange("description") {
		request := map[string]interface{}{
			"EventBusName": d.Id(),
		}
		request["Description"] = d.Get("description")
		action := "UpdateEventBus"
		conn, err := client.NewEventbridgeClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("UpdateEventBus failed, response: %v", response))
		}
	}
	return resourceAlicloudEventBridgeEventBusRead(d, meta)
}
func resourceAlicloudEventBridgeEventBusDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventBus"
	var response map[string]interface{}
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EventBusName": d.Id(),
	}

	request["ClientToken"] = buildClientToken("DeleteEventBus")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"EventBusNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("DeleteEventBus failed, response: %v", response))
	}
	return nil
}
