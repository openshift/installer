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

func resourceAlicloudCloudSsoUserAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoUserAttachmentCreate,
		Read:   resourceAlicloudCloudSsoUserAttachmentRead,
		Delete: resourceAlicloudCloudSsoUserAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCloudSsoUserAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddUserToGroup"
	request := make(map[string]interface{})
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}

	request["DirectoryId"] = d.Get("directory_id")
	request["GroupId"] = d.Get("group_id")
	request["UserId"] = d.Get("user_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_user_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DirectoryId"], ":", request["GroupId"], ":", request["UserId"]))

	return resourceAlicloudCloudSsoUserAttachmentRead(d, meta)
}
func resourceAlicloudCloudSsoUserAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	_, err := cloudssoService.DescribeCloudSsoUserAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_user_attachment cloudssoService.DescribeCloudSsoUserAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("directory_id", parts[0])
	d.Set("group_id", parts[1])
	d.Set("user_id", parts[2])
	return nil
}
func resourceAlicloudCloudSsoUserAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "RemoveUserFromGroup"
	var response map[string]interface{}
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId": parts[0],
		"GroupId":     parts[1],
		"UserId":      parts[2],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
