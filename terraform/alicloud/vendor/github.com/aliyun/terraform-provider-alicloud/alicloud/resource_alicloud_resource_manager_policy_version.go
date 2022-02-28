package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudResourceManagerPolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerPolicyVersionCreate,
		Read:   resourceAlicloudResourceManagerPolicyVersionRead,
		Update: resourceAlicloudResourceManagerPolicyVersionUpdate,
		Delete: resourceAlicloudResourceManagerPolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"is_default_version": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "Field 'is_default_version' has been deprecated from provider version 1.90.0",
			},
			"policy_document": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
				Removed:  "Field 'create_date' has been removed from provider version 1.100.0.",
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Removed:  "Field 'version_id' has been removed from provider version 1.100.0.",
			},
		},
	}
}

func resourceAlicloudResourceManagerPolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicyVersion"
	request := make(map[string]interface{})
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("is_default_version"); ok {
		request["SetAsDefault"] = v
	}

	request["PolicyDocument"] = d.Get("policy_document")
	request["PolicyName"] = d.Get("policy_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_policy_version", action, AlibabaCloudSdkGoERROR)
	}
	responsePolicyVersion := response["PolicyVersion"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["PolicyName"], ":", responsePolicyVersion["VersionId"]))

	return resourceAlicloudResourceManagerPolicyVersionRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerPolicyVersion(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy_version resourcemanagerService.DescribeResourceManagerPolicyVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("policy_name", parts[0])
	d.Set("is_default_version", object["IsDefaultVersion"])
	d.Set("policy_document", object["PolicyDocument"])
	return nil
}
func resourceAlicloudResourceManagerPolicyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("is_default_version") {
		object, err := resourcemanagerService.DescribeResourceManagerPolicyVersion(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := strconv.FormatBool(d.Get("is_default_version").(bool))
		if strconv.FormatBool(object["IsDefaultVersion"].(bool)) != target {
			if target == "true" {
				request := map[string]interface{}{
					"PolicyName": parts[0],
					"VersionId":  parts[1],
				}
				action := "SetDefaultPolicyVersion"
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
		}
	}
	return resourceAlicloudResourceManagerPolicyVersionRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeletePolicyVersion"
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PolicyName": parts[0],
		"VersionId":  parts[1],
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
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
