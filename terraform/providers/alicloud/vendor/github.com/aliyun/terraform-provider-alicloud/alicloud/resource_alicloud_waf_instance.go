package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudWafInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudWafInstanceCreate,
		Read:   resourceAlicloudWafInstanceRead,
		Update: resourceAlicloudWafInstanceUpdate,
		Delete: resourceAlicloudWafInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"big_screen": {
				Type:     schema.TypeString,
				Required: true,
			},
			"exclusive_ip_package": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_domain_package": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_storage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"prefessional_service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waf_log": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudWafInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}

	request["ProductCode"] = "waf"
	request["ProductType"] = "waf"
	if v, ok := d.GetOk("renew_period"); ok {
		request["RenewPeriod"] = v
	}

	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	region := client.RegionId
	if v, ok := d.GetOk("region"); ok && v.(string) != "" {
		region = v.(string)
	}
	request["SubscriptionType"] = d.Get("subscription_type")
	request["Parameter"] = []map[string]string{
		{
			"Code":  "BigScreen",
			"Value": d.Get("big_screen").(string),
		},
		{
			"Code":  "ExclusiveIpPackage",
			"Value": d.Get("exclusive_ip_package").(string),
		},
		{
			"Code":  "ExtBandwidth",
			"Value": d.Get("ext_bandwidth").(string),
		},
		{
			"Code":  "ExtDomainPackage",
			"Value": d.Get("ext_domain_package").(string),
		},
		{
			"Code":  "LogStorage",
			"Value": d.Get("log_storage").(string),
		},
		{
			"Code":  "LogTime",
			"Value": d.Get("log_time").(string),
		},
		{
			"Code":  "PackageCode",
			"Value": d.Get("package_code").(string),
		},
		{
			"Code":  "PrefessionalService",
			"Value": d.Get("prefessional_service").(string),
		},
		{
			"Code":  "Region",
			"Value": region,
		},
		{
			"Code":  "WafLog",
			"Value": d.Get("waf_log").(string),
		},
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_waf_instance", action, AlibabaCloudSdkGoERROR)
	}
	if response["Code"].(string) != "Success" {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_waf_instance", action, AlibabaCloudSdkGoERROR)
	}
	response = response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["InstanceId"]))

	return resourceAlicloudWafInstanceUpdate(d, meta)
}
func resourceAlicloudWafInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	waf_openapiService := Waf_openapiService{client}
	object, err := waf_openapiService.DescribeWafInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_waf_instance waf_openapiService.DescribeWafInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", formatInt(object["InstanceInfo"].(map[string]interface{})["Status"]))
	d.Set("subscription_type", object["InstanceInfo"].(map[string]interface{})["SubscriptionType"])
	return nil
}
func resourceAlicloudWafInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("subscription_type") {
		update = true
	}
	request["ProductType"] = "waf"
	request["SubscriptionType"] = d.Get("subscription_type")
	request["ProductCode"] = "waf"
	request["ModifyType"] = d.Get("modify_type")
	request["Parameter"] = []map[string]string{
		{
			"Code":  "BigScreen",
			"Value": d.Get("big_screen").(string),
		},
		{
			"Code":  "ExclusiveIpPackage",
			"Value": d.Get("exclusive_ip_package").(string),
		},
		{
			"Code":  "ExtBandwidth",
			"Value": d.Get("ext_bandwidth").(string),
		},
		{
			"Code":  "ExtDomainPackage",
			"Value": d.Get("ext_domain_package").(string),
		},
		{
			"Code":  "LogStorage",
			"Value": d.Get("log_storage").(string),
		},
		{
			"Code":  "LogTime",
			"Value": d.Get("log_time").(string),
		},
		{
			"Code":  "PackageCode",
			"Value": d.Get("package_code").(string),
		},
		{
			"Code":  "PrefessionalService",
			"Value": d.Get("prefessional_service").(string),
		},
		{
			"Code":  "WafLog",
			"Value": d.Get("waf_log").(string),
		},
	}
	if update {
		action := "ModifyInstance"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"NotApplicable"}) {
					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
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
		if response["Code"].(string) != "Success" {
			return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

	}
	return resourceAlicloudWafInstanceRead(d, meta)
}
func resourceAlicloudWafInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteInstance"
	var response map[string]interface{}
	conn, err := client.NewWafClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"ComboError"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
