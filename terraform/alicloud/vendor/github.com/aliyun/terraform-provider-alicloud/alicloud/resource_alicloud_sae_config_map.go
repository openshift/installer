package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSaeConfigMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeConfigMapCreate,
		Read:   resourceAlicloudSaeConfigMapRead,
		Update: resourceAlicloudSaeConfigMapUpdate,
		Delete: resourceAlicloudSaeConfigMapDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSaeConfigMapCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/sam/configmap/configMap"
	request := make(map[string]*string)
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request["Data"] = StringPointer(d.Get("data").(string))
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	request["Name"], request["NamespaceId"] = StringPointer(d.Get("name").(string)), StringPointer(d.Get("namespace_id").(string))
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
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_config_map", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["ConfigMapId"]))

	return resourceAlicloudSaeConfigMapRead(d, meta)
}
func resourceAlicloudSaeConfigMapRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeConfigMap(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_config_map saeService.DescribeSaeConfigMap Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	resp, err := json.Marshal(object["Data"])
	if err != nil {
		return WrapError(err)
	}
	d.Set("data", string(resp))
	d.Set("description", object["Description"])
	d.Set("name", object["Name"])
	d.Set("namespace_id", object["NamespaceId"])
	return nil
}
func resourceAlicloudSaeConfigMapUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]*string{
		"ConfigMapId": StringPointer(d.Id()),
	}
	if d.HasChange("data") {
		update = true
	}
	request["Data"] = StringPointer(d.Get("data").(string))
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	if update {
		action := "/pop/v1/sam/configmap/configMap"
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
			if IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
				return WrapErrorf(Error(GetNotFoundMessage("SAE:ConfigMap", d.Id())), NotFoundMsg, ProviderERROR)
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
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}
	return resourceAlicloudSaeConfigMapRead(d, meta)
}
func resourceAlicloudSaeConfigMapDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/configmap/configMap"
	var response map[string]interface{}
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"ConfigMapId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "DELETE "+action, response))
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"NotFound.ConfigMap"}) {
		return nil
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
