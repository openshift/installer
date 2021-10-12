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

func resourceAlicloudEventBridgeEventSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEventBridgeEventSourceCreate,
		Read:   resourceAlicloudEventBridgeEventSourceRead,
		Update: resourceAlicloudEventBridgeEventSourceUpdate,
		Delete: resourceAlicloudEventBridgeEventSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_bus_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event_source_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"external_source_config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"external_source_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"linked_external_source": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEventBridgeEventSourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateEventSource"
	request := make(map[string]interface{})
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["EventSourceName"] = d.Get("event_source_name")
	request["EventBusName"] = d.Get("event_bus_name")
	if v, ok := d.GetOk("external_source_config"); ok {
		if v, err := convertMaptoJsonString(v.(map[string]interface{})); err == nil {
			request["ExternalSourceConfig"] = v
		} else {
			return WrapError(err)
		}
	}
	if v, ok := d.GetOk("external_source_type"); ok {
		request["ExternalSourceType"] = v
	}
	if v, ok := d.GetOkExists("linked_external_source"); ok {
		request["LinkedExternalSource"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_event_bridge_event_source", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["EventSourceName"]))

	return resourceAlicloudEventBridgeEventSourceRead(d, meta)
}
func resourceAlicloudEventBridgeEventSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eventbridgeService := EventbridgeService{client}
	object, err := eventbridgeService.DescribeEventBridgeEventSource(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_event_bridge_event_source eventbridgeService.DescribeEventBridgeEventSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("event_source_name", d.Id())
	d.Set("description", object["Description"])
	d.Set("event_bus_name", object["EventBusName"])
	d.Set("external_source_config", object["ExternalSourceConfig"])
	d.Set("external_source_type", object["ExternalSourceType"])
	d.Set("linked_external_source", object["LinkedExternalSource"])
	return nil
}
func resourceAlicloudEventBridgeEventSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"EventSourceName": d.Id(),
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("external_source_config") {
		update = true
		if v, err := convertMaptoJsonString(d.Get("external_source_config").(map[string]interface{})); err == nil {
			request["ExternalSourceConfig"] = v
		} else {
			return WrapError(err)
		}
	}
	request["ExternalSourceType"] = d.Get("external_source_type")
	if d.HasChange("external_source_type") {
		update = true
	}
	request["LinkedExternalSource"] = d.Get("linked_external_source")
	if d.HasChange("linked_external_source") || d.IsNewResource() {
		update = true
	}

	if update {
		action := "UpdateEventSource"
		conn, err := client.NewEventbridgeClient()
		if err != nil {
			return WrapError(err)
		}
		if v, ok := d.GetOk("event_bus_name"); ok {
			request["EventBusName"] = v
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
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudEventBridgeEventSourceRead(d, meta)
}
func resourceAlicloudEventBridgeEventSourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventSource"
	var response map[string]interface{}
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"EventSourceName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
