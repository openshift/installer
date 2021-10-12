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

func resourceAlicloudEipanycastAnycastEipAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEipanycastAnycastEipAddressCreate,
		Read:   resourceAlicloudEipanycastAnycastEipAddressRead,
		Update: resourceAlicloudEipanycastAnycastEipAddressUpdate,
		Delete: resourceAlicloudEipanycastAnycastEipAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"anycast_eip_address_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internet_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "PayByTraffic",
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
				Default:      "PayAsYouGo",
			},
			"service_location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEipanycastAnycastEipAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipanycastService := EipanycastService{client}
	var response map[string]interface{}
	action := "AllocateAnycastEipAddress"
	request := make(map[string]interface{})
	conn, err := client.NewEipanycastClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("anycast_eip_address_name"); ok {
		request["Name"] = v
	}

	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request["InternetChargeType"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["InstanceChargeType"] = convertEipanycastAnycastEipAddressPaymentTypeRequest(v.(string))
	}

	request["ServiceLocation"] = d.Get("service_location")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("AllocateAnycastEipAddress")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eipanycast_anycast_eip_address", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["AnycastId"]))
	stateConf := BuildStateConf([]string{}, []string{"Allocated"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, eipanycastService.EipanycastAnycastEipAddressStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEipanycastAnycastEipAddressRead(d, meta)
}
func resourceAlicloudEipanycastAnycastEipAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eipanycastService := EipanycastService{client}
	object, err := eipanycastService.DescribeEipanycastAnycastEipAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eipanycast_anycast_eip_address eipanycastService.DescribeEipanycastAnycastEipAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("anycast_eip_address_name", object["Name"])
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("description", object["Description"])
	d.Set("internet_charge_type", object["InternetChargeType"])
	d.Set("payment_type", convertEipanycastAnycastEipAddressPaymentTypeResponse(object["InstanceChargeType"].(string)))
	d.Set("service_location", object["ServiceLocation"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudEipanycastAnycastEipAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("bandwidth") {
		request := map[string]interface{}{
			"AnycastId": d.Id(),
		}
		request["Bandwidth"] = d.Get("bandwidth")
		action := "ModifyAnycastEipAddressSpec"
		conn, err := client.NewEipanycastClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("bandwidth")
	}
	update := false
	request := map[string]interface{}{
		"AnycastId": d.Id(),
	}
	if d.HasChange("anycast_eip_address_name") {
		update = true
		request["Name"] = d.Get("anycast_eip_address_name")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if update {
		action := "ModifyAnycastEipAddressAttribute"
		conn, err := client.NewEipanycastClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("anycast_eip_address_name")
		d.SetPartial("description")
	}
	d.Partial(false)
	return resourceAlicloudEipanycastAnycastEipAddressRead(d, meta)
}
func resourceAlicloudEipanycastAnycastEipAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseAnycastEipAddress"
	var response map[string]interface{}
	conn, err := client.NewEipanycastClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AnycastId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
func convertEipanycastAnycastEipAddressPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	}
	return source
}

func convertEipanycastAnycastEipAddressPaymentTypeResponse(source string) string {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	}
	return source
}
