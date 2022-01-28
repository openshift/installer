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

func resourceAlicloudGaAccelerator() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaAcceleratorCreate,
		Read:   resourceAlicloudGaAcceleratorRead,
		Update: resourceAlicloudGaAcceleratorUpdate,
		Delete: resourceAlicloudGaAcceleratorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renew_duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duration": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"pricing_cycle": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
			},
			"spec": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"1", "10", "2", "3", "5", "8"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaAcceleratorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateAccelerator"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	// there is an api bug that the name can not effect
	//if v, ok := d.GetOk("accelerator_name"); ok {
	//	request["Name"] = v
	//}
	request["AutoPay"] = true
	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}
	request["Duration"] = d.Get("duration")
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	} else {
		request["PricingCycle"] = "Month"
	}
	request["RegionId"] = client.RegionId
	request["Spec"] = d.Get("spec")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateAccelerator")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_accelerator", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["AcceleratorId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaAcceleratorUpdate(d, meta)
}
func resourceAlicloudGaAcceleratorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaAccelerator(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_accelerator gaService.DescribeGaAccelerator Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("accelerator_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("spec", object["Spec"])
	d.Set("status", object["State"])
	describeAcceleratorAutoRenewAttributeObject, err := gaService.DescribeAcceleratorAutoRenewAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if v, ok := describeAcceleratorAutoRenewAttributeObject["AutoRenewDuration"]; ok && fmt.Sprint(v) != "0" {
		d.Set("auto_renew_duration", formatInt(v))
	}
	d.Set("renewal_status", describeAcceleratorAutoRenewAttributeObject["RenewalStatus"])
	return nil
}
func resourceAlicloudGaAcceleratorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	d.Partial(true)
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"AcceleratorId": d.Id(),
	}
	if d.HasChange("auto_renew_duration") {
		update = true
	}
	if v, ok := d.GetOk("auto_renew_duration"); ok {
		request["AutoRenewDuration"] = v
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("renewal_status") {
		update = true
	}
	if v, ok := d.GetOk("renewal_status"); ok {
		request["RenewalStatus"] = v
	}
	if update {
		action := "UpdateAcceleratorAutoRenewAttribute"
		request["ClientToken"] = buildClientToken("UpdateAcceleratorAutoRenewAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
		d.SetPartial("auto_renew_duration")
		d.SetPartial("renewal_status")
	}
	update = false
	updateAcceleratorReq := map[string]interface{}{
		"AcceleratorId": d.Id(),
	}
	if d.HasChange("accelerator_name") {
		update = true
		if v, ok := d.GetOk("accelerator_name"); ok {
			updateAcceleratorReq["Name"] = v
		}
	}
	updateAcceleratorReq["AutoPay"] = true
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			updateAcceleratorReq["Description"] = v
		}
	}
	updateAcceleratorReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("spec") {
		update = true
		updateAcceleratorReq["Spec"] = d.Get("spec")
	}
	if update {
		if v, ok := d.GetOkExists("auto_use_coupon"); ok {
			updateAcceleratorReq["AutoUseCoupon"] = v
		}
		action := "UpdateAccelerator"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		request["ClientToken"] = buildClientToken("UpdateAccelerator")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, updateAcceleratorReq, &runtime)
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("accelerator_name")
		d.SetPartial("description")
		d.SetPartial("spec")
	}
	d.Partial(false)
	return resourceAlicloudGaAcceleratorRead(d, meta)
}
func resourceAlicloudGaAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudGaAccelerator. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
