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

func resourceAlicloudSaeNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeNamespaceCreate,
		Read:   resourceAlicloudSaeNamespaceRead,
		Update: resourceAlicloudSaeNamespaceUpdate,
		Delete: resourceAlicloudSaeNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"namespace_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudSaeNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/paas/namespace"
	request := make(map[string]*string)
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("namespace_description"); ok {
		request["NamespaceDescription"] = StringPointer(v.(string))
	}
	request["NamespaceId"], request["NamespaceName"] = StringPointer(d.Get("namespace_id").(string)), StringPointer(d.Get("namespace_name").(string))
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_namespace", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	addDebug(action, response, request)
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	d.SetId(*request["NamespaceId"])

	return resourceAlicloudSaeNamespaceRead(d, meta)
}
func resourceAlicloudSaeNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeNamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_namespace saeService.DescribeSaeNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("namespace_id", d.Id())
	d.Set("namespace_description", object["NamespaceDescription"])
	d.Set("namespace_name", object["NamespaceName"])
	return nil
}
func resourceAlicloudSaeNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"NamespaceId": StringPointer(d.Id()),
	}
	if d.HasChange("namespace_name") {
		update = true
	}
	request["NamespaceName"] = StringPointer(d.Get("namespace_name").(string))
	if d.HasChange("namespace_description") {
		update = true
	}
	request["NamespaceDescription"] = StringPointer(d.Get("namespace_description").(string))
	if update {
		action := "/pop/v1/paas/namespace"
		conn, err := client.NewServerlessClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) {
				return WrapErrorf(Error(GetNotFoundMessage("SAE:Namespace", d.Id())), NotFoundMsg, ProviderERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
		}
		if respBody, isExist := response["body"]; isExist {
			response = respBody.(map[string]interface{})
		} else {
			return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
		}
		addDebug(action, response, request)
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", "Put "+action, response))
		}
	}
	return resourceAlicloudSaeNamespaceRead(d, meta)
}
func resourceAlicloudSaeNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/paas/namespace"
	var response map[string]interface{}
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"NamespaceId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 1*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidNamespaceId.NotFound"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "DELETE "+action, response))
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"InvalidNamespaceId.NotFound"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "AlicloudSaeNamespaceDelete", response))
	}
	return nil
}
