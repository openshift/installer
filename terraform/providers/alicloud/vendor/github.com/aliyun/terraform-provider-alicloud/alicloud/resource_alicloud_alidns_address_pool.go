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

func resourceAlicloudAlidnsAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsAddressPoolCreate,
		Read:   resourceAlicloudAlidnsAddressPoolRead,
		Update: resourceAlicloudAlidnsAddressPoolUpdate,
		Delete: resourceAlicloudAlidnsAddressPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 20,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"attribute_info": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.ValidateJsonString,
						},
						"lba_weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"SMART", "ONLINE", "OFFLINE"}, false),
						},
						"remark": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"address_pool_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lba_strategy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_RR", "RATIO"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6", "DOMAIN"}, false),
			},
		},
	}
}

func resourceAlicloudAlidnsAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddDnsGtmAddressPool"
	request := make(map[string]interface{})
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}

	request["Name"] = d.Get("address_pool_name")
	request["InstanceId"] = d.Get("instance_id")
	request["LbaStrategy"] = d.Get("lba_strategy")
	if v, ok := d.GetOk("address"); ok {
		addressMaps := make([]map[string]interface{}, 0)
		for _, address := range v.(*schema.Set).List() {
			addressArg := address.(map[string]interface{})
			addressMap := map[string]interface{}{}
			addressMap["Addr"] = addressArg["address"]
			addressMap["AttributeInfo"] = addressArg["attribute_info"]
			if v, ok := addressArg["remark"]; ok {
				addressMap["Remark"] = v
			}
			if v, ok := addressArg["lba_weight"]; ok && fmt.Sprint(v) != "0" {
				addressMap["LbaWeight"] = v
			}
			addressMap["Mode"] = addressArg["mode"]
			addressMaps = append(addressMaps, addressMap)
		}
		request["Addr"] = addressMaps
	}
	request["Type"] = d.Get("type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_address_pool", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AddrPoolId"]))

	return resourceAlicloudAlidnsAddressPoolRead(d, meta)
}
func resourceAlicloudAlidnsAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsAddressPool(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_address_pool alidnsService.DescribeAlidnsAddressPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("address_pool_name", object["Name"])
	d.Set("lba_strategy", object["LbaStrategy"])
	d.Set("type", object["Type"])

	return nil
}
func resourceAlicloudAlidnsAddressPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"AddrPoolId": d.Id(),
	}
	if d.HasChange("lba_strategy") {
		update = true
	}
	request["LbaStrategy"] = d.Get("lba_strategy")
	if d.HasChange("address_pool_name") {
		update = true
		request["Name"] = d.Get("address_pool_name")
	}
	if d.HasChange("address") {
		update = true
	}
	if v, ok := d.GetOk("address"); ok {
		addressMaps := make([]map[string]interface{}, 0)
		for _, address := range v.(*schema.Set).List() {
			addressArg := address.(map[string]interface{})
			addressMap := map[string]interface{}{}
			addressMap["Addr"] = addressArg["address"]
			addressMap["AttributeInfo"] = addressArg["attribute_info"]
			if v, ok := addressArg["remark"]; ok {
				addressMap["Remark"] = v
			}
			if v, ok := addressArg["lba_weight"]; ok && fmt.Sprint(v) != "0" {
				addressMap["LbaWeight"] = v
			}
			addressMap["Mode"] = addressArg["mode"]
			addressMaps = append(addressMaps, addressMap)
		}
		request["Addr"] = addressMaps
	}
	if update {
		action := "UpdateDnsGtmAddressPool"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudAlidnsAddressPoolRead(d, meta)
}
func resourceAlicloudAlidnsAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDnsGtmAddressPool"
	var response map[string]interface{}
	conn, err := client.NewAlidnsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AddrPoolId": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-01-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
