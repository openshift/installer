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

func resourceAlicloudGaBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaBandwidthPackageCreate,
		Read:   resourceAlicloudGaBandwidthPackageRead,
		Update: resourceAlicloudGaBandwidthPackageUpdate,
		Delete: resourceAlicloudGaBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_use_coupon": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bandwidth_package_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Advanced", "Basic", "Enhanced"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "Basic"
				},
			},
			"billing_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayBy95", "PayByTraffic"}, false),
			},
			"cbn_geographic_region_ida": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "China-mainland",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "CrossDomain"
				},
			},
			"cbn_geographic_region_idb": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Global",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "CrossDomain"
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
				Default:      "Subscription",
			},
			"ratio": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Basic", "CrossDomain"}, false),
			},
		},
	}
}

func resourceAlicloudGaBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	action := "CreateBandwidthPackage"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}

	if v, ok := d.GetOkExists("auto_use_coupon"); ok {
		request["AutoUseCoupon"] = v
	}

	request["Bandwidth"] = d.Get("bandwidth")
	if v, ok := d.GetOk("bandwidth_type"); ok {
		request["BandwidthType"] = v
	}

	if v, ok := d.GetOk("billing_type"); ok {
		request["BillingType"] = v
	}

	if v, ok := d.GetOk("cbn_geographic_region_ida"); ok {
		request["CbnGeographicRegionIdA"] = v
	}

	if v, ok := d.GetOk("cbn_geographic_region_idb"); ok {
		request["CbnGeographicRegionIdB"] = v
	}

	if v, ok := d.GetOk("duration"); ok {
		request["Duration"] = v
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertGaBandwidthPackagePaymentTypeRequest(v.(string))
		if request["ChargeType"].(string) == "PREPAY" {
			request["PricingCycle"] = "Month"
		}
	}

	if v, ok := d.GetOk("ratio"); ok {
		request["Ratio"] = v
	}

	request["RegionId"] = client.RegionId
	request["Type"] = d.Get("type")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("CreateBandwidthPackage")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_bandwidth_package", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	d.SetId(fmt.Sprint(response["BandwidthPackageId"]))
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, gaService.GaBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaBandwidthPackageUpdate(d, meta)
}
func resourceAlicloudGaBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_bandwidth_package gaService.DescribeGaBandwidthPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("bandwidth", formatInt(object["Bandwidth"]))
	d.Set("bandwidth_package_name", object["Name"])
	d.Set("bandwidth_type", object["BandwidthType"])
	d.Set("cbn_geographic_region_ida", object["CbnGeographicRegionIdA"])
	d.Set("cbn_geographic_region_idb", object["CbnGeographicRegionIdB"])
	d.Set("description", object["Description"])
	d.Set("payment_type", convertGaBandwidthPackagePaymentTypeResponse(object["ChargeType"].(string)))
	d.Set("status", object["State"])
	d.Set("type", object["Type"])
	if val, ok := d.GetOk("auto_use_coupon"); ok {
		d.Set("auto_use_coupon", val)
	}
	return nil
}
func resourceAlicloudGaBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"BandwidthPackageId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("bandwidth") {
		update = true
		request["Bandwidth"] = d.Get("bandwidth")
	}
	if d.HasChange("bandwidth_package_name") {
		update = true
		request["Name"] = d.Get("bandwidth_package_name")
	}
	if !d.IsNewResource() && d.HasChange("bandwidth_type") {
		update = true
		request["BandwidthType"] = d.Get("bandwidth_type")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if update {
		if _, ok := d.GetOkExists("auto_pay"); ok {
			request["AutoPay"] = d.Get("auto_pay")
		}
		if _, ok := d.GetOkExists("auto_use_coupon"); ok {
			request["AutoUseCoupon"] = d.Get("auto_use_coupon")
		}
		action := "UpdateBandwidthPackage"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"NotExist.BandwidthPackage", "StateError.BandwidthPackage", "UpgradeError.BandwidthPackage"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"active", "binded"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, gaService.GaBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudGaBandwidthPackageRead(d, meta)
}
func resourceAlicloudGaBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type") == "Subscription" {
		log.Printf("[WARN] Cannot destroy resourceAlicloudGaBandwidthPackage prepay type. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteBandwidthPackage"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BandwidthPackageId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["ClientToken"] = buildClientToken("DeleteBandwidthPackage")
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertGaBandwidthPackagePaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "POSTPAY"
	case "Subscription":
		return "PREPAY"
	}
	return source
}
func convertGaBandwidthPackagePaymentTypeResponse(source string) string {
	switch source {
	case "POSTPAY":
		return "PayAsYouGo"
	case "PREPAY":
		return "Subscription"
	}
	return source
}
