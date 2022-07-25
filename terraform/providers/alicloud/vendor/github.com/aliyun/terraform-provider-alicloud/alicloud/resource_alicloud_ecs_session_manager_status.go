package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcsSessionManagerStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsSessionManagerStatusCreate,
		Read:   resourceAlicloudEcsSessionManagerStatusRead,
		Update: resourceAlicloudEcsSessionManagerStatusUpdate,
		Delete: resourceAlicloudEcsSessionManagerStatusDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"session_manager_status_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"sessionManagerStatus"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
		},
	}
}

func resourceAlicloudEcsSessionManagerStatusCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyUserBusinessBehavior"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["statusKey"] = d.Get("session_manager_status_name")
	request["statusValue"] = convertEcsSessionManagerStatusStatusRequest(d.Get("status").(string))
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_session_manager_status", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["statusKey"]))
	return resourceAlicloudEcsSessionManagerStatusRead(d, meta)
}
func resourceAlicloudEcsSessionManagerStatusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsSessionManagerStatus(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_session_manager_status ecsService.DescribeEcsSessionManagerStatus Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("session_manager_status_name", d.Id())
	d.Set("status", convertEcsSessionManagerStatusStatusResponse(object["StatusValue"]))
	return nil
}
func resourceAlicloudEcsSessionManagerStatusUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := map[string]interface{}{
		"statusKey": d.Id(),
	}
	if d.HasChange("status") {
		request["statusValue"] = convertEcsSessionManagerStatusStatusRequest(d.Get("status").(string))
	}
	action := "ModifyUserBusinessBehavior"
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudEcsSessionManagerStatusRead(d, meta)
}
func resourceAlicloudEcsSessionManagerStatusDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudEcsSessionManagerStatus. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
func convertEcsSessionManagerStatusStatusRequest(source interface{}) interface{} {
	switch source {
	case "Disabled":
		return "disabled"
	case "Enabled":
		return "enabled"
	}
	return source
}
func convertEcsSessionManagerStatusStatusResponse(source interface{}) interface{} {
	switch source {
	case "disabled":
		return "Disabled"
	case "enabled":
		return "Enabled"
	}
	return source
}
