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

func resourceAlicloudEcdDesktop() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdDesktopCreate,
		Read:   resourceAlicloudEcdDesktopRead,
		Update: resourceAlicloudEcdDesktopUpdate,
		Delete: resourceAlicloudEcdDesktopDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"bundle_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desktop_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktop_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_user_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_disk_size_gib": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("status"); ok && v.(string) == "Stopped" {
						return false
					}
					return true
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped", "Pending", "Stopping", "Starting", "Expired", "Deleted"}, false),
			},
			"stopped_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"StopCharging", "KeepCharging"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("status"); ok && v.(string) == "Stopped" {
						return false
					}
					return true
				},
			},
			"tags": tagsSchema(),
			"user_assign_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL", "PER_USER"}, false),
			},
			"user_disk_size_gib": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("status"); ok && v.(string) == "Stopped" {
						return false
					}
					return true
				},
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEcdDesktopCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	var response map[string]interface{}
	action := "CreateDesktops"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("amount"); ok {
		request["Amount"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	request["BundleId"] = d.Get("bundle_id")
	if v, ok := d.GetOk("desktop_name"); ok {
		request["DesktopName"] = v
	}
	request["OfficeSiteId"] = d.Get("office_site_id")

	if m, ok := d.GetOk("end_user_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("EndUserId.%d", k+1)] = v.(string)
		}
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEcdDesktopPaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	request["PolicyGroupId"] = d.Get("policy_group_id")

	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("user_assign_mode"); ok {
		request["UserAssignMode"] = v
	}
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_desktop", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DesktopId"].([]interface{})[0]))

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdDesktopStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEcdDesktopUpdate(d, meta)
}
func resourceAlicloudEcdDesktopRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdDesktop(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_desktop ecdService.DescribeEcdDesktop Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("desktop_name", object["DesktopName"])
	d.Set("desktop_type", object["DesktopType"])
	d.Set("office_site_id", object["OfficeSiteId"])
	d.Set("end_user_ids", object["EndUserIds"])
	d.Set("payment_type", convertEcdDesktopPaymentTypeResponse(object["ChargeType"]))
	d.Set("policy_group_id", object["PolicyGroupId"])
	d.Set("status", object["DesktopStatus"])
	d.Set("tags", tagsToMap(object["Tags"]))
	return nil
}
func resourceAlicloudEcdDesktopUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecdService.SetResourceTags(d, "ALIYUN::GWS::INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	request := map[string]interface{}{
		"DesktopId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("desktop_name") {
		if v, ok := d.GetOk("desktop_name"); ok {
			request["NewDesktopName"] = v
		}
		request["RegionId"] = client.RegionId
		action := "ModifyDesktopName"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("desktop_name")
	}

	if !d.IsNewResource() && d.HasChange("policy_group_id") {
		request := map[string]interface{}{
			"DesktopId.1": d.Id(),
		}
		request["PolicyGroupId"] = d.Get("policy_group_id")
		request["RegionId"] = client.RegionId

		action := "ModifyDesktopsPolicyGroup"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("policy_group_id")
	}

	if !d.IsNewResource() && d.HasChange("end_user_ids") {
		request := map[string]interface{}{
			"DesktopId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		if m, ok := d.GetOk("end_user_ids"); ok {
			for k, v := range m.([]interface{}) {
				request[fmt.Sprintf("EndUserId.%d", k+1)] = v.(string)
			}
		}
		action := "ModifyEntitlement"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("end_user_ids")
	}

	if d.HasChange("status") {
		object, err := ecdService.DescribeEcdDesktop(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DesktopStatus"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"DesktopId.1": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "StartDesktops"
				conn, err := client.NewGwsecdClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			if target == "Stopped" {
				request := map[string]interface{}{
					"DesktopId.1": d.Id(),
				}
				request["RegionId"] = client.RegionId
				if v, ok := d.GetOk("stopped_mode"); ok {
					request["StoppedMode"] = v
				}
				action := "StopDesktops"
				conn, err := client.NewGwsecdClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			stateConf := BuildStateConf([]string{"Pending", "Starting", "Stopping"}, []string{"Running", "Stopped"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecdService.EcdDesktopStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("status")
		}
	}
	if d.HasChange("desktop_type") || d.HasChange("root_disk_size_gib") || d.HasChange("user_disk_size_gib") {
		request = map[string]interface{}{
			"DesktopId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		if v, ok := d.GetOk("desktop_type"); ok {
			request["DesktopType"] = v
		}
		if v, ok := d.GetOkExists("auto_pay"); ok {
			request["AutoPay"] = v
		}
		if v, ok := d.GetOk("root_disk_size_gib"); ok {
			request["RootDiskSizeGib"] = v
		}
		if v, ok := d.GetOk("user_disk_size_gib"); ok {
			request["UserDiskSizeGib"] = v
		}
		action := "ModifyDesktopSpec"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{request["DesktopType"].(string)}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecdService.EcdDesktopDesktopTypeRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("desktop_type")
	}

	if !d.IsNewResource() && d.HasChange("payment_type") {
		object, err := ecdService.DescribeEcdDesktop(d.Id())
		if err != nil {
			return WrapError(err)
		}
		if convertEcdDesktopPaymentTypeResponse(object["ChargeType"]) != "Subscription" {
			if object["DesktopStatus"].(string) == "Stopped" || object["DesktopStatus"].(string) == "Running" {
				request := map[string]interface{}{
					"DesktopId.1": d.Id(),
				}
				request["RegionId"] = client.RegionId
				if v, ok := d.GetOk("payment_type"); ok {
					request["ChargeType"] = convertEcdDesktopPaymentTypeRequest(v.(string))
				}
				if v, ok := d.GetOk("period"); ok {
					request["Period"] = v
				}
				if v, ok := d.GetOk("period_unit"); ok {
					request["PeriodUnit"] = v
				}
				if v, ok := d.GetOkExists("auto_pay"); ok {
					request["AutoPay"] = v
				}
				action := "ModifyDesktopChargeType"
				conn, err := client.NewGwsecdClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"Subscription"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecdService.EcdDesktopChargeTypeFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
		d.SetPartial("payment_type")
	}
	d.Partial(false)
	return resourceAlicloudEcdDesktopRead(d, meta)
}
func resourceAlicloudEcdDesktopDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	action := "DeleteDesktops"
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DesktopId.1": d.Id(),
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidChargeType.Unsupported"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecdService.EcdDesktopStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertEcdDesktopPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
func convertEcdDesktopPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
