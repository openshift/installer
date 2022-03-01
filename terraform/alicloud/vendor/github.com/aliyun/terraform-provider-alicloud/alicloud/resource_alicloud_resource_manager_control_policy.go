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

func resourceAlicloudResourceManagerControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerControlPolicyCreate,
		Read:   resourceAlicloudResourceManagerControlPolicyRead,
		Update: resourceAlicloudResourceManagerControlPolicyUpdate,
		Delete: resourceAlicloudResourceManagerControlPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"control_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"effect_scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"RAM"}, false),
			},
			"policy_document": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
		},
	}
}

func resourceAlicloudResourceManagerControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateControlPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	request["PolicyName"] = d.Get("control_policy_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	request["EffectScope"] = d.Get("effect_scope")
	request["PolicyDocument"] = d.Get("policy_document")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateControlPolicy")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_control_policy", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	responseControlPolicy := response["ControlPolicy"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseControlPolicy["PolicyId"]))

	return resourceAlicloudResourceManagerControlPolicyRead(d, meta)
}
func resourceAlicloudResourceManagerControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerControlPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_control_policy resourcemanagerService.DescribeResourceManagerControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("control_policy_name", object["PolicyName"])
	d.Set("description", object["Description"])
	d.Set("effect_scope", object["EffectScope"])
	d.Set("policy_document", object["PolicyDocument"])
	return nil
}
func resourceAlicloudResourceManagerControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PolicyId": d.Id(),
	}
	if d.HasChange("control_policy_name") {
		update = true
		request["NewPolicyName"] = d.Get("control_policy_name")
	}
	if d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}
	if d.HasChange("policy_document") {
		update = true
		request["NewPolicyDocument"] = d.Get("policy_document")
	}
	if update {
		action := "UpdateControlPolicy"
		conn, err := client.NewResourcemanagerClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudResourceManagerControlPolicyRead(d, meta)
}
func resourceAlicloudResourceManagerControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteControlPolicy"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PolicyId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"EntityNotExists.ControlPolicy"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
