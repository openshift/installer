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

func resourceAlicloudAmqpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAmqpInstanceCreate,
		Read:   resourceAlicloudAmqpInstanceRead,
		Update: resourceAlicloudAmqpInstanceUpdate,
		Delete: resourceAlicloudAmqpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"professional", "vip"}, false),
			},
			"logistics": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_eip_tps": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOkExists("support_eip"); ok && v.(bool) {
						return false
					}
					return true
				},
			},
			"max_tps": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Downgrade", "Upgrade"}, false),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Subscription"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 2, 24, 3, 6}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"queue_capacity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"renewal_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 2, 3, 6}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
							return false
						}
					}
					return true
				},
			},
			"renewal_duration_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
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
				ValidateFunc: validation.StringInSlice([]string{"AutoRenewal", "ManualRenewal", "NotRenewal"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("instance_type"); ok && v.(string) == "vip" {
						return false
					}
					return true
				},
			},
			"support_eip": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceAlicloudAmqpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	parameterMapList := make([]map[string]interface{}, 0)
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "InstanceType",
		"Value": d.Get("instance_type"),
	})
	if v, ok := d.GetOk("logistics"); ok {
		request["Logistics"] = v
	}
	if v, ok := d.GetOk("max_eip_tps"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "MaxEipTps",
			"Value": v,
		})
	} else if v, ok := d.GetOkExists("support_eip"); ok && v.(bool) {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "max_eip_tps", "support_eip", d.Get("support_eip")))
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "MaxTps",
		"Value": d.Get("max_tps"),
	})
	request["SubscriptionType"] = d.Get("payment_type")
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["ProductCode"] = "ons"
	request["ProductType"] = "ons_onsproxy_pre"
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "QueueCapacity",
		"Value": d.Get("queue_capacity"),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Region",
		"Value": client.RegionId,
	})
	if v, ok := d.GetOk("renewal_duration"); ok {
		request["RenewPeriod"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_duration", "renewal_status", d.Get("renewal_status")))
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if v, ok := d.GetOk("storage_size"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "StorageSize",
			"Value": v,
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "SupportEip",
		"Value": convertAmqpInstanceSupportEipRequest(d.Get("support_eip").(bool)),
	})
	request["Parameter"] = parameterMapList
	request["ClientToken"] = buildClientToken("CreateInstance")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_amqp_instance", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Code"]) != "Success" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))
	amqpOpenService := AmqpOpenService{client}
	stateConf := BuildStateConf([]string{}, []string{"SERVING"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, amqpOpenService.AmqpInstanceStateRefreshFunc(d.Id(), []string{"EXPIRED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAmqpInstanceUpdate(d, meta)
}
func resourceAlicloudAmqpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	amqpOpenService := AmqpOpenService{client}
	object, err := amqpOpenService.DescribeAmqpInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_amqp_instance amqpOpenService.DescribeAmqpInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_name", object["InstanceName"])
	d.Set("instance_type", convertAmqpInstanceInstanceTypeResponse(object["InstanceType"]))
	d.Set("status", object["Status"])
	d.Set("support_eip", object["SupportEIP"])
	bssOpenApiService := BssOpenApiService{client}
	queryAvailableInstancesObject, err := bssOpenApiService.QueryAvailableInstances(d.Id(), "ons", "ons_onsproxy_pre")
	if err != nil {
		return WrapError(err)
	}
	d.Set("payment_type", queryAvailableInstancesObject["SubscriptionType"])
	if v, ok := queryAvailableInstancesObject["RenewalDuration"]; ok && fmt.Sprint(v) != "0" {
		d.Set("renewal_duration", formatInt(v))
	}
	d.Set("renewal_duration_unit", convertAmqpInstanceRenewalDurationUnitResponse(object["RenewalDurationUnit"]))
	d.Set("renewal_status", queryAvailableInstancesObject["RenewStatus"])
	return nil
}
func resourceAlicloudAmqpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	}
	if update {
		action := "UpdateInstanceName"
		conn, err := client.NewOnsproxyClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("instance_name")
	}
	update = false
	setRenewalReq := map[string]interface{}{
		"InstanceIDs": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		setRenewalReq["RenewalStatus"] = v
	}
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		setRenewalReq["SubscriptionType"] = d.Get("payment_type")
	}
	setRenewalReq["ProductCode"] = "ons"
	setRenewalReq["ProductType"] = "ons_onsproxy_pre"
	if !d.IsNewResource() && d.HasChange("renewal_duration") {
		update = true
		if v, ok := d.GetOk("renewal_duration"); ok {
			setRenewalReq["RenewalPeriod"] = v
		} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_duration", "renewal_status", d.Get("renewal_status")))
		}
	}
	if d.HasChange("renewal_duration_unit") {
		update = true
		if v, ok := d.GetOk("renewal_duration_unit"); ok {
			setRenewalReq["RenewalPeriodUnit"] = convertAmqpInstanceRenewalDurationUnitRequest(v.(string))
		} else if v, ok := d.GetOk("renewal_status"); ok && v.(string) == "AutoRenewal" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "renewal_duration_unit", "renewal_status", d.Get("renewal_status")))
		}
	}
	if update {
		action := "SetRenewal"
		conn, err := client.NewBssopenapiClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, setRenewalReq, &util.RuntimeOptions{})
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
		addDebug(action, response, setRenewalReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("renewal_status")
		d.SetPartial("payment_type")
		d.SetPartial("renewal_duration")
		d.SetPartial("renewal_duration_unit")
	}
	update = false
	modifyInstanceReq := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	parameterMapList := make([]map[string]interface{}, 0)
	if !d.IsNewResource() && d.HasChange("max_tps") {
		update = true
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "MaxTps",
		"Value": d.Get("max_tps"),
	})

	modifyInstanceReq["SubscriptionType"] = d.Get("payment_type")
	modifyInstanceReq["ProductCode"] = "ons"
	if !d.IsNewResource() && d.HasChange("queue_capacity") {
		update = true
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "QueueCapacity",
		"Value": d.Get("queue_capacity"),
	})
	if !d.IsNewResource() && d.HasChange("support_eip") {
		update = true
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "SupportEip",
		"Value": convertAmqpInstanceSupportEipRequest(d.Get("support_eip").(bool)),
	})
	if !d.IsNewResource() && d.HasChange("max_eip_tps") {
		update = true
	}
	if v, ok := d.GetOk("max_eip_tps"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "MaxEipTps",
			"Value": v,
		})
	} else if v, ok := d.GetOkExists("support_eip"); ok && v.(bool) {
		return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v ", "max_eip_tps", "support_eip", d.Get("support_eip")))
	}
	modifyInstanceReq["ProductType"] = "ons_onsproxy_pre"
	if !d.IsNewResource() && d.HasChange("storage_size") {
		update = true
	}
	if v, ok := d.GetOk("storage_size"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "StorageSize",
			"Value": v,
		})
	}
	modifyInstanceReq["Parameter"] = parameterMapList
	if update {
		if v, ok := d.GetOk("modify_type"); ok {
			modifyInstanceReq["ModifyType"] = v
		}
		action := "ModifyInstance"
		conn, err := client.NewBssopenapiClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ModifyInstance")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, modifyInstanceReq, &runtime)
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
		addDebug(action, response, modifyInstanceReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Code"]) != "Success" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		d.SetPartial("max_tps")
		d.SetPartial("payment_type")
		d.SetPartial("queue_capacity")
		d.SetPartial("support_eip")
		d.SetPartial("max_eip_tps")
		d.SetPartial("storage_size")
	}
	d.Partial(false)
	return resourceAlicloudAmqpInstanceRead(d, meta)
}
func resourceAlicloudAmqpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudAmqpInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
func convertAmqpInstanceSupportEipRequest(source interface{}) interface{} {
	switch source {
	case false:
		return "eip_false"
	case true:
		return "eip_true"
	}
	return ""
}
func convertAmqpInstanceInstanceTypeResponse(source interface{}) interface{} {
	switch source {
	case "PROFESSIONAL":
		return "professional"
	case "VIP":
		return "vip"
	}
	return source
}
func convertAmqpInstanceRenewalDurationUnitResponse(source interface{}) interface{} {
	switch source {
	case "M":
		return "Month"
	case "Y":
		return "Year"
	}
	return source
}
func convertAmqpInstanceRenewalDurationUnitRequest(source interface{}) interface{} {
	switch source {
	case "Month":
		return "M"
	case "Year":
		return "Y"
	}
	return source
}
