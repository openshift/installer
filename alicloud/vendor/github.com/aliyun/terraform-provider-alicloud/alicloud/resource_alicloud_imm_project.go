package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudImmProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudImmProjectCreate,
		Read:   resourceAlicloudImmProjectRead,
		Update: resourceAlicloudImmProjectUpdate,
		Delete: resourceAlicloudImmProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-z][a-zA-Z0-9\\-]{1,49}$`), "The length of the project name is limited to 1-50 characters, can only contain English letters, numbers and dashes (-), must start with an English letter, and cannot be ALL (reserved words)"),
			},
			"service_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudImmProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutProject"
	request := make(map[string]interface{})
	conn, err := client.NewImmClient()
	if err != nil {
		return WrapError(err)
	}
	request["Project"] = d.Get("project")
	if v, ok := d.GetOk("service_role"); ok {
		request["ServiceRole"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_imm_project", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Project"]))

	return resourceAlicloudImmProjectUpdate(d, meta)
}
func resourceAlicloudImmProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	immService := ImmService{client}
	object, err := immService.DescribeImmProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_imm_project immService.DescribeImmProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("project", d.Id())
	if v, ok := object["CU"]; ok && fmt.Sprint(v) != "0" {
		d.Set("compute_unit", formatInt(v))
	}
	d.Set("service_role", object["ServiceRole"])
	return nil
}
func resourceAlicloudImmProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Project": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("service_role") {
		update = true
		if v, ok := d.GetOk("service_role"); ok {
			request["NewServiceRole"] = v
		}
	}
	if update {
		action := "UpdateProject"
		conn, err := client.NewImmClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudImmProjectRead(d, meta)
}
func resourceAlicloudImmProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteProject"
	var response map[string]interface{}
	conn, err := client.NewImmClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Project": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-06"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
