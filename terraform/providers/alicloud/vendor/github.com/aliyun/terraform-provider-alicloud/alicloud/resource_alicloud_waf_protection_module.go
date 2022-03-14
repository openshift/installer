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

func resourceAlicloudWafProtectionModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafProtectionModuleCreate,
		Read:   resourceAlicloudWafProtectionModuleRead,
		Update: resourceAlicloudWafProtectionModuleUpdate,
		Delete: resourceAlicloudWafProtectionModuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"defense_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ac_cc", "antifraud", "dld", "normalized", "waf"}, false),
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func resourceAlicloudWafProtectionModuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyProtectionModuleMode"
	request := make(map[string]interface{})
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	request["DefenseType"] = d.Get("defense_type")
	request["Domain"] = d.Get("domain")
	request["InstanceId"] = d.Get("instance_id")
	request["Mode"] = d.Get("mode")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_waf_protection_module", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"], ":", request["DefenseType"]))

	return resourceAlicloudWafProtectionModuleUpdate(d, meta)
}
func resourceAlicloudWafProtectionModuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := Waf_openapiService{client}
	object, err := wafOpenapiService.DescribeWafProtectionModule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_waf_protection_module wafOpenapiService.DescribeWafProtectionModule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("defense_type", parts[2])
	d.Set("domain", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("mode", formatInt(object["Mode"]))
	describeProtectionModuleStatusObject, err := wafOpenapiService.DescribeProtectionModuleStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("status", formatInt(describeProtectionModuleStatusObject["ModuleStatus"]))
	return nil
}
func resourceAlicloudWafProtectionModuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	request := map[string]interface{}{
		"DefenseType": parts[2],
		"Domain":      parts[1],
		"InstanceId":  parts[0],
	}

	update := false
	if !d.IsNewResource() && d.HasChange("mode") {
		update = true
		if v, ok := d.GetOkExists("mode"); ok {
			request["Mode"] = v
		}
	}

	if update {
		action := "ModifyProtectionModuleMode"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("mode")
	}

	update = false
	if d.IsNewResource() || d.HasChange("status") {
		update = true
		if v, ok := d.GetOkExists("status"); ok {
			request["ModuleStatus"] = v
		}
	}

	if update {
		action := "ModifyProtectionModuleStatus"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("status")
	}

	d.Partial(false)
	return resourceAlicloudWafProtectionModuleRead(d, meta)
}
func resourceAlicloudWafProtectionModuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudWafProtectionModule. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
