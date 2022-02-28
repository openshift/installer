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

func resourceAlicloudCloudFirewallInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudFirewallInstanceCreate,
		Read:   resourceAlicloudCloudFirewallInstanceRead,
		Update: resourceAlicloudCloudFirewallInstanceUpdate,
		Delete: resourceAlicloudCloudFirewallInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{6, 12, 24, 36}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"renew_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 12),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
							return false
						}
					}
					return true
				},
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"logistics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cfw_service": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"fw_vpc_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 500),
			},
			"ip_number": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(20, 4000),
			},
			"cfw_log_storage": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1000, 500000),
			},
			"cfw_log": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"band_width": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(10, 15000),
			},
			"instance_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(5, 5000),
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"premium_version", "enterprise_version", "ultimate_version"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"renewal_duration_unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Downgrade", "Upgrade"}, false),
			},
		},
	}
}
func resourceAlicloudCloudFirewallInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request["ClientToken"] = buildClientToken(action)
	request["ProductCode"] = "vipcloudfw"

	request["ProductType"] = "vipcloudfw"
	request["SubscriptionType"] = d.Get("payment_type")

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	} else if d.Get("payment_type").(string) == "Subscription" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v", "period", "payment_type", "Subscription"))
	}

	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}

	if v, ok := d.GetOk("renewal_duration"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_duration", "renewal_status", d.Get("renewal_status")))
	}
	if v, ok := d.GetOk("logistics"); ok {
		request["Logistics"] = v
	}

	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Spec",
		"Value": convertCloudFirewallInstanceVersion(d.Get("spec").(string)),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "IpNumber",
		"Value": d.Get("ip_number"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "BandWidth",
		"Value": d.Get("band_width"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "CfwLog",
		"Value": d.Get("cfw_log"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "CfwLogStorage",
		"Value": d.Get("cfw_log_storage"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "CfwService",
		"Value": d.Get("cfw_service"),
	})
	if v, ok := d.GetOk("fw_vpc_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "FwVpcNumber",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("instance_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceCount",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &runtime)
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
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))
	return resourceAlicloudCloudFirewallInstanceRead(d, meta)
}
func resourceAlicloudCloudFirewallInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	bssOpenApiService := BssOpenApiService{client}
	getQueryInstanceObject, err := bssOpenApiService.QueryAvailableInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("create_time", getQueryInstanceObject["CreateTime"])
	d.Set("renewal_status", getQueryInstanceObject["RenewStatus"])
	d.Set("renewal_duration_unit", convertCloudFirewallInstanceRenewalDurationUnitResponse(getQueryInstanceObject["RenewalDurationUnit"]))
	d.Set("status", getQueryInstanceObject["Status"])
	d.Set("subscription_type", getQueryInstanceObject["SubscriptionType"])
	d.Set("end_time", getQueryInstanceObject["EndTime"])
	return nil
}
func resourceAlicloudCloudFirewallInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)
	update := false
	renewInstancerequest := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	if d.HasChange("renew_period") && !d.IsNewResource() {
		update = true
		renewInstancerequest["RenewPeriod"] = d.Get("renew_period")
	}
	if update {
		action := "RenewInstance"
		renewInstancerequest["ClientToken"] = buildClientToken(action)
		renewInstancerequest["ProductCode"] = "vipcloudfw"
		renewInstancerequest["ProductType"] = "vipcloudfw"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, renewInstancerequest, &runtime)
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
			return nil
		})
		addDebug(action, response, renewInstancerequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	update = false
	modifyInstanceRequest := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	modifyInstanceRequest["ProductType"] = "vipcloudfw"
	modifyInstanceRequest["ProductCode"] = "vipcloudfw"
	modifyInstanceRequest["SubscriptionType"] = d.Get("payment_type")
	parameterMapList := make([]map[string]interface{}, 0)
	if d.HasChange("cfw_service") {
		update = true
	}
	if v, ok := d.GetOk("cfw_service"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwService",
			"Value": v,
		})
	}

	if d.HasChange("fw_vpc_number") {
		update = true
	}
	if v, ok := d.GetOk("fw_vpc_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "FwVpcNumber",
			"Value": v,
		})
	}
	if d.HasChange("ip_number") {
		update = true
	}
	if v, ok := d.GetOk("ip_number"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "IpNumber",
			"Value": v,
		})
	}
	if d.HasChange("cfw_log_storage") {
		update = true
	}
	if v, ok := d.GetOk("cfw_log_storage"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLogStorage",
			"Value": v,
		})
	}
	if d.HasChange("cfw_log") {
		update = true
	}
	if v, ok := d.GetOk("cfw_log"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "CfwLog",
			"Value": v,
		})
	}
	if d.HasChange("band_width") {
		update = true
	}
	if v, ok := d.GetOk("band_width"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BandWidth",
			"Value": v,
		})
	}
	if d.HasChange("spec") {
		update = true
	}
	if v, ok := d.GetOk("spec"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Spec",
			"Value": convertCloudFirewallInstanceVersion(v.(string)),
		})
	}
	if d.HasChange("instance_count") {
		update = true
	}
	if v, ok := d.GetOk("instance_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "InstanceCount",
			"Value": v,
		})
	}
	modifyInstanceRequest["Parameter"] = parameterMapList
	if update {
		if v, ok := d.GetOk("modify_type"); ok {
			modifyInstanceRequest["ModifyType"] = v
		}
		action := "ModifyInstance"
		modifyInstanceRequest["ClientToken"] = buildClientToken(action)
		conn, err := client.NewBssopenapiClient()
		if err != nil {
			return WrapError(err)
		}
		modifyInstanceRequest["ClientToken"] = buildClientToken(action)
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, modifyInstanceRequest, &runtime)
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
			return nil
		})
		addDebug(action, response, modifyInstanceRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("payment_type")
		d.SetPartial("cfw_service")
		d.SetPartial("fw_vpc_number")
		d.SetPartial("ip_number")
		d.SetPartial("cfw_log_storage")
		d.SetPartial("cfw_log")
		d.SetPartial("band_width")
		d.SetPartial("spec")
		d.SetPartial("instance_count")
	}
	d.Partial(false)
	return resourceAlicloudCloudFirewallInstanceRead(d, meta)
}
func resourceAlicloudCloudFirewallInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudCloudFirewallInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertCloudFirewallInstanceVersion(source string) interface{} {
	switch source {
	case "premium_version":
		return 2
	case "enterprise_version":
		return 3
	case "ultimate_version":
		return 4
	}
	return source
}

func convertCloudFirewallInstanceRenewalDurationUnitResponse(source interface{}) interface{} {
	switch source {
	case "M":
		return "Month"
	case "Y":
		return "Year"
	}
	return source
}
