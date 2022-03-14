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

func resourceAlicloudAlbSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbSecurityPolicyCreate,
		Read:   resourceAlicloudAlbSecurityPolicyRead,
		Update: resourceAlicloudAlbSecurityPolicyUpdate,
		Delete: resourceAlicloudAlbSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(16 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ciphers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"tls_versions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudAlbSecurityPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSecurityPolicy"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["SecurityPolicyName"] = d.Get("security_policy_name")
	request["TLSVersions"] = d.Get("tls_versions")
	request["Ciphers"] = d.Get("ciphers")
	request["ClientToken"] = buildClientToken("CreateSecurityPolicy")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_security_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SecurityPolicyId"]))

	return resourceAlicloudAlbSecurityPolicyRead(d, meta)
}
func resourceAlicloudAlbSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbSecurityPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_security_policy albService.DescribeAlbSecurityPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ciphers", object["Ciphers"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("security_policy_name", object["SecurityPolicyName"])
	d.Set("status", object["SecurityPolicyStatus"])
	d.Set("tls_versions", object["TLSVersions"])

	tagResp, err := albService.ListTagResources(d.Id(), "securitypolicy")
	d.Set("tags", tagsToMap(tagResp))

	return nil
}
func resourceAlicloudAlbSecurityPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := albService.SetResourceTags(d, "securitypolicy"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["NewResourceGroupId"] = v
	}
	request["ResourceType"] = "securitypolicy"
	if update {
		action := "MoveResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("resource_group_id")
	}
	update = false
	updateSecurityPolicyAttributeReq := map[string]interface{}{
		"SecurityPolicyId": d.Id(),
	}
	if d.HasChange("ciphers") {
		update = true
		updateSecurityPolicyAttributeReq[""] = d.Get("ciphers")
	}
	if d.HasChange("security_policy_name") {
		update = true
		updateSecurityPolicyAttributeReq["SecurityPolicyName"] = d.Get("security_policy_name")
	}
	if d.HasChange("tls_versions") {
		update = true
		updateSecurityPolicyAttributeReq[""] = d.Get("tls_versions")
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			updateSecurityPolicyAttributeReq["DryRun"] = v
		}
		if v, ok := d.GetOk("tls_versions"); ok {
			updateSecurityPolicyAttributeReq["TLSVersions"] = v
		}
		if v, ok := d.GetOk("ciphers"); ok {
			updateSecurityPolicyAttributeReq["Ciphers"] = v
		}
		action := "UpdateSecurityPolicyAttribute"
		request["ClientToken"] = buildClientToken("UpdateSecurityPolicyAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 30*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateSecurityPolicyAttributeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectStatus.SecurityPolicy"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateSecurityPolicyAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbSecurityPolicyStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ciphers")
		d.SetPartial("security_policy_name")
		d.SetPartial("tls_versions")
	}
	d.Partial(false)
	return resourceAlicloudAlbSecurityPolicyRead(d, meta)
}
func resourceAlicloudAlbSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecurityPolicy"
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SecurityPolicyId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteSecurityPolicy")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound.SecurityPolicy"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
