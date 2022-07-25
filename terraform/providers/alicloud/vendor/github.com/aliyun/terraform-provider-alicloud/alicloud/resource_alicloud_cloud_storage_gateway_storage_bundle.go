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

func resourceAlicloudCloudStorageGatewayStorageBundle() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudStorageGatewayStorageBundleCreate,
		Read:   resourceAlicloudCloudStorageGatewayStorageBundleRead,
		Update: resourceAlicloudCloudStorageGatewayStorageBundleUpdate,
		Delete: resourceAlicloudCloudStorageGatewayStorageBundleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_bundle_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudCloudStorageGatewayStorageBundleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateStorageBundle"
	request := make(map[string]interface{})
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	request["BackendBucketRegionId"] = client.RegionId
	request["Name"] = d.Get("storage_bundle_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_storage_bundle", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["StorageBundleId"]))

	return resourceAlicloudCloudStorageGatewayStorageBundleRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayStorageBundleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	object, err := sgwService.DescribeCloudStorageGatewayStorageBundle(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_storage_bundle sgwService.DescribeCloudStorageGatewayStorageBundle Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("storage_bundle_name", object["Name"])
	return nil
}
func resourceAlicloudCloudStorageGatewayStorageBundleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"StorageBundleId": d.Id(),
	}
	if d.HasChange("storage_bundle_name") {
		update = true
	}
	request["Name"] = d.Get("storage_bundle_name")
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if update {
		action := "ModifyStorageBundle"
		conn, err := client.NewHcsSgwClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCloudStorageGatewayStorageBundleRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayStorageBundleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteStorageBundle"
	var response map[string]interface{}
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"StorageBundleId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"StorageBundleNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
