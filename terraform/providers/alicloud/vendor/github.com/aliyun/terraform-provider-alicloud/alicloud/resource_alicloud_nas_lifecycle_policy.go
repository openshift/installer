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

func resourceAlicloudNasLifecyclePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudNasLifecyclePolicyCreate,
		Read:   resourceAlicloudNasLifecyclePolicyRead,
		Update: resourceAlicloudNasLifecyclePolicyUpdate,
		Delete: resourceAlicloudNasLifecyclePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lifecycle_policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lifecycle_rule_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT_ATIME_14", "DEFAULT_ATIME_30", "DEFAULT_ATIME_60", "DEFAULT_ATIME_90"}, false),
			},
			"paths": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 10,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"InfrequentAccess"}, false),
			},
		},
	}
}

func resourceAlicloudNasLifecyclePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateLifecyclePolicy"
	request := make(map[string]interface{})
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	request["FileSystemId"] = d.Get("file_system_id")
	request["LifecyclePolicyName"] = d.Get("lifecycle_policy_name")
	request["LifecycleRuleName"] = d.Get("lifecycle_rule_name")
	request["Paths"] = d.Get("paths")
	request["StorageType"] = d.Get("storage_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_lifecycle_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FileSystemId"], ":", request["LifecyclePolicyName"]))

	return resourceAlicloudNasLifecyclePolicyRead(d, meta)
}
func resourceAlicloudNasLifecyclePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasLifecyclePolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_lifecycle_policy nasService.DescribeNasLifecyclePolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("file_system_id", object["FileSystemId"])
	d.Set("lifecycle_policy_name", object["LifecyclePolicyName"])
	d.Set("lifecycle_rule_name", object["LifecycleRuleName"])
	d.Set("paths", object["Paths"])
	d.Set("storage_type", object["StorageType"])
	return nil
}
func resourceAlicloudNasLifecyclePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"FileSystemId":        parts[0],
		"LifecyclePolicyName": parts[1],
	}

	request["LifecycleRuleName"] = d.Get("lifecycle_rule_name")
	if d.HasChange("lifecycle_rule_name") {
		update = true
	}
	if update {
		action := "ModifyLifecyclePolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudNasLifecyclePolicyRead(d, meta)
}
func resourceAlicloudNasLifecyclePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteLifecyclePolicy"
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FileSystemId":        parts[0],
		"LifecyclePolicyName": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound", "InvalidLifecyclePolicy.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
