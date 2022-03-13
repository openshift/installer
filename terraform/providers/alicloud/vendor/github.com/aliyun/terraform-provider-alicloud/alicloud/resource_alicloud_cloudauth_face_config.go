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

func resourceAlicloudCloudauthFaceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudauthFaceConfigCreate,
		Read:   resourceAlicloudCloudauthFaceConfigRead,
		Update: resourceAlicloudCloudauthFaceConfigUpdate,
		Delete: resourceAlicloudCloudauthFaceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"biz_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"biz_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCloudauthFaceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateFaceConfig"
	request := make(map[string]interface{})
	conn, err := client.NewCloudauthClient()
	if err != nil {
		return WrapError(err)
	}
	request["BizName"] = d.Get("biz_name")
	request["BizType"] = d.Get("biz_type")
	request["SourceIp"] = buildClientToken(action)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloudauth_face_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["BizType"]))

	return resourceAlicloudCloudauthFaceConfigRead(d, meta)
}
func resourceAlicloudCloudauthFaceConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudauthService := CloudauthService{client}
	object, err := cloudauthService.DescribeCloudauthFaceConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloudauth_face_config cloudauthService.DescribeCloudauthFaceConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("biz_name", object["BizName"])
	d.Set("biz_type", object["BizType"])
	d.Set("gmt_modified", object["GmtUpdated"])
	return nil
}
func resourceAlicloudCloudauthFaceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"BizType": d.Id(),
	}
	if d.HasChange("biz_name") {
		update = true
	}
	request["BizName"] = d.Get("biz_name")
	if update {
		action := "UpdateFaceConfig"
		request["SourceIp"] = buildClientToken(action)
		conn, err := client.NewCloudauthClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudCloudauthFaceConfigRead(d, meta)
}
func resourceAlicloudCloudauthFaceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudCloudauthFaceConfig. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
