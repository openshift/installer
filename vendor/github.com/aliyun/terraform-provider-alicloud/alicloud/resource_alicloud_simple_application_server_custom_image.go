package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudSimpleApplicationServerCustomImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSimpleApplicationServerCustomImageCreate,
		Read:   resourceAlicloudSimpleApplicationServerCustomImageRead,
		Update: resourceAlicloudSimpleApplicationServerCustomImageUpdate,
		Delete: resourceAlicloudSimpleApplicationServerCustomImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"custom_image_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9a-zA-Z\u4E00-\u9FA5][\u4E00-\u9FA5A-Za-z0-9:_-]{2,128}$"), "The name must be `2` to `128` characters in length. It must start with a letter or a number. It can contain letters, digits, colons (:), underscores (_) and hyphens (-)."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Share", "UnShare"}, false),
			},
			"system_snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudSimpleApplicationServerCustomImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateCustomImage"
	request := make(map[string]interface{})
	conn, err := client.NewSwasClient()
	if err != nil {
		return WrapError(err)
	}
	request["ImageName"] = d.Get("custom_image_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	request["SystemSnapshotId"] = d.Get("system_snapshot_id")
	request["ClientToken"] = buildClientToken("CreateCustomImage")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_simple_application_server_custom_image", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ImageId"]))

	return resourceAlicloudSimpleApplicationServerCustomImageUpdate(d, meta)
}
func resourceAlicloudSimpleApplicationServerCustomImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	swasOpenService := SwasOpenService{client}
	object, err := swasOpenService.DescribeSimpleApplicationServerCustomImage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_simple_application_server_custom_image swasOpenService.DescribeSimpleApplicationServerCustomImage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("custom_image_name", object["ImageName"])
	d.Set("description", object["Description"])
	return nil
}
func resourceAlicloudSimpleApplicationServerCustomImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ImageId": d.Id(),
	}
	if d.HasChange("status") {
		update = true
		if v, ok := d.GetOk("status"); ok {
			request["Operation"] = v
		}
	}
	if update {
		action := "ModifyImageShareStatus"
		conn, err := client.NewSwasClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyImageShareStatus")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-01"), StringPointer("AK"), nil, request, &runtime)
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
	}
	return resourceAlicloudSimpleApplicationServerCustomImageRead(d, meta)
}
func resourceAlicloudSimpleApplicationServerCustomImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if v := d.Get("status"); v.(string) == "Share" {
		var response map[string]interface{}
		request := map[string]interface{}{
			"ImageId":   d.Id(),
			"Operation": "UnShare",
		}
		action := "ModifyImageShareStatus"
		conn, err := client.NewSwasClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyImageShareStatus")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-01"), StringPointer("AK"), nil, request, &runtime)
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
	}
	action := "DeleteCustomImage"
	var response map[string]interface{}
	conn, err := client.NewSwasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ImageId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeleteCustomImage")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-01"), StringPointer("AK"), nil, request, &runtime)
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
