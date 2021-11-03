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

func resourceAlicloudResourceManagerResourceShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerResourceShareCreate,
		Read:   resourceAlicloudResourceManagerResourceShareRead,
		Update: resourceAlicloudResourceManagerResourceShareUpdate,
		Delete: resourceAlicloudResourceManagerResourceShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"resource_share_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_share_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerResourceShareCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateResourceShare"
	request := make(map[string]interface{})
	conn, err := client.NewRessharingClient()
	if err != nil {
		return WrapError(err)
	}
	request["ResourceShareName"] = d.Get("resource_share_name")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_resource_share", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	response = response["ResourceShare"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["ResourceShareId"]))

	return resourceAlicloudResourceManagerResourceShareRead(d, meta)
}
func resourceAlicloudResourceManagerResourceShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcesharingService := ResourcesharingService{client}
	object, err := resourcesharingService.DescribeResourceManagerResourceShare(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_resource_share resourcesharingService.DescribeResourceManagerResourceShare Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("resource_share_name", object["ResourceShareName"])
	d.Set("resource_share_owner", object["ResourceShareOwner"])
	d.Set("status", object["ResourceShareStatus"])
	return nil
}
func resourceAlicloudResourceManagerResourceShareUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ResourceShareId": d.Id(),
	}
	if d.HasChange("resource_share_name") {
		update = true
	}
	request["ResourceShareName"] = d.Get("resource_share_name")
	if update {
		action := "UpdateResourceShare"
		conn, err := client.NewRessharingClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudResourceManagerResourceShareRead(d, meta)
}
func resourceAlicloudResourceManagerResourceShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcesharingService := ResourcesharingService{client}
	action := "DeleteResourceShare"
	var response map[string]interface{}
	conn, err := client.NewRessharingClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ResourceShareId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{"Deleted"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, resourcesharingService.ResourceManagerResourceShareStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
