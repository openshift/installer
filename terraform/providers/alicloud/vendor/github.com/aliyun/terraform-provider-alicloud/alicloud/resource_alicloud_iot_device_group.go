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

func resourceAlicloudIotDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudIotDeviceGroupCreate,
		Read:   resourceAlicloudIotDeviceGroupRead,
		Update: resourceAlicloudIotDeviceGroupUpdate,
		Delete: resourceAlicloudIotDeviceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 30),
			},
			"iot_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"super_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudIotDeviceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDeviceGroup"
	request := make(map[string]interface{})
	conn, err := client.NewIotClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("group_desc"); ok {
		request["GroupDesc"] = v
	}
	request["GroupName"] = d.Get("group_name")
	if v, ok := d.GetOk("iot_instance_id"); ok {
		request["IotInstanceId"] = v
	}
	if v, ok := d.GetOk("super_group_id"); ok {
		request["SuperGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_iot_device_group", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["GroupId"]))

	return resourceAlicloudIotDeviceGroupRead(d, meta)
}
func resourceAlicloudIotDeviceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	iotService := IotService{client}
	object, err := iotService.DescribeIotDeviceGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_iot_device_group iotService.DescribeIotDeviceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	descmap := object["Data"].(map[string]interface{})
	if val, exist := descmap["GroupDesc"]; exist {
		d.Set("group_desc", val)
	}

	d.Set("group_name", object["Data"].(map[string]interface{})["GroupName"])
	return nil
}
func resourceAlicloudIotDeviceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}
	if d.HasChange("group_desc") {
		update = true
		if v, ok := d.GetOk("group_desc"); ok {
			request["GroupDesc"] = v
		}
	}
	if d.HasChange("iot_instance_id") {
		update = true
		if v, ok := d.GetOk("iot_instance_id"); ok {
			request["IotInstanceId"] = v
		}
	}
	if update {
		action := "UpdateDeviceGroup"
		conn, err := client.NewIotClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudIotDeviceGroupRead(d, meta)
}
func resourceAlicloudIotDeviceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDeviceGroup"
	var response map[string]interface{}
	conn, err := client.NewIotClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GroupId": d.Id(),
	}

	if v, ok := d.GetOk("iot_instance_id"); ok {
		request["IotInstanceId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
