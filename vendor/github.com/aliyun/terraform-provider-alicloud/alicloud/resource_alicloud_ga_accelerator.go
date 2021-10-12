package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
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
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duration": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 9),
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
	request["AutoPay"] = true
	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}

	request["Duration"] = d.Get("duration")
	request["PricingCycle"] = "Month"
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
	if val, ok := d.GetOk("auto_use_coupon"); ok {
		d.Set("auto_use_coupon", val)
	}
	return nil
}
func resourceAlicloudGaAcceleratorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"AcceleratorId": d.Id(),
	}
	if d.HasChange("accelerator_name") {
		update = true
		request["Name"] = d.Get("accelerator_name")
	}
	request["AutoPay"] = true
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("spec") {
		update = true
		request["Spec"] = d.Get("spec")
	}
	if update {
		if _, ok := d.GetOkExists("auto_use_coupon"); ok {
			request["AutoUseCoupon"] = d.Get("auto_use_coupon")
		}
		action := "UpdateAccelerator"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		request["ClientToken"] = buildClientToken("UpdateAccelerator")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaAcceleratorStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudGaAcceleratorRead(d, meta)
}
func resourceAlicloudGaAcceleratorDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudGaAccelerator. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
